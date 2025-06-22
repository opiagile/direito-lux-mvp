# Direito Lux - Networking Module
# VPC, Subnets, Firewall Rules, and NAT Gateway configuration

# VPC Network
resource "google_compute_network" "vpc" {
  name                    = var.vpc_name
  auto_create_subnetworks = false
  mtu                     = 1460
  routing_mode           = "REGIONAL"

  lifecycle {
    prevent_destroy = true
  }
}

# Public Subnet (for load balancers and NAT gateway)
resource "google_compute_subnetwork" "public" {
  name          = "${var.vpc_name}-public"
  ip_cidr_range = var.subnet_cidrs.public
  region        = var.region
  network       = google_compute_network.vpc.id

  # Enable private Google access
  private_ip_google_access = true

  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata            = "INCLUDE_ALL_METADATA"
  }
}

# Private Subnet (for application servers)
resource "google_compute_subnetwork" "private" {
  name          = "${var.vpc_name}-private"
  ip_cidr_range = var.subnet_cidrs.private
  region        = var.region
  network       = google_compute_network.vpc.id

  # Enable private Google access
  private_ip_google_access = true

  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata            = "INCLUDE_ALL_METADATA"
  }
}

# Database Subnet (for Cloud SQL and Redis)
resource "google_compute_subnetwork" "database" {
  name          = "${var.vpc_name}-database"
  ip_cidr_range = var.subnet_cidrs.database
  region        = var.region
  network       = google_compute_network.vpc.id

  # Enable private Google access
  private_ip_google_access = true

  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata            = "INCLUDE_ALL_METADATA"
  }
}

# GKE Subnet with secondary ranges
resource "google_compute_subnetwork" "gke" {
  name          = "${var.vpc_name}-gke"
  ip_cidr_range = var.subnet_cidrs.private
  region        = var.region
  network       = google_compute_network.vpc.id

  # Enable private Google access
  private_ip_google_access = true

  # Secondary IP ranges for GKE pods and services
  secondary_ip_range {
    range_name    = "gke-pods"
    ip_cidr_range = var.subnet_cidrs.gke_pods
  }

  secondary_ip_range {
    range_name    = "gke-services"
    ip_cidr_range = var.subnet_cidrs.gke_services
  }

  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata            = "INCLUDE_ALL_METADATA"
  }
}

# Cloud Router for NAT Gateway
resource "google_compute_router" "router" {
  name    = "${var.vpc_name}-router"
  region  = var.region
  network = google_compute_network.vpc.id

  bgp {
    asn = 64514
  }
}

# NAT Gateway for outbound internet access from private subnets
resource "google_compute_router_nat" "nat" {
  name                               = "${var.vpc_name}-nat"
  router                             = google_compute_router.router.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

# Firewall Rules

# Allow internal communication
resource "google_compute_firewall" "allow_internal" {
  name    = "${var.vpc_name}-allow-internal"
  network = google_compute_network.vpc.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  source_ranges = [
    var.subnet_cidrs.public,
    var.subnet_cidrs.private,
    var.subnet_cidrs.database,
    var.subnet_cidrs.gke_pods,
    var.subnet_cidrs.gke_services
  ]
}

# Allow SSH from authorized networks
resource "google_compute_firewall" "allow_ssh" {
  name    = "${var.vpc_name}-allow-ssh"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = var.authorized_networks
  target_tags   = ["ssh-allowed"]
}

# Allow HTTP and HTTPS from internet
resource "google_compute_firewall" "allow_http_https" {
  name    = "${var.vpc_name}-allow-http-https"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["http-server", "https-server"]
}

# Allow health checks from Google Cloud Load Balancer
resource "google_compute_firewall" "allow_health_check" {
  name    = "${var.vpc_name}-allow-health-check"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["8080", "8000", "3000", "9090"]
  }

  # Health check source ranges
  source_ranges = [
    "130.211.0.0/22",
    "35.191.0.0/16"
  ]

  target_tags = ["allow-health-check"]
}

# Allow NodePort services (for GKE)
resource "google_compute_firewall" "allow_nodeport" {
  name    = "${var.vpc_name}-allow-nodeport"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["gke-node"]
}

# Deny all other inbound traffic (default)
resource "google_compute_firewall" "deny_all" {
  name     = "${var.vpc_name}-deny-all"
  network  = google_compute_network.vpc.name
  priority = 65534

  deny {
    protocol = "tcp"
  }

  deny {
    protocol = "udp"
  }

  deny {
    protocol = "icmp"
  }

  source_ranges = ["0.0.0.0/0"]
}

# Allow specific database ports
resource "google_compute_firewall" "allow_database" {
  name    = "${var.vpc_name}-allow-database"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["5432", "6379", "5672", "15672"] # PostgreSQL, Redis, RabbitMQ
  }

  source_ranges = [
    var.subnet_cidrs.private,
    var.subnet_cidrs.gke_pods
  ]

  target_tags = ["database"]
}

# Private Service Connection for Cloud SQL
resource "google_compute_global_address" "private_ip_address" {
  name          = "${var.vpc_name}-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

# DNS Configuration
resource "google_dns_managed_zone" "private_zone" {
  name        = "${var.environment}-private-zone"
  dns_name    = "${var.environment}.local."
  description = "Private DNS zone for ${var.environment} environment"

  visibility = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.vpc.id
    }
  }
}

# Add DNS records for services
resource "google_dns_record_set" "services" {
  for_each = var.service_dns_records

  name         = "${each.key}.${google_dns_managed_zone.private_zone.dns_name}"
  managed_zone = google_dns_managed_zone.private_zone.name
  type         = "A"
  ttl          = 300

  rrdatas = [each.value]
}