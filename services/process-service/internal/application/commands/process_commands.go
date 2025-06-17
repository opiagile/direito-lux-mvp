package commands

import (
	"time"
	"github.com/direito-lux/process-service/internal/domain"
)

// Command representa um comando na arquitetura CQRS
type Command interface {
	Validate() error
}

// CreateProcessCommand comando para criar processo
type CreateProcessCommand struct {
	TenantID         string                        `json:"tenant_id" validate:"required"`
	ClientID         string                        `json:"client_id" validate:"required"`
	Number           string                        `json:"number" validate:"required"`
	OriginalNumber   string                        `json:"original_number"`
	Title            string                        `json:"title" validate:"required,min=5,max=200"`
	Description      string                        `json:"description" validate:"max=1000"`
	Status           domain.ProcessStatus          `json:"status" validate:"required"`
	Stage            domain.ProcessStage           `json:"stage" validate:"required"`
	Subject          CreateProcessSubjectCommand   `json:"subject" validate:"required"`
	Value            *CreateProcessValueCommand    `json:"value"`
	CourtID          string                        `json:"court_id" validate:"required"`
	JudgeID          *string                       `json:"judge_id"`
	Parties          []CreatePartyCommand          `json:"parties"`
	Monitoring       CreateMonitoringCommand       `json:"monitoring"`
	Tags             []string                      `json:"tags"`
	CustomFields     map[string]string             `json:"custom_fields"`
	CreatedBy        string                        `json:"created_by" validate:"required"`
}

// CreateProcessSubjectCommand comando para assunto do processo
type CreateProcessSubjectCommand struct {
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	ParentCode  string `json:"parent_code"`
}

// CreateProcessValueCommand comando para valor da causa
type CreateProcessValueCommand struct {
	Amount   float64 `json:"amount" validate:"min=0"`
	Currency string  `json:"currency" validate:"required,len=3"`
}

// CreateMonitoringCommand comando para configuração de monitoramento
type CreateMonitoringCommand struct {
	Enabled              bool     `json:"enabled"`
	NotificationChannels []string `json:"notification_channels"`
	Keywords             []string `json:"keywords"`
	AutoSync             bool     `json:"auto_sync"`
	SyncIntervalHours    int      `json:"sync_interval_hours" validate:"min=1,max=168"`
}

// CreatePartyCommand comando para criar parte
type CreatePartyCommand struct {
	Type         domain.PartyType    `json:"type" validate:"required"`
	Name         string              `json:"name" validate:"required,min=2,max=200"`
	Document     string              `json:"document"`
	DocumentType string              `json:"document_type"`
	Role         domain.PartyRole    `json:"role" validate:"required"`
	IsActive     bool                `json:"is_active"`
	Lawyer       *CreateLawyerCommand `json:"lawyer"`
	Contact      CreateContactCommand `json:"contact"`
	Address      CreateAddressCommand `json:"address"`
}

// CreateLawyerCommand comando para criar advogado
type CreateLawyerCommand struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	OAB      string `json:"oab" validate:"required"`
	OABState string `json:"oab_state" validate:"required,len=2"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
}

// CreateContactCommand comando para criar contato
type CreateContactCommand struct {
	Email     string `json:"email" validate:"email"`
	Phone     string `json:"phone"`
	CellPhone string `json:"cell_phone"`
	Website   string `json:"website" validate:"url"`
}

// CreateAddressCommand comando para criar endereço
type CreateAddressCommand struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state" validate:"len=2"`
	ZipCode    string `json:"zip_code"`
	Country    string `json:"country"`
}

// UpdateProcessCommand comando para atualizar processo
type UpdateProcessCommand struct {
	ID               string                        `json:"id" validate:"required"`
	TenantID         string                        `json:"tenant_id" validate:"required"`
	Title            *string                       `json:"title,omitempty"`
	Description      *string                       `json:"description,omitempty"`
	Status           *domain.ProcessStatus         `json:"status,omitempty"`
	Stage            *domain.ProcessStage          `json:"stage,omitempty"`
	Subject          *CreateProcessSubjectCommand  `json:"subject,omitempty"`
	Value            *CreateProcessValueCommand    `json:"value,omitempty"`
	JudgeID          *string                       `json:"judge_id,omitempty"`
	Tags             []string                      `json:"tags,omitempty"`
	CustomFields     map[string]string             `json:"custom_fields,omitempty"`
	UpdatedBy        string                        `json:"updated_by" validate:"required"`
}

// ArchiveProcessCommand comando para arquivar processo
type ArchiveProcessCommand struct {
	ID         string `json:"id" validate:"required"`
	TenantID   string `json:"tenant_id" validate:"required"`
	Reason     string `json:"reason" validate:"required,min=5,max=200"`
	ArchivedBy string `json:"archived_by" validate:"required"`
}

// ReactivateProcessCommand comando para reativar processo
type ReactivateProcessCommand struct {
	ID            string `json:"id" validate:"required"`
	TenantID      string `json:"tenant_id" validate:"required"`
	ReactivatedBy string `json:"reactivated_by" validate:"required"`
}

// EnableMonitoringCommand comando para habilitar monitoramento
type EnableMonitoringCommand struct {
	ProcessID            string   `json:"process_id" validate:"required"`
	TenantID             string   `json:"tenant_id" validate:"required"`
	NotificationChannels []string `json:"notification_channels" validate:"required,min=1"`
	SyncIntervalHours    int      `json:"sync_interval_hours" validate:"min=1,max=168"`
	Keywords             []string `json:"keywords"`
	EnabledBy            string   `json:"enabled_by" validate:"required"`
}

// DisableMonitoringCommand comando para desabilitar monitoramento
type DisableMonitoringCommand struct {
	ProcessID  string `json:"process_id" validate:"required"`
	TenantID   string `json:"tenant_id" validate:"required"`
	DisabledBy string `json:"disabled_by" validate:"required"`
}

// SyncProcessCommand comando para sincronizar processo
type SyncProcessCommand struct {
	ProcessID  string `json:"process_id" validate:"required"`
	TenantID   string `json:"tenant_id" validate:"required"`
	ForceSync  bool   `json:"force_sync"`
	SyncSource string `json:"sync_source" validate:"required"` // datajud, manual
	SyncedBy   string `json:"synced_by" validate:"required"`
}

// BatchSyncProcessesCommand comando para sincronização em lote
type BatchSyncProcessesCommand struct {
	TenantID     string   `json:"tenant_id" validate:"required"`
	ProcessIDs   []string `json:"process_ids" validate:"required,min=1"`
	SyncType     string   `json:"sync_type" validate:"required"` // full, incremental
	ForceSync    bool     `json:"force_sync"`
	StartedBy    string   `json:"started_by" validate:"required"`
}

// CreateMovementCommand comando para criar movimentação
type CreateMovementCommand struct {
	ProcessID      string               `json:"process_id" validate:"required"`
	TenantID       string               `json:"tenant_id" validate:"required"`
	ExternalID     string               `json:"external_id"`
	Date           time.Time            `json:"date" validate:"required"`
	Type           domain.MovementType  `json:"type" validate:"required"`
	Code           string               `json:"code" validate:"required"`
	Title          string               `json:"title" validate:"required,min=5,max=200"`
	Description    string               `json:"description" validate:"required,min=10"`
	Content        string               `json:"content"`
	Judge          string               `json:"judge"`
	Responsible    string               `json:"responsible"`
	RelatedParties []string             `json:"related_parties"`
	IsImportant    bool                 `json:"is_important"`
	IsPublic       bool                 `json:"is_public"`
	Tags           []string             `json:"tags"`
	Source         string               `json:"source" validate:"required"` // datajud, manual, import
	CreatedBy      string               `json:"created_by" validate:"required"`
}

// UpdateMovementCommand comando para atualizar movimentação
type UpdateMovementCommand struct {
	ID          string               `json:"id" validate:"required"`
	TenantID    string               `json:"tenant_id" validate:"required"`
	Title       *string              `json:"title,omitempty"`
	Description *string              `json:"description,omitempty"`
	Content     *string              `json:"content,omitempty"`
	Judge       *string              `json:"judge,omitempty"`
	Responsible *string              `json:"responsible,omitempty"`
	IsImportant *bool                `json:"is_important,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	UpdatedBy   string               `json:"updated_by" validate:"required"`
}

// AnalyzeMovementCommand comando para analisar movimentação
type AnalyzeMovementCommand struct {
	MovementID  string  `json:"movement_id" validate:"required"`
	TenantID    string  `json:"tenant_id" validate:"required"`
	Sentiment   string  `json:"sentiment" validate:"required"`
	Importance  int     `json:"importance" validate:"min=1,max=5"`
	Category    string  `json:"category" validate:"required"`
	Confidence  float64 `json:"confidence" validate:"min=0,max=1"`
	ProcessedBy string  `json:"processed_by" validate:"required"`
}

// AddPartyCommand comando para adicionar parte ao processo
type AddPartyCommand struct {
	ProcessID string               `json:"process_id" validate:"required"`
	TenantID  string               `json:"tenant_id" validate:"required"`
	Party     CreatePartyCommand   `json:"party" validate:"required"`
	AddedBy   string               `json:"added_by" validate:"required"`
}

// UpdatePartyCommand comando para atualizar parte
type UpdatePartyCommand struct {
	ID           string               `json:"id" validate:"required"`
	ProcessID    string               `json:"process_id" validate:"required"`
	TenantID     string               `json:"tenant_id" validate:"required"`
	Name         *string              `json:"name,omitempty"`
	Document     *string              `json:"document,omitempty"`
	DocumentType *string              `json:"document_type,omitempty"`
	IsActive     *bool                `json:"is_active,omitempty"`
	Lawyer       *CreateLawyerCommand `json:"lawyer,omitempty"`
	Contact      *CreateContactCommand `json:"contact,omitempty"`
	Address      *CreateAddressCommand `json:"address,omitempty"`
	UpdatedBy    string               `json:"updated_by" validate:"required"`
}

// RemovePartyCommand comando para remover parte
type RemovePartyCommand struct {
	ID        string `json:"id" validate:"required"`
	ProcessID string `json:"process_id" validate:"required"`
	TenantID  string `json:"tenant_id" validate:"required"`
	RemovedBy string `json:"removed_by" validate:"required"`
}

// Validate implementações de validação básica

func (c *CreateProcessCommand) Validate() error {
	// Validações específicas do domínio podem ser implementadas aqui
	// Por enquanto, a validação estrutural é feita via tags
	return nil
}

func (c *UpdateProcessCommand) Validate() error {
	return nil
}

func (c *ArchiveProcessCommand) Validate() error {
	return nil
}

func (c *ReactivateProcessCommand) Validate() error {
	return nil
}

func (c *EnableMonitoringCommand) Validate() error {
	return nil
}

func (c *DisableMonitoringCommand) Validate() error {
	return nil
}

func (c *SyncProcessCommand) Validate() error {
	return nil
}

func (c *BatchSyncProcessesCommand) Validate() error {
	return nil
}

func (c *CreateMovementCommand) Validate() error {
	return nil
}

func (c *UpdateMovementCommand) Validate() error {
	return nil
}

func (c *AnalyzeMovementCommand) Validate() error {
	return nil
}

func (c *AddPartyCommand) Validate() error {
	return nil
}

func (c *UpdatePartyCommand) Validate() error {
	return nil
}

func (c *RemovePartyCommand) Validate() error {
	return nil
}