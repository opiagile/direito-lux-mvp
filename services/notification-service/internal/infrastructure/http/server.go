package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/application/services"
	"github.com/direito-lux/notification-service/internal/domain"
	"github.com/direito-lux/notification-service/internal/infrastructure/config"
	"github.com/direito-lux/notification-service/internal/infrastructure/http/middleware"
	"github.com/direito-lux/notification-service/internal/infrastructure/metrics"
)

// Server representa o servidor HTTP
type Server struct {
	httpServer          *http.Server
	gin                 *gin.Engine
	logger              *zap.Logger
	config              *config.Config
	metrics             *metrics.Metrics
	notificationService *services.NotificationService
	templateService     *services.TemplateService
	preferenceRepo      domain.NotificationPreferenceRepository
}

// NewServer cria novo servidor HTTP
func NewServer(
	cfg *config.Config,
	logger *zap.Logger,
	metrics *metrics.Metrics,
	notificationService *services.NotificationService,
	templateService *services.TemplateService,
	preferenceRepo domain.NotificationPreferenceRepository,
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
	r.Use(middleware.CORS())

	// Criar servidor
	server := &Server{
		gin:                 r,
		logger:              logger,
		config:              cfg,
		metrics:             metrics,
		notificationService: notificationService,
		templateService:     templateService,
		preferenceRepo:      preferenceRepo,
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
	// Setup comprehensive routes using routes.go
	routerConfig := &RouterConfig{
		NotificationService: s.notificationService,
		TemplateService:     s.templateService,
		PreferenceRepo:      s.preferenceRepo,
		Logger:              s.logger,
	}
	
	SetupRoutes(s.gin, routerConfig)
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