#!/bin/bash

# =============================================================================
# Script para verificar regiões e zonas válidas no GCP
# =============================================================================

echo "🌍 Verificando regiões disponíveis no GCP..."
echo "==========================================="

echo "📍 Regiões US (recomendadas para staging):"
gcloud compute regions list --filter="name:(us-central1 OR us-east1 OR us-west1)" --format="table(name,description,status)"

echo ""
echo "📍 Zonas em us-central1:"
gcloud compute zones list --filter="region:us-central1" --format="table(name,description,status)" --limit=5

echo ""
echo "🔧 Vou configurar uma região válida..."
gcloud config set compute/region us-central1-a
gcloud config set compute/zone us-central1-a

echo ""
echo "✅ Configuração atualizada!"