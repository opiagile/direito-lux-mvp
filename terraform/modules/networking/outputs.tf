# Networking Module Outputs

output "vpc_id" {
  description = "ID of the VPC network"
  value       = google_compute_network.vpc.id
}

output "vpc_name" {
  description = "Name of the VPC network"
  value       = google_compute_network.vpc.name
}

output "vpc_self_link" {
  description = "Self link of the VPC network"
  value       = google_compute_network.vpc.self_link
}

output "public_subnet_id" {
  description = "ID of the public subnet"
  value       = google_compute_subnetwork.public.id
}

output "public_subnet_name" {
  description = "Name of the public subnet"
  value       = google_compute_subnetwork.public.name
}

output "private_subnet_id" {
  description = "ID of the private subnet"
  value       = google_compute_subnetwork.private.id
}

output "private_subnet_name" {
  description = "Name of the private subnet"
  value       = google_compute_subnetwork.private.name
}

output "database_subnet_id" {
  description = "ID of the database subnet"
  value       = google_compute_subnetwork.database.id
}

output "database_subnet_name" {
  description = "Name of the database subnet"
  value       = google_compute_subnetwork.database.name
}

output "gke_subnet_id" {
  description = "ID of the GKE subnet"
  value       = google_compute_subnetwork.gke.id
}

output "gke_subnet_name" {
  description = "Name of the GKE subnet"
  value       = google_compute_subnetwork.gke.name
}

output "gke_pods_range_name" {
  description = "Name of the GKE pods secondary range"
  value       = google_compute_subnetwork.gke.secondary_ip_range[0].range_name
}

output "gke_services_range_name" {
  description = "Name of the GKE services secondary range"
  value       = google_compute_subnetwork.gke.secondary_ip_range[1].range_name
}

output "router_name" {
  description = "Name of the Cloud Router"
  value       = google_compute_router.router.name
}

output "nat_name" {
  description = "Name of the NAT Gateway"
  value       = google_compute_router_nat.nat.name
}

output "private_vpc_connection_id" {
  description = "ID of the private VPC connection"
  value       = google_service_networking_connection.private_vpc_connection.id
}

output "private_ip_address_name" {
  description = "Name of the private IP address range"
  value       = google_compute_global_address.private_ip_address.name
}

output "dns_zone_name" {
  description = "Name of the private DNS zone"
  value       = google_dns_managed_zone.private_zone.name
}

output "dns_zone_dns_name" {
  description = "DNS name of the private zone"
  value       = google_dns_managed_zone.private_zone.dns_name
}