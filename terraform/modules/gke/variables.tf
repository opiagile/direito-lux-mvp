# GKE Module Variables

variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "cluster_name" {
  description = "Name of the GKE cluster"
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

variable "regional_cluster" {
  description = "Create a regional cluster"
  type        = bool
  default     = true
}

variable "vpc_name" {
  description = "Name of the VPC network"
  type        = string
}

variable "subnet_name" {
  description = "Name of the subnet for GKE nodes"
  type        = string
}

variable "pods_range_name" {
  description = "Name of the secondary range for pods"
  type        = string
}

variable "services_range_name" {
  description = "Name of the secondary range for services"
  type        = string
}

variable "enable_private_nodes" {
  description = "Enable private nodes"
  type        = bool
  default     = true
}

variable "master_ipv4_cidr_block" {
  description = "CIDR block for GKE master nodes"
  type        = string
  default     = "172.16.0.0/28"
}

variable "authorized_networks" {
  description = "Authorized networks for GKE master access"
  type = list(object({
    cidr_block   = string
    display_name = string
  }))
  default = []
}

variable "enable_workload_identity" {
  description = "Enable Workload Identity"
  type        = bool
  default     = true
}

variable "enable_network_policy" {
  description = "Enable network policy"
  type        = bool
  default     = true
}

variable "enable_pod_security_policy" {
  description = "Enable pod security policy"
  type        = bool
  default     = false
}

variable "database_encryption_key" {
  description = "KMS key for database encryption"
  type        = string
  default     = null
}

variable "release_channel" {
  description = "GKE release channel"
  type        = string
  default     = "REGULAR"
  
  validation {
    condition     = contains(["RAPID", "REGULAR", "STABLE"], var.release_channel)
    error_message = "Release channel must be RAPID, REGULAR, or STABLE."
  }
}

variable "node_pools" {
  description = "Node pool configurations"
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
    taints = list(object({
      key    = string
      value  = string
      effect = string
    }))
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
      taints         = []
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
      taints         = []
    }
  }
}

variable "enable_ai_node_pool" {
  description = "Enable AI workload node pool"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Labels to apply to resources"
  type        = map(string)
  default     = {}
}