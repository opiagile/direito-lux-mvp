#!/bin/bash

# 📱 Configurar WhatsApp Business API - Direito Lux Staging

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
    echo "║    📱 CONFIGURAR WHATSAPP API - DIREITO LUX       ║"
    echo "╚═══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Solicitar credenciais
get_credentials() {
    echo ""
    log_step "Cole as credenciais obtidas no Meta for Developers:"
    echo ""
    
    echo -e "${YELLOW}1. Access Token (temporário - 24h):${NC}"
    echo "Exemplo: EAAxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    read -p "Access Token: " ACCESS_TOKEN
    
    echo ""
    echo -e "${YELLOW}2. Phone Number ID:${NC}"
    echo "Exemplo: 123456789012345"
    read -p "Phone Number ID: " PHONE_NUMBER_ID
    
    echo ""
    echo -e "${YELLOW}3. Business Account ID (opcional):${NC}"
    echo "Exemplo: 123456789012345"
    read -p "Business Account ID: " BUSINESS_ACCOUNT_ID
    
    # Validar formato básico
    if [[ ! "$ACCESS_TOKEN" =~ ^EAA[A-Za-z0-9_-]+$ ]]; then
        log_error "Formato de Access Token inválido!"
        log_warning "Deve começar com 'EAA' seguido de letras/números"
        exit 1
    fi
    
    if [[ ! "$PHONE_NUMBER_ID" =~ ^[0-9]+$ ]]; then
        log_error "Phone Number ID deve conter apenas números!"
        exit 1
    fi
    
    log_success "Credenciais recebidas com formato válido"
}

# Testar credenciais
test_credentials() {
    log_step "Testando credenciais com API do WhatsApp..."
    
    # Testar se pode acessar informações do número
    response=$(curl -s -X GET \
        "https://graph.facebook.com/v18.0/${PHONE_NUMBER_ID}" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}")
    
    if echo "$response" | grep -q '"verified_name"'; then
        log_success "Credenciais válidas! Acesso à API confirmado"
        
        # Extrair informações
        verified_name=$(echo "$response" | grep -o '"verified_name":"[^"]*"' | cut -d'"' -f4)
        display_number=$(echo "$response" | grep -o '"display_phone_number":"[^"]*"' | cut -d'"' -f4)
        
        echo ""
        log_info "📱 Informações do WhatsApp Business:"
        echo -e "   ${CYAN}Nome verificado:${NC} $verified_name"
        echo -e "   ${CYAN}Número:${NC} $display_number"
        echo -e "   ${CYAN}Phone Number ID:${NC} $PHONE_NUMBER_ID"
        echo ""
        
        return 0
    else
        log_error "Credenciais inválidas ou erro na API!"
        echo "Resposta: $response"
        exit 1
    fi
}

# Testar webhook
test_webhook() {
    log_step "Testando webhook..."
    
    WEBHOOK_URL="https://direito-lux-staging.loca.lt/webhook/whatsapp"
    
    # Verificar se túnel está ativo
    if curl -s "$WEBHOOK_URL" > /dev/null 2>&1; then
        log_success "Webhook endpoint está acessível"
    else
        log_error "Webhook endpoint não está acessível"
        log_warning "Certifique-se de que o túnel está rodando:"
        echo "npx localtunnel --port 8085 --subdomain direito-lux-staging"
        exit 1
    fi
    
    # Testar verificação do webhook
    verify_token="direito_lux_staging_2025"
    challenge="test_challenge_12345"
    
    response=$(curl -s "${WEBHOOK_URL}?hub.mode=subscribe&hub.challenge=${challenge}&hub.verify_token=${verify_token}")
    
    if [ "$response" = "$challenge" ]; then
        log_success "Verificação do webhook funcionando corretamente"
    else
        log_warning "Verificação do webhook pode precisar de ajustes"
        log_info "Resposta recebida: $response"
        log_info "Esperado: $challenge"
    fi
}

# Atualizar arquivo .env
update_env_file() {
    log_step "Atualizando arquivo de configuração..."
    
    ENV_FILE="services/notification-service/.env"
    
    if [ -f "$ENV_FILE" ]; then
        # Fazer backup
        cp "$ENV_FILE" "${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
        
        # Atualizar credenciais
        sed -i.bak "s|WHATSAPP_ACCESS_TOKEN=.*|WHATSAPP_ACCESS_TOKEN=$ACCESS_TOKEN|" "$ENV_FILE"
        sed -i.bak "s|WHATSAPP_PHONE_NUMBER_ID=.*|WHATSAPP_PHONE_NUMBER_ID=$PHONE_NUMBER_ID|" "$ENV_FILE"
        
        if [ ! -z "$BUSINESS_ACCOUNT_ID" ]; then
            # Adicionar business account ID se não existir
            if ! grep -q "WHATSAPP_BUSINESS_ACCOUNT_ID" "$ENV_FILE"; then
                echo "WHATSAPP_BUSINESS_ACCOUNT_ID=$BUSINESS_ACCOUNT_ID" >> "$ENV_FILE"
            else
                sed -i.bak "s|WHATSAPP_BUSINESS_ACCOUNT_ID=.*|WHATSAPP_BUSINESS_ACCOUNT_ID=$BUSINESS_ACCOUNT_ID|" "$ENV_FILE"
            fi
        fi
        
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
        
        # Aguardar serviço inicializar
        sleep 5
        
        # Verificar se está funcionando
        if curl -s http://localhost:8085/health > /dev/null; then
            log_success "Notification service está funcionando"
        else
            log_warning "Notification service pode estar com problemas"
        fi
    else
        log_warning "Notification service não está rodando"
        log_info "Execute: docker-compose up -d notification-service"
    fi
}

# Teste de envio de mensagem
test_message_sending() {
    log_step "Preparando teste de envio de mensagem..."
    
    echo ""
    echo -e "${YELLOW}Para testar o envio de mensagens:${NC}"
    echo "1. Certifique-se de que seu número está na lista de teste do Meta"
    echo "2. Acesse: https://developers.facebook.com/"
    echo "3. Vá para seu app → WhatsApp → API Setup"
    echo "4. Na seção 'Send and receive messages', teste enviar uma mensagem"
    echo ""
    
    echo -e "${CYAN}Exemplo de teste via curl:${NC}"
    echo "curl -X POST 'https://graph.facebook.com/v18.0/${PHONE_NUMBER_ID}/messages' \\"
    echo "  -H 'Authorization: Bearer ${ACCESS_TOKEN}' \\"
    echo "  -H 'Content-Type: application/json' \\"
    echo "  -d '{"
    echo "    \"messaging_product\": \"whatsapp\","
    echo "    \"to\": \"SEU_NUMERO_AQUI\","
    echo "    \"type\": \"text\","
    echo "    \"text\": {"
    echo "      \"body\": \"Teste WhatsApp API - Direito Lux Staging\""
    echo "    }"
    echo "  }'"
    echo ""
}

# Mostrar resumo final
show_summary() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}🎉 WHATSAPP API CONFIGURADO!${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo ""
    
    log_success "✅ WhatsApp Business API configurado com sucesso!"
    echo ""
    echo -e "   ${CYAN}Phone Number ID:${NC} $PHONE_NUMBER_ID"
    echo -e "   ${CYAN}Access Token:${NC} Salvo em .env"
    echo -e "   ${CYAN}Webhook:${NC} Configurado e testado"
    echo -e "   ${CYAN}Service:${NC} Reiniciado"
    echo ""
    echo -e "${CYAN}📋 PRÓXIMOS PASSOS:${NC}"
    echo "1. Configure o webhook no Meta for Developers:"
    echo "   URL: https://direito-lux-staging.loca.lt/webhook/whatsapp"
    echo "   Token: direito_lux_staging_2025"
    echo ""
    echo "2. Teste o envio de mensagens via Meta console"
    echo ""
    echo "3. Monitore os logs:"
    echo "   docker-compose logs -f notification-service"
    echo ""
    echo -e "${YELLOW}⚠️ LEMBRE-SE:${NC}"
    echo "- Access Token temporário expira em 24h"
    echo "- Para produção, configure um token permanente"
    echo "- Limite: 100 mensagens/dia no teste"
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
    
    # Verificar se túnel está ativo
    if ! curl -s https://direito-lux-staging.loca.lt/health > /dev/null; then
        log_error "Túnel HTTPS não está ativo!"
        log_warning "Execute primeiro: npx localtunnel --port 8085 --subdomain direito-lux-staging"
        exit 1
    fi
    
    # Executar passos
    get_credentials
    test_credentials
    test_webhook
    update_env_file
    restart_service
    test_message_sending
    show_summary
    
    # Salvar informações
    cat > whatsapp_api_info.txt << EOF
WhatsApp Business API - Direito Lux Staging
==========================================
Phone Number ID: $PHONE_NUMBER_ID
Business Account ID: $BUSINESS_ACCOUNT_ID
Access Token: [SALVO NO .ENV]
Webhook URL: https://direito-lux-staging.loca.lt/webhook/whatsapp
Verify Token: direito_lux_staging_2025
Configurado em: $(date)

LIMITAÇÕES TESTE:
- Token temporário: 24h
- Mensagens: 100/dia
- Destinatários: 5 números
EOF
    
    log_info "Informações salvas em: whatsapp_api_info.txt"
}

# Executar
main "$@"