# Direito Lux - GKE Module
# Google Kubernetes Engine cluster and node pool configuration

# Service Account for GKE nodes
resource "google_service_account" "gke_nodes" {
  account_id   = "${var.cluster_name}-nodes"
  display_name = "GKE Nodes Service Account for ${var.cluster_name}"
  description  = "Service account for GKE nodes in ${var.cluster_name} cluster"
}

# IAM bindings for GKE node service account
resource "google_project_iam_member" "gke_nodes_roles" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
    "roles/stackdriver.resourceMetadata.writer",
    "roles/storage.objectViewer",
    "roles/artifactregistry.reader"
  ])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}

# Workload Identity binding
resource "google_service_account_iam_member" "workload_identity" {
  count = var.enable_workload_identity ? 1 : 0

  service_account_id = google_service_account.gke_nodes.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[default/default]"
}

# GKE Cluster
resource "google_container_cluster" "primary" {
  name     = var.cluster_name
  location = var.regional_cluster ? var.region : var.zone

  # Remove default node pool and create custom ones
  remove_default_node_pool = true
  initial_node_count       = 1

  # Network configuration
  network    = var.vpc_name
  subnetwork = var.subnet_name

  # Private cluster configuration
  dynamic "private_cluster_config" {
    for_each = var.enable_private_nodes ? [1] : []
    content {
      enable_private_nodes    = true
      enable_private_endpoint = false
      master_ipv4_cidr_block  = var.master_ipv4_cidr_block

      master_global_access_config {
        enabled = true
      }
    }
  }

  # IP allocation policy for pods and services
  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_range_name
    services_secondary_range_name = var.services_range_name
  }

  # Master authorized networks
  dynamic "master_authorized_networks_config" {
    for_each = length(var.authorized_networks) > 0 ? [1] : []
    content {
      dynamic "cidr_blocks" {
        for_each = var.authorized_networks
        content {
          cidr_block   = cidr_blocks.value.cidr_block
          display_name = cidr_blocks.value.display_name
        }
      }
    }
  }

  # Workload Identity
  dynamic "workload_identity_config" {
    for_each = var.enable_workload_identity ? [1] : []
    content {
      workload_pool = "${var.project_id}.svc.id.goog"
    }
  }

  # Network policy
  dynamic "network_policy" {
    for_each = var.enable_network_policy ? [1] : []
    content {
      enabled  = true
      provider = "CALICO"
    }
  }

  # Add-on configuration
  addons_config {
    http_load_balancing {
      disabled = false
    }

    horizontal_pod_autoscaling {
      disabled = false
    }

    network_policy_config {
      disabled = !var.enable_network_policy
    }

    dns_cache_config {
      enabled = true
    }

    gcp_filestore_csi_driver_config {
      enabled = true
    }

    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
  }

  # Binary authorization
  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }

  # Resource usage export
  resource_usage_export_config {
    enable_network_egress_metering       = true
    enable_resource_consumption_metering = true

    bigquery_destination {
      dataset_id = google_bigquery_dataset.gke_usage.dataset_id
    }
  }

  # Telemetry is enabled by default in modern GKE versions

  # Logging and monitoring
  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"

  logging_config {
    enable_components = [
      "SYSTEM_COMPONENTS",
      "WORKLOADS",
      "APISERVER"
    ]
  }

  monitoring_config {
    enable_components = [
      "SYSTEM_COMPONENTS",
      "WORKLOADS"
    ]

    managed_prometheus {
      enabled = true
    }
  }

  # Pod Security Policy is deprecated, use Pod Security Standards instead

  # Database encryption
  database_encryption {
    state    = "ENCRYPTED"
    key_name = var.database_encryption_key
  }

  # Maintenance policy
  maintenance_policy {
    recurring_window {
      start_time = "2023-01-01T04:00:00Z"
      end_time   = "2023-01-01T08:00:00Z"
      recurrence = "FREQ=WEEKLY;BYDAY=SU"
    }
  }

  # Release channel
  release_channel {
    channel = var.release_channel
  }

  # Lifecycle rules
  lifecycle {
    ignore_changes = [
      node_pool,
      initial_node_count,
    ]
  }

  depends_on = [
    google_service_account.gke_nodes,
    google_project_iam_member.gke_nodes_roles
  ]
}

# BigQuery dataset for resource usage export
resource "google_bigquery_dataset" "gke_usage" {
  dataset_id  = "${replace(var.cluster_name, "-", "_")}_usage"
  description = "GKE resource usage data for ${var.cluster_name}"
  location    = var.region

  delete_contents_on_destroy = true

  labels = var.labels
}

# Primary node pool
resource "google_container_node_pool" "primary" {
  name       = "primary-pool"
  cluster    = google_container_cluster.primary.id
  location   = google_container_cluster.primary.location

  # Node pool size
  initial_node_count = var.node_pools["primary"].initial_count

  autoscaling {
    min_node_count = var.node_pools["primary"].min_count
    max_node_count = var.node_pools["primary"].max_count
  }

  # Node configuration
  node_config {
    machine_type = var.node_pools["primary"].machine_type
    
    # Service account
    service_account = google_service_account.gke_nodes.email
    
    # OAuth scopes
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    # Disk configuration
    disk_size_gb = var.node_pools["primary"].disk_size_gb
    disk_type    = var.node_pools["primary"].disk_type

    # Image type
    image_type = "COS_CONTAINERD"

    # Preemptible nodes
    preemptible = var.node_pools["primary"].preemptible

    # Labels
    labels = merge(var.labels, {
      node-pool = "primary"
      workload  = "general"
    })

    # Taints for dedicated workloads
    dynamic "taint" {
      for_each = var.node_pools["primary"].taints
      content {
        key    = taint.value.key
        value  = taint.value.value
        effect = taint.value.effect
      }
    }

    # Metadata
    metadata = {
      disable-legacy-endpoints = "true"
    }

    # Shielded instance config
    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    # Workload metadata config
    workload_metadata_config {
      mode = var.enable_workload_identity ? "GKE_METADATA" : "GCE_METADATA"
    }
  }

  # Management configuration
  management {
    auto_repair  = var.node_pools["primary"].auto_repair
    auto_upgrade = var.node_pools["primary"].auto_upgrade
  }

  # Upgrade settings
  upgrade_settings {
    max_surge       = 1
    max_unavailable = 0
  }

  lifecycle {
    ignore_changes = [initial_node_count]
  }
}

# AI workload node pool (for AI services)
resource "google_container_node_pool" "ai_workload" {
  count = var.enable_ai_node_pool ? 1 : 0

  name       = "ai-workload-pool"
  cluster    = google_container_cluster.primary.id
  location   = google_container_cluster.primary.location

  # Node pool size
  initial_node_count = var.node_pools["ai_workload"].initial_count

  autoscaling {
    min_node_count = var.node_pools["ai_workload"].min_count
    max_node_count = var.node_pools["ai_workload"].max_count
  }

  # Node configuration
  node_config {
    machine_type = var.node_pools["ai_workload"].machine_type
    
    # Service account
    service_account = google_service_account.gke_nodes.email
    
    # OAuth scopes
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    # Disk configuration
    disk_size_gb = var.node_pools["ai_workload"].disk_size_gb
    disk_type    = var.node_pools["ai_workload"].disk_type

    # Image type
    image_type = "COS_CONTAINERD"

    # Preemptible nodes for cost optimization
    preemptible = var.node_pools["ai_workload"].preemptible

    # Labels
    labels = merge(var.labels, {
      node-pool = "ai-workload"
      workload  = "ai-ml"
    })

    # Taints to ensure only AI workloads are scheduled
    taint {
      key    = "workload"
      value  = "ai"
      effect = "NO_SCHEDULE"
    }

    # Metadata
    metadata = {
      disable-legacy-endpoints = "true"
    }

    # Shielded instance config
    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    # Workload metadata config
    workload_metadata_config {
      mode = var.enable_workload_identity ? "GKE_METADATA" : "GCE_METADATA"
    }
  }

  # Management configuration
  management {
    auto_repair  = var.node_pools["ai_workload"].auto_repair
    auto_upgrade = var.node_pools["ai_workload"].auto_upgrade
  }

  # Upgrade settings
  upgrade_settings {
    max_surge       = 1
    max_unavailable = 0
  }

  lifecycle {
    ignore_changes = [initial_node_count]
  }
}