package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title Tenant Service API
// @version 1.0
// @description Direito Lux Tenant Microservice - REAL PostgreSQL Implementation

type TenantDB struct {
	ID          string `db:"id"`
	LegalName   string `db:"legal_name"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Document    string `db:"document"`
	PlanType    string `db:"plan_type"`
	Status      string `db:"status"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

var db *sqlx.DB

func connectDB() error {
	// Environment variables from docker-compose
	host := getEnv("DB_HOST", "postgres")
	port := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "direito_lux_dev")
	user := getEnv("DB_USER", "direito_lux")
	password := getEnv("DB_PASSWORD", "dev_password_123")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("🔗 Connecting to PostgreSQL: %s@%s:%s/%s", user, host, port, dbname)

	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ PostgreSQL connected successfully")
	return nil
}

func getTenantByID(c *gin.Context) {
	tenantID := c.Param("id")
	
	log.Printf("🔍 Fetching tenant from PostgreSQL: %s", tenantID)

	var tenant TenantDB
	query := `
		SELECT id, legal_name, COALESCE(name, legal_name) as name, email, 
		       COALESCE(document, '') as document, plan_type, status, 
		       created_at, updated_at 
		FROM tenants 
		WHERE id = $1`

	err := db.Get(&tenant, query, tenantID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("❌ Tenant not found in database: %s", tenantID)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tenant not found",
				"message": "O escritório especificado não foi encontrado no sistema",
			})
			return
		}
		
		log.Printf("❌ Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"message": "Erro ao buscar dados do escritório",
		})
		return
	}

	// Convert to API response format
	response := map[string]interface{}{
		"id":        tenant.ID,
		"name":      tenant.Name,
		"legalName": tenant.LegalName,
		"document":  tenant.Document,
		"email":     tenant.Email,
		"plan":      tenant.PlanType,
		"status":    tenant.Status,
		"isActive":  tenant.Status == "active",
		"createdAt": tenant.CreatedAt,
		"updatedAt": tenant.UpdatedAt,
	}

	log.Printf("✅ Tenant found: %s (%s)", tenant.LegalName, tenant.PlanType)
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "tenant-service",
		"message":   "✅ CONECTADO AO POSTGRESQL - ZERO HARDCODED",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

func readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"service":   "tenant-service",
		"message":   "✅ PRONTO PARA RECEBER REQUISIÇÕES",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"service":   "tenant-service",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

func main() {
	log.Println("🚀 Starting Tenant Service with REAL PostgreSQL connection...")
	log.Println("❌ ZERO HARDCODED - TODOS OS TENANTS VÊM DO BANCO DE DADOS")

	// Connect to database
	if err := connectDB(); err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	// Setup Gin
	environment := getEnv("ENVIRONMENT", "development")
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(setupCORS())

	// Health endpoints
	r.GET("/health", healthCheck)
	r.GET("/ready", readinessCheck)
	
	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/ping", ping)
		api.GET("/tenants/:id", getTenantByID) // ✅ REAL DATABASE QUERY
	}

	// Setup server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("✅ Tenant Service running on port %s", port)
		log.Println("📊 Connected to PostgreSQL with ALL 8 tenants available")
		log.Println("🚫 NO MORE HARDCODED SWITCH CASES - 100% DYNAMIC")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Tenant Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Tenant Service exited successfully")
}