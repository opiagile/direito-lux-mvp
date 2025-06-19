package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventType representa o tipo de evento
type EventType string

const (
	// Session Events
	EventSessionCreated  EventType = "mcp.session.created"
	EventSessionExpired  EventType = "mcp.session.expired"
	EventSessionClosed   EventType = "mcp.session.closed"
	
	// Message Events
	EventMessageReceived EventType = "mcp.message.received"
	EventMessageSent     EventType = "mcp.message.sent"
	
	// Tool Events
	EventToolExecuted    EventType = "mcp.tool.executed"
	EventToolFailed      EventType = "mcp.tool.failed"
	
	// Quota Events
	EventQuotaExceeded   EventType = "mcp.quota.exceeded"
	EventQuotaWarning    EventType = "mcp.quota.warning"
	
	// Bot Events
	EventBotConnected    EventType = "mcp.bot.connected"
	EventBotDisconnected EventType = "mcp.bot.disconnected"
	EventBotError        EventType = "mcp.bot.error"
)

// Event representa um evento de domínio
type Event struct {
	ID           uuid.UUID              `json:"id"`
	Type         EventType              `json:"type"`
	AggregateID  uuid.UUID              `json:"aggregate_id"`
	TenantID     uuid.UUID              `json:"tenant_id"`
	UserID       uuid.UUID              `json:"user_id"`
	Data         map[string]interface{} `json:"data"`
	Metadata     map[string]interface{} `json:"metadata"`
	OccurredAt   time.Time              `json:"occurred_at"`
	ProcessedAt  *time.Time             `json:"processed_at,omitempty"`
}

// NewEvent cria um novo evento
func NewEvent(eventType EventType, aggregateID, tenantID, userID uuid.UUID) *Event {
	return &Event{
		ID:          uuid.New(),
		Type:        eventType,
		AggregateID: aggregateID,
		TenantID:    tenantID,
		UserID:      userID,
		Data:        make(map[string]interface{}),
		Metadata:    make(map[string]interface{}),
		OccurredAt:  time.Now(),
	}
}

// SessionCreatedEvent representa o evento de criação de sessão
type SessionCreatedEvent struct {
	*Event
	Channel    string `json:"channel"`
	ExternalID string `json:"external_id"`
}

// NewSessionCreatedEvent cria um evento de sessão criada
func NewSessionCreatedEvent(session *MCPSession) *SessionCreatedEvent {
	event := NewEvent(EventSessionCreated, session.ID, session.TenantID, session.UserID)
	event.Data["channel"] = session.Channel
	event.Data["external_id"] = session.ExternalID
	
	return &SessionCreatedEvent{
		Event:      event,
		Channel:    session.Channel,
		ExternalID: session.ExternalID,
	}
}

// MessageReceivedEvent representa o evento de mensagem recebida
type MessageReceivedEvent struct {
	*Event
	SessionID    uuid.UUID `json:"session_id"`
	Channel      string    `json:"channel"`
	Message      string    `json:"message"`
	ExternalData interface{} `json:"external_data,omitempty"`
}

// NewMessageReceivedEvent cria um evento de mensagem recebida
func NewMessageReceivedEvent(sessionID, tenantID, userID uuid.UUID, channel, message string) *MessageReceivedEvent {
	event := NewEvent(EventMessageReceived, sessionID, tenantID, userID)
	event.Data["channel"] = channel
	event.Data["message"] = message
	
	return &MessageReceivedEvent{
		Event:     event,
		SessionID: sessionID,
		Channel:   channel,
		Message:   message,
	}
}

// ToolExecutedEvent representa o evento de ferramenta executada
type ToolExecutedEvent struct {
	*Event
	SessionID    uuid.UUID              `json:"session_id"`
	ToolName     string                 `json:"tool_name"`
	Parameters   map[string]interface{} `json:"parameters"`
	Result       interface{}            `json:"result,omitempty"`
	Duration     time.Duration          `json:"duration"`
	Success      bool                   `json:"success"`
	Error        string                 `json:"error,omitempty"`
}

// NewToolExecutedEvent cria um evento de ferramenta executada
func NewToolExecutedEvent(sessionID, tenantID, userID uuid.UUID, toolName string, params map[string]interface{}, result interface{}, duration time.Duration, err error) *ToolExecutedEvent {
	event := NewEvent(EventToolExecuted, sessionID, tenantID, userID)
	
	success := err == nil
	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}
	
	event.Data["tool_name"] = toolName
	event.Data["parameters"] = params
	event.Data["result"] = result
	event.Data["duration_ms"] = duration.Milliseconds()
	event.Data["success"] = success
	event.Data["error"] = errorMsg
	
	return &ToolExecutedEvent{
		Event:      event,
		SessionID:  sessionID,
		ToolName:   toolName,
		Parameters: params,
		Result:     result,
		Duration:   duration,
		Success:    success,
		Error:      errorMsg,
	}
}

// QuotaExceededEvent representa o evento de quota excedida
type QuotaExceededEvent struct {
	*Event
	Plan          string `json:"plan"`
	QuotaType     string `json:"quota_type"`
	CurrentUsage  int    `json:"current_usage"`
	QuotaLimit    int    `json:"quota_limit"`
}

// NewQuotaExceededEvent cria um evento de quota excedida
func NewQuotaExceededEvent(tenantID, userID uuid.UUID, plan, quotaType string, currentUsage, quotaLimit int) *QuotaExceededEvent {
	event := NewEvent(EventQuotaExceeded, tenantID, tenantID, userID)
	event.Data["plan"] = plan
	event.Data["quota_type"] = quotaType
	event.Data["current_usage"] = currentUsage
	event.Data["quota_limit"] = quotaLimit
	
	return &QuotaExceededEvent{
		Event:        event,
		Plan:         plan,
		QuotaType:    quotaType,
		CurrentUsage: currentUsage,
		QuotaLimit:   quotaLimit,
	}
}

// BotErrorEvent representa um erro do bot
type BotErrorEvent struct {
	*Event
	Channel      string `json:"channel"`
	ErrorMessage string `json:"error_message"`
	ErrorType    string `json:"error_type"`
	SessionID    *uuid.UUID `json:"session_id,omitempty"`
}

// NewBotErrorEvent cria um evento de erro do bot
func NewBotErrorEvent(tenantID, userID uuid.UUID, channel, errorMessage, errorType string) *BotErrorEvent {
	event := NewEvent(EventBotError, uuid.New(), tenantID, userID)
	event.Data["channel"] = channel
	event.Data["error_message"] = errorMessage
	event.Data["error_type"] = errorType
	
	return &BotErrorEvent{
		Event:        event,
		Channel:      channel,
		ErrorMessage: errorMessage,
		ErrorType:    errorType,
	}
}

// Generic event types used by services

// SessionEvent representa um evento genérico de sessão
type SessionEvent struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id"`
	Channel   string                 `json:"channel"`
	UserID    string                 `json:"user_id"`
	TenantID  string                 `json:"tenant_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// MCPEvent representa um evento genérico do MCP
type MCPEvent struct {
	Type      string                 `json:"type"`
	ID        string                 `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	UserID    string                 `json:"user_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ToolEvent representa um evento genérico de ferramenta
type ToolEvent struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id"`
	ToolName  string                 `json:"tool_name"`
	UserID    string                 `json:"user_id"`
	TenantID  string                 `json:"tenant_id"`
	Timestamp time.Time              `json:"timestamp"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Data      map[string]interface{} `json:"data"`
}