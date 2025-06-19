package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Health retorna status de saúde do serviço
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "search-service",
		"version":   "1.0.0",
		"timestamp": time.Now().Unix(),
	})
}

// Ready retorna se o serviço está pronto
func Ready(c *gin.Context) {
	// TODO: Add actual readiness checks (DB, Elasticsearch, etc)
	c.JSON(http.StatusOK, gin.H{
		"status":      "ready",
		"service":     "search-service",
		"timestamp":   time.Now().Unix(),
		"dependencies": gin.H{
			"database":      "ok",
			"elasticsearch": "ok",
			"redis":         "ok",
		},
	})
}
