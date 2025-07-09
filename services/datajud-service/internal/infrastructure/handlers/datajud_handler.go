package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/application"
	"github.com/direito-lux/datajud-service/internal/infrastructure/config"
)

// DataJudServiceInterface interface comum para ambos os services
type DataJudServiceInterface interface {
	QueryProcess(ctx context.Context, req *application.ProcessQueryRequest) (*application.ProcessQueryResponse, error)
	QueryMovements(ctx context.Context, req *application.MovementQueryRequest) (*application.MovementQueryResponse, error)
	BulkQuery(ctx context.Context, req *application.BulkQueryRequest) (*application.BulkQueryResponse, error)
}

// DataJudHandler handler HTTP para o serviço DataJud
type DataJudHandler struct {
	service DataJudServiceInterface
	config  *config.Config
}

// NewDataJudHandler cria nova instância do handler
func NewDataJudHandler(service DataJudServiceInterface, config *config.Config) *DataJudHandler {
	return &DataJudHandler{
		service: service,
		config:  config,
	}
}

// RegisterRoutes registra as rotas do handler
func (h *DataJudHandler) RegisterRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", h.healthCheck)
	
	// API routes
	api := router.Group("/api/v1")
	{
		// Consultas de processos
		api.POST("/process/query", h.requireAuth(), h.queryProcess)
		api.POST("/process/movements", h.requireAuth(), h.queryMovements)
		api.POST("/process/bulk", h.requireAuth(), h.bulkQuery)
		
		// Compatibilidade com API antiga
		api.POST("/search", h.requireAuth(), h.searchProcesses)
		api.GET("/process/:number", h.requireAuth(), h.getProcess)
		api.GET("/process/:number/movements", h.requireAuth(), h.getMovements)
		
		// Estatísticas e quota
		api.GET("/stats", h.requireAuth(), h.getStats)
		api.GET("/quota", h.requireAuth(), h.getQuota)
		
		// Tribunais
		api.GET("/tribunals", h.listTribunals)
		api.GET("/tribunals/:code", h.getTribunal)
	}
}

// healthCheck endpoint de health check
func (h *DataJudHandler) healthCheck(c *gin.Context) {
	health := gin.H{
		"status":     "healthy",
		"service":    "datajud-service",
		"version":    "2.0.0",
		"timestamp":  time.Now().Unix(),
		"environment": h.config.Environment,
		"datajud_mock": h.config.IsDataJudMockEnabled(),
	}
	
	c.JSON(http.StatusOK, health)
}

// requireAuth middleware de autenticação
func (h *DataJudHandler) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			c.Abort()
			return
		}
		
		// TODO: Validar token JWT quando Auth Service estiver integrado
		
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// queryProcess consulta processo específico
func (h *DataJudHandler) queryProcess(c *gin.Context) {
	var req application.ProcessQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de requisição inválido",
			"details": err.Error(),
		})
		return
	}
	
	// Definir tenant_id da requisição
	tenantIDStr := c.GetString("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant_id format",
		})
		return
	}
	req.TenantID = tenantID
	
	// Executar consulta
	response, err := h.service.QueryProcess(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao consultar processo",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// queryMovements consulta movimentações de processo
func (h *DataJudHandler) queryMovements(c *gin.Context) {
	var req application.MovementQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de requisição inválido",
			"details": err.Error(),
		})
		return
	}
	
	// Definir tenant_id da requisição
	tenantIDStr := c.GetString("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant_id format",
		})
		return
	}
	req.TenantID = tenantID
	
	// Executar consulta
	response, err := h.service.QueryMovements(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao consultar movimentações",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// bulkQuery consulta múltiplos processos
func (h *DataJudHandler) bulkQuery(c *gin.Context) {
	var req application.BulkQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de requisição inválido",
			"details": err.Error(),
		})
		return
	}
	
	// Definir tenant_id da requisição
	tenantIDStr := c.GetString("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant_id format",
		})
		return
	}
	req.TenantID = tenantID
	
	// Executar consulta
	response, err := h.service.BulkQuery(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao executar consulta em lote",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// searchProcesses compatibilidade com API antiga
func (h *DataJudHandler) searchProcesses(c *gin.Context) {
	var request struct {
		Query      string   `json:"query" binding:"required"`
		Tribunais  []string `json:"tribunais"`
		DataInicio string   `json:"data_inicio"`
		DataFim    string   `json:"data_fim"`
		Pagina     int      `json:"pagina"`
		Tamanho    int      `json:"tamanho"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Converter para nova estrutura
	// Converter tenant_id para UUID
	tenantIDStr := c.GetString("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id format"})
		return
	}

	// Gerar UUID para cliente legacy
	legacyClientID := uuid.New()

	req := application.ProcessQueryRequest{
		ProcessNumber: request.Query,
		TenantID:     tenantID,
		ClientID:     legacyClientID,
		UseCache:     true,
		Urgent:       false,
	}
	
	// Usar primeiro tribunal se especificado
	if len(request.Tribunais) > 0 {
		req.CourtID = request.Tribunais[0]
	}
	
	// Executar consulta
	response, err := h.service.QueryProcess(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar processos",
			"details": err.Error(),
		})
		return
	}
	
	// Converter para formato antigo
	legacyResponse := gin.H{
		"data": gin.H{
			"total": 1,
			"processos": []gin.H{},
		},
		"meta": gin.H{
			"tenant_id": req.TenantID,
			"timestamp": time.Now().Unix(),
		},
	}
	
	if response.Data != nil && response.Data.Number != "" {
		processo := gin.H{
			"numero":           response.Data.Number,
			"tribunal":         response.Data.Court,
			"assunto":          response.Data.Subject,
			"dataDistribuicao": response.Data.CreatedAt,
			"valorCausa":       0, // Not available in new structure
			"status":           response.Data.Status,
		}
		legacyResponse["data"].(gin.H)["processos"] = []gin.H{processo}
	}
	
	c.JSON(http.StatusOK, legacyResponse)
}

// getProcess compatibilidade com API antiga
func (h *DataJudHandler) getProcess(c *gin.Context) {
	number := c.Param("number")
	
	// Converter tenant_id para UUID
	tenantIDStr := c.GetString("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id format"})
		return
	}

	// Gerar UUID para cliente legacy
	legacyClientID := uuid.New()

	req := application.ProcessQueryRequest{
		ProcessNumber: number,
		TenantID:     tenantID,
		ClientID:     legacyClientID,
		UseCache:     true,
		Urgent:       false,
	}
	
	response, err := h.service.QueryProcess(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar processo",
			"details": err.Error(),
		})
		return
	}
	
	// Converter para formato antigo
	legacyResponse := gin.H{
		"data": gin.H{
			"numero": number,
		},
		"meta": gin.H{
			"tenant_id": req.TenantID,
			"timestamp": time.Now().Unix(),
		},
	}
	
	if response.Data != nil && response.Data.Number != "" {
		legacyResponse["data"] = gin.H{
			"numero":           response.Data.Number,
			"tribunal":         response.Data.Court,
			"assunto":          response.Data.Subject,
			"dataDistribuicao": response.Data.CreatedAt,
			"valorCausa":       0, // Not available in new structure
			"status":           response.Data.Status,
			"juizo":            response.Data.Stage,
		}
	}
	
	c.JSON(http.StatusOK, legacyResponse)
}

// getMovements compatibilidade com API antiga
func (h *DataJudHandler) getMovements(c *gin.Context) {
	number := c.Param("number")
	
	// Converter tenant_id para UUID
	tenantIDStr := c.GetString("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id format"})
		return
	}

	// Gerar UUID para cliente legacy
	legacyClientID := uuid.New()

	req := application.MovementQueryRequest{
		ProcessNumber: number,
		TenantID:     tenantID,
		ClientID:     legacyClientID,
		Page:         1,
		PageSize:     50,
		UseCache:     true,
		Urgent:       false,
	}
	
	response, err := h.service.QueryMovements(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar movimentações",
			"details": err.Error(),
		})
		return
	}
	
	// Converter para formato antigo
	legacyResponse := gin.H{
		"data": gin.H{
			"numero": number,
			"movimentacoes": []gin.H{},
		},
		"meta": gin.H{
			"tenant_id": req.TenantID,
			"timestamp": time.Now().Unix(),
		},
	}
	
	if response.Data != nil && len(response.Data.Movements) > 0 {
		movimentacoes := make([]gin.H, len(response.Data.Movements))
		for i, mov := range response.Data.Movements {
			movimentacoes[i] = gin.H{
				"data":        mov.Date,
				"descricao":   mov.Description,
				"complemento": mov.Content,
				"titulo":      mov.Title,
				"codigo":      mov.Code,
				"tipo":        mov.Type,
				"sequencia":   mov.Sequence,
			}
		}
		legacyResponse["data"].(gin.H)["movimentacoes"] = movimentacoes
	}
	
	c.JSON(http.StatusOK, legacyResponse)
}

// getStats endpoint de estatísticas
func (h *DataJudHandler) getStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar estatísticas reais
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"consultas_hoje":        45,
			"consultas_mes":         1250,
			"limite_diario":         1000,
			"usado_hoje":            45,
			"quota_disponivel":      955,
			"percentual_uso":        4.5,
			"cnpj_providers_ativos": 3,
			"cache_hits":            320,
			"cache_misses":          125,
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// getQuota endpoint de quota
func (h *DataJudHandler) getQuota(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar controle de quota real
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"usado_hoje":       45,
			"limite_diario":    1000,
			"percentual_uso":   4.5,
			"resets_em":        "00:00:00",
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// listTribunals lista tribunais disponíveis
func (h *DataJudHandler) listTribunals(c *gin.Context) {
	// TODO: Usar TribunalMapper para listar tribunais
	tribunais := []gin.H{
		{"code": "STF", "name": "Supremo Tribunal Federal", "type": "supremo"},
		{"code": "STJ", "name": "Superior Tribunal de Justiça", "type": "superior"},
		{"code": "TJSP", "name": "Tribunal de Justiça de São Paulo", "type": "estadual"},
		{"code": "TJRJ", "name": "Tribunal de Justiça do Rio de Janeiro", "type": "estadual"},
		{"code": "TRF1", "name": "Tribunal Regional Federal da 1ª Região", "type": "federal"},
		{"code": "TRF2", "name": "Tribunal Regional Federal da 2ª Região", "type": "federal"},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": tribunais,
		"meta": gin.H{
			"total": len(tribunais),
			"timestamp": time.Now().Unix(),
		},
	})
}

// getTribunal busca tribunal específico
func (h *DataJudHandler) getTribunal(c *gin.Context) {
	code := c.Param("code")
	
	// TODO: Usar TribunalMapper para buscar tribunal
	tribunal := gin.H{
		"code": code,
		"name": "Tribunal de Justiça de São Paulo",
		"type": "estadual",
		"active": true,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": tribunal,
		"meta": gin.H{
			"timestamp": time.Now().Unix(),
		},
	})
}

// parseIntParam converte parâmetro string para int
func parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	
	return defaultValue
}

// parseBoolParam converte parâmetro string para bool
func parseBoolParam(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	
	return value == "true" || value == "1"
}