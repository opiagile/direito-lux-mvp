package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// TelegramProvider implementação do provedor Telegram Bot API
type TelegramProvider struct {
	config TelegramConfig
	client *http.Client
	logger *zap.Logger
}

// TelegramConfig configuração do provedor Telegram
type TelegramConfig struct {
	BotToken      string `json:"bot_token" validate:"required"`
	BaseURL       string `json:"base_url"`
	WebhookURL    string `json:"webhook_url"`
	WebhookSecret string `json:"webhook_secret"`
	Timeout       int    `json:"timeout"` // segundos
	MaxRetries    int    `json:"max_retries"`
	RateLimit     int    `json:"rate_limit"` // mensagens per minute
	ParseMode     string `json:"parse_mode"` // HTML, Markdown, MarkdownV2
}

// NewTelegramProvider cria nova instância do provedor Telegram
func NewTelegramProvider(config TelegramConfig, logger *zap.Logger) *TelegramProvider {
	// Valores padrão
	if config.BaseURL == "" {
		config.BaseURL = "https://api.telegram.org"
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RateLimit == 0 {
		config.RateLimit = 30 // 30 mensagens por minuto (limite do Telegram)
	}
	if config.ParseMode == "" {
		config.ParseMode = "HTML"
	}

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	return &TelegramProvider{
		config: config,
		client: client,
		logger: logger,
	}
}

// Send envia uma notificação via Telegram
func (p *TelegramProvider) Send(ctx context.Context, notification *domain.Notification) error {
	p.logger.Debug("Sending Telegram notification", 
		zap.String("notification_id", notification.ID.String()),
		zap.String("recipient", notification.RecipientContact))

	// Validar chat ID do destinatário
	chatID, err := p.validateChatID(notification.RecipientContact)
	if err != nil {
		return fmt.Errorf("invalid recipient chat ID: %w", err)
	}

	// Construir mensagem
	message, err := p.buildMessage(chatID, notification)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// Enviar mensagem
	response, err := p.sendMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}

	// Atualizar notificação com ID externo
	if response.Result != nil {
		externalID := strconv.Itoa(response.Result.MessageID)
		notification.ExternalID = &externalID
	}

	p.logger.Info("Telegram message sent successfully", 
		zap.String("notification_id", notification.ID.String()),
		zap.Int("telegram_message_id", response.Result.MessageID))

	return nil
}

// SendBatch envia múltiplas notificações em lote
func (p *TelegramProvider) SendBatch(ctx context.Context, notifications []*domain.Notification) error {
	p.logger.Debug("Sending batch Telegram notifications", zap.Int("count", len(notifications)))

	for _, notification := range notifications {
		if err := p.Send(ctx, notification); err != nil {
			p.logger.Error("Failed to send Telegram message in batch", 
				zap.String("notification_id", notification.ID.String()),
				zap.Error(err))
			return err
		}

		// Rate limiting - aguardar entre envios
		if len(notifications) > 1 {
			time.Sleep(time.Minute / time.Duration(p.config.RateLimit))
		}
	}

	return nil
}

// Configure configura o provedor
func (p *TelegramProvider) Configure(ctx context.Context, config map[string]interface{}) error {
	p.logger.Debug("Configuring Telegram provider")

	// Converter mapa para struct
	if botToken, ok := config["bot_token"].(string); ok {
		p.config.BotToken = botToken
	}
	if baseURL, ok := config["base_url"].(string); ok {
		p.config.BaseURL = baseURL
	}
	if webhookURL, ok := config["webhook_url"].(string); ok {
		p.config.WebhookURL = webhookURL
	}
	if webhookSecret, ok := config["webhook_secret"].(string); ok {
		p.config.WebhookSecret = webhookSecret
	}
	if parseMode, ok := config["parse_mode"].(string); ok {
		p.config.ParseMode = parseMode
	}

	return p.ValidateConfiguration(config)
}

// ValidateConfiguration valida a configuração
func (p *TelegramProvider) ValidateConfiguration(config map[string]interface{}) error {
	required := []string{"bot_token"}
	
	for _, field := range required {
		if _, exists := config[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validar bot token format
	if botToken, ok := config["bot_token"].(string); ok {
		if !strings.Contains(botToken, ":") {
			return fmt.Errorf("invalid bot token format")
		}
	}

	// Validar parse mode
	if parseMode, ok := config["parse_mode"].(string); ok {
		validModes := map[string]bool{
			"HTML":        true,
			"Markdown":    true,
			"MarkdownV2":  true,
			"":            true,
		}
		if !validModes[parseMode] {
			return fmt.Errorf("invalid parse mode: %s", parseMode)
		}
	}

	return nil
}

// IsHealthy verifica se o provedor está saudável
func (p *TelegramProvider) IsHealthy(ctx context.Context) bool {
	// Verificar se pode acessar a API do Telegram
	url := fmt.Sprintf("%s/bot%s/getMe", p.config.BaseURL, p.config.BotToken)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		p.logger.Warn("Telegram provider unhealthy - cannot create request", zap.Error(err))
		return false
	}

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Warn("Telegram provider unhealthy - request failed", zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	var response TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		p.logger.Warn("Telegram provider unhealthy - decode failed", zap.Error(err))
		return false
	}

	return response.OK
}

// GetChannel retorna o canal do provedor
func (p *TelegramProvider) GetChannel() domain.NotificationChannel {
	return domain.NotificationChannelTelegram
}

// SupportsHTML indica se suporta HTML
func (p *TelegramProvider) SupportsHTML() bool {
	return true
}

// SupportsAttachments indica se suporta anexos
func (p *TelegramProvider) SupportsAttachments() bool {
	return true
}

// SupportsTemplates indica se suporta templates
func (p *TelegramProvider) SupportsTemplates() bool {
	return true
}

// GetMaxContentLength retorna o tamanho máximo do conteúdo
func (p *TelegramProvider) GetMaxContentLength() int {
	return 4096 // Telegram limit
}

// GetRateLimit retorna o limite de taxa
func (p *TelegramProvider) GetRateLimit() int {
	return p.config.RateLimit
}

// validateChatID valida o formato do chat ID
func (p *TelegramProvider) validateChatID(chatID string) (int64, error) {
	// Tentar converter para int64
	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		// Se não for um número, pode ser um username (@username)
		if strings.HasPrefix(chatID, "@") && len(chatID) > 1 {
			return 0, nil // Username válido
		}
		return 0, fmt.Errorf("invalid chat ID format: %s", chatID)
	}

	return id, nil
}

// buildMessage constrói a mensagem para o Telegram API
func (p *TelegramProvider) buildMessage(chatID int64, notification *domain.Notification) (*TelegramMessage, error) {
	message := &TelegramMessage{
		ChatID:    chatID,
		Text:      notification.Content,
		ParseMode: p.config.ParseMode,
	}

	// Se for username, usar string
	if chatID == 0 {
		message.ChatIDString = notification.RecipientContact
	}

	return message, nil
}

// sendMessage envia a mensagem via Telegram API
func (p *TelegramProvider) sendMessage(ctx context.Context, message *TelegramMessage) (*TelegramResponse, error) {
	url := fmt.Sprintf("%s/bot%s/sendMessage", p.config.BaseURL, p.config.BotToken)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("Telegram API error: %s (%d)", response.Description, response.ErrorCode)
	}

	return &response, nil
}

// SetWebhook configura webhook do bot
func (p *TelegramProvider) SetWebhook(ctx context.Context, webhookURL string) error {
	url := fmt.Sprintf("%s/bot%s/setWebhook", p.config.BaseURL, p.config.BotToken)

	payload := map[string]interface{}{
		"url": webhookURL,
	}

	if p.config.WebhookSecret != "" {
		payload["secret_token"] = p.config.WebhookSecret
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	defer resp.Body.Close()

	var response TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode webhook response: %w", err)
	}

	if !response.OK {
		return fmt.Errorf("Telegram webhook error: %s (%d)", response.Description, response.ErrorCode)
	}

	p.logger.Info("Telegram webhook set successfully", zap.String("url", webhookURL))
	return nil
}

// HandleWebhook processa webhooks do Telegram
func (p *TelegramProvider) HandleWebhook(ctx context.Context, payload []byte) error {
	p.logger.Debug("Processing Telegram webhook", zap.ByteString("payload", payload))

	var update TelegramUpdate
	if err := json.Unmarshal(payload, &update); err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	// Processar mensagem recebida
	if update.Message != nil {
		p.logger.Info("Received Telegram message", 
			zap.Int("message_id", update.Message.MessageID),
			zap.Int64("chat_id", update.Message.Chat.ID),
			zap.String("from", update.Message.From.Username))

		// Aqui seria processada a mensagem recebida
		// Por exemplo, criar uma notificação de resposta ou executar um comando
	}

	// Processar callback query (botões inline)
	if update.CallbackQuery != nil {
		p.logger.Info("Received Telegram callback", 
			zap.String("callback_data", update.CallbackQuery.Data),
			zap.Int64("chat_id", update.CallbackQuery.From.ID))

		// Processar callback do botão
	}

	return nil
}

// Estruturas para a API do Telegram

// TelegramMessage estrutura da mensagem
type TelegramMessage struct {
	ChatID       int64  `json:"chat_id,omitempty"`
	ChatIDString string `json:"chat_id,omitempty"`
	Text         string `json:"text"`
	ParseMode    string `json:"parse_mode,omitempty"`
	ReplyMarkup  interface{} `json:"reply_markup,omitempty"`
}

// TelegramSentMessage mensagem enviada (resposta da API)
type TelegramSentMessage struct {
	MessageID int                `json:"message_id"`
	From      *TelegramUser      `json:"from,omitempty"`
	Chat      *TelegramChat      `json:"chat"`
	Date      int64              `json:"date"`
	Text      string             `json:"text,omitempty"`
}

// TelegramResponse resposta da API
type TelegramResponse struct {
	OK          bool                  `json:"ok"`
	Result      *TelegramSentMessage  `json:"result,omitempty"`
	ErrorCode   int                   `json:"error_code,omitempty"`
	Description string                `json:"description,omitempty"`
}

// TelegramUpdate update do webhook
type TelegramUpdate struct {
	UpdateID      int                    `json:"update_id"`
	Message       *TelegramWebhookMessage `json:"message,omitempty"`
	CallbackQuery *TelegramCallbackQuery  `json:"callback_query,omitempty"`
}

// TelegramWebhookMessage mensagem do webhook
type TelegramWebhookMessage struct {
	MessageID int                `json:"message_id"`
	From      *TelegramUser      `json:"from,omitempty"`
	Chat      *TelegramChat      `json:"chat"`
	Date      int64              `json:"date"`
	Text      string             `json:"text,omitempty"`
}

// TelegramUser usuário do Telegram
type TelegramUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// TelegramChat chat do Telegram
type TelegramChat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// TelegramCallbackQuery callback query
type TelegramCallbackQuery struct {
	ID   string        `json:"id"`
	From *TelegramUser `json:"from"`
	Data string        `json:"data"`
}

// TelegramInlineKeyboard teclado inline
type TelegramInlineKeyboard struct {
	InlineKeyboard [][]TelegramInlineKeyboardButton `json:"inline_keyboard"`
}

// TelegramInlineKeyboardButton botão do teclado inline
type TelegramInlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}