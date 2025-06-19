package dto

import (
	"time"
)

// WhatsApp DTOs

// WhatsAppWebhookRequest webhook do WhatsApp
type WhatsAppWebhookRequest struct {
	Entry []WhatsAppEntry `json:"entry"`
}

// WhatsAppEntry entrada do webhook
type WhatsAppEntry struct {
	ID      string                  `json:"id"`
	Changes []WhatsAppChangeEntry `json:"changes"`
}

// WhatsAppChangeEntry mudança no webhook
type WhatsAppChangeEntry struct {
	Value WhatsAppValue `json:"value"`
	Field string        `json:"field"`
}

// WhatsAppValue valor da mudança
type WhatsAppValue struct {
	MessagingProduct string                `json:"messaging_product"`
	Metadata        WhatsAppMetadata      `json:"metadata"`
	Contacts        []WhatsAppContact     `json:"contacts"`
	Messages        []WhatsAppMessage     `json:"messages"`
	Statuses        []WhatsAppStatus      `json:"statuses"`
}

// WhatsAppMetadata metadados
type WhatsAppMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// WhatsAppContact contato
type WhatsAppContact struct {
	Profile WhatsAppProfile `json:"profile"`
	WaID    string          `json:"wa_id"`
}

// WhatsAppProfile perfil do contato
type WhatsAppProfile struct {
	Name string `json:"name"`
}

// WhatsAppMessage mensagem do WhatsApp
type WhatsAppMessage struct {
	From      string                     `json:"from"`
	ID        string                     `json:"id"`
	Timestamp string                     `json:"timestamp"`
	Type      string                     `json:"type"`
	Text      *WhatsAppTextMessage      `json:"text,omitempty"`
	Image     *WhatsAppMediaMessage     `json:"image,omitempty"`
	Document  *WhatsAppMediaMessage     `json:"document,omitempty"`
	Audio     *WhatsAppMediaMessage     `json:"audio,omitempty"`
	Video     *WhatsAppMediaMessage     `json:"video,omitempty"`
	Location  *WhatsAppLocationMessage `json:"location,omitempty"`
	Context   *WhatsAppContext          `json:"context,omitempty"`
}

// WhatsAppTextMessage mensagem de texto
type WhatsAppTextMessage struct {
	Body string `json:"body"`
}

// WhatsAppMediaMessage mensagem de mídia
type WhatsAppMediaMessage struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	Filename string `json:"filename,omitempty"`
	Caption  string `json:"caption,omitempty"`
}

// WhatsAppLocationMessage mensagem de localização
type WhatsAppLocationMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// WhatsAppContext contexto da mensagem
type WhatsAppContext struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

// WhatsAppStatus status da mensagem
type WhatsAppStatus struct {
	ID           string                    `json:"id"`
	Status       string                    `json:"status"`
	Timestamp    string                    `json:"timestamp"`
	RecipientID  string                    `json:"recipient_id"`
	Conversation *WhatsAppConversationInfo `json:"conversation,omitempty"`
	Pricing      *WhatsAppPricingInfo      `json:"pricing,omitempty"`
}

// WhatsAppConversationInfo informações da conversa
type WhatsAppConversationInfo struct {
	ID                  string               `json:"id"`
	ExpirationTimestamp string               `json:"expiration_timestamp,omitempty"`
	Origin              WhatsAppOriginInfo   `json:"origin"`
}

// WhatsAppOriginInfo origem da conversa
type WhatsAppOriginInfo struct {
	Type string `json:"type"`
}

// WhatsAppPricingInfo informações de preço
type WhatsAppPricingInfo struct {
	Billable     bool   `json:"billable"`
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
}

// Telegram DTOs

// TelegramWebhookRequest webhook do Telegram
type TelegramWebhookRequest struct {
	UpdateID int                  `json:"update_id"`
	Message  *TelegramMessage     `json:"message,omitempty"`
	Callback *TelegramCallbackQuery `json:"callback_query,omitempty"`
}

// TelegramMessage mensagem do Telegram
type TelegramMessage struct {
	MessageID int                `json:"message_id"`
	From      TelegramUser       `json:"from"`
	Chat      TelegramChat       `json:"chat"`
	Date      int64              `json:"date"`
	Text      string             `json:"text,omitempty"`
	Photo     []TelegramPhotoSize `json:"photo,omitempty"`
	Document  *TelegramDocument   `json:"document,omitempty"`
	Audio     *TelegramAudio      `json:"audio,omitempty"`
	Video     *TelegramVideo      `json:"video,omitempty"`
	Location  *TelegramLocation   `json:"location,omitempty"`
	Caption   string             `json:"caption,omitempty"`
}

// TelegramUser usuário do Telegram
type TelegramUser struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
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

// TelegramPhotoSize tamanho da foto
type TelegramPhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

// TelegramDocument documento do Telegram
type TelegramDocument struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

// TelegramAudio áudio do Telegram
type TelegramAudio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer,omitempty"`
	Title        string `json:"title,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

// TelegramVideo vídeo do Telegram
type TelegramVideo struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

// TelegramLocation localização do Telegram
type TelegramLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TelegramCallbackQuery callback query
type TelegramCallbackQuery struct {
	ID           string           `json:"id"`
	From         TelegramUser     `json:"from"`
	Message      *TelegramMessage `json:"message,omitempty"`
	Data         string           `json:"data,omitempty"`
	GameShortName string          `json:"game_short_name,omitempty"`
}

// Slack DTOs

// SlackEventRequest evento do Slack
type SlackEventRequest struct {
	Token     string      `json:"token"`
	TeamID    string      `json:"team_id"`
	APIAppID  string      `json:"api_app_id"`
	Event     SlackEvent  `json:"event"`
	Type      string      `json:"type"`
	EventID   string      `json:"event_id"`
	EventTime int64       `json:"event_time"`
	Challenge string      `json:"challenge,omitempty"` // Para verificação de URL
}

// SlackEvent evento específico
type SlackEvent struct {
	Type      string `json:"type"`
	Channel   string `json:"channel,omitempty"`
	User      string `json:"user,omitempty"`
	Text      string `json:"text,omitempty"`
	Timestamp string `json:"ts,omitempty"`
	ThreadTS  string `json:"thread_ts,omitempty"`
	EventTS   string `json:"event_ts,omitempty"`
}

// SlackCommandRequest comando slash do Slack
type SlackCommandRequest struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}

// Bot Message Sending DTOs

// SendBotMessageRequest request para enviar mensagem via bot
type SendBotMessageRequest struct {
	Channel     string                 `json:"channel" binding:"required"` // whatsapp, telegram, slack
	Recipient   string                 `json:"recipient" binding:"required"`
	Message     string                 `json:"message" binding:"required"`
	MessageType string                 `json:"message_type,omitempty"`
	Attachments []MessageAttachment   `json:"attachments,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// BotMessageResponse response do envio
type BotMessageResponse struct {
	MessageID   string    `json:"message_id"`
	Channel     string    `json:"channel"`
	Recipient   string    `json:"recipient"`
	Status      string    `json:"status"`
	SentAt      time.Time `json:"sent_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}