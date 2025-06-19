package domain

import (
	"time"

	"github.com/google/uuid"
)

// MCPSession representa uma sessão de conversação com contexto
type MCPSession struct {
	ID              uuid.UUID              `json:"id" db:"id"`
	TenantID        uuid.UUID              `json:"tenant_id" db:"tenant_id"`
	UserID          uuid.UUID              `json:"user_id" db:"user_id"`
	Channel         string                 `json:"channel" db:"channel"` // whatsapp, telegram, claude, slack
	ExternalID      string                 `json:"external_id" db:"external_id"` // ID externo do chat
	State           SessionState           `json:"state" db:"state"`
	Context         map[string]interface{} `json:"context" db:"context"`
	LastInteraction time.Time              `json:"last_interaction" db:"last_interaction"`
	MessageCount    int                    `json:"message_count" db:"message_count"`
	CommandCount    int                    `json:"command_count" db:"command_count"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	ExpiresAt       time.Time              `json:"expires_at" db:"expires_at"`
}

// SessionState representa o estado da sessão
type SessionState string

const (
	SessionStateActive   SessionState = "active"
	SessionStateIdle     SessionState = "idle"
	SessionStateExpired  SessionState = "expired"
	SessionStateClosed   SessionState = "closed"
)

// NewMCPSession cria uma nova sessão MCP
func NewMCPSession(tenantID, userID uuid.UUID, channel, externalID string) *MCPSession {
	now := time.Now()
	return &MCPSession{
		ID:              uuid.New(),
		TenantID:        tenantID,
		UserID:          userID,
		Channel:         channel,
		ExternalID:      externalID,
		State:           SessionStateActive,
		Context:         make(map[string]interface{}),
		LastInteraction: now,
		MessageCount:    0,
		CommandCount:    0,
		CreatedAt:       now,
		UpdatedAt:       now,
		ExpiresAt:       now.Add(30 * time.Minute), // 30 minutos de sessão
	}
}

// UpdateInteraction atualiza a última interação da sessão
func (s *MCPSession) UpdateInteraction() {
	s.LastInteraction = time.Now()
	s.UpdatedAt = time.Now()
	s.MessageCount++
	
	// Renovar sessão se estiver próxima de expirar
	if time.Until(s.ExpiresAt) < 10*time.Minute {
		s.ExpiresAt = time.Now().Add(30 * time.Minute)
	}
}

// IncrementCommandCount incrementa o contador de comandos
func (s *MCPSession) IncrementCommandCount() {
	s.CommandCount++
	s.UpdatedAt = time.Now()
}

// IsExpired verifica se a sessão expirou
func (s *MCPSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt) || s.State == SessionStateExpired
}

// Close fecha a sessão
func (s *MCPSession) Close() {
	s.State = SessionStateClosed
	s.UpdatedAt = time.Now()
}

// SetContext define um valor no contexto da sessão
func (s *MCPSession) SetContext(key string, value interface{}) {
	if s.Context == nil {
		s.Context = make(map[string]interface{})
	}
	s.Context[key] = value
	s.UpdatedAt = time.Now()
}

// GetContext obtém um valor do contexto da sessão
func (s *MCPSession) GetContext(key string) (interface{}, bool) {
	if s.Context == nil {
		return nil, false
	}
	value, exists := s.Context[key]
	return value, exists
}