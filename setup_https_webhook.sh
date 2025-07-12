#!/bin/bash

# 🌐 Setup HTTPS Webhook para Direito Lux Staging
# Cria túnel HTTPS para receber webhooks do Telegram e WhatsApp

set -e

echo "🌐 Configurando HTTPS Webhook para Direito Lux Staging"
echo "======================================================="

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Verificar se notification service está rodando
check_notification_service() {
    log_info "Verificando notification service..."
    
    if curl -s http://localhost:8085/health > /dev/null 2>&1; then
        log_success "Notification service está rodando na porta 8085"
        return 0
    else
        log_error "Notification service não está acessível"
        log_warning "Execute: docker-compose up -d notification-service"
        return 1
    fi
}

# Instalar localtunnel se não estiver instalado
install_localtunnel() {
    if ! command -v lt &> /dev/null; then
        log_info "Instalando localtunnel..."
        npm install -g localtunnel
        log_success "localtunnel instalado com sucesso"
    else
        log_success "localtunnel já está instalado"
    fi
}

# Iniciar túnel HTTPS
start_tunnel() {
    log_info "Iniciando túnel HTTPS..."
    
    # Matar processo existente se houver
    pkill -f "lt --port 8085" || true
    
    # Iniciar túnel em background
    lt --port 8085 --subdomain direito-lux-staging > tunnel_output.log 2>&1 &
    TUNNEL_PID=$!
    
    # Aguardar túnel inicializar
    log_info "Aguardando túnel inicializar..."
    sleep 5
    
    # Verificar se túnel está funcionando
    if kill -0 $TUNNEL_PID 2>/dev/null; then
        log_success "Túnel iniciado com PID: $TUNNEL_PID"
        
        # Tentar extrair URL do log
        if [ -f tunnel_output.log ]; then
            TUNNEL_URL=$(grep -o 'https://[^ ]*' tunnel_output.log | head -1)
            if [ -n "$TUNNEL_URL" ]; then
                log_success "URL do túnel: $TUNNEL_URL"
                echo "$TUNNEL_URL" > tunnel_url.txt
                return 0
            fi
        fi
        
        # Fallback para URL padrão
        TUNNEL_URL="https://direito-lux-staging.loca.lt"
        log_warning "Usando URL padrão: $TUNNEL_URL"
        echo "$TUNNEL_URL" > tunnel_url.txt
        return 0
    else
        log_error "Falha ao iniciar túnel"
        return 1
    fi
}

# Testar túnel
test_tunnel() {
    if [ -f tunnel_url.txt ]; then
        TUNNEL_URL=$(cat tunnel_url.txt)
        log_info "Testando túnel: $TUNNEL_URL"
        
        # Testar health endpoint
        if curl -s -L "$TUNNEL_URL/health" > /dev/null 2>&1; then
            log_success "Túnel está funcionando corretamente"
            return 0
        else
            log_error "Túnel não está respondendo"
            return 1
        fi
    else
        log_error "URL do túnel não encontrada"
        return 1
    fi
}

# Configurar webhook Telegram
configure_telegram_webhook() {
    if [ -f tunnel_url.txt ]; then
        TUNNEL_URL=$(cat tunnel_url.txt)
        WEBHOOK_URL="$TUNNEL_URL/webhook/telegram"
        
        log_info "Configurando webhook Telegram: $WEBHOOK_URL"
        
        # Nota: Precisa do token real do bot
        log_warning "Para configurar webhook Telegram, execute:"
        echo "curl -X POST \"https://api.telegram.org/bot<TOKEN>/setWebhook\" \\"
        echo "  -H \"Content-Type: application/json\" \\"
        echo "  -d '{\"url\": \"$WEBHOOK_URL\"}'"
        
        return 0
    else
        log_error "URL do túnel não encontrada"
        return 1
    fi
}

# Mostrar informações finais
show_summary() {
    if [ -f tunnel_url.txt ]; then
        TUNNEL_URL=$(cat tunnel_url.txt)
        
        echo ""
        echo "🎉 CONFIGURAÇÃO CONCLUÍDA!"
        echo "=========================="
        echo ""
        echo "📡 Túnel HTTPS: $TUNNEL_URL"
        echo "🔗 Health Check: $TUNNEL_URL/health"
        echo "📱 Webhook Telegram: $TUNNEL_URL/webhook/telegram"
        echo "📲 Webhook WhatsApp: $TUNNEL_URL/webhook/whatsapp"
        echo ""
        echo "📋 Próximos passos:"
        echo "1. Configurar webhook Telegram com token real"
        echo "2. Configurar webhook WhatsApp Business API"
        echo "3. Testar envio de mensagens"
        echo ""
        echo "🛑 Para parar o túnel: pkill -f 'lt --port 8085'"
        echo "📄 Logs do túnel: cat tunnel_output.log"
        echo ""
        echo "✅ Pronto para receber webhooks!"
    fi
}

# Função de cleanup
cleanup() {
    log_info "Limpando recursos..."
    pkill -f "lt --port 8085" || true
    rm -f tunnel_output.log tunnel_url.txt
}

# Trap para cleanup no exit
trap cleanup EXIT

# Executar setup
main() {
    # Verificar notification service
    if ! check_notification_service; then
        log_error "Notification service deve estar rodando primeiro"
        exit 1
    fi
    
    # Instalar localtunnel
    install_localtunnel
    
    # Iniciar túnel
    if ! start_tunnel; then
        log_error "Falha ao iniciar túnel"
        exit 1
    fi
    
    # Testar túnel
    if ! test_tunnel; then
        log_error "Falha no teste do túnel"
        exit 1
    fi
    
    # Configurar webhooks
    configure_telegram_webhook
    
    # Mostrar resumo
    show_summary
    
    # Manter script rodando
    log_info "Túnel rodando... Pressione Ctrl+C para parar"
    
    # Aguardar interrupção
    while true; do
        sleep 30
        if [ -f tunnel_url.txt ]; then
            TUNNEL_URL=$(cat tunnel_url.txt)
            if ! curl -s -L "$TUNNEL_URL/health" > /dev/null 2>&1; then
                log_warning "Túnel parece estar inativo, tentando reconectar..."
                start_tunnel
            fi
        fi
    done
}

# Executar main
main "$@"