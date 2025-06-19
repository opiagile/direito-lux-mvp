package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/application/services"
	"github.com/direito-lux/report-service/internal/infrastructure/http/handlers"
	"github.com/direito-lux/report-service/internal/infrastructure/http/middleware"
)

// RouterConfig configuração do router
type RouterConfig struct {
	ReportService    *services.ReportService
	DashboardService *services.DashboardService
	SchedulerService *services.SchedulerService
	Logger           *zap.Logger
}

// SetupRoutes configura as rotas do serviço
func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	// Middleware global
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(config.Logger))
	router.Use(middleware.CORS())

	// Handlers
	reportHandler := handlers.NewReportHandler(config.ReportService, config.Logger)
	dashboardHandler := handlers.NewDashboardHandler(config.DashboardService, config.Logger)
	schedulerHandler := handlers.NewSchedulerHandler(config.SchedulerService, config.Logger)
	healthHandler := handlers.NewHealthHandler(config.Logger)

	// Health check routes
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// API routes
	api := router.Group("/api/v1")
	{
		// Reports
		reports := api.Group("/reports")
		{
			reports.POST("", middleware.Auth(), reportHandler.CreateReport)
			reports.GET("", middleware.Auth(), reportHandler.ListReports)
			reports.GET("/:id", middleware.Auth(), reportHandler.GetReport)
			reports.DELETE("/:id", middleware.Auth(), reportHandler.DeleteReport)
			reports.GET("/:id/download", middleware.Auth(), reportHandler.DownloadReport)
			reports.GET("/stats", middleware.Auth(), reportHandler.GetStatistics)
		}

		// Dashboards
		dashboards := api.Group("/dashboards")
		{
			dashboards.POST("", middleware.Auth(), dashboardHandler.CreateDashboard)
			dashboards.GET("", middleware.Auth(), dashboardHandler.ListDashboards)
			dashboards.GET("/:id", middleware.Auth(), dashboardHandler.GetDashboard)
			dashboards.PUT("/:id", middleware.Auth(), dashboardHandler.UpdateDashboard)
			dashboards.DELETE("/:id", middleware.Auth(), dashboardHandler.DeleteDashboard)
			dashboards.GET("/:id/data", middleware.Auth(), dashboardHandler.GetDashboardData)
			
			// Widgets
			dashboards.POST("/:id/widgets", middleware.Auth(), dashboardHandler.AddWidget)
			dashboards.PUT("/:id/widgets/:widget_id", middleware.Auth(), dashboardHandler.UpdateWidget)
			dashboards.DELETE("/:id/widgets/:widget_id", middleware.Auth(), dashboardHandler.RemoveWidget)
		}

		// Schedules
		schedules := api.Group("/schedules")
		{
			schedules.POST("", middleware.Auth(), schedulerHandler.CreateSchedule)
			schedules.GET("", middleware.Auth(), schedulerHandler.ListSchedules)
			schedules.GET("/:id", middleware.Auth(), schedulerHandler.GetSchedule)
			schedules.PUT("/:id", middleware.Auth(), schedulerHandler.UpdateSchedule)
			schedules.DELETE("/:id", middleware.Auth(), schedulerHandler.DeleteSchedule)
		}

		// KPIs
		kpis := api.Group("/kpis")
		{
			kpis.GET("", middleware.Auth(), dashboardHandler.ListKPIs)
			kpis.POST("/calculate", middleware.Auth(), dashboardHandler.CalculateKPIs)
		}
	}

	// Metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		// TODO: Implementar endpoint de métricas Prometheus
		c.JSON(200, gin.H{"status": "metrics endpoint"})
	})
}