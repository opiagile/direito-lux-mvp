#!/bin/bash

# =============================================================================
# Script MASTER - Deploy Completo Staging - Direito Lux
# =============================================================================

set -e

echo "🎯 DEPLOY COMPLETO STAGING - DIREITO LUX"
echo "======================================"
echo ""

echo "📋 PRÉ-REQUISITOS:"
echo "   ✅ Billing vinculado ao projeto direito-lux-staging"
echo "   ✅ gcloud auth application-default login executado"
echo "   ✅ GitHub Secrets configurados"
echo ""

read -p "🚨 Todos os pré-requisitos estão OK? (y/N): " confirm
if [[ $confirm != [yY] ]]; then
    echo "❌ Complete os pré-requisitos primeiro!"
    exit 1
fi

echo ""
echo "🚀 INICIANDO DEPLOY COMPLETO..."
echo "============================="

echo ""
echo "FASE 1: 🏗️  TERRAFORM (Infraestrutura)"
echo "--------------------------------------"
/Users/franc/Opiagile/SAAS/direito-lux/scripts/terraform-deploy-staging.sh

echo ""
read -p "🎯 Terraform plan OK? Executar apply? (y/N): " terraform_confirm
if [[ $terraform_confirm == [yY] ]]; then
    cd /Users/franc/Opiagile/SAAS/direito-lux/terraform
    terraform apply -var-file=environments/staging.tfvars -auto-approve
    echo "✅ Infraestrutura criada!"
else
    echo "❌ Deploy cancelado pelo usuário"
    exit 1
fi

echo ""
echo "FASE 2: ☸️  KUBERNETES (Aplicações)"
echo "----------------------------------"
sleep 30  # Aguardar cluster ficar pronto
/Users/franc/Opiagile/SAAS/direito-lux/scripts/k8s-deploy-staging.sh

echo ""
echo "FASE 3: 🌍 DNS CONFIGURATION"
echo "----------------------------"
LB_IP=$(kubectl get ingress direito-lux-ingress -n direito-lux -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

echo "🎯 CONFIGURE DNS:"
echo "   Domain: staging.direitolux.com.br"
echo "   Type: A Record"
echo "   Value: $LB_IP"
echo "   TTL: 300"

echo ""
echo "🎉 DEPLOY STAGING COMPLETO!"
echo "=========================="
echo "📋 URLs disponíveis:"
echo "   🌐 Frontend: https://staging.direitolux.com.br"
echo "   🔧 API: https://staging.direitolux.com.br/api"
echo "   📊 Monitoring: https://staging.direitolux.com.br/grafana"
echo ""
echo "🎯 PRÓXIMOS PASSOS:"
echo "   1. Configure DNS staging.direitolux.com.br → $LB_IP"
echo "   2. Aguarde SSL provisioning (~5-10 min)"
echo "   3. Teste aplicação completa"
echo "   4. Deploy produção quando validado"

echo ""
echo "✅ DIREITO LUX STAGING LIVE! 🚀"