package tools

import (
	"fmt"

	"github.com/direito-lux/mcp-service/internal/domain"
)

// RegisterSearchTools registra as ferramentas de busca
func RegisterSearchTools(registry *domain.ToolRegistry) error {
	tools := []*domain.Tool{
		{
			Name:         "advanced_search",
			Description:  "Busca avançada com filtros complexos",
			Category:     "search",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "query", Type: "string", Description: "Consulta de busca", Required: true},
				{Name: "indices", Type: "array", Description: "Índices para buscar", 
					Required: false, Default: []string{"processes", "jurisprudence", "documents"}},
				{Name: "filters", Type: "object", Description: "Filtros complexos", Required: false},
				{Name: "aggregations", Type: "array", Description: "Agregações desejadas", Required: false},
				{Name: "sort_by", Type: "string", Description: "Campo para ordenação", Required: false},
				{Name: "page_size", Type: "number", Description: "Tamanho da página", 
					Required: false, Default: 20},
			},
			Handler: handleAdvancedSearch,
		},
		{
			Name:         "search_suggestions",
			Description:  "Sugestões automáticas de busca",
			Category:     "search",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "partial_query", Type: "string", Description: "Consulta parcial", Required: true},
				{Name: "context", Type: "string", Description: "Contexto da busca", 
					Required: false, Default: "general", 
					Enum: []string{"processes", "jurisprudence", "documents", "general"}},
				{Name: "max_suggestions", Type: "number", Description: "Número máximo de sugestões", 
					Required: false, Default: 5},
			},
			Handler: handleSearchSuggestions,
		},
	}
	
	// Registrar todas as ferramentas
	for _, tool := range tools {
		if err := registry.RegisterTool(tool); err != nil {
			return fmt.Errorf("erro ao registrar ferramenta %s: %w", tool.Name, err)
		}
	}
	
	return nil
}

// handleAdvancedSearch executa busca avançada
func handleAdvancedSearch(params map[string]interface{}) (interface{}, error) {
	query := params["query"].(string)
	indices := []string{"processes", "jurisprudence", "documents"}
	if val, ok := params["indices"].([]string); ok {
		indices = val
	}
	
	// Usar indices na busca
	_ = indices
	
	// TODO: Implementar integração real com Search Service
	
	result := map[string]interface{}{
		"query": query,
		"total_results": 342,
		"took_ms": 125,
		"hits": []map[string]interface{}{
			{
				"_index": "processes",
				"_id": "proc-789",
				"_score": 0.98,
				"_source": map[string]interface{}{
					"number": "3001234-90.2024.8.26.0300",
					"title": "Ação de cobrança com pedido de tutela",
					"client": "Empresa ABC Ltda",
					"status": "active",
					"highlight": map[string]string{
						"title": "Ação de <em>cobrança</em> com pedido de tutela",
					},
				},
			},
			{
				"_index": "jurisprudence",
				"_id": "jur-456",
				"_score": 0.92,
				"_source": map[string]interface{}{
					"court": "STJ",
					"number": "REsp 1.876.543/RS",
					"summary": "Cobrança. Título executivo. Prescrição.",
					"date": "2024-04-20",
					"highlight": map[string]string{
						"summary": "<em>Cobrança</em>. Título executivo. Prescrição.",
					},
				},
			},
		},
		"aggregations": map[string]interface{}{
			"by_index": map[string]interface{}{
				"processes": 156,
				"jurisprudence": 142,
				"documents": 44,
			},
			"by_date": map[string]interface{}{
				"last_week": 45,
				"last_month": 128,
				"last_year": 342,
			},
			"by_status": map[string]interface{}{
				"active": 234,
				"archived": 89,
				"suspended": 19,
			},
		},
		"facets": []map[string]interface{}{
			{
				"field": "court",
				"values": []map[string]interface{}{
					{"value": "TJ-SP", "count": 125},
					{"value": "STJ", "count": 89},
					{"value": "STF", "count": 34},
				},
			},
		},
	}
	
	return result, nil
}

// handleSearchSuggestions gera sugestões de busca
func handleSearchSuggestions(params map[string]interface{}) (interface{}, error) {
	partialQuery := params["partial_query"].(string)
	context := "general"
	if val, ok := params["context"].(string); ok {
		context = val
	}
	maxSuggestions := 5
	if val, ok := params["max_suggestions"].(float64); ok {
		maxSuggestions = int(val)
	}
	
	// TODO: Implementar integração real com Search Service
	
	suggestions := []map[string]interface{}{}
	
	// Gerar sugestões baseadas no contexto
	switch context {
	case "processes":
		suggestions = []map[string]interface{}{
			{
				"text": partialQuery + " em andamento",
				"score": 0.95,
				"type": "completion",
			},
			{
				"text": partialQuery + " arquivados",
				"score": 0.88,
				"type": "completion",
			},
			{
				"text": fmt.Sprintf("processo %s TJ-SP", partialQuery),
				"score": 0.85,
				"type": "suggestion",
			},
		}
		
	case "jurisprudence":
		suggestions = []map[string]interface{}{
			{
				"text": partialQuery + " responsabilidade civil",
				"score": 0.92,
				"type": "completion",
			},
			{
				"text": partialQuery + " STJ súmula",
				"score": 0.87,
				"type": "suggestion",
			},
			{
				"text": fmt.Sprintf("jurisprudência %s dano moral", partialQuery),
				"score": 0.84,
				"type": "suggestion",
			},
		}
		
	default:
		suggestions = []map[string]interface{}{
			{
				"text": partialQuery + " processo",
				"score": 0.90,
				"type": "completion",
			},
			{
				"text": partialQuery + " jurisprudência",
				"score": 0.85,
				"type": "completion",
			},
			{
				"text": partialQuery + " petição",
				"score": 0.80,
				"type": "completion",
			},
		}
	}
	
	// Limitar ao número máximo de sugestões
	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}
	
	result := map[string]interface{}{
		"partial_query": partialQuery,
		"context": context,
		"suggestions": suggestions,
		"related_searches": []string{
			"busca avançada " + partialQuery,
			partialQuery + " últimos 30 dias",
			partialQuery + " tribunal superior",
		},
		"popular_filters": []map[string]interface{}{
			{"field": "date_range", "value": "last_30_days", "label": "Últimos 30 dias"},
			{"field": "court", "value": "superior_courts", "label": "Tribunais Superiores"},
			{"field": "status", "value": "active", "label": "Processos Ativos"},
		},
	}
	
	return result, nil
}