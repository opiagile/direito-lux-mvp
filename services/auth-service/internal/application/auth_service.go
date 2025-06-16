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
)

// AuthService implementa os casos de uso de autenticação
type AuthService struct {
	userRepo         domain.UserRepository
	sessionRepo      domain.SessionRepository
	refreshTokenRepo domain.RefreshTokenRepository
	loginAttemptRepo domain.LoginAttemptRepository
	eventBus         EventBus
	jwtSecret        string
	jwtExpiryHours   int
	refreshExpiryDays int
}

// NewAuthService cria uma nova instância do serviço de autenticação
func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	loginAttemptRepo domain.LoginAttemptRepository,
	eventBus EventBus,
	jwtSecret string,
	jwtExpiryHours int,
	refreshExpiryDays int,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		refreshTokenRepo:  refreshTokenRepo,
		loginAttemptRepo:  loginAttemptRepo,
		eventBus:          eventBus,
		jwtSecret:         jwtSecret,
		jwtExpiryHours:    jwtExpiryHours,
		refreshExpiryDays: refreshExpiryDays,
	}
}

// LoginRequest representa uma solicitação de login
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	TenantID  string `json:"tenant_id" validate:"required"`
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
		
		// Publicar evento de login bloqueado
		event := domain.LoginFailedEvent{
			Email:     req.Email,
			TenantID:  req.TenantID,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
			Reason:    "too_many_failures",
			FailedAt:  time.Now(),
		}
		s.eventBus.Publish(ctx, event)
		
		return nil, domain.ErrTooManyFailures
	}
	
	// Buscar usuário por email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		event := domain.LoginFailedEvent{
			Email:     req.Email,
			TenantID:  req.TenantID,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
			Reason:    "user_not_found",
			FailedAt:  time.Now(),
		}
		s.eventBus.Publish(ctx, event)
		
		return nil, domain.ErrUserNotFound
	}
	
	// Verificar se o usuário pertence ao tenant correto
	if user.TenantID != req.TenantID {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		event := domain.LoginFailedEvent{
			Email:     req.Email,
			TenantID:  req.TenantID,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
			Reason:    "wrong_tenant",
			FailedAt:  time.Now(),
		}
		s.eventBus.Publish(ctx, event)
		
		return nil, domain.ErrUserNotFound
	}
	
	// Verificar se o usuário pode fazer login
	if !user.CanLogin() {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		event := domain.LoginFailedEvent{
			Email:     req.Email,
			TenantID:  req.TenantID,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
			Reason:    "user_inactive",
			FailedAt:  time.Now(),
		}
		s.eventBus.Publish(ctx, event)
		
		return nil, fmt.Errorf("usuário inativo")
	}
	
	// Verificar senha
	if !user.VerifyPassword(req.Password) {
		s.recordLoginAttempt(req.Email, req.TenantID, req.IPAddress, req.UserAgent, false)
		
		event := domain.LoginFailedEvent{
			Email:     req.Email,
			TenantID:  req.TenantID,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
			Reason:    "invalid_password",
			FailedAt:  time.Now(),
		}
		s.eventBus.Publish(ctx, event)
		
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
	
	// Publicar evento de login
	loginEvent := domain.UserLoggedInEvent{
		UserID:     user.ID,
		TenantID:   user.TenantID,
		Email:      user.Email,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
		LoggedInAt: now,
	}
	s.eventBus.Publish(ctx, loginEvent)
	
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
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// RefreshTokenRequest representa uma solicitação de refresh de token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
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
	
	// Publicar evento de logout
	event := domain.UserLoggedOutEvent{
		UserID:      session.UserID,
		TenantID:    session.TenantID,
		SessionID:   session.ID,
		LoggedOutAt: time.Now(),
	}
	s.eventBus.Publish(ctx, event)
	
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