package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/application"
	"github.com/direito-lux/datajud-service/internal/domain"
	"github.com/direito-lux/datajud-service/internal/infrastructure/config"
	"github.com/direito-lux/datajud-service/internal/infrastructure/http"
)

// Exemplo de uso do DataJud Service com cliente real
func main() {
	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Criar cliente HTTP real
	clientConfig := http.DataJudRealClientConfig{
		BaseURL:    cfg.DataJud.BaseURL,
		APIKey:     cfg.DataJud.APIKey,
		Timeout:    cfg.DataJud.Timeout,
		RetryCount: cfg.DataJud.RetryCount,
		RetryDelay: cfg.DataJud.RetryDelay,
	}

	httpClient := http.NewDataJudRealClient(clientConfig)
	defer httpClient.Close()

	// Teste de conexão
	ctx := context.Background()
	if err := httpClient.TestConnection(ctx); err != nil {
		log.Printf("Aviso: Teste de conexão falhou: %v", err)
	}

	// Criar serviço (assumindo que repositories, managers etc. estão configurados)
	repos := &domain.Repositories{} // Configurar repositories
	poolManager := &application.CNPJPoolManager{} // Configurar pool manager
	rateLimitManager := &application.RateLimitManager{} // Configurar rate limit manager
	circuitManager := &application.CircuitBreakerManager{} // Configurar circuit breaker
	cacheManager := &application.CacheManager{} // Configurar cache manager
	domainService := &mockDomainService{} // Mock domain service

	service := application.NewDataJudService(
		repos,
		poolManager,
		rateLimitManager,
		circuitManager,
		cacheManager,
		domainService,
		cfg.DataJud,
		httpClient,
	)

	// Exemplos de uso
	exemploConsultaProcesso(ctx, service)
	exemploConsultaMovimentacoes(ctx, service)
	exemploConsultaLote(ctx, service)
	exemploTribunalMapper()
	exemploQueryBuilder()
}

// exemploConsultaProcesso demonstra consulta de processo individual
func exemploConsultaProcesso(ctx context.Context, service *application.DataJudService) {
	fmt.Println("=== Exemplo: Consulta de Processo ===")
	
	req := &application.ProcessQueryRequest{
		ProcessNumber: "0001234-56.2023.8.26.0001",
		CourtID:       "TJSP",
		TenantID:      "tenant-123",
		ClientID:      "client-456",
		UseCache:      true,
		Urgent:        false,
	}

	response, err := service.QueryProcess(ctx, req)
	if err != nil {
		fmt.Printf("Erro na consulta: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("From Cache: %v\n", response.FromCache)
	if response.Duration != nil {
		fmt.Printf("Duração: %v\n", *response.Duration)
	}
	
	fmt.Println()
}

// exemploConsultaMovimentacoes demonstra consulta de movimentações
func exemploConsultaMovimentacoes(ctx context.Context, service *application.DataJudService) {
	fmt.Println("=== Exemplo: Consulta de Movimentações ===")
	
	dateFrom := time.Now().AddDate(0, 0, -30) // Últimos 30 dias
	dateTo := time.Now()
	
	req := &application.MovementQueryRequest{
		ProcessNumber: "0001234-56.2023.8.26.0001",
		CourtID:       "TJSP",
		TenantID:      "tenant-123",
		ClientID:      "client-456",
		Page:          1,
		PageSize:      50,
		DateFrom:      &dateFrom,
		DateTo:        &dateTo,
		UseCache:      true,
		Urgent:        false,
	}

	response, err := service.QueryMovements(ctx, req)
	if err != nil {
		fmt.Printf("Erro na consulta: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("From Cache: %v\n", response.FromCache)
	if response.Data != nil {
		fmt.Printf("Total de movimentações: %d\n", response.Data.Total)
		fmt.Printf("Movimentações retornadas: %d\n", len(response.Data.Movements))
	}
	
	fmt.Println()
}

// exemploConsultaLote demonstra consulta em lote
func exemploConsultaLote(ctx context.Context, service *application.DataJudService) {
	fmt.Println("=== Exemplo: Consulta em Lote ===")
	
	req := &application.BulkQueryRequest{
		TenantID: "tenant-123",
		ClientID: "client-456",
		Queries: []application.BulkQueryItem{
			{
				ProcessNumber: "0001234-56.2023.8.26.0001",
				CourtID:       "TJSP",
			},
			{
				ProcessNumber: "0001234-56.2023.8.26.0002",
				CourtID:       "TJSP",
			},
			{
				ProcessNumber: "0001234-56.2023.8.26.0003",
				CourtID:       "TJSP",
			},
		},
		UseCache: true,
		Urgent:   false,
	}

	response, err := service.BulkQuery(ctx, req)
	if err != nil {
		fmt.Printf("Erro na consulta: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Total de consultas: %d\n", len(response.Results))
	
	sucessos := 0
	for _, result := range response.Results {
		if result.Status == "completed" {
			sucessos++
		}
	}
	
	fmt.Printf("Sucessos: %d/%d\n", sucessos, len(response.Results))
	fmt.Println()
}

// exemploTribunalMapper demonstra uso do mapeador de tribunais
func exemploTribunalMapper() {
	fmt.Println("=== Exemplo: Tribunal Mapper ===")
	
	mapper := http.NewTribunalMapper()
	
	// Listar alguns tribunais
	tribunais := []string{"STF", "STJ", "TJSP", "TRF3", "TRT2"}
	
	for _, codigo := range tribunais {
		tribunal := mapper.GetTribunal(codigo)
		if tribunal != nil {
			fmt.Printf("%s: %s (Endpoint: %s)\n", 
				tribunal.Code, tribunal.Name, tribunal.Endpoint)
		}
	}
	
	// Estatísticas
	stats := mapper.GetTribunalStats()
	fmt.Printf("\nEstatísticas:\n")
	for tipo, count := range stats {
		fmt.Printf("  %s: %d\n", tipo, count)
	}
	
	fmt.Println()
}

// exemploQueryBuilder demonstra construção de queries
func exemploQueryBuilder() {
	fmt.Println("=== Exemplo: Query Builder ===")
	
	// Query simples por número
	builder := http.NewElasticsearchQueryBuilder()
	query := builder.ProcessByNumber("0001234-56.2023.8.26.0001").Build()
	
	fmt.Printf("Query por número: Size=%d, From=%d\n", query.Size, query.From)
	
	// Query com filtros
	from := time.Now().AddDate(0, 0, -90)
	to := time.Now()
	
	builder = http.NewElasticsearchQueryBuilder()
	query = builder.
		ProcessByNumber("0001234-56.2023.8.26.0001").
		CourtFilter("TJSP").
		DateRange("dataAjuizamento", &from, &to).
		Pagination(2, 25).
		Build()
	
	fmt.Printf("Query com filtros: Size=%d, From=%d\n", query.Size, query.From)
	
	// Templates prontos
	templates := http.NewQueryTemplates()
	
	// Query para processos por parte
	query = templates.GetProcessesByParty("João da Silva", "TJSP")
	fmt.Printf("Query por parte: Size=%d\n", query.Size)
	
	// Query para movimentações recentes
	query = templates.GetRecentMovements("0001234-56.2023.8.26.0001", 30)
	fmt.Printf("Query movimentações recentes: Size=%d\n", query.Size)
	
	fmt.Println()
}

// mockDomainService implementação mock para exemplo
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

// Exemplo de configuração de variáveis de ambiente
func exemploConfiguracaoEnvironment() {
	fmt.Println("=== Exemplo: Variáveis de Ambiente ===")
	fmt.Println("Configurações necessárias no .env:")
	fmt.Println("")
	fmt.Println("# DataJud API Configuration")
	fmt.Println("DATAJUD_BASE_URL=https://api-publica.datajud.cnj.jus.br")
	fmt.Println("DATAJUD_API_KEY=sua-chave-api-aqui")
	fmt.Println("DATAJUD_TIMEOUT=30s")
	fmt.Println("DATAJUD_RETRY_COUNT=3")
	fmt.Println("DATAJUD_RETRY_DELAY=1s")
	fmt.Println("")
	fmt.Println("# Rate Limiting")
	fmt.Println("DATAJUD_RATE_LIMIT_ENABLED=true")
	fmt.Println("DATAJUD_RATE_LIMIT_RPM=100")
	fmt.Println("DATAJUD_RATE_LIMIT_BURST=10")
	fmt.Println("")
	fmt.Println("# Cache")
	fmt.Println("DATAJUD_CACHE_ENABLED=true")
	fmt.Println("DATAJUD_CACHE_PROCESS_TTL=24h")
	fmt.Println("DATAJUD_CACHE_MOVEMENT_TTL=30m")
	fmt.Println("")
	fmt.Println("# Circuit Breaker")
	fmt.Println("DATAJUD_CIRCUIT_BREAKER_ENABLED=true")
	fmt.Println("DATAJUD_CIRCUIT_BREAKER_FAILURE_THRESHOLD=5")
	fmt.Println("")
	fmt.Println("# Development")
	fmt.Println("DATAJUD_MOCK_ENABLED=false")
	fmt.Println()
}