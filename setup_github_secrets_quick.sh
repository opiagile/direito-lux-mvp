#!/bin/bash

# 🔐 Setup Rápido GitHub Secrets - Direito Lux
# Configuração automática de todos os secrets necessários

set -e

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔═══════════════════════════════════════════════════╗"
echo "║    🔐 GITHUB SECRETS - SETUP RÁPIDO              ║"
echo "║    Direito Lux - Staging & Production Ready      ║"
echo "╚═══════════════════════════════════════════════════╝"
echo -e "${NC}"

# Verificar se está no diretório correto
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}❌ Execute na raiz do projeto Direito Lux${NC}"
    exit 1
fi

# Verificar GitHub CLI
if ! command -v gh &> /dev/null; then
    echo -e "${RED}❌ GitHub CLI não encontrado${NC}"
    echo "Instale: brew install gh"
    exit 1
fi

# Verificar autenticação
if ! gh auth status &> /dev/null; then
    echo -e "${YELLOW}🔑 Fazendo login no GitHub...${NC}"
    gh auth login --web
fi

echo -e "${GREEN}✅ GitHub CLI autenticado${NC}"
echo ""

# Coletar secrets
echo -e "${BLUE}📋 Configure os secrets necessários:${NC}"
echo ""

# Telegram (obrigatório)
echo -e "${YELLOW}1. 🤖 TELEGRAM BOT TOKEN (obrigatório):${NC}"
echo "Token atual: 7927061803:AAGC5GMerAe9CVegcl85o6BTFj2hqkcjO04"
read -p "Usar token atual? (y/n): " use_current_telegram
if [ "$use_current_telegram" = "y" ]; then
    TELEGRAM_TOKEN="7927061803:AAGC5GMerAe9CVegcl85o6BTFj2hqkcjO04"
else
    read -s -p "Digite novo Telegram Bot Token: " TELEGRAM_TOKEN
    echo ""
fi

# WhatsApp (opcional)
echo ""
echo -e "${YELLOW}2. 📱 WHATSAPP ACCESS TOKEN (opcional):${NC}"
read -p "Configurar WhatsApp agora? (y/n): " configure_whatsapp
if [ "$configure_whatsapp" = "y" ]; then
    read -s -p "WhatsApp Access Token: " WHATSAPP_TOKEN
    echo ""
    read -p "WhatsApp Phone Number ID: " WHATSAPP_PHONE_ID
    read -p "WhatsApp Business Account ID (opcional): " WHATSAPP_BUSINESS_ID
fi

# OpenAI (opcional)  
echo ""
echo -e "${YELLOW}3. 🤖 OPENAI API KEY (opcional):${NC}"
read -p "Configurar OpenAI? (y/n): " configure_openai
if [ "$configure_openai" = "y" ]; then
    read -s -p "OpenAI API Key: " OPENAI_KEY
    echo ""
fi

# Anthropic (opcional)
echo ""
echo -e "${YELLOW}4. 🧠 ANTHROPIC API KEY (opcional):${NC}"
read -p "Configurar Anthropic? (y/n): " configure_anthropic
if [ "$configure_anthropic" = "y" ]; then
    read -s -p "Anthropic API Key: " ANTHROPIC_KEY
    echo ""
fi

# Email
echo ""
echo -e "${YELLOW}5. 📧 EMAIL SMTP PASSWORD:${NC}"
echo "Email: contato@direitolux.com.br"
read -s -p "SMTP Password: " SMTP_PASSWORD
echo ""

# Configurar secrets no GitHub
echo ""
echo -e "${BLUE}🚀 Configurando secrets no GitHub...${NC}"

# Secrets obrigatórios
gh secret set TELEGRAM_BOT_TOKEN --body "$TELEGRAM_TOKEN"
echo -e "${GREEN}✅ TELEGRAM_BOT_TOKEN configurado${NC}"

gh secret set SMTP_PASSWORD --body "$SMTP_PASSWORD"
echo -e "${GREEN}✅ SMTP_PASSWORD configurado${NC}"

# WhatsApp (se configurado)
if [ "$configure_whatsapp" = "y" ]; then
    gh secret set WHATSAPP_ACCESS_TOKEN --body "$WHATSAPP_TOKEN"
    gh secret set WHATSAPP_PHONE_NUMBER_ID --body "$WHATSAPP_PHONE_ID"
    if [ ! -z "$WHATSAPP_BUSINESS_ID" ]; then
        gh secret set WHATSAPP_BUSINESS_ACCOUNT_ID --body "$WHATSAPP_BUSINESS_ID"
    fi
    echo -e "${GREEN}✅ WhatsApp secrets configurados${NC}"
fi

# OpenAI (se configurado)
if [ "$configure_openai" = "y" ]; then
    gh secret set OPENAI_API_KEY --body "$OPENAI_KEY"
    echo -e "${GREEN}✅ OPENAI_API_KEY configurado${NC}"
fi

# Anthropic (se configurado)
if [ "$configure_anthropic" = "y" ]; then
    gh secret set ANTHROPIC_API_KEY --body "$ANTHROPIC_KEY"
    echo -e "${GREEN}✅ ANTHROPIC_API_KEY configurado${NC}"
fi

# Secrets de infraestrutura (valores padrão seguros)
gh secret set DB_PASSWORD --body "direito_lux_prod_$(openssl rand -hex 8)"
gh secret set RABBITMQ_PASSWORD --body "rabbit_prod_$(openssl rand -hex 8)"

echo -e "${GREEN}✅ Secrets de infraestrutura configurados${NC}"

# Listar secrets configurados
echo ""
echo -e "${BLUE}📊 Secrets configurados no GitHub:${NC}"
gh secret list

echo ""
echo -e "${GREEN}🎉 GITHUB SECRETS CONFIGURADO COM SUCESSO!${NC}"
echo ""
echo -e "${YELLOW}📋 PRÓXIMOS PASSOS:${NC}"
echo "1. 🚀 Push para main → Deploy automático"
echo "2. 📱 Configure WhatsApp API (se não fez ainda)"
echo "3. 💰 Configure payment gateways (ASAAS/NOWPayments)"
echo "4. 🏛️ Configure DataJud CNJ (certificado digital)"
echo ""
echo -e "${BLUE}🔗 LINKS ÚTEIS:${NC}"
echo "• Repository: $(gh repo view --web --json url -q .url)"
echo "• Actions: $(gh repo view --web --json url -q .url)/actions"
echo "• Secrets: $(gh repo view --web --json url -q .url)/settings/secrets/actions"
echo ""
echo -e "${YELLOW}📖 Documentação: CONFIGURAR_GITHUB_SECRETS.md${NC}"

# Teste opcional
echo ""
read -p "🧪 Executar teste de validação? (y/n): " run_test
if [ "$run_test" = "y" ]; then
    echo ""
    echo -e "${BLUE}🧪 Testando configuração...${NC}"
    
    # Exportar variáveis para teste local
    export TELEGRAM_BOT_TOKEN="$TELEGRAM_TOKEN"
    export SMTP_PASSWORD="$SMTP_PASSWORD"
    
    if [ "$configure_whatsapp" = "y" ]; then
        export WHATSAPP_ACCESS_TOKEN="$WHATSAPP_TOKEN"
    fi
    
    # Verificar se serviços inicializam
    echo "Testando configuração do notification service..."
    if docker-compose config | grep -q "notification-service"; then
        echo -e "${GREEN}✅ Docker Compose configuração válida${NC}"
    else
        echo -e "${RED}❌ Problema na configuração Docker Compose${NC}"
    fi
    
    echo -e "${GREEN}✅ Teste de configuração concluído${NC}"
fi

echo ""
echo -e "${GREEN}🔐 GITHUB SECRETS TOTALMENTE OPERACIONAL!${NC}"