#!/bin/bash

# 🔧 CONFIGURAÇÃO DIRETA DO KUBECTL
# Configura kubectl sem usar gcloud get-credentials

echo "🔧 CONFIGURANDO KUBECTL DIRETAMENTE"
echo "==================================="

# Extrair informações do cluster via API REST
ACCESS_TOKEN=$(gcloud auth application-default print-access-token)

echo "📡 Obtendo informações do cluster..."
CLUSTER_INFO=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" \
  "https://container.googleapis.com/v1/projects/direito-lux-staging-2025/locations/us-central1/clusters/direito-lux-gke-staging")

# Extrair endpoint e certificado usando python
ENDPOINT=$(echo "$CLUSTER_INFO" | python3 -c "import json, sys; data=json.load(sys.stdin); print(data['endpoint'])")
CA_CERT=$(echo "$CLUSTER_INFO" | python3 -c "import json, sys; data=json.load(sys.stdin); print(data['masterAuth']['clusterCaCertificate'])")

echo "📡 Endpoint: $ENDPOINT"
echo "🔐 Configurando certificado..."

# Criar arquivo permanente para o certificado
mkdir -p $HOME/.kube/certs
echo "$CA_CERT" | base64 -d > $HOME/.kube/certs/cluster-ca.crt

# Configurar kubectl diretamente
echo "⚙️ Configurando kubectl..."

# Adicionar cluster
kubectl config set-cluster direito-lux-gke-staging \
  --server=https://$ENDPOINT \
  --certificate-authority=$HOME/.kube/certs/cluster-ca.crt

# Adicionar usuário com token de acesso
kubectl config set-credentials gke-user \
  --token=$ACCESS_TOKEN

# Criar contexto
kubectl config set-context direito-lux-gke-staging \
  --cluster=direito-lux-gke-staging \
  --user=gke-user

# Usar o contexto
kubectl config use-context direito-lux-gke-staging

# Testar conectividade
echo "🧪 Testando conectividade..."
if kubectl get nodes; then
    echo "✅ Kubectl configurado com sucesso!"
else
    echo "❌ Erro na configuração"
    exit 1
fi

# Certificado mantido em ~/.kube/certs/cluster-ca.crt

echo ""
echo "🎯 Kubectl pronto para uso!"