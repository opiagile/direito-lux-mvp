package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/direito-lux/billing-service/internal/infrastructure/config"
	"github.com/direito-lux/billing-service/internal/infrastructure/database"
	"github.com/direito-lux/billing-service/internal/infrastructure/events"
	"github.com/direito-lux/billing-service/internal/infrastructure/http"
	"github.com/direito-lux/billing-service/internal/infrastructure/logging"
	"github.com/direito-lux/billing-service/internal/infrastructure/metrics"
	"github.com/direito-lux/billing-service/internal/infrastructure/tracing"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title Billing Service API
// @version 1.0
// @description Direito Lux Billing Microservice
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.direitolux.com/support
// @contact.email support@direitolux.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8089
// @BasePath /billing

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

	logger.Info("Iniciando Billing Service",
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
			database.NewConnection,
			messaging.NewRabbitMQConnection,
			events.NewEventBus,
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

	logger.Info("Parando Billing Service...")

	// Parar aplicação
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		logger.Error("Erro ao parar aplicação", zap.Error(err))
	}

	logger.Info("Billing Service parado com sucesso")
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