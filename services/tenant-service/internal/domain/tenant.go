package domain

import (
	"time"
	"errors"
	"regexp"
	"strings"
)

// Tenant representa uma organização/escritório no sistema
type Tenant struct {
	ID            string      `json:"id" db:"id"`
	Name          string      `json:"name" db:"name"`
	LegalName     string      `json:"legal_name" db:"legal_name"`
	Document      string      `json:"document" db:"document"` // CNPJ
	Email         string      `json:"email" db:"email"`
	Phone         string      `json:"phone" db:"phone"`
	Website       string      `json:"website" db:"website"`
	Address       Address     `json:"address" db:"address"`
	Status        TenantStatus `json:"status" db:"status"`
	PlanType      PlanType    `json:"plan_type" db:"plan_type"`
	OwnerUserID   string      `json:"owner_user_id" db:"owner_user_id"`
	Settings      TenantSettings `json:"settings" db:"settings"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	ActivatedAt   *time.Time  `json:"activated_at" db:"activated_at"`
	SuspendedAt   *time.Time  `json:"suspended_at" db:"suspended_at"`
}

// Address representa o endereço do tenant
type Address struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
	Country    string `json:"country"`
}

// TenantSettings contém configurações específicas do tenant
type TenantSettings struct {
	Timezone         string            `json:"timezone"`
	Language         string            `json:"language"`
	Currency         string            `json:"currency"`
	DateFormat       string            `json:"date_format"`
	NotificationPrefs NotificationPrefs `json:"notification_prefs"`
	BrandingConfig   BrandingConfig    `json:"branding_config"`
	IntegrationKeys  map[string]string `json:"integration_keys"`
}

// NotificationPrefs configurações de notificação do tenant
type NotificationPrefs struct {
	WhatsAppEnabled   bool `json:"whatsapp_enabled"`
	EmailEnabled      bool `json:"email_enabled"`
	TelegramEnabled   bool `json:"telegram_enabled"`
	PushEnabled       bool `json:"push_enabled"`
	BusinessHoursOnly bool `json:"business_hours_only"`
}

// BrandingConfig configurações de marca/visual
type BrandingConfig struct {
	LogoURL      string `json:"logo_url"`
	PrimaryColor string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	FontFamily   string `json:"font_family"`
	CustomDomain string `json:"custom_domain"`
}

// TenantStatus define os status possíveis de um tenant
type TenantStatus string

const (
	TenantStatusPending   TenantStatus = "pending"
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusCanceled  TenantStatus = "canceled"
	TenantStatusBlocked   TenantStatus = "blocked"
)

// PlanType define os tipos de plano disponíveis
type PlanType string

const (
	PlanStarter      PlanType = "starter"
	PlanProfessional PlanType = "professional"
	PlanBusiness     PlanType = "business"
	PlanEnterprise   PlanType = "enterprise"
)

// TenantRepository define a interface para persistência de tenants
type TenantRepository interface {
	Create(tenant *Tenant) error
	GetByID(id string) (*Tenant, error)
	GetByDocument(document string) (*Tenant, error)
	GetByOwner(ownerUserID string) ([]*Tenant, error)
	GetAll(limit, offset int) ([]*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id string) error
	GetByStatus(status TenantStatus, limit, offset int) ([]*Tenant, error)
}

// Erros de domínio
var (
	ErrTenantNotFound      = errors.New("tenant não encontrado")
	ErrTenantExists        = errors.New("tenant já existe")
	ErrInvalidDocument     = errors.New("documento inválido")
	ErrInvalidEmail        = errors.New("email inválido")
	ErrInvalidName         = errors.New("nome inválido")
	ErrInvalidPlan         = errors.New("plano inválido")
	ErrInvalidStatus       = errors.New("status inválido")
	ErrTenantLimitExceeded = errors.New("limite de tenants excedido")
)

// ValidateName valida o nome do tenant
func (t *Tenant) ValidateName() error {
	if len(strings.TrimSpace(t.Name)) < 3 {
		return ErrInvalidName
	}
	if len(t.Name) > 50 {
		return ErrInvalidName
	}
	return nil
}

// ValidateDocument valida o CNPJ
func (t *Tenant) ValidateDocument() error {
	if t.Document == "" {
		return nil // Documento é opcional
	}
	
	// Remove caracteres não numéricos
	doc := regexp.MustCompile(`[^0-9]`).ReplaceAllString(t.Document, "")
	
	// CNPJ deve ter 14 dígitos
	if len(doc) != 14 {
		return ErrInvalidDocument
	}
	
	// Validação básica de CNPJ (algoritmo simplificado)
	if !isValidCNPJ(doc) {
		return ErrInvalidDocument
	}
	
	t.Document = doc
	return nil
}

// ValidateEmail valida o email do tenant
func (t *Tenant) ValidateEmail() error {
	if t.Email == "" {
		return ErrInvalidEmail
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(t.Email) {
		return ErrInvalidEmail
	}
	
	return nil
}

// ValidatePlan valida o tipo de plano
func (t *Tenant) ValidatePlan() error {
	validPlans := []PlanType{PlanStarter, PlanProfessional, PlanBusiness, PlanEnterprise}
	for _, plan := range validPlans {
		if t.PlanType == plan {
			return nil
		}
	}
	return ErrInvalidPlan
}

// ValidateStatus valida o status do tenant
func (t *Tenant) ValidateStatus() error {
	validStatuses := []TenantStatus{
		TenantStatusPending, TenantStatusActive, TenantStatusSuspended,
		TenantStatusCanceled, TenantStatusBlocked,
	}
	for _, status := range validStatuses {
		if t.Status == status {
			return nil
		}
	}
	return ErrInvalidStatus
}

// IsActive verifica se o tenant está ativo
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

// CanAccess verifica se o tenant pode acessar o sistema
func (t *Tenant) CanAccess() bool {
	return t.Status == TenantStatusActive || t.Status == TenantStatusPending
}

// Activate ativa o tenant
func (t *Tenant) Activate() {
	t.Status = TenantStatusActive
	now := time.Now()
	t.ActivatedAt = &now
	t.UpdatedAt = now
}

// Suspend suspende o tenant
func (t *Tenant) Suspend() {
	t.Status = TenantStatusSuspended
	now := time.Now()
	t.SuspendedAt = &now
	t.UpdatedAt = now
}

// Cancel cancela o tenant
func (t *Tenant) Cancel() {
	t.Status = TenantStatusCanceled
	t.UpdatedAt = time.Now()
}

// GetDisplayName retorna o nome para exibição
func (t *Tenant) GetDisplayName() string {
	if t.LegalName != "" {
		return t.LegalName
	}
	return t.Name
}

// GetDomain retorna o domínio para uso interno
func (t *Tenant) GetDomain() string {
	if t.Settings.BrandingConfig.CustomDomain != "" {
		return t.Settings.BrandingConfig.CustomDomain
	}
	return strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(t.Name, ""))
}

// SetDefaultSettings define configurações padrão para o tenant
func (t *Tenant) SetDefaultSettings() {
	t.Settings = TenantSettings{
		Timezone:    "America/Sao_Paulo",
		Language:    "pt-BR",
		Currency:    "BRL",
		DateFormat:  "DD/MM/YYYY",
		NotificationPrefs: NotificationPrefs{
			WhatsAppEnabled:   true,
			EmailEnabled:      true,
			TelegramEnabled:   false,
			PushEnabled:       true,
			BusinessHoursOnly: false,
		},
		BrandingConfig: BrandingConfig{
			PrimaryColor:   "#1976d2",
			SecondaryColor: "#424242",
			FontFamily:     "Roboto",
		},
		IntegrationKeys: make(map[string]string),
	}
}

// isValidCNPJ valida CNPJ usando algoritmo oficial
func isValidCNPJ(cnpj string) bool {
	// Implementação simplificada - em produção usar biblioteca específica
	if len(cnpj) != 14 {
		return false
	}
	
	// Verifica se todos os dígitos são iguais
	first := cnpj[0]
	allSame := true
	for _, digit := range cnpj {
		if digit != first {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	
	// Para simplificar, aceita qualquer CNPJ com 14 dígitos diferentes
	return true
}