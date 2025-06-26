#!/bin/bash

echo "🚀 INICIANDO DESENVOLVIMENTO LOCAL"
echo "=================================="

# 1. Subir apenas infraestrutura essencial
echo "📊 Subindo infraestrutura..."
docker-compose -f docker-compose.minimal.yml up -d

# 2. Aguardar PostgreSQL
echo "⏰ Aguardando PostgreSQL..."
sleep 15

# 3. Verificar se dados existem
echo "🔍 Verificando dados..."
PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c "\dt" | grep -q tenants

if [ $? -ne 0 ]; then
    echo "⚠️  Dados não encontrados. Execute o setup:"
    echo "   ./SETUP_MASTER_ONBOARDING.sh"
else
    echo "✅ Dados encontrados!"
fi

echo ""
echo "🎯 Para desenvolver:"
echo ""
echo "AI Service (Python):"
echo "   cd services/ai-service"
echo "   pip install -r requirements.txt"
echo "   uvicorn main:app --reload --port 8000"
echo ""
echo "Auth Service (Go):"
echo "   cd services/auth-service"
echo "   go run cmd/server/main.go"
echo ""
echo "Frontend (Next.js):"
echo "   cd frontend"
echo "   npm install && npm run dev"
echo ""
echo "🌐 URLs:"
echo "   - PostgreSQL: localhost:5432"
echo "   - Redis: localhost:6379"
echo "   - MailPit: http://localhost:8025"
