package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
)

// SchedulerService serviço de agendamento de relatórios
type SchedulerService struct {
	scheduleRepo  domain.ReportScheduleRepository
	reportService *ReportService
	eventBus      domain.EventBus
	cron          *cron.Cron
	logger        *zap.Logger
	scheduleMap   map[uuid.UUID]cron.EntryID
}

// NewSchedulerService cria nova instância do serviço
func NewSchedulerService(
	scheduleRepo domain.ReportScheduleRepository,
	reportService *ReportService,
	eventBus domain.EventBus,
	logger *zap.Logger,
) *SchedulerService {
	c := cron.New(cron.WithSeconds())
	
	return &SchedulerService{
		scheduleRepo:  scheduleRepo,
		reportService: reportService,
		eventBus:      eventBus,
		cron:          c,
		logger:        logger,
		scheduleMap:   make(map[uuid.UUID]cron.EntryID),
	}
}

// Start inicia o scheduler
func (s *SchedulerService) Start(ctx context.Context) error {
	s.logger.Info("Starting report scheduler")

	// Carregar agendamentos ativos
	schedules, err := s.scheduleRepo.GetActiveSchedules(ctx)
	if err != nil {
		return fmt.Errorf("failed to load active schedules: %w", err)
	}

	// Registrar cada agendamento
	for _, schedule := range schedules {
		if err := s.registerSchedule(schedule); err != nil {
			s.logger.Error("Failed to register schedule",
				zap.String("schedule_id", schedule.ID.String()),
				zap.Error(err))
			continue
		}
	}

	// Iniciar cron
	s.cron.Start()

	// Adicionar job para verificar agendamentos pendentes a cada minuto
	_, err = s.cron.AddFunc("@every 1m", func() {
		s.checkDueSchedules(context.Background())
	})

	if err != nil {
		return fmt.Errorf("failed to add due check job: %w", err)
	}

	s.logger.Info("Report scheduler started", zap.Int("active_schedules", len(schedules)))
	return nil
}

// Stop para o scheduler
func (s *SchedulerService) Stop() {
	s.logger.Info("Stopping report scheduler")
	s.cron.Stop()
}

// CreateScheduleRequest requisição para criar agendamento
type CreateScheduleRequest struct {
	ReportType     domain.ReportType       `json:"report_type" validate:"required"`
	Title          string                  `json:"title" validate:"required"`
	Description    string                  `json:"description,omitempty"`
	Format         domain.ReportFormat     `json:"format" validate:"required"`
	Frequency      domain.ReportFrequency  `json:"frequency" validate:"required"`
	CronExpression *string                 `json:"cron_expression,omitempty"`
	Parameters     map[string]interface{}  `json:"parameters,omitempty"`
	Filters        map[string]interface{}  `json:"filters,omitempty"`
	Recipients     []string                `json:"recipients,omitempty"`
}

// CreateSchedule cria um novo agendamento
func (s *SchedulerService) CreateSchedule(ctx context.Context, req *CreateScheduleRequest) (*domain.ReportSchedule, error) {
	tenantID := domain.MustGetTenantID(ctx)
	userID := domain.MustGetUserID(ctx)

	s.logger.Info("Creating report schedule",
		zap.String("tenant_id", tenantID.String()),
		zap.String("title", req.Title),
		zap.String("frequency", string(req.Frequency)))

	// Validar cron expression se fornecida
	if req.CronExpression != nil {
		if _, err := cron.ParseStandard(*req.CronExpression); err != nil {
			return nil, domain.ErrInvalidCronExpression.WithDetail("expression", *req.CronExpression)
		}
	}

	// Criar agendamento
	schedule := &domain.ReportSchedule{
		ID:             uuid.New(),
		TenantID:       tenantID,
		UserID:         userID,
		ReportType:     req.ReportType,
		Title:          req.Title,
		Description:    req.Description,
		Format:         req.Format,
		Frequency:      req.Frequency,
		CronExpression: req.CronExpression,
		Parameters:     req.Parameters,
		Filters:        req.Filters,
		Recipients:     req.Recipients,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Calcular próxima execução
	nextRun := s.calculateNextRun(schedule)
	schedule.NextRunAt = &nextRun

	// Salvar no banco
	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	// Registrar no cron
	if err := s.registerSchedule(schedule); err != nil {
		s.logger.Error("Failed to register schedule", zap.Error(err))
		// Não falhar a criação, pode registrar depois
	}

	// Publicar evento
	event := domain.ReportScheduledEvent{
		ReportEvent: domain.ReportEvent{
			EventID:    uuid.New(),
			EventType:  "report.scheduled",
			TenantID:   tenantID,
			UserID:     userID,
			OccurredAt: time.Now(),
		},
		ScheduleID: schedule.ID,
		Frequency:  schedule.Frequency,
		NextRunAt:  nextRun,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	return schedule, nil
}

// registerSchedule registra um agendamento no cron
func (s *SchedulerService) registerSchedule(schedule *domain.ReportSchedule) error {
	// Remover agendamento anterior se existir
	if entryID, exists := s.scheduleMap[schedule.ID]; exists {
		s.cron.Remove(entryID)
	}

	// Se não está ativo, não registrar
	if !schedule.IsActive {
		return nil
	}

	// Determinar expressão cron
	cronExpr := s.getCronExpression(schedule)
	if cronExpr == "" {
		return fmt.Errorf("cannot determine cron expression for schedule")
	}

	// Adicionar job
	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.executeSchedule(context.Background(), schedule.ID)
	})

	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	s.scheduleMap[schedule.ID] = entryID
	return nil
}

// executeSchedule executa um agendamento
func (s *SchedulerService) executeSchedule(ctx context.Context, scheduleID uuid.UUID) {
	s.logger.Info("Executing scheduled report", zap.String("schedule_id", scheduleID.String()))

	// Buscar agendamento atualizado
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		s.logger.Error("Failed to get schedule", zap.Error(err))
		return
	}

	// Verificar se ainda está ativo
	if !schedule.IsActive {
		s.logger.Warn("Schedule is inactive", zap.String("schedule_id", scheduleID.String()))
		return
	}

	// Criar contexto com autenticação
	authCtx := domain.WithAuthContext(ctx, domain.AuthContext{
		TenantID: schedule.TenantID,
		UserID:   schedule.UserID,
		TraceID:  uuid.New().String(),
	})

	// Criar relatório
	reportReq := &CreateReportRequest{
		Type:        schedule.ReportType,
		Title:       fmt.Sprintf("%s - %s", schedule.Title, time.Now().Format("2006-01-02 15:04")),
		Description: fmt.Sprintf("Scheduled report: %s", schedule.Description),
		Format:      schedule.Format,
		Parameters:  schedule.Parameters,
		Filters:     schedule.Filters,
		ScheduleID:  &schedule.ID,
	}

	report, err := s.reportService.CreateReport(authCtx, reportReq)
	if err != nil {
		s.logger.Error("Failed to create scheduled report",
			zap.String("schedule_id", scheduleID.String()),
			zap.Error(err))
		return
	}

	// Atualizar última execução e próxima
	now := time.Now()
	nextRun := s.calculateNextRun(schedule)
	
	if err := s.scheduleRepo.UpdateLastRun(ctx, schedule.ID, now, nextRun, &report.ID); err != nil {
		s.logger.Error("Failed to update schedule last run", zap.Error(err))
	}

	// Se tem destinatários, agendar envio após conclusão
	if len(schedule.Recipients) > 0 {
		go s.waitAndSendReport(context.Background(), report.ID, schedule.Recipients)
	}
}

// waitAndSendReport aguarda relatório ser gerado e envia
func (s *SchedulerService) waitAndSendReport(ctx context.Context, reportID uuid.UUID, recipients []string) {
	// Aguardar até 5 minutos pela conclusão
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timeout:
			s.logger.Warn("Timeout waiting for report", zap.String("report_id", reportID.String()))
			return
		case <-ticker.C:
			// Verificar status do relatório
			report, err := s.reportService.reportRepo.GetByID(ctx, reportID)
			if err != nil {
				s.logger.Error("Failed to get report status", zap.Error(err))
				return
			}

			if report.Status == domain.ReportStatusCompleted {
				// Publicar evento para envio
				event := domain.ReportEmailedEvent{
					ReportEvent: domain.ReportEvent{
						EventID:    uuid.New(),
						EventType:  "report.emailed",
						ReportID:   report.ID,
						TenantID:   report.TenantID,
						UserID:     report.UserID,
						OccurredAt: time.Now(),
					},
					Recipients: recipients,
					Subject:    fmt.Sprintf("Relatório: %s", report.Title),
				}

				if err := s.eventBus.Publish(ctx, event); err != nil {
					s.logger.Error("Failed to publish email event", zap.Error(err))
				}
				return
			} else if report.Status == domain.ReportStatusFailed {
				s.logger.Warn("Scheduled report failed", zap.String("report_id", reportID.String()))
				return
			}
		}
	}
}

// checkDueSchedules verifica agendamentos que devem ser executados
func (s *SchedulerService) checkDueSchedules(ctx context.Context) {
	now := time.Now()
	schedules, err := s.scheduleRepo.GetDueSchedules(ctx, now)
	if err != nil {
		s.logger.Error("Failed to get due schedules", zap.Error(err))
		return
	}

	for _, schedule := range schedules {
		// Executar se não registrado no cron (frequência once)
		if _, exists := s.scheduleMap[schedule.ID]; !exists {
			go s.executeSchedule(ctx, schedule.ID)
		}
	}
}

// getCronExpression obtém expressão cron baseada na frequência
func (s *SchedulerService) getCronExpression(schedule *domain.ReportSchedule) string {
	// Se tem expressão customizada, usar ela
	if schedule.CronExpression != nil && *schedule.CronExpression != "" {
		return *schedule.CronExpression
	}

	// Baseado na frequência
	switch schedule.Frequency {
	case domain.ReportFrequencyDaily:
		return "0 0 9 * * *" // Todo dia às 9h
	case domain.ReportFrequencyWeekly:
		return "0 0 9 * * 1" // Toda segunda às 9h
	case domain.ReportFrequencyMonthly:
		return "0 0 9 1 * *" // Dia 1 de cada mês às 9h
	case domain.ReportFrequencyOnce:
		return "" // Não usar cron
	default:
		return ""
	}
}

// calculateNextRun calcula próxima execução
func (s *SchedulerService) calculateNextRun(schedule *domain.ReportSchedule) time.Time {
	if schedule.Frequency == domain.ReportFrequencyOnce {
		// Para execução única, usar horário agendado ou agora + 1 minuto
		if schedule.NextRunAt != nil {
			return *schedule.NextRunAt
		}
		return time.Now().Add(1 * time.Minute)
	}

	// Para outras frequências, calcular baseado no cron
	cronExpr := s.getCronExpression(schedule)
	if cronExpr == "" {
		return time.Now().Add(24 * time.Hour) // Fallback para amanhã
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Now().Add(24 * time.Hour) // Fallback
	}

	return sched.Next(time.Now())
}

// GetSchedule busca agendamento por ID
func (s *SchedulerService) GetSchedule(ctx context.Context, id uuid.UUID) (*domain.ReportSchedule, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verificar permissão
	tenantID := domain.MustGetTenantID(ctx)
	if schedule.TenantID != tenantID {
		return nil, domain.ErrUnauthorized
	}

	return schedule, nil
}

// ListSchedules lista agendamentos do tenant
func (s *SchedulerService) ListSchedules(ctx context.Context) ([]*domain.ReportSchedule, error) {
	tenantID := domain.MustGetTenantID(ctx)
	return s.scheduleRepo.GetByTenantID(ctx, tenantID)
}

// UpdateSchedule atualiza um agendamento
func (s *SchedulerService) UpdateSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*domain.ReportSchedule, error) {
	schedule, err := s.GetSchedule(ctx, id)
	if err != nil {
		return nil, err
	}

	// Aplicar atualizações
	if isActive, ok := updates["is_active"].(bool); ok {
		schedule.IsActive = isActive
	}
	if recipients, ok := updates["recipients"].([]string); ok {
		schedule.Recipients = recipients
	}
	if filters, ok := updates["filters"].(map[string]interface{}); ok {
		schedule.Filters = filters
	}

	schedule.UpdatedAt = time.Now()

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	// Re-registrar no cron
	if err := s.registerSchedule(schedule); err != nil {
		s.logger.Error("Failed to re-register schedule", zap.Error(err))
	}

	return schedule, nil
}

// DeleteSchedule exclui um agendamento
func (s *SchedulerService) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	schedule, err := s.GetSchedule(ctx, id)
	if err != nil {
		return err
	}

	// Remover do cron
	if entryID, exists := s.scheduleMap[schedule.ID]; exists {
		s.cron.Remove(entryID)
		delete(s.scheduleMap, schedule.ID)
	}

	return s.scheduleRepo.Delete(ctx, id)
}