# Infraestrutura como Código (IaC) - GCP

## Stack IaC Escolhida
- **Terraform** para provisionamento de recursos GCP
- **Pulumi** como alternativa (TypeScript/Go)
- **Helm Charts** para deployments Kubernetes
- **ArgoCD** para GitOps

## Estrutura de Diretórios
```
infrastructure/
├── terraform/
│   ├── environments/
│   │   ├── dev/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   └── terraform.tfvars
│   │   ├── staging/
│   │   └── production/
│   ├── modules/
│   │   ├── gke/
│   │   ├── cloudsql/
│   │   ├── memorystore/
│   │   ├── pubsub/
│   │   ├── gcs/
│   │   ├── vpc/
│   │   └── iam/
│   └── global/
│       ├── dns/
│       └── cdn/
├── kubernetes/
│   ├── base/
│   ├── overlays/
│   └── helm/
└── scripts/
    ├── setup.sh
    └── destroy.sh
```

## Terraform - Recursos GCP Core

### 1. VPC e Network
```hcl
# modules/vpc/main.tf
resource "google_compute_network" "main" {
  name                    = "${var.project}-vpc"
  auto_create_subnetworks = false
  project                 = var.project_id
}

resource "google_compute_subnetwork" "gke" {
  name          = "${var.project}-gke-subnet"
  ip_cidr_range = "10.0.0.0/20"
  region        = var.region
  network       = google_compute_network.main.id
  
  secondary_ip_range {
    range_name    = "gke-pods"
    ip_cidr_range = "10.4.0.0/14"
  }
  
  secondary_ip_range {
    range_name    = "gke-services"
    ip_cidr_range = "10.8.0.0/20"
  }
  
  private_ip_google_access = true
}

# Cloud NAT para saída de internet
resource "google_compute_router" "main" {
  name    = "${var.project}-router"
  region  = var.region
  network = google_compute_network.main.id
}

resource "google_compute_router_nat" "main" {
  name                               = "${var.project}-nat"
  router                             = google_compute_router.main.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}
```

### 2. GKE Cluster
```hcl
# modules/gke/main.tf
resource "google_container_cluster" "primary" {
  name     = "${var.project}-gke"
  location = var.region
  
  # Autopilot para gestão simplificada
  enable_autopilot = true
  
  network    = var.network_id
  subnetwork = var.subnetwork_id
  
  ip_allocation_policy {
    cluster_secondary_range_name  = "gke-pods"
    services_secondary_range_name = "gke-services"
  }
  
  # Workload Identity
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }
  
  # Binary Authorization
  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }
  
  # Network Policy
  network_policy {
    enabled = true
  }
  
  # Monitoring
  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    managed_prometheus {
      enabled = true
    }
  }
}
```

### 3. Cloud SQL (PostgreSQL)
```hcl
# modules/cloudsql/main.tf
resource "google_sql_database_instance" "postgres" {
  name             = "${var.project}-postgres"
  database_version = "POSTGRES_15"
  region           = var.region
  
  settings {
    tier              = var.tier # db-custom-2-8192 para produção
    availability_type = "REGIONAL" # HA
    
    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
    }
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.network_id
      require_ssl     = true
    }
    
    database_flags {
      name  = "max_connections"
      value = "200"
    }
    
    insights_config {
      query_insights_enabled  = true
      query_string_length     = 1024
      record_application_tags = true
    }
  }
  
  deletion_protection = true
}

# Database por tenant (multi-tenant com isolamento)
resource "google_sql_database" "tenant_template" {
  name     = "tenant_template"
  instance = google_sql_database_instance.postgres.name
}
```

### 4. Memorystore (Redis)
```hcl
# modules/memorystore/main.tf
resource "google_redis_instance" "cache" {
  name               = "${var.project}-redis"
  memory_size_gb     = var.memory_size
  tier               = "STANDARD_HA"
  region             = var.region
  
  redis_version      = "REDIS_7_0"
  display_name       = "Direito Lux Cache"
  
  auth_enabled       = true
  transit_encryption_mode = "SERVER_AUTHENTICATION"
  
  persistence_config {
    persistence_mode = "RDB"
    rdb_snapshot_period = "ONE_HOUR"
  }
  
  maintenance_policy {
    weekly_maintenance_window {
      day = "SUNDAY"
      start_time {
        hours   = 3
        minutes = 0
      }
    }
  }
}
```

### 5. Pub/Sub (Message Broker)
```hcl
# modules/pubsub/main.tf
# Topics principais
locals {
  topics = [
    "process-events",
    "notification-events",
    "document-events",
    "tenant-events",
    "billing-events"
  ]
}

resource "google_pubsub_topic" "events" {
  for_each = toset(local.topics)
  
  name = each.value
  
  message_retention_duration = "604800s" # 7 dias
  
  schema_settings {
    schema = google_pubsub_schema.event_schema[each.key].id
    encoding = "JSON"
  }
}

# Dead Letter Queue
resource "google_pubsub_topic" "dlq" {
  for_each = toset(local.topics)
  name     = "${each.value}-dlq"
}

# Subscriptions com retry policy
resource "google_pubsub_subscription" "services" {
  for_each = toset(local.topics)
  
  name  = "${each.value}-subscription"
  topic = google_pubsub_topic.events[each.key].name
  
  ack_deadline_seconds = 60
  
  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }
  
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dlq[each.key].id
    max_delivery_attempts = 5
  }
  
  enable_exactly_once_delivery = true
}
```

### 6. Cloud Storage
```hcl
# modules/gcs/main.tf
# Bucket para documentos
resource "google_storage_bucket" "documents" {
  name          = "${var.project}-documents"
  location      = var.region
  storage_class = "STANDARD"
  
  uniform_bucket_level_access = true
  
  encryption {
    default_kms_key_name = var.kms_key
  }
  
  lifecycle_rule {
    condition {
      age = 90
      matches_storage_class = ["STANDARD"]
    }
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }
  
  versioning {
    enabled = true
  }
}

# Bucket para backups
resource "google_storage_bucket" "backups" {
  name          = "${var.project}-backups"
  location      = var.region
  storage_class = "NEARLINE"
  
  lifecycle_rule {
    condition {
      age = 365
    }
    action {
      type          = "SetStorageClass"
      storage_class = "ARCHIVE"
    }
  }
}
```

### 7. IAM e Service Accounts
```hcl
# modules/iam/main.tf
# Service Account para cada microserviço
locals {
  services = [
    "auth-service",
    "process-service",
    "datajud-service",
    "notification-service",
    "ai-service",
    "tenant-service"
  ]
}

resource "google_service_account" "services" {
  for_each = toset(local.services)
  
  account_id   = each.value
  display_name = "${each.value} Service Account"
}

# Workload Identity Binding
resource "google_service_account_iam_member" "workload_identity" {
  for_each = toset(local.services)
  
  service_account_id = google_service_account.services[each.key].name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[direito-lux/${each.key}]"
}

# Roles específicas por serviço
resource "google_project_iam_member" "cloudsql_client" {
  for_each = toset(["auth-service", "process-service", "tenant-service"])
  
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.services[each.key].email}"
}
```

### 8. Load Balancer e CDN
```hcl
# modules/cdn/main.tf
# Global Load Balancer
resource "google_compute_global_address" "default" {
  name = "${var.project}-global-ip"
}

# Cloud Armor Security Policy
resource "google_compute_security_policy" "default" {
  name = "${var.project}-security-policy"
  
  rule {
    action   = "rate_based_ban"
    priority = "1000"
    
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    
    rate_limit_options {
      rate_limit_threshold {
        count        = 100
        interval_sec = 60
      }
      ban_duration_sec = 600
    }
  }
  
  # OWASP Top 10 rules
  rule {
    action   = "deny(403)"
    priority = "2000"
    
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('sqli-stable')"
      }
    }
  }
}

# Cloud CDN
resource "google_compute_backend_bucket" "static" {
  name        = "${var.project}-static-backend"
  bucket_name = google_storage_bucket.static.name
  enable_cdn  = true
  
  cdn_policy {
    cache_mode = "CACHE_ALL_STATIC"
    default_ttl = 3600
    max_ttl     = 86400
  }
}
```

### 9. Monitoring e Logging
```hcl
# modules/monitoring/main.tf
# Log Sink para BigQuery
resource "google_logging_project_sink" "bigquery" {
  name        = "${var.project}-bigquery-sink"
  destination = "bigquery.googleapis.com/${google_bigquery_dataset.logs.id}"
  
  filter = <<EOF
    resource.type="k8s_container"
    severity >= "WARNING"
  EOF
  
  unique_writer_identity = true
}

# Uptime Checks
resource "google_monitoring_uptime_check_config" "api" {
  display_name = "API Health Check"
  timeout      = "10s"
  period       = "60s"
  
  http_check {
    path         = "/health"
    port         = "443"
    use_ssl      = true
    validate_ssl = true
  }
  
  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = var.project_id
      host       = "api.direitolux.com.br"
    }
  }
}

# Alert Policies
resource "google_monitoring_alert_policy" "high_latency" {
  display_name = "High API Latency"
  combiner     = "OR"
  
  conditions {
    display_name = "API Latency > 500ms"
    
    condition_threshold {
      filter          = "metric.type=\"serviceruntime.googleapis.com/api/request_latencies\""
      duration        = "300s"
      comparison      = "COMPARISON_GT"
      threshold_value = 500
      
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_PERCENTILE_99"
      }
    }
  }
  
  notification_channels = [google_monitoring_notification_channel.email.id]
}
```

## Helm Charts para Aplicações

### Auth Service Helm Chart
```yaml
# kubernetes/helm/auth-service/values.yaml
replicaCount: 3

image:
  repository: gcr.io/direito-lux/auth-service
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080
  grpc: 9090

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: api.direitolux.com.br
      paths:
        - path: /auth
          pathType: Prefix

resources:
  requests:
    cpu: 250m
    memory: 512Mi
  limits:
    cpu: 500m
    memory: 1Gi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

env:
  - name: KEYCLOAK_URL
    value: "https://auth.direitolux.com.br"
  - name: DATABASE_NAME
    value: "auth_db"

secrets:
  - name: database-credentials
    keys:
      - DB_USER
      - DB_PASSWORD
```

## GitOps com ArgoCD

### Application Manifest
```yaml
# kubernetes/argocd/applications/direito-lux.yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: direito-lux-services
  namespace: argocd
spec:
  generators:
  - list:
      elements:
      - service: auth-service
        namespace: direito-lux
      - service: process-service
        namespace: direito-lux
      - service: datajud-service
        namespace: direito-lux
  template:
    metadata:
      name: '{{service}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/direitolux/infrastructure
        targetRevision: HEAD
        path: kubernetes/helm/{{service}}
      destination:
        server: https://kubernetes.default.svc
        namespace: '{{namespace}}'
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
        - CreateNamespace=true
```

## Scripts de Automação

### Setup Inicial
```bash
#!/bin/bash
# scripts/setup.sh

# Configurar projeto GCP
export PROJECT_ID="direito-lux-prod"
export REGION="southamerica-east1"

# Criar projeto
gcloud projects create $PROJECT_ID --name="Direito Lux"

# Ativar APIs necessárias
gcloud services enable \
  compute.googleapis.com \
  container.googleapis.com \
  sqladmin.googleapis.com \
  redis.googleapis.com \
  storage.googleapis.com \
  pubsub.googleapis.com \
  cloudkms.googleapis.com \
  monitoring.googleapis.com \
  logging.googleapis.com \
  --project=$PROJECT_ID

# Terraform init e apply
cd infrastructure/terraform/environments/production
terraform init
terraform plan -out=tfplan
terraform apply tfplan

# Instalar ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Configurar Istio
istioctl install --set profile=demo -y
kubectl label namespace direito-lux istio-injection=enabled
```

## Estimativa de Custos GCP

### Recursos Base (Mensal)
- **GKE Autopilot**: ~$74/cluster + ~$0.10/vCPU-hora
- **Cloud SQL (2vCPU, 8GB)**: ~$140
- **Memorystore Redis (1GB HA)**: ~$90
- **Cloud Storage**: ~$20/TB
- **Pub/Sub**: ~$40 (primeiros 10GB free)
- **Load Balancer**: ~$25
- **Cloud CDN**: ~$0.08/GB
- **Monitoring**: ~$0.258/milhão de logs

### Total Estimado
- **Desenvolvimento**: ~$300/mês
- **Produção (inicial)**: ~$600/mês
- **Produção (1000 usuários)**: ~$1,200/mês

## CI/CD Pipeline

### GitHub Actions
```yaml
# .github/workflows/deploy.yml
name: Deploy to GKE

on:
  push:
    branches: [main]

env:
  PROJECT_ID: direito-lux-prod
  GKE_CLUSTER: direito-lux-gke
  GKE_ZONE: southamerica-east1

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - id: 'auth'
      uses: 'google-github-actions/auth@v1'
      with:
        credentials_json: ${{ secrets.GCP_SA_KEY }}
    
    - name: 'Set up Cloud SDK'
      uses: 'google-github-actions/setup-gcloud@v1'
    
    - name: 'Configure Docker'
      run: |
        gcloud auth configure-docker
    
    - name: 'Build and Push Image'
      run: |
        docker build -t gcr.io/$PROJECT_ID/$SERVICE_NAME:$GITHUB_SHA .
        docker push gcr.io/$PROJECT_ID/$SERVICE_NAME:$GITHUB_SHA
    
    - name: 'Update Helm Values'
      run: |
        sed -i "s|tag:.*|tag: $GITHUB_SHA|" kubernetes/helm/$SERVICE_NAME/values.yaml
        git add .
        git commit -m "Update image tag to $GITHUB_SHA"
        git push
```

## Disaster Recovery

### Backup Automático
```hcl
# Backup do Cloud SQL para GCS
resource "google_sql_database_instance" "postgres" {
  # ... configurações anteriores ...
  
  backup_configuration {
    enabled = true
    start_time = "03:00"
    location = "southamerica-east1"
    
    backup_retention_settings {
      retained_backups = 30
      retention_unit   = "COUNT"
    }
  }
}

# Backup cross-region
resource "google_storage_bucket" "backup_replica" {
  name     = "${var.project}-backup-replica"
  location = "US" # Multi-region para DR
}
```

## Segurança Adicional

### Secret Manager
```hcl
resource "google_secret_manager_secret" "api_keys" {
  secret_id = "datajud-api-key"
  
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "api_key_v1" {
  secret = google_secret_manager_secret.api_keys.id
  secret_data = var.datajud_api_key
}
```

### VPC Service Controls
```hcl
resource "google_access_context_manager_service_perimeter" "secure_perimeter" {
  parent = "accessPolicies/${var.access_policy}"
  name   = "accessPolicies/${var.access_policy}/servicePerimeters/secure_data"
  title  = "Secure Data Perimeter"
  
  status {
    restricted_services = [
      "storage.googleapis.com",
      "sqladmin.googleapis.com",
    ]
  }
}
```