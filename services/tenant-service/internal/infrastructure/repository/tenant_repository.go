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

// PostgresTenantRepository implementação PostgreSQL do TenantRepository
type PostgresTenantRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresTenantRepository cria nova instância do repositório
func NewPostgresTenantRepository(db *sqlx.DB, logger *zap.Logger) domain.TenantRepository {
	return &PostgresTenantRepository{
		db:     db,
		logger: logger,
	}
}

// tenantDB representa um tenant no banco de dados
type tenantDB struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	LegalName     string `db:"legal_name"`
	Document      string `db:"document"`
	Email         string `db:"email"`
	Phone         string `db:"phone"`
	Website       string `db:"website"`
	AddressJSON   string `db:"address"`
	Status        string `db:"status"`
	PlanType      string `db:"plan_type"`
	OwnerUserID   string `db:"owner_user_id"`
	SettingsJSON  string `db:"settings"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
	ActivatedAt   *string `db:"activated_at"`
	SuspendedAt   *string `db:"suspended_at"`
}

// Create cria um novo tenant
func (r *PostgresTenantRepository) Create(tenant *domain.Tenant) error {
	r.logger.Debug("Creating tenant", zap.String("tenant_id", tenant.ID))

	addressJSON, err := json.Marshal(tenant.Address)
	if err != nil {
		return fmt.Errorf("erro ao serializar endereço: %w", err)
	}

	settingsJSON, err := json.Marshal(tenant.Settings)
	if err != nil {
		return fmt.Errorf("erro ao serializar configurações: %w", err)
	}

	query := `
		INSERT INTO tenants (
			id, name, legal_name, document, email, phone, website, 
			address, status, plan_type, owner_user_id, settings,
			created_at, updated_at, activated_at, suspended_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)`

	var activatedAt, suspendedAt *string
	if tenant.ActivatedAt != nil {
		activatedAtStr := tenant.ActivatedAt.Format("2006-01-02 15:04:05")
		activatedAt = &activatedAtStr
	}
	if tenant.SuspendedAt != nil {
		suspendedAtStr := tenant.SuspendedAt.Format("2006-01-02 15:04:05")
		suspendedAt = &suspendedAtStr
	}

	_, err = r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
		tenant.LegalName,
		tenant.Document,
		tenant.Email,
		tenant.Phone,
		tenant.Website,
		string(addressJSON),
		string(tenant.Status),
		string(tenant.PlanType),
		tenant.OwnerUserID,
		string(settingsJSON),
		tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
		activatedAt,
		suspendedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create tenant", zap.Error(err))
		return fmt.Errorf("erro ao criar tenant: %w", err)
	}

	r.logger.Debug("Tenant created successfully", zap.String("tenant_id", tenant.ID))
	return nil
}

// GetByID busca tenant por ID
func (r *PostgresTenantRepository) GetByID(id string) (*domain.Tenant, error) {
	r.logger.Debug("Getting tenant by ID", zap.String("tenant_id", id))

	query := `
		SELECT id, name, legal_name, document, email, phone, website,
		       address, status, plan_type, owner_user_id, settings,
		       created_at, updated_at, activated_at, suspended_at
		FROM tenants
		WHERE id = $1`

	var tenantDB tenantDB
	err := r.db.Get(&tenantDB, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTenantNotFound
		}
		r.logger.Error("Failed to get tenant", zap.Error(err))
		return nil, err
	}

	tenant, err := r.toDomainTenant(&tenantDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Tenant found", zap.String("tenant_id", id))
	return tenant, nil
}

// GetByDocument busca tenant por documento
func (r *PostgresTenantRepository) GetByDocument(document string) (*domain.Tenant, error) {
	r.logger.Debug("Getting tenant by document", zap.String("document", document))

	query := `
		SELECT id, name, legal_name, document, email, phone, website,
		       address, status, plan_type, owner_user_id, settings,
		       created_at, updated_at, activated_at, suspended_at
		FROM tenants
		WHERE document = $1`

	var tenantDB tenantDB
	err := r.db.Get(&tenantDB, query, document)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTenantNotFound
		}
		r.logger.Error("Failed to get tenant by document", zap.Error(err))
		return nil, err
	}

	tenant, err := r.toDomainTenant(&tenantDB)
	if err != nil {
		return nil, err
	}

	r.logger.Debug("Tenant found by document", zap.String("document", document))
	return tenant, nil
}

// GetByOwner busca tenants por proprietário
func (r *PostgresTenantRepository) GetByOwner(ownerUserID string) ([]*domain.Tenant, error) {
	r.logger.Debug("Getting tenants by owner", zap.String("owner_user_id", ownerUserID))

	query := `
		SELECT id, name, legal_name, document, email, phone, website,
		       address, status, plan_type, owner_user_id, settings,
		       created_at, updated_at, activated_at, suspended_at
		FROM tenants
		WHERE owner_user_id = $1
		ORDER BY created_at DESC`

	var tenantsDB []tenantDB
	err := r.db.Select(&tenantsDB, query, ownerUserID)
	if err != nil {
		r.logger.Error("Failed to get tenants by owner", zap.Error(err))
		return nil, err
	}

	var tenants []*domain.Tenant
	for _, tenantDB := range tenantsDB {
		tenant, err := r.toDomainTenant(&tenantDB)
		if err != nil {
			r.logger.Error("Failed to convert tenant", zap.Error(err))
			continue
		}
		tenants = append(tenants, tenant)
	}

	r.logger.Debug("Tenants found by owner", 
		zap.String("owner_user_id", ownerUserID),
		zap.Int("count", len(tenants)),
	)
	return tenants, nil
}

// GetAll busca todos os tenants com paginação
func (r *PostgresTenantRepository) GetAll(limit, offset int) ([]*domain.Tenant, error) {
	r.logger.Debug("Getting all tenants", zap.Int("limit", limit), zap.Int("offset", offset))

	query := `
		SELECT id, name, legal_name, document, email, phone, website,
		       address, status, plan_type, owner_user_id, settings,
		       created_at, updated_at, activated_at, suspended_at
		FROM tenants
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	var tenantsDB []tenantDB
	err := r.db.Select(&tenantsDB, query, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get all tenants", zap.Error(err))
		return nil, err
	}

	var tenants []*domain.Tenant
	for _, tenantDB := range tenantsDB {
		tenant, err := r.toDomainTenant(&tenantDB)
		if err != nil {
			r.logger.Error("Failed to convert tenant", zap.Error(err))
			continue
		}
		tenants = append(tenants, tenant)
	}

	r.logger.Debug("All tenants retrieved", zap.Int("count", len(tenants)))
	return tenants, nil
}

// Update atualiza um tenant
func (r *PostgresTenantRepository) Update(tenant *domain.Tenant) error {
	r.logger.Debug("Updating tenant", zap.String("tenant_id", tenant.ID))

	addressJSON, err := json.Marshal(tenant.Address)
	if err != nil {
		return fmt.Errorf("erro ao serializar endereço: %w", err)
	}

	settingsJSON, err := json.Marshal(tenant.Settings)
	if err != nil {
		return fmt.Errorf("erro ao serializar configurações: %w", err)
	}

	query := `
		UPDATE tenants SET
			name = $2,
			legal_name = $3,
			document = $4,
			email = $5,
			phone = $6,
			website = $7,
			address = $8,
			status = $9,
			plan_type = $10,
			owner_user_id = $11,
			settings = $12,
			updated_at = $13,
			activated_at = $14,
			suspended_at = $15
		WHERE id = $1`

	var activatedAt, suspendedAt *string
	if tenant.ActivatedAt != nil {
		activatedAtStr := tenant.ActivatedAt.Format("2006-01-02 15:04:05")
		activatedAt = &activatedAtStr
	}
	if tenant.SuspendedAt != nil {
		suspendedAtStr := tenant.SuspendedAt.Format("2006-01-02 15:04:05")
		suspendedAt = &suspendedAtStr
	}

	result, err := r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
		tenant.LegalName,
		tenant.Document,
		tenant.Email,
		tenant.Phone,
		tenant.Website,
		string(addressJSON),
		string(tenant.Status),
		string(tenant.PlanType),
		tenant.OwnerUserID,
		string(settingsJSON),
		tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
		activatedAt,
		suspendedAt,
	)

	if err != nil {
		r.logger.Error("Failed to update tenant", zap.Error(err))
		return fmt.Errorf("erro ao atualizar tenant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrTenantNotFound
	}

	r.logger.Debug("Tenant updated successfully", zap.String("tenant_id", tenant.ID))
	return nil
}

// Delete exclui um tenant
func (r *PostgresTenantRepository) Delete(id string) error {
	r.logger.Debug("Deleting tenant", zap.String("tenant_id", id))

	query := `DELETE FROM tenants WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete tenant", zap.Error(err))
		return fmt.Errorf("erro ao excluir tenant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrTenantNotFound
	}

	r.logger.Debug("Tenant deleted successfully", zap.String("tenant_id", id))
	return nil
}

// GetByStatus busca tenants por status
func (r *PostgresTenantRepository) GetByStatus(status domain.TenantStatus, limit, offset int) ([]*domain.Tenant, error) {
	r.logger.Debug("Getting tenants by status", 
		zap.String("status", string(status)),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	query := `
		SELECT id, name, legal_name, document, email, phone, website,
		       address, status, plan_type, owner_user_id, settings,
		       created_at, updated_at, activated_at, suspended_at
		FROM tenants
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	var tenantsDB []tenantDB
	err := r.db.Select(&tenantsDB, query, string(status), limit, offset)
	if err != nil {
		r.logger.Error("Failed to get tenants by status", zap.Error(err))
		return nil, err
	}

	var tenants []*domain.Tenant
	for _, tenantDB := range tenantsDB {
		tenant, err := r.toDomainTenant(&tenantDB)
		if err != nil {
			r.logger.Error("Failed to convert tenant", zap.Error(err))
			continue
		}
		tenants = append(tenants, tenant)
	}

	r.logger.Debug("Tenants found by status", 
		zap.String("status", string(status)),
		zap.Int("count", len(tenants)),
	)
	return tenants, nil
}

// toDomainTenant converte tenantDB para domain.Tenant
func (r *PostgresTenantRepository) toDomainTenant(tenantDB *tenantDB) (*domain.Tenant, error) {
	var address domain.Address
	if err := json.Unmarshal([]byte(tenantDB.AddressJSON), &address); err != nil {
		r.logger.Error("Failed to unmarshal address", zap.Error(err))
		return nil, err
	}

	var settings domain.TenantSettings
	if err := json.Unmarshal([]byte(tenantDB.SettingsJSON), &settings); err != nil {
		r.logger.Error("Failed to unmarshal settings", zap.Error(err))
		return nil, err
	}

	createdAt, err := parseTime(tenantDB.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := parseTime(tenantDB.UpdatedAt)
	if err != nil {
		return nil, err
	}

	var activatedAt, suspendedAt *time.Time
	if tenantDB.ActivatedAt != nil {
		parsed, err := parseTime(*tenantDB.ActivatedAt)
		if err != nil {
			return nil, err
		}
		activatedAt = &parsed
	}

	if tenantDB.SuspendedAt != nil {
		parsed, err := parseTime(*tenantDB.SuspendedAt)
		if err != nil {
			return nil, err
		}
		suspendedAt = &parsed
	}

	return &domain.Tenant{
		ID:            tenantDB.ID,
		Name:          tenantDB.Name,
		LegalName:     tenantDB.LegalName,
		Document:      tenantDB.Document,
		Email:         tenantDB.Email,
		Phone:         tenantDB.Phone,
		Website:       tenantDB.Website,
		Address:       address,
		Status:        domain.TenantStatus(tenantDB.Status),
		PlanType:      domain.PlanType(tenantDB.PlanType),
		OwnerUserID:   tenantDB.OwnerUserID,
		Settings:      settings,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		ActivatedAt:   activatedAt,
		SuspendedAt:   suspendedAt,
	}, nil
}

// parseTime converte string para time.Time
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}