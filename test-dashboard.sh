#!/bin/bash

echo "🧪 Testando fluxo completo Dashboard - Direito Lux"
echo "================================================="

# 1. Testar Auth Service
echo "1️⃣ Testando Auth Service..."
AUTH_RESPONSE=$(curl -s -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@silvaassociados.com.br","password":"password"}')

if [[ $AUTH_RESPONSE == *"access_token"* ]]; then
  echo "✅ Auth Service: OK"
  TOKEN=$(echo $AUTH_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
  TENANT_ID=$(echo $AUTH_RESPONSE | grep -o '"tenant":{"id":"[^"]*"' | cut -d'"' -f6)
  echo "🔑 Token obtido: ${TOKEN:0:20}..."
  echo "🏢 Tenant ID: $TENANT_ID"
else
  echo "❌ Auth Service: FALHOU"
  echo "Response: $AUTH_RESPONSE"
  exit 1
fi

echo ""

# 2. Testar Tenant Service
echo "2️⃣ Testando Tenant Service..."
TENANT_RESPONSE=$(curl -s "http://localhost:8082/api/v1/tenants/$TENANT_ID" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Authorization: Bearer $TOKEN")

if [[ $TENANT_RESPONSE == *"Silva"* ]]; then
  echo "✅ Tenant Service: OK"
  echo "🏢 Tenant: $(echo $TENANT_RESPONSE | grep -o '"name":"[^"]*"' | cut -d'"' -f4)"
else
  echo "❌ Tenant Service: FALHOU"
  echo "Response: $TENANT_RESPONSE"
fi

echo ""

# 3. Testar Process Service (crítico para dashboard)
echo "3️⃣ Testando Process Service - Endpoint crítico /stats..."
STATS_RESPONSE=$(curl -s "http://localhost:8083/api/v1/processes/stats" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Authorization: Bearer $TOKEN")

if [[ $STATS_RESPONSE == *"total"* ]]; then
  echo "✅ Process Service: OK"
  echo "📊 Estatísticas:"
  echo "   Total: $(echo $STATS_RESPONSE | grep -o '"total":[0-9]*' | cut -d':' -f2)"
  echo "   Ativo: $(echo $STATS_RESPONSE | grep -o '"active":[0-9]*' | cut -d':' -f2)"
  echo "   Pausado: $(echo $STATS_RESPONSE | grep -o '"paused":[0-9]*' | cut -d':' -f2)"
  echo "   Este mês: $(echo $STATS_RESPONSE | grep -o '"this_month":[0-9]*' | cut -d':' -f2)"
else
  echo "❌ Process Service: FALHOU"
  echo "Response: $STATS_RESPONSE"
fi

echo ""

# 4. Testar Frontend
echo "4️⃣ Testando Frontend..."
FRONTEND_RESPONSE=$(curl -s "http://localhost:3000")

if [[ $FRONTEND_RESPONSE == *"Direito Lux"* ]]; then
  echo "✅ Frontend: OK"
  echo "🌐 Aplicação rodando em http://localhost:3000"
else
  echo "❌ Frontend: FALHOU"
fi

echo ""
echo "🎯 RESULTADO FINAL:"
echo "=================="

if [[ $AUTH_RESPONSE == *"access_token"* ]] && [[ $TENANT_RESPONSE == *"Silva"* ]] && [[ $STATS_RESPONSE == *"total"* ]] && [[ $FRONTEND_RESPONSE == *"Direito Lux"* ]]; then
  echo "🎉 SUCESSO! Dashboard deve estar funcionando com dados reais!"
  echo ""
  echo "📋 Para testar manualmente:"
  echo "1. Acesse: http://localhost:3000/login"
  echo "2. Faça login com: admin@silvaassociados.com.br / password"
  echo "3. Vá para: http://localhost:3000/dashboard"
  echo "4. As estatísticas devem aparecer:"
  echo "   • Total de Processos: 45"
  echo "   • Processos Ativos: 38"
  echo "   • Processos Pausados: 5"
  echo "   • Novos este mês: 12"
else
  echo "❌ FALHOU! Algum serviço não está funcionando corretamente."
fi

echo ""
echo "🔧 Serviços rodando:"
echo "   • Auth Service: http://localhost:8081"
echo "   • Tenant Service: http://localhost:8082"
echo "   • Process Service: http://localhost:8083"
echo "   • Frontend: http://localhost:3000"