package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/direito-lux/tenant-service/internal/application"
	"go.uber.org/zap"
)

// TenantHandler handler HTTP para tenants
type TenantHandler struct {
	tenantService *application.TenantService
	logger        *zap.Logger
}

// NewTenantHandler cria nova instância do handler
func NewTenantHandler(tenantService *application.TenantService, logger *zap.Logger) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		logger:        logger,
	}
}

// CreateTenant cria um novo tenant
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Creating tenant", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var req application.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tenant, err := h.tenantService.CreateTenant(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, tenant)
}

// GetTenant obtém tenant por ID
func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	h.logger.Info("Getting tenant", zap.String("tenant_id", tenantID))

	tenant, err := h.tenantService.GetTenant(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to get tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, tenant)
}

// GetTenantByDocument obtém tenant por documento
func (h *TenantHandler) GetTenantByDocument(w http.ResponseWriter, r *http.Request) {
	document := r.URL.Query().Get("document")
	if document == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Document parameter is required")
		return
	}

	h.logger.Info("Getting tenant by document", zap.String("document", document))

	tenant, err := h.tenantService.GetTenantByDocument(r.Context(), document)
	if err != nil {
		h.logger.Error("Failed to get tenant by document", zap.Error(err))
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, tenant)
}

// GetTenantsByOwner obtém tenants por proprietário
func (h *TenantHandler) GetTenantsByOwner(w http.ResponseWriter, r *http.Request) {
	ownerUserID := r.URL.Query().Get("owner_user_id")
	if ownerUserID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "owner_user_id parameter is required")
		return
	}

	h.logger.Info("Getting tenants by owner", zap.String("owner_user_id", ownerUserID))

	tenants, err := h.tenantService.GetTenantsByOwner(r.Context(), ownerUserID)
	if err != nil {
		h.logger.Error("Failed to get tenants by owner", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, tenants)
}

// UpdateTenant atualiza dados do tenant
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	h.logger.Info("Updating tenant", zap.String("tenant_id", tenantID))

	var req application.UpdateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Pega user ID do header ou token para auditoria
	updatedBy := r.Header.Get("X-User-ID")
	if updatedBy == "" {
		updatedBy = "system"
	}

	tenant, err := h.tenantService.UpdateTenant(r.Context(), tenantID, &req, updatedBy)
	if err != nil {
		h.logger.Error("Failed to update tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, tenant)
}

// ActivateTenant ativa um tenant
func (h *TenantHandler) ActivateTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	h.logger.Info("Activating tenant", zap.String("tenant_id", tenantID))

	// Pega user ID do header ou token para auditoria
	activatedBy := r.Header.Get("X-User-ID")
	if activatedBy == "" {
		activatedBy = "system"
	}

	err := h.tenantService.ActivateTenant(r.Context(), tenantID, activatedBy)
	if err != nil {
		h.logger.Error("Failed to activate tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Tenant activated successfully"})
}

// SuspendTenant suspende um tenant
func (h *TenantHandler) SuspendTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	h.logger.Info("Suspending tenant", zap.String("tenant_id", tenantID))

	var reqBody struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Pega user ID do header ou token para auditoria
	suspendedBy := r.Header.Get("X-User-ID")
	if suspendedBy == "" {
		suspendedBy = "system"
	}

	err := h.tenantService.SuspendTenant(r.Context(), tenantID, reqBody.Reason, suspendedBy)
	if err != nil {
		h.logger.Error("Failed to suspend tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Tenant suspended successfully"})
}

// CancelTenant cancela um tenant
func (h *TenantHandler) CancelTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	h.logger.Info("Canceling tenant", zap.String("tenant_id", tenantID))

	var reqBody struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Pega user ID do header ou token para auditoria
	canceledBy := r.Header.Get("X-User-ID")
	if canceledBy == "" {
		canceledBy = "system"
	}

	err := h.tenantService.CancelTenant(r.Context(), tenantID, reqBody.Reason, canceledBy)
	if err != nil {
		h.logger.Error("Failed to cancel tenant", zap.Error(err))
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Tenant canceled successfully"})
}

// RegisterRoutes registra as rotas do tenant
func (h *TenantHandler) RegisterRoutes(router *mux.Router) {
	tenantRouter := router.PathPrefix("/tenants").Subrouter()

	tenantRouter.HandleFunc("", h.CreateTenant).Methods("POST")
	tenantRouter.HandleFunc("/search", h.GetTenantByDocument).Methods("GET").Queries("document", "{document}")
	tenantRouter.HandleFunc("/search", h.GetTenantsByOwner).Methods("GET").Queries("owner_user_id", "{owner_user_id}")
	tenantRouter.HandleFunc("/{id}", h.GetTenant).Methods("GET")
	tenantRouter.HandleFunc("/{id}", h.UpdateTenant).Methods("PUT")
	tenantRouter.HandleFunc("/{id}/activate", h.ActivateTenant).Methods("POST")
	tenantRouter.HandleFunc("/{id}/suspend", h.SuspendTenant).Methods("POST")
	tenantRouter.HandleFunc("/{id}/cancel", h.CancelTenant).Methods("POST")
}

// writeJSONResponse escreve resposta JSON
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse escreve resposta de erro
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}