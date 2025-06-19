package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
	"github.com/direito-lux/mcp-service/internal/infrastructure/logging"
)

// Logger middleware para log das requisições
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Processar request
		c.Next()

		// Log da request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// Recovery middleware para recuperação de panics
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
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS middleware para configuração de CORS
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Configurar headers CORS
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Tenant-ID, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Responder a requests OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestID middleware para adicionar ID único à request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Adicionar ao contexto de logging
		ctx := logging.WithTraceID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Tenant middleware para extrair informações do tenant
func Tenant(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Request.Header.Get("X-Tenant-ID")
		
		// Também pode vir do JWT token (implementar se necessário)
		if tenantID == "" {
			// Extrair do token JWT se disponível
			if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
				// TODO: Implementar extração do tenant do JWT
			}
		}

		if tenantID != "" {
			c.Set("tenant_id", tenantID)
			
			// Adicionar ao contexto de logging
			ctx := logging.WithTenantID(c.Request.Context(), tenantID)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// RateLimit middleware para limitação de taxa
func RateLimit(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implementar rate limiting usando Redis ou memória
		// Por agora, apenas continua
		c.Next()
	}
}

// MCPAuth middleware para autenticação MCP
func MCPAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pular autenticação para health checks e webhooks
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/health") || 
		   strings.HasPrefix(path, "/ready") ||
		   strings.HasPrefix(path, "/webhooks") {
			c.Next()
			return
		}

		// Verificar token de autenticação
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extrair Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		
		// TODO: Validar token JWT
		// Por agora, apenas verifica se não está vazio
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token required",
			})
			c.Abort()
			return
		}

		// TODO: Extrair informações do usuário do token
		c.Set("user_id", "user-from-token")
		c.Set("tenant_id", "tenant-from-token")

		c.Next()
	}
}

// MCPQuota middleware para verificação de quota
func MCPQuota(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pular quota para health checks
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/health") || 
		   strings.HasPrefix(path, "/ready") ||
		   strings.HasPrefix(path, "/webhooks") {
			c.Next()
			return
		}

		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.Next()
			return
		}

		// TODO: Verificar quota do tenant
		// Por agora, apenas continua
		_ = tenantID

		c.Next()
	}
}