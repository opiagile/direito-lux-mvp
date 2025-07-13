#!/bin/bash

# =============================================================================
# Script para verificar ambiente GCP
# =============================================================================

echo "🔍 Verificando ambiente GCP atual..."
echo "=================================="

echo "📋 1. Projetos disponíveis:"
gcloud projects list --format="table(projectId,name,lifecycleState)"

echo ""
echo "💳 2. Billing accounts:"
gcloud beta billing accounts list --format="table(name,displayName,open)"

echo ""
echo "🌍 3. Região atual:"
gcloud config get-value compute/region

echo ""
echo "🎯 4. Projeto ativo:"
gcloud config get-value project

echo ""
echo "✅ Verificação completa!"