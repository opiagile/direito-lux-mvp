package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/domain"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http/dto"
	"github.com/direito-lux/mcp-service/internal/infrastructure/events"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
)

// SessionServiceFixed serviço para gerenciar sessões MCP (versão corrigida)
type SessionServiceFixed struct {
	logger    *zap.Logger
	metrics   *metrics.Metrics
	eventBus  *events.EventBus
	sessions  map[string]*domain.MCPSession // Mapear por string ID para compatibilidade com DTOs
}

// NewSessionServiceFixed cria nova instância do serviço
func NewSessionServiceFixed(
	logger *zap.Logger,
	metrics *metrics.Metrics,
	eventBus *events.EventBus,
) *SessionServiceFixed {
	return &SessionServiceFixed{
		logger:   logger,
		metrics:  metrics,
		eventBus: eventBus,
		sessions: make(map[string]*domain.MCPSession),
	}
}

// CreateSession cria uma nova sessão MCP
func (s *SessionServiceFixed) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	// Converter strings para UUIDs
	tenantUUID, err := uuid.Parse(req.TenantID)
	if err != nil {
		// Se não for UUID válido, gerar um novo
		tenantUUID = uuid.New()
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		// Se não for UUID válido, gerar um novo
		userUUID = uuid.New()
	}

	// Criar sessão usando o construtor do domain
	session := domain.NewMCPSession(tenantUUID, userUUID, req.Channel, "")

	// Mapear por string ID para compatibilidade
	sessionIDStr := session.ID.String()
	s.sessions[sessionIDStr] = session

	// Publicar evento
	event := domain.SessionEvent{
		Type:      "session_created",
		SessionID: sessionIDStr,
		Channel:   req.Channel,
		UserID:    req.UserID,
		TenantID:  req.TenantID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"settings": req.Settings,
			"metadata": req.Metadata,
		},
	}

	if s.eventBus != nil {
		if err := s.eventBus.PublishSessionEvent(ctx, event); err != nil {
			s.logger.Warn("Erro ao publicar evento de sessão", zap.Error(err))
		}
	}

	// Registrar métricas
	if s.metrics != nil {
		s.metrics.RecordMCPSession(req.Channel, req.TenantID, true)
		s.metrics.RecordMCPConversation(req.Channel, req.TenantID, "initiated")
	}

	s.logger.Info("Sessão MCP criada",
		zap.String("session_id", sessionIDStr),
		zap.String("channel", req.Channel),
		zap.String("user_id", req.UserID),
		zap.String("tenant_id", req.TenantID),
	)

	return s.mapSessionToResponse(session), nil
}

// GetSession obtém sessão por ID
func (s *SessionServiceFixed) GetSession(ctx context.Context, sessionID string) (*dto.SessionResponse, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	return s.mapSessionToResponse(session), nil
}

// CloseSession fecha uma sessão
func (s *SessionServiceFixed) CloseSession(ctx context.Context, sessionID string) error {
	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	// Atualizar estado
	session.State = domain.SessionStateClosed
	session.UpdatedAt = time.Now()

	// Publicar evento
	event := domain.SessionEvent{
		Type:      "session_closed",
		SessionID: sessionID,
		Channel:   session.Channel,
		UserID:    session.UserID.String(),
		TenantID:  session.TenantID.String(),
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"duration":        time.Since(session.CreatedAt).Seconds(),
			"messages_count":  session.MessageCount,
			"commands_count":  session.CommandCount,
		},
	}

	if s.eventBus != nil {
		if err := s.eventBus.PublishSessionEvent(ctx, event); err != nil {
			s.logger.Warn("Erro ao publicar evento de fechamento", zap.Error(err))
		}
	}

	// Registrar métricas
	if s.metrics != nil {
		s.metrics.RecordMCPSession(session.Channel, session.TenantID.String(), false)
		s.metrics.RecordMCPConversation(session.Channel, session.TenantID.String(), "completed")
	}

	s.logger.Info("Sessão MCP fechada",
		zap.String("session_id", sessionID),
		zap.String("channel", session.Channel),
		zap.Duration("duration", time.Since(session.CreatedAt)),
	)

	// Remover da memória
	delete(s.sessions, sessionID)

	return nil
}

// mapSessionToResponse converte sessão do domínio para DTO
func (s *SessionServiceFixed) mapSessionToResponse(session *domain.MCPSession) *dto.SessionResponse {
	return &dto.SessionResponse{
		ID:           session.ID.String(),
		Channel:      session.Channel,
		UserID:       session.UserID.String(),
		TenantID:     session.TenantID.String(),
		Status:       string(session.State),
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
		LastActivity: session.LastInteraction,
		Metadata:     session.Context,
		Settings: dto.SessionSettings{
			ClaudeModel:   "claude-3-sonnet-20241022", // Default
			MaxTokens:     4096,                       // Default
			Timeout:       30,                         // Default
			AutoSave:      true,
			Notifications: true,
		},
		Context: dto.ConversationContext{
			MessagesCount: session.MessageCount,
			TokensUsed:    session.CommandCount * 100, // Estimativa
			ContextData:   session.Context,
		},
	}
}