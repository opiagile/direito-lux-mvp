package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// PostgresSubscriptionRepository implementação PostgreSQL do SubscriptionRepository
type PostgresSubscriptionRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresSubscriptionRepository cria nova instância do repositório
func NewPostgresSubscriptionRepository(db *sqlx.DB, logger *zap.Logger) domain.SubscriptionRepository {
	return &PostgresSubscriptionRepository{
		db:     db,
		logger: logger,
	}
}

// subscriptionDB representa uma assinatura no banco de dados
type subscriptionDB struct {
	ID                 string  `db:"id"`
	TenantID           string  `db:"tenant_id"`
	PlanID             string  `db:"plan_id"`
	Status             string  `db:"status"`
	CurrentPeriodStart string  `db:"current_period_start"`
	CurrentPeriodEnd   string  `db:"current_period_end"`
	CancelAtPeriodEnd  bool    `db:"cancel_at_period_end"`
	TrialStart         *string `db:"trial_start"`
	TrialEnd           *string `db:"trial_end"`
	CreatedAt          string  `db:"created_at"`
	UpdatedAt          string  `db:"updated_at"`
	CanceledAt         *string `db:"canceled_at"`
}

// Create cria uma nova assinatura
func (r *PostgresSubscriptionRepository) Create(subscription *domain.Subscription) error {
	r.logger.Debug("Creating subscription", zap.String("subscription_id", subscription.ID))

	query := `
		INSERT INTO subscriptions (
			id, tenant_id, plan_id, status, current_period_start, current_period_end,
			cancel_at_period_end, trial_start, trial_end, created_at, updated_at, canceled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`

	var trialStart, trialEnd, canceledAt *string
	if subscription.TrialStart != nil {
		trialStartStr := subscription.TrialStart.Format("2006-01-02 15:04:05")
		trialStart = &trialStartStr
	}
	if subscription.TrialEnd != nil {
		trialEndStr := subscription.TrialEnd.Format("2006-01-02 15:04:05")
		trialEnd = &trialEndStr
	}
	if subscription.CanceledAt != nil {
		canceledAtStr := subscription.CanceledAt.Format("2006-01-02 15:04:05")
		canceledAt = &canceledAtStr
	}

	_, err := r.db.Exec(
		query,
		subscription.ID,
		subscription.TenantID,
		subscription.PlanID,
		string(subscription.Status),
		subscription.CurrentPeriodStart.Format("2006-01-02 15:04:05"),
		subscription.CurrentPeriodEnd.Format("2006-01-02 15:04:05"),
		subscription.CancelAtPeriodEnd,
		trialStart,
		trialEnd,
		subscription.CreatedAt.Format("2006-01-02 15:04:05"),
		subscription.UpdatedAt.Format("2006-01-02 15:04:05"),
		canceledAt,
	)

	if err != nil {
		r.logger.Error("Failed to create subscription", zap.Error(err))
		return fmt.Errorf("erro ao criar assinatura: %w", err)
	}

	r.logger.Debug("Subscription created successfully", zap.String("subscription_id", subscription.ID))
	return nil
}

// GetByID busca assinatura por ID
func (r *PostgresSubscriptionRepository) GetByID(id string) (*domain.Subscription, error) {
	r.logger.Debug("Getting subscription by ID", zap.String("subscription_id", id))

	query := `
		SELECT id, tenant_id, plan_id, status, current_period_start, current_period_end,
		       cancel_at_period_end, trial_start, trial_end, created_at, updated_at, canceled_at
		FROM subscriptions
		WHERE id = $1`

	var subscriptionDB subscriptionDB
	err := r.db.Get(&subscriptionDB, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		r.logger.Error("Failed to get subscription", zap.Error(err))
		return nil, err
	}

	subscription, err := r.toDomainSubscription(&subscriptionDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Subscription found", zap.String("subscription_id", id))
	return subscription, nil
}

// GetByTenantID busca assinatura ativa por tenant ID
func (r *PostgresSubscriptionRepository) GetByTenantID(tenantID string) (*domain.Subscription, error) {
	r.logger.Debug("Getting subscription by tenant ID", zap.String("tenant_id", tenantID))

	query := `
		SELECT id, tenant_id, plan_id, status, current_period_start, current_period_end,
		       cancel_at_period_end, trial_start, trial_end, created_at, updated_at, canceled_at
		FROM subscriptions
		WHERE tenant_id = $1 AND status IN ('active', 'trialing', 'past_due')
		ORDER BY created_at DESC
		LIMIT 1`

	var subscriptionDB subscriptionDB
	err := r.db.Get(&subscriptionDB, query, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		r.logger.Error("Failed to get subscription by tenant", zap.Error(err))
		return nil, err
	}

	subscription, err := r.toDomainSubscription(&subscriptionDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Subscription found by tenant", zap.String("tenant_id", tenantID))
	return subscription, nil
}

// GetByStatus busca assinaturas por status
func (r *PostgresSubscriptionRepository) GetByStatus(status domain.SubscriptionStatus, limit, offset int) ([]*domain.Subscription, error) {
	r.logger.Debug("Getting subscriptions by status", 
		zap.String("status", string(status)),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	query := `
		SELECT id, tenant_id, plan_id, status, current_period_start, current_period_end,
		       cancel_at_period_end, trial_start, trial_end, created_at, updated_at, canceled_at
		FROM subscriptions
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	var subscriptionsDB []subscriptionDB
	err := r.db.Select(&subscriptionsDB, query, string(status), limit, offset)
	if err != nil {
		r.logger.Error("Failed to get subscriptions by status", zap.Error(err))
		return nil, err
	}

	var subscriptions []*domain.Subscription
	for _, subscriptionDB := range subscriptionsDB {
		subscription, err := r.toDomainSubscription(&subscriptionDB)
		if err != nil {
			r.logger.Error("Failed to convert subscription", zap.Error(err))
			continue
		}
		subscriptions = append(subscriptions, subscription)
	}

	r.logger.Debug("Subscriptions found by status", 
		zap.String("status", string(status)),
		zap.Int("count", len(subscriptions)),
	)
	return subscriptions, nil
}

// Update atualiza uma assinatura
func (r *PostgresSubscriptionRepository) Update(subscription *domain.Subscription) error {
	r.logger.Debug("Updating subscription", zap.String("subscription_id", subscription.ID))

	query := `
		UPDATE subscriptions SET
			tenant_id = $2,
			plan_id = $3,
			status = $4,
			current_period_start = $5,
			current_period_end = $6,
			cancel_at_period_end = $7,
			trial_start = $8,
			trial_end = $9,
			updated_at = $10,
			canceled_at = $11
		WHERE id = $1`

	var trialStart, trialEnd, canceledAt *string
	if subscription.TrialStart != nil {
		trialStartStr := subscription.TrialStart.Format("2006-01-02 15:04:05")
		trialStart = &trialStartStr
	}
	if subscription.TrialEnd != nil {
		trialEndStr := subscription.TrialEnd.Format("2006-01-02 15:04:05")
		trialEnd = &trialEndStr
	}
	if subscription.CanceledAt != nil {
		canceledAtStr := subscription.CanceledAt.Format("2006-01-02 15:04:05")
		canceledAt = &canceledAtStr
	}

	result, err := r.db.Exec(
		query,
		subscription.ID,
		subscription.TenantID,
		subscription.PlanID,
		string(subscription.Status),
		subscription.CurrentPeriodStart.Format("2006-01-02 15:04:05"),
		subscription.CurrentPeriodEnd.Format("2006-01-02 15:04:05"),
		subscription.CancelAtPeriodEnd,
		trialStart,
		trialEnd,
		subscription.UpdatedAt.Format("2006-01-02 15:04:05"),
		canceledAt,
	)

	if err != nil {
		r.logger.Error("Failed to update subscription", zap.Error(err))
		return fmt.Errorf("erro ao atualizar assinatura: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrSubscriptionNotFound
	}

	r.logger.Debug("Subscription updated successfully", zap.String("subscription_id", subscription.ID))
	return nil
}

// Delete exclui uma assinatura
func (r *PostgresSubscriptionRepository) Delete(id string) error {
	r.logger.Debug("Deleting subscription", zap.String("subscription_id", id))

	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete subscription", zap.Error(err))
		return fmt.Errorf("erro ao excluir assinatura: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrSubscriptionNotFound
	}

	r.logger.Debug("Subscription deleted successfully", zap.String("subscription_id", id))
	return nil
}

// GetExpiring busca assinaturas que vão expirar em X dias
func (r *PostgresSubscriptionRepository) GetExpiring(days int) ([]*domain.Subscription, error) {
	r.logger.Debug("Getting expiring subscriptions", zap.Int("days", days))

	query := `
		SELECT id, tenant_id, plan_id, status, current_period_start, current_period_end,
		       cancel_at_period_end, trial_start, trial_end, created_at, updated_at, canceled_at
		FROM subscriptions
		WHERE status IN ('active', 'trialing') 
		  AND current_period_end <= NOW() + INTERVAL '%d days'
		  AND current_period_end > NOW()
		ORDER BY current_period_end ASC`

	var subscriptionsDB []subscriptionDB
	err := r.db.Select(&subscriptionsDB, fmt.Sprintf(query, days))
	if err != nil {
		r.logger.Error("Failed to get expiring subscriptions", zap.Error(err))
		return nil, err
	}

	var subscriptions []*domain.Subscription
	for _, subscriptionDB := range subscriptionsDB {
		subscription, err := r.toDomainSubscription(&subscriptionDB)
		if err != nil {
			r.logger.Error("Failed to convert subscription", zap.Error(err))
			continue
		}
		subscriptions = append(subscriptions, subscription)
	}

	r.logger.Debug("Expiring subscriptions found", 
		zap.Int("days", days),
		zap.Int("count", len(subscriptions)),
	)
	return subscriptions, nil
}

// toDomainSubscription converte subscriptionDB para domain.Subscription
func (r *PostgresSubscriptionRepository) toDomainSubscription(subscriptionDB *subscriptionDB) (*domain.Subscription, error) {
	currentPeriodStart, err := parseTime(subscriptionDB.CurrentPeriodStart)
	if err != nil {
		return nil, err
	}

	currentPeriodEnd, err := parseTime(subscriptionDB.CurrentPeriodEnd)
	if err != nil {
		return nil, err
	}

	createdAt, err := parseTime(subscriptionDB.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := parseTime(subscriptionDB.UpdatedAt)
	if err != nil {
		return nil, err
	}

	var trialStart, trialEnd, canceledAt *time.Time
	if subscriptionDB.TrialStart != nil {
		parsed, err := parseTime(*subscriptionDB.TrialStart)
		if err != nil {
			return nil, err
		}
		trialStart = &parsed
	}

	if subscriptionDB.TrialEnd != nil {
		parsed, err := parseTime(*subscriptionDB.TrialEnd)
		if err != nil {
			return nil, err
		}
		trialEnd = &parsed
	}

	if subscriptionDB.CanceledAt != nil {
		parsed, err := parseTime(*subscriptionDB.CanceledAt)
		if err != nil {
			return nil, err
		}
		canceledAt = &parsed
	}

	return &domain.Subscription{
		ID:                 subscriptionDB.ID,
		TenantID:           subscriptionDB.TenantID,
		PlanID:             subscriptionDB.PlanID,
		Status:             domain.SubscriptionStatus(subscriptionDB.Status),
		CurrentPeriodStart: currentPeriodStart,
		CurrentPeriodEnd:   currentPeriodEnd,
		CancelAtPeriodEnd:  subscriptionDB.CancelAtPeriodEnd,
		TrialStart:         trialStart,
		TrialEnd:           trialEnd,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		CanceledAt:         canceledAt,
	}, nil
}