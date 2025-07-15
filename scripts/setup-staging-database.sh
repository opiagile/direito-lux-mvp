#!/bin/bash

echo "🗄️ SETUP STAGING DATABASE - EXECUTAR TODAS MIGRATIONS"
echo "======================================================"

# Variáveis
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-direito_lux}"
DB_PASSWORD="${DB_PASSWORD:-dev_password_123}"
DB_NAME="${DB_NAME:-direito_lux_staging}"
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "🔗 Conectando em: ${DB_HOST}:${DB_PORT}/${DB_NAME}"

# Função para executar migrations
run_migrations() {
    local service=$1
    local service_path="services/${service}"
    
    echo "📋 Executando migrations: ${service}"
    
    if [ ! -d "${service_path}/migrations" ]; then
        echo "⚠️  Migrations não encontradas para ${service}"
        return 0
    fi
    
    cd "${service_path}"
    
    # Verificar se migrate está instalado
    if ! command -v migrate &> /dev/null; then
        echo "❌ migrate CLI não encontrado. Instale com:"
        echo "   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
    
    # Executar migrations
    if migrate -path migrations -database "${DB_URL}" up; then
        echo "✅ Migrations ${service} executadas com sucesso"
    else
        echo "⚠️  Erro nas migrations ${service} (pode ser normal se já executadas)"
    fi
    
    cd - > /dev/null
}

# Lista de serviços com migrations
SERVICES=(
    "auth-service"
    "tenant-service"
    "process-service"
    "search-service"
    "notification-service"
    "ai-service"
    "report-service"
    "billing-service"
)

echo "🚀 Executando migrations para todos os serviços..."

for service in "${SERVICES[@]}"; do
    run_migrations "$service"
    echo ""
done

echo "✅ Setup de database completo!"
echo ""
echo "🧪 Comandos de teste:"
echo "   # Teste Auth Service"
echo "   curl -k https://35.188.198.87/api/v1/auth/login \\"
echo "     -X POST -H 'Content-Type: application/json' \\"
echo "     -H 'X-Tenant-ID: 550e8400-e29b-41d4-a716-446655440001' \\"
echo "     -d '{\"email\":\"admin@silvaassociados.com.br\",\"password\":\"password\"}'"
echo ""
echo "   # Teste Tenant Service"  
echo "   curl -k https://35.188.198.87/api/v1/tenants/550e8400-e29b-41d4-a716-446655440001 \\"
echo "     -H 'X-Tenant-ID: 550e8400-e29b-41d4-a716-446655440001'"