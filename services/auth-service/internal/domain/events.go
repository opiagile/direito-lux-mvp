package domain

import (
	"time"
	"encoding/json"
)

// Eventos de domínio para autenticação

// UserCreatedEvent é disparado quando um usuário é criado
type UserCreatedEvent struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (e UserCreatedEvent) EventType() string {
	return "user.created"
}

func (e UserCreatedEvent) AggregateID() string {
	return e.UserID
}

func (e UserCreatedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// UserUpdatedEvent é disparado quando um usuário é atualizado
type UserUpdatedEvent struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      UserRole  `json:"role"`
	Status    UserStatus `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e UserUpdatedEvent) EventType() string {
	return "user.updated"
}

func (e UserUpdatedEvent) AggregateID() string {
	return e.UserID
}

func (e UserUpdatedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// UserDeletedEvent é disparado quando um usuário é removido
type UserDeletedEvent struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (e UserDeletedEvent) EventType() string {
	return "user.deleted"
}

func (e UserDeletedEvent) AggregateID() string {
	return e.UserID
}

func (e UserDeletedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// UserLoggedInEvent é disparado quando um usuário faz login
type UserLoggedInEvent struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	LoggedInAt time.Time `json:"logged_in_at"`
}

func (e UserLoggedInEvent) EventType() string {
	return "user.logged_in"
}

func (e UserLoggedInEvent) AggregateID() string {
	return e.UserID
}

func (e UserLoggedInEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// UserLoggedOutEvent é disparado quando um usuário faz logout
type UserLoggedOutEvent struct {
	UserID     string    `json:"user_id"`
	TenantID   string    `json:"tenant_id"`
	SessionID  string    `json:"session_id"`
	LoggedOutAt time.Time `json:"logged_out_at"`
}

func (e UserLoggedOutEvent) EventType() string {
	return "user.logged_out"
}

func (e UserLoggedOutEvent) AggregateID() string {
	return e.UserID
}

func (e UserLoggedOutEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// LoginFailedEvent é disparado quando uma tentativa de login falha
type LoginFailedEvent struct {
	Email     string    `json:"email"`
	TenantID  string    `json:"tenant_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

func (e LoginFailedEvent) EventType() string {
	return "login.failed"
}

func (e LoginFailedEvent) AggregateID() string {
	return e.Email
}

func (e LoginFailedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// UserStatusChangedEvent é disparado quando o status de um usuário muda
type UserStatusChangedEvent struct {
	UserID    string     `json:"user_id"`
	TenantID  string     `json:"tenant_id"`
	Email     string     `json:"email"`
	OldStatus UserStatus `json:"old_status"`
	NewStatus UserStatus `json:"new_status"`
	ChangedAt time.Time  `json:"changed_at"`
	ChangedBy string     `json:"changed_by"`
}

func (e UserStatusChangedEvent) EventType() string {
	return "user.status_changed"
}

func (e UserStatusChangedEvent) AggregateID() string {
	return e.UserID
}

func (e UserStatusChangedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// UserRoleChangedEvent é disparado quando o papel de um usuário muda
type UserRoleChangedEvent struct {
	UserID    string   `json:"user_id"`
	TenantID  string   `json:"tenant_id"`
	Email     string   `json:"email"`
	OldRole   UserRole `json:"old_role"`
	NewRole   UserRole `json:"new_role"`
	ChangedAt time.Time `json:"changed_at"`
	ChangedBy string   `json:"changed_by"`
}

func (e UserRoleChangedEvent) EventType() string {
	return "user.role_changed"
}

func (e UserRoleChangedEvent) AggregateID() string {
	return e.UserID
}

func (e UserRoleChangedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}

// PasswordChangedEvent é disparado quando um usuário muda a senha
type PasswordChangedEvent struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	ChangedAt time.Time `json:"changed_at"`
	IPAddress string    `json:"ip_address"`
}

func (e PasswordChangedEvent) EventType() string {
	return "password.changed"
}

func (e PasswordChangedEvent) AggregateID() string {
	return e.UserID
}

func (e PasswordChangedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}