# 🏗️ Direito Lux - Terraform Infrastructure as Code

Complete Infrastructure as Code (IaC) setup for deploying the Direito Lux SaaS platform on Google Cloud Platform using Terraform.

## 📋 Overview

This Terraform configuration provides a production-ready infrastructure for the Direito Lux legal platform, including:

- **Google Kubernetes Engine (GKE)** - Container orchestration
- **Cloud SQL PostgreSQL** - Primary database with read replicas
- **Redis** - Caching and session storage
- **VPC Networking** - Secure network isolation
- **Load Balancing** - Global HTTPS load balancer with SSL
- **DNS Management** - Cloud DNS with automatic SSL certificates
- **Monitoring & Logging** - Comprehensive observability stack
- **Security** - IAM, network policies, and encryption

## 🗂️ Directory Structure

```
terraform/
├── README.md                    # This documentation
├── main.tf                      # Main Terraform configuration
├── variables.tf                 # Variable definitions
├── outputs.tf                   # Output definitions
├── infrastructure.tf            # Infrastructure resources
├── deploy.sh                    # Deployment automation script
├── modules/                     # Reusable Terraform modules
│   ├── networking/             # VPC, subnets, firewall rules
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── gke/                    # Kubernetes cluster configuration
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── database/               # Cloud SQL and Redis
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
└── environments/               # Environment-specific configurations
    ├── staging.tfvars         # Staging environment config
    └── production.tfvars      # Production environment config
```

## 🚀 Quick Start

### Prerequisites

1. **Google Cloud SDK**: Install and authenticate
   ```bash
   gcloud auth login
   gcloud auth application-default login
   ```

2. **Terraform**: Version >= 1.0
   ```bash
   # macOS
   brew install terraform
   
   # Linux
   wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
   unzip terraform_1.6.0_linux_amd64.zip
   sudo mv terraform /usr/local/bin/
   ```

3. **Required Permissions**: Your GCP account needs:
   - Project Editor or Owner
   - Billing Account User (if creating projects)

### 🔧 Initial Setup

1. **Clone and navigate**:
   ```bash
   git clone <repository-url>
   cd direito-lux/terraform
   ```

2. **Make deployment script executable**:
   ```bash
   chmod +x deploy.sh
   ```

3. **Initialize infrastructure**:
   ```bash
   # For staging
   ./deploy.sh staging init
   
   # For production
   ./deploy.sh production init
   ```

### 🚀 Deployment

1. **Plan deployment**:
   ```bash
   ./deploy.sh staging plan
   ```

2. **Apply infrastructure**:
   ```bash
   ./deploy.sh staging apply
   ```

3. **View outputs**:
   ```bash
   ./deploy.sh staging output
   ```

## 🛠️ Deployment Script Usage

The `deploy.sh` script provides complete automation for infrastructure management:

```bash
./deploy.sh [environment] [action] [options]
```

### Parameters

- **Environment**: `staging` | `production`
- **Action**: `init` | `plan` | `apply` | `destroy` | `output`
- **Options**: `--auto-approve` (skip confirmation prompts)

### Examples

```bash
# Initialize staging environment
./deploy.sh staging init

# Plan changes for production
./deploy.sh production plan

# Apply with auto-approval
./deploy.sh staging apply --auto-approve

# Show current outputs
./deploy.sh production output

# Destroy infrastructure (careful!)
./deploy.sh staging destroy
```

## 🏗️ Infrastructure Components

### Networking

- **VPC**: Private network with multiple subnets
- **Subnets**: Segmented networks for different tiers
  - Public subnet (load balancers, NAT gateway)
  - Private subnet (application servers)
  - Database subnet (Cloud SQL, Redis)
  - GKE subnet (Kubernetes nodes with secondary ranges)
- **Firewall Rules**: Restrictive ingress/egress rules
- **NAT Gateway**: Outbound internet access for private resources
- **Private Google Access**: Access to Google APIs without public IPs

### Kubernetes (GKE)

- **Regional Cluster**: High availability across zones
- **Private Nodes**: No public IP addresses
- **Node Pools**:
  - Primary pool: General workloads
  - AI workload pool: High-memory instances for ML
- **Security Features**:
  - Workload Identity
  - Network Policy (Calico)
  - Shielded GKE nodes
  - Binary Authorization
- **Auto-scaling**: Horizontal Pod Autoscaler and Cluster Autoscaler

### Database

- **Cloud SQL PostgreSQL**:
  - Regional high availability (production)
  - Automated backups with point-in-time recovery
  - Read replicas for reporting
  - Private IP only
  - Encryption at rest and in transit
- **Redis**:
  - Standard HA tier (production) / Basic tier (staging)
  - AUTH enabled with automatic rotation
  - Persistence configuration

### Load Balancing & SSL

- **Global Load Balancer**: HTTP(S) with SSL termination
- **Managed SSL Certificates**: Automatic provisioning and renewal
- **Cloud DNS**: Authoritative DNS with health checks
- **CDN Integration**: Global content delivery

### Security

- **IAM**: Principle of least privilege
- **Service Accounts**: Dedicated accounts per service
- **Workload Identity**: Secure GKE to GCP API access
- **Secret Manager**: Encrypted secret storage
- **KMS**: Customer-managed encryption keys
- **Network Policies**: Microsegmentation
- **Private Clusters**: No public endpoints

### Monitoring & Observability

- **Cloud Monitoring**: Metrics, alerting, and dashboards
- **Cloud Logging**: Centralized log aggregation
- **BigQuery**: Log analytics and retention
- **Notification Channels**: Email and Slack alerts
- **Health Checks**: Application and infrastructure monitoring

## 🌍 Environment Configuration

### Staging Environment

- **Purpose**: Development, testing, and staging
- **Resources**: Cost-optimized with preemptible instances
- **High Availability**: Single zone deployment
- **Monitoring**: Basic monitoring and alerting
- **Retention**: 7-day backup and log retention

### Production Environment

- **Purpose**: Live production workloads
- **Resources**: Performance-optimized instances
- **High Availability**: Multi-zone with read replicas
- **Monitoring**: Comprehensive monitoring and alerting
- **Retention**: 30-day backup and 90-day log retention
- **Security**: Enhanced security policies and network restrictions

## 📊 Resource Specifications

### Staging Environment

| Component | Specification | Quantity |
|-----------|---------------|----------|
| **GKE Nodes** | e2-standard-2 (2 vCPU, 8GB) | 2-6 nodes |
| **Database** | db-custom-1-3840 (1 vCPU, 3.75GB) | 1 instance |
| **Redis** | 2GB Basic tier | 1 instance |
| **Storage** | 50GB SSD | Per instance |
| **Network** | 10.0.0.0/16 VPC | 1 VPC |

### Production Environment

| Component | Specification | Quantity |
|-----------|---------------|----------|
| **GKE Nodes** | e2-standard-4 (4 vCPU, 16GB) | 5-20 nodes |
| **AI Nodes** | n1-highmem-8 (8 vCPU, 52GB) | 2-10 nodes |
| **Database** | db-custom-4-16384 (4 vCPU, 16GB) | 1 primary + 1 replica |
| **Redis** | 8GB Standard HA tier | 1 instance + replica |
| **Storage** | 200GB SSD | Per instance |
| **Network** | 10.0.0.0/16 VPC | 1 VPC |

## 🔒 Security Features

### Network Security

- **Private GKE Cluster**: No public IP addresses on nodes
- **Authorized Networks**: Restricted master API access
- **Network Policies**: Pod-to-pod communication rules
- **VPC-native Networking**: IP alias ranges for pods and services
- **Private Google Access**: Google APIs without internet routing

### Identity & Access Management

- **Service Accounts**: Minimal privilege principle
- **Workload Identity**: Secure pod-to-GCP API access
- **IAM Roles**: Fine-grained permissions
- **Resource-level Permissions**: Service-specific access controls

### Data Security

- **Encryption at Rest**: All data encrypted with Google-managed keys
- **Encryption in Transit**: TLS 1.2+ for all communications
- **Secret Management**: Automatic secret rotation
- **Database Security**: Private IP, SSL connections, audit logging

## 📈 Monitoring & Alerting

### Key Metrics

- **Infrastructure**: CPU, memory, disk, network utilization
- **Application**: Response times, error rates, throughput
- **Database**: Connection count, query performance, replication lag
- **Security**: Failed authentication attempts, unauthorized access

### Alert Policies

- **Critical**: Service downtime, database failures
- **Warning**: High resource utilization, slow response times
- **Info**: Deployment events, configuration changes

### Dashboards

- **Infrastructure Overview**: Cluster health and resource usage
- **Application Performance**: Service-level metrics
- **Database Health**: Connection pools and query performance
- **Security Dashboard**: Access patterns and threats

## 💰 Cost Optimization

### Staging Cost Optimizations

- **Preemptible Instances**: 80% cost reduction for non-critical workloads
- **Right-sizing**: Minimal instance sizes for development
- **Auto-scaling**: Scale to zero during off-hours
- **Short Retention**: Reduced backup and log retention periods

### Production Optimizations

- **Sustained Use Discounts**: Automatic discounts for long-running instances
- **Committed Use Discounts**: 1-year commitments for predictable workloads
- **Regional Persistent Disks**: Cost-effective storage with replication
- **Log-based Metrics**: Reduce monitoring costs with selective metrics

## 🚨 Disaster Recovery

### Backup Strategy

- **Database Backups**: Daily automated backups with point-in-time recovery
- **Configuration Backup**: Terraform state stored in versioned GCS buckets
- **Application Data**: Regular exports to Cloud Storage
- **Cross-region Replication**: Critical data replicated to secondary region

### Recovery Procedures

1. **Database Recovery**: Point-in-time restore from automated backups
2. **Infrastructure Recovery**: Re-deploy from Terraform configuration
3. **Application Recovery**: Deploy from container registry
4. **Data Recovery**: Restore from cross-region backups

### Recovery Time Objectives (RTO)

- **Staging**: 4 hours
- **Production**: 1 hour for infrastructure, 15 minutes for applications

### Recovery Point Objectives (RPO)

- **Database**: 5 minutes (point-in-time recovery)
- **Application Data**: 1 hour (backup frequency)

## 🔧 Maintenance

### Regular Tasks

- **Security Updates**: Automated OS and container updates
- **Certificate Renewal**: Automatic SSL certificate renewal
- **Backup Verification**: Monthly backup restore tests
- **Capacity Planning**: Quarterly resource utilization reviews

### Upgrade Procedures

1. **GKE Upgrades**: Rolling updates during maintenance windows
2. **Database Upgrades**: Blue-green deployment with replica promotion
3. **Infrastructure Updates**: Terraform plan and apply with testing
4. **Application Deployments**: GitOps-based continuous deployment

## 🐛 Troubleshooting

### Common Issues

1. **Quota Exceeded**:
   ```bash
   # Check quotas
   gcloud compute project-info describe --project=PROJECT_ID
   
   # Request quota increase
   gcloud alpha compute quotas list --filter="service=compute.googleapis.com"
   ```

2. **Permission Denied**:
   ```bash
   # Check IAM permissions
   gcloud projects get-iam-policy PROJECT_ID
   
   # Add required roles
   gcloud projects add-iam-policy-binding PROJECT_ID \
     --member="user:email@domain.com" \
     --role="roles/editor"
   ```

3. **Network Connectivity**:
   ```bash
   # Test connectivity
   gcloud compute ssh INSTANCE_NAME --zone=ZONE
   
   # Check firewall rules
   gcloud compute firewall-rules list
   ```

### Debug Commands

```bash
# Terraform debugging
export TF_LOG=DEBUG
terraform apply

# GKE cluster info
kubectl cluster-info
kubectl get nodes -o wide

# Database connectivity
gcloud sql connect INSTANCE_NAME --user=USERNAME

# Network debugging
gcloud compute routes list
gcloud compute networks list
```

## 📚 Additional Resources

- [Terraform Google Provider Documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [Google Cloud Architecture Center](https://cloud.google.com/architecture)
- [GKE Best Practices](https://cloud.google.com/kubernetes-engine/docs/best-practices)
- [Cloud SQL Best Practices](https://cloud.google.com/sql/docs/postgres/best-practices)

## 🤝 Contributing

1. Follow existing code structure and naming conventions
2. Update variable descriptions and validation rules
3. Test changes in staging before production
4. Document any new features or significant changes
5. Update cost estimates for new resources

## 📄 License

This infrastructure configuration is proprietary to Direito Lux and is not licensed for external use.

---

**⚠️ Important**: Always test infrastructure changes in staging before applying to production. Keep Terraform state files secure and use remote state storage.