package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ProcessStats estrutura para estatísticas de processos
type ProcessStats struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Paused    int `json:"paused"`
	Archived  int `json:"archived"`
	ThisMonth int `json:"this_month"`
}

func main() {
	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Tenant-ID")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "process-service-temp",
			"timestamp": time.Now().UTC(),
		})
	})

	// Endpoint crítico: /api/v1/processes/stats
	r.GET("/api/v1/processes/stats", func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		// Dados temporários diferentes para cada tenant
		var stats ProcessStats
		
		switch tenantID {
		case "11111111-1111-1111-1111-111111111111": // Silva & Associados
			stats = ProcessStats{
				Total:     45,
				Active:    38,
				Paused:    5,
				Archived:  2,
				ThisMonth: 12,
			}
		case "22222222-2222-2222-2222-222222222222": // Costa Santos
			stats = ProcessStats{
				Total:     32,
				Active:    28,
				Paused:    3,
				Archived:  1,
				ThisMonth: 8,
			}
		case "33333333-3333-3333-3333-333333333333": // Barros Empresa
			stats = ProcessStats{
				Total:     67,
				Active:    58,
				Paused:    7,
				Archived:  2,
				ThisMonth: 15,
			}
		case "44444444-4444-4444-4444-444444444444": // Lima Advogados
			stats = ProcessStats{
				Total:     29,
				Active:    25,
				Paused:    2,
				Archived:  2,
				ThisMonth: 7,
			}
		case "55555555-5555-5555-5555-555555555555": // Pereira Advocacia
			stats = ProcessStats{
				Total:     18,
				Active:    16,
				Paused:    1,
				Archived:  1,
				ThisMonth: 4,
			}
		case "66666666-6666-6666-6666-666666666666": // Rodrigues Global
			stats = ProcessStats{
				Total:     89,
				Active:    76,
				Paused:    8,
				Archived:  5,
				ThisMonth: 22,
			}
		case "77777777-7777-7777-7777-777777777777": // Oliveira Partners
			stats = ProcessStats{
				Total:     156,
				Active:    142,
				Paused:    10,
				Archived:  4,
				ThisMonth: 35,
			}
		case "88888888-8888-8888-8888-888888888888": // Machado Advogados
			stats = ProcessStats{
				Total:     78,
				Active:    71,
				Paused:    5,
				Archived:  2,
				ThisMonth: 18,
			}
		default:
			stats = ProcessStats{
				Total:     15,
				Active:    12,
				Paused:    2,
				Archived:  1,
				ThisMonth: 3,
			}
		}

		log.Printf("📊 Estatísticas retornadas para tenant %s: %+v", tenantID, stats)

		c.JSON(200, gin.H{
			"data":      stats,
			"timestamp": time.Now().UTC(),
		})
	})

	// Endpoint básico para listar processos
	r.GET("/api/v1/processes", func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		// Processos de exemplo
		processes := []map[string]interface{}{
			{
				"id":                "proc-1",
				"number":            "5001234-12.2024.8.26.0100",
				"title":             "Ação de Cobrança - Cliente XYZ",
				"court":             "TJSP",
				"status":            "active",
				"monitoring_enabled": true,
				"created_at":        "2024-12-15T10:30:00Z",
			},
			{
				"id":                "proc-2",
				"number":            "5001235-45.2024.8.26.0224",
				"title":             "Ação Trabalhista - Rescisão Indireta",
				"court":             "TJSP",
				"status":            "active",
				"monitoring_enabled": true,
				"created_at":        "2024-12-20T14:15:00Z",
			},
		}

		c.JSON(200, gin.H{
			"data":  processes,
			"total": len(processes),
		})
	})

	// Endpoint para criar processo
	r.POST("/api/v1/processes", func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Dados inválidos"})
			return
		}

		c.JSON(201, gin.H{
			"data": gin.H{
				"id":      fmt.Sprintf("proc-%d", time.Now().Unix()),
				"message": "Processo criado com sucesso",
			},
		})
	})

	// Outros endpoints básicos
	r.GET("/api/v1/processes/:id", func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(400, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			return
		}

		processID := c.Param("id")
		c.JSON(200, gin.H{
			"data": gin.H{
				"id":     processID,
				"number": "5001234-12.2024.8.26.0100",
				"title":  "Processo encontrado",
				"status": "active",
			},
		})
	})

	r.PUT("/api/v1/processes/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": gin.H{
				"message": "Processo atualizado com sucesso",
			},
		})
	})

	r.DELETE("/api/v1/processes/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": gin.H{
				"message": "Processo excluído com sucesso",
			},
		})
	})

	// Log de início
	log.Printf("🚀 Process Stats Server iniciando na porta 8083")
	log.Printf("📊 Endpoint crítico: GET /api/v1/processes/stats")
	log.Printf("🔍 Health check: GET /health")
	log.Printf("📋 Processos: GET /api/v1/processes")

	// Iniciar servidor na porta 8083
	if err := r.Run(":8083"); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}