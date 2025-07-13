#!/bin/bash

# =============================================================================
# Script para configurar projeto GCP Staging - Direito Lux
# =============================================================================

set -e  # Exit on any error

BILLING_ACCOUNT="01B2F9-AD5BB4-BE339E"
PROJECT_ID="direito-lux-staging"
REGION="us-central1"
ZONE="us-central1-a"

echo "🚀 Configurando projeto Direito Lux Staging..."
echo "=============================================="

echo "📝 1. Criando projeto '$PROJECT_ID'..."
if gcloud projects describe $PROJECT_ID >/dev/null 2>&1; then
    echo "   ✅ Projeto já existe!"
else
    gcloud projects create $PROJECT_ID --name="Direito Lux Staging"
    echo "   ✅ Projeto criado!"
fi

echo ""
echo "💳 2. Vinculando billing account..."
gcloud beta billing projects link $PROJECT_ID --billing-account=$BILLING_ACCOUNT

echo ""
echo "🎯 3. Definindo como projeto ativo..."
gcloud config set project $PROJECT_ID

echo ""
echo "🌍 4. Configurando região padrão..."
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE

echo ""
echo "⚡ 5. Habilitando APIs essenciais..."
apis=(
    "compute.googleapis.com"
    "container.googleapis.com" 
    "sqladmin.googleapis.com"
    "redis.googleapis.com"
    "dns.googleapis.com"
    "certificatemanager.googleapis.com"
    "cloudresourcemanager.googleapis.com"
    "iam.googleapis.com"
    "logging.googleapis.com"
    "monitoring.googleapis.com"
)

for api in "${apis[@]}"; do
    echo "   🔌 Habilitando $api..."
    gcloud services enable $api
done

echo ""
echo "📋 6. Verificando configuração final..."
echo "   Projeto: $(gcloud config get-value project)"
echo "   Região: $(gcloud config get-value compute/region)"
echo "   Zone: $(gcloud config get-value compute/zone)"

echo ""
echo "🎉 Setup staging completo! Pronto para Terraform."