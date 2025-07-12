#!/bin/bash

# 🤖 Teste do Telegram Bot API - Direito Lux Staging
# Este script testa a configuração do bot do Telegram

echo "🤖 Testando configuração do Telegram Bot..."

# Token do bot (obtido via BotFather)
BOT_TOKEN="7404885967:AAGkfZJIr8zVpKfJ4PO6YaZzPCVYKYUzYE4"

# Função para testar getMe
test_bot_info() {
    echo "📋 Testando informações do bot..."
    
    response=$(curl -s "https://api.telegram.org/bot${BOT_TOKEN}/getMe")
    
    if echo "$response" | grep -q '"ok":true'; then
        echo "✅ Bot configurado com sucesso!"
        echo "📊 Informações do bot:"
        echo "$response" | jq '.result' 2>/dev/null || echo "$response"
    else
        echo "❌ Falha na configuração do bot"
        echo "🔍 Resposta da API:"
        echo "$response"
        return 1
    fi
}

# Função para configurar comandos do bot
set_bot_commands() {
    echo "📋 Configurando comandos do bot..."
    
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
        echo "✅ Comandos configurados com sucesso!"
    else
        echo "❌ Falha ao configurar comandos"
        echo "🔍 Resposta da API:"
        echo "$response"
        return 1
    fi
}

# Função para definir descrição do bot
set_bot_description() {
    echo "📋 Configurando descrição do bot..."
    
    description="🏛️ Direito Lux - Assistente Jurídico Inteligente

Bot oficial para monitoramento de processos, notificações automáticas e análises jurídicas.

✅ Ambiente STAGING - Para testes e validação"
    
    response=$(curl -s -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setMyDescription" \
        -H "Content-Type: application/json" \
        -d "{\"description\": \"$description\"}")
    
    if echo "$response" | grep -q '"ok":true'; then
        echo "✅ Descrição configurada com sucesso!"
    else
        echo "❌ Falha ao configurar descrição"
        echo "🔍 Resposta da API:"
        echo "$response"
        return 1
    fi
}

# Função para testar health do notification service
test_notification_service() {
    echo "📋 Testando notification service..."
    
    # Verificar se o serviço está rodando
    if curl -s http://localhost:8085/health > /dev/null; then
        echo "✅ Notification service está rodando"
    else
        echo "⚠️ Notification service não está acessível"
        echo "💡 Execute: docker-compose up -d notification-service"
        return 1
    fi
}

# Executar testes
echo "🎯 Direito Lux - Configuração Telegram Bot API"
echo "=============================================="

# Testar informações do bot
test_bot_info
if [ $? -ne 0 ]; then
    echo "❌ Falha no teste básico do bot"
    exit 1
fi

# Configurar comandos
set_bot_commands
if [ $? -ne 0 ]; then
    echo "⚠️ Falha ao configurar comandos (continuando...)"
fi

# Configurar descrição
set_bot_description
if [ $? -ne 0 ]; then
    echo "⚠️ Falha ao configurar descrição (continuando...)"
fi

# Testar notification service
test_notification_service

echo ""
echo "🎉 Configuração do bot concluída!"
echo "📱 Para testar o bot: https://t.me/direitolux_staging_bot"
echo "🔗 Username: @direitolux_staging_bot"
echo ""
echo "📋 Próximos passos:"
echo "1. Configurar webhook HTTPS (ngrok ou domínio real)"
echo "2. Testar envio de mensagens via notification service"
echo "3. Validar recebimento de mensagens do bot"
echo ""
echo "✅ Bot pronto para integração!"