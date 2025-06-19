package tools

import (
	"fmt"
	"time"

	"github.com/direito-lux/mcp-service/internal/domain"
)

// RegisterProcessTools registra as ferramentas de processos
func RegisterProcessTools(registry *domain.ToolRegistry) error {
	tools := []*domain.Tool{
		{
			Name:         "process_search",
			Description:  "Buscar processos por diversos critérios",
			Category:     "process",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "client_name", Type: "string", Description: "Nome do cliente", Required: false},
				{Name: "process_number", Type: "string", Description: "Número do processo", Required: false},
				{Name: "status", Type: "string", Description: "Status do processo", Required: false, 
					Enum: []string{"active", "archived", "suspended"}},
				{Name: "date_range", Type: "object", Description: "Intervalo de datas", Required: false},
				{Name: "court", Type: "string", Description: "Tribunal", Required: false},
				{Name: "lawyer", Type: "string", Description: "Advogado responsável", Required: false},
			},
			Handler: handleProcessSearch,
		},
		{
			Name:         "process_monitor",
			Description:  "Configurar monitoramento automático de processos",
			Category:     "process",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "process_id", Type: "string", Description: "ID do processo", Required: true},
				{Name: "notification_types", Type: "array", Description: "Tipos de notificação", 
					Required: true, Default: []string{"movement", "deadline"}},
				{Name: "channels", Type: "array", Description: "Canais de notificação", 
					Required: true, Default: []string{"whatsapp", "email"}},
				{Name: "frequency", Type: "string", Description: "Frequência de verificação", 
					Required: false, Default: "immediate", Enum: []string{"immediate", "daily", "weekly"}},
			},
			Handler: handleProcessMonitor,
		},
		{
			Name:         "process_create",
			Description:  "Adicionar novo processo ao monitoramento",
			Category:     "process",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "cnj_number", Type: "string", Description: "Número CNJ do processo", Required: true},
				{Name: "client_name", Type: "string", Description: "Nome do cliente", Required: true},
				{Name: "process_type", Type: "string", Description: "Tipo do processo", Required: false},
				{Name: "responsible_lawyer", Type: "string", Description: "Advogado responsável", Required: false},
				{Name: "tags", Type: "array", Description: "Tags do processo", Required: false},
			},
			Handler: handleProcessCreate,
		},
		{
			Name:         "process_details",
			Description:  "Obter detalhes completos de um processo",
			Category:     "process",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "process_id", Type: "string", Description: "ID ou número do processo", Required: true},
				{Name: "include_movements", Type: "boolean", Description: "Incluir movimentações", 
					Required: false, Default: true},
				{Name: "include_parties", Type: "boolean", Description: "Incluir partes", 
					Required: false, Default: true},
			},
			Handler: handleProcessDetails,
		},
		{
			Name:         "process_update_status",
			Description:  "Atualizar status de um processo",
			Category:     "process",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "process_id", Type: "string", Description: "ID do processo", Required: true},
				{Name: "status", Type: "string", Description: "Novo status", Required: true,
					Enum: []string{"active", "archived", "suspended"}},
				{Name: "reason", Type: "string", Description: "Motivo da mudança", Required: false},
			},
			Handler: handleProcessUpdateStatus,
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

// handleProcessSearch busca processos
func handleProcessSearch(params map[string]interface{}) (interface{}, error) {
	// TODO: Implementar integração real com Process Service
	// Por enquanto, retornar dados mock
	
	result := map[string]interface{}{
		"total": 15,
		"processes": []map[string]interface{}{
			{
				"id":            "proc-123",
				"number":        "1001234-56.2024.8.26.0100",
				"client_name":   "João Silva",
				"status":        "active",
				"court":         "TJ-SP",
				"last_movement": "Petição protocolada",
				"updated_at":    time.Now().Add(-2 * 24 * time.Hour).Format(time.RFC3339),
			},
			{
				"id":            "proc-124",
				"number":        "2001234-78.2024.8.26.0200",
				"client_name":   "Maria Santos",
				"status":        "active",
				"court":         "TJ-SP",
				"last_movement": "Audiência agendada",
				"updated_at":    time.Now().Add(-5 * 24 * time.Hour).Format(time.RFC3339),
			},
		},
		"filters_applied": params,
	}
	
	return result, nil
}

// handleProcessMonitor configura monitoramento
func handleProcessMonitor(params map[string]interface{}) (interface{}, error) {
	processID := params["process_id"].(string)
	
	// TODO: Implementar integração real com Process Service
	
	result := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Monitoramento configurado para processo %s", processID),
		"monitoring_config": map[string]interface{}{
			"process_id":         processID,
			"notification_types": params["notification_types"],
			"channels":          params["channels"],
			"frequency":         params["frequency"],
			"created_at":        time.Now().Format(time.RFC3339),
		},
	}
	
	return result, nil
}

// handleProcessCreate cria novo processo
func handleProcessCreate(params map[string]interface{}) (interface{}, error) {
	cnjNumber := params["cnj_number"].(string)
	clientName := params["client_name"].(string)
	
	// TODO: Implementar integração real com Process Service
	
	result := map[string]interface{}{
		"success": true,
		"process": map[string]interface{}{
			"id":                 "proc-" + time.Now().Format("20060102150405"),
			"cnj_number":        cnjNumber,
			"client_name":       clientName,
			"process_type":      params["process_type"],
			"responsible_lawyer": params["responsible_lawyer"],
			"tags":              params["tags"],
			"status":            "active",
			"created_at":        time.Now().Format(time.RFC3339),
		},
		"message": fmt.Sprintf("Processo %s adicionado com sucesso", cnjNumber),
	}
	
	return result, nil
}

// handleProcessDetails obtém detalhes do processo
func handleProcessDetails(params map[string]interface{}) (interface{}, error) {
	processID := params["process_id"].(string)
	includeMovements := true
	includeParties := true
	
	if val, ok := params["include_movements"].(bool); ok {
		includeMovements = val
	}
	if val, ok := params["include_parties"].(bool); ok {
		includeParties = val
	}
	
	// TODO: Implementar integração real com Process Service
	
	result := map[string]interface{}{
		"process": map[string]interface{}{
			"id":              processID,
			"number":          "1001234-56.2024.8.26.0100",
			"client_name":     "João Silva",
			"status":          "active",
			"court":           "TJ-SP",
			"judge":           "Dr. Carlos Pereira",
			"value":           150000.00,
			"subject":         "Indenização por danos morais",
			"distribution_date": "2024-01-15",
		},
	}
	
	if includeMovements {
		result["movements"] = []map[string]interface{}{
			{
				"id":          "mov-001",
				"date":        "2024-06-15",
				"description": "Petição protocolada",
				"content":     "Petição inicial protocolada com pedido de tutela de urgência",
			},
			{
				"id":          "mov-002",
				"date":        "2024-06-10",
				"description": "Decisão judicial",
				"content":     "Tutela de urgência deferida parcialmente",
			},
		}
	}
	
	if includeParties {
		result["parties"] = []map[string]interface{}{
			{
				"type":     "author",
				"name":     "João Silva",
				"document": "123.456.789-00",
				"lawyer":   "Dra. Ana Costa - OAB/SP 123456",
			},
			{
				"type":     "defendant",
				"name":     "Empresa XYZ Ltda",
				"document": "12.345.678/0001-90",
				"lawyer":   "Dr. Pedro Santos - OAB/SP 654321",
			},
		}
	}
	
	return result, nil
}

// handleProcessUpdateStatus atualiza status do processo
func handleProcessUpdateStatus(params map[string]interface{}) (interface{}, error) {
	processID := params["process_id"].(string)
	newStatus := params["status"].(string)
	
	// TODO: Implementar integração real com Process Service
	
	result := map[string]interface{}{
		"success": true,
		"process_id": processID,
		"old_status": "active",
		"new_status": newStatus,
		"updated_at": time.Now().Format(time.RFC3339),
		"message": fmt.Sprintf("Status do processo %s atualizado para %s", processID, newStatus),
	}
	
	if reason, ok := params["reason"].(string); ok && reason != "" {
		result["reason"] = reason
	}
	
	return result, nil
}