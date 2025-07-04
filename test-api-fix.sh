#!/bin/bash

echo "🔧 Testando correção da API - Process Service"
echo "============================================="

# 1. Testar auth para obter token
echo "1️⃣ Fazendo login para obter token..."
AUTH_RESPONSE=$(curl -s -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@silvaassociados.com.br","password":"password"}')

TOKEN=$(echo $AUTH_RESPONSE | jq -r '.access_token')
TENANT_ID="11111111-1111-1111-1111-111111111111"

echo "✅ Token obtido: ${TOKEN:0:20}..."
echo ""

# 2. Testar diretamente o Process Service
echo "2️⃣ Testando Process Service diretamente..."
DIRECT_RESPONSE=$(curl -s "http://127.0.0.1:8083/api/v1/processes/stats" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "📊 Resposta direta: $DIRECT_RESPONSE"
echo ""

# 3. Testar via frontend API
echo "3️⃣ Testando como o frontend faria a chamada..."
FRONTEND_RESPONSE=$(curl -s "http://127.0.0.1:8083/api/v1/processes/stats" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Origin: http://localhost:3000" \
  -H "Content-Type: application/json")

echo "🌐 Resposta simulando frontend: $FRONTEND_RESPONSE"
echo ""

# 4. Verificar se contém dados esperados
if [[ $FRONTEND_RESPONSE == *"total"* ]] && [[ $FRONTEND_RESPONSE == *"45"* ]]; then
  echo "🎉 SUCESSO! API corrigida e funcionando!"
  echo ""
  echo "📋 Dados retornados:"
  echo "   • Total: $(echo $FRONTEND_RESPONSE | jq -r '.data.total')"
  echo "   • Ativos: $(echo $FRONTEND_RESPONSE | jq -r '.data.active')"
  echo "   • Pausados: $(echo $FRONTEND_RESPONSE | jq -r '.data.paused')"
  echo "   • Este mês: $(echo $FRONTEND_RESPONSE | jq -r '.data.this_month')"
  echo ""
  echo "✅ O dashboard agora deve funcionar corretamente!"
  echo "   Acesse: http://localhost:3000/dashboard"
else
  echo "❌ Ainda há problemas com a API"
  echo "Response: $FRONTEND_RESPONSE"
fi