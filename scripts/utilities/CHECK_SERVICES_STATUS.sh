#!/bin/bash

# =============================================================================
# DIREITO LUX - CHECK SERVICES STATUS
# =============================================================================
# Script para verificar status de todos os serviços
# =============================================================================

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📊 DIREITO LUX - STATUS DOS SERVIÇOS${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# =============================================================================
# 1. DOCKER CONTAINERS
# =============================================================================
echo -e "${BLUE}🐳 Status dos Containers Docker:${NC}"
echo ""
docker-compose ps 2>/dev/null || echo -e "${RED}❌ Docker Compose não está rodando${NC}"
echo ""

# =============================================================================
# 2. HEALTH CHECKS DOS SERVIÇOS
# =============================================================================
echo -e "${BLUE}🏥 Health Checks dos Serviços:${NC}"
echo ""

check_service() {
    local name=$1
    local port=$2
    local endpoint=${3:-"/health"}
    local url="http://localhost:$port$endpoint"
    
    printf "%-20s" "$name:"
    
    # Verificar se porta está aberta
    if ! nc -z localhost $port 2>/dev/null; then
        echo -e "${RED}❌ Porta $port fechada${NC}"
        return
    fi
    
    # Fazer health check
    local response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    
    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✅ OK (HTTP 200)${NC}"
    elif [ "$response" = "000" ]; then
        echo -e "${YELLOW}⚠️  Sem resposta${NC}"
    else
        echo -e "${RED}❌ HTTP $response${NC}"
    fi
}

# Verificar todos os serviços
check_service "Auth Service" 8081
check_service "Tenant Service" 8082
check_service "Process Service" 8083
check_service "DataJud Service" 8084
check_service "Notification Service" 8085
check_service "AI Service" 8000
check_service "Search Service" 8086
check_service "Report Service" 8087
check_service "MCP Service" 8088

echo ""

# =============================================================================
# 3. INFRAESTRUTURA
# =============================================================================
echo -e "${BLUE}🏗️ Status da Infraestrutura:${NC}"
echo ""

# PostgreSQL
printf "%-20s" "PostgreSQL:"
if docker-compose exec -T postgres pg_isready -U postgres 2>/dev/null | grep -q "accepting connections"; then
    echo -e "${GREEN}✅ Conectado${NC}"
else
    echo -e "${RED}❌ Não conectado${NC}"
fi

# Redis
printf "%-20s" "Redis:"
if docker-compose exec -T redis redis-cli --no-auth-warning -a dev_redis_123 ping 2>/dev/null | grep -q PONG; then
    echo -e "${GREEN}✅ Respondendo${NC}"
else
    echo -e "${RED}❌ Não respondendo${NC}"
fi

# RabbitMQ
printf "%-20s" "RabbitMQ:"
if docker-compose exec -T rabbitmq rabbitmq-diagnostics ping 2>/dev/null | grep -q "Ping succeeded"; then
    echo -e "${GREEN}✅ Rodando${NC}"
elif docker-compose exec -T rabbitmq rabbitmqctl status 2>/dev/null | grep -q "Running"; then
    echo -e "${GREEN}✅ Rodando (status)${NC}"
else
    echo -e "${RED}❌ Não rodando${NC}"
fi

# Elasticsearch (se existir)
printf "%-20s" "Elasticsearch:"
if curl -s "http://localhost:9200/_cluster/health" 2>/dev/null | grep -q "green\|yellow"; then
    echo -e "${GREEN}✅ Cluster OK${NC}"
else
    echo -e "${RED}❌ Indisponível${NC}"
fi

echo ""

# =============================================================================
# 4. FRONTEND
# =============================================================================
echo -e "${BLUE}💻 Frontend Status:${NC}"
echo ""

printf "%-20s" "Next.js (3000):"
if nc -z localhost 3000 2>/dev/null; then
    echo -e "${GREEN}✅ Rodando${NC}"
else
    echo -e "${YELLOW}⚠️  Não iniciado${NC}"
    echo -e "    Para iniciar: ${GREEN}cd frontend && npm run dev${NC}"
fi

echo ""

# =============================================================================
# 5. RECURSOS DO SISTEMA
# =============================================================================
echo -e "${BLUE}⚡ Recursos do Sistema:${NC}"
echo ""

# Docker stats resumido
echo -e "${YELLOW}💾 Uso de Memória Docker:${NC}"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" 2>/dev/null | head -10

echo ""

# =============================================================================
# 6. URLS DISPONÍVEIS
# =============================================================================
echo -e "${BLUE}🌐 URLs Disponíveis:${NC}"
echo ""
echo -e "  Frontend:         ${GREEN}http://localhost:3000${NC}"
echo -e "  Auth Service:     ${GREEN}http://localhost:8081/health${NC}"
echo -e "  AI Service:       ${GREEN}http://localhost:8000/health${NC}"
echo -e "  pgAdmin:          ${GREEN}http://localhost:5050${NC}"
echo -e "  RabbitMQ:         ${GREEN}http://localhost:15672${NC}"
echo -e "  Mailhog:          ${GREEN}http://localhost:8025${NC}"

echo ""

# =============================================================================
# 7. LOGS RECENTES
# =============================================================================
echo -e "${BLUE}📄 Logs Recentes (últimas 5 linhas):${NC}"
echo ""

# Mostrar logs recentes dos serviços principais
for service in postgres redis rabbitmq; do
    echo -e "${YELLOW}$service:${NC}"
    docker-compose logs --tail=2 $service 2>/dev/null | tail -2 || echo "  Não disponível"
    echo ""
done

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}📊 Verificação de status concluída!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${YELLOW}💡 Para logs detalhados:${NC} docker-compose logs -f [serviço]"
echo -e "${YELLOW}💡 Para reiniciar:${NC} ./START_CLEAN_ENVIRONMENT.sh"
echo ""