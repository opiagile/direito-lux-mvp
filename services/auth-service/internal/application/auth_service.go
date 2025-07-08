package application

import (
	"context"
	"time"
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	
	"github.com/direito-lux/auth-service/internal/domain"
	"github.com/direito-lux/auth-service/internal/infrastructure/events"
)

// AuthService implementa os casos de uso de autenticação
type AuthService struct {
	userRepo              domain.UserRepository
	sessionRepo           domain.SessionRepository
	refreshTokenRepo      domain.RefreshTokenRepository
	loginAttemptRepo      domain.LoginAttemptRepository
	passwordResetTokenRepo domain.PasswordResetTokenRepository
	eventBus              events.EventBus
	jwtSecret             string
	jwtExpiryHours        int
	refreshExpiryDays     int
}

// NewAuthService cria uma nova instância do serviço de autenticação
func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	loginAttemptRepo domain.LoginAttemptRepository,
	passwordResetTokenRepo domain.PasswordResetTokenRepository,
	eventBus events.EventBus,
	jwtSecret string,
	jwtExpiryHours int,
	refreshExpiryDays int,
) *AuthService {
	return &AuthService{
		userRepo:               userRepo,
		sessionRepo:            sessionRepo,
		refreshTokenRepo:       refreshTokenRepo,
		loginAttemptRepo:       loginAttemptRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		eventBus:               eventBus,
		jwtSecret:              jwtSecret,
		jwtExpiryHours:         jwtExpiryHours,
		refreshExpiryDays:      refreshExpiryDays,
	}
}

// LoginRequest representa uma solicitação de login
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	TenantID  string `json:"tenant_id"` // Optional - will be resolved from user's tenant
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

// LoginResponse representa a resposta de um login bem-sucedido
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *UserDTO  `json:"user"`
}

// UserDTO representa um usuário para transferência de dados
type UserDTO struct {
	ID        string               `json:"id"`
	Email     string               `json:"email"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Role      domain.UserRole      `json:"role"`
	Status    domain.UserStatus    `json:"status"`
	TenantID  string               `json:"tenant_id"`
	CreatedAt time.Time            `json:"created_at"`
}

// Login autentica um usuário e cria uma sessão
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Verificar tentativas recentes de login
	recentFailures, err := s.loginAttemptRepo.GetRecentFailures(req.Email, 15)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar tentativas de login: %w", err)
	}
	
	if recentFailures >= 5 {
		// Registrar tentativa bloqueada
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		// TODO: Publicar evento de login bloqueado
		// event := domain.LoginFailedEvent{
		//	Email:     req.Email,
		//	TenantID:  req.TenantID,
		//	IPAddress: req.IPAddress,
		//	UserAgent: req.UserAgent,
		//	Reason:    "too_many_failures",
		//	FailedAt:  time.Now(),
		// }
		// // s.eventBus.Publish(ctx, event)
		
		return nil, domain.ErrTooManyFailures
	}
	
	// Buscar usuário por email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		// TODO: Publicar evento de login falhou
		// event := domain.LoginFailedEvent{
		//	Email:     req.Email,
		//	TenantID:  req.TenantID,
		//	IPAddress: req.IPAddress,
		//	UserAgent: req.UserAgent,
		//	Reason:    "user_not_found",
		//	FailedAt:  time.Now(),
		// }
		// // s.eventBus.Publish(ctx, event)
		
		return nil, domain.ErrUserNotFound
	}
	
	// Se TenantID não foi fornecido, usar o do usuário
	if req.TenantID == "" {
		req.TenantID = user.TenantID
	} else if user.TenantID != req.TenantID {
		// Se fornecido, verificar se corresponde ao usuário
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		// TODO: Publish login failed event
		
		return nil, domain.ErrUserNotFound
	}
	
	// Verificar se o usuário pode fazer login
	if !user.CanLogin() {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		// TODO: Publish event
		
		return nil, fmt.Errorf("usuário inativo")
	}
	
	// Verificar senha
	if !user.VerifyPassword(req.Password) {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		// TODO: Publish event
		
		return nil, fmt.Errorf("credenciais inválidas")
	}
	
	// Gerar tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}
	
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}
	
	// Criar sessão
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.jwtExpiryHours) * time.Hour)
	refreshExpiresAt := now.Add(time.Duration(s.refreshExpiryDays) * 24 * time.Hour)
	
	session := &domain.Session{
		ID:               uuid.New().String(),
		UserID:           user.ID,
		TenantID:         user.TenantID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		IPAddress:        req.IPAddress,
		UserAgent:        req.UserAgent,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("erro ao criar sessão: %w", err)
	}
	
	// Salvar refresh token
	refreshTokenEntity := &domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Token:     refreshToken,
		ExpiresAt: refreshExpiresAt,
		IsUsed:    false,
		CreatedAt: now,
	}
	
	if err := s.refreshTokenRepo.Create(refreshTokenEntity); err != nil {
		return nil, fmt.Errorf("erro ao salvar refresh token: %w", err)
	}
	
	// Atualizar último login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Não falhar o login por erro ao atualizar último login
		// Apenas log do erro
	}
	
	// Registrar tentativa de login bem-sucedida
	s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, true)
	
	// TODO: Publish event
	
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: &UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Status:    user.Status,
			TenantID:  user.TenantID,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// RefreshTokenRequest representa uma solicitação de refresh de token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TenantData representa dados do tenant para registro
type TenantData struct {
	Name     string `json:"name" validate:"required"`
	Document string `json:"document" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Website  string `json:"website"`
	Plan     string `json:"plan" validate:"required"`
	Address  struct {
		Street       string `json:"street" validate:"required"`
		Number       string `json:"number" validate:"required"`
		Complement   string `json:"complement"`
		Neighborhood string `json:"neighborhood" validate:"required"`
		City         string `json:"city" validate:"required"`
		State        string `json:"state" validate:"required"`
		ZipCode      string `json:"zipCode" validate:"required"`
	} `json:"address"`
}

// UserData representa dados do usuário para registro
type UserData struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

// TenantDTO representa um tenant para transferência de dados
type TenantDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Website   string    `json:"website"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest representa uma solicitação de registro
type RegisterRequest struct {
	Tenant TenantData `json:"tenant" validate:"required"`
	User   UserData   `json:"user" validate:"required"`
}

// RegisterResponse representa resposta de registro bem-sucedido
type RegisterResponse struct {
	Tenant TenantDTO `json:"tenant"`
	User   UserDTO   `json:"user"`
	Message string   `json:"message"`
}

// ForgotPasswordRequest representa solicitação de recuperação de senha
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest representa solicitação de reset de senha
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RefreshToken renova os tokens de acesso usando um refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error) {
	// Buscar refresh token
	refreshToken, err := s.refreshTokenRepo.GetByToken(req.RefreshToken)
	if err != nil {
		return nil, domain.ErrRefreshTokenNotFound
	}
	
	// Verificar se pode ser usado
	if !refreshToken.CanUse() {
		if refreshToken.IsUsed {
			return nil, domain.ErrTokenAlreadyUsed
		}
		return nil, domain.ErrTokenExpired
	}
	
	// Buscar usuário
	user, err := s.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	
	// Verificar se usuário ainda pode fazer login
	if !user.CanLogin() {
		return nil, fmt.Errorf("usuário inativo")
	}
	
	// Marcar refresh token como usado
	if err := s.refreshTokenRepo.MarkAsUsed(refreshToken.ID); err != nil {
		return nil, fmt.Errorf("erro ao marcar refresh token como usado: %w", err)
	}
	
	// Gerar novos tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}
	
	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}
	
	// Criar nova sessão
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.jwtExpiryHours) * time.Hour)
	refreshExpiresAt := now.Add(time.Duration(s.refreshExpiryDays) * 24 * time.Hour)
	
	session := &domain.Session{
		ID:               uuid.New().String(),
		UserID:           user.ID,
		TenantID:         user.TenantID,
		AccessToken:      accessToken,
		RefreshToken:     newRefreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("erro ao criar sessão: %w", err)
	}
	
	// Salvar novo refresh token
	newRefreshTokenEntity := &domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Token:     newRefreshToken,
		ExpiresAt: refreshExpiresAt,
		IsUsed:    false,
		CreatedAt: now,
	}
	
	if err := s.refreshTokenRepo.Create(newRefreshTokenEntity); err != nil {
		return nil, fmt.Errorf("erro ao salvar refresh token: %w", err)
	}
	
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User: &UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Status:    user.Status,
			TenantID:  user.TenantID,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// Logout invalida uma sessão
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	session, err := s.sessionRepo.GetByAccessToken(accessToken)
	if err != nil {
		return domain.ErrSessionNotFound
	}
	
	// Invalidar sessão
	session.Invalidate()
	if err := s.sessionRepo.Update(session); err != nil {
		return fmt.Errorf("erro ao invalidar sessão: %w", err)
	}
	
	// Invalidar refresh tokens do usuário
	if err := s.refreshTokenRepo.DeleteByUserID(session.UserID); err != nil {
		// Log do erro, mas não falhar o logout
	}
	
	// TODO: Publish event
	
	return nil
}

// ValidateToken valida um token de acesso
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*UserDTO, error) {
	// Parse do token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	
	if !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	
	// Extrair claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	// Buscar usuário
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	
	// Verificar se usuário ainda está ativo
	if !user.IsActive() {
		return nil, fmt.Errorf("usuário inativo")
	}
	
	return &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}

// generateAccessToken gera um JWT token
func (s *AuthService) generateAccessToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"email":     user.Email,
		"role":      user.Role,
		"iat":       now.Unix(),
		"exp":       now.Add(time.Duration(s.jwtExpiryHours) * time.Hour).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken gera um refresh token aleatório
func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// recordLoginAttempt registra uma tentativa de login
func (s *AuthService) recordLoginAttempt(email, tenantID, ipAddress, userAgent string, success bool) {
	attempt := &domain.LoginAttempt{
		ID:        uuid.New().String(),
		Email:     email,
		TenantID:  tenantID,
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	
	// Não falhar por erro ao registrar tentativa
	s.loginAttemptRepo.Create(attempt)
}

// Register cria novo tenant e usuário administrador
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// Verificar se email já existe
	existingUser, _ := s.userRepo.GetByEmail(req.User.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email já existe")
	}
	
	// Por enquanto, criar um tenant ID simples
	// Em produção, isso deveria integrar com o tenant-service
	tenantID := uuid.New().String()
	
	// Criar usuário administrador
	now := time.Now()
	user := &domain.User{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Email:     req.User.Email,
		FirstName: req.User.Name, // Usando o nome completo como firstName por simplicidade
		LastName:  "",
		Role:      domain.RoleAdmin,
		Status:    domain.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Validar dados do usuário
	if err := user.ValidateEmail(); err != nil {
		return nil, fmt.Errorf("dados inválidos")
	}
	
	// Definir senha
	if err := user.SetPassword(req.User.Password); err != nil {
		return nil, fmt.Errorf("dados inválidos")
	}
	
	// Salvar usuário
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}
	
	// TODO: Integrar com tenant-service para criar tenant real
	// Por enquanto, retornar dados mock do tenant
	tenantDTO := TenantDTO{
		ID:        tenantID,
		Name:      req.Tenant.Name,
		Document:  req.Tenant.Document,
		Email:     req.Tenant.Email,
		Phone:     req.Tenant.Phone,
		Website:   req.Tenant.Website,
		Plan:      req.Tenant.Plan,
		Status:    "active",
		CreatedAt: now,
	}
	
	userDTO := UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		TenantID:  user.TenantID,
		CreatedAt: user.CreatedAt,
	}
	
	return &RegisterResponse{
		Tenant:  tenantDTO,
		User:    userDTO,
		Message: "Conta criada com sucesso! Você pode fazer login agora.",
	}, nil
}

// ForgotPassword cria token de recuperação e envia email
func (s *AuthService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Buscar usuário por email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return domain.ErrUserNotFound
	}
	
	// Verificar se usuário pode recuperar senha
	if !user.IsActive() {
		return fmt.Errorf("usuário inativo")
	}
	
	// Gerar token de recuperação
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("erro ao gerar token: %w", err)
	}
	tokenString := hex.EncodeToString(tokenBytes)
	
	// Criar registro do token (válido por 1 hora)
	now := time.Now()
	resetToken := &domain.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Token:     tokenString,
		Email:     user.Email,
		ExpiresAt: now.Add(1 * time.Hour),
		IsUsed:    false,
		CreatedAt: now,
	}
	
	// Salvar token
	if err := s.passwordResetTokenRepo.Create(resetToken); err != nil {
		return fmt.Errorf("erro ao salvar token de recuperação: %w", err)
	}
	
	// TODO: Enviar email com link de recuperação
	// Por enquanto, apenas log que seria enviado
	// Em produção, integrar com notification-service
	
	// TODO: Publish event para envio de email
	// event := domain.PasswordResetRequestedEvent{
	//     UserID:    user.ID,
	//     Email:     user.Email,
	//     Token:     tokenString,
	//     ExpiresAt: resetToken.ExpiresAt,
	// }
	// s.eventBus.Publish(ctx, event)
	
	return nil
}

// ResetPassword valida token e redefine senha
func (s *AuthService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// Buscar token de recuperação
	resetToken, err := s.passwordResetTokenRepo.GetByToken(req.Token)
	if err != nil {
		return fmt.Errorf("token inválido")
	}
	
	// Verificar se token pode ser usado
	if !resetToken.CanUse() {
		if resetToken.IsUsed {
			return fmt.Errorf("token já utilizado")
		}
		return fmt.Errorf("token expirado")
	}
	
	// Buscar usuário
	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	
	// Verificar se usuário ainda está ativo
	if !user.IsActive() {
		return fmt.Errorf("usuário inativo")
	}
	
	// Definir nova senha
	if err := user.SetPassword(req.Password); err != nil {
		return fmt.Errorf("senha inválida: %w", err)
	}
	
	// Atualizar usuário
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("erro ao atualizar senha: %w", err)
	}
	
	// Marcar token como usado
	if err := s.passwordResetTokenRepo.MarkAsUsed(resetToken.ID); err != nil {
		// Log do erro, mas não falhar o reset
	}
	
	// Invalidar todas as sessões do usuário
	if err := s.sessionRepo.DeleteByUserID(user.ID); err != nil {
		// Log do erro, mas não falhar o reset
	}
	
	// Invalidar todos os refresh tokens do usuário
	if err := s.refreshTokenRepo.DeleteByUserID(user.ID); err != nil {
		// Log do erro, mas não falhar o reset
	}
	
	// TODO: Publish event de senha alterada
	// event := domain.PasswordResetCompletedEvent{
	//     UserID:   user.ID,
	//     Email:    user.Email,
	//     ResetAt:  time.Now(),
	// }
	// s.eventBus.Publish(ctx, event)
	
	return nil
}