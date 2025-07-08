package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/auth-service/internal/application"
	"github.com/direito-lux/auth-service/internal/infrastructure/config"
	"github.com/direito-lux/auth-service/internal/infrastructure/metrics"
	"github.com/direito-lux/auth-service/internal/infrastructure/http/middleware"

)

// Server estrutura do servidor HTTP
type Server struct {
	config      *config.Config
	logger      *zap.Logger
	metrics     *metrics.Metrics
	server      *http.Server
	router      *gin.Engine
	authHandler *AuthHandler
	userHandler *UserHandler
}

// NewServer cria nova instância do servidor HTTP
func NewServer(
	cfg *config.Config,
	logger *zap.Logger,
	metrics *metrics.Metrics,
	authService *application.AuthService,
	userService *application.UserService,
) *Server {
	// Configurar modo do Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Criar handlers
	authHandler := NewAuthHandler(authService, logger)
	userHandler := NewUserHandler(userService, logger)

	server := &Server{
		config:      cfg,
		logger:      logger,
		metrics:     metrics,
		router:      router,
		authHandler: authHandler,
		userHandler: userHandler,
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
}

// setupRoutes configura rotas do servidor
func (s *Server) setupRoutes() {
	// Health check routes
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/ready", s.readinessCheck)

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			// Endpoints existentes
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/refresh", s.authHandler.RefreshToken)
			auth.POST("/logout", s.authHandler.Logout)
			auth.GET("/validate", s.authHandler.ValidateToken)
			
			// Novos endpoints
			auth.POST("/register", s.authHandler.Register)          // Registro público
			auth.POST("/forgot-password", s.authHandler.ForgotPassword) // Recuperação de senha
			auth.POST("/reset-password", s.authHandler.ResetPassword)   // Reset de senha
		}
		
		// User management routes
		users := api.Group("/users")
		{
			users.POST("", s.userHandler.CreateUser)
			users.GET("", s.userHandler.GetUsersByTenant)
			users.GET("/:id", s.userHandler.GetUser)
			users.PUT("/:id", s.userHandler.UpdateUser)
			users.DELETE("/:id", s.userHandler.DeleteUser)
			users.POST("/change-password", s.userHandler.ChangePassword)
		}
	}

	// Swagger documentation
	if !s.config.IsProduction() {
		s.setupSwagger()
	}
}

// healthCheck handler para health check
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   s.config.ServiceName,
		"version":   s.config.Version,
		"timestamp": gin.H{},
	})
}

// readinessCheck handler para readiness check
func (s *Server) readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": s.config.ServiceName,
	})
}

// setupSwagger configura documentação Swagger
func (s *Server) setupSwagger() {
	// Implementar setup do Swagger aqui
	s.logger.Info("Swagger documentação disponível em /swagger/index.html")
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	s.logger.Info("Iniciando servidor HTTP",
		zap.String("addr", s.server.Addr),
		zap.String("environment", s.config.Environment),
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("erro ao iniciar servidor: %w", err)
	}

	return nil
}

// Shutdown para o servidor HTTP gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Parando servidor HTTP...")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("erro ao parar servidor: %w", err)
	}

	s.logger.Info("Servidor HTTP parado com sucesso")
	return nil
}