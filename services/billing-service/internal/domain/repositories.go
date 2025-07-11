package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PlanRepository interface do repositório de planos
type PlanRepository interface {
	// Create cria um novo plano
	Create(ctx context.Context, plan *Plan) error
	
	// GetByID busca um plano por ID
	GetByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	
	// GetByName busca um plano por nome
	GetByName(ctx context.Context, name string) (*Plan, error)
	
	// GetAll busca todos os planos
	GetAll(ctx context.Context) ([]*Plan, error)
	
	// GetActive busca apenas planos ativos
	GetActive(ctx context.Context) ([]*Plan, error)
	
	// Update atualiza um plano
	Update(ctx context.Context, plan *Plan) error
	
	// Delete deleta um plano
	Delete(ctx context.Context, id uuid.UUID) error
	
	// InitializeDefaultPlans inicializa os planos padrão
	InitializeDefaultPlans(ctx context.Context) error
}

// SubscriptionRepository interface do repositório de assinaturas
type SubscriptionRepository interface {
	// Create cria uma nova assinatura
	Create(ctx context.Context, subscription *Subscription) error
	
	// GetByID busca uma assinatura por ID
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	
	// GetByTenantID busca assinaturas por tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*Subscription, error)
	
	// GetCurrentByTenantID busca a assinatura atual do tenant
	GetCurrentByTenantID(ctx context.Context, tenantID uuid.UUID) (*Subscription, error)
	
	// GetByStatus busca assinaturas por status
	GetByStatus(ctx context.Context, status SubscriptionStatus) ([]*Subscription, error)
	
	// GetExpiring busca assinaturas expirando
	GetExpiring(ctx context.Context, days int) ([]*Subscription, error)
	
	// GetPendingBilling busca assinaturas pendentes de cobrança
	GetPendingBilling(ctx context.Context) ([]*Subscription, error)
	
	// GetByAsaasID busca assinatura por ID do ASAAS
	GetByAsaasID(ctx context.Context, asaasID string) (*Subscription, error)
	
	// Update atualiza uma assinatura
	Update(ctx context.Context, subscription *Subscription) error
	
	// Delete deleta uma assinatura
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetStats retorna estatísticas de assinaturas
	GetStats(ctx context.Context) (*SubscriptionStats, error)
}

// PaymentRepository interface do repositório de pagamentos
type PaymentRepository interface {
	// Create cria um novo pagamento
	Create(ctx context.Context, payment *Payment) error
	
	// GetByID busca um pagamento por ID
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	
	// GetBySubscriptionID busca pagamentos por assinatura
	GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]*Payment, error)
	
	// GetByTenantID busca pagamentos por tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*Payment, error)
	
	// GetByStatus busca pagamentos por status
	GetByStatus(ctx context.Context, status PaymentStatus) ([]*Payment, error)
	
	// GetPending busca pagamentos pendentes
	GetPending(ctx context.Context) ([]*Payment, error)
	
	// GetFailedRetriable busca pagamentos falhos que podem ser tentados novamente
	GetFailedRetriable(ctx context.Context) ([]*Payment, error)
	
	// GetByAsaasID busca pagamento por ID do ASAAS
	GetByAsaasID(ctx context.Context, asaasID string) (*Payment, error)
	
	// GetByNOWPaymentID busca pagamento por ID do NOWPayments
	GetByNOWPaymentID(ctx context.Context, nowPaymentID string) (*Payment, error)
	
	// Update atualiza um pagamento
	Update(ctx context.Context, payment *Payment) error
	
	// Delete deleta um pagamento
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetStats retorna estatísticas de pagamentos
	GetStats(ctx context.Context, tenantID *uuid.UUID) (*PaymentStats, error)
}

// InvoiceRepository interface do repositório de faturas
type InvoiceRepository interface {
	// Create cria uma nova fatura
	Create(ctx context.Context, invoice *Invoice) error
	
	// GetByID busca uma fatura por ID
	GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	
	// GetByNumber busca uma fatura por número
	GetByNumber(ctx context.Context, number string) (*Invoice, error)
	
	// GetBySubscriptionID busca faturas por assinatura
	GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]*Invoice, error)
	
	// GetByTenantID busca faturas por tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*Invoice, error)
	
	// GetByStatus busca faturas por status
	GetByStatus(ctx context.Context, status InvoiceStatus) ([]*Invoice, error)
	
	// GetOverdue busca faturas vencidas
	GetOverdue(ctx context.Context) ([]*Invoice, error)
	
	// GetByAsaasID busca fatura por ID do ASAAS
	GetByAsaasID(ctx context.Context, asaasID string) (*Invoice, error)
	
	// GetByNFE busca fatura por número da NF-e
	GetByNFE(ctx context.Context, nfeNumber string) (*Invoice, error)
	
	// Update atualiza uma fatura
	Update(ctx context.Context, invoice *Invoice) error
	
	// Delete deleta uma fatura
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetStats retorna estatísticas de faturas
	GetStats(ctx context.Context, tenantID *uuid.UUID) (*InvoiceStats, error)
}

// CustomerRepository interface do repositório de clientes
type CustomerRepository interface {
	// Create cria um novo cliente
	Create(ctx context.Context, customer *Customer) error
	
	// GetByID busca um cliente por ID
	GetByID(ctx context.Context, id uuid.UUID) (*Customer, error)
	
	// GetByTenantID busca clientes por tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*Customer, error)
	
	// GetByDocument busca um cliente por documento
	GetByDocument(ctx context.Context, document string) (*Customer, error)
	
	// GetByEmail busca um cliente por email
	GetByEmail(ctx context.Context, email string) (*Customer, error)
	
	// GetByAsaasID busca cliente por ID do ASAAS
	GetByAsaasID(ctx context.Context, asaasID string) (*Customer, error)
	
	// Update atualiza um cliente
	Update(ctx context.Context, customer *Customer) error
	
	// Delete deleta um cliente
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetStats retorna estatísticas de clientes
	GetStats(ctx context.Context, tenantID *uuid.UUID) (*CustomerStats, error)
}

// Estruturas de estatísticas

// SubscriptionStats estatísticas de assinaturas
type SubscriptionStats struct {
	TotalSubscriptions    int64 `json:"total_subscriptions"`
	ActiveSubscriptions   int64 `json:"active_subscriptions"`
	TrialSubscriptions    int64 `json:"trial_subscriptions"`
	CancelledSubscriptions int64 `json:"cancelled_subscriptions"`
	ExpiredSubscriptions  int64 `json:"expired_subscriptions"`
	
	// Por plano
	ByPlan map[string]int64 `json:"by_plan"`
	
	// Receita
	MonthlyRevenue        int64 `json:"monthly_revenue"`
	YearlyRevenue         int64 `json:"yearly_revenue"`
	TotalRevenue          int64 `json:"total_revenue"`
	
	// Métricas de crescimento
	NewThisMonth          int64 `json:"new_this_month"`
	ChurnThisMonth        int64 `json:"churn_this_month"`
	ChurnRate             float64 `json:"churn_rate"`
	
	// Previsões
	PredictedNextMonth    int64 `json:"predicted_next_month"`
	
	LastUpdated           time.Time `json:"last_updated"`
}

// PaymentStats estatísticas de pagamentos
type PaymentStats struct {
	TotalPayments         int64 `json:"total_payments"`
	SuccessfulPayments    int64 `json:"successful_payments"`
	FailedPayments        int64 `json:"failed_payments"`
	PendingPayments       int64 `json:"pending_payments"`
	
	// Por método
	ByPaymentMethod       map[string]int64 `json:"by_payment_method"`
	
	// Valores
	TotalAmount           int64 `json:"total_amount"`
	SuccessfulAmount      int64 `json:"successful_amount"`
	FailedAmount          int64 `json:"failed_amount"`
	PendingAmount         int64 `json:"pending_amount"`
	
	// Métricas
	SuccessRate           float64 `json:"success_rate"`
	AverageAmount         float64 `json:"average_amount"`
	
	// Tendências
	ThisMonthPayments     int64 `json:"this_month_payments"`
	LastMonthPayments     int64 `json:"last_month_payments"`
	GrowthRate            float64 `json:"growth_rate"`
	
	LastUpdated           time.Time `json:"last_updated"`
}

// InvoiceStats estatísticas de faturas
type InvoiceStats struct {
	TotalInvoices         int64 `json:"total_invoices"`
	PaidInvoices          int64 `json:"paid_invoices"`
	OverdueInvoices       int64 `json:"overdue_invoices"`
	CancelledInvoices     int64 `json:"cancelled_invoices"`
	
	// Valores
	TotalAmount           int64 `json:"total_amount"`
	PaidAmount            int64 `json:"paid_amount"`
	OverdueAmount         int64 `json:"overdue_amount"`
	
	// Métricas
	PaymentRate           float64 `json:"payment_rate"`
	AveragePaymentTime    float64 `json:"average_payment_time_days"`
	
	// NF-e
	WithNFE               int64 `json:"with_nfe"`
	NFESuccessRate        float64 `json:"nfe_success_rate"`
	
	LastUpdated           time.Time `json:"last_updated"`
}

// CustomerStats estatísticas de clientes
type CustomerStats struct {
	TotalCustomers        int64 `json:"total_customers"`
	ActiveCustomers       int64 `json:"active_customers"`
	InactiveCustomers     int64 `json:"inactive_customers"`
	BlockedCustomers      int64 `json:"blocked_customers"`
	
	// Por tipo
	CorporateCustomers    int64 `json:"corporate_customers"`
	IndividualCustomers   int64 `json:"individual_customers"`
	
	// Crescimento
	NewThisMonth          int64 `json:"new_this_month"`
	
	LastUpdated           time.Time `json:"last_updated"`
}

// Queries para filtros avançados

// SubscriptionQuery filtros para busca de assinaturas
type SubscriptionQuery struct {
	TenantID     *uuid.UUID         `json:"tenant_id"`
	PlanID       *uuid.UUID         `json:"plan_id"`
	Status       []SubscriptionStatus `json:"status"`
	PaymentMethod []PaymentMethod    `json:"payment_method"`
	
	// Filtros de data
	CreatedAfter  *time.Time `json:"created_after"`
	CreatedBefore *time.Time `json:"created_before"`
	
	// Filtros de valor
	MinAmount *int64 `json:"min_amount"`
	MaxAmount *int64 `json:"max_amount"`
	
	// Paginação
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	
	// Ordenação
	OrderBy   string `json:"order_by"`
	OrderDesc bool   `json:"order_desc"`
}

// PaymentQuery filtros para busca de pagamentos
type PaymentQuery struct {
	TenantID       *uuid.UUID       `json:"tenant_id"`
	SubscriptionID *uuid.UUID       `json:"subscription_id"`
	Status         []PaymentStatus  `json:"status"`
	PaymentMethod  []PaymentMethod  `json:"payment_method"`
	
	// Filtros de data
	CreatedAfter  *time.Time `json:"created_after"`
	CreatedBefore *time.Time `json:"created_before"`
	PaidAfter     *time.Time `json:"paid_after"`
	PaidBefore    *time.Time `json:"paid_before"`
	
	// Filtros de valor
	MinAmount *int64 `json:"min_amount"`
	MaxAmount *int64 `json:"max_amount"`
	
	// Paginação
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	
	// Ordenação
	OrderBy   string `json:"order_by"`
	OrderDesc bool   `json:"order_desc"`
}