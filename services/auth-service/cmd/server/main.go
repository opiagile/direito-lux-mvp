package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/direito-lux/auth-service/internal/application"
	"github.com/direito-lux/auth-service/internal/domain"
	"github.com/direito-lux/auth-service/internal/infrastructure/config"
	"github.com/direito-lux/auth-service/internal/infrastructure/database"
	"github.com/direito-lux/auth-service/internal/infrastructure/events"
	"github.com/direito-lux/auth-service/internal/infrastructure/http"
	"github.com/direito-lux/auth-service/internal/infrastructure/logging"
	"github.com/direito-lux/auth-service/internal/infrastructure/metrics"
	"github.com/direito-lux/auth-service/internal/infrastructure/tracing"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// eventBusAdapter adapta events.EventBus para application.EventBus
type eventBusAdapter struct {
	eventBus events.EventBus
}

// domainEventAdapter adapta events.DomainEvent para application.DomainEvent
type domainEventAdapter struct {
	event events.DomainEvent
}

func (a *domainEventAdapter) EventType() string {
	return a.event.GetEventType()
}

func (a *domainEventAdapter) AggregateID() string {
	return a.event.GetAggregateID()
}

func (a *domainEventAdapter) Payload() ([]byte, error) {
	// Simple JSON encoding for now
	return []byte(fmt.Sprintf(`{"event_type":"%s","aggregate_id":"%s"}`, a.event.GetEventType(), a.event.GetAggregateID())), nil
}

func (adapter *eventBusAdapter) Publish(ctx context.Context, event application.DomainEvent) error {
	// For now, we'll skip the event publishing to avoid interface complexity
	// TODO: Implement proper adapter when events are needed
	return nil
}

// @title Auth Service API
// @version 1.0
// @description Direito Lux Authentication & Authorization Microservice
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.direitolux.com/support
// @contact.email support@direitolux.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
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

	logger.Info("Iniciando Auth Service",
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
			func(conn *database.Connection) *sqlx.DB {
				return conn.GetDB()
			},
			events.NewEventBus,
		),
		
		// Repositories (cast to interfaces)
		fx.Provide(
			fx.Annotate(
				database.NewUserRepository,
				fx.As(new(domain.UserRepository)),
			),
			fx.Annotate(
				database.NewSessionRepository,
				fx.As(new(domain.SessionRepository)),
			),
			fx.Annotate(
				database.NewRefreshTokenRepository,
				fx.As(new(domain.RefreshTokenRepository)),
			),
			fx.Annotate(
				database.NewLoginAttemptRepository,
				fx.As(new(domain.LoginAttemptRepository)),
			),
			fx.Annotate(
				database.NewPasswordResetTokenRepository,
				fx.As(new(domain.PasswordResetTokenRepository)),
			),
		),
		
		// Services with configuration injection
		fx.Provide(
			func(
				userRepo domain.UserRepository,
				sessionRepo domain.SessionRepository,
				refreshTokenRepo domain.RefreshTokenRepository,
				loginAttemptRepo domain.LoginAttemptRepository,
				passwordResetTokenRepo domain.PasswordResetTokenRepository,
				eventBus events.EventBus,
				cfg *config.Config,
			) *application.AuthService {
				return application.NewAuthService(
					userRepo,
					sessionRepo,
					refreshTokenRepo,
					loginAttemptRepo,
					passwordResetTokenRepo,
					eventBus,
					cfg.JWT.Secret,
					cfg.JWT.ExpiryHours,
					cfg.JWT.RefreshExpiryDays,
				)
			},
			fx.Annotate(
				func(eventBus events.EventBus) application.EventBus {
					return &eventBusAdapter{eventBus: eventBus}
				},
				fx.As(new(application.EventBus)),
			),
			func(userRepo domain.UserRepository, eventBus application.EventBus) *application.UserService {
				return application.NewUserService(userRepo, eventBus)
			},
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

	logger.Info("Parando Auth Service...")

	// Parar aplicação
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		logger.Error("Erro ao parar aplicação", zap.Error(err))
	}

	logger.Info("Auth Service parado com sucesso")
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