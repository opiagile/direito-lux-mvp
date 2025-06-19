package domain

import (
	"fmt"
	"strings"
)

// Tool representa uma ferramenta MCP disponível
type Tool struct {
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	Category     string                       `json:"category"`
	Parameters   []ToolParameter              `json:"parameters"`
	RequiredPlan string                       `json:"required_plan"` // starter, professional, business, enterprise
	Handler      func(params map[string]interface{}) (interface{}, error) `json:"-"`
}

// ToolParameter representa um parâmetro de uma ferramenta
type ToolParameter struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` // string, number, boolean, array, object
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolRegistry gerencia o registro de ferramentas MCP
type ToolRegistry struct {
	tools       map[string]*Tool
	categories  map[string][]string
}

// NewToolRegistry cria um novo registro de ferramentas
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:      make(map[string]*Tool),
		categories: make(map[string][]string),
	}
}

// RegisterTool registra uma nova ferramenta
func (r *ToolRegistry) RegisterTool(tool *Tool) error {
	if tool.Name == "" {
		return fmt.Errorf("nome da ferramenta é obrigatório")
	}
	
	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("ferramenta %s já está registrada", tool.Name)
	}
	
	r.tools[tool.Name] = tool
	
	// Adicionar à categoria
	if tool.Category != "" {
		r.categories[tool.Category] = append(r.categories[tool.Category], tool.Name)
	}
	
	return nil
}

// GetTool obtém uma ferramenta pelo nome
func (r *ToolRegistry) GetTool(name string) (*Tool, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("ferramenta %s não encontrada", name)
	}
	return tool, nil
}

// GetToolsByCategory obtém todas as ferramentas de uma categoria
func (r *ToolRegistry) GetToolsByCategory(category string) []*Tool {
	toolNames, exists := r.categories[category]
	if !exists {
		return []*Tool{}
	}
	
	tools := make([]*Tool, 0, len(toolNames))
	for _, name := range toolNames {
		if tool, exists := r.tools[name]; exists {
			tools = append(tools, tool)
		}
	}
	
	return tools
}

// GetToolsByPlan obtém todas as ferramentas disponíveis para um plano
func (r *ToolRegistry) GetToolsByPlan(plan string) []*Tool {
	planHierarchy := map[string]int{
		"starter":      1,
		"professional": 2,
		"business":     3,
		"enterprise":   4,
	}
	
	userPlanLevel, exists := planHierarchy[strings.ToLower(plan)]
	if !exists {
		return []*Tool{}
	}
	
	tools := make([]*Tool, 0)
	for _, tool := range r.tools {
		toolPlanLevel, exists := planHierarchy[strings.ToLower(tool.RequiredPlan)]
		if !exists {
			continue
		}
		
		if userPlanLevel >= toolPlanLevel {
			tools = append(tools, tool)
		}
	}
	
	return tools
}

// GetAllTools obtém todas as ferramentas registradas
func (r *ToolRegistry) GetAllTools() []*Tool {
	tools := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetCategories obtém todas as categorias disponíveis
func (r *ToolRegistry) GetCategories() []string {
	categories := make([]string, 0, len(r.categories))
	for category := range r.categories {
		categories = append(categories, category)
	}
	return categories
}

// ValidateParameters valida os parâmetros de uma ferramenta
func (t *Tool) ValidateParameters(params map[string]interface{}) error {
	// Verificar parâmetros obrigatórios
	for _, param := range t.Parameters {
		if param.Required {
			if _, exists := params[param.Name]; !exists {
				return fmt.Errorf("parâmetro obrigatório '%s' não fornecido", param.Name)
			}
		}
	}
	
	// Validar tipos e enums
	for _, param := range t.Parameters {
		value, exists := params[param.Name]
		if !exists {
			continue
		}
		
		// Validar enum se existir
		if len(param.Enum) > 0 {
			strValue := fmt.Sprintf("%v", value)
			found := false
			for _, enum := range param.Enum {
				if strValue == enum {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("valor '%v' inválido para parâmetro '%s'. Valores permitidos: %v", 
					value, param.Name, param.Enum)
			}
		}
		
		// TODO: Adicionar validação de tipos mais robusta
	}
	
	return nil
}

// Execute executa a ferramenta com os parâmetros fornecidos
func (t *Tool) Execute(params map[string]interface{}) (interface{}, error) {
	// Validar parâmetros
	if err := t.ValidateParameters(params); err != nil {
		return nil, err
	}
	
	// Aplicar valores default
	for _, param := range t.Parameters {
		if _, exists := params[param.Name]; !exists && param.Default != nil {
			params[param.Name] = param.Default
		}
	}
	
	// Executar handler
	if t.Handler == nil {
		return nil, fmt.Errorf("handler não implementado para ferramenta %s", t.Name)
	}
	
	return t.Handler(params)
}