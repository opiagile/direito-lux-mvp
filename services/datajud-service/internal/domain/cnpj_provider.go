package domain

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// CNPJProvider representa um CNPJ habilitado para consultas na API DataJud
type CNPJProvider struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	CNPJ              string     `json:"cnpj"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	APIKey            string     `json:"api_key"`           // Chave de acesso DataJud
	Certificate       string     `json:"certificate"`       // Certificado digital
	CertificatePass   string     `json:"certificate_pass"`  // Senha do certificado
	DailyLimit        int        `json:"daily_limit"`       // Limite diário (padrão 10.000)
	DailyUsage        int        `json:"daily_usage"`       // Uso do dia atual
	UsageResetTime    time.Time  `json:"usage_reset_time"`  // Quando o contador reseta
	IsActive          bool       `json:"is_active"`
	Priority          int        `json:"priority"`          // Prioridade de uso (1 = mais alta)
	LastUsedAt        *time.Time `json:"last_used_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeactivatedAt     *time.Time `json:"deactivated_at"`
}

// CNPJProviderRepository interface para persistência
type CNPJProviderRepository interface {
	Save(provider *CNPJProvider) error
	FindByID(id uuid.UUID) (*CNPJProvider, error)
	FindByTenantID(tenantID uuid.UUID) ([]*CNPJProvider, error)
	FindActiveCNPJs() ([]*CNPJProvider, error)
	FindAvailableCNPJs(minQuota int) ([]*CNPJProvider, error)
	UpdateUsage(id uuid.UUID, usage int) error
	ResetDailyUsage() error
	Update(provider *CNPJProvider) error
	Delete(id uuid.UUID) error
}

// NewCNPJProvider cria um novo provedor CNPJ
func NewCNPJProvider(tenantID uuid.UUID, cnpj, name, email, apiKey string) (*CNPJProvider, error) {
	if !isValidCNPJ(cnpj) {
		return nil, errors.New("CNPJ inválido")
	}

	if name == "" {
		return nil, errors.New("nome é obrigatório")
	}

	if email == "" {
		return nil, errors.New("email é obrigatório")
	}

	if apiKey == "" {
		return nil, errors.New("chave da API é obrigatória")
	}

	now := time.Now()
	resetTime := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	return &CNPJProvider{
		ID:             uuid.New(),
		TenantID:       tenantID,
		CNPJ:           formatCNPJ(cnpj),
		Name:           name,
		Email:          email,
		APIKey:         apiKey,
		DailyLimit:     10000, // Limite padrão DataJud
		DailyUsage:     0,
		UsageResetTime: resetTime,
		IsActive:       true,
		Priority:       1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// CanMakeRequest verifica se o CNPJ pode fazer uma requisição
func (c *CNPJProvider) CanMakeRequest() bool {
	if !c.IsActive {
		return false
	}

	// Verifica se precisa resetar o contador diário
	if time.Now().After(c.UsageResetTime) {
		c.ResetDailyUsage()
	}

	return c.DailyUsage < c.DailyLimit
}

// GetAvailableQuota retorna quantas consultas ainda podem ser feitas hoje
func (c *CNPJProvider) GetAvailableQuota() int {
	if !c.IsActive {
		return 0
	}

	// Verifica se precisa resetar o contador diário
	if time.Now().After(c.UsageResetTime) {
		c.ResetDailyUsage()
	}

	available := c.DailyLimit - c.DailyUsage
	if available < 0 {
		return 0
	}
	return available
}

// UseQuota consome uma quota do CNPJ
func (c *CNPJProvider) UseQuota(amount int) error {
	if !c.CanMakeRequest() {
		return errors.New("CNPJ não pode fazer requisições no momento")
	}

	if amount <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}

	if c.DailyUsage+amount > c.DailyLimit {
		return errors.New("quota insuficiente para esta requisição")
	}

	c.DailyUsage += amount
	now := time.Now()
	c.LastUsedAt = &now
	c.UpdatedAt = now

	return nil
}

// ResetDailyUsage reseta o contador diário
func (c *CNPJProvider) ResetDailyUsage() {
	c.DailyUsage = 0
	now := time.Now()
	c.UsageResetTime = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	c.UpdatedAt = now
}

// Activate ativa o CNPJ provider
func (c *CNPJProvider) Activate() {
	c.IsActive = true
	c.DeactivatedAt = nil
	c.UpdatedAt = time.Now()
}

// Deactivate desativa o CNPJ provider
func (c *CNPJProvider) Deactivate(reason string) {
	c.IsActive = false
	now := time.Now()
	c.DeactivatedAt = &now
	c.UpdatedAt = now
}

// UpdateCertificate atualiza certificado e senha
func (c *CNPJProvider) UpdateCertificate(certificate, password string) error {
	if certificate == "" {
		return errors.New("certificado é obrigatório")
	}

	c.Certificate = certificate
	c.CertificatePass = password
	c.UpdatedAt = time.Now()

	return nil
}

// SetPriority define a prioridade de uso do CNPJ
func (c *CNPJProvider) SetPriority(priority int) {
	if priority < 1 {
		priority = 1
	}
	if priority > 10 {
		priority = 10
	}

	c.Priority = priority
	c.UpdatedAt = time.Now()
}

// GetUsagePercentage retorna o percentual de uso do dia
func (c *CNPJProvider) GetUsagePercentage() float64 {
	if c.DailyLimit == 0 {
		return 0
	}

	return float64(c.DailyUsage) / float64(c.DailyLimit) * 100
}

// isValidCNPJ valida formato do CNPJ
func isValidCNPJ(cnpj string) bool {
	// Remove caracteres não numéricos
	cnpjClean := regexp.MustCompile(`[^0-9]`).ReplaceAllString(cnpj, "")
	
	// Verifica se tem 14 dígitos
	if len(cnpjClean) != 14 {
		return false
	}

	// Verifica se não são todos dígitos iguais
	first := cnpjClean[0]
	allSame := true
	for _, char := range cnpjClean {
		if char != rune(first) {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	// Validação dos dígitos verificadores
	return validateCNPJDigits(cnpjClean)
}

// validateCNPJDigits valida os dígitos verificadores do CNPJ
func validateCNPJDigits(cnpj string) bool {
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// Primeiro dígito
	sum := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights1[i]
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	if firstDigit != int(cnpj[12]-'0') {
		return false
	}

	// Segundo dígito
	sum = 0
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights2[i]
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	return secondDigit == int(cnpj[13]-'0')
}

// formatCNPJ formata CNPJ no padrão XX.XXX.XXX/XXXX-XX
func formatCNPJ(cnpj string) string {
	cnpjClean := regexp.MustCompile(`[^0-9]`).ReplaceAllString(cnpj, "")
	
	if len(cnpjClean) != 14 {
		return cnpj // Retorna original se inválido
	}

	return cnpjClean[:2] + "." + cnpjClean[2:5] + "." + cnpjClean[5:8] + "/" + cnpjClean[8:12] + "-" + cnpjClean[12:14]
}