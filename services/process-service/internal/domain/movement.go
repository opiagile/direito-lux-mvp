package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Movement representa uma movimentação/andamento do processo
type Movement struct {
	ID             string           `json:"id" db:"id"`
	ProcessID      string           `json:"process_id" db:"process_id"`
	TenantID       string           `json:"tenant_id" db:"tenant_id"`
	Sequence       int              `json:"sequence" db:"sequence"` // Ordem sequencial
	ExternalID     string           `json:"external_id" db:"external_id"` // ID do DataJud
	Date           time.Time        `json:"date" db:"date"`
	Type           MovementType     `json:"type" db:"type"`
	Code           string           `json:"code" db:"code"` // Código CNJ da movimentação
	Title          string           `json:"title" db:"title"`
	Description    string           `json:"description" db:"description"`
	Content        string           `json:"content" db:"content"` // Texto completo
	Judge          string           `json:"judge" db:"judge"`
	Responsible    string           `json:"responsible" db:"responsible"` // Responsável pela movimentação
	Attachments    []Attachment     `json:"attachments"`
	RelatedParties []string         `json:"related_parties" db:"related_parties"` // IDs das partes relacionadas
	IsImportant    bool             `json:"is_important" db:"is_important"`
	IsPublic       bool             `json:"is_public" db:"is_public"`
	NotificationSent bool           `json:"notification_sent" db:"notification_sent"`
	Tags           []string         `json:"tags" db:"tags"`
	Metadata       MovementMetadata `json:"metadata" db:"metadata"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
	SyncedAt       time.Time        `json:"synced_at" db:"synced_at"`
}

// Attachment representa um anexo da movimentação
type Attachment struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Type        string    `json:"type"` // pdf, doc, image, etc.
	URL         string    `json:"url"`
	ExternalID  string    `json:"external_id"` // ID no DataJud
	IsDownloaded bool     `json:"is_downloaded"`
	CreatedAt   time.Time `json:"created_at"`
}

// MovementMetadata metadados adicionais da movimentação
type MovementMetadata struct {
	OriginalSource string            `json:"original_source"` // datajud, manual, import
	DataJudID      string            `json:"datajud_id"`
	ImportBatch    string            `json:"import_batch"`
	Keywords       []string          `json:"keywords"`
	Analysis       MovementAnalysis  `json:"analysis"`
	CustomFields   map[string]string `json:"custom_fields"`
}

// MovementAnalysis análise automática da movimentação
type MovementAnalysis struct {
	Sentiment      string    `json:"sentiment"`       // positive, negative, neutral
	Importance     int       `json:"importance"`      // 1-5, sendo 5 mais importante
	Category       string    `json:"category"`        // decisao, despacho, ato, etc.
	HasDeadline    bool      `json:"has_deadline"`
	DeadlineDate   *time.Time `json:"deadline_date"`
	RequiresAction bool      `json:"requires_action"`
	ActionType     string    `json:"action_type"`     // recurso, manifestacao, etc.
	Confidence     float64   `json:"confidence"`      // 0-1, confiança da análise
	ProcessedBy    string    `json:"processed_by"`    // ai, rule_based, manual
	ProcessedAt    time.Time `json:"processed_at"`
}

// MovementType tipo de movimentação
type MovementType string

const (
	MovementTypeDecision      MovementType = "decision"       // Decisão
	MovementTypeOrder         MovementType = "order"          // Despacho
	MovementTypeAct           MovementType = "act"            // Ato processual
	MovementTypeFiling        MovementType = "filing"         // Peticionamento
	MovementTypeHearing       MovementType = "hearing"        // Audiência
	MovementTypePublication   MovementType = "publication"    // Publicação
	MovementTypeRemittance    MovementType = "remittance"     // Remessa
	MovementTypeReturn        MovementType = "return"         // Retorno
	MovementTypeCitation      MovementType = "citation"       // Citação
	MovementTypeIntimation    MovementType = "intimation"     // Intimação
	MovementTypeArchiving     MovementType = "archiving"      // Arquivamento
	MovementTypeReactivation  MovementType = "reactivation"   // Desarquivamento
	MovementTypeAppeal        MovementType = "appeal"         // Recurso
	MovementTypeExecution     MovementType = "execution"      // Execução
	MovementTypeOther         MovementType = "other"          // Outros
)

// MovementRepository interface para persistência de movimentações
type MovementRepository interface {
	// Comandos (Write)
	Create(movement *Movement) error
	CreateBatch(movements []*Movement) error
	Update(movement *Movement) error
	Delete(id string) error
	
	// Queries (Read)
	GetByID(id string) (*Movement, error)
	GetByProcess(processID string, filters MovementFilters) ([]*Movement, error)
	GetByTenant(tenantID string, filters MovementFilters) ([]*Movement, error)
	GetByExternalID(externalID string) (*Movement, error)
	GetImportantMovements(processID string) ([]*Movement, error)
	GetRecentMovements(tenantID string, limit int) ([]*Movement, error)
	GetPendingNotifications(tenantID string) ([]*Movement, error)
	
	// Análise e busca
	SearchByContent(tenantID string, query string, filters MovementFilters) ([]*Movement, error)
	GetMovementsByKeywords(tenantID string, keywords []string) ([]*Movement, error)
	GetStatistics(tenantID string, dateFrom, dateTo time.Time) (*MovementStatistics, error)
}

// MovementFilters filtros para consultas de movimentações
type MovementFilters struct {
	ProcessID     string         `json:"process_id"`
	Type          []MovementType `json:"type"`
	IsImportant   *bool          `json:"is_important"`
	IsPublic      *bool          `json:"is_public"`
	HasNotification *bool        `json:"has_notification"`
	DateFrom      *time.Time     `json:"date_from"`
	DateTo        *time.Time     `json:"date_to"`
	Judge         string         `json:"judge"`
	Tags          []string       `json:"tags"`
	Search        string         `json:"search"`
	Limit         int            `json:"limit"`
	Offset        int            `json:"offset"`
	SortBy        string         `json:"sort_by"`    // date, sequence, importance
	SortOrder     string         `json:"sort_order"` // asc, desc
}

// MovementStatistics estatísticas de movimentações
type MovementStatistics struct {
	TotalMovements       int                        `json:"total_movements"`
	MovementsByType      map[MovementType]int       `json:"movements_by_type"`
	MovementsByMonth     map[string]int             `json:"movements_by_month"`
	AveragePerProcess    float64                    `json:"average_per_process"`
	ImportantMovements   int                        `json:"important_movements"`
	MostActiveProcesses  []ProcessActivity          `json:"most_active_processes"`
	TopKeywords          []KeywordCount             `json:"top_keywords"`
	NotificationsSent    int                        `json:"notifications_sent"`
}

// ProcessActivity atividade por processo
type ProcessActivity struct {
	ProcessID      string `json:"process_id"`
	ProcessNumber  string `json:"process_number"`
	MovementCount  int    `json:"movement_count"`
	LastMovement   time.Time `json:"last_movement"`
}

// KeywordCount contagem de palavras-chave
type KeywordCount struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

// Erros de domínio para movimentações
var (
	ErrMovementNotFound      = errors.New("movimentação não encontrada")
	ErrMovementExists        = errors.New("movimentação já existe")
	ErrInvalidMovementType   = errors.New("tipo de movimentação inválido")
	ErrInvalidSequence       = errors.New("sequência inválida")
	ErrDuplicateExternalID   = errors.New("ID externo duplicado")
	ErrInvalidAttachment     = errors.New("anexo inválido")
	ErrMovementInFuture      = errors.New("data da movimentação não pode ser futura")
)

// ValidateType valida o tipo da movimentação
func (m *Movement) ValidateType() error {
	validTypes := []MovementType{
		MovementTypeDecision, MovementTypeOrder, MovementTypeAct,
		MovementTypeFiling, MovementTypeHearing, MovementTypePublication,
		MovementTypeRemittance, MovementTypeReturn, MovementTypeCitation,
		MovementTypeIntimation, MovementTypeArchiving, MovementTypeReactivation,
		MovementTypeAppeal, MovementTypeExecution, MovementTypeOther,
	}
	
	for _, validType := range validTypes {
		if m.Type == validType {
			return nil
		}
	}
	
	return ErrInvalidMovementType
}

// ValidateDate valida a data da movimentação
func (m *Movement) ValidateDate() error {
	if m.Date.After(time.Now()) {
		return ErrMovementInFuture
	}
	return nil
}

// ValidateSequence valida a sequência da movimentação
func (m *Movement) ValidateSequence() error {
	if m.Sequence <= 0 {
		return ErrInvalidSequence
	}
	return nil
}

// IsDecision verifica se é uma decisão
func (m *Movement) IsDecision() bool {
	return m.Type == MovementTypeDecision
}

// IsAppeal verifica se é um recurso
func (m *Movement) IsAppeal() bool {
	return m.Type == MovementTypeAppeal
}

// IsHearing verifica se é uma audiência
func (m *Movement) IsHearing() bool {
	return m.Type == MovementTypeHearing
}

// HasAttachments verifica se tem anexos
func (m *Movement) HasAttachments() bool {
	return len(m.Attachments) > 0
}

// GetDisplayTitle retorna título para exibição
func (m *Movement) GetDisplayTitle() string {
	if m.Title != "" {
		return m.Title
	}
	return "Movimentação sem título"
}

// GetFormattedDate retorna data formatada
func (m *Movement) GetFormattedDate() string {
	return m.Date.Format("02/01/2006 15:04")
}

// GetSummary retorna resumo da movimentação
func (m *Movement) GetSummary() string {
	if len(m.Description) <= 150 {
		return m.Description
	}
	return m.Description[:147] + "..."
}

// MarkAsImportant marca como importante
func (m *Movement) MarkAsImportant() {
	m.IsImportant = true
	m.UpdatedAt = time.Now()
}

// UnmarkAsImportant desmarca como importante
func (m *Movement) UnmarkAsImportant() {
	m.IsImportant = false
	m.UpdatedAt = time.Now()
}

// AddTag adiciona uma tag
func (m *Movement) AddTag(tag string) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if tag == "" {
		return
	}
	
	// Verifica se a tag já existe
	for _, existingTag := range m.Tags {
		if existingTag == tag {
			return
		}
	}
	
	m.Tags = append(m.Tags, tag)
	m.UpdatedAt = time.Now()
}

// RemoveTag remove uma tag
func (m *Movement) RemoveTag(tag string) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	
	for i, existingTag := range m.Tags {
		if existingTag == tag {
			m.Tags = append(m.Tags[:i], m.Tags[i+1:]...)
			m.UpdatedAt = time.Now()
			return
		}
	}
}

// MarkNotificationSent marca notificação como enviada
func (m *Movement) MarkNotificationSent() {
	m.NotificationSent = true
	m.UpdatedAt = time.Now()
}

// AddAttachment adiciona um anexo
func (m *Movement) AddAttachment(name, url, externalID string, size int64, fileType string) {
	attachment := Attachment{
		ID:          generateUUID(), // Função auxiliar para gerar UUID
		Name:        name,
		Size:        size,
		Type:        fileType,
		URL:         url,
		ExternalID:  externalID,
		IsDownloaded: false,
		CreatedAt:   time.Now(),
	}
	
	m.Attachments = append(m.Attachments, attachment)
	m.UpdatedAt = time.Now()
}

// SetAnalysis define análise automática
func (m *Movement) SetAnalysis(sentiment string, importance int, category string, confidence float64, processedBy string) {
	m.Metadata.Analysis = MovementAnalysis{
		Sentiment:   sentiment,
		Importance:  importance,
		Category:    category,
		Confidence:  confidence,
		ProcessedBy: processedBy,
		ProcessedAt: time.Now(),
	}
	m.UpdatedAt = time.Now()
}

// SetCustomField define um campo customizado
func (m *Movement) SetCustomField(key, value string) {
	if m.Metadata.CustomFields == nil {
		m.Metadata.CustomFields = make(map[string]string)
	}
	
	m.Metadata.CustomFields[key] = value
	m.UpdatedAt = time.Now()
}

// GetImportanceLevel retorna nível de importância como string
func (m *Movement) GetImportanceLevel() string {
	switch m.Metadata.Analysis.Importance {
	case 5:
		return "Crítica"
	case 4:
		return "Alta"
	case 3:
		return "Média"
	case 2:
		return "Baixa"
	case 1:
		return "Mínima"
	default:
		return "Não analisada"
	}
}

// RequiresNotification verifica se requer notificação
func (m *Movement) RequiresNotification() bool {
	return m.IsImportant && !m.NotificationSent
}

// GetKeywords extrai palavras-chave da movimentação
func (m *Movement) GetKeywords() []string {
	if len(m.Metadata.Keywords) > 0 {
		return m.Metadata.Keywords
	}
	
	// Extração simples de palavras-chave do conteúdo
	text := strings.ToLower(m.Title + " " + m.Description)
	words := strings.Fields(text)
	
	keywords := []string{}
	for _, word := range words {
		if len(word) > 4 { // Apenas palavras com mais de 4 caracteres
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// generateUUID função auxiliar para gerar UUID (simplificada)
func generateUUID() string {
	// Implementação simplificada - em produção usar uuid.New().String()
	return fmt.Sprintf("mov_%d", time.Now().UnixNano())
}