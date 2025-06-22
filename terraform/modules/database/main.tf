# Direito Lux - Database Module
# Cloud SQL PostgreSQL and Redis instances

# Random password for database
resource "random_password" "db_password" {
  length  = 16
  special = true
}

# Cloud SQL PostgreSQL instance
resource "google_sql_database_instance" "postgres" {
  name             = var.db_instance_name
  database_version = var.database_version
  region           = var.region

  deletion_protection = var.database_config.deletion_protection

  settings {
    tier                        = var.database_config.tier
    availability_type           = var.database_config.availability_type
    disk_type                   = var.database_config.disk_type
    disk_size                   = var.database_config.disk_size
    disk_autoresize            = true
    disk_autoresize_limit      = var.database_config.disk_size * 10

    # Database flags for optimization
    database_flags {
      name  = "max_connections"
      value = var.environment == "production" ? "400" : "200"
    }

    database_flags {
      name  = "shared_preload_libraries"
      value = "pg_stat_statements"
    }

    database_flags {
      name  = "log_statement"
      value = "all"
    }

    database_flags {
      name  = "log_min_duration_statement"
      value = "1000"  # Log queries taking longer than 1 second
    }

    # Backup configuration
    backup_configuration {
      enabled                        = var.database_config.backup_enabled
      start_time                     = var.database_config.backup_start_time
      location                       = var.region
      point_in_time_recovery_enabled = var.backup_config.point_in_time_recovery_enabled
      transaction_log_retention_days = 7

      backup_retention_settings {
        retained_backups = var.backup_config.retention_period_days
        retention_unit   = "COUNT"
      }
    }

    # Maintenance window
    maintenance_window {
      day          = var.database_config.maintenance_window_day
      hour         = var.database_config.maintenance_window_hour
      update_track = "stable"
    }

    # IP configuration for private networking
    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = var.vpc_id
      enable_private_path_for_google_cloud_services = true
      
      # Authorized networks (only for development/staging)
      dynamic "authorized_networks" {
        for_each = var.environment == "production" ? [] : var.authorized_networks
        content {
          name  = authorized_networks.value.display_name
          value = authorized_networks.value.cidr_block
        }
      }
    }

    # Insights configuration
    insights_config {
      query_insights_enabled  = true
      query_string_length     = 1024
      record_application_tags = true
      record_client_address   = true
    }

    # User labels
    user_labels = var.labels
  }

  depends_on = [var.private_vpc_connection]
}

# Main application database
resource "google_sql_database" "main_db" {
  name     = "${var.environment}_direito_lux"
  instance = google_sql_database_instance.postgres.name
  charset  = "UTF8"
  collation = "en_US.UTF8"
}

# Database user for application
resource "google_sql_user" "app_user" {
  name     = "direito_lux"
  instance = google_sql_database_instance.postgres.name
  password = random_password.db_password.result
}

# Read replica for reporting (production only)
resource "google_sql_database_instance" "postgres_replica" {
  count = var.environment == "production" ? 1 : 0

  name                 = "${var.db_instance_name}-replica"
  database_version     = var.database_version
  region               = var.secondary_region
  master_instance_name = google_sql_database_instance.postgres.name

  replica_configuration {
    failover_target = false
  }

  settings {
    tier              = var.database_config.tier
    availability_type = "ZONAL"  # Replicas are always zonal
    disk_type         = var.database_config.disk_type
    disk_size         = var.database_config.disk_size

    # IP configuration
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }

    # User labels
    user_labels = merge(var.labels, {
      role = "read-replica"
    })
  }
}

# Redis instance for caching
resource "google_redis_instance" "cache" {
  name           = var.redis_instance_name
  tier           = var.redis_config.tier
  memory_size_gb = var.redis_config.memory_size_gb
  region         = var.region

  location_id             = var.zone
  alternative_location_id = var.redis_config.tier == "STANDARD_HA" ? "${substr(var.zone, 0, length(var.zone)-1)}b" : null

  authorized_network = var.vpc_id
  redis_version      = var.redis_config.redis_version

  # Security settings
  auth_enabled                = var.redis_config.auth_enabled
  transit_encryption_mode     = var.redis_config.transit_encryption ? "SERVER_AUTHENTICATION" : "DISABLED"
  
  # Replica count for high availability
  replica_count = var.redis_config.replica_count

  # Redis configuration
  redis_configs = {
    maxmemory-policy = "allkeys-lru"
    timeout          = "300"
    tcp-keepalive    = "60"
  }

  # Maintenance policy
  maintenance_policy {
    weekly_maintenance_window {
      day = "SUNDAY"
      start_time {
        hours   = 4
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
  }

  # Persistence config (for STANDARD_HA tier)
  dynamic "persistence_config" {
    for_each = var.redis_config.tier == "STANDARD_HA" ? [1] : []
    content {
      persistence_mode    = "RDB"
      rdb_snapshot_period = "TWENTY_FOUR_HOURS"
      rdb_snapshot_start_time = "03:00"
    }
  }

  labels = var.labels
}

# Cloud SQL Proxy service account (for secure connections)
resource "google_service_account" "sql_proxy" {
  account_id   = "${var.environment}-sql-proxy"
  display_name = "Cloud SQL Proxy Service Account"
  description  = "Service account for Cloud SQL Proxy connections"
}

resource "google_project_iam_member" "sql_proxy_roles" {
  for_each = toset([
    "roles/cloudsql.client",
    "roles/cloudsql.instanceUser"
  ])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.sql_proxy.email}"
}

# Secret for database password
resource "google_secret_manager_secret" "db_password" {
  secret_id = "${var.environment}-db-password"

  labels = var.labels

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

# Secret for Redis auth string
resource "google_secret_manager_secret" "redis_auth" {
  count = var.redis_config.auth_enabled ? 1 : 0

  secret_id = "${var.environment}-redis-auth"

  labels = var.labels

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "redis_auth" {
  count = var.redis_config.auth_enabled ? 1 : 0

  secret      = google_secret_manager_secret.redis_auth[0].id
  secret_data = google_redis_instance.cache.auth_string
}

# Database monitoring and alerting
resource "google_monitoring_alert_policy" "database_cpu" {
  display_name = "${var.environment} - Database High CPU"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "Database CPU usage"
    
    condition_threshold {
      filter          = "resource.type=\"cloudsql_database\" AND resource.labels.database_id=\"${var.project_id}:${google_sql_database_instance.postgres.name}\""
      duration        = "300s"
      comparison      = "COMPARISON_GREATER_THAN"
      threshold_value = 0.8

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }

  notification_channels = var.notification_channels

  alert_strategy {
    auto_close = "1800s"
  }
}

resource "google_monitoring_alert_policy" "database_connections" {
  display_name = "${var.environment} - Database High Connections"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "Database connection count"
    
    condition_threshold {
      filter          = "resource.type=\"cloudsql_database\" AND resource.labels.database_id=\"${var.project_id}:${google_sql_database_instance.postgres.name}\""
      duration        = "300s"
      comparison      = "COMPARISON_GREATER_THAN"
      threshold_value = var.environment == "production" ? 320 : 160

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_MEAN"
      }
    }
  }

  notification_channels = var.notification_channels

  alert_strategy {
    auto_close = "1800s"
  }
}

resource "google_monitoring_alert_policy" "redis_memory" {
  display_name = "${var.environment} - Redis High Memory Usage"
  combiner     = "OR"
  enabled      = true

  conditions {
    display_name = "Redis memory usage"
    
    condition_threshold {
      filter          = "resource.type=\"redis_instance\" AND resource.labels.instance_id=\"${google_redis_instance.cache.name}\""
      duration        = "300s"
      comparison      = "COMPARISON_GREATER_THAN"
      threshold_value = 0.9

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_MEAN"
      }
    }
  }

  notification_channels = var.notification_channels

  alert_strategy {
    auto_close = "1800s"
  }
}