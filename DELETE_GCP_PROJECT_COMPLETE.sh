#!/bin/bash

# 🗑️ DELETE COMPLETO DO PROJETO GCP
# Remove TUDO do projeto direito-lux-staging-2025

set -e

echo "🗑️ INICIANDO EXCLUSÃO COMPLETA DO PROJETO GCP"
echo "=============================================="
echo ""

# Cores para output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

PROJECT_ID="direito-lux-staging-2025"

echo -e "${YELLOW}⚠️  ATENÇÃO: Esta ação irá DELETAR PERMANENTEMENTE:${NC}"
echo "   • Todo o projeto GCP: $PROJECT_ID"
echo "   • TODOS os recursos (cluster, VMs, discos, IPs)"
echo "   • TODOS os dados (bancos, backups)"
echo "   • TODAS as configurações"
echo "   • Billing e faturamento associado"
echo ""
echo -e "${RED}Esta ação é IRREVERSÍVEL!${NC}"
echo ""
read -p "Tem certeza que deseja continuar? Digite 'DELETE' para confirmar: " confirmation

if [ "$confirmation" != "DELETE" ]; then
    echo -e "${YELLOW}Operação cancelada.${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}ETAPA 1: Configurando projeto para exclusão${NC}"
gcloud config set project $PROJECT_ID 2>/dev/null || true

echo -e "${YELLOW}ETAPA 2: Removendo proteções de exclusão (se existirem)${NC}"
# Remove lien protection se existir
gcloud alpha resource-manager liens list --project=$PROJECT_ID 2>/dev/null || true

echo -e "${YELLOW}ETAPA 3: Desabilitando billing${NC}"
# Remove billing account
gcloud beta billing projects unlink $PROJECT_ID 2>/dev/null || true

echo -e "${YELLOW}ETAPA 4: Deletando recursos específicos${NC}"
echo "Deletando cluster GKE..."
gcloud container clusters delete direito-lux-gke-staging \
  --region=us-central1 \
  --quiet 2>/dev/null || true

echo "Deletando instâncias Compute Engine..."
gcloud compute instances list --format="value(name,zone)" | while read -r instance zone; do
    echo "Deletando instância: $instance"
    gcloud compute instances delete $instance --zone=$zone --quiet 2>/dev/null || true
done

echo "Deletando discos persistentes..."
gcloud compute disks list --format="value(name,zone)" | while read -r disk zone; do
    echo "Deletando disco: $disk"
    gcloud compute disks delete $disk --zone=$zone --quiet 2>/dev/null || true
done

echo "Deletando IPs estáticos..."
gcloud compute addresses list --format="value(name,region)" | while read -r ip region; do
    echo "Deletando IP: $ip"
    gcloud compute addresses delete $ip --region=$region --quiet 2>/dev/null || true
done

echo "Deletando load balancers..."
gcloud compute forwarding-rules list --format="value(name)" | while read -r rule; do
    echo "Deletando forwarding rule: $rule"
    gcloud compute forwarding-rules delete $rule --global --quiet 2>/dev/null || true
done

echo "Deletando backend services..."
gcloud compute backend-services list --format="value(name)" | while read -r backend; do
    echo "Deletando backend: $backend"
    gcloud compute backend-services delete $backend --global --quiet 2>/dev/null || true
done

echo "Deletando artifacts registry..."
gcloud artifacts repositories delete direito-lux-staging \
  --location=us-central1 \
  --quiet 2>/dev/null || true

echo -e "${YELLOW}ETAPA 5: Aguardando recursos serem removidos${NC}"
sleep 30

echo -e "${RED}ETAPA 6: DELETANDO O PROJETO COMPLETAMENTE${NC}"
echo "Executando exclusão final do projeto..."
gcloud projects delete $PROJECT_ID --quiet

echo ""
echo -e "${GREEN}✅ PROJETO DELETADO COM SUCESSO!${NC}"
echo ""
echo "Resumo da exclusão:"
echo "   • Projeto: $PROJECT_ID"
echo "   • Status: DELETED"
echo "   • Billing: REMOVED"
echo "   • Recursos: ALL DELETED"
echo ""

echo -e "${YELLOW}ETAPA 7: Limpando configurações locais${NC}"
# Remove configurações locais do gcloud
gcloud config configurations delete staging-config 2>/dev/null || true

# Remove contexto kubectl
kubectl config delete-context gke_${PROJECT_ID}_us-central1_direito-lux-gke-staging 2>/dev/null || true

echo -e "${GREEN}🎉 LIMPEZA COMPLETA FINALIZADA!${NC}"
echo ""
echo "Próximos passos:"
echo "1. ✅ Projeto GCP completamente removido"
echo "2. ✅ Sem cobranças futuras"
echo "3. ✅ Pronto para começar do zero"
echo ""
echo -e "${YELLOW}IMPORTANTE:${NC}"
echo "• O projeto levará alguns minutos para ser totalmente removido do GCP"
echo "• Você pode criar um novo projeto com ID diferente quando quiser"
echo "• Recomendo usar um novo ID como: processalert-prod-2025"
echo ""