package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http/middleware"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http/handlers"
)

// Server estrutura do servidor HTTP
type Server struct {
	config     *config.Config
	logger     *zap.Logger
	metrics    *metrics.Metrics
	server     *http.Server
	router     *gin.Engine
}

// NewServer cria nova instância do servidor HTTP
func NewServer(cfg *config.Config, logger *zap.Logger, metrics *metrics.Metrics) *Server {
	// Configurar modo do Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	server := &Server{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		router:  router,
	}

	// Configurar middlewares
	server.setupMiddlewares()

	// Configurar rotas
	server.setupRoutes()

	// Configurar servidor HTTP
	server.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           router,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
	}

	return server
}

// setupMiddlewares configura middlewares do servidor
func (s *Server) setupMiddlewares() {
	// Logger middleware
	s.router.Use(middleware.Logger(s.logger))

	// Recovery middleware
	s.router.Use(middleware.Recovery(s.logger))

	// CORS middleware
	s.router.Use(middleware.CORS(s.config))

	// Request ID middleware
	s.router.Use(middleware.RequestID())

	// Tenant middleware
	s.router.Use(middleware.Tenant(s.logger))

	// Rate limiting middleware
	if s.config.HTTP.RateLimitEnabled {
		s.router.Use(middleware.RateLimit(s.config))
	}

	// Metrics middleware
	if s.metrics != nil {
		s.router.Use(s.metrics.HTTPMiddleware())
	}

	// MCP specific middlewares
	s.router.Use(middleware.MCPAuth(s.config))
	s.router.Use(middleware.MCPQuota(s.config))
}

// setupRoutes configura rotas do servidor
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", handlers.HealthCheck(s.config))
	s.router.GET("/ready", handlers.ReadinessCheck(s.config))

	// API routes
	api := s.router.Group("/api/v1")
	{
		// MCP Session Management
		sessions := api.Group("/sessions")
		{
			sessions.POST("", handlers.CreateMCPSession())
			sessions.GET("/:id", handlers.GetMCPSession())
			sessions.DELETE("/:id", handlers.CloseMCPSession())
			sessions.GET("/:id/status", handlers.GetSessionStatus())
			sessions.POST("/:id/messages", handlers.SendMessage())
			sessions.GET("/:id/history", handlers.GetConversationHistory())
		}

		// MCP Tools
		tools := api.Group("/tools")
		{
			tools.GET("", handlers.ListAvailableTools())
			tools.POST("/execute", handlers.ExecuteTool())
			tools.GET("/executions/:id", handlers.GetToolExecution())
			tools.GET("/executions", handlers.ListToolExecutions())
		}

		// MCP Conversations
		conversations := api.Group("/conversations")
		{
			conversations.GET("", handlers.ListConversations())
			conversations.GET("/:id", handlers.GetConversation())
			conversations.POST("/:id/messages", handlers.AddMessage())
			conversations.DELETE("/:id", handlers.DeleteConversation())
		}

		// Bot Integrations
		bots := api.Group("/bots")
		{
			// WhatsApp Bot
			whatsapp := bots.Group("/whatsapp")
			{
				whatsapp.POST("/webhook", handlers.WhatsAppWebhook())
				whatsapp.GET("/webhook", handlers.WhatsAppWebhookVerification())
				whatsapp.POST("/send", handlers.SendWhatsAppMessage())
			}

			// Telegram Bot
			telegram := bots.Group("/telegram")
			{
				telegram.POST("/webhook", handlers.TelegramWebhook())
				telegram.POST("/send", handlers.SendTelegramMessage())
			}

			// Slack Bot
			slack := bots.Group("/slack")
			{
				slack.POST("/events", handlers.SlackEvents())
				slack.POST("/commands", handlers.SlackCommands())
				slack.POST("/interactive", handlers.SlackInteractive())
			}
		}

		// Claude API Integration
		claude := api.Group("/claude")
		{
			claude.POST("/chat", handlers.ClaudeChat())
			claude.POST("/tools", handlers.ClaudeWithTools())
			claude.GET("/models", handlers.ListClaudeModels())
		}

		// Quota Management
		quota := api.Group("/quota")
		{
			quota.GET("", handlers.GetQuotaUsage())
			quota.GET("/limits", handlers.GetQuotaLimits())
			quota.POST("/reset", handlers.ResetQuota())
		}

		// Analytics and Reports
		analytics := api.Group("/analytics")
		{
			analytics.GET("/usage", handlers.GetUsageAnalytics())
			analytics.GET("/tools", handlers.GetToolAnalytics())
			analytics.GET("/conversations", handlers.GetConversationAnalytics())
			analytics.GET("/performance", handlers.GetPerformanceMetrics())
		}

		// Configuration
		config := api.Group("/config")
		{
			config.GET("/channels", handlers.GetChannelConfig())
			config.PUT("/channels/:channel", handlers.UpdateChannelConfig())
			config.GET("/tools", handlers.GetToolsConfig())
			config.PUT("/tools/:tool", handlers.UpdateToolConfig())
		}
	}

	// WebSocket endpoints for real-time communication
	ws := s.router.Group("/ws")
	{
		ws.GET("/sessions/:id", handlers.WebSocketSession())
		ws.GET("/tools/:session_id", handlers.WebSocketTools())
	}

	// Webhook endpoints (sem autenticação para alguns)
	webhooks := s.router.Group("/webhooks")
	{
		webhooks.POST("/whatsapp", handlers.WhatsAppWebhook())
		webhooks.GET("/whatsapp", handlers.WhatsAppWebhookVerification())
		webhooks.POST("/telegram/:token", handlers.TelegramWebhook())
		webhooks.POST("/slack", handlers.SlackEvents())
	}

	// Swagger documentation
	if !s.config.IsProduction() {
		s.setupSwagger()
	}
}

// setupSwagger configura documentação Swagger
func (s *Server) setupSwagger() {
	// Implementar setup do Swagger aqui
	s.logger.Info("Swagger documentação disponível em /swagger/index.html")
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	s.logger.Info("Iniciando servidor HTTP MCP",
		zap.String("addr", s.server.Addr),
		zap.String("environment", s.config.Environment),
		zap.String("service", "mcp-service"),
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("erro ao iniciar servidor MCP: %w", err)
	}

	return nil
}

// Shutdown para o servidor HTTP gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Parando servidor HTTP MCP...")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("erro ao parar servidor MCP: %w", err)
	}

	s.logger.Info("Servidor HTTP MCP parado com sucesso")
	return nil
}