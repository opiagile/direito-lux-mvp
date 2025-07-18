#!/bin/bash

# 🚀 REATIVAR STAGING GCP
# Reativa todos os recursos para continuar desenvolvimento

set -e

echo "🚀 REATIVANDO AMBIENTE STAGING GCP"
echo "=================================="

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_ID="direito-lux-staging-2025"
CLUSTER_NAME="direito-lux-gke-staging"
REGION="us-central1"

echo -e "${YELLOW}ETAPA 1: Configurando projeto${NC}"
gcloud config set project $PROJECT_ID

echo -e "${YELLOW}ETAPA 2: Reativando cluster para 5 nodes${NC}"
gcloud container clusters resize $CLUSTER_NAME \
  --region=$REGION \
  --num-nodes=5 \
  --quiet

echo -e "${GREEN}✅ Aguardando nodes ficarem prontos...${NC}"
sleep 60

echo -e "${YELLOW}ETAPA 3: Verificando status${NC}"
kubectl get nodes
kubectl get pods -n direito-lux-staging

echo -e "${YELLOW}ETAPA 4: Verificando serviços críticos${NC}"
echo "Aguardando emergency-auth-proxy..."
kubectl wait --for=condition=ready pod -l app=emergency-auth-proxy -n direito-lux-staging --timeout=300s

echo "Testando APIs:"
kubectl port-forward -n direito-lux-staging $(kubectl get pods -n direito-lux-staging -l app=emergency-auth-proxy -o name) 8080:8080 &
PORT_FORWARD_PID=$!

sleep 5
curl -s http://localhost:8080/health && echo " ✅ Auth API OK"
curl -s http://localhost:8080/api/v1/tenants/health && echo " ✅ Tenant API OK"

kill $PORT_FORWARD_PID

echo ""
echo -e "${GREEN}🎉 STAGING REATIVADO COM SUCESSO!${NC}"
echo ""
echo "🌐 URLs disponíveis:"
echo "   • Frontend: https://35.188.198.87"
echo "   • Auth API: https://35.188.198.87/api/v1/auth/health"
echo "   • Tenant API: https://35.188.198.87/api/v1/tenants/health"
echo ""
echo "🧪 Para testar via port-forward:"
echo "   kubectl port-forward -n direito-lux-staging \\$(kubectl get pods -n direito-lux-staging -l app=emergency-auth-proxy -o name) 8080:8080"
echo ""
echo "💰 CUSTO REATIVADO: ~$0.45/hora"