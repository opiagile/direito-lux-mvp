package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/direito-lux/process-service/internal/infrastructure/config"
	"github.com/direito-lux/process-service/internal/infrastructure/metrics"
	"github.com/direito-lux/process-service/internal/infrastructure/http/middleware"
	"github.com/direito-lux/process-service/internal/infrastructure/http/handlers"
)

// Server estrutura do servidor HTTP
type Server struct {
	config     *config.Config
	logger     *zap.Logger
	metrics    *metrics.Metrics
	server     *http.Server
	router     *gin.Engine
	db         *sqlx.DB
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

	// Conectar ao banco de dados
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		"localhost", 5432, "direito_lux", "direito_lux_pass_dev", "direito_lux_dev", "disable")
	
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}

	server := &Server{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		router:  router,
		db:      db,
	}

	// Configurar middlewares
	server.setupMiddlewares()

	// Configurar rotas
	server.setupRoutes()

	// Configurar servidor HTTP
	server.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
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
	s.router.Use(middleware.RateLimit(s.config))

	// Metrics middleware
	if s.metrics != nil {
		s.router.Use(s.metrics.HTTPMiddleware())
	}
}

// setupRoutes configura rotas do servidor
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", handlers.HealthCheck(s.config))
	s.router.GET("/ready", handlers.ReadinessCheck(s.config))

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Example routes
		api.GET("/ping", handlers.Ping())
		
		// Process endpoints (REAL - substituindo templates)
		processes := api.Group("/processes")
		{
			// Inicializar handlers de processos (será injetado via DI depois)
			// Por enquanto, endpoints básicos sem DB
			processes.GET("", s.listProcesses())
			processes.POST("", s.createProcess())
			processes.GET("/:id", s.getProcess())
			processes.PUT("/:id", s.updateProcess())
			processes.DELETE("/:id", s.deleteProcess())
			processes.GET("/stats", s.getProcessStats()) // CRÍTICO para dashboard
		}
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

	// Fechar conexão com banco de dados
	if s.db != nil {
		s.db.Close()
	}

	s.logger.Info("Servidor HTTP parado com sucesso")
	return nil
}

// Métodos para endpoints de processos (implementação temporária simples)

// getProcessStats retorna estatísticas de processos
func (s *Server) getProcessStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		// Conectar com banco de dados para dados reais
		db := s.db
		if db == nil {
			c.JSON(500, gin.H{"error": "Database connection not available"})
			return
		}

		var stats struct {
			Total    int `db:"total"`
			Active   int `db:"active"`
			Archived int `db:"archived"`
		}

		// Query para buscar estatísticas reais
		query := `
			SELECT 
				COUNT(*) as total,
				COUNT(*) FILTER (WHERE status = 'active') as active,
				COUNT(*) FILTER (WHERE status = 'archived') as archived
			FROM processes 
			WHERE tenant_id = $1
		`
		
		err := db.Get(&stats, query, tenantID)
		if err != nil {
			s.logger.Error("Erro ao buscar estatísticas", zap.Error(err))
			c.JSON(500, gin.H{"error": "Database query failed"})
			return
		}

		statsResponse := gin.H{
			"total":         stats.Total,
			"active":        stats.Active,
			"paused":        0, // TODO: implementar status paused
			"archived":      stats.Archived,
			"this_month":    stats.Total, // TODO: implementar filtro por mês
			"todayMovements":     0, // TODO: implementar
			"upcomingDeadlines":  0, // TODO: implementar
		}

		s.logger.Info("Retornando estatísticas de processos",
			zap.String("tenant_id", tenantID),
			zap.Any("stats", stats),
		)

		c.JSON(200, gin.H{
			"data": statsResponse,
			"timestamp": "2025-01-05T10:30:00Z",
		})
	}
}

// listProcesses lista processos
func (s *Server) listProcesses() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		// Implementação temporária - TODO: conectar DB
		processes := []gin.H{
			{
				"id":     "proc-1",
				"number": "5001234-12.2024.8.26.0100",
				"title":  "Ação de Cobrança - Cliente XYZ",
				"court":  "TJSP",
				"status": "active",
			},
		}

		c.JSON(200, gin.H{
			"data":  processes,
			"total": len(processes),
		})
	}
}

// createProcess cria novo processo
func (s *Server) createProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		// Implementação temporária - TODO: conectar DB
		c.JSON(201, gin.H{
			"data": gin.H{
				"id":      "proc-new",
				"message": "Processo criado com sucesso",
			},
		})
	}
}

// getProcess busca processo por ID
func (s *Server) getProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		processID := c.Param("id")
		
		// Implementação temporária - TODO: conectar DB
		c.JSON(200, gin.H{
			"data": gin.H{
				"id":     processID,
				"number": "5001234-12.2024.8.26.0100",
				"title":  "Processo encontrado",
			},
		})
	}
}

// updateProcess atualiza processo
func (s *Server) updateProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		processID := c.Param("id")

		// Implementação temporária - TODO: conectar DB
		c.JSON(200, gin.H{
			"data": gin.H{
				"id":      processID,
				"message": "Processo atualizado com sucesso",
			},
		})
	}
}

// deleteProcess exclui processo
func (s *Server) deleteProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		processID := c.Param("id")

		// Implementação temporária - TODO: conectar DB
		c.JSON(200, gin.H{
			"data": gin.H{
				"id":      processID,
				"message": "Processo excluído com sucesso",
			},
		})
	}
}