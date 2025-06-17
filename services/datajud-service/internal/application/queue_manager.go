package application

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// QueueManager gerencia filas de requisições DataJud
type QueueManager struct {
	repos            *domain.Repositories
	poolManager      *CNPJPoolManager
	rateLimitManager *RateLimitManager
	circuitManager   *CircuitBreakerManager
	cacheManager     *CacheManager
	
	// Filas por prioridade
	urgentQueue   *RequestQueue
	highQueue     *RequestQueue
	normalQueue   *RequestQueue
	lowQueue      *RequestQueue
	
	// Workers
	workers       []*QueueWorker
	workerCount   int
	
	// Controle
	isRunning     bool
	stopChan      chan bool
	wg            sync.WaitGroup
	mu            sync.RWMutex
	
	// Configuração
	config        domain.DataJudConfig
	
	// Estatísticas
	stats         *QueueStatistics
}

// RequestQueue fila de requisições com prioridade
type RequestQueue struct {
	requests  []*domain.DataJudRequest
	mu        sync.RWMutex
	priority  domain.RequestPriority
	maxSize   int
}

// QueueWorker processa requisições da fila
type QueueWorker struct {
	id       int
	manager  *QueueManager
	stopChan chan bool
	isActive bool
	mu       sync.RWMutex
}

// QueueStatistics estatísticas das filas
type QueueStatistics struct {
	TotalQueued     int64 `json:"total_queued"`
	TotalProcessed  int64 `json:"total_processed"`
	TotalFailed     int64 `json:"total_failed"`
	TotalCompleted  int64 `json:"total_completed"`
	TotalRetries    int64 `json:"total_retries"`
	AvgProcessTime  float64 `json:"avg_process_time_ms"`
	StartTime       time.Time `json:"start_time"`
	mu              sync.RWMutex
}

// NewQueueManager cria novo gerenciador de filas
func NewQueueManager(
	repos *domain.Repositories,
	poolManager *CNPJPoolManager,
	rateLimitManager *RateLimitManager,
	circuitManager *CircuitBreakerManager,
	cacheManager *CacheManager,
	config domain.DataJudConfig,
) *QueueManager {
	manager := &QueueManager{
		repos:            repos,
		poolManager:      poolManager,
		rateLimitManager: rateLimitManager,
		circuitManager:   circuitManager,
		cacheManager:     cacheManager,
		config:           config,
		stopChan:         make(chan bool),
		workerCount:      5, // Configurável
		stats: &QueueStatistics{
			StartTime: time.Now(),
		},
	}

	// Inicializar filas
	manager.urgentQueue = NewRequestQueue(domain.PriorityUrgent, 1000)
	manager.highQueue = NewRequestQueue(domain.PriorityHigh, 5000)
	manager.normalQueue = NewRequestQueue(domain.PriorityNormal, 10000)
	manager.lowQueue = NewRequestQueue(domain.PriorityLow, 20000)

	return manager
}

// NewRequestQueue cria nova fila de requisições
func NewRequestQueue(priority domain.RequestPriority, maxSize int) *RequestQueue {
	return &RequestQueue{
		requests: make([]*domain.DataJudRequest, 0),
		priority: priority,
		maxSize:  maxSize,
	}
}

// Start inicia o processamento das filas
func (qm *QueueManager) Start(ctx context.Context) error {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if qm.isRunning {
		return fmt.Errorf("queue manager já está rodando")
	}

	// Carregar requisições pendentes do banco
	if err := qm.loadPendingRequests(ctx); err != nil {
		return fmt.Errorf("erro ao carregar requisições pendentes: %w", err)
	}

	// Criar workers
	qm.workers = make([]*QueueWorker, qm.workerCount)
	for i := 0; i < qm.workerCount; i++ {
		qm.workers[i] = NewQueueWorker(i, qm)
		qm.wg.Add(1)
		go qm.workers[i].Run(ctx)
	}

	qm.isRunning = true
	return nil
}

// Stop para o processamento das filas
func (qm *QueueManager) Stop() {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if !qm.isRunning {
		return
	}

	// Parar workers
	close(qm.stopChan)
	
	// Parar workers individuais
	for _, worker := range qm.workers {
		worker.Stop()
	}

	// Aguardar workers terminarem
	qm.wg.Wait()

	qm.isRunning = false
}

// Enqueue adiciona requisição à fila apropriada
func (qm *QueueManager) Enqueue(ctx context.Context, request *domain.DataJudRequest) error {
	// Salvar requisição no banco
	if err := qm.repos.DataJudRequest.Save(request); err != nil {
		return fmt.Errorf("erro ao salvar requisição: %w", err)
	}

	// Adicionar à fila baseada na prioridade
	queue := qm.getQueueByPriority(request.Priority)
	if err := queue.Enqueue(request); err != nil {
		return fmt.Errorf("erro ao adicionar à fila: %w", err)
	}

	// Atualizar estatísticas
	qm.stats.mu.Lock()
	qm.stats.TotalQueued++
	qm.stats.mu.Unlock()

	return nil
}

// Dequeue obtém próxima requisição para processamento
func (qm *QueueManager) Dequeue() *domain.DataJudRequest {
	// Processar por ordem de prioridade: urgent -> high -> normal -> low
	queues := []*RequestQueue{
		qm.urgentQueue,
		qm.highQueue,
		qm.normalQueue,
		qm.lowQueue,
	}

	for _, queue := range queues {
		if request := queue.Dequeue(); request != nil {
			return request
		}
	}

	return nil
}

// GetQueueSize retorna tamanho total das filas
func (qm *QueueManager) GetQueueSize() int {
	return qm.urgentQueue.Size() + 
		   qm.highQueue.Size() + 
		   qm.normalQueue.Size() + 
		   qm.lowQueue.Size()
}

// GetQueueSizeByPriority retorna tamanho de fila específica
func (qm *QueueManager) GetQueueSizeByPriority(priority domain.RequestPriority) int {
	return qm.getQueueByPriority(priority).Size()
}

// GetQueueStats retorna estatísticas das filas
func (qm *QueueManager) GetQueueStats() map[string]interface{} {
	qm.stats.mu.RLock()
	defer qm.stats.mu.RUnlock()

	return map[string]interface{}{
		"total_queued":       qm.stats.TotalQueued,
		"total_processed":    qm.stats.TotalProcessed,
		"total_completed":    qm.stats.TotalCompleted,
		"total_failed":       qm.stats.TotalFailed,
		"total_retries":      qm.stats.TotalRetries,
		"avg_process_time":   qm.stats.AvgProcessTime,
		"queue_sizes": map[string]int{
			"urgent":  qm.urgentQueue.Size(),
			"high":    qm.highQueue.Size(),
			"normal":  qm.normalQueue.Size(),
			"low":     qm.lowQueue.Size(),
		},
		"total_queue_size":   qm.GetQueueSize(),
		"active_workers":     qm.getActiveWorkerCount(),
		"uptime":            time.Since(qm.stats.StartTime).String(),
		"is_running":        qm.isRunning,
	}
}

// ProcessRequest processa uma requisição específica
func (qm *QueueManager) ProcessRequest(ctx context.Context, request *domain.DataJudRequest) error {
	startTime := time.Now()

	// Marcar como processando
	request.StartProcessing()
	qm.repos.DataJudRequest.Update(request)

	// Verificar se pode usar cache
	if request.UseCache {
		cachedResponse, err := qm.cacheManager.Get(request.CacheKey)
		if err == nil && cachedResponse != nil {
			// Cache hit - completar requisição
			response := cachedResponse.Value.(*domain.DataJudResponse)
			response.FromCache = true
			request.Complete(response)
			qm.repos.DataJudRequest.Update(request)
			
			qm.updateProcessingStats(true, time.Since(startTime))
			return nil
		}
	}

	// Obter pool de CNPJs
	pool, err := qm.poolManager.GetTenantPool(request.TenantID)
	if err != nil {
		qm.failRequest(request, "POOL_ERROR", err.Error())
		qm.updateProcessingStats(false, time.Since(startTime))
		return err
	}

	// Selecionar CNPJ provider
	provider, err := pool.GetNextProvider()
	if err != nil {
		qm.failRequest(request, "NO_PROVIDER", err.Error())
		qm.updateProcessingStats(false, time.Since(startTime))
		return err
	}

	// Verificar rate limiting
	status, err := qm.rateLimitManager.CheckAllowance(provider.CNPJ, request.TenantID)
	if err != nil || !status.Allowed {
		// Tentar agendar para mais tarde
		if request.CanRetry() {
			retryDelay := status.RetryAfter
			if retryDelay < time.Minute {
				retryDelay = time.Minute
			}
			request.Retry(retryDelay)
			qm.repos.DataJudRequest.Update(request)
			qm.stats.mu.Lock()
			qm.stats.TotalRetries++
			qm.stats.mu.Unlock()
			return qm.Enqueue(ctx, request)
		}
		
		qm.failRequest(request, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded")
		qm.updateProcessingStats(false, time.Since(startTime))
		return fmt.Errorf("rate limit exceeded")
	}

	// Executar com circuit breaker
	circuitKey := fmt.Sprintf("datajud:%s", request.CourtID)
	result := qm.circuitManager.ExecuteWithBreaker(ctx, circuitKey, func() error {
		// Simular execução da requisição HTTP
		return qm.executeDataJudRequest(ctx, request, provider)
	})

	if result.Success {
		// Sucesso - usar quota do provider
		provider.UseQuota(1)
		qm.repos.CNPJProvider.Update(provider)
		qm.updateProcessingStats(true, time.Since(startTime))
	} else {
		// Falha - verificar se deve retry
		if request.CanRetry() {
			retryDelay := qm.calculateRetryDelay(request.RetryCount)
			request.Retry(retryDelay)
			qm.repos.DataJudRequest.Update(request)
			qm.stats.mu.Lock()
			qm.stats.TotalRetries++
			qm.stats.mu.Unlock()
			return qm.Enqueue(ctx, request)
		}
		
		qm.failRequest(request, "EXECUTION_FAILED", result.Error.Error())
		qm.updateProcessingStats(false, time.Since(startTime))
	}

	return nil
}

// RetryFailedRequests reprocessa requisições que falharam
func (qm *QueueManager) RetryFailedRequests(ctx context.Context) error {
	// Buscar requisições que podem ser reprocessadas
	requests, err := qm.repos.DataJudRequest.FindByStatus(domain.StatusRetrying, 100, 0)
	if err != nil {
		return err
	}

	for _, request := range requests {
		// Verificar se já é hora de tentar novamente
		if request.RetryAfter != nil && time.Now().Before(*request.RetryAfter) {
			continue
		}

		// Recolocar na fila
		request.Status = domain.StatusPending
		qm.repos.DataJudRequest.Update(request)
		qm.Enqueue(ctx, request)
	}

	return nil
}

// CleanupOldRequests limpa requisições antigas
func (qm *QueueManager) CleanupOldRequests(ctx context.Context, maxAge time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-maxAge)
	return qm.repos.DataJudRequest.CleanupOldRequests(cutoffTime)
}

// getQueueByPriority retorna fila baseada na prioridade
func (qm *QueueManager) getQueueByPriority(priority domain.RequestPriority) *RequestQueue {
	switch priority {
	case domain.PriorityUrgent:
		return qm.urgentQueue
	case domain.PriorityHigh:
		return qm.highQueue
	case domain.PriorityNormal:
		return qm.normalQueue
	case domain.PriorityLow:
		return qm.lowQueue
	default:
		return qm.normalQueue
	}
}

// loadPendingRequests carrega requisições pendentes do banco
func (qm *QueueManager) loadPendingRequests(ctx context.Context) error {
	requests, err := qm.repos.DataJudRequest.FindPendingRequests(1000)
	if err != nil {
		return err
	}

	for _, request := range requests {
		queue := qm.getQueueByPriority(request.Priority)
		queue.Enqueue(request)
	}

	return nil
}

// failRequest marca requisição como falhada
func (qm *QueueManager) failRequest(request *domain.DataJudRequest, errorCode, errorMessage string) {
	request.Fail(errorCode, errorMessage)
	qm.repos.DataJudRequest.Update(request)
}

// calculateRetryDelay calcula delay para retry
func (qm *QueueManager) calculateRetryDelay(retryCount int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s, 8s, 16s, max 5min
	delay := time.Duration(1<<uint(retryCount)) * time.Second
	maxDelay := 5 * time.Minute
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

// executeDataJudRequest simula execução da requisição DataJud
func (qm *QueueManager) executeDataJudRequest(ctx context.Context, request *domain.DataJudRequest, provider *domain.CNPJProvider) error {
	// Esta seria a implementação real da chamada HTTP para DataJud
	// Por simplicidade, simulando sucesso/falha
	
	// Simular tempo de processamento
	time.Sleep(100 * time.Millisecond)
	
	// Criar resposta simulada
	response := &domain.DataJudResponse{
		ID:         uuid.New(),
		RequestID:  request.ID,
		StatusCode: 200,
		Headers:    make(map[string]string),
		Body:       []byte(`{"status": "success"}`),
		Size:       100,
		Duration:   100,
		FromCache:  false,
		ReceivedAt: time.Now(),
	}

	// Completar requisição
	request.Complete(response)
	qm.repos.DataJudRequest.Update(request)

	// Armazenar no cache se configurado
	if request.UseCache {
		qm.cacheManager.Set(
			request.CacheKey, 
			response, 
			request.CacheTTL, 
			request.TenantID, 
			request.Type,
		)
	}

	return nil
}

// updateProcessingStats atualiza estatísticas de processamento
func (qm *QueueManager) updateProcessingStats(success bool, duration time.Duration) {
	qm.stats.mu.Lock()
	defer qm.stats.mu.Unlock()

	qm.stats.TotalProcessed++
	
	if success {
		qm.stats.TotalCompleted++
	} else {
		qm.stats.TotalFailed++
	}

	// Atualizar tempo médio de processamento
	durationMs := float64(duration.Nanoseconds()) / 1000000
	if qm.stats.AvgProcessTime == 0 {
		qm.stats.AvgProcessTime = durationMs
	} else {
		qm.stats.AvgProcessTime = (qm.stats.AvgProcessTime + durationMs) / 2
	}
}

// getActiveWorkerCount conta workers ativos
func (qm *QueueManager) getActiveWorkerCount() int {
	count := 0
	for _, worker := range qm.workers {
		if worker.IsActive() {
			count++
		}
	}
	return count
}

// ========================================
// MÉTODOS DA REQUESTQUEUE
// ========================================

// Enqueue adiciona requisição à fila
func (rq *RequestQueue) Enqueue(request *domain.DataJudRequest) error {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if len(rq.requests) >= rq.maxSize {
		return fmt.Errorf("fila cheia (max: %d)", rq.maxSize)
	}

	rq.requests = append(rq.requests, request)
	
	// Ordenar por prioridade e timestamp
	sort.Slice(rq.requests, func(i, j int) bool {
		if rq.requests[i].Priority == rq.requests[j].Priority {
			return rq.requests[i].RequestedAt.Before(rq.requests[j].RequestedAt)
		}
		return rq.requests[i].Priority > rq.requests[j].Priority
	})

	return nil
}

// Dequeue remove e retorna próxima requisição
func (rq *RequestQueue) Dequeue() *domain.DataJudRequest {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if len(rq.requests) == 0 {
		return nil
	}

	request := rq.requests[0]
	rq.requests = rq.requests[1:]
	return request
}

// Size retorna tamanho da fila
func (rq *RequestQueue) Size() int {
	rq.mu.RLock()
	defer rq.mu.RUnlock()
	return len(rq.requests)
}

// ========================================
// MÉTODOS DO QUEUEWORKER
// ========================================

// NewQueueWorker cria novo worker
func NewQueueWorker(id int, manager *QueueManager) *QueueWorker {
	return &QueueWorker{
		id:       id,
		manager:  manager,
		stopChan: make(chan bool),
		isActive: false,
	}
}

// Run executa o worker
func (qw *QueueWorker) Run(ctx context.Context) {
	defer qw.manager.wg.Done()
	
	qw.mu.Lock()
	qw.isActive = true
	qw.mu.Unlock()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-qw.stopChan:
			qw.mu.Lock()
			qw.isActive = false
			qw.mu.Unlock()
			return
		case <-qw.manager.stopChan:
			qw.mu.Lock()
			qw.isActive = false
			qw.mu.Unlock()
			return
		case <-ticker.C:
			// Processar próxima requisição
			request := qw.manager.Dequeue()
			if request != nil {
				qw.manager.ProcessRequest(ctx, request)
			}
		}
	}
}

// Stop para o worker
func (qw *QueueWorker) Stop() {
	close(qw.stopChan)
}

// IsActive retorna se worker está ativo
func (qw *QueueWorker) IsActive() bool {
	qw.mu.RLock()
	defer qw.mu.RUnlock()
	return qw.isActive
}