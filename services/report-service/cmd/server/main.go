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
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/application/services"
	"github.com/direito-lux/report-service/internal/domain"
	"github.com/direito-lux/report-service/internal/infrastructure/config"
	"github.com/direito-lux/report-service/internal/infrastructure/database"
	"github.com/direito-lux/report-service/internal/infrastructure/events"
	httpinfra "github.com/direito-lux/report-service/internal/infrastructure/http"
	"github.com/direito-lux/report-service/internal/infrastructure/report"
)

func main() {
	app := fx.New(
		// Configuração
		fx.Provide(config.Load),

		// Logger
		fx.Provide(NewLogger),

		// Database
		fx.Provide(NewDatabase),
		fx.Provide(NewRedis),

		// Repositories
		fx.Provide(NewReportRepository),
		fx.Provide(NewScheduleRepository),
		fx.Provide(NewDashboardRepository),
		fx.Provide(NewKPIRepository),
		fx.Provide(NewTemplateRepository),

		// Infrastructure
		fx.Provide(NewReportGenerator),
		fx.Provide(NewDataCollector),
		fx.Provide(NewEventBus),

		// Services
		fx.Provide(NewReportService),
		fx.Provide(NewDashboardService),
		fx.Provide(NewSchedulerService),

		// HTTP
		fx.Provide(NewGinRouter),
		fx.Provide(NewHTTPServer),

		// Lifecycle
		fx.Invoke(StartScheduler),
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
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully")
	return db, nil
}

// NewRedis cria conexão com Redis
func NewRedis(cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	})

	// Testar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully")
	return client, nil
}

// NewReportRepository cria repositório de relatórios
func NewReportRepository(db *sqlx.DB, logger *zap.Logger) domain.ReportRepository {
	return database.NewPostgresReportRepository(db, logger)
}

// NewScheduleRepository cria repositório de agendamentos
func NewScheduleRepository(db *sqlx.DB, logger *zap.Logger) domain.ReportScheduleRepository {
	return database.NewPostgresScheduleRepository(db, logger)
}

// NewDashboardRepository cria repositório de dashboards
func NewDashboardRepository(db *sqlx.DB, logger *zap.Logger) domain.DashboardRepository {
	return database.NewPostgresDashboardRepository(db, logger)
}

// NewKPIRepository cria repositório de KPIs
func NewKPIRepository(db *sqlx.DB, logger *zap.Logger) domain.KPIRepository {
	return database.NewPostgresKPIRepository(db, logger)
}

// NewTemplateRepository cria repositório de templates
func NewTemplateRepository(db *sqlx.DB, logger *zap.Logger) domain.ReportTemplateRepository {
	return database.NewPostgresTemplateRepository(db, logger)
}

// NewReportGenerator cria gerador de relatórios
func NewReportGenerator(logger *zap.Logger) domain.ReportGenerator {
	return report.NewGenerator(logger)
}

// NewDataCollector cria coletor de dados
func NewDataCollector(cfg *config.Config, logger *zap.Logger) domain.DataCollector {
	return database.NewDataCollector(cfg, logger)
}

// NewEventBus cria event bus
func NewEventBus(logger *zap.Logger) domain.EventBus {
	return events.NewInMemoryEventBus(logger)
}

// NewReportService cria serviço de relatórios
func NewReportService(
	cfg *config.Config,
	reportRepo domain.ReportRepository,
	scheduleRepo domain.ReportScheduleRepository,
	templateRepo domain.ReportTemplateRepository,
	generator domain.ReportGenerator,
	dataCollector domain.DataCollector,
	eventBus domain.EventBus,
	logger *zap.Logger,
) *services.ReportService {
	return services.NewReportService(
		cfg,
		reportRepo,
		scheduleRepo,
		templateRepo,
		generator,
		dataCollector,
		eventBus,
		logger,
	)
}

// NewDashboardService cria serviço de dashboards
func NewDashboardService(
	dashboardRepo domain.DashboardRepository,
	kpiRepo domain.KPIRepository,
	dataCollector domain.DataCollector,
	eventBus domain.EventBus,
	logger *zap.Logger,
) *services.DashboardService {
	return services.NewDashboardService(
		dashboardRepo,
		kpiRepo,
		dataCollector,
		eventBus,
		logger,
	)
}

// NewSchedulerService cria serviço de agendamento
func NewSchedulerService(
	scheduleRepo domain.ReportScheduleRepository,
	reportService *services.ReportService,
	eventBus domain.EventBus,
	logger *zap.Logger,
) *services.SchedulerService {
	return services.NewSchedulerService(
		scheduleRepo,
		reportService,
		eventBus,
		logger,
	)
}

// NewGinRouter cria router Gin
func NewGinRouter(
	cfg *config.Config,
	reportService *services.ReportService,
	dashboardService *services.DashboardService,
	schedulerService *services.SchedulerService,
	logger *zap.Logger,
) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Configurar rotas
	routerConfig := &httpinfra.RouterConfig{
		ReportService:    reportService,
		DashboardService: dashboardService,
		SchedulerService: schedulerService,
		Logger:           logger,
	}

	httpinfra.SetupRoutes(router, routerConfig)

	return router
}

// NewHTTPServer cria servidor HTTP
func NewHTTPServer(cfg *config.Config, router *gin.Engine, logger *zap.Logger) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

// StartScheduler inicia o scheduler
func StartScheduler(
	lc fx.Lifecycle,
	scheduler *services.SchedulerService,
	logger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting report scheduler")
			return scheduler.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping report scheduler")
			scheduler.Stop()
			return nil
		},
	})
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