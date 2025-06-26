#!/bin/bash

# =============================================================================
# DIREITO LUX - START CLEAN ENVIRONMENT
# =============================================================================
# Script completo para limpar cache, derrubar serviços e subir ambiente limpo
# =============================================================================

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Diretório base
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$BASE_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🚀 DIREITO LUX - CLEAN START${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# =============================================================================
# 1. PARAR TODOS OS SERVIÇOS
# =============================================================================
echo -e "${YELLOW}📋 Parando todos os serviços...${NC}"

# Parar Docker Compose (ignora erro se não estiver rodando)
docker-compose down -v --remove-orphans 2>/dev/null || true
docker-compose -f docker-compose.simple.yml down -v --remove-orphans 2>/dev/null || true
docker-compose -f services/docker-compose.dev.yml down -v --remove-orphans 2>/dev/null || true

# Parar containers órfãos
echo -e "${YELLOW}🔍 Procurando containers órfãos...${NC}"
ORPHAN_CONTAINERS=$(docker ps -a -q --filter "label=com.docker.compose.project=direito-lux" 2>/dev/null || true)
if [ ! -z "$ORPHAN_CONTAINERS" ]; then
    echo -e "${YELLOW}Removendo containers órfãos...${NC}"
    docker rm -f $ORPHAN_CONTAINERS 2>/dev/null || true
fi

# =============================================================================
# 2. LIMPAR CACHE E VOLUMES
# =============================================================================
echo -e "${YELLOW}🧹 Limpando cache e volumes...${NC}"

# Remover redes antigas do projeto
echo -e "${YELLOW}Removendo redes antigas...${NC}"
docker network ls -q --filter name=direito | xargs -r docker network rm 2>/dev/null || true

# Remover volumes específicos do projeto
echo -e "${YELLOW}Removendo volumes antigos...${NC}"
docker volume ls -q | grep -E "direito-lux|direito_lux" | xargs -r docker volume rm 2>/dev/null || true

# Limpar cache de build
echo -e "${YELLOW}Limpando cache de build...${NC}"
docker builder prune -f 2>/dev/null || true

# Limpar imagens não utilizadas
echo -e "${YELLOW}Limpando imagens não utilizadas...${NC}"
docker image prune -f 2>/dev/null || true

# Limpar sistema (mantém imagens base para economizar tempo)
echo -e "${YELLOW}Limpeza final do sistema...${NC}"
docker system prune -f --volumes 2>/dev/null || true

# Limpar logs locais
echo -e "${YELLOW}📄 Limpando logs...${NC}"
rm -rf logs/*.log 2>/dev/null || true
find services -name "*.log" -type f -delete 2>/dev/null || true

# Limpar binários compilados
echo -e "${YELLOW}🗑️ Limpando binários...${NC}"
find services -name "main" -type f -delete 2>/dev/null || true
find services -name "server" -type f -delete 2>/dev/null || true
find services -name "*-service" -type f -delete 2>/dev/null || true

# =============================================================================
# 3. VERIFICAR PORTAS
# =============================================================================
echo -e "${YELLOW}🔍 Verificando portas...${NC}"

check_port() {
    local port=$1
    local service=$2
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${RED}⚠️  Porta $port ($service) está em uso!${NC}"
        echo -e "${YELLOW}   Matando processo na porta $port...${NC}"
        lsof -ti:$port | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# Verificar principais portas
check_port 5432 "PostgreSQL"
check_port 6379 "Redis"
check_port 15672 "RabbitMQ Management"
check_port 5672 "RabbitMQ"
check_port 9200 "Elasticsearch"
check_port 3000 "Frontend"
check_port 8000 "AI Service"
check_port 8080 "API Gateway"
check_port 8081 "Auth Service"
check_port 8082 "Tenant Service"
check_port 8083 "Process Service"
check_port 8084 "DataJud Service"
check_port 8085 "Notification Service"
check_port 8086 "Search Service"
check_port 8087 "Report Service"
check_port 8088 "MCP Service"

echo -e "${GREEN}✅ Portas liberadas${NC}"

# =============================================================================
# 4. CRIAR REDE DOCKER
# =============================================================================
echo -e "${YELLOW}🌐 Configurando rede Docker...${NC}"

# Remover rede antiga se existir
docker network rm direito-lux-network 2>/dev/null || true

# Aguardar um pouco
sleep 1

# Criar nova rede
if docker network create direito-lux-network 2>/dev/null; then
    echo -e "${GREEN}✅ Rede direito-lux-network criada${NC}"
else
    echo -e "${YELLOW}⚠️  Usando rede existente${NC}"
fi

# =============================================================================
# 5. SUBIR INFRAESTRUTURA BASE
# =============================================================================
echo -e "${BLUE}🏗️ Subindo infraestrutura base...${NC}"

# Usar docker-compose simplificado sem builds
COMPOSE_FILE="docker-compose.simple.yml"

# Subir apenas serviços essenciais primeiro
docker-compose -f $COMPOSE_FILE up -d postgres redis rabbitmq

# Aguardar PostgreSQL estar pronto
echo -e "${YELLOW}⏳ Aguardando PostgreSQL...${NC}"
until docker-compose -f $COMPOSE_FILE exec -T postgres pg_isready -U postgres 2>/dev/null; do
    echo -n "."
    sleep 2
done
echo -e "\n${GREEN}✅ PostgreSQL pronto${NC}"

# Aguardar Redis (usando healthcheck)
echo -e "${YELLOW}⏳ Aguardando Redis...${NC}"
REDIS_TIMEOUT=60
REDIS_COUNT=0

# Primeiro tenta usar o healthcheck do Docker
while [ "$(docker inspect --format='{{.State.Health.Status}}' direito-lux-redis 2>/dev/null)" != "healthy" ]; do
    echo -n "."
    sleep 2
    REDIS_COUNT=$((REDIS_COUNT + 2))
    
    if [ $REDIS_COUNT -ge $REDIS_TIMEOUT ]; then
        echo -e "\n${YELLOW}Healthcheck timeout, tentando conexão direta...${NC}"
        
        # Tenta conexão direta como fallback
        if docker-compose exec -T redis redis-cli --no-auth-warning -a dev_redis_123 ping 2>/dev/null | grep -q PONG; then
            echo -e "${GREEN}✅ Redis conectado (conexão direta)${NC}"
            break
        else
            echo -e "\n${RED}❌ Redis não respondeu após ${REDIS_TIMEOUT}s${NC}"
            echo -e "${YELLOW}Logs do Redis:${NC}"
            docker-compose logs --tail=5 redis
            echo -e "${YELLOW}Status do container:${NC}"
            docker inspect --format='{{.State.Status}} - {{.State.Health.Status}}' direito-lux-redis
        fi
        break
    fi
done

if [ $REDIS_COUNT -lt $REDIS_TIMEOUT ]; then
    echo -e "\n${GREEN}✅ Redis pronto${NC}"
fi

# Aguardar RabbitMQ (usando healthcheck)
echo -e "${YELLOW}⏳ Aguardando RabbitMQ...${NC}"
RABBITMQ_TIMEOUT=120  # RabbitMQ demora mais para iniciar
RABBITMQ_COUNT=0

# Primeiro tenta usar o healthcheck do Docker
while [ "$(docker inspect --format='{{.State.Health.Status}}' direito-lux-rabbitmq 2>/dev/null)" != "healthy" ]; do
    echo -n "."
    sleep 3
    RABBITMQ_COUNT=$((RABBITMQ_COUNT + 3))
    
    if [ $RABBITMQ_COUNT -ge $RABBITMQ_TIMEOUT ]; then
        echo -e "\n${YELLOW}Healthcheck timeout, tentando verificação direta...${NC}"
        
        # Tenta verificação direta como fallback
        if docker-compose exec -T rabbitmq rabbitmqctl status 2>/dev/null | grep -q "Running"; then
            echo -e "${GREEN}✅ RabbitMQ conectado (verificação direta)${NC}"
            break
        elif docker-compose exec -T rabbitmq rabbitmq-diagnostics ping 2>/dev/null | grep -q "Ping succeeded"; then
            echo -e "${GREEN}✅ RabbitMQ conectado (ping)${NC}"
            break
        else
            echo -e "\n${RED}❌ RabbitMQ não respondeu após ${RABBITMQ_TIMEOUT}s${NC}"
            echo -e "${YELLOW}Logs do RabbitMQ:${NC}"
            docker-compose logs --tail=10 rabbitmq
            echo -e "${YELLOW}Status do container:${NC}"
            docker inspect --format='{{.State.Status}} - {{.State.Health.Status}}' direito-lux-rabbitmq
            echo -e "${YELLOW}Tentando continuar mesmo assim...${NC}"
        fi
        break
    fi
done

if [ $RABBITMQ_COUNT -lt $RABBITMQ_TIMEOUT ]; then
    echo -e "\n${GREEN}✅ RabbitMQ pronto${NC}"
fi

# =============================================================================
# 6. EXECUTAR SETUP MASTER
# =============================================================================
echo -e "${BLUE}🔧 Executando setup master...${NC}"

if [ -f "./SETUP_MASTER_ONBOARDING.sh" ]; then
    chmod +x ./SETUP_MASTER_ONBOARDING.sh
    ./SETUP_MASTER_ONBOARDING.sh
else
    echo -e "${RED}❌ SETUP_MASTER_ONBOARDING.sh não encontrado!${NC}"
    echo -e "${YELLOW}Executando setup básico...${NC}"
    
    # Setup básico se o master não existir
    docker-compose exec -T postgres psql -U postgres <<EOF
CREATE USER direito_lux WITH PASSWORD 'direito_lux_pass_dev';
CREATE DATABASE direito_lux_dev OWNER direito_lux;
GRANT ALL PRIVILEGES ON DATABASE direito_lux_dev TO direito_lux;
EOF
fi

# =============================================================================
# 7. SUBIR SERVIÇOS DE APOIO
# =============================================================================
echo -e "${BLUE}🚀 Subindo serviços de apoio...${NC}"

# Subir resto da infraestrutura (pgAdmin, Mailhog)
docker-compose -f $COMPOSE_FILE up -d

# Aguardar serviços de apoio
echo -e "${YELLOW}⏳ Aguardando serviços de apoio...${NC}"
sleep 5

# =============================================================================
# 8. VERIFICAR SAÚDE DA INFRAESTRUTURA
# =============================================================================
echo -e "${BLUE}🏥 Verificando saúde da infraestrutura...${NC}"

# Verificar PostgreSQL
echo -n "PostgreSQL: "
if docker-compose -f $COMPOSE_FILE exec -T postgres pg_isready -U postgres 2>/dev/null | grep -q "accepting connections"; then
    echo -e "${GREEN}✅ OK${NC}"
else
    echo -e "${RED}❌ FALHOU${NC}"
fi

# Verificar Redis
echo -n "Redis: "
if docker-compose -f $COMPOSE_FILE exec -T redis redis-cli --no-auth-warning -a dev_redis_123 ping 2>/dev/null | grep -q PONG; then
    echo -e "${GREEN}✅ OK${NC}"
else
    echo -e "${RED}❌ FALHOU${NC}"
fi

# Verificar RabbitMQ
echo -n "RabbitMQ: "
if docker-compose -f $COMPOSE_FILE exec -T rabbitmq rabbitmq-diagnostics ping 2>/dev/null | grep -q "Ping succeeded"; then
    echo -e "${GREEN}✅ OK${NC}"
else
    echo -e "${RED}❌ FALHOU${NC}"
fi

# Verificar pgAdmin
echo -n "pgAdmin: "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:5050" | grep -q "200\|302"; then
    echo -e "${GREEN}✅ OK${NC}"
else
    echo -e "${YELLOW}⚠️  Pode demorar um pouco${NC}"
fi

# Verificar Mailhog
echo -n "Mailhog: "
if curl -s -o /dev/null -w "%{http_code}" "http://localhost:8025" | grep -q "200"; then
    echo -e "${GREEN}✅ OK${NC}"
else
    echo -e "${YELLOW}⚠️  Não crítico${NC}"
fi

# =============================================================================
# 9. INICIAR FRONTEND (OPCIONAL)
# =============================================================================
echo -e "${BLUE}💻 Frontend...${NC}"
echo -e "${YELLOW}Para iniciar o frontend, execute em outro terminal:${NC}"
echo -e "${GREEN}cd frontend && npm install && npm run dev${NC}"

# =============================================================================
# 10. RESUMO FINAL
# =============================================================================
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ AMBIENTE INICIADO COM SUCESSO!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${BLUE}📊 Infraestrutura Disponível:${NC}"
echo -e "  PostgreSQL:       ${GREEN}localhost:5432${NC}"
echo -e "  Redis:            ${GREEN}localhost:6379${NC}"
echo -e "  RabbitMQ:         ${GREEN}http://localhost:15672${NC}"
echo -e "  pgAdmin:          ${GREEN}http://localhost:5050${NC}"
echo -e "  Mailhog:          ${GREEN}http://localhost:8025${NC}"
echo ""
echo -e "${BLUE}🔐 Credenciais:${NC}"
echo -e "  PostgreSQL:       ${GREEN}direito_lux / direito_lux_pass_dev${NC}"
echo -e "  Redis:            ${GREEN}dev_redis_123${NC}"
echo -e "  RabbitMQ:         ${GREEN}direito_lux / dev_rabbit_123${NC}"
echo -e "  pgAdmin:          ${GREEN}admin@direitolux.com / dev_pgadmin_123${NC}"
echo ""
echo -e "${BLUE}🛠️ Comandos Úteis:${NC}"
echo -e "  Ver logs:         ${YELLOW}docker-compose -f docker-compose.simple.yml logs -f [serviço]${NC}"
echo -e "  Parar tudo:       ${YELLOW}docker-compose -f docker-compose.simple.yml down${NC}"
echo -e "  Limpar tudo:      ${YELLOW}docker-compose -f docker-compose.simple.yml down -v${NC}"
echo -e "  Status:           ${YELLOW}docker-compose -f docker-compose.simple.yml ps${NC}"
echo ""
echo -e "${BLUE}🚀 Para subir os microserviços:${NC}"
echo -e "  1. Crie/corrija os Dockerfile.dev que faltam"
echo -e "  2. Use: ${GREEN}docker-compose up -d${NC}"
echo -e "  3. Ou execute serviços individualmente com ${GREEN}go run${NC}"
echo ""
echo -e "${GREEN}🎉 Infraestrutura básica pronta!${NC}"