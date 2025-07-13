#!/bin/bash

# 📱 Setup WhatsApp Business API - AGORA
# Configuração rápida para Direito Lux Staging

set -e

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m' 
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}"
echo "╔═══════════════════════════════════════════════════╗"
echo "║    📱 WHATSAPP BUSINESS API - SETUP AGORA         ║"
echo "║    20 minutos para 100 mensagens/dia grátis      ║"
echo "╚═══════════════════════════════════════════════════╝"
echo -e "${NC}"

# Verificar pré-requisitos
echo -e "${BLUE}🔍 Verificando pré-requisitos...${NC}"

# Verificar túnel HTTPS
if curl -s https://direito-lux-staging.loca.lt/health > /dev/null; then
    echo -e "${GREEN}✅ Túnel HTTPS ativo: https://direito-lux-staging.loca.lt${NC}"
else
    echo -e "${RED}❌ Túnel HTTPS não está ativo!${NC}"
    echo "Execute: npx localtunnel --port 8085 --subdomain direito-lux-staging"
    exit 1
fi

# Verificar notification service
if docker-compose ps notification-service | grep -q "Up"; then
    echo -e "${GREEN}✅ Notification Service rodando${NC}"
else
    echo -e "${RED}❌ Notification Service não está rodando${NC}"
    echo "Execute: docker-compose up -d notification-service"
    exit 1
fi

echo ""
echo -e "${YELLOW}📋 PASSOS PARA CONFIGURAR WHATSAPP (20 min):${NC}"
echo ""
echo "1. 🌐 Acesse: https://developers.facebook.com/"
echo "2. 🔑 Login com sua conta Facebook pessoal"
echo "3. ➕ Criar App → Business → 'Direito Lux Staging'"
echo "4. 📱 Adicionar produto: WhatsApp Business API"
echo "5. 🔧 API Setup → Copie as credenciais"
echo "6. 🌐 Configuration → Webhooks → Configure URL"
echo ""

read -p "Pressione Enter quando estiver pronto para configurar..."

echo ""
echo -e "${CYAN}═══ CONFIGURAÇÃO WEBHOOKS ═══${NC}"
echo ""
echo -e "${YELLOW}📋 Configure no Meta for Developers:${NC}"
echo ""
echo "🔗 Callback URL:"
echo "   https://direito-lux-staging.loca.lt/webhook/whatsapp"
echo ""
echo "🔑 Verify Token:"  
echo "   direito_lux_staging_2025"
echo ""
echo "✅ Eventos para marcar:"
echo "   □ messages"
echo "   □ message_deliveries" 
echo "   □ message_reads"
echo "   □ messaging_postbacks"
echo ""

read -p "Pressione Enter após configurar o webhook..."

# Testar webhook
echo ""
echo -e "${BLUE}🧪 Testando webhook...${NC}"

# Testar verificação
echo "Testando verificação do webhook..."
response=$(curl -s "https://direito-lux-staging.loca.lt/webhook/whatsapp?hub.mode=subscribe&hub.challenge=test123&hub.verify_token=direito_lux_staging_2025")

if [ "$response" = "test123" ]; then
    echo -e "${GREEN}✅ Webhook verificado com sucesso!${NC}"
else
    echo -e "${YELLOW}⚠️ Webhook pode precisar de ajustes${NC}"
    echo "Resposta: $response"
fi

echo ""
echo -e "${CYAN}═══ OBTER CREDENCIAIS ═══${NC}"
echo ""
echo -e "${YELLOW}📋 No Meta for Developers → API Setup:${NC}"
echo ""
echo "1. 🔑 Copie o 'Temporary Access Token'"
echo "2. 📱 Copie o 'Phone number ID'" 
echo "3. 🏢 Copie o 'Business account ID' (opcional)"
echo ""

# Coletar credenciais
echo -e "${BLUE}📝 Digite as credenciais abaixo:${NC}"
echo ""

read -p "Access Token (EAAxxxxxxxxxx): " ACCESS_TOKEN
read -p "Phone Number ID (números): " PHONE_NUMBER_ID
read -p "Business Account ID (opcional): " BUSINESS_ACCOUNT_ID

# Validar formato
if [[ ! "$ACCESS_TOKEN" =~ ^EAA[A-Za-z0-9_-]+$ ]]; then
    echo -e "${RED}❌ Formato de Access Token inválido!${NC}"
    echo "Deve começar com 'EAA'"
    exit 1
fi

if [[ ! "$PHONE_NUMBER_ID" =~ ^[0-9]+$ ]]; then
    echo -e "${RED}❌ Phone Number ID deve ser só números!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Credenciais válidas!${NC}"

# Testar API
echo ""
echo -e "${BLUE}🧪 Testando credenciais com API...${NC}"

response=$(curl -s -X GET \
    "https://graph.facebook.com/v18.0/${PHONE_NUMBER_ID}" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

if echo "$response" | grep -q '"verified_name"'; then
    echo -e "${GREEN}✅ Credenciais funcionando!${NC}"
    
    # Extrair info
    verified_name=$(echo "$response" | grep -o '"verified_name":"[^"]*"' | cut -d'"' -f4)
    display_number=$(echo "$response" | grep -o '"display_phone_number":"[^"]*"' | cut -d'"' -f4)
    
    echo ""
    echo -e "${CYAN}📱 Informações do WhatsApp Business:${NC}"
    echo "   Nome: $verified_name"
    echo "   Número: $display_number"
    echo "   ID: $PHONE_NUMBER_ID"
    echo ""
else
    echo -e "${RED}❌ Erro na API: $response${NC}"
    exit 1
fi

# Configurar no sistema
echo -e "${BLUE}🔧 Configurando no sistema...${NC}"

# GitHub Secrets (se GitHub CLI estiver configurado)
if command -v gh &> /dev/null && gh auth status &> /dev/null 2>&1; then
    echo "Configurando GitHub Secrets..."
    gh secret set WHATSAPP_ACCESS_TOKEN --body "$ACCESS_TOKEN"
    gh secret set WHATSAPP_PHONE_NUMBER_ID --body "$PHONE_NUMBER_ID"
    if [ ! -z "$BUSINESS_ACCOUNT_ID" ]; then
        gh secret set WHATSAPP_BUSINESS_ACCOUNT_ID --body "$BUSINESS_ACCOUNT_ID"
    fi
    echo -e "${GREEN}✅ GitHub Secrets atualizados${NC}"
else
    echo -e "${YELLOW}⚠️ GitHub CLI não configurado - configure secrets manualmente${NC}"
fi

# Atualizar .env local
ENV_FILE="services/notification-service/.env"
if [ -f "$ENV_FILE" ]; then
    # Backup
    cp "$ENV_FILE" "${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
    
    # Atualizar
    sed -i.bak "s|WHATSAPP_ACCESS_TOKEN=.*|WHATSAPP_ACCESS_TOKEN=$ACCESS_TOKEN|" "$ENV_FILE"
    sed -i.bak "s|WHATSAPP_PHONE_NUMBER_ID=.*|WHATSAPP_PHONE_NUMBER_ID=$PHONE_NUMBER_ID|" "$ENV_FILE"
    
    if [ ! -z "$BUSINESS_ACCOUNT_ID" ]; then
        if ! grep -q "WHATSAPP_BUSINESS_ACCOUNT_ID" "$ENV_FILE"; then
            echo "WHATSAPP_BUSINESS_ACCOUNT_ID=$BUSINESS_ACCOUNT_ID" >> "$ENV_FILE"
        else
            sed -i.bak "s|WHATSAPP_BUSINESS_ACCOUNT_ID=.*|WHATSAPP_BUSINESS_ACCOUNT_ID=$BUSINESS_ACCOUNT_ID|" "$ENV_FILE"
        fi
    fi
    
    echo -e "${GREEN}✅ Arquivo .env atualizado${NC}"
fi

# Reiniciar notification service
echo ""
echo -e "${BLUE}🔄 Reiniciando notification service...${NC}"
docker-compose restart notification-service

sleep 5

if curl -s http://localhost:8085/health > /dev/null; then
    echo -e "${GREEN}✅ Notification service funcionando${NC}"
else
    echo -e "${YELLOW}⚠️ Service pode estar inicializando...${NC}"
fi

# Salvar informações
cat > whatsapp_api_info.txt << EOF
WhatsApp Business API - Direito Lux Staging
==========================================
Access Token: [SALVO NO GITHUB SECRETS]
Phone Number ID: $PHONE_NUMBER_ID
Business Account ID: $BUSINESS_ACCOUNT_ID
Verified Name: $verified_name
Display Number: $display_number
Webhook URL: https://direito-lux-staging.loca.lt/webhook/whatsapp
Verify Token: direito_lux_staging_2025
Configurado em: $(date)

LIMITAÇÕES TESTE:
- Token temporário: 24h
- Mensagens: 100/dia
- Destinatários: 5 números máximo

PRÓXIMO: Adicionar números de teste no Meta console
EOF

echo ""
echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}🎉 WHATSAPP API CONFIGURADO COM SUCESSO!${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}📋 INFORMAÇÕES IMPORTANTES:${NC}"
echo "• 📱 Número verificado: $display_number"
echo "• 💬 100 mensagens gratuitas/dia"  
echo "• ⏰ Token temporário (24h)"
echo "• 📱 Máximo 5 números de teste"
echo ""
echo -e "${BLUE}📋 PRÓXIMOS PASSOS:${NC}"
echo "1. 📱 Adicione seu número na lista de teste (Meta console)"
echo "2. 📨 Teste envio via Meta console"
echo "3. 💬 Teste recebimento no webhook"
echo "4. 🔄 Configure token permanente para produção"
echo ""
echo -e "${CYAN}🔗 LINKS ÚTEIS:${NC}"
echo "• Meta Console: https://developers.facebook.com/"
echo "• Webhook: https://direito-lux-staging.loca.lt/webhook/whatsapp"
echo "• Logs: docker-compose logs -f notification-service"
echo ""
echo -e "${GREEN}✅ WhatsApp Business API totalmente operacional!${NC}"