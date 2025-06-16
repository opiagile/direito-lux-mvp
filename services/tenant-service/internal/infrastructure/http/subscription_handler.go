package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/direito-lux/tenant-service/internal/application"
	"github.com/direito-lux/tenant-service/internal/domain"
	"go.uber.org/zap"
)

// SubscriptionHandler handler HTTP para assinaturas
type SubscriptionHandler struct {
	subscriptionService *application.SubscriptionService
	logger              *zap.Logger
}

// NewSubscriptionHandler cria nova instância do handler
func NewSubscriptionHandler(subscriptionService *application.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// CreateSubscription cria uma nova assinatura
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Creating subscription", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var req application.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscription, err := h.subscriptionService.CreateSubscription(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, subscription)
}

// GetSubscription obtém assinatura por ID
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID := vars["id"]

	h.logger.Info("Getting subscription", zap.String("subscription_id", subscriptionID))

	subscription, err := h.subscriptionService.GetSubscription(r.Context(), subscriptionID)
	if err != nil {
		h.logger.Error("Failed to get subscription", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, subscription)
}

// GetSubscriptionByTenant obtém assinatura ativa do tenant
func (h *SubscriptionHandler) GetSubscriptionByTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Getting subscription by tenant", zap.String("tenant_id", tenantID))

	subscription, err := h.subscriptionService.GetSubscriptionByTenant(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to get subscription by tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, subscription)
}

// ActivateSubscription ativa uma assinatura
func (h *SubscriptionHandler) ActivateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID := vars["id"]

	h.logger.Info("Activating subscription", zap.String("subscription_id", subscriptionID))

	err := h.subscriptionService.ActivateSubscription(r.Context(), subscriptionID)
	if err != nil {
		h.logger.Error("Failed to activate subscription", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Subscription activated successfully"})
}

// CancelSubscription cancela uma assinatura
func (h *SubscriptionHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID := vars["id"]

	h.logger.Info("Canceling subscription", zap.String("subscription_id", subscriptionID))

	var reqBody struct {
		Reason      string `json:"reason"`
		Immediately bool   `json:"immediately"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.subscriptionService.CancelSubscription(r.Context(), subscriptionID, reqBody.Reason, reqBody.Immediately)
	if err != nil {
		h.logger.Error("Failed to cancel subscription", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Subscription canceled successfully"})
}

// ReactivateSubscription reativa uma assinatura cancelada
func (h *SubscriptionHandler) ReactivateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID := vars["id"]

	h.logger.Info("Reactivating subscription", zap.String("subscription_id", subscriptionID))

	err := h.subscriptionService.ReactivateSubscription(r.Context(), subscriptionID)
	if err != nil {
		h.logger.Error("Failed to reactivate subscription", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Subscription reactivated successfully"})
}

// ChangePlan altera o plano de uma assinatura
func (h *SubscriptionHandler) ChangePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID := vars["id"]

	h.logger.Info("Changing subscription plan", zap.String("subscription_id", subscriptionID))

	var req application.ChangePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscription, err := h.subscriptionService.ChangePlan(r.Context(), subscriptionID, &req)
	if err != nil {
		h.logger.Error("Failed to change subscription plan", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, subscription)
}

// RenewSubscription renova uma assinatura
func (h *SubscriptionHandler) RenewSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID := vars["id"]

	h.logger.Info("Renewing subscription", zap.String("subscription_id", subscriptionID))

	var reqBody struct {
		BillingInterval string `json:"billing_interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Converte string para BillingInterval
	var billingInterval domain.BillingInterval
	switch reqBody.BillingInterval {
	case "monthly":
		billingInterval = domain.BillingMonthly
	case "yearly":
		billingInterval = domain.BillingYearly
	default:
		writeErrorResponse(w, http.StatusBadRequest, "Invalid billing interval")
		return
	}

	err := h.subscriptionService.RenewSubscription(r.Context(), subscriptionID, billingInterval)
	if err != nil {
		h.logger.Error("Failed to renew subscription", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Subscription renewed successfully"})
}

// GetExpiringSubscriptions obtém assinaturas que vão expirar
func (h *SubscriptionHandler) GetExpiringSubscriptions(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	days := 7 // Padrão: 7 dias
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil {
			days = parsedDays
		}
	}

	h.logger.Info("Getting expiring subscriptions", zap.Int("days", days))

	subscriptions, err := h.subscriptionService.GetExpiringSubscriptions(r.Context(), days)
	if err != nil {
		h.logger.Error("Failed to get expiring subscriptions", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, subscriptions)
}

// ListPlans lista todos os planos disponíveis
func (h *SubscriptionHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	h.logger.Info("Listing plans", zap.Bool("active_only", activeOnly))

	plans, err := h.subscriptionService.ListPlans(r.Context(), activeOnly)
	if err != nil {
		h.logger.Error("Failed to list plans", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, plans)
}

// GetPlan obtém plano por ID
func (h *SubscriptionHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	planID := vars["plan_id"]

	h.logger.Info("Getting plan", zap.String("plan_id", planID))

	plan, err := h.subscriptionService.GetPlan(r.Context(), planID)
	if err != nil {
		h.logger.Error("Failed to get plan", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, plan)
}

// RegisterRoutes registra as rotas de assinatura
func (h *SubscriptionHandler) RegisterRoutes(router *mux.Router) {
	// Rotas de assinaturas
	subscriptionRouter := router.PathPrefix("/subscriptions").Subrouter()
	subscriptionRouter.HandleFunc("", h.CreateSubscription).Methods("POST")
	subscriptionRouter.HandleFunc("/expiring", h.GetExpiringSubscriptions).Methods("GET")
	subscriptionRouter.HandleFunc("/{id}", h.GetSubscription).Methods("GET")
	subscriptionRouter.HandleFunc("/{id}/activate", h.ActivateSubscription).Methods("POST")
	subscriptionRouter.HandleFunc("/{id}/cancel", h.CancelSubscription).Methods("POST")
	subscriptionRouter.HandleFunc("/{id}/reactivate", h.ReactivateSubscription).Methods("POST")
	subscriptionRouter.HandleFunc("/{id}/change-plan", h.ChangePlan).Methods("POST")
	subscriptionRouter.HandleFunc("/{id}/renew", h.RenewSubscription).Methods("POST")

	// Rotas para buscar assinatura por tenant
	router.HandleFunc("/tenants/{tenant_id}/subscription", h.GetSubscriptionByTenant).Methods("GET")

	// Rotas de planos
	planRouter := router.PathPrefix("/plans").Subrouter()
	planRouter.HandleFunc("", h.ListPlans).Methods("GET")
	planRouter.HandleFunc("/{plan_id}", h.GetPlan).Methods("GET")
}