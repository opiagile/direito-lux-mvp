package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/application/services"
)

// DashboardHandler handler para dashboards
type DashboardHandler struct {
	dashboardService *services.DashboardService
	logger           *zap.Logger
}

// NewDashboardHandler cria nova instância do handler
func NewDashboardHandler(dashboardService *services.DashboardService, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		logger:           logger,
	}
}

// Métodos stub para compilação
func (h *DashboardHandler) CreateDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "CreateDashboard not implemented"})
}

func (h *DashboardHandler) ListDashboards(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListDashboards not implemented"})
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetDashboard not implemented"})
}

func (h *DashboardHandler) UpdateDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateDashboard not implemented"})
}

func (h *DashboardHandler) DeleteDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteDashboard not implemented"})
}

func (h *DashboardHandler) GetDashboardData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetDashboardData not implemented"})
}

func (h *DashboardHandler) AddWidget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "AddWidget not implemented"})
}

func (h *DashboardHandler) UpdateWidget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateWidget not implemented"})
}

func (h *DashboardHandler) RemoveWidget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "RemoveWidget not implemented"})
}

func (h *DashboardHandler) ListKPIs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListKPIs not implemented"})
}

func (h *DashboardHandler) CalculateKPIs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "CalculateKPIs not implemented"})
}