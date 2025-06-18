package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/direito-lux/notification-service/internal/infrastructure/config"
)

// Connection representa uma conexão com PostgreSQL
type Connection struct {
	DB *sqlx.DB
}

// NewConnection cria nova conexão com PostgreSQL
func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no PostgreSQL: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pingar PostgreSQL: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close fecha a conexão
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// Health verifica saúde da conexão
func (c *Connection) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.DB.PingContext(ctx)
}

// GetDB retorna instância do sqlx.DB
func (c *Connection) GetDB() *sqlx.DB {
	return c.DB
}