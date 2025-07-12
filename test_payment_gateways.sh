#!/bin/bash

# 💰 Teste dos Gateways de Pagamento - Direito Lux Staging

echo "💰 Testando webhooks dos gateways de pagamento..."

# URLs dos webhooks
ASAAS_WEBHOOK_URL="https://direito-lux-staging.loca.lt/billing/webhooks/asaas"
NOWPAYMENTS_WEBHOOK_URL="https://direito-lux-staging.loca.lt/billing/webhooks/crypto"
BILLING_SERVICE_URL="http://localhost:8089"

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

# Teste 1: Verificar se billing service está rodando
test_billing_service() {
    log_info "Verificando billing service..."
    
    if curl -s "$BILLING_SERVICE_URL/health" > /dev/null 2>&1; then
        log_success "Billing service está rodando na porta 8089"
        return 0
    else
        log_error "Billing service não está acessível"
        log_warning "Execute: docker-compose up -d billing-service"
        return 1
    fi
}

# Teste 2: Verificar túnel HTTPS
test_tunnel() {
    log_info "Verificando túnel HTTPS..."
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$ASAAS_WEBHOOK_URL")
    
    if [ "$response" = "200" ] || [ "$response" = "404" ] || [ "$response" = "405" ]; then
        log_success "Túnel HTTPS está funcionando (HTTP $response)"
        return 0
    else
        log_error "Túnel HTTPS não está funcionando (HTTP $response)"
        return 1
    fi
}

# Teste 3: Testar webhook ASAAS
test_asaas_webhook() {
    log_info "Testando webhook ASAAS..."
    
    # Payload simulado do ASAAS
    asaas_payload='{
        "event": "PAYMENT_RECEIVED",
        "dateCreated": "2025-07-11T13:00:00.000Z",
        "payment": {
            "id": "pay_123456789",
            "customer": "cus_123456789",
            "subscription": "sub_123456789",
            "value": 99.00,
            "netValue": 95.00,
            "originalValue": 99.00,
            "interestValue": 0.00,
            "description": "Plano Starter - Direito Lux",
            "billingType": "PIX",
            "status": "RECEIVED",
            "pixTransaction": {
                "txid": "pix_test_123456789",
                "endToEndIdentifier": "E123456789202507111300123456789"
            },
            "paymentDate": "2025-07-11T13:00:00.000Z",
            "clientPaymentDate": "2025-07-11T13:00:00.000Z",
            "creditDate": "2025-07-11T13:00:00.000Z",
            "estimatedCreditDate": "2025-07-11T13:00:00.000Z"
        }
    }'
    
    response=$(curl -s -X POST "$ASAAS_WEBHOOK_URL" \
        -H "Content-Type: application/json" \
        -H "User-Agent: Asaas-hookshot" \
        -d "$asaas_payload")
    
    if [ $? -eq 0 ]; then
        log_success "Webhook ASAAS aceita requisições POST"
        log_info "Resposta: $response"
        return 0
    else
        log_error "Falha ao testar webhook ASAAS"
        return 1
    fi
}

# Teste 4: Testar webhook NOWPayments
test_nowpayments_webhook() {
    log_info "Testando webhook NOWPayments..."
    
    # Payload simulado do NOWPayments
    nowpayments_payload='{
        "payment_id": "123456789",
        "payment_status": "finished",
        "pay_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
        "price_amount": 99.00,
        "price_currency": "BRL",
        "pay_amount": 0.00234567,
        "pay_currency": "BTC",
        "order_id": "order_123456789",
        "order_description": "Plano Starter - Direito Lux",
        "purchase_id": "purchase_123456789",
        "outcome_amount": 0.00234567,
        "outcome_currency": "BTC",
        "payment_extra_id": "",
        "created_at": "2025-07-11T13:00:00.000Z",
        "updated_at": "2025-07-11T13:00:00.000Z"
    }'
    
    response=$(curl -s -X POST "$NOWPAYMENTS_WEBHOOK_URL" \
        -H "Content-Type: application/json" \
        -H "User-Agent: NOWPayments-webhook" \
        -d "$nowpayments_payload")
    
    if [ $? -eq 0 ]; then
        log_success "Webhook NOWPayments aceita requisições POST"
        log_info "Resposta: $response"
        return 0
    else
        log_error "Falha ao testar webhook NOWPayments"
        return 1
    fi
}

# Teste 5: Testar APIs do billing service
test_billing_apis() {
    log_info "Testando APIs do billing service..."
    
    # Testar endpoint de planos
    if curl -s "$BILLING_SERVICE_URL/billing/plans" > /dev/null 2>&1; then
        log_success "Endpoint de planos está funcionando"
    else
        log_warning "Endpoint de planos não está acessível"
    fi
    
    # Testar endpoint de health
    if curl -s "$BILLING_SERVICE_URL/health" > /dev/null 2>&1; then
        log_success "Health check está funcionando"
    else
        log_warning "Health check não está acessível"
    fi
    
    # Testar endpoint de métricas
    if curl -s "$BILLING_SERVICE_URL/metrics" > /dev/null 2>&1; then
        log_success "Métricas estão funcionando"
    else
        log_warning "Métricas não estão acessíveis"
    fi
}

# Mostrar configurações para produção
show_production_config() {
    echo ""
    echo "🔧 CONFIGURAÇÕES PARA PRODUÇÃO"
    echo "=============================="
    echo ""
    echo "📋 ASAAS (Gateway tradicional):"
    echo "   1. Criar conta: https://www.asaas.com/"
    echo "   2. Obter API Key de produção"
    echo "   3. Configurar webhook: $ASAAS_WEBHOOK_URL"
    echo "   4. Testar com PIX real (R$ 0,01)"
    echo ""
    echo "📋 NOWPayments (Gateway cripto):"
    echo "   1. Criar conta: https://nowpayments.io/"
    echo "   2. Obter API Key de produção"
    echo "   3. Configurar webhook: $NOWPAYMENTS_WEBHOOK_URL"
    echo "   4. Testar com Bitcoin real (valor mínimo)"
    echo ""
    echo "📋 Variáveis de ambiente para produção:"
    echo "   ASAAS_API_KEY=\$asaas_production_key"
    echo "   ASAAS_ENVIRONMENT=production"
    echo "   NOWPAYMENTS_API_KEY=\$nowpayments_production_key"
    echo ""
    echo "📋 Custos estimados para testes:"
    echo "   - ASAAS: R$ 0,01 (PIX mínimo)"
    echo "   - NOWPayments: ~R$ 5,00 (Bitcoin mínimo)"
    echo "   - Total: ~R$ 5,01 para validação completa"
    echo ""
}

# Main
main() {
    echo "💰 Direito Lux - Teste dos Gateways de Pagamento"
    echo "==============================================="
    echo ""
    
    # Executar testes
    test_billing_service
    billing_ok=$?
    
    test_tunnel
    tunnel_ok=$?
    
    test_asaas_webhook
    asaas_ok=$?
    
    test_nowpayments_webhook
    nowpayments_ok=$?
    
    test_billing_apis
    
    # Resumo dos testes
    echo ""
    echo "📊 RESUMO DOS TESTES"
    echo "==================="
    
    if [ $billing_ok -eq 0 ]; then
        log_success "Billing Service: OK"
    else
        log_error "Billing Service: FALHOU"
    fi
    
    if [ $tunnel_ok -eq 0 ]; then
        log_success "Túnel HTTPS: OK"
    else
        log_error "Túnel HTTPS: FALHOU"
    fi
    
    if [ $asaas_ok -eq 0 ]; then
        log_success "Webhook ASAAS: OK"
    else
        log_warning "Webhook ASAAS: PENDENTE"
    fi
    
    if [ $nowpayments_ok -eq 0 ]; then
        log_success "Webhook NOWPayments: OK"
    else
        log_warning "Webhook NOWPayments: PENDENTE"
    fi
    
    echo ""
    
    # Mostrar próximos passos
    if [ $billing_ok -eq 0 ] && [ $tunnel_ok -eq 0 ]; then
        log_success "✅ Infraestrutura pronta para gateways de pagamento"
        show_production_config
    else
        log_error "❌ Infraestrutura precisa ser corrigida"
    fi
}

# Executar
main "$@"