package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReportType define o tipo de relatório
type ReportType string

const (
	ReportTypeExecutiveSummary ReportType = "executive_summary"
	ReportTypeProcessAnalysis  ReportType = "process_analysis"
	ReportTypeProductivity     ReportType = "productivity"
	ReportTypeFinancial        ReportType = "financial"
	ReportTypeLegalTimeline    ReportType = "legal_timeline"
	ReportTypeJurisprudence    ReportType = "jurisprudence_analysis"
	ReportTypeCustom           ReportType = "custom"
)

// ReportFormat define o formato de saída do relatório
type ReportFormat string

const (
	ReportFormatPDF   ReportFormat = "pdf"
	ReportFormatExcel ReportFormat = "excel"
	ReportFormatCSV   ReportFormat = "csv"
	ReportFormatJSON  ReportFormat = "json"
	ReportFormatHTML  ReportFormat = "html"
)

// ReportStatus define o status do relatório
type ReportStatus string

const (
	ReportStatusPending    ReportStatus = "pending"
	ReportStatusProcessing ReportStatus = "processing"
	ReportStatusCompleted  ReportStatus = "completed"
	ReportStatusFailed     ReportStatus = "failed"
	ReportStatusCancelled  ReportStatus = "cancelled"
)

// ReportFrequency define a frequência de geração
type ReportFrequency string

const (
	ReportFrequencyOnce    ReportFrequency = "once"
	ReportFrequencyDaily   ReportFrequency = "daily"
	ReportFrequencyWeekly  ReportFrequency = "weekly"
	ReportFrequencyMonthly ReportFrequency = "monthly"
)

// Report representa um relatório no sistema
type Report struct {
	ID               uuid.UUID                      `json:"id" db:"id"`
	TenantID         uuid.UUID                      `json:"tenant_id" db:"tenant_id"`
	UserID           uuid.UUID                      `json:"user_id" db:"user_id"`
	Type             ReportType                     `json:"type" db:"type"`
	Title            string                         `json:"title" db:"title"`
	Description      string                         `json:"description,omitempty" db:"description"`
	Format           ReportFormat                   `json:"format" db:"format"`
	Status           ReportStatus                   `json:"status" db:"status"`
	Parameters       map[string]interface{}         `json:"parameters,omitempty" db:"parameters"`
	Filters          map[string]interface{}         `json:"filters,omitempty" db:"filters"`
	FileURL          *string                        `json:"file_url,omitempty" db:"file_url"`
	FileSize         *int64                         `json:"file_size,omitempty" db:"file_size"`
	ProcessingTime   *int64                         `json:"processing_time,omitempty" db:"processing_time"`
	ErrorMessage     *string                        `json:"error_message,omitempty" db:"error_message"`
	Metadata         map[string]interface{}         `json:"metadata,omitempty" db:"metadata"`
	ScheduleID       *uuid.UUID                     `json:"schedule_id,omitempty" db:"schedule_id"`
	CreatedAt        time.Time                      `json:"created_at" db:"created_at"`
	StartedAt        *time.Time                     `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time                     `json:"completed_at,omitempty" db:"completed_at"`
	ExpiresAt        *time.Time                     `json:"expires_at,omitempty" db:"expires_at"`
}

// ReportSchedule representa um agendamento de relatório
type ReportSchedule struct {
	ID               uuid.UUID                      `json:"id" db:"id"`
	TenantID         uuid.UUID                      `json:"tenant_id" db:"tenant_id"`
	UserID           uuid.UUID                      `json:"user_id" db:"user_id"`
	ReportType       ReportType                     `json:"report_type" db:"report_type"`
	Title            string                         `json:"title" db:"title"`
	Description      string                         `json:"description,omitempty" db:"description"`
	Format           ReportFormat                   `json:"format" db:"format"`
	Frequency        ReportFrequency                `json:"frequency" db:"frequency"`
	CronExpression   *string                        `json:"cron_expression,omitempty" db:"cron_expression"`
	Parameters       map[string]interface{}         `json:"parameters,omitempty" db:"parameters"`
	Filters          map[string]interface{}         `json:"filters,omitempty" db:"filters"`
	Recipients       []string                       `json:"recipients,omitempty" db:"recipients"`
	IsActive         bool                           `json:"is_active" db:"is_active"`
	NextRunAt        *time.Time                     `json:"next_run_at,omitempty" db:"next_run_at"`
	LastRunAt        *time.Time                     `json:"last_run_at,omitempty" db:"last_run_at"`
	LastReportID     *uuid.UUID                     `json:"last_report_id,omitempty" db:"last_report_id"`
	CreatedAt        time.Time                      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                      `json:"updated_at" db:"updated_at"`
}

// Dashboard representa um dashboard executivo
type Dashboard struct {
	ID               uuid.UUID                      `json:"id" db:"id"`
	TenantID         uuid.UUID                      `json:"tenant_id" db:"tenant_id"`
	Title            string                         `json:"title" db:"title"`
	Description      string                         `json:"description,omitempty" db:"description"`
	IsPublic         bool                           `json:"is_public" db:"is_public"`
	IsDefault        bool                           `json:"is_default" db:"is_default"`
	Layout           map[string]interface{}         `json:"layout,omitempty" db:"layout"`
	Widgets          []DashboardWidget              `json:"widgets" db:"widgets"`
	RefreshInterval  *int                           `json:"refresh_interval,omitempty" db:"refresh_interval"`
	CreatedBy        uuid.UUID                      `json:"created_by" db:"created_by"`
	UpdatedBy        uuid.UUID                      `json:"updated_by" db:"updated_by"`
	CreatedAt        time.Time                      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                      `json:"updated_at" db:"updated_at"`
}

// DashboardWidget representa um widget no dashboard
type DashboardWidget struct {
	ID               uuid.UUID                      `json:"id" db:"id"`
	DashboardID      uuid.UUID                      `json:"dashboard_id" db:"dashboard_id"`
	Type             string                         `json:"type" db:"type"`
	Title            string                         `json:"title" db:"title"`
	DataSource       string                         `json:"data_source" db:"data_source"`
	ChartType        string                         `json:"chart_type,omitempty" db:"chart_type"`
	Parameters       map[string]interface{}         `json:"parameters,omitempty" db:"parameters"`
	Filters          map[string]interface{}         `json:"filters,omitempty" db:"filters"`
	Position         WidgetPosition                 `json:"position" db:"position"`
	Size             WidgetSize                     `json:"size" db:"size"`
	RefreshInterval  *int                           `json:"refresh_interval,omitempty" db:"refresh_interval"`
	IsVisible        bool                           `json:"is_visible" db:"is_visible"`
	Order            int                            `json:"order" db:"order"`
}

// WidgetPosition define a posição do widget
type WidgetPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// WidgetSize define o tamanho do widget
type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// KPI representa um indicador chave de performance
type KPI struct {
	ID               uuid.UUID                      `json:"id" db:"id"`
	TenantID         uuid.UUID                      `json:"tenant_id" db:"tenant_id"`
	Name             string                         `json:"name" db:"name"`
	DisplayName      string                         `json:"display_name" db:"display_name"`
	Description      string                         `json:"description,omitempty" db:"description"`
	Category         string                         `json:"category" db:"category"`
	Unit             string                         `json:"unit,omitempty" db:"unit"`
	CurrentValue     float64                        `json:"current_value" db:"current_value"`
	PreviousValue    *float64                       `json:"previous_value,omitempty" db:"previous_value"`
	Target           *float64                       `json:"target,omitempty" db:"target"`
	Trend            string                         `json:"trend,omitempty" db:"trend"`
	TrendPercentage  *float64                       `json:"trend_percentage,omitempty" db:"trend_percentage"`
	LastCalculatedAt time.Time                      `json:"last_calculated_at" db:"last_calculated_at"`
	UpdatedAt        time.Time                      `json:"updated_at" db:"updated_at"`
}

// ReportTemplate representa um template de relatório personalizado
type ReportTemplate struct {
	ID               uuid.UUID                      `json:"id" db:"id"`
	TenantID         uuid.UUID                      `json:"tenant_id" db:"tenant_id"`
	Name             string                         `json:"name" db:"name"`
	Description      string                         `json:"description,omitempty" db:"description"`
	Type             ReportType                     `json:"type" db:"type"`
	Template         string                         `json:"template" db:"template"`
	DefaultFormat    ReportFormat                   `json:"default_format" db:"default_format"`
	Variables        []TemplateVariable             `json:"variables,omitempty" db:"variables"`
	IsPublic         bool                           `json:"is_public" db:"is_public"`
	CreatedBy        uuid.UUID                      `json:"created_by" db:"created_by"`
	CreatedAt        time.Time                      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                      `json:"updated_at" db:"updated_at"`
}

// TemplateVariable representa uma variável no template
type TemplateVariable struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Required     bool   `json:"required"`
}