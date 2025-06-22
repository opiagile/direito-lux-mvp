# Database Module Variables

variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "zone" {
  description = "GCP zone"
  type        = string
}

variable "secondary_region" {
  description = "Secondary region for replicas"
  type        = string
}

variable "db_instance_name" {
  description = "Name of the database instance"
  type        = string
}

variable "redis_instance_name" {
  description = "Name of the Redis instance"
  type        = string
}

variable "database_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "POSTGRES_15"
}

variable "vpc_id" {
  description = "VPC network ID"
  type        = string
}

variable "private_vpc_connection" {
  description = "Private VPC connection for Cloud SQL"
  type        = any
}

variable "database_config" {
  description = "Database configuration"
  type = object({
    tier                = string
    disk_size           = number
    disk_type           = string
    availability_type   = string
    backup_enabled      = bool
    backup_start_time   = string
    maintenance_window_day  = number
    maintenance_window_hour = number
    deletion_protection = bool
  })
}

variable "redis_config" {
  description = "Redis configuration"
  type = object({
    memory_size_gb     = number
    tier               = string
    redis_version      = string
    auth_enabled       = bool
    transit_encryption = bool
    replica_count      = number
  })
}

variable "backup_config" {
  description = "Backup configuration"
  type = object({
    enabled                     = bool
    retention_period_days       = number
    automated_backup_enabled    = bool
    point_in_time_recovery_enabled = bool
  })
}

variable "authorized_networks" {
  description = "Authorized networks for database access"
  type = list(object({
    cidr_block   = string
    display_name = string
  }))
  default = []
}

variable "notification_channels" {
  description = "Notification channels for alerts"
  type        = list(string)
  default     = []
}

variable "labels" {
  description = "Labels to apply to resources"
  type        = map(string)
  default     = {}
}