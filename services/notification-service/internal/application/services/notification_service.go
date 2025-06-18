package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// NotificationService serviço de aplicação para notificações
type NotificationService struct {
	notificationRepo domain.NotificationRepository
	templateRepo     domain.NotificationTemplateRepository
	preferenceRepo   domain.NotificationPreferenceRepository
	queue           domain.NotificationQueue
	providers       map[domain.NotificationChannel]domain.NotificationProvider
	logger          *zap.Logger
}

// NewNotificationService cria nova instância do serviço
func NewNotificationService(
	notificationRepo domain.NotificationRepository,
	templateRepo domain.NotificationTemplateRepository,
	preferenceRepo domain.NotificationPreferenceRepository,
	queue domain.NotificationQueue,
	providers map[domain.NotificationChannel]domain.NotificationProvider,
	logger *zap.Logger,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		templateRepo:     templateRepo,
		preferenceRepo:   preferenceRepo,
		queue:           queue,
		providers:       providers,
		logger:          logger,
	}
}

// CreateNotificationRequest request para criação de notificação
type CreateNotificationRequest struct {
	TenantID         uuid.UUID                          `json:"tenant_id" validate:"required"`
	Type             domain.NotificationType            `json:"type" validate:"required"`
	Channel          *domain.NotificationChannel        `json:"channel,omitempty"`
	Priority         domain.NotificationPriority        `json:"priority" validate:"required"`
	Subject          string                             `json:"subject,omitempty"`
	Content          string                             `json:"content,omitempty"`
	RecipientID      string                             `json:"recipient_id" validate:"required"`
	RecipientType    string                             `json:"recipient_type" validate:"required"`
	RecipientContact string                             `json:"recipient_contact" validate:"required"`
	TemplateID       *uuid.UUID                         `json:"template_id,omitempty"`
	Variables        map[string]interface{}             `json:"variables,omitempty"`
	Metadata         map[string]interface{}             `json:"metadata,omitempty"`
	ScheduledAt      *time.Time                         `json:"scheduled_at,omitempty"`
	Tags             []string                           `json:"tags,omitempty"`
}

// CreateNotification cria uma nova notificação
func (s *NotificationService) CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*domain.Notification, error) {
	s.logger.Debug("Creating notification", 
		zap.String("tenant_id", req.TenantID.String()),
		zap.String("type", string(req.Type)))

	// Validar canais preferidos do usuário se não especificado
	channels := []domain.NotificationChannel{}
	if req.Channel != nil {
		channels = []domain.NotificationChannel{*req.Channel}
	} else {
		userID, err := uuid.Parse(req.RecipientID)
		if err == nil {
			prefs, err := s.preferenceRepo.GetUserPreferences(ctx, userID)
			if err == nil {
				if userChannels, exists := prefs[req.Type]; exists {
					channels = userChannels
				}
			}
		}
		
		// Fallback para canais padrão
		if len(channels) == 0 {
			channels = []domain.NotificationChannel{
				domain.NotificationChannelEmail,
				domain.NotificationChannelWhatsApp,
			}
		}
	}

	// Criar notificações para cada canal
	notifications := make([]*domain.Notification, 0, len(channels))
	for _, channel := range channels {
		notification, err := s.createSingleNotification(ctx, req, channel)
		if err != nil {
			s.logger.Error("Failed to create notification for channel", 
				zap.String("channel", string(channel)), 
				zap.Error(err))
			continue
		}
		notifications = append(notifications, notification)
	}

	if len(notifications) == 0 {
		return nil, fmt.Errorf("failed to create any notification")
	}

	// Salvar todas as notificações
	if err := s.notificationRepo.CreateBatch(ctx, notifications); err != nil {
		return nil, fmt.Errorf("failed to save notifications: %w", err)
	}

	// Enfileirar para processamento
	for _, notification := range notifications {
		if notification.ScheduledAt != nil {
			if err := s.queue.EnqueueScheduled(ctx, notification, *notification.ScheduledAt); err != nil {
				s.logger.Error("Failed to enqueue scheduled notification", 
					zap.String("notification_id", notification.ID.String()),
					zap.Error(err))
			}
		} else {
			if err := s.queue.Enqueue(ctx, notification); err != nil {
				s.logger.Error("Failed to enqueue notification", 
					zap.String("notification_id", notification.ID.String()),
					zap.Error(err))
			}
		}
	}

	// Retornar a primeira notificação criada
	return notifications[0], nil
}

// createSingleNotification cria uma notificação para um canal específico
func (s *NotificationService) createSingleNotification(ctx context.Context, req *CreateNotificationRequest, channel domain.NotificationChannel) (*domain.Notification, error) {
	notification := &domain.Notification{
		ID:               uuid.New(),
		TenantID:         req.TenantID,
		Type:             req.Type,
		Channel:          channel,
		Priority:         req.Priority,
		Status:           domain.NotificationStatusPending,
		Subject:          req.Subject,
		Content:          req.Content,
		RecipientID:      req.RecipientID,
		RecipientType:    req.RecipientType,
		RecipientContact: req.RecipientContact,
		TemplateID:       req.TemplateID,
		Variables:        req.Variables,
		Metadata:         req.Metadata,
		RetryCount:       0,
		MaxRetries:       3,
		ScheduledAt:      req.ScheduledAt,
		Tags:             req.Tags,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Se agendada, marcar como scheduled
	if req.ScheduledAt != nil {
		notification.Status = domain.NotificationStatusScheduled
	}

	// Processar template se especificado ou se conteúdo está vazio
	if req.TemplateID != nil || (req.Subject == "" && req.Content == "") {
		if err := s.processTemplate(ctx, notification); err != nil {
			return nil, fmt.Errorf("failed to process template: %w", err)
		}
	}

	return notification, nil
}

// processTemplate processa template para a notificação
func (s *NotificationService) processTemplate(ctx context.Context, notification *domain.Notification) error {
	var template *domain.NotificationTemplate
	var err error

	// Buscar template específico ou padrão
	if notification.TemplateID != nil {
		template, err = s.templateRepo.GetByID(ctx, *notification.TemplateID)
	} else {
		template, err = s.templateRepo.FindByTypeAndChannel(ctx, notification.Type, notification.Channel, notification.TenantID)
		if err == domain.ErrTemplateNotFound {
			template, err = s.templateRepo.GetDefaultTemplate(ctx, notification.Type, notification.Channel)
		}
	}

	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// Processar variáveis no template
	notification.Subject = s.processVariables(template.Subject, notification.Variables)
	notification.Content = s.processVariables(template.Content, notification.Variables)
	notification.TemplateID = &template.ID

	return nil
}

// processVariables substitui variáveis no texto
func (s *NotificationService) processVariables(text string, variables map[string]interface{}) string {
	if variables == nil {
		return text
	}

	result := text
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

// SendNotification envia uma notificação imediatamente
func (s *NotificationService) SendNotification(ctx context.Context, notification *domain.Notification) error {
	s.logger.Debug("Sending notification", 
		zap.String("notification_id", notification.ID.String()),
		zap.String("channel", string(notification.Channel)))

	// Verificar se o provedor está disponível
	provider, exists := s.providers[notification.Channel]
	if !exists {
		return fmt.Errorf("provider not available for channel: %s", notification.Channel)
	}

	if !provider.IsHealthy(ctx) {
		return fmt.Errorf("provider unhealthy for channel: %s", notification.Channel)
	}

	// Marcar como em processamento
	notification.Status = domain.NotificationStatusProcessing
	notification.ProcessingStartedAt = &time.Time{}
	*notification.ProcessingStartedAt = time.Now()

	if err := s.notificationRepo.Update(ctx, notification); err != nil {
		s.logger.Error("Failed to update notification status", zap.Error(err))
	}

	// Enviar através do provedor
	err := provider.Send(ctx, notification)
	
	// Atualizar status baseado no resultado
	now := time.Now()
	notification.ProcessingFinishedAt = &now

	if err != nil {
		notification.Status = domain.NotificationStatusFailed
		notification.FailedAt = &now
		errMsg := err.Error()
		notification.ErrorMessage = &errMsg
		notification.RetryCount++
		
		// Calcular próximo retry se ainda há tentativas disponíveis
		if notification.RetryCount < notification.MaxRetries {
			nextRetry := s.calculateNextRetry(notification.RetryCount)
			notification.NextRetryAt = &nextRetry
		}
	} else {
		notification.Status = domain.NotificationStatusSent
		notification.SentAt = &now
	}

	// Salvar mudanças
	if updateErr := s.notificationRepo.Update(ctx, notification); updateErr != nil {
		s.logger.Error("Failed to update notification after send", zap.Error(updateErr))
	}

	return err
}

// calculateNextRetry calcula o próximo momento de retry
func (s *NotificationService) calculateNextRetry(retryCount int) time.Time {
	// Backoff exponencial: 1min, 5min, 15min
	delays := []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		15 * time.Minute,
	}

	if retryCount-1 >= len(delays) {
		return time.Now().Add(delays[len(delays)-1])
	}

	return time.Now().Add(delays[retryCount-1])
}

// GetNotification busca notificação por ID
func (s *NotificationService) GetNotification(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	return s.notificationRepo.GetByID(ctx, id)
}

// ListNotifications lista notificações com filtros
func (s *NotificationService) ListNotifications(ctx context.Context, tenantID uuid.UUID, filters domain.NotificationFilters) ([]*domain.Notification, error) {
	return s.notificationRepo.FindByTenantID(ctx, tenantID, filters)
}

// GetNotificationStatistics obtém estatísticas de notificações
func (s *NotificationService) GetNotificationStatistics(ctx context.Context, tenantID uuid.UUID, period time.Duration) (*domain.NotificationStatistics, error) {
	return s.notificationRepo.GetStatistics(ctx, tenantID, period)
}

// ProcessPendingNotifications processa notificações pendentes
func (s *NotificationService) ProcessPendingNotifications(ctx context.Context, limit int) error {
	s.logger.Debug("Processing pending notifications", zap.Int("limit", limit))

	notifications, err := s.notificationRepo.FindPending(ctx, limit)
	if err != nil {
		return fmt.Errorf("failed to find pending notifications: %w", err)
	}

	for _, notification := range notifications {
		if err := s.SendNotification(ctx, notification); err != nil {
			s.logger.Error("Failed to send pending notification", 
				zap.String("notification_id", notification.ID.String()),
				zap.Error(err))
		}
	}

	return nil
}

// ProcessScheduledNotifications processa notificações agendadas
func (s *NotificationService) ProcessScheduledNotifications(ctx context.Context, limit int) error {
	s.logger.Debug("Processing scheduled notifications", zap.Int("limit", limit))

	notifications, err := s.notificationRepo.FindScheduled(ctx, time.Now(), limit)
	if err != nil {
		return fmt.Errorf("failed to find scheduled notifications: %w", err)
	}

	for _, notification := range notifications {
		notification.Status = domain.NotificationStatusPending
		if err := s.notificationRepo.Update(ctx, notification); err != nil {
			s.logger.Error("Failed to update scheduled notification", 
				zap.String("notification_id", notification.ID.String()),
				zap.Error(err))
			continue
		}

		if err := s.queue.Enqueue(ctx, notification); err != nil {
			s.logger.Error("Failed to enqueue scheduled notification", 
				zap.String("notification_id", notification.ID.String()),
				zap.Error(err))
		}
	}

	return nil
}

// ProcessRetries processa notificações para reenvio
func (s *NotificationService) ProcessRetries(ctx context.Context, limit int) error {
	s.logger.Debug("Processing retry notifications", zap.Int("limit", limit))

	notifications, err := s.notificationRepo.FindForRetry(ctx, limit)
	if err != nil {
		return fmt.Errorf("failed to find notifications for retry: %w", err)
	}

	for _, notification := range notifications {
		if err := s.SendNotification(ctx, notification); err != nil {
			s.logger.Error("Failed to retry notification", 
				zap.String("notification_id", notification.ID.String()),
				zap.Error(err))
		}
	}

	return nil
}

// CleanupExpiredNotifications remove notificações expiradas
func (s *NotificationService) CleanupExpiredNotifications(ctx context.Context, retentionDays int) (int64, error) {
	s.logger.Debug("Cleaning up expired notifications", zap.Int("retention_days", retentionDays))

	expiredBefore := time.Now().AddDate(0, 0, -retentionDays)
	deleted, err := s.notificationRepo.DeleteExpired(ctx, expiredBefore)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired notifications: %w", err)
	}

	s.logger.Info("Expired notifications cleaned up", zap.Int64("deleted_count", deleted))
	return deleted, nil
}

// CancelNotification cancela uma notificação
func (s *NotificationService) CancelNotification(ctx context.Context, id uuid.UUID) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("notification not found: %w", err)
	}

	if notification.Status != domain.NotificationStatusPending && notification.Status != domain.NotificationStatusScheduled {
		return fmt.Errorf("cannot cancel notification with status: %s", notification.Status)
	}

	notification.Status = domain.NotificationStatusCancelled
	notification.UpdatedAt = time.Now()

	return s.notificationRepo.Update(ctx, notification)
}