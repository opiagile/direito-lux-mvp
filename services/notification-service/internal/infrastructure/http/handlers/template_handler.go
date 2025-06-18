package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/application/services"
	"github.com/direito-lux/notification-service/internal/domain"
)

// TemplateHandler handler para templates
type TemplateHandler struct {
	templateService *services.TemplateService
	logger          *zap.Logger
}

// NewTemplateHandler cria nova instância do handler
func NewTemplateHandler(templateService *services.TemplateService, logger *zap.Logger) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		logger:          logger,
	}
}

// CreateTemplate cria um novo template
// @Summary Criar template
// @Description Cria um novo template de notificação
// @Tags templates
// @Accept json
// @Produce json
// @Param request body services.CreateTemplateRequest true "Dados do template"
// @Success 201 {object} domain.NotificationTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates [post]
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req services.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind template request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// Extrair tenant ID do contexto
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	req.TenantID = &tenantID

	// Validar template
	if err := h.templateService.ValidateTemplate(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template",
			Message: err.Error(),
		})
		return
	}

	template, err := h.templateService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// GetTemplate busca template por ID
// @Summary Buscar template
// @Description Busca um template por ID
// @Tags templates
// @Produce json
// @Param id path string true "ID do template"
// @Success 200 {object} domain.NotificationTemplate
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id} [get]
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	template, err := h.templateService.GetTemplate(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate atualiza um template
// @Summary Atualizar template
// @Description Atualiza um template existente
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "ID do template"
// @Param request body services.UpdateTemplateRequest true "Dados para atualização"
// @Success 200 {object} domain.NotificationTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id} [put]
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	var req services.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	template, err := h.templateService.UpdateTemplate(c.Request.Context(), id, &req)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to update template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate remove um template
// @Summary Deletar template
// @Description Remove um template
// @Tags templates
// @Produce json
// @Param id path string true "ID do template"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id} [delete]
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	err = h.templateService.DeleteTemplate(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to delete template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Template deleted successfully",
	})
}

// ListTemplates lista templates com filtros
// @Summary Listar templates
// @Description Lista templates com filtros opcionais
// @Tags templates
// @Produce json
// @Param type query string false "Tipo de notificação"
// @Param channel query string false "Canal de notificação"
// @Param status query string false "Status do template"
// @Param is_system query bool false "Se é template do sistema"
// @Param search query string false "Busca por nome ou conteúdo"
// @Param limit query int false "Limite de resultados" default(50)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} ListTemplatesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates [get]
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Parse query parameters
	filters := domain.TemplateFilters{}

	if typeStr := c.Query("type"); typeStr != "" {
		notificationType := domain.NotificationType(typeStr)
		filters.Type = &notificationType
	}

	if channelStr := c.Query("channel"); channelStr != "" {
		channel := domain.NotificationChannel(channelStr)
		filters.Channel = &channel
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.TemplateStatus(statusStr)
		filters.Status = &status
	}

	if isSystemStr := c.Query("is_system"); isSystemStr != "" {
		if isSystem, err := strconv.ParseBool(isSystemStr); err == nil {
			filters.IsSystem = &isSystem
		}
	}

	filters.Search = c.Query("search")

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

	templates, err := h.templateService.ListTemplates(c.Request.Context(), &tenantID, filters)
	if err != nil {
		h.logger.Error("Failed to list templates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list templates",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListTemplatesResponse{
		Templates: templates,
		Total:     len(templates),
		Limit:     limit,
		Offset:    offset,
	})
}

// PreviewTemplate faz preview de um template
// @Summary Preview de template
// @Description Faz preview de um template com variáveis
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "ID do template"
// @Param request body TemplatePreviewRequest true "Variáveis para o preview"
// @Success 200 {object} services.TemplatePreview
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id}/preview [post]
func (h *TemplateHandler) PreviewTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	var req TemplatePreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind preview request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	preview, err := h.templateService.PreviewTemplate(c.Request.Context(), id, req.Variables)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to preview template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to preview template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, preview)
}

// DuplicateTemplate duplica um template
// @Summary Duplicar template
// @Description Cria uma cópia de um template existente
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "ID do template"
// @Param request body DuplicateTemplateRequest true "Dados para duplicação"
// @Success 201 {object} domain.NotificationTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id}/duplicate [post]
func (h *TemplateHandler) DuplicateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	var req DuplicateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind duplicate request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	template, err := h.templateService.DuplicateTemplate(c.Request.Context(), id, req.Name, &tenantID)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to duplicate template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to duplicate template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// ActivateTemplate ativa um template
// @Summary Ativar template
// @Description Ativa um template
// @Tags templates
// @Produce json
// @Param id path string true "ID do template"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id}/activate [post]
func (h *TemplateHandler) ActivateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	err = h.templateService.ActivateTemplate(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to activate template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to activate template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Template activated successfully",
	})
}

// DeactivateTemplate desativa um template
// @Summary Desativar template
// @Description Desativa um template
// @Tags templates
// @Produce json
// @Param id path string true "ID do template"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates/{id}/deactivate [post]
func (h *TemplateHandler) DeactivateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid template ID",
			Message: err.Error(),
		})
		return
	}

	err = h.templateService.DeactivateTemplate(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Template not found",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to deactivate template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to deactivate template",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Template deactivated successfully",
	})
}

// Response types
type ListTemplatesResponse struct {
	Templates []*domain.NotificationTemplate `json:"templates"`
	Total     int                            `json:"total"`
	Limit     int                            `json:"limit"`
	Offset    int                            `json:"offset"`
}

type TemplatePreviewRequest struct {
	Variables map[string]interface{} `json:"variables"`
}

type DuplicateTemplateRequest struct {
	Name string `json:"name" validate:"required"`
}