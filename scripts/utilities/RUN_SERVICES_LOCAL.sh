#!/bin/bash

# =============================================================================
# DIREITO LUX - RUN SERVICES LOCALLY
# =============================================================================
# Script para executar os microserviços Go localmente (sem Docker)
# =============================================================================

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🚀 DIREITO LUX - EXECUTAR SERVIÇOS${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Diretório base
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$BASE_DIR"

# Verificar se infraestrutura está rodando
echo -e "${YELLOW}🔍 Verificando infraestrutura...${NC}"

check_service() {
    local service=$1
    local port=$2
    local name=$3
    
    if nc -z localhost $port 2>/dev/null; then
        echo -e "  $name: ${GREEN}✅ OK${NC}"
        return 0
    else
        echo -e "  $name: ${RED}❌ NÃO ENCONTRADO${NC}"
        return 1
    fi
}

INFRA_OK=true
check_service postgres 5432 "PostgreSQL" || INFRA_OK=false
check_service redis 6379 "Redis" || INFRA_OK=false
check_service rabbitmq 5672 "RabbitMQ" || INFRA_OK=false

if [ "$INFRA_OK" = false ]; then
    echo ""
    echo -e "${RED}❌ Infraestrutura não está completa!${NC}"
    echo -e "${YELLOW}Execute primeiro: ./START_CLEAN_ENVIRONMENT.sh${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Infraestrutura OK${NC}"
echo ""

# Função para executar serviço
run_service() {
    local service_name=$1
    local service_dir=$2
    local port=$3
    
    echo -e "${BLUE}Starting $service_name on port $port...${NC}"
    
    cd "$BASE_DIR/services/$service_dir"
    
    # Verificar se go.mod existe
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ go.mod não encontrado em $service_dir${NC}"
        return
    fi
    
    # Configurar variáveis de ambiente
    export PORT=$port
    export SERVER_PORT=$port
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_NAME=direito_lux_dev
    export DB_USER=direito_lux
    export DB_PASSWORD=direito_lux_pass_dev
    export REDIS_HOST=localhost
    export REDIS_PORT=6379
    export REDIS_PASSWORD=dev_redis_123
    export RABBITMQ_URL=amqp://direito_lux:dev_rabbit_123@localhost:5672/direito_lux
    export LOG_LEVEL=debug
    
    # Baixar dependências
    echo -e "${YELLOW}Baixando dependências...${NC}"
    go mod download
    
    # Executar serviço em background
    echo -e "${GREEN}Iniciando $service_name...${NC}"
    go run cmd/server/main.go > "../../logs/${service_dir}.log" 2>&1 &
    
    echo -e "${GREEN}✅ $service_name iniciado (PID: $!)${NC}"
    echo ""
}

# Menu de seleção
echo -e "${BLUE}Selecione os serviços para executar:${NC}"
echo "1) Auth Service (8081)"
echo "2) Tenant Service (8082)"
echo "3) Process Service (8083)"
echo "4) DataJud Service (8084)"
echo "5) Notification Service (8085)"
echo "6) Search Service (8086)"
echo "7) Report Service (8087)"
echo "8) MCP Service (8088)"
echo "9) Todos os serviços"
echo "0) Sair"
echo ""

read -p "Opção: " choice

case $choice in
    1) run_service "Auth Service" "auth-service" 8081 ;;
    2) run_service "Tenant Service" "tenant-service" 8082 ;;
    3) run_service "Process Service" "process-service" 8083 ;;
    4) run_service "DataJud Service" "datajud-service" 8084 ;;
    5) run_service "Notification Service" "notification-service" 8085 ;;
    6) run_service "Search Service" "search-service" 8086 ;;
    7) run_service "Report Service" "report-service" 8087 ;;
    8) run_service "MCP Service" "mcp-service" 8088 ;;
    9)
        echo -e "${BLUE}Iniciando todos os serviços...${NC}"
        run_service "Auth Service" "auth-service" 8081
        sleep 2
        run_service "Tenant Service" "tenant-service" 8082
        sleep 2
        run_service "Process Service" "process-service" 8083
        sleep 2
        run_service "DataJud Service" "datajud-service" 8084
        sleep 2
        run_service "Notification Service" "notification-service" 8085
        sleep 2
        run_service "Search Service" "search-service" 8086
        sleep 2
        run_service "Report Service" "report-service" 8087
        sleep 2
        run_service "MCP Service" "mcp-service" 8088
        ;;
    0)
        echo -e "${YELLOW}Saindo...${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}Opção inválida!${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ Serviços iniciados!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${BLUE}📋 Comandos úteis:${NC}"
echo -e "  Ver logs:     ${YELLOW}tail -f logs/[service-name].log${NC}"
echo -e "  Matar todos:  ${YELLOW}pkill -f 'go run'${NC}"
echo -e "  Status:       ${YELLOW}ps aux | grep 'go run'${NC}"
echo ""
echo -e "${BLUE}🌐 Testar serviços:${NC}"
echo -e "  Auth:         ${GREEN}curl http://localhost:8081/health${NC}"
echo -e "  Tenant:       ${GREEN}curl http://localhost:8082/health${NC}"
echo -e "  Process:      ${GREEN}curl http://localhost:8083/health${NC}"
echo ""