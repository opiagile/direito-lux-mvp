package commands

import (
	"context"
	"time"
	"fmt"
	"github.com/google/uuid"
	"github.com/direito-lux/process-service/internal/domain"
)

// ProcessCommandHandler handler principal para comandos de processo
type ProcessCommandHandler struct {
	processRepo    domain.ProcessRepository
	movementRepo   domain.MovementRepository
	partyRepo      domain.PartyRepository
	eventPublisher domain.EventPublisher
}

// NewProcessCommandHandler cria novo handler de comandos
func NewProcessCommandHandler(
	processRepo domain.ProcessRepository,
	movementRepo domain.MovementRepository,
	partyRepo domain.PartyRepository,
	eventPublisher domain.EventPublisher,
) *ProcessCommandHandler {
	return &ProcessCommandHandler{
		processRepo:    processRepo,
		movementRepo:   movementRepo,
		partyRepo:      partyRepo,
		eventPublisher: eventPublisher,
	}
}

// HandleCreateProcess processa comando de criação de processo
func (h *ProcessCommandHandler) HandleCreateProcess(ctx context.Context, cmd *CreateProcessCommand) error {
	// Validar comando
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("comando inválido: %w", err)
	}

	// Verificar se processo já existe
	existing, _ := h.processRepo.GetByNumber(cmd.Number)
	if existing != nil {
		return domain.ErrProcessExists
	}

	// Criar entidade processo
	process := &domain.Process{
		ID:             uuid.New().String(),
		TenantID:       cmd.TenantID,
		ClientID:       cmd.ClientID,
		Number:         cmd.Number,
		OriginalNumber: cmd.OriginalNumber,
		Title:          cmd.Title,
		Description:    cmd.Description,
		Status:         cmd.Status,
		Stage:          cmd.Stage,
		Subject: domain.ProcessSubject{
			Code:        cmd.Subject.Code,
			Description: cmd.Subject.Description,
			ParentCode:  cmd.Subject.ParentCode,
		},
		CourtID:      cmd.CourtID,
		JudgeID:      cmd.JudgeID,
		Tags:         cmd.Tags,
		CustomFields: cmd.CustomFields,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Adicionar valor da causa se fornecido
	if cmd.Value != nil {
		process.Value = &domain.ProcessValue{
			Amount:   cmd.Value.Amount,
			Currency: cmd.Value.Currency,
		}
	}

	// Configurar monitoramento
	process.Monitoring = domain.ProcessMonitoring{
		Enabled:              cmd.Monitoring.Enabled,
		NotificationChannels: cmd.Monitoring.NotificationChannels,
		Keywords:             cmd.Monitoring.Keywords,
		AutoSync:             cmd.Monitoring.AutoSync,
		SyncIntervalHours:    cmd.Monitoring.SyncIntervalHours,
	}

	// Validar processo
	if err := process.ValidateNumber(); err != nil {
		return err
	}
	if err := process.ValidateStatus(); err != nil {
		return err
	}
	if err := process.ValidateStage(); err != nil {
		return err
	}

	// Persistir processo
	if err := h.processRepo.Create(process); err != nil {
		return fmt.Errorf("erro ao criar processo: %w", err)
	}

	// Criar partes se fornecidas
	for _, partyCmd := range cmd.Parties {
		party := h.createPartyFromCommand(process.ID, &partyCmd)
		if err := h.partyRepo.Create(party); err != nil {
			return fmt.Errorf("erro ao criar parte: %w", err)
		}
	}

	// Publicar evento
	event := domain.NewProcessCreatedEvent(process, cmd.CreatedBy)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleUpdateProcess processa comando de atualização de processo
func (h *ProcessCommandHandler) HandleUpdateProcess(ctx context.Context, cmd *UpdateProcessCommand) error {
	// Buscar processo
	process, err := h.processRepo.GetByID(cmd.ID)
	if err != nil {
		return err
	}

	// Verificar se pode ser modificado
	if !process.CanBeModified() {
		return domain.ErrCannotModifyArchived
	}

	// Rastrear mudanças para evento
	changes := make(map[string]interface{})

	// Aplicar mudanças
	if cmd.Title != nil && *cmd.Title != process.Title {
		changes["title"] = map[string]string{"from": process.Title, "to": *cmd.Title}
		process.Title = *cmd.Title
	}

	if cmd.Description != nil && *cmd.Description != process.Description {
		changes["description"] = map[string]string{"from": process.Description, "to": *cmd.Description}
		process.Description = *cmd.Description
	}

	if cmd.Status != nil && *cmd.Status != process.Status {
		changes["status"] = map[string]string{"from": string(process.Status), "to": string(*cmd.Status)}
		process.Status = *cmd.Status
	}

	if cmd.Stage != nil && *cmd.Stage != process.Stage {
		changes["stage"] = map[string]string{"from": string(process.Stage), "to": string(*cmd.Stage)}
		process.Stage = *cmd.Stage
	}

	if cmd.Subject != nil {
		process.Subject = domain.ProcessSubject{
			Code:        cmd.Subject.Code,
			Description: cmd.Subject.Description,
			ParentCode:  cmd.Subject.ParentCode,
		}
		changes["subject"] = "updated"
	}

	if cmd.Value != nil {
		process.Value = &domain.ProcessValue{
			Amount:   cmd.Value.Amount,
			Currency: cmd.Value.Currency,
		}
		changes["value"] = "updated"
	}

	if cmd.JudgeID != nil {
		changes["judge_id"] = map[string]interface{}{"from": process.JudgeID, "to": *cmd.JudgeID}
		process.JudgeID = cmd.JudgeID
	}

	if cmd.Tags != nil {
		changes["tags"] = map[string]interface{}{"from": process.Tags, "to": cmd.Tags}
		process.Tags = cmd.Tags
	}

	if cmd.CustomFields != nil {
		changes["custom_fields"] = "updated"
		process.CustomFields = cmd.CustomFields
	}

	// Atualizar timestamp
	process.UpdatedAt = time.Now()

	// Validar se necessário
	if err := process.ValidateStatus(); err != nil {
		return err
	}
	if err := process.ValidateStage(); err != nil {
		return err
	}

	// Persistir mudanças
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao atualizar processo: %w", err)
	}

	// Publicar evento se houve mudanças
	if len(changes) > 0 {
		event := domain.NewProcessUpdatedEvent(process.ID, process.TenantID, process.Number, changes, cmd.UpdatedBy)
		if err := h.eventPublisher.Publish(event); err != nil {
			return fmt.Errorf("erro ao publicar evento: %w", err)
		}
	}

	return nil
}

// HandleArchiveProcess processa comando de arquivamento
func (h *ProcessCommandHandler) HandleArchiveProcess(ctx context.Context, cmd *ArchiveProcessCommand) error {
	// Buscar processo
	process, err := h.processRepo.GetByID(cmd.ID)
	if err != nil {
		return err
	}

	// Arquivar processo
	if err := process.Archive(); err != nil {
		return err
	}

	// Persistir mudanças
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao arquivar processo: %w", err)
	}

	// Publicar evento
	event := domain.NewProcessArchivedEvent(process, cmd.Reason, cmd.ArchivedBy)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleReactivateProcess processa comando de reativação
func (h *ProcessCommandHandler) HandleReactivateProcess(ctx context.Context, cmd *ReactivateProcessCommand) error {
	// Buscar processo
	process, err := h.processRepo.GetByID(cmd.ID)
	if err != nil {
		return err
	}

	// Reativar processo
	if err := process.Reactivate(); err != nil {
		return err
	}

	// Persistir mudanças
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao reativar processo: %w", err)
	}

	// Publicar evento
	event := domain.NewProcessReactivatedEvent(process, cmd.ReactivatedBy)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleEnableMonitoring processa comando de habilitação do monitoramento
func (h *ProcessCommandHandler) HandleEnableMonitoring(ctx context.Context, cmd *EnableMonitoringCommand) error {
	// Buscar processo
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return err
	}

	// Habilitar monitoramento
	process.EnableMonitoring(cmd.NotificationChannels, cmd.SyncIntervalHours)
	process.Monitoring.Keywords = cmd.Keywords

	// Persistir mudanças
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao habilitar monitoramento: %w", err)
	}

	// Publicar evento
	event := domain.NewProcessMonitoringEnabledEvent(process, cmd.EnabledBy)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// HandleDisableMonitoring processa comando de desabilitação do monitoramento
func (h *ProcessCommandHandler) HandleDisableMonitoring(ctx context.Context, cmd *DisableMonitoringCommand) error {
	// Buscar processo
	process, err := h.processRepo.GetByID(cmd.ProcessID)
	if err != nil {
		return err
	}

	// Desabilitar monitoramento
	process.DisableMonitoring()

	// Persistir mudanças
	if err := h.processRepo.Update(process); err != nil {
		return fmt.Errorf("erro ao desabilitar monitoramento: %w", err)
	}

	// Publicar evento
	event := domain.NewProcessMonitoringDisabledEvent(process, cmd.DisabledBy)
	if err := h.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	return nil
}

// createPartyFromCommand cria entidade Party a partir do comando
func (h *ProcessCommandHandler) createPartyFromCommand(processID string, cmd *CreatePartyCommand) *domain.Party {
	party := &domain.Party{
		ID:           uuid.New().String(),
		ProcessID:    processID,
		Type:         cmd.Type,
		Name:         cmd.Name,
		Document:     cmd.Document,
		DocumentType: cmd.DocumentType,
		Role:         cmd.Role,
		IsActive:     cmd.IsActive,
		Contact: domain.PartyContact{
			Email:     cmd.Contact.Email,
			Phone:     cmd.Contact.Phone,
			CellPhone: cmd.Contact.CellPhone,
			Website:   cmd.Contact.Website,
		},
		Address: domain.PartyAddress{
			Street:     cmd.Address.Street,
			Number:     cmd.Address.Number,
			Complement: cmd.Address.Complement,
			District:   cmd.Address.District,
			City:       cmd.Address.City,
			State:      cmd.Address.State,
			ZipCode:    cmd.Address.ZipCode,
			Country:    cmd.Address.Country,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Adicionar advogado se fornecido
	if cmd.Lawyer != nil {
		party.Lawyer = &domain.Lawyer{
			Name:     cmd.Lawyer.Name,
			OAB:      cmd.Lawyer.OAB,
			OABState: cmd.Lawyer.OABState,
			Email:    cmd.Lawyer.Email,
			Phone:    cmd.Lawyer.Phone,
		}
	}

	return party
}