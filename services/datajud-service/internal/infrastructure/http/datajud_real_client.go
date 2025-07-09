package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// DataJudRealClient cliente HTTP real para API DataJud CNJ
type DataJudRealClient struct {
	httpClient     *http.Client
	baseURL        string
	apiKey         string
	timeout        time.Duration
	retryCount     int
	retryDelay     time.Duration
	tribunalMapper *TribunalMapper
}

// DataJudRealClientConfig configuração do cliente real
type DataJudRealClientConfig struct {
	BaseURL    string
	APIKey     string
	Timeout    time.Duration
	RetryCount int
	RetryDelay time.Duration
}

// ElasticsearchQuery estrutura para queries Elasticsearch
type ElasticsearchQuery struct {
	Query ElasticsearchQueryClause `json:"query"`
	Size  int                      `json:"size"`
	From  int                      `json:"from"`
	Sort  []map[string]interface{} `json:"sort,omitempty"`
}

// ElasticsearchQueryClause cláusulas de query
type ElasticsearchQueryClause struct {
	Bool  *ElasticsearchBoolQuery  `json:"bool,omitempty"`
	Match map[string]interface{}   `json:"match,omitempty"`
	Term  map[string]interface{}   `json:"term,omitempty"`
	Range map[string]interface{}   `json:"range,omitempty"`
}

// ElasticsearchBoolQuery query booleana
type ElasticsearchBoolQuery struct {
	Must   []ElasticsearchQueryClause `json:"must,omitempty"`
	Should []ElasticsearchQueryClause `json:"should,omitempty"`
	Filter []ElasticsearchQueryClause `json:"filter,omitempty"`
}

// ElasticsearchResponse resposta da API
type ElasticsearchResponse struct {
	Took     int                           `json:"took"`
	TimedOut bool                          `json:"timed_out"`
	Shards   ElasticsearchShards           `json:"_shards"`
	Hits     ElasticsearchHits             `json:"hits"`
}

// ElasticsearchShards informações de shards
type ElasticsearchShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// ElasticsearchHits hits da resposta
type ElasticsearchHits struct {
	Total    ElasticsearchTotal `json:"total"`
	MaxScore *float64           `json:"max_score"`
	Hits     []ElasticsearchHit `json:"hits"`
}

// ElasticsearchTotal total de hits
type ElasticsearchTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// ElasticsearchHit hit individual
type ElasticsearchHit struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	ID     string                 `json:"_id"`
	Score  *float64               `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// NewDataJudRealClient cria novo cliente real
func NewDataJudRealClient(config DataJudRealClientConfig) *DataJudRealClient {
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:       20,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: false,
		},
	}

	return &DataJudRealClient{
		httpClient:     httpClient,
		baseURL:        config.BaseURL,
		apiKey:         config.APIKey,
		timeout:        config.Timeout,
		retryCount:     config.RetryCount,
		retryDelay:     config.RetryDelay,
		tribunalMapper: NewTribunalMapper(),
	}
}

// QueryProcess consulta dados de um processo específico
func (c *DataJudRealClient) QueryProcess(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	tribunal := c.tribunalMapper.GetTribunal(req.CourtID)
	if tribunal == nil {
		return nil, fmt.Errorf("tribunal não encontrado: %s", req.CourtID)
	}

	// Construir query para buscar processo específico
	query := &ElasticsearchQuery{
		Query: ElasticsearchQueryClause{
			Bool: &ElasticsearchBoolQuery{
				Must: []ElasticsearchQueryClause{
					{
						Match: map[string]interface{}{
							"numeroProcesso": req.ProcessNumber,
						},
					},
				},
			},
		},
		Size: 1,
		From: 0,
		Sort: []map[string]interface{}{
			{"dataHoraUltimaAtualizacao": map[string]string{"order": "desc"}},
		},
	}

	endpoint := fmt.Sprintf("%s/%s/_search", c.baseURL, tribunal.Endpoint)
	
	response, err := c.executeQuery(ctx, endpoint, query, req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query de processo: %w", err)
	}

	// Processar dados do processo
	if err := c.parseProcessData(response); err != nil {
		return nil, fmt.Errorf("erro ao processar dados do processo: %w", err)
	}

	return response, nil
}

// QueryMovements consulta movimentações de um processo
func (c *DataJudRealClient) QueryMovements(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	tribunal := c.tribunalMapper.GetTribunal(req.CourtID)
	if tribunal == nil {
		return nil, fmt.Errorf("tribunal não encontrado: %s", req.CourtID)
	}

	// Construir query para movimentações
	query := &ElasticsearchQuery{
		Query: ElasticsearchQueryClause{
			Bool: &ElasticsearchBoolQuery{
				Must: []ElasticsearchQueryClause{
					{
						Match: map[string]interface{}{
							"numeroProcesso": req.ProcessNumber,
						},
					},
					{
						Term: map[string]interface{}{
							"tipoDocumento": "movimentacao",
						},
					},
				},
			},
		},
		Size: c.getPageSize(req),
		From: c.getPageOffset(req),
		Sort: []map[string]interface{}{
			{"dataHora": map[string]string{"order": "desc"}},
		},
	}

	// Filtros por data se especificados
	if dateFrom, ok := req.Parameters["date_from"].(time.Time); ok {
		if dateTo, ok := req.Parameters["date_to"].(time.Time); ok {
			query.Query.Bool.Filter = append(query.Query.Bool.Filter, ElasticsearchQueryClause{
				Range: map[string]interface{}{
					"dataHora": map[string]interface{}{
						"gte": dateFrom.Format("2006-01-02"),
						"lte": dateTo.Format("2006-01-02"),
					},
				},
			})
		}
	}

	endpoint := fmt.Sprintf("%s/%s/_search", c.baseURL, tribunal.Endpoint)
	
	response, err := c.executeQuery(ctx, endpoint, query, req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query de movimentações: %w", err)
	}

	// Processar dados de movimentações
	if err := c.parseMovementData(response); err != nil {
		return nil, fmt.Errorf("erro ao processar dados de movimentações: %w", err)
	}

	return response, nil
}

// QueryParties consulta partes de um processo
func (c *DataJudRealClient) QueryParties(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	tribunal := c.tribunalMapper.GetTribunal(req.CourtID)
	if tribunal == nil {
		return nil, fmt.Errorf("tribunal não encontrado: %s", req.CourtID)
	}

	query := &ElasticsearchQuery{
		Query: ElasticsearchQueryClause{
			Bool: &ElasticsearchBoolQuery{
				Must: []ElasticsearchQueryClause{
					{
						Match: map[string]interface{}{
							"numeroProcesso": req.ProcessNumber,
						},
					},
					{
						Term: map[string]interface{}{
							"tipoDocumento": "parte",
						},
					},
				},
			},
		},
		Size: 100, // Máximo de partes esperadas
		From: 0,
	}

	endpoint := fmt.Sprintf("%s/%s/_search", c.baseURL, tribunal.Endpoint)
	
	response, err := c.executeQuery(ctx, endpoint, query, req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query de partes: %w", err)
	}

	// Processar dados de partes
	if err := c.parsePartyData(response); err != nil {
		return nil, fmt.Errorf("erro ao processar dados de partes: %w", err)
	}

	return response, nil
}

// BulkQuery executa consultas em lote
func (c *DataJudRealClient) BulkQuery(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	tribunal := c.tribunalMapper.GetTribunal(req.CourtID)
	if tribunal == nil {
		return nil, fmt.Errorf("tribunal não encontrado: %s", req.CourtID)
	}

	processNumbers, ok := req.Parameters["process_numbers"].([]string)
	if !ok {
		return nil, fmt.Errorf("números de processo não fornecidos para consulta em lote")
	}

	// Construir query para múltiplos processos
	should := make([]ElasticsearchQueryClause, 0, len(processNumbers))
	for _, number := range processNumbers {
		should = append(should, ElasticsearchQueryClause{
			Match: map[string]interface{}{
				"numeroProcesso": number,
			},
		})
	}

	query := &ElasticsearchQuery{
		Query: ElasticsearchQueryClause{
			Bool: &ElasticsearchBoolQuery{
				Should: should,
			},
		},
		Size: len(processNumbers),
		From: 0,
		Sort: []map[string]interface{}{
			{"numeroProcesso": map[string]string{"order": "asc"}},
		},
	}

	endpoint := fmt.Sprintf("%s/%s/_search", c.baseURL, tribunal.Endpoint)
	
	response, err := c.executeQuery(ctx, endpoint, query, req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query em lote: %w", err)
	}

	// Processar dados em lote
	if err := c.parseBulkData(response, processNumbers); err != nil {
		return nil, fmt.Errorf("erro ao processar dados em lote: %w", err)
	}

	return response, nil
}

// executeQuery executa uma query Elasticsearch com retry
func (c *DataJudRealClient) executeQuery(ctx context.Context, endpoint string, query *ElasticsearchQuery, req *domain.DataJudRequest) (*domain.DataJudResponse, error) {
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

		response, err := c.executeHTTPRequest(ctx, endpoint, query, req)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Verificar se deve fazer retry
		if !c.shouldRetry(err, attempt) {
			break
		}
	}

	return nil, fmt.Errorf("query falhou após %d tentativas: %w", c.retryCount+1, lastErr)
}

// executeHTTPRequest executa uma única requisição HTTP
func (c *DataJudRealClient) executeHTTPRequest(ctx context.Context, endpoint string, query *ElasticsearchQuery, req *domain.DataJudRequest) (*domain.DataJudResponse, error) {
	startTime := time.Now()

	// Serializar query
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar query: %w", err)
	}

	// Criar requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Configurar headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Direito-Lux/1.0")
	httpReq.Header.Set("Authorization", fmt.Sprintf("APIKey %s", c.apiKey))
	httpReq.Header.Set("X-Request-ID", uuid.New().String())

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

	// Verificar status da resposta
	if httpResp.StatusCode >= 400 {
		return nil, fmt.Errorf("API DataJud retornou erro %d: %s", httpResp.StatusCode, string(responseBody))
	}

	// Criar resposta DataJud
	response := &domain.DataJudResponse{
		ID:         uuid.New(),
		RequestID:  req.ID,
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

	return response, nil
}

// parseProcessData processa dados específicos de processo
func (c *DataJudRealClient) parseProcessData(response *domain.DataJudResponse) error {
	var esResponse ElasticsearchResponse
	if err := json.Unmarshal(response.Body, &esResponse); err != nil {
		return fmt.Errorf("erro ao parsear resposta Elasticsearch: %w", err)
	}

	if len(esResponse.Hits.Hits) == 0 {
		// Processo não encontrado
		response.ProcessData = &domain.ProcessResponseData{
			Number: "",
		}
		return nil
	}

	// Pegar primeiro hit (processo)
	hit := esResponse.Hits.Hits[0]
	
	processData := &domain.ProcessResponseData{
		Number:    c.extractString(hit.Source, "numeroProcesso"),
		Title:     c.extractString(hit.Source, "classe"),
		Subject:   c.extractMapInterface(hit.Source, "assunto"),
		Court:     c.extractString(hit.Source, "tribunal"),
		Status:    c.extractString(hit.Source, "situacao"),
		Stage:     c.extractString(hit.Source, "grauOrigem"),
		CreatedAt: c.extractTimeWithDefault(hit.Source, "dataAjuizamento"),
		UpdatedAt: c.extractTimeWithDefault(hit.Source, "dataHoraUltimaAtualizacao"),
		Parties:   []domain.PartyData{},
		Movements: []domain.MovementData{},
	}

	response.ProcessData = processData
	return nil
}

// parseMovementData processa dados de movimentações
func (c *DataJudRealClient) parseMovementData(response *domain.DataJudResponse) error {
	var esResponse ElasticsearchResponse
	if err := json.Unmarshal(response.Body, &esResponse); err != nil {
		return fmt.Errorf("erro ao parsear resposta Elasticsearch: %w", err)
	}

	movements := make([]domain.MovementData, 0, len(esResponse.Hits.Hits))
	
	for i, hit := range esResponse.Hits.Hits {
		movement := domain.MovementData{
			Sequence:    i + 1,
			Date:        c.extractTimeWithDefault(hit.Source, "dataHora"),
			Code:        c.extractString(hit.Source, "codigoMovimento"),
			Type:        c.extractString(hit.Source, "tipoMovimento"),
			Title:       c.extractString(hit.Source, "descricaoMovimento"),
			Description: c.extractString(hit.Source, "complementoMovimento"),
			Content:     c.extractString(hit.Source, "conteudoMovimento"),
			IsPublic:    !c.extractBool(hit.Source, "sigiloso"),
			Metadata:    c.extractMapInterface(hit.Source, "metadados"),
		}
		movements = append(movements, movement)
	}

	response.MovementData = &domain.MovementResponseData{
		Total:     esResponse.Hits.Total.Value,
		Page:      c.calculateCurrentPage(response),
		PageSize:  len(movements),
		Movements: movements,
	}

	return nil
}

// parsePartyData processa dados de partes
func (c *DataJudRealClient) parsePartyData(response *domain.DataJudResponse) error {
	var esResponse ElasticsearchResponse
	if err := json.Unmarshal(response.Body, &esResponse); err != nil {
		return fmt.Errorf("erro ao parsear resposta Elasticsearch: %w", err)
	}

	parties := make([]domain.PartyData, 0, len(esResponse.Hits.Hits))
	
	for _, hit := range esResponse.Hits.Hits {
		party := domain.PartyData{
			Type:     c.extractString(hit.Source, "tipoParte"),
			Name:     c.extractString(hit.Source, "nomeParte"),
			Document: c.extractString(hit.Source, "numeroDocumentoParte"),
			Role:     c.extractString(hit.Source, "papelParte"),
			Contact:  c.extractMapInterface(hit.Source, "contato"),
			Address:  c.extractMapInterface(hit.Source, "endereco"),
			Lawyer:   c.extractMapInterface(hit.Source, "advogado"),
		}
		parties = append(parties, party)
	}

	response.PartyData = &domain.PartyResponseData{
		Total:   esResponse.Hits.Total.Value,
		Parties: parties,
	}

	return nil
}

// parseBulkData processa dados de consulta em lote
func (c *DataJudRealClient) parseBulkData(response *domain.DataJudResponse, requestedNumbers []string) error {
	var esResponse ElasticsearchResponse
	if err := json.Unmarshal(response.Body, &esResponse); err != nil {
		return fmt.Errorf("erro ao parsear resposta Elasticsearch: %w", err)
	}

	// Organizar resultados por número de processo
	processMap := make(map[string]*domain.ProcessInfo)
	
	for _, hit := range esResponse.Hits.Hits {
		number := c.extractString(hit.Source, "numeroProcesso")
		if number != "" {
			processMap[number] = &domain.ProcessInfo{
				Number:     number,
				Class:      c.extractString(hit.Source, "classe"),
				Subject:    c.extractString(hit.Source, "assunto"),
				Court:      c.extractString(hit.Source, "tribunal"),
				Status:     c.extractString(hit.Source, "situacao"),
				LastUpdate: c.extractDate(hit.Source, "dataHoraUltimaAtualizacao"),
			}
		}
	}

	// Criar array ordenado conforme solicitado
	results := make([]*domain.BulkProcessResult, 0, len(requestedNumbers))
	for i, number := range requestedNumbers {
		result := &domain.BulkProcessResult{
			Index:         i,
			ProcessNumber: number,
			Found:         false,
			Process:       nil,
		}
		
		if process, exists := processMap[number]; exists {
			result.Found = true
			result.Process = process
		}
		
		results = append(results, result)
	}

	response.BulkData = &domain.BulkResponseData{
		Total:     len(requestedNumbers),
		Found:     len(processMap),
		NotFound:  len(requestedNumbers) - len(processMap),
		Processes: results,
	}

	return nil
}

// Métodos auxiliares para extração de dados
func (c *DataJudRealClient) extractString(source map[string]interface{}, key string) string {
	if value, ok := source[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func (c *DataJudRealClient) extractFloat(source map[string]interface{}, key string) float64 {
	if value, ok := source[key]; ok {
		switch v := value.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0
}

func (c *DataJudRealClient) extractBool(source map[string]interface{}, key string) bool {
	if value, ok := source[key]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
		if str, ok := value.(string); ok {
			return strings.ToLower(str) == "true" || str == "1"
		}
	}
	return false
}

func (c *DataJudRealClient) extractDate(source map[string]interface{}, key string) *time.Time {
	if value, ok := source[key]; ok {
		if str, ok := value.(string); ok {
			// Tentar diferentes formatos de data
			formats := []string{
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05",
				"2006-01-02 15:04:05",
				"2006-01-02",
			}
			
			for _, format := range formats {
				if t, err := time.Parse(format, str); err == nil {
					return &t
				}
			}
		}
	}
	return nil
}

func (c *DataJudRealClient) extractTimeWithDefault(source map[string]interface{}, key string) time.Time {
	if datePtr := c.extractDate(source, key); datePtr != nil {
		return *datePtr
	}
	return time.Now()
}

func (c *DataJudRealClient) extractMapInterface(source map[string]interface{}, key string) map[string]interface{} {
	if value, ok := source[key]; ok {
		if mapValue, ok := value.(map[string]interface{}); ok {
			return mapValue
		}
	}
	return make(map[string]interface{})
}

func (c *DataJudRealClient) getPageSize(req *domain.DataJudRequest) int {
	if pageSize, ok := req.Parameters["page_size"].(int); ok && pageSize > 0 {
		if pageSize > 1000 { // Limite máximo
			return 1000
		}
		return pageSize
	}
	return 100 // Padrão
}

func (c *DataJudRealClient) getPageOffset(req *domain.DataJudRequest) int {
	if page, ok := req.Parameters["page"].(int); ok && page > 0 {
		pageSize := c.getPageSize(req)
		return (page - 1) * pageSize
	}
	return 0
}

func (c *DataJudRealClient) calculateCurrentPage(response *domain.DataJudResponse) int {
	// Implementar lógica baseada nos parâmetros da requisição original
	return 1
}

func (c *DataJudRealClient) shouldRetry(err error, attempt int) bool {
	if attempt >= c.retryCount {
		return false
	}

	// Retry para erros de rede, timeout, 5xx
	if strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "500") ||
		strings.Contains(err.Error(), "502") ||
		strings.Contains(err.Error(), "503") ||
		strings.Contains(err.Error(), "504") {
		return true
	}

	return false
}

// TestConnection testa conectividade com a API DataJud
func (c *DataJudRealClient) TestConnection(ctx context.Context) error {
	// Usar endpoint de um tribunal para teste
	endpoint := fmt.Sprintf("%s/api_publica_tjsp/_search", c.baseURL)
	
	simpleQuery := &ElasticsearchQuery{
		Query: ElasticsearchQueryClause{
			Match: map[string]interface{}{
				"tribunal": "TJSP",
			},
		},
		Size: 1,
		From: 0,
	}

	queryBytes, err := json.Marshal(simpleQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar query de teste: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(queryBytes))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição de teste: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("APIKey %s", c.apiKey))

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

// Close fecha o cliente e limpa recursos
func (c *DataJudRealClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}