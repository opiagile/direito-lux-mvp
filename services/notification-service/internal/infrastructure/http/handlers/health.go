package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handler para health checks
type HealthHandler struct{}

// NewHealthHandler cria novo health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health retorna status de saúde do serviço
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
		"mode":   "full",
	})
}

// Ready retorna se o serviço está pronto
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"time":   time.Now().Unix(),
	})
}
