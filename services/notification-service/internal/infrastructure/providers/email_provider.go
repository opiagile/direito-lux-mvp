package providers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// EmailProvider implementação do provedor de email SMTP
type EmailProvider struct {
	config EmailConfig
	logger *zap.Logger
}

// EmailConfig configuração do provedor de email
type EmailConfig struct {
	Host         string `json:"host" validate:"required"`
	Port         int    `json:"port" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	FromEmail    string `json:"from_email" validate:"required,email"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"`
	UseStartTLS  bool   `json:"use_starttls"`
	Timeout      int    `json:"timeout"` // segundos
	MaxRetries   int    `json:"max_retries"`
	RateLimit    int    `json:"rate_limit"` // emails per minute
}

// NewEmailProvider cria nova instância do provedor de email
func NewEmailProvider(config EmailConfig, logger *zap.Logger) *EmailProvider {
	// Valores padrão
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RateLimit == 0 {
		config.RateLimit = 60 // 1 email per second
	}

	return &EmailProvider{
		config: config,
		logger: logger,
	}
}

// Send envia uma notificação por email
func (p *EmailProvider) Send(ctx context.Context, notification *domain.Notification) error {
	p.logger.Debug("Sending email notification", 
		zap.String("notification_id", notification.ID.String()),
		zap.String("recipient", notification.RecipientContact))

	// Validar email do destinatário
	if err := p.validateEmail(notification.RecipientContact); err != nil {
		return fmt.Errorf("invalid recipient email: %w", err)
	}

	// Construir mensagem
	message, err := p.buildMessage(notification)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// Enviar email
	if err := p.sendEmail(ctx, notification.RecipientContact, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	p.logger.Info("Email sent successfully", 
		zap.String("notification_id", notification.ID.String()),
		zap.String("recipient", notification.RecipientContact))

	return nil
}

// SendBatch envia múltiplas notificações em lote
func (p *EmailProvider) SendBatch(ctx context.Context, notifications []*domain.Notification) error {
	p.logger.Debug("Sending batch email notifications", zap.Int("count", len(notifications)))

	for _, notification := range notifications {
		if err := p.Send(ctx, notification); err != nil {
			p.logger.Error("Failed to send email in batch", 
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
func (p *EmailProvider) Configure(ctx context.Context, config map[string]interface{}) error {
	p.logger.Debug("Configuring email provider")

	// Converter mapa para struct
	if host, ok := config["host"].(string); ok {
		p.config.Host = host
	}
	if port, ok := config["port"].(int); ok {
		p.config.Port = port
	}
	if username, ok := config["username"].(string); ok {
		p.config.Username = username
	}
	if password, ok := config["password"].(string); ok {
		p.config.Password = password
	}
	if fromEmail, ok := config["from_email"].(string); ok {
		p.config.FromEmail = fromEmail
	}
	if fromName, ok := config["from_name"].(string); ok {
		p.config.FromName = fromName
	}
	if useTLS, ok := config["use_tls"].(bool); ok {
		p.config.UseTLS = useTLS
	}
	if useStartTLS, ok := config["use_starttls"].(bool); ok {
		p.config.UseStartTLS = useStartTLS
	}

	return p.ValidateConfiguration(config)
}

// ValidateConfiguration valida a configuração
func (p *EmailProvider) ValidateConfiguration(config map[string]interface{}) error {
	required := []string{"host", "port", "username", "password", "from_email"}
	
	for _, field := range required {
		if _, exists := config[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validar email
	if fromEmail, ok := config["from_email"].(string); ok {
		if err := p.validateEmail(fromEmail); err != nil {
			return fmt.Errorf("invalid from_email: %w", err)
		}
	}

	// Validar porta
	if port, ok := config["port"].(int); ok {
		if port <= 0 || port > 65535 {
			return fmt.Errorf("invalid port: %d", port)
		}
	}

	return nil
}

// IsHealthy verifica se o provedor está saudável
func (p *EmailProvider) IsHealthy(ctx context.Context) bool {
	// Testar conexão SMTP
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.config.Host, p.config.Port), 
		time.Duration(p.config.Timeout)*time.Second)
	if err != nil {
		p.logger.Warn("Email provider unhealthy - cannot connect", zap.Error(err))
		return false
	}
	defer conn.Close()

	return true
}

// GetChannel retorna o canal do provedor
func (p *EmailProvider) GetChannel() domain.NotificationChannel {
	return domain.NotificationChannelEmail
}

// SupportsHTML indica se suporta HTML
func (p *EmailProvider) SupportsHTML() bool {
	return true
}

// SupportsAttachments indica se suporta anexos
func (p *EmailProvider) SupportsAttachments() bool {
	return true
}

// SupportsTemplates indica se suporta templates
func (p *EmailProvider) SupportsTemplates() bool {
	return true
}

// GetMaxContentLength retorna o tamanho máximo do conteúdo
func (p *EmailProvider) GetMaxContentLength() int {
	return 64 * 1024 * 1024 // 64MB
}

// GetRateLimit retorna o limite de taxa
func (p *EmailProvider) GetRateLimit() int {
	return p.config.RateLimit
}

// validateEmail valida formato do email
func (p *EmailProvider) validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// buildMessage constrói a mensagem de email
func (p *EmailProvider) buildMessage(notification *domain.Notification) ([]byte, error) {
	from := p.config.FromEmail
	if p.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", p.config.FromName, p.config.FromEmail)
	}

	// Headers básicos
	headers := map[string]string{
		"From":         from,
		"To":           notification.RecipientContact,
		"Subject":      notification.Subject,
		"MIME-Version": "1.0",
		"Date":         time.Now().Format(time.RFC1123Z),
		"Message-ID":   fmt.Sprintf("<%s@%s>", notification.ID.String(), p.extractDomain(p.config.FromEmail)),
	}

	// Detectar se o conteúdo é HTML
	isHTML := p.isHTMLContent(notification.Content)
	if isHTML {
		headers["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
	}

	// Construir mensagem
	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(notification.Content)

	return []byte(message.String()), nil
}

// sendEmail envia o email via SMTP
func (p *EmailProvider) sendEmail(ctx context.Context, to string, message []byte) error {
	// Conectar ao servidor SMTP
	addr := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)
	
	var client *smtp.Client
	var err error

	if p.config.UseTLS {
		// Conexão TLS direta
		tlsConfig := &tls.Config{
			ServerName: p.config.Host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect with TLS: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, p.config.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	} else {
		// Conexão normal
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
	}
	defer client.Quit()

	// STARTTLS se configurado
	if p.config.UseStartTLS {
		tlsConfig := &tls.Config{
			ServerName: p.config.Host,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Autenticação
	auth := smtp.PlainAuth("", p.config.Username, p.config.Password, p.config.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Enviar email
	if err := client.Mail(p.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// isHTMLContent verifica se o conteúdo é HTML
func (p *EmailProvider) isHTMLContent(content string) bool {
	content = strings.ToLower(strings.TrimSpace(content))
	htmlTags := []string{"<html", "<body", "<div", "<p", "<br", "<span", "<h1", "<h2", "<h3", "<table", "<tr", "<td"}
	
	for _, tag := range htmlTags {
		if strings.Contains(content, tag) {
			return true
		}
	}
	
	return false
}

// extractDomain extrai o domínio do email
func (p *EmailProvider) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return "localhost"
}