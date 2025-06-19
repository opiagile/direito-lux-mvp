package dto

import (
	"time"
)

// ExecuteToolRequest request para executar ferramenta
type ExecuteToolRequest struct {
	SessionID   string                 `json:"session_id" binding:"required"`
	ToolName    string                 `json:"tool_name" binding:"required"`
	Parameters  map[string]interface{} `json:"parameters" binding:"required"`
	Async       bool                   `json:"async,omitempty"`
	Timeout     int                    `json:"timeout,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ToolExecutionResponse response da execução
type ToolExecutionResponse struct {
	ID          string                 `json:"id"`
	SessionID   string                 `json:"session_id"`
	ToolName    string                 `json:"tool_name"`
	Status      string                 `json:"status"` // pending, running, completed, failed
	Result      interface{}            `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ToolDefinition definição de ferramenta
type ToolDefinition struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Category    string                `json:"category"`
	Parameters  ToolParameterSchema   `json:"parameters"`
	Permissions []string              `json:"permissions,omitempty"`
	QuotaCost   int                   `json:"quota_cost,omitempty"`
	Timeout     time.Duration         `json:"timeout,omitempty"`
	Async       bool                  `json:"async"`
	Examples    []ToolExample        `json:"examples,omitempty"`
}

// ToolParameterSchema schema dos parâmetros
type ToolParameterSchema struct {
	Type       string                         `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
	Required   []string                      `json:"required,omitempty"`
}

// ParameterProperty propriedade do parâmetro
type ParameterProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Minimum     *float64    `json:"minimum,omitempty"`
	Maximum     *float64    `json:"maximum,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
}

// ToolExample exemplo de uso da ferramenta
type ToolExample struct {
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	ExpectedResult interface{}         `json:"expected_result,omitempty"`
}

// ListToolsResponse lista de ferramentas
type ListToolsResponse struct {
	Tools    []ToolDefinition `json:"tools"`
	Count    int              `json:"count"`
	Category []string         `json:"categories"`
}

// ToolExecutionsListResponse lista de execuções
type ToolExecutionsListResponse struct {
	Executions []ToolExecutionResponse `json:"executions"`
	TotalCount int                      `json:"total_count"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	HasMore    bool                     `json:"has_more"`
}

// ToolUsageStats estatísticas de uso de ferramentas
type ToolUsageStats struct {
	ToolName      string        `json:"tool_name"`
	TotalUses     int           `json:"total_uses"`
	SuccessRate   float64       `json:"success_rate"`
	AvgDuration   time.Duration `json:"avg_duration"`
	LastUsed      time.Time     `json:"last_used"`
	QuotaConsumed int           `json:"quota_consumed"`
}