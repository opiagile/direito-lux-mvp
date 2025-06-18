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
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// PostgresNotificationRepository implementação PostgreSQL do NotificationRepository
type PostgresNotificationRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresNotificationRepository cria nova instância do repositório
func NewPostgresNotificationRepository(db *sqlx.DB, logger *zap.Logger) domain.NotificationRepository {
	return &PostgresNotificationRepository{
		db:     db,
		logger: logger,
	}
}

// notificationDB representa uma notificação no banco de dados
type notificationDB struct {
	ID                   string         `db:"id"`
	TenantID             string         `db:"tenant_id"`
	Type                 string         `db:"type"`
	Channel              string         `db:"channel"`
	Priority             string         `db:"priority"`
	Status               string         `db:"status"`
	Subject              string         `db:"subject"`
	Content              string         `db:"content"`
	RecipientID          string         `db:"recipient_id"`
	RecipientType        string         `db:"recipient_type"`
	RecipientContact     string         `db:"recipient_contact"`
	TemplateID           *string        `db:"template_id"`
	VariablesJSON        string         `db:"variables"`
	MetadataJSON         string         `db:"metadata"`
	RetryCount           int            `db:"retry_count"`
	MaxRetries           int            `db:"max_retries"`
	NextRetryAt          *time.Time     `db:"next_retry_at"`
	ScheduledAt          *time.Time     `db:"scheduled_at"`
	SentAt               *time.Time     `db:"sent_at"`
	DeliveredAt          *time.Time     `db:"delivered_at"`
	FailedAt             *time.Time     `db:"failed_at"`
	ErrorMessage         *string        `db:"error_message"`
	ExternalID           *string        `db:"external_id"`
	ExternalStatus       *string        `db:"external_status"`
	ProcessingStartedAt  *time.Time     `db:"processing_started_at"`
	ProcessingFinishedAt *time.Time     `db:"processing_finished_at"`
	Tags                 pq.StringArray `db:"tags"`
	CreatedAt            time.Time      `db:"created_at"`
	UpdatedAt            time.Time      `db:"updated_at"`
}

// Create cria uma nova notificação
func (r *PostgresNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	r.logger.Debug("Creating notification", zap.String("notification_id", notification.ID.String()))

	notifDB := r.toDatabase(notification)

	query := `
		INSERT INTO notifications (
			id, tenant_id, type, channel, priority, status, subject, content,
			recipient_id, recipient_type, recipient_contact, template_id,
			variables, metadata, retry_count, max_retries, next_retry_at,
			scheduled_at, tags, created_at, updated_at
		) VALUES (
			:id, :tenant_id, :type, :channel, :priority, :status, :subject, :content,
			:recipient_id, :recipient_type, :recipient_contact, :template_id,
			:variables, :metadata, :retry_count, :max_retries, :next_retry_at,
			:scheduled_at, :tags, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, notifDB)
	if err != nil {
		r.logger.Error("Failed to create notification", zap.Error(err))
		return fmt.Errorf("failed to create notification: %w", err)
	}

	r.logger.Info("Notification created successfully", zap.String("notification_id", notification.ID.String()))
	return nil
}

// GetByID busca uma notificação por ID
func (r *PostgresNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	r.logger.Debug("Getting notification by ID", zap.String("notification_id", id.String()))

	var notifDB notificationDB
	query := `SELECT * FROM notifications WHERE id = $1`

	err := r.db.GetContext(ctx, &notifDB, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotificationNotFound
		}
		r.logger.Error("Failed to get notification", zap.Error(err))
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	notification, err := r.fromDatabase(&notifDB)
	if err != nil {
		return nil, fmt.Errorf("failed to convert notification from database: %w", err)
	}

	return notification, nil
}

// Update atualiza uma notificação
func (r *PostgresNotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	r.logger.Debug("Updating notification", zap.String("notification_id", notification.ID.String()))

	notifDB := r.toDatabase(notification)
	notifDB.UpdatedAt = time.Now()

	query := `
		UPDATE notifications SET
			status = :status, retry_count = :retry_count, next_retry_at = :next_retry_at,
			sent_at = :sent_at, delivered_at = :delivered_at, failed_at = :failed_at,
			error_message = :error_message, external_id = :external_id,
			external_status = :external_status, processing_started_at = :processing_started_at,
			processing_finished_at = :processing_finished_at, metadata = :metadata,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, notifDB)
	if err != nil {
		r.logger.Error("Failed to update notification", zap.Error(err))
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotificationNotFound
	}

	r.logger.Info("Notification updated successfully", zap.String("notification_id", notification.ID.String()))
	return nil
}

// Delete remove uma notificação
func (r *PostgresNotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Deleting notification", zap.String("notification_id", id.String()))

	query := `DELETE FROM notifications WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		r.logger.Error("Failed to delete notification", zap.Error(err))
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotificationNotFound
	}

	r.logger.Info("Notification deleted successfully", zap.String("notification_id", id.String()))
	return nil
}

// FindByTenantID busca notificações por tenant ID com filtros
func (r *PostgresNotificationRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID, filters domain.NotificationFilters) ([]*domain.Notification, error) {
	r.logger.Debug("Finding notifications by tenant ID", zap.String("tenant_id", tenantID.String()))

	query, args := r.buildFilterQuery("SELECT * FROM notifications WHERE tenant_id = $1", []interface{}{tenantID.String()}, filters)

	var notificationsDB []notificationDB
	err := r.db.SelectContext(ctx, &notificationsDB, query, args...)
	if err != nil {
		r.logger.Error("Failed to find notifications by tenant", zap.Error(err))
		return nil, fmt.Errorf("failed to find notifications by tenant: %w", err)
	}

	notifications := make([]*domain.Notification, len(notificationsDB))
	for i, notifDB := range notificationsDB {
		notification, err := r.fromDatabase(&notifDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification from database: %w", err)
		}
		notifications[i] = notification
	}

	return notifications, nil
}

// FindByRecipient busca notificações por recipient ID
func (r *PostgresNotificationRepository) FindByRecipient(ctx context.Context, recipientID string, filters domain.NotificationFilters) ([]*domain.Notification, error) {
	r.logger.Debug("Finding notifications by recipient", zap.String("recipient_id", recipientID))

	query, args := r.buildFilterQuery("SELECT * FROM notifications WHERE recipient_id = $1", []interface{}{recipientID}, filters)

	var notificationsDB []notificationDB
	err := r.db.SelectContext(ctx, &notificationsDB, query, args...)
	if err != nil {
		r.logger.Error("Failed to find notifications by recipient", zap.Error(err))
		return nil, fmt.Errorf("failed to find notifications by recipient: %w", err)
	}

	notifications := make([]*domain.Notification, len(notificationsDB))
	for i, notifDB := range notificationsDB {
		notification, err := r.fromDatabase(&notifDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification from database: %w", err)
		}
		notifications[i] = notification
	}

	return notifications, nil
}

// FindByStatus busca notificações por status
func (r *PostgresNotificationRepository) FindByStatus(ctx context.Context, status domain.NotificationStatus, limit int) ([]*domain.Notification, error) {
	r.logger.Debug("Finding notifications by status", zap.String("status", string(status)))

	query := `SELECT * FROM notifications WHERE status = $1 ORDER BY created_at ASC LIMIT $2`

	var notificationsDB []notificationDB
	err := r.db.SelectContext(ctx, &notificationsDB, query, string(status), limit)
	if err != nil {
		r.logger.Error("Failed to find notifications by status", zap.Error(err))
		return nil, fmt.Errorf("failed to find notifications by status: %w", err)
	}

	notifications := make([]*domain.Notification, len(notificationsDB))
	for i, notifDB := range notificationsDB {
		notification, err := r.fromDatabase(&notifDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification from database: %w", err)
		}
		notifications[i] = notification
	}

	return notifications, nil
}

// FindPending busca notificações pendentes
func (r *PostgresNotificationRepository) FindPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	return r.FindByStatus(ctx, domain.NotificationStatusPending, limit)
}

// FindScheduled busca notificações agendadas para antes de um tempo específico
func (r *PostgresNotificationRepository) FindScheduled(ctx context.Context, scheduledBefore time.Time, limit int) ([]*domain.Notification, error) {
	r.logger.Debug("Finding scheduled notifications", zap.Time("scheduled_before", scheduledBefore))

	query := `
		SELECT * FROM notifications 
		WHERE status = $1 AND scheduled_at IS NOT NULL AND scheduled_at <= $2 
		ORDER BY scheduled_at ASC LIMIT $3`

	var notificationsDB []notificationDB
	err := r.db.SelectContext(ctx, &notificationsDB, query, string(domain.NotificationStatusScheduled), scheduledBefore, limit)
	if err != nil {
		r.logger.Error("Failed to find scheduled notifications", zap.Error(err))
		return nil, fmt.Errorf("failed to find scheduled notifications: %w", err)
	}

	notifications := make([]*domain.Notification, len(notificationsDB))
	for i, notifDB := range notificationsDB {
		notification, err := r.fromDatabase(&notifDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification from database: %w", err)
		}
		notifications[i] = notification
	}

	return notifications, nil
}

// FindExpired busca notificações expiradas
func (r *PostgresNotificationRepository) FindExpired(ctx context.Context, limit int) ([]*domain.Notification, error) {
	r.logger.Debug("Finding expired notifications")

	query := `
		SELECT * FROM notifications 
		WHERE status IN ($1, $2) AND created_at < NOW() - INTERVAL '7 days'
		ORDER BY created_at ASC LIMIT $3`

	var notificationsDB []notificationDB
	err := r.db.SelectContext(ctx, &notificationsDB, query, 
		string(domain.NotificationStatusPending), 
		string(domain.NotificationStatusScheduled), 
		limit)
	if err != nil {
		r.logger.Error("Failed to find expired notifications", zap.Error(err))
		return nil, fmt.Errorf("failed to find expired notifications: %w", err)
	}

	notifications := make([]*domain.Notification, len(notificationsDB))
	for i, notifDB := range notificationsDB {
		notification, err := r.fromDatabase(&notifDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification from database: %w", err)
		}
		notifications[i] = notification
	}

	return notifications, nil
}

// FindForRetry busca notificações prontas para reenvio
func (r *PostgresNotificationRepository) FindForRetry(ctx context.Context, limit int) ([]*domain.Notification, error) {
	r.logger.Debug("Finding notifications for retry")

	query := `
		SELECT * FROM notifications 
		WHERE status = $1 AND retry_count < max_retries 
		AND (next_retry_at IS NULL OR next_retry_at <= NOW())
		ORDER BY priority DESC, next_retry_at ASC LIMIT $2`

	var notificationsDB []notificationDB
	err := r.db.SelectContext(ctx, &notificationsDB, query, string(domain.NotificationStatusFailed), limit)
	if err != nil {
		r.logger.Error("Failed to find notifications for retry", zap.Error(err))
		return nil, fmt.Errorf("failed to find notifications for retry: %w", err)
	}

	notifications := make([]*domain.Notification, len(notificationsDB))
	for i, notifDB := range notificationsDB {
		notification, err := r.fromDatabase(&notifDB)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification from database: %w", err)
		}
		notifications[i] = notification
	}

	return notifications, nil
}

// CountByTenant conta notificações por tenant
func (r *PostgresNotificationRepository) CountByTenant(ctx context.Context, tenantID uuid.UUID, filters domain.NotificationFilters) (int64, error) {
	r.logger.Debug("Counting notifications by tenant", zap.String("tenant_id", tenantID.String()))

	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM notifications WHERE tenant_id = $1", []interface{}{tenantID.String()}, filters)

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		r.logger.Error("Failed to count notifications by tenant", zap.Error(err))
		return 0, fmt.Errorf("failed to count notifications by tenant: %w", err)
	}

	return count, nil
}

// CountByStatus conta notificações por status
func (r *PostgresNotificationRepository) CountByStatus(ctx context.Context, tenantID uuid.UUID, status domain.NotificationStatus) (int64, error) {
	r.logger.Debug("Counting notifications by status", 
		zap.String("tenant_id", tenantID.String()),
		zap.String("status", string(status)))

	query := `SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND status = $2`

	var count int64
	err := r.db.GetContext(ctx, &count, query, tenantID.String(), string(status))
	if err != nil {
		r.logger.Error("Failed to count notifications by status", zap.Error(err))
		return 0, fmt.Errorf("failed to count notifications by status: %w", err)
	}

	return count, nil
}

// CountByChannel conta notificações por canal
func (r *PostgresNotificationRepository) CountByChannel(ctx context.Context, tenantID uuid.UUID, channel domain.NotificationChannel) (int64, error) {
	r.logger.Debug("Counting notifications by channel", 
		zap.String("tenant_id", tenantID.String()),
		zap.String("channel", string(channel)))

	query := `SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND channel = $2`

	var count int64
	err := r.db.GetContext(ctx, &count, query, tenantID.String(), string(channel))
	if err != nil {
		r.logger.Error("Failed to count notifications by channel", zap.Error(err))
		return 0, fmt.Errorf("failed to count notifications by channel: %w", err)
	}

	return count, nil
}

// GetStatistics obtém estatísticas de notificações
func (r *PostgresNotificationRepository) GetStatistics(ctx context.Context, tenantID uuid.UUID, period time.Duration) (*domain.NotificationStatistics, error) {
	r.logger.Debug("Getting notification statistics", zap.String("tenant_id", tenantID.String()))

	startTime := time.Now().Add(-period)

	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'sent' THEN 1 END) as sent,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN channel = 'whatsapp' AND status = 'sent' THEN 1 END) as whatsapp_sent,
			COUNT(CASE WHEN channel = 'email' AND status = 'sent' THEN 1 END) as email_sent,
			COUNT(CASE WHEN channel = 'telegram' AND status = 'sent' THEN 1 END) as telegram_sent,
			COUNT(CASE WHEN channel = 'push' AND status = 'sent' THEN 1 END) as push_sent,
			COUNT(CASE WHEN channel = 'sms' AND status = 'sent' THEN 1 END) as sms_sent,
			COUNT(CASE WHEN type = 'process_update' THEN 1 END) as process_updates,
			COUNT(CASE WHEN type = 'movement_alert' THEN 1 END) as movement_alerts,
			COUNT(CASE WHEN type = 'deadline_reminder' THEN 1 END) as deadline_reminders,
			COUNT(CASE WHEN type = 'system_alert' THEN 1 END) as system_alerts,
			AVG(retry_count) as avg_retries,
			AVG(EXTRACT(EPOCH FROM (processing_finished_at - processing_started_at))) as avg_processing_time
		FROM notifications 
		WHERE tenant_id = $1 AND created_at >= $2`

	type statsRow struct {
		Total              int64   `db:"total"`
		Sent               int64   `db:"sent"`
		Failed             int64   `db:"failed"`
		Pending            int64   `db:"pending"`
		WhatsAppSent       int64   `db:"whatsapp_sent"`
		EmailSent          int64   `db:"email_sent"`
		TelegramSent       int64   `db:"telegram_sent"`
		PushSent           int64   `db:"push_sent"`
		SMSSent            int64   `db:"sms_sent"`
		ProcessUpdates     int64   `db:"process_updates"`
		MovementAlerts     int64   `db:"movement_alerts"`
		DeadlineReminders  int64   `db:"deadline_reminders"`
		SystemAlerts       int64   `db:"system_alerts"`
		AvgRetries         *float64 `db:"avg_retries"`
		AvgProcessingTime  *float64 `db:"avg_processing_time"`
	}

	var row statsRow
	err := r.db.GetContext(ctx, &row, query, tenantID.String(), startTime)
	if err != nil {
		r.logger.Error("Failed to get notification statistics", zap.Error(err))
		return nil, fmt.Errorf("failed to get notification statistics: %w", err)
	}

	stats := &domain.NotificationStatistics{
		TenantID:           tenantID,
		Period:             period.String(),
		TotalSent:          row.Sent,
		TotalFailed:        row.Failed,
		TotalPending:       row.Pending,
		WhatsAppSent:       row.WhatsAppSent,
		EmailSent:          row.EmailSent,
		TelegramSent:       row.TelegramSent,
		PushSent:           row.PushSent,
		SMSSent:            row.SMSSent,
		ProcessUpdates:     row.ProcessUpdates,
		MovementAlerts:     row.MovementAlerts,
		DeadlineReminders:  row.DeadlineReminders,
		SystemAlerts:       row.SystemAlerts,
		GeneratedAt:        time.Now(),
	}

	// Calcular taxas
	if row.Total > 0 {
		stats.DeliveryRate = float64(row.Sent) / float64(row.Total) * 100
		stats.ErrorRate = float64(row.Failed) / float64(row.Total) * 100
	}

	if row.AvgRetries != nil {
		stats.AverageRetries = *row.AvgRetries
		stats.RetryRate = *row.AvgRetries / float64(row.Total) * 100
	}

	if row.AvgProcessingTime != nil {
		stats.AverageProcessingTime = time.Duration(*row.AvgProcessingTime * float64(time.Second))
	}

	return stats, nil
}

// CreateBatch cria múltiplas notificações em lote
func (r *PostgresNotificationRepository) CreateBatch(ctx context.Context, notifications []*domain.Notification) error {
	if len(notifications) == 0 {
		return nil
	}

	r.logger.Debug("Creating notification batch", zap.Int("count", len(notifications)))

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO notifications (
			id, tenant_id, type, channel, priority, status, subject, content,
			recipient_id, recipient_type, recipient_contact, template_id,
			variables, metadata, retry_count, max_retries, next_retry_at,
			scheduled_at, tags, created_at, updated_at
		) VALUES (
			:id, :tenant_id, :type, :channel, :priority, :status, :subject, :content,
			:recipient_id, :recipient_type, :recipient_contact, :template_id,
			:variables, :metadata, :retry_count, :max_retries, :next_retry_at,
			:scheduled_at, :tags, :created_at, :updated_at
		)`

	for _, notification := range notifications {
		notifDB := r.toDatabase(notification)
		_, err := tx.NamedExecContext(ctx, query, notifDB)
		if err != nil {
			r.logger.Error("Failed to create notification in batch", zap.Error(err))
			return fmt.Errorf("failed to create notification in batch: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Notification batch created successfully", zap.Int("count", len(notifications)))
	return nil
}

// UpdateBatch atualiza múltiplas notificações em lote
func (r *PostgresNotificationRepository) UpdateBatch(ctx context.Context, notifications []*domain.Notification) error {
	if len(notifications) == 0 {
		return nil
	}

	r.logger.Debug("Updating notification batch", zap.Int("count", len(notifications)))

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE notifications SET
			status = :status, retry_count = :retry_count, next_retry_at = :next_retry_at,
			sent_at = :sent_at, delivered_at = :delivered_at, failed_at = :failed_at,
			error_message = :error_message, external_id = :external_id,
			external_status = :external_status, processing_started_at = :processing_started_at,
			processing_finished_at = :processing_finished_at, metadata = :metadata,
			updated_at = :updated_at
		WHERE id = :id`

	for _, notification := range notifications {
		notifDB := r.toDatabase(notification)
		notifDB.UpdatedAt = time.Now()
		
		result, err := tx.NamedExecContext(ctx, query, notifDB)
		if err != nil {
			r.logger.Error("Failed to update notification in batch", zap.Error(err))
			return fmt.Errorf("failed to update notification in batch: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("notification not found: %s", notification.ID.String())
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Notification batch updated successfully", zap.Int("count", len(notifications)))
	return nil
}

// DeleteExpired remove notificações expiradas
func (r *PostgresNotificationRepository) DeleteExpired(ctx context.Context, expiredBefore time.Time) (int64, error) {
	r.logger.Debug("Deleting expired notifications", zap.Time("expired_before", expiredBefore))

	query := `
		DELETE FROM notifications 
		WHERE status IN ($1, $2, $3) AND created_at < $4`

	result, err := r.db.ExecContext(ctx, query, 
		string(domain.NotificationStatusSent),
		string(domain.NotificationStatusFailed),
		string(domain.NotificationStatusExpired),
		expiredBefore)
	if err != nil {
		r.logger.Error("Failed to delete expired notifications", zap.Error(err))
		return 0, fmt.Errorf("failed to delete expired notifications: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.logger.Info("Expired notifications deleted", zap.Int64("count", rowsAffected))
	return rowsAffected, nil
}

// buildFilterQuery constrói query com filtros dinâmicos
func (r *PostgresNotificationRepository) buildFilterQuery(baseQuery string, baseArgs []interface{}, filters domain.NotificationFilters) (string, []interface{}) {
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

	if filters.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, string(*filters.Priority))
		argIndex++
	}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(*filters.Status))
		argIndex++
	}

	if filters.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("recipient_id = $%d", argIndex))
		args = append(args, filters.UserID.String())
		argIndex++
	}

	if filters.CreatedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filters.CreatedAfter)
		argIndex++
	}

	if filters.CreatedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filters.CreatedBefore)
		argIndex++
	}

	if filters.SentAfter != nil {
		conditions = append(conditions, fmt.Sprintf("sent_at >= $%d", argIndex))
		args = append(args, *filters.SentAfter)
		argIndex++
	}

	if filters.SentBefore != nil {
		conditions = append(conditions, fmt.Sprintf("sent_at <= $%d", argIndex))
		args = append(args, *filters.SentBefore)
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

// toDatabase converte domain.Notification para notificationDB
func (r *PostgresNotificationRepository) toDatabase(notification *domain.Notification) *notificationDB {
	variablesJSON, _ := json.Marshal(notification.Variables)
	metadataJSON, _ := json.Marshal(notification.Metadata)

	var templateID *string
	if notification.TemplateID != nil {
		id := notification.TemplateID.String()
		templateID = &id
	}

	return &notificationDB{
		ID:                   notification.ID.String(),
		TenantID:             notification.TenantID.String(),
		Type:                 string(notification.Type),
		Channel:              string(notification.Channel),
		Priority:             string(notification.Priority),
		Status:               string(notification.Status),
		Subject:              notification.Subject,
		Content:              notification.Content,
		RecipientID:          notification.RecipientID,
		RecipientType:        notification.RecipientType,
		RecipientContact:     notification.RecipientContact,
		TemplateID:           templateID,
		VariablesJSON:        string(variablesJSON),
		MetadataJSON:         string(metadataJSON),
		RetryCount:           notification.RetryCount,
		MaxRetries:           notification.MaxRetries,
		NextRetryAt:          notification.NextRetryAt,
		ScheduledAt:          notification.ScheduledAt,
		SentAt:               notification.SentAt,
		DeliveredAt:          notification.DeliveredAt,
		FailedAt:             notification.FailedAt,
		ErrorMessage:         notification.ErrorMessage,
		ExternalID:           notification.ExternalID,
		ExternalStatus:       notification.ExternalStatus,
		ProcessingStartedAt:  notification.ProcessingStartedAt,
		ProcessingFinishedAt: notification.ProcessingFinishedAt,
		Tags:                 pq.StringArray(notification.Tags),
		CreatedAt:            notification.CreatedAt,
		UpdatedAt:            notification.UpdatedAt,
	}
}

// fromDatabase converte notificationDB para domain.Notification
func (r *PostgresNotificationRepository) fromDatabase(notifDB *notificationDB) (*domain.Notification, error) {
	id, err := uuid.Parse(notifDB.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID: %w", err)
	}

	tenantID, err := uuid.Parse(notifDB.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	var templateID *uuid.UUID
	if notifDB.TemplateID != nil {
		id, err := uuid.Parse(*notifDB.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("invalid template ID: %w", err)
		}
		templateID = &id
	}

	var variables map[string]interface{}
	if notifDB.VariablesJSON != "" {
		if err := json.Unmarshal([]byte(notifDB.VariablesJSON), &variables); err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}
	}

	var metadata map[string]interface{}
	if notifDB.MetadataJSON != "" {
		if err := json.Unmarshal([]byte(notifDB.MetadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &domain.Notification{
		ID:                   id,
		TenantID:             tenantID,
		Type:                 domain.NotificationType(notifDB.Type),
		Channel:              domain.NotificationChannel(notifDB.Channel),
		Priority:             domain.NotificationPriority(notifDB.Priority),
		Status:               domain.NotificationStatus(notifDB.Status),
		Subject:              notifDB.Subject,
		Content:              notifDB.Content,
		RecipientID:          notifDB.RecipientID,
		RecipientType:        notifDB.RecipientType,
		RecipientContact:     notifDB.RecipientContact,
		TemplateID:           templateID,
		Variables:            variables,
		Metadata:             metadata,
		RetryCount:           notifDB.RetryCount,
		MaxRetries:           notifDB.MaxRetries,
		NextRetryAt:          notifDB.NextRetryAt,
		ScheduledAt:          notifDB.ScheduledAt,
		SentAt:               notifDB.SentAt,
		DeliveredAt:          notifDB.DeliveredAt,
		FailedAt:             notifDB.FailedAt,
		ErrorMessage:         notifDB.ErrorMessage,
		ExternalID:           notifDB.ExternalID,
		ExternalStatus:       notifDB.ExternalStatus,
		ProcessingStartedAt:  notifDB.ProcessingStartedAt,
		ProcessingFinishedAt: notifDB.ProcessingFinishedAt,
		Tags:                 []string(notifDB.Tags),
		CreatedAt:            notifDB.CreatedAt,
		UpdatedAt:            notifDB.UpdatedAt,
	}, nil
}