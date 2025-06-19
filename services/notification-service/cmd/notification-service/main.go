package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/application/services"
	"github.com/direito-lux/notification-service/internal/domain"
	"github.com/direito-lux/notification-service/internal/infrastructure/config"
	httpinfra "github.com/direito-lux/notification-service/internal/infrastructure/http"
	"github.com/direito-lux/notification-service/internal/infrastructure/providers"
	"github.com/direito-lux/notification-service/internal/infrastructure/repository"
)

func main() {
	app := fx.New(
		// Configuração
		fx.Provide(config.NewConfig),
		
		// Logger
		fx.Provide(NewLogger),
		
		// Database
		fx.Provide(NewDatabase),
		
		// Repositories
		fx.Provide(NewNotificationRepository),
		fx.Provide(NewTemplateRepository),
		fx.Provide(NewPreferenceRepository),
		
		// Providers
		fx.Provide(NewEmailProvider),
		fx.Provide(NewWhatsAppProvider),
		fx.Provide(NewTelegramProvider),
		fx.Provide(NewProviderMap),
		
		// Queue (implementação simples em memória por enquanto)
		fx.Provide(NewInMemoryQueue),
		
		// Services
		fx.Provide(NewNotificationService),
		fx.Provide(NewTemplateService),
		
		// HTTP Server
		fx.Provide(NewGinRouter),
		fx.Provide(NewHTTPServer),
		
		// Lifecycle
		fx.Invoke(StartHTTPServer),
	)

	app.Run()
}

// NewLogger cria instância do logger
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.IsDevelopment() {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

// NewDatabase cria conexão com o banco de dados
func NewDatabase(cfg *config.Config, logger *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetime) * time.Second)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully")
	return db, nil
}

// NewNotificationRepository cria repositório de notificações
func NewNotificationRepository(db *sqlx.DB, logger *zap.Logger) domain.NotificationRepository {
	return repository.NewPostgresNotificationRepository(db, logger)
}

// NewTemplateRepository cria repositório de templates
func NewTemplateRepository(db *sqlx.DB, logger *zap.Logger) domain.NotificationTemplateRepository {
	return repository.NewPostgresTemplateRepository(db, logger)
}

// NewPreferenceRepository cria repositório de preferências
func NewPreferenceRepository(db *sqlx.DB, logger *zap.Logger) domain.NotificationPreferenceRepository {
	return repository.NewPostgresPreferenceRepository(db, logger)
}

// NewEmailProvider cria provedor de email
func NewEmailProvider(cfg *config.Config, logger *zap.Logger) domain.NotificationProvider {
	emailConfig := providers.EmailConfig{
		Host:        cfg.SMTP.Host,
		Port:        cfg.SMTP.Port,
		Username:    cfg.SMTP.Username,
		Password:    cfg.SMTP.Password,
		FromEmail:   cfg.SMTP.FromEmail,
		FromName:    cfg.SMTP.FromName,
		UseTLS:      cfg.SMTP.UseTLS,
		UseStartTLS: cfg.SMTP.UseStartTLS,
		Timeout:     30,
		MaxRetries:  3,
		RateLimit:   60,
	}
	return providers.NewEmailProvider(emailConfig, logger)
}

// NewWhatsAppProvider cria provedor WhatsApp
func NewWhatsAppProvider(cfg *config.Config, logger *zap.Logger) domain.NotificationProvider {
	whatsappConfig := providers.WhatsAppConfig{
		AccessToken:   cfg.WhatsApp.AccessToken,
		PhoneNumberID: cfg.WhatsApp.PhoneNumberID,
		WebhookURL:    cfg.WhatsApp.WebhookURL,
		VerifyToken:   cfg.WhatsApp.VerifyToken,
		Timeout:       30,
		MaxRetries:    3,
		RateLimit:     60,
	}
	return providers.NewWhatsAppProvider(whatsappConfig, logger)
}

// NewTelegramProvider cria provedor Telegram
func NewTelegramProvider(cfg *config.Config, logger *zap.Logger) domain.NotificationProvider {
	telegramConfig := providers.TelegramConfig{
		BotToken:      cfg.Telegram.BotToken,
		WebhookURL:    cfg.Telegram.WebhookURL,
		WebhookSecret: cfg.Telegram.WebhookSecret,
		Timeout:       30,
		MaxRetries:    3,
		RateLimit:     30,
		ParseMode:     "HTML",
	}
	return providers.NewTelegramProvider(telegramConfig, logger)
}

// NewProviderMap cria mapa de provedores
func NewProviderMap(
	emailProvider domain.NotificationProvider,
	whatsappProvider domain.NotificationProvider,
	telegramProvider domain.NotificationProvider,
) map[domain.NotificationChannel]domain.NotificationProvider {
	return map[domain.NotificationChannel]domain.NotificationProvider{
		domain.NotificationChannelEmail:    emailProvider,
		domain.NotificationChannelWhatsApp: whatsappProvider,
		domain.NotificationChannelTelegram: telegramProvider,
	}
}

// InMemoryQueue implementação simples de fila em memória
type InMemoryQueue struct {
	queues map[domain.NotificationChannel]chan *domain.Notification
	logger *zap.Logger
}

// NewInMemoryQueue cria nova fila em memória
func NewInMemoryQueue(logger *zap.Logger) domain.NotificationQueue {
	return &InMemoryQueue{
		queues: make(map[domain.NotificationChannel]chan *domain.Notification),
		logger: logger,
	}
}

func (q *InMemoryQueue) Enqueue(ctx context.Context, notification *domain.Notification) error {
	if _, exists := q.queues[notification.Channel]; !exists {
		q.queues[notification.Channel] = make(chan *domain.Notification, 1000)
	}
	
	select {
	case q.queues[notification.Channel] <- notification:
		return nil
	default:
		return fmt.Errorf("queue full for channel: %s", notification.Channel)
	}
}

func (q *InMemoryQueue) EnqueueBatch(ctx context.Context, notifications []*domain.Notification) error {
	for _, notification := range notifications {
		if err := q.Enqueue(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

func (q *InMemoryQueue) EnqueueScheduled(ctx context.Context, notification *domain.Notification, scheduledAt time.Time) error {
	// Para simplicidade, enfileirar imediatamente
	// Em produção, usaria um scheduler
	return q.Enqueue(ctx, notification)
}

func (q *InMemoryQueue) Dequeue(ctx context.Context, channel domain.NotificationChannel) (*domain.Notification, error) {
	if queue, exists := q.queues[channel]; exists {
		select {
		case notification := <-queue:
			return notification, nil
		default:
			return nil, fmt.Errorf("queue empty for channel: %s", channel)
		}
	}
	return nil, fmt.Errorf("queue not found for channel: %s", channel)
}

func (q *InMemoryQueue) DequeueBatch(ctx context.Context, channel domain.NotificationChannel, limit int) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	for i := 0; i < limit; i++ {
		notification, err := q.Dequeue(ctx, channel)
		if err != nil {
			break
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (q *InMemoryQueue) DequeueByPriority(ctx context.Context, channel domain.NotificationChannel) (*domain.Notification, error) {
	return q.Dequeue(ctx, channel)
}

func (q *InMemoryQueue) QueueSize(ctx context.Context, channel domain.NotificationChannel) (int64, error) {
	if queue, exists := q.queues[channel]; exists {
		return int64(len(queue)), nil
	}
	return 0, nil
}

func (q *InMemoryQueue) QueueSizeByPriority(ctx context.Context, channel domain.NotificationChannel, priority domain.NotificationPriority) (int64, error) {
	return q.QueueSize(ctx, channel)
}

func (q *InMemoryQueue) Clear(ctx context.Context, channel domain.NotificationChannel) error {
	if queue, exists := q.queues[channel]; exists {
		for len(queue) > 0 {
			<-queue
		}
	}
	return nil
}

func (q *InMemoryQueue) Health(ctx context.Context) error {
	return nil
}

// NewNotificationService cria serviço de notificações
func NewNotificationService(
	notificationRepo domain.NotificationRepository,
	templateRepo domain.NotificationTemplateRepository,
	preferenceRepo domain.NotificationPreferenceRepository,
	queue domain.NotificationQueue,
	providers map[domain.NotificationChannel]domain.NotificationProvider,
	logger *zap.Logger,
) *services.NotificationService {
	return services.NewNotificationService(
		notificationRepo,
		templateRepo,
		preferenceRepo,
		queue,
		providers,
		logger,
	)
}

// NewTemplateService cria serviço de templates
func NewTemplateService(
	templateRepo domain.NotificationTemplateRepository,
	logger *zap.Logger,
) *services.TemplateService {
	return services.NewTemplateService(templateRepo, logger)
}

// NewGinRouter cria router Gin
func NewGinRouter(
	notificationService *services.NotificationService,
	templateService *services.TemplateService,
	preferenceRepo domain.NotificationPreferenceRepository,
	logger *zap.Logger,
) *gin.Engine {
	// Configurar Gin para produção
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	
	// Configurar rotas
	routerConfig := &httpinfra.RouterConfig{
		NotificationService: notificationService,
		TemplateService:     templateService,
		PreferenceRepo:      preferenceRepo,
		Logger:              logger,
	}
	
	httpinfra.SetupRoutes(router, routerConfig)
	
	return router
}

// NewHTTPServer cria servidor HTTP
func NewHTTPServer(cfg *config.Config, router *gin.Engine, logger *zap.Logger) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}
}

// StartHTTPServer inicia o servidor HTTP
func StartHTTPServer(
	lc fx.Lifecycle,
	server *http.Server,
	logger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting HTTP server", zap.String("addr", server.Addr))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})

	// Aguardar sinal de interrupção
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		logger.Info("Received shutdown signal")
		os.Exit(0)
	}()
}