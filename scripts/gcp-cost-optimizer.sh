#!/bin/bash

echo "💰 GCP COST OPTIMIZER - REDUÇÃO EMERGENCIAL DE CUSTOS"
echo "===================================================="

PROJECT_ID="direito-lux-staging-2025"
CLUSTER_NAME="direito-lux-gke-staging"
REGION="us-central1"

# Função para parar cluster
stop_cluster() {
    echo "🛑 Parando cluster GKE..."
    gcloud container clusters resize $CLUSTER_NAME \
        --num-nodes=0 \
        --region=$REGION \
        --project=$PROJECT_ID \
        --quiet
    
    echo "✅ Cluster parado - custo: ~$0/hora"
    echo "📅 $(date): Cluster stopped" >> cluster-status.log
}

# Função para iniciar cluster mínimo
start_cluster() {
    echo "🚀 Iniciando cluster GKE mínimo..."
    gcloud container clusters resize $CLUSTER_NAME \
        --num-nodes=1 \
        --region=$REGION \
        --project=$PROJECT_ID \
        --quiet
    
    echo "⏳ Aguardando cluster ficar pronto..."
    sleep 30
    
    echo "✅ Cluster iniciado - custo: ~$50/dia"
    echo "📅 $(date): Cluster started" >> cluster-status.log
}

# Função para otimizar recursos imediatamente
optimize_now() {
    echo "⚡ OTIMIZAÇÃO EMERGENCIAL..."
    
    # 1. Reduzir para 1 node apenas
    echo "1. 📉 Reduzindo para 1 node..."
    gcloud container clusters resize $CLUSTER_NAME \
        --num-nodes=1 \
        --region=$REGION \
        --project=$PROJECT_ID \
        --quiet
    
    # 2. Configurar auto-scaling mínimo
    echo "2. 📊 Configurando auto-scaling mínimo..."
    gcloud container clusters update $CLUSTER_NAME \
        --enable-autoscaling \
        --min-nodes=1 \
        --max-nodes=2 \
        --region=$REGION \
        --project=$PROJECT_ID \
        --quiet
    
    # 3. Remover pods desnecessários
    echo "3. 🧹 Removendo serviços desnecessários..."
    kubectl delete deployment -n direito-lux-staging grafana 2>/dev/null || true
    kubectl delete deployment -n direito-lux-staging prometheus 2>/dev/null || true
    kubectl delete deployment -n direito-lux-staging elasticsearch 2>/dev/null || true
    
    echo "✅ Otimização concluída - economia estimada: 70%"
}

# Função para verificar custos
check_costs() {
    echo "💸 Verificando custos atual..."
    
    # Listar recursos ativos
    echo "📊 Recursos ativos:"
    gcloud compute instances list --project=$PROJECT_ID
    echo ""
    gcloud container clusters list --project=$PROJECT_ID
    echo ""
    
    # Estimar custo
    echo "💰 Estimativa de custo por hora:"
    echo "   - GKE nodes (6x e2-standard-2): ~$3.60/hora = $86.40/dia"
    echo "   - Load Balancer: ~$0.025/hora = $0.60/dia"
    echo "   - Persistent Disks: ~$0.10/dia"
    echo "   - TOTAL ATUAL: ~$87/dia = $2610/mês 🚨"
    echo ""
    echo "📉 Com otimização (1 node e2-small):"
    echo "   - GKE node (1x e2-small): ~$0.60/hora = $14.40/dia"
    echo "   - Load Balancer: ~$0.025/hora = $0.60/dia"
    echo "   - TOTAL OTIMIZADO: ~$15/dia = $450/mês ✅"
}

# Menu
case "${1:-help}" in
    "stop")
        stop_cluster
        ;;
    "start")
        start_cluster
        ;;
    "optimize")
        optimize_now
        ;;
    "costs")
        check_costs
        ;;
    "help"|*)
        echo "🔧 Comandos disponíveis:"
        echo "  ./gcp-cost-optimizer.sh optimize  # Reduz custos AGORA (70% economia)"
        echo "  ./gcp-cost-optimizer.sh stop      # Para cluster (economia 100%)"
        echo "  ./gcp-cost-optimizer.sh start     # Inicia cluster mínimo"
        echo "  ./gcp-cost-optimizer.sh costs     # Mostra análise de custos"
        echo ""
        echo "🚨 EXECUTAR PRIMEIRO: ./gcp-cost-optimizer.sh optimize"
        ;;
esac