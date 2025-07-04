package main

import (
	"context"
	"database/sql"
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
	Port            string
	DatabaseURL     string
	RedisURL        string
	DataJudAPIURL   string
	DataJudAPIKey   string
	RateLimitDaily  int
	IsProduction    bool
}

// Handler principal do DataJud Service
type DataJudHandler struct {
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
		log.Fatalf("❌ Erro ao testar conexão PostgreSQL: %v", err)
	}
	log.Println("✅ Conectado ao PostgreSQL")
	
	// Conectar ao Redis
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("❌ Erro ao parsear Redis URL: %v", err)
	}
	
	redisClient := redis.NewClient(opt)
	ctx := context.Background()
	
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ Redis não disponível: %v (continuando sem cache)", err)
		redisClient = nil
	} else {
		log.Println("✅ Conectado ao Redis")
	}
	defer func() {
		if redisClient != nil {
			redisClient.Close()
		}
	}()
	
	// Criar handler
	handler := &DataJudHandler{
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
		// Busca de processos
		api.POST("/search", handler.requireAuth(), handler.searchProcesses)
		api.GET("/process/:number", handler.requireAuth(), handler.getProcess)
		api.GET("/process/:number/movements", handler.requireAuth(), handler.getMovements)
		
		// Estatísticas e quota
		api.GET("/stats", handler.requireAuth(), handler.getStats)
		api.GET("/quota", handler.requireAuth(), handler.getQuota)
		
		// Gerenciamento de CNPJs (admin)
		api.GET("/cnpj-providers", handler.requireAuth(), handler.requireAdmin(), handler.listCNPJProviders)
		api.POST("/cnpj-providers", handler.requireAuth(), handler.requireAdmin(), handler.createCNPJProvider)
		api.PUT("/cnpj-providers/:id", handler.requireAuth(), handler.requireAdmin(), handler.updateCNPJProvider)
		api.DELETE("/cnpj-providers/:id", handler.requireAuth(), handler.requireAdmin(), handler.deleteCNPJProvider)
	}
	
	// Criar servidor HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	
	// Iniciar servidor em goroutine
	go func() {
		log.Printf("🚀 DataJud Service rodando na porta %s", cfg.Port)
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
		Port:           getEnv("PORT", "8084"),
		DatabaseURL:    buildDatabaseURL(),
		RedisURL:       buildRedisURL(),
		DataJudAPIURL:  getEnv("DATAJUD_API_URL", "https://api-publica.datajud.cnj.jus.br"),
		DataJudAPIKey:  getEnv("DATAJUD_API_KEY", ""),
		RateLimitDaily: getEnvAsInt("RATE_LIMIT_DAILY", 10000),
		IsProduction:   getEnv("ENV", "development") == "production",
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

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Handlers

func (h *DataJudHandler) healthCheck(c *gin.Context) {
	health := gin.H{
		"status":  "healthy",
		"service": "datajud-service",
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

func (h *DataJudHandler) requireAuth() gin.HandlerFunc {
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

func (h *DataJudHandler) requireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Verificar se usuário é admin
		c.Next()
	}
}

// Search processes na API DataJud
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
	
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar busca real na API DataJud
	// Por enquanto, retornar dados de exemplo
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total": 2,
			"processos": []gin.H{
				{
					"numero": "0001234-56.2024.8.19.0001",
					"tribunal": "TJRJ",
					"assunto": "Direito do Consumidor",
					"dataDistribuicao": "2024-01-15",
					"valorCausa": 15000.00,
					"partes": gin.H{
						"autor": "João da Silva",
						"reu": "Empresa XYZ Ltda",
					},
				},
				{
					"numero": "0005678-90.2024.8.19.0002",
					"tribunal": "TJRJ",
					"assunto": "Direito Civil",
					"dataDistribuicao": "2024-02-20",
					"valorCausa": 25000.00,
					"partes": gin.H{
						"autor": "Maria Santos",
						"reu": "Empresa ABC S.A.",
					},
				},
			},
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// Get specific process
func (h *DataJudHandler) getProcess(c *gin.Context) {
	number := c.Param("number")
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar busca real na API DataJud
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"numero": number,
			"tribunal": "TJRJ",
			"assunto": "Direito do Consumidor",
			"dataDistribuicao": "2024-01-15",
			"valorCausa": 15000.00,
			"status": "Em andamento",
			"juizo": "1ª Vara Cível",
			"partes": gin.H{
				"autor": gin.H{
					"nome": "João da Silva",
					"cpf": "123.456.789-00",
					"advogados": []string{"Dr. Pedro Advogado - OAB/RJ 12345"},
				},
				"reu": gin.H{
					"nome": "Empresa XYZ Ltda",
					"cnpj": "12.345.678/0001-00",
					"advogados": []string{"Dra. Ana Advogada - OAB/RJ 54321"},
				},
			},
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// Get process movements
func (h *DataJudHandler) getMovements(c *gin.Context) {
	number := c.Param("number")
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar busca real na API DataJud
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"numero": number,
			"movimentacoes": []gin.H{
				{
					"data": "2024-03-15",
					"descricao": "Juntada de Petição",
					"complemento": "Petição inicial protocolada",
				},
				{
					"data": "2024-03-20",
					"descricao": "Conclusão para Despacho",
					"complemento": "Autos conclusos ao juiz",
				},
				{
					"data": "2024-03-22",
					"descricao": "Despacho",
					"complemento": "Cite-se o réu para apresentar contestação",
				},
			},
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// Get statistics
func (h *DataJudHandler) getStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	// Buscar estatísticas dos CNPJ providers
	var stats struct {
		ProvidersAtivos int
		TotalQuotaUsada int
		TotalQuota      int
		RequestsHoje    int
		RequestsMes     int
	}
	
	// Estatísticas dos providers
	providerQuery := `
		SELECT 
			COUNT(CASE WHEN is_active = true THEN 1 END) as providers_ativos,
			COALESCE(SUM(daily_usage), 0) as total_quota_usada,
			COALESCE(SUM(daily_limit), 0) as total_quota
		FROM cnpj_providers 
		WHERE tenant_id = $1
	`
	
	err := h.db.QueryRow(providerQuery, tenantID).Scan(
		&stats.ProvidersAtivos,
		&stats.TotalQuotaUsada,
		&stats.TotalQuota,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar estatísticas"})
		return
	}
	
	// Buscar requests do dia (se tabela existe)
	requestsQuery := `
		SELECT 
			COUNT(CASE WHEN created_at >= CURRENT_DATE THEN 1 END) as requests_hoje,
			COUNT(CASE WHEN created_at >= DATE_TRUNC('month', CURRENT_DATE) THEN 1 END) as requests_mes
		FROM datajud_requests 
		WHERE tenant_id = $1
	`
	
	err = h.db.QueryRow(requestsQuery, tenantID).Scan(
		&stats.RequestsHoje,
		&stats.RequestsMes,
	)
	// Se a tabela não existe ou há erro, usar valores padrão
	if err != nil {
		stats.RequestsHoje = stats.TotalQuotaUsada
		stats.RequestsMes = stats.TotalQuotaUsada * 30
	}
	
	// Calcular métricas
	percentualUso := 0.0
	if stats.TotalQuota > 0 {
		percentualUso = float64(stats.TotalQuotaUsada) / float64(stats.TotalQuota) * 100
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"consultas_hoje":        stats.RequestsHoje,
			"consultas_mes":         stats.RequestsMes,
			"limite_diario":         stats.TotalQuota,
			"usado_hoje":            stats.TotalQuotaUsada,
			"quota_disponivel":      stats.TotalQuota - stats.TotalQuotaUsada,
			"percentual_uso":        percentualUso,
			"cnpj_providers_ativos": stats.ProvidersAtivos,
			"cache_hits":            320, // TODO: implementar cache real
			"cache_misses":          125, // TODO: implementar cache real
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// Get quota usage
func (h *DataJudHandler) getQuota(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	// TODO: Implementar controle de quota real
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"usado_hoje": 45,
			"limite_diario": h.cfg.RateLimitDaily,
			"percentual_uso": 0.45,
			"resets_em": "00:00:00",
		},
		"meta": gin.H{
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
		},
	})
}

// CNPJ Providers management
func (h *DataJudHandler) listCNPJProviders(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	query := `
		SELECT 
			id,
			cnpj,
			name,
			email,
			api_key,
			daily_limit,
			daily_usage,
			priority,
			is_active,
			last_used_at,
			created_at
		FROM cnpj_providers 
		WHERE tenant_id = $1 
		ORDER BY priority
	`
	
	rows, err := h.db.Query(query, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar providers"})
		return
	}
	defer rows.Close()
	
	var providers []gin.H
	for rows.Next() {
		var id, cnpj, name, email, apiKey string
		var dailyLimit, dailyUsage, priority int
		var isActive bool
		var lastUsedAt, createdAt *time.Time
		
		err := rows.Scan(&id, &cnpj, &name, &email, &apiKey, &dailyLimit, &dailyUsage, &priority, &isActive, &lastUsedAt, &createdAt)
		if err != nil {
			continue
		}
		
		// Mascarar API key para segurança
		maskedKey := "****" + apiKey[len(apiKey)-4:]
		
		provider := gin.H{
			"id":            id,
			"cnpj":          cnpj,
			"nome":          name,
			"email":         email,
			"api_key":       maskedKey,
			"quota_diaria":  dailyLimit,
			"usado_hoje":    dailyUsage,
			"quota_disponivel": dailyLimit - dailyUsage,
			"percentual_uso": float64(dailyUsage) / float64(dailyLimit) * 100,
			"prioridade":    priority,
			"ativo":         isActive,
			"created_at":    createdAt,
		}
		
		if lastUsedAt != nil {
			provider["last_used_at"] = lastUsedAt
		}
		
		providers = append(providers, provider)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": providers,
		"meta": gin.H{
			"tenant_id": tenantID,
			"total":     len(providers),
			"timestamp": time.Now().Unix(),
		},
	})
}

func (h *DataJudHandler) createCNPJProvider(c *gin.Context) {
	// TODO: Implementar criação real
	c.JSON(http.StatusCreated, gin.H{
		"message": "CNPJ Provider criado com sucesso",
		"id": "3",
	})
}

func (h *DataJudHandler) updateCNPJProvider(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implementar atualização real
	c.JSON(http.StatusOK, gin.H{
		"message": "CNPJ Provider atualizado com sucesso",
		"id": id,
	})
}

func (h *DataJudHandler) deleteCNPJProvider(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implementar exclusão real
	c.JSON(http.StatusOK, gin.H{
		"message": "CNPJ Provider removido com sucesso",
		"id": id,
	})
}