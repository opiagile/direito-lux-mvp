package database

import (
	"fmt"
	"time"
	"github.com/jmoiron/sqlx"
	
	"github.com/direito-lux/auth-service/internal/domain"
)

// LoginAttemptRepository implementa a interface domain.LoginAttemptRepository
type LoginAttemptRepository struct {
	db *sqlx.DB
}

// NewLoginAttemptRepository cria uma nova instância do repositório de tentativas de login
func NewLoginAttemptRepository(db *sqlx.DB) *LoginAttemptRepository {
	return &LoginAttemptRepository{db: db}
}

// Create cria uma nova tentativa de login
func (r *LoginAttemptRepository) Create(attempt *domain.LoginAttempt) error {
	query := `
		INSERT INTO login_attempts (
			id, email, tenant_id, success, ip_address, 
			user_agent, created_at
		) VALUES (
			:id, :email, :tenant_id, :success, :ip_address,
			:user_agent, :created_at
		)`
	
	_, err := r.db.NamedExec(query, attempt)
	if err != nil {
		return fmt.Errorf("erro ao criar tentativa de login: %w", err)
	}
	
	return nil
}

// GetRecentFailures busca tentativas de login falhadas recentes
func (r *LoginAttemptRepository) GetRecentFailures(email string, minutes int) (int, error) {
	var count int
	
	query := `
		SELECT COUNT(*) 
		FROM login_attempts 
		WHERE email = $1 
		  AND success = false 
		  AND created_at > $2`
	
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	
	err := r.db.Get(&count, query, email, since)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar tentativas de login recentes: %w", err)
	}
	
	return count, nil
}