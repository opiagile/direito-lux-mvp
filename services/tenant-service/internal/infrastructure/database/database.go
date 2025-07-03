package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/direito-lux/tenant-service/internal/infrastructure/config"
)

// NewPostgreSQLConnection cria nova conexão PostgreSQL
func NewPostgreSQLConnection(cfg *config.Config, logger *zap.Logger) (*sqlx.DB, error) {
	logger.Info("📊 Connecting to PostgreSQL database...",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
		zap.String("user", cfg.Database.User),
	)

	// Create connection string
	dsn := cfg.GetDatabaseDSN()
	
	// Connect to database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("✅ PostgreSQL database connected successfully",
		zap.Int("max_open_conns", cfg.Database.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Database.MaxIdleConns),
		zap.Duration("conn_max_lifetime", cfg.Database.ConnMaxLifetime),
	)

	return db, nil
}

// TestConnection testa a conexão com o banco
func TestConnection(db *sqlx.DB, logger *zap.Logger) error {
	logger.Info("🔍 Testing database connection...")

	// Test basic connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test basic query
	var result int
	if err := db.Get(&result, "SELECT 1"); err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	// Test tenants table exists
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM tenants"); err != nil {
		logger.Warn("⚠️ Tenants table not accessible", zap.Error(err))
		return fmt.Errorf("tenants table not accessible: %w", err)
	}

	logger.Info("✅ Database connection test successful",
		zap.Int("tenants_count", count),
	)

	return nil
}

// CloseConnection fecha a conexão com o banco
func CloseConnection(db *sqlx.DB, logger *zap.Logger) error {
	if db == nil {
		return nil
	}

	logger.Info("📊 Closing database connection...")
	
	if err := db.Close(); err != nil {
		logger.Error("❌ Error closing database connection", zap.Error(err))
		return err
	}

	logger.Info("✅ Database connection closed successfully")
	return nil
}