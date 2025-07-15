#!/bin/bash

echo "🔍 PROCURANDO CLUSTER GKE NO PROJETO"
echo "===================================="

echo "1. 📋 Listando todos os clusters GKE no projeto..."
gcloud container clusters list --project=direito-lux-staging-2025

echo ""
echo "2. 🌍 Listando clusters em todas as regiões/zonas..."
gcloud container clusters list --project=direito-lux-staging-2025 --format="table(name,location,status)"

echo ""
echo "3. 🔍 Buscando especificamente por 'direito-lux'..."
gcloud container clusters list --project=direito-lux-staging-2025 --filter="name~direito" --format="table(name,location,status,currentMasterVersion)"