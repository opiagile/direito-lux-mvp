package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Logger middleware para logging de requisições
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithWriter(os.Stdout)
}

// CORS middleware para Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Em produção, configurar origins específicos
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001", 
			"https://app.direito-lux.com",
		}

		// Verificar se origin é permitido
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed || origin == "" { // Permitir requests sem origin (Postman, etc.)
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Requested-With, X-Tenant-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Recovery middleware para capturar panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// AuthRequired middleware para verificar autenticação
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Verificar formato Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized", 
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// Aqui seria feita a validação do JWT token
		// Por simplicidade, vamos aceitar qualquer token válido por enquanto
		// Em produção, isso deveria validar o token JWT e extrair user_id

		// Simular extração de user_id do token
		// Em produção, isso viria do JWT payload
		userID := uuid.New() // Placeholder - em produção viria do token
		c.Set("user_id", userID)
		c.Set("authenticated", true)

		c.Next()
	}
}

// TenantRequired middleware para verificar tenant
func TenantRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantHeader := c.GetHeader("X-Tenant-ID")
		if tenantHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "X-Tenant-ID header required",
			})
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(tenantHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid tenant ID format",
			})
			c.Abort()
			return
		}

		// Aqui seria feita a validação se o usuário tem acesso ao tenant
		// Por simplicidade, vamos aceitar qualquer tenant válido por enquanto

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// AdminAuthRequired middleware para rotas administrativas
func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Verificar formato Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization format", 
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// Aqui seria feita a validação do token de admin
		// Por simplicidade, vamos aceitar um token específico
		if token != "admin-token-123" { // Em produção seria um JWT com role admin
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Set("admin", true)
		c.Next()
	}
}

// RateLimit middleware para limitação de taxa
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	// Implementação simples baseada em IP
	// Em produção, usaria Redis ou similar
	ipRequests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Limpar requests antigos
		if requests, exists := ipRequests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < time.Minute {
					validRequests = append(validRequests, reqTime)
				}
			}
			ipRequests[clientIP] = validRequests
		}

		// Verificar limite
		if len(ipRequests[clientIP]) >= requestsPerMinute {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		// Adicionar request atual
		ipRequests[clientIP] = append(ipRequests[clientIP], now)
		c.Next()
	}
}

// RequestTimeout middleware para timeout de requisições
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Configurar timeout no contexto
		ctx := c.Request.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// APIKeyAuth middleware para autenticação via API key
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "API key required",
			})
			c.Abort()
			return
		}

		// Validar API key
		// Em produção, isso seria validado contra base de dados
		validAPIKeys := []string{
			"notification-api-key-123",
			"external-service-key-456",
		}

		valid := false
		for _, validKey := range validAPIKeys {
			if apiKey == validKey {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Set("api_authenticated", true)
		c.Next()
	}
}
