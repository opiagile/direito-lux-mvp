package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/search-service/internal/infrastructure/config"
	"github.com/direito-lux/search-service/internal/infrastructure/http/handlers"
	"github.com/direito-lux/search-service/internal/infrastructure/http/middleware"
	"github.com/direito-lux/search-service/internal/infrastructure/metrics"
)

// Server wraps HTTP server
type Server struct {
	server *http.Server
	router *gin.Engine
	config *config.Config
	logger *zap.Logger
}

// NewServer creates a new HTTP server
func NewServer(
	cfg *config.Config,
	logger *zap.Logger,
	metrics *metrics.Metrics,
) *Server {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.Metrics(metrics))

	// Health endpoints
	router.GET("/health", handlers.Health)
	router.GET("/ready", handlers.Ready)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Search endpoints
		v1.POST("/search", handlers.Search)
		v1.POST("/index", handlers.IndexDocument)
		v1.GET("/indices", handlers.ListIndices)
		v1.DELETE("/indices/:index", handlers.DeleteIndex)
		
		// Advanced search
		v1.POST("/search/advanced", handlers.AdvancedSearch)
		v1.POST("/search/aggregate", handlers.SearchWithAggregations)
		v1.GET("/search/suggestions", handlers.SearchSuggestions)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		server: server,
		router: router,
		config: cfg,
		logger: logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.Int("port", s.config.Port))
	
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}