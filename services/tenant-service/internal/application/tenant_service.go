package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// TenantService serviço de aplicação para tenants
type TenantService struct {
	tenantRepo       domain.TenantRepository
	subscriptionRepo domain.SubscriptionRepository
	planRepo         domain.PlanRepository
	quotaRepo        domain.QuotaRepository
	eventPublisher   domain.EventPublisher
	logger           *zap.Logger
}

// NewTenantService cria nova instância do serviço
func NewTenantService(
	tenantRepo domain.TenantRepository,
	subscriptionRepo domain.SubscriptionRepository,
	planRepo domain.PlanRepository,
	quotaRepo domain.QuotaRepository,
	eventPublisher domain.EventPublisher,
	logger *zap.Logger,
) *TenantService {
	return &TenantService{
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		quotaRepo:        quotaRepo,
		eventPublisher:   eventPublisher,
		logger:           logger,
	}
}

// CreateTenantRequest request para criação de tenant
type CreateTenantRequest struct {
	Name        string                    `json:"name" validate:"required,min=3,max=50"`
	LegalName   string                    `json:"legal_name"`
	Document    string                    `json:"document"`
	Email       string                    `json:"email" validate:"required,email"`
	Phone       string                    `json:"phone"`
	Website     string                    `json:"website"`
	Address     *CreateAddressRequest     `json:"address"`
	PlanType    domain.PlanType           `json:"plan_type" validate:"required"`
	OwnerUserID string                    `json:"owner_user_id" validate:"required"`
	Settings    *CreateTenantSettings     `json:"settings"`
}

// CreateAddressRequest request para criação de endereço
type CreateAddressRequest struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
	Country    string `json:"country"`
}

// CreateTenantSettings request para configurações do tenant
type CreateTenantSettings struct {
	Timezone      string                       `json:"timezone"`
	Language      string                       `json:"language"`
	Currency      string                       `json:"currency"`
	NotificationPrefs *CreateNotificationPrefs `json:"notification_prefs"`
	BrandingConfig    *CreateBrandingConfig    `json:"branding_config"`
}

// CreateNotificationPrefs request para preferências de notificação
type CreateNotificationPrefs struct {
	WhatsAppEnabled   bool `json:"whatsapp_enabled"`
	EmailEnabled      bool `json:"email_enabled"`
	TelegramEnabled   bool `json:"telegram_enabled"`
	PushEnabled       bool `json:"push_enabled"`
	BusinessHoursOnly bool `json:"business_hours_only"`
}

// CreateBrandingConfig request para configuração de marca
type CreateBrandingConfig struct {
	LogoURL        string `json:"logo_url"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	FontFamily     string `json:"font_family"`
	CustomDomain   string `json:"custom_domain"`
}

// UpdateTenantRequest request para atualização de tenant
type UpdateTenantRequest struct {
	Name        *string                   `json:"name,omitempty"`
	LegalName   *string                   `json:"legal_name,omitempty"`
	Document    *string                   `json:"document,omitempty"`
	Email       *string                   `json:"email,omitempty"`
	Phone       *string                   `json:"phone,omitempty"`
	Website     *string                   `json:"website,omitempty"`
	Address     *CreateAddressRequest     `json:"address,omitempty"`
	Settings    *CreateTenantSettings     `json:"settings,omitempty"`
}

// TenantResponse response com dados do tenant
type TenantResponse struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	LegalName     string                     `json:"legal_name"`
	Document      string                     `json:"document"`
	Email         string                     `json:"email"`
	Phone         string                     `json:"phone"`
	Website       string                     `json:"website"`
	Address       domain.Address             `json:"address"`
	Status        domain.TenantStatus        `json:"status"`
	PlanType      domain.PlanType            `json:"plan_type"`
	OwnerUserID   string                     `json:"owner_user_id"`
	Settings      domain.TenantSettings      `json:"settings"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	ActivatedAt   *time.Time                 `json:"activated_at"`
	SuspendedAt   *time.Time                 `json:"suspended_at"`
}

// CreateTenant cria um novo tenant
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error) {
	s.logger.Info("Creating tenant", zap.String("name", req.Name), zap.String("email", req.Email))

	// Verifica se já existe tenant com mesmo documento ou email
	if req.Document != "" {
		existing, _ := s.tenantRepo.GetByDocument(req.Document)
		if existing != nil {
			return nil, domain.ErrTenantExists
		}
	}

	// Cria ID único
	tenantID := uuid.New().String()

	// Cria entidade tenant
	tenant := &domain.Tenant{
		ID:          tenantID,
		Name:        req.Name,
		LegalName:   req.LegalName,
		Document:    req.Document,
		Email:       req.Email,
		Phone:       req.Phone,
		Website:     req.Website,
		Status:      domain.TenantStatusPending,
		PlanType:    req.PlanType,
		OwnerUserID: req.OwnerUserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Define endereço se fornecido
	if req.Address != nil {
		tenant.Address = domain.Address{
			Street:     req.Address.Street,
			Number:     req.Address.Number,
			Complement: req.Address.Complement,
			District:   req.Address.District,
			City:       req.Address.City,
			State:      req.Address.State,
			ZipCode:    req.Address.ZipCode,
			Country:    req.Address.Country,
		}
	}

	// Define configurações
	if req.Settings != nil {
		tenant.Settings = s.buildTenantSettings(req.Settings)
	} else {
		tenant.SetDefaultSettings()
	}

	// Valida tenant
	if err := s.validateTenant(tenant); err != nil {
		s.logger.Error("Tenant validation failed", zap.Error(err))
		return nil, err
	}

	// Salva tenant
	if err := s.tenantRepo.Create(tenant); err != nil {
		s.logger.Error("Failed to create tenant", zap.Error(err))
		return nil, fmt.Errorf("erro ao criar tenant: %w", err)
	}

	// Cria subscription inicial
	_, subErr := s.createInitialSubscription(tenantID, req.PlanType)
	if subErr != nil {
		s.logger.Error("Failed to create initial subscription", zap.Error(subErr))
		return nil, fmt.Errorf("erro ao criar assinatura inicial: %w", subErr)
	}

	// Cria quota usage inicial
	quotaUsage := domain.NewQuotaUsage(tenantID)
	if err := s.quotaRepo.UpdateUsage(quotaUsage); err != nil {
		s.logger.Error("Failed to create initial quota usage", zap.Error(err))
	}

	// Publica evento
	event := domain.NewTenantCreatedEvent(tenant)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish tenant created event", zap.Error(err))
	}

	s.logger.Info("Tenant created successfully", zap.String("tenant_id", tenantID))

	return s.toTenantResponse(tenant), nil
}

// GetTenant obtém tenant por ID
func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*TenantResponse, error) {
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, err
	}

	return s.toTenantResponse(tenant), nil
}

// GetTenantByDocument obtém tenant por documento
func (s *TenantService) GetTenantByDocument(ctx context.Context, document string) (*TenantResponse, error) {
	tenant, err := s.tenantRepo.GetByDocument(document)
	if err != nil {
		return nil, err
	}

	return s.toTenantResponse(tenant), nil
}

// GetTenantsByOwner obtém tenants de um proprietário
func (s *TenantService) GetTenantsByOwner(ctx context.Context, ownerUserID string) ([]*TenantResponse, error) {
	tenants, err := s.tenantRepo.GetByOwner(ownerUserID)
	if err != nil {
		return nil, err
	}

	var responses []*TenantResponse
	for _, tenant := range tenants {
		responses = append(responses, s.toTenantResponse(tenant))
	}

	return responses, nil
}

// UpdateTenant atualiza dados do tenant
func (s *TenantService) UpdateTenant(ctx context.Context, tenantID string, req *UpdateTenantRequest, updatedBy string) (*TenantResponse, error) {
	s.logger.Info("Updating tenant", zap.String("tenant_id", tenantID))

	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, err
	}

	// Coleta mudanças para o evento
	changes := make(map[string]interface{})

	// Atualiza campos se fornecidos
	if req.Name != nil && *req.Name != tenant.Name {
		changes["name"] = map[string]string{"from": tenant.Name, "to": *req.Name}
		tenant.Name = *req.Name
	}

	if req.LegalName != nil && *req.LegalName != tenant.LegalName {
		changes["legal_name"] = map[string]string{"from": tenant.LegalName, "to": *req.LegalName}
		tenant.LegalName = *req.LegalName
	}

	if req.Document != nil && *req.Document != tenant.Document {
		changes["document"] = map[string]string{"from": tenant.Document, "to": *req.Document}
		tenant.Document = *req.Document
	}

	if req.Email != nil && *req.Email != tenant.Email {
		changes["email"] = map[string]string{"from": tenant.Email, "to": *req.Email}
		tenant.Email = *req.Email
	}

	if req.Phone != nil && *req.Phone != tenant.Phone {
		changes["phone"] = map[string]string{"from": tenant.Phone, "to": *req.Phone}
		tenant.Phone = *req.Phone
	}

	if req.Website != nil && *req.Website != tenant.Website {
		changes["website"] = map[string]string{"from": tenant.Website, "to": *req.Website}
		tenant.Website = *req.Website
	}

	if req.Address != nil {
		tenant.Address = domain.Address{
			Street:     req.Address.Street,
			Number:     req.Address.Number,
			Complement: req.Address.Complement,
			District:   req.Address.District,
			City:       req.Address.City,
			State:      req.Address.State,
			ZipCode:    req.Address.ZipCode,
			Country:    req.Address.Country,
		}
		changes["address"] = "updated"
	}

	if req.Settings != nil {
		tenant.Settings = s.buildTenantSettings(req.Settings)
		changes["settings"] = "updated"
	}

	tenant.UpdatedAt = time.Now()

	// Valida tenant
	if err := s.validateTenant(tenant); err != nil {
		s.logger.Error("Tenant validation failed on update", zap.Error(err))
		return nil, err
	}

	// Salva alterações
	if err := s.tenantRepo.Update(tenant); err != nil {
		s.logger.Error("Failed to update tenant", zap.Error(err))
		return nil, fmt.Errorf("erro ao atualizar tenant: %w", err)
	}

	// Publica evento se houve mudanças
	if len(changes) > 0 {
		event := domain.NewTenantUpdatedEvent(tenantID, changes, updatedBy)
		if err := s.eventPublisher.Publish(event); err != nil {
			s.logger.Error("Failed to publish tenant updated event", zap.Error(err))
		}
	}

	s.logger.Info("Tenant updated successfully", zap.String("tenant_id", tenantID))

	return s.toTenantResponse(tenant), nil
}

// ActivateTenant ativa um tenant
func (s *TenantService) ActivateTenant(ctx context.Context, tenantID, activatedBy string) error {
	s.logger.Info("Activating tenant", zap.String("tenant_id", tenantID))

	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return err
	}

	if tenant.Status == domain.TenantStatusActive {
		return fmt.Errorf("tenant já está ativo")
	}

	tenant.Activate()

	if err := s.tenantRepo.Update(tenant); err != nil {
		s.logger.Error("Failed to activate tenant", zap.Error(err))
		return fmt.Errorf("erro ao ativar tenant: %w", err)
	}

	// Publica evento
	event := domain.NewTenantActivatedEvent(tenantID, activatedBy)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish tenant activated event", zap.Error(err))
	}

	s.logger.Info("Tenant activated successfully", zap.String("tenant_id", tenantID))

	return nil
}

// SuspendTenant suspende um tenant
func (s *TenantService) SuspendTenant(ctx context.Context, tenantID, reason, suspendedBy string) error {
	s.logger.Info("Suspending tenant", zap.String("tenant_id", tenantID), zap.String("reason", reason))

	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return err
	}

	tenant.Suspend()

	if err := s.tenantRepo.Update(tenant); err != nil {
		s.logger.Error("Failed to suspend tenant", zap.Error(err))
		return fmt.Errorf("erro ao suspender tenant: %w", err)
	}

	// Publica evento
	event := domain.NewTenantSuspendedEvent(tenantID, reason, suspendedBy)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish tenant suspended event", zap.Error(err))
	}

	s.logger.Info("Tenant suspended successfully", zap.String("tenant_id", tenantID))

	return nil
}

// CancelTenant cancela um tenant
func (s *TenantService) CancelTenant(ctx context.Context, tenantID, reason, canceledBy string) error {
	s.logger.Info("Canceling tenant", zap.String("tenant_id", tenantID), zap.String("reason", reason))

	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return err
	}

	tenant.Cancel()

	if err := s.tenantRepo.Update(tenant); err != nil {
		s.logger.Error("Failed to cancel tenant", zap.Error(err))
		return fmt.Errorf("erro ao cancelar tenant: %w", err)
	}

	// Publica evento
	event := domain.NewTenantCanceledEvent(tenantID, reason, canceledBy)
	if err := s.eventPublisher.Publish(event); err != nil {
		s.logger.Error("Failed to publish tenant canceled event", zap.Error(err))
	}

	s.logger.Info("Tenant canceled successfully", zap.String("tenant_id", tenantID))

	return nil
}

// validateTenant valida entidade tenant
func (s *TenantService) validateTenant(tenant *domain.Tenant) error {
	if err := tenant.ValidateName(); err != nil {
		return err
	}
	if err := tenant.ValidateEmail(); err != nil {
		return err
	}
	if err := tenant.ValidateDocument(); err != nil {
		return err
	}
	if err := tenant.ValidatePlan(); err != nil {
		return err
	}
	if err := tenant.ValidateStatus(); err != nil {
		return err
	}
	return nil
}

// buildTenantSettings constrói TenantSettings a partir do request
func (s *TenantService) buildTenantSettings(req *CreateTenantSettings) domain.TenantSettings {
	settings := domain.TenantSettings{
		Timezone: "America/Sao_Paulo",
		Language: "pt-BR",
		Currency: "BRL",
		IntegrationKeys: make(map[string]string),
	}

	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}
	if req.Language != "" {
		settings.Language = req.Language
	}
	if req.Currency != "" {
		settings.Currency = req.Currency
	}

	// Configurações de notificação
	if req.NotificationPrefs != nil {
		settings.NotificationPrefs = domain.NotificationPrefs{
			WhatsAppEnabled:   req.NotificationPrefs.WhatsAppEnabled,
			EmailEnabled:      req.NotificationPrefs.EmailEnabled,
			TelegramEnabled:   req.NotificationPrefs.TelegramEnabled,
			PushEnabled:       req.NotificationPrefs.PushEnabled,
			BusinessHoursOnly: req.NotificationPrefs.BusinessHoursOnly,
		}
	} else {
		settings.NotificationPrefs = domain.NotificationPrefs{
			WhatsAppEnabled: true,
			EmailEnabled:    true,
			PushEnabled:     true,
		}
	}

	// Configurações de marca
	if req.BrandingConfig != nil {
		settings.BrandingConfig = domain.BrandingConfig{
			LogoURL:       req.BrandingConfig.LogoURL,
			PrimaryColor:  req.BrandingConfig.PrimaryColor,
			SecondaryColor: req.BrandingConfig.SecondaryColor,
			FontFamily:    req.BrandingConfig.FontFamily,
			CustomDomain:  req.BrandingConfig.CustomDomain,
		}
	} else {
		settings.BrandingConfig = domain.BrandingConfig{
			PrimaryColor:   "#1976d2",
			SecondaryColor: "#424242",
			FontFamily:     "Roboto",
		}
	}

	return settings
}

// createInitialSubscription cria assinatura inicial para o tenant
func (s *TenantService) createInitialSubscription(tenantID string, planType domain.PlanType) (*domain.Subscription, error) {
	// Busca o plano
	plan, err := s.planRepo.GetByType(planType)
	if err != nil {
		return nil, fmt.Errorf("plano não encontrado: %w", err)
	}

	// Cria assinatura
	subscription := &domain.Subscription{
		ID:                 uuid.New().String(),
		TenantID:           tenantID,
		PlanID:             plan.ID,
		Status:             domain.SubscriptionStatusTrialing,
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 0, 7), // 7 dias de trial
		TrialStart:         &time.Time{},
		TrialEnd:           &time.Time{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	now := time.Now()
	trialEnd := now.AddDate(0, 0, 7)
	subscription.TrialStart = &now
	subscription.TrialEnd = &trialEnd

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, err
	}

	// Cria limites de quota
	quotas := domain.GetDefaultPlanQuotas(planType)
	quotaLimits := domain.GetQuotaLimitFromPlan(tenantID, quotas)
	if err := s.quotaRepo.UpdateLimits(quotaLimits); err != nil {
		s.logger.Error("Failed to create quota limits", zap.Error(err))
	}

	return subscription, nil
}

// toTenantResponse converte domain.Tenant para TenantResponse
func (s *TenantService) toTenantResponse(tenant *domain.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:            tenant.ID,
		Name:          tenant.Name,
		LegalName:     tenant.LegalName,
		Document:      tenant.Document,
		Email:         tenant.Email,
		Phone:         tenant.Phone,
		Website:       tenant.Website,
		Address:       tenant.Address,
		Status:        tenant.Status,
		PlanType:      tenant.PlanType,
		OwnerUserID:   tenant.OwnerUserID,
		Settings:      tenant.Settings,
		CreatedAt:     tenant.CreatedAt,
		UpdatedAt:     tenant.UpdatedAt,
		ActivatedAt:   tenant.ActivatedAt,
		SuspendedAt:   tenant.SuspendedAt,
	}
}