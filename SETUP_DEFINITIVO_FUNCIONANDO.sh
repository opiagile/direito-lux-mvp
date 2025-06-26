#!/bin/bash

echo "🚀 SETUP DEFINITIVO FUNCIONANDO - DIREITO LUX"
echo "============================================="

# 1. Limpar completamente
echo "🧹 Limpando ambiente..."
docker-compose down -v
docker volume prune -f

# 2. Remover script antigo problemático
rm -f infrastructure/sql/init/00-create-user.sh 2>/dev/null

# 3. Subir PostgreSQL com usuário postgres padrão
echo "🐘 Iniciando PostgreSQL..."
docker-compose up -d postgres

# 4. Aguardar PostgreSQL estar completamente pronto
echo "⏰ Aguardando PostgreSQL (30s)..."
sleep 30

# 5. Verificar se está rodando
echo "🔍 Verificando PostgreSQL..."
docker exec direito-lux-postgres pg_isready -U postgres
if [ $? -ne 0 ]; then
    echo "⏰ Aguardando mais 15s..."
    sleep 15
fi

# 6. Testar conexão com usuário direito_lux
echo "🔗 Testando conexão com direito_lux..."
PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -c "SELECT version();"

if [ $? -eq 0 ]; then
    echo "✅ Usuário direito_lux já existe e está funcionando!"
else
    echo "❌ Usuário não existe. Criando manualmente..."
    
    # Criar usuário e banco via superusuário postgres
    PGPASSWORD=postgres psql -h localhost -U postgres << EOF
-- Criar role direito_lux
CREATE ROLE direito_lux WITH LOGIN PASSWORD 'dev_password_123' CREATEDB SUPERUSER;

-- Criar database
CREATE DATABASE direito_lux_dev OWNER direito_lux;

-- Conectar ao novo banco
\c direito_lux_dev

-- Criar extensões
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Verificar
SELECT 'Setup manual concluído!' as status;
EOF

    # Testar novamente
    echo "🔗 Testando conexão após setup manual..."
    PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -c "SELECT 'Conexão OK!' as status;"
    
    if [ $? -ne 0 ]; then
        echo "❌ Erro persistente. Verifique os logs:"
        docker logs direito-lux-postgres | tail -20
        exit 1
    fi
fi

# 7. Executar migrations
echo ""
echo "📋 Executando migrations..."

# Verificar golang-migrate
if ! command -v migrate &> /dev/null; then
    echo "📦 Instalando golang-migrate..."
    brew install golang-migrate
fi

# Tenant Service
echo "  📊 [1/3] Tenant Service..."
cd services/tenant-service
DATABASE_URL="postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" make migrate-up
if [ $? -eq 0 ]; then
    echo "  ✅ Tenant Service OK"
else
    echo "  ❌ Erro no Tenant Service"
    exit 1
fi

# Auth Service
echo "  🔐 [2/3] Auth Service..."
cd ../auth-service
DATABASE_URL="postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" make migrate-up
if [ $? -eq 0 ]; then
    echo "  ✅ Auth Service OK"
else
    echo "  ❌ Erro no Auth Service"
    exit 1
fi

# Process Service
echo "  ⚖️ [3/3] Process Service..."
cd ../process-service
DATABASE_URL="postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" make migrate-up
if [ $? -eq 0 ]; then
    echo "  ✅ Process Service OK"
else
    echo "  ❌ Erro no Process Service"
    exit 1
fi

cd ../..

# 8. Verificar dados
echo ""
echo "📊 Dados carregados:"
PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -t << EOF
SELECT 'Tenants: ' || COUNT(*) FROM tenants
UNION ALL
SELECT 'Users: ' || COUNT(*) FROM users  
UNION ALL
SELECT 'Processes: ' || COUNT(*) FROM processes;
EOF

# 9. Mostrar credenciais
echo ""
echo "🔑 Credenciais de teste (senha: 123456):"
PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -t << EOF
SELECT email || ' - ' || role || ' (' || t.plan || ')'
FROM users u
JOIN tenants t ON u.tenant_id = t.id
WHERE u.role = 'admin'
ORDER BY t.plan
LIMIT 4;
EOF

echo ""
echo "🎊 SETUP COMPLETO!"
echo "=================="
echo ""
echo "✅ PostgreSQL: Rodando com usuário postgres como superuser"
echo "✅ Aplicação: Usando usuário direito_lux"
echo "✅ Database: direito_lux_dev"
echo "✅ Migrations: Executadas com sucesso"
echo "✅ Dados de teste: Carregados"
echo ""
echo "🚀 Próximos comandos:"
echo "   docker-compose up -d     # Subir todos os serviços"
echo "   cd frontend && npm run dev   # Iniciar frontend"
echo ""
echo "📝 Documentação: cat DOCUMENTO_TESTE_VALIDACAO.md"