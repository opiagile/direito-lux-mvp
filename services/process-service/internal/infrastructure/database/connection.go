package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/direito-lux/process-service/internal/infrastructure/config"
	"github.com/direito-lux/process-service/internal/infrastructure/logging"
	"github.com/direito-lux/process-service/internal/infrastructure/metrics"
	"github.com/direito-lux/process-service/internal/infrastructure/tracing"
)

// Connection estrutura para conexão com banco de dados
type Connection struct {
	db      *sqlx.DB
	config  *config.Config
	logger  *zap.Logger
	metrics *metrics.Metrics
}

// NewConnection cria uma nova conexão com o banco de dados
func NewConnection(cfg *config.Config, logger *zap.Logger, metrics *metrics.Metrics) (*Connection, error) {
	// Conectar ao banco
	db, err := sqlx.Connect("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Testar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão: %w", err)
	}

	conn := &Connection{
		db:      db,
		config:  cfg,
		logger:  logger,
		metrics: metrics,
	}

	// Executar migrações se em desenvolvimento
	if cfg.IsDevelopment() {
		if err := conn.runMigrations(); err != nil {
			logger.Warn("Erro ao executar migrações", zap.Error(err))
		}
	}

	// Iniciar coleta de métricas de conexão
	go conn.collectConnectionMetrics()

	logger.Info("Conexão com banco de dados estabelecida",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
		zap.Int("max_open_conns", cfg.Database.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Database.MaxIdleConns),
	)

	return conn, nil
}

// GetDB retorna a instância do banco de dados
func (c *Connection) GetDB() *sqlx.DB {
	return c.db
}

// Close fecha a conexão com o banco
func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// runMigrations executa as migrações do banco
func (c *Connection) runMigrations() error {
	if c.config.Database.MigrationsPath == "" {
		return nil
	}

	driver, err := postgres.WithInstance(c.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("erro ao criar driver de migração: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		c.config.Database.MigrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar instância de migração: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("erro ao executar migrações: %w", err)
	}

	c.logger.Info("Migrações executadas com sucesso")
	return nil
}

// collectConnectionMetrics coleta métricas de conexão periodicamente
func (c *Connection) collectConnectionMetrics() {
	if c.metrics == nil {
		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := c.db.Stats()
		
		if c.metrics.DatabaseConnections != nil {
			c.metrics.DatabaseConnections.WithLabelValues("open").Set(float64(stats.OpenConnections))
			c.metrics.DatabaseConnections.WithLabelValues("idle").Set(float64(stats.Idle))
			c.metrics.DatabaseConnections.WithLabelValues("in_use").Set(float64(stats.InUse))
		}
	}
}

// QueryContext executa uma query com contexto e métricas
func (c *Connection) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.queryWithMetrics(ctx, "SELECT", query, func() (*sql.Rows, error) {
		return c.db.QueryContext(ctx, query, args...)
	})
}

// QueryxContext executa uma query com contexto usando sqlx
func (c *Connection) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return c.queryxWithMetrics(ctx, "SELECT", query, func() (*sqlx.Rows, error) {
		return c.db.QueryxContext(ctx, query, args...)
	})
}

// QueryRowContext executa uma query que retorna uma linha
func (c *Connection) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, "SELECT", "unknown")
	defer span.Finish()

	row := c.db.QueryRowContext(ctx, query, args...)
	
	c.recordQueryMetrics(ctx, "SELECT", "unknown", time.Since(start), nil)
	
	return row
}

// QueryRowxContext executa uma query que retorna uma linha usando sqlx
func (c *Connection) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, "SELECT", "unknown")
	defer span.Finish()

	row := c.db.QueryRowxContext(ctx, query, args...)
	
	c.recordQueryMetrics(ctx, "SELECT", "unknown", time.Since(start), nil)
	
	return row
}

// ExecContext executa uma query de modificação (INSERT, UPDATE, DELETE)
func (c *Connection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	operation := "EXEC"
	table := "unknown"
	
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, operation, table)
	defer span.Finish()

	result, err := c.db.ExecContext(ctx, query, args...)
	
	c.recordQueryMetrics(ctx, operation, table, time.Since(start), err)
	
	if err != nil {
		tracing.SetSpanError(span, err)
	}
	
	return result, err
}

// GetContext busca um registro usando sqlx
func (c *Connection) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, "GET", "unknown")
	defer span.Finish()

	err := c.db.GetContext(ctx, dest, query, args...)
	
	c.recordQueryMetrics(ctx, "GET", "unknown", time.Since(start), err)
	
	if err != nil {
		tracing.SetSpanError(span, err)
	}
	
	return err
}

// SelectContext busca múltiplos registros usando sqlx
func (c *Connection) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, "SELECT", "unknown")
	defer span.Finish()

	err := c.db.SelectContext(ctx, dest, query, args...)
	
	c.recordQueryMetrics(ctx, "SELECT", "unknown", time.Since(start), err)
	
	if err != nil {
		tracing.SetSpanError(span, err)
	}
	
	return err
}

// NamedExecContext executa query com parâmetros nomeados
func (c *Connection) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, "NAMED_EXEC", "unknown")
	defer span.Finish()

	result, err := c.db.NamedExecContext(ctx, query, arg)
	
	c.recordQueryMetrics(ctx, "NAMED_EXEC", "unknown", time.Since(start), err)
	
	if err != nil {
		tracing.SetSpanError(span, err)
	}
	
	return result, err
}

// Transaction executa uma função dentro de uma transação
func (c *Connection) Transaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	span, ctx := tracing.TraceDatabase(ctx, "TRANSACTION", "multiple")
	defer span.Finish()

	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		tracing.SetSpanError(span, err)
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	if err != nil {
		tracing.SetSpanError(span, err)
	}

	return err
}

// queryWithMetrics executa query com métricas
func (c *Connection) queryWithMetrics(ctx context.Context, operation, query string, fn func() (*sql.Rows, error)) (*sql.Rows, error) {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, operation, "unknown")
	defer span.Finish()

	rows, err := fn()
	
	c.recordQueryMetrics(ctx, operation, "unknown", time.Since(start), err)
	
	if err != nil {
		tracing.SetSpanError(span, err)
	}
	
	return rows, err
}

// queryxWithMetrics executa query sqlx com métricas
func (c *Connection) queryxWithMetrics(ctx context.Context, operation, query string, fn func() (*sqlx.Rows, error)) (*sqlx.Rows, error) {
	start := time.Now()
	span, ctx := tracing.TraceDatabase(ctx, operation, "unknown")
	defer span.Finish()

	rows, err := fn()
	
	c.recordQueryMetrics(ctx, operation, "unknown", time.Since(start), err)
	
	if err != nil {
		tracing.SetSpanError(span, err)
	}
	
	return rows, err
}

// recordQueryMetrics registra métricas da query
func (c *Connection) recordQueryMetrics(ctx context.Context, operation, table string, duration time.Duration, err error) {
	if c.metrics != nil {
		tenantID := logging.GetTenantID(ctx)
		c.metrics.RecordDatabaseQuery(operation, table, tenantID, duration, err)
	}
}

// HealthCheck verifica a saúde da conexão
func (c *Connection) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result int
	err := c.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("health check falhou: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("health check retornou valor inesperado: %d", result)
	}

	return nil
}

// SetTenantSchema define o schema do tenant no contexto da sessão
func (c *Connection) SetTenantSchema(ctx context.Context, tenantID string) error {
	schema := fmt.Sprintf("tenant_%s", tenantID)
	
	// Validar nome do schema para evitar SQL injection
	if !isValidSchemaName(schema) {
		return fmt.Errorf("nome do schema inválido: %s", schema)
	}

	_, err := c.db.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s, public", schema))
	if err != nil {
		return fmt.Errorf("erro ao definir schema do tenant: %w", err)
	}

	return nil
}

// CreateTenantSchema cria schema para um novo tenant
func (c *Connection) CreateTenantSchema(ctx context.Context, tenantID string) error {
	schema := fmt.Sprintf("tenant_%s", tenantID)
	
	// Validar nome do schema
	if !isValidSchemaName(schema) {
		return fmt.Errorf("nome do schema inválido: %s", schema)
	}

	return c.Transaction(ctx, func(tx *sqlx.Tx) error {
		// Criar schema
		_, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
		if err != nil {
			return fmt.Errorf("erro ao criar schema: %w", err)
		}

		// Aqui você adicionaria a criação das tabelas específicas do tenant
		// Por enquanto, apenas o schema vazio

		return nil
	})
}

// isValidSchemaName valida se o nome do schema é seguro
func isValidSchemaName(name string) bool {
	// Implementar validação adequada
	// Por simplicidade, verificar apenas caracteres alfanuméricos e underscore
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 63 // Limite do PostgreSQL
}