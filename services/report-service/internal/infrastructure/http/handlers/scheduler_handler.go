package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/application/services"
)

// SchedulerHandler handler para agendamentos
type SchedulerHandler struct {
	schedulerService *services.SchedulerService
	logger           *zap.Logger
}

// NewSchedulerHandler cria nova instância do handler
func NewSchedulerHandler(schedulerService *services.SchedulerService, logger *zap.Logger) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerService: schedulerService,
		logger:           logger,
	}
}

// Métodos stub para compilação
func (h *SchedulerHandler) CreateSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "CreateSchedule not implemented"})
}

func (h *SchedulerHandler) ListSchedules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListSchedules not implemented"})
}

func (h *SchedulerHandler) GetSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetSchedule not implemented"})
}

func (h *SchedulerHandler) UpdateSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateSchedule not implemented"})
}

func (h *SchedulerHandler) DeleteSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteSchedule not implemented"})
}