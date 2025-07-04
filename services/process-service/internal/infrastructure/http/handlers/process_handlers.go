package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ProcessStats estrutura para estatísticas de processos
type ProcessStats struct {
	Total     int `json:"total" db:"total"`
	Active    int `json:"active" db:"active"`
	Paused    int `json:"paused" db:"paused"`
	Archived  int `json:"archived" db:"archived"`
	ThisMonth int `json:"this_month" db:"this_month"`
}

// ProcessStatsHandler handler para endpoint de estatísticas
type ProcessStatsHandler struct {
	db *sqlx.DB
}

// NewProcessStatsHandler cria novo handler de estatísticas
func NewProcessStatsHandler(db *sqlx.DB) *ProcessStatsHandler {
	return &ProcessStatsHandler{db: db}
}

// GetProcessStats retorna estatísticas dos processos
// GET /api/v1/processes/stats
func (h *ProcessStatsHandler) GetProcessStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter tenant ID do header
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			return
		}

		// Query para buscar estatísticas
		query := `
			SELECT 
				COUNT(*) as total,
				COUNT(*) FILTER (WHERE status = 'active') as active,
				COUNT(*) FILTER (WHERE status = 'paused') as paused,
				COUNT(*) FILTER (WHERE status = 'archived') as archived,
				COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as this_month
			FROM processes 
			WHERE tenant_id = $1
		`

		var stats ProcessStats
		err := h.db.Get(&stats, query, tenantID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Se não há dados, retornar zeros
				stats = ProcessStats{
					Total:     0,
					Active:    0,
					Paused:    0,
					Archived:  0,
					ThisMonth: 0,
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erro ao buscar estatísticas de processos",
				})
				return
			}
		}

		// Retornar estatísticas
		c.JSON(http.StatusOK, gin.H{
			"data": stats,
			"timestamp": time.Now().UTC(),
		})
	}
}

// Process estrutura básica de processo
type Process struct {
	ID               string    `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	Number           string    `json:"number" db:"number"`
	Title            string    `json:"title" db:"title"`
	Description      string    `json:"description" db:"description"`
	Court            string    `json:"court" db:"court"`
	Status           string    `json:"status" db:"status"`
	MonitoringEnabled bool     `json:"monitoring_enabled" db:"monitoring_enabled"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ProcessHandler handler para operações de processos
type ProcessHandler struct {
	db *sqlx.DB
}

// NewProcessHandler cria novo handler de processos
func NewProcessHandler(db *sqlx.DB) *ProcessHandler {
	return &ProcessHandler{db: db}
}

// ListProcesses lista processos do tenant
// GET /api/v1/processes
func (h *ProcessHandler) ListProcesses() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			return
		}

		var processes []Process
		query := `
			SELECT id, tenant_id, number, title, description, court, status, 
			       monitoring_enabled, created_at, updated_at
			FROM processes 
			WHERE tenant_id = $1 
			ORDER BY created_at DESC
			LIMIT 50
		`

		err := h.db.Select(&processes, query, tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao listar processos",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": processes,
			"total": len(processes),
			"timestamp": time.Now().UTC(),
		})
	}
}

// GetProcess busca processo por ID
// GET /api/v1/processes/:id
func (h *ProcessHandler) GetProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			return
		}

		processID := c.Param("id")
		if processID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID do processo é obrigatório",
			})
			return
		}

		var process Process
		query := `
			SELECT id, tenant_id, number, title, description, court, status, 
			       monitoring_enabled, created_at, updated_at
			FROM processes 
			WHERE id = $1 AND tenant_id = $2
		`

		err := h.db.Get(&process, query, processID, tenantID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Processo não encontrado",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erro ao buscar processo",
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": process,
			"timestamp": time.Now().UTC(),
		})
	}
}

// CreateProcessRequest estrutura para criação de processo
type CreateProcessRequest struct {
	Number      string `json:"number" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Court       string `json:"court" binding:"required"`
}

// CreateProcess cria novo processo
// POST /api/v1/processes
func (h *ProcessHandler) CreateProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			return
		}

		var req CreateProcessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Dados inválidos: " + err.Error(),
			})
			return
		}

		query := `
			INSERT INTO processes (tenant_id, number, title, description, court, status)
			VALUES ($1, $2, $3, $4, $5, 'active')
			RETURNING id
		`

		var processID string
		err := h.db.Get(&processID, query, tenantID, req.Number, req.Title, req.Description, req.Court)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao criar processo",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": gin.H{
				"id": processID,
				"message": "Processo criado com sucesso",
			},
			"timestamp": time.Now().UTC(),
		})
	}
}

// UpdateProcessRequest estrutura para atualização de processo
type UpdateProcessRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Court       string `json:"court"`
	Status      string `json:"status"`
}

// UpdateProcess atualiza processo
// PUT /api/v1/processes/:id
func (h *ProcessHandler) UpdateProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			return
		}

		processID := c.Param("id")
		if processID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID do processo é obrigatório",
			})
			return
		}

		var req UpdateProcessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Dados inválidos: " + err.Error(),
			})
			return
		}

		query := `
			UPDATE processes 
			SET title = COALESCE(NULLIF($3, ''), title),
			    description = COALESCE(NULLIF($4, ''), description),
			    court = COALESCE(NULLIF($5, ''), court),
			    status = COALESCE(NULLIF($6, ''), status),
			    updated_at = NOW()
			WHERE id = $1 AND tenant_id = $2
		`

		result, err := h.db.Exec(query, processID, tenantID, req.Title, req.Description, req.Court, req.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao atualizar processo",
			})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Processo não encontrado",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"message": "Processo atualizado com sucesso",
			},
			"timestamp": time.Now().UTC(),
		})
	}
}

// DeleteProcess exclui processo
// DELETE /api/v1/processes/:id
func (h *ProcessHandler) DeleteProcess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-ID header é obrigatório",
			})
			return
		}

		processID := c.Param("id")
		if processID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID do processo é obrigatório",
			})
			return
		}

		query := `DELETE FROM processes WHERE id = $1 AND tenant_id = $2`

		result, err := h.db.Exec(query, processID, tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao excluir processo",
			})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Processo não encontrado",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"message": "Processo excluído com sucesso",
			},
			"timestamp": time.Now().UTC(),
		})
	}
}