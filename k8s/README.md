# 🚀 Direito Lux - Kubernetes Infrastructure

Complete Kubernetes deployment configuration for the Direito Lux SaaS legal platform.

## 📋 Overview

This directory contains all Kubernetes manifests and deployment scripts for both staging and production environments of the Direito Lux platform.

### 🏗️ Architecture

- **Multi-tenant SaaS Platform**: Legal process management system
- **Microservices Architecture**: 9 independent services
- **High Availability**: Auto-scaling, load balancing, health checks
- **Security**: Network policies, RBAC, secret management
- **Observability**: Prometheus, Grafana, distributed tracing

## 📁 Directory Structure

```
k8s/
├── README.md                    # This file
├── deploy.sh                    # Main deployment script
├── namespace.yaml               # Kubernetes namespaces
├── databases/                   # Database configurations
│   ├── postgres.yaml           # PostgreSQL deployment
│   ├── redis.yaml              # Redis cache deployment
│   └── rabbitmq.yaml           # RabbitMQ message broker
├── services/                    # Microservice deployments
│   ├── auth-service.yaml       # Authentication service
│   ├── tenant-service.yaml     # Multi-tenancy management
│   ├── process-service.yaml    # Legal process management
│   ├── ai-service.yaml         # AI/ML processing
│   ├── notification-service.yaml # Multi-channel notifications
│   ├── search-service.yaml     # Search and indexing
│   ├── report-service.yaml     # Analytics and reporting
│   ├── mcp-service.yaml        # Claude MCP integration
│   ├── datajud-service.yaml    # Court data integration
│   └── frontend.yaml           # Next.js web application
├── ingress/                     # Load balancer configuration
│   └── ingress.yaml            # NGINX ingress with SSL
├── security/                    # Security policies
│   └── network-policies.yaml  # Network segmentation
└── monitoring/                  # Observability stack
    └── prometheus.yaml         # Metrics and alerting
```

## 🚀 Quick Start

### Prerequisites

1. **Kubernetes Cluster**: GKE, EKS, AKS, or self-managed
2. **kubectl**: Kubernetes CLI tool
3. **NGINX Ingress Controller**: For load balancing
4. **cert-manager**: For SSL certificate management

### 🔧 Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd direito-lux/k8s
   ```

2. **Make deployment script executable**:
   ```bash
   chmod +x deploy.sh
   ```

3. **Deploy to staging**:
   ```bash
   ./deploy.sh staging --apply
   ```

4. **Deploy to production**:
   ```bash
   ./deploy.sh production --apply
   ```

## 🛠️ Deployment Script Usage

The `deploy.sh` script provides a complete deployment automation solution:

```bash
./deploy.sh [staging|production] [--apply|--delete|--dry-run]
```

### Options

- **Environment**:
  - `staging`: Development/testing environment
  - `production`: Live production environment

- **Actions**:
  - `--dry-run`: Validate configurations without applying (default)
  - `--apply`: Deploy resources to the cluster
  - `--delete`: Remove all resources from the cluster

### Examples

```bash
# Dry run for staging (validate only)
./deploy.sh staging --dry-run

# Deploy to staging
./deploy.sh staging --apply

# Deploy to production
./deploy.sh production --apply

# Delete staging environment
./deploy.sh staging --delete
```

## 🌐 Service Architecture

### Core Services

| Service | Port | Description |
|---------|------|-------------|
| **auth-service** | 8080 | JWT authentication, user management |
| **tenant-service** | 8080 | Multi-tenant isolation, billing |
| **process-service** | 8080 | Legal process lifecycle management |
| **ai-service** | 8000 | AI/ML processing, document analysis |
| **notification-service** | 8080 | Email, SMS, WhatsApp, Telegram |
| **search-service** | 8080 | Elasticsearch integration |
| **report-service** | 8080 | Analytics, dashboards, exports |
| **mcp-service** | 8080 | Claude MCP tool integration |
| **datajud-service** | 8080 | Court system data integration |
| **frontend** | 3000 | Next.js web application |

### Infrastructure Services

| Service | Port | Description |
|---------|------|-------------|
| **PostgreSQL** | 5432 | Primary database |
| **Redis** | 6379 | Caching and sessions |
| **RabbitMQ** | 5672/15672 | Message broker |
| **Prometheus** | 9090 | Metrics collection |
| **Grafana** | 3000 | Monitoring dashboards |

## 📊 Resource Requirements

### Staging Environment

| Component | CPU Request | CPU Limit | Memory Request | Memory Limit |
|-----------|-------------|-----------|----------------|--------------|
| PostgreSQL | 250m | 500m | 256Mi | 1Gi |
| Redis | 50m | 200m | 64Mi | 256Mi |
| RabbitMQ | 100m | 300m | 256Mi | 512Mi |
| Go Services | 100m | 200m | 128Mi | 256Mi |
| AI Service | 500m | 2000m | 1Gi | 4Gi |
| Frontend | 100m | 300m | 256Mi | 512Mi |

### Production Environment

| Component | CPU Request | CPU Limit | Memory Request | Memory Limit |
|-----------|-------------|-----------|----------------|--------------|
| PostgreSQL | 500m | 1000m | 512Mi | 2Gi |
| Redis | 100m | 500m | 256Mi | 1Gi |
| RabbitMQ | 200m | 500m | 512Mi | 1Gi |
| Go Services | 200m | 500m | 256Mi | 512Mi |
| AI Service | 1000m | 4000m | 2Gi | 8Gi |
| Frontend | 200m | 500m | 512Mi | 1Gi |

## 🔒 Security Features

### Network Policies

- **Database Isolation**: Only backend services can access databases
- **Service Mesh**: Controlled inter-service communication
- **Ingress Control**: Frontend accessible only via load balancer
- **Default Deny**: All traffic denied by default

### Secret Management

- **Database Credentials**: Stored in Kubernetes secrets
- **API Keys**: Encrypted at rest and in transit
- **JWT Secrets**: Environment-specific signing keys
- **TLS Certificates**: Automated with cert-manager

### RBAC Configuration

- **Service Accounts**: Minimal required permissions
- **Role-Based Access**: Environment-specific roles
- **Namespace Isolation**: Staging and production separation

## 📈 Auto-Scaling Configuration

### Horizontal Pod Autoscaler (HPA)

All services include HPA configuration:

- **Staging**: 2-10 replicas based on CPU/memory usage
- **Production**: 3-20 replicas based on CPU/memory usage
- **Metrics**: CPU utilization (70%), Memory utilization (80%)

### Scaling Policies

```yaml
# Example HPA configuration
minReplicas: 3
maxReplicas: 20
metrics:
- type: Resource
  resource:
    name: cpu
    target:
      type: Utilization
      averageUtilization: 70
```

## 🌍 DNS and SSL Configuration

### Domain Setup

**Staging**:
- Frontend: `https://staging.direitolux.com`
- API: `https://api-staging.direitolux.com`

**Production**:
- Frontend: `https://app.direitolux.com`
- API: `https://api.direitolux.com`

### SSL Certificates

- **cert-manager**: Automated Let's Encrypt certificates
- **NGINX Ingress**: SSL termination and HTTP/2
- **HSTS**: HTTP Strict Transport Security headers

## 📊 Monitoring and Observability

### Prometheus Metrics

- **Service Metrics**: Response times, error rates, throughput
- **Infrastructure Metrics**: CPU, memory, disk, network
- **Business Metrics**: User registrations, process counts
- **Database Metrics**: Connection pools, query performance

### Grafana Dashboards

- **Service Overview**: Health status of all services
- **Performance**: Response times and error rates
- **Infrastructure**: Resource utilization and capacity
- **Business KPIs**: Key performance indicators

### Alerting Rules

- **Service Down**: Critical alert when service unavailable
- **High CPU/Memory**: Warning when resources exceed thresholds
- **Database Issues**: Connection problems or slow queries
- **API Performance**: Response times above acceptable limits

## 🚨 Health Checks and Probes

### Liveness Probes

- **HTTP Health Checks**: `/health` endpoint for all services
- **Database Connections**: `pg_isready`, `redis-cli ping`
- **Failure Threshold**: 3 consecutive failures trigger restart

### Readiness Probes

- **Service Dependencies**: Check database and cache connectivity
- **Initialization**: Ensure service is ready to handle requests
- **Load Balancer**: Only route traffic to ready pods

## 🔄 Deployment Strategies

### Rolling Updates

- **Zero Downtime**: Gradual replacement of old pods
- **Health Checks**: Only replace healthy pods
- **Rollback**: Automatic rollback on failure

### Blue-Green Deployment

For critical production updates:

1. Deploy new version alongside current
2. Switch traffic gradually
3. Monitor health and metrics
4. Complete switch or rollback

## 🛠️ Troubleshooting

### Common Issues

1. **Pod CrashLoopBackOff**:
   ```bash
   kubectl logs <pod-name> -n direito-lux-staging
   kubectl describe pod <pod-name> -n direito-lux-staging
   ```

2. **Service Not Accessible**:
   ```bash
   kubectl get endpoints -n direito-lux-staging
   kubectl get ingress -n direito-lux-staging
   ```

3. **Database Connection Issues**:
   ```bash
   kubectl exec -it postgres-pod -n direito-lux-staging -- psql -U direito_lux
   ```

### Debug Commands

```bash
# Check all pods status
kubectl get pods -n direito-lux-staging

# View service logs
kubectl logs -f deployment/auth-service -n direito-lux-staging

# Check resource usage
kubectl top pods -n direito-lux-staging

# Debug networking
kubectl exec -it <pod-name> -n direito-lux-staging -- /bin/sh
```

## 🔧 Maintenance Tasks

### Database Backups

```bash
# PostgreSQL backup
kubectl exec postgres-pod -n direito-lux-production -- pg_dump -U direito_lux direito_lux_production > backup.sql
```

### Scaling Services

```bash
# Manual scaling
kubectl scale deployment auth-service --replicas=5 -n direito-lux-production
```

### Certificate Renewal

```bash
# Check certificate status
kubectl get certificates -n direito-lux-production

# Force renewal
kubectl delete certificate direito-lux-production-tls -n direito-lux-production
```

## 📚 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [cert-manager Documentation](https://cert-manager.io/docs/)
- [Prometheus Operator](https://prometheus-operator.dev/)

## 🤝 Contributing

1. Follow the existing naming conventions
2. Update resource limits based on actual usage
3. Test changes in staging before production
4. Document any configuration changes
5. Update monitoring and alerting rules

## 📄 License

This configuration is part of the Direito Lux platform and is proprietary to the organization.

---

**Note**: Remember to update secrets, domain names, and resource requirements according to your specific environment and requirements.