package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	
	// Para testes rápidos, usar serviço simplificado
	dataJudService := &SimpleDataJudService{
		httpClient: httpClientInstance,
		config:     cfg.GetDataJudDomainConfig(),
	}
	
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
		Addr:           ":" + strconv.Itoa(cfg.Port),
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


// SimpleDataJudService implementação simplificada para testes rápidos
type SimpleDataJudService struct {
	httpClient application.HTTPClient
	config     domain.DataJudConfig
}

// QueryProcess implementação simplificada de consulta de processo
func (s *SimpleDataJudService) QueryProcess(ctx context.Context, req *application.ProcessQueryRequest) (*application.ProcessQueryResponse, error) {
	// Validar entrada
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validação falhou: %w", err)
	}

	// Criar requisição DataJud básica
	datajudReq := domain.NewDataJudRequest(
		req.TenantID,
		req.ClientID,
		domain.RequestTypeProcess,
		domain.PriorityNormal,
	)
	datajudReq.SetProcessNumber(req.ProcessNumber)
	datajudReq.SetCourtID(req.CourtID)

	// Criar provider mock para testes
	provider := &domain.CNPJProvider{
		ID:   uuid.New(),
		CNPJ: "00000000000000", // CNPJ fictício para testes
	}

	// Executar consulta HTTP diretamente
	response, err := s.httpClient.QueryProcess(ctx, datajudReq, provider)
	if err != nil {
		return &application.ProcessQueryResponse{
			RequestID: datajudReq.ID,
			Status:    "failed",
			Error:     err.Error(),
		}, err
	}

	return &application.ProcessQueryResponse{
		RequestID: datajudReq.ID,
		Status:    "completed",
		Data:      response.ProcessData,
		FromCache: false,
		Duration:  response.Duration,
	}, nil
}

// QueryMovements implementação simplificada de consulta de movimentações
func (s *SimpleDataJudService) QueryMovements(ctx context.Context, req *application.MovementQueryRequest) (*application.MovementQueryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validação falhou: %w", err)
	}

	datajudReq := domain.NewDataJudRequest(
		req.TenantID,
		req.ClientID,
		domain.RequestTypeMovement,
		domain.PriorityNormal,
	)
	datajudReq.SetProcessNumber(req.ProcessNumber)
	datajudReq.SetCourtID(req.CourtID)

	provider := &domain.CNPJProvider{
		ID:   uuid.New(),
		CNPJ: "00000000000000",
	}

	response, err := s.httpClient.QueryMovements(ctx, datajudReq, provider)
	if err != nil {
		return &application.MovementQueryResponse{
			RequestID: datajudReq.ID,
			Status:    "failed",
			Error:     err.Error(),
		}, err
	}

	return &application.MovementQueryResponse{
		RequestID: datajudReq.ID,
		Status:    "completed",
		Data:      response.MovementData,
		FromCache: false,
		Duration:  response.Duration,
	}, nil
}

// BulkQuery implementação simplificada de consulta em lote
func (s *SimpleDataJudService) BulkQuery(ctx context.Context, req *application.BulkQueryRequest) (*application.BulkQueryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validação falhou: %w", err)
	}

	response := &application.BulkQueryResponse{
		RequestID: uuid.New(),
		Status:    "completed",
		Results:   make([]application.BulkQueryResult, 0, len(req.Queries)),
		StartedAt: time.Now(),
	}

	// Simular processamento simples
	for i, query := range req.Queries {
		result := application.BulkQueryResult{
			Index:         i,
			ProcessNumber: query.ProcessNumber,
			Status:        "completed",
		}
		response.Results = append(response.Results, result)
	}

	return response, nil
}

