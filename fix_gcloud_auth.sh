#!/bin/bash

# 🔑 CORREÇÃO DE AUTENTICAÇÃO GCLOUD
# Resolve problemas de token expirado e configura acesso ao GKE

echo "🔑 CORRIGINDO AUTENTICAÇÃO GCLOUD"
echo "=================================="

# 1. VERIFICAR STATUS ATUAL
echo "📋 Status atual:"
echo "Projeto: $(gcloud config get-value project 2>/dev/null || echo 'Não configurado')"
echo "Conta: $(gcloud config get-value account 2>/dev/null || echo 'Não configurado')"
echo "Região: $(gcloud config get-value compute/region 2>/dev/null || echo 'Não configurado')"
echo ""

# 2. FAZER LOGIN (REQUER INTERAÇÃO)
echo "🔐 Executando nova autenticação..."
echo "⚠️  ATENÇÃO: Uma janela do browser será aberta para autenticação."
echo ""

gcloud auth login

# 3. CONFIGURAR APPLICATION DEFAULT CREDENTIALS
echo "🔧 Configurando Application Default Credentials..."
gcloud auth application-default login

# 4. VERIFICAR PROJETO
echo "🎯 Configurando projeto..."
gcloud config set project direito-lux-staging-2025
gcloud config set compute/region us-central1
gcloud config set compute/zone us-central1-c

# 5. RECONFIGURAR CLUSTER GKE
echo "🚀 Reconfigurando acesso ao cluster GKE..."
gcloud container clusters get-credentials direito-lux-gke-staging \
    --region=us-central1 \
    --project=direito-lux-staging-2025

# 6. TESTAR CONECTIVIDADE
echo "🏥 Testando conectividade..."
kubectl get nodes

echo ""
echo "✅ AUTENTICAÇÃO CORRIGIDA!"
echo "========================="
echo "🎯 Agora você pode executar:"
echo "   ./deploy_gke_apps.sh"
echo ""