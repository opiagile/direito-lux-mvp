package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/infrastructure/config"
	"github.com/direito-lux/notification-service/internal/infrastructure/http/handlers"
	"github.com/direito-lux/notification-service/internal/infrastructure/http/middleware"
	"github.com/direito-lux/notification-service/internal/infrastructure/metrics"
)

// Server representa o servidor HTTP
type Server struct {
	httpServer *http.Server
	gin        *gin.Engine
	logger     *zap.Logger
	config     *config.Config
	metrics    *metrics.Metrics
}

// NewServer cria novo servidor HTTP
func NewServer(
	cfg *config.Config,
	logger *zap.Logger,
	metrics *metrics.Metrics,
) *Server {
	// Configurar Gin baseado no ambiente
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar instância do Gin
	r := gin.New()

	// Middlewares globais
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg))

	// Criar servidor
	server := &Server{
		gin:     r,
		logger:  logger,
		config:  cfg,
		metrics: metrics,
	}

	// Configurar rotas
	server.setupRoutes()

	// Configurar servidor HTTP
	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

// setupRoutes configura as rotas da API
func (s *Server) setupRoutes() {
	// Health checks
	health := handlers.NewHealthHandler()
	s.gin.GET("/health", health.Health)
	s.gin.GET("/ready", health.Ready)

	// API versioning
	v1 := s.gin.Group("/api/v1")
	{
		// Rotas de notificações serão adicionadas aqui
		// Placeholder para futuras rotas
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Notification Service is running",
				"version": s.config.Version,
			})
		})
	}
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	s.logger.Info("Iniciando servidor HTTP", 
		zap.String("address", s.httpServer.Addr),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("erro ao iniciar servidor HTTP: %w", err)
	}

	return nil
}

// Shutdown para o servidor HTTP gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Parando servidor HTTP...")

	return s.httpServer.Shutdown(ctx)
}

// GetGin retorna instância do Gin (para testes)
func (s *Server) GetGin() *gin.Engine {
	return s.gin
}