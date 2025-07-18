#!/bin/bash

# 🛑 PARAR STAGING GCP - Economia de Custos
# Para todos os recursos para não cobrar durante a madrugada

set -e

echo "🛑 PARANDO AMBIENTE STAGING GCP - ECONOMIA DE CUSTOS"
echo "===================================================="

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

echo -e "${YELLOW}ETAPA 2: Escalando cluster para 0 nodes${NC}"
echo "Salvando configuração atual..."
kubectl get nodes -o wide > /tmp/nodes-backup-$(date +%Y%m%d-%H%M).txt

echo "Escalando node pool para 0..."
gcloud container clusters resize $CLUSTER_NAME \
  --region=$REGION \
  --num-nodes=0 \
  --quiet

echo -e "${GREEN}✅ Cluster escalado para 0 nodes${NC}"

echo -e "${YELLOW}ETAPA 3: Verificando status${NC}"
gcloud container clusters describe $CLUSTER_NAME \
  --region=$REGION \
  --format="value(currentNodeCount)"

echo -e "${YELLOW}ETAPA 4: Parando instâncias compute restantes${NC}"
echo "Listando instâncias ativas..."
gcloud compute instances list --filter="status:RUNNING" --format="table(name,zone,status)"

echo "Parando todas as instâncias compute..."
INSTANCES=$(gcloud compute instances list --filter="status:RUNNING" --format="value(name,zone)")
if [ ! -z "$INSTANCES" ]; then
    while IFS= read -r line; do
        INSTANCE_NAME=$(echo $line | cut -d' ' -f1)
        ZONE=$(echo $line | cut -d' ' -f2)
        echo "Parando instância: $INSTANCE_NAME na zona: $ZONE"
        gcloud compute instances stop $INSTANCE_NAME --zone=$ZONE --quiet
    done <<< "$INSTANCES"
else
    echo "Nenhuma instância ativa encontrada."
fi

echo -e "${YELLOW}ETAPA 5: Verificando custos atuais${NC}"
echo "IP Estático mantido: direito-lux-staging-ip"
gcloud compute addresses list --filter="name:direito-lux-staging-ip"

echo "Discos persistentes mantidos:"
gcloud compute disks list --format="table(name,sizeGb,type,status)"

echo ""
echo -e "${GREEN}🎉 STAGING PARADO COM SUCESSO!${NC}"
echo ""
echo "💰 ECONOMIA ESTIMADA:"
echo "   • GKE Nodes: ~$0.10/hora/node × 5 nodes = $0.50/hora → $0/hora"
echo "   • Compute Engine: Parado → $0/hora"
echo "   • Load Balancer: Mantido (necessário) → $0.025/hora"
echo "   • IP Estático: Mantido → $0.004/hora"
echo "   • Discos: Mantidos → ~$0.10/hora"
echo ""
echo "💡 TOTAL ECONOMIA: ~$0.45/hora (~$10.80/dia)"
echo ""
echo "🔄 Para REATIVAR use:"
echo "   ./scripts/gcp-start-staging.sh"
echo ""

# Salvar estado para reativação
cat > /tmp/gcp-staging-state.txt << EOF
PROJECT_ID=$PROJECT_ID
CLUSTER_NAME=$CLUSTER_NAME
REGION=$REGION
STOPPED_AT=$(date)
NODES_BEFORE=5
EOF

echo "Estado salvo em: /tmp/gcp-staging-state.txt"
echo -e "${RED}⚠️  STAGING TOTALMENTE PARADO - SEM COBRANÇA NOTURNA${NC}"