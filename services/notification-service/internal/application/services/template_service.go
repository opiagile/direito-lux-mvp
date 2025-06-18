package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/domain"
)

// TemplateService serviço de aplicação para templates
type TemplateService struct {
	templateRepo domain.NotificationTemplateRepository
	logger       *zap.Logger
}

// NewTemplateService cria nova instância do serviço
func NewTemplateService(
	templateRepo domain.NotificationTemplateRepository,
	logger *zap.Logger,
) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
		logger:       logger,
	}
}

// CreateTemplateRequest request para criação de template
type CreateTemplateRequest struct {
	Name      string                      `json:"name" validate:"required"`
	Type      domain.NotificationType     `json:"type" validate:"required"`
	Channel   domain.NotificationChannel  `json:"channel" validate:"required"`
	Subject   string                      `json:"subject"`
	Content   string                      `json:"content" validate:"required"`
	Variables []string                    `json:"variables,omitempty"`
	TenantID  *uuid.UUID                  `json:"tenant_id,omitempty"`
}

// UpdateTemplateRequest request para atualização de template
type UpdateTemplateRequest struct {
	Name      *string                     `json:"name,omitempty"`
	Subject   *string                     `json:"subject,omitempty"`
	Content   *string                     `json:"content,omitempty"`
	Variables []string                    `json:"variables,omitempty"`
	Status    *domain.TemplateStatus      `json:"status,omitempty"`
}

// CreateTemplate cria um novo template
func (s *TemplateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*domain.NotificationTemplate, error) {
	s.logger.Debug("Creating template", 
		zap.String("name", req.Name),
		zap.String("type", string(req.Type)),
		zap.String("channel", string(req.Channel)))

	// Validar que não é um template do sistema sendo criado por tenant
	if req.TenantID == nil {
		return nil, fmt.Errorf("cannot create system template through this endpoint")
	}

	// Extrair variáveis do conteúdo se não fornecidas
	variables := req.Variables
	if len(variables) == 0 {
		variables = s.extractVariables(req.Subject, req.Content)
	}

	template := &domain.NotificationTemplate{
		ID:        uuid.New(),
		Name:      req.Name,
		Type:      req.Type,
		Channel:   req.Channel,
		Status:    domain.TemplateStatusActive,
		Subject:   req.Subject,
		Content:   req.Content,
		Variables: variables,
		TenantID:  req.TenantID,
		IsSystem:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	s.logger.Info("Template created successfully", zap.String("template_id", template.ID.String()))
	return template, nil
}

// extractVariables extrai variáveis do formato {{variavel}} do conteúdo
func (s *TemplateService) extractVariables(subject, content string) []string {
	variables := make(map[string]bool)
	text := subject + " " + content

	// Buscar padrão {{variavel}}
	start := 0
	for {
		openIndex := strings.Index(text[start:], "{{")
		if openIndex == -1 {
			break
		}
		openIndex += start

		closeIndex := strings.Index(text[openIndex:], "}}")
		if closeIndex == -1 {
			break
		}
		closeIndex += openIndex

		if closeIndex > openIndex+2 {
			variable := text[openIndex+2 : closeIndex]
			variable = strings.TrimSpace(variable)
			if variable != "" {
				variables[variable] = true
			}
		}

		start = closeIndex + 2
	}

	// Converter para slice
	result := make([]string, 0, len(variables))
	for variable := range variables {
		result = append(result, variable)
	}

	return result
}

// GetTemplate busca template por ID
func (s *TemplateService) GetTemplate(ctx context.Context, id uuid.UUID) (*domain.NotificationTemplate, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// UpdateTemplate atualiza um template
func (s *TemplateService) UpdateTemplate(ctx context.Context, id uuid.UUID, req *UpdateTemplateRequest) (*domain.NotificationTemplate, error) {
	s.logger.Debug("Updating template", zap.String("template_id", id.String()))

	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	// Não permitir edição de templates do sistema
	if template.IsSystem {
		return nil, fmt.Errorf("cannot update system template")
	}

	// Aplicar mudanças
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Subject != nil {
		template.Subject = *req.Subject
	}
	if req.Content != nil {
		template.Content = *req.Content
	}
	if req.Status != nil {
		template.Status = *req.Status
	}

	// Atualizar variáveis se conteúdo mudou
	if req.Subject != nil || req.Content != nil {
		if req.Variables != nil {
			template.Variables = req.Variables
		} else {
			template.Variables = s.extractVariables(template.Subject, template.Content)
		}
	}

	template.UpdatedAt = time.Now()

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	s.logger.Info("Template updated successfully", zap.String("template_id", template.ID.String()))
	return template, nil
}

// DeleteTemplate remove um template
func (s *TemplateService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Deleting template", zap.String("template_id", id.String()))

	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// Não permitir exclusão de templates do sistema
	if template.IsSystem {
		return fmt.Errorf("cannot delete system template")
	}

	if err := s.templateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	s.logger.Info("Template deleted successfully", zap.String("template_id", id.String()))
	return nil
}

// ListTemplates lista templates com filtros
func (s *TemplateService) ListTemplates(ctx context.Context, tenantID *uuid.UUID, filters domain.TemplateFilters) ([]*domain.NotificationTemplate, error) {
	return s.templateRepo.FindByTenantID(ctx, tenantID, filters)
}

// ListSystemTemplates lista templates do sistema
func (s *TemplateService) ListSystemTemplates(ctx context.Context, filters domain.TemplateFilters) ([]*domain.NotificationTemplate, error) {
	return s.templateRepo.FindSystemTemplates(ctx, filters)
}

// GetTemplateByTypeAndChannel busca template por tipo e canal
func (s *TemplateService) GetTemplateByTypeAndChannel(ctx context.Context, notificationType domain.NotificationType, channel domain.NotificationChannel, tenantID uuid.UUID) (*domain.NotificationTemplate, error) {
	return s.templateRepo.FindByTypeAndChannel(ctx, notificationType, channel, tenantID)
}

// PreviewTemplate faz preview de um template com variáveis
func (s *TemplateService) PreviewTemplate(ctx context.Context, id uuid.UUID, variables map[string]interface{}) (*TemplatePreview, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	preview := &TemplatePreview{
		TemplateID: template.ID,
		Subject:    s.processVariables(template.Subject, variables),
		Content:    s.processVariables(template.Content, variables),
		Variables:  template.Variables,
	}

	// Verificar variáveis não fornecidas
	missingVars := []string{}
	for _, variable := range template.Variables {
		if _, exists := variables[variable]; !exists {
			missingVars = append(missingVars, variable)
		}
	}
	preview.MissingVariables = missingVars

	return preview, nil
}

// processVariables substitui variáveis no texto
func (s *TemplateService) processVariables(text string, variables map[string]interface{}) string {
	if variables == nil {
		return text
	}

	result := text
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

// ValidateTemplate valida a estrutura de um template
func (s *TemplateService) ValidateTemplate(ctx context.Context, req *CreateTemplateRequest) error {
	// Validar que o conteúdo não está vazio
	if strings.TrimSpace(req.Content) == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Validar que variáveis no conteúdo estão bem formadas
	if err := s.validateVariables(req.Subject, req.Content); err != nil {
		return fmt.Errorf("invalid variables in template: %w", err)
	}

	return nil
}

// validateVariables valida se as variáveis estão bem formadas
func (s *TemplateService) validateVariables(subject, content string) error {
	text := subject + " " + content
	
	// Contar abertura e fechamento de variáveis
	openCount := strings.Count(text, "{{")
	closeCount := strings.Count(text, "}}")

	if openCount != closeCount {
		return fmt.Errorf("mismatched variable delimiters")
	}

	// Verificar se todas as variáveis estão bem formadas
	start := 0
	for {
		openIndex := strings.Index(text[start:], "{{")
		if openIndex == -1 {
			break
		}
		openIndex += start

		closeIndex := strings.Index(text[openIndex:], "}}")
		if closeIndex == -1 {
			return fmt.Errorf("unclosed variable at position %d", openIndex)
		}
		closeIndex += openIndex

		if closeIndex <= openIndex+2 {
			return fmt.Errorf("empty variable at position %d", openIndex)
		}

		variable := text[openIndex+2 : closeIndex]
		if strings.TrimSpace(variable) == "" {
			return fmt.Errorf("empty variable name at position %d", openIndex)
		}

		start = closeIndex + 2
	}

	return nil
}

// DuplicateTemplate cria uma cópia de um template
func (s *TemplateService) DuplicateTemplate(ctx context.Context, id uuid.UUID, name string, tenantID *uuid.UUID) (*domain.NotificationTemplate, error) {
	s.logger.Debug("Duplicating template", 
		zap.String("template_id", id.String()),
		zap.String("new_name", name))

	original, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	duplicate := &domain.NotificationTemplate{
		ID:        uuid.New(),
		Name:      name,
		Type:      original.Type,
		Channel:   original.Channel,
		Status:    domain.TemplateStatusDraft,
		Subject:   original.Subject,
		Content:   original.Content,
		Variables: make([]string, len(original.Variables)),
		TenantID:  tenantID,
		IsSystem:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	copy(duplicate.Variables, original.Variables)

	if err := s.templateRepo.Create(ctx, duplicate); err != nil {
		return nil, fmt.Errorf("failed to duplicate template: %w", err)
	}

	s.logger.Info("Template duplicated successfully", 
		zap.String("original_id", original.ID.String()),
		zap.String("duplicate_id", duplicate.ID.String()))

	return duplicate, nil
}

// ActivateTemplate ativa um template
func (s *TemplateService) ActivateTemplate(ctx context.Context, id uuid.UUID) error {
	return s.updateTemplateStatus(ctx, id, domain.TemplateStatusActive)
}

// DeactivateTemplate desativa um template
func (s *TemplateService) DeactivateTemplate(ctx context.Context, id uuid.UUID) error {
	return s.updateTemplateStatus(ctx, id, domain.TemplateStatusInactive)
}

// updateTemplateStatus atualiza status do template
func (s *TemplateService) updateTemplateStatus(ctx context.Context, id uuid.UUID, status domain.TemplateStatus) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if template.IsSystem {
		return fmt.Errorf("cannot change status of system template")
	}

	template.Status = status
	template.UpdatedAt = time.Now()

	return s.templateRepo.Update(ctx, template)
}

// TemplatePreview resultado do preview de template
type TemplatePreview struct {
	TemplateID       uuid.UUID `json:"template_id"`
	Subject          string    `json:"subject"`
	Content          string    `json:"content"`
	Variables        []string  `json:"variables"`
	MissingVariables []string  `json:"missing_variables"`
}