# Direito Lux - Infrastructure Configuration
# Main infrastructure deployment using modules

# Enable required Google Cloud APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "compute.googleapis.com",
    "container.googleapis.com",
    "servicenetworking.googleapis.com",
    "sqladmin.googleapis.com",
    "redis.googleapis.com",
    "secretmanager.googleapis.com",
    "monitoring.googleapis.com",
    "logging.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "artifactregistry.googleapis.com",
    "dns.googleapis.com",
    "cloudkms.googleapis.com",
    "bigquery.googleapis.com",
    "storage.googleapis.com"
  ])

  project = var.project_id
  service = each.value

  disable_dependent_services = false
  disable_on_destroy        = false
}

# KMS key for encryption
resource "google_kms_key_ring" "direito_lux" {
  name     = "direito-lux-${var.environment}"
  location = var.region

  depends_on = [google_project_service.required_apis]
}

resource "google_kms_crypto_key" "database" {
  name     = "database-encryption"
  key_ring = google_kms_key_ring.direito_lux.id

  rotation_period = "7776000s" # 90 days

  lifecycle {
    prevent_destroy = true
  }
}

# Networking module
module "networking" {
  source = "./modules/networking"

  vpc_name           = local.vpc_name
  region             = var.region
  environment        = var.environment
  subnet_cidrs       = var.subnet_cidrs
  authorized_networks = [for network in var.authorized_networks : network.cidr_block]
  labels             = local.common_labels

  depends_on = [google_project_service.required_apis]
}

# Environment-specific configuration
locals {
  # Select configuration based on environment
  env_config = var.environment == "production" ? var.production_config : var.staging_config
  
  # Database configuration
  database_config = {
    tier                = var.environment == "production" ? var.production_config.db_tier : var.staging_config.db_tier
    disk_size           = var.environment == "production" ? 200 : 100
    disk_type           = "PD_SSD"
    availability_type   = var.environment == "production" ? "REGIONAL" : "ZONAL"
    backup_enabled      = true
    backup_start_time   = "03:00"
    maintenance_window_day  = 7
    maintenance_window_hour = 4
    deletion_protection = var.environment == "production"
  }

  # Redis configuration
  redis_config = {
    memory_size_gb     = var.environment == "production" ? var.production_config.redis_memory : var.staging_config.redis_memory
    tier               = var.environment == "production" ? "STANDARD_HA" : "BASIC"
    redis_version      = "REDIS_7_0"
    auth_enabled       = true
    transit_encryption = var.environment == "production"
    replica_count      = var.environment == "production" ? 1 : 0
  }

  # GKE node pools configuration
  gke_node_pools = {
    primary = {
      machine_type   = local.env_config.machine_type
      min_count      = 1
      max_count      = var.environment == "production" ? 20 : 10
      initial_count  = local.env_config.node_count
      disk_size_gb   = local.env_config.disk_size
      disk_type      = "pd-ssd"
      preemptible    = var.environment != "production"
      auto_upgrade   = true
      auto_repair    = true
      taints         = []
    }
    ai_workload = {
      machine_type   = var.environment == "production" ? "n1-highmem-8" : "n1-highmem-4"
      min_count      = 0
      max_count      = var.environment == "production" ? 10 : 5
      initial_count  = var.environment == "production" ? 2 : 1
      disk_size_gb   = 200
      disk_type      = "pd-ssd"
      preemptible    = var.spot_instances
      auto_upgrade   = true
      auto_repair    = true
      taints         = []
    }
  }
}

# Database module
module "database" {
  source = "./modules/database"

  project_id            = var.project_id
  environment           = var.environment
  region                = var.region
  zone                  = var.zone
  secondary_region      = var.secondary_region
  db_instance_name      = local.db_instance_name
  redis_instance_name   = local.redis_instance_name
  database_version      = var.database_version
  vpc_id                = module.networking.vpc_id
  private_vpc_connection = module.networking.private_vpc_connection_id
  database_config       = local.database_config
  redis_config          = local.redis_config
  backup_config         = var.backup_config
  authorized_networks   = var.authorized_networks
  notification_channels = []  # Will be populated later with monitoring setup
  labels                = local.common_labels

  depends_on = [
    module.networking,
    google_project_service.required_apis
  ]
}

# GKE module
module "gke" {
  source = "./modules/gke"

  project_id               = var.project_id
  cluster_name             = local.cluster_name
  region                   = var.region
  zone                     = var.zone
  regional_cluster         = true
  vpc_name                 = module.networking.vpc_name
  subnet_name              = module.networking.gke_subnet_name
  pods_range_name          = module.networking.gke_pods_range_name
  services_range_name      = module.networking.gke_services_range_name
  enable_private_nodes     = var.enable_private_nodes
  master_ipv4_cidr_block   = var.master_ipv4_cidr_block
  authorized_networks      = var.authorized_networks
  enable_workload_identity = var.enable_workload_identity
  enable_network_policy    = var.enable_network_policy
  enable_pod_security_policy = var.enable_pod_security_policy
  database_encryption_key  = google_kms_crypto_key.database.id
  release_channel          = "REGULAR"
  node_pools              = local.gke_node_pools
  enable_ai_node_pool     = var.app_config.ai_service_enabled
  labels                  = local.common_labels

  depends_on = [
    module.networking,
    google_project_service.required_apis,
    google_kms_crypto_key.database
  ]
}

# Artifact Registry for container images
resource "google_artifact_registry_repository" "direito_lux" {
  location      = var.region
  repository_id = "direito-lux-${var.environment}"
  description   = "Container registry for Direito Lux ${var.environment} environment"
  format        = "DOCKER"

  labels = local.common_labels

  depends_on = [google_project_service.required_apis]
}

# Load Balancer and SSL certificates
resource "google_compute_global_address" "default" {
  name = "${var.environment}-lb-ip"

  depends_on = [google_project_service.required_apis]
}

# Managed SSL certificate
resource "google_compute_managed_ssl_certificate" "default" {
  count = var.enable_ssl ? 1 : 0

  name = "${var.environment}-ssl-cert"

  managed {
    domains = var.environment == "production" ? [
      var.domain_name,
      "app.${var.domain_name}",
      "api.${var.domain_name}"
    ] : [
      "${var.subdomain_staging}.${var.domain_name}",
      "${var.subdomain_api}-${var.subdomain_staging}.${var.domain_name}"
    ]
  }

  depends_on = [google_project_service.required_apis]
}

# Cloud DNS zone
resource "google_dns_managed_zone" "default" {
  name        = "${var.environment}-dns-zone"
  dns_name    = var.environment == "production" ? "${var.domain_name}." : "${var.subdomain_staging}.${var.domain_name}."
  description = "DNS zone for Direito Lux ${var.environment} environment"

  labels = local.common_labels

  depends_on = [google_project_service.required_apis]
}

# DNS A record for load balancer
resource "google_dns_record_set" "default" {
  name = google_dns_managed_zone.default.dns_name
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.default.name

  rrdatas = [google_compute_global_address.default.address]
}

# DNS A record for API subdomain
resource "google_dns_record_set" "api" {
  name = var.environment == "production" ? "api.${var.domain_name}." : "${var.subdomain_api}-${var.subdomain_staging}.${var.domain_name}."
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.default.name

  rrdatas = [google_compute_global_address.default.address]
}

# IAM service accounts for applications
resource "google_service_account" "app_services" {
  for_each = toset([
    "auth-service",
    "tenant-service", 
    "process-service",
    "ai-service",
    "notification-service",
    "search-service",
    "report-service",
    "mcp-service",
    "datajud-service"
  ])

  account_id   = "${var.environment}-${each.value}"
  display_name = "Service Account for ${each.value} in ${var.environment}"
  description  = "Service account for ${each.value} microservice"
}

# IAM bindings for application services
resource "google_project_iam_member" "app_services_roles" {
  for_each = toset([
    "auth-service",
    "tenant-service",
    "process-service", 
    "ai-service",
    "notification-service",
    "search-service",
    "report-service",
    "mcp-service",
    "datajud-service"
  ])

  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.app_services[each.value].email}"
}

# Workload Identity bindings
resource "google_service_account_iam_binding" "workload_identity" {
  for_each = var.enable_workload_identity ? toset([
    "auth-service",
    "tenant-service",
    "process-service",
    "ai-service", 
    "notification-service",
    "search-service",
    "report-service",
    "mcp-service",
    "datajud-service"
  ]) : toset([])

  service_account_id = google_service_account.app_services[each.value].name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[direito-lux-${var.environment}/${each.value}]"
  ]
}

# Monitoring notification channels
resource "google_monitoring_notification_channel" "email" {
  display_name = "${var.environment} Email Alerts"
  type         = "email"

  labels = {
    email_address = "devops@direitolux.com"
  }

  depends_on = [google_project_service.required_apis]
}

resource "google_monitoring_notification_channel" "slack" {
  count = var.environment == "production" ? 1 : 0

  display_name = "${var.environment} Slack Alerts"
  type         = "slack"

  labels = {
    channel_name = "#alerts"
    url          = "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  }

  depends_on = [google_project_service.required_apis]
}

# Log sinks for centralized logging
resource "google_logging_project_sink" "gke_logs" {
  name = "${var.environment}-gke-logs"

  destination = "bigquery.googleapis.com/projects/${var.project_id}/datasets/${google_bigquery_dataset.logs.dataset_id}"

  filter = "resource.type=\"k8s_container\" OR resource.type=\"gke_cluster\""

  unique_writer_identity = true

  bigquery_options {
    use_partitioned_tables = true
  }

  depends_on = [google_project_service.required_apis]
}

# BigQuery dataset for logs
resource "google_bigquery_dataset" "logs" {
  dataset_id  = "direito_lux_${var.environment}_logs"
  description = "Centralized logs for Direito Lux ${var.environment} environment"
  location    = var.region

  default_table_expiration_ms = var.log_retention_days * 24 * 60 * 60 * 1000

  labels = local.common_labels

  depends_on = [google_project_service.required_apis]
}

# BigQuery dataset access for log sink
resource "google_bigquery_dataset_iam_member" "log_sink_writer" {
  dataset_id = google_bigquery_dataset.logs.dataset_id
  role       = "roles/bigquery.dataEditor"
  member     = google_logging_project_sink.gke_logs.writer_identity
}

# Storage bucket for Terraform state (created separately)
data "google_storage_bucket" "terraform_state" {
  name = "direito-lux-terraform-state"
}

# Output important information
output "infrastructure_summary" {
  description = "Summary of deployed infrastructure"
  value = {
    environment         = var.environment
    region             = var.region
    gke_cluster_name   = module.gke.cluster_name
    gke_cluster_endpoint = module.gke.cluster_endpoint
    database_instance  = module.database.postgres_instance_name
    redis_instance     = module.database.redis_instance_name
    load_balancer_ip   = google_compute_global_address.default.address
    domain_name        = var.environment == "production" ? var.domain_name : "${var.subdomain_staging}.${var.domain_name}"
    artifact_registry  = google_artifact_registry_repository.direito_lux.name
  }
  sensitive = true
}