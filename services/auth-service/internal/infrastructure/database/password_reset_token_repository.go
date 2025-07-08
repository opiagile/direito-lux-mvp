package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/direito-lux/auth-service/internal/domain"
)

// passwordResetTokenRepository implementa domain.PasswordResetTokenRepository
type passwordResetTokenRepository struct {
	db *sqlx.DB
}

// NewPasswordResetTokenRepository cria nova instância do repositório
func NewPasswordResetTokenRepository(db *sqlx.DB) domain.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{
		db: db,
	}
}

// Create salva um novo token de recuperação de senha
func (r *passwordResetTokenRepository) Create(token *domain.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (
			id, user_id, tenant_id, token, email, expires_at, 
			is_used, created_at, used_at
		) VALUES (
			:id, :user_id, :tenant_id, :token, :email, :expires_at,
			:is_used, :created_at, :used_at
		)`
	
	_, err := r.db.NamedExec(query, token)
	return err
}

// GetByToken busca token de recuperação pelo token string
func (r *passwordResetTokenRepository) GetByToken(tokenStr string) (*domain.PasswordResetToken, error) {
	var token domain.PasswordResetToken
	
	query := `
		SELECT id, user_id, tenant_id, token, email, expires_at,
			   is_used, created_at, used_at
		FROM password_reset_tokens 
		WHERE token = $1`
	
	err := r.db.Get(&token, query, tokenStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPasswordResetTokenNotFound
		}
		return nil, err
	}
	
	return &token, nil
}

// MarkAsUsed marca token como usado
func (r *passwordResetTokenRepository) MarkAsUsed(id string) error {
	now := time.Now()
	
	query := `
		UPDATE password_reset_tokens 
		SET is_used = true, used_at = $1 
		WHERE id = $2`
	
	_, err := r.db.Exec(query, now, id)
	return err
}

// DeleteExpired remove tokens expirados
func (r *passwordResetTokenRepository) DeleteExpired() error {
	query := `
		DELETE FROM password_reset_tokens 
		WHERE expires_at < NOW()`
	
	_, err := r.db.Exec(query)
	return err
}

// DeleteByUserID remove todos os tokens de um usuário
func (r *passwordResetTokenRepository) DeleteByUserID(userID string) error {
	query := `
		DELETE FROM password_reset_tokens 
		WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	return err
}