# Direito Lux - Terraform Infrastructure as Code for GCP
# Main configuration file

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

  # Backend configuration for state management
  backend "gcs" {
    bucket = "direito-lux-terraform-state"
    prefix = "infrastructure"
  }
}

# Configure the Google Cloud Provider
provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

# Local values for common configurations
locals {
  # Common labels for all resources
  common_labels = {
    project     = "direito-lux"
    environment = var.environment
    managed_by  = "terraform"
    team        = "devops"
  }

  # Network configuration
  vpc_name = "direito-lux-vpc-${var.environment}"
  
  # Cluster configuration
  cluster_name = "direito-lux-gke-${var.environment}"
  
  # Database configuration
  db_instance_name = "direito-lux-db-${var.environment}"
  
  # Redis configuration
  redis_instance_name = "direito-lux-redis-${var.environment}"
}

# Data sources
data "google_client_config" "default" {}

data "google_container_engine_versions" "gke_version" {
  location       = var.region
  version_prefix = var.gke_version
}

# Random suffix for unique resource names
resource "random_id" "suffix" {
  byte_length = 4
}