package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
)

// Connection estrutura para conexão com banco de dados
type Connection struct {
	db     *sqlx.DB
	config *config.Config
	logger *zap.Logger
}

// NewConnection cria uma nova conexão com o banco de dados
func NewConnection(cfg *config.Config, logger *zap.Logger) (*Connection, error) {
	// Conectar ao banco
	db, err := sqlx.Connect("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Testar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão: %w", err)
	}

	conn := &Connection{
		db:     db,
		config: cfg,
		logger: logger,
	}

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

// QueryRowContext executa uma query que retorna uma linha
func (c *Connection) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// QueryRowxContext executa uma query que retorna uma linha usando sqlx
func (c *Connection) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return c.db.QueryRowxContext(ctx, query, args...)
}

// ExecContext executa uma query de modificação
func (c *Connection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

// GetContext busca um registro usando sqlx
func (c *Connection) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

// SelectContext busca múltiplos registros usando sqlx
func (c *Connection) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

// Transaction executa uma função dentro de uma transação
func (c *Connection) Transaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
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
	return err
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

// Health alias for HealthCheck - maintains compatibility
func (c *Connection) Health(ctx context.Context) error {
	return c.HealthCheck(ctx)
}