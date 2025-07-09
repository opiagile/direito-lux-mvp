package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/direito-lux/notification-service/internal/application/services"
	"github.com/direito-lux/notification-service/internal/domain"
	"github.com/direito-lux/notification-service/internal/infrastructure/config"
	"github.com/direito-lux/notification-service/internal/infrastructure/database"
	"github.com/direito-lux/notification-service/internal/infrastructure/events"
	"github.com/direito-lux/notification-service/internal/infrastructure/http"
	"github.com/direito-lux/notification-service/internal/infrastructure/logging"
	"github.com/direito-lux/notification-service/internal/infrastructure/metrics"
	"github.com/direito-lux/notification-service/internal/infrastructure/repository"
	"github.com/direito-lux/notification-service/internal/infrastructure/tracing"

	"github.com/jmoiron/sqlx"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title notificationservice Service API
// @version 1.0
// @description Direito Lux notificationservice Microservice
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.direitolux.com/support
// @contact.email support@direitolux.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8085
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Erro ao carregar configurações: %v\n", err)
		os.Exit(1)
	}

	// Configurar logger
	logger, err := logging.NewLogger(cfg.LogLevel, cfg.Environment)
	if err != nil {
		fmt.Printf("Erro ao configurar logger: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Erro ao sincronizar logger: %v\n", err)
		}
	}()

	logger.Info("Iniciando notificationservice Service",
		zap.String("version", cfg.Version),
		zap.String("environment", cfg.Environment),
		zap.Int("port", cfg.Port),
	)

	// Usar Fx para injeção de dependência
	app := fx.New(
		fx.Supply(cfg),
		fx.Supply(logger),
		
		// Infraestrutura
		fx.Provide(
			tracing.NewTracer,
			metrics.NewMetrics,
			provideDatabaseConfig,
			database.NewConnection,
			provideDatabaseConnection,
			events.NewEventBus,
		),

		// Repositories
		fx.Provide(
			repository.NewPostgresNotificationRepository,
			repository.NewPostgresTemplateRepository,
			repository.NewPostgresPreferenceRepository,
		),

		// Application Services
		fx.Provide(
			services.NewNotificationService,
			services.NewTemplateService,
			provideNotificationProviders,
			provideNotificationQueue,
		),

		// HTTP Server
		fx.Provide(http.NewServer),
		
		// Lifecycle hooks
		fx.Invoke(registerHooks),
	)

	// Iniciar aplicação
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		logger.Fatal("Erro ao iniciar aplicação", zap.Error(err))
	}

	// Aguardar sinal de parada
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Parando notificationservice Service...")

	// Parar aplicação
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		logger.Error("Erro ao parar aplicação", zap.Error(err))
	}

	logger.Info("notificationservice Service parado com sucesso")
}

// provideDatabaseConfig provides *config.DatabaseConfig from main config
func provideDatabaseConfig(cfg *config.Config) *config.DatabaseConfig {
	return &cfg.Database
}

// provideDatabaseConnection provides *sqlx.DB from database connection
func provideDatabaseConnection(conn *database.Connection) *sqlx.DB {
	return conn.GetDB()
}

// provideNotificationProviders provides notification providers map
func provideNotificationProviders() map[domain.NotificationChannel]domain.NotificationProvider {
	// Mock providers for now - will be replaced with real implementations
	return map[domain.NotificationChannel]domain.NotificationProvider{
		// Will be implemented with real providers
	}
}

// provideNotificationQueue provides notification queue
func provideNotificationQueue() domain.NotificationQueue {
	// Mock queue for now - will be replaced with real implementation
	return nil
}

// registerHooks registra hooks do ciclo de vida da aplicação
func registerHooks(
	lc fx.Lifecycle,
	logger *zap.Logger,
	server *http.Server,
	metrics *metrics.Metrics,
	tracer *tracing.Tracer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Iniciando componentes...")
			
			// Iniciar servidor HTTP em goroutine
			go func() {
				if err := server.Start(); err != nil {
					logger.Error("Erro no servidor HTTP", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Parando componentes...")
			
			// Parar servidor HTTP
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Erro ao parar servidor HTTP", zap.Error(err))
				return err
			}

			// Fechar tracer
			if err := tracer.Close(); err != nil {
				logger.Error("Erro ao fechar tracer", zap.Error(err))
			}

			return nil
		},
	})
}