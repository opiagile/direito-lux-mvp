package domain

import (
	"time"

	"github.com/google/uuid"
)

// ConversationContext representa o contexto de uma conversa
type ConversationContext struct {
	ID              uuid.UUID        `json:"id" db:"id"`
	SessionID       uuid.UUID        `json:"session_id" db:"session_id"`
	CurrentTool     string           `json:"current_tool" db:"current_tool"`
	CurrentStep     string           `json:"current_step" db:"current_step"`
	WaitingFor      string           `json:"waiting_for" db:"waiting_for"`
	Variables       map[string]interface{} `json:"variables" db:"variables"`
	History         []Message        `json:"history" db:"history"`
	PendingActions  []PendingAction  `json:"pending_actions" db:"pending_actions"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
}

// Message representa uma mensagem no histórico
type Message struct {
	ID          uuid.UUID `json:"id"`
	Role        string    `json:"role"` // user, assistant, system
	Content     string    `json:"content"`
	ToolCalls   []ToolCall `json:"tool_calls,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// ToolCall representa uma chamada de ferramenta
type ToolCall struct {
	ID         string                 `json:"id"`
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters"`
	Result     interface{}            `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	ExecutedAt time.Time              `json:"executed_at"`
}

// PendingAction representa uma ação pendente
type PendingAction struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	CreatedAt   time.Time              `json:"created_at"`
}

// NewConversationContext cria um novo contexto de conversa
func NewConversationContext(sessionID uuid.UUID) *ConversationContext {
	now := time.Now()
	return &ConversationContext{
		ID:             uuid.New(),
		SessionID:      sessionID,
		CurrentTool:    "",
		CurrentStep:    "",
		WaitingFor:     "",
		Variables:      make(map[string]interface{}),
		History:        []Message{},
		PendingActions: []PendingAction{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// AddMessage adiciona uma mensagem ao histórico
func (c *ConversationContext) AddMessage(role, content string) {
	message := Message{
		ID:        uuid.New(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
	c.History = append(c.History, message)
	c.UpdatedAt = time.Now()
}

// AddToolCall adiciona uma chamada de ferramenta ao histórico
func (c *ConversationContext) AddToolCall(toolName string, parameters map[string]interface{}) string {
	toolCall := ToolCall{
		ID:         uuid.New().String(),
		ToolName:   toolName,
		Parameters: parameters,
		ExecutedAt: time.Now(),
	}
	
	// Adicionar à última mensagem do assistente
	if len(c.History) > 0 && c.History[len(c.History)-1].Role == "assistant" {
		c.History[len(c.History)-1].ToolCalls = append(c.History[len(c.History)-1].ToolCalls, toolCall)
	}
	
	c.UpdatedAt = time.Now()
	return toolCall.ID
}

// UpdateToolCallResult atualiza o resultado de uma chamada de ferramenta
func (c *ConversationContext) UpdateToolCallResult(toolCallID string, result interface{}, err error) {
	for i := len(c.History) - 1; i >= 0; i-- {
		for j := range c.History[i].ToolCalls {
			if c.History[i].ToolCalls[j].ID == toolCallID {
				c.History[i].ToolCalls[j].Result = result
				if err != nil {
					c.History[i].ToolCalls[j].Error = err.Error()
				}
				c.UpdatedAt = time.Now()
				return
			}
		}
	}
}

// SetVariable define uma variável no contexto
func (c *ConversationContext) SetVariable(key string, value interface{}) {
	if c.Variables == nil {
		c.Variables = make(map[string]interface{})
	}
	c.Variables[key] = value
	c.UpdatedAt = time.Now()
}

// GetVariable obtém uma variável do contexto
func (c *ConversationContext) GetVariable(key string) (interface{}, bool) {
	if c.Variables == nil {
		return nil, false
	}
	value, exists := c.Variables[key]
	return value, exists
}

// AddPendingAction adiciona uma ação pendente
func (c *ConversationContext) AddPendingAction(actionType, description string, parameters map[string]interface{}) {
	action := PendingAction{
		ID:          uuid.New(),
		Type:        actionType,
		Description: description,
		Parameters:  parameters,
		CreatedAt:   time.Now(),
	}
	c.PendingActions = append(c.PendingActions, action)
	c.UpdatedAt = time.Now()
}

// ClearPendingActions limpa as ações pendentes
func (c *ConversationContext) ClearPendingActions() {
	c.PendingActions = []PendingAction{}
	c.UpdatedAt = time.Now()
}

// GetRecentHistory obtém as N mensagens mais recentes
func (c *ConversationContext) GetRecentHistory(n int) []Message {
	if len(c.History) <= n {
		return c.History
	}
	return c.History[len(c.History)-n:]
}