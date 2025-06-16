package domain

import (
	"time"
	"errors"
	"regexp"
	"strings"
	"golang.org/x/crypto/bcrypt"
)

// User representa um usuário do sistema
type User struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password_hash"` // Não expor senha no JSON
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Role        UserRole  `json:"role" db:"role"`
	Status      UserStatus `json:"status" db:"status"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole define os papéis possíveis de um usuário
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleManager   UserRole = "manager"
	RoleOperator  UserRole = "operator"
	RoleClient    UserRole = "client"
	RoleReadOnly  UserRole = "readonly"
)

// UserStatus define os status possíveis de um usuário
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusPending   UserStatus = "pending"
	StatusSuspended UserStatus = "suspended"
	StatusBlocked   UserStatus = "blocked"
)

// UserRepository define a interface para persistência de usuários
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByTenant(tenantID string, limit, offset int) ([]*User, error)
	Update(user *User) error
	Delete(id string) error
	UpdateLastLogin(id string) error
}

// LoginAttempt representa uma tentativa de login
type LoginAttempt struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Success   bool      `json:"success" db:"success"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// LoginAttemptRepository define interface para tentativas de login
type LoginAttemptRepository interface {
	Create(attempt *LoginAttempt) error
	GetRecentFailures(email string, minutes int) (int, error)
}

// Validação de domínio

var (
	ErrInvalidEmail    = errors.New("email inválido")
	ErrWeakPassword    = errors.New("senha muito fraca")
	ErrInvalidRole     = errors.New("papel de usuário inválido")
	ErrInvalidStatus   = errors.New("status de usuário inválido")
	ErrUserNotFound    = errors.New("usuário não encontrado")
	ErrEmailExists     = errors.New("email já está em uso")
	ErrTooManyFailures = errors.New("muitas tentativas de login falharam")
)

// ValidateEmail valida se o email está em formato válido
func (u *User) ValidateEmail() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	
	return nil
}

// ValidatePassword valida se a senha atende aos critérios mínimos
func (u *User) ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	
	// Verificar se contém ao menos uma letra maiúscula
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return ErrWeakPassword
	}
	
	// Verificar se contém ao menos uma letra minúscula
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return ErrWeakPassword
	}
	
	// Verificar se contém ao menos um número
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return ErrWeakPassword
	}
	
	// Verificar se contém ao menos um símbolo
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		return ErrWeakPassword
	}
	
	return nil
}

// ValidateRole valida se o papel é válido
func (u *User) ValidateRole() error {
	validRoles := []UserRole{RoleAdmin, RoleManager, RoleOperator, RoleClient, RoleReadOnly}
	for _, role := range validRoles {
		if u.Role == role {
			return nil
		}
	}
	return ErrInvalidRole
}

// ValidateStatus valida se o status é válido
func (u *User) ValidateStatus() error {
	validStatuses := []UserStatus{StatusActive, StatusInactive, StatusPending, StatusSuspended, StatusBlocked}
	for _, status := range validStatuses {
		if u.Status == status {
			return nil
		}
	}
	return ErrInvalidStatus
}

// SetPassword gera hash da senha usando bcrypt
func (u *User) SetPassword(password string) error {
	if err := u.ValidatePassword(password); err != nil {
		return err
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword verifica se a senha fornecida confere com o hash armazenado
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// FullName retorna o nome completo do usuário
func (u *User) FullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// IsActive verifica se o usuário está ativo
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// CanLogin verifica se o usuário pode fazer login
func (u *User) CanLogin() bool {
	return u.Status == StatusActive || u.Status == StatusPending
}

// HasRole verifica se o usuário tem o papel especificado
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// HasPermission verifica se o usuário tem permissão baseada na hierarquia de papéis
func (u *User) HasPermission(requiredRole UserRole) bool {
	roleHierarchy := map[UserRole]int{
		RoleReadOnly: 1,
		RoleClient:   2,
		RoleOperator: 3,
		RoleManager:  4,
		RoleAdmin:    5,
	}
	
	userLevel := roleHierarchy[u.Role]
	requiredLevel := roleHierarchy[requiredRole]
	
	return userLevel >= requiredLevel
}