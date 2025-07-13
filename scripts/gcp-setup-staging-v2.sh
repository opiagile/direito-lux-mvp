#!/bin/bash

# =============================================================================
# Script para configurar projeto GCP Staging - Direito Lux (v2 - Sem billing)
# =============================================================================

set -e  # Exit on any error

PROJECT_ID="direito-lux-staging"
REGION="us-central1"
ZONE="us-central1-a"

echo "🚀 Configurando projeto Direito Lux Staging (v2)..."
echo "================================================="

echo "🎯 1. Definindo como projeto ativo..."
gcloud config set project $PROJECT_ID

echo ""
echo "🌍 2. Configurando região padrão..."
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE

echo ""
echo "⚡ 3. Habilitando APIs essenciais (pode demorar)..."
apis=(
    "compute.googleapis.com"
    "container.googleapis.com" 
    "sqladmin.googleapis.com"
    "redis.googleapis.com"
    "dns.googleapis.com"
    "cloudresourcemanager.googleapis.com"
    "iam.googleapis.com"
)

for api in "${apis[@]}"; do
    echo "   🔌 Habilitando $api..."
    if gcloud services enable $api 2>/dev/null; then
        echo "   ✅ $api habilitada!"
    else
        echo "   ⚠️  $api - permissão necessária (configure no console)"
    fi
done

echo ""
echo "📋 4. Verificando configuração final..."
echo "   Projeto: $(gcloud config get-value project)"
echo "   Região: $(gcloud config get-value compute/region)"
echo "   Zone: $(gcloud config get-value compute/zone)"

echo ""
echo "💡 BILLING: Vincule manualmente no console GCP:"
echo "   https://console.cloud.google.com/billing/linkedaccount?project=$PROJECT_ID"

echo ""
echo "✅ Setup básico completo!"