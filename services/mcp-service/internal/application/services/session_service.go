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

// SessionService serviço para gerenciar sessões MCP
type SessionService struct {
	logger    *zap.Logger
	metrics   *metrics.Metrics
	eventBus  *events.EventBus
	sessions  map[string]*domain.MCPSession // Simulação em memória por agora
}

// NewSessionService cria nova instância do serviço
func NewSessionService(
	logger *zap.Logger,
	metrics *metrics.Metrics,
	eventBus *events.EventBus,
) *SessionService {
	return &SessionService{
		logger:   logger,
		metrics:  metrics,
		eventBus: eventBus,
		sessions: make(map[string]*domain.MCPSession),
	}
}

// CreateSession cria uma nova sessão MCP
func (s *SessionService) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	// Converter strings para UUIDs
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	// Criar sessão usando o construtor do domain
	session := domain.NewMCPSession(tenantID, userID, req.Channel, "")
	
	// Adicionar metadata ao contexto se existir
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			session.SetContext(key, value)
		}
	}

	// Salvar sessão (simulação em memória)
	sessionID := session.ID.String()
	s.sessions[sessionID] = session

	// Publicar evento
	event := domain.SessionEvent{
		Type:      "session_created",
		SessionID: sessionID,
		Channel:   req.Channel,
		UserID:    req.UserID,
		TenantID:  req.TenantID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"settings": req.Settings,
			"metadata": req.Metadata,
		},
	}

	if err := s.eventBus.PublishSessionEvent(ctx, event); err != nil {
		s.logger.Warn("Erro ao publicar evento de sessão", zap.Error(err))
	}

	// Registrar métricas
	if s.metrics != nil {
		s.metrics.RecordMCPSession(req.Channel, req.TenantID, true)
		s.metrics.RecordMCPConversation(req.Channel, req.TenantID, "initiated")
	}

	s.logger.Info("Sessão MCP criada",
		zap.String("session_id", sessionID),
		zap.String("channel", req.Channel),
		zap.String("user_id", req.UserID),
		zap.String("tenant_id", req.TenantID),
	)

	return s.mapSessionToResponse(session), nil
}

// GetSession obtém sessão por ID
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*dto.SessionResponse, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	return s.mapSessionToResponse(session), nil
}

// CloseSession fecha uma sessão
func (s *SessionService) CloseSession(ctx context.Context, sessionID string) error {
	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	// Atualizar status
	session.Close()
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
			"duration": time.Since(session.CreatedAt).Seconds(),
			"messages_count": session.MessageCount,
			"command_count": session.CommandCount,
		},
	}

	if err := s.eventBus.PublishSessionEvent(ctx, event); err != nil {
		s.logger.Warn("Erro ao publicar evento de fechamento", zap.Error(err))
	}

	// Registrar métricas
	if s.metrics != nil {
		s.metrics.RecordMCPSession(session.TenantID.String(), session.Channel, false)
		s.metrics.RecordMCPConversation(session.TenantID.String(), session.Channel, "session_closed")
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

// GetSessionStatus obtém status da sessão
func (s *SessionService) GetSessionStatus(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	return &dto.SessionStatusResponse{
		ID:           session.ID.String(),
		Status:       string(session.State),
		IsActive:     session.State == domain.SessionStateActive,
		LastActivity: session.LastInteraction,
		Context: dto.ConversationContext{
			MessagesCount:  session.MessageCount,
			TokensUsed:     0, // TODO: implementar contagem de tokens
			LastTopic:      "", // TODO: extrair do contexto se necessário
			ActiveTools:    []string{}, // TODO: extrair do contexto se necessário
			CurrentProcess: "", // TODO: extrair do contexto se necessário
			ContextData:    session.Context,
		},
		QuotaUsage: dto.QuotaUsage{
			TokensUsed:    0, // TODO: implementar contagem de tokens
			TokensLimit:   10000, // TODO: Buscar do sistema de quotas
			RequestsUsed:  session.MessageCount,
			RequestsLimit: 100, // TODO: Buscar do sistema de quotas
			UsagePercent:  0, // TODO: calcular baseado em dados reais
		},
	}, nil
}

// SendMessage envia mensagem para a sessão
func (s *SessionService) SendMessage(ctx context.Context, sessionID string, req dto.SendMessageRequest) (*dto.MessageResponse, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	if session.State != domain.SessionStateActive {
		return nil, fmt.Errorf("sessão não está ativa")
	}

	// Criar mensagem
	messageID := uuid.New().String()
	message := &dto.MessageResponse{
		ID:        messageID,
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
		TokensUsed: len(req.Message) / 4, // Aproximação simples
		Metadata:  req.Metadata,
	}

	// Atualizar contexto da sessão
	session.UpdateInteraction()
	session.SetContext("last_message_tokens", message.TokensUsed)

	// TODO: Aqui seria integrado com Claude API e execução de ferramentas

	s.logger.Debug("Mensagem enviada para sessão",
		zap.String("session_id", sessionID),
		zap.String("message_id", messageID),
		zap.Int("tokens_used", message.TokensUsed),
	)

	return message, nil
}

// GetConversationHistory obtém histórico da conversa
func (s *SessionService) GetConversationHistory(ctx context.Context, sessionID string, page, pageSize int) (*dto.ConversationHistoryResponse, error) {
	_, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	// TODO: Implementar busca real no banco de dados
	// Por agora, retorna estrutura vazia
	return &dto.ConversationHistoryResponse{
		SessionID:  sessionID,
		Messages:   []dto.MessageResponse{},
		TotalCount: 0,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    false,
	}, nil
}

// mapSessionToResponse converte sessão do domínio para DTO
func (s *SessionService) mapSessionToResponse(session *domain.MCPSession) *dto.SessionResponse {
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
			ClaudeModel:   "claude-3-sonnet-20241022", // Valor padrão
			MaxTokens:     4096, // Valor padrão
			Timeout:       30, // Valor padrão
			AutoSave:      true, // Valor padrão
			Notifications: true, // Valor padrão
		},
		Context: dto.ConversationContext{
			MessagesCount:  session.MessageCount,
			TokensUsed:     0, // TODO: implementar contagem de tokens
			LastTopic:      "", // TODO: extrair do contexto se necessário
			ActiveTools:    []string{}, // TODO: extrair do contexto se necessário
			CurrentProcess: "", // TODO: extrair do contexto se necessário
			ContextData:    session.Context,
		},
	}
}

// Helpers
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getOrDefaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}