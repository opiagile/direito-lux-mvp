# Direito Lux - Terraform Variables
# Configuration variables for infrastructure deployment

# Project Configuration
variable "project_id" {
  description = "GCP Project ID"
  type        = string
  validation {
    condition     = length(var.project_id) > 0
    error_message = "Project ID cannot be empty."
  }
}

variable "organization_id" {
  description = "GCP Organization ID"
  type        = string
  default     = ""
}

variable "billing_account" {
  description = "GCP Billing Account ID"
  type        = string
  default     = ""
}

# Environment Configuration
variable "environment" {
  description = "Environment name (staging, production)"
  type        = string
  default     = "staging"
  
  validation {
    condition     = contains(["staging", "production"], var.environment)
    error_message = "Environment must be either 'staging' or 'production'."
  }
}

# Region and Zone Configuration
variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "GCP Zone"
  type        = string
  default     = "us-central1-a"
}

variable "secondary_region" {
  description = "Secondary region for disaster recovery"
  type        = string
  default     = "us-east1"
}

# Network Configuration
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_cidrs" {
  description = "CIDR blocks for subnets"
  type = object({
    public    = string
    private   = string
    database  = string
    gke_pods  = string
    gke_services = string
  })
  default = {
    public    = "10.0.1.0/24"
    private   = "10.0.2.0/24"
    database  = "10.0.3.0/24"
    gke_pods  = "10.1.0.0/16"
    gke_services = "10.2.0.0/16"
  }
}

# GKE Configuration
variable "gke_version" {
  description = "GKE cluster version"
  type        = string
  default     = "1.28."
}

variable "gke_node_pools" {
  description = "GKE node pool configurations"
  type = map(object({
    machine_type   = string
    min_count      = number
    max_count      = number
    initial_count  = number
    disk_size_gb   = number
    disk_type      = string
    preemptible    = bool
    auto_upgrade   = bool
    auto_repair    = bool
  }))
  default = {
    primary = {
      machine_type   = "e2-standard-4"
      min_count      = 1
      max_count      = 10
      initial_count  = 3
      disk_size_gb   = 100
      disk_type      = "pd-ssd"
      preemptible    = false
      auto_upgrade   = true
      auto_repair    = true
    }
    ai_workload = {
      machine_type   = "n1-highmem-4"
      min_count      = 0
      max_count      = 5
      initial_count  = 1
      disk_size_gb   = 200
      disk_type      = "pd-ssd"
      preemptible    = true
      auto_upgrade   = true
      auto_repair    = true
    }
  }
}

variable "enable_private_nodes" {
  description = "Enable private GKE nodes"
  type        = bool
  default     = true
}

variable "master_ipv4_cidr_block" {
  description = "CIDR block for GKE master nodes"
  type        = string
  default     = "172.16.0.0/28"
}

# Database Configuration
variable "database_config" {
  description = "Cloud SQL database configuration"
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
  default = {
    tier                = "db-custom-2-8192"  # 2 vCPU, 8GB RAM
    disk_size           = 100
    disk_type           = "PD_SSD"
    availability_type   = "REGIONAL"
    backup_enabled      = true
    backup_start_time   = "03:00"
    maintenance_window_day  = 7  # Sunday
    maintenance_window_hour = 4  # 4 AM
    deletion_protection = true
  }
}

variable "database_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "POSTGRES_15"
}

# Redis Configuration
variable "redis_config" {
  description = "Redis instance configuration"
  type = object({
    memory_size_gb     = number
    tier               = string
    redis_version      = string
    auth_enabled       = bool
    transit_encryption = bool
    replica_count      = number
  })
  default = {
    memory_size_gb     = 4
    tier               = "STANDARD_HA"
    redis_version      = "REDIS_7_0"
    auth_enabled       = true
    transit_encryption = true
    replica_count      = 1
  }
}

# Load Balancer Configuration
variable "enable_ssl" {
  description = "Enable SSL certificates"
  type        = bool
  default     = true
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
  default     = "direitolux.com"
}

variable "subdomain_staging" {
  description = "Subdomain for staging environment"
  type        = string
  default     = "staging"
}

variable "subdomain_api" {
  description = "Subdomain for API"
  type        = string
  default     = "api"
}

# Monitoring Configuration
variable "enable_monitoring" {
  description = "Enable Google Cloud Monitoring"
  type        = bool
  default     = true
}

variable "enable_logging" {
  description = "Enable Google Cloud Logging"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "Log retention period in days"
  type        = number
  default     = 30
}

# Security Configuration
variable "enable_network_policy" {
  description = "Enable network policy for GKE"
  type        = bool
  default     = true
}

variable "enable_pod_security_policy" {
  description = "Enable pod security policy"
  type        = bool
  default     = true
}

variable "authorized_networks" {
  description = "Authorized networks for GKE master access"
  type = list(object({
    cidr_block   = string
    display_name = string
  }))
  default = [
    {
      cidr_block   = "0.0.0.0/0"
      display_name = "All networks"
    }
  ]
}

# Backup Configuration
variable "backup_config" {
  description = "Backup configuration"
  type = object({
    enabled                     = bool
    retention_period_days       = number
    automated_backup_enabled    = bool
    point_in_time_recovery_enabled = bool
  })
  default = {
    enabled                     = true
    retention_period_days       = 30
    automated_backup_enabled    = true
    point_in_time_recovery_enabled = true
  }
}

# Cost Optimization
variable "preemptible_nodes" {
  description = "Use preemptible nodes for cost optimization"
  type        = bool
  default     = false
}

variable "spot_instances" {
  description = "Use spot instances for AI workloads"
  type        = bool
  default     = true
}

# Feature Flags
variable "enable_istio" {
  description = "Enable Istio service mesh"
  type        = bool
  default     = false
}

variable "enable_anthos" {
  description = "Enable Anthos features"
  type        = bool
  default     = false
}

variable "enable_workload_identity" {
  description = "Enable Workload Identity"
  type        = bool
  default     = true
}

# Application Configuration
variable "app_config" {
  description = "Application-specific configuration"
  type = object({
    ai_service_enabled    = bool
    search_service_enabled = bool
    notification_providers = list(string)
    max_tenants           = number
  })
  default = {
    ai_service_enabled    = true
    search_service_enabled = true
    notification_providers = ["email", "sms", "whatsapp", "telegram"]
    max_tenants           = 1000
  }
}

# Environment-specific overrides
variable "staging_config" {
  description = "Staging environment specific configuration"
  type = object({
    node_count     = number
    machine_type   = string
    disk_size      = number
    db_tier        = string
    redis_memory   = number
  })
  default = {
    node_count     = 3
    machine_type   = "e2-standard-2"
    disk_size      = 50
    db_tier        = "db-custom-1-3840"
    redis_memory   = 2
  }
}

variable "production_config" {
  description = "Production environment specific configuration"
  type = object({
    node_count     = number
    machine_type   = string
    disk_size      = number
    db_tier        = string
    redis_memory   = number
  })
  default = {
    node_count     = 5
    machine_type   = "e2-standard-4"
    disk_size      = 100
    db_tier        = "db-custom-4-16384"
    redis_memory   = 8
  }
}