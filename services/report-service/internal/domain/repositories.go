package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ReportRepository define as operações de persistência para relatórios
type ReportRepository interface {
	// Report operations
	Create(ctx context.Context, report *Report) error
	Update(ctx context.Context, report *Report) error
	GetByID(ctx context.Context, id uuid.UUID) (*Report, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, filters ReportFilters) ([]*Report, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, filters ReportFilters) ([]*Report, error)
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Batch operations
	CreateBatch(ctx context.Context, reports []*Report) error
	GetExpiredReports(ctx context.Context, before time.Time) ([]*Report, error)
	DeleteExpired(ctx context.Context, before time.Time) (int64, error)
	
	// Statistics
	GetStatsByTenant(ctx context.Context, tenantID uuid.UUID, period time.Duration) (*ReportStatistics, error)
	GetProcessingMetrics(ctx context.Context, tenantID uuid.UUID) (*ProcessingMetrics, error)
}

// ReportScheduleRepository define as operações para agendamentos
type ReportScheduleRepository interface {
	Create(ctx context.Context, schedule *ReportSchedule) error
	Update(ctx context.Context, schedule *ReportSchedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*ReportSchedule, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*ReportSchedule, error)
	GetActiveSchedules(ctx context.Context) ([]*ReportSchedule, error)
	GetDueSchedules(ctx context.Context, before time.Time) ([]*ReportSchedule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastRun(ctx context.Context, id uuid.UUID, lastRunAt time.Time, nextRunAt time.Time, reportID *uuid.UUID) error
}

// DashboardRepository define as operações para dashboards
type DashboardRepository interface {
	Create(ctx context.Context, dashboard *Dashboard) error
	Update(ctx context.Context, dashboard *Dashboard) error
	GetByID(ctx context.Context, id uuid.UUID) (*Dashboard, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*Dashboard, error)
	GetPublicDashboards(ctx context.Context, tenantID uuid.UUID) ([]*Dashboard, error)
	GetDefaultDashboard(ctx context.Context, tenantID uuid.UUID) (*Dashboard, error)
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Widget operations
	AddWidget(ctx context.Context, widget *DashboardWidget) error
	UpdateWidget(ctx context.Context, widget *DashboardWidget) error
	RemoveWidget(ctx context.Context, dashboardID, widgetID uuid.UUID) error
	GetWidgets(ctx context.Context, dashboardID uuid.UUID) ([]DashboardWidget, error)
}

// KPIRepository define as operações para KPIs
type KPIRepository interface {
	Create(ctx context.Context, kpi *KPI) error
	Update(ctx context.Context, kpi *KPI) error
	GetByID(ctx context.Context, id uuid.UUID) (*KPI, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*KPI, error)
	GetByCategory(ctx context.Context, tenantID uuid.UUID, category string) ([]*KPI, error)
	UpdateValue(ctx context.Context, id uuid.UUID, value float64) error
	GetHistoricalData(ctx context.Context, kpiID uuid.UUID, from, to time.Time) ([]KPIHistoryPoint, error)
}

// ReportTemplateRepository define as operações para templates
type ReportTemplateRepository interface {
	Create(ctx context.Context, template *ReportTemplate) error
	Update(ctx context.Context, template *ReportTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*ReportTemplate, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*ReportTemplate, error)
	GetByType(ctx context.Context, tenantID uuid.UUID, reportType ReportType) ([]*ReportTemplate, error)
	GetPublicTemplates(ctx context.Context) ([]*ReportTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ReportGenerator define a interface para geração de relatórios
type ReportGenerator interface {
	GeneratePDF(ctx context.Context, report *Report, data interface{}) ([]byte, error)
	GenerateExcel(ctx context.Context, report *Report, data interface{}) ([]byte, error)
	GenerateCSV(ctx context.Context, report *Report, data interface{}) ([]byte, error)
	GenerateHTML(ctx context.Context, report *Report, data interface{}) ([]byte, error)
}

// DataCollector define a interface para coleta de dados
type DataCollector interface {
	CollectProcessData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error)
	CollectProductivityData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error)
	CollectFinancialData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error)
	CollectJurisprudenceData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error)
	CollectDashboardData(ctx context.Context, widget *DashboardWidget) (interface{}, error)
	CalculateKPIs(ctx context.Context, tenantID uuid.UUID) ([]*KPI, error)
}

// ReportFilters define filtros para busca de relatórios
type ReportFilters struct {
	Type      *ReportType   `json:"type,omitempty"`
	Status    *ReportStatus `json:"status,omitempty"`
	Format    *ReportFormat `json:"format,omitempty"`
	StartDate *time.Time    `json:"start_date,omitempty"`
	EndDate   *time.Time    `json:"end_date,omitempty"`
	Limit     int           `json:"limit,omitempty"`
	Offset    int           `json:"offset,omitempty"`
	OrderBy   string        `json:"order_by,omitempty"`
}

// ReportStatistics representa estatísticas de relatórios
type ReportStatistics struct {
	TotalReports         int64                    `json:"total_reports"`
	ReportsByType        map[string]int64         `json:"reports_by_type"`
	ReportsByFormat      map[string]int64         `json:"reports_by_format"`
	AverageProcessingTime float64                 `json:"average_processing_time"`
	SuccessRate          float64                  `json:"success_rate"`
}

// ProcessingMetrics representa métricas de processamento
type ProcessingMetrics struct {
	TotalProcessed       int64   `json:"total_processed"`
	TotalFailed          int64   `json:"total_failed"`
	AverageTime          float64 `json:"average_time"`
	MedianTime           float64 `json:"median_time"`
	P95Time              float64 `json:"p95_time"`
	LastProcessedAt      *time.Time `json:"last_processed_at,omitempty"`
}

// KPIHistoryPoint representa um ponto histórico de KPI
type KPIHistoryPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}