package domain

import (
	"time"
	"errors"
	"regexp"
	"strings"
	"fmt"
)

// Process representa um processo jurídico
type Process struct {
	ID                  string              `json:"id" db:"id"`
	TenantID            string              `json:"tenant_id" db:"tenant_id"`
	ClientID            string              `json:"client_id" db:"client_id"`
	Number              string              `json:"number" db:"number"` // Número CNJ
	OriginalNumber      string              `json:"original_number" db:"original_number"` // Número original antes da unificação
	Title               string              `json:"title" db:"title"`
	Description         string              `json:"description" db:"description"`
	Status              ProcessStatus       `json:"status" db:"status"`
	Stage               ProcessStage        `json:"stage" db:"stage"`
	Subject             ProcessSubject      `json:"subject" db:"subject"`
	Value               *ProcessValue       `json:"value" db:"value"`
	CourtID             string              `json:"court_id" db:"court_id"`
	JudgeID             *string             `json:"judge_id" db:"judge_id"`
	Parties             []Party             `json:"parties"`
	Monitoring          ProcessMonitoring   `json:"monitoring" db:"monitoring"`
	Tags                []string            `json:"tags" db:"tags"`
	CustomFields        map[string]string   `json:"custom_fields" db:"custom_fields"`
	LastMovementAt      *time.Time          `json:"last_movement_at" db:"last_movement_at"`
	LastSyncAt          *time.Time          `json:"last_sync_at" db:"last_sync_at"`
	CreatedAt           time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at" db:"updated_at"`
	ArchivedAt          *time.Time          `json:"archived_at" db:"archived_at"`
}

// ProcessValue representa o valor da causa
type ProcessValue struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// ProcessMonitoring configurações de monitoramento
type ProcessMonitoring struct {
	Enabled              bool     `json:"enabled"`
	NotificationChannels []string `json:"notification_channels"` // whatsapp, email, telegram
	Keywords             []string `json:"keywords"`              // Palavras-chave para alertas
	AutoSync             bool     `json:"auto_sync"`
	SyncIntervalHours    int      `json:"sync_interval_hours"`
	LastNotificationAt   *time.Time `json:"last_notification_at"`
}

// ProcessStatus status do processo
type ProcessStatus string

const (
	ProcessStatusActive     ProcessStatus = "active"      // Ativo/Em andamento
	ProcessStatusSuspended  ProcessStatus = "suspended"   // Suspenso
	ProcessStatusFinished   ProcessStatus = "finished"    // Finalizado
	ProcessStatusArchived   ProcessStatus = "archived"    // Arquivado
	ProcessStatusCanceled   ProcessStatus = "canceled"    // Cancelado
)

// ProcessStage fase processual
type ProcessStage string

const (
	ProcessStageKnowledge     ProcessStage = "knowledge"     // Conhecimento
	ProcessStageExecution     ProcessStage = "execution"     // Execução
	ProcessStageAppeal        ProcessStage = "appeal"        // Recursal
	ProcessStageEmergency     ProcessStage = "emergency"     // Cautelar/Urgência
	ProcessStagePreventive    ProcessStage = "preventive"    // Preventiva
	ProcessStageIncidental    ProcessStage = "incidental"    // Incidental
)

// ProcessSubject assunto do processo (baseado em tabela CNJ)
type ProcessSubject struct {
	Code        string `json:"code"`        // Código CNJ do assunto
	Description string `json:"description"` // Descrição do assunto
	ParentCode  string `json:"parent_code"` // Código do assunto pai
}

// ProcessRepository interface para persistência de processos
type ProcessRepository interface {
	// Comandos (Write)
	Create(process *Process) error
	Update(process *Process) error
	Delete(id string) error
	Archive(id string) error
	
	// Queries (Read)
	GetByID(id string) (*Process, error)
	GetByNumber(number string) (*Process, error)
	GetByTenant(tenantID string, filters ProcessFilters) ([]*Process, error)
	GetByClient(clientID string, filters ProcessFilters) ([]*Process, error)
	GetActiveForMonitoring() ([]*Process, error)
	GetNeedingSync() ([]*Process, error)
	
	// Estatísticas
	CountByTenant(tenantID string) (int, error)
	CountByStatus(tenantID string, status ProcessStatus) (int, error)
}

// ProcessFilters filtros para consultas
type ProcessFilters struct {
	Status    []ProcessStatus `json:"status"`
	Stage     []ProcessStage  `json:"stage"`
	CourtID   string          `json:"court_id"`
	JudgeID   string          `json:"judge_id"`
	Tags      []string        `json:"tags"`
	DateFrom  *time.Time      `json:"date_from"`
	DateTo    *time.Time      `json:"date_to"`
	Search    string          `json:"search"`    // Busca textual
	Limit     int             `json:"limit"`
	Offset    int             `json:"offset"`
	SortBy    string          `json:"sort_by"`   // created_at, updated_at, last_movement_at
	SortOrder string          `json:"sort_order"` // asc, desc
}

// Erros de domínio
var (
	ErrProcessNotFound        = errors.New("processo não encontrado")
	ErrProcessExists          = errors.New("processo já existe")
	ErrInvalidProcessNumber   = errors.New("número de processo inválido")
	ErrInvalidProcessStatus   = errors.New("status de processo inválido")
	ErrInvalidProcessStage    = errors.New("fase processual inválida")
	ErrProcessAlreadyArchived = errors.New("processo já está arquivado")
	ErrCannotModifyArchived   = errors.New("não é possível modificar processo arquivado")
	ErrInvalidCourt           = errors.New("tribunal inválido")
	ErrTenantQuotaExceeded    = errors.New("quota de processos do tenant excedida")
)

// ValidateNumber valida o número CNJ do processo
func (p *Process) ValidateNumber() error {
	if p.Number == "" {
		return ErrInvalidProcessNumber
	}
	
	// Remove caracteres não numéricos para validação
	cleanNumber := regexp.MustCompile(`[^0-9]`).ReplaceAllString(p.Number, "")
	
	// Número CNJ deve ter 20 dígitos
	if len(cleanNumber) != 20 {
		return ErrInvalidProcessNumber
	}
	
	// Valida padrão CNJ: NNNNNNN-DD.AAAA.J.TR.OOOO
	cnpjPattern := regexp.MustCompile(`^\d{7}-\d{2}\.\d{4}\.\d{1}\.\d{2}\.\d{4}$`)
	if !cnpjPattern.MatchString(p.Number) {
		return ErrInvalidProcessNumber
	}
	
	return nil
}

// ValidateStatus valida o status do processo
func (p *Process) ValidateStatus() error {
	validStatuses := []ProcessStatus{
		ProcessStatusActive, ProcessStatusSuspended, ProcessStatusFinished,
		ProcessStatusArchived, ProcessStatusCanceled,
	}
	
	for _, status := range validStatuses {
		if p.Status == status {
			return nil
		}
	}
	
	return ErrInvalidProcessStatus
}

// ValidateStage valida a fase processual
func (p *Process) ValidateStage() error {
	validStages := []ProcessStage{
		ProcessStageKnowledge, ProcessStageExecution, ProcessStageAppeal,
		ProcessStageEmergency, ProcessStagePreventive, ProcessStageIncidental,
	}
	
	for _, stage := range validStages {
		if p.Stage == stage {
			return nil
		}
	}
	
	return ErrInvalidProcessStage
}

// IsActive verifica se o processo está ativo
func (p *Process) IsActive() bool {
	return p.Status == ProcessStatusActive
}

// IsArchived verifica se o processo está arquivado
func (p *Process) IsArchived() bool {
	return p.Status == ProcessStatusArchived || p.ArchivedAt != nil
}

// CanBeModified verifica se o processo pode ser modificado
func (p *Process) CanBeModified() bool {
	return !p.IsArchived()
}

// Archive arquiva o processo
func (p *Process) Archive() error {
	if p.IsArchived() {
		return ErrProcessAlreadyArchived
	}
	
	p.Status = ProcessStatusArchived
	now := time.Now()
	p.ArchivedAt = &now
	p.UpdatedAt = now
	
	return nil
}

// Reactivate reativa um processo arquivado
func (p *Process) Reactivate() error {
	if !p.IsArchived() {
		return errors.New("processo não está arquivado")
	}
	
	p.Status = ProcessStatusActive
	p.ArchivedAt = nil
	p.UpdatedAt = time.Now()
	
	return nil
}

// UpdateLastMovement atualiza data da última movimentação
func (p *Process) UpdateLastMovement() {
	now := time.Now()
	p.LastMovementAt = &now
	p.UpdatedAt = now
}

// UpdateLastSync atualiza data da última sincronização
func (p *Process) UpdateLastSync() {
	now := time.Now()
	p.LastSyncAt = &now
	p.UpdatedAt = now
}

// AddTag adiciona uma tag ao processo
func (p *Process) AddTag(tag string) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if tag == "" {
		return
	}
	
	// Verifica se a tag já existe
	for _, existingTag := range p.Tags {
		if existingTag == tag {
			return
		}
	}
	
	p.Tags = append(p.Tags, tag)
	p.UpdatedAt = time.Now()
}

// RemoveTag remove uma tag do processo
func (p *Process) RemoveTag(tag string) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	
	for i, existingTag := range p.Tags {
		if existingTag == tag {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			p.UpdatedAt = time.Now()
			return
		}
	}
}

// SetCustomField define um campo customizado
func (p *Process) SetCustomField(key, value string) {
	if p.CustomFields == nil {
		p.CustomFields = make(map[string]string)
	}
	
	p.CustomFields[key] = value
	p.UpdatedAt = time.Now()
}

// GetCustomField obtém um campo customizado
func (p *Process) GetCustomField(key string) (string, bool) {
	if p.CustomFields == nil {
		return "", false
	}
	
	value, exists := p.CustomFields[key]
	return value, exists
}

// EnableMonitoring habilita monitoramento automático
func (p *Process) EnableMonitoring(channels []string, syncIntervalHours int) {
	p.Monitoring.Enabled = true
	p.Monitoring.NotificationChannels = channels
	p.Monitoring.AutoSync = true
	p.Monitoring.SyncIntervalHours = syncIntervalHours
	p.UpdatedAt = time.Now()
}

// DisableMonitoring desabilita monitoramento automático
func (p *Process) DisableMonitoring() {
	p.Monitoring.Enabled = false
	p.Monitoring.AutoSync = false
	p.UpdatedAt = time.Now()
}

// AddKeyword adiciona palavra-chave para alertas
func (p *Process) AddKeyword(keyword string) {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return
	}
	
	// Verifica se a palavra-chave já existe
	for _, existingKeyword := range p.Monitoring.Keywords {
		if existingKeyword == keyword {
			return
		}
	}
	
	p.Monitoring.Keywords = append(p.Monitoring.Keywords, keyword)
	p.UpdatedAt = time.Now()
}

// NeedsSync verifica se o processo precisa ser sincronizado
func (p *Process) NeedsSync() bool {
	if !p.Monitoring.Enabled || !p.Monitoring.AutoSync {
		return false
	}
	
	if p.LastSyncAt == nil {
		return true
	}
	
	syncInterval := time.Duration(p.Monitoring.SyncIntervalHours) * time.Hour
	return time.Since(*p.LastSyncAt) >= syncInterval
}

// GetDisplayTitle retorna título para exibição
func (p *Process) GetDisplayTitle() string {
	if p.Title != "" {
		return p.Title
	}
	return p.Subject.Description
}

// GetFormattedValue retorna valor formatado da causa
func (p *Process) GetFormattedValue() string {
	if p.Value == nil {
		return "Não informado"
	}
	
	if p.Value.Currency == "BRL" {
		return fmt.Sprintf("R$ %.2f", p.Value.Amount)
	}
	
	return fmt.Sprintf("%.2f %s", p.Value.Amount, p.Value.Currency)
}