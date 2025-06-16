package database

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/jmoiron/sqlx"
	
	"github.com/direito-lux/auth-service/internal/domain"
)

// RefreshTokenRepository implementa a interface domain.RefreshTokenRepository
type RefreshTokenRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenRepository cria uma nova instância do repositório de refresh tokens
func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create cria um novo refresh token
func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			id, user_id, tenant_id, token, expires_at, 
			is_used, created_at, used_at
		) VALUES (
			:id, :user_id, :tenant_id, :token, :expires_at,
			:is_used, :created_at, :used_at
		)`
	
	_, err := r.db.NamedExec(query, token)
	if err != nil {
		return fmt.Errorf("erro ao criar refresh token: %w", err)
	}
	
	return nil
}

// GetByToken busca um refresh token pelo valor do token
func (r *RefreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
	refreshToken := &domain.RefreshToken{}
	
	query := `
		SELECT id, user_id, tenant_id, token, expires_at,
			   is_used, created_at, used_at
		FROM refresh_tokens 
		WHERE token = $1`
	
	err := r.db.Get(refreshToken, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("erro ao buscar refresh token: %w", err)
	}
	
	return refreshToken, nil
}

// MarkAsUsed marca um refresh token como usado
func (r *RefreshTokenRepository) MarkAsUsed(id string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_used = true, used_at = $1
		WHERE id = $2`
	
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao marcar refresh token como usado: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}
	
	return nil
}

// DeleteExpired remove refresh tokens expirados
func (r *RefreshTokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	
	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("erro ao deletar refresh tokens expirados: %w", err)
	}
	
	return nil
}

// DeleteByUserID remove todos os refresh tokens de um usuário
func (r *RefreshTokenRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("erro ao deletar refresh tokens do usuário: %w", err)
	}
	
	return nil
}