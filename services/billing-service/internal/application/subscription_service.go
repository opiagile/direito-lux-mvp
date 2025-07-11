package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/billing-service/internal/domain"
)

// SubscriptionService serviço de aplicação para assinaturas
type SubscriptionService struct {
	subscriptionRepo domain.SubscriptionRepository
	planRepo         domain.PlanRepository
	customerRepo     domain.CustomerRepository
	paymentGateway   PaymentGatewayInterface
	eventBus         domain.EventBus
}

// NewSubscriptionService cria um novo serviço de assinatura
func NewSubscriptionService(
	subscriptionRepo domain.SubscriptionRepository,
	planRepo domain.PlanRepository,
	customerRepo domain.CustomerRepository,
	paymentGateway PaymentGatewayInterface,
	eventBus domain.EventBus,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		customerRepo:     customerRepo,
		paymentGateway:   paymentGateway,
		eventBus:         eventBus,
	}
}

// CreateSubscriptionCommand comando para criar assinatura
type CreateSubscriptionCommand struct {
	TenantID         uuid.UUID            `json:"tenant_id"`
	PlanID           uuid.UUID            `json:"plan_id"`
	BillingCycle     domain.BillingCycle  `json:"billing_cycle"`
	PaymentMethod    domain.PaymentMethod `json:"payment_method"`
	CustomerName     string               `json:"customer_name"`
	CustomerEmail    string               `json:"customer_email"`
	CustomerDocument string               `json:"customer_document"`
	CustomerPhone    *string              `json:"customer_phone"`
	CustomerAddress  *domain.Address      `json:"customer_address"`
}

// CreateSubscription cria uma nova assinatura
func (s *SubscriptionService) CreateSubscription(ctx context.Context, cmd CreateSubscriptionCommand) (*domain.Subscription, error) {
	// Validar se já existe assinatura ativa para o tenant
	existingSubscription, err := s.subscriptionRepo.GetCurrentByTenantID(ctx, cmd.TenantID)
	if err == nil && existingSubscription != nil && existingSubscription.IsActive() {
		return nil, errors.New("tenant already has an active subscription")
	}

	// Buscar o plano
	plan, err := s.planRepo.GetByID(ctx, cmd.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	if !plan.Active {
		return nil, errors.New("plan is not active")
	}

	// Criar ou buscar cliente
	customer, err := s.getOrCreateCustomer(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Calcular valor baseado no ciclo de cobrança
	amount := plan.GetPrice(cmd.BillingCycle)

	// Criar assinatura
	subscription := domain.NewSubscription(
		cmd.TenantID,
		cmd.PlanID,
		cmd.BillingCycle,
		amount,
		cmd.PaymentMethod,
	)

	// Salvar assinatura
	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Publicar evento
	event := domain.NewSubscriptionCreatedEvent(
		cmd.TenantID,
		subscription.ID,
		plan.ID,
		plan.DisplayName,
		amount,
		cmd.BillingCycle,
		cmd.PaymentMethod,
		plan.TrialDays,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha a criação
		// TODO: implementar retry ou dead letter queue
	}

	return subscription, nil
}

// getOrCreateCustomer busca ou cria um cliente
func (s *SubscriptionService) getOrCreateCustomer(ctx context.Context, cmd CreateSubscriptionCommand) (*domain.Customer, error) {
	// Tentar buscar cliente existente por documento
	customer, err := s.customerRepo.GetByDocument(ctx, cmd.CustomerDocument)
	if err == nil && customer != nil {
		return customer, nil
	}

	// Determinar tipo de documento
	documentType := domain.DocumentTypeCPF
	if len(cmd.CustomerDocument) == 14 {
		documentType = domain.DocumentTypeCNPJ
	}

	// Criar novo cliente
	customer = domain.NewCustomer(
		cmd.TenantID,
		cmd.CustomerName,
		cmd.CustomerEmail,
		cmd.CustomerDocument,
		documentType,
	)

	if cmd.CustomerPhone != nil {
		customer.SetPhone(*cmd.CustomerPhone)
	}

	if cmd.CustomerAddress != nil {
		customer.SetAddress(cmd.CustomerAddress)
	}

	// Salvar cliente
	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// ActivateSubscription ativa uma assinatura após pagamento
func (s *SubscriptionService) ActivateSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if subscription.IsActive() {
		return errors.New("subscription is already active")
	}

	// Ativar assinatura
	subscription.Activate()

	// Salvar alterações
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Buscar plano para evento
	plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return fmt.Errorf("failed to get plan: %w", err)
	}

	// Publicar evento
	event := domain.NewSubscriptionActivatedEvent(
		subscription.TenantID,
		subscription.ID,
		plan.DisplayName,
		subscription.Amount,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha a ativação
	}

	return nil
}

// CancelSubscription cancela uma assinatura
func (s *SubscriptionService) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, reason string, cancelledBy uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if subscription.IsCancelled() {
		return errors.New("subscription is already cancelled")
	}

	// Cancelar assinatura
	subscription.Cancel(reason, cancelledBy)

	// Salvar alterações
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Buscar plano para evento
	plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return fmt.Errorf("failed to get plan: %w", err)
	}

	// Publicar evento
	event := domain.NewSubscriptionCancelledEvent(
		subscription.TenantID,
		subscription.ID,
		plan.DisplayName,
		reason,
		cancelledBy,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha o cancelamento
	}

	return nil
}

// ChangeSubscriptionPlan muda o plano da assinatura
func (s *SubscriptionService) ChangeSubscriptionPlan(ctx context.Context, subscriptionID, newPlanID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if !subscription.CanUpgrade() {
		return errors.New("subscription cannot be upgraded")
	}

	// Buscar novo plano
	newPlan, err := s.planRepo.GetByID(ctx, newPlanID)
	if err != nil {
		return fmt.Errorf("failed to get new plan: %w", err)
	}

	if !newPlan.Active {
		return errors.New("new plan is not active")
	}

	// Calcular prorrata se necessário
	remainingValue := subscription.GetRemainingValue()
	newAmount := newPlan.GetPrice(subscription.BillingCycle)

	// Atualizar assinatura
	subscription.PlanID = newPlan.ID
	subscription.Amount = newAmount
	subscription.UpdatedAt = time.Now()

	// Salvar alterações
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Se há diferença de valor, criar cobrança/crédito
	if newAmount > remainingValue {
		// Criar cobrança complementar
		additionalAmount := newAmount - remainingValue
		// TODO: implementar cobrança complementar
		_ = additionalAmount
	}

	return nil
}

// GetSubscription busca uma assinatura por ID
func (s *SubscriptionService) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*domain.Subscription, error) {
	return s.subscriptionRepo.GetByID(ctx, subscriptionID)
}

// GetSubscriptionByTenant busca assinatura por tenant
func (s *SubscriptionService) GetSubscriptionByTenant(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	return s.subscriptionRepo.GetCurrentByTenantID(ctx, tenantID)
}

// ListSubscriptions lista assinaturas com paginação
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, query domain.SubscriptionQuery) ([]*domain.Subscription, error) {
	// TODO: implementar busca com query
	return nil, errors.New("not implemented")
}

// GetSubscriptionStats obtém estatísticas de assinaturas
func (s *SubscriptionService) GetSubscriptionStats(ctx context.Context) (*domain.SubscriptionStats, error) {
	return s.subscriptionRepo.GetStats(ctx)
}

// ProcessExpiredTrials processa trials expirados
func (s *SubscriptionService) ProcessExpiredTrials(ctx context.Context) error {
	// Buscar assinaturas em trial que expiraram
	expiredTrials, err := s.subscriptionRepo.GetByStatus(ctx, domain.SubscriptionStatusTrial)
	if err != nil {
		return fmt.Errorf("failed to get trial subscriptions: %w", err)
	}

	for _, subscription := range expiredTrials {
		if subscription.IsExpired() {
			// Finalizar trial
			subscription.EndTrial()

			// Salvar alterações
			if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
				// Log erro mas continua processando
				continue
			}

			// Buscar plano para evento
			plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
			if err != nil {
				continue
			}

			// Publicar evento
			event := domain.NewSubscriptionExpiredEvent(
				subscription.TenantID,
				subscription.ID,
				plan.DisplayName,
				subscription.LastSuccessfulPayment,
			)

			if err := s.eventBus.Publish(ctx, event); err != nil {
				// Log erro mas continua processando
			}
		}
	}

	return nil
}

// ProcessTrialEndingNotifications processa notificações de fim de trial
func (s *SubscriptionService) ProcessTrialEndingNotifications(ctx context.Context) error {
	// Buscar assinaturas em trial que expiram em 3 dias
	expiring, err := s.subscriptionRepo.GetExpiring(ctx, 3)
	if err != nil {
		return fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}

	for _, subscription := range expiring {
		if subscription.IsInTrial() && subscription.TrialEndDate != nil {
			daysRemaining := subscription.GetDaysUntilExpiry()
			
			// Buscar plano para evento
			plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
			if err != nil {
				continue
			}

			// Publicar evento
			event := domain.NewSubscriptionTrialEndingEvent(
				subscription.TenantID,
				subscription.ID,
				plan.DisplayName,
				daysRemaining,
				*subscription.TrialEndDate,
			)

			if err := s.eventBus.Publish(ctx, event); err != nil {
				// Log erro mas continua processando
			}
		}
	}

	return nil
}

// ProcessBillingCycle processa ciclo de cobrança
func (s *SubscriptionService) ProcessBillingCycle(ctx context.Context) error {
	// Buscar assinaturas pendentes de cobrança
	pendingBilling, err := s.subscriptionRepo.GetPendingBilling(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending billing subscriptions: %w", err)
	}

	for _, subscription := range pendingBilling {
		if subscription.ShouldBeBilled() {
			// Processar cobrança
			if err := s.processBilling(ctx, subscription); err != nil {
				// Log erro mas continua processando
				continue
			}
		}
	}

	return nil
}

// processBilling processa cobrança de uma assinatura
func (s *SubscriptionService) processBilling(ctx context.Context, subscription *domain.Subscription) error {
	// Buscar cliente
	customer, err := s.customerRepo.GetByTenantID(ctx, subscription.TenantID)
	if err != nil || len(customer) == 0 {
		return fmt.Errorf("customer not found for tenant: %s", subscription.TenantID)
	}

	// Criar cobrança via gateway
	paymentRequest := &PaymentRequest{
		SubscriptionID: subscription.ID,
		TenantID:       subscription.TenantID,
		Amount:         subscription.Amount,
		PaymentMethod:  subscription.PaymentMethod,
		CustomerData:   customer[0],
	}

	payment, err := s.paymentGateway.CreatePayment(ctx, paymentRequest)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// Atualizar próximo ciclo de cobrança
	subscription.SetNextBillingPeriod()

	// Salvar alterações
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Publicar evento de pagamento criado
	event := domain.NewPaymentCreatedEvent(
		subscription.TenantID,
		payment.ID,
		subscription.ID,
		subscription.Amount,
		subscription.PaymentMethod,
		"BRL",
		payment.DueDate,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha o processamento
	}

	return nil
}

// PaymentGatewayInterface interface para gateway de pagamento
type PaymentGatewayInterface interface {
	CreatePayment(ctx context.Context, request *PaymentRequest) (*domain.Payment, error)
	ProcessPayment(ctx context.Context, paymentID uuid.UUID) error
	RefundPayment(ctx context.Context, paymentID uuid.UUID, amount int64) error
	GetPaymentStatus(ctx context.Context, paymentID uuid.UUID) (domain.PaymentStatus, error)
}

// PaymentRequest estrutura de requisição de pagamento
type PaymentRequest struct {
	SubscriptionID uuid.UUID           `json:"subscription_id"`
	TenantID       uuid.UUID           `json:"tenant_id"`
	Amount         int64               `json:"amount"`
	PaymentMethod  domain.PaymentMethod `json:"payment_method"`
	CustomerData   *domain.Customer    `json:"customer_data"`
}