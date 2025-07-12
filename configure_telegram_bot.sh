#!/bin/bash

# 🤖 Configurar Bot Telegram - Direito Lux Staging
# Este script valida e configura o bot do Telegram após criação via BotFather

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Função para log colorido
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_step() {
    echo -e "${PURPLE}📋 $1${NC}"
}

# Banner
show_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════╗"
    echo "║     🤖 CONFIGURAR BOT TELEGRAM - DIREITO LUX      ║"
    echo "╚═══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Solicitar token
get_token() {
    echo ""
    log_step "Por favor, cole o token do bot obtido via @BotFather:"
    echo -e "${YELLOW}Exemplo: 7458394857:AAHKz9XjB8vK_2QxYz0-fG8kNvM_xQz7890${NC}"
    echo ""
    read -p "Token: " BOT_TOKEN
    
    # Validar formato básico
    if [[ ! "$BOT_TOKEN" =~ ^[0-9]+:[A-Za-z0-9_-]+$ ]]; then
        log_error "Formato de token inválido!"
        log_warning "O token deve ter o formato: NÚMEROS:LETRAS_E_NÚMEROS"
        exit 1
    fi
    
    log_success "Token recebido com formato válido"
}

# Testar token com API do Telegram
test_token() {
    log_step "Testando token com API do Telegram..."
    
    response=$(curl -s "https://api.telegram.org/bot${BOT_TOKEN}/getMe")
    
    if echo "$response" | grep -q '"ok":true'; then
        log_success "Token válido! Bot encontrado na API do Telegram"
        
        # Extrair informações do bot
        bot_username=$(echo "$response" | grep -o '"username":"[^"]*"' | cut -d'"' -f4)
        bot_name=$(echo "$response" | grep -o '"first_name":"[^"]*"' | cut -d'"' -f4)
        bot_id=$(echo "$response" | grep -o '"id":[0-9]*' | cut -d':' -f2 | head -1)
        
        echo ""
        log_info "📱 Informações do Bot:"
        echo -e "   ${CYAN}Nome:${NC} $bot_name"
        echo -e "   ${CYAN}Username:${NC} @$bot_username"
        echo -e "   ${CYAN}ID:${NC} $bot_id"
        echo -e "   ${CYAN}Link:${NC} https://t.me/$bot_username"
        echo ""
        
        return 0
    else
        log_error "Token inválido ou erro na API!"
        echo "Resposta: $response"
        exit 1
    fi
}

# Configurar comandos do bot
configure_commands() {
    log_step "Configurando comandos do bot..."
    
    commands='[
        {"command": "start", "description": "🚀 Iniciar conversa com o bot"},
        {"command": "help", "description": "🆘 Ajuda e comandos disponíveis"},
        {"command": "status", "description": "📊 Status dos seus processos"},
        {"command": "agenda", "description": "📅 Agenda de prazos importantes"},
        {"command": "busca", "description": "🔍 Buscar processos jurídicos"},
        {"command": "relatorio", "description": "📈 Relatório rápido dos processos"},
        {"command": "configurar", "description": "⚙️ Configurações do bot"}
    ]'
    
    response=$(curl -s -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setMyCommands" \
        -H "Content-Type: application/json" \
        -d "{\"commands\": $commands}")
    
    if echo "$response" | grep -q '"ok":true'; then
        log_success "Comandos configurados com sucesso"
    else
        log_warning "Não foi possível configurar comandos automaticamente"
        log_info "Configure manualmente via @BotFather com /setcommands"
    fi
}

# Configurar webhook
configure_webhook() {
    log_step "Configurando webhook..."
    
    # Verificar se túnel está ativo
    WEBHOOK_URL="https://direito-lux-staging.loca.lt/webhook/telegram"
    
    log_info "Testando túnel HTTPS..."
    if curl -s -o /dev/null -w "%{http_code}" "$WEBHOOK_URL" | grep -q "404"; then
        log_success "Túnel HTTPS está funcionando"
        
        # Configurar webhook
        response=$(curl -s -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook" \
            -H "Content-Type: application/json" \
            -d "{\"url\": \"$WEBHOOK_URL\", \"allowed_updates\": [\"message\", \"callback_query\"]}")
        
        if echo "$response" | grep -q '"ok":true'; then
            log_success "Webhook configurado com sucesso!"
            log_info "URL: $WEBHOOK_URL"
        else
            log_error "Falha ao configurar webhook"
            echo "Resposta: $response"
        fi
    else
        log_warning "Túnel HTTPS não está ativo"
        log_info "Execute primeiro: npx localtunnel --port 8085 --subdomain direito-lux-staging"
        WEBHOOK_CONFIGURED=false
    fi
}

# Atualizar arquivo .env
update_env_file() {
    log_step "Atualizando arquivo de configuração..."
    
    ENV_FILE="services/notification-service/.env"
    
    if [ -f "$ENV_FILE" ]; then
        # Fazer backup
        cp "$ENV_FILE" "${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
        
        # Atualizar token
        sed -i.bak "s|TELEGRAM_BOT_TOKEN=.*|TELEGRAM_BOT_TOKEN=$BOT_TOKEN|" "$ENV_FILE"
        
        log_success "Arquivo .env atualizado com sucesso"
        log_info "Backup salvo em: ${ENV_FILE}.backup.*"
    else
        log_error "Arquivo .env não encontrado em: $ENV_FILE"
        exit 1
    fi
}

# Reiniciar notification service
restart_service() {
    log_step "Reiniciando notification service..."
    
    if docker-compose ps | grep -q "direito-lux-notification.*Up"; then
        docker-compose restart notification-service
        log_success "Notification service reiniciado"
    else
        log_warning "Notification service não está rodando"
        log_info "Execute: docker-compose up -d notification-service"
    fi
}

# Testar bot enviando mensagem
test_bot_message() {
    if [ "$WEBHOOK_CONFIGURED" != "false" ]; then
        log_step "Testando envio de mensagem..."
        
        echo ""
        log_info "📱 TESTE MANUAL DO BOT:"
        echo "1. Abra o Telegram"
        echo "2. Acesse: https://t.me/$bot_username"
        echo "3. Clique em 'START'"
        echo "4. Envie: /help"
        echo ""
        log_warning "O bot deve responder se tudo estiver configurado corretamente!"
    fi
}

# Mostrar resumo final
show_summary() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}🎉 CONFIGURAÇÃO CONCLUÍDA!${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo ""
    
    log_success "✅ Bot Telegram configurado com sucesso!"
    echo ""
    echo -e "   ${CYAN}Bot:${NC} @$bot_username"
    echo -e "   ${CYAN}Link:${NC} https://t.me/$bot_username"
    echo -e "   ${CYAN}Token:${NC} Salvo em .env"
    
    if [ "$WEBHOOK_CONFIGURED" != "false" ]; then
        echo -e "   ${CYAN}Webhook:${NC} Configurado e ativo"
    else
        echo -e "   ${YELLOW}Webhook:${NC} Pendente (ative o túnel primeiro)"
    fi
    
    echo ""
    echo -e "${CYAN}📋 PRÓXIMOS PASSOS:${NC}"
    
    if [ "$WEBHOOK_CONFIGURED" = "false" ]; then
        echo "1. Ative o túnel HTTPS:"
        echo "   npx localtunnel --port 8085 --subdomain direito-lux-staging"
        echo ""
        echo "2. Execute novamente este script para configurar webhook"
        echo ""
    fi
    
    echo "3. Teste o bot no Telegram:"
    echo "   - Envie /start"
    echo "   - Envie /help"
    echo "   - Verifique os logs: docker-compose logs -f notification-service"
    echo ""
}

# Main
main() {
    show_banner
    
    # Verificar se está no diretório correto
    if [ ! -f "docker-compose.yml" ]; then
        log_error "Execute este script na raiz do projeto Direito Lux"
        exit 1
    fi
    
    # Executar passos
    get_token
    test_token
    configure_commands
    update_env_file
    restart_service
    configure_webhook
    test_bot_message
    show_summary
    
    # Salvar informações do bot
    cat > telegram_bot_info.txt << EOF
Bot Telegram - Direito Lux Staging
==================================
Nome: $bot_name
Username: @$bot_username
ID: $bot_id
Link: https://t.me/$bot_username
Token: [SALVO NO .ENV]
Webhook: $WEBHOOK_URL
Configurado em: $(date)
EOF
    
    log_info "Informações salvas em: telegram_bot_info.txt"
}

# Executar
WEBHOOK_CONFIGURED=true
main "$@"