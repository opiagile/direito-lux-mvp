#!/bin/bash

# =============================================================================
# Script de Deploy Kubernetes - Direito Lux Staging
# =============================================================================

set -e

PROJECT_ID="direito-lux-staging"
CLUSTER_NAME="direito-lux-staging-gke"
REGION="us-central1"

echo "☸️  DEPLOY KUBERNETES - DIREITO LUX STAGING"
echo "========================================="

echo "🔐 1. Conectando ao cluster GKE..."
gcloud container clusters get-credentials $CLUSTER_NAME --region=$REGION --project=$PROJECT_ID

echo ""
echo "📋 2. Verificando conexão k8s..."
kubectl cluster-info

echo ""
echo "🏗️  3. Aplicando namespace..."
kubectl apply -f /Users/franc/Opiagile/SAAS/direito-lux/k8s/namespace.yaml

echo ""
echo "💾 4. Deployando databases..."
kubectl apply -f /Users/franc/Opiagile/SAAS/direito-lux/k8s/databases/

echo ""
echo "⚡ 5. Deployando microserviços..."
kubectl apply -f /Users/franc/Opiagile/SAAS/direito-lux/k8s/services/

echo ""
echo "🌐 6. Configurando ingress..."
kubectl apply -f /Users/franc/Opiagile/SAAS/direito-lux/k8s/ingress/

echo ""
echo "📊 7. Configurando monitoring..."
kubectl apply -f /Users/franc/Opiagile/SAAS/direito-lux/k8s/monitoring/

echo ""
echo "🔐 8. Configurando secrets do GitHub..."
kubectl create secret generic github-secrets \
  --from-literal=telegram-bot-token="$TELEGRAM_BOT_TOKEN" \
  --from-literal=openai-api-key="$STAGING_OPENAI_API_KEY" \
  --from-literal=asaas-api-key="$ASAAS_SANDBOX_API_KEY" \
  --dry-run=client -o yaml | kubectl apply -f -

echo ""
echo "📋 9. Verificando status dos pods..."
kubectl get pods -n direito-lux

echo ""
echo "🌍 10. Obtendo IP do Load Balancer..."
kubectl get ingress -n direito-lux

echo ""
echo "✅ Deploy K8s completo!"
echo "🎯 Configure DNS: staging.direitolux.com.br → Load Balancer IP"