package queries

import (
	"time"
	"github.com/direito-lux/process-service/internal/domain"
)

// Query interface base para queries
type Query interface {
	Validate() error
}

// ProcessQuery consulta básica de processo
type ProcessQuery struct {
	ID       string `json:"id" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required"`
}

func (q *ProcessQuery) Validate() error {
	// Validação usando tags seria implementada aqui
	return nil
}

// ProcessListQuery consulta lista de processos
type ProcessListQuery struct {
	TenantID  string                   `json:"tenant_id" validate:"required"`
	ClientID  string                   `json:"client_id"`
	Status    []domain.ProcessStatus   `json:"status"`
	Stage     []domain.ProcessStage    `json:"stage"`
	CourtID   string                   `json:"court_id"`
	JudgeID   string                   `json:"judge_id"`
	Tags      []string                 `json:"tags"`
	DateFrom  *time.Time               `json:"date_from"`
	DateTo    *time.Time               `json:"date_to"`
	Search    string                   `json:"search"`
	Page      int                      `json:"page" validate:"min=1"`
	PageSize  int                      `json:"page_size" validate:"min=1,max=100"`
	SortBy    string                   `json:"sort_by"`
	SortOrder string                   `json:"sort_order"`
}

func (q *ProcessListQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	if q.SortBy == "" {
		q.SortBy = "updated_at"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	return nil
}

// ProcessStatsQuery consulta estatísticas de processos
type ProcessStatsQuery struct {
	TenantID string     `json:"tenant_id" validate:"required"`
	ClientID string     `json:"client_id"`
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}

func (q *ProcessStatsQuery) Validate() error {
	return nil
}

// ProcessMonitoringQuery consulta processos para monitoramento
type ProcessMonitoringQuery struct {
	TenantID         string `json:"tenant_id" validate:"required"`
	OnlyNeedingSync  bool   `json:"only_needing_sync"`
	OnlyWithAlerts   bool   `json:"only_with_alerts"`
	Page             int    `json:"page" validate:"min=1"`
	PageSize         int    `json:"page_size" validate:"min=1,max=100"`
}

func (q *ProcessMonitoringQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	return nil
}

// MovementQuery consulta básica de movimentação
type MovementQuery struct {
	ID       string `json:"id" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required"`
}

func (q *MovementQuery) Validate() error {
	return nil
}

// MovementListQuery consulta lista de movimentações
type MovementListQuery struct {
	TenantID          string               `json:"tenant_id" validate:"required"`
	ProcessID         string               `json:"process_id"`
	Type              []domain.MovementType `json:"type"`
	IsImportant       *bool                `json:"is_important"`
	IsPublic          *bool                `json:"is_public"`
	HasNotification   *bool                `json:"has_notification"`
	DateFrom          *time.Time           `json:"date_from"`
	DateTo            *time.Time           `json:"date_to"`
	Judge             string               `json:"judge"`
	Tags              []string             `json:"tags"`
	Search            string               `json:"search"`
	Page              int                  `json:"page" validate:"min=1"`
	PageSize          int                  `json:"page_size" validate:"min=1,max=100"`
	SortBy            string               `json:"sort_by"`
	SortOrder         string               `json:"sort_order"`
	IncludeContent    bool                 `json:"include_content"`
	IncludeAttachments bool                `json:"include_attachments"`
}

func (q *MovementListQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	if q.SortBy == "" {
		q.SortBy = "date"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	return nil
}

// MovementSearchQuery consulta de busca textual em movimentações
type MovementSearchQuery struct {
	TenantID    string     `json:"tenant_id" validate:"required"`
	Query       string     `json:"query" validate:"required,min=3"`
	ProcessID   string     `json:"process_id"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	OnlyPublic  bool       `json:"only_public"`
	Page        int        `json:"page" validate:"min=1"`
	PageSize    int        `json:"page_size" validate:"min=1,max=100"`
	Highlights  bool       `json:"highlights"`
}

func (q *MovementSearchQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	return nil
}

// MovementStatsQuery consulta estatísticas de movimentações
type MovementStatsQuery struct {
	TenantID  string     `json:"tenant_id" validate:"required"`
	ProcessID string     `json:"process_id"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
}

func (q *MovementStatsQuery) Validate() error {
	return nil
}

// PartyQuery consulta básica de parte
type PartyQuery struct {
	ID       string `json:"id" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required"`
}

func (q *PartyQuery) Validate() error {
	return nil
}

// PartyListQuery consulta lista de partes
type PartyListQuery struct {
	TenantID     string              `json:"tenant_id" validate:"required"`
	ProcessID    string              `json:"process_id"`
	Type         []domain.PartyType  `json:"type"`
	Role         []domain.PartyRole  `json:"role"`
	IsActive     *bool               `json:"is_active"`
	DocumentType string              `json:"document_type"`
	Search       string              `json:"search"`
	Page         int                 `json:"page" validate:"min=1"`
	PageSize     int                 `json:"page_size" validate:"min=1,max=100"`
}

func (q *PartyListQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	return nil
}

// PartySearchQuery consulta busca de partes por documento
type PartySearchQuery struct {
	TenantID string `json:"tenant_id" validate:"required"`
	Document string `json:"document" validate:"required"`
}

func (q *PartySearchQuery) Validate() error {
	return nil
}

// DashboardQuery consulta para dashboard
type DashboardQuery struct {
	TenantID string `json:"tenant_id" validate:"required"`
	ClientID string `json:"client_id"`
	Period   string `json:"period"` // week, month, quarter, year
}

func (q *DashboardQuery) Validate() error {
	if q.Period == "" {
		q.Period = "month"
	}
	return nil
}

// ProcessTimelineQuery consulta timeline do processo
type ProcessTimelineQuery struct {
	ProcessID string     `json:"process_id" validate:"required"`
	TenantID  string     `json:"tenant_id" validate:"required"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	EventTypes []string  `json:"event_types"` // movements, parties, status_changes
	Page      int        `json:"page" validate:"min=1"`
	PageSize  int        `json:"page_size" validate:"min=1,max=100"`
}

func (q *ProcessTimelineQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 50
	}
	if len(q.EventTypes) == 0 {
		q.EventTypes = []string{"movements", "parties", "status_changes"}
	}
	return nil
}

// CalendarQuery consulta eventos de calendário
type CalendarQuery struct {
	TenantID  string     `json:"tenant_id" validate:"required"`
	ClientID  string     `json:"client_id"`
	DateFrom  time.Time  `json:"date_from" validate:"required"`
	DateTo    time.Time  `json:"date_to" validate:"required"`
	EventTypes []string  `json:"event_types"` // hearings, deadlines, notifications
}

func (q *CalendarQuery) Validate() error {
	if len(q.EventTypes) == 0 {
		q.EventTypes = []string{"hearings", "deadlines", "notifications"}
	}
	return nil
}

// AlertsQuery consulta alertas e notificações
type AlertsQuery struct {
	TenantID   string `json:"tenant_id" validate:"required"`
	ClientID   string `json:"client_id"`
	OnlyUnread bool   `json:"only_unread"`
	Severity   string `json:"severity"` // low, medium, high, critical
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

func (q *AlertsQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	return nil
}

// ReportsQuery consulta para relatórios
type ReportsQuery struct {
	TenantID   string     `json:"tenant_id" validate:"required"`
	ClientID   string     `json:"client_id"`
	ReportType string     `json:"report_type" validate:"required"` // productivity, movements, parties, timeline
	DateFrom   time.Time  `json:"date_from" validate:"required"`
	DateTo     time.Time  `json:"date_to" validate:"required"`
	Filters    map[string]interface{} `json:"filters"`
	Format     string     `json:"format"` // json, csv, excel, pdf
}

func (q *ReportsQuery) Validate() error {
	if q.Format == "" {
		q.Format = "json"
	}
	return nil
}