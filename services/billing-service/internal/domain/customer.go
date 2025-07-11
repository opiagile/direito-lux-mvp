package domain

import (
	"time"

	"github.com/google/uuid"
)

// Customer representa um cliente para faturamento
type Customer struct {
	ID                uuid.UUID     `json:"id"`
	TenantID          uuid.UUID     `json:"tenant_id"`
	
	// Dados pessoais
	Name              string        `json:"name"`
	Email             string        `json:"email"`
	Phone             *string       `json:"phone"`
	Document          string        `json:"document"`      // CPF ou CNPJ
	DocumentType      DocumentType  `json:"document_type"` // CPF ou CNPJ
	
	// Dados da empresa (se CNPJ)
	CompanyName       *string       `json:"company_name"`
	TradingName       *string       `json:"trading_name"`
	StateRegistration *string       `json:"state_registration"`
	
	// Endereço
	Address           *Address      `json:"address"`
	
	// Dados de cobrança
	BillingEmail      *string       `json:"billing_email"`
	BillingPhone      *string       `json:"billing_phone"`
	
	// Integrações externas
	AsaasCustomerID   *string       `json:"asaas_customer_id"`
	
	// Status
	Status            CustomerStatus `json:"status"`
	
	// Controle
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// DocumentType enumeração dos tipos de documento
type DocumentType string

const (
	DocumentTypeCPF  DocumentType = "cpf"
	DocumentTypeCNPJ DocumentType = "cnpj"
)

// CustomerStatus enumeração dos status do cliente
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusBlocked   CustomerStatus = "blocked"
)

// NewCustomer cria um novo cliente
func NewCustomer(tenantID uuid.UUID, name, email, document string, documentType DocumentType) *Customer {
	now := time.Now()
	
	return &Customer{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         name,
		Email:        email,
		Document:     document,
		DocumentType: documentType,
		Status:       CustomerStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// SetAddress define o endereço do cliente
func (c *Customer) SetAddress(address *Address) {
	c.Address = address
	c.UpdatedAt = time.Now()
}

// SetPhone define o telefone do cliente
func (c *Customer) SetPhone(phone string) {
	c.Phone = &phone
	c.UpdatedAt = time.Now()
}

// SetBillingData define os dados de cobrança
func (c *Customer) SetBillingData(email, phone string) {
	c.BillingEmail = &email
	c.BillingPhone = &phone
	c.UpdatedAt = time.Now()
}

// SetCompanyData define os dados da empresa (para CNPJ)
func (c *Customer) SetCompanyData(companyName, tradingName, stateRegistration string) {
	c.CompanyName = &companyName
	c.TradingName = &tradingName
	c.StateRegistration = &stateRegistration
	c.UpdatedAt = time.Now()
}

// SetAsaasCustomerID define o ID do cliente no ASAAS
func (c *Customer) SetAsaasCustomerID(asaasID string) {
	c.AsaasCustomerID = &asaasID
	c.UpdatedAt = time.Now()
}

// Activate ativa o cliente
func (c *Customer) Activate() {
	c.Status = CustomerStatusActive
	c.UpdatedAt = time.Now()
}

// Deactivate desativa o cliente
func (c *Customer) Deactivate() {
	c.Status = CustomerStatusInactive
	c.UpdatedAt = time.Now()
}

// Block bloqueia o cliente
func (c *Customer) Block() {
	c.Status = CustomerStatusBlocked
	c.UpdatedAt = time.Now()
}

// IsActive verifica se o cliente está ativo
func (c *Customer) IsActive() bool {
	return c.Status == CustomerStatusActive
}

// IsBlocked verifica se o cliente está bloqueado
func (c *Customer) IsBlocked() bool {
	return c.Status == CustomerStatusBlocked
}

// IsCorporate verifica se é pessoa jurídica
func (c *Customer) IsCorporate() bool {
	return c.DocumentType == DocumentTypeCNPJ
}

// GetDisplayName retorna o nome para exibição
func (c *Customer) GetDisplayName() string {
	if c.IsCorporate() && c.TradingName != nil && *c.TradingName != "" {
		return *c.TradingName
	}
	
	if c.IsCorporate() && c.CompanyName != nil && *c.CompanyName != "" {
		return *c.CompanyName
	}
	
	return c.Name
}

// GetBillingEmail retorna o email de cobrança
func (c *Customer) GetBillingEmail() string {
	if c.BillingEmail != nil && *c.BillingEmail != "" {
		return *c.BillingEmail
	}
	return c.Email
}

// GetBillingPhone retorna o telefone de cobrança
func (c *Customer) GetBillingPhone() string {
	if c.BillingPhone != nil && *c.BillingPhone != "" {
		return *c.BillingPhone
	}
	
	if c.Phone != nil {
		return *c.Phone
	}
	
	return ""
}

// ValidateDocument valida o documento (CPF/CNPJ)
func (c *Customer) ValidateDocument() bool {
	if c.Document == "" {
		return false
	}
	
	switch c.DocumentType {
	case DocumentTypeCPF:
		return c.validateCPF()
	case DocumentTypeCNPJ:
		return c.validateCNPJ()
	default:
		return false
	}
}

// validateCPF valida CPF (implementação básica)
func (c *Customer) validateCPF() bool {
	cpf := c.cleanDocument()
	
	if len(cpf) != 11 {
		return false
	}
	
	// Verificar se todos os dígitos são iguais
	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	
	if allSame {
		return false
	}
	
	// Implementar validação completa do CPF se necessário
	return true
}

// validateCNPJ valida CNPJ (implementação básica)
func (c *Customer) validateCNPJ() bool {
	cnpj := c.cleanDocument()
	
	if len(cnpj) != 14 {
		return false
	}
	
	// Implementar validação completa do CNPJ se necessário
	return true
}

// cleanDocument remove caracteres especiais do documento
func (c *Customer) cleanDocument() string {
	clean := ""
	for _, char := range c.Document {
		if char >= '0' && char <= '9' {
			clean += string(char)
		}
	}
	return clean
}

// GetFormattedDocument retorna o documento formatado
func (c *Customer) GetFormattedDocument() string {
	clean := c.cleanDocument()
	
	switch c.DocumentType {
	case DocumentTypeCPF:
		if len(clean) == 11 {
			return clean[:3] + "." + clean[3:6] + "." + clean[6:9] + "-" + clean[9:]
		}
	case DocumentTypeCNPJ:
		if len(clean) == 14 {
			return clean[:2] + "." + clean[2:5] + "." + clean[5:8] + "/" + clean[8:12] + "-" + clean[12:]
		}
	}
	
	return c.Document
}

// HasCompleteAddress verifica se tem endereço completo
func (c *Customer) HasCompleteAddress() bool {
	return c.Address != nil &&
		c.Address.Street != "" &&
		c.Address.Number != "" &&
		c.Address.City != "" &&
		c.Address.State != "" &&
		c.Address.ZipCode != ""
}

// CanReceiveInvoice verifica se pode receber fatura
func (c *Customer) CanReceiveInvoice() bool {
	return c.IsActive() &&
		c.Email != "" &&
		c.ValidateDocument() &&
		c.HasCompleteAddress()
}

// GetAsaasCustomerData retorna os dados formatados para o ASAAS
func (c *Customer) GetAsaasCustomerData() map[string]interface{} {
	data := map[string]interface{}{
		"name":                c.GetDisplayName(),
		"email":               c.GetBillingEmail(),
		"cpfCnpj":             c.cleanDocument(),
		"mobilePhone":         c.GetBillingPhone(),
		"externalReference":   c.ID.String(),
		"notificationDisabled": false,
	}
	
	if c.HasCompleteAddress() {
		data["address"] = c.Address.Street
		data["addressNumber"] = c.Address.Number
		data["complement"] = c.Address.Complement
		data["province"] = c.Address.Neighborhood
		data["city"] = c.Address.City
		data["state"] = c.Address.State
		data["postalCode"] = c.Address.ZipCode
	}
	
	if c.IsCorporate() {
		if c.CompanyName != nil {
			data["company"] = *c.CompanyName
		}
		if c.StateRegistration != nil {
			data["stateInscription"] = *c.StateRegistration
		}
	}
	
	return data
}