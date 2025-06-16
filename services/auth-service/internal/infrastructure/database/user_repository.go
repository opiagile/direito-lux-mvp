package database

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/jmoiron/sqlx"
	
	"github.com/direito-lux/auth-service/internal/domain"
)

// UserRepository implementa a interface domain.UserRepository
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository cria uma nova instância do repositório de usuários
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create cria um novo usuário
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (
			id, tenant_id, email, password_hash, first_name, 
			last_name, role, status, created_at, updated_at
		) VALUES (
			:id, :tenant_id, :email, :password_hash, :first_name,
			:last_name, :role, :status, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}
	
	return nil
}

// GetByID busca um usuário por ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	user := &domain.User{}
	
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, 
			   last_name, role, status, last_login_at, 
			   created_at, updated_at
		FROM users 
		WHERE id = $1`
	
	err := r.db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("erro ao buscar usuário por ID: %w", err)
	}
	
	return user, nil
}

// GetByEmail busca um usuário por email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, 
			   last_name, role, status, last_login_at, 
			   created_at, updated_at
		FROM users 
		WHERE email = $1`
	
	err := r.db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("erro ao buscar usuário por email: %w", err)
	}
	
	return user, nil
}

// GetByTenant busca usuários de um tenant com paginação
func (r *UserRepository) GetByTenant(tenantID string, limit, offset int) ([]*domain.User, error) {
	users := []*domain.User{}
	
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, 
			   last_name, role, status, last_login_at, 
			   created_at, updated_at
		FROM users 
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	err := r.db.Select(&users, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários por tenant: %w", err)
	}
	
	return users, nil
}

// Update atualiza um usuário
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users SET
			email = :email,
			password_hash = :password_hash,
			first_name = :first_name,
			last_name = :last_name,
			role = :role,
			status = :status,
			updated_at = :updated_at
		WHERE id = :id`
	
	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	
	return nil
}

// Delete remove um usuário
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	
	return nil
}

// UpdateLastLogin atualiza o último login do usuário
func (r *UserRepository) UpdateLastLogin(id string) error {
	query := `
		UPDATE users 
		SET last_login_at = $1, updated_at = $2
		WHERE id = $3`
	
	now := time.Now()
	result, err := r.db.Exec(query, now, now, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar último login: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	
	return nil
}