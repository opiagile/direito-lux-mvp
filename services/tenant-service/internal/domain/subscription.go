package domain

import (
	"time"
	"errors"
)

// Subscription representa uma assinatura de plano
type Subscription struct {
	ID                string              `json:"id" db:"id"`
	TenantID          string              `json:"tenant_id" db:"tenant_id"`
	PlanID            string              `json:"plan_id" db:"plan_id"`
	Status            SubscriptionStatus  `json:"status" db:"status"`
	CurrentPeriodStart time.Time          `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end" db:"current_period_end"`
	CancelAtPeriodEnd  bool               `json:"cancel_at_period_end" db:"cancel_at_period_end"`
	TrialStart         *time.Time         `json:"trial_start" db:"trial_start"`
	TrialEnd           *time.Time         `json:"trial_end" db:"trial_end"`
	CreatedAt          time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" db:"updated_at"`
	CanceledAt         *time.Time         `json:"canceled_at" db:"canceled_at"`
}

// Plan representa um plano de assinatura
type Plan struct {
	ID              string            `json:"id" db:"id"`
	Name            string            `json:"name" db:"name"`
	Type            PlanType          `json:"type" db:"type"`
	Description     string            `json:"description" db:"description"`
	Price           int64             `json:"price" db:"price"` // Em centavos
	Currency        string            `json:"currency" db:"currency"`
	BillingInterval BillingInterval   `json:"billing_interval" db:"billing_interval"`
	Features        PlanFeatures      `json:"features" db:"features"`
	Quotas          PlanQuotas        `json:"quotas" db:"quotas"`
	IsActive        bool              `json:"is_active" db:"is_active"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// PlanFeatures define as funcionalidades disponíveis em cada plano
type PlanFeatures struct {
	WhatsAppEnabled       bool `json:"whatsapp_enabled"`
	AIEnabled             bool `json:"ai_enabled"`
	AdvancedAI            bool `json:"advanced_ai"`
	JurisprudenceEnabled  bool `json:"jurisprudence_enabled"`
	WhiteLabelEnabled     bool `json:"white_label_enabled"`
	CustomIntegrations    bool `json:"custom_integrations"`
	PrioritySupport       bool `json:"priority_support"`
	CustomReports         bool `json:"custom_reports"`
	APIAccess             bool `json:"api_access"`
	WebhooksEnabled       bool `json:"webhooks_enabled"`
}

// PlanQuotas define os limites de cada plano
type PlanQuotas struct {
	MaxProcesses        int `json:"max_processes"`
	MaxUsers            int `json:"max_users"`
	MaxClients          int `json:"max_clients"`
	DataJudQueriesDaily int `json:"datajud_queries_daily"`
	AIQueriesMonthly    int `json:"ai_queries_monthly"`
	StorageGB           int `json:"storage_gb"`
	MaxWebhooks         int `json:"max_webhooks"`
	MaxAPICallsDaily    int `json:"max_api_calls_daily"`
}

// BillingInterval define o intervalo de cobrança
type BillingInterval string

const (
	BillingMonthly BillingInterval = "monthly"
	BillingYearly  BillingInterval = "yearly"
)

// SubscriptionStatus define os status possíveis de uma assinatura
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
	SubscriptionStatusUnpaid    SubscriptionStatus = "unpaid"
)

// SubscriptionRepository define a interface para persistência de assinaturas
type SubscriptionRepository interface {
	Create(subscription *Subscription) error
	GetByID(id string) (*Subscription, error)
	GetByTenantID(tenantID string) (*Subscription, error)
	GetByStatus(status SubscriptionStatus, limit, offset int) ([]*Subscription, error)
	Update(subscription *Subscription) error
	Delete(id string) error
	GetExpiring(days int) ([]*Subscription, error)
}

// PlanRepository define a interface para persistência de planos
type PlanRepository interface {
	Create(plan *Plan) error
	GetByID(id string) (*Plan, error)
	GetByType(planType PlanType) (*Plan, error)
	GetAll(activeOnly bool) ([]*Plan, error)
	Update(plan *Plan) error
	Delete(id string) error
}

// Erros de domínio para assinaturas
var (
	ErrSubscriptionNotFound    = errors.New("assinatura não encontrada")
	ErrSubscriptionExists      = errors.New("assinatura já existe")
	ErrPlanNotFound           = errors.New("plano não encontrado")
	ErrInvalidBillingInterval = errors.New("intervalo de cobrança inválido")
	ErrInvalidSubscriptionStatus = errors.New("status de assinatura inválido")
	ErrTrialAlreadyUsed       = errors.New("período de teste já foi utilizado")
	ErrCannotCancelTrial      = errors.New("não é possível cancelar durante período de teste")
)

// IsActive verifica se a assinatura está ativa
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive || s.Status == SubscriptionStatusTrialing
}

// IsInTrial verifica se a assinatura está em período de teste
func (s *Subscription) IsInTrial() bool {
	if s.TrialEnd == nil {
		return false
	}
	return s.Status == SubscriptionStatusTrialing && time.Now().Before(*s.TrialEnd)
}

// IsExpired verifica se a assinatura expirou
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.CurrentPeriodEnd)
}

// IsExpiring verifica se a assinatura vai expirar em X dias
func (s *Subscription) IsExpiring(days int) bool {
	return time.Now().Add(time.Duration(days)*24*time.Hour).After(s.CurrentPeriodEnd)
}

// Cancel cancela a assinatura
func (s *Subscription) Cancel() {
	s.Status = SubscriptionStatusCanceled
	now := time.Now()
	s.CanceledAt = &now
	s.UpdatedAt = now
}

// ScheduleCancellation agenda cancelamento no final do período
func (s *Subscription) ScheduleCancellation() {
	s.CancelAtPeriodEnd = true
	s.UpdatedAt = time.Now()
}

// Reactivate reativa uma assinatura cancelada
func (s *Subscription) Reactivate() {
	if s.Status == SubscriptionStatusCanceled {
		s.Status = SubscriptionStatusActive
		s.CancelAtPeriodEnd = false
		s.CanceledAt = nil
		s.UpdatedAt = time.Now()
	}
}

// Renew renova a assinatura para o próximo período
func (s *Subscription) Renew(billingInterval BillingInterval) {
	now := time.Now()
	s.CurrentPeriodStart = s.CurrentPeriodEnd
	
	switch billingInterval {
	case BillingMonthly:
		s.CurrentPeriodEnd = s.CurrentPeriodStart.AddDate(0, 1, 0)
	case BillingYearly:
		s.CurrentPeriodEnd = s.CurrentPeriodStart.AddDate(1, 0, 0)
	}
	
	s.Status = SubscriptionStatusActive
	s.UpdatedAt = now
}

// GetDefaultPlanQuotas retorna as quotas padrão para cada tipo de plano
func GetDefaultPlanQuotas(planType PlanType) PlanQuotas {
	switch planType {
	case PlanStarter:
		return PlanQuotas{
			MaxProcesses:        50,
			MaxUsers:           2,
			MaxClients:         20,
			DataJudQueriesDaily: 100,
			AIQueriesMonthly:   10,
			StorageGB:          1,
			MaxWebhooks:        3,
			MaxAPICallsDaily:   1000,
		}
	case PlanProfessional:
		return PlanQuotas{
			MaxProcesses:        200,
			MaxUsers:           5,
			MaxClients:         100,
			DataJudQueriesDaily: 500,
			AIQueriesMonthly:   50,
			StorageGB:          5,
			MaxWebhooks:        10,
			MaxAPICallsDaily:   5000,
		}
	case PlanBusiness:
		return PlanQuotas{
			MaxProcesses:        500,
			MaxUsers:           15,
			MaxClients:         500,
			DataJudQueriesDaily: 2000,
			AIQueriesMonthly:   200,
			StorageGB:          20,
			MaxWebhooks:        25,
			MaxAPICallsDaily:   15000,
		}
	case PlanEnterprise:
		return PlanQuotas{
			MaxProcesses:        -1, // Ilimitado
			MaxUsers:           -1, // Ilimitado
			MaxClients:         -1, // Ilimitado
			DataJudQueriesDaily: 10000,
			AIQueriesMonthly:   -1, // Ilimitado
			StorageGB:          100,
			MaxWebhooks:        -1, // Ilimitado
			MaxAPICallsDaily:   -1, // Ilimitado
		}
	default:
		return PlanQuotas{}
	}
}

// GetDefaultPlanFeatures retorna as funcionalidades padrão para cada tipo de plano
func GetDefaultPlanFeatures(planType PlanType) PlanFeatures {
	switch planType {
	case PlanStarter:
		return PlanFeatures{
			WhatsAppEnabled:      true,
			AIEnabled:            false,
			AdvancedAI:           false,
			JurisprudenceEnabled: false,
			WhiteLabelEnabled:    false,
			CustomIntegrations:   false,
			PrioritySupport:      false,
			CustomReports:        false,
			APIAccess:            false,
			WebhooksEnabled:      false,
		}
	case PlanProfessional:
		return PlanFeatures{
			WhatsAppEnabled:      true,
			AIEnabled:            true,
			AdvancedAI:           false,
			JurisprudenceEnabled: false,
			WhiteLabelEnabled:    false,
			CustomIntegrations:   false,
			PrioritySupport:      false,
			CustomReports:        true,
			APIAccess:            true,
			WebhooksEnabled:      true,
		}
	case PlanBusiness:
		return PlanFeatures{
			WhatsAppEnabled:      true,
			AIEnabled:            true,
			AdvancedAI:           true,
			JurisprudenceEnabled: true,
			WhiteLabelEnabled:    false,
			CustomIntegrations:   true,
			PrioritySupport:      true,
			CustomReports:        true,
			APIAccess:            true,
			WebhooksEnabled:      true,
		}
	case PlanEnterprise:
		return PlanFeatures{
			WhatsAppEnabled:      true,
			AIEnabled:            true,
			AdvancedAI:           true,
			JurisprudenceEnabled: true,
			WhiteLabelEnabled:    true,
			CustomIntegrations:   true,
			PrioritySupport:      true,
			CustomReports:        true,
			APIAccess:            true,
			WebhooksEnabled:      true,
		}
	default:
		return PlanFeatures{}
	}
}

// GetPlanPrice retorna o preço padrão para cada tipo de plano
func GetPlanPrice(planType PlanType) int64 {
	switch planType {
	case PlanStarter:
		return 9900 // R$ 99,00
	case PlanProfessional:
		return 29900 // R$ 299,00
	case PlanBusiness:
		return 69900 // R$ 699,00
	case PlanEnterprise:
		return 199900 // R$ 1.999,00
	default:
		return 0
	}
}