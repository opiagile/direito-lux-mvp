package domain

import (
	"time"

	"github.com/google/uuid"
)

// Invoice representa uma fatura/nota fiscal
type Invoice struct {
	ID                uuid.UUID     `json:"id"`
	SubscriptionID    uuid.UUID     `json:"subscription_id"`
	TenantID          uuid.UUID     `json:"tenant_id"`
	PaymentID         *uuid.UUID    `json:"payment_id"`
	
	// Dados da fatura
	Number            string        `json:"number"`
	Amount            int64         `json:"amount"`        // em centavos
	TaxAmount         int64         `json:"tax_amount"`    // ISS, etc.
	DiscountAmount    int64         `json:"discount_amount"`
	TotalAmount       int64         `json:"total_amount"`
	
	// Período da fatura
	PeriodStart       time.Time     `json:"period_start"`
	PeriodEnd         time.Time     `json:"period_end"`
	
	// Dados de cobrança
	DueDate           time.Time     `json:"due_date"`
	IssuedAt          time.Time     `json:"issued_at"`
	PaidAt            *time.Time    `json:"paid_at"`
	
	// Status
	Status            InvoiceStatus `json:"status"`
	
	// Nota fiscal eletrônica
	NFENumber         *string       `json:"nfe_number"`
	NFEKey            *string       `json:"nfe_key"`
	NFEURL            *string       `json:"nfe_url"`
	NFEStatus         *string       `json:"nfe_status"`
	NFEIssuedAt       *time.Time    `json:"nfe_issued_at"`
	
	// Dados do cliente (para NF-e)
	CustomerName      string        `json:"customer_name"`
	CustomerEmail     string        `json:"customer_email"`
	CustomerDocument  string        `json:"customer_document"`
	CustomerPhone     *string       `json:"customer_phone"`
	CustomerAddress   *Address      `json:"customer_address"`
	
	// Dados da empresa (Curitiba)
	CompanyName       string        `json:"company_name"`
	CompanyDocument   string        `json:"company_document"`
	CompanyAddress    *Address      `json:"company_address"`
	
	// Integrações externas
	AsaasInvoiceID    *string       `json:"asaas_invoice_id"`
	
	// Controle
	RetryCount        int           `json:"retry_count"`
	LastRetryAt       *time.Time    `json:"last_retry_at"`
	ErrorMessage      *string       `json:"error_message"`
	
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// InvoiceStatus enumeração dos status da fatura
type InvoiceStatus string

const (
	InvoiceStatusDraft      InvoiceStatus = "draft"
	InvoiceStatusIssued     InvoiceStatus = "issued"
	InvoiceStatusSent       InvoiceStatus = "sent"
	InvoiceStatusPaid       InvoiceStatus = "paid"
	InvoiceStatusOverdue    InvoiceStatus = "overdue"
	InvoiceStatusCancelled  InvoiceStatus = "cancelled"
	InvoiceStatusRefunded   InvoiceStatus = "refunded"
)

// Address representa um endereço
type Address struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
	Country      string `json:"country"`
}

// NewInvoice cria uma nova fatura
func NewInvoice(subscriptionID, tenantID uuid.UUID, amount int64, periodStart, periodEnd time.Time) *Invoice {
	now := time.Now()
	
	return &Invoice{
		ID:             uuid.New(),
		SubscriptionID: subscriptionID,
		TenantID:       tenantID,
		Amount:         amount,
		TaxAmount:      calculateISS(amount), // ISS 2% para Curitiba
		DiscountAmount: 0,
		TotalAmount:    amount + calculateISS(amount),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		DueDate:        now.AddDate(0, 0, 7), // 7 dias para pagamento
		IssuedAt:       now,
		Status:         InvoiceStatusDraft,
		
		// Dados da empresa (Curitiba)
		CompanyName:     "Sua Empresa LTDA",
		CompanyDocument: "12.345.678/0001-90",
		CompanyAddress: &Address{
			Street:       "Rua das Flores",
			Number:       "123",
			Complement:   "Sala 456",
			Neighborhood: "Centro",
			City:         "Curitiba",
			State:        "PR",
			ZipCode:      "80000-000",
			Country:      "Brasil",
		},
		
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// calculateISS calcula o ISS (2% para Curitiba)
func calculateISS(amount int64) int64 {
	return int64(float64(amount) * 0.02)
}

// SetCustomerData define os dados do cliente
func (i *Invoice) SetCustomerData(name, email, document string, phone *string, address *Address) {
	i.CustomerName = name
	i.CustomerEmail = email
	i.CustomerDocument = document
	i.CustomerPhone = phone
	i.CustomerAddress = address
	i.UpdatedAt = time.Now()
}

// Issue emite a fatura
func (i *Invoice) Issue() {
	i.Status = InvoiceStatusIssued
	i.IssuedAt = time.Now()
	i.UpdatedAt = time.Now()
}

// MarkAsSent marca a fatura como enviada
func (i *Invoice) MarkAsSent() {
	i.Status = InvoiceStatusSent
	i.UpdatedAt = time.Now()
}

// MarkAsPaid marca a fatura como paga
func (i *Invoice) MarkAsPaid(paymentID uuid.UUID) {
	now := time.Now()
	i.Status = InvoiceStatusPaid
	i.PaidAt = &now
	i.PaymentID = &paymentID
	i.UpdatedAt = now
}

// MarkAsOverdue marca a fatura como vencida
func (i *Invoice) MarkAsOverdue() {
	i.Status = InvoiceStatusOverdue
	i.UpdatedAt = time.Now()
}

// Cancel cancela a fatura
func (i *Invoice) Cancel() {
	i.Status = InvoiceStatusCancelled
	i.UpdatedAt = time.Now()
}

// MarkAsRefunded marca a fatura como reembolsada
func (i *Invoice) MarkAsRefunded() {
	i.Status = InvoiceStatusRefunded
	i.UpdatedAt = time.Now()
}

// SetNFEData define os dados da nota fiscal eletrônica
func (i *Invoice) SetNFEData(number, key, url, status string) {
	now := time.Now()
	i.NFENumber = &number
	i.NFEKey = &key
	i.NFEURL = &url
	i.NFEStatus = &status
	i.NFEIssuedAt = &now
	i.UpdatedAt = now
}

// SetAsaasData define os dados do ASAAS
func (i *Invoice) SetAsaasData(invoiceID string) {
	i.AsaasInvoiceID = &invoiceID
	i.UpdatedAt = time.Now()
}

// IncrementRetry incrementa o contador de tentativas
func (i *Invoice) IncrementRetry(errorMsg string) {
	i.RetryCount++
	now := time.Now()
	i.LastRetryAt = &now
	i.ErrorMessage = &errorMsg
	i.UpdatedAt = now
}

// IsOverdue verifica se a fatura está vencida
func (i *Invoice) IsOverdue() bool {
	return time.Now().After(i.DueDate) && i.Status != InvoiceStatusPaid
}

// IsPaid verifica se a fatura está paga
func (i *Invoice) IsPaid() bool {
	return i.Status == InvoiceStatusPaid
}

// CanRetry verifica se pode tentar emitir a NF-e novamente
func (i *Invoice) CanRetry() bool {
	return i.RetryCount < 3 && i.Status != InvoiceStatusPaid
}

// GetFormattedAmount retorna o valor formatado
func (i *Invoice) GetFormattedAmount() float64 {
	return float64(i.Amount) / 100
}

// GetFormattedTaxAmount retorna o valor do imposto formatado
func (i *Invoice) GetFormattedTaxAmount() float64 {
	return float64(i.TaxAmount) / 100
}

// GetFormattedTotalAmount retorna o valor total formatado
func (i *Invoice) GetFormattedTotalAmount() float64 {
	return float64(i.TotalAmount) / 100
}

// GetDaysUntilDue retorna quantos dias restam até o vencimento
func (i *Invoice) GetDaysUntilDue() int {
	diff := i.DueDate.Sub(time.Now())
	days := int(diff.Hours() / 24)
	
	if days < 0 {
		return 0
	}
	
	return days
}

// GetDaysOverdue retorna quantos dias está vencida
func (i *Invoice) GetDaysOverdue() int {
	if !i.IsOverdue() {
		return 0
	}
	
	diff := time.Now().Sub(i.DueDate)
	return int(diff.Hours() / 24)
}

// GenerateNumber gera um número único para a fatura
func (i *Invoice) GenerateNumber() {
	// Formato: YYYYMM-NNNNNN
	now := time.Now()
	i.Number = now.Format("200601") + "-" + i.ID.String()[:6]
	i.UpdatedAt = time.Now()
}

// HasNFE verifica se a fatura tem nota fiscal eletrônica
func (i *Invoice) HasNFE() bool {
	return i.NFENumber != nil && *i.NFENumber != ""
}

// GetNFEDownloadURL retorna a URL para download da NF-e
func (i *Invoice) GetNFEDownloadURL() string {
	if i.NFEURL != nil {
		return *i.NFEURL
	}
	return ""
}

// ApplyDiscount aplica um desconto na fatura
func (i *Invoice) ApplyDiscount(discountAmount int64) {
	i.DiscountAmount = discountAmount
	i.TotalAmount = i.Amount + i.TaxAmount - i.DiscountAmount
	i.UpdatedAt = time.Now()
}

// GetInvoiceItems retorna os itens da fatura
func (i *Invoice) GetInvoiceItems() []InvoiceItem {
	return []InvoiceItem{
		{
			Description: "Assinatura Direito Lux",
			Quantity:    1,
			UnitPrice:   i.Amount,
			TotalPrice:  i.Amount,
			PeriodStart: i.PeriodStart,
			PeriodEnd:   i.PeriodEnd,
		},
	}
}

// InvoiceItem representa um item da fatura
type InvoiceItem struct {
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UnitPrice   int64     `json:"unit_price"`
	TotalPrice  int64     `json:"total_price"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}