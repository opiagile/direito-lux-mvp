package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/domain"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http/dto"
	"github.com/direito-lux/mcp-service/internal/infrastructure/events"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
)

// ToolService serviço para gerenciar ferramentas MCP
type ToolService struct {
	logger       *zap.Logger
	metrics      *metrics.Metrics
	eventBus     *events.EventBus
	toolRegistry *domain.ToolRegistry
	executions   map[string]*domain.ToolExecution // Simulação em memória
}

// NewToolService cria nova instância do serviço
func NewToolService(
	logger *zap.Logger,
	metrics *metrics.Metrics,
	eventBus *events.EventBus,
	toolRegistry *domain.ToolRegistry,
) *ToolService {
	return &ToolService{
		logger:       logger,
		metrics:      metrics,
		eventBus:     eventBus,
		toolRegistry: toolRegistry,
		executions:   make(map[string]*domain.ToolExecution),
	}
}

// ListAvailableTools lista ferramentas disponíveis
func (s *ToolService) ListAvailableTools(ctx context.Context, tenantID string) (*dto.ListToolsResponse, error) {
	// Obter ferramentas do registry
	tools := s.toolRegistry.GetAllTools()
	
	// Mapear para DTOs
	toolDTOs := make([]dto.ToolDefinition, 0, len(tools))
	categories := make(map[string]bool)
	
	for _, tool := range tools {
		// Mapear parâmetros
		properties := make(map[string]dto.ParameterProperty)
		required := make([]string, 0)
		
		for _, param := range tool.Parameters {
			properties[param.Name] = dto.ParameterProperty{
				Type:        param.Type,
				Description: param.Description,
				Enum:        param.Enum,
				Default:     param.Default,
			}
			if param.Required {
				required = append(required, param.Name)
			}
		}
		
		toolDTO := dto.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Category:    tool.Category,
			Parameters: dto.ToolParameterSchema{
				Type:       "object",
				Properties: properties,
				Required:   required,
			},
			Permissions: []string{}, // TODO: implementar sistema de permissões
			QuotaCost:   1, // Valor padrão
			Timeout:     30 * time.Second, // Valor padrão
			Async:       false, // Valor padrão
			Examples:    []dto.ToolExample{}, // TODO: implementar exemplos
		}
		
		toolDTOs = append(toolDTOs, toolDTO)
		categories[tool.Category] = true
	}
	
	// Extrair categorias únicas
	categoryList := make([]string, 0, len(categories))
	for category := range categories {
		categoryList = append(categoryList, category)
	}
	
	s.logger.Debug("Ferramentas listadas",
		zap.String("tenant_id", tenantID),
		zap.Int("count", len(toolDTOs)),
	)
	
	return &dto.ListToolsResponse{
		Tools:    toolDTOs,
		Count:    len(toolDTOs),
		Category: categoryList,
	}, nil
}

// ExecuteTool executa uma ferramenta
func (s *ToolService) ExecuteTool(ctx context.Context, req dto.ExecuteToolRequest) (*dto.ToolExecutionResponse, error) {
	// Validar sessão (simulação)
	// TODO: Integrar com SessionService para validar sessão
	
	// Obter ferramenta do registry
	tool, err := s.toolRegistry.GetTool(req.ToolName)
	if err != nil {
		return nil, fmt.Errorf("ferramenta não encontrada: %s", req.ToolName)
	}
	
	// Converter SessionID string para UUID
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}
	
	// Criar execução usando o construtor do domain
	execution := domain.NewToolExecution(sessionID, req.ToolName, req.Parameters)
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			execution.SetMetadata(key, value)
		}
	}
	
	// Salvar execução
	executionID := execution.ID.String()
	s.executions[executionID] = execution
	
	// Publicar evento de início
	event := domain.ToolEvent{
		Type:      "tool_execution_started",
		SessionID: req.SessionID,
		ToolName:  req.ToolName,
		UserID:    "", // TODO: obter do contexto da sessão
		TenantID:  "", // TODO: obter do contexto da sessão
		Timestamp: time.Now(),
		Success:   false, // Ainda não executado
		Data: map[string]interface{}{
			"execution_id": executionID,
			"parameters":   req.Parameters,
			"async":        req.Async,
		},
	}
	
	if err := s.eventBus.PublishToolEvent(ctx, event); err != nil {
		s.logger.Warn("Erro ao publicar evento de ferramenta", zap.Error(err))
	}
	
	// Se não for asíncrono, executar imediatamente
	if !req.Async {
		return s.executeToolSync(ctx, execution, tool)
	}
	
	// Se for asíncrono, executar em background
	go s.executeToolAsync(ctx, execution, tool)
	
	// Retornar resposta imediata
	return s.mapExecutionToResponse(execution), nil
}

// GetToolExecution obtém status de execução
func (s *ToolService) GetToolExecution(ctx context.Context, executionID string) (*dto.ToolExecutionResponse, error) {
	execution, exists := s.executions[executionID]
	if !exists {
		return nil, fmt.Errorf("execução não encontrada: %s", executionID)
	}
	
	return s.mapExecutionToResponse(execution), nil
}

// ListToolExecutions lista execuções de ferramentas
func (s *ToolService) ListToolExecutions(ctx context.Context, sessionID string, page, pageSize int) (*dto.ToolExecutionsListResponse, error) {
	// Filtrar execuções por sessão
	filteredExecutions := make([]*domain.ToolExecution, 0)
	for _, execution := range s.executions {
		if sessionID == "" || execution.SessionID.String() == sessionID {
			filteredExecutions = append(filteredExecutions, execution)
		}
	}
	
	// Simular paginação
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(filteredExecutions) {
		end = len(filteredExecutions)
	}
	
	paginatedExecutions := filteredExecutions[start:end]
	
	// Mapear para DTOs
	executionDTOs := make([]dto.ToolExecutionResponse, 0, len(paginatedExecutions))
	for _, execution := range paginatedExecutions {
		executionDTOs = append(executionDTOs, *s.mapExecutionToResponse(execution))
	}
	
	return &dto.ToolExecutionsListResponse{
		Executions: executionDTOs,
		TotalCount: len(filteredExecutions),
		Page:       page,
		PageSize:   pageSize,
		HasMore:    end < len(filteredExecutions),
	}, nil
}

// executeToolSync executa ferramenta sincronamente
func (s *ToolService) executeToolSync(ctx context.Context, execution *domain.ToolExecution, tool *domain.Tool) (*dto.ToolExecutionResponse, error) {
	start := time.Now()
	execution.Status = "running"
	
	// TODO: Implementar execução real da ferramenta
	// Por agora, simular execução
	time.Sleep(100 * time.Millisecond) // Simular processamento
	
	// Simular resultado baseado no nome da ferramenta
	result := s.simulateToolExecution(tool.Name, execution.Parameters)
	
	execution.Status = "completed"
	execution.Result = result
	completedAt := time.Now()
	execution.CompletedAt = &completedAt
	execution.Duration = time.Since(start)
	
	// Registrar métricas
	if s.metrics != nil {
		// TODO: Extrair tenant_id do contexto
		s.metrics.RecordMCPToolExecution(tool.Name, "unknown", "unknown", execution.Duration, nil)
	}
	
	// Publicar evento de conclusão
	event := domain.ToolEvent{
		Type:      "tool_execution_completed",
		SessionID: execution.SessionID.String(),
		ToolName:  execution.ToolName,
		UserID:    "", // TODO: obter do contexto da sessão
		TenantID:  "", // TODO: obter do contexto da sessão
		Timestamp: time.Now(),
		Success:   true,
		Data: map[string]interface{}{
			"execution_id": execution.ID.String(),
			"result":       result,
			"duration":     execution.Duration,
			"status":       "completed",
		},
	}
	
	if err := s.eventBus.PublishToolEvent(ctx, event); err != nil {
		s.logger.Warn("Erro ao publicar evento de conclusão", zap.Error(err))
	}
	
	s.logger.Info("Ferramenta executada",
		zap.String("execution_id", execution.ID.String()),
		zap.String("tool_name", tool.Name),
		zap.Duration("duration", execution.Duration),
	)
	
	return s.mapExecutionToResponse(execution), nil
}

// executeToolAsync executa ferramenta assincronamente
func (s *ToolService) executeToolAsync(ctx context.Context, execution *domain.ToolExecution, tool *domain.Tool) {
	// Executar em background usando executeToolSync
	_, err := s.executeToolSync(ctx, execution, tool)
	if err != nil {
		execution.Status = "failed"
		execution.Error = err.Error()
		completedAt := time.Now()
		execution.CompletedAt = &completedAt
		
		s.logger.Error("Erro na execução assíncrona",
			zap.String("execution_id", execution.ID.String()),
			zap.Error(err),
		)
	}
}

// simulateToolExecution simula execução de ferramenta
func (s *ToolService) simulateToolExecution(toolName string, parameters map[string]interface{}) interface{} {
	switch toolName {
	case "process_search":
		return map[string]interface{}{
			"processes": []map[string]interface{}{
				{
					"id":     "123",
					"number": "0001234-56.2024.8.02.0001",
					"client": "João Silva",
					"status": "active",
				},
			},
			"total": 1,
		}
	case "document_analyze":
		return map[string]interface{}{
			"analysis": "Documento analisado com sucesso. Encontradas 3 cláusulas importantes.",
			"clauses": []string{
				"Cláusula de pagamento",
				"Cláusula de rescisão",
				"Cláusula de vigência",
			},
			"confidence": 0.95,
		}
	default:
		return map[string]interface{}{
			"message": fmt.Sprintf("Ferramenta %s executada com sucesso", toolName),
			"parameters": parameters,
			"timestamp": time.Now().UTC(),
		}
	}
}

// mapExecutionToResponse mapeia execução para DTO
func (s *ToolService) mapExecutionToResponse(execution *domain.ToolExecution) *dto.ToolExecutionResponse {
	response := &dto.ToolExecutionResponse{
		ID:         execution.ID.String(),
		SessionID:  execution.SessionID.String(),
		ToolName:   execution.ToolName,
		Status:     string(execution.Status),
		Result:     execution.Result,
		Error:      execution.Error,
		StartedAt:  execution.StartedAt,
		Parameters: execution.Parameters,
		Metadata:   execution.Metadata,
	}
	
	if execution.CompletedAt != nil {
		response.CompletedAt = execution.CompletedAt
		response.Duration = execution.Duration
	}
	
	return response
}

// mapParameters mapeia parâmetros da ferramenta
func (s *ToolService) mapParameters(params map[string]domain.ParameterDefinition) map[string]dto.ParameterProperty {
	result := make(map[string]dto.ParameterProperty)
	
	for name, param := range params {
		result[name] = dto.ParameterProperty{
			Type:        param.Type,
			Description: param.Description,
			Enum:        param.Enum,
			Default:     param.Default,
			Minimum:     param.Minimum,
			Maximum:     param.Maximum,
			Pattern:     param.Pattern,
		}
	}
	
	return result
}

// mapExamples mapeia exemplos da ferramenta
func (s *ToolService) mapExamples(examples []domain.ToolExample) []dto.ToolExample {
	result := make([]dto.ToolExample, 0, len(examples))
	
	for _, example := range examples {
		result = append(result, dto.ToolExample{
			Description:    example.Description,
			Parameters:     example.Parameters,
			ExpectedResult: example.ExpectedResult,
		})
	}
	
	return result
}