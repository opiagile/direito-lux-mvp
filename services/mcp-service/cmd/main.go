package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
	"github.com/direito-lux/mcp-service/internal/infrastructure/database"
	"github.com/direito-lux/mcp-service/internal/infrastructure/logging"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
	"github.com/direito-lux/mcp-service/internal/infrastructure/events"
	"github.com/direito-lux/mcp-service/internal/infrastructure/messaging"
	"github.com/direito-lux/mcp-service/internal/infrastructure/tracing"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http"
	"github.com/direito-lux/mcp-service/internal/application/services"
)

func main() {
	// Criar contexto para graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Criar aplicação com Fx
	app := fx.New(
		// Configuração
		fx.Provide(config.New),

		// Logging
		fx.Provide(func(cfg *config.Config) (*zap.Logger, error) {
			return logging.NewLogger(cfg.LogLevel, cfg.Environment)
		}),

		// Database
		fx.Provide(func(cfg *config.Config, logger *zap.Logger) (*database.Connection, error) {
			return database.NewConnection(cfg, logger)
		}),

		// Messaging (RabbitMQ)
		fx.Provide(func(cfg *config.Config, logger *zap.Logger) (*messaging.RabbitMQConnection, error) {
			if cfg.Environment == "test" {
				// Para testes, retornar nil sem erro
				return nil, nil
			}
			return messaging.NewRabbitMQConnection(cfg, logger)
		}),

		// Metrics
		fx.Provide(func(cfg *config.Config, logger *zap.Logger) (*metrics.Metrics, error) {
			return metrics.NewMetrics(cfg, logger)
		}),

		// Tracing
		fx.Provide(func(cfg *config.Config, logger *zap.Logger) (*tracing.Tracer, error) {
			return tracing.NewTracer(cfg, logger)
		}),

		// Event Bus
		fx.Provide(func(
			logger *zap.Logger,
			messaging *messaging.RabbitMQConnection,
			metrics *metrics.Metrics,
		) *events.EventBus {
			return events.NewEventBus(logger, messaging, metrics)
		}),

		// Application Services
		fx.Provide(func(
			logger *zap.Logger,
			metrics *metrics.Metrics,
			eventBus *events.EventBus,
		) *services.SessionService {
			return services.NewSessionService(logger, metrics, eventBus)
		}),

		fx.Provide(func(
			logger *zap.Logger,
			metrics *metrics.Metrics,
			eventBus *events.EventBus,
		) *services.ToolService {
			return services.NewToolService(logger, metrics, eventBus)
		}),

		// HTTP Server
		fx.Provide(func(
			cfg *config.Config,
			logger *zap.Logger,
			sessionService *services.SessionService,
			toolService *services.ToolService,
		) (*http.Server, error) {
			return http.NewServer(cfg, logger, sessionService, toolService)
		}),

		// Lifecycle hooks
		fx.Invoke(func(
			lifecycle fx.Lifecycle,
			cfg *config.Config,
			logger *zap.Logger,
			server *http.Server,
			db *database.Connection,
			tracer *tracing.Tracer,
		) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("🚀 Iniciando MCP Service",
						zap.String("version", cfg.ServiceName),
						zap.String("environment", cfg.Environment),
						zap.Int("port", cfg.Port),
					)

					// Testar conexões
					if err := testConnections(ctx, logger, db); err != nil {
						logger.Error("❌ Falha nos testes de conexão", zap.Error(err))
						return err
					}

					// Iniciar servidor HTTP em goroutine
					go func() {
						if err := server.Start(ctx); err != nil {
							logger.Error("❌ Erro ao iniciar servidor HTTP", zap.Error(err))
						}
					}()

					logger.Info("✅ MCP Service iniciado com sucesso",
						zap.String("http_address", fmt.Sprintf(":%d", cfg.Port)),
						zap.String("metrics_address", fmt.Sprintf(":%d", cfg.Metrics.Port)),
					)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("🛑 Parando MCP Service...")

					// Parar servidor HTTP
					if err := server.Stop(ctx); err != nil {
						logger.Error("Erro ao parar servidor HTTP", zap.Error(err))
					}

					// Fechar tracer
					if err := tracer.Close(); err != nil {
						logger.Error("Erro ao fechar tracer", zap.Error(err))
					}

					// Fechar conexão com banco
					if err := db.Close(); err != nil {
						logger.Error("Erro ao fechar conexão com banco", zap.Error(err))
					}

					logger.Info("✅ MCP Service parado com sucesso")
					return nil
				},
			})
		}),
	)

	// Executar aplicação
	if err := app.Start(ctx); err != nil {
		log.Fatalf("❌ Falha ao iniciar aplicação: %v", err)
	}

	// Aguardar sinal de shutdown
	<-ctx.Done()

	// Graceful shutdown
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatalf("❌ Falha no shutdown graceful: %v", err)
	}

	fmt.Println("👋 MCP Service finalizado")
}

// testConnections testa as conexões com os serviços externos
func testConnections(ctx context.Context, logger *zap.Logger, db *database.Connection) error {
	logger.Info("🔍 Testando conexões...")

	// Testar conexão com banco de dados
	if err := db.Health(ctx); err != nil {
		return fmt.Errorf("falha na conexão com PostgreSQL: %w", err)
	}
	logger.Info("✅ PostgreSQL conectado")

	// TODO: Adicionar testes para Redis e RabbitMQ quando necessário

	return nil
}