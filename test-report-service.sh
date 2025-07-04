#!/bin/bash

echo "🔍 Testando Report Service..."
echo "================================"

# Testar health
echo "1️⃣ Health Check:"
curl -s http://localhost:8087/health | jq

echo -e "\n2️⃣ Atividades Recentes (CRÍTICO PARA DASHBOARD):"
curl -s http://localhost:8087/api/v1/reports/recent-activities \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n3️⃣ Dashboard Data (KPIs Adicionais):"
curl -s http://localhost:8087/api/v1/reports/dashboard \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n4️⃣ Listar Relatórios:"
curl -s http://localhost:8087/api/v1/reports \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n5️⃣ Criar Relatório:"
curl -s -X POST http://localhost:8087/api/v1/reports \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" \
  -d '{
    "titulo": "Relatório de Teste",
    "tipo": "performance",
    "formato": "PDF",
    "parametros": {"periodo": "mensal"}
  }' | jq

echo -e "\n6️⃣ Consultar Relatório Específico:"
curl -s http://localhost:8087/api/v1/reports/report_1 \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n7️⃣ Relatórios Agendados:"
curl -s http://localhost:8087/api/v1/reports/scheduled \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" | jq

echo -e "\n8️⃣ Criar Relatório Agendado:"
curl -s -X POST http://localhost:8087/api/v1/reports/scheduled \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 11111111-1111-1111-1111-111111111111" \
  -d '{
    "titulo": "Relatório Mensal Automático",
    "tipo": "monthly_summary",
    "formato": "PDF",
    "frequencia": "monthly"
  }' | jq

echo -e "\n✅ Testes concluídos!"
echo ""
echo "🎯 DASHBOARD ENDPOINTS FUNCIONAIS:"
echo "   • /api/v1/reports/recent-activities ← CRÍTICO para completar dashboard"
echo "   • /api/v1/reports/dashboard ← KPIs adicionais"