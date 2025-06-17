package commands

import (
	"context"
	"time"
	"fmt"
	"github.com/google/uuid"
	"github.com/direito-lux/process-service/internal/domain"
)

// MovementCommandHandler handler para comandos de movimentação
type MovementCommandHandler struct {
	movementRepo   domain.MovementRepository
	processRepo    domain.ProcessRepository
	eventPublisher domain.EventPublisher
}

// NewMovementCommandHandler cria novo handler de comandos de movimentação
func NewMovementCommandHandler(
	movementRepo domain.MovementRepository,
	processRepo domain.ProcessRepository,
	eventPublisher domain.EventPublisher,
) *MovementCommandHandler {
	return &MovementCommandHandler{
		movementRepo:   movementRepo,
		processRepo:    processRepo,
		eventPublisher: eventPublisher,
	}
}

// HandleCreateMovement processa comando de criação de movimentação
func (h *MovementCommandHandler) HandleCreateMovement(ctx context.Context, cmd *CreateMovementCommand) error {
	// Validar comando
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("comando inválido: %w", err)
	}

	// Verificar se processo existe
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return fmt.Errorf("processo não encontrado: %w", err)
	}

	// Obter próximo número de sequência
	existingMovements, _ := h.movementRepo.GetByProcess(cmd.ProcessID, domain.MovementFilters{
		SortBy:    "sequence",
		SortOrder: "desc",
		Limit:     1,
	})

	sequence := 1
	if len(existingMovements) > 0 {
		sequence = existingMovements[0].Sequence + 1
	}

	// Criar entidade movimentação
	movement := &domain.Movement{
		ID:             uuid.New().String(),
		ProcessID:      cmd.ProcessID,
		TenantID:       cmd.TenantID,
		Sequence:       sequence,
		ExternalID:     cmd.ExternalID,
		Date:           cmd.Date,
		Type:           cmd.Type,
		Code:           cmd.Code,
		Title:          cmd.Title,
		Description:    cmd.Description,
		Content:        cmd.Content,
		Judge:          cmd.Judge,
		Responsible:    cmd.Responsible,
		RelatedParties: cmd.RelatedParties,
		IsImportant:    cmd.IsImportant,
		IsPublic:       cmd.IsPublic,
		Tags:           cmd.Tags,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		SyncedAt:       time.Now(),
	}

	// Configurar metadados
	movement.Metadata = domain.MovementMetadata{
		OriginalSource: cmd.Source,
		CustomFields:   make(map[string]string),
	}

	// Validar movimentação
	if err := movement.ValidateType(); err != nil {
		return err
	}
	if err := movement.ValidateDate(); err != nil {
		return err
	}
	if err := movement.ValidateSequence(); err != nil {
		return err
	}

	// Persistir movimentação
	if err := h.movementRepo.Create(movement); err != nil {
		return fmt.Errorf("erro ao criar movimentação: %w", err)
	}

	// Atualizar última movimentação do processo
	process.UpdateLastMovement()
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao atualizar processo: %w", err)
	}

	// Publicar evento
	event := domain.NewMovementCreatedEvent(movement, process.Number, cmd.Source)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleUpdateMovement processa comando de atualização de movimentação
func (h *MovementCommandHandler) HandleUpdateMovement(ctx context.Context, cmd *UpdateMovementCommand) error {
	// Buscar movimentação
	movement, err := h.movementRepo.GetByID(cmd.ID)
	if err != nil {
		return err
	}

	// Aplicar mudanças
	if cmd.Title != nil {
		movement.Title = *cmd.Title
	}

	if cmd.Description != nil {
		movement.Description = *cmd.Description
	}

	if cmd.Content != nil {
		movement.Content = *cmd.Content
	}

	if cmd.Judge != nil {
		movement.Judge = *cmd.Judge
	}

	if cmd.Responsible != nil {
		movement.Responsible = *cmd.Responsible
	}

	if cmd.IsImportant != nil {
		movement.IsImportant = *cmd.IsImportant
	}

	if cmd.Tags != nil {
		movement.Tags = cmd.Tags
	}

	// Atualizar timestamp
	movement.UpdatedAt = time.Now()

	// Persistir mudanças
	if err := h.movementRepo.Update(movement); err != nil {
		return fmt.Errorf("erro ao atualizar movimentação: %w", err)
	}

	return nil
}

// HandleAnalyzeMovement processa comando de análise de movimentação
func (h *MovementCommandHandler) HandleAnalyzeMovement(ctx context.Context, cmd *AnalyzeMovementCommand) error {
	// Buscar movimentação
	movement, err := h.movementRepo.GetByID(cmd.MovementID)
	if err != nil {
		return err
	}

	// Aplicar análise
	movement.SetAnalysis(cmd.Sentiment, cmd.Importance, cmd.Category, cmd.Confidence, cmd.ProcessedBy)

	// Persistir mudanças
	if err := h.movementRepo.Update(movement); err != nil {
		return fmt.Errorf("erro ao atualizar análise: %w", err)
	}

	// Buscar processo para evento
	process, err := h.processRepo.GetByID(movement.ProcessID)
	if err != nil {
		return fmt.Errorf("processo não encontrado: %w", err)
	}

	// Publicar evento
	event := domain.NewMovementAnalyzedEvent(movement, process.Number)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	// Se movimentação é importante, publicar evento específico
	if movement.IsImportant && process.Monitoring.Enabled {
		importantEvent := domain.NewImportantMovementDetectedEvent(movement, process.Number, process.Monitoring.NotificationChannels)
		if err := h.eventPublisher.Publish(importantEvent); err != nil {
			return fmt.Errorf("erro ao publicar evento de movimentação importante: %w", err)
		}
	}

	return nil
}

// SyncCommandHandler handler para comandos de sincronização
type SyncCommandHandler struct {
	processRepo    domain.ProcessRepository
	movementRepo   domain.MovementRepository
	eventPublisher domain.EventPublisher
}

// NewSyncCommandHandler cria novo handler de comandos de sincronização
func NewSyncCommandHandler(
	processRepo domain.ProcessRepository,
	movementRepo domain.MovementRepository,
	eventPublisher domain.EventPublisher,
) *SyncCommandHandler {
	return &SyncCommandHandler{
		processRepo:    processRepo,
		movementRepo:   movementRepo,
		eventPublisher: eventPublisher,
	}
}

// HandleSyncProcess processa comando de sincronização de processo
func (h *SyncCommandHandler) HandleSyncProcess(ctx context.Context, cmd *SyncProcessCommand) error {
	// Buscar processo
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return err
	}

	// Simular sincronização (integração com DataJud seria aqui)
	newMovements := 0
	updatedMovements := 0

	// TODO: Implementar integração real com DataJud
	// Por enquanto, apenas simular resultado da sincronização

	// Atualizar última sincronização do processo
	process.UpdateLastSync()
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao atualizar processo: %w", err)
	}

	// Publicar evento
	event := domain.NewProcessSyncedEvent(process, newMovements, updatedMovements, cmd.SyncSource)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleBatchSyncProcesses processa comando de sincronização em lote
func (h *SyncCommandHandler) HandleBatchSyncProcesses(ctx context.Context, cmd *BatchSyncProcessesCommand) error {
	batchID := uuid.New().String()
	startTime := time.Now()

	// Publicar evento de início
	startEvent := domain.NewProcessBatchSyncStartedEvent(cmd.TenantID, batchID, cmd.ProcessIDs, cmd.SyncType, cmd.StartedBy)
	if err := h.eventPublisher.Publish(startEvent); err != nil {
		return fmt.Errorf("erro ao publicar evento de início: %w", err)
	}

	// Processar cada processo
	successCount := 0
	errorCount := 0
	totalNewMovements := 0
	totalUpdatedMovements := 0

	for _, processID := range cmd.ProcessIDs {
		// Criar comando de sincronização individual
		syncCmd := &SyncProcessCommand{
			ProcessID:  processID,
			TenantID:   cmd.TenantID,
			ForceSync:  cmd.ForceSync,
			SyncSource: "batch_" + cmd.SyncType,
			SyncedBy:   cmd.StartedBy,
		}

		// Executar sincronização
		if err := h.HandleSyncProcess(ctx, syncCmd); err != nil {
			errorCount++
			continue
		}

		successCount++
		// Em uma implementação real, contaríamos as movimentações novas/atualizadas
		totalNewMovements += 0
		totalUpdatedMovements += 0
	}

	// Publicar evento de conclusão
	duration := time.Since(startTime)
	completedEvent := domain.NewProcessBatchSyncCompletedEvent(
		cmd.TenantID, batchID, len(cmd.ProcessIDs), successCount, errorCount,
		totalNewMovements, totalUpdatedMovements, duration,
	)
	if err := h.eventPublisher.Publish(completedEvent); err != nil {
		return fmt.Errorf("erro ao publicar evento de conclusão: %w", err)
	}

	return nil
}

// PartyCommandHandler handler para comandos de partes
type PartyCommandHandler struct {
	partyRepo      domain.PartyRepository
	processRepo    domain.ProcessRepository
	eventPublisher domain.EventPublisher
}

// NewPartyCommandHandler cria novo handler de comandos de partes
func NewPartyCommandHandler(
	partyRepo domain.PartyRepository,
	processRepo domain.ProcessRepository,
	eventPublisher domain.EventPublisher,
) *PartyCommandHandler {
	return &PartyCommandHandler{
		partyRepo:      partyRepo,
		processRepo:    processRepo,
		eventPublisher: eventPublisher,
	}
}

// HandleAddParty processa comando de adição de parte
func (h *PartyCommandHandler) HandleAddParty(ctx context.Context, cmd *AddPartyCommand) error {
	// Verificar se processo existe
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return fmt.Errorf("processo não encontrado: %w", err)
	}

	// Criar entidade parte
	party := &domain.Party{
		ID:           uuid.New().String(),
		ProcessID:    cmd.ProcessID,
		Type:         cmd.Party.Type,
		Name:         cmd.Party.Name,
		Document:     cmd.Party.Document,
		DocumentType: cmd.Party.DocumentType,
		Role:         cmd.Party.Role,
		IsActive:     cmd.Party.IsActive,
		Contact: domain.PartyContact{
			Email:     cmd.Party.Contact.Email,
			Phone:     cmd.Party.Contact.Phone,
			CellPhone: cmd.Party.Contact.CellPhone,
			Website:   cmd.Party.Contact.Website,
		},
		Address: domain.PartyAddress{
			Street:     cmd.Party.Address.Street,
			Number:     cmd.Party.Address.Number,
			Complement: cmd.Party.Address.Complement,
			District:   cmd.Party.Address.District,
			City:       cmd.Party.Address.City,
			State:      cmd.Party.Address.State,
			ZipCode:    cmd.Party.Address.ZipCode,
			Country:    cmd.Party.Address.Country,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Adicionar advogado se fornecido
	if cmd.Party.Lawyer != nil {
		party.Lawyer = &domain.Lawyer{
			Name:     cmd.Party.Lawyer.Name,
			OAB:      cmd.Party.Lawyer.OAB,
			OABState: cmd.Party.Lawyer.OABState,
			Email:    cmd.Party.Lawyer.Email,
			Phone:    cmd.Party.Lawyer.Phone,
		}
	}

	// Validar parte
	if err := party.ValidateDocument(); err != nil {
		return err
	}
	if err := party.ValidateType(); err != nil {
		return err
	}
	if err := party.ValidateRole(); err != nil {
		return err
	}
	if err := party.ValidateLawyer(); err != nil {
		return err
	}

	// Persistir parte
	if err := h.partyRepo.Create(party); err != nil {
		return fmt.Errorf("erro ao criar parte: %w", err)
	}

	// Publicar evento
	event := domain.NewPartyAddedEvent(party, process.Number, cmd.AddedBy)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleUpdateParty processa comando de atualização de parte
func (h *PartyCommandHandler) HandleUpdateParty(ctx context.Context, cmd *UpdatePartyCommand) error {
	// Buscar parte
	party, err := h.partyRepo.GetByID(cmd.ID)
	if err != nil {
		return err
	}

	// Buscar processo para validações
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return fmt.Errorf("processo não encontrado: %w", err)
	}

	// Verificar se pode ser modificado
	if !process.CanBeModified() {
		return domain.ErrCannotModifyArchived
	}

	// Rastrear mudanças
	changes := make(map[string]interface{})

	// Aplicar mudanças
	if cmd.Name != nil && *cmd.Name != party.Name {
		changes["name"] = map[string]string{"from": party.Name, "to": *cmd.Name}
		party.Name = *cmd.Name
	}

	if cmd.Document != nil && *cmd.Document != party.Document {
		changes["document"] = map[string]string{"from": party.Document, "to": *cmd.Document}
		party.Document = *cmd.Document
	}

	if cmd.DocumentType != nil && *cmd.DocumentType != party.DocumentType {
		changes["document_type"] = map[string]string{"from": party.DocumentType, "to": *cmd.DocumentType}
		party.DocumentType = *cmd.DocumentType
	}

	if cmd.IsActive != nil && *cmd.IsActive != party.IsActive {
		changes["is_active"] = map[string]bool{"from": party.IsActive, "to": *cmd.IsActive}
		party.IsActive = *cmd.IsActive
	}

	if cmd.Lawyer != nil {
		party.Lawyer = &domain.Lawyer{
			Name:     cmd.Lawyer.Name,
			OAB:      cmd.Lawyer.OAB,
			OABState: cmd.Lawyer.OABState,
			Email:    cmd.Lawyer.Email,
			Phone:    cmd.Lawyer.Phone,
		}
		changes["lawyer"] = "updated"
	}

	if cmd.Contact != nil {
		party.Contact = domain.PartyContact{
			Email:     cmd.Contact.Email,
			Phone:     cmd.Contact.Phone,
			CellPhone: cmd.Contact.CellPhone,
			Website:   cmd.Contact.Website,
		}
		changes["contact"] = "updated"
	}

	if cmd.Address != nil {
		party.Address = domain.PartyAddress{
			Street:     cmd.Address.Street,
			Number:     cmd.Address.Number,
			Complement: cmd.Address.Complement,
			District:   cmd.Address.District,
			City:       cmd.Address.City,
			State:      cmd.Address.State,
			ZipCode:    cmd.Address.ZipCode,
			Country:    cmd.Address.Country,
		}
		changes["address"] = "updated"
	}

	// Atualizar timestamp
	party.UpdatedAt = time.Now()

	// Validar se necessário
	if err := party.ValidateDocument(); err != nil {
		return err
	}
	if err := party.ValidateLawyer(); err != nil {
		return err
	}

	// Persistir mudanças
	if err := h.partyRepo.Update(party); err != nil {
		return fmt.Errorf("erro ao atualizar parte: %w", err)
	}

	// Publicar evento se houve mudanças
	if len(changes) > 0 {
		event := domain.NewPartyUpdatedEvent(party.ID, process.Number, party.Name, changes, cmd.UpdatedBy)
		if err := h.eventPublisher.Publish(event); err != nil {
			return fmt.Errorf("erro ao publicar evento: %w", err)
		}
	}

	return nil
}

// HandleRemoveParty processa comando de remoção de parte
func (h *PartyCommandHandler) HandleRemoveParty(ctx context.Context, cmd *RemovePartyCommand) error {
	// Buscar parte
	party, err := h.partyRepo.GetByID(cmd.ID)
	if err != nil {
		return err
	}

	// Buscar processo para validações
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return fmt.Errorf("processo não encontrado: %w", err)
	}

	// Verificar se pode ser modificado
	if !process.CanBeModified() {
		return domain.ErrCannotModifyArchived
	}

	// Desativar parte (soft delete)
	party.Deactivate()

	// Persistir mudanças
	if err := h.partyRepo.Update(party); err != nil {
		return fmt.Errorf("erro ao remover parte: %w", err)
	}

	return nil
}