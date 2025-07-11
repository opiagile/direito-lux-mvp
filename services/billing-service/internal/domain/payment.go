package domain

import (
	"time"

	"github.com/google/uuid"
)

// Payment representa um pagamento
type Payment struct {
	ID                uuid.UUID     `json:"id"`
	SubscriptionID    uuid.UUID     `json:"subscription_id"`
	TenantID          uuid.UUID     `json:"tenant_id"`
	InvoiceID         *uuid.UUID    `json:"invoice_id"`
	
	// Dados do pagamento
	Amount            int64         `json:"amount"`        // em centavos
	Currency          string        `json:"currency"`      // BRL, BTC, etc.
	PaymentMethod     PaymentMethod `json:"payment_method"`
	Status            PaymentStatus `json:"status"`
	
	// Integrações externas
	AsaasPaymentID    *string       `json:"asaas_payment_id"`
	AsaasChargeID     *string       `json:"asaas_charge_id"`
	NOWPaymentID      *string       `json:"now_payment_id"`
	
	// Dados específicos do pagamento
	GatewayResponse   *string       `json:"gateway_response"`
	GatewayReference  *string       `json:"gateway_reference"`
	TransactionID     *string       `json:"transaction_id"`
	
	// Dados de cobrança
	DueDate           *time.Time    `json:"due_date"`
	PaidAt            *time.Time    `json:"paid_at"`
	
	// Controle de retry
	RetryCount        int           `json:"retry_count"`
	LastAttemptAt     *time.Time    `json:"last_attempt_at"`
	NextRetryAt       *time.Time    `json:"next_retry_at"`
	
	// Cancelamento/Estorno
	CancelledAt       *time.Time    `json:"cancelled_at"`
	CancelReason      *string       `json:"cancel_reason"`
	RefundedAt        *time.Time    `json:"refunded_at"`
	RefundAmount      *int64        `json:"refund_amount"`
	
	// Cripto específico
	CryptoAddress     *string       `json:"crypto_address"`
	CryptoAmount      *string       `json:"crypto_amount"`
	CryptoTxHash      *string       `json:"crypto_tx_hash"`
	ExchangeRate      *float64      `json:"exchange_rate"`
	
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// PaymentStatus enumeração dos status do pagamento
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusExpired   PaymentStatus = "expired"
	PaymentStatusPartial   PaymentStatus = "partial"
)

// NewPayment cria um novo pagamento
func NewPayment(subscriptionID, tenantID uuid.UUID, amount int64, method PaymentMethod, currency string) *Payment {
	now := time.Now()
	
	payment := &Payment{
		ID:             uuid.New(),
		SubscriptionID: subscriptionID,
		TenantID:       tenantID,
		Amount:         amount,
		Currency:       currency,
		PaymentMethod:  method,
		Status:         PaymentStatusPending,
		RetryCount:     0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	
	// Definir data de vencimento baseada no método
	switch method {
	case PaymentMethodBoleto:
		dueDate := now.AddDate(0, 0, 7) // 7 dias para boleto
		payment.DueDate = &dueDate
	case PaymentMethodPix:
		dueDate := now.AddDate(0, 0, 1) // 1 dia para PIX
		payment.DueDate = &dueDate
	default:
		// Cartão e cripto processam imediatamente
		payment.DueDate = &now
	}
	
	return payment
}

// MarkAsPaid marca o pagamento como pago
func (p *Payment) MarkAsPaid(transactionID string) {
	now := time.Now()
	p.Status = PaymentStatusPaid
	p.PaidAt = &now
	p.TransactionID = &transactionID
	p.UpdatedAt = now
}

// MarkAsFailed marca o pagamento como falha
func (p *Payment) MarkAsFailed(reason string) {
	now := time.Now()
	p.Status = PaymentStatusFailed
	p.CancelReason = &reason
	p.LastAttemptAt = &now
	p.RetryCount++
	p.UpdatedAt = now
	
	// Definir próxima tentativa (backoff exponencial)
	nextRetry := now.Add(time.Duration(p.RetryCount*p.RetryCount) * time.Hour)
	p.NextRetryAt = &nextRetry
}

// MarkAsCancelled marca o pagamento como cancelado
func (p *Payment) MarkAsCancelled(reason string) {
	now := time.Now()
	p.Status = PaymentStatusCancelled
	p.CancelledAt = &now
	p.CancelReason = &reason
	p.UpdatedAt = now
}

// MarkAsRefunded marca o pagamento como reembolsado
func (p *Payment) MarkAsRefunded(refundAmount int64) {
	now := time.Now()
	p.Status = PaymentStatusRefunded
	p.RefundedAt = &now
	p.RefundAmount = &refundAmount
	p.UpdatedAt = now
}

// SetAsaasData define os dados do ASAAS
func (p *Payment) SetAsaasData(paymentID, chargeID string) {
	p.AsaasPaymentID = &paymentID
	p.AsaasChargeID = &chargeID
	p.UpdatedAt = time.Now()
}

// SetNOWPaymentData define os dados do NOWPayments
func (p *Payment) SetNOWPaymentData(paymentID string) {
	p.NOWPaymentID = &paymentID
	p.UpdatedAt = time.Now()
}

// SetCryptoData define os dados específicos de criptomoedas
func (p *Payment) SetCryptoData(address, amount, txHash string, exchangeRate float64) {
	p.CryptoAddress = &address
	p.CryptoAmount = &amount
	p.CryptoTxHash = &txHash
	p.ExchangeRate = &exchangeRate
	p.UpdatedAt = time.Now()
}

// IsExpired verifica se o pagamento está expirado
func (p *Payment) IsExpired() bool {
	return p.DueDate != nil && time.Now().After(*p.DueDate) && p.Status == PaymentStatusPending
}

// CanRetry verifica se pode tentar o pagamento novamente
func (p *Payment) CanRetry() bool {
	return p.Status == PaymentStatusFailed && 
		   p.RetryCount < 3 && 
		   (p.NextRetryAt == nil || time.Now().After(*p.NextRetryAt))
}

// IsSuccessful verifica se o pagamento foi bem-sucedido
func (p *Payment) IsSuccessful() bool {
	return p.Status == PaymentStatusPaid
}

// IsPending verifica se o pagamento está pendente
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsCrypto verifica se é um pagamento em criptomoeda
func (p *Payment) IsCrypto() bool {
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
		if p.PaymentMethod == method {
			return true
		}
	}
	return false
}

// GetFormattedAmount retorna o valor formatado
func (p *Payment) GetFormattedAmount() float64 {
	if p.Currency == "BRL" {
		return float64(p.Amount) / 100
	}
	return float64(p.Amount)
}

// GetGateway retorna qual gateway deve ser usado
func (p *Payment) GetGateway() string {
	if p.IsCrypto() {
		return "nowpayments"
	}
	return "asaas"
}

// ExpirePayment marca o pagamento como expirado
func (p *Payment) ExpirePayment() {
	if p.IsExpired() {
		p.Status = PaymentStatusExpired
		p.UpdatedAt = time.Now()
	}
}

// GetDaysUntilDue retorna quantos dias restam até o vencimento
func (p *Payment) GetDaysUntilDue() int {
	if p.DueDate == nil {
		return 0
	}
	
	diff := p.DueDate.Sub(time.Now())
	days := int(diff.Hours() / 24)
	
	if days < 0 {
		return 0
	}
	
	return days
}

// GetRetryDelay retorna o delay para próxima tentativa
func (p *Payment) GetRetryDelay() time.Duration {
	// Backoff exponencial: 1h, 4h, 9h
	delay := time.Duration(p.RetryCount*p.RetryCount) * time.Hour
	if delay > 24*time.Hour {
		delay = 24 * time.Hour
	}
	return delay
}