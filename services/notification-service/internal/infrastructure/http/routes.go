package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/application/services"
	"github.com/direito-lux/notification-service/internal/domain"
	"github.com/direito-lux/notification-service/internal/infrastructure/http/handlers"
	"github.com/direito-lux/notification-service/internal/infrastructure/http/middleware"
)

// RouterConfig configuração do router
type RouterConfig struct {
	NotificationService *services.NotificationService
	TemplateService     *services.TemplateService
	PreferenceRepo      domain.NotificationPreferenceRepository
	Logger              *zap.Logger
}

// SetupRoutes configura todas as rotas da API
func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	// Criar handlers
	notificationHandler := handlers.NewNotificationHandler(config.NotificationService, config.Logger)
	templateHandler := handlers.NewTemplateHandler(config.TemplateService, config.Logger)
	preferenceHandler := handlers.NewPreferenceHandler(config.PreferenceRepo, config.Logger)

	// Middleware global
	router.Use(middleware.Logger(config.Logger))
	router.Use(middleware.CORS())
	router.Use(middleware.Recovery(config.Logger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "notification-service",
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Middleware de autenticação e tenant
		v1.Use(middleware.AuthRequired())
		v1.Use(middleware.TenantRequired())

		// Rotas de notificações
		notifications := v1.Group("/notifications")
		{
			notifications.POST("", notificationHandler.CreateNotification)
			notifications.GET("", notificationHandler.ListNotifications)
			notifications.GET("/:id", notificationHandler.GetNotification)
			notifications.POST("/:id/cancel", notificationHandler.CancelNotification)
			notifications.GET("/statistics", notificationHandler.GetNotificationStatistics)
			notifications.POST("/bulk", notificationHandler.SendBulkNotifications)
		}

		// Rotas de templates
		templates := v1.Group("/templates")
		{
			templates.POST("", templateHandler.CreateTemplate)
			templates.GET("", templateHandler.ListTemplates)
			templates.GET("/:id", templateHandler.GetTemplate)
			templates.PUT("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)
			templates.POST("/:id/preview", templateHandler.PreviewTemplate)
			templates.POST("/:id/duplicate", templateHandler.DuplicateTemplate)
			templates.POST("/:id/activate", templateHandler.ActivateTemplate)
			templates.POST("/:id/deactivate", templateHandler.DeactivateTemplate)
		}

		// Rotas de preferências
		preferences := v1.Group("/preferences")
		{
			preferences.POST("/defaults", preferenceHandler.CreateDefaultPreferences)
			
			// Preferências por usuário
			userPrefs := preferences.Group("/users/:user_id")
			{
				userPrefs.GET("", preferenceHandler.GetUserPreferences)
				userPrefs.PUT("", preferenceHandler.UpdateUserPreferences)
				userPrefs.POST("/disable-all", preferenceHandler.DisableAllNotifications)
				userPrefs.POST("/enable-all", preferenceHandler.EnableAllNotifications)
				
				// Preferências por tipo
				userPrefs.GET("/types/:type", preferenceHandler.GetPreferenceByType)
				userPrefs.PUT("/types/:type", preferenceHandler.UpdatePreferenceByType)
			}
		}
	}

	// Rotas administrativas (sem autenticação de tenant)
	admin := router.Group("/admin")
	{
		admin.Use(middleware.AdminAuthRequired())
		
		// Templates do sistema
		systemTemplates := admin.Group("/templates")
		{
			systemTemplates.GET("", templateHandler.ListTemplates) // Lista templates do sistema
		}
	}

	// Webhooks (sem autenticação)
	_ = router.Group("/webhooks")
	{
		// WhatsApp webhook seria adicionado aqui no futuro
		// webhooks.POST("/whatsapp", whatsappHandler.HandleWebhook)
		
		// Email bounce/complaint webhooks
		// webhooks.POST("/email/bounce", emailHandler.HandleBounce)
		// webhooks.POST("/email/complaint", emailHandler.HandleComplaint)
	}
}