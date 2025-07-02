package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/direito-lux/tenant-service/internal/infrastructure/config"
)

// HealthResponse estrutura de resposta do health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HealthCheck handler para health check
func HealthCheck(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
			Service:   cfg.ServiceName,
			Version:   cfg.Version,
		}

		c.JSON(http.StatusOK, response)
	}
}

// ReadinessCheck handler para readiness check
func ReadinessCheck(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Aqui você verificaria dependências (banco, redis, etc.)
		response := HealthResponse{
			Status:    "ready",
			Timestamp: time.Now().UTC(),
			Service:   cfg.ServiceName,
			Version:   cfg.Version,
		}

		c.JSON(http.StatusOK, response)
	}
}

// Ping handler simples para teste
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"timestamp": time.Now().UTC(),
		})
	}
}

// Tenant handlers

// CreateTenant cria novo tenant
func CreateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Tenant criado",
		})
	}
}

// ❌ HANDLER MOCK REMOVIDO - USAR IMPLEMENTAÇÃO REAL NO server.go
// Este handler estava duplicado e 100% hardcoded
// A implementação real está em server.go conectada ao banco PostgreSQL

// GetTenantByDocument busca tenant por documento
func GetTenantByDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		document := c.Query("document")
		if document == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Document parameter is required"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":       "13333333-3333-3333-3333-333333333333",
				"name":     "Costa Santos Advogados 2",
				"document": document,
				"plan":     "professional",
			},
		})
	}
}

// UpdateTenant atualiza tenant
func UpdateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Tenant atualizado",
		})
	}
}

// ActivateTenant ativa tenant
func ActivateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Tenant ativado",
		})
	}
}

// SuspendTenant suspende tenant
func SuspendTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Tenant suspenso",
		})
	}
}

// CancelTenant cancela tenant
func CancelTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Tenant cancelado",
		})
	}
}

// Subscription handlers

// ListSubscriptions lista assinaturas
func ListSubscriptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
	}
}

// CreateSubscription cria assinatura
func CreateSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "Assinatura criada"})
	}
}

// GetSubscription busca assinatura por ID
func GetSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id}})
	}
}

// UpdateSubscription atualiza assinatura
func UpdateSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "Assinatura atualizada"})
	}
}

// CancelSubscription cancela assinatura
func CancelSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "Assinatura cancelada"})
	}
}

// Quota handlers

// GetQuotas busca quotas do tenant
func GetQuotas() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
	}
}

// UpdateQuotas atualiza quotas
func UpdateQuotas() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Quotas atualizadas"})
	}
}

// GetUsage busca uso atual
func GetUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
	}
}

// UpdateUsage atualiza uso
func UpdateUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Uso atualizado"})
	}
}
