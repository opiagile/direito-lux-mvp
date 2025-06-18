package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NotificationType tipo de notificação
type NotificationType string

const (
	NotificationTypeProcessUpdate     NotificationType = "process_update"
	NotificationTypeMovementAlert     NotificationType = "movement_alert"
	NotificationTypeDeadlineReminder  NotificationType = "deadline_reminder"
	NotificationTypeTrialExpiring     NotificationType = "trial_expiring"
	NotificationTypeSubscriptionDue   NotificationType = "subscription_due"
	NotificationTypeSystemAlert       NotificationType = "system_alert"
	NotificationTypeWelcome          NotificationType = "welcome"
	NotificationTypePasswordReset    NotificationType = "password_reset"
)

// NotificationChannel canal de notificação
type NotificationChannel string

const (
	NotificationChannelWhatsApp NotificationChannel = "whatsapp"
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelTelegram NotificationChannel = "telegram"
	NotificationChannelPush     NotificationChannel = "push"
	NotificationChannelSMS      NotificationChannel = "sms"
)

// NotificationStatus status da notificação
type NotificationStatus string

const (
	NotificationStatusPending    NotificationStatus = "pending"
	NotificationStatusScheduled  NotificationStatus = "scheduled"
	NotificationStatusProcessing NotificationStatus = "processing"
	NotificationStatusSent       NotificationStatus = "sent"
	NotificationStatusDelivered  NotificationStatus = "delivered"
	NotificationStatusFailed     NotificationStatus = "failed"
	NotificationStatusExpired    NotificationStatus = "expired"
	NotificationStatusCancelled  NotificationStatus = "cancelled"
)

// NotificationPriority prioridade da notificação
type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityNormal   NotificationPriority = "normal"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"
)

// Notification representa uma notificação
type Notification struct {
	ID               uuid.UUID            `json:"id" db:"id"`
	Type             NotificationType     `json:"type" db:"type"`
	Channel          NotificationChannel  `json:"channel" db:"channel"`
	Priority         NotificationPriority `json:"priority" db:"priority"`
	Status           NotificationStatus   `json:"status" db:"status"`
	Subject          string               `json:"subject" db:"subject"`
	Content          string               `json:"content" db:"content"`
	ContentHTML      *string              `json:"content_html,omitempty" db:"content_html"`
	RecipientID      string               `json:"recipient_id" db:"recipient_id"`
	RecipientType    string               `json:"recipient_type" db:"recipient_type"`
	RecipientContact string               `json:"recipient_contact" db:"recipient_contact"`
	TenantID         uuid.UUID            `json:"tenant_id" db:"tenant_id"`
	UserID           *uuid.UUID           `json:"user_id,omitempty" db:"user_id"`
	ProcessID        *uuid.UUID           `json:"process_id,omitempty" db:"process_id"`
	TemplateID       *uuid.UUID           `json:"template_id,omitempty" db:"template_id"`
	Variables        map[string]interface{} `json:"variables,omitempty" db:"variables"`
	Metadata         map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	RetryCount       int                  `json:"retry_count" db:"retry_count"`
	MaxRetries       int                  `json:"max_retries" db:"max_retries"`
	NextRetryAt      *time.Time           `json:"next_retry_at,omitempty" db:"next_retry_at"`
	ScheduledAt      *time.Time           `json:"scheduled_at,omitempty" db:"scheduled_at"`
	SentAt           *time.Time           `json:"sent_at,omitempty" db:"sent_at"`
	DeliveredAt      *time.Time           `json:"delivered_at,omitempty" db:"delivered_at"`
	FailedAt         *time.Time           `json:"failed_at,omitempty" db:"failed_at"`
	ErrorMessage     *string              `json:"error_message,omitempty" db:"error_message"`
	ExternalID       *string              `json:"external_id,omitempty" db:"external_id"`
	ExternalStatus   *string              `json:"external_status,omitempty" db:"external_status"`
	ProcessingStartedAt  *time.Time       `json:"processing_started_at,omitempty" db:"processing_started_at"`
	ProcessingFinishedAt *time.Time       `json:"processing_finished_at,omitempty" db:"processing_finished_at"`
	Tags             []string             `json:"tags,omitempty" db:"tags"`
	ExpiresAt        *time.Time           `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt        time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at" db:"updated_at"`
}

// NewNotification cria nova notificação
func NewNotification(
	notificationType NotificationType,
	channel NotificationChannel,
	priority NotificationPriority,
	subject, content string,
	recipientID, recipientType, recipientContact string,
	tenantID uuid.UUID,
) (*Notification, error) {
	// Validações
	if err := ValidateNotificationType(notificationType); err != nil {
		return nil, err
	}

	if err := ValidateChannel(channel); err != nil {
		return nil, err
	}

	if err := ValidatePriority(priority); err != nil {
		return nil, err
	}

	if subject == "" {
		return nil, errors.New("subject é obrigatório")
	}

	if content == "" {
		return nil, errors.New("content é obrigatório")
	}

	if recipientID == "" {
		return nil, errors.New("recipient_id é obrigatório")
	}

	if recipientType == "" {
		return nil, errors.New("recipient_type é obrigatório")
	}

	if recipientContact == "" {
		return nil, errors.New("recipient_contact é obrigatório")
	}

	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}

	now := time.Now()
	
	return &Notification{
		ID:               uuid.New(),
		Type:             notificationType,
		Channel:          channel,
		Priority:         priority,
		Status:           NotificationStatusPending,
		Subject:          subject,
		Content:          content,
		RecipientID:      recipientID,
		RecipientType:    recipientType,
		RecipientContact: recipientContact,
		TenantID:         tenantID,
		Variables:        make(map[string]interface{}),
		Metadata:         make(map[string]interface{}),
		RetryCount:       0,
		MaxRetries:       3,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// MarkAsProcessing marca notificação como processando
func (n *Notification) MarkAsProcessing() error {
	if n.Status != NotificationStatusPending {
		return fmt.Errorf("não é possível processar notificação com status %s", n.Status)
	}

	n.Status = NotificationStatusProcessing
	n.UpdatedAt = time.Now()
	
	return nil
}

// MarkAsSent marca notificação como enviada
func (n *Notification) MarkAsSent() error {
	if n.Status != NotificationStatusProcessing {
		return fmt.Errorf("não é possível marcar como enviada notificação com status %s", n.Status)
	}

	now := time.Now()
	n.Status = NotificationStatusSent
	n.SentAt = &now
	n.UpdatedAt = now
	
	return nil
}

// MarkAsFailed marca notificação como falhada
func (n *Notification) MarkAsFailed(errorMessage string) error {
	if n.Status != NotificationStatusProcessing {
		return fmt.Errorf("não é possível marcar como falhada notificação com status %s", n.Status)
	}

	now := time.Now()
	n.Status = NotificationStatusFailed
	n.FailedAt = &now
	n.ErrorMessage = &errorMessage
	n.RetryCount++
	n.UpdatedAt = now
	
	return nil
}

// CanRetry verifica se a notificação pode ser reenvida
func (n *Notification) CanRetry() bool {
	return n.Status == NotificationStatusFailed && n.RetryCount < n.MaxRetries
}

// ResetForRetry reseta notificação para reenvio
func (n *Notification) ResetForRetry() error {
	if !n.CanRetry() {
		return errors.New("notificação não pode ser reenviada")
	}

	n.Status = NotificationStatusPending
	n.FailedAt = nil
	n.ErrorMessage = nil
	n.UpdatedAt = time.Now()
	
	return nil
}

// Cancel cancela notificação
func (n *Notification) Cancel() error {
	if n.Status == NotificationStatusSent {
		return errors.New("não é possível cancelar notificação já enviada")
	}

	n.Status = NotificationStatusCancelled
	n.UpdatedAt = time.Now()
	
	return nil
}

// IsExpired verifica se a notificação expirou
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// SetScheduledAt agenda notificação
func (n *Notification) SetScheduledAt(scheduledAt time.Time) {
	n.ScheduledAt = &scheduledAt
	n.UpdatedAt = time.Now()
}

// SetExpiresAt define expiração
func (n *Notification) SetExpiresAt(expiresAt time.Time) {
	n.ExpiresAt = &expiresAt
	n.UpdatedAt = time.Now()
}

// SetTemplate define template
func (n *Notification) SetTemplate(templateID uuid.UUID) {
	n.TemplateID = &templateID
	n.UpdatedAt = time.Now()
}

// SetVariables define variáveis do template
func (n *Notification) SetVariables(variables map[string]interface{}) {
	n.Variables = variables
	n.UpdatedAt = time.Now()
}

// SetMetadata define metadados
func (n *Notification) SetMetadata(metadata map[string]interface{}) {
	n.Metadata = metadata
	n.UpdatedAt = time.Now()
}

// SetContentHTML define conteúdo HTML
func (n *Notification) SetContentHTML(contentHTML string) {
	n.ContentHTML = &contentHTML
	n.UpdatedAt = time.Now()
}

// SetUserID define ID do usuário
func (n *Notification) SetUserID(userID uuid.UUID) {
	n.UserID = &userID
	n.UpdatedAt = time.Now()
}

// SetProcessID define ID do processo
func (n *Notification) SetProcessID(processID uuid.UUID) {
	n.ProcessID = &processID
	n.UpdatedAt = time.Now()
}

// ValidateNotificationType valida tipo de notificação
func ValidateNotificationType(notificationType NotificationType) error {
	switch notificationType {
	case NotificationTypeProcessUpdate,
		 NotificationTypeMovementAlert,
		 NotificationTypeDeadlineReminder,
		 NotificationTypeTrialExpiring,
		 NotificationTypeSubscriptionDue,
		 NotificationTypeSystemAlert,
		 NotificationTypeWelcome,
		 NotificationTypePasswordReset:
		return nil
	default:
		return fmt.Errorf("tipo de notificação inválido: %s", notificationType)
	}
}

// ValidateChannel valida canal de notificação
func ValidateChannel(channel NotificationChannel) error {
	switch channel {
	case NotificationChannelWhatsApp,
		 NotificationChannelEmail,
		 NotificationChannelTelegram,
		 NotificationChannelPush,
		 NotificationChannelSMS:
		return nil
	default:
		return fmt.Errorf("canal de notificação inválido: %s", channel)
	}
}

// ValidatePriority valida prioridade
func ValidatePriority(priority NotificationPriority) error {
	switch priority {
	case NotificationPriorityLow,
		 NotificationPriorityNormal,
		 NotificationPriorityHigh,
		 NotificationPriorityCritical:
		return nil
	default:
		return fmt.Errorf("prioridade inválida: %s", priority)
	}
}