# Database Module Outputs

# PostgreSQL outputs
output "postgres_instance_name" {
  description = "Name of the PostgreSQL instance"
  value       = google_sql_database_instance.postgres.name
}

output "postgres_instance_connection_name" {
  description = "Connection name of the PostgreSQL instance"
  value       = google_sql_database_instance.postgres.connection_name
}

output "postgres_instance_ip_address" {
  description = "IP address of the PostgreSQL instance"
  value       = google_sql_database_instance.postgres.ip_address
  sensitive   = true
}

output "postgres_instance_private_ip_address" {
  description = "Private IP address of the PostgreSQL instance"
  value       = google_sql_database_instance.postgres.private_ip_address
  sensitive   = true
}

output "postgres_database_name" {
  description = "Name of the main database"
  value       = google_sql_database.main_db.name
}

output "postgres_user_name" {
  description = "Name of the database user"
  value       = google_sql_user.app_user.name
}

output "postgres_user_password" {
  description = "Password of the database user"
  value       = google_sql_user.app_user.password
  sensitive   = true
}

output "postgres_replica_connection_name" {
  description = "Connection name of the PostgreSQL read replica"
  value       = var.environment == "production" ? google_sql_database_instance.postgres_replica[0].connection_name : null
}

output "postgres_replica_ip_address" {
  description = "IP address of the PostgreSQL read replica"
  value       = var.environment == "production" ? google_sql_database_instance.postgres_replica[0].ip_address : null
  sensitive   = true
}

# Redis outputs
output "redis_instance_id" {
  description = "ID of the Redis instance"
  value       = google_redis_instance.cache.id
}

output "redis_instance_name" {
  description = "Name of the Redis instance"
  value       = google_redis_instance.cache.name
}

output "redis_host" {
  description = "Host of the Redis instance"
  value       = google_redis_instance.cache.host
  sensitive   = true
}

output "redis_port" {
  description = "Port of the Redis instance"
  value       = google_redis_instance.cache.port
}

output "redis_auth_string" {
  description = "Auth string for Redis instance"
  value       = var.redis_config.auth_enabled ? google_redis_instance.cache.auth_string : null
  sensitive   = true
}

output "redis_current_location_id" {
  description = "Current location of the Redis instance"
  value       = google_redis_instance.cache.current_location_id
}

# Service account outputs
output "sql_proxy_service_account_email" {
  description = "Email of the Cloud SQL Proxy service account"
  value       = google_service_account.sql_proxy.email
}

output "sql_proxy_service_account_name" {
  description = "Name of the Cloud SQL Proxy service account"
  value       = google_service_account.sql_proxy.name
}

# Secret Manager outputs
output "db_password_secret_id" {
  description = "Secret Manager secret ID for database password"
  value       = google_secret_manager_secret.db_password.secret_id
}

output "redis_auth_secret_id" {
  description = "Secret Manager secret ID for Redis auth string"
  value       = var.redis_config.auth_enabled ? google_secret_manager_secret.redis_auth[0].secret_id : null
}

# Connection strings
output "postgres_connection_string" {
  description = "PostgreSQL connection string"
  value       = "postgresql://${google_sql_user.app_user.name}:${google_sql_user.app_user.password}@${google_sql_database_instance.postgres.private_ip_address}:5432/${google_sql_database.main_db.name}"
  sensitive   = true
}

output "redis_connection_string" {
  description = "Redis connection string"
  value       = var.redis_config.auth_enabled ? "redis://:${google_redis_instance.cache.auth_string}@${google_redis_instance.cache.host}:${google_redis_instance.cache.port}" : "redis://${google_redis_instance.cache.host}:${google_redis_instance.cache.port}"
  sensitive   = true
}

# Database URLs for applications
output "database_urls" {
  description = "Database connection URLs for different environments"
  value = {
    postgres_primary = "postgresql://${google_sql_user.app_user.name}:${google_sql_user.app_user.password}@${google_sql_database_instance.postgres.private_ip_address}:5432/${google_sql_database.main_db.name}"
    postgres_replica = var.environment == "production" ? "postgresql://${google_sql_user.app_user.name}:${google_sql_user.app_user.password}@${google_sql_database_instance.postgres_replica[0].private_ip_address}:5432/${google_sql_database.main_db.name}" : null
    redis           = var.redis_config.auth_enabled ? "redis://:${google_redis_instance.cache.auth_string}@${google_redis_instance.cache.host}:${google_redis_instance.cache.port}" : "redis://${google_redis_instance.cache.host}:${google_redis_instance.cache.port}"
  }
  sensitive = true
}