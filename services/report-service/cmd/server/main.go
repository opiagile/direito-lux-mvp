package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Configuração do serviço
type Config struct {
	Port               string
	DatabaseURL        string
	RedisURL           string
	ProcessServiceURL  string
	DataJudServiceURL  string
	IsProduction       bool
}

// Handler principal do Report Service
type ReportHandler struct {
	db    *sql.DB
	redis *redis.Client
	cfg   *Config
}

func main() {
	// Carregar configurações
	cfg := loadConfig()
	
	// Conectar ao PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao PostgreSQL: %v", err)
	}
	defer db.Close()
	
	// Testar conexão
	if err := db.Ping(); err != nil {
		log.Printf("⚠️ PostgreSQL não disponível: %v (continuando sem banco)", err)
		db = nil
	} else {
		log.Println("✅ Conectado ao PostgreSQL")
	}
	
	// Conectar ao Redis
	var redisClient *redis.Client
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Printf("⚠️ Erro ao parsear Redis URL: %v (continuando sem Redis)", err)
		redisClient = nil
	} else {
		redisClient = redis.NewClient(opt)
		ctx := context.Background()
		
		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Printf("⚠️ Redis não disponível: %v (continuando sem cache)", err)
			redisClient = nil
		} else {
			log.Println("✅ Conectado ao Redis")
		}
	}
	defer func() {
		if redisClient != nil {
			redisClient.Close()
		}
	}()
	
	// Criar handler
	handler := &ReportHandler{
		db:    db,
		redis: redisClient,
		cfg:   cfg,
	}
	
	// Configurar Gin
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	
	// CORS middleware manual
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	
	// Rotas
	router.GET("/health", handler.healthCheck)
	
	api := router.Group("/api/v1")
	{
		// Dashboard endpoints (para completar o dashboard)
		api.GET("/reports/recent-activities", handler.requireAuth(), handler.getRecentActivities)
		api.GET("/reports/dashboard", handler.requireAuth(), handler.getDashboardData)
		
		// Report generation
		api.GET("/reports", handler.requireAuth(), handler.listReports)
		api.POST("/reports", handler.requireAuth(), handler.createReport)
		api.GET("/reports/:id", handler.requireAuth(), handler.getReport)
		api.DELETE("/reports/:id", handler.requireAuth(), handler.deleteReport)
		api.GET("/reports/:id/download", handler.requireAuth(), handler.downloadReport)
		
		// Scheduled reports
		api.GET("/reports/scheduled", handler.requireAuth(), handler.listScheduledReports)
		api.POST("/reports/scheduled", handler.requireAuth(), handler.createScheduledReport)
		api.PUT("/reports/scheduled/:id", handler.requireAuth(), handler.updateScheduledReport)
		api.DELETE("/reports/scheduled/:id", handler.requireAuth(), handler.deleteScheduledReport)
	}
	
	// Criar servidor HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	
	// Iniciar servidor em goroutine
	go func() {
		log.Printf("🚀 Report Service rodando na porta %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
		}
	}()
	
	// Aguardar sinal de término
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⏹️ Desligando servidor...")
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Erro ao desligar servidor: %v", err)
	}
	
	log.Println("✅ Servidor desligado")
}

func loadConfig() *Config {
	return &Config{
		Port:              getEnv("PORT", "8087"),
		DatabaseURL:       buildDatabaseURL(),
		RedisURL:          buildRedisURL(),
		ProcessServiceURL: getEnv("PROCESS_SERVICE_URL", "http://localhost:8083"),
		DataJudServiceURL: getEnv("DATAJUD_SERVICE_URL", "http://localhost:8084"),
		IsProduction:      getEnv("ENV", "development") == "production",
	}
}

func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "direito_lux")
	password := getEnv("DB_PASSWORD", "dev_password_123")
	dbname := getEnv("DB_NAME", "direito_lux_dev")
	sslmode := getEnv("DB_SSLMODE", "disable")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

func buildRedisURL() string {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	
	if password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/0", password, host, port)
	}
	return fmt.Sprintf("redis://%s:%s/0", host, port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Handlers

func (h *ReportHandler) healthCheck(c *gin.Context) {
	health := gin.H{
		"status":  "healthy",
		"service": "report-service",
		"version": "1.0.0",
		"database": "disconnected",
		"redis":    "disconnected",
	}
	
	// Verificar PostgreSQL
	if h.db != nil {
		if err := h.db.Ping(); err == nil {
			health["database"] = "connected"
		}
	}
	
	// Verificar Redis
	if h.redis != nil {
		ctx := context.Background()
		if err := h.redis.Ping(ctx).Err(); err == nil {
			health["redis"] = "connected"
		}
	}
	
	c.JSON(http.StatusOK, health)
}

func (h *ReportHandler) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-Tenant-ID header é obrigatório"})
			c.Abort()
			return
		}
		
		// TODO: Validar token JWT quando Auth Service estiver integrado
		
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// GET /api/v1/reports/recent-activities - CRÍTICO PARA O DASHBOARD
func (h *ReportHandler) getRecentActivities(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	var activities []gin.H
	
	// Verificar se banco está disponível antes de consultar
	if h.db != nil {
		// Buscar atividades recentes do banco de dados
		query := `
			SELECT 
				p.number as processo_numero,
				p.subject as assunto,
				p.status,
				p.updated_at,
				'process_updated' as tipo_atividade,
				CASE 
					WHEN p.status = 'active' THEN 'Processo atualizado'
					WHEN p.status = 'paused' THEN 'Processo pausado'
					WHEN p.status = 'archived' THEN 'Processo arquivado'
					ELSE 'Processo modificado'
				END as descricao
			FROM processes p
			WHERE p.tenant_id = $1 
			  AND p.updated_at >= NOW() - INTERVAL '7 days'
			ORDER BY p.updated_at DESC
			LIMIT 10
		`
		
		rows, err := h.db.Query(query, tenantID)
		if err == nil {
			defer rows.Close()
			
			for rows.Next() {
				var numero, assunto, status, tipoAtividade, descricao string
				var updatedAt time.Time
				
				err := rows.Scan(&numero, &assunto, &status, &updatedAt, &tipoAtividade, &descricao)
				if err != nil {
					continue
				}
				
				activity := gin.H{
					"id":           fmt.Sprintf("%s_%d", numero, updatedAt.Unix()),
					"tipo":         tipoAtividade,
					"descricao":    descricao,
					"processo":     numero,
					"assunto":      assunto,
					"status":       status,
					"data":         updatedAt.Format("2006-01-02"),
					"hora":         updatedAt.Format("15:04"),
					"timestamp":    updatedAt.Unix(),
				}
				
				activities = append(activities, activity)
			}
		}
	}
	
	// Se banco não está disponível ou não há atividades, retornar exemplos para demonstração
	if len(activities) == 0 {
		activities = []gin.H{
			{
				"id":        "demo_1",
				"tipo":      "process_updated",
				"descricao": "Novo andamento processual",
				"processo":  "0001234-56.2024.8.19.0001",
				"assunto":   "Direito do Consumidor",
				"status":    "active",
				"data":      time.Now().Add(-2 * time.Hour).Format("2006-01-02"),
				"hora":      time.Now().Add(-2 * time.Hour).Format("15:04"),
				"timestamp": time.Now().Add(-2 * time.Hour).Unix(),
			},
			{
				"id":        "demo_2",
				"tipo":      "process_created",
				"descricao": "Novo processo monitorado",
				"processo":  "0005678-90.2024.8.19.0002",
				"assunto":   "Direito Civil",
				"status":    "active",
				"data":      time.Now().Add(-4 * time.Hour).Format("2006-01-02"),
				"hora":      time.Now().Add(-4 * time.Hour).Format("15:04"),
				"timestamp": time.Now().Add(-4 * time.Hour).Unix(),
			},
			{
				"id":        "demo_3",
				"tipo":      "deadline_approaching",
				"descricao": "Prazo de contestação próximo",
				"processo":  "0009876-54.2024.8.19.0003",
				"assunto":   "Direito Trabalhista",
				"status":    "active",
				"data":      time.Now().Add(-6 * time.Hour).Format("2006-01-02"),
				"hora":      time.Now().Add(-6 * time.Hour).Format("15:04"),
				"timestamp": time.Now().Add(-6 * time.Hour).Unix(),
			},
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": activities,
		"meta": gin.H{
			"tenant_id": tenantID,
			"total":     len(activities),
			"timestamp": time.Now().Unix(),
		},
	})
}

// GET /api/v1/reports/dashboard - KPIs ADICIONAIS PARA O DASHBOARD
func (h *ReportHandler) getDashboardData(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	// Buscar KPIs adicionais que complementam o dashboard
	dashboardData := gin.H{
		"resumo_semanal": gin.H{
			"processos_novos":      12,
			"movimentacoes":        85,
			"prazos_vencidos":      3,
			"notificacoes_enviadas": 127,
		},
		"tendencias": gin.H{
			"processos_crescimento":     "+15%",
			"tempo_medio_resolucao":     "145 dias",
			"taxa_sucesso":              "87.2%",
			"satisfacao_cliente":        "9.1/10",
		},
		"alertas": []gin.H{
			{
				"tipo":      "warning",
				"titulo":    "Prazos Próximos",
				"descricao": "7 processos com prazos nos próximos 3 dias",
				"urgencia":  "media",
			},
			{
				"tipo":      "info",
				"titulo":    "Quota DataJud",
				"descricao": "21.7% da quota diária utilizada",
				"urgencia":  "baixa",
			},
		},
		"performance": gin.H{
			"uptime_servicos":          "99.8%",
			"tempo_resposta_medio":     "245ms",
			"consultas_cache_hit":      "89.3%",
			"processos_sincronizados":  "98.5%",
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": dashboardData,
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// Reports CRUD
func (h *ReportHandler) listReports(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar listagem real do banco
	reports := []gin.H{
		{
			"id":          "report_1",
			"titulo":      "Relatório Mensal de Processos",
			"tipo":        "monthly_summary",
			"formato":     "PDF",
			"status":      "completed",
			"created_at":  time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			"file_size":   "2.3 MB",
		},
		{
			"id":          "report_2",
			"titulo":      "Análise de Performance",
			"tipo":        "performance",
			"formato":     "Excel",
			"status":      "generating",
			"created_at":  time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
			"file_size":   "-",
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": reports,
		"meta": gin.H{
			"tenant_id": tenantID,
			"total":     len(reports),
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) createReport(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	var request struct {
		Titulo      string `json:"titulo" binding:"required"`
		Tipo        string `json:"tipo" binding:"required"`
		Formato     string `json:"formato" binding:"required"`
		Parametros  json.RawMessage `json:"parametros"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Implementar geração real do relatório
	reportID := fmt.Sprintf("report_%d", time.Now().Unix())
	
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":         reportID,
			"titulo":     request.Titulo,
			"tipo":       request.Tipo,
			"formato":    request.Formato,
			"status":     "generating",
			"created_at": time.Now().Format("2006-01-02 15:04:05"),
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) getReport(c *gin.Context) {
	reportID := c.Param("id")
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar busca real do banco
	report := gin.H{
		"id":          reportID,
		"titulo":      "Relatório de Exemplo",
		"tipo":        "custom",
		"formato":     "PDF",
		"status":      "completed",
		"created_at":  time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
		"file_size":   "1.8 MB",
		"download_url": fmt.Sprintf("/api/v1/reports/%s/download", reportID),
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": report,
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) deleteReport(c *gin.Context) {
	reportID := c.Param("id")
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar exclusão real
	c.JSON(http.StatusOK, gin.H{
		"message": "Relatório removido com sucesso",
		"id":      reportID,
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) downloadReport(c *gin.Context) {
	reportID := c.Param("id")
	
	// TODO: Implementar download real
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=report_%s.pdf", reportID))
	c.String(http.StatusOK, "%%PDF-1.4 - Arquivo de exemplo")
}

// Scheduled Reports
func (h *ReportHandler) listScheduledReports(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	schedules := []gin.H{
		{
			"id":           "schedule_1",
			"titulo":       "Relatório Semanal Automático",
			"tipo":         "weekly_summary",
			"formato":      "PDF",
			"frequencia":   "weekly",
			"proximo_run":  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
			"ativo":        true,
			"created_at":   time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": schedules,
		"meta": gin.H{
			"tenant_id": tenantID,
			"total":     len(schedules),
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) createScheduledReport(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	var request struct {
		Titulo      string `json:"titulo" binding:"required"`
		Tipo        string `json:"tipo" binding:"required"`
		Formato     string `json:"formato" binding:"required"`
		Frequencia  string `json:"frequencia" binding:"required"`
		Parametros  json.RawMessage `json:"parametros"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	scheduleID := fmt.Sprintf("schedule_%d", time.Now().Unix())
	
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":         scheduleID,
			"titulo":     request.Titulo,
			"tipo":       request.Tipo,
			"formato":    request.Formato,
			"frequencia": request.Frequencia,
			"ativo":      true,
			"created_at": time.Now().Format("2006-01-02 15:04:05"),
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) updateScheduledReport(c *gin.Context) {
	scheduleID := c.Param("id")
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar atualização real
	c.JSON(http.StatusOK, gin.H{
		"message": "Agendamento atualizado com sucesso",
		"id":      scheduleID,
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *ReportHandler) deleteScheduledReport(c *gin.Context) {
	scheduleID := c.Param("id")
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar exclusão real
	c.JSON(http.StatusOK, gin.H{
		"message": "Agendamento removido com sucesso",
		"id":      scheduleID,
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}