package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// PostgresTemplateRepository implementação PostgreSQL do NotificationTemplateRepository
type PostgresTemplateRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresTemplateRepository cria nova instância do repositório
func NewPostgresTemplateRepository(db *sqlx.DB, logger *zap.Logger) domain.NotificationTemplateRepository {
	return &PostgresTemplateRepository{
		db:     db,
		logger: logger,
	}
}

// templateDB representa um template no banco de dados
type templateDB struct {
	ID           string         `db:"id"`
	Name         string         `db:"name"`
	Type         string         `db:"type"`
	Channel      string         `db:"channel"`
	Status       string         `db:"status"`
	Subject      string         `db:"subject"`
	Content      string         `db:"content"`
	VariablesJSON string        `db:"variables"`
	TenantID     *string        `db:"tenant_id"`
	IsSystem     bool           `db:"is_system"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

// Create cria um novo template
func (r *PostgresTemplateRepository) Create(ctx context.Context, template *domain.NotificationTemplate) error {
	r.logger.Debug("Creating template", zap.String("template_id", template.ID.String()))

	templateDB := r.toDatabase(template)

	query := `
		INSERT INTO notification_templates (
			id, name, type, channel, status, subject, content,
			variables, tenant_id, is_system, created_at, updated_at
		) VALUES (
			:id, :name, :type, :channel, :status, :subject, :content,
			:variables, :tenant_id, :is_system, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, templateDB)
	if err != nil {
		r.logger.Error("Failed to create template", zap.Error(err))
		return fmt.Errorf("failed to create template: %w", err)
	}

	r.logger.Info("Template created successfully", zap.String("template_id", template.ID.String()))
	return nil
}

// GetByID busca um template por ID
func (r *PostgresTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.NotificationTemplate, error) {
	r.logger.Debug("Getting template by ID", zap.String("template_id", id.String()))

	var templateDB templateDB
	query := `SELECT * FROM notification_templates WHERE id = $1`

	err := r.db.GetContext(ctx, &templateDB, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTemplateNotFound
		}
		r.logger.Error("Failed to get template", zap.Error(err))
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	template, err := r.fromDatabase(&templateDB)
	if err != nil {
		return nil, fmt.Errorf("failed to convert template from database: %w", err)
	}

	return template, nil
}

// Update atualiza um template
func (r *PostgresTemplateRepository) Update(ctx context.Context, template *domain.NotificationTemplate) error {
	r.logger.Debug("Updating template", zap.String("template_id", template.ID.String()))

	templateDB := r.toDatabase(template)
	templateDB.UpdatedAt = time.Now()

	query := `
		UPDATE notification_templates SET
			name = :name, type = :type, channel = :channel, status = :status,
			subject = :subject, content = :content, variables = :variables,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, templateDB)
	if err != nil {
		r.logger.Error("Failed to update template", zap.Error(err))
		return fmt.Errorf("failed to update template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTemplateNotFound
	}

	r.logger.Info("Template updated successfully", zap.String("template_id", template.ID.String()))
	return nil
}

// Delete remove um template
func (r *PostgresTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Deleting template", zap.String("template_id", id.String()))

	query := `DELETE FROM notification_templates WHERE id = $1 AND is_system = FALSE`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		r.logger.Error("Failed to delete template", zap.Error(err))
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTemplateNotFound
	}

	r.logger.Info("Template deleted successfully", zap.String("template_id", id.String()))
	return nil
}

// FindByTenantID busca templates por tenant ID com filtros
func (r *PostgresTemplateRepository) FindByTenantID(ctx context.Context, tenantID *uuid.UUID, filters domain.TemplateFilters) ([]*domain.NotificationTemplate, error) {
	r.logger.Debug("Finding templates by tenant ID")

	var baseQuery string
	var baseArgs []interface{}

	if tenantID != nil {
		baseQuery = "SELECT * FROM notification_templates WHERE (tenant_id = $1 OR is_system = TRUE)"
		baseArgs = []interface{}{tenantID.String()}
	} else {
		baseQuery = "SELECT * FROM notification_templates WHERE is_system = TRUE"
		baseArgs = []interface{}{}
	}

	query, args := r.buildFilterQuery(baseQuery, baseArgs, filters)

	var templatesDB []templateDB
	err := r.db.SelectContext(ctx, &templatesDB, query, args...)
	if err != nil {
		r.logger.Error("Failed to find templates by tenant", zap.Error(err))
		return nil, fmt.Errorf("failed to find templates by tenant: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(templatesDB))
	for i, templateDB := range templatesDB {
		template, err := r.fromDatabase(&templateDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert template from database: %w", err)
		}
		templates[i] = template
	}

	return templates, nil
}

// FindByType busca templates por tipo e canal
func (r *PostgresTemplateRepository) FindByType(ctx context.Context, notificationType domain.NotificationType, channel domain.NotificationChannel, tenantID *uuid.UUID) ([]*domain.NotificationTemplate, error) {
	r.logger.Debug("Finding templates by type and channel",
		zap.String("type", string(notificationType)),
		zap.String("channel", string(channel)))

	var query string
	var args []interface{}

	if tenantID != nil {
		query = `
			SELECT * FROM notification_templates 
			WHERE type = $1 AND channel = $2 AND (tenant_id = $3 OR is_system = TRUE) AND status = 'active'
			ORDER BY is_system ASC, created_at DESC`
		args = []interface{}{string(notificationType), string(channel), tenantID.String()}
	} else {
		query = `
			SELECT * FROM notification_templates 
			WHERE type = $1 AND channel = $2 AND is_system = TRUE AND status = 'active'
			ORDER BY created_at DESC`
		args = []interface{}{string(notificationType), string(channel)}
	}

	var templatesDB []templateDB
	err := r.db.SelectContext(ctx, &templatesDB, query, args...)
	if err != nil {
		r.logger.Error("Failed to find templates by type", zap.Error(err))
		return nil, fmt.Errorf("failed to find templates by type: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(templatesDB))
	for i, templateDB := range templatesDB {
		template, err := r.fromDatabase(&templateDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert template from database: %w", err)
		}
		templates[i] = template
	}

	return templates, nil
}

// FindSystemTemplates busca templates do sistema
func (r *PostgresTemplateRepository) FindSystemTemplates(ctx context.Context, filters domain.TemplateFilters) ([]*domain.NotificationTemplate, error) {
	r.logger.Debug("Finding system templates")

	baseQuery := "SELECT * FROM notification_templates WHERE is_system = TRUE"
	baseArgs := []interface{}{}

	query, args := r.buildFilterQuery(baseQuery, baseArgs, filters)

	var templatesDB []templateDB
	err := r.db.SelectContext(ctx, &templatesDB, query, args...)
	if err != nil {
		r.logger.Error("Failed to find system templates", zap.Error(err))
		return nil, fmt.Errorf("failed to find system templates: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(templatesDB))
	for i, templateDB := range templatesDB {
		template, err := r.fromDatabase(&templateDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert template from database: %w", err)
		}
		templates[i] = template
	}

	return templates, nil
}

// FindActiveTemplates busca templates ativos
func (r *PostgresTemplateRepository) FindActiveTemplates(ctx context.Context, tenantID *uuid.UUID) ([]*domain.NotificationTemplate, error) {
	r.logger.Debug("Finding active templates")

	var query string
	var args []interface{}

	if tenantID != nil {
		query = `
			SELECT * FROM notification_templates 
			WHERE (tenant_id = $1 OR is_system = TRUE) AND status = 'active'
			ORDER BY is_system ASC, name ASC`
		args = []interface{}{tenantID.String()}
	} else {
		query = `
			SELECT * FROM notification_templates 
			WHERE is_system = TRUE AND status = 'active'
			ORDER BY name ASC`
		args = []interface{}{}
	}

	var templatesDB []templateDB
	err := r.db.SelectContext(ctx, &templatesDB, query, args...)
	if err != nil {
		r.logger.Error("Failed to find active templates", zap.Error(err))
		return nil, fmt.Errorf("failed to find active templates: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(templatesDB))
	for i, templateDB := range templatesDB {
		template, err := r.fromDatabase(&templateDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert template from database: %w", err)
		}
		templates[i] = template
	}

	return templates, nil
}

// FindByTypeAndChannel busca template específico por tipo e canal para um tenant
func (r *PostgresTemplateRepository) FindByTypeAndChannel(ctx context.Context, notificationType domain.NotificationType, channel domain.NotificationChannel, tenantID uuid.UUID) (*domain.NotificationTemplate, error) {
	r.logger.Debug("Finding template by type and channel for tenant",
		zap.String("type", string(notificationType)),
		zap.String("channel", string(channel)),
		zap.String("tenant_id", tenantID.String()))

	query := `
		SELECT * FROM notification_templates 
		WHERE type = $1 AND channel = $2 AND (tenant_id = $3 OR is_system = TRUE) AND status = 'active'
		ORDER BY is_system ASC, created_at DESC
		LIMIT 1`

	var templateDB templateDB
	err := r.db.GetContext(ctx, &templateDB, query, string(notificationType), string(channel), tenantID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTemplateNotFound
		}
		r.logger.Error("Failed to find template by type and channel", zap.Error(err))
		return nil, fmt.Errorf("failed to find template by type and channel: %w", err)
	}

	template, err := r.fromDatabase(&templateDB)
	if err != nil {
		return nil, fmt.Errorf("failed to convert template from database: %w", err)
	}

	return template, nil
}

// GetDefaultTemplate busca template padrão do sistema
func (r *PostgresTemplateRepository) GetDefaultTemplate(ctx context.Context, notificationType domain.NotificationType, channel domain.NotificationChannel) (*domain.NotificationTemplate, error) {
	r.logger.Debug("Getting default template",
		zap.String("type", string(notificationType)),
		zap.String("channel", string(channel)))

	query := `
		SELECT * FROM notification_templates 
		WHERE type = $1 AND channel = $2 AND is_system = TRUE AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1`

	var templateDB templateDB
	err := r.db.GetContext(ctx, &templateDB, query, string(notificationType), string(channel))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTemplateNotFound
		}
		r.logger.Error("Failed to get default template", zap.Error(err))
		return nil, fmt.Errorf("failed to get default template: %w", err)
	}

	template, err := r.fromDatabase(&templateDB)
	if err != nil {
		return nil, fmt.Errorf("failed to convert template from database: %w", err)
	}

	return template, nil
}

// CountByTenant conta templates por tenant
func (r *PostgresTemplateRepository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	r.logger.Debug("Counting templates by tenant", zap.String("tenant_id", tenantID.String()))

	query := `SELECT COUNT(*) FROM notification_templates WHERE tenant_id = $1`

	var count int64
	err := r.db.GetContext(ctx, &count, query, tenantID.String())
	if err != nil {
		r.logger.Error("Failed to count templates by tenant", zap.Error(err))
		return 0, fmt.Errorf("failed to count templates by tenant: %w", err)
	}

	return count, nil
}

// buildFilterQuery constrói query com filtros dinâmicos
func (r *PostgresTemplateRepository) buildFilterQuery(baseQuery string, baseArgs []interface{}, filters domain.TemplateFilters) (string, []interface{}) {
	query := baseQuery
	args := baseArgs
	conditions := []string{}
	argIndex := len(baseArgs) + 1

	if filters.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, string(*filters.Type))
		argIndex++
	}

	if filters.Channel != nil {
		conditions = append(conditions, fmt.Sprintf("channel = $%d", argIndex))
		args = append(args, string(*filters.Channel))
		argIndex++
	}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(*filters.Status))
		argIndex++
	}

	if filters.IsSystem != nil {
		conditions = append(conditions, fmt.Sprintf("is_system = $%d", argIndex))
		args = append(args, *filters.IsSystem)
		argIndex++
	}

	if filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR subject ILIKE $%d OR content ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Ordenação
	orderBy := "created_at"
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}

	orderDir := "DESC"
	if filters.OrderDir != "" {
		orderDir = strings.ToUpper(filters.OrderDir)
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)

	// Paginação
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	return query, args
}

// toDatabase converte domain.NotificationTemplate para templateDB
func (r *PostgresTemplateRepository) toDatabase(template *domain.NotificationTemplate) *templateDB {
	variablesJSON, _ := json.Marshal(template.Variables)

	var tenantID *string
	if template.TenantID != nil {
		id := template.TenantID.String()
		tenantID = &id
	}

	return &templateDB{
		ID:           template.ID.String(),
		Name:         template.Name,
		Type:         string(template.Type),
		Channel:      string(template.Channel),
		Status:       string(template.Status),
		Subject:      template.Subject,
		Content:      template.Content,
		VariablesJSON: string(variablesJSON),
		TenantID:     tenantID,
		IsSystem:     template.IsSystem,
		CreatedAt:    template.CreatedAt,
		UpdatedAt:    template.UpdatedAt,
	}
}

// fromDatabase converte templateDB para domain.NotificationTemplate
func (r *PostgresTemplateRepository) fromDatabase(templateDB *templateDB) (*domain.NotificationTemplate, error) {
	id, err := uuid.Parse(templateDB.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID: %w", err)
	}

	var tenantID *uuid.UUID
	if templateDB.TenantID != nil {
		id, err := uuid.Parse(*templateDB.TenantID)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant ID: %w", err)
		}
		tenantID = &id
	}

	var variables []string
	if templateDB.VariablesJSON != "" {
		if err := json.Unmarshal([]byte(templateDB.VariablesJSON), &variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}
	}

	return &domain.NotificationTemplate{
		ID:        id,
		Name:      templateDB.Name,
		Type:      domain.NotificationType(templateDB.Type),
		Channel:   domain.NotificationChannel(templateDB.Channel),
		Status:    domain.TemplateStatus(templateDB.Status),
		Subject:   templateDB.Subject,
		Content:   templateDB.Content,
		Variables: variables,
		TenantID:  tenantID,
		IsSystem:  templateDB.IsSystem,
		CreatedAt: templateDB.CreatedAt,
		UpdatedAt: templateDB.UpdatedAt,
	}, nil
}