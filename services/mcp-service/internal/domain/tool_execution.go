package domain

import (
	"time"

	"github.com/google/uuid"
)

// ToolExecution representa uma execução de ferramenta
type ToolExecution struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	SessionID   uuid.UUID              `json:"session_id" db:"session_id"`
	ToolName    string                 `json:"tool_name" db:"tool_name"`
	Status      ExecutionStatus        `json:"status" db:"status"`
	Parameters  map[string]interface{} `json:"parameters" db:"parameters"`
	Result      interface{}            `json:"result,omitempty" db:"result"`
	Error       string                 `json:"error,omitempty" db:"error"`
	StartedAt   time.Time              `json:"started_at" db:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	Duration    time.Duration          `json:"duration,omitempty" db:"duration"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// ExecutionStatus status de execução da ferramenta
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// ParameterDefinition definição de parâmetro
type ParameterDefinition struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Minimum     *float64    `json:"minimum,omitempty"`
	Maximum     *float64    `json:"maximum,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
}

// ToolExample exemplo de uso da ferramenta
type ToolExample struct {
	Description    string                 `json:"description"`
	Parameters     map[string]interface{} `json:"parameters"`
	ExpectedResult interface{}            `json:"expected_result,omitempty"`
}

// NewToolExecution cria nova execução de ferramenta
func NewToolExecution(sessionID uuid.UUID, toolName string, parameters map[string]interface{}) *ToolExecution {
	now := time.Now()
	return &ToolExecution{
		ID:         uuid.New(),
		SessionID:  sessionID,
		ToolName:   toolName,
		Status:     ExecutionStatusPending,
		Parameters: parameters,
		StartedAt:  now,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Start marca a execução como iniciada
func (e *ToolExecution) Start() {
	e.Status = ExecutionStatusRunning
	e.StartedAt = time.Now()
	e.UpdatedAt = time.Now()
}

// Complete marca a execução como concluída
func (e *ToolExecution) Complete(result interface{}) {
	e.Status = ExecutionStatusCompleted
	e.Result = result
	now := time.Now()
	e.CompletedAt = &now
	e.Duration = now.Sub(e.StartedAt)
	e.UpdatedAt = now
}

// Fail marca a execução como falhada
func (e *ToolExecution) Fail(errorMsg string) {
	e.Status = ExecutionStatusFailed
	e.Error = errorMsg
	now := time.Now()
	e.CompletedAt = &now
	e.Duration = now.Sub(e.StartedAt)
	e.UpdatedAt = now
}

// Cancel cancela a execução
func (e *ToolExecution) Cancel() {
	e.Status = ExecutionStatusCancelled
	now := time.Now()
	e.CompletedAt = &now
	e.Duration = now.Sub(e.StartedAt)
	e.UpdatedAt = now
}

// IsCompleted verifica se a execução foi concluída
func (e *ToolExecution) IsCompleted() bool {
	return e.Status == ExecutionStatusCompleted || 
		   e.Status == ExecutionStatusFailed || 
		   e.Status == ExecutionStatusCancelled
}

// SetMetadata define um valor de metadados
func (e *ToolExecution) SetMetadata(key string, value interface{}) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	e.UpdatedAt = time.Now()
}