package tools

import (
	"fmt"
	"time"

	"github.com/direito-lux/mcp-service/internal/domain"
)

// RegisterAITools registra as ferramentas de IA
func RegisterAITools(registry *domain.ToolRegistry) error {
	tools := []*domain.Tool{
		{
			Name:         "jurisprudence_search",
			Description:  "Busca semântica em decisões judiciais",
			Category:     "ai",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "query", Type: "string", Description: "Consulta de busca", Required: true},
				{Name: "court_types", Type: "array", Description: "Tipos de tribunal", 
					Required: false, Default: []string{"STF", "STJ", "TJ"}},
				{Name: "similarity_threshold", Type: "number", Description: "Limiar de similaridade (0.0-1.0)", 
					Required: false, Default: 0.7},
				{Name: "max_results", Type: "number", Description: "Número máximo de resultados", 
					Required: false, Default: 10},
				{Name: "date_range", Type: "object", Description: "Intervalo de datas", Required: false},
			},
			Handler: handleJurisprudenceSearch,
		},
		{
			Name:         "case_similarity_analysis",
			Description:  "Análise de similaridade entre casos",
			Category:     "ai",
			RequiredPlan: "business",
			Parameters: []domain.ToolParameter{
				{Name: "base_case", Type: "string", Description: "Descrição do caso base", Required: true},
				{Name: "comparison_cases", Type: "array", Description: "IDs ou descrições dos casos para comparar", 
					Required: true},
				{Name: "analysis_dimensions", Type: "array", Description: "Dimensões de análise", 
					Required: false, Default: []string{"semantic", "legal", "factual", "procedural"}},
			},
			Handler: handleCaseSimilarityAnalysis,
		},
		{
			Name:         "document_analysis",
			Description:  "Análise completa de documentos legais",
			Category:     "ai",
			RequiredPlan: "business",
			Parameters: []domain.ToolParameter{
				{Name: "document_text", Type: "string", Description: "Texto do documento", Required: true},
				{Name: "analysis_type", Type: "string", Description: "Tipo de análise", 
					Required: true, Enum: []string{"risk", "compliance", "entities", "classification"}},
				{Name: "legal_area", Type: "string", Description: "Área jurídica", Required: false},
			},
			Handler: handleDocumentAnalysis,
		},
		{
			Name:         "legal_document_generation",
			Description:  "Geração de documentos jurídicos",
			Category:     "ai",
			RequiredPlan: "business",
			Parameters: []domain.ToolParameter{
				{Name: "document_type", Type: "string", Description: "Tipo de documento", 
					Required: true, Enum: []string{"contract", "petition", "opinion", "motion"}},
				{Name: "template_variables", Type: "object", Description: "Variáveis do template", Required: true},
				{Name: "quality_level", Type: "string", Description: "Nível de qualidade", 
					Required: false, Default: "standard", Enum: []string{"draft", "standard", "professional", "premium"}},
				{Name: "customizations", Type: "object", Description: "Personalizações específicas", Required: false},
			},
			Handler: handleLegalDocumentGeneration,
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

// handleJurisprudenceSearch busca jurisprudências
func handleJurisprudenceSearch(params map[string]interface{}) (interface{}, error) {
	query := params["query"].(string)
	
	// TODO: Implementar integração real com AI Service
	
	result := map[string]interface{}{
		"query": query,
		"total_results": 47,
		"jurisprudences": []map[string]interface{}{
			{
				"id": "jur-001",
				"court": "STJ",
				"case_number": "REsp 1.234.567/SP",
				"similarity_score": 0.95,
				"date": "2024-03-15",
				"summary": "RESPONSABILIDADE CIVIL. DANO MORAL. CONFIGURAÇÃO. " +
					"A responsabilidade civil do médico é subjetiva, dependendo da comprovação de culpa...",
				"relevant_excerpt": "Nos termos da jurisprudência consolidada desta Corte, " +
					"a responsabilidade civil do médico...",
				"judge": "Min. Ricardo Villas Bôas Cueva",
			},
			{
				"id": "jur-002",
				"court": "TJ-SP",
				"case_number": "Apelação 1001234-56.2023.8.26.0100",
				"similarity_score": 0.92,
				"date": "2024-02-08",
				"summary": "INDENIZAÇÃO. ERRO MÉDICO. NEXO CAUSAL. " +
					"Em casos de responsabilidade médica, é necessário comprovar...",
				"relevant_excerpt": "O nexo de causalidade entre a conduta do profissional " +
					"e o dano alegado deve ser claramente demonstrado...",
				"judge": "Des. João Carlos Saletti",
			},
		},
		"search_metadata": map[string]interface{}{
			"execution_time_ms": 245,
			"index_version": "2024.06",
			"similarity_threshold": params["similarity_threshold"],
		},
	}
	
	return result, nil
}

// handleCaseSimilarityAnalysis analisa similaridade entre casos
func handleCaseSimilarityAnalysis(params map[string]interface{}) (interface{}, error) {
	baseCase := params["base_case"].(string)
	comparisonCases := params["comparison_cases"].([]interface{})
	
	// TODO: Implementar integração real com AI Service
	
	similarities := []map[string]interface{}{}
	for i, caseDesc := range comparisonCases {
		similarities = append(similarities, map[string]interface{}{
			"case_id": fmt.Sprintf("case-%d", i+1),
			"case_description": caseDesc,
			"overall_similarity": 0.85 - float64(i)*0.1,
			"dimensions": map[string]float64{
				"semantic":   0.90 - float64(i)*0.05,
				"legal":      0.85 - float64(i)*0.08,
				"factual":    0.80 - float64(i)*0.12,
				"procedural": 0.88 - float64(i)*0.07,
			},
			"key_similarities": []string{
				"Mesmo tipo de ação judicial",
				"Argumentação jurídica similar",
				"Contexto fático comparável",
			},
			"key_differences": []string{
				"Valor da causa diferente",
				"Jurisdição distinta",
				"Período temporal diverso",
			},
		})
	}
	
	result := map[string]interface{}{
		"base_case": baseCase,
		"analysis_date": time.Now().Format(time.RFC3339),
		"similarities": similarities,
		"recommendations": []string{
			"Considerar precedente do STJ no caso 1 como argumento principal",
			"Adaptar estratégia processual baseada no caso 2",
			"Evitar abordagem utilizada no caso 3 devido às diferenças contextuais",
		},
		"confidence_score": 0.88,
	}
	
	return result, nil
}

// handleDocumentAnalysis analisa documentos
func handleDocumentAnalysis(params map[string]interface{}) (interface{}, error) {
	documentText := params["document_text"].(string)
	analysisType := params["analysis_type"].(string)
	
	// TODO: Implementar integração real com AI Service
	
	var analysisResult map[string]interface{}
	
	switch analysisType {
	case "risk":
		analysisResult = map[string]interface{}{
			"risk_level": "medium",
			"risk_score": 0.65,
			"risk_factors": []map[string]interface{}{
				{
					"factor": "Cláusula de responsabilidade ambígua",
					"severity": "high",
					"location": "Parágrafo 5.2",
					"recommendation": "Reformular para maior clareza jurídica",
				},
				{
					"factor": "Prazo prescricional próximo",
					"severity": "medium",
					"location": "Cláusula 8",
					"recommendation": "Incluir salvaguardas processuais",
				},
			},
			"overall_assessment": "Documento apresenta riscos moderados que podem ser mitigados com ajustes pontuais",
		}
		
	case "entities":
		analysisResult = map[string]interface{}{
			"entities": []map[string]interface{}{
				{"type": "person", "name": "João Silva", "role": "autor", "document": "123.456.789-00"},
				{"type": "organization", "name": "Empresa XYZ Ltda", "role": "réu", "document": "12.345.678/0001-90"},
				{"type": "court", "name": "Tribunal de Justiça de São Paulo", "abbreviation": "TJ-SP"},
				{"type": "lawyer", "name": "Dra. Ana Costa", "oab": "OAB/SP 123456", "representing": "autor"},
			},
			"locations": []string{"São Paulo", "Campinas"},
			"dates": []string{"15/01/2024", "30/06/2024"},
			"monetary_values": []map[string]interface{}{
				{"value": 50000.00, "context": "valor da causa"},
				{"value": 15000.00, "context": "danos morais"},
			},
		}
		
	default:
		analysisResult = map[string]interface{}{
			"analysis_type": analysisType,
			"status": "completed",
			"summary": "Análise concluída com sucesso",
		}
	}
	
	result := map[string]interface{}{
		"document_length": len(documentText),
		"analysis_type": analysisType,
		"analysis_result": analysisResult,
		"processed_at": time.Now().Format(time.RFC3339),
		"confidence": 0.92,
	}
	
	return result, nil
}

// handleLegalDocumentGeneration gera documentos jurídicos
func handleLegalDocumentGeneration(params map[string]interface{}) (interface{}, error) {
	documentType := params["document_type"].(string)
	templateVars := params["template_variables"].(map[string]interface{})
	qualityLevel := "standard"
	if val, ok := params["quality_level"].(string); ok {
		qualityLevel = val
	}
	
	// TODO: Implementar integração real com AI Service
	
	var generatedDocument string
	var documentMetadata map[string]interface{}
	
	switch documentType {
	case "petition":
		generatedDocument = fmt.Sprintf(`EXCELENTÍSSIMO SENHOR DOUTOR JUIZ DE DIREITO DA %s VARA CÍVEL DA COMARCA DE %s

%s, brasileiro(a), %s, portador(a) do RG nº %s e CPF nº %s, residente e domiciliado(a) na %s, 
vem, respeitosamente, à presença de Vossa Excelência, por seu(sua) advogado(a) que esta subscreve, 
propor a presente

AÇÃO DE %s

em face de %s, pessoa jurídica de direito privado, inscrita no CNPJ sob nº %s, 
com sede na %s, pelos fatos e fundamentos que passa a expor:

I - DOS FATOS

%s

II - DO DIREITO

%s

III - DOS PEDIDOS

Ante o exposto, requer:

a) A citação do réu para, querendo, contestar a presente ação;
b) %s;
c) A condenação do réu ao pagamento das custas processuais e honorários advocatícios.

Dá-se à causa o valor de R$ %s.

Termos em que,
Pede deferimento.

%s, %s

_________________________
%s
OAB/%s %s`,
			getTemplateVar(templateVars, "vara", "1ª"),
			getTemplateVar(templateVars, "comarca", "São Paulo"),
			getTemplateVar(templateVars, "autor_nome", "NOME DO AUTOR"),
			getTemplateVar(templateVars, "autor_profissao", "profissão"),
			getTemplateVar(templateVars, "autor_rg", "XX.XXX.XXX-X"),
			getTemplateVar(templateVars, "autor_cpf", "XXX.XXX.XXX-XX"),
			getTemplateVar(templateVars, "autor_endereco", "endereço completo"),
			getTemplateVar(templateVars, "tipo_acao", "INDENIZAÇÃO POR DANOS MORAIS E MATERIAIS"),
			getTemplateVar(templateVars, "reu_nome", "NOME DO RÉU"),
			getTemplateVar(templateVars, "reu_cnpj", "XX.XXX.XXX/XXXX-XX"),
			getTemplateVar(templateVars, "reu_endereco", "endereço do réu"),
			getTemplateVar(templateVars, "fatos", "Descrição detalhada dos fatos..."),
			getTemplateVar(templateVars, "direito", "Fundamentação jurídica..."),
			getTemplateVar(templateVars, "pedidos", "Pedidos específicos"),
			getTemplateVar(templateVars, "valor_causa", "50.000,00"),
			getTemplateVar(templateVars, "cidade", "São Paulo"),
			time.Now().Format("02 de January de 2006"),
			getTemplateVar(templateVars, "advogado_nome", "NOME DO ADVOGADO"),
			getTemplateVar(templateVars, "advogado_uf", "SP"),
			getTemplateVar(templateVars, "advogado_oab", "XXXXX"),
		)
		
		documentMetadata = map[string]interface{}{
			"type": "Petição Inicial",
			"pages": 3,
			"sections": []string{"Qualificação", "Fatos", "Direito", "Pedidos"},
		}
		
	default:
		generatedDocument = "Documento gerado com base no template fornecido..."
		documentMetadata = map[string]interface{}{
			"type": documentType,
			"pages": 1,
		}
	}
	
	result := map[string]interface{}{
		"success": true,
		"document_type": documentType,
		"quality_level": qualityLevel,
		"document": generatedDocument,
		"metadata": documentMetadata,
		"generation_time_ms": 1250,
		"estimated_review_time": "15-20 minutos",
		"suggestions": []string{
			"Revisar valores e datas antes de protocolar",
			"Verificar documentação anexa necessária",
			"Conferir competência do juízo",
		},
	}
	
	return result, nil
}

// getTemplateVar obtém uma variável do template com valor padrão
func getTemplateVar(vars map[string]interface{}, key, defaultValue string) string {
	if val, ok := vars[key].(string); ok {
		return val
	}
	return defaultValue
}