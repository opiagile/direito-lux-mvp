package domain

import (
	"time"
	"encoding/json"
)

// DomainEvent representa um evento de domínio
type DomainEvent interface {
	EventType() string
	EventVersion() string
	AggregateID() string
	Timestamp() time.Time
	ToJSON() ([]byte, error)
}

// BaseDomainEvent implementação base para eventos de domínio
type BaseDomainEvent struct {
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	ID          string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	TenantID    string    `json:"tenant_id"`
}

func (e BaseDomainEvent) EventType() string {
	return e.Type
}

func (e BaseDomainEvent) EventVersion() string {
	return e.Version
}

func (e BaseDomainEvent) AggregateID() string {
	return e.ID
}

func (e BaseDomainEvent) Timestamp() time.Time {
	return e.OccurredAt
}

func (e BaseDomainEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ProcessCreatedEvent evento quando um processo é criado
type ProcessCreatedEvent struct {
	BaseDomainEvent
	ProcessNumber string        `json:"process_number"`
	ClientID      string        `json:"client_id"`
	Title         string        `json:"title"`
	Status        ProcessStatus `json:"status"`
	Stage         ProcessStage  `json:"stage"`
	CourtID       string        `json:"court_id"`
	CreatedBy     string        `json:"created_by"`
}

func NewProcessCreatedEvent(process *Process, createdBy string) *ProcessCreatedEvent {
	return &ProcessCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.created",
			Version:    "1.0",
			ID:         process.ID,
			OccurredAt: time.Now(),
			TenantID:   process.TenantID,
		},
		ProcessNumber: process.Number,
		ClientID:      process.ClientID,
		Title:         process.Title,
		Status:        process.Status,
		Stage:         process.Stage,
		CourtID:       process.CourtID,
		CreatedBy:     createdBy,
	}
}

// ProcessUpdatedEvent evento quando um processo é atualizado
type ProcessUpdatedEvent struct {
	BaseDomainEvent
	ProcessNumber string                 `json:"process_number"`
	Changes       map[string]interface{} `json:"changes"`
	UpdatedBy     string                 `json:"updated_by"`
}

func NewProcessUpdatedEvent(processID, tenantID, processNumber string, changes map[string]interface{}, updatedBy string) *ProcessUpdatedEvent {
	return &ProcessUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.updated",
			Version:    "1.0",
			ID:         processID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		ProcessNumber: processNumber,
		Changes:       changes,
		UpdatedBy:     updatedBy,
	}
}

// ProcessArchivedEvent evento quando um processo é arquivado
type ProcessArchivedEvent struct {
	BaseDomainEvent
	ProcessNumber string `json:"process_number"`
	Reason        string `json:"reason"`
	ArchivedBy    string `json:"archived_by"`
}

func NewProcessArchivedEvent(process *Process, reason, archivedBy string) *ProcessArchivedEvent {
	return &ProcessArchivedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.archived",
			Version:    "1.0",
			ID:         process.ID,
			OccurredAt: time.Now(),
			TenantID:   process.TenantID,
		},
		ProcessNumber: process.Number,
		Reason:        reason,
		ArchivedBy:    archivedBy,
	}
}

// ProcessReactivatedEvent evento quando um processo é reativado
type ProcessReactivatedEvent struct {
	BaseDomainEvent
	ProcessNumber string `json:"process_number"`
	ReactivatedBy string `json:"reactivated_by"`
}

func NewProcessReactivatedEvent(process *Process, reactivatedBy string) *ProcessReactivatedEvent {
	return &ProcessReactivatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.reactivated",
			Version:    "1.0",
			ID:         process.ID,
			OccurredAt: time.Now(),
			TenantID:   process.TenantID,
		},
		ProcessNumber: process.Number,
		ReactivatedBy: reactivatedBy,
	}
}

// ProcessMonitoringEnabledEvent evento quando monitoramento é habilitado
type ProcessMonitoringEnabledEvent struct {
	BaseDomainEvent
	ProcessNumber        string   `json:"process_number"`
	NotificationChannels []string `json:"notification_channels"`
	SyncIntervalHours    int      `json:"sync_interval_hours"`
	EnabledBy            string   `json:"enabled_by"`
}

func NewProcessMonitoringEnabledEvent(process *Process, enabledBy string) *ProcessMonitoringEnabledEvent {
	return &ProcessMonitoringEnabledEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.monitoring.enabled",
			Version:    "1.0",
			ID:         process.ID,
			OccurredAt: time.Now(),
			TenantID:   process.TenantID,
		},
		ProcessNumber:        process.Number,
		NotificationChannels: process.Monitoring.NotificationChannels,
		SyncIntervalHours:    process.Monitoring.SyncIntervalHours,
		EnabledBy:            enabledBy,
	}
}

// ProcessMonitoringDisabledEvent evento quando monitoramento é desabilitado
type ProcessMonitoringDisabledEvent struct {
	BaseDomainEvent
	ProcessNumber string `json:"process_number"`
	DisabledBy    string `json:"disabled_by"`
}

func NewProcessMonitoringDisabledEvent(process *Process, disabledBy string) *ProcessMonitoringDisabledEvent {
	return &ProcessMonitoringDisabledEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.monitoring.disabled",
			Version:    "1.0",
			ID:         process.ID,
			OccurredAt: time.Now(),
			TenantID:   process.TenantID,
		},
		ProcessNumber: process.Number,
		DisabledBy:    disabledBy,
	}
}

// ProcessSyncedEvent evento quando um processo é sincronizado com DataJud
type ProcessSyncedEvent struct {
	BaseDomainEvent
	ProcessNumber    string    `json:"process_number"`
	NewMovements     int       `json:"new_movements"`
	UpdatedMovements int       `json:"updated_movements"`
	LastSyncAt       time.Time `json:"last_sync_at"`
	SyncSource       string    `json:"sync_source"` // datajud, manual
}

func NewProcessSyncedEvent(process *Process, newMovements, updatedMovements int, syncSource string) *ProcessSyncedEvent {
	return &ProcessSyncedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.synced",
			Version:    "1.0",
			ID:         process.ID,
			OccurredAt: time.Now(),
			TenantID:   process.TenantID,
		},
		ProcessNumber:    process.Number,
		NewMovements:     newMovements,
		UpdatedMovements: updatedMovements,
		LastSyncAt:       *process.LastSyncAt,
		SyncSource:       syncSource,
	}
}

// MovementCreatedEvent evento quando uma movimentação é criada
type MovementCreatedEvent struct {
	BaseDomainEvent
	ProcessNumber string       `json:"process_number"`
	MovementID    string       `json:"movement_id"`
	Sequence      int          `json:"sequence"`
	Type          MovementType `json:"type"`
	Title         string       `json:"title"`
	Date          time.Time    `json:"date"`
	IsImportant   bool         `json:"is_important"`
	Source        string       `json:"source"` // datajud, manual, import
}

func NewMovementCreatedEvent(movement *Movement, processNumber, source string) *MovementCreatedEvent {
	return &MovementCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "movement.created",
			Version:    "1.0",
			ID:         movement.ID,
			OccurredAt: time.Now(),
			TenantID:   movement.TenantID,
		},
		ProcessNumber: processNumber,
		MovementID:    movement.ID,
		Sequence:      movement.Sequence,
		Type:          movement.Type,
		Title:         movement.Title,
		Date:          movement.Date,
		IsImportant:   movement.IsImportant,
		Source:        source,
	}
}

// MovementAnalyzedEvent evento quando uma movimentação é analisada
type MovementAnalyzedEvent struct {
	BaseDomainEvent
	MovementID      string  `json:"movement_id"`
	ProcessNumber   string  `json:"process_number"`
	Sentiment       string  `json:"sentiment"`
	Importance      int     `json:"importance"`
	Category        string  `json:"category"`
	RequiresAction  bool    `json:"requires_action"`
	ActionType      string  `json:"action_type"`
	Confidence      float64 `json:"confidence"`
	ProcessedBy     string  `json:"processed_by"`
}

func NewMovementAnalyzedEvent(movement *Movement, processNumber string) *MovementAnalyzedEvent {
	return &MovementAnalyzedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "movement.analyzed",
			Version:    "1.0",
			ID:         movement.ID,
			OccurredAt: time.Now(),
			TenantID:   movement.TenantID,
		},
		MovementID:     movement.ID,
		ProcessNumber:  processNumber,
		Sentiment:      movement.Metadata.Analysis.Sentiment,
		Importance:     movement.Metadata.Analysis.Importance,
		Category:       movement.Metadata.Analysis.Category,
		RequiresAction: movement.Metadata.Analysis.RequiresAction,
		ActionType:     movement.Metadata.Analysis.ActionType,
		Confidence:     movement.Metadata.Analysis.Confidence,
		ProcessedBy:    movement.Metadata.Analysis.ProcessedBy,
	}
}

// ImportantMovementDetectedEvent evento quando movimentação importante é detectada
type ImportantMovementDetectedEvent struct {
	BaseDomainEvent
	ProcessNumber        string   `json:"process_number"`
	MovementID           string   `json:"movement_id"`
	MovementTitle        string   `json:"movement_title"`
	MovementType         string   `json:"movement_type"`
	ImportanceLevel      int      `json:"importance_level"`
	NotificationChannels []string `json:"notification_channels"`
	Keywords             []string `json:"keywords"`
}

func NewImportantMovementDetectedEvent(movement *Movement, processNumber string, notificationChannels []string) *ImportantMovementDetectedEvent {
	return &ImportantMovementDetectedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "movement.important.detected",
			Version:    "1.0",
			ID:         movement.ID,
			OccurredAt: time.Now(),
			TenantID:   movement.TenantID,
		},
		ProcessNumber:        processNumber,
		MovementID:           movement.ID,
		MovementTitle:        movement.Title,
		MovementType:         string(movement.Type),
		ImportanceLevel:      movement.Metadata.Analysis.Importance,
		NotificationChannels: notificationChannels,
		Keywords:             movement.Metadata.Keywords,
	}
}

// PartyAddedEvent evento quando uma parte é adicionada ao processo
type PartyAddedEvent struct {
	BaseDomainEvent
	ProcessNumber string    `json:"process_number"`
	PartyID       string    `json:"party_id"`
	PartyName     string    `json:"party_name"`
	PartyType     PartyType `json:"party_type"`
	PartyRole     PartyRole `json:"party_role"`
	AddedBy       string    `json:"added_by"`
}

func NewPartyAddedEvent(party *Party, processNumber, addedBy string) *PartyAddedEvent {
	return &PartyAddedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "party.added",
			Version:    "1.0",
			ID:         party.ID,
			OccurredAt: time.Now(),
			TenantID:   "", // Será preenchido pelo serviço
		},
		ProcessNumber: processNumber,
		PartyID:       party.ID,
		PartyName:     party.Name,
		PartyType:     party.Type,
		PartyRole:     party.Role,
		AddedBy:       addedBy,
	}
}

// PartyUpdatedEvent evento quando uma parte é atualizada
type PartyUpdatedEvent struct {
	BaseDomainEvent
	ProcessNumber string                 `json:"process_number"`
	PartyID       string                 `json:"party_id"`
	PartyName     string                 `json:"party_name"`
	Changes       map[string]interface{} `json:"changes"`
	UpdatedBy     string                 `json:"updated_by"`
}

func NewPartyUpdatedEvent(partyID, processNumber, partyName string, changes map[string]interface{}, updatedBy string) *PartyUpdatedEvent {
	return &PartyUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "party.updated",
			Version:    "1.0",
			ID:         partyID,
			OccurredAt: time.Now(),
			TenantID:   "", // Será preenchido pelo serviço
		},
		ProcessNumber: processNumber,
		PartyID:       partyID,
		PartyName:     partyName,
		Changes:       changes,
		UpdatedBy:     updatedBy,
	}
}

// ProcessBatchSyncStartedEvent evento quando sincronização em lote inicia
type ProcessBatchSyncStartedEvent struct {
	BaseDomainEvent
	BatchID       string   `json:"batch_id"`
	ProcessCount  int      `json:"process_count"`
	ProcessIDs    []string `json:"process_ids"`
	SyncType      string   `json:"sync_type"` // full, incremental, manual
	StartedBy     string   `json:"started_by"`
}

func NewProcessBatchSyncStartedEvent(tenantID, batchID string, processIDs []string, syncType, startedBy string) *ProcessBatchSyncStartedEvent {
	return &ProcessBatchSyncStartedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.batch.sync.started",
			Version:    "1.0",
			ID:         batchID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		BatchID:      batchID,
		ProcessCount: len(processIDs),
		ProcessIDs:   processIDs,
		SyncType:     syncType,
		StartedBy:    startedBy,
	}
}

// ProcessBatchSyncCompletedEvent evento quando sincronização em lote é concluída
type ProcessBatchSyncCompletedEvent struct {
	BaseDomainEvent
	BatchID           string        `json:"batch_id"`
	ProcessCount      int           `json:"process_count"`
	SuccessCount      int           `json:"success_count"`
	ErrorCount        int           `json:"error_count"`
	NewMovements      int           `json:"new_movements"`
	UpdatedMovements  int           `json:"updated_movements"`
	Duration          time.Duration `json:"duration"`
	CompletedAt       time.Time     `json:"completed_at"`
}

func NewProcessBatchSyncCompletedEvent(tenantID, batchID string, processCount, successCount, errorCount, newMovements, updatedMovements int, duration time.Duration) *ProcessBatchSyncCompletedEvent {
	return &ProcessBatchSyncCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "process.batch.sync.completed",
			Version:    "1.0",
			ID:         batchID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		BatchID:          batchID,
		ProcessCount:     processCount,
		SuccessCount:     successCount,
		ErrorCount:       errorCount,
		NewMovements:     newMovements,
		UpdatedMovements: updatedMovements,
		Duration:         duration,
		CompletedAt:      time.Now(),
	}
}

// EventPublisher interface para publicação de eventos
type EventPublisher interface {
	Publish(event DomainEvent) error
	PublishBatch(events []DomainEvent) error
}