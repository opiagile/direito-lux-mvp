package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NotificationRepository interface para repositório de notificações
type NotificationRepository interface {
	// CRUD básico
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*Notification, error)
	Update(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Busca e listagem
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, filters NotificationFilters) ([]*Notification, error)
	FindByRecipient(ctx context.Context, recipientID string, filters NotificationFilters) ([]*Notification, error)
	FindByStatus(ctx context.Context, status NotificationStatus, limit int) ([]*Notification, error)
	FindPending(ctx context.Context, limit int) ([]*Notification, error)
	FindScheduled(ctx context.Context, scheduledBefore time.Time, limit int) ([]*Notification, error)
	FindExpired(ctx context.Context, limit int) ([]*Notification, error)
	FindForRetry(ctx context.Context, limit int) ([]*Notification, error)

	// Estatísticas
	CountByTenant(ctx context.Context, tenantID uuid.UUID, filters NotificationFilters) (int64, error)
	CountByStatus(ctx context.Context, tenantID uuid.UUID, status NotificationStatus) (int64, error)
	CountByChannel(ctx context.Context, tenantID uuid.UUID, channel NotificationChannel) (int64, error)
	GetStatistics(ctx context.Context, tenantID uuid.UUID, period time.Duration) (*NotificationStatistics, error)

	// Operações em lote
	CreateBatch(ctx context.Context, notifications []*Notification) error
	UpdateBatch(ctx context.Context, notifications []*Notification) error
	DeleteExpired(ctx context.Context, expiredBefore time.Time) (int64, error)
}

// NotificationTemplateRepository interface para repositório de templates
type NotificationTemplateRepository interface {
	// CRUD básico
	Create(ctx context.Context, template *NotificationTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*NotificationTemplate, error)
	Update(ctx context.Context, template *NotificationTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Busca e listagem
	FindByTenantID(ctx context.Context, tenantID *uuid.UUID, filters TemplateFilters) ([]*NotificationTemplate, error)
	FindByType(ctx context.Context, notificationType NotificationType, channel NotificationChannel, tenantID *uuid.UUID) ([]*NotificationTemplate, error)
	FindSystemTemplates(ctx context.Context, filters TemplateFilters) ([]*NotificationTemplate, error)
	FindActiveTemplates(ctx context.Context, tenantID *uuid.UUID) ([]*NotificationTemplate, error)

	// Operações específicas
	FindByTypeAndChannel(ctx context.Context, notificationType NotificationType, channel NotificationChannel, tenantID uuid.UUID) (*NotificationTemplate, error)
	GetDefaultTemplate(ctx context.Context, notificationType NotificationType, channel NotificationChannel) (*NotificationTemplate, error)
	
	// Contadores
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// NotificationFilters filtros para busca de notificações
type NotificationFilters struct {
	Type      *NotificationType     `json:"type,omitempty"`
	Channel   *NotificationChannel  `json:"channel,omitempty"`
	Priority  *NotificationPriority `json:"priority,omitempty"`
	Status    *NotificationStatus   `json:"status,omitempty"`
	UserID    *uuid.UUID           `json:"user_id,omitempty"`
	ProcessID *uuid.UUID           `json:"process_id,omitempty"`
	CreatedAfter  *time.Time       `json:"created_after,omitempty"`
	CreatedBefore *time.Time       `json:"created_before,omitempty"`
	SentAfter     *time.Time       `json:"sent_after,omitempty"`
	SentBefore    *time.Time       `json:"sent_before,omitempty"`
	Limit         int              `json:"limit,omitempty"`
	Offset        int              `json:"offset,omitempty"`
	OrderBy       string           `json:"order_by,omitempty"`
	OrderDir      string           `json:"order_dir,omitempty"`
}

// TemplateFilters filtros para busca de templates
type TemplateFilters struct {
	Type     *NotificationType    `json:"type,omitempty"`
	Channel  *NotificationChannel `json:"channel,omitempty"`
	Status   *TemplateStatus      `json:"status,omitempty"`
	IsSystem *bool                `json:"is_system,omitempty"`
	Search   string               `json:"search,omitempty"`
	Limit    int                  `json:"limit,omitempty"`
	Offset   int                  `json:"offset,omitempty"`
	OrderBy  string               `json:"order_by,omitempty"`
	OrderDir string               `json:"order_dir,omitempty"`
}

// NotificationStatistics estatísticas de notificações
type NotificationStatistics struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	Period       string    `json:"period"`
	TotalSent    int64     `json:"total_sent"`
	TotalFailed  int64     `json:"total_failed"`
	TotalPending int64     `json:"total_pending"`
	
	// Por canal
	WhatsAppSent int64 `json:"whatsapp_sent"`
	EmailSent    int64 `json:"email_sent"`
	TelegramSent int64 `json:"telegram_sent"`
	PushSent     int64 `json:"push_sent"`
	SMSSent      int64 `json:"sms_sent"`
	
	// Por tipo
	ProcessUpdates     int64 `json:"process_updates"`
	MovementAlerts     int64 `json:"movement_alerts"`
	DeadlineReminders  int64 `json:"deadline_reminders"`
	SystemAlerts       int64 `json:"system_alerts"`
	
	// Métricas de qualidade
	DeliveryRate    float64 `json:"delivery_rate"`    // Taxa de entrega (%)
	ErrorRate       float64 `json:"error_rate"`       // Taxa de erro (%)
	RetryRate       float64 `json:"retry_rate"`       // Taxa de reenvio (%)
	AverageRetries  float64 `json:"average_retries"`  // Média de tentativas
	
	// Tempos
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	MedianProcessingTime  time.Duration `json:"median_processing_time"`
	
	// Horários de pico
	PeakHour   int `json:"peak_hour"`   // Hora com mais notificações (0-23)
	PeakDay    int `json:"peak_day"`    // Dia da semana com mais notificações (0-6)
	
	GeneratedAt time.Time `json:"generated_at"`
}

// NotificationPreference preferências de notificação do usuário
type NotificationPreference struct {
	ID        uuid.UUID                                  `json:"id" db:"id"`
	TenantID  uuid.UUID                                  `json:"tenant_id" db:"tenant_id"`
	UserID    uuid.UUID                                  `json:"user_id" db:"user_id"`
	Type      NotificationType                           `json:"type" db:"type"`
	Channels  []NotificationChannel                      `json:"channels" db:"channels"`
	Enabled   bool                                       `json:"enabled" db:"enabled"`
	CreatedAt time.Time                                  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time                                  `json:"updated_at" db:"updated_at"`
}

// NotificationPreferenceRepository interface para repositório de preferências
type NotificationPreferenceRepository interface {
	// CRUD básico
	Create(ctx context.Context, preference *NotificationPreference) error
	GetByID(ctx context.Context, id uuid.UUID) (*NotificationPreference, error)
	Update(ctx context.Context, preference *NotificationPreference) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Busca específica
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*NotificationPreference, error)
	FindByUserAndType(ctx context.Context, userID uuid.UUID, notificationType NotificationType) (*NotificationPreference, error)
	
	// Operações em lote
	UpsertUserPreferences(ctx context.Context, userID uuid.UUID, preferences []*NotificationPreference) error
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (map[NotificationType][]NotificationChannel, error)
}

// NotificationQueue interface para fila de notificações
type NotificationQueue interface {
	// Enfileirar
	Enqueue(ctx context.Context, notification *Notification) error
	EnqueueBatch(ctx context.Context, notifications []*Notification) error
	EnqueueScheduled(ctx context.Context, notification *Notification, scheduledAt time.Time) error

	// Desenfileirar
	Dequeue(ctx context.Context, channel NotificationChannel) (*Notification, error)
	DequeueBatch(ctx context.Context, channel NotificationChannel, limit int) ([]*Notification, error)
	DequeueByPriority(ctx context.Context, channel NotificationChannel) (*Notification, error)

	// Status da fila
	QueueSize(ctx context.Context, channel NotificationChannel) (int64, error)
	QueueSizeByPriority(ctx context.Context, channel NotificationChannel, priority NotificationPriority) (int64, error)
	
	// Limpar fila
	Clear(ctx context.Context, channel NotificationChannel) error
	
	// Health check
	Health(ctx context.Context) error
}

// NotificationProvider interface para provedores de notificação (WhatsApp, Email, etc.)
type NotificationProvider interface {
	// Envio
	Send(ctx context.Context, notification *Notification) error
	SendBatch(ctx context.Context, notifications []*Notification) error
	
	// Configuração
	Configure(ctx context.Context, config map[string]interface{}) error
	ValidateConfiguration(config map[string]interface{}) error
	
	// Status
	IsHealthy(ctx context.Context) bool
	GetChannel() NotificationChannel
	
	// Capacidades
	SupportsHTML() bool
	SupportsAttachments() bool
	SupportsTemplates() bool
	GetMaxContentLength() int
	GetRateLimit() int // mensagens por minuto
}