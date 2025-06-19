package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ReportEvent representa um evento base de relatório
type ReportEvent struct {
	EventID     uuid.UUID              `json:"event_id"`
	EventType   string                 `json:"event_type"`
	ReportID    uuid.UUID              `json:"report_id"`
	TenantID    uuid.UUID              `json:"tenant_id"`
	UserID      uuid.UUID              `json:"user_id"`
	OccurredAt  time.Time              `json:"occurred_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ReportRequestedEvent evento quando um relatório é solicitado
type ReportRequestedEvent struct {
	ReportEvent
	ReportType ReportType             `json:"report_type"`
	Format     ReportFormat           `json:"format"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ReportGeneratedEvent evento quando um relatório é gerado com sucesso
type ReportGeneratedEvent struct {
	ReportEvent
	FileURL        string        `json:"file_url"`
	FileSize       int64         `json:"file_size"`
	ProcessingTime int64         `json:"processing_time"`
	Format         ReportFormat  `json:"format"`
}

// ReportFailedEvent evento quando a geração falha
type ReportFailedEvent struct {
	ReportEvent
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code,omitempty"`
}

// ReportScheduledEvent evento quando um relatório é agendado
type ReportScheduledEvent struct {
	ReportEvent
	ScheduleID     uuid.UUID       `json:"schedule_id"`
	Frequency      ReportFrequency `json:"frequency"`
	NextRunAt      time.Time       `json:"next_run_at"`
}

// ReportEmailedEvent evento quando um relatório é enviado por email
type ReportEmailedEvent struct {
	ReportEvent
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
}

// DashboardEvent representa um evento base de dashboard
type DashboardEvent struct {
	EventID     uuid.UUID              `json:"event_id"`
	EventType   string                 `json:"event_type"`
	DashboardID uuid.UUID              `json:"dashboard_id"`
	TenantID    uuid.UUID              `json:"tenant_id"`
	UserID      uuid.UUID              `json:"user_id"`
	OccurredAt  time.Time              `json:"occurred_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DashboardCreatedEvent evento quando um dashboard é criado
type DashboardCreatedEvent struct {
	DashboardEvent
	Title       string `json:"title"`
	IsPublic    bool   `json:"is_public"`
	WidgetCount int    `json:"widget_count"`
}

// DashboardUpdatedEvent evento quando um dashboard é atualizado
type DashboardUpdatedEvent struct {
	DashboardEvent
	Changes map[string]interface{} `json:"changes"`
}

// DashboardSharedEvent evento quando um dashboard é compartilhado
type DashboardSharedEvent struct {
	DashboardEvent
	SharedWith []uuid.UUID `json:"shared_with"`
	Permission string      `json:"permission"`
}

// WidgetAddedEvent evento quando um widget é adicionado
type WidgetAddedEvent struct {
	DashboardEvent
	WidgetID   uuid.UUID `json:"widget_id"`
	WidgetType string    `json:"widget_type"`
	DataSource string    `json:"data_source"`
}

// WidgetRemovedEvent evento quando um widget é removido
type WidgetRemovedEvent struct {
	DashboardEvent
	WidgetID uuid.UUID `json:"widget_id"`
}

// KPIEvent representa um evento base de KPI
type KPIEvent struct {
	EventID    uuid.UUID              `json:"event_id"`
	EventType  string                 `json:"event_type"`
	KPIID      uuid.UUID              `json:"kpi_id"`
	TenantID   uuid.UUID              `json:"tenant_id"`
	OccurredAt time.Time              `json:"occurred_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// KPICalculatedEvent evento quando um KPI é calculado
type KPICalculatedEvent struct {
	KPIEvent
	Name          string   `json:"name"`
	CurrentValue  float64  `json:"current_value"`
	PreviousValue *float64 `json:"previous_value,omitempty"`
	Trend         string   `json:"trend,omitempty"`
}

// KPIAlertEvent evento quando um KPI dispara um alerta
type KPIAlertEvent struct {
	KPIEvent
	AlertType    string  `json:"alert_type"`
	CurrentValue float64 `json:"current_value"`
	Threshold    float64 `json:"threshold"`
	Message      string  `json:"message"`
}

// EventBus interface para publicação de eventos
type EventBus interface {
	Publish(ctx context.Context, event interface{}) error
	Subscribe(ctx context.Context, eventType string, handler EventHandler) error
}

// EventHandler função para processar eventos
type EventHandler func(ctx context.Context, event interface{}) error