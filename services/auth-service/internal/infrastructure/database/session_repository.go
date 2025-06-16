package database

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/jmoiron/sqlx"
	
	"github.com/direito-lux/auth-service/internal/domain"
)

// SessionRepository implementa a interface domain.SessionRepository
type SessionRepository struct {
	db *sqlx.DB
}

// NewSessionRepository cria uma nova instância do repositório de sessões
func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create cria uma nova sessão
func (r *SessionRepository) Create(session *domain.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, tenant_id, access_token, refresh_token,
			expires_at, refresh_expires_at, ip_address, user_agent,
			is_active, created_at, updated_at
		) VALUES (
			:id, :user_id, :tenant_id, :access_token, :refresh_token,
			:expires_at, :refresh_expires_at, :ip_address, :user_agent,
			:is_active, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("erro ao criar sessão: %w", err)
	}
	
	return nil
}

// GetByAccessToken busca uma sessão pelo access token
func (r *SessionRepository) GetByAccessToken(token string) (*domain.Session, error) {
	session := &domain.Session{}
	
	query := `
		SELECT id, user_id, tenant_id, access_token, refresh_token,
			   expires_at, refresh_expires_at, ip_address, user_agent,
			   is_active, created_at, updated_at
		FROM sessions 
		WHERE access_token = $1 AND is_active = true`
	
	err := r.db.Get(session, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("erro ao buscar sessão por access token: %w", err)
	}
	
	return session, nil
}

// GetByRefreshToken busca uma sessão pelo refresh token
func (r *SessionRepository) GetByRefreshToken(token string) (*domain.Session, error) {
	session := &domain.Session{}
	
	query := `
		SELECT id, user_id, tenant_id, access_token, refresh_token,
			   expires_at, refresh_expires_at, ip_address, user_agent,
			   is_active, created_at, updated_at
		FROM sessions 
		WHERE refresh_token = $1 AND is_active = true`
	
	err := r.db.Get(session, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("erro ao buscar sessão por refresh token: %w", err)
	}
	
	return session, nil
}

// GetByUserID busca todas as sessões de um usuário
func (r *SessionRepository) GetByUserID(userID string) ([]*domain.Session, error) {
	sessions := []*domain.Session{}
	
	query := `
		SELECT id, user_id, tenant_id, access_token, refresh_token,
			   expires_at, refresh_expires_at, ip_address, user_agent,
			   is_active, created_at, updated_at
		FROM sessions 
		WHERE user_id = $1
		ORDER BY created_at DESC`
	
	err := r.db.Select(&sessions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar sessões por usuário: %w", err)
	}
	
	return sessions, nil
}

// Update atualiza uma sessão
func (r *SessionRepository) Update(session *domain.Session) error {
	query := `
		UPDATE sessions SET
			access_token = :access_token,
			refresh_token = :refresh_token,
			expires_at = :expires_at,
			refresh_expires_at = :refresh_expires_at,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id`
	
	result, err := r.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("erro ao atualizar sessão: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrSessionNotFound
	}
	
	return nil
}

// Delete remove uma sessão
func (r *SessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar sessão: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrSessionNotFound
	}
	
	return nil
}

// DeleteByUserID remove todas as sessões de um usuário
func (r *SessionRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("erro ao deletar sessões do usuário: %w", err)
	}
	
	return nil
}

// DeleteExpired remove sessões expiradas
func (r *SessionRepository) DeleteExpired() error {
	query := `
		DELETE FROM sessions 
		WHERE expires_at < $1 OR refresh_expires_at < $1`
	
	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("erro ao deletar sessões expiradas: %w", err)
	}
	
	return nil
}