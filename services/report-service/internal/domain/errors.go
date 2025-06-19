package domain

import "fmt"

// DomainError representa um erro de domínio
type DomainError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Erros de domínio comuns
var (
	// Report errors
	ErrReportNotFound        = DomainError{Code: "REPORT_NOT_FOUND", Message: "Report not found"}
	ErrReportAlreadyExists   = DomainError{Code: "REPORT_ALREADY_EXISTS", Message: "Report already exists"}
	ErrInvalidReportType     = DomainError{Code: "INVALID_REPORT_TYPE", Message: "Invalid report type"}
	ErrInvalidReportFormat   = DomainError{Code: "INVALID_REPORT_FORMAT", Message: "Invalid report format"}
	ErrReportGenerationFailed = DomainError{Code: "REPORT_GENERATION_FAILED", Message: "Report generation failed"}
	ErrReportExpired         = DomainError{Code: "REPORT_EXPIRED", Message: "Report has expired"}
	ErrReportInProgress      = DomainError{Code: "REPORT_IN_PROGRESS", Message: "Report is already being processed"}
	
	// Schedule errors
	ErrScheduleNotFound      = DomainError{Code: "SCHEDULE_NOT_FOUND", Message: "Schedule not found"}
	ErrInvalidSchedule       = DomainError{Code: "INVALID_SCHEDULE", Message: "Invalid schedule configuration"}
	ErrScheduleInactive      = DomainError{Code: "SCHEDULE_INACTIVE", Message: "Schedule is inactive"}
	ErrInvalidCronExpression = DomainError{Code: "INVALID_CRON_EXPRESSION", Message: "Invalid cron expression"}
	
	// Dashboard errors
	ErrDashboardNotFound     = DomainError{Code: "DASHBOARD_NOT_FOUND", Message: "Dashboard not found"}
	ErrDashboardAlreadyExists = DomainError{Code: "DASHBOARD_ALREADY_EXISTS", Message: "Dashboard already exists"}
	ErrWidgetNotFound        = DomainError{Code: "WIDGET_NOT_FOUND", Message: "Widget not found"}
	ErrInvalidWidgetType     = DomainError{Code: "INVALID_WIDGET_TYPE", Message: "Invalid widget type"}
	ErrMaxWidgetsReached     = DomainError{Code: "MAX_WIDGETS_REACHED", Message: "Maximum number of widgets reached"}
	
	// KPI errors
	ErrKPINotFound           = DomainError{Code: "KPI_NOT_FOUND", Message: "KPI not found"}
	ErrInvalidKPICategory    = DomainError{Code: "INVALID_KPI_CATEGORY", Message: "Invalid KPI category"}
	ErrKPICalculationFailed  = DomainError{Code: "KPI_CALCULATION_FAILED", Message: "KPI calculation failed"}
	
	// Template errors
	ErrTemplateNotFound      = DomainError{Code: "TEMPLATE_NOT_FOUND", Message: "Template not found"}
	ErrInvalidTemplate       = DomainError{Code: "INVALID_TEMPLATE", Message: "Invalid template format"}
	ErrTemplateRenderFailed  = DomainError{Code: "TEMPLATE_RENDER_FAILED", Message: "Template rendering failed"}
	
	// Permission errors
	ErrUnauthorized          = DomainError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrInsufficientPlan      = DomainError{Code: "INSUFFICIENT_PLAN", Message: "Current plan does not support this feature"}
	ErrQuotaExceeded         = DomainError{Code: "QUOTA_EXCEEDED", Message: "Report quota exceeded"}
	
	// Data errors
	ErrNoDataAvailable       = DomainError{Code: "NO_DATA_AVAILABLE", Message: "No data available for the requested period"}
	ErrInvalidDateRange      = DomainError{Code: "INVALID_DATE_RANGE", Message: "Invalid date range"}
	ErrDataCollectionFailed  = DomainError{Code: "DATA_COLLECTION_FAILED", Message: "Failed to collect data"}
	
	// File errors
	ErrFileGenerationFailed  = DomainError{Code: "FILE_GENERATION_FAILED", Message: "File generation failed"}
	ErrFileTooLarge          = DomainError{Code: "FILE_TOO_LARGE", Message: "Generated file is too large"}
	ErrFileUploadFailed      = DomainError{Code: "FILE_UPLOAD_FAILED", Message: "File upload failed"}
)

// NewDomainError cria um novo erro de domínio personalizado
func NewDomainError(code, message string, details map[string]interface{}) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WithDetails adiciona detalhes ao erro
func (e DomainError) WithDetails(details map[string]interface{}) DomainError {
	e.Details = details
	return e
}

// WithDetail adiciona um detalhe específico ao erro
func (e DomainError) WithDetail(key string, value interface{}) DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}