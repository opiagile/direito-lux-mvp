package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/direito-lux/datajud-service/internal/application"
	"github.com/direito-lux/datajud-service/internal/domain"
	"github.com/direito-lux/datajud-service/internal/infrastructure/config"
	httpClient "github.com/direito-lux/datajud-service/internal/infrastructure/http"
	"github.com/direito-lux/datajud-service/internal/infrastructure/handlers"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Erro ao carregar configuração: %v", err)
	}
	
	// Configurar Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// Inicializar dependências
	if err := initializeDependencies(cfg); err != nil {
		log.Fatalf("❌ Erro ao inicializar dependências: %v", err)
	}
	
	// Criar cliente HTTP DataJud
	var httpClientInstance application.HTTPClient
	if cfg.IsDataJudMockEnabled() {
		log.Println("⚠️  Usando cliente MOCK para DataJud (desenvolvimento)")
		httpClientInstance = httpClient.NewMockClient()
	} else {
		log.Println("🔗 Usando cliente HTTP real para DataJud")
		clientConfig := httpClient.DataJudRealClientConfig{
			BaseURL:    cfg.DataJud.BaseURL,
			APIKey:     cfg.DataJud.APIKey,
			Timeout:    cfg.DataJud.Timeout,
			RetryCount: cfg.DataJud.RetryCount,
			RetryDelay: cfg.DataJud.RetryDelay,
		}
		httpClientInstance = httpClient.NewDataJudRealClient(clientConfig)
	}
	defer httpClientInstance.Close()
	
	// Inicializar repositórios e managers (em produção, usar implementações reais)
	repos := &domain.Repositories{} // TODO: Implementar repositories reais
	poolManager := &application.CNPJPoolManager{} // TODO: Implementar pool manager
	rateLimitManager := &application.RateLimitManager{} // TODO: Implementar rate limit manager
	circuitManager := &application.CircuitBreakerManager{} // TODO: Implementar circuit breaker
	cacheManager := &application.CacheManager{} // TODO: Implementar cache manager
	domainService := &mockDomainService{} // Mock domain service
	
	// Criar serviço DataJud
	dataJudService := application.NewDataJudService(
		repos,
		poolManager,
		rateLimitManager,
		circuitManager,
		cacheManager,
		domainService,
		cfg.DataJud,
		httpClientInstance,
	)
	
	// Criar handler HTTP
	handler := handlers.NewDataJudHandler(dataJudService, cfg)
	
	// Configurar router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	
	// Registrar rotas
	handler.RegisterRoutes(router)
	
	// Criar servidor HTTP
	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
		MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes,
	}
	
	// Iniciar servidor em goroutine
	go func() {
		log.Printf("🚀 DataJud Service rodando na porta %d", cfg.Port)
		log.Printf("📊 Ambiente: %s", cfg.Environment)
		log.Printf("🔑 DataJud Mock: %v", cfg.IsDataJudMockEnabled())
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
		}
	}()
	
	// Aguardar sinal de término
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⏹️ Desligando servidor...")
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Erro ao desligar servidor: %v", err)
	}
	
	log.Println("✅ Servidor desligado")
}

// initializeDependencies inicializa dependências do sistema
func initializeDependencies(cfg *config.Config) error {
	log.Println("🔄 Inicializando dependências...")
	
	// Em produção, aqui seria onde conectaríamos ao banco de dados,
	// Redis, RabbitMQ, etc. Por enquanto, apenas log
	log.Println("📦 Dependências inicializadas")
	return nil
}

// mockDomainService implementação mock do domain service
type mockDomainService struct{}

func (m *mockDomainService) CalculateRequestPriority(requestType domain.RequestType, urgent bool) domain.Priority {
	if urgent {
		return domain.PriorityHigh
	}
	return domain.PriorityNormal
}

func (m *mockDomainService) ShouldUseCache(requestType domain.RequestType, age time.Duration) bool {
	switch requestType {
	case domain.RequestTypeProcess:
		return age < 24*time.Hour
	case domain.RequestTypeMovement:
		return age < 30*time.Minute
	default:
		return age < 1*time.Hour
	}
}

