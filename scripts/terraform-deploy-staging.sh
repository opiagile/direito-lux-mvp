#!/bin/bash

# =============================================================================
# Script de Deploy Terraform - Direito Lux Staging
# =============================================================================

set -e

PROJECT_ID="direito-lux-staging"
ENVIRONMENT="staging"

echo "🚀 DEPLOY TERRAFORM - DIREITO LUX STAGING"
echo "========================================"

echo "📋 1. Verificando pré-requisitos..."
echo "   Projeto: $(gcloud config get-value project)"
echo "   Terraform: $(terraform version -json | jq -r '.terraform_version')"

echo ""
echo "🔐 2. Verificando autenticação..."
if gcloud auth application-default print-access-token >/dev/null 2>&1; then
    echo "   ✅ Application Default Credentials OK"
else
    echo "   ❌ Execute: gcloud auth application-default login"
    exit 1
fi

echo ""
echo "🏗️  3. Inicializando Terraform..."
cd /Users/franc/Opiagile/SAAS/direito-lux/terraform
terraform init

echo ""
echo "📋 4. Validando configuração..."
terraform validate

echo ""
echo "🎯 5. Executando plan..."
terraform plan -var-file=environments/staging.tfvars -out=staging.tfplan

echo ""
echo "💰 6. Estimativa de custos:"
echo "   - GKE Cluster: ~$73/mês"
echo "   - Cloud SQL: ~$25/mês"
echo "   - Redis: ~$15/mês"
echo "   - Load Balancer: ~$18/mês"
echo "   - TOTAL: ~$130/mês"

echo ""
echo "🚨 7. PRONTO PARA APPLY!"
echo "   Execute: terraform apply staging.tfplan"
echo "   Ou: terraform apply -var-file=environments/staging.tfvars"

echo ""
echo "🌍 8. Após deploy, configure DNS:"
echo "   staging.direitolux.com.br → Load Balancer IP"

echo ""
echo "✅ Deploy plan completo!"