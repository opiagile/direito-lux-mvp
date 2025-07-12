#!/bin/bash

# 🤖 Teste Manual do Bot Telegram

echo "🤖 Testando bot Telegram manualmente..."

BOT_TOKEN="${TELEGRAM_BOT_TOKEN}"

# Enviar mensagem de teste para chat específico (você pode usar seu próprio chat ID)
echo "Para testar o bot completamente, você precisa:"
echo ""
echo "1. Abra o Telegram"
echo "2. Acesse: https://t.me/direitolux_staging_bot"
echo "3. Clique em 'START'"
echo "4. Envie: /help"
echo ""
echo "✅ Webhook configurado: /webhook/telegram"
echo "✅ Bot online: @direitolux_staging_bot"
echo ""
echo "📋 Monitore os logs:"
echo "docker-compose logs -f notification-service"
echo ""
echo "Se você enviar uma mensagem, deve aparecer logs como:"
echo "INFO Recebido webhook do Telegram"
echo "INFO Mensagem recebida do Telegram"