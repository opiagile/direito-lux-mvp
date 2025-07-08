package http

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	
	"github.com/direito-lux/auth-service/internal/application"
	"github.com/direito-lux/auth-service/internal/infrastructure/logging"
)

// AuthHandler implementa os endpoints de autenticação
type AuthHandler struct {
	authService *application.AuthService
	logger      *zap.Logger
}

// NewAuthHandler cria uma nova instância do handler de autenticação
func NewAuthHandler(authService *application.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Login godoc
// @Summary Login de usuário
// @Description Autentica um usuário e retorna tokens de acesso
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.LoginRequest true "Dados de login"
// @Success 200 {object} application.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req application.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao fazer bind da requisição", err,
			zap.String("operation", "login"),
		)
		
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	// Extrair informações da requisição
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")
	
	// Tenant ID será obtido do banco baseado no email do usuário
	// Não é necessário fornecer X-Tenant-ID no login
	tenantID := c.GetHeader("X-Tenant-ID")
	req.TenantID = tenantID // Pode ser vazio, será resolvido no service
	
	logging.LogInfo(c.Request.Context(), h.logger, "Tentativa de login",
		zap.String("email", req.Email),
		zap.String("tenant_id", req.TenantID),
		zap.String("ip_address", req.IPAddress),
		zap.String("operation", "login"),
	)
	
	// Realizar login
	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro no login", err,
			zap.String("email", req.Email),
			zap.String("tenant_id", req.TenantID),
			zap.String("operation", "login"),
		)
		
		// Mapear erros específicos para códigos HTTP
		switch err.Error() {
		case "muitas tentativas de login falharam":
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Muitas tentativas",
				Message: "Muitas tentativas de login falharam. Tente novamente em 15 minutos.",
			})
		case "usuário não encontrado", "credenciais inválidas":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Credenciais inválidas",
				Message: "Email ou senha incorretos",
			})
		case "usuário inativo":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Usuário inativo",
				Message: "Sua conta está inativa. Entre em contato com o suporte.",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Login realizado com sucesso",
		zap.String("user_id", response.User.ID),
		zap.String("email", response.User.Email),
		zap.String("tenant_id", req.TenantID),
		zap.String("operation", "login"),
	)
	
	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Renovar token de acesso
// @Description Renova o token de acesso usando um refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} application.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req application.RefreshTokenRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Tentativa de refresh token",
		zap.String("operation", "refresh_token"),
	)
	
	response, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro no refresh token", err,
			zap.String("operation", "refresh_token"),
		)
		
		switch err.Error() {
		case "refresh token não encontrado", "token inválido", "token expirado", "token já foi utilizado":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Token inválido",
				Message: "Refresh token inválido ou expirado",
			})
		case "usuário inativo":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Usuário inativo",
				Message: "Sua conta está inativa",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Refresh token realizado com sucesso",
		zap.String("user_id", response.User.ID),
		zap.String("operation", "refresh_token"),
	)
	
	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout de usuário
// @Description Invalida a sessão atual do usuário
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token obrigatório",
			Message: "Header Authorization é obrigatório",
		})
		return
	}
	
	// Extrair token do header Authorization
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token inválido",
			Message: "Token deve estar no formato 'Bearer <token>'",
		})
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Tentativa de logout",
		zap.String("operation", "logout"),
	)
	
	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro no logout", err,
			zap.String("operation", "logout"),
		)
		
		if err.Error() == "sessão não encontrada" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Sessão inválida",
				Message: "Sessão não encontrada ou já expirada",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Logout realizado com sucesso",
		zap.String("operation", "logout"),
	)
	
	c.JSON(http.StatusOK, MessageResponse{
		Message: "Logout realizado com sucesso",
	})
}

// ValidateToken godoc
// @Summary Validar token de acesso
// @Description Valida um token de acesso e retorna informações do usuário
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} application.UserDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token obrigatório",
			Message: "Header Authorization é obrigatório",
		})
		return
	}
	
	// Extrair token do header Authorization
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token inválido",
			Message: "Token deve estar no formato 'Bearer <token>'",
		})
		return
	}
	
	user, err := h.authService.ValidateToken(c.Request.Context(), token)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro na validação de token", err,
			zap.String("operation", "validate_token"),
		)
		
		switch err.Error() {
		case "token inválido", "usuário não encontrado", "usuário inativo":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Token inválido",
				Message: "Token inválido ou usuário inativo",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// Register godoc
// @Summary Registro de novo usuário e tenant
// @Description Cria novo tenant e usuário administrador (endpoint público)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.RegisterRequest true "Dados de registro"
// @Success 201 {object} application.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req application.RegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro ao fazer bind da requisição de registro", err,
			zap.String("operation", "register"),
		)
		
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Tentativa de registro",
		zap.String("email", req.User.Email),
		zap.String("tenant_name", req.Tenant.Name),
		zap.String("operation", "register"),
	)
	
	// Realizar registro
	response, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro no registro", err,
			zap.String("email", req.User.Email),
			zap.String("operation", "register"),
		)
		
		switch err.Error() {
		case "email já existe", "documento já existe":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Dados já existem",
				Message: "Email ou documento já cadastrados",
			})
		case "dados inválidos":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Dados inválidos",
				Message: "Verifique os dados fornecidos",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Registro realizado com sucesso",
		zap.String("user_id", response.User.ID),
		zap.String("tenant_id", response.Tenant.ID),
		zap.String("operation", "register"),
	)
	
	c.JSON(http.StatusCreated, response)
}

// ForgotPassword godoc
// @Summary Solicitar recuperação de senha
// @Description Envia email com link para recuperação de senha
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.ForgotPasswordRequest true "Email para recuperação"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req application.ForgotPasswordRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Solicitação de recuperação de senha",
		zap.String("email", req.Email),
		zap.String("operation", "forgot_password"),
	)
	
	err := h.authService.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro na recuperação de senha", err,
			zap.String("email", req.Email),
			zap.String("operation", "forgot_password"),
		)
		
		switch err.Error() {
		case "usuário não encontrado":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Usuário não encontrado",
				Message: "Email não cadastrado no sistema",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, MessageResponse{
		Message: "Se o email existir, você receberá instruções para recuperação da senha",
	})
}

// ResetPassword godoc
// @Summary Resetar senha com token
// @Description Redefine senha do usuário usando token de recuperação
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.ResetPasswordRequest true "Token e nova senha"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req application.ResetPasswordRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	
	logging.LogInfo(c.Request.Context(), h.logger, "Tentativa de reset de senha",
		zap.String("operation", "reset_password"),
	)
	
	err := h.authService.ResetPassword(c.Request.Context(), req)
	if err != nil {
		logging.LogError(c.Request.Context(), h.logger, "Erro no reset de senha", err,
			zap.String("operation", "reset_password"),
		)
		
		switch err.Error() {
		case "token inválido", "token expirado", "token já utilizado":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Token inválido",
				Message: "Token de recuperação inválido ou expirado",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro interno",
				Message: "Erro interno do servidor",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, MessageResponse{
		Message: "Senha alterada com sucesso",
	})
}