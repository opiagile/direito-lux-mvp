package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event interface base para todos os eventos
type Event interface {
	GetID() uuid.UUID
	GetType() string
	GetTenantID() uuid.UUID
	GetTimestamp() time.Time
	GetData() interface{}
	GetVersion() int
}

// BaseEvent implementação base para eventos
type BaseEvent struct {
	ID        uuid.UUID   `json:"id"`
	Type      string      `json:"type"`
	TenantID  uuid.UUID   `json:"tenant_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Version   int         `json:"version"`
}

func (e BaseEvent) GetID() uuid.UUID      { return e.ID }
func (e BaseEvent) GetType() string       { return e.Type }
func (e BaseEvent) GetTenantID() uuid.UUID { return e.TenantID }
func (e BaseEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e BaseEvent) GetData() interface{}  { return e.Data }
func (e BaseEvent) GetVersion() int       { return e.Version }

// NewBaseEvent cria um novo evento base
func NewBaseEvent(eventType string, tenantID uuid.UUID, data interface{}) BaseEvent {
	return BaseEvent{
		ID:        uuid.New(),
		Type:      eventType,
		TenantID:  tenantID,
		Timestamp: time.Now(),
		Data:      data,
		Version:   1,
	}
}

// Eventos de Assinatura

// SubscriptionCreatedEvent evento de criação de assinatura
type SubscriptionCreatedEvent struct {
	BaseEvent
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PlanID         uuid.UUID `json:"plan_id"`
	PlanName       string    `json:"plan_name"`
	Amount         int64     `json:"amount"`
	BillingCycle   string    `json:"billing_cycle"`
	PaymentMethod  string    `json:"payment_method"`
	TrialDays      int       `json:"trial_days"`
}

func NewSubscriptionCreatedEvent(tenantID, subscriptionID, planID uuid.UUID, planName string, amount int64, billingCycle BillingCycle, paymentMethod PaymentMethod, trialDays int) *SubscriptionCreatedEvent {
	return &SubscriptionCreatedEvent{
		BaseEvent:      NewBaseEvent("subscription.created", tenantID, nil),
		SubscriptionID: subscriptionID,
		PlanID:         planID,
		PlanName:       planName,
		Amount:         amount,
		BillingCycle:   string(billingCycle),
		PaymentMethod:  string(paymentMethod),
		TrialDays:      trialDays,
	}
}

// SubscriptionActivatedEvent evento de ativação de assinatura
type SubscriptionActivatedEvent struct {
	BaseEvent
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PlanName       string    `json:"plan_name"`
	Amount         int64     `json:"amount"`
}

func NewSubscriptionActivatedEvent(tenantID, subscriptionID uuid.UUID, planName string, amount int64) *SubscriptionActivatedEvent {
	return &SubscriptionActivatedEvent{
		BaseEvent:      NewBaseEvent("subscription.activated", tenantID, nil),
		SubscriptionID: subscriptionID,
		PlanName:       planName,
		Amount:         amount,
	}
}

// SubscriptionCancelledEvent evento de cancelamento de assinatura
type SubscriptionCancelledEvent struct {
	BaseEvent
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PlanName       string    `json:"plan_name"`
	CancelReason   string    `json:"cancel_reason"`
	CancelledBy    uuid.UUID `json:"cancelled_by"`
}

func NewSubscriptionCancelledEvent(tenantID, subscriptionID uuid.UUID, planName, cancelReason string, cancelledBy uuid.UUID) *SubscriptionCancelledEvent {
	return &SubscriptionCancelledEvent{
		BaseEvent:      NewBaseEvent("subscription.cancelled", tenantID, nil),
		SubscriptionID: subscriptionID,
		PlanName:       planName,
		CancelReason:   cancelReason,
		CancelledBy:    cancelledBy,
	}
}

// SubscriptionExpiredEvent evento de expiração de assinatura
type SubscriptionExpiredEvent struct {
	BaseEvent
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PlanName       string    `json:"plan_name"`
	LastPaymentDate *time.Time `json:"last_payment_date"`
}

func NewSubscriptionExpiredEvent(tenantID, subscriptionID uuid.UUID, planName string, lastPaymentDate *time.Time) *SubscriptionExpiredEvent {
	return &SubscriptionExpiredEvent{
		BaseEvent:       NewBaseEvent("subscription.expired", tenantID, nil),
		SubscriptionID:  subscriptionID,
		PlanName:        planName,
		LastPaymentDate: lastPaymentDate,
	}
}

// SubscriptionTrialEndingEvent evento de fim de trial se aproximando
type SubscriptionTrialEndingEvent struct {
	BaseEvent
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PlanName       string    `json:"plan_name"`
	DaysRemaining  int       `json:"days_remaining"`
	TrialEndDate   time.Time `json:"trial_end_date"`
}

func NewSubscriptionTrialEndingEvent(tenantID, subscriptionID uuid.UUID, planName string, daysRemaining int, trialEndDate time.Time) *SubscriptionTrialEndingEvent {
	return &SubscriptionTrialEndingEvent{
		BaseEvent:      NewBaseEvent("subscription.trial_ending", tenantID, nil),
		SubscriptionID: subscriptionID,
		PlanName:       planName,
		DaysRemaining:  daysRemaining,
		TrialEndDate:   trialEndDate,
	}
}

// Eventos de Pagamento

// PaymentCreatedEvent evento de criação de pagamento
type PaymentCreatedEvent struct {
	BaseEvent
	PaymentID      uuid.UUID `json:"payment_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Amount         int64     `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	Currency       string    `json:"currency"`
	DueDate        *time.Time `json:"due_date"`
}

func NewPaymentCreatedEvent(tenantID, paymentID, subscriptionID uuid.UUID, amount int64, paymentMethod PaymentMethod, currency string, dueDate *time.Time) *PaymentCreatedEvent {
	return &PaymentCreatedEvent{
		BaseEvent:      NewBaseEvent("payment.created", tenantID, nil),
		PaymentID:      paymentID,
		SubscriptionID: subscriptionID,
		Amount:         amount,
		PaymentMethod:  string(paymentMethod),
		Currency:       currency,
		DueDate:        dueDate,
	}
}

// PaymentSuccessEvent evento de pagamento bem-sucedido
type PaymentSuccessEvent struct {
	BaseEvent
	PaymentID      uuid.UUID `json:"payment_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Amount         int64     `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	TransactionID  string    `json:"transaction_id"`
	PaidAt         time.Time `json:"paid_at"`
}

func NewPaymentSuccessEvent(tenantID, paymentID, subscriptionID uuid.UUID, amount int64, paymentMethod PaymentMethod, transactionID string, paidAt time.Time) *PaymentSuccessEvent {
	return &PaymentSuccessEvent{
		BaseEvent:      NewBaseEvent("payment.success", tenantID, nil),
		PaymentID:      paymentID,
		SubscriptionID: subscriptionID,
		Amount:         amount,
		PaymentMethod:  string(paymentMethod),
		TransactionID:  transactionID,
		PaidAt:         paidAt,
	}
}

// PaymentFailedEvent evento de falha de pagamento
type PaymentFailedEvent struct {
	BaseEvent
	PaymentID      uuid.UUID `json:"payment_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Amount         int64     `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	FailureReason  string    `json:"failure_reason"`
	RetryCount     int       `json:"retry_count"`
	NextRetryAt    *time.Time `json:"next_retry_at"`
}

func NewPaymentFailedEvent(tenantID, paymentID, subscriptionID uuid.UUID, amount int64, paymentMethod PaymentMethod, failureReason string, retryCount int, nextRetryAt *time.Time) *PaymentFailedEvent {
	return &PaymentFailedEvent{
		BaseEvent:      NewBaseEvent("payment.failed", tenantID, nil),
		PaymentID:      paymentID,
		SubscriptionID: subscriptionID,
		Amount:         amount,
		PaymentMethod:  string(paymentMethod),
		FailureReason:  failureReason,
		RetryCount:     retryCount,
		NextRetryAt:    nextRetryAt,
	}
}

// PaymentRefundedEvent evento de reembolso de pagamento
type PaymentRefundedEvent struct {
	BaseEvent
	PaymentID      uuid.UUID `json:"payment_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	OriginalAmount int64     `json:"original_amount"`
	RefundAmount   int64     `json:"refund_amount"`
	RefundReason   string    `json:"refund_reason"`
	RefundedAt     time.Time `json:"refunded_at"`
}

func NewPaymentRefundedEvent(tenantID, paymentID, subscriptionID uuid.UUID, originalAmount, refundAmount int64, refundReason string, refundedAt time.Time) *PaymentRefundedEvent {
	return &PaymentRefundedEvent{
		BaseEvent:      NewBaseEvent("payment.refunded", tenantID, nil),
		PaymentID:      paymentID,
		SubscriptionID: subscriptionID,
		OriginalAmount: originalAmount,
		RefundAmount:   refundAmount,
		RefundReason:   refundReason,
		RefundedAt:     refundedAt,
	}
}

// Eventos de Fatura

// InvoiceCreatedEvent evento de criação de fatura
type InvoiceCreatedEvent struct {
	BaseEvent
	InvoiceID      uuid.UUID `json:"invoice_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Number         string    `json:"number"`
	Amount         int64     `json:"amount"`
	DueDate        time.Time `json:"due_date"`
	CustomerEmail  string    `json:"customer_email"`
}

func NewInvoiceCreatedEvent(tenantID, invoiceID, subscriptionID uuid.UUID, number string, amount int64, dueDate time.Time, customerEmail string) *InvoiceCreatedEvent {
	return &InvoiceCreatedEvent{
		BaseEvent:      NewBaseEvent("invoice.created", tenantID, nil),
		InvoiceID:      invoiceID,
		SubscriptionID: subscriptionID,
		Number:         number,
		Amount:         amount,
		DueDate:        dueDate,
		CustomerEmail:  customerEmail,
	}
}

// InvoicePaidEvent evento de pagamento de fatura
type InvoicePaidEvent struct {
	BaseEvent
	InvoiceID      uuid.UUID `json:"invoice_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PaymentID      uuid.UUID `json:"payment_id"`
	Number         string    `json:"number"`
	Amount         int64     `json:"amount"`
	PaidAt         time.Time `json:"paid_at"`
}

func NewInvoicePaidEvent(tenantID, invoiceID, subscriptionID, paymentID uuid.UUID, number string, amount int64, paidAt time.Time) *InvoicePaidEvent {
	return &InvoicePaidEvent{
		BaseEvent:      NewBaseEvent("invoice.paid", tenantID, nil),
		InvoiceID:      invoiceID,
		SubscriptionID: subscriptionID,
		PaymentID:      paymentID,
		Number:         number,
		Amount:         amount,
		PaidAt:         paidAt,
	}
}

// InvoiceOverdueEvent evento de fatura vencida
type InvoiceOverdueEvent struct {
	BaseEvent
	InvoiceID      uuid.UUID `json:"invoice_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Number         string    `json:"number"`
	Amount         int64     `json:"amount"`
	DueDate        time.Time `json:"due_date"`
	DaysOverdue    int       `json:"days_overdue"`
	CustomerEmail  string    `json:"customer_email"`
}

func NewInvoiceOverdueEvent(tenantID, invoiceID, subscriptionID uuid.UUID, number string, amount int64, dueDate time.Time, daysOverdue int, customerEmail string) *InvoiceOverdueEvent {
	return &InvoiceOverdueEvent{
		BaseEvent:      NewBaseEvent("invoice.overdue", tenantID, nil),
		InvoiceID:      invoiceID,
		SubscriptionID: subscriptionID,
		Number:         number,
		Amount:         amount,
		DueDate:        dueDate,
		DaysOverdue:    daysOverdue,
		CustomerEmail:  customerEmail,
	}
}

// NFEIssuedEvent evento de emissão de NF-e
type NFEIssuedEvent struct {
	BaseEvent
	InvoiceID      uuid.UUID `json:"invoice_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	NFENumber      string    `json:"nfe_number"`
	NFEKey         string    `json:"nfe_key"`
	NFEURL         string    `json:"nfe_url"`
	Amount         int64     `json:"amount"`
	CustomerEmail  string    `json:"customer_email"`
}

func NewNFEIssuedEvent(tenantID, invoiceID, subscriptionID uuid.UUID, nfeNumber, nfeKey, nfeURL string, amount int64, customerEmail string) *NFEIssuedEvent {
	return &NFEIssuedEvent{
		BaseEvent:      NewBaseEvent("nfe.issued", tenantID, nil),
		InvoiceID:      invoiceID,
		SubscriptionID: subscriptionID,
		NFENumber:      nfeNumber,
		NFEKey:         nfeKey,
		NFEURL:         nfeURL,
		Amount:         amount,
		CustomerEmail:  customerEmail,
	}
}

// Eventos de Cliente

// CustomerCreatedEvent evento de criação de cliente
type CustomerCreatedEvent struct {
	BaseEvent
	CustomerID   uuid.UUID `json:"customer_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
}

func NewCustomerCreatedEvent(tenantID, customerID uuid.UUID, name, email, document string, documentType DocumentType) *CustomerCreatedEvent {
	return &CustomerCreatedEvent{
		BaseEvent:    NewBaseEvent("customer.created", tenantID, nil),
		CustomerID:   customerID,
		Name:         name,
		Email:        email,
		Document:     document,
		DocumentType: string(documentType),
	}
}

// CustomerUpdatedEvent evento de atualização de cliente
type CustomerUpdatedEvent struct {
	BaseEvent
	CustomerID uuid.UUID `json:"customer_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Changes    []string  `json:"changes"`
}

func NewCustomerUpdatedEvent(tenantID, customerID uuid.UUID, name, email string, changes []string) *CustomerUpdatedEvent {
	return &CustomerUpdatedEvent{
		BaseEvent:  NewBaseEvent("customer.updated", tenantID, nil),
		CustomerID: customerID,
		Name:       name,
		Email:      email,
		Changes:    changes,
	}
}

// Eventos de Plano

// PlanCreatedEvent evento de criação de plano
type PlanCreatedEvent struct {
	BaseEvent
	PlanID       uuid.UUID `json:"plan_id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	PriceMonthly int64     `json:"price_monthly"`
	PriceYearly  int64     `json:"price_yearly"`
}

func NewPlanCreatedEvent(tenantID, planID uuid.UUID, name, displayName string, priceMonthly, priceYearly int64) *PlanCreatedEvent {
	return &PlanCreatedEvent{
		BaseEvent:    NewBaseEvent("plan.created", tenantID, nil),
		PlanID:       planID,
		Name:         name,
		DisplayName:  displayName,
		PriceMonthly: priceMonthly,
		PriceYearly:  priceYearly,
	}
}

// PlanUpdatedEvent evento de atualização de plano
type PlanUpdatedEvent struct {
	BaseEvent
	PlanID       uuid.UUID `json:"plan_id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	PriceMonthly int64     `json:"price_monthly"`
	PriceYearly  int64     `json:"price_yearly"`
	Changes      []string  `json:"changes"`
}

func NewPlanUpdatedEvent(tenantID, planID uuid.UUID, name, displayName string, priceMonthly, priceYearly int64, changes []string) *PlanUpdatedEvent {
	return &PlanUpdatedEvent{
		BaseEvent:    NewBaseEvent("plan.updated", tenantID, nil),
		PlanID:       planID,
		Name:         name,
		DisplayName:  displayName,
		PriceMonthly: priceMonthly,
		PriceYearly:  priceYearly,
		Changes:      changes,
	}
}

// EventBus interface para publicação de eventos
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, eventType string, handler EventHandler) error
}

// EventHandler interface para manipulação de eventos
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventHandlerFunc função que implementa EventHandler
type EventHandlerFunc func(ctx context.Context, event Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// EventStore interface para armazenamento de eventos
type EventStore interface {
	Save(ctx context.Context, event Event) error
	Load(ctx context.Context, aggregateID uuid.UUID) ([]Event, error)
	LoadByType(ctx context.Context, eventType string) ([]Event, error)
	LoadByTenant(ctx context.Context, tenantID uuid.UUID) ([]Event, error)
}