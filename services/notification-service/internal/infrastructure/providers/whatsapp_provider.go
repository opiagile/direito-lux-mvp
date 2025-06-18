package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// WhatsAppProvider implementação do provedor WhatsApp Business API
type WhatsAppProvider struct {
	config WhatsAppConfig
	client *http.Client
	logger *zap.Logger
}

// WhatsAppConfig configuração do provedor WhatsApp
type WhatsAppConfig struct {
	AccessToken   string `json:"access_token" validate:"required"`
	PhoneNumberID string `json:"phone_number_id" validate:"required"`
	WebhookURL    string `json:"webhook_url"`
	VerifyToken   string `json:"verify_token"`
	APIVersion    string `json:"api_version"`
	BaseURL       string `json:"base_url"`
	Timeout       int    `json:"timeout"` // segundos
	MaxRetries    int    `json:"max_retries"`
	RateLimit     int    `json:"rate_limit"` // mensagens por minuto
}

// NewWhatsAppProvider cria nova instância do provedor WhatsApp
func NewWhatsAppProvider(config WhatsAppConfig, logger *zap.Logger) *WhatsAppProvider {
	// Valores padrão
	if config.APIVersion == "" {
		config.APIVersion = "v18.0"
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://graph.facebook.com"
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RateLimit == 0 {
		config.RateLimit = 60 // 1 mensagem por segundo
	}

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	return &WhatsAppProvider{
		config: config,
		client: client,
		logger: logger,
	}
}

// Send envia uma notificação via WhatsApp
func (p *WhatsAppProvider) Send(ctx context.Context, notification *domain.Notification) error {
	p.logger.Debug("Sending WhatsApp notification", 
		zap.String("notification_id", notification.ID.String()),
		zap.String("recipient", notification.RecipientContact))

	// Validar número do destinatário
	if err := p.validatePhoneNumber(notification.RecipientContact); err != nil {
		return fmt.Errorf("invalid recipient phone number: %w", err)
	}

	// Construir mensagem
	message, err := p.buildMessage(notification)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// Enviar mensagem
	response, err := p.sendMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	// Atualizar notificação com ID externo
	if response.Messages != nil && len(response.Messages) > 0 {
		externalID := response.Messages[0].ID
		notification.ExternalID = &externalID
	}

	p.logger.Info("WhatsApp message sent successfully", 
		zap.String("notification_id", notification.ID.String()),
		zap.String("whatsapp_message_id", response.Messages[0].ID))

	return nil
}

// SendBatch envia múltiplas notificações em lote
func (p *WhatsAppProvider) SendBatch(ctx context.Context, notifications []*domain.Notification) error {
	p.logger.Debug("Sending batch WhatsApp notifications", zap.Int("count", len(notifications)))

	for _, notification := range notifications {
		if err := p.Send(ctx, notification); err != nil {
			p.logger.Error("Failed to send WhatsApp message in batch", 
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
func (p *WhatsAppProvider) Configure(ctx context.Context, config map[string]interface{}) error {
	p.logger.Debug("Configuring WhatsApp provider")

	// Converter mapa para struct
	if accessToken, ok := config["access_token"].(string); ok {
		p.config.AccessToken = accessToken
	}
	if phoneNumberID, ok := config["phone_number_id"].(string); ok {
		p.config.PhoneNumberID = phoneNumberID
	}
	if webhookURL, ok := config["webhook_url"].(string); ok {
		p.config.WebhookURL = webhookURL
	}
	if verifyToken, ok := config["verify_token"].(string); ok {
		p.config.VerifyToken = verifyToken
	}

	return p.ValidateConfiguration(config)
}

// ValidateConfiguration valida a configuração
func (p *WhatsAppProvider) ValidateConfiguration(config map[string]interface{}) error {
	required := []string{"access_token", "phone_number_id"}
	
	for _, field := range required {
		if _, exists := config[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validar access token format
	if accessToken, ok := config["access_token"].(string); ok {
		if !strings.HasPrefix(accessToken, "EAA") {
			return fmt.Errorf("invalid access token format")
		}
	}

	return nil
}

// IsHealthy verifica se o provedor está saudável
func (p *WhatsAppProvider) IsHealthy(ctx context.Context) bool {
	// Verificar se pode acessar a API do WhatsApp
	url := fmt.Sprintf("%s/%s/%s", p.config.BaseURL, p.config.APIVersion, p.config.PhoneNumberID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		p.logger.Warn("WhatsApp provider unhealthy - cannot create request", zap.Error(err))
		return false
	}

	req.Header.Set("Authorization", "Bearer "+p.config.AccessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Warn("WhatsApp provider unhealthy - request failed", zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GetChannel retorna o canal do provedor
func (p *WhatsAppProvider) GetChannel() domain.NotificationChannel {
	return domain.NotificationChannelWhatsApp
}

// SupportsHTML indica se suporta HTML
func (p *WhatsAppProvider) SupportsHTML() bool {
	return false
}

// SupportsAttachments indica se suporta anexos
func (p *WhatsAppProvider) SupportsAttachments() bool {
	return true
}

// SupportsTemplates indica se suporta templates
func (p *WhatsAppProvider) SupportsTemplates() bool {
	return true
}

// GetMaxContentLength retorna o tamanho máximo do conteúdo
func (p *WhatsAppProvider) GetMaxContentLength() int {
	return 4096 // WhatsApp limit
}

// GetRateLimit retorna o limite de taxa
func (p *WhatsAppProvider) GetRateLimit() int {
	return p.config.RateLimit
}

// validatePhoneNumber valida o formato do número de telefone
func (p *WhatsAppProvider) validatePhoneNumber(phoneNumber string) error {
	// Remover caracteres não numéricos
	cleaned := strings.ReplaceAll(phoneNumber, "+", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Verificar se tem pelo menos 10 dígitos
	if len(cleaned) < 10 {
		return fmt.Errorf("phone number too short: %s", phoneNumber)
	}

	// Verificar se tem no máximo 15 dígitos (padrão internacional)
	if len(cleaned) > 15 {
		return fmt.Errorf("phone number too long: %s", phoneNumber)
	}

	return nil
}

// buildMessage constrói a mensagem para o WhatsApp API
func (p *WhatsAppProvider) buildMessage(notification *domain.Notification) (*WhatsAppMessage, error) {
	// Limpar número de telefone
	phoneNumber := strings.ReplaceAll(notification.RecipientContact, "+", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")

	message := &WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               phoneNumber,
		Type:             "text",
		Text: &WhatsAppText{
			Body: notification.Content,
		},
	}

	return message, nil
}

// sendMessage envia a mensagem via WhatsApp API
func (p *WhatsAppProvider) sendMessage(ctx context.Context, message *WhatsAppMessage) (*WhatsAppResponse, error) {
	url := fmt.Sprintf("%s/%s/%s/messages", p.config.BaseURL, p.config.APIVersion, p.config.PhoneNumberID)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.AccessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response WhatsAppResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if response.Error != nil {
			return nil, fmt.Errorf("WhatsApp API error: %s (%d)", response.Error.Message, response.Error.Code)
		}
		return nil, fmt.Errorf("WhatsApp API returned status: %d", resp.StatusCode)
	}

	return &response, nil
}

// Estruturas para a API do WhatsApp

// WhatsAppMessage estrutura da mensagem
type WhatsAppMessage struct {
	MessagingProduct string         `json:"messaging_product"`
	To               string         `json:"to"`
	Type             string         `json:"type"`
	Text             *WhatsAppText  `json:"text,omitempty"`
	Template         *WhatsAppTemplate `json:"template,omitempty"`
}

// WhatsAppText conteúdo de texto
type WhatsAppText struct {
	Body string `json:"body"`
}

// WhatsAppTemplate template de mensagem
type WhatsAppTemplate struct {
	Name       string                   `json:"name"`
	Language   WhatsAppLanguage         `json:"language"`
	Components []WhatsAppComponent      `json:"components,omitempty"`
}

// WhatsAppLanguage idioma do template
type WhatsAppLanguage struct {
	Code string `json:"code"`
}

// WhatsAppComponent componente do template
type WhatsAppComponent struct {
	Type       string               `json:"type"`
	Parameters []WhatsAppParameter  `json:"parameters,omitempty"`
}

// WhatsAppParameter parâmetro do template
type WhatsAppParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// WhatsAppResponse resposta da API
type WhatsAppResponse struct {
	Messages []WhatsAppMessageResponse `json:"messages,omitempty"`
	Error    *WhatsAppError            `json:"error,omitempty"`
}

// WhatsAppMessageResponse resposta de mensagem enviada
type WhatsAppMessageResponse struct {
	ID string `json:"id"`
}

// WhatsAppError erro da API
type WhatsAppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// HandleWebhook processa webhooks do WhatsApp
func (p *WhatsAppProvider) HandleWebhook(ctx context.Context, payload []byte) error {
	p.logger.Debug("Processing WhatsApp webhook", zap.ByteString("payload", payload))

	var webhook WhatsAppWebhook
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	// Processar cada entrada no webhook
	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			if change.Value.Messages != nil {
				for _, message := range change.Value.Messages {
					p.logger.Info("Received WhatsApp message", 
						zap.String("message_id", message.ID),
						zap.String("from", message.From))
				}
			}

			if change.Value.Statuses != nil {
				for _, status := range change.Value.Statuses {
					p.logger.Info("Received WhatsApp status update", 
						zap.String("message_id", status.ID),
						zap.String("status", status.Status))
					
					// Aqui seria atualizado o status da notificação no banco de dados
					// baseado no status recebido (sent, delivered, read, failed)
				}
			}
		}
	}

	return nil
}

// Estruturas para webhooks

// WhatsAppWebhook estrutura do webhook
type WhatsAppWebhook struct {
	Object string              `json:"object"`
	Entry  []WhatsAppWebhookEntry `json:"entry"`
}

// WhatsAppWebhookEntry entrada do webhook
type WhatsAppWebhookEntry struct {
	ID      string                 `json:"id"`
	Changes []WhatsAppWebhookChange `json:"changes"`
}

// WhatsAppWebhookChange mudança no webhook
type WhatsAppWebhookChange struct {
	Value WhatsAppWebhookValue `json:"value"`
	Field string               `json:"field"`
}

// WhatsAppWebhookValue valor do webhook
type WhatsAppWebhookValue struct {
	Messages []WhatsAppWebhookMessage `json:"messages,omitempty"`
	Statuses []WhatsAppWebhookStatus  `json:"statuses,omitempty"`
}

// WhatsAppWebhookMessage mensagem do webhook
type WhatsAppWebhookMessage struct {
	ID   string `json:"id"`
	From string `json:"from"`
	Type string `json:"type"`
	Text struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`
}

// WhatsAppWebhookStatus status do webhook
type WhatsAppWebhookStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}