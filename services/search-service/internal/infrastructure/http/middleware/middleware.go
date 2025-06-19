package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/search-service/internal/infrastructure/logging"
	"github.com/direito-lux/search-service/internal/infrastructure/metrics"
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
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
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

// Metrics middleware para coletar métricas das requisições
func Metrics(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		duration := time.Since(start)
		status := c.Writer.Status()
		
		// Incrementar contador de requisições
		m.HTTPRequests.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			http.StatusText(status),
		).Inc()
		
		// Observar duração
		m.HTTPDuration.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
		).Observe(duration.Seconds())
	}
}
