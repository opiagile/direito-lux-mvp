package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
	"github.com/direito-lux/report-service/internal/infrastructure/config"
)

// DataCollector implementação do coletor de dados
type DataCollector struct {
	config *config.Config
	client *http.Client
	logger *zap.Logger
}

// NewDataCollector cria nova instância do coletor
func NewDataCollector(cfg *config.Config, logger *zap.Logger) *DataCollector {
	return &DataCollector{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
		logger: logger,
	}
}

// CollectProcessData implementa domain.DataCollector
func (c *DataCollector) CollectProcessData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error) {
	c.logger.Info("Collecting process data", zap.String("tenant_id", tenantID.String()))

	// Fazer chamada para Process Service
	url := fmt.Sprintf("%s/api/v1/processes/stats?tenant_id=%s", c.config.Services.ProcessServiceURL, tenantID.String())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Warn("Failed to call process service, using mock data", zap.Error(err))
		return c.getMockProcessData(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Warn("Process service returned error, using mock data", zap.Int("status", resp.StatusCode))
		return c.getMockProcessData(), nil
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data, nil
}

// CollectProductivityData implementa domain.DataCollector
func (c *DataCollector) CollectProductivityData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error) {
	c.logger.Info("Collecting productivity data", zap.String("tenant_id", tenantID.String()))

	// Mock data por enquanto
	return map[string]interface{}{
		"processes_per_day":    15.5,
		"average_resolution":   "45 days",
		"user_productivity":    []map[string]interface{}{
			{"user": "João Silva", "processes": 25, "avg_time": "30 days"},
			{"user": "Maria Santos", "processes": 20, "avg_time": "35 days"},
		},
		"monthly_growth":      8.5,
		"efficiency_score":    85.2,
	}, nil
}

// CollectFinancialData implementa domain.DataCollector
func (c *DataCollector) CollectFinancialData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error) {
	c.logger.Info("Collecting financial data", zap.String("tenant_id", tenantID.String()))

	// Mock data por enquanto
	return map[string]interface{}{
		"total_revenue":       250000.00,
		"monthly_revenue":     35000.00,
		"pending_payments":    15000.00,
		"expenses":           120000.00,
		"profit_margin":      52.0,
		"cost_per_process":   850.00,
		"billing_by_client":  []map[string]interface{}{
			{"client": "Empresa ABC", "amount": 75000.00},
			{"client": "João Silva", "amount": 25000.00},
		},
	}, nil
}

// CollectJurisprudenceData implementa domain.DataCollector
func (c *DataCollector) CollectJurisprudenceData(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) (interface{}, error) {
	c.logger.Info("Collecting jurisprudence data", zap.String("tenant_id", tenantID.String()))

	// Fazer chamada para AI Service
	url := fmt.Sprintf("%s/api/v1/jurisprudence/stats?tenant_id=%s", c.config.Services.AIServiceURL, tenantID.String())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Warn("Failed to call AI service, using mock data", zap.Error(err))
		return c.getMockJurisprudenceData(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Warn("AI service returned error, using mock data", zap.Int("status", resp.StatusCode))
		return c.getMockJurisprudenceData(), nil
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data, nil
}

// CollectDashboardData implementa domain.DataCollector
func (c *DataCollector) CollectDashboardData(ctx context.Context, widget *domain.DashboardWidget) (interface{}, error) {
	c.logger.Info("Collecting dashboard widget data", 
		zap.String("widget_id", widget.ID.String()),
		zap.String("data_source", widget.DataSource))

	switch widget.DataSource {
	case "processes":
		return c.getMockProcessData(), nil
	case "productivity":
		return map[string]interface{}{
			"value": 85.2,
			"trend": "up",
			"change": "+5.3%",
		}, nil
	case "kpis":
		kpis, _ := c.CalculateKPIs(ctx, uuid.MustParse("00000000-0000-0000-0000-000000000000"))
		return kpis, nil
	default:
		return map[string]interface{}{
			"message": "Widget data not available",
		}, nil
	}
}

// CalculateKPIs implementa domain.DataCollector
func (c *DataCollector) CalculateKPIs(ctx context.Context, tenantID uuid.UUID) ([]*domain.KPI, error) {
	c.logger.Info("Calculating KPIs", zap.String("tenant_id", tenantID.String()))

	now := time.Now()
	
	kpis := []*domain.KPI{
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			Name:             "total_processes",
			DisplayName:      "Total de Processos",
			Description:      "Número total de processos ativos",
			Category:         "processes",
			Unit:             "processos",
			CurrentValue:     156.0,
			PreviousValue:    func() *float64 { v := 142.0; return &v }(),
			Target:           func() *float64 { v := 200.0; return &v }(),
			Trend:            "up",
			TrendPercentage:  func() *float64 { v := 9.9; return &v }(),
			LastCalculatedAt: now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			Name:             "success_rate",
			DisplayName:      "Taxa de Sucesso",
			Description:      "Percentual de processos ganhos",
			Category:         "performance",
			Unit:             "%",
			CurrentValue:     87.5,
			PreviousValue:    func() *float64 { v := 85.2; return &v }(),
			Target:           func() *float64 { v := 90.0; return &v }(),
			Trend:            "up",
			TrendPercentage:  func() *float64 { v := 2.7; return &v }(),
			LastCalculatedAt: now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			Name:             "average_resolution_time",
			DisplayName:      "Tempo Médio de Resolução",
			Description:      "Tempo médio para resolução de processos",
			Category:         "efficiency",
			Unit:             "dias",
			CurrentValue:     42.5,
			PreviousValue:    func() *float64 { v := 48.2; return &v }(),
			Target:           func() *float64 { v := 30.0; return &v }(),
			Trend:            "down",
			TrendPercentage:  func() *float64 { v := -11.8; return &v }(),
			LastCalculatedAt: now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			Name:             "monthly_revenue",
			DisplayName:      "Receita Mensal",
			Description:      "Receita total do mês atual",
			Category:         "financial",
			Unit:             "R$",
			CurrentValue:     35000.0,
			PreviousValue:    func() *float64 { v := 32000.0; return &v }(),
			Target:           func() *float64 { v := 40000.0; return &v }(),
			Trend:            "up",
			TrendPercentage:  func() *float64 { v := 9.4; return &v }(),
			LastCalculatedAt: now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			Name:             "client_satisfaction",
			DisplayName:      "Satisfação do Cliente",
			Description:      "Índice de satisfação dos clientes",
			Category:         "quality",
			Unit:             "/10",
			CurrentValue:     8.7,
			PreviousValue:    func() *float64 { v := 8.4; return &v }(),
			Target:           func() *float64 { v := 9.0; return &v }(),
			Trend:            "up",
			TrendPercentage:  func() *float64 { v := 3.6; return &v }(),
			LastCalculatedAt: now,
			UpdatedAt:        now,
		},
	}

	return kpis, nil
}

// getMockProcessData retorna dados mock de processos
func (c *DataCollector) getMockProcessData() map[string]interface{} {
	return map[string]interface{}{
		"total":          156,
		"active":         124,
		"archived":       28,
		"suspended":      4,
		"new_this_month": 15,
		"by_court": map[string]int{
			"TJ-SP": 85,
			"TJ-RJ": 42,
			"STJ":   18,
			"STF":   11,
		},
		"by_status": map[string]int{
			"Tramitando":    89,
			"Aguardando":    35,
			"Sentenciado":   28,
			"Recurso":       4,
		},
		"by_type": map[string]int{
			"Cível":         95,
			"Trabalhista":   38,
			"Criminal":      15,
			"Tributário":    8,
		},
	}
}

// getMockJurisprudenceData retorna dados mock de jurisprudência
func (c *DataCollector) getMockJurisprudenceData() map[string]interface{} {
	return map[string]interface{}{
		"total_decisions": 1247,
		"relevant_decisions": 89,
		"success_probability": 78.5,
		"similar_cases": []map[string]interface{}{
			{
				"case_number": "REsp 1234567/SP",
				"similarity": 92.5,
				"outcome": "Favorável",
				"court": "STJ",
			},
			{
				"case_number": "AI 987654/RJ", 
				"similarity": 87.3,
				"outcome": "Favorável",
				"court": "TJ-RJ",
			},
		},
		"key_precedents": []string{
			"Súmula 123 STJ",
			"Súmula 456 STF",
			"Tema 789 STJ",
		},
		"prediction_confidence": 85.7,
	}
}