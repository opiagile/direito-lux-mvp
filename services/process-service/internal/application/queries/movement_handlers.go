package queries

import (
	"context"
	"fmt"
	"math"
	"time"
	"github.com/direito-lux/process-service/internal/domain"
)

// MovementQueryHandler handler para queries de movimentação
type MovementQueryHandler struct {
	movementRepo domain.MovementRepository
	processRepo  domain.ProcessRepository
}

// NewMovementQueryHandler cria novo handler de queries de movimentação
func NewMovementQueryHandler(
	movementRepo domain.MovementRepository,
	processRepo domain.ProcessRepository,
) *MovementQueryHandler {
	return &MovementQueryHandler{
		movementRepo: movementRepo,
		processRepo:  processRepo,
	}
}

// HandleGetMovement busca movimentação por ID
func (h *MovementQueryHandler) HandleGetMovement(ctx context.Context, query *MovementQuery) (*MovementDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Buscar movimentação
	movement, err := h.movementRepo.GetByID(query.ID)
	if err != nil {
		return nil, err
	}

	// Verificar tenant
	if movement.TenantID != query.TenantID {
		return nil, domain.ErrMovementNotFound
	}

	// Buscar processo para obter número
	process, err := h.processRepo.GetByID(movement.ProcessID)
	if err != nil {
		return nil, fmt.Errorf("processo não encontrado: %w", err)
	}

	// Converter para DTO
	dto := h.movementToDTO(movement, process.Number)

	return dto, nil
}

// HandleListMovements lista movimentações com filtros
func (h *MovementQueryHandler) HandleListMovements(ctx context.Context, query *MovementListQuery) (*MovementListDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Converter para filtros do domínio
	filters := domain.MovementFilters{
		ProcessID:       query.ProcessID,
		Type:            query.Type,
		IsImportant:     query.IsImportant,
		IsPublic:        query.IsPublic,
		HasNotification: query.HasNotification,
		DateFrom:        query.DateFrom,
		DateTo:          query.DateTo,
		Judge:           query.Judge,
		Tags:            query.Tags,
		Search:          query.Search,
		Limit:           query.PageSize,
		Offset:          (query.Page - 1) * query.PageSize,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
	}

	// Buscar movimentações
	var movements []*domain.Movement
	var err error

	if query.ProcessID != "" {
		movements, err = h.movementRepo.GetByProcess(query.ProcessID, filters)
	} else {
		movements, err = h.movementRepo.GetByTenant(query.TenantID, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações: %w", err)
	}

	// Buscar números dos processos para DTOs
	processNumbers := make(map[string]string)
	for _, movement := range movements {
		if _, exists := processNumbers[movement.ProcessID]; !exists {
			process, err := h.processRepo.GetByID(movement.ProcessID)
			if err == nil {
				processNumbers[movement.ProcessID] = process.Number
			}
		}
	}

	// Converter para DTOs
	movementDTOs := make([]MovementDTO, len(movements))
	for i, movement := range movements {
		processNumber := processNumbers[movement.ProcessID]
		dto := h.movementToDTO(movement, processNumber)
		
		// Incluir conteúdo completo se solicitado
		if !query.IncludeContent {
			dto.Content = ""
		}
		
		// Incluir anexos se solicitado
		if !query.IncludeAttachments {
			dto.Attachments = nil
		}
		
		movementDTOs[i] = *dto
	}

	// Calcular estatísticas para o total (sem paginação)
	allFilters := filters
	allFilters.Limit = 0
	allFilters.Offset = 0
	
	var allMovements []*domain.Movement
	if query.ProcessID != "" {
		allMovements, _ = h.movementRepo.GetByProcess(query.ProcessID, allFilters)
	} else {
		allMovements, _ = h.movementRepo.GetByTenant(query.TenantID, allFilters)
	}

	totalCount := len(allMovements)

	// Criar paginação
	totalPages := int(math.Ceil(float64(totalCount) / float64(query.PageSize)))
	pagination := PaginationDTO{
		Page:        query.Page,
		PageSize:    query.PageSize,
		TotalPages:  totalPages,
		TotalItems:  int64(totalCount),
		HasPrevious: query.Page > 1,
		HasNext:     query.Page < totalPages,
	}

	// Criar resumo
	summary := h.createMovementSummary(allMovements)

	return &MovementListDTO{
		Movements:  movementDTOs,
		Pagination: pagination,
		Summary:    summary,
	}, nil
}

// HandleSearchMovements busca textual em movimentações
func (h *MovementQueryHandler) HandleSearchMovements(ctx context.Context, query *MovementSearchQuery) (*SearchResultsDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	startTime := time.Now()

	// Converter para filtros do domínio
	filters := domain.MovementFilters{
		ProcessID: query.ProcessID,
		Search:    query.Query,
		DateFrom:  query.DateFrom,
		DateTo:    query.DateTo,
		IsPublic:  &query.OnlyPublic,
		Limit:     query.PageSize,
		Offset:    (query.Page - 1) * query.PageSize,
		SortBy:    "date",
		SortOrder: "desc",
	}

	// Buscar movimentações
	movements, err := h.movementRepo.SearchByContent(query.TenantID, query.Query, filters)
	if err != nil {
		return nil, fmt.Errorf("erro na busca: %w", err)
	}

	// Contar total sem paginação
	allFilters := filters
	allFilters.Limit = 0
	allFilters.Offset = 0
	allMovements, _ := h.movementRepo.SearchByContent(query.TenantID, query.Query, allFilters)
	totalCount := len(allMovements)

	// Converter para resultados de busca
	results := make([]SearchResultDTO, len(movements))
	for i, movement := range movements {
		results[i] = SearchResultDTO{
			ID:          movement.ID,
			Type:        "movement",
			Title:       movement.GetDisplayTitle(),
			Description: movement.GetSummary(),
			Score:       1.0, // Em produção seria calculado pelo algoritmo de busca
			Data: map[string]interface{}{
				"process_id":     movement.ProcessID,
				"date":          movement.Date,
				"type":          movement.Type,
				"is_important":  movement.IsImportant,
				"sequence":      movement.Sequence,
			},
		}

		// Adicionar highlights se solicitado
		if query.Highlights {
			results[i].Highlights = map[string][]string{
				"title":       {movement.Title},
				"description": {movement.Description},
			}
		}
	}

	// Criar paginação
	totalPages := int(math.Ceil(float64(totalCount) / float64(query.PageSize)))
	pagination := PaginationDTO{
		Page:        query.Page,
		PageSize:    query.PageSize,
		TotalPages:  totalPages,
		TotalItems:  int64(totalCount),
		HasPrevious: query.Page > 1,
		HasNext:     query.Page < totalPages,
	}

	return &SearchResultsDTO{
		Query:      query.Query,
		Results:    results,
		Pagination: pagination,
		TookMs:     time.Since(startTime).Milliseconds(),
	}, nil
}

// HandleGetMovementStats busca estatísticas de movimentações
func (h *MovementQueryHandler) HandleGetMovementStats(ctx context.Context, query *MovementStatsQuery) (*MovementSummaryDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Criar filtros
	filters := domain.MovementFilters{
		ProcessID: query.ProcessID,
		DateFrom:  query.DateFrom,
		DateTo:    query.DateTo,
	}

	// Buscar movimentações
	var movements []*domain.Movement
	var err error

	if query.ProcessID != "" {
		movements, err = h.movementRepo.GetByProcess(query.ProcessID, filters)
	} else {
		movements, err = h.movementRepo.GetByTenant(query.TenantID, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações: %w", err)
	}

	// Criar resumo
	summary := h.createMovementSummary(movements)

	return &summary, nil
}

// movementToDTO converte entidade Movement para DTO
func (h *MovementQueryHandler) movementToDTO(movement *domain.Movement, processNumber string) *MovementDTO {
	dto := &MovementDTO{
		ID:               movement.ID,
		ProcessID:        movement.ProcessID,
		ProcessNumber:    processNumber,
		TenantID:         movement.TenantID,
		Sequence:         movement.Sequence,
		ExternalID:       movement.ExternalID,
		Date:             movement.Date,
		Type:             movement.Type,
		Code:             movement.Code,
		Title:            movement.Title,
		Description:      movement.Description,
		Content:          movement.Content,
		Summary:          movement.GetSummary(),
		Judge:            movement.Judge,
		Responsible:      movement.Responsible,
		RelatedParties:   movement.RelatedParties,
		IsImportant:      movement.IsImportant,
		IsPublic:         movement.IsPublic,
		NotificationSent: movement.NotificationSent,
		Tags:             movement.Tags,
		FormattedDate:    movement.GetFormattedDate(),
		DisplayTitle:     movement.GetDisplayTitle(),
		ImportanceLevel:  movement.GetImportanceLevel(),
		CreatedAt:        movement.CreatedAt,
		UpdatedAt:        movement.UpdatedAt,
		SyncedAt:         movement.SyncedAt,
	}

	// Adicionar anexos
	if movement.Attachments != nil {
		dto.Attachments = make([]AttachmentDTO, len(movement.Attachments))
		for i, attachment := range movement.Attachments {
			dto.Attachments[i] = AttachmentDTO{
				ID:           attachment.ID,
				Name:         attachment.Name,
				Size:         attachment.Size,
				Type:         attachment.Type,
				URL:          attachment.URL,
				ExternalID:   attachment.ExternalID,
				IsDownloaded: attachment.IsDownloaded,
				CreatedAt:    attachment.CreatedAt,
			}
		}
	}

	// Adicionar metadados
	dto.Metadata = MovementMetadataDTO{
		OriginalSource: movement.Metadata.OriginalSource,
		DataJudID:      movement.Metadata.DataJudID,
		ImportBatch:    movement.Metadata.ImportBatch,
		Keywords:       movement.Metadata.Keywords,
		CustomFields:   movement.Metadata.CustomFields,
		Analysis: MovementAnalysisDTO{
			Sentiment:      movement.Metadata.Analysis.Sentiment,
			Importance:     movement.Metadata.Analysis.Importance,
			Category:       movement.Metadata.Analysis.Category,
			HasDeadline:    movement.Metadata.Analysis.HasDeadline,
			DeadlineDate:   movement.Metadata.Analysis.DeadlineDate,
			RequiresAction: movement.Metadata.Analysis.RequiresAction,
			ActionType:     movement.Metadata.Analysis.ActionType,
			Confidence:     movement.Metadata.Analysis.Confidence,
			ProcessedBy:    movement.Metadata.Analysis.ProcessedBy,
			ProcessedAt:    movement.Metadata.Analysis.ProcessedAt,
		},
	}

	return dto
}

// createMovementSummary cria resumo das movimentações
func (h *MovementQueryHandler) createMovementSummary(movements []*domain.Movement) MovementSummaryDTO {
	summary := MovementSummaryDTO{
		TotalMovements: len(movements),
		ByType:         make(map[domain.MovementType]int),
		ByMonth:        make(map[string]int),
	}

	for _, movement := range movements {
		// Contar por tipo
		summary.ByType[movement.Type]++

		// Contar movimentações importantes
		if movement.IsImportant {
			summary.ImportantMovements++
		}

		// Contar notificações pendentes
		if movement.RequiresNotification() {
			summary.PendingNotifications++
		}

		// Contar por mês
		monthKey := movement.Date.Format("2006-01")
		summary.ByMonth[monthKey]++
	}

	return summary
}

// DashboardQueryHandler handler para queries do dashboard
type DashboardQueryHandler struct {
	processRepo  domain.ProcessRepository
	movementRepo domain.MovementRepository
	partyRepo    domain.PartyRepository
}

// NewDashboardQueryHandler cria novo handler de queries do dashboard
func NewDashboardQueryHandler(
	processRepo domain.ProcessRepository,
	movementRepo domain.MovementRepository,
	partyRepo domain.PartyRepository,
) *DashboardQueryHandler {
	return &DashboardQueryHandler{
		processRepo:  processRepo,
		movementRepo: movementRepo,
		partyRepo:    partyRepo,
	}
}

// HandleGetDashboard busca dados do dashboard
func (h *DashboardQueryHandler) HandleGetDashboard(ctx context.Context, query *DashboardQuery) (*DashboardDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Calcular período baseado na query
	now := time.Now()
	var dateFrom time.Time
	
	switch query.Period {
	case "week":
		dateFrom = now.AddDate(0, 0, -7)
	case "month":
		dateFrom = now.AddDate(0, -1, 0)
	case "quarter":
		dateFrom = now.AddDate(0, -3, 0)
	case "year":
		dateFrom = now.AddDate(-1, 0, 0)
	default:
		dateFrom = now.AddDate(0, -1, 0) // padrão: último mês
	}

	// Buscar processos do período
	processFilters := domain.ProcessFilters{
		DateFrom: &dateFrom,
		DateTo:   &now,
	}

	var processes []*domain.Process
	var err error

	if query.ClientID != "" {
		processes, err = h.processRepo.GetByClient(query.ClientID, processFilters)
	} else {
		processes, err = h.processRepo.GetByTenant(query.TenantID, processFilters)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos: %w", err)
	}

	// Buscar movimentações do período
	movementFilters := domain.MovementFilters{
		DateFrom: &dateFrom,
		DateTo:   &now,
	}

	movements, err := h.movementRepo.GetByTenant(query.TenantID, movementFilters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações: %w", err)
	}

	// Construir dashboard
	dashboard := &DashboardDTO{
		TenantID:    query.TenantID,
		ClientID:    query.ClientID,
		Period:      query.Period,
		GeneratedAt: now,
	}

	// Visão geral
	dashboard.Overview = h.createOverview(processes, movements)

	// Dados de processos
	dashboard.Processes = h.createProcessesDashboard(processes)

	// Dados de movimentações
	dashboard.Movements = h.createMovementsDashboard(movements)

	// Dados de monitoramento
	dashboard.Monitoring = h.createMonitoringDashboard(processes)

	// Alertas (simulado)
	dashboard.Alerts = DashboardAlertsDTO{
		TotalAlerts:  0,
		UnreadAlerts: 0,
		BySeverity:   make(map[string]int),
		RecentAlerts: []AlertDTO{},
	}

	// Calendário (simplificado)
	dashboard.Calendar = []DashboardEventDTO{}

	// Atividade recente (simplificada)
	dashboard.RecentActivity = []DashboardActivityDTO{}

	return dashboard, nil
}

// createOverview cria visão geral do dashboard
func (h *DashboardQueryHandler) createOverview(processes []*domain.Process, movements []*domain.Movement) DashboardOverviewDTO {
	overview := DashboardOverviewDTO{
		TotalProcesses:   len(processes),
		TotalMovements:   len(movements),
		ProcessesGrowth:  0.0, // Seria calculado comparando com período anterior
		MovementsGrowth:  0.0,
	}

	// Contar processos por status
	for _, process := range processes {
		switch process.Status {
		case domain.ProcessStatusActive:
			overview.ActiveProcesses++
		case domain.ProcessStatusArchived:
			overview.ArchivedProcesses++
		}

		if process.Monitoring.Enabled {
			overview.MonitoringEnabled++
		}
	}

	// Contar movimentações importantes
	for _, movement := range movements {
		if movement.IsImportant {
			overview.ImportantMovements++
		}
	}

	return overview
}

// createProcessesDashboard cria dados de processos do dashboard
func (h *DashboardQueryHandler) createProcessesDashboard(processes []*domain.Process) DashboardProcessesDTO {
	dashboard := DashboardProcessesDTO{
		ByStatus:        make(map[domain.ProcessStatus]int),
		ByStage:         make(map[domain.ProcessStage]int),
		ByCourt:         make(map[string]int),
		MostActive:      []ProcessActivityDTO{},
		RecentlyCreated: []ProcessDTO{},
		NeedingAttention: []ProcessDTO{},
	}

	for _, process := range processes {
		dashboard.ByStatus[process.Status]++
		dashboard.ByStage[process.Stage]++
		dashboard.ByCourt[process.CourtID]++
	}

	return dashboard
}

// createMovementsDashboard cria dados de movimentações do dashboard
func (h *DashboardQueryHandler) createMovementsDashboard(movements []*domain.Movement) DashboardMovementsDTO {
	dashboard := DashboardMovementsDTO{
		ByType:       make(map[domain.MovementType]int),
		ByImportance: make(map[int]int),
		Timeline:     []MovementTimelineDTO{},
		TopKeywords:  []KeywordCountDTO{},
		RecentImportant: []MovementDTO{},
	}

	for _, movement := range movements {
		dashboard.ByType[movement.Type]++
		dashboard.ByImportance[movement.Metadata.Analysis.Importance]++
	}

	return dashboard
}

// createMonitoringDashboard cria dados de monitoramento do dashboard
func (h *DashboardQueryHandler) createMonitoringDashboard(processes []*domain.Process) DashboardMonitoringDTO {
	dashboard := DashboardMonitoringDTO{
		ChannelStats: make(map[string]int),
		SyncHistory:  []SyncHistoryDTO{},
	}

	for _, process := range processes {
		if process.Monitoring.Enabled {
			dashboard.EnabledProcesses++
			
			for _, channel := range process.Monitoring.NotificationChannels {
				dashboard.ChannelStats[channel]++
			}
		}

		if process.NeedsSync() {
			dashboard.NeedingSync++
		}
	}

	return dashboard
}