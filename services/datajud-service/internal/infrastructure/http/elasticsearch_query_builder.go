package http

import (
	"strings"
	"time"
)

// ElasticsearchQueryBuilder construtor de queries otimizadas para DataJud
type ElasticsearchQueryBuilder struct {
	query *ElasticsearchQuery
}

// NewElasticsearchQueryBuilder cria novo construtor de queries
func NewElasticsearchQueryBuilder() *ElasticsearchQueryBuilder {
	return &ElasticsearchQueryBuilder{
		query: &ElasticsearchQuery{
			Query: ElasticsearchQueryClause{
				Bool: &ElasticsearchBoolQuery{
					Must:   []ElasticsearchQueryClause{},
					Should: []ElasticsearchQueryClause{},
					Filter: []ElasticsearchQueryClause{},
				},
			},
			Size: 100,
			From: 0,
			Sort: []map[string]interface{}{},
		},
	}
}

// ProcessByNumber query otimizada para buscar processo por número
func (qb *ElasticsearchQueryBuilder) ProcessByNumber(processNumber string) *ElasticsearchQueryBuilder {
	// Normalizar número do processo
	normalizedNumber := qb.normalizeProcessNumber(processNumber)
	
	// Query principal por número exato
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"numeroProcesso.keyword": normalizedNumber,
		},
	})
	
	// Fallback para busca fuzzy se necessário
	if len(normalizedNumber) >= 20 { // Número CNJ completo
		qb.query.Query.Bool.Should = append(qb.query.Query.Bool.Should, ElasticsearchQueryClause{
			Match: map[string]interface{}{
				"numeroProcesso": map[string]interface{}{
					"query":     normalizedNumber,
					"fuzziness": "AUTO",
					"operator":  "AND",
				},
			},
		})
	}
	
	// Garantir que busca apenas processos ativos
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"tipoDocumento": "processo",
		},
	})
	
	return qb
}

// MovementsByProcess query otimizada para movimentações de processo
func (qb *ElasticsearchQueryBuilder) MovementsByProcess(processNumber string) *ElasticsearchQueryBuilder {
	normalizedNumber := qb.normalizeProcessNumber(processNumber)
	
	// Buscar movimentações do processo
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"numeroProcesso.keyword": normalizedNumber,
		},
	})
	
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"tipoDocumento": "movimentacao",
		},
	})
	
	// Ordenar por data mais recente
	qb.query.Sort = []map[string]interface{}{
		{"dataHora": map[string]string{"order": "desc"}},
		{"sequenciaMovimento": map[string]string{"order": "desc"}},
	}
	
	return qb
}

// PartiesByProcess query otimizada para partes de processo
func (qb *ElasticsearchQueryBuilder) PartiesByProcess(processNumber string) *ElasticsearchQueryBuilder {
	normalizedNumber := qb.normalizeProcessNumber(processNumber)
	
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"numeroProcesso.keyword": normalizedNumber,
		},
	})
	
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"tipoDocumento": "parte",
		},
	})
	
	// Ordenar por tipo de parte (autor, réu, etc.)
	qb.query.Sort = []map[string]interface{}{
		{"tipoParte.keyword": map[string]string{"order": "asc"}},
		{"nomeParte.keyword": map[string]string{"order": "asc"}},
	}
	
	return qb
}

// BulkProcesses query otimizada para múltiplos processos
func (qb *ElasticsearchQueryBuilder) BulkProcesses(processNumbers []string) *ElasticsearchQueryBuilder {
	if len(processNumbers) == 0 {
		return qb
	}
	
	// Normalizar todos os números
	normalizedNumbers := make([]string, len(processNumbers))
	for i, number := range processNumbers {
		normalizedNumbers[i] = qb.normalizeProcessNumber(number)
	}
	
	// Query com terms para múltiplos valores
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"numeroProcesso.keyword": normalizedNumbers,
		},
	})
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"tipoDocumento": "processo",
		},
	})
	
	// Ordenar por número do processo
	qb.query.Sort = []map[string]interface{}{
		{"numeroProcesso.keyword": map[string]string{"order": "asc"}},
	}
	
	// Aumentar tamanho da página para comportar todos os processos
	qb.query.Size = len(processNumbers)
	
	return qb
}

// SearchProcesses query otimizada para busca textual
func (qb *ElasticsearchQueryBuilder) SearchProcesses(searchTerm string) *ElasticsearchQueryBuilder {
	if searchTerm == "" {
		return qb
	}
	
	// Multi-match query em campos relevantes
	qb.query.Query.Bool.Must = append(qb.query.Query.Bool.Must, ElasticsearchQueryClause{
		Match: map[string]interface{}{
			"_all": map[string]interface{}{
				"query":                searchTerm,
				"fuzziness":           "AUTO",
				"minimum_should_match": "75%",
			},
		},
	})
	
	// Boost para campos mais importantes
	qb.query.Query.Bool.Should = append(qb.query.Query.Bool.Should, []ElasticsearchQueryClause{
		{
			Match: map[string]interface{}{
				"assunto": map[string]interface{}{
					"query": searchTerm,
					"boost": 3.0,
				},
			},
		},
		{
			Match: map[string]interface{}{
				"classe": map[string]interface{}{
					"query": searchTerm,
					"boost": 2.0,
				},
			},
		},
		{
			Match: map[string]interface{}{
				"nomeParte": map[string]interface{}{
					"query": searchTerm,
					"boost": 2.5,
				},
			},
		},
	}...)
	
	// Filtrar apenas processos
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"tipoDocumento": "processo",
		},
	})
	
	// Ordenar por relevância e data
	qb.query.Sort = []map[string]interface{}{
		{"_score": map[string]string{"order": "desc"}},
		{"dataHoraUltimaAtualizacao": map[string]string{"order": "desc"}},
	}
	
	return qb
}

// DateRange adiciona filtro de data
func (qb *ElasticsearchQueryBuilder) DateRange(field string, from, to *time.Time) *ElasticsearchQueryBuilder {
	if from == nil && to == nil {
		return qb
	}
	
	dateFilter := make(map[string]interface{})
	
	if from != nil {
		dateFilter["gte"] = from.Format("2006-01-02")
	}
	
	if to != nil {
		dateFilter["lte"] = to.Format("2006-01-02")
	}
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Range: map[string]interface{}{
			field: dateFilter,
		},
	})
	
	return qb
}

// CourtFilter adiciona filtro por tribunal
func (qb *ElasticsearchQueryBuilder) CourtFilter(courtCode string) *ElasticsearchQueryBuilder {
	if courtCode == "" {
		return qb
	}
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"tribunal.keyword": strings.ToUpper(courtCode),
		},
	})
	
	return qb
}

// ClassFilter adiciona filtro por classe processual
func (qb *ElasticsearchQueryBuilder) ClassFilter(classCode string) *ElasticsearchQueryBuilder {
	if classCode == "" {
		return qb
	}
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"classe.keyword": classCode,
		},
	})
	
	return qb
}

// SubjectFilter adiciona filtro por assunto
func (qb *ElasticsearchQueryBuilder) SubjectFilter(subjectCode string) *ElasticsearchQueryBuilder {
	if subjectCode == "" {
		return qb
	}
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"assunto.keyword": subjectCode,
		},
	})
	
	return qb
}

// StatusFilter adiciona filtro por status do processo
func (qb *ElasticsearchQueryBuilder) StatusFilter(status string) *ElasticsearchQueryBuilder {
	if status == "" {
		return qb
	}
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"situacao.keyword": status,
		},
	})
	
	return qb
}

// ValueRange adiciona filtro por valor da causa
func (qb *ElasticsearchQueryBuilder) ValueRange(minValue, maxValue *float64) *ElasticsearchQueryBuilder {
	if minValue == nil && maxValue == nil {
		return qb
	}
	
	valueFilter := make(map[string]interface{})
	
	if minValue != nil {
		valueFilter["gte"] = *minValue
	}
	
	if maxValue != nil {
		valueFilter["lte"] = *maxValue
	}
	
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Range: map[string]interface{}{
			"valorCausa": valueFilter,
		},
	})
	
	return qb
}

// PartyFilter adiciona filtro por parte (nome ou documento)
func (qb *ElasticsearchQueryBuilder) PartyFilter(partyIdentifier string) *ElasticsearchQueryBuilder {
	if partyIdentifier == "" {
		return qb
	}
	
	// Buscar por nome ou documento da parte
	qb.query.Query.Bool.Should = append(qb.query.Query.Bool.Should, []ElasticsearchQueryClause{
		{
			Match: map[string]interface{}{
				"nomeParte": map[string]interface{}{
					"query":     partyIdentifier,
					"fuzziness": "AUTO",
				},
			},
		},
		{
			Term: map[string]interface{}{
				"numeroDocumentoParte.keyword": partyIdentifier,
			},
		},
	}...)
	
	return qb
}

// LawyerFilter adiciona filtro por advogado
func (qb *ElasticsearchQueryBuilder) LawyerFilter(lawyerName string, oabNumber string) *ElasticsearchQueryBuilder {
	if lawyerName == "" && oabNumber == "" {
		return qb
	}
	
	if lawyerName != "" {
		qb.query.Query.Bool.Should = append(qb.query.Query.Bool.Should, ElasticsearchQueryClause{
			Match: map[string]interface{}{
				"nomeAdvogado": map[string]interface{}{
					"query":     lawyerName,
					"fuzziness": "AUTO",
				},
			},
		})
	}
	
	if oabNumber != "" {
		qb.query.Query.Bool.Should = append(qb.query.Query.Bool.Should, ElasticsearchQueryClause{
			Term: map[string]interface{}{
				"oabAdvogado.keyword": oabNumber,
			},
		})
	}
	
	return qb
}

// Pagination configura paginação
func (qb *ElasticsearchQueryBuilder) Pagination(page, pageSize int) *ElasticsearchQueryBuilder {
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 100
	}
	
	if pageSize > 1000 {
		pageSize = 1000 // Limite máximo da API
	}
	
	qb.query.From = (page - 1) * pageSize
	qb.query.Size = pageSize
	
	return qb
}

// SortByDate adiciona ordenação por data
func (qb *ElasticsearchQueryBuilder) SortByDate(field string, ascending bool) *ElasticsearchQueryBuilder {
	order := "desc"
	if ascending {
		order = "asc"
	}
	
	qb.query.Sort = append(qb.query.Sort, map[string]interface{}{
		field: map[string]string{"order": order},
	})
	
	return qb
}

// SortByRelevance adiciona ordenação por relevância
func (qb *ElasticsearchQueryBuilder) SortByRelevance() *ElasticsearchQueryBuilder {
	qb.query.Sort = append(qb.query.Sort, map[string]interface{}{
		"_score": map[string]string{"order": "desc"},
	})
	
	return qb
}

// OnlyActive filtrar apenas processos ativos
func (qb *ElasticsearchQueryBuilder) OnlyActive() *ElasticsearchQueryBuilder {
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"ativo": true,
		},
	})
	
	return qb
}

// ExcludeSecrets excluir processos em segredo de justiça
func (qb *ElasticsearchQueryBuilder) ExcludeSecrets() *ElasticsearchQueryBuilder {
	qb.query.Query.Bool.Filter = append(qb.query.Query.Bool.Filter, ElasticsearchQueryClause{
		Term: map[string]interface{}{
			"nivelSigilo": "0", // 0 = público
		},
	})
	
	return qb
}

// Build constrói a query final
func (qb *ElasticsearchQueryBuilder) Build() *ElasticsearchQuery {
	// Otimizações finais
	qb.optimizeQuery()
	
	return qb.query
}

// normalizeProcessNumber normaliza número do processo
func (qb *ElasticsearchQueryBuilder) normalizeProcessNumber(processNumber string) string {
	// Remove caracteres especiais
	normalized := strings.ReplaceAll(processNumber, ".", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	
	return strings.ToUpper(normalized)
}

// optimizeQuery otimiza a query final
func (qb *ElasticsearchQueryBuilder) optimizeQuery() {
	// Se não há must clauses, mas há should, promover should para must
	if len(qb.query.Query.Bool.Must) == 0 && len(qb.query.Query.Bool.Should) > 0 {
		qb.query.Query.Bool.Must = qb.query.Query.Bool.Should
		qb.query.Query.Bool.Should = []ElasticsearchQueryClause{}
	}
	
	// Garantir ordenação padrão se não especificada
	if len(qb.query.Sort) == 0 {
		qb.query.Sort = []map[string]interface{}{
			{"dataHoraUltimaAtualizacao": map[string]string{"order": "desc"}},
		}
	}
}

// QueryTemplates templates de queries pré-definidas
type QueryTemplates struct{}

// NewQueryTemplates cria nova instância de templates
func NewQueryTemplates() *QueryTemplates {
	return &QueryTemplates{}
}

// GetProcessByNumber template para buscar processo por número
func (qt *QueryTemplates) GetProcessByNumber(processNumber string) *ElasticsearchQuery {
	return NewElasticsearchQueryBuilder().
		ProcessByNumber(processNumber).
		ExcludeSecrets().
		Build()
}

// GetRecentMovements template para movimentações recentes
func (qt *QueryTemplates) GetRecentMovements(processNumber string, days int) *ElasticsearchQuery {
	from := time.Now().AddDate(0, 0, -days)
	
	return NewElasticsearchQueryBuilder().
		MovementsByProcess(processNumber).
		DateRange("dataHora", &from, nil).
		Pagination(1, 100).
		Build()
}

// GetProcessesByParty template para processos por parte
func (qt *QueryTemplates) GetProcessesByParty(partyIdentifier string, courtCode string) *ElasticsearchQuery {
	builder := NewElasticsearchQueryBuilder().
		PartyFilter(partyIdentifier).
		ExcludeSecrets().
		SortByDate("dataHoraUltimaAtualizacao", false).
		Pagination(1, 50)
	
	if courtCode != "" {
		builder = builder.CourtFilter(courtCode)
	}
	
	return builder.Build()
}

// GetProcessesByLawyer template para processos por advogado
func (qt *QueryTemplates) GetProcessesByLawyer(lawyerName, oabNumber string) *ElasticsearchQuery {
	return NewElasticsearchQueryBuilder().
		LawyerFilter(lawyerName, oabNumber).
		ExcludeSecrets().
		SortByDate("dataHoraUltimaAtualizacao", false).
		Pagination(1, 50).
		Build()
}

// GetProcessesByClassAndSubject template para processos por classe e assunto
func (qt *QueryTemplates) GetProcessesByClassAndSubject(classCode, subjectCode string) *ElasticsearchQuery {
	return NewElasticsearchQueryBuilder().
		ClassFilter(classCode).
		SubjectFilter(subjectCode).
		ExcludeSecrets().
		SortByDate("dataAjuizamento", false).
		Pagination(1, 100).
		Build()
}

// GetProcessesUpdatedAfter template para processos atualizados após data
func (qt *QueryTemplates) GetProcessesUpdatedAfter(afterDate time.Time, courtCode string) *ElasticsearchQuery {
	builder := NewElasticsearchQueryBuilder().
		DateRange("dataHoraUltimaAtualizacao", &afterDate, nil).
		ExcludeSecrets().
		SortByDate("dataHoraUltimaAtualizacao", false).
		Pagination(1, 100)
	
	if courtCode != "" {
		builder = builder.CourtFilter(courtCode)
	}
	
	return builder.Build()
}

// GetStatisticsQuery template para estatísticas
func (qt *QueryTemplates) GetStatisticsQuery(courtCode string, dateFrom, dateTo *time.Time) *ElasticsearchQuery {
	builder := NewElasticsearchQueryBuilder().
		CourtFilter(courtCode).
		ExcludeSecrets().
		Pagination(1, 0) // Apenas agregações
	
	if dateFrom != nil || dateTo != nil {
		builder = builder.DateRange("dataAjuizamento", dateFrom, dateTo)
	}
	
	return builder.Build()
}