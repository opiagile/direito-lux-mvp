package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/application/services"
	"github.com/direito-lux/report-service/internal/domain"
)

// ReportHandler handler para relatórios
type ReportHandler struct {
	reportService *services.ReportService
	logger        *zap.Logger
}

// NewReportHandler cria nova instância do handler
func NewReportHandler(reportService *services.ReportService, logger *zap.Logger) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		logger:        logger,
	}
}

// CreateReport cria um novo relatório
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req services.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Adicionar contexto de autenticação
	ctx := h.addAuthContext(c)

	report, err := h.reportService.CreateReport(ctx, &req)
	if err != nil {
		h.logger.Error("Failed to create report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// ListReports lista relatórios
func (h *ReportHandler) ListReports(c *gin.Context) {
	ctx := h.addAuthContext(c)

	// Parse query parameters
	filters := domain.ReportFilters{
		Limit:  10,
		Offset: 0,
	}

	reports, err := h.reportService.ListReports(ctx, filters)
	if err != nil {
		h.logger.Error("Failed to list reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

// GetReport obtém um relatório
func (h *ReportHandler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	ctx := h.addAuthContext(c)

	report, err := h.reportService.GetReport(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// DeleteReport exclui um relatório
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	ctx := h.addAuthContext(c)

	if err := h.reportService.DeleteReport(ctx, id); err != nil {
		h.logger.Error("Failed to delete report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DownloadReport faz download de um relatório
func (h *ReportHandler) DownloadReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	ctx := h.addAuthContext(c)

	report, err := h.reportService.GetReport(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	if report.Status != domain.ReportStatusCompleted || report.FileURL == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report not ready for download"})
		return
	}

	// TODO: Implementar download real do arquivo
	c.JSON(http.StatusOK, gin.H{
		"download_url": *report.FileURL,
		"expires_at":   report.ExpiresAt,
	})
}

// GetStatistics obtém estatísticas de relatórios
func (h *ReportHandler) GetStatistics(c *gin.Context) {
	ctx := h.addAuthContext(c)

	// TODO: Parse period from query
	stats, err := h.reportService.GetStatistics(ctx, 30*24*time.Hour)
	if err != nil {
		h.logger.Error("Failed to get statistics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// addAuthContext adiciona contexto de autenticação
func (h *ReportHandler) addAuthContext(c *gin.Context) context.Context {
	// TODO: Extrair dados reais de autenticação dos headers
	tenantID := uuid.New()
	userID := uuid.New()
	plan := "professional"

	authCtx := domain.AuthContext{
		TenantID: tenantID,
		UserID:   userID,
		Plan:     plan,
		TraceID:  c.GetHeader("X-Trace-ID"),
	}

	return domain.WithAuthContext(c.Request.Context(), authCtx)
}