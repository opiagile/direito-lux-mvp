package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/direito-lux/mcp-service/internal/infrastructure/http/dto"
	"github.com/direito-lux/mcp-service/internal/infrastructure/logging"
)

// ListAvailableTools lista ferramentas disponíveis
func ListAvailableTools() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Buscar ferramentas do ToolService
		// Por agora, simular lista de ferramentas
		tools := []dto.ToolDefinition{
			{
				Name:        "process_search",
				Description: "Buscar processos por diversos critérios",
				Category:    "process_management",
				Parameters: dto.ToolParameterSchema{
					Type: "object",
					Properties: map[string]dto.ParameterProperty{
						"client_name": {
							Type:        "string",
							Description: "Nome do cliente",
						},
						"process_number": {
							Type:        "string",
							Description: "Número do processo",
						},
						"status": {
							Type:        "string",
							Description: "Status do processo",
							Enum:        []string{"active", "archived", "suspended"},
						},
					},
				},
				QuotaCost: 1,
				Timeout:   30 * time.Second,
				Async:     false,
			},
			{
				Name:        "document_analyze",
				Description: "Analisar documentos com IA",
				Category:    "ai_analysis",
				Parameters: dto.ToolParameterSchema{
					Type: "object",
					Properties: map[string]dto.ParameterProperty{
						"document_url": {
							Type:        "string",
							Description: "URL do documento",
						},
						"analysis_type": {
							Type:        "string",
							Description: "Tipo de análise",
							Enum:        []string{"summary", "clauses", "risks", "compliance"},
						},
					},
					Required: []string{"document_url", "analysis_type"},
				},
				QuotaCost: 3,
				Timeout:   60 * time.Second,
				Async:     true,
			},
			{
				Name:        "jurisprudence_search",
				Description: "Buscar jurisprudência relevante",
				Category:    "ai_analysis",
				Parameters: dto.ToolParameterSchema{
					Type: "object",
					Properties: map[string]dto.ParameterProperty{
						"query": {
							Type:        "string",
							Description: "Consulta de busca",
						},
						"court": {
							Type:        "string",
							Description: "Tribunal específico",
						},
						"limit": {
							Type:        "integer",
							Description: "Limite de resultados",
							Minimum:     &[]float64{1}[0],
							Maximum:     &[]float64{50}[0],
						},
					},
					Required: []string{"query"},
				},
				QuotaCost: 2,
				Timeout:   45 * time.Second,
				Async:     false,
			},
		}

		response := &dto.ListToolsResponse{
			Tools:    tools,
			Count:    len(tools),
			Category: []string{"process_management", "ai_analysis"},
		}

		c.JSON(http.StatusOK, response)
	}
}

// ExecuteTool executa uma ferramenta
func ExecuteTool() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ExecuteToolRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Validar se a sessão existe (simulação)
		if req.SessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		// TODO: Executar ferramenta via ToolService
		// Por agora, simular execução
		executionID := uuid.New().String()
		
		// Simular execução baseada no nome da ferramenta
		result := simulateToolExecution(req.ToolName, req.Parameters)
		
		response := &dto.ToolExecutionResponse{
			ID:          executionID,
			SessionID:   req.SessionID,
			ToolName:    req.ToolName,
			Status:      "completed",
			Result:      result,
			StartedAt:   time.Now(),
			CompletedAt: &[]time.Time{time.Now()}[0],
			Duration:    150 * time.Millisecond,
			Parameters:  req.Parameters,
			Metadata:    req.Metadata,
		}

		logging.LogInfo(c.Request.Context(), nil, "Ferramenta executada",
			logging.String("execution_id", executionID),
			logging.String("tool_name", req.ToolName),
			logging.String("session_id", req.SessionID),
		)

		c.JSON(http.StatusOK, response)
	}
}

// GetToolExecution obtém resultado de execução
func GetToolExecution() gin.HandlerFunc {
	return func(c *gin.Context) {
		executionID := c.Param("id")
		if executionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Execution ID is required",
			})
			return
		}

		// TODO: Buscar execução real do banco
		// Por agora, simular resposta
		response := &dto.ToolExecutionResponse{
			ID:        executionID,
			SessionID: "session123",
			ToolName:  "process_search",
			Status:    "completed",
			Result: map[string]interface{}{
				"processes": []map[string]interface{}{
					{
						"id":     "123",
						"number": "0001234-56.2024.8.02.0001",
						"client": "João Silva",
						"status": "active",
					},
				},
				"total": 1,
			},
			StartedAt:   time.Now().Add(-2 * time.Minute),
			CompletedAt: &[]time.Time{time.Now().Add(-1 * time.Minute)}[0],
			Duration:    1 * time.Minute,
			Parameters: map[string]interface{}{
				"client_name": "João Silva",
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// ListToolExecutions lista execuções de ferramentas
func ListToolExecutions() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parâmetros de busca
		sessionID := c.Query("session_id")
		page := 1
		pageSize := 20

		if p := c.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}

		if ps := c.Query("page_size"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil {
				pageSize = parsed
			}
		}

		// TODO: Buscar execuções reais do banco
		// Por agora, simular lista
		executions := []dto.ToolExecutionResponse{
			{
				ID:        "exec1",
				SessionID: sessionID,
				ToolName:  "process_search",
				Status:    "completed",
				StartedAt: time.Now().Add(-10 * time.Minute),
				CompletedAt: &[]time.Time{time.Now().Add(-9 * time.Minute)}[0],
				Duration:    1 * time.Minute,
			},
			{
				ID:        "exec2",
				SessionID: sessionID,
				ToolName:  "document_analyze",
				Status:    "running",
				StartedAt: time.Now().Add(-2 * time.Minute),
			},
		}

		response := &dto.ToolExecutionsListResponse{
			Executions: executions,
			TotalCount: len(executions),
			Page:       page,
			PageSize:   pageSize,
			HasMore:    false,
		}

		c.JSON(http.StatusOK, response)
	}
}

// simulateToolExecution simula execução de ferramenta
func simulateToolExecution(toolName string, parameters map[string]interface{}) interface{} {
	switch toolName {
	case "process_search":
		return map[string]interface{}{
			"processes": []map[string]interface{}{
				{
					"id":          "123",
					"number":      "0001234-56.2024.8.02.0001",
					"client":      "João Silva",
					"status":      "active",
					"court":       "TRT 2ª Região",
					"created_at":  time.Now().Add(-30 * 24 * time.Hour),
					"description": "Ação trabalhista por horas extras",
				},
				{
					"id":          "124",
					"number":      "0001235-57.2024.8.02.0001",
					"client":      "Maria Santos",
					"status":      "active",
					"court":       "TRT 2ª Região",
					"created_at":  time.Now().Add(-15 * 24 * time.Hour),
					"description": "Rescisão indireta",
				},
			},
			"total": 2,
			"page":  1,
		}
	case "document_analyze":
		return map[string]interface{}{
			"analysis": "Documento analisado com sucesso. Encontradas 5 cláusulas importantes e 2 pontos de atenção.",
			"clauses": []map[string]interface{}{
				{
					"type":        "payment",
					"description": "Cláusula de pagamento em 30 dias",
					"risk_level":  "low",
				},
				{
					"type":        "termination",
					"description": "Cláusula de rescisão unilateral",
					"risk_level":  "medium",
				},
				{
					"type":        "penalty",
					"description": "Multa por descumprimento",
					"risk_level":  "high",
				},
			},
			"confidence": 0.92,
			"warnings": []string{
				"Cláusula de multa pode ser considerada abusiva",
				"Falta especificação de juros por atraso",
			},
		}
	case "jurisprudence_search":
		return map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"court":     "STJ",
					"number":    "REsp 1.234.567/SP",
					"date":      "2024-01-15",
					"summary":   "Horas extras. Cálculo. Acordo coletivo prevalece sobre lei",
					"relevance": 0.95,
				},
				{
					"court":     "TST",
					"number":    "RR 2.345.678/RJ",
					"date":      "2024-02-20",
					"summary":   "Adicional noturno. Base de cálculo. Horas extras",
					"relevance": 0.88,
				},
			},
			"total": 2,
			"query_analyzed": parameters["query"],
		}
	default:
		return map[string]interface{}{
			"message":    "Ferramenta executada com sucesso",
			"tool_name":  toolName,
			"parameters": parameters,
			"timestamp":  time.Now().UTC(),
			"status":     "success",
		}
	}
}