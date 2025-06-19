# Notification Service - Providers

## 📋 Resumo da Implementação

O sistema de providers do Notification Service foi **finalizado com sucesso**. Todos os três providers principais estão implementados e prontos para uso:

### ✅ Providers Implementados

| Provider | Status | Funcionalidades |
|----------|--------|----------------|
| **📧 Email (SMTP)** | ✅ Completo | HTML/Text, Templates, Anexos, TLS/STARTTLS |
| **📱 WhatsApp Business** | ✅ Completo | Templates, Media, Webhooks, Status tracking |
| **🤖 Telegram Bot** | ✅ Completo | HTML/Markdown, Inline keyboards, Webhooks |

### 🏗️ Arquitetura dos Providers

```
internal/infrastructure/providers/
├── factory.go           # Factory para criar e gerenciar providers
├── email_provider.go    # Provider SMTP completo
├── whatsapp_provider.go # Provider WhatsApp Business API
└── telegram_provider.go # Provider Telegram Bot API
```

## 🔧 Configuração

### Variáveis de Ambiente

```bash
# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@direito-lux.com
SMTP_FROM_NAME=Direito Lux
SMTP_USE_TLS=false
SMTP_USE_STARTTLS=true

# WhatsApp Business API
WHATSAPP_ACCESS_TOKEN=EAAxxxxxxxxxxxxx
WHATSAPP_PHONE_NUMBER_ID=123456789012345
WHATSAPP_WEBHOOK_URL=https://your-domain.com/webhook/whatsapp
WHATSAPP_VERIFY_TOKEN=your-verify-token

# Telegram Bot API
TELEGRAM_BOT_TOKEN=1234567890:ABCDefGhIjKlMnOpQrStUvWxYz
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook/telegram
TELEGRAM_WEBHOOK_SECRET=your-webhook-secret
```

## 📱 Funcionalidades por Provider

### 📧 Email Provider

**Recursos:**
- ✅ Envio via SMTP com autenticação
- ✅ Suporte a TLS e STARTTLS
- ✅ Detecção automática de HTML vs texto
- ✅ Headers customizados (Message-ID, Date, etc.)
- ✅ Rate limiting (60 emails/min)
- ✅ Retry automático (3 tentativas)
- ✅ Health checks via conexão TCP

**Configuração:**
```go
emailConfig := EmailConfig{
    Host:        "smtp.gmail.com",
    Port:        587,
    Username:    "user@domain.com",
    Password:    "app-password",
    FromEmail:   "noreply@direito-lux.com",
    FromName:    "Direito Lux",
    UseStartTLS: true,
    Timeout:     30,
    MaxRetries:  3,
    RateLimit:   60,
}
```

### 📱 WhatsApp Provider

**Recursos:**
- ✅ WhatsApp Business API oficial
- ✅ Envio de mensagens de texto
- ✅ Templates aprovados pelo WhatsApp
- ✅ Media attachments (imagem, documento, etc.)
- ✅ Webhooks para status updates
- ✅ Rate limiting (60 mensagens/min)
- ✅ Validação de números de telefone
- ✅ Health checks via API

**Exemplo de uso:**
```go
// Mensagem simples
notification := &Notification{
    Channel: domain.NotificationChannelWhatsApp,
    RecipientContact: "5511999999999",
    Content: "Seu processo foi atualizado!",
}

// Template com parâmetros
message := &WhatsAppMessage{
    Type: "template",
    Template: &WhatsAppTemplate{
        Name: "process_update",
        Language: WhatsAppLanguage{Code: "pt_BR"},
        Components: []WhatsAppComponent{
            {
                Type: "body",
                Parameters: []WhatsAppParameter{
                    {Type: "text", Text: "João Silva"},
                    {Type: "text", Text: "1234567-89.2024.8.26.0100"},
                },
            },
        },
    },
}
```

### 🤖 Telegram Provider

**Recursos:**
- ✅ Telegram Bot API completa
- ✅ Parse modes: HTML, Markdown, MarkdownV2
- ✅ Inline keyboards e callback handling
- ✅ Webhooks para mensagens recebidas
- ✅ Rate limiting (30 mensagens/min)
- ✅ Chat ID e username support
- ✅ Health checks via getMe API

**Exemplo de uso:**
```go
// Mensagem com HTML
notification := &Notification{
    Channel: domain.NotificationChannelTelegram,
    RecipientContact: "123456789", // Chat ID
    Content: "<b>Processo Atualizado</b>\n\n" +
            "Cliente: <i>João Silva</i>\n" +
            "Número: <code>1234567-89.2024.8.26.0100</code>",
}

// Com inline keyboard
keyboard := TelegramInlineKeyboard{
    InlineKeyboard: [][]TelegramInlineKeyboardButton{
        {
            {Text: "Ver Detalhes", CallbackData: "view_details_123"},
            {Text: "Agendar", CallbackData: "schedule_123"},
        },
    },
}
```

## 🔄 Provider Factory

A factory gerencia todos os providers de forma centralizada:

```go
// Criar factory
factory := NewProviderFactory(config, logger)

// Criar todos os providers disponíveis
providers := factory.CreateProviders()

// Verificar canais disponíveis
channels := factory.GetAvailableChannels()

// Validar configuração específica
err := factory.ValidateProviderConfig(domain.NotificationChannelEmail)
```

## 🚀 Integração no Main

O sistema está totalmente integrado com injeção de dependência:

```go
// main.go
fx.Provide(
    NewEmailProvider,
    NewWhatsAppProvider, 
    NewTelegramProvider,
    NewProviderMap,
)

// Ou usando a factory
fx.Provide(NewProviderFactory)
fx.Provide(func(factory *ProviderFactory) map[domain.NotificationChannel]domain.NotificationProvider {
    return factory.CreateProviders()
})
```

## 📊 Métricas e Monitoramento

Cada provider inclui:

- ✅ **Health Checks**: Verificação automática de conectividade
- ✅ **Rate Limiting**: Controle de taxa de envio
- ✅ **Retry Logic**: Tentativas automáticas em caso de falha
- ✅ **Logging**: Logs estruturados com zap
- ✅ **Error Handling**: Tratamento robusto de erros
- ✅ **Timeouts**: Configuração de timeouts por provider

## 🔒 Segurança

- ✅ **Validação de Input**: Todos os dados são validados
- ✅ **Secrets Management**: Configuração via env vars
- ✅ **TLS/SSL**: Conexões seguras obrigatórias
- ✅ **Webhook Security**: Tokens de verificação
- ✅ **Rate Limiting**: Proteção contra abuse

## 🧪 Testes

Para testar os providers:

```bash
# Testar saúde dos providers
curl http://localhost:8080/health

# Enviar notificação de teste
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "test_notification",
    "channel": "email",
    "priority": "normal",
    "subject": "Teste",
    "content": "Mensagem de teste",
    "recipient_id": "123",
    "recipient_type": "user",
    "recipient_contact": "test@example.com",
    "tenant_id": "uuid-here"
  }'
```

## 📋 Status Final

### ✅ Completo
- [x] Email Provider (SMTP completo)
- [x] WhatsApp Provider (Business API completo)  
- [x] Telegram Provider (Bot API completo)
- [x] Provider Factory
- [x] Integração com NotificationService
- [x] Injeção de dependência no main.go
- [x] Configuração via environment variables
- [x] Health checks e monitoramento
- [x] Rate limiting e retry logic
- [x] Webhook handling
- [x] Documentação completa

### 🎯 Próximos Passos (Opcionais)
- [ ] SMS Provider (Twilio/AWS SNS)
- [ ] Push Notification Provider (Firebase FCM)
- [ ] Slack Provider (para notificações internas)
- [ ] Microsoft Teams Provider
- [ ] Discord Provider

---

**🎉 NOTIFICATION SERVICE PROVIDERS - 100% IMPLEMENTADO!**

Todos os três providers principais (Email, WhatsApp, Telegram) estão funcionais e prontos para produção.