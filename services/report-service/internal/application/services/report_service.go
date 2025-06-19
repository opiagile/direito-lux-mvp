package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
	"github.com/direito-lux/report-service/internal/infrastructure/config"
)

// ReportService serviço de aplicação para relatórios
type ReportService struct {
	config           *config.Config
	reportRepo       domain.ReportRepository
	scheduleRepo     domain.ReportScheduleRepository
	templateRepo     domain.ReportTemplateRepository
	generator        domain.ReportGenerator
	dataCollector    domain.DataCollector
	eventBus         domain.EventBus
	logger           *zap.Logger
}

// NewReportService cria nova instância do serviço
func NewReportService(
	config *config.Config,
	reportRepo domain.ReportRepository,
	scheduleRepo domain.ReportScheduleRepository,
	templateRepo domain.ReportTemplateRepository,
	generator domain.ReportGenerator,
	dataCollector domain.DataCollector,
	eventBus domain.EventBus,
	logger *zap.Logger,
) *ReportService {
	return &ReportService{
		config:        config,
		reportRepo:    reportRepo,
		scheduleRepo:  scheduleRepo,
		templateRepo:  templateRepo,
		generator:     generator,
		dataCollector: dataCollector,
		eventBus:      eventBus,
		logger:        logger,
	}
}

// CreateReportRequest requisição para criar relatório
type CreateReportRequest struct {
	Type        domain.ReportType      `json:"type" validate:"required"`
	Title       string                 `json:"title" validate:"required"`
	Description string                 `json:"description,omitempty"`
	Format      domain.ReportFormat    `json:"format" validate:"required"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	TemplateID  *uuid.UUID             `json:"template_id,omitempty"`
	ScheduleID  *uuid.UUID             `json:"schedule_id,omitempty"`
}

// CreateReport cria um novo relatório
func (s *ReportService) CreateReport(ctx context.Context, req *CreateReportRequest) (*domain.Report, error) {
	tenantID := domain.MustGetTenantID(ctx)
	userID := domain.MustGetUserID(ctx)
	plan, _ := domain.GetPlan(ctx)

	s.logger.Info("Creating report",
		zap.String("tenant_id", tenantID.String()),
		zap.String("user_id", userID.String()),
		zap.String("type", string(req.Type)),
		zap.String("format", string(req.Format)))

	// Verificar quota do plano
	if err := s.checkPlanQuota(ctx, tenantID, plan); err != nil {
		return nil, err
	}

	// Criar relatório
	report := &domain.Report{
		ID:          uuid.New(),
		TenantID:    tenantID,
		UserID:      userID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Format:      req.Format,
		Status:      domain.ReportStatusPending,
		Parameters:  req.Parameters,
		Filters:     req.Filters,
		ScheduleID:  req.ScheduleID,
		CreatedAt:   time.Now(),
	}

	// Salvar no banco
	if err := s.reportRepo.Create(ctx, report); err != nil {
		s.logger.Error("Failed to create report", zap.Error(err))
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Publicar evento
	event := domain.ReportRequestedEvent{
		ReportEvent: domain.ReportEvent{
			EventID:    uuid.New(),
			EventType:  "report.requested",
			ReportID:   report.ID,
			TenantID:   tenantID,
			UserID:     userID,
			OccurredAt: time.Now(),
		},
		ReportType: req.Type,
		Format:     req.Format,
		Parameters: req.Parameters,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	// Processar relatório em background
	go s.processReport(context.Background(), report)

	return report, nil
}

// processReport processa a geração do relatório
func (s *ReportService) processReport(ctx context.Context, report *domain.Report) {
	startTime := time.Now()

	// Atualizar status para processando
	report.Status = domain.ReportStatusProcessing
	report.StartedAt = &startTime
	if err := s.reportRepo.Update(ctx, report); err != nil {
		s.logger.Error("Failed to update report status", zap.Error(err))
		return
	}

	// Coletar dados baseado no tipo
	data, err := s.collectReportData(ctx, report)
	if err != nil {
		s.handleReportError(ctx, report, err)
		return
	}

	// Gerar arquivo do relatório
	fileData, err := s.generateReportFile(ctx, report, data)
	if err != nil {
		s.handleReportError(ctx, report, err)
		return
	}

	// Salvar arquivo e atualizar relatório
	fileURL, fileSize, err := s.saveReportFile(ctx, report, fileData)
	if err != nil {
		s.handleReportError(ctx, report, err)
		return
	}

	// Atualizar relatório como completo
	completedAt := time.Now()
	processingTime := int64(completedAt.Sub(startTime).Seconds())
	expiresAt := completedAt.Add(time.Duration(s.config.Storage.RetentionDays) * 24 * time.Hour)

	report.Status = domain.ReportStatusCompleted
	report.FileURL = &fileURL
	report.FileSize = &fileSize
	report.ProcessingTime = &processingTime
	report.CompletedAt = &completedAt
	report.ExpiresAt = &expiresAt

	if err := s.reportRepo.Update(ctx, report); err != nil {
		s.logger.Error("Failed to update completed report", zap.Error(err))
		return
	}

	// Publicar evento de sucesso
	event := domain.ReportGeneratedEvent{
		ReportEvent: domain.ReportEvent{
			EventID:    uuid.New(),
			EventType:  "report.generated",
			ReportID:   report.ID,
			TenantID:   report.TenantID,
			UserID:     report.UserID,
			OccurredAt: time.Now(),
		},
		FileURL:        fileURL,
		FileSize:       fileSize,
		ProcessingTime: processingTime,
		Format:         report.Format,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish generated event", zap.Error(err))
	}

	s.logger.Info("Report generated successfully",
		zap.String("report_id", report.ID.String()),
		zap.Int64("processing_time", processingTime),
		zap.Int64("file_size", fileSize))
}

// collectReportData coleta dados para o relatório
func (s *ReportService) collectReportData(ctx context.Context, report *domain.Report) (interface{}, error) {
	switch report.Type {
	case domain.ReportTypeProcessAnalysis:
		return s.dataCollector.CollectProcessData(ctx, report.TenantID, report.Filters)
	
	case domain.ReportTypeProductivity:
		return s.dataCollector.CollectProductivityData(ctx, report.TenantID, report.Filters)
	
	case domain.ReportTypeFinancial:
		return s.dataCollector.CollectFinancialData(ctx, report.TenantID, report.Filters)
	
	case domain.ReportTypeJurisprudence:
		return s.dataCollector.CollectJurisprudenceData(ctx, report.TenantID, report.Filters)
	
	case domain.ReportTypeExecutiveSummary:
		// Coletar dados de múltiplas fontes
		processData, _ := s.dataCollector.CollectProcessData(ctx, report.TenantID, report.Filters)
		productivityData, _ := s.dataCollector.CollectProductivityData(ctx, report.TenantID, report.Filters)
		kpis, _ := s.dataCollector.CalculateKPIs(ctx, report.TenantID)
		
		return map[string]interface{}{
			"processes":    processData,
			"productivity": productivityData,
			"kpis":         kpis,
		}, nil
	
	default:
		return nil, fmt.Errorf("unsupported report type: %s", report.Type)
	}
}

// generateReportFile gera o arquivo do relatório
func (s *ReportService) generateReportFile(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	// Se tem template, aplicar primeiro
	if report.Parameters != nil && report.Parameters["template_id"] != nil {
		templateID, ok := report.Parameters["template_id"].(uuid.UUID)
		if ok {
			template, err := s.templateRepo.GetByID(ctx, templateID)
			if err == nil {
				// Aplicar template aos dados
				data = s.applyTemplate(template, data, report.Parameters)
			}
		}
	}

	switch report.Format {
	case domain.ReportFormatPDF:
		return s.generator.GeneratePDF(ctx, report, data)
	
	case domain.ReportFormatExcel:
		return s.generator.GenerateExcel(ctx, report, data)
	
	case domain.ReportFormatCSV:
		return s.generator.GenerateCSV(ctx, report, data)
	
	case domain.ReportFormatHTML:
		return s.generator.GenerateHTML(ctx, report, data)
	
	default:
		return nil, fmt.Errorf("unsupported format: %s", report.Format)
	}
}

// saveReportFile salva o arquivo do relatório
func (s *ReportService) saveReportFile(ctx context.Context, report *domain.Report, fileData []byte) (string, int64, error) {
	// Por enquanto, salvando localmente
	// TODO: Implementar upload para GCS/S3
	
	fileName := fmt.Sprintf("%s/%s_%s.%s",
		report.TenantID.String(),
		report.ID.String(),
		time.Now().Format("20060102_150405"),
		s.getFileExtension(report.Format))
	
	// Aqui seria o upload real
	fileURL := fmt.Sprintf("/reports/%s", fileName)
	fileSize := int64(len(fileData))
	
	return fileURL, fileSize, nil
}

// handleReportError trata erros na geração do relatório
func (s *ReportService) handleReportError(ctx context.Context, report *domain.Report, err error) {
	s.logger.Error("Report generation failed",
		zap.String("report_id", report.ID.String()),
		zap.Error(err))
	
	errorMsg := err.Error()
	report.Status = domain.ReportStatusFailed
	report.ErrorMessage = &errorMsg
	
	if updateErr := s.reportRepo.Update(ctx, report); updateErr != nil {
		s.logger.Error("Failed to update failed report", zap.Error(updateErr))
	}
	
	// Publicar evento de falha
	event := domain.ReportFailedEvent{
		ReportEvent: domain.ReportEvent{
			EventID:    uuid.New(),
			EventType:  "report.failed",
			ReportID:   report.ID,
			TenantID:   report.TenantID,
			UserID:     report.UserID,
			OccurredAt: time.Now(),
		},
		ErrorMessage: errorMsg,
	}
	
	if pubErr := s.eventBus.Publish(ctx, event); pubErr != nil {
		s.logger.Error("Failed to publish failed event", zap.Error(pubErr))
	}
}

// checkPlanQuota verifica a quota do plano
func (s *ReportService) checkPlanQuota(ctx context.Context, tenantID uuid.UUID, plan string) error {
	limit := s.config.GetReportLimitByPlan(plan)
	if limit == -1 {
		return nil // unlimited
	}
	
	// Contar relatórios do mês atual
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	filters := domain.ReportFilters{
		StartDate: &startOfMonth,
		EndDate:   &now,
	}
	
	reports, err := s.reportRepo.GetByTenantID(ctx, tenantID, filters)
	if err != nil {
		return fmt.Errorf("failed to check quota: %w", err)
	}
	
	if len(reports) >= limit {
		return domain.ErrQuotaExceeded.WithDetail("limit", limit).WithDetail("used", len(reports))
	}
	
	return nil
}

// getFileExtension retorna a extensão do arquivo baseado no formato
func (s *ReportService) getFileExtension(format domain.ReportFormat) string {
	switch format {
	case domain.ReportFormatPDF:
		return "pdf"
	case domain.ReportFormatExcel:
		return "xlsx"
	case domain.ReportFormatCSV:
		return "csv"
	case domain.ReportFormatHTML:
		return "html"
	default:
		return "bin"
	}
}

// applyTemplate aplica um template aos dados
func (s *ReportService) applyTemplate(template *domain.ReportTemplate, data interface{}, params map[string]interface{}) interface{} {
	// TODO: Implementar engine de template
	return data
}

// GetReport busca um relatório por ID
func (s *ReportService) GetReport(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Verificar permissão
	tenantID := domain.MustGetTenantID(ctx)
	if report.TenantID != tenantID {
		return nil, domain.ErrUnauthorized
	}
	
	return report, nil
}

// ListReports lista relatórios com filtros
func (s *ReportService) ListReports(ctx context.Context, filters domain.ReportFilters) ([]*domain.Report, error) {
	tenantID := domain.MustGetTenantID(ctx)
	return s.reportRepo.GetByTenantID(ctx, tenantID, filters)
}

// DeleteReport exclui um relatório
func (s *ReportService) DeleteReport(ctx context.Context, id uuid.UUID) error {
	report, err := s.GetReport(ctx, id)
	if err != nil {
		return err
	}
	
	// Não permitir excluir relatórios em processamento
	if report.Status == domain.ReportStatusProcessing {
		return domain.ErrReportInProgress
	}
	
	return s.reportRepo.Delete(ctx, id)
}

// GetStatistics obtém estatísticas de relatórios
func (s *ReportService) GetStatistics(ctx context.Context, period time.Duration) (*domain.ReportStatistics, error) {
	tenantID := domain.MustGetTenantID(ctx)
	return s.reportRepo.GetStatsByTenant(ctx, tenantID, period)
}