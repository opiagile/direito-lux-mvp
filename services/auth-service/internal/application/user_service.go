package application

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	
	"github.com/direito-lux/auth-service/internal/domain"
)

// UserService implementa os casos de uso de gerenciamento de usuários
type UserService struct {
	userRepo domain.UserRepository
	eventBus EventBus
}

// NewUserService cria uma nova instância do serviço de usuários
func NewUserService(userRepo domain.UserRepository, eventBus EventBus) *UserService {
	return &UserService{
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

// CreateUserRequest representa uma solicitação de criação de usuário
type CreateUserRequest struct {
	TenantID  string            `json:"tenant_id" validate:"required"`
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password" validate:"required,min=8"`
	FirstName string            `json:"first_name" validate:"required"`
	LastName  string            `json:"last_name" validate:"required"`
	Role      domain.UserRole   `json:"role" validate:"required"`
	Status    domain.UserStatus `json:"status"`
}

// UpdateUserRequest representa uma solicitação de atualização de usuário
type UpdateUserRequest struct {
	FirstName *string            `json:"first_name,omitempty"`
	LastName  *string            `json:"last_name,omitempty"`
	Role      *domain.UserRole   `json:"role,omitempty"`
	Status    *domain.UserStatus `json:"status,omitempty"`
}

// ChangePasswordRequest representa uma solicitação de mudança de senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// CreateUser cria um novo usuário
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	// Verificar se email já existe
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, domain.ErrEmailExists
	}
	
	// Criar usuário
	user := &domain.User{
		ID:        uuid.New().String(),
		TenantID:  req.TenantID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Se status não foi especificado, usar pending
	if user.Status == "" {
		user.Status = domain.StatusPending
	}
	
	// Validar dados
	if err := user.ValidateEmail(); err != nil {
		return nil, err
	}
	
	if err := user.ValidateRole(); err != nil {
		return nil, err
	}
	
	if err := user.ValidateStatus(); err != nil {
		return nil, err
	}
	
	// Definir senha
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}
	
	// Salvar usuário
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}
	
	// Publicar evento
	event := domain.UserCreatedEvent{
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	s.eventBus.Publish(ctx, event)
	
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

// GetUser busca um usuário por ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*UserDTO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
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

// GetUserByEmail busca um usuário por email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, domain.ErrUserNotFound
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

// GetUsersByTenant busca usuários de um tenant
func (s *UserService) GetUsersByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*UserDTO, error) {
	users, err := s.userRepo.GetByTenant(tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}
	
	userDTOs := make([]*UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = &UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		}
	}
	
	return userDTOs, nil
}

// UpdateUser atualiza um usuário
func (s *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*UserDTO, error) {
	// Buscar usuário atual
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	
	// Guardar valores antigos para eventos
	oldRole := user.Role
	oldStatus := user.Status
	
	// Atualizar campos
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	
	if req.Role != nil {
		user.Role = *req.Role
		if err := user.ValidateRole(); err != nil {
			return nil, err
		}
	}
	
	if req.Status != nil {
		user.Status = *req.Status
		if err := user.ValidateStatus(); err != nil {
			return nil, err
		}
	}
	
	user.UpdatedAt = time.Now()
	
	// Salvar usuário
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}
	
	// Publicar evento de atualização
	updateEvent := domain.UserUpdatedEvent{
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		UpdatedAt: user.UpdatedAt,
	}
	s.eventBus.Publish(ctx, updateEvent)
	
	// Publicar eventos específicos se houve mudanças
	if req.Role != nil && oldRole != *req.Role {
		roleEvent := domain.UserRoleChangedEvent{
			UserID:    user.ID,
			TenantID:  user.TenantID,
			Email:     user.Email,
			OldRole:   oldRole,
			NewRole:   *req.Role,
			ChangedAt: time.Now(),
		}
		s.eventBus.Publish(ctx, roleEvent)
	}
	
	if req.Status != nil && oldStatus != *req.Status {
		statusEvent := domain.UserStatusChangedEvent{
			UserID:    user.ID,
			TenantID:  user.TenantID,
			Email:     user.Email,
			OldStatus: oldStatus,
			NewStatus: *req.Status,
			ChangedAt: time.Now(),
		}
		s.eventBus.Publish(ctx, statusEvent)
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

// ChangePassword altera a senha do usuário
func (s *UserService) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	// Buscar usuário
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	
	// Verificar senha atual
	if !user.VerifyPassword(req.CurrentPassword) {
		return fmt.Errorf("senha atual incorreta")
	}
	
	// Definir nova senha
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}
	
	user.UpdatedAt = time.Now()
	
	// Salvar usuário
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("erro ao atualizar senha: %w", err)
	}
	
	// Publicar evento
	event := domain.PasswordChangedEvent{
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		ChangedAt: time.Now(),
	}
	s.eventBus.Publish(ctx, event)
	
	return nil
}

// DeleteUser remove um usuário
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// Buscar usuário para pegar dados antes de deletar
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	
	// Deletar usuário
	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	
	// Publicar evento
	event := domain.UserDeletedEvent{
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		DeletedAt: time.Now(),
	}
	s.eventBus.Publish(ctx, event)
	
	return nil
}

// ActivateUser ativa um usuário
func (s *UserService) ActivateUser(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	
	oldStatus := user.Status
	user.Status = domain.StatusActive
	user.UpdatedAt = time.Now()
	
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("erro ao ativar usuário: %w", err)
	}
	
	// Publicar evento
	event := domain.UserStatusChangedEvent{
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		OldStatus: oldStatus,
		NewStatus: domain.StatusActive,
		ChangedAt: time.Now(),
	}
	s.eventBus.Publish(ctx, event)
	
	return nil
}

// DeactivateUser desativa um usuário
func (s *UserService) DeactivateUser(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	
	oldStatus := user.Status
	user.Status = domain.StatusInactive
	user.UpdatedAt = time.Now()
	
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("erro ao desativar usuário: %w", err)
	}
	
	// Publicar evento
	event := domain.UserStatusChangedEvent{
		UserID:    user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		OldStatus: oldStatus,
		NewStatus: domain.StatusInactive,
		ChangedAt: time.Now(),
	}
	s.eventBus.Publish(ctx, event)
	
	return nil
}