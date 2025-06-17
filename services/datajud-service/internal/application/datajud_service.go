package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// DataJudService serviço principal para consultas DataJud
type DataJudService struct {
	repos           *domain.Repositories
	poolManager     *CNPJPoolManager
	rateLimitManager *RateLimitManager
	circuitManager  *CircuitBreakerManager
	cacheManager    *CacheManager
	domainService   domain.DomainService
	config          domain.DataJudConfig
}

// NewDataJudService cria nova instância do serviço
func NewDataJudService(
	repos *domain.Repositories,
	poolManager *CNPJPoolManager,
	rateLimitManager *RateLimitManager,
	circuitManager *CircuitBreakerManager,
	cacheManager *CacheManager,
	domainService domain.DomainService,
	config domain.DataJudConfig,
) *DataJudService {
	return &DataJudService{
		repos:            repos,
		poolManager:      poolManager,
		rateLimitManager: rateLimitManager,
		circuitManager:   circuitManager,
		cacheManager:     cacheManager,
		domainService:    domainService,
		config:           config,
	}
}

// QueryProcess consulta um processo específico
func (s *DataJudService) QueryProcess(ctx context.Context, req *ProcessQueryRequest) (*ProcessQueryResponse, error) {
	// Validar entrada
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validação falhou: %w", err)
	}

	// Criar requisição DataJud
	datajudReq := domain.NewDataJudRequest(
		req.TenantID,
		req.ClientID,
		domain.RequestTypeProcess,
		s.domainService.CalculateRequestPriority(domain.RequestTypeProcess, req.Urgent),
	)
	datajudReq.SetProcessNumber(req.ProcessNumber)
	datajudReq.SetCourtID(req.CourtID)
	if req.ProcessID != nil {
		datajudReq.ProcessID = req.ProcessID
	}

	// Verificar cache primeiro
	if req.UseCache {
		if cachedResp, err := s.checkCache(ctx, datajudReq); err == nil && cachedResp != nil {
			return &ProcessQueryResponse{
				RequestID: datajudReq.ID,
				Status:    "completed",
				Data:      cachedResp.ProcessData,
				FromCache: true,
				CachedAt:  &cachedResp.ReceivedAt,
			}, nil
		}
	}

	// Executar consulta
	response, err := s.executeRequest(ctx, datajudReq)
	if err != nil {
		return &ProcessQueryResponse{
			RequestID: datajudReq.ID,
			Status:    "failed",
			Error:     err.Error(),
		}, err
	}

	return &ProcessQueryResponse{
		RequestID: datajudReq.ID,
		Status:    "completed",
		Data:      response.ProcessData,
		FromCache: false,
		Duration:  response.Duration,
	}, nil
}

// QueryMovements consulta movimentações de um processo
func (s *DataJudService) QueryMovements(ctx context.Context, req *MovementQueryRequest) (*MovementQueryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validação falhou: %w", err)
	}

	datajudReq := domain.NewDataJudRequest(
		req.TenantID,
		req.ClientID,
		domain.RequestTypeMovement,
		s.domainService.CalculateRequestPriority(domain.RequestTypeMovement, req.Urgent),
	)
	datajudReq.SetProcessNumber(req.ProcessNumber)
	datajudReq.SetCourtID(req.CourtID)
	datajudReq.SetParameter("page", req.Page)
	datajudReq.SetParameter("page_size", req.PageSize)
	datajudReq.SetParameter("date_from", req.DateFrom)
	datajudReq.SetParameter("date_to", req.DateTo)

	// Verificar cache
	if req.UseCache {
		if cachedResp, err := s.checkCache(ctx, datajudReq); err == nil && cachedResp != nil {
			return &MovementQueryResponse{
				RequestID: datajudReq.ID,
				Status:    "completed",
				Data:      cachedResp.MovementData,
				FromCache: true,
				CachedAt:  &cachedResp.ReceivedAt,
			}, nil
		}
	}

	response, err := s.executeRequest(ctx, datajudReq)
	if err != nil {
		return &MovementQueryResponse{
			RequestID: datajudReq.ID,
			Status:    "failed",
			Error:     err.Error(),
		}, err
	}

	return &MovementQueryResponse{
		RequestID: datajudReq.ID,
		Status:    "completed",
		Data:      response.MovementData,
		FromCache: false,
		Duration:  response.Duration,
	}, nil
}

// BulkQuery executa consultas em lote
func (s *DataJudService) BulkQuery(ctx context.Context, req *BulkQueryRequest) (*BulkQueryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validação falhou: %w", err)
	}

	response := &BulkQueryResponse{
		RequestID: uuid.New(),
		Status:    "processing",
		Results:   make([]BulkQueryResult, 0, len(req.Queries)),
		StartedAt: time.Now(),
	}

	// Processar cada consulta
	for i, query := range req.Queries {
		result := BulkQueryResult{
			Index:         i,
			ProcessNumber: query.ProcessNumber,
			Status:        "processing",
		}

		// Criar requisição individual
		datajudReq := domain.NewDataJudRequest(
			req.TenantID,
			req.ClientID,
			domain.RequestTypeBulk,
			domain.PriorityLow, // Bulk sempre baixa prioridade
		)
		datajudReq.SetProcessNumber(query.ProcessNumber)
		datajudReq.SetCourtID(query.CourtID)

		// Executar consulta
		datajudResp, err := s.executeRequest(ctx, datajudReq)
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
		} else {
			result.Status = "completed"
			result.Data = datajudResp.ProcessData
			result.Duration = datajudResp.Duration
		}

		response.Results = append(response.Results, result)

		// Throttling entre consultas para não sobrecarregar
		if i < len(req.Queries)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	response.Status = "completed"
	response.CompletedAt = &[]time.Time{time.Now()}[0]
	response.Duration = time.Since(response.StartedAt)

	return response, nil
}

// executeRequest executa uma requisição DataJud com todas as proteções
func (s *DataJudService) executeRequest(ctx context.Context, req *domain.DataJudRequest) (*domain.DataJudResponse, error) {
	// 1. Verificar circuit breaker
	circuitKey := fmt.Sprintf("datajud:%s", req.CourtID)
	cb := s.circuitManager.GetOrCreate(circuitKey)
	
	if !cb.CanExecute() {
		event := domain.NewCircuitBreakerOpened(cb.ID, cb.Name, cb.FailureCount, cb.Config.FailureThreshold)
		s.repos.EventStore.SaveEvent(ctx, event)
		return nil, domain.ErrCircuitBreakerOpen
	}

	// 2. Selecionar CNPJ provider
	pool, err := s.poolManager.GetTenantPool(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter pool de CNPJs: %w", err)
	}

	provider, err := pool.GetNextProvider()
	if err != nil {
		return nil, fmt.Errorf("nenhum CNPJ disponível: %w", err)
	}

	// 3. Verificar rate limiting do CNPJ
	rateLimiter := s.rateLimitManager.GetCNPJLimiter(provider.CNPJ)
	status := rateLimiter.Allow()
	if !status.Allowed {
		// Tentar outro CNPJ
		provider, err = s.findAlternativeCNPJ(ctx, pool, req)
		if err != nil {
			event := &domain.RateLimitExceeded{
				BaseEvent: domain.BaseEvent{
					ID:          uuid.New(),
					Type:        "datajud.rate_limit.exceeded",
					AggregateID: provider.ID,
					OccurredAt:  time.Now(),
					Version:     1,
					Metadata:    make(map[string]interface{}),
				},
				LimitType:     domain.RateLimitCNPJ,
				Key:           provider.CNPJ,
				RequestsUsed:  status.RequestsUsed,
				RequestsLimit: status.RequestsLimit,
				ResetTime:     status.ResetTime,
			}
			s.repos.EventStore.SaveEvent(ctx, event)
			return nil, domain.ErrRateLimitExceeded
		}
	}

	// 4. Atualizar requisição com provider selecionado
	req.SetCNPJProvider(provider.ID)
	req.SetCircuitBreakerKey(circuitKey)

	// 5. Salvar requisição
	if err := s.repos.DataJudRequest.Save(req); err != nil {
		return nil, fmt.Errorf("erro ao salvar requisição: %w", err)
	}

	// 6. Executar com circuit breaker
	var response *domain.DataJudResponse
	result := cb.Execute(func() error {
		req.StartProcessing()
		s.repos.DataJudRequest.Update(req)

		// Publicar evento de início
		event := &domain.DataJudRequestStarted{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.request.started",
				AggregateID: req.ID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata:    make(map[string]interface{}),
			},
			CNPJProviderID: provider.ID,
			CNPJ:           provider.CNPJ,
			ProcessingAt:   *req.ProcessingAt,
		}
		s.repos.EventStore.SaveEvent(ctx, event)

		// Executar consulta real (seria implementado na infraestrutura)
		var err error
		response, err = s.executeHTTPRequest(ctx, req, provider)
		return err
	})

	// 7. Processar resultado
	if result.Success {
		// Usar quota do provider
		provider.UseQuota(1)
		s.repos.CNPJProvider.Update(provider)

		// Completar requisição
		req.Complete(response)
		s.repos.DataJudRequest.Update(req)

		// Salvar no cache se configurado
		if req.UseCache && response != nil {
			s.cacheManager.Set(req.CacheKey, response, req.CacheTTL, req.TenantID, req.Type)
		}

		// Publicar evento de sucesso
		event := domain.NewDataJudRequestCompleted(
			req.ID,
			provider.ID,
			response.StatusCode,
			response.Size,
			result.Duration,
			response.FromCache,
		)
		s.repos.EventStore.SaveEvent(ctx, event)

		return response, nil
	} else {
		// Processar falha
		req.Fail("REQUEST_FAILED", result.Error.Error())
		s.repos.DataJudRequest.Update(req)

		// Verificar se deve fazer retry
		if req.CanRetry() {
			retryDelay := s.calculateRetryDelay(req.RetryCount)
			req.Retry(retryDelay)
			s.repos.DataJudRequest.Update(req)

			// Publicar evento de retry
			event := &domain.DataJudRequestRetrying{
				BaseEvent: domain.BaseEvent{
					ID:          uuid.New(),
					Type:        "datajud.request.retrying",
					AggregateID: req.ID,
					OccurredAt:  time.Now(),
					Version:     1,
					Metadata:    make(map[string]interface{}),
				},
				RetryCount: req.RetryCount,
				MaxRetries: req.MaxRetries,
				RetryAfter: retryDelay,
			}
			s.repos.EventStore.SaveEvent(ctx, event)
		}

		// Publicar evento de falha
		event := &domain.DataJudRequestFailed{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.request.failed",
				AggregateID: req.ID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata:    make(map[string]interface{}),
			},
			CNPJProviderID: &provider.ID,
			ErrorCode:      "REQUEST_FAILED",
			ErrorMessage:   result.Error.Error(),
			RetryCount:     req.RetryCount,
			WillRetry:      req.CanRetry(),
		}
		s.repos.EventStore.SaveEvent(ctx, event)

		return nil, result.Error
	}
}

// checkCache verifica se existe entrada no cache
func (s *DataJudService) checkCache(ctx context.Context, req *domain.DataJudRequest) (*domain.DataJudResponse, error) {
	entry, err := s.cacheManager.Get(req.CacheKey)
	if err != nil || entry == nil {
		// Publicar evento de cache miss
		event := &domain.DataJudCacheMiss{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.cache.miss",
				AggregateID: req.ID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata:    make(map[string]interface{}),
			},
			CacheKey:    req.CacheKey,
			RequestType: req.Type,
		}
		s.repos.EventStore.SaveEvent(ctx, event)
		return nil, err
	}

	// Verificar se o cache ainda é válido
	if !s.domainService.ShouldUseCache(req.Type, entry.GetAge()) {
		s.cacheManager.Delete(req.CacheKey)
		return nil, nil
	}

	// Publicar evento de cache hit
	event := &domain.DataJudCacheHit{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.cache.hit",
			AggregateID: req.ID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		CacheKey:    req.CacheKey,
		RequestType: req.Type,
		HitCount:    entry.HitCount,
		Age:         entry.GetAge(),
	}
	s.repos.EventStore.SaveEvent(ctx, event)

	// Converter entrada do cache para response
	response, ok := entry.Value.(*domain.DataJudResponse)
	if !ok {
		return nil, fmt.Errorf("formato inválido no cache")
	}

	response.FromCache = true
	return response, nil
}

// findAlternativeCNPJ busca um CNPJ alternativo com quota disponível
func (s *DataJudService) findAlternativeCNPJ(ctx context.Context, pool *domain.CNPJPool, req *domain.DataJudRequest) (*domain.CNPJProvider, error) {
	providers := pool.GetAllProviders()
	
	for _, provider := range providers {
		if !provider.IsActive {
			continue
		}

		rateLimiter := s.rateLimitManager.GetCNPJLimiter(provider.CNPJ)
		status := rateLimiter.GetStatus()
		
		if status.Allowed && provider.CanMakeRequest() {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("nenhum CNPJ alternativo disponível")
}

// calculateRetryDelay calcula delay exponencial para retry
func (s *DataJudService) calculateRetryDelay(retryCount int) time.Duration {
	baseDelay := s.config.APIRetryDelay
	exponentialDelay := time.Duration(1<<uint(retryCount)) * baseDelay
	
	// Máximo de 5 minutos
	maxDelay := 5 * time.Minute
	if exponentialDelay > maxDelay {
		exponentialDelay = maxDelay
	}
	
	return exponentialDelay
}

// executeHTTPRequest executa a requisição HTTP real (seria implementado na infraestrutura)
func (s *DataJudService) executeHTTPRequest(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	// Esta implementação seria feita na camada de infraestrutura
	// Aqui é apenas um placeholder
	return &domain.DataJudResponse{
		ID:         uuid.New(),
		RequestID:  req.ID,
		StatusCode: 200,
		Headers:    make(map[string]string),
		Body:       []byte(`{"status": "success"}`),
		Size:       100,
		Duration:   2000, // 2 segundos
		FromCache:  false,
		ReceivedAt: time.Now(),
	}, nil
}