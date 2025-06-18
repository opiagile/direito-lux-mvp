package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TemplateStatus status do template
type TemplateStatus string

const (
	TemplateStatusActive   TemplateStatus = "active"
	TemplateStatusInactive TemplateStatus = "inactive"
	TemplateStatusDraft    TemplateStatus = "draft"
)

// NotificationTemplate template de notificação
type NotificationTemplate struct {
	ID          uuid.UUID            `json:"id" db:"id"`
	Name        string               `json:"name" db:"name"`
	Type        NotificationType     `json:"type" db:"type"`
	Channel     NotificationChannel  `json:"channel" db:"channel"`
	Status      TemplateStatus       `json:"status" db:"status"`
	Subject     string               `json:"subject" db:"subject"`
	Content     string               `json:"content" db:"content"`
	ContentHTML *string              `json:"content_html,omitempty" db:"content_html"`
	Variables   []string             `json:"variables" db:"variables"`
	TenantID    *uuid.UUID           `json:"tenant_id,omitempty" db:"tenant_id"`
	IsSystem    bool                 `json:"is_system" db:"is_system"`
	CreatedAt   time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" db:"updated_at"`
}

// NewNotificationTemplate cria novo template
func NewNotificationTemplate(
	name string,
	notificationType NotificationType,
	channel NotificationChannel,
	subject, content string,
	tenantID *uuid.UUID,
) (*NotificationTemplate, error) {
	// Validações
	if name == "" {
		return nil, errors.New("name é obrigatório")
	}

	if err := ValidateNotificationType(notificationType); err != nil {
		return nil, err
	}

	if err := ValidateChannel(channel); err != nil {
		return nil, err
	}

	if subject == "" {
		return nil, errors.New("subject é obrigatório")
	}

	if content == "" {
		return nil, errors.New("content é obrigatório")
	}

	now := time.Now()
	
	template := &NotificationTemplate{
		ID:        uuid.New(),
		Name:      name,
		Type:      notificationType,
		Channel:   channel,
		Status:    TemplateStatusDraft,
		Subject:   subject,
		Content:   content,
		Variables: extractVariables(subject + " " + content),
		TenantID:  tenantID,
		IsSystem:  tenantID == nil, // Templates sem tenant são do sistema
		CreatedAt: now,
		UpdatedAt: now,
	}

	return template, nil
}

// Activate ativa o template
func (t *NotificationTemplate) Activate() error {
	if t.Status == TemplateStatusActive {
		return errors.New("template já está ativo")
	}

	t.Status = TemplateStatusActive
	t.UpdatedAt = time.Now()
	
	return nil
}

// Deactivate desativa o template
func (t *NotificationTemplate) Deactivate() error {
	if t.Status == TemplateStatusInactive {
		return errors.New("template já está inativo")
	}

	t.Status = TemplateStatusInactive
	t.UpdatedAt = time.Now()
	
	return nil
}

// UpdateContent atualiza conteúdo do template
func (t *NotificationTemplate) UpdateContent(subject, content string, contentHTML *string) error {
	if subject == "" {
		return errors.New("subject é obrigatório")
	}

	if content == "" {
		return errors.New("content é obrigatório")
	}

	t.Subject = subject
	t.Content = content
	t.ContentHTML = contentHTML
	t.Variables = extractVariables(subject + " " + content)
	t.UpdatedAt = time.Now()
	
	return nil
}

// IsActive verifica se template está ativo
func (t *NotificationTemplate) IsActive() bool {
	return t.Status == TemplateStatusActive
}

// CanBeUsedByTenant verifica se template pode ser usado pelo tenant
func (t *NotificationTemplate) CanBeUsedByTenant(tenantID uuid.UUID) bool {
	// Templates do sistema podem ser usados por qualquer tenant
	if t.IsSystem {
		return true
	}
	
	// Templates específicos só podem ser usados pelo próprio tenant
	return t.TenantID != nil && *t.TenantID == tenantID
}

// RenderSubject renderiza subject com variáveis
func (t *NotificationTemplate) RenderSubject(variables map[string]interface{}) string {
	return renderTemplate(t.Subject, variables)
}

// RenderContent renderiza content com variáveis
func (t *NotificationTemplate) RenderContent(variables map[string]interface{}) string {
	return renderTemplate(t.Content, variables)
}

// RenderContentHTML renderiza content HTML com variáveis
func (t *NotificationTemplate) RenderContentHTML(variables map[string]interface{}) *string {
	if t.ContentHTML == nil {
		return nil
	}
	
	rendered := renderTemplate(*t.ContentHTML, variables)
	return &rendered
}

// extractVariables extrai variáveis do template (formato {{variavel}})
func extractVariables(text string) []string {
	var variables []string
	
	start := 0
	for {
		// Procurar início da variável
		startIndex := strings.Index(text[start:], "{{")
		if startIndex == -1 {
			break
		}
		startIndex += start
		
		// Procurar fim da variável
		endIndex := strings.Index(text[startIndex:], "}}")
		if endIndex == -1 {
			break
		}
		endIndex += startIndex
		
		// Extrair nome da variável
		variable := strings.TrimSpace(text[startIndex+2 : endIndex])
		if variable != "" {
			// Verificar se já existe
			found := false
			for _, v := range variables {
				if v == variable {
					found = true
					break
				}
			}
			if !found {
				variables = append(variables, variable)
			}
		}
		
		start = endIndex + 2
	}
	
	return variables
}

// renderTemplate renderiza template substituindo variáveis
func renderTemplate(template string, variables map[string]interface{}) string {
	result := template
	
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}
	
	return result
}

// ValidateTemplateStatus valida status do template
func ValidateTemplateStatus(status TemplateStatus) error {
	switch status {
	case TemplateStatusActive,
		 TemplateStatusInactive,
		 TemplateStatusDraft:
		return nil
	default:
		return fmt.Errorf("status de template inválido: %s", status)
	}
}

// GetDefaultTemplates retorna templates padrão do sistema
func GetDefaultTemplates() []*NotificationTemplate {
	templates := []*NotificationTemplate{
		{
			ID:        uuid.New(),
			Name:      "Bem-vindo",
			Type:      NotificationTypeWelcome,
			Channel:   NotificationChannelEmail,
			Status:    TemplateStatusActive,
			Subject:   "Bem-vindo ao Direito Lux, {{nome}}!",
			Content:   "Olá {{nome}}, seja bem-vindo ao Direito Lux! Sua conta foi criada com sucesso.",
			Variables: []string{"nome"},
			IsSystem:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Movimentação Processual",
			Type:      NotificationTypeProcessUpdate,
			Channel:   NotificationChannelWhatsApp,
			Status:    TemplateStatusActive,
			Subject:   "Nova movimentação no processo {{numero_processo}}",
			Content:   "Olá {{nome}}, houve uma nova movimentação no processo {{numero_processo}}: {{descricao_movimento}}. Data: {{data_movimento}}",
			Variables: []string{"nome", "numero_processo", "descricao_movimento", "data_movimento"},
			IsSystem:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Lembrete de Prazo",
			Type:      NotificationTypeDeadlineReminder,
			Channel:   NotificationChannelWhatsApp,
			Status:    TemplateStatusActive,
			Subject:   "Prazo vencendo em {{dias}} dias - Processo {{numero_processo}}",
			Content:   "Atenção {{nome}}, o prazo para o processo {{numero_processo}} vence em {{dias}} dias. Descrição: {{descricao_prazo}}",
			Variables: []string{"nome", "numero_processo", "dias", "descricao_prazo"},
			IsSystem:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Trial Expirando",
			Type:      NotificationTypeTrialExpiring,
			Channel:   NotificationChannelEmail,
			Status:    TemplateStatusActive,
			Subject:   "Seu período de teste expira em {{dias}} dias",
			Content:   "Olá {{nome}}, seu período de teste do Direito Lux expira em {{dias}} dias. Não perca tempo e assine já!",
			Variables: []string{"nome", "dias"},
			IsSystem:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return templates
}