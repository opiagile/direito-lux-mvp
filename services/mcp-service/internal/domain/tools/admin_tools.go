package tools

import (
	"fmt"
	"time"

	"github.com/direito-lux/mcp-service/internal/domain"
)

// RegisterAdminTools registra as ferramentas administrativas
func RegisterAdminTools(registry *domain.ToolRegistry) error {
	tools := []*domain.Tool{
		{
			Name:         "user_management",
			Description:  "Gerenciamento de usuários e permissões",
			Category:     "admin",
			RequiredPlan: "business",
			Parameters: []domain.ToolParameter{
				{Name: "action", Type: "string", Description: "Ação a executar", Required: true,
					Enum: []string{"list", "create", "update", "deactivate", "permissions"}},
				{Name: "user_data", Type: "object", Description: "Dados do usuário", 
					Required: false},
				{Name: "filters", Type: "object", Description: "Filtros para listagem", 
					Required: false},
				{Name: "permissions", Type: "array", Description: "Permissões do usuário", 
					Required: false},
			},
			Handler: handleUserManagement,
		},
		{
			Name:         "tenant_analytics",
			Description:  "Análise de uso por escritório/tenant",
			Category:     "admin",
			RequiredPlan: "business",
			Parameters: []domain.ToolParameter{
				{Name: "tenant_id", Type: "string", Description: "ID do tenant (admin only)", 
					Required: false},
				{Name: "metrics", Type: "array", Description: "Métricas desejadas", Required: true,
					Default: []string{"usage", "quotas", "performance", "costs"}},
				{Name: "period", Type: "string", Description: "Período de análise", 
					Required: false, Default: "month", 
					Enum: []string{"day", "week", "month", "quarter"}},
				{Name: "export_format", Type: "string", Description: "Formato de exportação", 
					Required: false, Default: "json", Enum: []string{"json", "csv", "pdf"}},
			},
			Handler: handleTenantAnalytics,
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

// handleUserManagement gerencia usuários
func handleUserManagement(params map[string]interface{}) (interface{}, error) {
	action := params["action"].(string)
	
	// TODO: Implementar integração real com Auth Service
	
	var result map[string]interface{}
	
	switch action {
	case "list":
		filters := make(map[string]interface{})
		if f, ok := params["filters"].(map[string]interface{}); ok {
			filters = f
		}
		
		result = map[string]interface{}{
			"total_users": 25,
			"users": []map[string]interface{}{
				{
					"id": "user-001",
					"name": "Dr. João Silva",
					"email": "joao.silva@escritorio.com",
					"role": "lawyer",
					"status": "active",
					"last_login": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
					"processes_count": 43,
				},
				{
					"id": "user-002",
					"name": "Dra. Maria Santos",
					"email": "maria.santos@escritorio.com",
					"role": "lawyer",
					"status": "active",
					"last_login": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
					"processes_count": 38,
				},
				{
					"id": "user-003",
					"name": "Ana Costa",
					"email": "ana.costa@escritorio.com",
					"role": "assistant",
					"status": "active",
					"last_login": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					"processes_count": 0,
				},
			},
			"filters_applied": filters,
			"page": 1,
			"per_page": 20,
		}
		
	case "create":
		userData := params["user_data"].(map[string]interface{})
		userID := fmt.Sprintf("user-%d", time.Now().Unix())
		
		result = map[string]interface{}{
			"success": true,
			"user": map[string]interface{}{
				"id": userID,
				"name": userData["name"],
				"email": userData["email"],
				"role": userData["role"],
				"status": "active",
				"created_at": time.Now().Format(time.RFC3339),
				"temporary_password": generateTempPassword(),
			},
			"message": "Usuário criado com sucesso. E-mail de boas-vindas enviado.",
		}
		
	case "update":
		userData := params["user_data"].(map[string]interface{})
		
		result = map[string]interface{}{
			"success": true,
			"user": userData,
			"updated_fields": []string{"name", "role", "permissions"},
			"message": "Usuário atualizado com sucesso",
		}
		
	case "deactivate":
		userData := params["user_data"].(map[string]interface{})
		userID := userData["id"].(string)
		
		result = map[string]interface{}{
			"success": true,
			"user_id": userID,
			"status": "deactivated",
			"deactivated_at": time.Now().Format(time.RFC3339),
			"message": "Usuário desativado. Processos transferidos para o responsável padrão.",
		}
		
	case "permissions":
		userData := params["user_data"].(map[string]interface{})
		permissions := params["permissions"].([]interface{})
		
		result = map[string]interface{}{
			"success": true,
			"user_id": userData["id"],
			"permissions": permissions,
			"effective_permissions": []string{
				"process.read",
				"process.write",
				"client.read",
				"notification.manage",
				"report.generate",
			},
			"message": "Permissões atualizadas com sucesso",
		}
		
	default:
		return nil, fmt.Errorf("ação inválida: %s", action)
	}
	
	return result, nil
}

// handleTenantAnalytics analisa uso do tenant
func handleTenantAnalytics(params map[string]interface{}) (interface{}, error) {
	metrics := params["metrics"].([]interface{})
	period := "month"
	if val, ok := params["period"].(string); ok {
		period = val
	}
	
	// TODO: Implementar integração real com Analytics Service
	
	analyticsData := make(map[string]interface{})
	
	for _, metric := range metrics {
		switch metric.(string) {
		case "usage":
			analyticsData["usage"] = map[string]interface{}{
				"active_users": 15,
				"total_logins": 487,
				"api_calls": 12543,
				"storage_gb": 23.4,
				"processes_monitored": 387,
				"notifications_sent": 1234,
				"trend": map[string]string{
					"users": "+2",
					"api_calls": "+15%",
					"storage": "+8%",
				},
			}
			
		case "quotas":
			analyticsData["quotas"] = map[string]interface{}{
				"plan": "Business",
				"limits": map[string]interface{}{
					"users": map[string]int{"used": 15, "limit": 20},
					"processes": map[string]int{"used": 387, "limit": 500},
					"storage_gb": map[string]float64{"used": 23.4, "limit": 200},
					"api_calls": map[string]int{"used": 12543, "limit": 100000},
					"mcp_commands": map[string]int{"used": 743, "limit": 1000},
				},
				"usage_percentage": map[string]float64{
					"users": 75,
					"processes": 77.4,
					"storage": 11.7,
					"api_calls": 12.5,
					"mcp_commands": 74.3,
				},
				"alerts": []string{
					"Uso de processos próximo do limite (77%)",
					"Comandos MCP em 74% do limite mensal",
				},
			}
			
		case "performance":
			analyticsData["performance"] = map[string]interface{}{
				"average_response_time_ms": 156,
				"uptime_percentage": 99.95,
				"error_rate": 0.002,
				"peak_usage_hour": "14:00-15:00",
				"slowest_endpoints": []map[string]interface{}{
					{"endpoint": "/api/v1/jurisprudence/search", "avg_ms": 890},
					{"endpoint": "/api/v1/reports/generate", "avg_ms": 2340},
				},
				"busiest_users": []map[string]interface{}{
					{"name": "Dr. João Silva", "requests": 3421},
					{"name": "Dra. Maria Santos", "requests": 2987},
				},
			}
			
		case "costs":
			analyticsData["costs"] = map[string]interface{}{
				"monthly_cost": 699.00,
				"cost_breakdown": map[string]float64{
					"base_plan": 699.00,
					"extra_storage": 0,
					"extra_api_calls": 0,
					"extra_users": 0,
				},
				"projected_monthly": 699.00,
				"cost_per_user": 46.60,
				"cost_per_process": 1.81,
				"optimization_suggestions": []string{
					"Considere arquivar processos antigos para liberar espaço",
					"Uso eficiente de API está dentro do plano",
				},
			}
		}
	}
	
	result := map[string]interface{}{
		"tenant_id": params["tenant_id"],
		"tenant_name": "Silva & Associados Advogados",
		"period": period,
		"period_start": getPeriodStart(period).Format("2006-01-02"),
		"period_end": time.Now().Format("2006-01-02"),
		"analytics": analyticsData,
		"summary": map[string]interface{}{
			"health_score": 0.92,
			"growth_rate": "+12%",
			"efficiency_score": 0.88,
			"recommendations": []string{
				"Excelente taxa de utilização dos recursos",
				"Considere upgrade para Enterprise se crescimento continuar",
				"Implemente rotação de logs para otimizar storage",
			},
		},
		"generated_at": time.Now().Format(time.RFC3339),
		"export_url": fmt.Sprintf("/api/v1/analytics/export/%s", params["export_format"]),
	}
	
	return result, nil
}

// generateTempPassword gera senha temporária
func generateTempPassword() string {
	return fmt.Sprintf("DirLux%d!", time.Now().Unix()%10000)
}