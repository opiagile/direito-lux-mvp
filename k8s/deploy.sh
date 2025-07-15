#!/bin/bash

# Direito Lux Kubernetes Deployment Script
# Usage: ./deploy.sh [staging|production] [--apply|--delete|--dry-run]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENVIRONMENT=${1:-staging}
ACTION=${2:---dry-run}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Validate environment
if [[ "$ENVIRONMENT" != "staging" && "$ENVIRONMENT" != "production" ]]; then
    echo -e "${RED}❌ Invalid environment. Use 'staging' or 'production'${NC}"
    exit 1
fi

# Validate action
if [[ "$ACTION" != "--apply" && "$ACTION" != "--delete" && "$ACTION" != "--dry-run" ]]; then
    echo -e "${RED}❌ Invalid action. Use '--apply', '--delete', or '--dry-run'${NC}"
    exit 1
fi

echo -e "${BLUE}🚀 Direito Lux Kubernetes Deployment${NC}"
echo -e "${BLUE}====================================${NC}"
echo -e "Environment: ${YELLOW}$ENVIRONMENT${NC}"
echo -e "Action: ${YELLOW}$ACTION${NC}"
echo ""

# Function to check prerequisites
check_prerequisites() {
    echo -e "${YELLOW}🔍 Checking prerequisites...${NC}"
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is not installed${NC}"
        exit 1
    fi
    
    # Check cluster connection
    if ! kubectl cluster-info &> /dev/null; then
        echo -e "${RED}❌ Cannot connect to Kubernetes cluster${NC}"
        exit 1
    fi
    
    # Check if namespace exists
    NAMESPACE="direito-lux-$ENVIRONMENT"
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null && [[ "$ACTION" != "--delete" ]]; then
        echo -e "${YELLOW}⚠️  Namespace $NAMESPACE does not exist${NC}"
        if [[ "$ACTION" == "--apply" ]]; then
            echo -e "${BLUE}📝 Creating namespace...${NC}"
            kubectl apply -f "$SCRIPT_DIR/namespace.yaml"
        fi
    fi
    
    echo -e "${GREEN}✅ Prerequisites check completed${NC}"
    echo ""
}

# Function to deploy resources in order
deploy_resources() {
    echo -e "${YELLOW}🔧 Deploying resources...${NC}"
    
    # Define deployment order
    DEPLOY_ORDER=(
        "namespace.yaml"
        "databases/postgres.yaml"
        "databases/redis.yaml"
        "databases/rabbitmq.yaml"
        "services/auth-service.yaml"
        "services/tenant-service.yaml"
        "services/ai-service.yaml"
        "services/frontend.yaml"
        "ingress/ingress.yaml"
        "security/network-policies.yaml"
        "monitoring/prometheus.yaml"
    )
    
    for resource in "${DEPLOY_ORDER[@]}"; do
        file_path="$SCRIPT_DIR/$resource"
        
        if [[ -f "$file_path" ]]; then
            echo -e "${BLUE}📄 Processing $resource...${NC}"
            
            case $ACTION in
                "--apply")
                    kubectl apply -f "$file_path"
                    ;;
                "--delete")
                    kubectl delete -f "$file_path" --ignore-not-found=true
                    ;;
                "--dry-run")
                    kubectl apply -f "$file_path" --dry-run=client
                    ;;
            esac
            
            echo -e "${GREEN}✅ $resource processed${NC}"
        else
            echo -e "${YELLOW}⚠️  File not found: $resource${NC}"
        fi
        
        echo ""
    done
}

# Function to wait for deployments
wait_for_deployments() {
    if [[ "$ACTION" == "--apply" ]]; then
        echo -e "${YELLOW}⏳ Waiting for deployments to be ready...${NC}"
        
        NAMESPACE="direito-lux-$ENVIRONMENT"
        
        # Wait for databases first
        echo -e "${BLUE}📊 Waiting for databases...${NC}"
        kubectl wait --for=condition=available --timeout=300s deployment/postgres -n "$NAMESPACE" || true
        kubectl wait --for=condition=available --timeout=300s deployment/redis -n "$NAMESPACE" || true
        kubectl wait --for=condition=available --timeout=300s deployment/rabbitmq -n "$NAMESPACE" || true
        
        # Wait for services
        echo -e "${BLUE}🔧 Waiting for services...${NC}"
        kubectl wait --for=condition=available --timeout=300s deployment/auth-service -n "$NAMESPACE" || true
        kubectl wait --for=condition=available --timeout=300s deployment/tenant-service -n "$NAMESPACE" || true
        kubectl wait --for=condition=available --timeout=300s deployment/ai-service -n "$NAMESPACE" || true
        kubectl wait --for=condition=available --timeout=300s deployment/frontend -n "$NAMESPACE" || true
        
        echo -e "${GREEN}✅ Deployments are ready${NC}"
    fi
}

# Function to show deployment status
show_status() {
    if [[ "$ACTION" != "--delete" ]]; then
        echo -e "${YELLOW}📊 Deployment Status${NC}"
        echo -e "${YELLOW}===================${NC}"
        
        NAMESPACE="direito-lux-$ENVIRONMENT"
        
        echo -e "${BLUE}Pods:${NC}"
        kubectl get pods -n "$NAMESPACE" -o wide
        echo ""
        
        echo -e "${BLUE}Services:${NC}"
        kubectl get services -n "$NAMESPACE"
        echo ""
        
        echo -e "${BLUE}Ingress:${NC}"
        kubectl get ingress -n "$NAMESPACE"
        echo ""
        
        if [[ "$ENVIRONMENT" == "staging" ]]; then
            echo -e "${GREEN}🌐 Staging URLs:${NC}"
            echo -e "  Frontend: https://staging.direitolux.com"
            echo -e "  API: https://api-staging.direitolux.com"
        else
            echo -e "${GREEN}🌐 Production URLs:${NC}"
            echo -e "  Frontend: https://app.direitolux.com"
            echo -e "  API: https://api.direitolux.com"
        fi
        echo ""
    fi
}

# Function to run health checks
run_health_checks() {
    if [[ "$ACTION" == "--apply" ]]; then
        echo -e "${YELLOW}🏥 Running health checks...${NC}"
        
        NAMESPACE="direito-lux-$ENVIRONMENT"
        
        # Check if all pods are running
        echo -e "${BLUE}📋 Checking pod status...${NC}"
        kubectl get pods -n "$NAMESPACE" --field-selector=status.phase!=Running --no-headers | wc -l | {
            read count
            if [[ "$count" -gt 0 ]]; then
                echo -e "${YELLOW}⚠️  $count pods are not running${NC}"
                kubectl get pods -n "$NAMESPACE" --field-selector=status.phase!=Running
            else
                echo -e "${GREEN}✅ All pods are running${NC}"
            fi
        }
        
        # Check service endpoints
        echo -e "${BLUE}🔗 Checking service endpoints...${NC}"
        services=("auth-service" "tenant-service" "ai-service" "frontend")
        for service in "${services[@]}"; do
            if kubectl get endpoints "$service" -n "$NAMESPACE" &> /dev/null; then
                endpoint_count=$(kubectl get endpoints "$service" -n "$NAMESPACE" -o jsonpath='{.subsets[0].addresses}' | jq length 2>/dev/null || echo "0")
                if [[ "$endpoint_count" -gt 0 ]]; then
                    echo -e "${GREEN}✅ $service has $endpoint_count endpoints${NC}"
                else
                    echo -e "${YELLOW}⚠️  $service has no ready endpoints${NC}"
                fi
            else
                echo -e "${RED}❌ $service endpoints not found${NC}"
            fi
        done
        
        echo ""
    fi
}

# Main execution
main() {
    echo -e "${BLUE}Starting deployment process...${NC}"
    echo ""
    
    check_prerequisites
    deploy_resources
    
    if [[ "$ACTION" == "--apply" ]]; then
        wait_for_deployments
        show_status
        run_health_checks
        
        echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
        echo -e "${BLUE}💡 Next steps:${NC}"
        echo -e "  1. Update DNS records to point to the load balancer"
        echo -e "  2. Configure SSL certificates"
        echo -e "  3. Run database migrations"
        echo -e "  4. Set up monitoring alerts"
        echo ""
    elif [[ "$ACTION" == "--delete" ]]; then
        echo -e "${GREEN}🧹 Resources deleted successfully!${NC}"
        echo ""
    else
        echo -e "${GREEN}👀 Dry run completed successfully!${NC}"
        echo -e "${BLUE}💡 To apply changes, run:${NC}"
        echo -e "  ./deploy.sh $ENVIRONMENT --apply"
        echo ""
    fi
}

# Error handling
trap 'echo -e "${RED}❌ Deployment failed!${NC}"; exit 1' ERR

# Run main function
main

echo -e "${BLUE}🏁 Script execution completed${NC}"