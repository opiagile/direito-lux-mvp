package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// PostgresPreferenceRepository implementação PostgreSQL do NotificationPreferenceRepository
type PostgresPreferenceRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresPreferenceRepository cria nova instância do repositório
func NewPostgresPreferenceRepository(db *sqlx.DB, logger *zap.Logger) domain.NotificationPreferenceRepository {
	return &PostgresPreferenceRepository{
		db:     db,
		logger: logger,
	}
}

// preferenceDB representa uma preferência no banco de dados
type preferenceDB struct {
	ID        string         `db:"id"`
	TenantID  string         `db:"tenant_id"`
	UserID    string         `db:"user_id"`
	Type      string         `db:"type"`
	Channels  pq.StringArray `db:"channels"`
	Enabled   bool           `db:"enabled"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

// Create cria uma nova preferência
func (r *PostgresPreferenceRepository) Create(ctx context.Context, preference *domain.NotificationPreference) error {
	r.logger.Debug("Creating preference", zap.String("preference_id", preference.ID.String()))

	prefDB := r.toDatabase(preference)

	query := `
		INSERT INTO notification_preferences (
			id, tenant_id, user_id, type, channels, enabled, created_at, updated_at
		) VALUES (
			:id, :tenant_id, :user_id, :type, :channels, :enabled, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, prefDB)
	if err != nil {
		r.logger.Error("Failed to create preference", zap.Error(err))
		return fmt.Errorf("failed to create preference: %w", err)
	}

	r.logger.Info("Preference created successfully", zap.String("preference_id", preference.ID.String()))
	return nil
}

// GetByID busca uma preferência por ID
func (r *PostgresPreferenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.NotificationPreference, error) {
	r.logger.Debug("Getting preference by ID", zap.String("preference_id", id.String()))

	var prefDB preferenceDB
	query := `SELECT * FROM notification_preferences WHERE id = $1`

	err := r.db.GetContext(ctx, &prefDB, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPreferenceNotFound
		}
		r.logger.Error("Failed to get preference", zap.Error(err))
		return nil, fmt.Errorf("failed to get preference: %w", err)
	}

	preference, err := r.fromDatabase(&prefDB)
	if err != nil {
		return nil, fmt.Errorf("failed to convert preference from database: %w", err)
	}

	return preference, nil
}

// Update atualiza uma preferência
func (r *PostgresPreferenceRepository) Update(ctx context.Context, preference *domain.NotificationPreference) error {
	r.logger.Debug("Updating preference", zap.String("preference_id", preference.ID.String()))

	prefDB := r.toDatabase(preference)
	prefDB.UpdatedAt = time.Now()

	query := `
		UPDATE notification_preferences SET
			channels = :channels, enabled = :enabled, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, prefDB)
	if err != nil {
		r.logger.Error("Failed to update preference", zap.Error(err))
		return fmt.Errorf("failed to update preference: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPreferenceNotFound
	}

	r.logger.Info("Preference updated successfully", zap.String("preference_id", preference.ID.String()))
	return nil
}

// Delete remove uma preferência
func (r *PostgresPreferenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Deleting preference", zap.String("preference_id", id.String()))

	query := `DELETE FROM notification_preferences WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		r.logger.Error("Failed to delete preference", zap.Error(err))
		return fmt.Errorf("failed to delete preference: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPreferenceNotFound
	}

	r.logger.Info("Preference deleted successfully", zap.String("preference_id", id.String()))
	return nil
}

// FindByUserID busca preferências por usuário
func (r *PostgresPreferenceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.NotificationPreference, error) {
	r.logger.Debug("Finding preferences by user ID", zap.String("user_id", userID.String()))

	query := `
		SELECT * FROM notification_preferences 
		WHERE user_id = $1 
		ORDER BY type ASC`

	var preferencesDB []preferenceDB
	err := r.db.SelectContext(ctx, &preferencesDB, query, userID.String())
	if err != nil {
		r.logger.Error("Failed to find preferences by user", zap.Error(err))
		return nil, fmt.Errorf("failed to find preferences by user: %w", err)
	}

	preferences := make([]*domain.NotificationPreference, len(preferencesDB))
	for i, prefDB := range preferencesDB {
		preference, err := r.fromDatabase(&prefDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert preference from database: %w", err)
		}
		preferences[i] = preference
	}

	return preferences, nil
}

// FindByUserAndType busca preferência específica por usuário e tipo
func (r *PostgresPreferenceRepository) FindByUserAndType(ctx context.Context, userID uuid.UUID, notificationType domain.NotificationType) (*domain.NotificationPreference, error) {
	r.logger.Debug("Finding preference by user and type",
		zap.String("user_id", userID.String()),
		zap.String("type", string(notificationType)))

	query := `
		SELECT * FROM notification_preferences 
		WHERE user_id = $1 AND type = $2`

	var prefDB preferenceDB
	err := r.db.GetContext(ctx, &prefDB, query, userID.String(), string(notificationType))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPreferenceNotFound
		}
		r.logger.Error("Failed to find preference by user and type", zap.Error(err))
		return nil, fmt.Errorf("failed to find preference by user and type: %w", err)
	}

	preference, err := r.fromDatabase(&prefDB)
	if err != nil {
		return nil, fmt.Errorf("failed to convert preference from database: %w", err)
	}

	return preference, nil
}

// UpsertUserPreferences cria ou atualiza preferências do usuário em lote
func (r *PostgresPreferenceRepository) UpsertUserPreferences(ctx context.Context, userID uuid.UUID, preferences []*domain.NotificationPreference) error {
	if len(preferences) == 0 {
		return nil
	}

	r.logger.Debug("Upserting user preferences", 
		zap.String("user_id", userID.String()),
		zap.Int("count", len(preferences)))

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Query para upsert (INSERT ... ON CONFLICT ... DO UPDATE)
	query := `
		INSERT INTO notification_preferences (
			id, tenant_id, user_id, type, channels, enabled, created_at, updated_at
		) VALUES (
			:id, :tenant_id, :user_id, :type, :channels, :enabled, :created_at, :updated_at
		) ON CONFLICT (user_id, type) DO UPDATE SET
			channels = EXCLUDED.channels,
			enabled = EXCLUDED.enabled,
			updated_at = EXCLUDED.updated_at`

	for _, preference := range preferences {
		prefDB := r.toDatabase(preference)
		_, err := tx.NamedExecContext(ctx, query, prefDB)
		if err != nil {
			r.logger.Error("Failed to upsert preference", zap.Error(err))
			return fmt.Errorf("failed to upsert preference: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("User preferences upserted successfully", 
		zap.String("user_id", userID.String()),
		zap.Int("count", len(preferences)))
	return nil
}

// GetUserPreferences obtém preferências do usuário em formato de mapa
func (r *PostgresPreferenceRepository) GetUserPreferences(ctx context.Context, userID uuid.UUID) (map[domain.NotificationType][]domain.NotificationChannel, error) {
	r.logger.Debug("Getting user preferences", zap.String("user_id", userID.String()))

	query := `
		SELECT type, channels FROM notification_preferences 
		WHERE user_id = $1 AND enabled = TRUE`

	type preferenceResult struct {
		Type     string         `db:"type"`
		Channels pq.StringArray `db:"channels"`
	}

	var results []preferenceResult
	err := r.db.SelectContext(ctx, &results, query, userID.String())
	if err != nil {
		r.logger.Error("Failed to get user preferences", zap.Error(err))
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	preferences := make(map[domain.NotificationType][]domain.NotificationChannel)
	for _, result := range results {
		notificationType := domain.NotificationType(result.Type)
		channels := make([]domain.NotificationChannel, len(result.Channels))
		for i, channel := range result.Channels {
			channels[i] = domain.NotificationChannel(channel)
		}
		preferences[notificationType] = channels
	}

	return preferences, nil
}

// CreateDefaultPreferences cria preferências padrão para um usuário
func (r *PostgresPreferenceRepository) CreateDefaultPreferences(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error {
	r.logger.Debug("Creating default preferences",
		zap.String("tenant_id", tenantID.String()),
		zap.String("user_id", userID.String()))

	// Definir preferências padrão
	defaultPreferences := []*domain.NotificationPreference{
		{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   userID,
			Type:     domain.NotificationTypeProcessUpdate,
			Channels: []domain.NotificationChannel{
				domain.NotificationChannelEmail,
				domain.NotificationChannelWhatsApp,
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   userID,
			Type:     domain.NotificationTypeMovementAlert,
			Channels: []domain.NotificationChannel{
				domain.NotificationChannelEmail,
				domain.NotificationChannelWhatsApp,
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   userID,
			Type:     domain.NotificationTypeDeadlineReminder,
			Channels: []domain.NotificationChannel{
				domain.NotificationChannelEmail,
				domain.NotificationChannelWhatsApp,
				domain.NotificationChannelTelegram,
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   userID,
			Type:     domain.NotificationTypeSystemAlert,
			Channels: []domain.NotificationChannel{
				domain.NotificationChannelEmail,
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return r.UpsertUserPreferences(ctx, userID, defaultPreferences)
}

// toDatabase converte domain.NotificationPreference para preferenceDB
func (r *PostgresPreferenceRepository) toDatabase(preference *domain.NotificationPreference) *preferenceDB {
	channels := make([]string, len(preference.Channels))
	for i, channel := range preference.Channels {
		channels[i] = string(channel)
	}

	return &preferenceDB{
		ID:        preference.ID.String(),
		TenantID:  preference.TenantID.String(),
		UserID:    preference.UserID.String(),
		Type:      string(preference.Type),
		Channels:  pq.StringArray(channels),
		Enabled:   preference.Enabled,
		CreatedAt: preference.CreatedAt,
		UpdatedAt: preference.UpdatedAt,
	}
}

// fromDatabase converte preferenceDB para domain.NotificationPreference
func (r *PostgresPreferenceRepository) fromDatabase(prefDB *preferenceDB) (*domain.NotificationPreference, error) {
	id, err := uuid.Parse(prefDB.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid preference ID: %w", err)
	}

	tenantID, err := uuid.Parse(prefDB.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	userID, err := uuid.Parse(prefDB.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	channels := make([]domain.NotificationChannel, len(prefDB.Channels))
	for i, channel := range prefDB.Channels {
		channels[i] = domain.NotificationChannel(channel)
	}

	return &domain.NotificationPreference{
		ID:        id,
		TenantID:  tenantID,
		UserID:    userID,
		Type:      domain.NotificationType(prefDB.Type),
		Channels:  channels,
		Enabled:   prefDB.Enabled,
		CreatedAt: prefDB.CreatedAt,
		UpdatedAt: prefDB.UpdatedAt,
	}, nil
}