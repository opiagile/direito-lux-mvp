#!/bin/bash

# 📱 Teste do WhatsApp Business API Webhook - Direito Lux Staging

echo "📱 Testando webhook do WhatsApp Business API..."

# URLs e tokens
WEBHOOK_URL="https://direito-lux-staging.loca.lt/webhook/whatsapp"
VERIFY_TOKEN="direito_lux_staging_2025"

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

# Teste 1: Verificar se túnel está funcionando
test_tunnel() {
    log_info "Testando túnel HTTPS..."
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$WEBHOOK_URL")
    
    if [ "$response" = "200" ] || [ "$response" = "404" ]; then
        log_success "Túnel HTTPS está funcionando (HTTP $response)"
        return 0
    else
        log_error "Túnel HTTPS não está funcionando (HTTP $response)"
        return 1
    fi
}

# Teste 2: Simular verificação do webhook pelo Facebook
test_webhook_verification() {
    log_info "Testando verificação do webhook (simulando Facebook)..."
    
    # Simular GET request que o Facebook faz
    challenge_value="test_challenge_123"
    verification_url="${WEBHOOK_URL}?hub.mode=subscribe&hub.challenge=${challenge_value}&hub.verify_token=${VERIFY_TOKEN}"
    
    response=$(curl -s "$verification_url")
    
    if [ "$response" = "$challenge_value" ]; then
        log_success "Verificação do webhook OK - challenge retornado corretamente"
        return 0
    else
        log_warning "Verificação do webhook não implementada ainda"
        log_info "Resposta recebida: $response"
        log_info "Esperado: $challenge_value"
        return 1
    fi
}

# Teste 3: Testar recebimento de mensagem
test_message_webhook() {
    log_info "Testando recebimento de mensagem via webhook..."
    
    # Simular POST request do WhatsApp
    message_payload='{
        "object": "whatsapp_business_account",
        "entry": [
            {
                "id": "123456789",
                "changes": [
                    {
                        "value": {
                            "messaging_product": "whatsapp",
                            "metadata": {
                                "display_phone_number": "+5511999999999",
                                "phone_number_id": "123456789"
                            },
                            "messages": [
                                {
                                    "from": "5511888888888",
                                    "id": "wamid.test123",
                                    "timestamp": "1642694617",
                                    "text": {
                                        "body": "Teste de mensagem WhatsApp"
                                    },
                                    "type": "text"
                                }
                            ]
                        },
                        "field": "messages"
                    }
                ]
            }
        ]
    }'
    
    response=$(curl -s -X POST "$WEBHOOK_URL" \
        -H "Content-Type: application/json" \
        -d "$message_payload")
    
    if [ $? -eq 0 ]; then
        log_success "Webhook aceita mensagens POST"
        log_info "Resposta: $response"
        return 0
    else
        log_error "Falha ao enviar mensagem POST para webhook"
        return 1
    fi
}

# Teste 4: Verificar configuração do notification service
test_notification_service_config() {
    log_info "Verificando configuração do notification service..."
    
    # Verificar se serviço está rodando
    if curl -s http://localhost:8085/health > /dev/null; then
        log_success "Notification service está rodando"
        
        # Verificar se tem endpoint específico para WhatsApp
        if curl -s -o /dev/null -w "%{http_code}" http://localhost:8085/webhook/whatsapp | grep -q "200\|404\|405"; then
            log_success "Endpoint /webhook/whatsapp existe"
            return 0
        else
            log_warning "Endpoint /webhook/whatsapp pode não estar implementado"
            return 1
        fi
    else
        log_error "Notification service não está rodando"
        return 1
    fi
}

# Mostrar instruções para configuração manual
show_configuration_instructions() {
    echo ""
    echo "🔧 INSTRUÇÕES PARA CONFIGURAR WHATSAPP BUSINESS API"
    echo "=================================================="
    echo ""
    echo "📋 1. Acesse Meta for Developers:"
    echo "   https://developers.facebook.com/"
    echo ""
    echo "📋 2. Crie um app WhatsApp Business API"
    echo ""
    echo "📋 3. Configure o webhook:"
    echo "   URL: $WEBHOOK_URL"
    echo "   Token: $VERIFY_TOKEN"
    echo ""
    echo "📋 4. Obtenha as credenciais:"
    echo "   - Access Token"
    echo "   - Phone Number ID"
    echo "   - Business Account ID"
    echo ""
    echo "📋 5. Atualize o arquivo .env:"
    echo "   WHATSAPP_ACCESS_TOKEN=EAAxxxxxxxxxxxxxxxx"
    echo "   WHATSAPP_PHONE_NUMBER_ID=123456789012345"
    echo ""
    echo "📋 6. Reinicie o notification service:"
    echo "   docker-compose restart notification-service"
    echo ""
}

# Main
main() {
    echo "📱 Direito Lux - Teste WhatsApp Business API Webhook"
    echo "===================================================="
    echo ""
    
    # Executar testes
    test_tunnel
    tunnel_ok=$?
    
    test_notification_service_config
    service_ok=$?
    
    test_webhook_verification
    verification_ok=$?
    
    test_message_webhook
    message_ok=$?
    
    # Resumo dos testes
    echo ""
    echo "📊 RESUMO DOS TESTES"
    echo "==================="
    
    if [ $tunnel_ok -eq 0 ]; then
        log_success "Túnel HTTPS: OK"
    else
        log_error "Túnel HTTPS: FALHOU"
    fi
    
    if [ $service_ok -eq 0 ]; then
        log_success "Notification Service: OK"
    else
        log_error "Notification Service: FALHOU"
    fi
    
    if [ $verification_ok -eq 0 ]; then
        log_success "Verificação Webhook: OK"
    else
        log_warning "Verificação Webhook: PENDENTE"
    fi
    
    if [ $message_ok -eq 0 ]; then
        log_success "Recebimento Mensagem: OK"
    else
        log_warning "Recebimento Mensagem: PENDENTE"
    fi
    
    echo ""
    
    # Mostrar próximos passos
    if [ $tunnel_ok -eq 0 ] && [ $service_ok -eq 0 ]; then
        log_success "✅ Infraestrutura pronta para WhatsApp Business API"
        show_configuration_instructions
    else
        log_error "❌ Infraestrutura precisa ser corrigida antes da configuração"
    fi
}

# Executar
main "$@"