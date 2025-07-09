package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// MockClient implementação mock do cliente HTTP para desenvolvimento
type MockClient struct {
	tribunalMapper *TribunalMapper
}

// NewMockClient cria nova instância do mock client
func NewMockClient() *MockClient {
	return &MockClient{
		tribunalMapper: NewTribunalMapper(),
	}
}

// QueryProcess mock para consulta de processo
func (m *MockClient) QueryProcess(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	// Simular delay da API
	time.Sleep(100 * time.Millisecond)
	
	// Gerar dados mock baseados no número do processo
	processNumber := req.ProcessNumber
	tribunal := m.tribunalMapper.GetTribunal(req.CourtID)
	
	mockProcess := &domain.ProcessInfo{
		Number:     processNumber,
		Class:      "PROCEDIMENTO COMUM",
		Subject:    "Direito do Consumidor",
		Court:      req.CourtID,
		Instance:   "1ª Vara Cível",
		Status:     "Em andamento",
		Judge:      "Dr. João Silva",
		StartDate:  timePtr(time.Now().AddDate(0, -6, 0)),
		LastUpdate: timePtr(time.Now().AddDate(0, 0, -7)),
		Value:      15000.50,
		SecretLevel: "Público",
		Priority:   "Normal",
		ElectronicJudge: true,
	}
	
	// Adaptar dados baseados no tribunal
	if tribunal != nil {
		mockProcess.Court = tribunal.Name
		switch tribunal.Type {
		case "supremo":
			mockProcess.Instance = "Supremo Tribunal Federal"
			mockProcess.Class = "AÇÃO DIRETA DE INCONSTITUCIONALIDADE"
		case "superior":
			mockProcess.Instance = "Superior Tribunal de Justiça"
			mockProcess.Class = "RECURSO ESPECIAL"
		case "federal":
			mockProcess.Instance = "1ª Vara Federal"
			mockProcess.Class = "MANDADO DE SEGURANÇA"
		case "trabalhista":
			mockProcess.Instance = "1ª Vara do Trabalho"
			mockProcess.Class = "RECLAMAÇÃO TRABALHISTA"
		}
	}
	
	// Criar resposta mock
	processData := &domain.ProcessResponseData{
		Number:    mockProcess.Number,
		Title:     mockProcess.Class,
		Subject:   map[string]interface{}{"description": mockProcess.Subject},
		Court:     mockProcess.Court,
		Status:    mockProcess.Status,
		Stage:     mockProcess.Instance,
		CreatedAt: func() time.Time {
			if mockProcess.StartDate != nil {
				return *mockProcess.StartDate
			}
			return time.Now()
		}(),
		UpdatedAt: func() time.Time {
			if mockProcess.LastUpdate != nil {
				return *mockProcess.LastUpdate
			}
			return time.Now()
		}(),
		Parties:   []domain.PartyData{},
		Movements: []domain.MovementData{},
	}
	
	mockBody, _ := json.Marshal(map[string]interface{}{
		"hits": map[string]interface{}{
			"total": map[string]interface{}{
				"value": 1,
			},
			"hits": []map[string]interface{}{
				{
					"_source": map[string]interface{}{
						"numeroProcesso":    processNumber,
						"classeProcessual":  mockProcess.Class,
						"assunto":           mockProcess.Subject,
						"orgaoJulgador":     mockProcess.Court,
						"dataAjuizamento":   mockProcess.StartDate,
						"valorCausa":        mockProcess.Value,
						"grauSigilo":        mockProcess.SecretLevel,
						"statusProcesso":    mockProcess.Status,
						"nomeJuiz":          mockProcess.Judge,
						"instancia":         mockProcess.Instance,
						"prioridade":        mockProcess.Priority,
						"juizEletronico":    mockProcess.ElectronicJudge,
					},
				},
			},
		},
	})
	
	response := &domain.DataJudResponse{
		ID:          uuid.New(),
		StatusCode:  200,
		Body:        mockBody,
		Headers:     map[string]string{"Content-Type": "application/json"},
		ProcessData: processData,
		FromCache:   false,
		ReceivedAt:  time.Now(),
		Size:        int64(len(mockBody)),
		Duration:    150, // milliseconds
	}
	
	return response, nil
}

// QueryMovements mock para consulta de movimentações
func (m *MockClient) QueryMovements(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	// Simular delay da API
	time.Sleep(80 * time.Millisecond)
	
	// Gerar movimentações mock
	movements := []domain.MovementData{
		{
			Sequence:    1,
			Date:        time.Now().AddDate(0, 0, -1),
			Code:        "123",
			Type:        "JUNTADA",
			Title:       "Juntada de Petição",
			Description: "Petição de esclarecimentos protocolada",
			Content:     "Conteúdo da movimentação de juntada",
			IsPublic:    true,
			Metadata:    map[string]interface{}{"responsible": "Advogado da parte autora", "has_document": true},
		},
		{
			Sequence:    2,
			Date:        time.Now().AddDate(0, 0, -7),
			Code:        "456",
			Type:        "CONCLUSAO",
			Title:       "Conclusão para Despacho",
			Description: "Autos conclusos ao magistrado",
			Content:     "Processo em andamento",
			IsPublic:    true,
			Metadata:    map[string]interface{}{"responsible": "Cartório", "has_document": false},
		},
		{
			Sequence:    3,
			Date:        time.Now().AddDate(0, 0, -14),
			Code:        "789",
			Type:        "AUDIENCIA",
			Title:       "Audiência de Conciliação",
			Description: "Audiência realizada sem acordo",
			Content:     "Audiência registrada no sistema",
			IsPublic:    true,
			Metadata:    map[string]interface{}{"responsible": "Conciliador", "has_document": true},
		},
	}
	
	// Aplicar paginação
	page := getPageFromRequest(req)
	pageSize := getPageSizeFromRequest(req)
	
	start := (page - 1) * pageSize
	end := start + pageSize
	
	var paginatedMovements []domain.MovementData
	if start >= len(movements) {
		paginatedMovements = []domain.MovementData{}
	} else if end > len(movements) {
		paginatedMovements = movements[start:]
	} else {
		paginatedMovements = movements[start:end]
	}
	
	movementData := &domain.MovementResponseData{
		Total:     3,
		Page:      page,
		PageSize:  pageSize,
		Movements: paginatedMovements,
	}
	
	mockBody, _ := json.Marshal(map[string]interface{}{
		"hits": map[string]interface{}{
			"total": map[string]interface{}{
				"value": 3,
			},
			"hits": buildMovementHits(paginatedMovements),
		},
	})
	
	response := &domain.DataJudResponse{
		ID:           uuid.New(),
		StatusCode:   200,
		Body:         mockBody,
		Headers:      map[string]string{"Content-Type": "application/json"},
		MovementData: movementData,
		FromCache:    false,
		ReceivedAt:   time.Now(),
		Size:         int64(len(mockBody)),
		Duration:     120, // milliseconds
	}
	
	return response, nil
}

// QueryParties mock para consulta de partes
func (m *MockClient) QueryParties(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	// Simular delay da API
	time.Sleep(90 * time.Millisecond)
	
	parties := []domain.PartyData{
		{
			Type:     "PESSOA_FISICA",
			Name:     "João da Silva",
			Document: "123.456.789-00",
			Role:     "REQUERENTE",
			Contact:  map[string]interface{}{"email": "joao@email.com", "phone": "11999999999"},
			Address:  map[string]interface{}{"street": "Rua A, 123", "city": "São Paulo", "state": "SP"},
			Lawyer:   map[string]interface{}{"name": "Dr. Pedro Advogado", "oab": "OAB/SP 123456", "uf": "SP"},
		},
		{
			Type:     "PESSOA_JURIDICA",
			Name:     "Empresa XYZ Ltda",
			Document: "12.345.678/0001-90",
			Role:     "REQUERIDO",
			Contact:  map[string]interface{}{"email": "contato@empresa.com", "phone": "1133333333"},
			Address:  map[string]interface{}{"street": "Av. Principal, 456", "city": "São Paulo", "state": "SP"},
			Lawyer:   map[string]interface{}{"name": "Dra. Ana Advogada", "oab": "OAB/SP 654321", "uf": "SP"},
		},
	}
	
	partyData := &domain.PartyResponseData{
		Total:   len(parties),
		Parties: parties,
	}
	
	mockBody, _ := json.Marshal(map[string]interface{}{
		"hits": map[string]interface{}{
			"total": map[string]interface{}{
				"value": len(parties),
			},
			"hits": buildPartyHits(parties),
		},
	})
	
	response := &domain.DataJudResponse{
		ID:         uuid.New(),
		StatusCode: 200,
		Body:       mockBody,
		Headers:    map[string]string{"Content-Type": "application/json"},
		PartyData:  partyData,
		FromCache:  false,
		ReceivedAt: time.Now(),
		Size:       int64(len(mockBody)),
		Duration:   130, // milliseconds
	}
	
	return response, nil
}

// BulkQuery mock para consulta em lote
func (m *MockClient) BulkQuery(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
	// Simular delay maior para bulk
	time.Sleep(300 * time.Millisecond)
	
	processNumbers, ok := req.Parameters["process_numbers"].([]string)
	if !ok {
		return nil, fmt.Errorf("parâmetro process_numbers não encontrado")
	}
	
	var processes []*domain.BulkProcessResult
	found := 0
	
	for i, number := range processNumbers {
		// Simular alguns processos não encontrados
		isFound := (i+1)%4 != 0 // 75% encontrados
		
		if isFound {
			found++
			processes = append(processes, &domain.BulkProcessResult{
				Index:         i,
				ProcessNumber: number,
				Found:         true,
				Process: &domain.ProcessInfo{
					Number:     number,
					Class:      "PROCEDIMENTO COMUM",
					Subject:    "Direito do Consumidor",
					Court:      req.CourtID,
					Instance:   "1ª Vara Cível",
					Status:     "Em andamento",
					StartDate:  timePtr(time.Now().AddDate(0, -6, 0)),
					Value:      15000.50,
				},
			})
		} else {
			processes = append(processes, &domain.BulkProcessResult{
				Index:         i,
				ProcessNumber: number,
				Found:         false,
			})
		}
	}
	
	bulkData := &domain.BulkResponseData{
		Total:     len(processNumbers),
		Found:     found,
		NotFound:  len(processNumbers) - found,
		Processes: processes,
	}
	
	mockBody, _ := json.Marshal(map[string]interface{}{
		"hits": map[string]interface{}{
			"total": map[string]interface{}{
				"value": len(processNumbers),
			},
			"hits": buildBulkHits(processes),
		},
	})
	
	response := &domain.DataJudResponse{
		ID:         uuid.New(),
		StatusCode: 200,
		Body:       mockBody,
		Headers:    map[string]string{"Content-Type": "application/json"},
		BulkData:   bulkData,
		FromCache:  false,
		ReceivedAt: time.Now(),
		Size:       int64(len(mockBody)),
		Duration:   450, // milliseconds
	}
	
	return response, nil
}

// TestConnection mock para teste de conexão
func (m *MockClient) TestConnection(ctx context.Context) error {
	// Simular delay do teste
	time.Sleep(50 * time.Millisecond)
	
	// Mock sempre retorna sucesso
	return nil
}

// Close mock para fechamento do cliente
func (m *MockClient) Close() error {
	return nil
}

// Funções auxiliares

func timePtr(t time.Time) *time.Time {
	return &t
}

func getPageFromRequest(req *domain.DataJudRequest) int {
	if page, ok := req.Parameters["page"]; ok {
		if pageInt, ok := page.(int); ok {
			return pageInt
		}
	}
	return 1
}

func getPageSizeFromRequest(req *domain.DataJudRequest) int {
	if pageSize, ok := req.Parameters["page_size"]; ok {
		if pageSizeInt, ok := pageSize.(int); ok {
			return pageSizeInt
		}
	}
	return 50
}

func buildMovementHits(movements []domain.MovementData) []map[string]interface{} {
	hits := make([]map[string]interface{}, len(movements))
	for i, mov := range movements {
		hits[i] = map[string]interface{}{
			"_source": map[string]interface{}{
				"sequencia":          mov.Sequence,
				"dataMovimento":      mov.Date,
				"codigoMovimento":    mov.Code,
				"tipoMovimento":      mov.Type,
				"tituloMovimento":    mov.Title,
				"descricaoMovimento": mov.Description,
				"conteudoMovimento":  mov.Content,
				"publico":            mov.IsPublic,
				"metadados":          mov.Metadata,
			},
		}
	}
	return hits
}

func buildPartyHits(parties []domain.PartyData) []map[string]interface{} {
	hits := make([]map[string]interface{}, len(parties))
	for i, party := range parties {
		hits[i] = map[string]interface{}{
			"_source": map[string]interface{}{
				"nome":        party.Name,
				"documento":   party.Document,
				"tipoPessoa":  party.Type,
				"polo":        party.Role,
				"contato":     party.Contact,
				"endereco":    party.Address,
				"advogado":    party.Lawyer,
			},
		}
	}
	return hits
}

func buildBulkHits(processes []*domain.BulkProcessResult) []map[string]interface{} {
	hits := make([]map[string]interface{}, 0)
	for _, proc := range processes {
		if proc.Found && proc.Process != nil {
			hits = append(hits, map[string]interface{}{
				"_source": map[string]interface{}{
					"numeroProcesso":   proc.Process.Number,
					"classeProcessual": proc.Process.Class,
					"assunto":          proc.Process.Subject,
					"orgaoJulgador":    proc.Process.Court,
					"dataAjuizamento":  proc.Process.StartDate,
					"valorCausa":       proc.Process.Value,
					"statusProcesso":   proc.Process.Status,
					"instancia":        proc.Process.Instance,
				},
			})
		}
	}
	return hits
}