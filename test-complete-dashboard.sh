#!/bin/bash

echo "🎯 Teste Completo do Dashboard - Todos os KPIs"
echo "=============================================="

# Fazer login para obter token
AUTH_RESPONSE=$(curl -s -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@silvaassociados.com.br","password":"password"}')

TOKEN=$(echo $AUTH_RESPONSE | jq -r '.access_token')
TENANT_ID="11111111-1111-1111-1111-111111111111"

echo "🔑 Login realizado para Silva & Associados"
echo ""

# Testar stats completos
STATS_RESPONSE=$(curl -s "http://127.0.0.1:8083/api/v1/processes/stats" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "📊 DASHBOARD - Todos os KPIs:"
echo "================================"

if command -v jq &> /dev/null; then
  echo "📋 Total de Processos: $(echo $STATS_RESPONSE | jq -r '.data.total')"
  echo "🟢 Processos Ativos: $(echo $STATS_RESPONSE | jq -r '.data.active')"
  echo "📈 Movimentações Hoje: $(echo $STATS_RESPONSE | jq -r '.data.todayMovements')"
  echo "⚠️  Prazos Próximos: $(echo $STATS_RESPONSE | jq -r '.data.upcomingDeadlines')"
  echo ""
  echo "📊 Dados adicionais:"
  echo "⏸️  Processos Pausados: $(echo $STATS_RESPONSE | jq -r '.data.paused')"
  echo "📦 Processos Arquivados: $(echo $STATS_RESPONSE | jq -r '.data.archived')"
  echo "📅 Novos este Mês: $(echo $STATS_RESPONSE | jq -r '.data.this_month')"
else
  echo "Raw response: $STATS_RESPONSE"
fi

echo ""
echo "✅ Todos os 4 cards do dashboard devem estar preenchidos agora!"
echo ""
echo "🌐 Para ver no navegador:"
echo "   1. Acesse: http://localhost:3000/dashboard"
echo "   2. Faça login com: admin@silvaassociados.com.br / password"
echo "   3. Observe os 4 cards principais:"
echo "      • Total de Processos: 45"
echo "      • Processos Ativos: 38" 
echo "      • Movimentações Hoje: 3"
echo "      • Prazos Próximos: 7"