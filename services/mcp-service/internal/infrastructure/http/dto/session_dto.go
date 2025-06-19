package dto

import (
	"time"
)

// CreateSessionRequest request para criar sessão MCP
type CreateSessionRequest struct {
	Channel   string                 `json:"channel" binding:"required"` // whatsapp, telegram, slack, web
	UserID    string                 `json:"user_id" binding:"required"`
	TenantID  string                 `json:"tenant_id" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Settings  SessionSettings        `json:"settings,omitempty"`
}

// SessionSettings configurações da sessão
type SessionSettings struct {
	ClaudeModel     string `json:"claude_model,omitempty"`
	MaxTokens       int    `json:"max_tokens,omitempty"`
	Timeout         int    `json:"timeout,omitempty"`
	AutoSave        bool   `json:"auto_save,omitempty"`
	Notifications   bool   `json:"notifications,omitempty"`
}

// SessionResponse response com dados da sessão
type SessionResponse struct {
	ID           string                 `json:"id"`
	Channel      string                 `json:"channel"`
	UserID       string                 `json:"user_id"`
	TenantID     string                 `json:"tenant_id"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	LastActivity time.Time              `json:"last_activity"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Settings     SessionSettings        `json:"settings"`
	Context      ConversationContext    `json:"context"`
}

// ConversationContext contexto da conversa
type ConversationContext struct {
	MessagesCount   int                    `json:"messages_count"`
	TokensUsed      int                    `json:"tokens_used"`
	LastTopic       string                 `json:"last_topic,omitempty"`
	ActiveTools     []string               `json:"active_tools,omitempty"`
	CurrentProcess  string                 `json:"current_process,omitempty"`
	ContextData     map[string]interface{} `json:"context_data,omitempty"`
}

// SendMessageRequest request para enviar mensagem
type SendMessageRequest struct {
	Message    string                 `json:"message" binding:"required"`
	MessageType string                `json:"message_type,omitempty"` // text, image, document
	Attachments []MessageAttachment   `json:"attachments,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Options    MessageOptions         `json:"options,omitempty"`
}

// MessageAttachment anexo da mensagem
type MessageAttachment struct {
	Type string `json:"type"` // image, document, audio
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

// MessageOptions opções da mensagem
type MessageOptions struct {
	UseTools       bool     `json:"use_tools,omitempty"`
	AllowedTools   []string `json:"allowed_tools,omitempty"`
	Stream         bool     `json:"stream,omitempty"`
	SaveToHistory  bool     `json:"save_to_history,omitempty"`
}

// MessageResponse response da mensagem
type MessageResponse struct {
	ID          string                 `json:"id"`
	SessionID   string                 `json:"session_id"`
	Role        string                 `json:"role"` // user, assistant, system
	Content     string                 `json:"content"`
	Timestamp   time.Time              `json:"timestamp"`
	TokensUsed  int                    `json:"tokens_used"`
	ToolCalls   []ToolCall            `json:"tool_calls,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ToolCall chamada de ferramenta
type ToolCall struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Args     map[string]interface{} `json:"args"`
	Result   interface{}            `json:"result,omitempty"`
	Status   string                 `json:"status"` // pending, running, completed, failed
	Duration time.Duration          `json:"duration,omitempty"`
	Error    string                 `json:"error,omitempty"`
}

// ConversationHistoryResponse histórico da conversa
type ConversationHistoryResponse struct {
	SessionID    string            `json:"session_id"`
	Messages     []MessageResponse `json:"messages"`
	TotalCount   int               `json:"total_count"`
	Page         int               `json:"page"`
	PageSize     int               `json:"page_size"`
	HasMore      bool              `json:"has_more"`
}

// SessionStatusResponse status da sessão
type SessionStatusResponse struct {
	ID           string              `json:"id"`
	Status       string              `json:"status"`
	IsActive     bool                `json:"is_active"`
	LastActivity time.Time           `json:"last_activity"`
	Context      ConversationContext `json:"context"`
	QuotaUsage   QuotaUsage         `json:"quota_usage"`
}

// QuotaUsage uso de quota
type QuotaUsage struct {
	TokensUsed    int     `json:"tokens_used"`
	TokensLimit   int     `json:"tokens_limit"`
	RequestsUsed  int     `json:"requests_used"`
	RequestsLimit int     `json:"requests_limit"`
	UsagePercent  float64 `json:"usage_percent"`
}