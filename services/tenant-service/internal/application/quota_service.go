package application

import (
	"context"
	"fmt"
	"time"

	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// QuotaService serviço de aplicação para quotas
type QuotaService struct {
	quotaRepo        domain.QuotaRepository
	tenantRepo       domain.TenantRepository
	subscriptionRepo domain.SubscriptionRepository
	planRepo         domain.PlanRepository
	eventPublisher   domain.EventPublisher
	logger           *zap.Logger
}

// NewQuotaService cria nova instância do serviço
func NewQuotaService(
	quotaRepo domain.QuotaRepository,
	tenantRepo domain.TenantRepository,
	subscriptionRepo domain.SubscriptionRepository,
	planRepo domain.PlanRepository,
	eventPublisher domain.EventPublisher,
	logger *zap.Logger,
) *QuotaService {
	return &QuotaService{
		quotaRepo:        quotaRepo,
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		eventPublisher:   eventPublisher,
		logger:           logger,
	}
}

// QuotaUsageResponse response com dados de uso de quota
type QuotaUsageResponse struct {
	TenantID              string    `json:"tenant_id"`
	ProcessesCount        int       `json:"processes_count"`
	UsersCount            int       `json:"users_count"`
	ClientsCount          int       `json:"clients_count"`
	DataJudQueriesDaily   int       `json:"datajud_queries_daily"`
	DataJudQueriesMonth   int       `json:"datajud_queries_month"`
	AIQueriesMonthly      int       `json:"ai_queries_monthly"`
	StorageUsedGB         float64   `json:"storage_used_gb"`
	WebhooksCount         int       `json:"webhooks_count"`
	APICallsDaily         int       `json:"api_calls_daily"`
	APICallsMonthly       int       `json:"api_calls_monthly"`
	LastUpdated           time.Time `json:"last_updated"`
	LastResetDaily        time.Time `json:"last_reset_daily"`
	LastResetMonthly      time.Time `json:"last_reset_monthly"`
}

// QuotaLimitResponse response com limites de quota
type QuotaLimitResponse struct {
	TenantID              string    `json:"tenant_id"`
	MaxProcesses          int       `json:"max_processes"`
	MaxUsers              int       `json:"max_users"`
	MaxClients            int       `json:"max_clients"`
	DataJudQueriesDaily   int       `json:"datajud_queries_daily"`
	AIQueriesMonthly      int       `json:"ai_queries_monthly"`
	StorageGB             int       `json:"storage_gb"`
	MaxWebhooks           int       `json:"max_webhooks"`
	MaxAPICallsDaily      int       `json:"max_api_calls_daily"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// QuotaCheckResponse response com verificação de quota
type QuotaCheckResponse struct {
	TenantID    string                 `json:"tenant_id"`
	Checks      []*domain.QuotaCheck   `json:"checks"`
	CanProceed  bool                   `json:"can_proceed"`
	Warnings    []*domain.QuotaCheck   `json:"warnings"`
	Exceeded    []*domain.QuotaCheck   `json:"exceeded"`
}

// IncrementQuotaRequest request para incrementar quota
type IncrementQuotaRequest struct {
	QuotaType string  `json:"quota_type" validate:"required"`
	Amount    int     `json:"amount" validate:"min=1"`
	AmountGB  float64 `json:"amount_gb,omitempty"`
}

// UpdateStorageRequest request para atualizar armazenamento
type UpdateStorageRequest struct {
	UsageGB float64 `json:"usage_gb" validate:"min=0"`
}

// GetQuotaUsage obtém uso atual de quotas do tenant
func (s *QuotaService) GetQuotaUsage(ctx context.Context, tenantID string) (*QuotaUsageResponse, error) {
	usage, err := s.quotaRepo.GetUsage(tenantID)
	if err != nil {
		if err == domain.ErrQuotaNotFound {
			// Cria usage inicial se não existir
			usage = domain.NewQuotaUsage(tenantID)
			if err := s.quotaRepo.UpdateUsage(usage); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return s.toQuotaUsageResponse(usage), nil
}

// GetQuotaLimits obtém limites de quota do tenant
func (s *QuotaService) GetQuotaLimits(ctx context.Context, tenantID string) (*QuotaLimitResponse, error) {
	limits, err := s.quotaRepo.GetLimits(tenantID)
	if err != nil {
		return nil, err
	}

	return s.toQuotaLimitResponse(limits), nil
}

// CheckQuotas verifica todas as quotas do tenant
func (s *QuotaService) CheckQuotas(ctx context.Context, tenantID string) (*QuotaCheckResponse, error) {
	// Obtém uso atual
	usage, err := s.quotaRepo.GetUsage(tenantID)
	if err != nil {
		return nil, err
	}

	// Obtém limites
	limits, err := s.quotaRepo.GetLimits(tenantID)
	if err != nil {
		return nil, err
	}

	// Executa verificações
	checks := []*domain.QuotaCheck{
		usage.CheckProcessesQuota(limits.MaxProcesses),
		usage.CheckUsersQuota(limits.MaxUsers),
		usage.CheckClientsQuota(limits.MaxClients),
		usage.CheckDataJudDailyQuota(limits.DataJudQueriesDaily),
		usage.CheckAIMonthlyQuota(limits.AIQueriesMonthly),
		usage.CheckStorageQuota(limits.StorageGB),
		usage.CheckWebhooksQuota(limits.MaxWebhooks),
		usage.CheckAPIDailyQuota(limits.MaxAPICallsDaily),
	}

	// Separa warnings e exceeded
	var warnings, exceeded []*domain.QuotaCheck
	canProceed := true

	for _, check := range checks {
		if check.IsExceeded {
			exceeded = append(exceeded, check)
			canProceed = false
		} else if check.IsWarning {
			warnings = append(warnings, check)
		}
	}

	return &QuotaCheckResponse{
		TenantID:   tenantID,
		Checks:     checks,
		CanProceed: canProceed,
		Warnings:   warnings,
		Exceeded:   exceeded,
	}, nil
}

// CanIncrementQuota verifica se pode incrementar uma quota específica
func (s *QuotaService) CanIncrementQuota(ctx context.Context, tenantID, quotaType string, amount int) (bool, error) {
	// Obtém uso atual
	usage, err := s.quotaRepo.GetUsage(tenantID)
	if err != nil {
		return false, err
	}

	// Obtém limites
	limits, err := s.quotaRepo.GetLimits(tenantID)
	if err != nil {
		return false, err
	}

	// Verifica quota específica
	switch quotaType {
	case "processes":
		return usage.CanIncrementProcesses(limits.MaxProcesses, amount), nil
	case "users":
		return usage.CanIncrementUsers(limits.MaxUsers, amount), nil
	case "clients":
		return usage.CanIncrementClients(limits.MaxClients, amount), nil
	case "datajud_daily":
		return usage.CanIncrementDataJudDaily(limits.DataJudQueriesDaily, amount), nil
	case "ai_monthly":
		return usage.CanIncrementAIMonthly(limits.AIQueriesMonthly, amount), nil
	case "webhooks":
		return usage.CanIncrementWebhooks(limits.MaxWebhooks, amount), nil
	case "api_daily":
		return usage.CanIncrementAPIDaily(limits.MaxAPICallsDaily, amount), nil
	default:
		return false, domain.ErrInvalidQuotaType
	}
}

// IncrementQuota incrementa uma quota específica
func (s *QuotaService) IncrementQuota(ctx context.Context, tenantID string, req *IncrementQuotaRequest) error {
	s.logger.Info("Incrementing quota", 
		zap.String("tenant_id", tenantID),
		zap.String("quota_type", req.QuotaType),
		zap.Int("amount", req.Amount),
	)

	// Verifica se pode incrementar
	canIncrement, err := s.CanIncrementQuota(ctx, tenantID, req.QuotaType, req.Amount)
	if err != nil {
		return err
	}

	if !canIncrement {
		// Publica evento de quota excedida
		usage, _ := s.quotaRepo.GetUsage(tenantID)
		limits, _ := s.quotaRepo.GetLimits(tenantID)
		
		if usage != nil && limits != nil {
			var currentUsage, limit int
			switch req.QuotaType {
			case "processes":
				currentUsage, limit = usage.ProcessesCount, limits.MaxProcesses
			case "users":
				currentUsage, limit = usage.UsersCount, limits.MaxUsers
			case "clients":
				currentUsage, limit = usage.ClientsCount, limits.MaxClients
			case "datajud_daily":
				currentUsage, limit = usage.DataJudQueriesDaily, limits.DataJudQueriesDaily
			case "ai_monthly":
				currentUsage, limit = usage.AIQueriesMonthly, limits.AIQueriesMonthly
			case "webhooks":
				currentUsage, limit = usage.WebhooksCount, limits.MaxWebhooks
			case "api_daily":
				currentUsage, limit = usage.APICallsDaily, limits.MaxAPICallsDaily
			}

			event := domain.NewQuotaExceededEvent(tenantID, req.QuotaType, currentUsage, limit)
			if err := s.eventPublisher.Publish(event); err != nil {
				s.logger.Error("Failed to publish quota exceeded event", zap.Error(err))
			}
		}

		return domain.ErrQuotaExceeded
	}

	// Incrementa quota no repositório
	if err := s.quotaRepo.IncrementCounter(tenantID, req.QuotaType, req.Amount); err != nil {
		s.logger.Error("Failed to increment quota", zap.Error(err))
		return fmt.Errorf("erro ao incrementar quota: %w", err)
	}

	// Verifica se deve publicar warning
	if err := s.checkAndPublishWarning(tenantID, req.QuotaType); err != nil {
		s.logger.Error("Failed to check quota warning", zap.Error(err))
	}

	s.logger.Info("Quota incremented successfully", 
		zap.String("tenant_id", tenantID),
		zap.String("quota_type", req.QuotaType),
		zap.Int("amount", req.Amount),
	)

	return nil
}

// UpdateStorageUsage atualiza uso de armazenamento
func (s *QuotaService) UpdateStorageUsage(ctx context.Context, tenantID string, req *UpdateStorageRequest) error {
	s.logger.Info("Updating storage usage", 
		zap.String("tenant_id", tenantID),
		zap.Float64("usage_gb", req.UsageGB),
	)

	// Obtém uso atual
	usage, err := s.quotaRepo.GetUsage(tenantID)
	if err != nil {
		return err
	}

	// Atualiza armazenamento
	if err := usage.UpdateStorageUsage(req.UsageGB); err != nil {
		return err
	}

	// Salva alterações
	if err := s.quotaRepo.UpdateUsage(usage); err != nil {
		s.logger.Error("Failed to update storage usage", zap.Error(err))
		return fmt.Errorf("erro ao atualizar uso de armazenamento: %w", err)
	}

	// Verifica quota de armazenamento
	limits, err := s.quotaRepo.GetLimits(tenantID)
	if err == nil {
		check := usage.CheckStorageQuota(limits.StorageGB)
		if check.IsExceeded {
			event := domain.NewQuotaExceededEvent(tenantID, "storage", int(req.UsageGB*100), limits.StorageGB*100)
			if err := s.eventPublisher.Publish(event); err != nil {
				s.logger.Error("Failed to publish storage quota exceeded event", zap.Error(err))
			}
		} else if check.IsWarning {
			event := domain.NewQuotaWarningEvent(tenantID, "storage", int(req.UsageGB*100), limits.StorageGB*100, 80.0)
			if err := s.eventPublisher.Publish(event); err != nil {
				s.logger.Error("Failed to publish storage quota warning event", zap.Error(err))
			}
		}
	}

	s.logger.Info("Storage usage updated successfully", 
		zap.String("tenant_id", tenantID),
		zap.Float64("usage_gb", req.UsageGB),
	)

	return nil
}

// ResetDailyQuotas reseta quotas diárias
func (s *QuotaService) ResetDailyQuotas(ctx context.Context, tenantID string) error {
	s.logger.Info("Resetting daily quotas", zap.String("tenant_id", tenantID))

	if err := s.quotaRepo.ResetDailyCounters(tenantID); err != nil {
		s.logger.Error("Failed to reset daily quotas", zap.Error(err))
		return fmt.Errorf("erro ao resetar quotas diárias: %w", err)
	}

	s.logger.Info("Daily quotas reset successfully", zap.String("tenant_id", tenantID))
	return nil
}

// ResetMonthlyQuotas reseta quotas mensais
func (s *QuotaService) ResetMonthlyQuotas(ctx context.Context, tenantID string) error {
	s.logger.Info("Resetting monthly quotas", zap.String("tenant_id", tenantID))

	if err := s.quotaRepo.ResetMonthlyCounters(tenantID); err != nil {
		s.logger.Error("Failed to reset monthly quotas", zap.Error(err))
		return fmt.Errorf("erro ao resetar quotas mensais: %w", err)
	}

	s.logger.Info("Monthly quotas reset successfully", zap.String("tenant_id", tenantID))
	return nil
}

// UpdateQuotaLimits atualiza limites de quota (quando plano é alterado)
func (s *QuotaService) UpdateQuotaLimits(ctx context.Context, tenantID string, quotas domain.PlanQuotas) error {
	s.logger.Info("Updating quota limits", zap.String("tenant_id", tenantID))

	limits := domain.GetQuotaLimitFromPlan(tenantID, quotas)
	
	if err := s.quotaRepo.UpdateLimits(limits); err != nil {
		s.logger.Error("Failed to update quota limits", zap.Error(err))
		return fmt.Errorf("erro ao atualizar limites de quota: %w", err)
	}

	s.logger.Info("Quota limits updated successfully", zap.String("tenant_id", tenantID))
	return nil
}

// checkAndPublishWarning verifica se deve publicar warning de quota
func (s *QuotaService) checkAndPublishWarning(tenantID, quotaType string) error {
	usage, err := s.quotaRepo.GetUsage(tenantID)
	if err != nil {
		return err
	}

	limits, err := s.quotaRepo.GetLimits(tenantID)
	if err != nil {
		return err
	}

	var check *domain.QuotaCheck
	switch quotaType {
	case "processes":
		check = usage.CheckProcessesQuota(limits.MaxProcesses)
	case "users":
		check = usage.CheckUsersQuota(limits.MaxUsers)
	case "clients":
		check = usage.CheckClientsQuota(limits.MaxClients)
	case "datajud_daily":
		check = usage.CheckDataJudDailyQuota(limits.DataJudQueriesDaily)
	case "ai_monthly":
		check = usage.CheckAIMonthlyQuota(limits.AIQueriesMonthly)
	case "webhooks":
		check = usage.CheckWebhooksQuota(limits.MaxWebhooks)
	case "api_daily":
		check = usage.CheckAPIDailyQuota(limits.MaxAPICallsDaily)
	default:
		return nil
	}

	if check != nil && check.IsWarning {
		event := domain.NewQuotaWarningEvent(tenantID, quotaType, check.Current, check.Limit, 80.0)
		if err := s.eventPublisher.Publish(event); err != nil {
			return err
		}
	}

	return nil
}

// toQuotaUsageResponse converte domain.QuotaUsage para QuotaUsageResponse
func (s *QuotaService) toQuotaUsageResponse(usage *domain.QuotaUsage) *QuotaUsageResponse {
	return &QuotaUsageResponse{
		TenantID:              usage.TenantID,
		ProcessesCount:        usage.ProcessesCount,
		UsersCount:            usage.UsersCount,
		ClientsCount:          usage.ClientsCount,
		DataJudQueriesDaily:   usage.DataJudQueriesDaily,
		DataJudQueriesMonth:   usage.DataJudQueriesMonth,
		AIQueriesMonthly:      usage.AIQueriesMonthly,
		StorageUsedGB:         usage.StorageUsedGB,
		WebhooksCount:         usage.WebhooksCount,
		APICallsDaily:         usage.APICallsDaily,
		APICallsMonthly:       usage.APICallsMonthly,
		LastUpdated:           usage.LastUpdated,
		LastResetDaily:        usage.LastResetDaily,
		LastResetMonthly:      usage.LastResetMonthly,
	}
}

// toQuotaLimitResponse converte domain.QuotaLimit para QuotaLimitResponse
func (s *QuotaService) toQuotaLimitResponse(limits *domain.QuotaLimit) *QuotaLimitResponse {
	return &QuotaLimitResponse{
		TenantID:              limits.TenantID,
		MaxProcesses:          limits.MaxProcesses,
		MaxUsers:              limits.MaxUsers,
		MaxClients:            limits.MaxClients,
		DataJudQueriesDaily:   limits.DataJudQueriesDaily,
		AIQueriesMonthly:      limits.AIQueriesMonthly,
		StorageGB:             limits.StorageGB,
		MaxWebhooks:           limits.MaxWebhooks,
		MaxAPICallsDaily:      limits.MaxAPICallsDaily,
		UpdatedAt:             limits.UpdatedAt,
	}
}