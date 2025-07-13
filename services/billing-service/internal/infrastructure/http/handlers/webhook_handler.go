package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/direito-lux/billing-service/internal/application"
	"github.com/direito-lux/billing-service/internal/domain"
)

// WebhookHandler handler para webhooks dos gateways de pagamento
type WebhookHandler struct {
	paymentService      *application.PaymentService
	subscriptionService *application.SubscriptionService
	paymentRepo         domain.PaymentRepository
	subscriptionRepo    domain.SubscriptionRepository
}

// NewWebhookHandler cria novo handler de webhooks
func NewWebhookHandler(
	paymentService *application.PaymentService,
	subscriptionService *application.SubscriptionService,
	paymentRepo domain.PaymentRepository,
	subscriptionRepo domain.SubscriptionRepository,
) *WebhookHandler {
	return &WebhookHandler{
		paymentService:      paymentService,
		subscriptionService: subscriptionService,
		paymentRepo:         paymentRepo,
		subscriptionRepo:    subscriptionRepo,
	}
}

// AsaasWebhook processa webhooks do ASAAS
func (h *WebhookHandler) AsaasWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Ler body do webhook
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Parse do webhook ASAAS
	var webhook AsaasWebhookPayload
	if err := json.Unmarshal(body, &webhook); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validar webhook (opcional - implementar verificação de assinatura)
	if !h.validateAsaasWebhook(webhook, r.Header.Get("X-Asaas-Signature")) {
		h.writeError(w, http.StatusUnauthorized, "Invalid webhook signature")
		return
	}

	// Processar evento baseado no tipo
	switch webhook.Event {
	case "PAYMENT_CONFIRMED":
		err = h.processAsaasPaymentConfirmed(ctx, webhook)
	case "PAYMENT_RECEIVED":
		err = h.processAsaasPaymentReceived(ctx, webhook)
	case "PAYMENT_OVERDUE":
		err = h.processAsaasPaymentOverdue(ctx, webhook)
	case "PAYMENT_DELETED":
		err = h.processAsaasPaymentDeleted(ctx, webhook)
	case "PAYMENT_RESTORED":
		err = h.processAsaasPaymentRestored(ctx, webhook)
	case "PAYMENT_REFUNDED":
		err = h.processAsaasPaymentRefunded(ctx, webhook)
	case "PAYMENT_RECEIVED_IN_CASH":
		err = h.processAsaasPaymentReceived(ctx, webhook)
	default:
		// Evento não reconhecido, mas retorna 200 para não retry
		h.writeSuccess(w, "Event not handled")
		return
	}

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeSuccess(w, "Webhook processed successfully")
}

// CryptoWebhook processa webhooks de criptomoedas (NOWPayments)
func (h *WebhookHandler) CryptoWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Ler body do webhook
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Parse do webhook NOWPayments
	var webhook CryptoWebhookPayload
	if err := json.Unmarshal(body, &webhook); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validar webhook (verificar assinatura HMAC)
	if !h.validateCryptoWebhook(body, r.Header.Get("X-Nowpayments-Sig")) {
		h.writeError(w, http.StatusUnauthorized, "Invalid webhook signature")
		return
	}

	// Processar evento baseado no status
	switch webhook.PaymentStatus {
	case "confirmed":
		err = h.processCryptoPaymentConfirmed(ctx, webhook)
	case "partially_paid":
		err = h.processCryptoPaymentPartial(ctx, webhook)
	case "failed":
		err = h.processCryptoPaymentFailed(ctx, webhook)
	case "refunded":
		err = h.processCryptoPaymentRefunded(ctx, webhook)
	case "expired":
		err = h.processCryptoPaymentExpired(ctx, webhook)
	default:
		// Status não reconhecido, mas retorna 200
		h.writeSuccess(w, "Status not handled")
		return
	}

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeSuccess(w, "Webhook processed successfully")
}

// processAsaasPaymentConfirmed processa pagamento confirmado no ASAAS
func (h *WebhookHandler) processAsaasPaymentConfirmed(ctx context.Context, webhook AsaasWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByAsaasID(ctx, webhook.Payment.ID)
	if err != nil {
		return err
	}

	// Processar sucesso do pagamento
	return h.paymentService.ProcessPaymentSuccess(ctx, payment.ID, webhook.Payment.ID)
}

// processAsaasPaymentReceived processa pagamento recebido no ASAAS
func (h *WebhookHandler) processAsaasPaymentReceived(ctx context.Context, webhook AsaasWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByAsaasID(ctx, webhook.Payment.ID)
	if err != nil {
		return err
	}

	// Processar sucesso do pagamento
	return h.paymentService.ProcessPaymentSuccess(ctx, payment.ID, webhook.Payment.ID)
}

// processAsaasPaymentOverdue processa pagamento vencido no ASAAS
func (h *WebhookHandler) processAsaasPaymentOverdue(ctx context.Context, webhook AsaasWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByAsaasID(ctx, webhook.Payment.ID)
	if err != nil {
		return err
	}

	// Processar falha do pagamento
	return h.paymentService.ProcessPaymentFailure(ctx, payment.ID, "Payment overdue")
}

// processAsaasPaymentDeleted processa pagamento deletado no ASAAS
func (h *WebhookHandler) processAsaasPaymentDeleted(ctx context.Context, webhook AsaasWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByAsaasID(ctx, webhook.Payment.ID)
	if err != nil {
		return err
	}

	// Marcar pagamento como cancelado
	return h.paymentService.ProcessPaymentFailure(ctx, payment.ID, "Payment deleted")
}

// processAsaasPaymentRestored processa pagamento restaurado no ASAAS
func (h *WebhookHandler) processAsaasPaymentRestored(ctx context.Context, webhook AsaasWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByAsaasID(ctx, webhook.Payment.ID)
	if err != nil {
		return err
	}

	// Reativar pagamento se necessário
	if payment.Status == domain.PaymentStatusCancelled {
		payment.Status = domain.PaymentStatusPending
		payment.UpdatedAt = time.Now()
		return h.paymentRepo.Update(ctx, payment)
	}

	return nil
}

// processAsaasPaymentRefunded processa reembolso no ASAAS
func (h *WebhookHandler) processAsaasPaymentRefunded(ctx context.Context, webhook AsaasWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByAsaasID(ctx, webhook.Payment.ID)
	if err != nil {
		return err
	}

	// Processar reembolso
	refundAmount := int64(webhook.Payment.Value * 100) // Converter para centavos
	return h.paymentService.RefundPayment(ctx, payment.ID, refundAmount, "Refunded by gateway")
}

// processCryptoPaymentConfirmed processa pagamento cripto confirmado
func (h *WebhookHandler) processCryptoPaymentConfirmed(ctx context.Context, webhook CryptoWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByNOWPaymentID(ctx, webhook.PaymentID)
	if err != nil {
		return err
	}

	// Atualizar dados cripto
	payment.SetCryptoData(
		webhook.PaymentAddress,
		webhook.PaymentAmount,
		webhook.OutcomeTxHash,
		webhook.ExchangeRate,
	)

	// Atualizar no banco
	if err := h.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	// Processar sucesso do pagamento
	return h.paymentService.ProcessPaymentSuccess(ctx, payment.ID, webhook.OutcomeTxHash)
}

// processCryptoPaymentPartial processa pagamento cripto parcial
func (h *WebhookHandler) processCryptoPaymentPartial(ctx context.Context, webhook CryptoWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByNOWPaymentID(ctx, webhook.PaymentID)
	if err != nil {
		return err
	}

	// Marcar como parcial
	payment.Status = domain.PaymentStatusPartial
	payment.UpdatedAt = time.Now()

	// Atualizar dados cripto
	payment.SetCryptoData(
		webhook.PaymentAddress,
		webhook.PaymentAmount,
		webhook.OutcomeTxHash,
		webhook.ExchangeRate,
	)

	return h.paymentRepo.Update(ctx, payment)
}

// processCryptoPaymentFailed processa pagamento cripto falho
func (h *WebhookHandler) processCryptoPaymentFailed(ctx context.Context, webhook CryptoWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByNOWPaymentID(ctx, webhook.PaymentID)
	if err != nil {
		return err
	}

	// Processar falha
	return h.paymentService.ProcessPaymentFailure(ctx, payment.ID, "Crypto payment failed")
}

// processCryptoPaymentRefunded processa reembolso cripto
func (h *WebhookHandler) processCryptoPaymentRefunded(ctx context.Context, webhook CryptoWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByNOWPaymentID(ctx, webhook.PaymentID)
	if err != nil {
		return err
	}

	// Processar reembolso
	refundAmount := payment.Amount // Reembolso total para cripto
	return h.paymentService.RefundPayment(ctx, payment.ID, refundAmount, "Refunded by gateway")
}

// processCryptoPaymentExpired processa pagamento cripto expirado
func (h *WebhookHandler) processCryptoPaymentExpired(ctx context.Context, webhook CryptoWebhookPayload) error {
	// Buscar pagamento por ID externo
	payment, err := h.paymentRepo.GetByNOWPaymentID(ctx, webhook.PaymentID)
	if err != nil {
		return err
	}

	// Marcar como expirado
	payment.ExpirePayment()
	return h.paymentRepo.Update(ctx, payment)
}

// validateAsaasWebhook valida webhook do ASAAS
func (h *WebhookHandler) validateAsaasWebhook(webhook AsaasWebhookPayload, signature string) bool {
	// TODO: Implementar validação de assinatura do ASAAS
	// Por enquanto, sempre retorna true
	return true
}

// validateCryptoWebhook valida webhook de criptomoedas
func (h *WebhookHandler) validateCryptoWebhook(body []byte, signature string) bool {
	// TODO: Implementar validação HMAC do NOWPayments
	// Por enquanto, sempre retorna true
	return true
}

// writeSuccess escreve resposta de sucesso
func (h *WebhookHandler) writeSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// writeError escreve resposta de erro
func (h *WebhookHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RegisterRoutes registra as rotas de webhook
func (h *WebhookHandler) RegisterRoutes(router *mux.Router) {
	// Webhooks
	router.HandleFunc("/webhooks/asaas", h.AsaasWebhook).Methods("POST")
	router.HandleFunc("/webhooks/crypto", h.CryptoWebhook).Methods("POST")
	router.HandleFunc("/webhooks/crypto/{payment_id}", h.CryptoWebhook).Methods("POST")
}

// Estruturas de payload dos webhooks

// AsaasWebhookPayload payload do webhook ASAAS
type AsaasWebhookPayload struct {
	Event   string `json:"event"`
	Payment struct {
		ID               string    `json:"id"`
		Status           string    `json:"status"`
		Value            float64   `json:"value"`
		NetValue         float64   `json:"netValue"`
		DueDate          time.Time `json:"dueDate"`
		PaymentDate      *time.Time `json:"paymentDate"`
		ExternalReference string   `json:"externalReference"`
		Customer         struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			CpfCnpj  string `json:"cpfCnpj"`
		} `json:"customer"`
		BillingType string `json:"billingType"`
		Description string `json:"description"`
	} `json:"payment"`
}

// CryptoWebhookPayload payload do webhook de criptomoedas
type CryptoWebhookPayload struct {
	PaymentID       string  `json:"payment_id"`
	PaymentStatus   string  `json:"payment_status"`
	PaymentAddress  string  `json:"payment_address"`
	PaymentAmount   string  `json:"payment_amount"`
	PriceAmount     float64 `json:"price_amount"`
	PriceCurrency   string  `json:"price_currency"`
	PaymentCurrency string  `json:"payment_currency"`
	OrderID         string  `json:"order_id"`
	OrderDescription string `json:"order_description"`
	OutcomeTxHash   string  `json:"outcome_tx_hash"`
	ExchangeRate    float64 `json:"exchange_rate"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}