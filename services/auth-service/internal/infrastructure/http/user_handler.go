package http

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	
	"github.com/direito-lux/auth-service/internal/application"
	"github.com/direito-lux/auth-service/internal/infrastructure/logging"
)

// UserHandler implementa os endpoints de gerenciamento de usuários
type UserHandler struct {
	userService *application.UserService
	logger      *zap.Logger
}

// NewUserHandler cria uma nova instância do handler de usuários
func NewUserHandler(userService *application.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser godoc
// @Summary Criar usuário
// @Description Cria um novo usuário no sistema
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body application.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} application.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req application.CreateUserRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	// Extrair tenant ID do contexto (middleware deve definir isso)
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Tenant ID obrigatório",
			Message: "Tenant ID não encontrado no contexto",
		})
		return
	}
	req.TenantID = tenantID.(string)
	
	logging.LogInfo(c.Request.Context(), h.logger, "Criando usuário",
		zap.String("email", req.Email),
		zap.String("tenant_id", req.TenantID),
		zap.String("operation", "create_user"),
	)
	
	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao criar usuário",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("tenant_id", req.TenantID),
			zap.String("operation", "create_user"),
		)
		
		switch err.Error() {
		case "email já está em uso":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Email já existe",
				Message: "Este email já está em uso",
			})
		case "email inválido", "senha muito fraca", "papel de usuário inválido", "status de usuário inválido":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Dados inválidos",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Usuário criado com sucesso",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("operation", "create_user"),
	)
	
	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Buscar usuário
// @Description Busca um usuário por ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do usuário"
// @Success 200 {object} application.UserDTO
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	
	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Usuário não encontrado",
				Message: "Usuário não foi encontrado",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// GetUsersByTenant godoc
// @Summary Listar usuários do tenant
// @Description Lista usuários do tenant atual com paginação
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} application.UserDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsersByTenant(c *gin.Context) {
	// Extrair tenant ID do contexto
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Tenant ID obrigatório",
			Message: "Tenant ID não encontrado no contexto",
		})
		return
	}
	
	// Parâmetros de paginação
	limit := 20
	offset := 0
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	
	users, err := h.userService.GetUsersByTenant(c.Request.Context(), tenantID.(string), limit, offset)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao buscar usuários do tenant",
			zap.Error(err),
			zap.String("tenant_id", tenantID.(string)),
			zap.String("operation", "get_users_by_tenant"),
		)
		
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro interno",
			Message: "Erro interno do servidor",
		})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Atualizar usuário
// @Description Atualiza dados de um usuário
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do usuário"
// @Param request body application.UpdateUserRequest true "Dados para atualização"
// @Success 200 {object} application.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	
	var req application.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	user, err := h.userService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao atualizar usuário",
			zap.Error(err),
			zap.String("user_id", userID),
			zap.String("operation", "update_user"),
		)
		
		switch err.Error() {
		case "usuário não encontrado":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Usuário não encontrado",
				Message: "Usuário não foi encontrado",
			})
		case "papel de usuário inválido", "status de usuário inválido":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Dados inválidos",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Usuário atualizado com sucesso",
		zap.String("user_id", userID),
		zap.String("operation", "update_user"),
	)
	
	c.JSON(http.StatusOK, user)
}

// ChangePassword godoc
// @Summary Alterar senha
// @Description Altera a senha do usuário atual
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body application.ChangePasswordRequest true "Dados para alteração de senha"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Extrair user ID do contexto (middleware deve definir isso)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Usuário não autenticado",
			Message: "User ID não encontrado no contexto",
		})
		return
	}
	
	var req application.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	err := h.userService.ChangePassword(c.Request.Context(), userID.(string), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao alterar senha",
			zap.Error(err),
			zap.String("user_id", userID.(string)),
			zap.String("operation", "change_password"),
		)
		
		switch err.Error() {
		case "usuário não encontrado":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Usuário não encontrado",
				Message: "Usuário não foi encontrado",
			})
		case "senha atual incorreta":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Senha incorreta",
				Message: "Senha atual está incorreta",
			})
		case "senha muito fraca":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Senha fraca",
				Message: "A nova senha não atende aos critérios de segurança",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Senha alterada com sucesso",
		zap.String("user_id", userID.(string)),
		zap.String("operation", "change_password"),
	)
	
	c.JSON(http.StatusOK, MessageResponse{
		Message: "Senha alterada com sucesso",
	})
}

// DeleteUser godoc
// @Summary Deletar usuário
// @Description Remove um usuário do sistema
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do usuário"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao deletar usuário",
			zap.Error(err),
			zap.String("user_id", userID),
			zap.String("operation", "delete_user"),
		)
		
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Usuário não encontrado",
				Message: "Usuário não foi encontrado",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Usuário deletado com sucesso",
		zap.String("user_id", userID),
		zap.String("operation", "delete_user"),
	)
	
	c.JSON(http.StatusOK, MessageResponse{
		Message: "Usuário deletado com sucesso",
	})
}