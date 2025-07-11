package domain

import (
	"time"

	"github.com/google/uuid"
)

// Subscription representa uma assinatura de um tenant
type Subscription struct {
	ID                     uuid.UUID    `json:"id"`
	TenantID               uuid.UUID    `json:"tenant_id"`
	PlanID                 uuid.UUID    `json:"plan_id"`
	Status                 SubscriptionStatus `json:"status"`
	BillingCycle           BillingCycle `json:"billing_cycle"`
	
	// Dados do período
	TrialStartDate         *time.Time   `json:"trial_start_date"`
	TrialEndDate           *time.Time   `json:"trial_end_date"`
	CurrentPeriodStart     time.Time    `json:"current_period_start"`
	CurrentPeriodEnd       time.Time    `json:"current_period_end"`
	
	// Dados de pagamento
	Amount                 int64        `json:"amount"`          // em centavos
	PaymentMethod          PaymentMethod `json:"payment_method"`
	
	// Integrações externas
	AsaasSubscriptionID    *string      `json:"asaas_subscription_id"`
	AsaasCustomerID        *string      `json:"asaas_customer_id"`
	
	// Controle de cobrança
	NextBillingDate        *time.Time   `json:"next_billing_date"`
	RetryCount             int          `json:"retry_count"`
	LastPaymentAttempt     *time.Time   `json:"last_payment_attempt"`
	LastSuccessfulPayment  *time.Time   `json:"last_successful_payment"`
	
	// Cancelamento
	CancelledAt            *time.Time   `json:"cancelled_at"`
	CancelReason           *string      `json:"cancel_reason"`
	CancelledBy            *uuid.UUID   `json:"cancelled_by"`
	
	// Controle de versão
	Version                int          `json:"version"`
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
}

// SubscriptionStatus enumeração dos status da assinatura
type SubscriptionStatus string

const (
	SubscriptionStatusTrial            SubscriptionStatus = "trial"
	SubscriptionStatusActive           SubscriptionStatus = "active"
	SubscriptionStatusPastDue          SubscriptionStatus = "past_due"
	SubscriptionStatusSuspended        SubscriptionStatus = "suspended"
	SubscriptionStatusCancelled        SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired          SubscriptionStatus = "expired"
	SubscriptionStatusPaymentPending   SubscriptionStatus = "payment_pending"
)

// PaymentMethod enumeração dos métodos de pagamento
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodBoleto     PaymentMethod = "boleto"
	PaymentMethodBitcoin    PaymentMethod = "bitcoin"
	PaymentMethodXRP        PaymentMethod = "xrp"
	PaymentMethodXLM        PaymentMethod = "xlm"
	PaymentMethodXDC        PaymentMethod = "xdc"
	PaymentMethodCardano    PaymentMethod = "cardano"
	PaymentMethodHBAR       PaymentMethod = "hbar"
	PaymentMethodXCN        PaymentMethod = "xcn"
	PaymentMethodEthereum   PaymentMethod = "ethereum"
	PaymentMethodSolana     PaymentMethod = "solana"
)

// NewSubscription cria uma nova assinatura
func NewSubscription(tenantID, planID uuid.UUID, cycle BillingCycle, amount int64, method PaymentMethod) *Subscription {
	now := time.Now()
	
	subscription := &Subscription{
		ID:                 uuid.New(),
		TenantID:           tenantID,
		PlanID:             planID,
		Status:             SubscriptionStatusTrial,
		BillingCycle:       cycle,
		Amount:             amount,
		PaymentMethod:      method,
		CurrentPeriodStart: now,
		Version:            1,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	// Configurar período de trial
	subscription.StartTrial(15) // 15 dias padrão
	
	return subscription
}

// StartTrial inicia o período de trial
func (s *Subscription) StartTrial(days int) {
	now := time.Now()
	s.TrialStartDate = &now
	trialEnd := now.AddDate(0, 0, days)
	s.TrialEndDate = &trialEnd
	s.CurrentPeriodEnd = trialEnd
	s.Status = SubscriptionStatusTrial
	s.UpdatedAt = now
}

// EndTrial finaliza o período de trial
func (s *Subscription) EndTrial() {
	if s.IsInTrial() {
		s.Status = SubscriptionStatusPaymentPending
		s.UpdatedAt = time.Now()
	}
}

// Activate ativa a assinatura
func (s *Subscription) Activate() {
	s.Status = SubscriptionStatusActive
	s.LastSuccessfulPayment = &time.Time{}
	*s.LastSuccessfulPayment = time.Now()
	s.UpdatedAt = time.Now()
}

// Suspend suspende a assinatura
func (s *Subscription) Suspend(reason string) {
	s.Status = SubscriptionStatusSuspended
	s.CancelReason = &reason
	s.UpdatedAt = time.Now()
}

// Cancel cancela a assinatura
func (s *Subscription) Cancel(reason string, cancelledBy uuid.UUID) {
	now := time.Now()
	s.Status = SubscriptionStatusCancelled
	s.CancelledAt = &now
	s.CancelReason = &reason
	s.CancelledBy = &cancelledBy
	s.UpdatedAt = now
}

// MarkPaymentFailed marca uma tentativa de pagamento como falha
func (s *Subscription) MarkPaymentFailed() {
	s.RetryCount++
	now := time.Now()
	s.LastPaymentAttempt = &now
	s.UpdatedAt = now
	
	// Após 3 tentativas, marcar como em atraso
	if s.RetryCount >= 3 {
		s.Status = SubscriptionStatusPastDue
	}
}

// MarkPaymentSuccess marca um pagamento como bem-sucedido
func (s *Subscription) MarkPaymentSuccess() {
	s.RetryCount = 0
	now := time.Now()
	s.LastSuccessfulPayment = &now
	s.LastPaymentAttempt = &now
	s.Status = SubscriptionStatusActive
	s.UpdatedAt = now
	
	// Atualizar próximo período
	s.SetNextBillingPeriod()
}

// SetNextBillingPeriod define o próximo período de cobrança
func (s *Subscription) SetNextBillingPeriod() {
	var nextPeriod time.Time
	
	switch s.BillingCycle {
	case BillingCycleMonthly:
		nextPeriod = s.CurrentPeriodEnd.AddDate(0, 1, 0)
	case BillingCycleYearly:
		nextPeriod = s.CurrentPeriodEnd.AddDate(1, 0, 0)
	default:
		nextPeriod = s.CurrentPeriodEnd.AddDate(0, 1, 0)
	}
	
	s.CurrentPeriodStart = s.CurrentPeriodEnd
	s.CurrentPeriodEnd = nextPeriod
	s.NextBillingDate = &nextPeriod
	s.UpdatedAt = time.Now()
}

// IsInTrial verifica se a assinatura está em período de trial
func (s *Subscription) IsInTrial() bool {
	return s.Status == SubscriptionStatusTrial && 
		   s.TrialEndDate != nil && 
		   time.Now().Before(*s.TrialEndDate)
}

// IsActive verifica se a assinatura está ativa
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive || s.IsInTrial()
}

// IsCancelled verifica se a assinatura está cancelada
func (s *Subscription) IsCancelled() bool {
	return s.Status == SubscriptionStatusCancelled
}

// IsExpired verifica se a assinatura está expirada
func (s *Subscription) IsExpired() bool {
	return s.Status == SubscriptionStatusExpired || 
		   (s.TrialEndDate != nil && time.Now().After(*s.TrialEndDate) && s.Status == SubscriptionStatusTrial)
}

// ShouldBeBilled verifica se a assinatura deve ser cobrada
func (s *Subscription) ShouldBeBilled() bool {
	return s.NextBillingDate != nil && 
		   time.Now().After(*s.NextBillingDate) && 
		   s.IsActive()
}

// GetDaysUntilExpiry retorna quantos dias restam até a expiração
func (s *Subscription) GetDaysUntilExpiry() int {
	if s.IsInTrial() && s.TrialEndDate != nil {
		diff := s.TrialEndDate.Sub(time.Now())
		return int(diff.Hours() / 24)
	}
	
	diff := s.CurrentPeriodEnd.Sub(time.Now())
	return int(diff.Hours() / 24)
}

// GetRemainingValue retorna o valor restante da assinatura (prorrata)
func (s *Subscription) GetRemainingValue() int64 {
	if !s.IsActive() {
		return 0
	}
	
	totalDays := s.CurrentPeriodEnd.Sub(s.CurrentPeriodStart).Hours() / 24
	remainingDays := s.CurrentPeriodEnd.Sub(time.Now()).Hours() / 24
	
	if remainingDays <= 0 {
		return 0
	}
	
	return int64(float64(s.Amount) * (remainingDays / totalDays))
}

// CanUpgrade verifica se pode fazer upgrade
func (s *Subscription) CanUpgrade() bool {
	return s.IsActive() && s.Status != SubscriptionStatusPastDue
}

// CanDowngrade verifica se pode fazer downgrade
func (s *Subscription) CanDowngrade() bool {
	return s.IsActive() && s.Status != SubscriptionStatusPastDue
}

// IsCrypto verifica se o método de pagamento é criptomoeda
func (s *Subscription) IsCrypto() bool {
	cryptoMethods := []PaymentMethod{
		PaymentMethodBitcoin,
		PaymentMethodXRP,
		PaymentMethodXLM,
		PaymentMethodXDC,
		PaymentMethodCardano,
		PaymentMethodHBAR,
		PaymentMethodXCN,
		PaymentMethodEthereum,
		PaymentMethodSolana,
	}
	
	for _, method := range cryptoMethods {
		if s.PaymentMethod == method {
			return true
		}
	}
	return false
}