#!/bin/bash

# Script para monitorar teste de usuário em tempo real
# Uso: ./monitor-teste-usuario.sh [start|stop|status|errors|db]

PROJECT_ID="direito-lux-staging-2025"
NAMESPACE="direito-lux-staging"
LOG_FILE="teste-usuario-$(date +%Y%m%d-%H%M%S).log"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para log com timestamp
log_message() {
    echo -e "[$(date '+%H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Função para verificar se sistema está pronto
check_system_ready() {
    log_message "${BLUE}🔍 Verificando se sistema está pronto...${NC}"
    
    # Verificar pods
    local pods_running=$(kubectl get pods -n $NAMESPACE --no-headers | grep Running | wc -l)
    local pods_total=$(kubectl get pods -n $NAMESPACE --no-headers | wc -l)
    
    log_message "📦 Pods rodando: $pods_running/$pods_total"
    
    # Verificar serviços críticos
    local critical_services=("auth-service" "tenant-service" "frontend" "postgres")
    local all_ready=true
    
    for service in "${critical_services[@]}"; do
        local status=$(kubectl get pods -n $NAMESPACE -l app=$service --no-headers | grep Running | wc -l)
        if [ "$status" -eq 0 ]; then
            log_message "${RED}❌ $service não está rodando${NC}"
            all_ready=false
        else
            log_message "${GREEN}✅ $service está rodando${NC}"
        fi
    done
    
    # Testar conectividade
    if curl -sk https://35.188.198.87/api/health > /dev/null 2>&1; then
        log_message "${GREEN}✅ Frontend acessível${NC}"
    else
        log_message "${RED}❌ Frontend não está acessível${NC}"
        all_ready=false
    fi
    
    if [ "$all_ready" = true ]; then
        log_message "${GREEN}🎉 Sistema pronto para teste!${NC}"
        return 0
    else
        log_message "${RED}⚠️ Sistema não está completamente pronto${NC}"
        return 1
    fi
}

# Função para monitorar erros em tempo real
monitor_errors() {
    log_message "${YELLOW}🔍 Iniciando monitoramento de erros...${NC}"
    
    # Criar função para capturar erros
    kubectl logs -n $NAMESPACE --all-containers=true -f | \
    grep --color=always -E "(ERROR|ERRO|error|Error|failed|Failed|FAILED|panic|PANIC|Exception|exception|401|403|404|500|502|503)" | \
    while read -r line; do
        log_message "${RED}🚨 ERRO: $line${NC}"
    done
}

# Função para verificar status dos serviços
check_services_status() {
    log_message "${BLUE}📊 Status dos serviços:${NC}"
    
    # Verificar cada serviço
    local services=("auth-service" "tenant-service" "frontend" "postgres" "redis" "rabbitmq")
    
    for service in "${services[@]}"; do
        local pod_count=$(kubectl get pods -n $NAMESPACE -l app=$service --no-headers | wc -l)
        local running_count=$(kubectl get pods -n $NAMESPACE -l app=$service --no-headers | grep Running | wc -l)
        
        if [ "$pod_count" -eq 0 ]; then
            log_message "${YELLOW}⚠️ $service: Nenhum pod encontrado${NC}"
        elif [ "$running_count" -eq "$pod_count" ]; then
            log_message "${GREEN}✅ $service: $running_count/$pod_count pods Running${NC}"
        else
            log_message "${RED}❌ $service: $running_count/$pod_count pods Running${NC}"
            
            # Mostrar detalhes dos pods com problema
            kubectl get pods -n $NAMESPACE -l app=$service | grep -v Running | while read -r line; do
                log_message "${RED}   📋 $line${NC}"
            done
        fi
    done
}

# Função para verificar dados no banco
check_database_data() {
    log_message "${BLUE}💾 Verificando dados no banco:${NC}"
    
    # Verificar conectividade com PostgreSQL
    if kubectl exec -n $NAMESPACE deploy/postgres -- pg_isready -U direito_lux > /dev/null 2>&1; then
        log_message "${GREEN}✅ PostgreSQL acessível${NC}"
        
        # Estatísticas gerais
        local tenants=$(kubectl exec -n $NAMESPACE deploy/postgres -- psql -U direito_lux -d direito_lux_staging -t -c "SELECT COUNT(*) FROM tenants;" 2>/dev/null | tr -d ' ' | head -1)
        local users=$(kubectl exec -n $NAMESPACE deploy/postgres -- psql -U direito_lux -d direito_lux_staging -t -c "SELECT COUNT(*) FROM users;" 2>/dev/null | tr -d ' ' | head -1)
        local processes=$(kubectl exec -n $NAMESPACE deploy/postgres -- psql -U direito_lux -d direito_lux_staging -t -c "SELECT COUNT(*) FROM processes;" 2>/dev/null | tr -d ' ' | head -1)
        
        log_message "📊 Estatísticas:"
        log_message "   👥 Tenants: $tenants"
        log_message "   👤 Users: $users"
        log_message "   📋 Processes: $processes"
    else
        log_message "${RED}❌ PostgreSQL não está acessível${NC}"
    fi
}

# Função para verificar performance
check_performance() {
    log_message "${BLUE}🚀 Verificando performance:${NC}"
    
    # Verificar uso de recursos
    log_message "💾 Uso de recursos (Top 5):"
    kubectl top pods -n $NAMESPACE 2>/dev/null | head -6 | while read -r line; do
        log_message "   $line"
    done
    
    # Verificar requests por minuto
    local requests_1m=$(kubectl logs -n $NAMESPACE -l app=frontend --since=1m 2>/dev/null | grep -c "GET\|POST" || echo "0")
    log_message "📈 Requests/minuto: $requests_1m"
    
    # Verificar latência média (últimos 10 requests)
    log_message "⏱️ Latências recentes:"
    kubectl logs -n $NAMESPACE -l app=frontend --since=5m 2>/dev/null | grep -o "[0-9]*ms" | tail -10 | while read -r latency; do
        log_message "   $latency"
    done
}

# Função para monitorar cadastro
monitor_signup() {
    log_message "${YELLOW}👀 Monitorando fluxo de cadastro...${NC}"
    
    # Monitorar tenant-service
    kubectl logs -n $NAMESPACE -l app=tenant-service -f | grep -E "(POST.*tenants|tenant.*created|error)" | while read -r line; do
        log_message "${BLUE}🏢 TENANT: $line${NC}"
    done &
    
    # Monitorar auth-service
    kubectl logs -n $NAMESPACE -l app=auth-service -f | grep -E "(POST.*register|user.*created|password|error)" | while read -r line; do
        log_message "${GREEN}🔐 AUTH: $line${NC}"
    done &
    
    # Aguardar Ctrl+C
    wait
}

# Função para monitorar login
monitor_login() {
    local email=${1:-""}
    log_message "${YELLOW}🔐 Monitorando tentativas de login${NC}"
    
    if [ -n "$email" ]; then
        log_message "   👤 Focando em: $email"
        kubectl logs -n $NAMESPACE -l app=auth-service -f | grep -E "(login.*$email|JWT|token|401|403|success|failed)" | while read -r line; do
            log_message "${GREEN}🔐 LOGIN: $line${NC}"
        done
    else
        kubectl logs -n $NAMESPACE -l app=auth-service -f | grep -E "(login|JWT|token|401|403|success|failed)" | while read -r line; do
            log_message "${GREEN}🔐 LOGIN: $line${NC}"
        done
    fi
}

# Função para dashboard completo
dashboard() {
    while true; do
        clear
        echo -e "${BLUE}=== DIREITO LUX - DASHBOARD DE MONITORAMENTO ===${NC}"
        echo -e "${BLUE}Horário: $(date)${NC}"
        echo ""
        
        # Status dos pods
        echo -e "${GREEN}PODS STATUS:${NC}"
        kubectl get pods -n $NAMESPACE | grep -E "(NAME|Running|Error|Crash|Pending)" | head -10
        echo ""
        
        # Últimos erros
        echo -e "${RED}ÚLTIMOS ERROS (2 min):${NC}"
        kubectl logs -n $NAMESPACE --all-containers=true --since=2m 2>/dev/null | grep -i error | tail -3
        echo ""
        
        # Requests por minuto
        echo -e "${YELLOW}REQUESTS/MIN:${NC}"
        kubectl logs -n $NAMESPACE -l app=frontend --since=1m 2>/dev/null | grep -c "GET\|POST"
        echo ""
        
        # Uso de recursos
        echo -e "${BLUE}TOP 5 RECURSOS:${NC}"
        kubectl top pods -n $NAMESPACE 2>/dev/null | head -6
        echo ""
        
        # Próxima atualização
        echo -e "${BLUE}Próxima atualização em 10 segundos... (Ctrl+C para parar)${NC}"
        sleep 10
    done
}

# Função para teste completo
run_full_test() {
    log_message "${BLUE}🧪 Iniciando teste completo de usuário...${NC}"
    
    # 1. Verificar sistema
    if ! check_system_ready; then
        log_message "${RED}❌ Sistema não está pronto. Abortando teste.${NC}"
        return 1
    fi
    
    # 2. Backup de dados atuais
    log_message "${YELLOW}💾 Fazendo backup dos dados atuais...${NC}"
    check_database_data
    
    # 3. Iniciar monitoramento em background
    log_message "${YELLOW}🔍 Iniciando monitoramento de erros...${NC}"
    monitor_errors > "errors-${LOG_FILE}" &
    MONITOR_PID=$!
    
    # 4. Aguardar teste manual
    log_message "${GREEN}✅ Sistema monitorado e pronto para teste!${NC}"
    log_message "${YELLOW}📋 Execute o teste manual conforme ROTEIRO_TESTE_USUARIO_COMPLETO.md${NC}"
    log_message "${YELLOW}🔍 Monitoramento de erros ativo (PID: $MONITOR_PID)${NC}"
    log_message "${YELLOW}📝 Logs salvos em: $LOG_FILE${NC}"
    
    # 5. Aguardar comando para finalizar
    echo "Pressione ENTER para finalizar o monitoramento..."
    read -r
    
    # 6. Finalizar monitoramento
    kill $MONITOR_PID 2>/dev/null
    
    # 7. Relatório final
    log_message "${BLUE}📊 Relatório final:${NC}"
    check_services_status
    check_database_data
    check_performance
    
    log_message "${GREEN}✅ Teste completo finalizado!${NC}"
    log_message "${YELLOW}📁 Logs salvos em: $LOG_FILE${NC}"
}

# Menu principal
case "${1:-help}" in
    "start")
        log_message "${GREEN}🚀 Iniciando sistema e monitoramento...${NC}"
        ./scripts/gcp-cost-optimizer.sh start
        sleep 60
        check_system_ready
        ;;
    "stop")
        log_message "${RED}🛑 Parando sistema...${NC}"
        ./scripts/gcp-cost-optimizer.sh stop
        ;;
    "status")
        check_system_ready
        check_services_status
        ;;
    "errors")
        monitor_errors
        ;;
    "db")
        check_database_data
        ;;
    "performance")
        check_performance
        ;;
    "signup")
        monitor_signup
        ;;
    "login")
        monitor_login "$2"
        ;;
    "dashboard")
        dashboard
        ;;
    "test")
        run_full_test
        ;;
    "help"|*)
        echo -e "${BLUE}🧪 MONITOR DE TESTE DE USUÁRIO${NC}"
        echo "=================================="
        echo ""
        echo "Comandos disponíveis:"
        echo "  start      - Iniciar sistema e verificar"
        echo "  stop       - Parar sistema"
        echo "  status     - Verificar status geral"
        echo "  errors     - Monitorar erros em tempo real"
        echo "  db         - Verificar dados no banco"
        echo "  performance - Verificar performance"
        echo "  signup     - Monitorar fluxo de cadastro"
        echo "  login [email] - Monitorar tentativas de login"
        echo "  dashboard  - Dashboard completo"
        echo "  test       - Teste completo automatizado"
        echo ""
        echo "Exemplo de uso:"
        echo "  ./monitor-teste-usuario.sh start"
        echo "  ./monitor-teste-usuario.sh dashboard"
        echo "  ./monitor-teste-usuario.sh login joao@costaadvogados.com.br"
        echo ""
        echo "📝 Logs são salvos em: teste-usuario-YYYYMMDD-HHMMSS.log"
        ;;
esac