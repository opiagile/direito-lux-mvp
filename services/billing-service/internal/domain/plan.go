package domain

import (
	"time"

	"github.com/google/uuid"
)

// Plan representa um plano de assinatura
type Plan struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	DisplayName      string    `json:"display_name"`
	Description      string    `json:"description"`
	PriceMonthly     int64     `json:"price_monthly"`     // em centavos
	PriceYearly      int64     `json:"price_yearly"`      // em centavos
	TrialDays        int       `json:"trial_days"`
	Active           bool      `json:"active"`
	
	// Limites do plano
	MaxProcesses     int `json:"max_processes"`
	MaxUsers         int `json:"max_users"`
	MaxAIRequests    int `json:"max_ai_requests"`
	MaxBotCommands   int `json:"max_bot_commands"`
	
	// Funcionalidades
	HasWhatsApp      bool `json:"has_whatsapp"`
	HasTelegram      bool `json:"has_telegram"`
	HasMCPBot        bool `json:"has_mcp_bot"`
	HasUnlimitedSearch bool `json:"has_unlimited_search"`
	HasAIAnalysis    bool `json:"has_ai_analysis"`
	HasPredictions   bool `json:"has_predictions"`
	HasDocGeneration bool `json:"has_doc_generation"`
	HasWhiteLabel    bool `json:"has_white_label"`
	HasAPIAccess     bool `json:"has_api_access"`
	HasPrioritySupport bool `json:"has_priority_support"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PlanType enumeração dos tipos de planos
type PlanType string

const (
	PlanTypeStarter      PlanType = "starter"
	PlanTypeProfessional PlanType = "professional"
	PlanTypeBusiness     PlanType = "business"
	PlanTypeEnterprise   PlanType = "enterprise"
)

// BillingCycle enumeração dos ciclos de cobrança
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

// NewPlan cria um novo plano
func NewPlan(name, displayName, description string, priceMonthly, priceYearly int64, trialDays int) *Plan {
	return &Plan{
		ID:               uuid.New(),
		Name:             name,
		DisplayName:      displayName,
		Description:      description,
		PriceMonthly:     priceMonthly,
		PriceYearly:      priceYearly,
		TrialDays:        trialDays,
		Active:           true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// GetPrice retorna o preço do plano baseado no ciclo de cobrança
func (p *Plan) GetPrice(cycle BillingCycle) int64 {
	switch cycle {
	case BillingCycleYearly:
		return p.PriceYearly
	default:
		return p.PriceMonthly
	}
}

// GetPriceFormatted retorna o preço formatado em reais
func (p *Plan) GetPriceFormatted(cycle BillingCycle) float64 {
	return float64(p.GetPrice(cycle)) / 100
}

// IsEnterprise verifica se é um plano enterprise
func (p *Plan) IsEnterprise() bool {
	return p.Name == string(PlanTypeEnterprise)
}

// HasFeature verifica se o plano tem uma funcionalidade específica
func (p *Plan) HasFeature(feature string) bool {
	switch feature {
	case "whatsapp":
		return p.HasWhatsApp
	case "telegram":
		return p.HasTelegram
	case "mcp_bot":
		return p.HasMCPBot
	case "unlimited_search":
		return p.HasUnlimitedSearch
	case "ai_analysis":
		return p.HasAIAnalysis
	case "predictions":
		return p.HasPredictions
	case "doc_generation":
		return p.HasDocGeneration
	case "white_label":
		return p.HasWhiteLabel
	case "api_access":
		return p.HasAPIAccess
	case "priority_support":
		return p.HasPrioritySupport
	default:
		return false
	}
}

// GetDefaultPlans retorna os planos padrão do sistema
func GetDefaultPlans() []*Plan {
	return []*Plan{
		{
			ID:               uuid.New(),
			Name:             string(PlanTypeStarter),
			DisplayName:      "Starter",
			Description:      "Ideal para advogados autônomos",
			PriceMonthly:     9900,  // R$ 99,00
			PriceYearly:      99000, // R$ 990,00 (2 meses grátis)
			TrialDays:        15,
			Active:           true,
			MaxProcesses:     50,
			MaxUsers:         2,
			MaxAIRequests:    10,
			MaxBotCommands:   0,
			HasWhatsApp:      true,
			HasTelegram:      true,
			HasMCPBot:        false,
			HasUnlimitedSearch: true,
			HasAIAnalysis:    true,
			HasPredictions:   false,
			HasDocGeneration: false,
			HasWhiteLabel:    false,
			HasAPIAccess:     false,
			HasPrioritySupport: false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			Name:             string(PlanTypeProfessional),
			DisplayName:      "Professional",
			Description:      "Para pequenos escritórios",
			PriceMonthly:     29900, // R$ 299,00
			PriceYearly:      299000, // R$ 2.990,00 (2 meses grátis)
			TrialDays:        15,
			Active:           true,
			MaxProcesses:     200,
			MaxUsers:         5,
			MaxAIRequests:    50,
			MaxBotCommands:   200,
			HasWhatsApp:      true,
			HasTelegram:      true,
			HasMCPBot:        true,
			HasUnlimitedSearch: true,
			HasAIAnalysis:    true,
			HasPredictions:   true,
			HasDocGeneration: true,
			HasWhiteLabel:    false,
			HasAPIAccess:     false,
			HasPrioritySupport: false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			Name:             string(PlanTypeBusiness),
			DisplayName:      "Business",
			Description:      "Para escritórios médios",
			PriceMonthly:     69900, // R$ 699,00
			PriceYearly:      699000, // R$ 6.990,00 (2 meses grátis)
			TrialDays:        15,
			Active:           true,
			MaxProcesses:     500,
			MaxUsers:         15,
			MaxAIRequests:    200,
			MaxBotCommands:   1000,
			HasWhatsApp:      true,
			HasTelegram:      true,
			HasMCPBot:        true,
			HasUnlimitedSearch: true,
			HasAIAnalysis:    true,
			HasPredictions:   true,
			HasDocGeneration: true,
			HasWhiteLabel:    false,
			HasAPIAccess:     true,
			HasPrioritySupport: true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			Name:             string(PlanTypeEnterprise),
			DisplayName:      "Enterprise",
			Description:      "Para grandes escritórios",
			PriceMonthly:     0, // Sob consulta
			PriceYearly:      0, // Sob consulta
			TrialDays:        30,
			Active:           true,
			MaxProcesses:     -1, // Ilimitado
			MaxUsers:         -1, // Ilimitado
			MaxAIRequests:    -1, // Ilimitado
			MaxBotCommands:   -1, // Ilimitado
			HasWhatsApp:      true,
			HasTelegram:      true,
			HasMCPBot:        true,
			HasUnlimitedSearch: true,
			HasAIAnalysis:    true,
			HasPredictions:   true,
			HasDocGeneration: true,
			HasWhiteLabel:    true,
			HasAPIAccess:     true,
			HasPrioritySupport: true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}
}