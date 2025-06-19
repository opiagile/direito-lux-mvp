package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/direito-lux/mcp-service/internal/infrastructure/http/dto"
)

// CreateMCPSessionV2 cria uma nova sessão MCP
func CreateMCPSessionV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Obter tenant_id do contexto se não fornecido
		if req.TenantID == "" {
			if tenantID, exists := c.Get("tenant_id"); exists {
				req.TenantID = tenantID.(string)
			}
		}

		// Obter user_id do contexto se não fornecido
		if req.UserID == "" {
			if userID, exists := c.Get("user_id"); exists {
				req.UserID = userID.(string)
			}
		}

		// TODO: Injetar SessionService via dependency injection
		// Por agora, simular resposta
		sessionID := uuid.New().String()
		response := &dto.SessionResponse{
			ID:           sessionID,
			Channel:      req.Channel,
			UserID:       req.UserID,
			TenantID:     req.TenantID,
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LastActivity: time.Now(),
			Metadata:     req.Metadata,
			Settings:     req.Settings,
			Context: dto.ConversationContext{
				MessagesCount: 0,
				TokensUsed:    0,
				ContextData:   make(map[string]interface{}),
			},
		}

		c.JSON(http.StatusCreated, response)
	}
}

// GetMCPSessionV2 obtém informações de uma sessão
func GetMCPSessionV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		// TODO: Buscar sessão do banco de dados
		// Por agora, simular resposta
		response := &dto.SessionResponse{
			ID:           sessionID,
			Channel:      "web",
			UserID:       "user123",
			TenantID:     "tenant123",
			Status:       "active",
			CreatedAt:    time.Now().Add(-1 * time.Hour),
			UpdatedAt:    time.Now(),
			LastActivity: time.Now().Add(-5 * time.Minute),
			Metadata:     make(map[string]interface{}),
			Settings: dto.SessionSettings{
				ClaudeModel: "claude-3-sonnet-20241022",
				MaxTokens:   4096,
				Timeout:     30,
			},
			Context: dto.ConversationContext{
				MessagesCount: 5,
				TokensUsed:    1250,
				LastTopic:     "Consulta sobre processo trabalhista",
				ActiveTools:   []string{"process_search", "document_analyze"},
				ContextData:   make(map[string]interface{}),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// CloseMCPSessionV2 fecha uma sessão MCP
func CloseMCPSessionV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		// TODO: Fechar sessão no banco de dados
		// Por agora, simular fechamento

		c.JSON(http.StatusOK, gin.H{
			"message": "Session closed successfully",
			"session_id": sessionID,
			"closed_at": time.Now(),
		})
	}
}

// GetSessionStatusV2 obtém status de uma sessão
func GetSessionStatusV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		// TODO: Buscar status real da sessão
		// Por agora, simular status
		response := &dto.SessionStatusResponse{
			ID:           sessionID,
			Status:       "active",
			IsActive:     true,
			LastActivity: time.Now().Add(-2 * time.Minute),
			Context: dto.ConversationContext{
				MessagesCount: 8,
				TokensUsed:    2100,
				LastTopic:     "Análise de contrato",
				ActiveTools:   []string{"contract_review"},
				ContextData:   make(map[string]interface{}),
			},
			QuotaUsage: dto.QuotaUsage{
				TokensUsed:    2100,
				TokensLimit:   10000,
				RequestsUsed:  8,
				RequestsLimit: 100,
				UsagePercent:  21.0,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// SendMessageV2 envia mensagem para uma sessão
func SendMessageV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		var req dto.SendMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// TODO: Processar mensagem com Claude e ferramentas
		// Por agora, simular resposta
		messageID := uuid.New().String()
		response := &dto.MessageResponse{
			ID:        messageID,
			SessionID: sessionID,
			Role:      "assistant",
			Content:   "Entendi sua solicitação. Como posso ajudar você hoje?",
			Timestamp: time.Now(),
			TokensUsed: 25,
			Metadata:  make(map[string]interface{}),
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetConversationHistoryV2 obtém histórico de conversa
func GetConversationHistoryV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		// Parâmetros de paginação
		page := 1
		pageSize := 20

		if p := c.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}

		if ps := c.Query("page_size"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil {
				pageSize = parsed
			}
		}

		// TODO: Buscar histórico real do banco
		// Por agora, simular histórico
		messages := []dto.MessageResponse{
			{
				ID:        "msg1",
				SessionID: sessionID,
				Role:      "user",
				Content:   "Olá, preciso de ajuda com um processo trabalhista",
				Timestamp: time.Now().Add(-30 * time.Minute),
				TokensUsed: 20,
			},
			{
				ID:        "msg2",
				SessionID: sessionID,
				Role:      "assistant",
				Content:   "Claro! Posso ajudar você com questões trabalhistas. Qual é a sua dúvida específica?",
				Timestamp: time.Now().Add(-29 * time.Minute),
				TokensUsed: 30,
			},
		}

		response := &dto.ConversationHistoryResponse{
			SessionID:  sessionID,
			Messages:   messages,
			TotalCount: len(messages),
			Page:       page,
			PageSize:   pageSize,
			HasMore:    false,
		}

		c.JSON(http.StatusOK, response)
	}
}