package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/direito-lux/tenant-service/internal/application"
	"go.uber.org/zap"
)

// QuotaHandler handler HTTP para quotas
type QuotaHandler struct {
	quotaService *application.QuotaService
	logger       *zap.Logger
}

// NewQuotaHandler cria nova instância do handler
func NewQuotaHandler(quotaService *application.QuotaService, logger *zap.Logger) *QuotaHandler {
	return &QuotaHandler{
		quotaService: quotaService,
		logger:       logger,
	}
}

// GetQuotaUsage obtém uso atual de quotas do tenant
func (h *QuotaHandler) GetQuotaUsage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Getting quota usage", zap.String("tenant_id", tenantID))

	usage, err := h.quotaService.GetQuotaUsage(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to get quota usage", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, usage)
}

// GetQuotaLimits obtém limites de quota do tenant
func (h *QuotaHandler) GetQuotaLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Getting quota limits", zap.String("tenant_id", tenantID))

	limits, err := h.quotaService.GetQuotaLimits(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to get quota limits", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, limits)
}

// CheckQuotas verifica todas as quotas do tenant
func (h *QuotaHandler) CheckQuotas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Checking quotas", zap.String("tenant_id", tenantID))

	check, err := h.quotaService.CheckQuotas(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to check quotas", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, check)
}

// IncrementQuota incrementa uma quota específica
func (h *QuotaHandler) IncrementQuota(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Incrementing quota", zap.String("tenant_id", tenantID))

	var req application.IncrementQuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.quotaService.IncrementQuota(r.Context(), tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to increment quota", zap.Error(err))
		statusCode := http.StatusInternalServerError
		
		// Verifica se é erro de quota excedida
		if err.Error() == "quota excedida" {
			statusCode = http.StatusForbidden
		}
		
		writeErrorResponse(w, statusCode, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Quota incremented successfully"})
}

// UpdateStorageUsage atualiza uso de armazenamento
func (h *QuotaHandler) UpdateStorageUsage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Updating storage usage", zap.String("tenant_id", tenantID))

	var req application.UpdateStorageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.quotaService.UpdateStorageUsage(r.Context(), tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to update storage usage", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Storage usage updated successfully"})
}

// ResetDailyQuotas reseta quotas diárias
func (h *QuotaHandler) ResetDailyQuotas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Resetting daily quotas", zap.String("tenant_id", tenantID))

	err := h.quotaService.ResetDailyQuotas(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to reset daily quotas", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Daily quotas reset successfully"})
}

// ResetMonthlyQuotas reseta quotas mensais
func (h *QuotaHandler) ResetMonthlyQuotas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	h.logger.Info("Resetting monthly quotas", zap.String("tenant_id", tenantID))

	err := h.quotaService.ResetMonthlyQuotas(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to reset monthly quotas", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Monthly quotas reset successfully"})
}

// CanIncrementQuota verifica se pode incrementar uma quota específica
func (h *QuotaHandler) CanIncrementQuota(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	quotaType := vars["quota_type"]

	amountStr := r.URL.Query().Get("amount")
	amount := 1 // Padrão: 1
	if amountStr != "" {
		if parsedAmount, err := strconv.Atoi(amountStr); err == nil {
			amount = parsedAmount
		}
	}

	h.logger.Info("Checking if can increment quota", 
		zap.String("tenant_id", tenantID),
		zap.String("quota_type", quotaType),
		zap.Int("amount", amount),
	)

	canIncrement, err := h.quotaService.CanIncrementQuota(r.Context(), tenantID, quotaType, amount)
	if err != nil {
		h.logger.Error("Failed to check quota increment", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"can_increment": canIncrement,
		"quota_type":    quotaType,
		"amount":        amount,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// RegisterRoutes registra as rotas de quota
func (h *QuotaHandler) RegisterRoutes(router *mux.Router) {
	quotaRouter := router.PathPrefix("/tenants/{tenant_id}/quotas").Subrouter()

	quotaRouter.HandleFunc("/usage", h.GetQuotaUsage).Methods("GET")
	quotaRouter.HandleFunc("/limits", h.GetQuotaLimits).Methods("GET")
	quotaRouter.HandleFunc("/check", h.CheckQuotas).Methods("GET")
	quotaRouter.HandleFunc("/increment", h.IncrementQuota).Methods("POST")
	quotaRouter.HandleFunc("/storage", h.UpdateStorageUsage).Methods("PUT")
	quotaRouter.HandleFunc("/reset/daily", h.ResetDailyQuotas).Methods("POST")
	quotaRouter.HandleFunc("/reset/monthly", h.ResetMonthlyQuotas).Methods("POST")
	quotaRouter.HandleFunc("/can-increment/{quota_type}", h.CanIncrementQuota).Methods("GET")
}