package domain

import (
	"time"
	"errors"
)

// Session representa uma sessão de usuário
type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" db:"refresh_expires_at"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SessionRepository define interface para persistência de sessões
type SessionRepository interface {
	Create(session *Session) error
	GetByAccessToken(token string) (*Session, error)
	GetByRefreshToken(token string) (*Session, error)
	GetByUserID(userID string) ([]*Session, error)
	Update(session *Session) error
	Delete(id string) error
	DeleteByUserID(userID string) error
	DeleteExpired() error
}

// RefreshToken representa um token de refresh
type RefreshToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
}

// RefreshTokenRepository define interface para tokens de refresh
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	MarkAsUsed(id string) error
	DeleteExpired() error
	DeleteByUserID(userID string) error
}

// Erros de domínio para sessão
var (
	ErrSessionNotFound    = errors.New("sessão não encontrada")
	ErrSessionExpired     = errors.New("sessão expirada")
	ErrInvalidToken       = errors.New("token inválido")
	ErrTokenExpired       = errors.New("token expirado")
	ErrTokenAlreadyUsed   = errors.New("token já foi utilizado")
	ErrRefreshTokenNotFound = errors.New("refresh token não encontrado")
)

// IsExpired verifica se a sessão expirou
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRefreshExpired verifica se o refresh token expirou
func (s *Session) IsRefreshExpired() bool {
	return time.Now().After(s.RefreshExpiresAt)
}

// CanRefresh verifica se a sessão pode ser renovada
func (s *Session) CanRefresh() bool {
	return s.IsActive && !s.IsRefreshExpired()
}

// Invalidate marca a sessão como inativa
func (s *Session) Invalidate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// IsExpired verifica se o refresh token expirou
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// CanUse verifica se o refresh token pode ser usado
func (rt *RefreshToken) CanUse() bool {
	return !rt.IsUsed && !rt.IsExpired()
}

// MarkUsed marca o refresh token como usado
func (rt *RefreshToken) MarkUsed() {
	rt.IsUsed = true
	now := time.Now()
	rt.UsedAt = &now
}