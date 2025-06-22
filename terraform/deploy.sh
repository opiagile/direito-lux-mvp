#!/bin/bash

# Direito Lux - Terraform Deployment Script
# Usage: ./deploy.sh [staging|production] [plan|apply|destroy|output] [--auto-approve]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENVIRONMENT=${1:-staging}
ACTION=${2:-plan}
AUTO_APPROVE=${3}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Validate environment
if [[ "$ENVIRONMENT" != "staging" && "$ENVIRONMENT" != "production" ]]; then
    echo -e "${RED}❌ Invalid environment. Use 'staging' or 'production'${NC}"
    exit 1
fi

# Validate action
if [[ "$ACTION" != "plan" && "$ACTION" != "apply" && "$ACTION" != "destroy" && "$ACTION" != "output" && "$ACTION" != "init" ]]; then
    echo -e "${RED}❌ Invalid action. Use 'init', 'plan', 'apply', 'destroy', or 'output'${NC}"
    exit 1
fi

echo -e "${BLUE}🚀 Direito Lux Terraform Deployment${NC}"
echo -e "${BLUE}===================================${NC}"
echo -e "Environment: ${YELLOW}$ENVIRONMENT${NC}"
echo -e "Action: ${YELLOW}$ACTION${NC}"
echo ""

# Function to check prerequisites
check_prerequisites() {
    echo -e "${YELLOW}🔍 Checking prerequisites...${NC}"
    
    # Check Terraform
    if ! command -v terraform &> /dev/null; then
        echo -e "${RED}❌ Terraform is not installed${NC}"
        exit 1
    fi
    
    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        echo -e "${RED}❌ Google Cloud SDK is not installed${NC}"
        exit 1
    fi
    
    # Check if authenticated
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        echo -e "${RED}❌ Not authenticated with Google Cloud. Run: gcloud auth login${NC}"
        exit 1
    fi
    
    # Check if required files exist
    if [[ ! -f "$SCRIPT_DIR/environments/${ENVIRONMENT}.tfvars" ]]; then
        echo -e "${RED}❌ Environment file not found: environments/${ENVIRONMENT}.tfvars${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Prerequisites check completed${NC}"
    echo ""
}

# Function to setup backend
setup_backend() {
    echo -e "${YELLOW}🔧 Setting up Terraform backend...${NC}"
    
    # Get project ID from tfvars file
    PROJECT_ID=$(grep '^project_id' "$SCRIPT_DIR/environments/${ENVIRONMENT}.tfvars" | cut -d'"' -f2)
    
    if [[ -z "$PROJECT_ID" ]]; then
        echo -e "${RED}❌ Could not extract project_id from tfvars file${NC}"
        exit 1
    fi
    
    # Set current project
    gcloud config set project "$PROJECT_ID"
    
    # Check if state bucket exists, create if not
    BUCKET_NAME="direito-lux-terraform-state-${ENVIRONMENT}"
    
    if ! gsutil ls -b "gs://${BUCKET_NAME}" &> /dev/null; then
        echo -e "${BLUE}📦 Creating Terraform state bucket: ${BUCKET_NAME}${NC}"
        
        # Create bucket
        gsutil mb -p "$PROJECT_ID" -l us-central1 "gs://${BUCKET_NAME}"
        
        # Enable versioning
        gsutil versioning set on "gs://${BUCKET_NAME}"
        
        # Set lifecycle policy
        cat > "${SCRIPT_DIR}/lifecycle.json" << EOF
{
  "rule": [
    {
      "action": {"type": "Delete"},
      "condition": {
        "age": 30,
        "isLive": false
      }
    }
  ]
}
EOF
        
        gsutil lifecycle set "${SCRIPT_DIR}/lifecycle.json" "gs://${BUCKET_NAME}"
        rm "${SCRIPT_DIR}/lifecycle.json"
        
        echo -e "${GREEN}✅ Terraform state bucket created${NC}"
    else
        echo -e "${GREEN}✅ Terraform state bucket exists${NC}"
    fi
    
    # Update backend configuration
    sed -i.bak "s/bucket = \".*\"/bucket = \"${BUCKET_NAME}\"/" "$SCRIPT_DIR/main.tf"
    sed -i.bak "s/prefix = \".*\"/prefix = \"${ENVIRONMENT}\/infrastructure\"/" "$SCRIPT_DIR/main.tf"
    
    echo ""
}

# Function to enable required APIs
enable_apis() {
    echo -e "${YELLOW}🔌 Enabling required Google Cloud APIs...${NC}"
    
    # Get project ID
    PROJECT_ID=$(grep '^project_id' "$SCRIPT_DIR/environments/${ENVIRONMENT}.tfvars" | cut -d'"' -f2)
    
    # Required APIs
    APIS=(
        "compute.googleapis.com"
        "container.googleapis.com"
        "servicenetworking.googleapis.com"
        "sqladmin.googleapis.com"
        "redis.googleapis.com"
        "secretmanager.googleapis.com"
        "monitoring.googleapis.com"
        "logging.googleapis.com"
        "cloudresourcemanager.googleapis.com"
        "iam.googleapis.com"
        "artifactregistry.googleapis.com"
        "dns.googleapis.com"
        "cloudkms.googleapis.com"
        "bigquery.googleapis.com"
        "storage.googleapis.com"
    )
    
    # Enable APIs
    for api in "${APIS[@]}"; do
        echo -e "${BLUE}  Enabling ${api}...${NC}"
        gcloud services enable "$api" --project="$PROJECT_ID" --quiet
    done
    
    echo -e "${GREEN}✅ APIs enabled${NC}"
    echo ""
}

# Function to initialize Terraform
terraform_init() {
    echo -e "${YELLOW}🏗️ Initializing Terraform...${NC}"
    
    cd "$SCRIPT_DIR"
    
    # Initialize with backend config
    terraform init \
        -reconfigure \
        -upgrade
    
    echo -e "${GREEN}✅ Terraform initialized${NC}"
    echo ""
}

# Function to validate Terraform configuration
terraform_validate() {
    echo -e "${YELLOW}✅ Validating Terraform configuration...${NC}"
    
    cd "$SCRIPT_DIR"
    
    terraform validate
    
    echo -e "${GREEN}✅ Terraform configuration is valid${NC}"
    echo ""
}

# Function to run Terraform plan
terraform_plan() {
    echo -e "${YELLOW}📋 Running Terraform plan...${NC}"
    
    cd "$SCRIPT_DIR"
    
    terraform plan \
        -var-file="environments/${ENVIRONMENT}.tfvars" \
        -out="${ENVIRONMENT}.tfplan"
    
    echo -e "${GREEN}✅ Terraform plan completed${NC}"
    echo ""
}

# Function to run Terraform apply
terraform_apply() {
    echo -e "${YELLOW}🚀 Applying Terraform configuration...${NC}"
    
    cd "$SCRIPT_DIR"
    
    # Check if plan file exists
    if [[ ! -f "${ENVIRONMENT}.tfplan" ]]; then
        echo -e "${YELLOW}⚠️  No plan file found. Running plan first...${NC}"
        terraform_plan
    fi
    
    # Apply configuration
    if [[ "$AUTO_APPROVE" == "--auto-approve" ]]; then
        terraform apply "${ENVIRONMENT}.tfplan"
    else
        echo -e "${YELLOW}🤔 Do you want to apply these changes? (y/N)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            terraform apply "${ENVIRONMENT}.tfplan"
        else
            echo -e "${YELLOW}⏸️  Apply cancelled${NC}"
            exit 0
        fi
    fi
    
    echo -e "${GREEN}✅ Terraform apply completed${NC}"
    echo ""
}

# Function to run Terraform destroy
terraform_destroy() {
    echo -e "${RED}💥 Destroying Terraform infrastructure...${NC}"
    echo -e "${RED}⚠️  This will destroy ALL infrastructure in ${ENVIRONMENT}!${NC}"
    
    if [[ "$AUTO_APPROVE" != "--auto-approve" ]]; then
        echo -e "${YELLOW}🤔 Are you sure you want to destroy? Type 'yes' to confirm:${NC}"
        read -r response
        if [[ "$response" != "yes" ]]; then
            echo -e "${YELLOW}⏸️  Destroy cancelled${NC}"
            exit 0
        fi
    fi
    
    cd "$SCRIPT_DIR"
    
    terraform destroy \
        -var-file="environments/${ENVIRONMENT}.tfvars" \
        ${AUTO_APPROVE:+--auto-approve}
    
    echo -e "${GREEN}✅ Terraform destroy completed${NC}"
    echo ""
}

# Function to show outputs
terraform_output() {
    echo -e "${YELLOW}📊 Terraform outputs...${NC}"
    
    cd "$SCRIPT_DIR"
    
    terraform output -json > "${ENVIRONMENT}_outputs.json"
    terraform output
    
    echo ""
    echo -e "${BLUE}💾 Outputs saved to: ${ENVIRONMENT}_outputs.json${NC}"
    echo ""
}

# Function to get kubeconfig
get_kubeconfig() {
    echo -e "${YELLOW}🔧 Configuring kubectl...${NC}"
    
    # Get project ID and cluster info from outputs
    PROJECT_ID=$(terraform output -raw project_id 2>/dev/null || echo "")
    CLUSTER_NAME=$(terraform output -raw gke_cluster_name 2>/dev/null || echo "")
    REGION=$(terraform output -raw region 2>/dev/null || echo "")
    
    if [[ -n "$PROJECT_ID" && -n "$CLUSTER_NAME" && -n "$REGION" ]]; then
        echo -e "${BLUE}📡 Getting GKE credentials...${NC}"
        gcloud container clusters get-credentials "$CLUSTER_NAME" \
            --region="$REGION" \
            --project="$PROJECT_ID"
        
        echo -e "${GREEN}✅ kubectl configured for ${CLUSTER_NAME}${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not get cluster information from outputs${NC}"
    fi
    
    echo ""
}

# Function to show connection information
show_connection_info() {
    echo -e "${BLUE}🔗 Connection Information${NC}"
    echo -e "${BLUE}=========================${NC}"
    
    cd "$SCRIPT_DIR"
    
    # Get URLs from outputs
    if terraform output service_urls &>/dev/null; then
        echo -e "${GREEN}🌐 Service URLs:${NC}"
        terraform output -json service_urls | jq -r 'to_entries[] | "  \(.key): \(.value)"'
        echo ""
    fi
    
    # Get database info
    if terraform output connection_info &>/dev/null; then
        echo -e "${GREEN}🗄️  Database Connection:${NC}"
        echo -e "  Host: $(terraform output -raw postgres_private_ip 2>/dev/null || echo 'N/A')"
        echo -e "  Database: $(terraform output -raw postgres_database_name 2>/dev/null || echo 'N/A')"
        echo -e "  User: $(terraform output -raw postgres_user_name 2>/dev/null || echo 'N/A')"
        echo ""
    fi
    
    # Get cluster info
    if terraform output gke_cluster_name &>/dev/null; then
        echo -e "${GREEN}☸️  Kubernetes Cluster:${NC}"
        echo -e "  Name: $(terraform output -raw gke_cluster_name 2>/dev/null || echo 'N/A')"
        echo -e "  Location: $(terraform output -raw gke_cluster_location 2>/dev/null || echo 'N/A')"
        echo ""
    fi
}

# Main execution
main() {
    echo -e "${BLUE}Starting deployment process...${NC}"
    echo ""
    
    check_prerequisites
    
    case $ACTION in
        "init")
            setup_backend
            enable_apis
            terraform_init
            terraform_validate
            ;;
        "plan")
            terraform_init
            terraform_validate
            terraform_plan
            ;;
        "apply")
            terraform_init
            terraform_validate
            terraform_plan
            terraform_apply
            terraform_output
            get_kubeconfig
            show_connection_info
            ;;
        "destroy")
            terraform_destroy
            ;;
        "output")
            terraform_output
            show_connection_info
            ;;
    esac
    
    echo -e "${GREEN}🎉 Operation completed successfully!${NC}"
    
    if [[ "$ACTION" == "apply" ]]; then
        echo -e "${BLUE}💡 Next steps:${NC}"
        echo -e "  1. Deploy Kubernetes applications: cd ../k8s && ./deploy.sh $ENVIRONMENT --apply"
        echo -e "  2. Update DNS records if needed"
        echo -e "  3. Configure monitoring and alerting"
        echo -e "  4. Run application database migrations"
        echo ""
    fi
}

# Error handling
trap 'echo -e "${RED}❌ Deployment failed!${NC}"; exit 1' ERR

# Run main function
main

echo -e "${BLUE}🏁 Script execution completed${NC}"