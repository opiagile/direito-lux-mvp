package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/billing-service/internal/domain"
)

// PaymentService serviço de aplicação para pagamentos
type PaymentService struct {
	paymentRepo      domain.PaymentRepository
	subscriptionRepo domain.SubscriptionRepository
	invoiceRepo      domain.InvoiceRepository
	customerRepo     domain.CustomerRepository
	asaasGateway     AsaasGatewayInterface
	cryptoGateway    CryptoGatewayInterface
	eventBus         domain.EventBus
}

// NewPaymentService cria um novo serviço de pagamento
func NewPaymentService(
	paymentRepo domain.PaymentRepository,
	subscriptionRepo domain.SubscriptionRepository,
	invoiceRepo domain.InvoiceRepository,
	customerRepo domain.CustomerRepository,
	asaasGateway AsaasGatewayInterface,
	cryptoGateway CryptoGatewayInterface,
	eventBus domain.EventBus,
) *PaymentService {
	return &PaymentService{
		paymentRepo:      paymentRepo,
		subscriptionRepo: subscriptionRepo,
		invoiceRepo:      invoiceRepo,
		customerRepo:     customerRepo,
		asaasGateway:     asaasGateway,
		cryptoGateway:    cryptoGateway,
		eventBus:         eventBus,
	}
}

// CreatePaymentCommand comando para criar pagamento
type CreatePaymentCommand struct {
	SubscriptionID uuid.UUID           `json:"subscription_id"`
	TenantID       uuid.UUID           `json:"tenant_id"`
	Amount         int64               `json:"amount"`
	PaymentMethod  domain.PaymentMethod `json:"payment_method"`
	Currency       string              `json:"currency"`
	InvoiceID      *uuid.UUID          `json:"invoice_id"`
}

// CreatePayment cria um novo pagamento
func (s *PaymentService) CreatePayment(ctx context.Context, cmd CreatePaymentCommand) (*domain.Payment, error) {
	// Validar assinatura
	subscription, err := s.subscriptionRepo.GetByID(ctx, cmd.SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if !subscription.IsActive() {
		return nil, errors.New("subscription is not active")
	}

	// Criar pagamento
	payment := domain.NewPayment(
		cmd.SubscriptionID,
		cmd.TenantID,
		cmd.Amount,
		cmd.PaymentMethod,
		cmd.Currency,
	)

	if cmd.InvoiceID != nil {
		payment.InvoiceID = cmd.InvoiceID
	}

	// Salvar pagamento
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Processar pagamento via gateway apropriado
	if err := s.processPaymentViaGateway(ctx, payment); err != nil {
		// Marcar pagamento como falha
		payment.MarkAsFailed(err.Error())
		s.paymentRepo.Update(ctx, payment)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Publicar evento
	event := domain.NewPaymentCreatedEvent(
		cmd.TenantID,
		payment.ID,
		cmd.SubscriptionID,
		cmd.Amount,
		cmd.PaymentMethod,
		cmd.Currency,
		payment.DueDate,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha a criação
	}

	return payment, nil
}

// processPaymentViaGateway processa pagamento via gateway apropriado
func (s *PaymentService) processPaymentViaGateway(ctx context.Context, payment *domain.Payment) error {
	// Buscar dados do cliente
	customers, err := s.customerRepo.GetByTenantID(ctx, payment.TenantID)
	if err != nil || len(customers) == 0 {
		return fmt.Errorf("customer not found for tenant: %s", payment.TenantID)
	}
	customer := customers[0]

	if payment.IsCrypto() {
		return s.processCryptoPayment(ctx, payment, customer)
	} else {
		return s.processAsaasPayment(ctx, payment, customer)
	}
}

// processAsaasPayment processa pagamento via ASAAS
func (s *PaymentService) processAsaasPayment(ctx context.Context, payment *domain.Payment, customer *domain.Customer) error {
	// Preparar dados do pagamento
	asaasPayment := &AsaasPaymentRequest{
		CustomerID:    customer.AsaasCustomerID,
		BillingType:   s.mapPaymentMethodToAsaas(payment.PaymentMethod),
		Value:         payment.GetFormattedAmount(),
		DueDate:       payment.DueDate,
		Description:   fmt.Sprintf("Assinatura Direito Lux - %s", payment.ID.String()[:8]),
		ExternalReference: payment.ID.String(),
	}

	// Criar cobrança no ASAAS
	asaasResponse, err := s.asaasGateway.CreateCharge(ctx, asaasPayment)
	if err != nil {
		return fmt.Errorf("failed to create ASAAS charge: %w", err)
	}

	// Atualizar pagamento com dados do ASAAS
	payment.SetAsaasData(asaasResponse.ID, asaasResponse.ID)
	payment.GatewayReference = &asaasResponse.ID

	// Salvar alterações
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// processCryptoPayment processa pagamento via criptomoedas
func (s *PaymentService) processCryptoPayment(ctx context.Context, payment *domain.Payment, customer *domain.Customer) error {
	// Preparar dados do pagamento cripto
	cryptoPayment := &CryptoPaymentRequest{
		Currency:      string(payment.PaymentMethod),
		Amount:        payment.Amount,
		OrderID:       payment.ID.String(),
		Description:   fmt.Sprintf("Assinatura Direito Lux - %s", payment.ID.String()[:8]),
		CustomerEmail: customer.Email,
		CallbackURL:   fmt.Sprintf("https://api.direitolux.com.br/billing/webhooks/crypto/%s", payment.ID.String()),
	}

	// Criar pagamento cripto
	cryptoResponse, err := s.cryptoGateway.CreatePayment(ctx, cryptoPayment)
	if err != nil {
		return fmt.Errorf("failed to create crypto payment: %w", err)
	}

	// Atualizar pagamento com dados cripto
	payment.SetNOWPaymentData(cryptoResponse.PaymentID)
	payment.SetCryptoData(
		cryptoResponse.PaymentAddress,
		cryptoResponse.PaymentAmount,
		"", // TxHash será preenchido após confirmação
		cryptoResponse.ExchangeRate,
	)

	// Salvar alterações
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// mapPaymentMethodToAsaas mapeia método de pagamento para ASAAS
func (s *PaymentService) mapPaymentMethodToAsaas(method domain.PaymentMethod) string {
	switch method {
	case domain.PaymentMethodCreditCard:
		return "CREDIT_CARD"
	case domain.PaymentMethodDebitCard:
		return "DEBIT_CARD"
	case domain.PaymentMethodPix:
		return "PIX"
	case domain.PaymentMethodBoleto:
		return "BOLETO"
	default:
		return "BOLETO"
	}
}

// ProcessPaymentSuccess processa sucesso de pagamento
func (s *PaymentService) ProcessPaymentSuccess(ctx context.Context, paymentID uuid.UUID, transactionID string) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment.IsSuccessful() {
		return errors.New("payment is already successful")
	}

	// Marcar pagamento como pago
	payment.MarkAsPaid(transactionID)

	// Salvar alterações
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Ativar assinatura
	if err := s.activateSubscription(ctx, payment.SubscriptionID); err != nil {
		return fmt.Errorf("failed to activate subscription: %w", err)
	}

	// Marcar fatura como paga se existir
	if payment.InvoiceID != nil {
		if err := s.markInvoiceAsPaid(ctx, *payment.InvoiceID, paymentID); err != nil {
			// Log erro mas não falha o processamento
		}
	}

	// Publicar evento
	event := domain.NewPaymentSuccessEvent(
		payment.TenantID,
		payment.ID,
		payment.SubscriptionID,
		payment.Amount,
		payment.PaymentMethod,
		transactionID,
		time.Now(),
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha o processamento
	}

	return nil
}

// ProcessPaymentFailure processa falha de pagamento
func (s *PaymentService) ProcessPaymentFailure(ctx context.Context, paymentID uuid.UUID, reason string) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Marcar pagamento como falha
	payment.MarkAsFailed(reason)

	// Salvar alterações
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Marcar assinatura como falha de pagamento
	if err := s.markSubscriptionPaymentFailed(ctx, payment.SubscriptionID); err != nil {
		return fmt.Errorf("failed to mark subscription payment failed: %w", err)
	}

	// Publicar evento
	event := domain.NewPaymentFailedEvent(
		payment.TenantID,
		payment.ID,
		payment.SubscriptionID,
		payment.Amount,
		payment.PaymentMethod,
		reason,
		payment.RetryCount,
		payment.NextRetryAt,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha o processamento
	}

	return nil
}

// activateSubscription ativa assinatura após pagamento
func (s *PaymentService) activateSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	subscription.MarkPaymentSuccess()

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// markSubscriptionPaymentFailed marca falha de pagamento na assinatura
func (s *PaymentService) markSubscriptionPaymentFailed(ctx context.Context, subscriptionID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	subscription.MarkPaymentFailed()

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// markInvoiceAsPaid marca fatura como paga
func (s *PaymentService) markInvoiceAsPaid(ctx context.Context, invoiceID, paymentID uuid.UUID) error {
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to get invoice: %w", err)
	}

	invoice.MarkAsPaid(paymentID)

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}

// RetryFailedPayments tenta novamente pagamentos falhos
func (s *PaymentService) RetryFailedPayments(ctx context.Context) error {
	// Buscar pagamentos falhos que podem ser tentados novamente
	failedPayments, err := s.paymentRepo.GetFailedRetriable(ctx)
	if err != nil {
		return fmt.Errorf("failed to get failed payments: %w", err)
	}

	for _, payment := range failedPayments {
		if payment.CanRetry() {
			// Tentar processar novamente
			if err := s.processPaymentViaGateway(ctx, payment); err != nil {
				// Incrementar contador de tentativas
				payment.MarkAsFailed(err.Error())
				s.paymentRepo.Update(ctx, payment)
				continue
			}

			// Atualizar status
			if err := s.paymentRepo.Update(ctx, payment); err != nil {
				// Log erro mas continua processando
			}
		}
	}

	return nil
}

// RefundPayment processa reembolso de pagamento
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uuid.UUID, refundAmount int64, reason string) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if !payment.IsSuccessful() {
		return errors.New("payment is not successful")
	}

	// Processar reembolso via gateway
	if payment.IsCrypto() {
		// Criptomoedas geralmente não permitem reembolso automático
		return errors.New("crypto payments cannot be refunded automatically")
	}

	// Processar reembolso via ASAAS
	if payment.AsaasPaymentID != nil {
		if err := s.asaasGateway.RefundPayment(ctx, *payment.AsaasPaymentID, refundAmount); err != nil {
			return fmt.Errorf("failed to refund payment via ASAAS: %w", err)
		}
	}

	// Marcar pagamento como reembolsado
	payment.MarkAsRefunded(refundAmount)

	// Salvar alterações
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Publicar evento
	event := domain.NewPaymentRefundedEvent(
		payment.TenantID,
		payment.ID,
		payment.SubscriptionID,
		payment.Amount,
		refundAmount,
		reason,
		time.Now(),
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha o processamento
	}

	return nil
}

// GetPayment busca pagamento por ID
func (s *PaymentService) GetPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	return s.paymentRepo.GetByID(ctx, paymentID)
}

// GetPaymentsBySubscription busca pagamentos por assinatura
func (s *PaymentService) GetPaymentsBySubscription(ctx context.Context, subscriptionID uuid.UUID) ([]*domain.Payment, error) {
	return s.paymentRepo.GetBySubscriptionID(ctx, subscriptionID)
}

// GetPaymentsByTenant busca pagamentos por tenant
func (s *PaymentService) GetPaymentsByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Payment, error) {
	return s.paymentRepo.GetByTenantID(ctx, tenantID)
}

// GetPaymentStats obtém estatísticas de pagamentos
func (s *PaymentService) GetPaymentStats(ctx context.Context, tenantID *uuid.UUID) (*domain.PaymentStats, error) {
	return s.paymentRepo.GetStats(ctx, tenantID)
}

// Interfaces dos gateways

// AsaasGatewayInterface interface para gateway ASAAS
type AsaasGatewayInterface interface {
	CreateCharge(ctx context.Context, request *AsaasPaymentRequest) (*AsaasPaymentResponse, error)
	GetCharge(ctx context.Context, chargeID string) (*AsaasPaymentResponse, error)
	RefundPayment(ctx context.Context, paymentID string, amount int64) error
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*AsaasCustomerResponse, error)
}

// CryptoGatewayInterface interface para gateway de criptomoedas
type CryptoGatewayInterface interface {
	CreatePayment(ctx context.Context, request *CryptoPaymentRequest) (*CryptoPaymentResponse, error)
	GetPayment(ctx context.Context, paymentID string) (*CryptoPaymentResponse, error)
	GetSupportedCurrencies(ctx context.Context) ([]string, error)
}

// Estruturas de requisição e resposta

// AsaasPaymentRequest requisição de pagamento ASAAS
type AsaasPaymentRequest struct {
	CustomerID        *string    `json:"customer"`
	BillingType       string     `json:"billingType"`
	Value             float64    `json:"value"`
	DueDate           *time.Time `json:"dueDate"`
	Description       string     `json:"description"`
	ExternalReference string     `json:"externalReference"`
}

// AsaasPaymentResponse resposta de pagamento ASAAS
type AsaasPaymentResponse struct {
	ID                string     `json:"id"`
	Status            string     `json:"status"`
	Value             float64    `json:"value"`
	NetValue          float64    `json:"netValue"`
	DueDate           time.Time  `json:"dueDate"`
	PaymentDate       *time.Time `json:"paymentDate"`
	InvoiceURL        string     `json:"invoiceUrl"`
	BankSlipURL       string     `json:"bankSlipUrl"`
	TransactionReceiptURL string `json:"transactionReceiptUrl"`
}

// AsaasCustomerResponse resposta de cliente ASAAS
type AsaasCustomerResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	CpfCnpj  string `json:"cpfCnpj"`
}

// CryptoPaymentRequest requisição de pagamento cripto
type CryptoPaymentRequest struct {
	Currency      string  `json:"currency"`
	Amount        int64   `json:"amount"`
	OrderID       string  `json:"order_id"`
	Description   string  `json:"description"`
	CustomerEmail string  `json:"customer_email"`
	CallbackURL   string  `json:"callback_url"`
}

// CryptoPaymentResponse resposta de pagamento cripto
type CryptoPaymentResponse struct {
	PaymentID      string  `json:"payment_id"`
	PaymentAddress string  `json:"payment_address"`
	PaymentAmount  string  `json:"payment_amount"`
	ExchangeRate   float64 `json:"exchange_rate"`
	PaymentURL     string  `json:"payment_url"`
	Status         string  `json:"status"`
}