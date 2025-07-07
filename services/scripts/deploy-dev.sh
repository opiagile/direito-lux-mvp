#!/bin/bash

# Direito Lux - Deploy Development Environment
# Deploys AI Service, Search Service, MCP Service and all dependencies

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$PROJECT_DIR/docker-compose.dev.yml"
SERVICES=("postgres" "redis" "rabbitmq" "elasticsearch" "jaeger" "ai-service" "search-service")
MCP_SERVICES=("mcp-postgres" "mcp-redis" "mcp-rabbitmq")

echo -e "${BLUE}🚀 Direito Lux - Deploy Development Environment${NC}"
echo -e "${CYAN}=====================================================${NC}"

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_step() {
    echo -e "${PURPLE}🔄 $1${NC}"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker não está rodando. Por favor, inicie o Docker primeiro."
        exit 1
    fi
    print_status "Docker está rodando"
}

# Function to check if docker-compose is available
check_compose() {
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        print_error "Docker Compose não encontrado. Por favor, instale o Docker Compose."
        exit 1
    fi
    print_status "Docker Compose disponível: $COMPOSE_CMD"
}

# Function to cleanup existing containers
cleanup() {
    print_step "Limpando containers existentes..."
    $COMPOSE_CMD -f "$COMPOSE_FILE" down --remove-orphans || true
    
    if [ "$1" = "--clean" ]; then
        print_step "Removendo volumes existentes..."
        $COMPOSE_CMD -f "$COMPOSE_FILE" down -v
        docker system prune -f --volumes
        print_warning "Volumes removidos - dados serão perdidos!"
    fi
}

# Function to pull latest images
pull_images() {
    print_step "Baixando imagens mais recentes..."
    $COMPOSE_CMD -f "$COMPOSE_FILE" pull
    print_status "Imagens atualizadas"
}

# Function to build services
build_services() {
    print_step "Construindo serviços customizados..."
    $COMPOSE_CMD -f "$COMPOSE_FILE" build --no-cache
    print_status "Serviços construídos"
}

# Function to start infrastructure services
start_infrastructure() {
    print_step "Iniciando infraestrutura (PostgreSQL, Redis, RabbitMQ, Elasticsearch, Jaeger)..."
    
    # Start infrastructure services first
    $COMPOSE_CMD -f "$COMPOSE_FILE" up -d postgres redis rabbitmq elasticsearch jaeger
    $COMPOSE_CMD -f "$COMPOSE_FILE" up -d mcp-postgres mcp-redis mcp-rabbitmq
    
    print_info "Aguardando serviços de infraestrutura ficarem prontos..."
    
    # Wait for services to be healthy
    wait_for_service "postgres" "PostgreSQL (Main)"
    wait_for_service "redis" "Redis (Main)"
    wait_for_service "rabbitmq" "RabbitMQ (Main)"
    wait_for_service "elasticsearch" "Elasticsearch"
    wait_for_service "mcp-postgres" "PostgreSQL (MCP)"
    wait_for_service "mcp-redis" "Redis (MCP)"
    wait_for_service "mcp-rabbitmq" "RabbitMQ (MCP)"
    
    print_status "Infraestrutura pronta!"
}

# Function to wait for a service to be healthy
wait_for_service() {
    local service=$1
    local name=$2
    local max_attempts=30
    local attempt=1
    
    print_info "   • Aguardando $name..."
    
    while [ $attempt -le $max_attempts ]; do
        if $COMPOSE_CMD -f "$COMPOSE_FILE" ps "$service" | grep -q "healthy\|Up"; then
            print_status "   • $name está pronto!"
            return 0
        fi
        
        if [ $((attempt % 5)) -eq 0 ]; then
            print_info "   • $name ainda não está pronto (tentativa $attempt/$max_attempts)..."
        fi
        
        sleep 2
        ((attempt++))
    done
    
    print_error "Timeout aguardando $name ficar pronto"
    return 1
}

# Function to start application services
start_applications() {
    print_step "Iniciando serviços de aplicação..."
    
    # Start AI Service
    print_info "   • Iniciando AI Service..."
    $COMPOSE_CMD -f "$COMPOSE_FILE" up -d ai-service
    
    # Start Search Service
    print_info "   • Iniciando Search Service..."
    $COMPOSE_CMD -f "$COMPOSE_FILE" up -d search-service
    
    # Wait for services to be healthy
    print_info "Aguardando serviços de aplicação ficarem prontos..."
    wait_for_service "ai-service" "AI Service"
    wait_for_service "search-service" "Search Service"
    
    print_status "Serviços de aplicação prontos!"
}

# Function to show service status
show_status() {
    echo ""
    print_step "Status dos serviços:"
    $COMPOSE_CMD -f "$COMPOSE_FILE" ps
    
    echo ""
    print_step "Logs recentes:"
    $COMPOSE_CMD -f "$COMPOSE_FILE" logs --tail=5 ai-service search-service
}

# Function to show service endpoints
show_endpoints() {
    echo ""
    print_info "🌐 Endpoints disponíveis:"
    echo ""
    echo -e "${CYAN}📊 Serviços Principais:${NC}"
    echo -e "   • AI Service:           ${GREEN}http://localhost:8000${NC}"
    echo -e "   • Search Service:       ${GREEN}http://localhost:8086${NC}"
    echo -e "   • AI Service Health:    ${GREEN}http://localhost:8000/health${NC}"
    echo -e "   • Search Service Health:${GREEN}http://localhost:8086/health${NC}"
    echo ""
    echo -e "${CYAN}🗄️  Infraestrutura:${NC}"
    echo -e "   • PostgreSQL (Main):    ${GREEN}localhost:5432${NC} (direito_lux/direito_lux_pass_dev)"
    echo -e "   • PostgreSQL (MCP):     ${GREEN}localhost:5434${NC} (mcp_user/mcp_pass_dev)"
    echo -e "   • Redis (Main):         ${GREEN}localhost:6379${NC} (redis_pass_dev)"
    echo -e "   • Redis (MCP):          ${GREEN}localhost:6380${NC} (redis_pass_dev)"
    echo -e "   • RabbitMQ (Main):      ${GREEN}localhost:5672${NC} (direito_lux/rabbit_pass_dev)"
    echo -e "   • RabbitMQ (MCP):       ${GREEN}localhost:5673${NC} (mcp_user/rabbit_pass_dev)"
    echo -e "   • Elasticsearch:        ${GREEN}http://localhost:9200${NC}"
    echo ""
    echo -e "${CYAN}📈 Monitoramento:${NC}"
    echo -e "   • RabbitMQ Management:  ${GREEN}http://localhost:15672${NC} (direito_lux/rabbit_pass_dev)"
    echo -e "   • RabbitMQ Mgmt (MCP):  ${GREEN}http://localhost:15673${NC} (mcp_user/rabbit_pass_dev)"
    echo -e "   • Jaeger UI:            ${GREEN}http://localhost:16686${NC}"
    echo ""
}

# Function to run tests
run_tests() {
    print_step "Executando testes de conectividade..."
    
    # Test AI Service
    if curl -s -f http://localhost:8000/health > /dev/null; then
        print_status "AI Service respondendo"
    else
        print_warning "AI Service pode não estar respondendo ainda"
    fi
    
    # Test Search Service
    if curl -s -f http://localhost:8086/health > /dev/null; then
        print_status "Search Service respondendo"
    else
        print_warning "Search Service pode não estar respondendo ainda"
    fi
    
    # Test Elasticsearch
    if curl -s -f http://localhost:9200/_health > /dev/null; then
        print_status "Elasticsearch respondendo"
    else
        print_warning "Elasticsearch pode não estar respondendo ainda"
    fi
}

# Function to show logs
show_logs() {
    local service=${1:-""}
    
    if [ -n "$service" ]; then
        print_info "Mostrando logs do serviço: $service"
        $COMPOSE_CMD -f "$COMPOSE_FILE" logs -f "$service"
    else
        print_info "Mostrando logs de todos os serviços:"
        $COMPOSE_CMD -f "$COMPOSE_FILE" logs -f
    fi
}

# Function to stop services
stop_services() {
    print_step "Parando todos os serviços..."
    $COMPOSE_CMD -f "$COMPOSE_FILE" down
    print_status "Serviços parados"
}

# Function to show help
show_help() {
    echo "Uso: $0 [opções] [comando]"
    echo ""
    echo "Comandos:"
    echo "  start     - Inicia todos os serviços (padrão)"
    echo "  stop      - Para todos os serviços"
    echo "  restart   - Reinicia todos os serviços"
    echo "  status    - Mostra status dos serviços"
    echo "  logs      - Mostra logs (use logs <service> para serviço específico)"
    echo "  test      - Executa testes de conectividade"
    echo "  endpoints - Mostra endpoints disponíveis"
    echo ""
    echo "Opções:"
    echo "  --clean   - Remove volumes existentes (CUIDADO: apaga dados!)"
    echo "  --build   - Reconstrói as imagens antes de iniciar"
    echo "  --pull    - Baixa imagens mais recentes antes de iniciar"
    echo "  --no-test - Não executa testes após inicialização"
    echo "  --help    - Mostra esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0                    # Inicia todos os serviços"
    echo "  $0 --clean start     # Limpa volumes e inicia"
    echo "  $0 --build start     # Reconstrói e inicia"
    echo "  $0 logs ai-service   # Mostra logs do AI Service"
    echo "  $0 stop              # Para todos os serviços"
}

# Parse command line arguments
CLEAN=false
BUILD=false
PULL=false
NO_TEST=false
COMMAND="start"

while [[ $# -gt 0 ]]; do
    case $1 in
        --clean)
            CLEAN=true
            shift
            ;;
        --build)
            BUILD=true
            shift
            ;;
        --pull)
            PULL=true
            shift
            ;;
        --no-test)
            NO_TEST=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        start|stop|restart|status|logs|test|endpoints)
            COMMAND=$1
            shift
            ;;
        *)
            if [ "$COMMAND" = "logs" ] && [ -n "$1" ]; then
                LOGS_SERVICE=$1
                shift
            else
                print_error "Opção desconhecida: $1"
                show_help
                exit 1
            fi
            ;;
    esac
done

# Main execution
cd "$PROJECT_DIR"

case $COMMAND in
    start)
        check_docker
        check_compose
        
        if [ "$CLEAN" = true ]; then
            cleanup --clean
        else
            cleanup
        fi
        
        if [ "$PULL" = true ]; then
            pull_images
        fi
        
        if [ "$BUILD" = true ]; then
            build_services
        fi
        
        start_infrastructure
        start_applications
        
        if [ "$NO_TEST" = false ]; then
            run_tests
        fi
        
        show_endpoints
        show_status
        
        print_status "Deploy concluído com sucesso!"
        print_info "Use '$0 logs' para ver logs em tempo real"
        print_info "Use '$0 stop' para parar todos os serviços"
        ;;
        
    stop)
        stop_services
        ;;
        
    restart)
        check_docker
        check_compose
        cleanup
        start_infrastructure
        start_applications
        show_endpoints
        print_status "Restart concluído!"
        ;;
        
    status)
        show_status
        ;;
        
    logs)
        show_logs "$LOGS_SERVICE"
        ;;
        
    test)
        run_tests
        ;;
        
    endpoints)
        show_endpoints
        ;;
        
    *)
        print_error "Comando desconhecido: $COMMAND"
        show_help
        exit 1
        ;;
esac