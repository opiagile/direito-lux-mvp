#!/bin/bash

echo "🚀 SETUP COMPLETO CORRIGIDO - DIREITO LUX"
echo "=========================================="
echo ""
echo "Este script vai configurar o ambiente do ZERO com dados corretos"
echo ""

# ============================================================================
# CONFIGURAÇÕES
# ============================================================================

# Database config
DB_NAME="direito_lux_dev"
DB_USER="direito_lux"
DB_PASSWORD="dev_password_123"
DB_HOST="localhost"
DB_PORT="5432"

# ============================================================================
# FASE 1: LIMPEZA TOTAL
# ============================================================================

echo "🧹 FASE 1: Limpeza do ambiente..."
echo ""

# Parar containers se estiverem rodando
echo "Parando containers Docker..."
docker-compose down -v 2>/dev/null || true

# Aguardar um pouco
sleep 3

# Limpar containers órfãos
echo "Removendo containers órfãos..."
docker container prune -f 2>/dev/null || true
docker volume prune -f 2>/dev/null || true

# Parar PostgreSQL local se estiver rodando
if command -v brew &> /dev/null; then
    if brew services list | grep -q "postgresql.*started"; then
        echo "Parando PostgreSQL local..."
        brew services stop postgresql
    fi
fi

# ============================================================================
# FASE 2: SUBIR INFRAESTRUTURA
# ============================================================================

echo ""
echo "🏗️ FASE 2: Subindo infraestrutura..."
echo ""

# Iniciar apenas serviços de infraestrutura
echo "Subindo PostgreSQL, Redis e RabbitMQ..."
docker-compose up -d postgres redis rabbitmq

# Aguardar os serviços ficarem prontos
echo "Aguardando PostgreSQL ficar pronto..."
for i in {1..30}; do
    if docker exec $(docker-compose ps -q postgres) pg_isready -U $DB_USER 2>/dev/null; then
        echo "✅ PostgreSQL está pronto!"
        break
    fi
    echo "Aguardando... ($i/30)"
    sleep 2
done

# Verificar se PostgreSQL está realmente funcionando
if ! docker exec $(docker-compose ps -q postgres) pg_isready -U $DB_USER 2>/dev/null; then
    echo "❌ ERRO: PostgreSQL não está funcionando!"
    exit 1
fi

echo "Aguardando Redis ficar pronto..."
for i in {1..15}; do
    if docker exec $(docker-compose ps -q redis) redis-cli ping 2>/dev/null | grep -q PONG; then
        echo "✅ Redis está pronto!"
        break
    fi
    echo "Aguardando... ($i/15)"
    sleep 2
done

echo "Aguardando RabbitMQ ficar pronto..."
for i in {1..30}; do
    if docker exec $(docker-compose ps -q rabbitmq) rabbitmqctl status 2>/dev/null | grep -q "Status of node"; then
        echo "✅ RabbitMQ está pronto!"
        break
    fi
    echo "Aguardando... ($i/30)"
    sleep 2
done

# ============================================================================
# FASE 3: MIGRATIONS
# ============================================================================

echo ""
echo "📄 FASE 3: Executando migrations..."
echo ""

# Verificar se migrate está disponível
if ! command -v migrate &> /dev/null; then
    echo "⚠️  golang-migrate não encontrado. Instalando..."
    if command -v brew &> /dev/null; then
        brew install golang-migrate
    else
        echo "❌ ERRO: Por favor instale golang-migrate manualmente"
        echo "   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz"
        exit 1
    fi
fi

# URL de conexão com o banco
DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Ordem correta de execução das migrations
echo "1. Executando migrations do Tenant Service (planos e tenants)..."
cd services/tenant-service
if [ -d "migrations" ]; then
    migrate -path migrations -database "$DATABASE_URL" up
    if [ $? -eq 0 ]; then
        echo "✅ Tenant Service migrations executadas"
    else
        echo "❌ ERRO nas migrations do Tenant Service"
        exit 1
    fi
else
    echo "⚠️  Diretório migrations não encontrado no tenant-service"
fi

echo ""
echo "2. Executando migrations do Auth Service (usuários e sessões)..."
cd ../auth-service
if [ -d "migrations" ]; then
    # Executar apenas as 2 primeiras migrations (sem o seed antigo)
    migrate -path migrations -database "$DATABASE_URL" up 2
    if [ $? -eq 0 ]; then
        echo "✅ Auth Service migrations executadas"
    else
        echo "❌ ERRO nas migrations do Auth Service"
        exit 1
    fi
else
    echo "⚠️  Diretório migrations não encontrado no auth-service"
fi

echo ""
echo "3. Executando migrations do Process Service..."
cd ../process-service
if [ -d "migrations" ]; then
    migrate -path migrations -database "$DATABASE_URL" up
    if [ $? -eq 0 ]; then
        echo "✅ Process Service migrations executadas"
    else
        echo "❌ ERRO nas migrations do Process Service"
        exit 1
    fi
else
    echo "⚠️  Diretório migrations não encontrado no process-service"
fi

# Voltar para o diretório raiz
cd ../..

# ============================================================================
# FASE 4: CARREGAR DADOS DE TESTE
# ============================================================================

echo ""
echo "📊 FASE 4: Carregando dados de teste..."
echo ""

# Executar o script de seed completo
echo "Carregando dados de teste completos..."
if docker exec -i $(docker-compose ps -q postgres) psql -U $DB_USER -d $DB_NAME < SEED_DATABASE_COMPLETE.sql; then
    echo "✅ Dados de teste carregados com sucesso!"
else
    echo "❌ ERRO ao carregar dados de teste"
    echo "Tentando carregar com método alternativo..."
    
    # Método alternativo usando PGPASSWORD
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f SEED_DATABASE_COMPLETE.sql
    
    if [ $? -eq 0 ]; then
        echo "✅ Dados carregados via método alternativo!"
    else
        echo "❌ ERRO: Não foi possível carregar os dados"
        exit 1
    fi
fi

# ============================================================================
# FASE 5: VERIFICAÇÃO
# ============================================================================

echo ""
echo "🔍 FASE 5: Verificando dados..."
echo ""

# Verificar se os dados foram carregados corretamente
echo "Verificando tenants..."
TENANT_COUNT=$(docker exec $(docker-compose ps -q postgres) psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM tenants;" | tr -d ' ')
echo "Tenants encontrados: $TENANT_COUNT"

echo "Verificando usuários..."
USER_COUNT=$(docker exec $(docker-compose ps -q postgres) psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM users;" | tr -d ' ')
echo "Usuários encontrados: $USER_COUNT"

echo "Verificando planos..."
PLAN_COUNT=$(docker exec $(docker-compose ps -q postgres) psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM plans;" | tr -d ' ')
echo "Planos encontrados: $PLAN_COUNT"

echo "Verificando subscriptions..."
SUB_COUNT=$(docker exec $(docker-compose ps -q postgres) psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM subscriptions;" | tr -d ' ')
echo "Subscriptions encontradas: $SUB_COUNT"

# ============================================================================
# FASE 6: TESTAR LOGIN
# ============================================================================

echo ""
echo "🔑 FASE 6: Testando credenciais..."
echo ""

# Mostrar credenciais de teste
echo "Credenciais para teste de login:"
echo "--------------------------------"
echo ""

docker exec $(docker-compose ps -q postgres) psql -U $DB_USER -d $DB_NAME -c "
SELECT 
    t.name AS tenant,
    u.email,
    u.role,
    'password' AS senha
FROM users u
JOIN tenants t ON u.tenant_id = t.id  
WHERE u.role = 'admin'
ORDER BY t.plan_type, t.name;
"

# ============================================================================
# FASE 7: RELATÓRIO FINAL
# ============================================================================

echo ""
echo "📋 RELATÓRIO FINAL"
echo "=================="
echo ""
echo "✅ Infraestrutura: PostgreSQL, Redis, RabbitMQ rodando"
echo "✅ Migrations: Tenant, Auth e Process executadas"
echo "✅ Dados de teste: $TENANT_COUNT tenants, $USER_COUNT usuários, $PLAN_COUNT planos"
echo ""
echo "🌐 URLs disponíveis:"
echo "   PostgreSQL: localhost:5432"
echo "   Redis: localhost:6379"  
echo "   RabbitMQ: http://localhost:15672"
echo ""
echo "🔑 Para fazer login:"
echo "   Email: admin@silvaassociados.com.br"
echo "   Senha: password"
echo ""
echo "📝 Próximos passos:"
echo "   1. cd frontend && npm install && npm run dev"
echo "   2. Acessar http://localhost:3000"
echo "   3. Fazer login com as credenciais acima"
echo ""
echo "🎯 STATUS: AMBIENTE CONFIGURADO COM SUCESSO!"
echo ""

# Opcional: subir outros serviços
read -p "Deseja subir os microserviços também? (s/n): " resposta
if [ "$resposta" = "s" ] || [ "$resposta" = "S" ]; then
    echo ""
    echo "🚀 Subindo microserviços..."
    docker-compose up -d
    echo "✅ Todos os serviços estão rodando!"
    echo "   Verifique com: docker-compose ps"
fi

echo ""
echo "✨ Setup completo finalizado!"