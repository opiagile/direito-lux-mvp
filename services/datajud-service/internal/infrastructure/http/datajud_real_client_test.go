package http

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// MockDataJudRealClient mock do cliente para testes
type MockDataJudRealClient struct {
	*DataJudRealClient
	mockEnabled bool
}

// NewMockDataJudRealClient cria cliente mock para testes
func NewMockDataJudRealClient() *MockDataJudRealClient {
	config := DataJudRealClientConfig{
		BaseURL:    "https://api-publica.datajud.cnj.jus.br",
		APIKey:     "test-api-key",
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: 1 * time.Second,
	}
	
	realClient := NewDataJudRealClient(config)
	
	return &MockDataJudRealClient{
		DataJudRealClient: realClient,
		mockEnabled:       true,
	}
}

// TestNewDataJudRealClient testa criação do cliente
func TestNewDataJudRealClient(t *testing.T) {
	config := DataJudRealClientConfig{
		BaseURL:    "https://api-publica.datajud.cnj.jus.br",
		APIKey:     "test-api-key",
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: 1 * time.Second,
	}
	
	client := NewDataJudRealClient(config)
	
	assert.NotNil(t, client)
	assert.Equal(t, config.BaseURL, client.baseURL)
	assert.Equal(t, config.APIKey, client.apiKey)
	assert.Equal(t, config.Timeout, client.timeout)
	assert.Equal(t, config.RetryCount, client.retryCount)
	assert.Equal(t, config.RetryDelay, client.retryDelay)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.tribunalMapper)
}

// TestTribunalMapper testa o mapeamento de tribunais
func TestTribunalMapper(t *testing.T) {
	mapper := NewTribunalMapper()
	
	// Testar tribunais superiores
	stf := mapper.GetTribunal("STF")
	require.NotNil(t, stf)
	assert.Equal(t, "STF", stf.Code)
	assert.Equal(t, "api_publica_stf", stf.Endpoint)
	assert.Equal(t, "supremo", stf.Type)
	assert.True(t, stf.Active)
	
	// Testar tribunais estaduais
	tjsp := mapper.GetTribunal("TJSP")
	require.NotNil(t, tjsp)
	assert.Equal(t, "TJSP", tjsp.Code)
	assert.Equal(t, "api_publica_tjsp", tjsp.Endpoint)
	assert.Equal(t, "estadual", tjsp.Type)
	assert.Equal(t, "SP", tjsp.Region)
	assert.True(t, tjsp.Active)
	
	// Testar tribunais federais
	trf3 := mapper.GetTribunal("TRF3")
	require.NotNil(t, trf3)
	assert.Equal(t, "TRF3", trf3.Code)
	assert.Equal(t, "api_publica_trf3", trf3.Endpoint)
	assert.Equal(t, "federal", trf3.Type)
	assert.True(t, trf3.Active)
	
	// Testar tribunal inexistente
	invalid := mapper.GetTribunal("INVALID")
	assert.Nil(t, invalid)
}

// TestTribunalMapperAlternativeCodes testa códigos alternativos
func TestTribunalMapperAlternativeCodes(t *testing.T) {
	mapper := NewTribunalMapper()
	
	// Testar diferentes formatos do mesmo tribunal
	testCases := []string{
		"TRF1",
		"TRF01",
		"trf1",
		"trf01",
		"TRF-1",
		"TRF_1",
		"TRF 1",
	}
	
	for _, code := range testCases {
		tribunal := mapper.GetTribunal(code)
		assert.NotNil(t, tribunal, "Código %s deveria ser válido", code)
		if tribunal != nil {
			assert.Equal(t, "TRF1", tribunal.Code)
		}
	}
}

// TestElasticsearchQueryBuilder testa construção de queries
func TestElasticsearchQueryBuilder(t *testing.T) {
	builder := NewElasticsearchQueryBuilder()
	
	// Testar query por número de processo
	query := builder.ProcessByNumber("0001234-56.2023.8.26.0001").Build()
	
	assert.NotNil(t, query)
	assert.NotNil(t, query.Query.Bool)
	assert.True(t, len(query.Query.Bool.Must) > 0)
	assert.Equal(t, 100, query.Size)
	assert.Equal(t, 0, query.From)
}

// TestElasticsearchQueryBuilderMovements testa query de movimentações
func TestElasticsearchQueryBuilderMovements(t *testing.T) {
	builder := NewElasticsearchQueryBuilder()
	
	query := builder.MovementsByProcess("0001234-56.2023.8.26.0001").Build()
	
	assert.NotNil(t, query)
	assert.NotNil(t, query.Query.Bool)
	assert.True(t, len(query.Query.Bool.Must) >= 2) // Número do processo + tipo documento
	assert.True(t, len(query.Sort) > 0)
}

// TestElasticsearchQueryBuilderFilters testa filtros
func TestElasticsearchQueryBuilderFilters(t *testing.T) {
	builder := NewElasticsearchQueryBuilder()
	
	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()
	
	query := builder.
		ProcessByNumber("0001234-56.2023.8.26.0001").
		CourtFilter("TJSP").
		ClassFilter("PROCEDIMENTO COMUM").
		DateRange("dataAjuizamento", &from, &to).
		Pagination(2, 50).
		Build()
	
	assert.NotNil(t, query)
	assert.True(t, len(query.Query.Bool.Must) > 0)
	assert.True(t, len(query.Query.Bool.Filter) > 0)
	assert.Equal(t, 50, query.Size)
	assert.Equal(t, 50, query.From) // Página 2
}

// TestQueryTemplates testa templates de queries
func TestQueryTemplates(t *testing.T) {
	templates := NewQueryTemplates()
	
	// Testar template de processo por número
	query := templates.GetProcessByNumber("0001234-56.2023.8.26.0001")
	assert.NotNil(t, query)
	assert.NotNil(t, query.Query.Bool)
	
	// Testar template de movimentações recentes
	query = templates.GetRecentMovements("0001234-56.2023.8.26.0001", 30)
	assert.NotNil(t, query)
	assert.True(t, len(query.Query.Bool.Filter) > 0) // Filtro de data
	
	// Testar template de processos por parte
	query = templates.GetProcessesByParty("João da Silva", "TJSP")
	assert.NotNil(t, query)
	assert.True(t, len(query.Query.Bool.Should) > 0) // Busca por nome ou documento
}

// TestNormalizeProcessNumber testa normalização de números
func TestNormalizeProcessNumber(t *testing.T) {
	builder := NewElasticsearchQueryBuilder()
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"0001234-56.2023.8.26.0001", "0001234562023826001"},
		{"0001234 56 2023 8 26 0001", "0001234562023826001"},
		{"0001234.56.2023.8.26.0001", "0001234562023826001"},
		{"0001234-56-2023-8-26-0001", "0001234562023826001"},
		{"123456", "123456"},
		{"", ""},
	}
	
	for _, tc := range testCases {
		result := builder.normalizeProcessNumber(tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}

// TestDataJudResponseParsing testa parsing de respostas
func TestDataJudResponseParsing(t *testing.T) {
	client := NewMockDataJudRealClient()
	
	// Mock response data
	mockResponse := &domain.DataJudResponse{
		ID:         uuid.New(),
		StatusCode: 200,
		Body:       []byte(`{"hits": {"hits": [{"_source": {"numeroProcesso": "0001234-56.2023.8.26.0001", "assunto": "Teste", "tribunal": "TJSP"}}]}}`),
		Headers:    make(map[string]string),
	}
	
	// Testar parse de dados de processo
	err := client.parseProcessData(mockResponse)
	assert.NoError(t, err)
	assert.NotNil(t, mockResponse.ProcessData)
	assert.True(t, mockResponse.ProcessData.Found)
	assert.NotNil(t, mockResponse.ProcessData.Process)
	assert.Equal(t, "0001234-56.2023.8.26.0001", mockResponse.ProcessData.Process.Number)
	assert.Equal(t, "Teste", mockResponse.ProcessData.Process.Subject)
	assert.Equal(t, "TJSP", mockResponse.ProcessData.Process.Court)
}

// TestDataJudResponseParsingEmpty testa parsing de resposta vazia
func TestDataJudResponseParsingEmpty(t *testing.T) {
	client := NewMockDataJudRealClient()
	
	// Mock empty response
	mockResponse := &domain.DataJudResponse{
		ID:         uuid.New(),
		StatusCode: 200,
		Body:       []byte(`{"hits": {"hits": []}}`),
		Headers:    make(map[string]string),
	}
	
	// Testar parse de dados vazios
	err := client.parseProcessData(mockResponse)
	assert.NoError(t, err)
	assert.NotNil(t, mockResponse.ProcessData)
	assert.False(t, mockResponse.ProcessData.Found)
	assert.Nil(t, mockResponse.ProcessData.Process)
}

// TestExtractMethods testa métodos de extração
func TestExtractMethods(t *testing.T) {
	client := NewMockDataJudRealClient()
	
	source := map[string]interface{}{
		"stringField":  "test string",
		"intField":     42,
		"floatField":   3.14,
		"boolField":    true,
		"dateField":    "2023-12-01T10:30:00Z",
		"emptyField":   "",
		"nullField":    nil,
	}
	
	// Testar extração de string
	assert.Equal(t, "test string", client.extractString(source, "stringField"))
	assert.Equal(t, "", client.extractString(source, "emptyField"))
	assert.Equal(t, "", client.extractString(source, "nullField"))
	assert.Equal(t, "", client.extractString(source, "nonExistentField"))
	
	// Testar extração de float
	assert.Equal(t, 3.14, client.extractFloat(source, "floatField"))
	assert.Equal(t, 42.0, client.extractFloat(source, "intField"))
	assert.Equal(t, 0.0, client.extractFloat(source, "nonExistentField"))
	
	// Testar extração de bool
	assert.True(t, client.extractBool(source, "boolField"))
	assert.False(t, client.extractBool(source, "nonExistentField"))
	
	// Testar extração de data
	date := client.extractDate(source, "dateField")
	assert.NotNil(t, date)
	assert.Equal(t, 2023, date.Year())
	assert.Equal(t, time.December, date.Month())
	assert.Equal(t, 1, date.Day())
	
	// Testar data inválida
	invalidDate := client.extractDate(source, "stringField")
	assert.Nil(t, invalidDate)
}

// TestPagination testa configuração de paginação
func TestPagination(t *testing.T) {
	client := NewMockDataJudRealClient()
	
	// Mock request
	req := &domain.DataJudRequest{
		Parameters: make(map[string]interface{}),
	}
	
	// Testar página size padrão
	assert.Equal(t, 100, client.getPageSize(req))
	
	// Testar página size customizado
	req.Parameters["page_size"] = 50
	assert.Equal(t, 50, client.getPageSize(req))
	
	// Testar limite máximo
	req.Parameters["page_size"] = 2000
	assert.Equal(t, 1000, client.getPageSize(req))
	
	// Testar offset padrão
	req.Parameters = make(map[string]interface{})
	assert.Equal(t, 0, client.getPageOffset(req))
	
	// Testar offset customizado
	req.Parameters["page"] = 3
	req.Parameters["page_size"] = 50
	assert.Equal(t, 100, client.getPageOffset(req)) // (3-1) * 50 = 100
}

// TestShouldRetry testa lógica de retry
func TestShouldRetry(t *testing.T) {
	client := NewMockDataJudRealClient()
	
	// Testar limite de tentativas
	assert.False(t, client.shouldRetry(nil, 5)) // Excede limit
	
	// Testar erros que devem ter retry
	retryableErrors := []string{
		"timeout",
		"connection refused",
		"500 Internal Server Error",
		"502 Bad Gateway",
		"503 Service Unavailable",
		"504 Gateway Timeout",
	}
	
	for _, errMsg := range retryableErrors {
		err := &mockError{message: errMsg}
		assert.True(t, client.shouldRetry(err, 1), "Erro %s deveria ter retry", errMsg)
	}
	
	// Testar erro que não deve ter retry
	err := &mockError{message: "400 Bad Request"}
	assert.False(t, client.shouldRetry(err, 1))
}

// mockError erro mock para testes
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

// TestBulkQuery testa consulta em lote
func TestBulkQuery(t *testing.T) {
	client := NewMockDataJudRealClient()
	
	// Mock request
	req := &domain.DataJudRequest{
		CourtID: "TJSP",
		Parameters: map[string]interface{}{
			"process_numbers": []string{
				"0001234-56.2023.8.26.0001",
				"0001234-56.2023.8.26.0002",
				"0001234-56.2023.8.26.0003",
			},
		},
	}
	
	// Mock provider
	provider := &domain.CNPJProvider{
		ID:   uuid.New(),
		CNPJ: "12345678000123",
	}
	
	// Testar construção da query
	builder := NewElasticsearchQueryBuilder()
	query := builder.BulkProcesses(req.Parameters["process_numbers"].([]string)).Build()
	
	assert.NotNil(t, query)
	assert.Equal(t, 3, query.Size) // Tamanho igual ao número de processos
	assert.True(t, len(query.Sort) > 0)
}

// TestConnectionTest testa teste de conexão
func TestConnectionTest(t *testing.T) {
	// Pular se não há chave de API real
	if testing.Short() {
		t.Skip("Pulando teste de conexão em modo short")
	}
	
	client := NewMockDataJudRealClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Este teste precisaria de uma API key real para funcionar
	// Em ambiente de teste, podemos mockar ou pular
	err := client.TestConnection(ctx)
	
	// Em ambiente mock, esperamos erro de conexão
	assert.Error(t, err)
}

// BenchmarkQueryBuilder benchmark para construção de queries
func BenchmarkQueryBuilder(b *testing.B) {
	processNumber := "0001234-56.2023.8.26.0001"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewElasticsearchQueryBuilder()
		query := builder.ProcessByNumber(processNumber).Build()
		_ = query
	}
}

// BenchmarkTribunalMapper benchmark para mapeamento de tribunais
func BenchmarkTribunalMapper(b *testing.B) {
	mapper := NewTribunalMapper()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tribunal := mapper.GetTribunal("TJSP")
		_ = tribunal
	}
}

// BenchmarkNormalizeProcessNumber benchmark para normalização
func BenchmarkNormalizeProcessNumber(b *testing.B) {
	builder := NewElasticsearchQueryBuilder()
	processNumber := "0001234-56.2023.8.26.0001"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalized := builder.normalizeProcessNumber(processNumber)
		_ = normalized
	}
}