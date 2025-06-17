package queries

import (
	"context"
	"fmt"
	"math"
	"time"
	"github.com/direito-lux/process-service/internal/domain"
)

// ProcessQueryHandler handler para queries de processo
type ProcessQueryHandler struct {
	processRepo  domain.ProcessRepository
	movementRepo domain.MovementRepository
	partyRepo    domain.PartyRepository
}

// NewProcessQueryHandler cria novo handler de queries de processo
func NewProcessQueryHandler(
	processRepo domain.ProcessRepository,
	movementRepo domain.MovementRepository,
	partyRepo domain.PartyRepository,
) *ProcessQueryHandler {
	return &ProcessQueryHandler{
		processRepo:  processRepo,
		movementRepo: movementRepo,
		partyRepo:    partyRepo,
	}
}

// HandleGetProcess busca processo por ID
func (h *ProcessQueryHandler) HandleGetProcess(ctx context.Context, query *ProcessQuery) (*ProcessDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Buscar processo
	process, err := h.processRepo.GetByID(query.ID)
	if err != nil {
		return nil, err
	}

	// Verificar tenant
	if process.TenantID != query.TenantID {
		return nil, domain.ErrProcessNotFound
	}

	// Buscar partes
	parties, _ := h.partyRepo.GetByProcess(process.ID)

	// Buscar estatísticas
	stats, _ := h.getProcessStats(process.ID)

	// Converter para DTO
	dto := h.processToDTO(process, parties, stats)

	return dto, nil
}

// HandleListProcesses lista processos com filtros
func (h *ProcessQueryHandler) HandleListProcesses(ctx context.Context, query *ProcessListQuery) (*ProcessListDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Converter para filtros do domínio
	filters := domain.ProcessFilters{
		Status:    query.Status,
		Stage:     query.Stage,
		CourtID:   query.CourtID,
		JudgeID:   query.JudgeID,
		Tags:      query.Tags,
		DateFrom:  query.DateFrom,
		DateTo:    query.DateTo,
		Search:    query.Search,
		Limit:     query.PageSize,
		Offset:    (query.Page - 1) * query.PageSize,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}

	// Buscar processos
	var processes []*domain.Process
	var err error

	if query.ClientID != "" {
		processes, err = h.processRepo.GetByClient(query.ClientID, filters)
	} else {
		processes, err = h.processRepo.GetByTenant(query.TenantID, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos: %w", err)
	}

	// Contar total
	totalCount, err := h.processRepo.CountByTenant(query.TenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar processos: %w", err)
	}

	// Converter para DTOs
	processDTOs := make([]ProcessDTO, len(processes))
	for i, process := range processes {
		// Buscar partes para cada processo (pode ser otimizado com batch loading)
		parties, _ := h.partyRepo.GetByProcess(process.ID)
		stats, _ := h.getProcessStats(process.ID)
		processDTOs[i] = *h.processToDTO(process, parties, stats)
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

	// Criar resumo
	summary := h.createProcessSummary(processes)

	return &ProcessListDTO{
		Processes:  processDTOs,
		Pagination: pagination,
		Summary:    summary,
	}, nil
}

// HandleGetProcessStats busca estatísticas de processos
func (h *ProcessQueryHandler) HandleGetProcessStats(ctx context.Context, query *ProcessStatsQuery) (*ProcessSummaryDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	// Criar filtros baseados na query
	filters := domain.ProcessFilters{
		DateFrom: query.DateFrom,
		DateTo:   query.DateTo,
	}

	// Buscar processos
	var processes []*domain.Process
	var err error

	if query.ClientID != "" {
		processes, err = h.processRepo.GetByClient(query.ClientID, filters)
	} else {
		processes, err = h.processRepo.GetByTenant(query.TenantID, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos: %w", err)
	}

	// Criar resumo
	summary := h.createProcessSummary(processes)

	return &summary, nil
}

// HandleGetMonitoringProcesses busca processos para monitoramento
func (h *ProcessQueryHandler) HandleGetMonitoringProcesses(ctx context.Context, query *ProcessMonitoringQuery) (*ProcessListDTO, error) {
	// Validar query
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query inválida: %w", err)
	}

	var processes []*domain.Process
	var err error

	if query.OnlyNeedingSync {
		// Buscar processos que precisam ser sincronizados
		processes, err = h.processRepo.GetNeedingSync()
	} else {
		// Buscar processos com monitoramento ativo
		processes, err = h.processRepo.GetActiveForMonitoring()
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos de monitoramento: %w", err)
	}

	// Filtrar por tenant
	filteredProcesses := make([]*domain.Process, 0)
	for _, process := range processes {
		if process.TenantID == query.TenantID {
			filteredProcesses = append(filteredProcesses, process)
		}
	}

	// Aplicar paginação
	offset := (query.Page - 1) * query.PageSize
	end := offset + query.PageSize
	if end > len(filteredProcesses) {
		end = len(filteredProcesses)
	}

	paginatedProcesses := filteredProcesses[offset:end]

	// Converter para DTOs
	processDTOs := make([]ProcessDTO, len(paginatedProcesses))
	for i, process := range paginatedProcesses {
		parties, _ := h.partyRepo.GetByProcess(process.ID)
		stats, _ := h.getProcessStats(process.ID)
		processDTOs[i] = *h.processToDTO(process, parties, stats)
	}

	// Criar paginação
	totalCount := len(filteredProcesses)
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
	summary := h.createProcessSummary(filteredProcesses)

	return &ProcessListDTO{
		Processes:  processDTOs,
		Pagination: pagination,
		Summary:    summary,
	}, nil
}

// processToDTO converte entidade Process para DTO
func (h *ProcessQueryHandler) processToDTO(process *domain.Process, parties []*domain.Party, stats *ProcessStatsDTO) *ProcessDTO {
	dto := &ProcessDTO{
		ID:             process.ID,
		TenantID:       process.TenantID,
		ClientID:       process.ClientID,
		Number:         process.Number,
		OriginalNumber: process.OriginalNumber,
		Title:          process.Title,
		Description:    process.Description,
		Status:         process.Status,
		Stage:          process.Stage,
		Subject: ProcessSubjectDTO{
			Code:        process.Subject.Code,
			Description: process.Subject.Description,
			ParentCode:  process.Subject.ParentCode,
		},
		CourtID:        process.CourtID,
		JudgeID:        process.JudgeID,
		Tags:           process.Tags,
		CustomFields:   process.CustomFields,
		LastMovementAt: process.LastMovementAt,
		LastSyncAt:     process.LastSyncAt,
		CreatedAt:      process.CreatedAt,
		UpdatedAt:      process.UpdatedAt,
		ArchivedAt:     process.ArchivedAt,
	}

	// Adicionar valor se existir
	if process.Value != nil {
		dto.Value = &ProcessValueDTO{
			Amount:         process.Value.Amount,
			Currency:       process.Value.Currency,
			FormattedValue: process.GetFormattedValue(),
		}
	}

	// Adicionar monitoramento
	dto.Monitoring = ProcessMonitoringDTO{
		Enabled:              process.Monitoring.Enabled,
		NotificationChannels: process.Monitoring.NotificationChannels,
		Keywords:             process.Monitoring.Keywords,
		AutoSync:             process.Monitoring.AutoSync,
		SyncIntervalHours:    process.Monitoring.SyncIntervalHours,
		LastNotificationAt:   process.Monitoring.LastNotificationAt,
		NeedsSync:            process.NeedsSync(),
	}

	// Adicionar partes
	if parties != nil {
		dto.Parties = make([]PartyDTO, len(parties))
		for i, party := range parties {
			dto.Parties[i] = *h.partyToDTO(party)
		}
	}

	// Adicionar estatísticas
	dto.Stats = stats

	return dto
}

// partyToDTO converte entidade Party para DTO
func (h *ProcessQueryHandler) partyToDTO(party *domain.Party) *PartyDTO {
	dto := &PartyDTO{
		ID:                party.ID,
		ProcessID:         party.ProcessID,
		Type:              party.Type,
		Name:              party.Name,
		Document:          party.Document,
		DocumentType:      party.DocumentType,
		FormattedDocument: party.GetFormattedDocument(),
		Role:              party.Role,
		IsActive:          party.IsActive,
		Contact: PartyContactDTO{
			Email:     party.Contact.Email,
			Phone:     party.Contact.Phone,
			CellPhone: party.Contact.CellPhone,
			Website:   party.Contact.Website,
		},
		Address: PartyAddressDTO{
			Street:     party.Address.Street,
			Number:     party.Address.Number,
			Complement: party.Address.Complement,
			District:   party.Address.District,
			City:       party.Address.City,
			State:      party.Address.State,
			ZipCode:    party.Address.ZipCode,
			Country:    party.Address.Country,
		},
		DisplayName:   party.GetDisplayName(),
		LawyerInfo:    party.GetLawyerInfo(),
		IsMainParty:   party.IsMainParty(),
		IsLegalEntity: party.IsLegalEntity(),
		CreatedAt:     party.CreatedAt,
		UpdatedAt:     party.UpdatedAt,
	}

	// Adicionar advogado se existir
	if party.Lawyer != nil {
		dto.Lawyer = &LawyerDTO{
			Name:     party.Lawyer.Name,
			OAB:      party.Lawyer.OAB,
			OABState: party.Lawyer.OABState,
			Email:    party.Lawyer.Email,
			Phone:    party.Lawyer.Phone,
		}
	}

	return dto
}

// getProcessStats busca estatísticas do processo
func (h *ProcessQueryHandler) getProcessStats(processID string) (*ProcessStatsDTO, error) {
	// Buscar movimentações
	movements, err := h.movementRepo.GetByProcess(processID, domain.MovementFilters{})
	if err != nil {
		return nil, err
	}

	stats := &ProcessStatsDTO{
		TotalMovements: len(movements),
	}

	// Contar movimentações importantes
	for _, movement := range movements {
		if movement.IsImportant {
			stats.ImportantMovements++
		}
	}

	// Encontrar última movimentação
	if len(movements) > 0 {
		latestDate := movements[0].Date
		for _, movement := range movements {
			if movement.Date.After(latestDate) {
				latestDate = movement.Date
			}
		}
		stats.LastMovementDate = &latestDate
	}

	// Calcular dias ativo (simplificado)
	// Em produção, seria baseado na data de criação vs última movimentação
	if stats.LastMovementDate != nil {
		stats.DaysActive = int(time.Since(*stats.LastMovementDate).Hours() / 24)
	}

	stats.AverageResponseTime = "N/A" // Seria calculado baseado em regras de negócio

	return stats, nil
}

// createProcessSummary cria resumo dos processos
func (h *ProcessQueryHandler) createProcessSummary(processes []*domain.Process) ProcessSummaryDTO {
	summary := ProcessSummaryDTO{
		TotalProcesses: len(processes),
		ByStatus:       make(map[domain.ProcessStatus]int),
		ByStage:        make(map[domain.ProcessStage]int),
	}

	for _, process := range processes {
		// Contar por status
		summary.ByStatus[process.Status]++

		// Contar por fase
		summary.ByStage[process.Stage]++

		// Contar monitoramento ativo
		if process.Monitoring.Enabled {
			summary.ActiveMonitoring++
		}

		// Contar que precisam sincronização
		if process.NeedsSync() {
			summary.NeedingSync++
		}
	}

	return summary
}