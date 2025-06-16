package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// SubscriptionService serviço de aplicação para assinaturas
type SubscriptionService struct {
	subscriptionRepo domain.SubscriptionRepository
	planRepo         domain.PlanRepository
	tenantRepo       domain.TenantRepository
	quotaRepo        domain.QuotaRepository
	eventPublisher   domain.EventPublisher
	logger           *zap.Logger
}

// NewSubscriptionService cria nova instância do serviço
func NewSubscriptionService(
	subscriptionRepo domain.SubscriptionRepository,
	planRepo domain.PlanRepository,
	tenantRepo domain.TenantRepository,
	quotaRepo domain.QuotaRepository,
	eventPublisher domain.EventPublisher,
	logger *zap.Logger,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		tenantRepo:       tenantRepo,
		quotaRepo:        quotaRepo,
		eventPublisher:   eventPublisher,
		logger:           logger,
	}
}

// CreateSubscriptionRequest request para criação de assinatura
type CreateSubscriptionRequest struct {
	TenantID        string                    `json:"tenant_id" validate:"required"`
	PlanID          string                    `json:"plan_id" validate:"required"`
	BillingInterval domain.BillingInterval   `json:"billing_interval" validate:"required"`
	TrialDays       int                       `json:"trial_days"`
}

// SubscriptionResponse response com dados da assinatura
type SubscriptionResponse struct {
	ID                 string                    `json:"id"`
	TenantID           string                    `json:"tenant_id"`
	PlanID             string                    `json:"plan_id"`
	Plan               *PlanResponse             `json:"plan,omitempty"`
	Status             domain.SubscriptionStatus `json:"status"`
	CurrentPeriodStart time.Time                 `json:"current_period_start"`
	CurrentPeriodEnd   time.Time                 `json:"current_period_end"`
	CancelAtPeriodEnd  bool                      `json:"cancel_at_period_end"`
	TrialStart         *time.Time                `json:"trial_start"`
	TrialEnd           *time.Time                `json:"trial_end"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
	CanceledAt         *time.Time                `json:"canceled_at"`
}

// PlanResponse response com dados do plano
type PlanResponse struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Type            domain.PlanType          `json:"type"`
	Description     string                   `json:"description"`
	Price           int64                    `json:"price"`
	Currency        string                   `json:"currency"`
	BillingInterval domain.BillingInterval   `json:"billing_interval"`
	Features        domain.PlanFeatures      `json:"features"`
	Quotas          domain.PlanQuotas        `json:"quotas"`
	IsActive        bool                     `json:"is_active"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// ChangePlanRequest request para mudança de plano
type ChangePlanRequest struct {
	NewPlanID string `json:"new_plan_id" validate:"required"`
	Reason    string `json:"reason"`
}

// CreateSubscription cria nova assinatura
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	s.logger.Info("Creating subscription", 
		zap.String("tenant_id", req.TenantID),
		zap.String("plan_id", req.PlanID),
	)

	// Verifica se tenant existe
	tenant, err := s.tenantRepo.GetByID(req.TenantID)
	if err != nil {
		return nil, err
	}

	// Verifica se plano existe
	plan, err := s.planRepo.GetByID(req.PlanID)
	if err != nil {
		return nil, err
	}

	if !plan.IsActive {
		return nil, fmt.Errorf("plano não está ativo")
	}

	// Verifica se já existe assinatura ativa
	existing, _ := s.subscriptionRepo.GetByTenantID(req.TenantID)
	if existing != nil && existing.IsActive() {
		return nil, domain.ErrSubscriptionExists
	}

	// Cria nova assinatura
	subscription := &domain.Subscription{
		ID:                 uuid.New().String(),
		TenantID:           req.TenantID,
		PlanID:             req.PlanID,
		Status:             domain.SubscriptionStatusTrialing,
		CurrentPeriodStart: time.Now(),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Configura período de trial se especificado
	if req.TrialDays > 0 {
		now := time.Now()
		trialEnd := now.AddDate(0, 0, req.TrialDays)
		subscription.TrialStart = &now
		subscription.TrialEnd = &trialEnd
		subscription.CurrentPeriodEnd = trialEnd
	} else {
		// Sem trial, inicia período pago imediatamente
		subscription.Status = domain.SubscriptionStatusActive
		switch req.BillingInterval {
		case domain.BillingMonthly:
			subscription.CurrentPeriodEnd = time.Now().AddDate(0, 1, 0)
		case domain.BillingYearly:
			subscription.CurrentPeriodEnd = time.Now().AddDate(1, 0, 0)
		}
	}

	// Salva assinatura
	if err := s.subscriptionRepo.Create(subscription); err != nil {
		s.logger.Error("Failed to create subscription", zap.Error(err))
		return nil, fmt.Errorf("erro ao criar assinatura: %w", err)
	}

	// Atualiza plano do tenant
	tenant.PlanType = plan.Type
	if err := s.tenantRepo.Update(tenant); err != nil {
		s.logger.Error("Failed to update tenant plan", zap.Error(err))
	}

	// Atualiza limites de quota
	quotas := domain.GetDefaultPlanQuotas(plan.Type)
	quotaLimits := domain.GetQuotaLimitFromPlan(req.TenantID, quotas)
	if err := s.quotaRepo.UpdateLimits(quotaLimits); err != nil {
		s.logger.Error("Failed to update quota limits", zap.Error(err))
	}

	// Publica evento
	event := domain.NewSubscriptionCreatedEvent(subscription, req.TrialDays)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish subscription created event", zap.Error(err))
	}

	s.logger.Info("Subscription created successfully", zap.String("subscription_id", subscription.ID))

	return s.toSubscriptionResponse(subscription, plan), nil
}

// GetSubscription obtém assinatura por ID
func (s *SubscriptionService) GetSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(subscription, plan), nil
}

// GetSubscriptionByTenant obtém assinatura ativa do tenant
func (s *SubscriptionService) GetSubscriptionByTenant(ctx context.Context, tenantID string) (*SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetByTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(subscription, plan), nil
}

// ActivateSubscription ativa uma assinatura
func (s *SubscriptionService) ActivateSubscription(ctx context.Context, subscriptionID string) error {
	s.logger.Info("Activating subscription", zap.String("subscription_id", subscriptionID))

	subscription, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return err
	}

	if subscription.Status == domain.SubscriptionStatusActive {
		return fmt.Errorf("assinatura já está ativa")
	}

	subscription.Status = domain.SubscriptionStatusActive
	subscription.UpdatedAt = time.Now()

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		s.logger.Error("Failed to activate subscription", zap.Error(err))
		return fmt.Errorf("erro ao ativar assinatura: %w", err)
	}

	// Publica evento
	event := domain.NewSubscriptionActivatedEvent(subscription)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish subscription activated event", zap.Error(err))
	}

	s.logger.Info("Subscription activated successfully", zap.String("subscription_id", subscriptionID))

	return nil
}

// CancelSubscription cancela uma assinatura
func (s *SubscriptionService) CancelSubscription(ctx context.Context, subscriptionID, reason string, immediately bool) error {
	s.logger.Info("Canceling subscription", 
		zap.String("subscription_id", subscriptionID),
		zap.String("reason", reason),
		zap.Bool("immediately", immediately),
	)

	subscription, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return err
	}

	if subscription.Status == domain.SubscriptionStatusCanceled {
		return fmt.Errorf("assinatura já está cancelada")
	}

	if immediately || subscription.IsInTrial() {
		subscription.Cancel()
	} else {
		subscription.ScheduleCancellation()
	}

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		s.logger.Error("Failed to cancel subscription", zap.Error(err))
		return fmt.Errorf("erro ao cancelar assinatura: %w", err)
	}

	// Publica evento
	event := domain.NewSubscriptionCanceledEvent(subscription, reason)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish subscription canceled event", zap.Error(err))
	}

	s.logger.Info("Subscription canceled successfully", zap.String("subscription_id", subscriptionID))

	return nil
}

// ReactivateSubscription reativa uma assinatura cancelada
func (s *SubscriptionService) ReactivateSubscription(ctx context.Context, subscriptionID string) error {
	s.logger.Info("Reactivating subscription", zap.String("subscription_id", subscriptionID))

	subscription, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return err
	}

	subscription.Reactivate()

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		s.logger.Error("Failed to reactivate subscription", zap.Error(err))
		return fmt.Errorf("erro ao reativar assinatura: %w", err)
	}

	// Publica evento
	event := domain.NewSubscriptionActivatedEvent(subscription)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish subscription reactivated event", zap.Error(err))
	}

	s.logger.Info("Subscription reactivated successfully", zap.String("subscription_id", subscriptionID))

	return nil
}

// ChangePlan altera o plano de uma assinatura
func (s *SubscriptionService) ChangePlan(ctx context.Context, subscriptionID string, req *ChangePlanRequest) (*SubscriptionResponse, error) {
	s.logger.Info("Changing subscription plan", 
		zap.String("subscription_id", subscriptionID),
		zap.String("new_plan_id", req.NewPlanID),
	)

	subscription, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return nil, err
	}

	if subscription.PlanID == req.NewPlanID {
		return nil, fmt.Errorf("plano já é o mesmo")
	}

	// Busca novo plano
	newPlan, err := s.planRepo.GetByID(req.NewPlanID)
	if err != nil {
		return nil, err
	}

	if !newPlan.IsActive {
		return nil, fmt.Errorf("novo plano não está ativo")
	}

	// Busca plano atual
	oldPlan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil {
		return nil, err
	}

	// Atualiza assinatura
	oldPlanID := subscription.PlanID
	subscription.PlanID = req.NewPlanID
	subscription.UpdatedAt = time.Now()

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		s.logger.Error("Failed to change subscription plan", zap.Error(err))
		return nil, fmt.Errorf("erro ao alterar plano: %w", err)
	}

	// Atualiza plano do tenant
	tenant, err := s.tenantRepo.GetByID(subscription.TenantID)
	if err == nil {
		tenant.PlanType = newPlan.Type
		tenant.UpdatedAt = time.Now()
		if err := s.tenantRepo.Update(tenant); err != nil {
			s.logger.Error("Failed to update tenant plan type", zap.Error(err))
		}
	}

	// Atualiza limites de quota
	quotas := domain.GetDefaultPlanQuotas(newPlan.Type)
	quotaLimits := domain.GetQuotaLimitFromPlan(subscription.TenantID, quotas)
	if err := s.quotaRepo.UpdateLimits(quotaLimits); err != nil {
		s.logger.Error("Failed to update quota limits", zap.Error(err))
	}

	// Publica evento
	event := domain.NewPlanChangedEvent(subscriptionID, subscription.TenantID, oldPlanID, req.NewPlanID, req.Reason)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish plan changed event", zap.Error(err))
	}

	s.logger.Info("Subscription plan changed successfully", 
		zap.String("subscription_id", subscriptionID),
		zap.String("from_plan", oldPlan.Name),
		zap.String("to_plan", newPlan.Name),
	)

	return s.toSubscriptionResponse(subscription, newPlan), nil
}

// RenewSubscription renova uma assinatura
func (s *SubscriptionService) RenewSubscription(ctx context.Context, subscriptionID string, billingInterval domain.BillingInterval) error {
	s.logger.Info("Renewing subscription", 
		zap.String("subscription_id", subscriptionID),
		zap.String("billing_interval", string(billingInterval)),
	)

	subscription, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return err
	}

	previousPeriodEnd := subscription.CurrentPeriodEnd
	subscription.Renew(billingInterval)

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		s.logger.Error("Failed to renew subscription", zap.Error(err))
		return fmt.Errorf("erro ao renovar assinatura: %w", err)
	}

	// Publica evento
	event := domain.NewSubscriptionRenewedEvent(subscription, previousPeriodEnd, billingInterval)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish subscription renewed event", zap.Error(err))
	}

	s.logger.Info("Subscription renewed successfully", zap.String("subscription_id", subscriptionID))

	return nil
}

// GetExpiringSubscriptions obtém assinaturas que vão expirar
func (s *SubscriptionService) GetExpiringSubscriptions(ctx context.Context, days int) ([]*SubscriptionResponse, error) {
	subscriptions, err := s.subscriptionRepo.GetExpiring(days)
	if err != nil {
		return nil, err
	}

	var responses []*SubscriptionResponse
	for _, subscription := range subscriptions {
		plan, err := s.planRepo.GetByID(subscription.PlanID)
		if err != nil {
			s.logger.Error("Failed to get plan for expiring subscription", 
				zap.String("subscription_id", subscription.ID),
				zap.Error(err),
			)
			continue
		}
		responses = append(responses, s.toSubscriptionResponse(subscription, plan))
	}

	return responses, nil
}

// ListPlans lista todos os planos disponíveis
func (s *SubscriptionService) ListPlans(ctx context.Context, activeOnly bool) ([]*PlanResponse, error) {
	plans, err := s.planRepo.GetAll(activeOnly)
	if err != nil {
		return nil, err
	}

	var responses []*PlanResponse
	for _, plan := range plans {
		responses = append(responses, s.toPlanResponse(plan))
	}

	return responses, nil
}

// GetPlan obtém plano por ID
func (s *SubscriptionService) GetPlan(ctx context.Context, planID string) (*PlanResponse, error) {
	plan, err := s.planRepo.GetByID(planID)
	if err != nil {
		return nil, err
	}

	return s.toPlanResponse(plan), nil
}

// toSubscriptionResponse converte domain.Subscription para SubscriptionResponse
func (s *SubscriptionService) toSubscriptionResponse(subscription *domain.Subscription, plan *domain.Plan) *SubscriptionResponse {
	response := &SubscriptionResponse{
		ID:                 subscription.ID,
		TenantID:           subscription.TenantID,
		PlanID:             subscription.PlanID,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		TrialStart:         subscription.TrialStart,
		TrialEnd:           subscription.TrialEnd,
		CreatedAt:          subscription.CreatedAt,
		UpdatedAt:          subscription.UpdatedAt,
		CanceledAt:         subscription.CanceledAt,
	}

	if plan != nil {
		response.Plan = s.toPlanResponse(plan)
	}

	return response
}

// toPlanResponse converte domain.Plan para PlanResponse
func (s *SubscriptionService) toPlanResponse(plan *domain.Plan) *PlanResponse {
	return &PlanResponse{
		ID:              plan.ID,
		Name:            plan.Name,
		Type:            plan.Type,
		Description:     plan.Description,
		Price:           plan.Price,
		Currency:        plan.Currency,
		BillingInterval: plan.BillingInterval,
		Features:        plan.Features,
		Quotas:          plan.Quotas,
		IsActive:        plan.IsActive,
		CreatedAt:       plan.CreatedAt,
		UpdatedAt:       plan.UpdatedAt,
	}
}