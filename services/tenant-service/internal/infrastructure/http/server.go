package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/tenant-service/internal/infrastructure/config"
	"github.com/direito-lux/tenant-service/internal/infrastructure/metrics"
	"github.com/direito-lux/tenant-service/internal/infrastructure/http/middleware"
	"github.com/direito-lux/tenant-service/internal/infrastructure/database"
	"github.com/direito-lux/tenant-service/internal/infrastructure/repository"
	"github.com/direito-lux/tenant-service/internal/domain"
)

// Server estrutura do servidor HTTP
type Server struct {
	config         *config.Config
	logger         *zap.Logger
	metrics        *metrics.Metrics
	server         *http.Server
	router         *gin.Engine
	tenantRepo     domain.TenantRepository
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
		config:     cfg,
		logger:     logger,
		metrics:    metrics,
		router:     router,
		tenantRepo: nil, // Will be initialized in Start()
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
	s.logger.Info("🚀🚀🚀 SETUPROUTES EXECUTING - NEW VERSION 🚀🚀🚀")
	
	// Health check
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/ready", s.readinessCheck)
	
	// SIMPLE TEST ROUTE
	s.router.GET("/simple-test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test works"})
	})

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Health check
		api.GET("/ping", s.ping)
		
		// Tenant endpoints - USING CORRECTED HANDLERS
		tenants := api.Group("/tenants")
		{
			tenants.GET("/:id", s.getTenantByID)          // REAL DATABASE QUERY
		}
		
		// Tenant-specific endpoints (direct routes) - FORCE COMPILATION
		api.GET("/tenant/current", s.getCurrentTenant)     // Get current tenant from X-Tenant-ID
		api.GET("/tenant/subscription", s.getSubscription) // Get subscription details
		api.GET("/tenant/quotas", s.getQuotas)            // Get usage quotas
	}

	// Swagger documentation
	if !s.config.IsProduction() {
		s.setupSwagger()
	}
}

// healthCheck endpoint
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "tenant-service", // Fixed name!
		"timestamp": "2025-07-02T00:00:00Z",
		"version":   s.config.Version,
	})
}

// readinessCheck endpoint
func (s *Server) readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"service":   "tenant-service", // Fixed name!
		"timestamp": "2025-07-02T00:00:00Z",
		"version":   s.config.Version,
	})
}

// ping endpoint
func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"service":   "tenant-service", // Fixed name!
		"timestamp": "2025-07-02T00:00:00Z",
	})
}

// getTenantByID - REAL POSTGRESQL DATABASE QUERY
func (s *Server) getTenantByID(c *gin.Context) {
	tenantID := c.Param("id")
	
	s.logger.Info("🔍 Fetching tenant from PostgreSQL database", zap.String("tenant_id", tenantID))
	
	// ✅ REAL DATABASE QUERY - NO MORE HARDCODED SWITCH CASES!
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		if err == domain.ErrTenantNotFound {
			s.logger.Warn("❌ Tenant not found in database", zap.String("tenant_id", tenantID))
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tenant not found",
				"message": "O escritório especificado não foi encontrado no sistema",
			})
			return
		}
		
		s.logger.Error("❌ Database error fetching tenant", 
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"message": "Erro ao buscar dados do escritório",
		})
		return
	}
	
	// Convert domain.Tenant to API response format
	response := map[string]interface{}{
		"id":        tenant.ID,
		"name":      tenant.Name,
		"legalName": tenant.LegalName,
		"document":  tenant.Document,
		"email":     tenant.Email,
		"phone":     tenant.Phone,
		"website":   tenant.Website,
		"address":   tenant.Address,
		"status":    tenant.Status,
		"plan":      tenant.PlanType,
		"isActive":  tenant.Status == domain.TenantStatusActive,
		"createdAt": tenant.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt": tenant.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		"settings":  tenant.Settings,
	}
	
	if tenant.ActivatedAt != nil {
		response["activatedAt"] = tenant.ActivatedAt.Format("2006-01-02T15:04:05Z")
	}
	
	if tenant.SuspendedAt != nil {
		response["suspendedAt"] = tenant.SuspendedAt.Format("2006-01-02T15:04:05Z")
	}
	
	s.logger.Info("✅ Tenant retrieved successfully from PostgreSQL", 
		zap.String("tenant_id", tenantID),
		zap.String("tenant_name", tenant.Name),
		zap.String("legal_name", tenant.LegalName),
		zap.String("plan", string(tenant.PlanType)),
	)
	
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// getCurrentTenant - Get current tenant from X-Tenant-ID header
func (s *Server) getCurrentTenant(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		s.logger.Warn("❌ X-Tenant-ID header missing")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "X-Tenant-ID header é obrigatório",
			"message": "O identificador do escritório deve ser fornecido no header",
		})
		return
	}
	
	s.logger.Info("🔍 Fetching current tenant", zap.String("tenant_id", tenantID))
	
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		if err == domain.ErrTenantNotFound {
			s.logger.Warn("❌ Current tenant not found", zap.String("tenant_id", tenantID))
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tenant not found",
				"message": "Escritório atual não encontrado",
			})
			return
		}
		
		s.logger.Error("❌ Database error fetching current tenant", 
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"message": "Erro ao buscar dados do escritório atual",
		})
		return
	}
	
	// Same response format as getTenantByID
	response := map[string]interface{}{
		"id":        tenant.ID,
		"name":      tenant.Name,
		"legalName": tenant.LegalName,
		"document":  tenant.Document,
		"email":     tenant.Email,
		"phone":     tenant.Phone,
		"website":   tenant.Website,
		"address":   tenant.Address,
		"status":    tenant.Status,
		"plan":      tenant.PlanType,
		"isActive":  tenant.Status == domain.TenantStatusActive,
		"createdAt": tenant.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt": tenant.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		"settings":  tenant.Settings,
	}
	
	if tenant.ActivatedAt != nil {
		response["activatedAt"] = tenant.ActivatedAt.Format("2006-01-02T15:04:05Z")
	}
	
	if tenant.SuspendedAt != nil {
		response["suspendedAt"] = tenant.SuspendedAt.Format("2006-01-02T15:04:05Z")
	}
	
	s.logger.Info("✅ Current tenant retrieved successfully", 
		zap.String("tenant_id", tenantID),
		zap.String("tenant_name", tenant.Name),
	)
	
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// getSubscription - Get subscription details for current tenant
func (s *Server) getSubscription(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "X-Tenant-ID header é obrigatório",
		})
		return
	}
	
	s.logger.Info("🔍 Fetching subscription data", zap.String("tenant_id", tenantID))
	
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		if err == domain.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tenant not found",
			})
			return
		}
		
		s.logger.Error("❌ Database error fetching tenant for subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	
	// Build subscription response based on tenant plan
	subscription := gin.H{
		"plan":         tenant.PlanType,
		"status":       "active", // TODO: implementar status real da assinatura
		"startDate":    tenant.ActivatedAt.Format("2006-01-02T15:04:05Z"),
		"renewalDate":  "2025-12-31T23:59:59Z", // TODO: calcular data real
		"autoRenewal":  true,
		"billingCycle": "monthly",
		"amount":       s.getPlanAmount(tenant.PlanType),
		"currency":     "BRL",
		"features":     s.getPlanFeatures(tenant.PlanType),
	}
	
	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

// getQuotas - Get usage quotas for current tenant
func (s *Server) getQuotas(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "X-Tenant-ID header é obrigatório",
		})
		return
	}
	
	s.logger.Info("🔍 Fetching quotas data", zap.String("tenant_id", tenantID))
	
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		if err == domain.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tenant not found",
			})
			return
		}
		
		s.logger.Error("❌ Database error fetching tenant for quotas", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	
	// Build quotas response based on tenant plan
	quotas := s.getPlanQuotas(tenant.PlanType)
	quotas["tenant_id"] = tenantID
	quotas["plan"] = tenant.PlanType
	quotas["period"] = "monthly"
	quotas["resetDate"] = "2025-08-01T00:00:00Z" // TODO: calcular data real
	
	c.JSON(http.StatusOK, gin.H{"data": quotas})
}

// Helper methods for plan data
func (s *Server) getPlanAmount(plan domain.PlanType) float64 {
	switch plan {
	case "starter":
		return 99.00
	case "professional":
		return 299.00
	case "business":
		return 699.00
	case "enterprise":
		return 1999.00
	default:
		return 0.00
	}
}

func (s *Server) getPlanFeatures(plan domain.PlanType) []string {
	switch plan {
	case "starter":
		return []string{"50 processos", "20 clientes", "100 consultas/dia", "WhatsApp", "Email"}
	case "professional":
		return []string{"200 processos", "100 clientes", "500 consultas/dia", "WhatsApp", "Email", "Telegram"}
	case "business":
		return []string{"500 processos", "500 clientes", "2000 consultas/dia", "WhatsApp", "Email", "Telegram", "Relatórios avançados"}
	case "enterprise":
		return []string{"Ilimitado", "Ilimitado", "10k consultas/dia", "Todos os canais", "White-label", "API", "Suporte prioritário"}
	default:
		return []string{}
	}
}

func (s *Server) getPlanQuotas(plan domain.PlanType) gin.H {
	switch plan {
	case "starter":
		return gin.H{
			"processes": gin.H{"limit": 50, "used": 2, "remaining": 48},
			"clients": gin.H{"limit": 20, "used": 8, "remaining": 12},
			"datajud_queries": gin.H{"limit": 100, "used": 15, "remaining": 85},
			"storage_gb": gin.H{"limit": 1.0, "used": 0.2, "remaining": 0.8},
		}
	case "professional":
		return gin.H{
			"processes": gin.H{"limit": 200, "used": 2, "remaining": 198},
			"clients": gin.H{"limit": 100, "used": 8, "remaining": 92},
			"datajud_queries": gin.H{"limit": 500, "used": 15, "remaining": 485},
			"storage_gb": gin.H{"limit": 5.0, "used": 0.2, "remaining": 4.8},
		}
	case "business":
		return gin.H{
			"processes": gin.H{"limit": 500, "used": 2, "remaining": 498},
			"clients": gin.H{"limit": 500, "used": 8, "remaining": 492},
			"datajud_queries": gin.H{"limit": 2000, "used": 15, "remaining": 1985},
			"storage_gb": gin.H{"limit": 20.0, "used": 0.2, "remaining": 19.8},
		}
	case "enterprise":
		return gin.H{
			"processes": gin.H{"limit": -1, "used": 2, "remaining": -1}, // -1 = unlimited
			"clients": gin.H{"limit": -1, "used": 8, "remaining": -1},
			"datajud_queries": gin.H{"limit": 10000, "used": 15, "remaining": 9985},
			"storage_gb": gin.H{"limit": 100.0, "used": 0.2, "remaining": 99.8},
		}
	default:
		return gin.H{
			"processes": gin.H{"limit": 0, "used": 0, "remaining": 0},
			"clients": gin.H{"limit": 0, "used": 0, "remaining": 0},
			"datajud_queries": gin.H{"limit": 0, "used": 0, "remaining": 0},
			"storage_gb": gin.H{"limit": 0.0, "used": 0.0, "remaining": 0.0},
		}
	}
}

// setupSwagger configura documentação Swagger
func (s *Server) setupSwagger() {
	s.logger.Info("Swagger documentação disponível em /swagger/index.html")
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	// Initialize database connection
	s.logger.Info("🚀 Initializing database connection...")
	
	db, err := database.NewPostgreSQLConnection(s.config, s.logger)
	if err != nil {
		s.logger.Error("❌ Failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test database connection
	if err := database.TestConnection(db, s.logger); err != nil {
		s.logger.Error("❌ Database connection test failed", zap.Error(err))
		return fmt.Errorf("database connection test failed: %w", err)
	}
	
	// Initialize repository
	s.tenantRepo = repository.NewPostgresTenantRepository(db, s.logger)
	s.logger.Info("✅ Database and repository initialized successfully")

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
