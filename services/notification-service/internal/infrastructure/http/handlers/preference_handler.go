package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
	"github.com/direito-lux/notification-service/internal/infrastructure/repository"
)

// PreferenceHandler handler para preferências de notificação
type PreferenceHandler struct {
	preferenceRepo domain.NotificationPreferenceRepository
	logger         *zap.Logger
}

// NewPreferenceHandler cria nova instância do handler
func NewPreferenceHandler(preferenceRepo domain.NotificationPreferenceRepository, logger *zap.Logger) *PreferenceHandler {
	return &PreferenceHandler{
		preferenceRepo: preferenceRepo,
		logger:         logger,
	}
}

// GetUserPreferences obtém preferências do usuário
// @Summary Obter preferências do usuário
// @Description Obtém todas as preferências de notificação do usuário
// @Tags preferences
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {object} GetUserPreferencesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/users/{user_id} [get]
func (h *PreferenceHandler) GetUserPreferences(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	preferences, err := h.preferenceRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preferences",
			Message: err.Error(),
		})
		return
	}

	// Converter para mapa para facilitar uso
	preferencesMap, err := h.preferenceRepo.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user preferences map", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preferences",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetUserPreferencesResponse{
		UserID:      userID,
		Preferences: preferences,
		ChannelMap:  preferencesMap,
	})
}

// UpdateUserPreferences atualiza preferências do usuário
// @Summary Atualizar preferências do usuário
// @Description Atualiza as preferências de notificação do usuário
// @Tags preferences
// @Accept json
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Param request body UpdateUserPreferencesRequest true "Novas preferências"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/users/{user_id} [put]
func (h *PreferenceHandler) UpdateUserPreferences(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var req UpdateUserPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind preferences request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Converter request para domain objects
	preferences := make([]*domain.NotificationPreference, len(req.Preferences))
	for i, prefReq := range req.Preferences {
		preference := &domain.NotificationPreference{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   userID,
			Type:     prefReq.Type,
			Channels: prefReq.Channels,
			Enabled:  prefReq.Enabled,
		}
		preferences[i] = preference
	}

	err = h.preferenceRepo.UpsertUserPreferences(c.Request.Context(), userID, preferences)
	if err != nil {
		h.logger.Error("Failed to update user preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update preferences",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Preferences updated successfully",
	})
}

// CreateDefaultPreferences cria preferências padrão para um usuário
// @Summary Criar preferências padrão
// @Description Cria preferências padrão para um novo usuário
// @Tags preferences
// @Accept json
// @Produce json
// @Param request body CreateDefaultPreferencesRequest true "Dados do usuário"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/defaults [post]
func (h *PreferenceHandler) CreateDefaultPreferences(c *gin.Context) {
	var req CreateDefaultPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind default preferences request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Usar método do repositório que cria preferências padrão
	repo, ok := h.preferenceRepo.(*repository.PostgresPreferenceRepository)
	if !ok {
		h.logger.Error("Repository does not support creating default preferences")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal error",
			Message: "Repository implementation issue",
		})
		return
	}

	err := repo.CreateDefaultPreferences(c.Request.Context(), tenantID, req.UserID)
	if err != nil {
		h.logger.Error("Failed to create default preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create default preferences",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Default preferences created successfully",
	})
}

// GetPreferenceByType obtém preferência específica por tipo
// @Summary Obter preferência por tipo
// @Description Obtém preferência específica do usuário por tipo de notificação
// @Tags preferences
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Param type path string true "Tipo de notificação"
// @Success 200 {object} domain.NotificationPreference
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/users/{user_id}/types/{type} [get]
func (h *PreferenceHandler) GetPreferenceByType(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	typeStr := c.Param("type")
	notificationType := domain.NotificationType(typeStr)

	preference, err := h.preferenceRepo.FindByUserAndType(c.Request.Context(), userID, notificationType)
	if err != nil {
		if err == domain.ErrPreferenceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Preference not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get preference by type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preference",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, preference)
}

// UpdatePreferenceByType atualiza preferência específica por tipo
// @Summary Atualizar preferência por tipo
// @Description Atualiza preferência específica do usuário por tipo de notificação
// @Tags preferences
// @Accept json
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Param type path string true "Tipo de notificação"
// @Param request body UpdatePreferenceRequest true "Dados da preferência"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/users/{user_id}/types/{type} [put]
func (h *PreferenceHandler) UpdatePreferenceByType(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	typeStr := c.Param("type")
	notificationType := domain.NotificationType(typeStr)

	var req UpdatePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind preference update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// Buscar preferência existente
	preference, err := h.preferenceRepo.FindByUserAndType(c.Request.Context(), userID, notificationType)
	if err != nil {
		if err == domain.ErrPreferenceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Preference not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to find preference", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to find preference",
			Message: err.Error(),
		})
		return
	}

	// Aplicar mudanças
	if req.Channels != nil {
		preference.Channels = *req.Channels
	}
	if req.Enabled != nil {
		preference.Enabled = *req.Enabled
	}

	err = h.preferenceRepo.Update(c.Request.Context(), preference)
	if err != nil {
		h.logger.Error("Failed to update preference", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update preference",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Preference updated successfully",
	})
}

// DisableAllNotifications desabilita todas as notificações do usuário
// @Summary Desabilitar todas as notificações
// @Description Desabilita todas as notificações para um usuário
// @Tags preferences
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/users/{user_id}/disable-all [post]
func (h *PreferenceHandler) DisableAllNotifications(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	preferences, err := h.preferenceRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preferences",
			Message: err.Error(),
		})
		return
	}

	// Desabilitar todas as preferências
	for _, preference := range preferences {
		preference.Enabled = false
	}

	err = h.preferenceRepo.UpsertUserPreferences(c.Request.Context(), userID, preferences)
	if err != nil {
		h.logger.Error("Failed to disable all notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to disable notifications",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "All notifications disabled successfully",
	})
}

// EnableAllNotifications habilita todas as notificações do usuário
// @Summary Habilitar todas as notificações
// @Description Habilita todas as notificações para um usuário
// @Tags preferences
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /preferences/users/{user_id}/enable-all [post]
func (h *PreferenceHandler) EnableAllNotifications(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	preferences, err := h.preferenceRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preferences",
			Message: err.Error(),
		})
		return
	}

	// Habilitar todas as preferências
	for _, preference := range preferences {
		preference.Enabled = true
	}

	err = h.preferenceRepo.UpsertUserPreferences(c.Request.Context(), userID, preferences)
	if err != nil {
		h.logger.Error("Failed to enable all notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to enable notifications",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "All notifications enabled successfully",
	})
}

// Request/Response types
type GetUserPreferencesResponse struct {
	UserID      uuid.UUID                                              `json:"user_id"`
	Preferences []*domain.NotificationPreference                       `json:"preferences"`
	ChannelMap  map[domain.NotificationType][]domain.NotificationChannel `json:"channel_map"`
}

type UpdateUserPreferencesRequest struct {
	Preferences []PreferenceRequest `json:"preferences"`
}

type PreferenceRequest struct {
	Type     domain.NotificationType     `json:"type" validate:"required"`
	Channels []domain.NotificationChannel `json:"channels" validate:"required"`
	Enabled  bool                        `json:"enabled"`
}

type CreateDefaultPreferencesRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type UpdatePreferenceRequest struct {
	Channels *[]domain.NotificationChannel `json:"channels,omitempty"`
	Enabled  *bool                         `json:"enabled,omitempty"`
}