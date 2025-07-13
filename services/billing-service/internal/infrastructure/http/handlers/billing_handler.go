package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/direito-lux/billing-service/internal/application"
)

// BillingHandler handler para endpoints de billing
type BillingHandler struct {
	subscriptionService *application.SubscriptionService
	paymentService      *application.PaymentService
	onboardingService   *application.OnboardingService
}

// NewBillingHandler cria um novo handler de billing
func NewBillingHandler(
	subscriptionService *application.SubscriptionService,
	paymentService *application.PaymentService,
	onboardingService *application.OnboardingService,
) *BillingHandler {
	return &BillingHandler{
		subscriptionService: subscriptionService,
		paymentService:      paymentService,
		onboardingService:   onboardingService,
	}
}

// CreateSubscription cria uma nova assinatura
func (h *BillingHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var cmd application.CreateSubscriptionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Obter tenant ID do contexto (middleware de autenticação)
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Tenant ID not found")
		return
	}
	cmd.TenantID = tenantID

	subscription, err := h.subscriptionService.CreateSubscription(ctx, cmd)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, subscription)
}

// GetSubscription busca assinatura por ID
func (h *BillingHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	subscriptionID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(ctx, subscriptionID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	h.writeJSON(w, http.StatusOK, subscription)
}

// GetCurrentSubscription busca assinatura atual do tenant
func (h *BillingHandler) GetCurrentSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Tenant ID not found")
		return
	}

	subscription, err := h.subscriptionService.GetSubscriptionByTenant(ctx, tenantID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	h.writeJSON(w, http.StatusOK, subscription)
}

// CancelSubscription cancela assinatura
func (h *BillingHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	subscriptionID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Obter user ID do contexto (quem está cancelando)
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "User ID not found")
		return
	}

	err = h.subscriptionService.CancelSubscription(ctx, subscriptionID, req.Reason, userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"message": "Subscription cancelled successfully"})
}

// ChangeSubscriptionPlan muda plano da assinatura
func (h *BillingHandler) ChangeSubscriptionPlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	subscriptionID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var req struct {
		NewPlanID uuid.UUID `json:"new_plan_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.subscriptionService.ChangeSubscriptionPlan(ctx, subscriptionID, req.NewPlanID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"message": "Plan changed successfully"})
}

// CreatePayment cria um novo pagamento
func (h *BillingHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var cmd application.CreatePaymentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Obter tenant ID do contexto
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Tenant ID not found")
		return
	}
	cmd.TenantID = tenantID

	payment, err := h.paymentService.CreatePayment(ctx, cmd)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, payment)
}

// GetPayment busca pagamento por ID
func (h *BillingHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	paymentID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.paymentService.GetPayment(ctx, paymentID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Payment not found")
		return
	}

	h.writeJSON(w, http.StatusOK, payment)
}

// GetPaymentsByTenant busca pagamentos do tenant
func (h *BillingHandler) GetPaymentsByTenant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Tenant ID not found")
		return
	}

	payments, err := h.paymentService.GetPaymentsByTenant(ctx, tenantID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, payments)
}

// RefundPayment processa reembolso
func (h *BillingHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	paymentID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var req struct {
		Amount int64  `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.paymentService.RefundPayment(ctx, paymentID, req.Amount, req.Reason)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"message": "Payment refunded successfully"})
}

// StartOnboarding inicia processo de onboarding
func (h *BillingHandler) StartOnboarding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data application.OnboardingData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Obter tenant ID do contexto
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Tenant ID not found")
		return
	}
	data.TenantID = tenantID

	result, err := h.onboardingService.StartOnboarding(ctx, data)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, result)
}

// GetOnboardingStatus retorna status do onboarding
func (h *BillingHandler) GetOnboardingStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Tenant ID not found")
		return
	}

	status, err := h.onboardingService.GetOnboardingStatus(ctx, tenantID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, status)
}

// GetAvailablePlans retorna planos disponíveis
func (h *BillingHandler) GetAvailablePlans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plans, err := h.onboardingService.GetAvailablePlans(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, plans)
}

// ValidateDocument valida documento CPF/CNPJ
func (h *BillingHandler) ValidateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Document string `json:"document"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	valid, err := h.onboardingService.ValidateDocument(ctx, req.Document)
	if err != nil {
		h.writeJSON(w, http.StatusOK, map[string]interface{}{
			"valid":   false,
			"message": err.Error(),
		})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"valid":   valid,
		"message": "Document is valid",
	})
}

// GetSubscriptionStats retorna estatísticas de assinaturas
func (h *BillingHandler) GetSubscriptionStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.subscriptionService.GetSubscriptionStats(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, stats)
}

// GetPaymentStats retorna estatísticas de pagamentos
func (h *BillingHandler) GetPaymentStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verificar se é admin ou usar tenant específico
	var tenantID *uuid.UUID
	if tenantIDStr := r.URL.Query().Get("tenant_id"); tenantIDStr != "" {
		// Verificar se user tem permissão de admin
		if !h.isAdmin(ctx) {
			h.writeError(w, http.StatusForbidden, "Admin access required")
			return
		}
		
		if tid, err := uuid.Parse(tenantIDStr); err == nil {
			tenantID = &tid
		}
	} else {
		// Usar tenant do contexto
		if tid, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
			tenantID = &tid
		}
	}

	stats, err := h.paymentService.GetPaymentStats(ctx, tenantID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, stats)
}

// isAdmin verifica se o usuário é admin
func (h *BillingHandler) isAdmin(ctx context.Context) bool {
	// TODO: Implementar verificação de admin
	return false
}

// writeJSON escreve resposta JSON
func (h *BillingHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError escreve erro JSON
func (h *BillingHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Registrar rotas
func (h *BillingHandler) RegisterRoutes(router *mux.Router) {
	// Assinaturas
	router.HandleFunc("/subscriptions", h.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions/current", h.GetCurrentSubscription).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", h.GetSubscription).Methods("GET")
	router.HandleFunc("/subscriptions/{id}/cancel", h.CancelSubscription).Methods("POST")
	router.HandleFunc("/subscriptions/{id}/change-plan", h.ChangeSubscriptionPlan).Methods("POST")
	router.HandleFunc("/subscriptions/stats", h.GetSubscriptionStats).Methods("GET")

	// Pagamentos
	router.HandleFunc("/payments", h.CreatePayment).Methods("POST")
	router.HandleFunc("/payments", h.GetPaymentsByTenant).Methods("GET")
	router.HandleFunc("/payments/{id}", h.GetPayment).Methods("GET")
	router.HandleFunc("/payments/{id}/refund", h.RefundPayment).Methods("POST")
	router.HandleFunc("/payments/stats", h.GetPaymentStats).Methods("GET")

	// Onboarding
	router.HandleFunc("/onboarding", h.StartOnboarding).Methods("POST")
	router.HandleFunc("/onboarding/status", h.GetOnboardingStatus).Methods("GET")
	router.HandleFunc("/onboarding/plans", h.GetAvailablePlans).Methods("GET")
	router.HandleFunc("/onboarding/validate-document", h.ValidateDocument).Methods("POST")

	// Health check
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

// HealthCheck endpoint de health check
func (h *BillingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "billing-service",
		"version": "1.0.0",
	})
}