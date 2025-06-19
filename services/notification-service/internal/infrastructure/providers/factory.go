package providers

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
	"github.com/direito-lux/notification-service/internal/infrastructure/config"
)

// ProviderFactory factory para criar provedores de notificação
type ProviderFactory struct {
	config *config.Config
	logger *zap.Logger
}

// NewProviderFactory cria nova instância da factory
func NewProviderFactory(config *config.Config, logger *zap.Logger) *ProviderFactory {
	return &ProviderFactory{
		config: config,
		logger: logger,
	}
}

// CreateProviders cria mapa com todos os provedores disponíveis
func (f *ProviderFactory) CreateProviders() map[domain.NotificationChannel]domain.NotificationProvider {
	providers := make(map[domain.NotificationChannel]domain.NotificationProvider)

	// Email Provider
	emailProvider := f.createEmailProvider()
	if emailProvider != nil {
		providers[domain.NotificationChannelEmail] = emailProvider
	}

	// WhatsApp Provider
	whatsappProvider := f.createWhatsAppProvider()
	if whatsappProvider != nil {
		providers[domain.NotificationChannelWhatsApp] = whatsappProvider
	}

	// Telegram Provider
	telegramProvider := f.createTelegramProvider()
	if telegramProvider != nil {
		providers[domain.NotificationChannelTelegram] = telegramProvider
	}

	f.logger.Info("Providers created successfully", 
		zap.Int("email_available", boolToInt(emailProvider != nil)),
		zap.Int("whatsapp_available", boolToInt(whatsappProvider != nil)),
		zap.Int("telegram_available", boolToInt(telegramProvider != nil)),
		zap.Int("total_providers", len(providers)))

	return providers
}

// createEmailProvider cria provedor de email
func (f *ProviderFactory) createEmailProvider() domain.NotificationProvider {
	if f.config.SMTP.Host == "" || f.config.SMTP.Username == "" {
		f.logger.Warn("Email provider not configured - missing SMTP settings")
		return nil
	}

	emailConfig := EmailConfig{
		Host:        f.config.SMTP.Host,
		Port:        f.config.SMTP.Port,
		Username:    f.config.SMTP.Username,
		Password:    f.config.SMTP.Password,
		FromEmail:   f.config.SMTP.FromEmail,
		FromName:    f.config.SMTP.FromName,
		UseTLS:      f.config.SMTP.UseTLS,
		UseStartTLS: f.config.SMTP.UseStartTLS,
		Timeout:     30,
		MaxRetries:  3,
		RateLimit:   60,
	}

	provider := NewEmailProvider(emailConfig, f.logger)
	f.logger.Info("Email provider created successfully")
	return provider
}

// createWhatsAppProvider cria provedor WhatsApp
func (f *ProviderFactory) createWhatsAppProvider() domain.NotificationProvider {
	if f.config.WhatsApp.AccessToken == "" || f.config.WhatsApp.PhoneNumberID == "" {
		f.logger.Warn("WhatsApp provider not configured - missing access token or phone number ID")
		return nil
	}

	whatsappConfig := WhatsAppConfig{
		AccessToken:   f.config.WhatsApp.AccessToken,
		PhoneNumberID: f.config.WhatsApp.PhoneNumberID,
		WebhookURL:    f.config.WhatsApp.WebhookURL,
		VerifyToken:   f.config.WhatsApp.VerifyToken,
		Timeout:       30,
		MaxRetries:    3,
		RateLimit:     60,
	}

	provider := NewWhatsAppProvider(whatsappConfig, f.logger)
	f.logger.Info("WhatsApp provider created successfully")
	return provider
}

// createTelegramProvider cria provedor Telegram
func (f *ProviderFactory) createTelegramProvider() domain.NotificationProvider {
	if f.config.Telegram.BotToken == "" {
		f.logger.Warn("Telegram provider not configured - missing bot token")
		return nil
	}

	telegramConfig := TelegramConfig{
		BotToken:      f.config.Telegram.BotToken,
		WebhookURL:    f.config.Telegram.WebhookURL,
		WebhookSecret: f.config.Telegram.WebhookSecret,
		Timeout:       30,
		MaxRetries:    3,
		RateLimit:     30,
		ParseMode:     "HTML",
	}

	provider := NewTelegramProvider(telegramConfig, f.logger)
	f.logger.Info("Telegram provider created successfully")
	return provider
}

// GetAvailableChannels retorna canais disponíveis
func (f *ProviderFactory) GetAvailableChannels() []domain.NotificationChannel {
	var channels []domain.NotificationChannel

	if f.config.SMTP.Host != "" && f.config.SMTP.Username != "" {
		channels = append(channels, domain.NotificationChannelEmail)
	}

	if f.config.WhatsApp.AccessToken != "" && f.config.WhatsApp.PhoneNumberID != "" {
		channels = append(channels, domain.NotificationChannelWhatsApp)
	}

	if f.config.Telegram.BotToken != "" {
		channels = append(channels, domain.NotificationChannelTelegram)
	}

	return channels
}

// ValidateProviderConfig valida configuração de um provedor específico
func (f *ProviderFactory) ValidateProviderConfig(channel domain.NotificationChannel) error {
	switch channel {
	case domain.NotificationChannelEmail:
		if f.config.SMTP.Host == "" {
			return fmt.Errorf("SMTP host is required for email provider")
		}
		if f.config.SMTP.Username == "" {
			return fmt.Errorf("SMTP username is required for email provider")
		}
		if f.config.SMTP.Password == "" {
			return fmt.Errorf("SMTP password is required for email provider")
		}
		if f.config.SMTP.FromEmail == "" {
			return fmt.Errorf("SMTP from email is required for email provider")
		}

	case domain.NotificationChannelWhatsApp:
		if f.config.WhatsApp.AccessToken == "" {
			return fmt.Errorf("access token is required for WhatsApp provider")
		}
		if f.config.WhatsApp.PhoneNumberID == "" {
			return fmt.Errorf("phone number ID is required for WhatsApp provider")
		}
		if f.config.WhatsApp.VerifyToken == "" {
			return fmt.Errorf("verify token is required for WhatsApp provider")
		}

	case domain.NotificationChannelTelegram:
		if f.config.Telegram.BotToken == "" {
			return fmt.Errorf("bot token is required for Telegram provider")
		}

	default:
		return fmt.Errorf("unsupported notification channel: %s", channel)
	}

	return nil
}

// boolToInt converte bool para int para logging
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}