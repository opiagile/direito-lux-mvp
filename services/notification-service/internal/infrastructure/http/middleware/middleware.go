package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/infrastructure/config"
	"github.com/direito-lux/notification-service/internal/infrastructure/logging"
)

// Logger middleware para logging de requisições
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithWriter(os.Stdout)
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

// RequestID middleware para adicionar ID único às requisições
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header("X-Request-ID", requestID)
		
		// Adicionar ao contexto
		ctx := logging.WithOperation(c.Request.Context(), c.Request.Method+" "+c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
	}
}

// CORS middleware
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {		
		// Em desenvolvimento, permitir qualquer origem
		allowedOrigin := "*"
		if !cfg.IsDevelopment() {
			// Em produção, verificar origens permitidas
			allowedOrigin = "https://app.direitolux.com"
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
