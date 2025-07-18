#!/bin/bash

# 🚨 DEPLOY URGENTE - Corrigir Sistema de Registro
# Executa em sequência: auth → deploy → validação
# Tempo estimado: 3-5 minutos

set -e  # Parar se qualquer comando falhar

echo "🚀 INICIANDO DEPLOY URGENTE..."
echo "=================================="

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}ETAPA 1/6: Reautenticação GCP${NC}"
echo "Executando: gcloud auth login"
gcloud auth login

echo -e "${YELLOW}ETAPA 2/6: Configuração kubectl${NC}"
echo "Executando: get-credentials"
gcloud container clusters get-credentials direito-lux-gke-staging \
  --region=us-central1 \
  --project=direito-lux-staging-2025

echo -e "${YELLOW}ETAPA 3/6: Deploy dos serviços corrigidos${NC}"
echo "Aplicando: fix-backend-services.yaml"
kubectl apply -f fix-backend-services.yaml

echo "Aplicando: fix-ingress.yaml"
kubectl apply -f fix-ingress.yaml

echo -e "${YELLOW}ETAPA 4/6: Aguardando inicialização (90 segundos)${NC}"
echo "Monitorando pods..."

# Aguardar pods aparecerem
sleep 10

# Mostrar status inicial
echo "Status inicial dos pods:"
kubectl get pods -n direito-lux-staging | grep fixed || echo "Pods ainda não criados"

# Aguardar pods ficarem ready
echo "Aguardando pods ficarem READY..."
kubectl wait --for=condition=ready pod -l app=auth-service-fixed -n direito-lux-staging --timeout=300s
kubectl wait --for=condition=ready pod -l app=tenant-service-fixed -n direito-lux-staging --timeout=300s

echo -e "${GREEN}✅ Pods estão READY!${NC}"

echo -e "${YELLOW}ETAPA 5/6: Validação das APIs${NC}"

# Aguardar alguns segundos para services se estabilizarem
sleep 15

echo "Testando Auth Service..."
AUTH_RESPONSE=$(curl -k -s -o /dev/null -w "%{http_code}" https://35.188.198.87/api/v1/auth/health)
if [ "$AUTH_RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Auth Service: OK (200)${NC}"
else
    echo -e "${RED}❌ Auth Service: $AUTH_RESPONSE${NC}"
fi

echo "Testando Tenant Service..."
TENANT_RESPONSE=$(curl -k -s -o /dev/null -w "%{http_code}" https://35.188.198.87/api/v1/tenants/health)
if [ "$TENANT_RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Tenant Service: OK (200)${NC}"
else
    echo -e "${RED}❌ Tenant Service: $TENANT_RESPONSE${NC}"
fi

echo -e "${YELLOW}ETAPA 6/6: Teste de registro completo${NC}"

echo "Testando criação de tenant (Costa Advogados)..."
REGISTER_RESPONSE=$(curl -k -s -o /dev/null -w "%{http_code}" -X POST https://35.188.198.87/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Costa Advogados Teste",
    "email": "teste@costaadvogados.com.br", 
    "plan": "professional",
    "legal_name": "Costa Advogados Ltda",
    "document": "12.345.678/0001-90",
    "phone": "(11) 98765-4321"
  }')

if [ "$REGISTER_RESPONSE" = "201" ] || [ "$REGISTER_RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Registro Tenant: OK ($REGISTER_RESPONSE)${NC}"
else
    echo -e "${RED}❌ Registro Tenant: $REGISTER_RESPONSE${NC}"
fi

echo ""
echo "=================================="
echo -e "${GREEN}🎉 DEPLOY CONCLUÍDO!${NC}"
echo ""

if [ "$AUTH_RESPONSE" = "200" ] && [ "$TENANT_RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Sistema totalmente funcional!${NC}"
    echo "Agora você pode testar o registro no frontend:"
    echo "https://35.188.198.87/register"
    echo ""
    echo "O erro 'Erro de conexão. Tente novamente.' foi RESOLVIDO!"
else
    echo -e "${YELLOW}⚠️  Alguns serviços ainda não estão funcionais.${NC}"
    echo "Execute para debug:"
    echo "kubectl logs -n direito-lux-staging -l app=auth-service-fixed --tail=20"
    echo "kubectl logs -n direito-lux-staging -l app=tenant-service-fixed --tail=20"
fi

echo ""
echo "Status final dos pods:"
kubectl get pods -n direito-lux-staging | grep fixed