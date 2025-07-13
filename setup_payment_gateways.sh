#!/bin/bash

# 💰 Setup Gateways de Pagamento - Direito Lux
# ASAAS (PIX + Cartão) + NOWPayments (Crypto)

set -e

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m'

# Banner
show_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════╗"
    echo "║    💰 GATEWAYS DE PAGAMENTO - DIREITO LUX         ║"
    echo "║    ASAAS (PIX+Cartão) + NOWPayments (Crypto)     ║"
    echo "╚═══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }
log_step() { echo -e "${PURPLE}📋 $1${NC}"; }

# Verificar pré-requisitos
check_prerequisites() {
    log_step "Verificando pré-requisitos..."
    
    # Verificar se está no diretório correto
    if [ ! -f "docker-compose.yml" ]; then
        log_error "Execute na raiz do projeto Direito Lux"
        exit 1
    fi
    
    # Verificar se billing service existe
    if [ ! -d "services/billing-service" ]; then
        log_error "Billing service não encontrado"
        exit 1
    fi
    
    log_success "Pré-requisitos OK"
}

# Mostrar informações dos gateways
show_gateway_info() {
    echo ""
    log_step "Informações dos Gateways de Pagamento"
    echo ""
    
    echo -e "${CYAN}🏦 ASAAS (PIX + Cartão + Boleto):${NC}"
    echo "• Site: https://www.asaas.com/"
    echo "• PIX: Taxa 0,99%"
    echo "• Cartão: Taxa 4,99% + R\$0,39"
    echo "• Boleto: Taxa R\$3,49"
    echo "• Aprovação: 24-48h"
    echo ""
    
    echo -e "${CYAN}₿ NOWPayments (Crypto):${NC}"
    echo "• Site: https://nowpayments.io/"
    echo "• Bitcoin, Ethereum, USDT/USDC"
    echo "• Taxa: 0.5-1.5%"
    echo "• +150 criptomoedas"
    echo "• Aprovação: 2-5 dias"
    echo ""
}

# Configurar ASAAS
setup_asaas() {
    log_step "Configurando ASAAS..."
    echo ""
    
    echo -e "${YELLOW}📋 Para configurar ASAAS:${NC}"
    echo "1. Acesse: https://www.asaas.com/"
    echo "2. Cadastro → Pessoa Jurídica (ou Física inicial)"
    echo "3. Dados: Direito Lux Tecnologia"
    echo "4. Aguarde aprovação (24-48h)"
    echo "5. Dashboard → Integrações → API"
    echo "6. Gerar API Key (Sandbox + Produção)"
    echo ""
    
    read -p "Já tem conta ASAAS? (y/n): " has_asaas
    
    if [ "$has_asaas" = "y" ]; then
        echo ""
        echo -e "${YELLOW}🔑 Digite suas API Keys ASAAS:${NC}"
        
        read -p "Sandbox API Key: " ASAAS_SANDBOX_KEY
        read -p "Production API Key (se tiver): " ASAAS_PROD_KEY
        
        # Validar formato
        if [[ ! "$ASAAS_SANDBOX_KEY" =~ ^\$aact_[A-Za-z0-9]+$ ]]; then
            log_warning "Formato de API Key ASAAS pode estar incorreto"
            log_info "Deve começar com '\$aact_'"
        fi
        
        # Configurar no GitHub Secrets
        if command -v gh &> /dev/null && gh auth status &> /dev/null 2>&1; then
            gh secret set ASAAS_API_KEY_SANDBOX --body "$ASAAS_SANDBOX_KEY"
            if [ ! -z "$ASAAS_PROD_KEY" ]; then
                gh secret set ASAAS_API_KEY_PRODUCTION --body "$ASAAS_PROD_KEY"
            fi
            log_success "ASAAS API Keys configuradas no GitHub Secrets"
        else
            log_warning "Configure manualmente no GitHub: ASAAS_API_KEY_SANDBOX"
        fi
        
        # Testar API
        test_asaas_api "$ASAAS_SANDBOX_KEY"
        
    else
        echo ""
        log_info "Abra uma nova aba e configure ASAAS:"
        echo "https://www.asaas.com/"
        echo ""
        log_warning "Após criar a conta, execute este script novamente"
    fi
}

# Testar API ASAAS
test_asaas_api() {
    local api_key="$1"
    log_step "Testando API ASAAS..."
    
    response=$(curl -s -X GET \
        "https://sandbox.asaas.com/api/v3/customers?limit=1" \
        -H "access_token: $api_key")
    
    if echo "$response" | grep -q '"object":"list"'; then
        log_success "API ASAAS funcionando!"
    else
        log_warning "Erro na API ASAAS: $response"
    fi
}

# Configurar NOWPayments
setup_nowpayments() {
    log_step "Configurando NOWPayments..."
    echo ""
    
    echo -e "${YELLOW}📋 Para configurar NOWPayments:${NC}"
    echo "1. Acesse: https://nowpayments.io/"
    echo "2. Sign Up → Business Account"
    echo "3. Complete KYC (documentos)"
    echo "4. Aguarde aprovação (2-5 dias)"
    echo "5. Dashboard → Settings → API"
    echo "6. Generate API Key"
    echo ""
    
    read -p "Já tem conta NOWPayments? (y/n): " has_nowpayments
    
    if [ "$has_nowpayments" = "y" ]; then
        echo ""
        echo -e "${YELLOW}🔑 Digite suas API Keys NOWPayments:${NC}"
        
        read -p "Sandbox API Key: " NOWPAYMENTS_SANDBOX_KEY
        read -p "Production API Key (se tiver): " NOWPAYMENTS_PROD_KEY
        
        # Validar formato
        if [[ ! "$NOWPAYMENTS_SANDBOX_KEY" =~ ^NP-[A-Za-z0-9]+$ ]]; then
            log_warning "Formato de API Key NOWPayments pode estar incorreto"
            log_info "Deve começar com 'NP-'"
        fi
        
        # Configurar no GitHub Secrets
        if command -v gh &> /dev/null && gh auth status &> /dev/null 2>&1; then
            gh secret set NOWPAYMENTS_API_KEY_SANDBOX --body "$NOWPAYMENTS_SANDBOX_KEY"
            if [ ! -z "$NOWPAYMENTS_PROD_KEY" ]; then
                gh secret set NOWPAYMENTS_API_KEY_PRODUCTION --body "$NOWPAYMENTS_PROD_KEY"
            fi
            log_success "NOWPayments API Keys configuradas no GitHub Secrets"
        else
            log_warning "Configure manualmente no GitHub: NOWPAYMENTS_API_KEY_SANDBOX"
        fi
        
        # Testar API
        test_nowpayments_api "$NOWPAYMENTS_SANDBOX_KEY"
        
    else
        echo ""
        log_info "Abra uma nova aba e configure NOWPayments:"
        echo "https://nowpayments.io/"
        echo ""
        log_warning "Após criar a conta, execute este script novamente"
    fi
}

# Testar API NOWPayments
test_nowpayments_api() {
    local api_key="$1"
    log_step "Testando API NOWPayments..."
    
    response=$(curl -s -X GET \
        "https://api-sandbox.nowpayments.io/v1/status" \
        -H "x-api-key: $api_key")
    
    if echo "$response" | grep -q '"message":"OK"'; then
        log_success "API NOWPayments funcionando!"
    else
        log_warning "Erro na API NOWPayments: $response"
    fi
}

# Configurar webhooks
setup_webhooks() {
    log_step "Configurando webhooks..."
    echo ""
    
    # Verificar se túnel está ativo
    if curl -s https://direito-lux-staging.loca.lt/health > /dev/null; then
        WEBHOOK_BASE="https://direito-lux-staging.loca.lt"
        log_success "Túnel HTTPS ativo: $WEBHOOK_BASE"
    else
        log_error "Túnel HTTPS não está ativo!"
        log_warning "Execute: npx localtunnel --port 8085 --subdomain direito-lux-staging"
        exit 1
    fi
    
    echo ""
    echo -e "${CYAN}🌐 Configure os webhooks nos dashboards:${NC}"
    echo ""
    
    echo -e "${YELLOW}ASAAS Webhook:${NC}"
    echo "URL: $WEBHOOK_BASE/webhook/asaas"
    echo "Events: payment_created, payment_confirmed, payment_overdue"
    echo ""
    
    echo -e "${YELLOW}NOWPayments Webhook:${NC}"
    echo "URL: $WEBHOOK_BASE/webhook/nowpayments"
    echo "Events: payment_waiting, payment_confirmed, payment_failed"
    echo ""
    
    # Testar endpoints
    log_step "Testando endpoints de webhook..."
    
    # Teste ASAAS
    if curl -s "$WEBHOOK_BASE/webhook/asaas" > /dev/null; then
        log_success "Endpoint ASAAS acessível"
    else
        log_warning "Endpoint ASAAS pode precisar ser implementado"
    fi
    
    # Teste NOWPayments
    if curl -s "$WEBHOOK_BASE/webhook/nowpayments" > /dev/null; then
        log_success "Endpoint NOWPayments acessível"
    else
        log_warning "Endpoint NOWPayments pode precisar ser implementado"
    fi
}

# Configurar billing service
setup_billing_service() {
    log_step "Configurando Billing Service..."
    
    # Verificar se está rodando
    if docker-compose ps billing-service | grep -q "Up"; then
        log_success "Billing Service já está rodando"
        
        # Restart para carregar novas variáveis
        docker-compose restart billing-service
        sleep 5
        
        if curl -s http://localhost:8087/health > /dev/null; then
            log_success "Billing Service funcionando"
        else
            log_warning "Billing Service pode estar inicializando..."
        fi
    else
        log_warning "Billing Service não está rodando"
        log_info "Execute: docker-compose up -d billing-service"
    fi
}

# Criar planos de teste
create_test_plans() {
    log_step "Criando planos de teste..."
    echo ""
    
    echo -e "${CYAN}💰 Planos Direito Lux:${NC}"
    echo ""
    echo "• Starter: R\$ 99,00/mês"
    echo "• Professional: R\$ 299,00/mês" 
    echo "• Business: R\$ 699,00/mês"
    echo "• Enterprise: R\$ 1.999,00/mês"
    echo ""
    
    if [ ! -z "$ASAAS_SANDBOX_KEY" ]; then
        log_info "Criando planos no ASAAS..."
        
        # Criar cliente de teste
        customer_response=$(curl -s -X POST \
            "https://sandbox.asaas.com/api/v3/customers" \
            -H "access_token: $ASAAS_SANDBOX_KEY" \
            -H "Content-Type: application/json" \
            -d '{
                "name": "Cliente Teste Direito Lux",
                "email": "teste@direitolux.com.br",
                "cpfCnpj": "12345678901"
            }')
        
        customer_id=$(echo "$customer_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        
        if [ ! -z "$customer_id" ]; then
            log_success "Cliente teste criado: $customer_id"
            
            # Criar cobrança de teste
            payment_response=$(curl -s -X POST \
                "https://sandbox.asaas.com/api/v3/payments" \
                -H "access_token: $ASAAS_SANDBOX_KEY" \
                -H "Content-Type: application/json" \
                -d "{
                    \"customer\": \"$customer_id\",
                    \"billingType\": \"PIX\",
                    \"dueDate\": \"$(date -d '+1 day' '+%Y-%m-%d')\",
                    \"value\": 99.00,
                    \"description\": \"Direito Lux - Plano Starter (Teste)\"
                }")
            
            payment_id=$(echo "$payment_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
            
            if [ ! -z "$payment_id" ]; then
                log_success "Cobrança teste criada: $payment_id"
            fi
        fi
    fi
    
    echo ""
    log_info "Para testes em produção, use dados reais dos clientes"
}

# Mostrar resumo final
show_summary() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}🎉 GATEWAYS DE PAGAMENTO CONFIGURADOS!${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo ""
    
    echo -e "${YELLOW}📊 STATUS:${NC}"
    if [ ! -z "$ASAAS_SANDBOX_KEY" ]; then
        echo "✅ ASAAS: Configurado"
    else
        echo "⏸️ ASAAS: Pendente"
    fi
    
    if [ ! -z "$NOWPAYMENTS_SANDBOX_KEY" ]; then
        echo "✅ NOWPayments: Configurado" 
    else
        echo "⏸️ NOWPayments: Pendente"
    fi
    
    echo ""
    echo -e "${BLUE}🔗 WEBHOOKS:${NC}"
    echo "• ASAAS: https://direito-lux-staging.loca.lt/webhook/asaas"
    echo "• NOWPayments: https://direito-lux-staging.loca.lt/webhook/nowpayments"
    echo ""
    
    echo -e "${YELLOW}📋 PRÓXIMOS PASSOS:${NC}"
    echo "1. Complete aprovação das contas nos gateways"
    echo "2. Configure webhooks nos dashboards"
    echo "3. Teste pagamentos em sandbox"
    echo "4. Implemente billing service completo"
    echo "5. Deploy em produção"
    echo ""
    
    echo -e "${CYAN}📚 DOCUMENTAÇÃO:${NC}"
    echo "• CONFIGURAR_GATEWAYS_PAGAMENTO.md"
    echo "• services/billing-service/README.md"
    echo ""
    
    echo -e "${GREEN}💰 SISTEMA PRONTO PARA MONETIZAÇÃO!${NC}"
}

# Main
main() {
    show_banner
    check_prerequisites
    show_gateway_info
    
    echo ""
    log_step "Configurando gateways de pagamento..."
    echo ""
    
    # Menu de opções
    echo -e "${YELLOW}Escolha uma opção:${NC}"
    echo "1. Configurar ASAAS (PIX + Cartão)"
    echo "2. Configurar NOWPayments (Crypto)"
    echo "3. Configurar ambos"
    echo "4. Apenas configurar webhooks"
    echo ""
    read -p "Opção (1-4): " choice
    
    case $choice in
        1)
            setup_asaas
            setup_webhooks
            setup_billing_service
            ;;
        2)
            setup_nowpayments
            setup_webhooks
            setup_billing_service
            ;;
        3)
            setup_asaas
            setup_nowpayments
            setup_webhooks
            setup_billing_service
            create_test_plans
            ;;
        4)
            setup_webhooks
            ;;
        *)
            log_error "Opção inválida"
            exit 1
            ;;
    esac
    
    show_summary
}

# Executar se chamado diretamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi