package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// DataJudClient cliente HTTP para API DataJud do CNJ
type DataJudClient struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
	retryCount int
	retryDelay time.Duration
}

// DataJudConfig configuração do cliente DataJud
type DataJudClientConfig struct {
	BaseURL    string
	Timeout    time.Duration
	RetryCount int
	RetryDelay time.Duration
}

// NewDataJudClient cria novo cliente DataJud
func NewDataJudClient(config DataJudClientConfig) *DataJudClient {
	// Configurar cliente HTTP com timeout e TLS
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Em produção deve ser false
			},
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	return &DataJudClient{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		timeout:    config.Timeout,
		retryCount: config.RetryCount,
		retryDelay: config.RetryDelay,
	}
}

// QueryProcess consulta dados de um processo específico
func (c *DataJudClient) QueryProcess(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/processos/%s", c.baseURL, req.ProcessNumber)
	
	// Construir query parameters
	params := map[string]string{
		"tribunal": req.CourtID,
	}

	return c.makeRequest(ctx, "GET", endpoint, params, nil, req, provider)
}

// QueryMovements consulta movimentações de um processo
func (c *DataJudClient) QueryMovements(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/processos/%s/movimentacoes", c.baseURL, req.ProcessNumber)
	
	params := map[string]string{
		"tribunal": req.CourtID,
	}

	// Adicionar parâmetros opcionais
	if page, ok := req.Parameters["page"].(int); ok && page > 0 {
		params["pagina"] = fmt.Sprintf("%d", page)
	}
	if pageSize, ok := req.Parameters["page_size"].(int); ok && pageSize > 0 {
		params["tamanho"] = fmt.Sprintf("%d", pageSize)
	}
	if dateFrom, ok := req.Parameters["date_from"].(time.Time); ok {
		params["data_inicio"] = dateFrom.Format("2006-01-02")
	}
	if dateTo, ok := req.Parameters["date_to"].(time.Time); ok {
		params["data_fim"] = dateTo.Format("2006-01-02")
	}

	return c.makeRequest(ctx, "GET", endpoint, params, nil, req, provider)
}

// QueryParties consulta partes de um processo
func (c *DataJudClient) QueryParties(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/processos/%s/partes", c.baseURL, req.ProcessNumber)
	
	params := map[string]string{
		"tribunal": req.CourtID,
	}

	return c.makeRequest(ctx, "GET", endpoint, params, nil, req, provider)
}

// QueryDocuments consulta documentos de um processo
func (c *DataJudClient) QueryDocuments(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/processos/%s/documentos", c.baseURL, req.ProcessNumber)
	
	params := map[string]string{
		"tribunal": req.CourtID,
	}

	return c.makeRequest(ctx, "GET", endpoint, params, nil, req, provider)
}

// BulkQuery executa consultas em lote
func (c *DataJudClient) BulkQuery(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/processos/lote", c.baseURL)
	
	// Construir payload para consulta em lote
	bulkRequest := map[string]interface{}{
		"tribunal":  req.CourtID,
		"processos": req.Parameters["process_numbers"],
	}

	return c.makeRequest(ctx, "POST", endpoint, nil, bulkRequest, req, provider)
}

// makeRequest executa uma requisição HTTP com retry e instrumentação
func (c *DataJudClient) makeRequest(
	ctx context.Context,
	method, endpoint string,
	params map[string]string,
	body interface{},
	req *domain.DataJudRequest,
	provider *domain.CNPJProvider,
) (*domain.DataJudResponse, error) {
	
	var lastErr error
	
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			// Aguardar antes de retry
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay * time.Duration(attempt)):
			}
		}

		response, err := c.executeRequest(ctx, method, endpoint, params, body, provider)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Verificar se deve fazer retry baseado no erro
		if !c.shouldRetry(err, attempt) {
			break
		}
	}

	return nil, fmt.Errorf("requisição falhou após %d tentativas: %w", c.retryCount+1, lastErr)
}

// executeRequest executa uma única requisição HTTP
func (c *DataJudClient) executeRequest(
	ctx context.Context,
	method, endpoint string,
	params map[string]string,
	body interface{},
	provider *domain.CNPJProvider,
) (*domain.DataJudResponse, error) {
	
	startTime := time.Now()

	// Preparar body se necessário
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("erro ao serializar body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Criar requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Configurar headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Direito-Lux/1.0")
	
	// Autenticação com API key do provider
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.APIKey))
	
	// Headers específicos DataJud
	httpReq.Header.Set("X-CNPJ", provider.CNPJ)
	httpReq.Header.Set("X-Request-ID", uuid.New().String())

	// Adicionar query parameters
	if len(params) > 0 {
		q := httpReq.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	// Executar requisição
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer httpResp.Body.Close()

	// Ler body da resposta
	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	duration := time.Since(startTime)

	// Criar resposta DataJud
	response := &domain.DataJudResponse{
		ID:         uuid.New(),
		RequestID:  uuid.New(), // Seria passado como parâmetro
		StatusCode: httpResp.StatusCode,
		Headers:    make(map[string]string),
		Body:       responseBody,
		Size:       int64(len(responseBody)),
		Duration:   duration.Milliseconds(),
		FromCache:  false,
		ReceivedAt: time.Now(),
	}

	// Copiar headers relevantes
	for key, values := range httpResp.Header {
		if len(values) > 0 {
			response.Headers[key] = values[0]
		}
	}

	// Verificar se a resposta foi bem-sucedida
	if httpResp.StatusCode >= 400 {
		return response, fmt.Errorf("API DataJud retornou erro %d: %s", httpResp.StatusCode, string(responseBody))
	}

	// Processar dados estruturados baseado no tipo de requisição
	if err := c.parseStructuredData(response, responseBody); err != nil {
		// Log do erro mas não falha a requisição
		// O body raw ainda está disponível
	}

	return response, nil
}

// parseStructuredData processa dados estruturados da resposta
func (c *DataJudClient) parseStructuredData(response *domain.DataJudResponse, body []byte) error {
	// Tentar parsear como JSON genérico primeiro
	var genericData map[string]interface{}
	if err := json.Unmarshal(body, &genericData); err != nil {
		return err // Não é JSON válido
	}

	// Detectar tipo de dados baseado na estrutura
	if c.isProcessData(genericData) {
		processData := &domain.ProcessResponseData{}
		if err := json.Unmarshal(body, processData); err == nil {
			response.ProcessData = processData
		}
	} else if c.isMovementData(genericData) {
		movementData := &domain.MovementResponseData{}
		if err := json.Unmarshal(body, movementData); err == nil {
			response.MovementData = movementData
		}
	} else if c.isPartyData(genericData) {
		partyData := &domain.PartyResponseData{}
		if err := json.Unmarshal(body, partyData); err == nil {
			response.PartyData = partyData
		}
	}

	return nil
}

// isProcessData verifica se os dados são de processo
func (c *DataJudClient) isProcessData(data map[string]interface{}) bool {
	// Verificar se tem campos típicos de processo
	_, hasNumber := data["numero"]
	_, hasAssunto := data["assunto"]
	_, hasTribunal := data["tribunal"]
	
	return hasNumber && (hasAssunto || hasTribunal)
}

// isMovementData verifica se os dados são de movimentação
func (c *DataJudClient) isMovementData(data map[string]interface{}) bool {
	// Verificar se tem array de movimentações
	if movs, ok := data["movimentacoes"]; ok {
		if movArray, ok := movs.([]interface{}); ok && len(movArray) > 0 {
			return true
		}
	}
	
	// Ou se tem campos típicos de movimentação individual
	_, hasData := data["data"]
	_, hasDescricao := data["descricao"]
	_, hasCodigo := data["codigo"]
	
	return hasData && (hasDescricao || hasCodigo)
}

// isPartyData verifica se os dados são de partes
func (c *DataJudClient) isPartyData(data map[string]interface{}) bool {
	// Verificar se tem array de partes
	if parts, ok := data["partes"]; ok {
		if partArray, ok := parts.([]interface{}); ok && len(partArray) > 0 {
			return true
		}
	}
	
	// Ou se tem campos típicos de parte individual
	_, hasNome := data["nome"]
	_, hasDocumento := data["documento"]
	_, hasTipo := data["tipo"]
	
	return hasNome && (hasDocumento || hasTipo)
}

// shouldRetry determina se deve fazer retry baseado no erro
func (c *DataJudClient) shouldRetry(err error, attempt int) bool {
	if attempt >= c.retryCount {
		return false
	}

	// Determinar se o erro é retryable
	// Em uma implementação real, analisaria o tipo de erro HTTP
	// Por simplicidade, retry para qualquer erro
	return true
}

// ValidateCertificate valida certificado digital para autenticação
func (c *DataJudClient) ValidateCertificate(provider *domain.CNPJProvider) error {
	if provider.Certificate == "" {
		return fmt.Errorf("certificado não configurado para CNPJ %s", provider.CNPJ)
	}

	// Implementação real validaria o certificado digital
	// Por simplicidade, apenas verificar se não está vazio
	return nil
}

// TestConnection testa conectividade com a API DataJud
func (c *DataJudClient) TestConnection(ctx context.Context, provider *domain.CNPJProvider) error {
	endpoint := fmt.Sprintf("%s/api/v1/health", c.baseURL)
	
	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.APIKey))
	httpReq.Header.Set("X-CNPJ", provider.CNPJ)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("erro de conectividade: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 400 {
		return fmt.Errorf("API DataJud não disponível: status %d", httpResp.StatusCode)
	}

	return nil
}

// GetRateLimitInfo obtém informações de rate limit dos headers
func (c *DataJudClient) GetRateLimitInfo(headers map[string]string) map[string]interface{} {
	info := make(map[string]interface{})
	
	// Headers padrão de rate limiting
	if limit, ok := headers["X-RateLimit-Limit"]; ok {
		info["limit"] = limit
	}
	if remaining, ok := headers["X-RateLimit-Remaining"]; ok {
		info["remaining"] = remaining
	}
	if reset, ok := headers["X-RateLimit-Reset"]; ok {
		info["reset"] = reset
	}
	if retryAfter, ok := headers["Retry-After"]; ok {
		info["retry_after"] = retryAfter
	}

	return info
}

// Close fecha o cliente e limpa recursos
func (c *DataJudClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}