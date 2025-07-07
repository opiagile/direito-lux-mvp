#!/bin/bash

# Script para executar testes E2E do Direito Lux
# Usage: ./run-tests.sh [options]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Diretório base
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

echo -e "${BLUE}🧪 Direito Lux - E2E Test Runner${NC}"
echo -e "${BLUE}===================================${NC}"

# Função para mostrar help
show_help() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --auth          Executar apenas testes de autenticação"
    echo "  --processes     Executar apenas testes de processos"  
    echo "  --dashboard     Executar apenas testes de dashboard"
    echo "  --full-flow     Executar apenas testes de fluxo completo"
    echo "  --setup         Apenas executar setup e verificações"
    echo "  --install       Instalar dependências npm"
    echo "  --verbose       Output verboso"
    echo "  --help          Mostrar esta ajuda"
    echo ""
    echo "Examples:"
    echo "  $0                    # Executar todos os testes"
    echo "  $0 --auth            # Apenas testes de auth"
    echo "  $0 --setup           # Apenas verificar ambiente"
    echo "  $0 --install         # Instalar dependências"
}

# Parse arguments
VERBOSE=false
TEST_SUITE="all"
INSTALL_DEPS=false
SETUP_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --auth)
            TEST_SUITE="auth"
            shift
            ;;
        --processes)
            TEST_SUITE="processes"
            shift
            ;;
        --dashboard)
            TEST_SUITE="dashboard"
            shift
            ;;
        --full-flow)
            TEST_SUITE="full-flow"
            shift
            ;;
        --setup)
            SETUP_ONLY=true
            shift
            ;;
        --install)
            INSTALL_DEPS=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Opção desconhecida: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Função para log
log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${BLUE}[$(date +'%H:%M:%S')] $1${NC}"
    fi
}

# Função para executar comando com log
run_cmd() {
    local cmd="$1"
    local desc="$2"
    
    echo -e "${YELLOW}⏳ $desc...${NC}"
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${BLUE}💻 Executando: $cmd${NC}"
        eval "$cmd"
    else
        eval "$cmd" > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✅ $desc - Concluído${NC}"
    else
        echo -e "${RED}❌ $desc - Falhou${NC}"
        exit 1
    fi
}

# Navegar para diretório dos testes
cd "$SCRIPT_DIR"

# 1. Instalar dependências se solicitado
if [[ "$INSTALL_DEPS" == "true" ]]; then
    echo -e "\n${BLUE}📦 Instalando dependências...${NC}"
    run_cmd "npm install" "Instalação de dependências"
fi

# 2. Verificar se dependências estão instaladas
if [[ ! -d "node_modules" ]]; then
    echo -e "${YELLOW}⚠️  Dependências não encontradas. Instalando...${NC}"
    run_cmd "npm install" "Instalação automática de dependências"
fi

# 3. Verificar se serviços estão rodando
echo -e "\n${BLUE}🔍 Verificando disponibilidade dos serviços...${NC}"

# Lista de serviços críticos
SERVICES=(
    "localhost:8081"  # auth-service
    "localhost:8083"  # process-service  
    "localhost:8087"  # report-service
)

FAILED_SERVICES=()

for service in "${SERVICES[@]}"; do
    log "Verificando $service"
    
    if curl -s --max-time 5 "http://$service/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $service - OK${NC}"
    else
        echo -e "${RED}❌ $service - INDISPONÍVEL${NC}"
        FAILED_SERVICES+=("$service")
    fi
done

# Se serviços críticos estão indisponíveis, orientar
if [[ ${#FAILED_SERVICES[@]} -gt 0 ]]; then
    echo -e "\n${RED}💥 Serviços indisponíveis detectados!${NC}"
    echo -e "${YELLOW}📋 Para resolver, execute:${NC}"
    echo -e "   cd $PROJECT_ROOT"
    echo -e "   ./services/scripts/deploy-dev.sh start"
    echo -e ""
    echo -e "${YELLOW}💡 Ou para reiniciar tudo:${NC}"
    echo -e "   ./services/scripts/deploy-dev.sh --clean --build start"
    
    # Permitir continuar se não for crítico
    read -p "Continuar mesmo assim? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 4. Se apenas setup, parar aqui
if [[ "$SETUP_ONLY" == "true" ]]; then
    echo -e "\n${GREEN}✅ Verificação de ambiente concluída${NC}"
    echo -e "${BLUE}💡 Execute sem --setup para rodar os testes${NC}"
    exit 0
fi

# 5. Executar testes
echo -e "\n${BLUE}🧪 Executando testes E2E...${NC}"

case $TEST_SUITE in
    "auth")
        echo -e "${YELLOW}🔑 Executando testes de autenticação...${NC}"
        run_cmd "npm run test:auth" "Testes de autenticação"
        ;;
    "processes")
        echo -e "${YELLOW}📋 Executando testes de processos...${NC}"
        run_cmd "npm run test:processes" "Testes de processos"
        ;;
    "dashboard")
        echo -e "${YELLOW}📊 Executando testes de dashboard...${NC}"
        run_cmd "npm run test:dashboard" "Testes de dashboard"
        ;;
    "full-flow")
        echo -e "${YELLOW}🔄 Executando testes de fluxo completo...${NC}"
        run_cmd "npm run test:full-flow" "Testes de fluxo completo"
        ;;
    "all")
        echo -e "${YELLOW}🎯 Executando todos os testes...${NC}"
        
        # Executar em ordem lógica
        echo -e "\n${BLUE}1/4 - Testes de Autenticação${NC}"
        run_cmd "npm run test:auth" "Testes de autenticação"
        
        echo -e "\n${BLUE}2/4 - Testes de Dashboard${NC}"
        run_cmd "npm run test:dashboard" "Testes de dashboard"
        
        echo -e "\n${BLUE}3/4 - Testes de Processos${NC}"
        run_cmd "npm run test:processes" "Testes de processos"
        
        echo -e "\n${BLUE}4/4 - Testes de Fluxo Completo${NC}"
        run_cmd "npm run test:full-flow" "Testes de fluxo completo"
        ;;
esac

# 6. Resultado final
echo -e "\n${GREEN}🎉 Todos os testes executados com sucesso!${NC}"
echo -e "${BLUE}📊 Resumo:${NC}"
echo -e "   - Ambiente: ✅ Verificado"
echo -e "   - Testes: ✅ Executados"
echo -e "   - Suite: $TEST_SUITE"

# 7. Próximos passos
echo -e "\n${BLUE}🚀 Próximos passos sugeridos:${NC}"
echo -e "   1. Revisar logs acima para verificar detalhes"
echo -e "   2. Executar testes específicos se necessário"
echo -e "   3. Implementar correções se algum teste falhou"
echo -e "   4. Continuar com desenvolvimento ou deploy"

echo -e "\n${BLUE}💡 Comandos úteis:${NC}"
echo -e "   $0 --auth           # Testar apenas autenticação"
echo -e "   $0 --setup          # Verificar ambiente"
echo -e "   $0 --verbose        # Output detalhado"

echo -e "\n${GREEN}✅ E2E Test Runner finalizado${NC}"