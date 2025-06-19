package tools

import (
	"fmt"
	"time"

	"github.com/direito-lux/mcp-service/internal/domain"
)

// RegisterReportTools registra as ferramentas de relatórios
func RegisterReportTools(registry *domain.ToolRegistry) error {
	tools := []*domain.Tool{
		{
			Name:         "generate_report",
			Description:  "Geração de relatórios personalizados",
			Category:     "report",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "report_type", Type: "string", Description: "Tipo de relatório", Required: true,
					Enum: []string{"productivity", "financial", "processes", "performance", "analytics"}},
				{Name: "date_range", Type: "object", Description: "Intervalo de datas", Required: true},
				{Name: "filters", Type: "object", Description: "Filtros específicos", Required: false},
				{Name: "format", Type: "string", Description: "Formato de saída", 
					Required: false, Default: "pdf", Enum: []string{"pdf", "excel", "json", "csv"}},
				{Name: "include_charts", Type: "boolean", Description: "Incluir gráficos", 
					Required: false, Default: true},
			},
			Handler: handleGenerateReport,
		},
		{
			Name:         "dashboard_metrics",
			Description:  "Métricas para dashboard em tempo real",
			Category:     "report",
			RequiredPlan: "professional",
			Parameters: []domain.ToolParameter{
				{Name: "metric_types", Type: "array", Description: "Tipos de métricas", Required: true,
					Default: []string{"processes_count", "notifications_sent", "search_volume"}},
				{Name: "period", Type: "string", Description: "Período de análise", 
					Required: false, Default: "today", 
					Enum: []string{"today", "week", "month", "quarter", "year"}},
				{Name: "tenant_filter", Type: "string", Description: "Filtro por tenant", Required: false},
				{Name: "user_filter", Type: "string", Description: "Filtro por usuário", Required: false},
			},
			Handler: handleDashboardMetrics,
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

// handleGenerateReport gera relatórios
func handleGenerateReport(params map[string]interface{}) (interface{}, error) {
	reportType := params["report_type"].(string)
	dateRange := params["date_range"].(map[string]interface{})
	format := "pdf"
	if val, ok := params["format"].(string); ok {
		format = val
	}
	includeCharts := true
	if val, ok := params["include_charts"].(bool); ok {
		includeCharts = val
	}
	
	// TODO: Implementar integração real com Report Service
	
	reportID := fmt.Sprintf("report-%s-%d", reportType, time.Now().Unix())
	
	// Simular dados do relatório baseado no tipo
	var reportData map[string]interface{}
	
	switch reportType {
	case "productivity":
		reportData = map[string]interface{}{
			"summary": map[string]interface{}{
				"total_processes": 127,
				"new_processes": 43,
				"completed_processes": 31,
				"average_resolution_time_days": 74,
				"success_rate": 0.87,
			},
			"by_lawyer": []map[string]interface{}{
				{
					"name": "Dr. João Silva",
					"new_processes": 18,
					"completed": 14,
					"in_progress": 25,
					"success_rate": 0.92,
				},
				{
					"name": "Dra. Maria Santos",
					"new_processes": 15,
					"completed": 11,
					"in_progress": 22,
					"success_rate": 0.85,
				},
			},
			"by_process_type": map[string]int{
				"civil": 45,
				"trabalhista": 38,
				"criminal": 22,
				"tributário": 22,
			},
		}
		
	case "financial":
		reportData = map[string]interface{}{
			"summary": map[string]interface{}{
				"total_revenue": 425000.00,
				"pending_payments": 87500.00,
				"average_process_value": 15750.00,
				"payment_rate": 0.82,
			},
			"by_month": []map[string]interface{}{
				{"month": "Janeiro", "revenue": 145000.00, "expenses": 42000.00},
				{"month": "Fevereiro", "revenue": 132000.00, "expenses": 38500.00},
				{"month": "Março", "revenue": 148000.00, "expenses": 44200.00},
			},
			"by_client": []map[string]interface{}{
				{"client": "Empresa ABC", "total": 85000.00, "status": "paid"},
				{"client": "João Silva", "total": 25000.00, "status": "pending"},
			},
		}
		
	default:
		reportData = map[string]interface{}{
			"total_items": 100,
			"period": dateRange,
			"generated_at": time.Now().Format(time.RFC3339),
		}
	}
	
	result := map[string]interface{}{
		"success": true,
		"report_id": reportID,
		"report_type": reportType,
		"format": format,
		"status": "generating",
		"data": reportData,
		"metadata": map[string]interface{}{
			"generated_at": time.Now().Format(time.RFC3339),
			"period": dateRange,
			"include_charts": includeCharts,
			"page_count": calculatePageCount(reportType),
			"estimated_time_seconds": 30,
		},
		"download_url": fmt.Sprintf("/api/v1/reports/%s/download", reportID),
		"preview_url": fmt.Sprintf("/api/v1/reports/%s/preview", reportID),
	}
	
	return result, nil
}

// handleDashboardMetrics obtém métricas do dashboard
func handleDashboardMetrics(params map[string]interface{}) (interface{}, error) {
	metricTypes := params["metric_types"].([]interface{})
	period := "today"
	if val, ok := params["period"].(string); ok {
		period = val
	}
	
	// TODO: Implementar integração real com Analytics Service
	
	metrics := make(map[string]interface{})
	
	for _, metricType := range metricTypes {
		switch metricType.(string) {
		case "processes_count":
			metrics["processes_count"] = map[string]interface{}{
				"total": 523,
				"active": 387,
				"archived": 112,
				"suspended": 24,
				"trend": "+12%",
				"comparison": "vs período anterior",
			}
			
		case "notifications_sent":
			metrics["notifications_sent"] = map[string]interface{}{
				"total": 1847,
				"whatsapp": 923,
				"email": 615,
				"telegram": 309,
				"delivery_rate": 0.96,
				"trend": "+8%",
			}
			
		case "search_volume":
			metrics["search_volume"] = map[string]interface{}{
				"total_searches": 3421,
				"unique_users": 187,
				"average_per_user": 18.3,
				"popular_terms": []string{"processo", "audiência", "prazo", "sentença"},
				"trend": "+15%",
			}
			
		case "user_activity":
			metrics["user_activity"] = map[string]interface{}{
				"active_users": 142,
				"sessions": 489,
				"average_session_minutes": 23,
				"peak_hour": "14:00-15:00",
				"trend": "+5%",
			}
			
		case "system_performance":
			metrics["system_performance"] = map[string]interface{}{
				"uptime": 0.999,
				"response_time_ms": 145,
				"error_rate": 0.002,
				"api_calls": 12847,
				"status": "healthy",
			}
		}
	}
	
	result := map[string]interface{}{
		"period": period,
		"period_start": getPeriodStart(period).Format(time.RFC3339),
		"period_end": time.Now().Format(time.RFC3339),
		"metrics": metrics,
		"summary": map[string]interface{}{
			"health_score": 0.95,
			"alerts": []map[string]interface{}{
				{
					"type": "info",
					"message": "Aumento de 15% no volume de buscas",
					"metric": "search_volume",
				},
			},
			"recommendations": []string{
				"Considere aumentar a capacidade do servidor devido ao crescimento de uso",
				"87% dos usuários preferem notificações via WhatsApp",
			},
		},
		"last_updated": time.Now().Format(time.RFC3339),
		"refresh_interval_seconds": 300,
	}
	
	return result, nil
}

// calculatePageCount calcula número de páginas do relatório
func calculatePageCount(reportType string) int {
	switch reportType {
	case "productivity":
		return 8
	case "financial":
		return 12
	case "processes":
		return 15
	case "performance":
		return 6
	default:
		return 5
	}
}

// getPeriodStart obtém início do período
func getPeriodStart(period string) time.Time {
	now := time.Now()
	
	switch period {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		return now.AddDate(0, 0, -7)
	case "month":
		return now.AddDate(0, -1, 0)
	case "quarter":
		return now.AddDate(0, -3, 0)
	case "year":
		return now.AddDate(-1, 0, 0)
	default:
		return now
	}
}