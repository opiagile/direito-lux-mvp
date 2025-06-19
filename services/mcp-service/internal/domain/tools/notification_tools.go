package tools

import (
	"fmt"
	"time"

	"github.com/direito-lux/mcp-service/internal/domain"
)

// RegisterNotificationTools registra as ferramentas de notificação
func RegisterNotificationTools(registry *domain.ToolRegistry) error {
	tools := []*domain.Tool{
		{
			Name:         "notification_setup",
			Description:  "Configurar notificações personalizadas",
			Category:     "notification",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "trigger_type", Type: "string", Description: "Tipo de gatilho", Required: true,
					Enum: []string{"process_update", "deadline", "keyword_match", "court_decision"}},
				{Name: "conditions", Type: "object", Description: "Condições específicas", Required: true},
				{Name: "channels", Type: "array", Description: "Canais de notificação", 
					Required: true, Default: []string{"whatsapp", "email"}},
				{Name: "template_id", Type: "string", Description: "ID do template", Required: false},
				{Name: "priority", Type: "string", Description: "Prioridade da notificação", 
					Required: false, Default: "normal", 
					Enum: []string{"low", "normal", "high", "critical"}},
			},
			Handler: handleNotificationSetup,
		},
		{
			Name:         "bulk_notification",
			Description:  "Envio em massa de notificações",
			Category:     "notification",
			RequiredPlan: "business",
			Parameters: []domain.ToolParameter{
				{Name: "recipient_filter", Type: "object", Description: "Filtros de destinatários", Required: true},
				{Name: "message_template", Type: "string", Description: "Template da mensagem", Required: true},
				{Name: "variables", Type: "object", Description: "Variáveis do template", Required: false},
				{Name: "schedule", Type: "string", Description: "Agendamento (datetime)", Required: false},
				{Name: "channels", Type: "array", Description: "Canais de envio", 
					Required: true, Default: []string{"whatsapp", "email"}},
			},
			Handler: handleBulkNotification,
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

// handleNotificationSetup configura notificações
func handleNotificationSetup(params map[string]interface{}) (interface{}, error) {
	triggerType := params["trigger_type"].(string)
	conditions := params["conditions"].(map[string]interface{})
	channels := params["channels"].([]interface{})
	
	// TODO: Implementar integração real com Notification Service
	
	notificationID := fmt.Sprintf("notif-%d", time.Now().Unix())
	
	result := map[string]interface{}{
		"success": true,
		"notification_id": notificationID,
		"configuration": map[string]interface{}{
			"id": notificationID,
			"trigger_type": triggerType,
			"conditions": conditions,
			"channels": channels,
			"priority": params["priority"],
			"status": "active",
			"created_at": time.Now().Format(time.RFC3339),
		},
		"message": fmt.Sprintf("Notificação configurada com sucesso para %s", triggerType),
		"test_notification_sent": true,
		"estimated_triggers_per_day": calculateEstimatedTriggers(triggerType, conditions),
	}
	
	if templateID, ok := params["template_id"].(string); ok {
		result["template"] = map[string]interface{}{
			"id": templateID,
			"name": "Template Personalizado",
			"preview": "Olá {{nome}}, houve uma atualização no processo {{numero_processo}}...",
		}
	}
	
	return result, nil
}

// handleBulkNotification envia notificações em massa
func handleBulkNotification(params map[string]interface{}) (interface{}, error) {
	recipientFilter := params["recipient_filter"].(map[string]interface{})
	messageTemplate := params["message_template"].(string)
	channels := params["channels"].([]interface{})
	
	// TODO: Implementar integração real com Notification Service
	
	// Simular contagem de destinatários
	recipientCount := 0
	if filter, ok := recipientFilter["type"].(string); ok {
		switch filter {
		case "all_clients":
			recipientCount = 150
		case "active_processes":
			recipientCount = 87
		case "specific_tag":
			recipientCount = 32
		default:
			recipientCount = 25
		}
	}
	
	bulkID := fmt.Sprintf("bulk-%d", time.Now().Unix())
	
	result := map[string]interface{}{
		"success": true,
		"bulk_notification_id": bulkID,
		"status": "processing",
		"summary": map[string]interface{}{
			"total_recipients": recipientCount,
			"channels": channels,
			"estimated_delivery_time": "15-30 minutos",
		},
		"message_preview": truncateMessage(messageTemplate, 100),
		"created_at": time.Now().Format(time.RFC3339),
	}
	
	// Se agendado
	if schedule, ok := params["schedule"].(string); ok {
		result["scheduled_for"] = schedule
		result["status"] = "scheduled"
	} else {
		// Simulação de envio imediato
		result["progress"] = map[string]interface{}{
			"sent": 0,
			"pending": recipientCount,
			"failed": 0,
			"status_url": fmt.Sprintf("/api/v1/notifications/bulk/%s/status", bulkID),
		}
	}
	
	// Breakdown por canal
	channelBreakdown := make(map[string]int)
	for _, channel := range channels {
		channelStr := channel.(string)
		channelBreakdown[channelStr] = recipientCount
	}
	result["channel_breakdown"] = channelBreakdown
	
	return result, nil
}

// calculateEstimatedTriggers calcula estimativa de triggers
func calculateEstimatedTriggers(triggerType string, conditions map[string]interface{}) int {
	// Simulação baseada no tipo de trigger
	switch triggerType {
	case "process_update":
		return 15 // média de 15 atualizações por dia
	case "deadline":
		return 3 // média de 3 prazos por dia
	case "keyword_match":
		if keywords, ok := conditions["keywords"].([]interface{}); ok {
			return len(keywords) * 2 // 2 matches por palavra-chave por dia
		}
		return 5
	case "court_decision":
		return 8 // média de 8 decisões por dia
	default:
		return 10
	}
}

// truncateMessage trunca mensagem para preview
func truncateMessage(message string, maxLength int) string {
	if len(message) <= maxLength {
		return message
	}
	return message[:maxLength-3] + "..."
}