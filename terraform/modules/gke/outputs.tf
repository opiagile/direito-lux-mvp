# GKE Module Outputs

output "cluster_id" {
  description = "ID of the GKE cluster"
  value       = google_container_cluster.primary.id
}

output "cluster_name" {
  description = "Name of the GKE cluster"
  value       = google_container_cluster.primary.name
}

output "cluster_location" {
  description = "Location of the GKE cluster"
  value       = google_container_cluster.primary.location
}

output "cluster_endpoint" {
  description = "Endpoint of the GKE cluster"
  value       = google_container_cluster.primary.endpoint
  sensitive   = true
}

output "cluster_ca_certificate" {
  description = "CA certificate of the GKE cluster"
  value       = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  sensitive   = true
}

output "cluster_master_version" {
  description = "Master version of the GKE cluster"
  value       = google_container_cluster.primary.master_version
}

output "cluster_node_version" {
  description = "Node version of the GKE cluster"
  value       = google_container_cluster.primary.node_version
}

output "node_service_account" {
  description = "Service account used by GKE nodes"
  value       = google_service_account.gke_nodes.email
}

output "node_service_account_name" {
  description = "Name of the service account used by GKE nodes"
  value       = google_service_account.gke_nodes.name
}

output "primary_node_pool_name" {
  description = "Name of the primary node pool"
  value       = google_container_node_pool.primary.name
}

output "ai_node_pool_name" {
  description = "Name of the AI workload node pool"
  value       = var.enable_ai_node_pool ? google_container_node_pool.ai_workload[0].name : null
}

output "cluster_resource_labels" {
  description = "The combination of labels configured directly on the resource and default labels"
  value       = google_container_cluster.primary.resource_labels
}

output "bigquery_dataset_id" {
  description = "BigQuery dataset ID for cluster usage export"
  value       = google_bigquery_dataset.gke_usage.dataset_id
}

# Kubernetes provider configuration outputs
output "kubernetes_config" {
  description = "Kubernetes provider configuration"
  value = {
    host                   = "https://${google_container_cluster.primary.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  }
  sensitive = true
}

# Data source for client config
data "google_client_config" "default" {}