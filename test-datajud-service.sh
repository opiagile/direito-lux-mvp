#!/bin/bash

echo "🔍 Testando DataJud Service..."
echo "================================"

# Testar health
echo "1️⃣ Health Check:"
curl -s http://localhost:8084/health | jq

echo -e "\n2️⃣ Buscar Processos (requer X-Tenant-ID):"
curl -s -X POST http://localhost:8084/api/v1/search \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" \
  -d '{
    "query": "direito consumidor",
    "tribunais": ["TJRJ", "TJSP"],
    "pagina": 1,
    "tamanho": 10
  }' | jq

echo -e "\n3️⃣ Consultar Processo Específico:"
curl -s http://localhost:8084/api/v1/process/0001234-56.2024.8.19.0001 \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n4️⃣ Movimentações do Processo:"
curl -s http://localhost:8084/api/v1/process/0001234-56.2024.8.19.0001/movements \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n5️⃣ Estatísticas:"
curl -s http://localhost:8084/api/v1/stats \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n6️⃣ Quota de Uso:"
curl -s http://localhost:8084/api/v1/quota \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n7️⃣ CNPJ Providers (Admin):"
curl -s http://localhost:8084/api/v1/cnpj-providers \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n✅ Testes concluídos!"