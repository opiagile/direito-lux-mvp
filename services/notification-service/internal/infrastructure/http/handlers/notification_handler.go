package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/application/services"
	"github.com/direito-lux/notification-service/internal/domain"
)

// NotificationHandler handler para notificações
type NotificationHandler struct {
	notificationService *services.NotificationService
	logger              *zap.Logger
}

// NewNotificationHandler cria nova instância do handler
func NewNotificationHandler(notificationService *services.NotificationService, logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:              logger,
	}
}

// CreateNotification cria uma nova notificação
// @Summary Criar notificação
// @Description Cria uma nova notificação
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body services.CreateNotificationRequest true "Dados da notificação"
// @Success 201 {object} domain.Notification
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req services.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// Extrair tenant ID do contexto (definido por middleware)
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Tenant ID not found in context",
		})
		return
	}
	req.TenantID = tenantID.(uuid.UUID)

	notification, err := h.notificationService.CreateNotification(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create notification", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create notification",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetNotification busca notificação por ID
// @Summary Buscar notificação
// @Description Busca uma notificação por ID
// @Tags notifications
// @Produce json
// @Param id path string true "ID da notificação"
// @Success 200 {object} domain.Notification
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id} [get]
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid notification ID",
			Message: err.Error(),
		})
		return
	}

	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Notification not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get notification", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get notification",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// ListNotifications lista notificações com filtros
// @Summary Listar notificações
// @Description Lista notificações com filtros opcionais
// @Tags notifications
// @Produce json
// @Param type query string false "Tipo de notificação"
// @Param channel query string false "Canal de notificação"
// @Param status query string false "Status da notificação"
// @Param priority query string false "Prioridade"
// @Param limit query int false "Limite de resultados" default(50)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} ListNotificationsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications [get]
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Parse query parameters
	filters := domain.NotificationFilters{}

	if typeStr := c.Query("type"); typeStr != "" {
		notificationType := domain.NotificationType(typeStr)
		filters.Type = &notificationType
	}

	if channelStr := c.Query("channel"); channelStr != "" {
		channel := domain.NotificationChannel(channelStr)
		filters.Channel = &channel
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.NotificationStatus(statusStr)
		filters.Status = &status
	}

	if priorityStr := c.Query("priority"); priorityStr != "" {
		priority := domain.NotificationPriority(priorityStr)
		filters.Priority = &priority
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			filters.UserID = &userID
		}
	}

	// Parse dates
	if createdAfterStr := c.Query("created_after"); createdAfterStr != "" {
		if createdAfter, err := time.Parse(time.RFC3339, createdAfterStr); err == nil {
			filters.CreatedAfter = &createdAfter
		}
	}

	if createdBeforeStr := c.Query("created_before"); createdBeforeStr != "" {
		if createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr); err == nil {
			filters.CreatedBefore = &createdBefore
		}
	}

	// Parse pagination
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	filters.Limit = limit

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	filters.Offset = offset

	// Parse ordering
	filters.OrderBy = c.DefaultQuery("order_by", "created_at")
	filters.OrderDir = c.DefaultQuery("order_dir", "DESC")

	notifications, err := h.notificationService.ListNotifications(c.Request.Context(), tenantID, filters)
	if err != nil {
		h.logger.Error("Failed to list notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list notifications",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListNotificationsResponse{
		Notifications: notifications,
		Total:         len(notifications),
		Limit:         limit,
		Offset:        offset,
	})
}

// GetNotificationStatistics obtém estatísticas de notificações
// @Summary Estatísticas de notificações
// @Description Obtém estatísticas de notificações para um período
// @Tags notifications
// @Produce json
// @Param period query string false "Período (24h, 7d, 30d)" default(7d)
// @Success 200 {object} domain.NotificationStatistics
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/statistics [get]
func (h *NotificationHandler) GetNotificationStatistics(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	periodStr := c.DefaultQuery("period", "7d")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		// Tentar formatos customizados
		switch periodStr {
		case "24h", "1d":
			period = 24 * time.Hour
		case "7d", "1w":
			period = 7 * 24 * time.Hour
		case "30d", "1m":
			period = 30 * 24 * time.Hour
		case "90d", "3m":
			period = 90 * 24 * time.Hour
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid period format",
				Message: "Use formats like 24h, 7d, 30d",
			})
			return
		}
	}

	stats, err := h.notificationService.GetNotificationStatistics(c.Request.Context(), tenantID, period)
	if err != nil {
		h.logger.Error("Failed to get notification statistics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get statistics",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// CancelNotification cancela uma notificação
// @Summary Cancelar notificação
// @Description Cancela uma notificação pendente ou agendada
// @Tags notifications
// @Produce json
// @Param id path string true "ID da notificação"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id}/cancel [post]
func (h *NotificationHandler) CancelNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid notification ID",
			Message: err.Error(),
		})
		return
	}

	err = h.notificationService.CancelNotification(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Notification not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to cancel notification", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to cancel notification",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notification cancelled successfully",
	})
}

// SendBulkNotifications envia notificações em lote
// @Summary Enviar notificações em lote
// @Description Cria e envia múltiplas notificações
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body BulkNotificationRequest true "Dados das notificações"
// @Success 202 {object} BulkNotificationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/bulk [post]
func (h *NotificationHandler) SendBulkNotifications(c *gin.Context) {
	var req BulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind bulk request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Processar cada notificação
	response := BulkNotificationResponse{
		Requested: len(req.Notifications),
		Processed: 0,
		Failed:    0,
		Results:   make([]BulkNotificationResult, len(req.Notifications)),
	}

	for i, notifReq := range req.Notifications {
		notifReq.TenantID = tenantID
		
		notification, err := h.notificationService.CreateNotification(c.Request.Context(), &notifReq)
		if err != nil {
			response.Failed++
			response.Results[i] = BulkNotificationResult{
				Index:   i,
				Success: false,
				Error:   err.Error(),
			}
			continue
		}

		response.Processed++
		response.Results[i] = BulkNotificationResult{
			Index:          i,
			Success:        true,
			NotificationID: &notification.ID,
		}
	}

	statusCode := http.StatusAccepted
	if response.Failed > 0 && response.Processed == 0 {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, response)
}

// Response types
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ListNotificationsResponse struct {
	Notifications []*domain.Notification `json:"notifications"`
	Total         int                    `json:"total"`
	Limit         int                    `json:"limit"`
	Offset        int                    `json:"offset"`
}

type BulkNotificationRequest struct {
	Notifications []services.CreateNotificationRequest `json:"notifications"`
}

type BulkNotificationResponse struct {
	Requested int                       `json:"requested"`
	Processed int                       `json:"processed"`
	Failed    int                       `json:"failed"`
	Results   []BulkNotificationResult  `json:"results"`
}

type BulkNotificationResult struct {
	Index          int        `json:"index"`
	Success        bool       `json:"success"`
	NotificationID *uuid.UUID `json:"notification_id,omitempty"`
	Error          string     `json:"error,omitempty"`
}