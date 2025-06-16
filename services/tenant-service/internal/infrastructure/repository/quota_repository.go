package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// PostgresQuotaRepository implementação PostgreSQL do QuotaRepository
type PostgresQuotaRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresQuotaRepository cria nova instância do repositório
func NewPostgresQuotaRepository(db *sqlx.DB, logger *zap.Logger) domain.QuotaRepository {
	return &PostgresQuotaRepository{
		db:     db,
		logger: logger,
	}
}

// quotaUsageDB representa uso de quota no banco de dados
type quotaUsageDB struct {
	ID                  string  `db:"id"`
	TenantID            string  `db:"tenant_id"`
	ProcessesCount      int     `db:"processes_count"`
	UsersCount          int     `db:"users_count"`
	ClientsCount        int     `db:"clients_count"`
	DataJudQueriesDaily int     `db:"datajud_queries_daily"`
	DataJudQueriesMonth int     `db:"datajud_queries_month"`
	AIQueriesMonthly    int     `db:"ai_queries_monthly"`
	StorageUsedGB       float64 `db:"storage_used_gb"`
	WebhooksCount       int     `db:"webhooks_count"`
	APICallsDaily       int     `db:"api_calls_daily"`
	APICallsMonthly     int     `db:"api_calls_monthly"`
	LastUpdated         string  `db:"last_updated"`
	LastResetDaily      string  `db:"last_reset_daily"`
	LastResetMonthly    string  `db:"last_reset_monthly"`
}

// quotaLimitDB representa limites de quota no banco de dados
type quotaLimitDB struct {
	TenantID              string `db:"tenant_id"`
	MaxProcesses          int    `db:"max_processes"`
	MaxUsers              int    `db:"max_users"`
	MaxClients            int    `db:"max_clients"`
	DataJudQueriesDaily   int    `db:"datajud_queries_daily"`
	AIQueriesMonthly      int    `db:"ai_queries_monthly"`
	StorageGB             int    `db:"storage_gb"`
	MaxWebhooks           int    `db:"max_webhooks"`
	MaxAPICallsDaily      int    `db:"max_api_calls_daily"`
	UpdatedAt             string `db:"updated_at"`
}

// GetUsage obtém uso atual de quotas
func (r *PostgresQuotaRepository) GetUsage(tenantID string) (*domain.QuotaUsage, error) {
	r.logger.Debug("Getting quota usage", zap.String("tenant_id", tenantID))

	query := `
		SELECT id, tenant_id, processes_count, users_count, clients_count,
		       datajud_queries_daily, datajud_queries_month, ai_queries_monthly,
		       storage_used_gb, webhooks_count, api_calls_daily, api_calls_monthly,
		       last_updated, last_reset_daily, last_reset_monthly
		FROM quota_usage
		WHERE tenant_id = $1`

	var usageDB quotaUsageDB
	err := r.db.Get(&usageDB, query, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrQuotaNotFound
		}
		r.logger.Error("Failed to get quota usage", zap.Error(err))
		return nil, err
	}

	usage, err := r.toDomainQuotaUsage(&usageDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Quota usage found", zap.String("tenant_id", tenantID))
	return usage, nil
}

// UpdateUsage atualiza uso de quotas
func (r *PostgresQuotaRepository) UpdateUsage(usage *domain.QuotaUsage) error {
	r.logger.Debug("Updating quota usage", zap.String("tenant_id", usage.TenantID))

	// Verifica se já existe
	_, err := r.GetUsage(usage.TenantID)
	if err == domain.ErrQuotaNotFound {
		// Cria novo registro
		return r.createUsage(usage)
	} else if err != nil {
		return err
	}

	// Atualiza registro existente
	query := `
		UPDATE quota_usage SET
			processes_count = $2,
			users_count = $3,
			clients_count = $4,
			datajud_queries_daily = $5,
			datajud_queries_month = $6,
			ai_queries_monthly = $7,
			storage_used_gb = $8,
			webhooks_count = $9,
			api_calls_daily = $10,
			api_calls_monthly = $11,
			last_updated = $12,
			last_reset_daily = $13,
			last_reset_monthly = $14
		WHERE tenant_id = $1`

	_, err = r.db.Exec(
		query,
		usage.TenantID,
		usage.ProcessesCount,
		usage.UsersCount,
		usage.ClientsCount,
		usage.DataJudQueriesDaily,
		usage.DataJudQueriesMonth,
		usage.AIQueriesMonthly,
		usage.StorageUsedGB,
		usage.WebhooksCount,
		usage.APICallsDaily,
		usage.APICallsMonthly,
		usage.LastUpdated.Format("2006-01-02 15:04:05"),
		usage.LastResetDaily.Format("2006-01-02 15:04:05"),
		usage.LastResetMonthly.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		r.logger.Error("Failed to update quota usage", zap.Error(err))
		return fmt.Errorf("erro ao atualizar uso de quota: %w", err)
	}

	r.logger.Debug("Quota usage updated successfully", zap.String("tenant_id", usage.TenantID))
	return nil
}

// IncrementCounter incrementa contador específico
func (r *PostgresQuotaRepository) IncrementCounter(tenantID, counterType string, amount int) error {
	r.logger.Debug("Incrementing counter", 
		zap.String("tenant_id", tenantID),
		zap.String("counter_type", counterType),
		zap.Int("amount", amount),
	)

	var query string
	switch counterType {
	case "processes":
		query = `UPDATE quota_usage SET processes_count = processes_count + $2, last_updated = NOW() WHERE tenant_id = $1`
	case "users":
		query = `UPDATE quota_usage SET users_count = users_count + $2, last_updated = NOW() WHERE tenant_id = $1`
	case "clients":
		query = `UPDATE quota_usage SET clients_count = clients_count + $2, last_updated = NOW() WHERE tenant_id = $1`
	case "datajud_daily":
		query = `UPDATE quota_usage SET datajud_queries_daily = datajud_queries_daily + $2, datajud_queries_month = datajud_queries_month + $2, last_updated = NOW() WHERE tenant_id = $1`
	case "ai_monthly":
		query = `UPDATE quota_usage SET ai_queries_monthly = ai_queries_monthly + $2, last_updated = NOW() WHERE tenant_id = $1`
	case "webhooks":
		query = `UPDATE quota_usage SET webhooks_count = webhooks_count + $2, last_updated = NOW() WHERE tenant_id = $1`
	case "api_daily":
		query = `UPDATE quota_usage SET api_calls_daily = api_calls_daily + $2, api_calls_monthly = api_calls_monthly + $2, last_updated = NOW() WHERE tenant_id = $1`
	default:
		return domain.ErrInvalidQuotaType
	}

	result, err := r.db.Exec(query, tenantID, amount)
	if err != nil {
		r.logger.Error("Failed to increment counter", zap.Error(err))
		return fmt.Errorf("erro ao incrementar contador: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrQuotaNotFound
	}

	r.logger.Debug("Counter incremented successfully", 
		zap.String("tenant_id", tenantID),
		zap.String("counter_type", counterType),
		zap.Int("amount", amount),
	)
	return nil
}

// ResetDailyCounters reseta contadores diários
func (r *PostgresQuotaRepository) ResetDailyCounters(tenantID string) error {
	r.logger.Debug("Resetting daily counters", zap.String("tenant_id", tenantID))

	query := `
		UPDATE quota_usage SET
			datajud_queries_daily = 0,
			api_calls_daily = 0,
			last_reset_daily = NOW(),
			last_updated = NOW()
		WHERE tenant_id = $1`

	result, err := r.db.Exec(query, tenantID)
	if err != nil {
		r.logger.Error("Failed to reset daily counters", zap.Error(err))
		return fmt.Errorf("erro ao resetar contadores diários: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrQuotaNotFound
	}

	r.logger.Debug("Daily counters reset successfully", zap.String("tenant_id", tenantID))
	return nil
}

// ResetMonthlyCounters reseta contadores mensais
func (r *PostgresQuotaRepository) ResetMonthlyCounters(tenantID string) error {
	r.logger.Debug("Resetting monthly counters", zap.String("tenant_id", tenantID))

	query := `
		UPDATE quota_usage SET
			datajud_queries_month = 0,
			ai_queries_monthly = 0,
			api_calls_monthly = 0,
			last_reset_monthly = NOW(),
			last_updated = NOW()
		WHERE tenant_id = $1`

	result, err := r.db.Exec(query, tenantID)
	if err != nil {
		r.logger.Error("Failed to reset monthly counters", zap.Error(err))
		return fmt.Errorf("erro ao resetar contadores mensais: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrQuotaNotFound
	}

	r.logger.Debug("Monthly counters reset successfully", zap.String("tenant_id", tenantID))
	return nil
}

// GetLimits obtém limites de quota
func (r *PostgresQuotaRepository) GetLimits(tenantID string) (*domain.QuotaLimit, error) {
	r.logger.Debug("Getting quota limits", zap.String("tenant_id", tenantID))

	query := `
		SELECT tenant_id, max_processes, max_users, max_clients,
		       datajud_queries_daily, ai_queries_monthly, storage_gb,
		       max_webhooks, max_api_calls_daily, updated_at
		FROM quota_limits
		WHERE tenant_id = $1`

	var limitDB quotaLimitDB
	err := r.db.Get(&limitDB, query, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrQuotaNotFound
		}
		r.logger.Error("Failed to get quota limits", zap.Error(err))
		return nil, err
	}

	limit, err := r.toDomainQuotaLimit(&limitDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Quota limits found", zap.String("tenant_id", tenantID))
	return limit, nil
}

// UpdateLimits atualiza limites de quota
func (r *PostgresQuotaRepository) UpdateLimits(limits *domain.QuotaLimit) error {
	r.logger.Debug("Updating quota limits", zap.String("tenant_id", limits.TenantID))

	// Verifica se já existe
	_, err := r.GetLimits(limits.TenantID)
	if err == domain.ErrQuotaNotFound {
		// Cria novo registro
		return r.createLimits(limits)
	} else if err != nil {
		return err
	}

	// Atualiza registro existente
	query := `
		UPDATE quota_limits SET
			max_processes = $2,
			max_users = $3,
			max_clients = $4,
			datajud_queries_daily = $5,
			ai_queries_monthly = $6,
			storage_gb = $7,
			max_webhooks = $8,
			max_api_calls_daily = $9,
			updated_at = $10
		WHERE tenant_id = $1`

	_, err = r.db.Exec(
		query,
		limits.TenantID,
		limits.MaxProcesses,
		limits.MaxUsers,
		limits.MaxClients,
		limits.DataJudQueriesDaily,
		limits.AIQueriesMonthly,
		limits.StorageGB,
		limits.MaxWebhooks,
		limits.MaxAPICallsDaily,
		limits.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		r.logger.Error("Failed to update quota limits", zap.Error(err))
		return fmt.Errorf("erro ao atualizar limites de quota: %w", err)
	}

	r.logger.Debug("Quota limits updated successfully", zap.String("tenant_id", limits.TenantID))
	return nil
}

// createUsage cria novo registro de uso de quota
func (r *PostgresQuotaRepository) createUsage(usage *domain.QuotaUsage) error {
	query := `
		INSERT INTO quota_usage (
			id, tenant_id, processes_count, users_count, clients_count,
			datajud_queries_daily, datajud_queries_month, ai_queries_monthly,
			storage_used_gb, webhooks_count, api_calls_daily, api_calls_monthly,
			last_updated, last_reset_daily, last_reset_monthly
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	_, err := r.db.Exec(
		query,
		usage.ID,
		usage.TenantID,
		usage.ProcessesCount,
		usage.UsersCount,
		usage.ClientsCount,
		usage.DataJudQueriesDaily,
		usage.DataJudQueriesMonth,
		usage.AIQueriesMonthly,
		usage.StorageUsedGB,
		usage.WebhooksCount,
		usage.APICallsDaily,
		usage.APICallsMonthly,
		usage.LastUpdated.Format("2006-01-02 15:04:05"),
		usage.LastResetDaily.Format("2006-01-02 15:04:05"),
		usage.LastResetMonthly.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		r.logger.Error("Failed to create quota usage", zap.Error(err))
		return fmt.Errorf("erro ao criar uso de quota: %w", err)
	}

	return nil
}

// createLimits cria novo registro de limites de quota
func (r *PostgresQuotaRepository) createLimits(limits *domain.QuotaLimit) error {
	query := `
		INSERT INTO quota_limits (
			tenant_id, max_processes, max_users, max_clients,
			datajud_queries_daily, ai_queries_monthly, storage_gb,
			max_webhooks, max_api_calls_daily, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	_, err := r.db.Exec(
		query,
		limits.TenantID,
		limits.MaxProcesses,
		limits.MaxUsers,
		limits.MaxClients,
		limits.DataJudQueriesDaily,
		limits.AIQueriesMonthly,
		limits.StorageGB,
		limits.MaxWebhooks,
		limits.MaxAPICallsDaily,
		limits.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		r.logger.Error("Failed to create quota limits", zap.Error(err))
		return fmt.Errorf("erro ao criar limites de quota: %w", err)
	}

	return nil
}

// toDomainQuotaUsage converte quotaUsageDB para domain.QuotaUsage
func (r *PostgresQuotaRepository) toDomainQuotaUsage(usageDB *quotaUsageDB) (*domain.QuotaUsage, error) {
	lastUpdated, err := parseTime(usageDB.LastUpdated)
	if err != nil {
		return nil, err
	}

	lastResetDaily, err := parseTime(usageDB.LastResetDaily)
	if err != nil {
		return nil, err
	}

	lastResetMonthly, err := parseTime(usageDB.LastResetMonthly)
	if err != nil {
		return nil, err
	}

	return &domain.QuotaUsage{
		ID:                  usageDB.ID,
		TenantID:            usageDB.TenantID,
		ProcessesCount:      usageDB.ProcessesCount,
		UsersCount:          usageDB.UsersCount,
		ClientsCount:        usageDB.ClientsCount,
		DataJudQueriesDaily: usageDB.DataJudQueriesDaily,
		DataJudQueriesMonth: usageDB.DataJudQueriesMonth,
		AIQueriesMonthly:    usageDB.AIQueriesMonthly,
		StorageUsedGB:       usageDB.StorageUsedGB,
		WebhooksCount:       usageDB.WebhooksCount,
		APICallsDaily:       usageDB.APICallsDaily,
		APICallsMonthly:     usageDB.APICallsMonthly,
		LastUpdated:         lastUpdated,
		LastResetDaily:      lastResetDaily,
		LastResetMonthly:    lastResetMonthly,
	}, nil
}

// toDomainQuotaLimit converte quotaLimitDB para domain.QuotaLimit
func (r *PostgresQuotaRepository) toDomainQuotaLimit(limitDB *quotaLimitDB) (*domain.QuotaLimit, error) {
	updatedAt, err := parseTime(limitDB.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.QuotaLimit{
		TenantID:              limitDB.TenantID,
		MaxProcesses:          limitDB.MaxProcesses,
		MaxUsers:              limitDB.MaxUsers,
		MaxClients:            limitDB.MaxClients,
		DataJudQueriesDaily:   limitDB.DataJudQueriesDaily,
		AIQueriesMonthly:      limitDB.AIQueriesMonthly,
		StorageGB:             limitDB.StorageGB,
		MaxWebhooks:           limitDB.MaxWebhooks,
		MaxAPICallsDaily:      limitDB.MaxAPICallsDaily,
		UpdatedAt:             updatedAt,
	}, nil
}