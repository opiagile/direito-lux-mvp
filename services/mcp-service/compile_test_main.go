package main

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
	"github.com/direito-lux/mcp-service/internal/infrastructure/database"
	"github.com/direito-lux/mcp-service/internal/infrastructure/logging"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
	"github.com/direito-lux/mcp-service/internal/infrastructure/events"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http"
	"github.com/direito-lux/mcp-service/internal/application/services"
)

func main() {
	app := fx.New(
		// Configuration
		fx.Provide(config.New),

		// Logging
		fx.Provide(logging.NewLogger),

		// Database
		fx.Provide(database.NewConnection),

		// Metrics
		fx.Provide(metrics.NewMetrics),

		// Event Bus (without messaging for now)
		fx.Provide(func(logger *zap.Logger, metrics *metrics.Metrics) *events.EventBus {
			return events.NewEventBus(logger, nil, metrics) // nil messaging for testing
		}),

		// Services
		fx.Provide(services.NewSessionService),
		fx.Provide(services.NewToolService),

		// HTTP Server
		fx.Provide(http.NewServer),

		// Lifecycle
		fx.Invoke(func(server *http.Server) {
			// Start server in lifecycle
		}),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatal("Erro ao iniciar aplicação:", err)
	}

	fmt.Println("✅ MCP Service - Teste de compilação executado com sucesso!")

	if err := app.Stop(ctx); err != nil {
		log.Fatal("Erro ao parar aplicação:", err)
	}
}