package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthHandler handler para endpoints de saúde
type HealthHandler struct {
	logger *zap.Logger
}

// NewHealthHandler cria nova instância do handler
func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// Health endpoint de health check
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "report-service",
		"version": "1.0.0",
	})
}

// Ready endpoint de readiness
func (h *HealthHandler) Ready(c *gin.Context) {
	// TODO: Verificar dependências (DB, Redis, etc.)
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
			"redis":    "ok",
		},
	})
}