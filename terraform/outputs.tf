# Direito Lux - Terraform Outputs
# Main infrastructure outputs

# Project Information
output "project_id" {
  description = "GCP Project ID"
  value       = var.project_id
}

output "environment" {
  description = "Environment name"
  value       = var.environment
}

output "region" {
  description = "GCP Region"
  value       = var.region
}

# Networking Outputs
output "vpc_name" {
  description = "VPC network name"
  value       = module.networking.vpc_name
}

output "vpc_id" {
  description = "VPC network ID"
  value       = module.networking.vpc_id
}

output "public_subnet_name" {
  description = "Public subnet name"
  value       = module.networking.public_subnet_name
}

output "private_subnet_name" {
  description = "Private subnet name"
  value       = module.networking.private_subnet_name
}

output "gke_subnet_name" {
  description = "GKE subnet name"
  value       = module.networking.gke_subnet_name
}

# GKE Cluster Outputs
output "gke_cluster_name" {
  description = "GKE cluster name"
  value       = module.gke.cluster_name
}

output "gke_cluster_endpoint" {
  description = "GKE cluster endpoint"
  value       = module.gke.cluster_endpoint
  sensitive   = true
}

output "gke_cluster_ca_certificate" {
  description = "GKE cluster CA certificate"
  value       = module.gke.cluster_ca_certificate
  sensitive   = true
}

output "gke_cluster_location" {
  description = "GKE cluster location"
  value       = module.gke.cluster_location
}

output "gke_node_service_account" {
  description = "GKE node service account email"
  value       = module.gke.node_service_account
}

# Database Outputs
output "postgres_instance_name" {
  description = "PostgreSQL instance name"
  value       = module.database.postgres_instance_name
}

output "postgres_connection_name" {
  description = "PostgreSQL connection name"
  value       = module.database.postgres_instance_connection_name
}

output "postgres_private_ip" {
  description = "PostgreSQL private IP address"
  value       = module.database.postgres_instance_private_ip_address
  sensitive   = true
}

output "postgres_database_name" {
  description = "PostgreSQL database name"
  value       = module.database.postgres_database_name
}

output "postgres_user_name" {
  description = "PostgreSQL user name"
  value       = module.database.postgres_user_name
}

output "redis_instance_name" {
  description = "Redis instance name"
  value       = module.database.redis_instance_name
}

output "redis_host" {
  description = "Redis host"
  value       = module.database.redis_host
  sensitive   = true
}

output "redis_port" {
  description = "Redis port"
  value       = module.database.redis_port
}

# Load Balancer and DNS
output "load_balancer_ip" {
  description = "Load balancer IP address"
  value       = google_compute_global_address.default.address
}

output "domain_name" {
  description = "Domain name"
  value       = var.environment == "production" ? var.domain_name : "${var.subdomain_staging}.${var.domain_name}"
}

output "api_domain" {
  description = "API domain name"
  value       = var.environment == "production" ? "api.${var.domain_name}" : "${var.subdomain_api}-${var.subdomain_staging}.${var.domain_name}"
}

output "dns_zone_name" {
  description = "DNS managed zone name"
  value       = google_dns_managed_zone.default.name
}

output "dns_zone_dns_name" {
  description = "DNS zone DNS name"
  value       = google_dns_managed_zone.default.dns_name
}

# SSL Certificate
output "ssl_certificate_name" {
  description = "SSL certificate name"
  value       = var.enable_ssl ? google_compute_managed_ssl_certificate.default[0].name : null
}

# Container Registry
output "artifact_registry_repository" {
  description = "Artifact Registry repository name"
  value       = google_artifact_registry_repository.direito_lux.name
}

output "artifact_registry_location" {
  description = "Artifact Registry location"
  value       = google_artifact_registry_repository.direito_lux.location
}

# Service Accounts
output "app_service_accounts" {
  description = "Application service account emails"
  value = {
    for service, account in google_service_account.app_services :
    service => account.email
  }
}

output "sql_proxy_service_account" {
  description = "Cloud SQL Proxy service account email"
  value       = module.database.sql_proxy_service_account_email
}

# Secrets
output "database_password_secret" {
  description = "Database password secret ID"
  value       = module.database.db_password_secret_id
}

output "redis_auth_secret" {
  description = "Redis auth secret ID"
  value       = module.database.redis_auth_secret_id
}

# Connection Information for Applications
output "connection_info" {
  description = "Connection information for applications"
  value = {
    postgres = {
      host     = module.database.postgres_instance_private_ip_address
      port     = 5432
      database = module.database.postgres_database_name
      user     = module.database.postgres_user_name
    }
    redis = {
      host = module.database.redis_host
      port = module.database.redis_port
    }
    gke = {
      cluster_name = module.gke.cluster_name
      endpoint     = module.gke.cluster_endpoint
      location     = module.gke.cluster_location
    }
  }
  sensitive = true
}

# Monitoring
output "monitoring_notification_channels" {
  description = "Monitoring notification channels"
  value = [
    google_monitoring_notification_channel.email.name
  ]
}

# BigQuery Datasets
output "logs_dataset_id" {
  description = "BigQuery logs dataset ID"
  value       = google_bigquery_dataset.logs.dataset_id
}

output "gke_usage_dataset_id" {
  description = "GKE usage BigQuery dataset ID"
  value       = module.gke.bigquery_dataset_id
}

# KMS
output "kms_key_ring_name" {
  description = "KMS key ring name"
  value       = google_kms_key_ring.direito_lux.name
}

output "database_kms_key_name" {
  description = "Database encryption KMS key name"
  value       = google_kms_crypto_key.database.name
}

# Kubernetes Configuration
output "kubernetes_config" {
  description = "Kubernetes provider configuration"
  value       = module.gke.kubernetes_config
  sensitive   = true
}

# URLs for accessing services
output "service_urls" {
  description = "Service URLs"
  value = {
    frontend    = var.environment == "production" ? "https://app.${var.domain_name}" : "https://${var.subdomain_staging}.${var.domain_name}"
    api         = var.environment == "production" ? "https://api.${var.domain_name}" : "https://${var.subdomain_api}-${var.subdomain_staging}.${var.domain_name}"
    grafana     = "${var.environment == "production" ? "https://api.${var.domain_name}" : "https://${var.subdomain_api}-${var.subdomain_staging}.${var.domain_name}"}/grafana"
    prometheus  = "${var.environment == "production" ? "https://api.${var.domain_name}" : "https://${var.subdomain_api}-${var.subdomain_staging}.${var.domain_name}"}/prometheus"
  }
}

# Environment Configuration Summary
output "environment_summary" {
  description = "Environment configuration summary"
  value = {
    environment_type = var.environment
    high_availability = var.environment == "production"
    node_count = var.environment == "production" ? var.production_config.node_count : var.staging_config.node_count
    machine_type = var.environment == "production" ? var.production_config.machine_type : var.staging_config.machine_type
    database_tier = var.environment == "production" ? var.production_config.db_tier : var.staging_config.db_tier
    redis_memory = var.environment == "production" ? var.production_config.redis_memory : var.staging_config.redis_memory
    ssl_enabled = var.enable_ssl
    monitoring_enabled = var.enable_monitoring
    backup_retention_days = var.backup_config.retention_period_days
  }
}