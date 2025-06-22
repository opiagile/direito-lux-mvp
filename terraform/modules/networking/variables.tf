# Networking Module Variables

variable "vpc_name" {
  description = "Name of the VPC network"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "subnet_cidrs" {
  description = "CIDR blocks for subnets"
  type = object({
    public       = string
    private      = string
    database     = string
    gke_pods     = string
    gke_services = string
  })
}

variable "authorized_networks" {
  description = "Authorized networks for SSH access"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "service_dns_records" {
  description = "DNS records for internal services"
  type        = map(string)
  default = {
    postgres  = "10.0.3.10"
    redis     = "10.0.3.20"
    rabbitmq  = "10.0.3.30"
  }
}

variable "labels" {
  description = "Labels to apply to resources"
  type        = map(string)
  default     = {}
}