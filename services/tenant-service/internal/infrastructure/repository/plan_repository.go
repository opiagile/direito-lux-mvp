package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// PostgresPlanRepository implementação PostgreSQL do PlanRepository
type PostgresPlanRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresPlanRepository cria nova instância do repositório
func NewPostgresPlanRepository(db *sqlx.DB, logger *zap.Logger) domain.PlanRepository {
	return &PostgresPlanRepository{
		db:     db,
		logger: logger,
	}
}

// planDB representa um plano no banco de dados
type planDB struct {
	ID              string `db:"id"`
	Name            string `db:"name"`
	Type            string `db:"type"`
	Description     string `db:"description"`
	Price           int64  `db:"price"`
	Currency        string `db:"currency"`
	BillingInterval string `db:"billing_interval"`
	FeaturesJSON    string `db:"features"`
	QuotasJSON      string `db:"quotas"`
	IsActive        bool   `db:"is_active"`
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
}

// Create cria um novo plano
func (r *PostgresPlanRepository) Create(plan *domain.Plan) error {
	r.logger.Debug("Creating plan", zap.String("plan_id", plan.ID))

	featuresJSON, err := json.Marshal(plan.Features)
	if err != nil {
		return fmt.Errorf("erro ao serializar features: %w", err)
	}

	quotasJSON, err := json.Marshal(plan.Quotas)
	if err != nil {
		return fmt.Errorf("erro ao serializar quotas: %w", err)
	}

	query := `
		INSERT INTO plans (
			id, name, type, description, price, currency, billing_interval,
			features, quotas, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`

	_, err = r.db.Exec(
		query,
		plan.ID,
		plan.Name,
		string(plan.Type),
		plan.Description,
		plan.Price,
		plan.Currency,
		string(plan.BillingInterval),
		string(featuresJSON),
		string(quotasJSON),
		plan.IsActive,
		plan.CreatedAt.Format("2006-01-02 15:04:05"),
		plan.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		r.logger.Error("Failed to create plan", zap.Error(err))
		return fmt.Errorf("erro ao criar plano: %w", err)
	}

	r.logger.Debug("Plan created successfully", zap.String("plan_id", plan.ID))
	return nil
}

// GetByID busca plano por ID
func (r *PostgresPlanRepository) GetByID(id string) (*domain.Plan, error) {
	r.logger.Debug("Getting plan by ID", zap.String("plan_id", id))

	query := `
		SELECT id, name, type, description, price, currency, billing_interval,
		       features, quotas, is_active, created_at, updated_at
		FROM plans
		WHERE id = $1`

	var planDB planDB
	err := r.db.Get(&planDB, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPlanNotFound
		}
		r.logger.Error("Failed to get plan", zap.Error(err))
		return nil, err
	}

	plan, err := r.toDomainPlan(&planDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Plan found", zap.String("plan_id", id))
	return plan, nil
}

// GetByType busca plano por tipo
func (r *PostgresPlanRepository) GetByType(planType domain.PlanType) (*domain.Plan, error) {
	r.logger.Debug("Getting plan by type", zap.String("plan_type", string(planType)))

	query := `
		SELECT id, name, type, description, price, currency, billing_interval,
		       features, quotas, is_active, created_at, updated_at
		FROM plans
		WHERE type = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT 1`

	var planDB planDB
	err := r.db.Get(&planDB, query, string(planType))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPlanNotFound
		}
		r.logger.Error("Failed to get plan by type", zap.Error(err))
		return nil, err
	}

	plan, err := r.toDomainPlan(&planDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Plan found by type", zap.String("plan_type", string(planType)))
	return plan, nil
}

// GetAll busca todos os planos
func (r *PostgresPlanRepository) GetAll(activeOnly bool) ([]*domain.Plan, error) {
	r.logger.Debug("Getting all plans", zap.Bool("active_only", activeOnly))

	query := `
		SELECT id, name, type, description, price, currency, billing_interval,
		       features, quotas, is_active, created_at, updated_at
		FROM plans`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY price ASC"

	var plansDB []planDB
	err := r.db.Select(&plansDB, query)
	if err != nil {
		r.logger.Error("Failed to get all plans", zap.Error(err))
		return nil, err
	}

	var plans []*domain.Plan
	for _, planDB := range plansDB {
		plan, err := r.toDomainPlan(&planDB)
		if err != nil {
			r.logger.Error("Failed to convert plan", zap.Error(err))
			continue
		}
		plans = append(plans, plan)
	}

	r.logger.Debug("All plans retrieved", 
		zap.Bool("active_only", activeOnly),
		zap.Int("count", len(plans)),
	)
	return plans, nil
}

// Update atualiza um plano
func (r *PostgresPlanRepository) Update(plan *domain.Plan) error {
	r.logger.Debug("Updating plan", zap.String("plan_id", plan.ID))

	featuresJSON, err := json.Marshal(plan.Features)
	if err != nil {
		return fmt.Errorf("erro ao serializar features: %w", err)
	}

	quotasJSON, err := json.Marshal(plan.Quotas)
	if err != nil {
		return fmt.Errorf("erro ao serializar quotas: %w", err)
	}

	query := `
		UPDATE plans SET
			name = $2,
			type = $3,
			description = $4,
			price = $5,
			currency = $6,
			billing_interval = $7,
			features = $8,
			quotas = $9,
			is_active = $10,
			updated_at = $11
		WHERE id = $1`

	result, err := r.db.Exec(
		query,
		plan.ID,
		plan.Name,
		string(plan.Type),
		plan.Description,
		plan.Price,
		plan.Currency,
		string(plan.BillingInterval),
		string(featuresJSON),
		string(quotasJSON),
		plan.IsActive,
		plan.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		r.logger.Error("Failed to update plan", zap.Error(err))
		return fmt.Errorf("erro ao atualizar plano: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrPlanNotFound
	}

	r.logger.Debug("Plan updated successfully", zap.String("plan_id", plan.ID))
	return nil
}

// Delete exclui um plano (marca como inativo)
func (r *PostgresPlanRepository) Delete(id string) error {
	r.logger.Debug("Deleting plan", zap.String("plan_id", id))

	query := `UPDATE plans SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete plan", zap.Error(err))
		return fmt.Errorf("erro ao excluir plano: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrPlanNotFound
	}

	r.logger.Debug("Plan deleted successfully", zap.String("plan_id", id))
	return nil
}

// toDomainPlan converte planDB para domain.Plan
func (r *PostgresPlanRepository) toDomainPlan(planDB *planDB) (*domain.Plan, error) {
	var features domain.PlanFeatures
	if err := json.Unmarshal([]byte(planDB.FeaturesJSON), &features); err != nil {
		r.logger.Error("Failed to unmarshal features", zap.Error(err))
		return nil, err
	}

	var quotas domain.PlanQuotas
	if err := json.Unmarshal([]byte(planDB.QuotasJSON), &quotas); err != nil {
		r.logger.Error("Failed to unmarshal quotas", zap.Error(err))
		return nil, err
	}

	createdAt, err := parseTime(planDB.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := parseTime(planDB.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Plan{
		ID:              planDB.ID,
		Name:            planDB.Name,
		Type:            domain.PlanType(planDB.Type),
		Description:     planDB.Description,
		Price:           planDB.Price,
		Currency:        planDB.Currency,
		BillingInterval: domain.BillingInterval(planDB.BillingInterval),
		Features:        features,
		Quotas:          quotas,
		IsActive:        planDB.IsActive,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}