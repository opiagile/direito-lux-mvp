package domain

import (
	"time"

	"github.com/google/uuid"
)

// Eventos de Domínio do Notification Service

// NotificationCreatedEvent evento de notificação criada
type NotificationCreatedEvent struct {
	NotificationID   uuid.UUID            `json:"notification_id"`
	Type             NotificationType     `json:"type"`
	Channel          NotificationChannel  `json:"channel"`
	Priority         NotificationPriority `json:"priority"`
	Subject          string               `json:"subject"`
	RecipientID      string               `json:"recipient_id"`
	RecipientType    string               `json:"recipient_type"`
	RecipientContact string               `json:"recipient_contact"`
	TenantID         uuid.UUID            `json:"tenant_id"`
	UserID           *uuid.UUID           `json:"user_id,omitempty"`
	ProcessID        *uuid.UUID           `json:"process_id,omitempty"`
	TemplateID       *uuid.UUID           `json:"template_id,omitempty"`
	ScheduledAt      *time.Time           `json:"scheduled_at,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
}

// NotificationSentEvent evento de notificação enviada
type NotificationSentEvent struct {
	NotificationID   uuid.UUID           `json:"notification_id"`
	Type             NotificationType    `json:"type"`
	Channel          NotificationChannel `json:"channel"`
	RecipientID      string              `json:"recipient_id"`
	RecipientContact string              `json:"recipient_contact"`
	TenantID         uuid.UUID           `json:"tenant_id"`
	UserID           *uuid.UUID          `json:"user_id,omitempty"`
	ProcessID        *uuid.UUID          `json:"process_id,omitempty"`
	SentAt           time.Time           `json:"sent_at"`
	Duration         time.Duration       `json:"duration"`
}

// NotificationFailedEvent evento de notificação falhada
type NotificationFailedEvent struct {
	NotificationID   uuid.UUID           `json:"notification_id"`
	Type             NotificationType    `json:"type"`
	Channel          NotificationChannel `json:"channel"`
	RecipientID      string              `json:"recipient_id"`
	RecipientContact string              `json:"recipient_contact"`
	TenantID         uuid.UUID           `json:"tenant_id"`
	UserID           *uuid.UUID          `json:"user_id,omitempty"`
	ProcessID        *uuid.UUID          `json:"process_id,omitempty"`
	ErrorMessage     string              `json:"error_message"`
	RetryCount       int                 `json:"retry_count"`
	MaxRetries       int                 `json:"max_retries"`
	FailedAt         time.Time           `json:"failed_at"`
}

// NotificationRetriedEvent evento de notificação reenviada
type NotificationRetriedEvent struct {
	NotificationID   uuid.UUID           `json:"notification_id"`
	Type             NotificationType    `json:"type"`
	Channel          NotificationChannel `json:"channel"`
	RecipientID      string              `json:"recipient_id"`
	RecipientContact string              `json:"recipient_contact"`
	TenantID         uuid.UUID           `json:"tenant_id"`
	RetryCount       int                 `json:"retry_count"`
	MaxRetries       int                 `json:"max_retries"`
	RetriedAt        time.Time           `json:"retried_at"`
}

// NotificationCancelledEvent evento de notificação cancelada
type NotificationCancelledEvent struct {
	NotificationID   uuid.UUID           `json:"notification_id"`
	Type             NotificationType    `json:"type"`
	Channel          NotificationChannel `json:"channel"`
	RecipientID      string              `json:"recipient_id"`
	TenantID         uuid.UUID           `json:"tenant_id"`
	Reason           string              `json:"reason"`
	CancelledAt      time.Time           `json:"cancelled_at"`
}

// NotificationExpiredEvent evento de notificação expirada
type NotificationExpiredEvent struct {
	NotificationID   uuid.UUID           `json:"notification_id"`
	Type             NotificationType    `json:"type"`
	Channel          NotificationChannel `json:"channel"`
	RecipientID      string              `json:"recipient_id"`
	TenantID         uuid.UUID           `json:"tenant_id"`
	ExpiresAt        time.Time           `json:"expires_at"`
	ExpiredAt        time.Time           `json:"expired_at"`
}

// TemplateCreatedEvent evento de template criado
type TemplateCreatedEvent struct {
	TemplateID uuid.UUID           `json:"template_id"`
	Name       string              `json:"name"`
	Type       NotificationType    `json:"type"`
	Channel    NotificationChannel `json:"channel"`
	TenantID   *uuid.UUID          `json:"tenant_id,omitempty"`
	IsSystem   bool                `json:"is_system"`
	CreatedAt  time.Time           `json:"created_at"`
}

// TemplateUpdatedEvent evento de template atualizado
type TemplateUpdatedEvent struct {
	TemplateID uuid.UUID           `json:"template_id"`
	Name       string              `json:"name"`
	Type       NotificationType    `json:"type"`
	Channel    NotificationChannel `json:"channel"`
	Status     TemplateStatus      `json:"status"`
	TenantID   *uuid.UUID          `json:"tenant_id,omitempty"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// TemplateActivatedEvent evento de template ativado
type TemplateActivatedEvent struct {
	TemplateID uuid.UUID           `json:"template_id"`
	Name       string              `json:"name"`
	Type       NotificationType    `json:"type"`
	Channel    NotificationChannel `json:"channel"`
	TenantID   *uuid.UUID          `json:"tenant_id,omitempty"`
	ActivatedAt time.Time          `json:"activated_at"`
}

// TemplateDeactivatedEvent evento de template desativado
type TemplateDeactivatedEvent struct {
	TemplateID    uuid.UUID           `json:"template_id"`
	Name          string              `json:"name"`
	Type          NotificationType    `json:"type"`
	Channel       NotificationChannel `json:"channel"`
	TenantID      *uuid.UUID          `json:"tenant_id,omitempty"`
	DeactivatedAt time.Time           `json:"deactivated_at"`
}

// BulkNotificationCreatedEvent evento de notificações em lote criadas
type BulkNotificationCreatedEvent struct {
	BatchID        uuid.UUID            `json:"batch_id"`
	Type           NotificationType     `json:"type"`
	Channel        NotificationChannel  `json:"channel"`
	Priority       NotificationPriority `json:"priority"`
	TenantID       uuid.UUID            `json:"tenant_id"`
	UserID         *uuid.UUID           `json:"user_id,omitempty"`
	Count          int                  `json:"count"`
	TemplateID     *uuid.UUID           `json:"template_id,omitempty"`
	ScheduledAt    *time.Time           `json:"scheduled_at,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
}

// BulkNotificationCompletedEvent evento de notificações em lote completadas
type BulkNotificationCompletedEvent struct {
	BatchID     uuid.UUID `json:"batch_id"`
	TotalCount  int       `json:"total_count"`
	SentCount   int       `json:"sent_count"`
	FailedCount int       `json:"failed_count"`
	Duration    time.Duration `json:"duration"`
	CompletedAt time.Time `json:"completed_at"`
}

// NotificationQuotaExceededEvent evento de quota de notificações excedida
type NotificationQuotaExceededEvent struct {
	TenantID      uuid.UUID           `json:"tenant_id"`
	Channel       NotificationChannel `json:"channel"`
	CurrentUsage  int                 `json:"current_usage"`
	QuotaLimit    int                 `json:"quota_limit"`
	Period        string              `json:"period"`
	ExceededAt    time.Time           `json:"exceeded_at"`
}

// NotificationChannelErrorEvent evento de erro no canal de notificação
type NotificationChannelErrorEvent struct {
	Channel      NotificationChannel `json:"channel"`
	ErrorType    string              `json:"error_type"`
	ErrorMessage string              `json:"error_message"`
	TenantID     *uuid.UUID          `json:"tenant_id,omitempty"`
	ErrorCount   int                 `json:"error_count"`
	OccurredAt   time.Time           `json:"occurred_at"`
}

// WhatsAppConnectionEstablishedEvent evento de conexão WhatsApp estabelecida
type WhatsAppConnectionEstablishedEvent struct {
	TenantID      uuid.UUID `json:"tenant_id"`
	PhoneNumber   string    `json:"phone_number"`
	InstanceID    string    `json:"instance_id"`
	ConnectedAt   time.Time `json:"connected_at"`
}

// WhatsAppConnectionLostEvent evento de conexão WhatsApp perdida
type WhatsAppConnectionLostEvent struct {
	TenantID       uuid.UUID `json:"tenant_id"`
	PhoneNumber    string    `json:"phone_number"`
	InstanceID     string    `json:"instance_id"`
	Reason         string    `json:"reason"`
	DisconnectedAt time.Time `json:"disconnected_at"`
}

// EmailDeliveryStatusEvent evento de status de entrega de email
type EmailDeliveryStatusEvent struct {
	NotificationID uuid.UUID `json:"notification_id"`
	MessageID      string    `json:"message_id"`
	RecipientEmail string    `json:"recipient_email"`
	Status         string    `json:"status"` // delivered, bounced, complained, etc.
	TenantID       uuid.UUID `json:"tenant_id"`
	Timestamp      time.Time `json:"timestamp"`
	Details        string    `json:"details,omitempty"`
}

// TelegramBotRegisteredEvent evento de bot Telegram registrado
type TelegramBotRegisteredEvent struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	BotToken    string    `json:"bot_token"`
	BotUsername string    `json:"bot_username"`
	RegisteredAt time.Time `json:"registered_at"`
}

// NotificationPreferencesUpdatedEvent evento de preferências de notificação atualizadas
type NotificationPreferencesUpdatedEvent struct {
	TenantID    uuid.UUID                        `json:"tenant_id"`
	UserID      uuid.UUID                        `json:"user_id"`
	Preferences map[NotificationType][]NotificationChannel `json:"preferences"`
	UpdatedAt   time.Time                        `json:"updated_at"`
}