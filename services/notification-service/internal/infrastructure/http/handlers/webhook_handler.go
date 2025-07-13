package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WebhookHandler handler para webhooks de provedores externos
type WebhookHandler struct {
	logger    *zap.Logger
	botToken  string
}

// NewWebhookHandler cria novo handler de webhooks
func NewWebhookHandler(logger *zap.Logger) *WebhookHandler {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	return &WebhookHandler{
		logger:   logger,
		botToken: botToken,
	}
}

// TelegramUpdate estrutura do update do Telegram
type TelegramUpdate struct {
	UpdateID int              `json:"update_id"`
	Message  *TelegramMessage `json:"message,omitempty"`
}

// TelegramMessage estrutura da mensagem do Telegram
type TelegramMessage struct {
	MessageID int           `json:"message_id"`
	From      *TelegramUser `json:"from,omitempty"`
	Chat      *TelegramChat `json:"chat"`
	Date      int64         `json:"date"`
	Text      string        `json:"text,omitempty"`
}

// TelegramUser estrutura do usuário do Telegram
type TelegramUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// TelegramChat estrutura do chat do Telegram
type TelegramChat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// HandleTelegramWebhook processa webhooks do Telegram
func (h *WebhookHandler) HandleTelegramWebhook(c *gin.Context) {
	h.logger.Info("Recebido webhook do Telegram", 
		zap.String("method", c.Request.Method),
		zap.String("user_agent", c.GetHeader("User-Agent")),
		zap.String("content_type", c.GetHeader("Content-Type")))

	var update TelegramUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		h.logger.Error("Erro ao fazer parse do webhook Telegram", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	h.logger.Info("Update do Telegram processado",
		zap.Int("update_id", update.UpdateID),
		zap.Any("message", update.Message))

	// Processar mensagem recebida
	if update.Message != nil {
		h.processMessage(update.Message)
	}

	// Telegram espera 200 OK
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// processMessage processa mensagem recebida do Telegram
func (h *WebhookHandler) processMessage(message *TelegramMessage) {
	h.logger.Info("Mensagem recebida do Telegram",
		zap.Int("message_id", message.MessageID),
		zap.Int64("chat_id", message.Chat.ID),
		zap.String("text", message.Text))

	var responseText string

	switch message.Text {
	case "/start":
		h.logger.Info("Comando /start recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `🏛️ *Bem-vindo ao Direito Lux!*

🤖 Sou seu assistente jurídico inteligente para monitoramento de processos.

📋 *Comandos disponíveis:*
/help - Lista todos os comandos
/status - Status dos seus processos
/agenda - Prazos importantes
/busca - Buscar processos
/relatorio - Relatório rápido
/configurar - Configurações

✅ *Ambiente STAGING* - Para testes e validação

Digite /help para ver mais opções!`

	case "/help":
		h.logger.Info("Comando /help recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `🆘 *Central de Ajuda - Direito Lux*

📋 *Comandos principais:*
/start - Iniciar conversa
/help - Esta mensagem de ajuda
/status - 📊 Status dos processos monitorados
/agenda - 📅 Próximos prazos e audiências
/busca - 🔍 Buscar processos por número/parte
/relatorio - 📈 Relatório resumido dos processos
/configurar - ⚙️ Configurar notificações

💡 *Exemplos de uso:*
• Digite um número de processo para buscar
• Use /status para ver resumo dos seus processos
• Configure alertas com /configurar

🔗 *Links úteis:*
• Web App: https://direito-lux-staging.loca.lt
• Suporte: Entre em contato com nossa equipe

✅ *Versão STAGING* - Ambiente de testes`

	case "/status":
		h.logger.Info("Comando /status recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `📊 *Status dos Processos*

🟢 *Ativos:* 0 processos monitorados
🟡 *Pendentes:* 0 atualizações
🔴 *Alertas:* 0 prazos vencendo

📈 *Resumo último mês:*
• Movimentações: 0
• Notificações enviadas: 0
• Última atualização: --

⚠️ *Ambiente STAGING*
Para adicionar processos reais, acesse:
https://direito-lux-staging.loca.lt

Use /help para ver mais comandos.`

	case "/agenda":
		h.logger.Info("Comando /agenda recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `📅 *Agenda de Prazos*

🗓️ *Próximos 7 dias:*
• Nenhum prazo cadastrado

⏰ *Alertas configurados:*
• 3 dias antes do vencimento
• 1 dia antes do vencimento
• No dia do vencimento

📋 *Para adicionar prazos:*
1. Acesse a plataforma web
2. Cadastre seus processos
3. Configure alertas automáticos

🔗 Acessar: https://direito-lux-staging.loca.lt

Use /configurar para ajustar notificações.`

	case "/busca":
		h.logger.Info("Comando /busca recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `🔍 *Busca de Processos*

📝 *Como buscar:*
• Digite o número do processo
• Formato: 1234567-89.2023.8.26.0001
• Ou envie apenas os números

🎯 *Exemplo:*
1234567892023826001

📊 *Informações retornadas:*
• Dados do processo
• Última movimentação
• Partes envolvidas
• Status atual

⚠️ *Ambiente STAGING*
Alguns processos podem não estar disponíveis.

💡 *Dica:* Para busca avançada, use a plataforma web:
https://direito-lux-staging.loca.lt`

	case "/relatorio":
		h.logger.Info("Comando /relatorio recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `📈 *Relatório Rápido*

📊 *Resumo Geral:*
• Total de processos: 0
• Ativos: 0
• Arquivados: 0
• Suspensos: 0

📅 *Últimos 30 dias:*
• Movimentações: 0
• Notificações: 0
• Alertas enviados: 0

🏛️ *Por tribunal:*
• TJSP: 0 processos
• TJRJ: 0 processos
• Outros: 0 processos

📄 *Relatório completo:*
Para relatórios detalhados, acesse:
https://direito-lux-staging.loca.lt/relatorios

Use /status para informações atualizadas.`

	case "/configurar":
		h.logger.Info("Comando /configurar recebido", zap.Int64("chat_id", message.Chat.ID))
		responseText = `⚙️ *Configurações*

🔔 *Notificações ativas:*
• Telegram: ✅ Ativo
• Email: ❌ Não configurado
• WhatsApp: ❌ Não configurado

⏰ *Frequência de alertas:*
• Movimentações: Instantâneo
• Prazos: 3, 1 dia antes e no dia
• Audiências: 7, 3, 1 dia antes

🎯 *Personalizar:*
Para configurações detalhadas, acesse:
https://direito-lux-staging.loca.lt/configuracoes

📱 *Canais disponíveis:*
• Telegram (atual)
• Email
• WhatsApp
• SMS

Use a plataforma web para configurações avançadas.`

	default:
		// Verifica se é um número de processo
		if len(message.Text) > 10 && h.isProcessNumber(message.Text) {
			h.logger.Info("Possível número de processo recebido", 
				zap.String("text", message.Text),
				zap.Int64("chat_id", message.Chat.ID))
			responseText = fmt.Sprintf(`🔍 *Buscando processo:* %s

⏳ Consultando DataJud CNJ...

⚠️ *Ambiente STAGING*
Em produção, aqui apareceriam:
• Dados do processo
• Última movimentação
• Partes envolvidas
• Próximos prazos

💡 Para busca real, acesse:
https://direito-lux-staging.loca.lt`, message.Text)
		} else {
			h.logger.Info("Mensagem de texto livre recebida", 
				zap.String("text", message.Text),
				zap.Int64("chat_id", message.Chat.ID))
			responseText = `🤖 *Direito Lux Bot*

Não entendi sua mensagem. 

📋 *Comandos disponíveis:*
/help - Ver todos os comandos
/status - Status dos processos
/busca - Buscar processo

💡 *Dicas:*
• Use os comandos com /
• Envie um número de processo para buscar
• Digite /help para ver todas as opções

✅ *Bot ativo* - Ambiente STAGING`
		}
	}

	// Enviar resposta
	if responseText != "" {
		h.sendTelegramMessage(message.Chat.ID, responseText)
	}
}

// isProcessNumber verifica se o texto parece ser um número de processo
func (h *WebhookHandler) isProcessNumber(text string) bool {
	// Verificação básica para números de processo CNJ
	// Formato: NNNNNNN-DD.AAAA.J.TR.OOOO ou só números
	return len(text) >= 15 && len(text) <= 25
}

// TelegramSendMessage estrutura para envio de mensagem
type TelegramSendMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// sendTelegramMessage envia mensagem de resposta via API do Telegram
func (h *WebhookHandler) sendTelegramMessage(chatID int64, text string) {
	if h.botToken == "" {
		h.logger.Error("Bot token não configurado")
		return
	}

	message := TelegramSendMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Erro ao serializar mensagem", zap.Error(err))
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.botToken)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		h.logger.Error("Erro ao enviar mensagem para Telegram", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		h.logger.Info("Mensagem enviada com sucesso para Telegram", 
			zap.Int64("chat_id", chatID),
			zap.String("text_preview", text[:min(50, len(text))]))
	} else {
		h.logger.Error("Erro ao enviar mensagem para Telegram", 
			zap.Int("status_code", resp.StatusCode),
			zap.Int64("chat_id", chatID))
	}
}

// min retorna o menor entre dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HandleWhatsAppWebhook processa webhooks do WhatsApp (placeholder)
func (h *WebhookHandler) HandleWhatsAppWebhook(c *gin.Context) {
	h.logger.Info("Recebido webhook do WhatsApp", 
		zap.String("method", c.Request.Method),
		zap.String("user_agent", c.GetHeader("User-Agent")))

	// Verificação do webhook (GET request)
	if c.Request.Method == "GET" {
		hubMode := c.Query("hub.mode")
		hubChallenge := c.Query("hub.challenge")
		hubVerifyToken := c.Query("hub.verify_token")

		h.logger.Info("Verificação do webhook WhatsApp",
			zap.String("hub.mode", hubMode),
			zap.String("hub.challenge", hubChallenge),
			zap.String("hub.verify_token", hubVerifyToken))

		// Validar token de verificação
		expectedToken := "direito_lux_staging_2025"
		if hubMode == "subscribe" && hubVerifyToken == expectedToken {
			h.logger.Info("Webhook WhatsApp verificado com sucesso")
			c.String(http.StatusOK, hubChallenge)
			return
		}

		h.logger.Warn("Token de verificação inválido",
			zap.String("expected", expectedToken),
			zap.String("received", hubVerifyToken))
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verify token"})
		return
	}

	// Processar webhook (POST request)
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Erro ao fazer parse do webhook WhatsApp", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	h.logger.Info("Webhook WhatsApp processado", zap.Any("payload", payload))

	// WhatsApp espera 200 OK
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}