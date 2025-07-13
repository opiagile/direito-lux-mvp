#!/bin/bash

# =============================================================================
# Script FINAL para configurar GCP - Direito Lux Staging
# =============================================================================

PROJECT_ID="direito-lux-staging"
REGION="us-central1"
ZONE="us-central1-c"

echo "🎯 CONFIGURAÇÃO FINAL GCP - DIREITO LUX"
echo "======================================"

echo "1. 🎛️  Configurando projeto ativo..."
gcloud config set project $PROJECT_ID

echo "2. 🌍 Configurando região/zona corretas..."
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE

echo "3. 📋 Status atual:"
echo "   Projeto: $(gcloud config get-value project)"
echo "   Região: $REGION" 
echo "   Zona: $ZONE"

echo ""
echo "4. 🚨 AÇÃO MANUAL NECESSÁRIA:"
echo "   ▶️  Acesse: https://console.cloud.google.com/billing/linkedaccount?project=$PROJECT_ID"
echo "   ▶️  Vincule a billing account ao projeto"
echo "   ▶️  Depois voltamos para habilitar APIs"

echo ""
echo "5. 🎯 PRÓXIMO PASSO - DEPLOY TERRAFORM:"
echo "   cd terraform/"
echo "   terraform init"
echo "   terraform plan -var-file=environments/staging.tfvars"

echo ""
echo "✅ GCP configurado! Faça billing link e vamos pro Terraform!"