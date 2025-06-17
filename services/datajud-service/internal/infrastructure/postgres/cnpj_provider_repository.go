package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// CNPJProviderRepository implementação PostgreSQL do repositório de CNPJProvider
type CNPJProviderRepository struct {
	db *sqlx.DB
}

// NewCNPJProviderRepository cria nova instância do repositório
func NewCNPJProviderRepository(db *sqlx.DB) *CNPJProviderRepository {
	return &CNPJProviderRepository{db: db}
}

// Save salva um CNPJ provider
func (r *CNPJProviderRepository) Save(provider *domain.CNPJProvider) error {
	query := `
		INSERT INTO cnpj_providers (
			id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			daily_limit, daily_usage, usage_reset_time, is_active, priority,
			last_used_at, created_at, updated_at, deactivated_at
		) VALUES (
			:id, :tenant_id, :cnpj, :name, :email, :api_key, :certificate, :certificate_pass,
			:daily_limit, :daily_usage, :usage_reset_time, :is_active, :priority,
			:last_used_at, :created_at, :updated_at, :deactivated_at
		)`

	_, err := r.db.NamedExec(query, provider)
	return err
}

// FindByID encontra provider por ID
func (r *CNPJProviderRepository) FindByID(id uuid.UUID) (*domain.CNPJProvider, error) {
	query := `
		SELECT id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			   daily_limit, daily_usage, usage_reset_time, is_active, priority,
			   last_used_at, created_at, updated_at, deactivated_at
		FROM cnpj_providers 
		WHERE id = $1`

	provider := &domain.CNPJProvider{}
	err := r.db.Get(provider, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return provider, nil
}

// FindByTenantID encontra providers de um tenant
func (r *CNPJProviderRepository) FindByTenantID(tenantID uuid.UUID) ([]*domain.CNPJProvider, error) {
	query := `
		SELECT id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			   daily_limit, daily_usage, usage_reset_time, is_active, priority,
			   last_used_at, created_at, updated_at, deactivated_at
		FROM cnpj_providers 
		WHERE tenant_id = $1
		ORDER BY priority ASC, created_at ASC`

	providers := []*domain.CNPJProvider{}
	err := r.db.Select(&providers, query, tenantID)
	return providers, err
}

// FindActiveCNPJs encontra CNPJs ativos
func (r *CNPJProviderRepository) FindActiveCNPJs() ([]*domain.CNPJProvider, error) {
	query := `
		SELECT id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			   daily_limit, daily_usage, usage_reset_time, is_active, priority,
			   last_used_at, created_at, updated_at, deactivated_at
		FROM cnpj_providers 
		WHERE is_active = true
		ORDER BY priority ASC, last_used_at ASC NULLS FIRST`

	providers := []*domain.CNPJProvider{}
	err := r.db.Select(&providers, query)
	return providers, err
}

// FindAvailableCNPJs encontra CNPJs com quota disponível
func (r *CNPJProviderRepository) FindAvailableCNPJs(minQuota int) ([]*domain.CNPJProvider, error) {
	query := `
		SELECT id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			   daily_limit, daily_usage, usage_reset_time, is_active, priority,
			   last_used_at, created_at, updated_at, deactivated_at
		FROM cnpj_providers 
		WHERE is_active = true 
		  AND (daily_limit - daily_usage) >= $1
		  AND (usage_reset_time > NOW() OR usage_reset_time IS NULL)
		ORDER BY priority ASC, (daily_limit - daily_usage) DESC`

	providers := []*domain.CNPJProvider{}
	err := r.db.Select(&providers, query, minQuota)
	return providers, err
}

// UpdateUsage atualiza o uso de um provider
func (r *CNPJProviderRepository) UpdateUsage(id uuid.UUID, usage int) error {
	query := `
		UPDATE cnpj_providers 
		SET daily_usage = $2,
			last_used_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Exec(query, id, usage)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider %s não encontrado", id)
	}

	return nil
}

// ResetDailyUsage reseta o uso diário de todos os providers
func (r *CNPJProviderRepository) ResetDailyUsage() error {
	query := `
		UPDATE cnpj_providers 
		SET daily_usage = 0,
			usage_reset_time = DATE_TRUNC('day', NOW()) + INTERVAL '1 day',
			updated_at = NOW()
		WHERE usage_reset_time <= NOW() OR usage_reset_time IS NULL`

	_, err := r.db.Exec(query)
	return err
}

// Update atualiza um provider
func (r *CNPJProviderRepository) Update(provider *domain.CNPJProvider) error {
	query := `
		UPDATE cnpj_providers 
		SET tenant_id = :tenant_id,
			cnpj = :cnpj,
			name = :name,
			email = :email,
			api_key = :api_key,
			certificate = :certificate,
			certificate_pass = :certificate_pass,
			daily_limit = :daily_limit,
			daily_usage = :daily_usage,
			usage_reset_time = :usage_reset_time,
			is_active = :is_active,
			priority = :priority,
			last_used_at = :last_used_at,
			updated_at = :updated_at,
			deactivated_at = :deactivated_at
		WHERE id = :id`

	result, err := r.db.NamedExec(query, provider)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider %s não encontrado", provider.ID)
	}

	return nil
}

// Delete remove um provider
func (r *CNPJProviderRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM cnpj_providers WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider %s não encontrado", id)
	}

	return nil
}

// FindByCNPJ encontra provider por CNPJ
func (r *CNPJProviderRepository) FindByCNPJ(cnpj string) (*domain.CNPJProvider, error) {
	query := `
		SELECT id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			   daily_limit, daily_usage, usage_reset_time, is_active, priority,
			   last_used_at, created_at, updated_at, deactivated_at
		FROM cnpj_providers 
		WHERE cnpj = $1`

	provider := &domain.CNPJProvider{}
	err := r.db.Get(provider, query, cnpj)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return provider, nil
}

// GetUsageStats obtém estatísticas de uso
func (r *CNPJProviderRepository) GetUsageStats() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_providers,
			COUNT(CASE WHEN is_active THEN 1 END) as active_providers,
			SUM(daily_limit) as total_daily_limit,
			SUM(daily_usage) as total_daily_usage,
			AVG(daily_usage::float / NULLIF(daily_limit, 0) * 100) as avg_usage_percentage,
			COUNT(CASE WHEN daily_usage >= daily_limit THEN 1 END) as exhausted_providers,
			MIN(last_used_at) as oldest_usage,
			MAX(last_used_at) as newest_usage
		FROM cnpj_providers`

	stats := make(map[string]interface{})
	row := r.db.QueryRow(query)

	var totalProviders, activeProviders, totalDailyLimit, totalDailyUsage, exhaustedProviders int
	var avgUsagePercentage sql.NullFloat64
	var oldestUsage, newestUsage sql.NullTime

	err := row.Scan(
		&totalProviders, &activeProviders, &totalDailyLimit, &totalDailyUsage,
		&avgUsagePercentage, &exhaustedProviders, &oldestUsage, &newestUsage,
	)
	if err != nil {
		return nil, err
	}

	stats["total_providers"] = totalProviders
	stats["active_providers"] = activeProviders
	stats["total_daily_limit"] = totalDailyLimit
	stats["total_daily_usage"] = totalDailyUsage
	stats["available_quota"] = totalDailyLimit - totalDailyUsage
	stats["exhausted_providers"] = exhaustedProviders

	if avgUsagePercentage.Valid {
		stats["avg_usage_percentage"] = avgUsagePercentage.Float64
	} else {
		stats["avg_usage_percentage"] = 0.0
	}

	if oldestUsage.Valid {
		stats["oldest_usage"] = oldestUsage.Time
	}

	if newestUsage.Valid {
		stats["newest_usage"] = newestUsage.Time
	}

	return stats, nil
}

// GetTopUsedProviders obtém providers mais utilizados
func (r *CNPJProviderRepository) GetTopUsedProviders(limit int) ([]*domain.CNPJProvider, error) {
	query := `
		SELECT id, tenant_id, cnpj, name, email, api_key, certificate, certificate_pass,
			   daily_limit, daily_usage, usage_reset_time, is_active, priority,
			   last_used_at, created_at, updated_at, deactivated_at
		FROM cnpj_providers 
		WHERE is_active = true
		ORDER BY daily_usage DESC, last_used_at DESC
		LIMIT $1`

	providers := []*domain.CNPJProvider{}
	err := r.db.Select(&providers, query, limit)
	return providers, err
}

// ActivateProvider ativa um provider
func (r *CNPJProviderRepository) ActivateProvider(id uuid.UUID) error {
	query := `
		UPDATE cnpj_providers 
		SET is_active = true,
			deactivated_at = NULL,
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider %s não encontrado", id)
	}

	return nil
}

// DeactivateProvider desativa um provider
func (r *CNPJProviderRepository) DeactivateProvider(id uuid.UUID) error {
	query := `
		UPDATE cnpj_providers 
		SET is_active = false,
			deactivated_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider %s não encontrado", id)
	}

	return nil
}

// CleanupOldProviders remove providers inativos há muito tempo
func (r *CNPJProviderRepository) CleanupOldProviders(olderThan time.Time) (int, error) {
	query := `
		DELETE FROM cnpj_providers 
		WHERE is_active = false 
		  AND deactivated_at < $1`

	result, err := r.db.Exec(query, olderThan)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}