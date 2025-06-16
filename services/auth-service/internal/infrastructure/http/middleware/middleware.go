package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/auth-service/internal/infrastructure/config"
	"github.com/direito-lux/auth-service/internal/infrastructure/logging"
)

// Logger middleware para logging de requisições
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithWriter(logger.Sugar().Desugar().Core())
}

// Recovery middleware para recuperação de panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recuperado",
				zap.String("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// RequestID middleware para adicionar ID único à requisição
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		
		// Adicionar ao contexto para logging
		ctx := logging.WithOperation(c.Request.Context(), c.Request.Method+" "+c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
	}
}

// Tenant middleware para extrair tenant ID
func Tenant(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		
		if tenantID != "" {
			c.Set("tenant_id", tenantID)
			
			// Adicionar ao contexto para logging
			ctx := logging.WithTenantID(c.Request.Context(), tenantID)
			c.Request = c.Request.WithContext(ctx)
		}
		
		c.Next()
	}
}

// CORS middleware para Cross-Origin Resource Sharing
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Verificar se origem é permitida
		allowed := false
		for _, allowedOrigin := range cfg.HTTP.CORSAllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", joinStrings(cfg.HTTP.CORSAllowedMethods, ","))
		c.Header("Access-Control-Allow-Headers", joinStrings(cfg.HTTP.CORSAllowedHeaders, ","))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RateLimit middleware para limitação de taxa
func RateLimit(cfg *config.Config) gin.HandlerFunc {
	// Implementação simplificada - em produção usar redis-based rate limiter
	return func(c *gin.Context) {
		// Implementar rate limiting aqui
		c.Next()
	}
}

// joinStrings utilitário para juntar strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	
	return result
}