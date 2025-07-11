package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/billing-service/internal/domain"
)

// OnboardingService serviço para onboarding de novos clientes
type OnboardingService struct {
	subscriptionRepo domain.SubscriptionRepository
	planRepo         domain.PlanRepository
	customerRepo     domain.CustomerRepository
	paymentService   *PaymentService
	eventBus         domain.EventBus
}

// NewOnboardingService cria novo serviço de onboarding
func NewOnboardingService(
	subscriptionRepo domain.SubscriptionRepository,
	planRepo domain.PlanRepository,
	customerRepo domain.CustomerRepository,
	paymentService *PaymentService,
	eventBus domain.EventBus,
) *OnboardingService {
	return &OnboardingService{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		customerRepo:     customerRepo,
		paymentService:   paymentService,
		eventBus:         eventBus,
	}
}

// OnboardingData dados para onboarding
type OnboardingData struct {
	TenantID         uuid.UUID            `json:"tenant_id"`
	
	// Dados do cliente
	CustomerName     string               `json:"customer_name"`
	CustomerEmail    string               `json:"customer_email"`
	CustomerPhone    string               `json:"customer_phone"`
	CustomerDocument string               `json:"customer_document"`
	
	// Dados da empresa (opcional)
	CompanyName      string               `json:"company_name"`
	TradingName      string               `json:"trading_name"`
	StateRegistration string              `json:"state_registration"`
	
	// Endereço
	Address          *domain.Address      `json:"address"`
	
	// Plano selecionado
	PlanID           uuid.UUID            `json:"plan_id"`
	BillingCycle     domain.BillingCycle  `json:"billing_cycle"`
	
	// Método de pagamento
	PaymentMethod    domain.PaymentMethod `json:"payment_method"`
	
	// Dados do cartão (se aplicável)
	CardData         *CardData            `json:"card_data"`
	
	// Dados cripto (se aplicável)
	CryptoWallet     *string              `json:"crypto_wallet"`
	
	// Preferências
	AcceptedTerms    bool                 `json:"accepted_terms"`
	AcceptedPrivacy  bool                 `json:"accepted_privacy"`
	NewsletterOptIn  bool                 `json:"newsletter_opt_in"`
	
	// Dados opcionais
	ReferralCode     *string              `json:"referral_code"`
	CouponCode       *string              `json:"coupon_code"`
	UTMSource        *string              `json:"utm_source"`
	UTMCampaign      *string              `json:"utm_campaign"`
}

// CardData dados do cartão de crédito
type CardData struct {
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	HolderName  string `json:"holder_name"`
}

// OnboardingResult resultado do onboarding
type OnboardingResult struct {
	SubscriptionID   uuid.UUID `json:"subscription_id"`
	CustomerID       uuid.UUID `json:"customer_id"`
	PaymentID        *uuid.UUID `json:"payment_id"`
	TrialEndDate     *time.Time `json:"trial_end_date"`
	NextBillingDate  *time.Time `json:"next_billing_date"`
	PaymentURL       *string    `json:"payment_url"`
	BoletoURL        *string    `json:"boleto_url"`
	CryptoAddress    *string    `json:"crypto_address"`
	CryptoAmount     *string    `json:"crypto_amount"`
	Success          bool       `json:"success"`
	Message          string     `json:"message"`
}

// StartOnboarding inicia o processo de onboarding
func (s *OnboardingService) StartOnboarding(ctx context.Context, data OnboardingData) (*OnboardingResult, error) {
	// Validar dados obrigatórios
	if err := s.validateOnboardingData(data); err != nil {
		return &OnboardingResult{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Verificar se já existe assinatura ativa
	existingSubscription, err := s.subscriptionRepo.GetCurrentByTenantID(ctx, data.TenantID)
	if err == nil && existingSubscription != nil && existingSubscription.IsActive() {
		return &OnboardingResult{
			Success: false,
			Message: "Tenant já possui assinatura ativa",
		}, errors.New("tenant already has active subscription")
	}

	// Buscar plano
	plan, err := s.planRepo.GetByID(ctx, data.PlanID)
	if err != nil {
		return &OnboardingResult{
			Success: false,
			Message: "Plano não encontrado",
		}, err
	}

	// Criar cliente
	customer, err := s.createCustomer(ctx, data)
	if err != nil {
		return &OnboardingResult{
			Success: false,
			Message: "Erro ao criar cliente: " + err.Error(),
		}, err
	}

	// Criar assinatura
	subscription, err := s.createSubscription(ctx, data, plan, customer)
	if err != nil {
		return &OnboardingResult{
			Success: false,
			Message: "Erro ao criar assinatura: " + err.Error(),
		}, err
	}

	result := &OnboardingResult{
		SubscriptionID:  subscription.ID,
		CustomerID:      customer.ID,
		TrialEndDate:    subscription.TrialEndDate,
		NextBillingDate: subscription.NextBillingDate,
		Success:         true,
		Message:         "Onboarding concluído com sucesso",
	}

	// Se não é trial gratuito, criar primeiro pagamento
	if !s.isFreeTrial(plan) {
		payment, err := s.createFirstPayment(ctx, subscription, customer, data)
		if err != nil {
			return &OnboardingResult{
				Success: false,
				Message: "Erro ao processar pagamento: " + err.Error(),
			}, err
		}

		result.PaymentID = &payment.ID
		result.PaymentURL = s.getPaymentURL(payment)
		result.BoletoURL = s.getBoletoURL(payment)
		result.CryptoAddress = payment.CryptoAddress
		result.CryptoAmount = payment.CryptoAmount
	}

	return result, nil
}

// validateOnboardingData valida dados de onboarding
func (s *OnboardingService) validateOnboardingData(data OnboardingData) error {
	if data.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}

	if data.CustomerName == "" {
		return errors.New("customer_name é obrigatório")
	}

	if data.CustomerEmail == "" {
		return errors.New("customer_email é obrigatório")
	}

	if data.CustomerDocument == "" {
		return errors.New("customer_document é obrigatório")
	}

	if data.PlanID == uuid.Nil {
		return errors.New("plan_id é obrigatório")
	}

	if !data.AcceptedTerms {
		return errors.New("é necessário aceitar os termos de uso")
	}

	if !data.AcceptedPrivacy {
		return errors.New("é necessário aceitar a política de privacidade")
	}

	// Validar método de pagamento
	if data.PaymentMethod == domain.PaymentMethodCreditCard && data.CardData == nil {
		return errors.New("dados do cartão são obrigatórios para pagamento com cartão")
	}

	// Validar endereço (obrigatório para NF-e)
	if data.Address == nil {
		return errors.New("endereço é obrigatório para emissão de nota fiscal")
	}

	if data.Address.Street == "" || data.Address.City == "" || data.Address.State == "" {
		return errors.New("dados do endereço incompletos")
	}

	return nil
}

// createCustomer cria o cliente
func (s *OnboardingService) createCustomer(ctx context.Context, data OnboardingData) (*domain.Customer, error) {
	// Determinar tipo de documento
	documentType := domain.DocumentTypeCPF
	if len(data.CustomerDocument) == 14 {
		documentType = domain.DocumentTypeCNPJ
	}

	// Criar cliente
	customer := domain.NewCustomer(
		data.TenantID,
		data.CustomerName,
		data.CustomerEmail,
		data.CustomerDocument,
		documentType,
	)

	// Definir telefone
	if data.CustomerPhone != "" {
		customer.SetPhone(data.CustomerPhone)
	}

	// Definir endereço
	if data.Address != nil {
		customer.SetAddress(data.Address)
	}

	// Dados da empresa (se CNPJ)
	if documentType == domain.DocumentTypeCNPJ {
		customer.SetCompanyData(data.CompanyName, data.TradingName, data.StateRegistration)
	}

	// Validar documento
	if !customer.ValidateDocument() {
		return nil, errors.New("documento inválido")
	}

	// Salvar cliente
	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("falha ao criar cliente: %w", err)
	}

	// Publicar evento
	event := domain.NewCustomerCreatedEvent(
		data.TenantID,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Document,
		customer.DocumentType,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha a criação
	}

	return customer, nil
}

// createSubscription cria a assinatura
func (s *OnboardingService) createSubscription(ctx context.Context, data OnboardingData, plan *domain.Plan, customer *domain.Customer) (*domain.Subscription, error) {
	// Calcular valor baseado no ciclo
	amount := plan.GetPrice(data.BillingCycle)

	// Criar assinatura
	subscription := domain.NewSubscription(
		data.TenantID,
		data.PlanID,
		data.BillingCycle,
		amount,
		data.PaymentMethod,
	)

	// Configurar trial
	if plan.TrialDays > 0 {
		subscription.StartTrial(plan.TrialDays)
	}

	// Salvar assinatura
	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("falha ao criar assinatura: %w", err)
	}

	// Publicar evento
	event := domain.NewSubscriptionCreatedEvent(
		data.TenantID,
		subscription.ID,
		plan.ID,
		plan.DisplayName,
		amount,
		data.BillingCycle,
		data.PaymentMethod,
		plan.TrialDays,
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log erro mas não falha a criação
	}

	return subscription, nil
}

// createFirstPayment cria primeiro pagamento (se não for trial gratuito)
func (s *OnboardingService) createFirstPayment(ctx context.Context, subscription *domain.Subscription, customer *domain.Customer, data OnboardingData) (*domain.Payment, error) {
	// Criar comando de pagamento
	cmd := CreatePaymentCommand{
		SubscriptionID: subscription.ID,
		TenantID:       subscription.TenantID,
		Amount:         subscription.Amount,
		PaymentMethod:  subscription.PaymentMethod,
		Currency:       "BRL",
	}

	// Se método cripto, usar moeda específica
	if subscription.IsCrypto() {
		cmd.Currency = string(subscription.PaymentMethod)
	}

	// Criar pagamento
	payment, err := s.paymentService.CreatePayment(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar pagamento: %w", err)
	}

	return payment, nil
}

// isFreeTrial verifica se é trial gratuito
func (s *OnboardingService) isFreeTrial(plan *domain.Plan) bool {
	return plan.TrialDays > 0 && plan.GetPrice(domain.BillingCycleMonthly) > 0
}

// getPaymentURL retorna URL de pagamento
func (s *OnboardingService) getPaymentURL(payment *domain.Payment) *string {
	if payment.AsaasPaymentID != nil {
		url := fmt.Sprintf("https://www.asaas.com/c/%s", *payment.AsaasPaymentID)
		return &url
	}
	return nil
}

// getBoletoURL retorna URL do boleto
func (s *OnboardingService) getBoletoURL(payment *domain.Payment) *string {
	if payment.PaymentMethod == domain.PaymentMethodBoleto && payment.AsaasPaymentID != nil {
		url := fmt.Sprintf("https://www.asaas.com/b/%s", *payment.AsaasPaymentID)
		return &url
	}
	return nil
}

// CompleteOnboarding finaliza onboarding após pagamento
func (s *OnboardingService) CompleteOnboarding(ctx context.Context, subscriptionID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("falha ao buscar assinatura: %w", err)
	}

	// Ativar assinatura
	subscription.Activate()

	// Salvar alterações
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("falha ao atualizar assinatura: %w", err)
	}

	// Buscar plano para evento
	plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return fmt.Errorf("falha ao buscar plano: %w", err)
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

// GetOnboardingStatus retorna status do onboarding
func (s *OnboardingService) GetOnboardingStatus(ctx context.Context, tenantID uuid.UUID) (*OnboardingStatus, error) {
	subscription, err := s.subscriptionRepo.GetCurrentByTenantID(ctx, tenantID)
	if err != nil {
		return &OnboardingStatus{
			TenantID:  tenantID,
			Completed: false,
			Stage:     "not_started",
		}, nil
	}

	status := &OnboardingStatus{
		TenantID:       tenantID,
		SubscriptionID: &subscription.ID,
		Completed:      subscription.IsActive(),
		Stage:          s.getOnboardingStage(subscription),
		TrialEndDate:   subscription.TrialEndDate,
	}

	return status, nil
}

// getOnboardingStage determina o estágio do onboarding
func (s *OnboardingService) getOnboardingStage(subscription *domain.Subscription) string {
	switch subscription.Status {
	case domain.SubscriptionStatusTrial:
		return "trial"
	case domain.SubscriptionStatusActive:
		return "active"
	case domain.SubscriptionStatusPaymentPending:
		return "payment_pending"
	case domain.SubscriptionStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// OnboardingStatus status do onboarding
type OnboardingStatus struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	SubscriptionID *uuid.UUID `json:"subscription_id"`
	Completed      bool       `json:"completed"`
	Stage          string     `json:"stage"`
	TrialEndDate   *time.Time `json:"trial_end_date"`
}

// ValidateDocument valida documento CPF/CNPJ
func (s *OnboardingService) ValidateDocument(ctx context.Context, document string) (bool, error) {
	// Limpar documento
	clean := ""
	for _, char := range document {
		if char >= '0' && char <= '9' {
			clean += string(char)
		}
	}

	// Validar comprimento
	if len(clean) != 11 && len(clean) != 14 {
		return false, errors.New("documento deve ter 11 (CPF) ou 14 (CNPJ) dígitos")
	}

	// Validar se todos os dígitos são iguais
	allSame := true
	for i := 1; i < len(clean); i++ {
		if clean[i] != clean[0] {
			allSame = false
			break
		}
	}

	if allSame {
		return false, errors.New("documento inválido")
	}

	// TODO: Implementar validação completa de CPF/CNPJ
	return true, nil
}

// GetAvailablePlans retorna planos disponíveis
func (s *OnboardingService) GetAvailablePlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.GetActive(ctx)
}

// ApplyCoupon aplica cupom de desconto
func (s *OnboardingService) ApplyCoupon(ctx context.Context, couponCode string, planID uuid.UUID) (*CouponResult, error) {
	// TODO: Implementar sistema de cupons
	return &CouponResult{
		Valid:            false,
		DiscountAmount:   0,
		DiscountPercent:  0,
		Message:          "Sistema de cupons não implementado",
	}, nil
}

// CouponResult resultado da aplicação do cupom
type CouponResult struct {
	Valid            bool    `json:"valid"`
	DiscountAmount   int64   `json:"discount_amount"`
	DiscountPercent  float64 `json:"discount_percent"`
	Message          string  `json:"message"`
}