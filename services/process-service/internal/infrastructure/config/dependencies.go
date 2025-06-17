package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/streadway/amqp"
	"github.com/direito-lux/process-service/internal/application"
	"github.com/direito-lux/process-service/internal/domain"
	"github.com/direito-lux/process-service/internal/infrastructure/postgres"
)

// Dependencies container com todas as dependências
type Dependencies struct {
	// Repositories
	ProcessRepo  domain.ProcessRepository
	MovementRepo domain.MovementRepository
	PartyRepo    domain.PartyRepository

	// Event Publisher
	EventPublisher domain.EventPublisher

	// Application Service
	ProcessService *application.ProcessService

	// Infrastructure
	DB             *sql.DB
	RabbitMQConn   *amqp.Connection
}

// NewDependencies cria e configura todas as dependências
func NewDependencies(cfg *Config) (*Dependencies, error) {
	deps := &Dependencies{}

	// Configurar banco de dados
	if err := deps.setupDatabase(cfg); err != nil {
		return nil, fmt.Errorf("erro ao configurar banco: %w", err)
	}

	// Configurar RabbitMQ
	if err := deps.setupRabbitMQ(cfg); err != nil {
		return nil, fmt.Errorf("erro ao configurar RabbitMQ: %w", err)
	}

	// Configurar repositórios
	deps.setupRepositories()

	// Configurar event publisher
	if err := deps.setupEventPublisher(cfg); err != nil {
		return nil, fmt.Errorf("erro ao configurar event publisher: %w", err)
	}

	// Configurar application service
	deps.setupApplicationService()

	log.Println("Dependências configuradas com sucesso")
	return deps, nil
}

// setupDatabase configura conexão com PostgreSQL
func (d *Dependencies) setupDatabase(cfg *Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return fmt.Errorf("erro ao testar conexão: %w", err)
	}

	d.DB = db
	log.Printf("Conectado ao PostgreSQL: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return nil
}

// setupRabbitMQ configura conexão com RabbitMQ
func (d *Dependencies) setupRabbitMQ(cfg *Config) error {
	connStr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.VHost,
	)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return fmt.Errorf("erro ao conectar RabbitMQ: %w", err)
	}

	d.RabbitMQConn = conn
	log.Printf("Conectado ao RabbitMQ: %s:%d", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	return nil
}

// setupRepositories configura repositórios
func (d *Dependencies) setupRepositories() {
	d.ProcessRepo = postgres.NewProcessRepository(d.DB)
	d.MovementRepo = postgres.NewMovementRepository(d.DB)
	d.PartyRepo = postgres.NewPartyRepository(d.DB)
	
	log.Println("Repositórios configurados")
}

// setupEventPublisher configura event publisher
func (d *Dependencies) setupEventPublisher(cfg *Config) error {
	// TODO: Implement event publisher setup
	// For now, just return nil to break import cycle
	log.Println("Event publisher configurado (placeholder)")
	return nil
}

// setupApplicationService configura application service
func (d *Dependencies) setupApplicationService() {
	d.ProcessService = application.NewProcessService(
		d.ProcessRepo,
		d.MovementRepo,
		d.PartyRepo,
		d.EventPublisher,
	)
	
	log.Println("Application service configurado")
}

// Close fecha todas as conexões
func (d *Dependencies) Close() error {
	var errors []error

	// Fechar event publisher
	if closer, ok := d.EventPublisher.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			errors = append(errors, fmt.Errorf("erro ao fechar event publisher: %w", err))
		}
	}

	// Fechar RabbitMQ
	if d.RabbitMQConn != nil {
		if err := d.RabbitMQConn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("erro ao fechar RabbitMQ: %w", err))
		}
	}

	// Fechar banco de dados
	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("erro ao fechar banco: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("erros ao fechar dependências: %v", errors)
	}

	log.Println("Dependências fechadas com sucesso")
	return nil
}

// Health verifica saúde das dependências
func (d *Dependencies) Health() error {
	// Verificar banco de dados
	if err := d.DB.Ping(); err != nil {
		return fmt.Errorf("banco indisponível: %w", err)
	}

	// Verificar RabbitMQ
	if d.RabbitMQConn.IsClosed() {
		return fmt.Errorf("RabbitMQ desconectado")
	}

	return nil
}

// GetMetrics retorna métricas das dependências
func (d *Dependencies) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Métricas do banco
	if d.DB != nil {
		stats := d.DB.Stats()
		metrics["database"] = map[string]interface{}{
			"open_connections": stats.OpenConnections,
			"idle_connections": stats.Idle,
			"in_use":          stats.InUse,
			"wait_count":      stats.WaitCount,
			"wait_duration":   stats.WaitDuration.String(),
		}
	}

	// Métricas de eventos (placeholder)
	metrics["events"] = map[string]interface{}{
		"total_published": 0,
		"total_failed":    0,
	}

	// Métricas do RabbitMQ
	metrics["rabbitmq"] = map[string]interface{}{
		"connected": !d.RabbitMQConn.IsClosed(),
	}

	return metrics
}

// Repository patterns para testes
type TestDependencies struct {
	*Dependencies
	// Mocks para testes podem ser adicionados aqui
}

// NewTestDependencies cria dependências para testes
func NewTestDependencies() *TestDependencies {
	// Em ambiente de teste, pode usar mocks ou banco em memória
	return &TestDependencies{
		Dependencies: &Dependencies{
			// Configurar mocks aqui
		},
	}
}

// DatabaseMigrator gerenciador de migrações
type DatabaseMigrator struct {
	db *sql.DB
}

// NewDatabaseMigrator cria novo migrator
func NewDatabaseMigrator(db *sql.DB) *DatabaseMigrator {
	return &DatabaseMigrator{db: db}
}

// Migrate executa migrações
func (m *DatabaseMigrator) Migrate() error {
	// Lista de migrações em ordem
	migrations := []Migration{
		{
			Version:     "001",
			Description: "Create processes table",
			SQL:         createProcessesTableSQL,
		},
		{
			Version:     "002",
			Description: "Create movements table",
			SQL:         createMovementsTableSQL,
		},
		{
			Version:     "003",
			Description: "Create parties table",
			SQL:         createPartiesTableSQL,
		},
		{
			Version:     "004",
			Description: "Create indexes",
			SQL:         createIndexesSQL,
		},
	}

	// Criar tabela de migrações se não existir
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("erro ao criar tabela de migrações: %w", err)
	}

	// Executar migrações pendentes
	for _, migration := range migrations {
		if err := m.runMigration(migration); err != nil {
			return fmt.Errorf("erro na migração %s: %w", migration.Version, err)
		}
	}

	log.Println("Migrações executadas com sucesso")
	return nil
}

// Migration representa uma migração
type Migration struct {
	Version     string
	Description string
	SQL         string
}

// createMigrationsTable cria tabela de controle de migrações
func (m *DatabaseMigrator) createMigrationsTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			description TEXT,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := m.db.Exec(sql)
	return err
}

// runMigration executa uma migração
func (m *DatabaseMigrator) runMigration(migration Migration) error {
	// Verificar se já foi executada
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", migration.Version).Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar migração: %w", err)
	}

	if count > 0 {
		log.Printf("Migração %s já executada", migration.Version)
		return nil
	}

	// Executar migração em transação
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Executar SQL da migração
	if _, err := tx.Exec(migration.SQL); err != nil {
		return fmt.Errorf("erro ao executar SQL: %w", err)
	}

	// Registrar migração como executada
	_, err = tx.Exec(
		"INSERT INTO schema_migrations (version, description) VALUES ($1, $2)",
		migration.Version,
		migration.Description,
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar migração: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar migração: %w", err)
	}

	log.Printf("Migração %s executada: %s", migration.Version, migration.Description)
	return nil
}

// SQLs das migrações
const createProcessesTableSQL = `
CREATE TABLE IF NOT EXISTS processes (
	id UUID PRIMARY KEY,
	tenant_id UUID NOT NULL,
	client_id UUID NOT NULL,
	number VARCHAR(30) UNIQUE NOT NULL,
	original_number VARCHAR(30),
	title VARCHAR(500) NOT NULL,
	description TEXT,
	status VARCHAR(20) NOT NULL,
	stage VARCHAR(20) NOT NULL,
	subject JSONB NOT NULL,
	value JSONB,
	court_id VARCHAR(50) NOT NULL,
	judge_id VARCHAR(50),
	monitoring JSONB NOT NULL DEFAULT '{}',
	tags TEXT[] DEFAULT '{}',
	custom_fields JSONB DEFAULT '{}',
	last_movement_at TIMESTAMP,
	last_sync_at TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	archived_at TIMESTAMP
);
`

const createMovementsTableSQL = `
CREATE TABLE IF NOT EXISTS movements (
	id UUID PRIMARY KEY,
	process_id UUID NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
	tenant_id UUID NOT NULL,
	sequence INTEGER NOT NULL,
	external_id VARCHAR(100),
	date TIMESTAMP NOT NULL,
	type VARCHAR(20) NOT NULL,
	code VARCHAR(20) NOT NULL,
	title VARCHAR(500) NOT NULL,
	description TEXT NOT NULL,
	content TEXT,
	judge VARCHAR(200),
	responsible VARCHAR(200),
	attachments JSONB DEFAULT '[]',
	related_parties TEXT[] DEFAULT '{}',
	is_important BOOLEAN DEFAULT FALSE,
	is_public BOOLEAN DEFAULT TRUE,
	notification_sent BOOLEAN DEFAULT FALSE,
	tags TEXT[] DEFAULT '{}',
	metadata JSONB DEFAULT '{}',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	synced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	UNIQUE(process_id, sequence),
	UNIQUE(external_id) WHERE external_id IS NOT NULL
);
`

const createPartiesTableSQL = `
CREATE TABLE IF NOT EXISTS parties (
	id UUID PRIMARY KEY,
	process_id UUID NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
	type VARCHAR(20) NOT NULL,
	name VARCHAR(500) NOT NULL,
	document VARCHAR(20),
	document_type VARCHAR(10),
	role VARCHAR(20) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE,
	lawyer JSONB,
	contact JSONB DEFAULT '{}',
	address JSONB DEFAULT '{}',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

const createIndexesSQL = `
-- Índices para processes
CREATE INDEX IF NOT EXISTS idx_processes_tenant_id ON processes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_processes_client_id ON processes(client_id);
CREATE INDEX IF NOT EXISTS idx_processes_status ON processes(status);
CREATE INDEX IF NOT EXISTS idx_processes_court_id ON processes(court_id);
CREATE INDEX IF NOT EXISTS idx_processes_created_at ON processes(created_at);
CREATE INDEX IF NOT EXISTS idx_processes_updated_at ON processes(updated_at);
CREATE INDEX IF NOT EXISTS idx_processes_last_movement_at ON processes(last_movement_at);
CREATE INDEX IF NOT EXISTS idx_processes_monitoring_enabled ON processes((monitoring->>'enabled')) WHERE monitoring->>'enabled' = 'true';

-- Índices para movements
CREATE INDEX IF NOT EXISTS idx_movements_process_id ON movements(process_id);
CREATE INDEX IF NOT EXISTS idx_movements_tenant_id ON movements(tenant_id);
CREATE INDEX IF NOT EXISTS idx_movements_date ON movements(date);
CREATE INDEX IF NOT EXISTS idx_movements_type ON movements(type);
CREATE INDEX IF NOT EXISTS idx_movements_sequence ON movements(process_id, sequence);
CREATE INDEX IF NOT EXISTS idx_movements_important ON movements(is_important) WHERE is_important = true;
CREATE INDEX IF NOT EXISTS idx_movements_notification_pending ON movements(is_important, notification_sent) WHERE is_important = true AND notification_sent = false;
CREATE INDEX IF NOT EXISTS idx_movements_external_id ON movements(external_id) WHERE external_id IS NOT NULL;

-- Índices para parties
CREATE INDEX IF NOT EXISTS idx_parties_process_id ON parties(process_id);
CREATE INDEX IF NOT EXISTS idx_parties_document ON parties(document) WHERE document IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_parties_role ON parties(role);
CREATE INDEX IF NOT EXISTS idx_parties_active ON parties(is_active) WHERE is_active = true;

-- Índices para busca textual
CREATE INDEX IF NOT EXISTS idx_processes_search ON processes USING gin(to_tsvector('portuguese', title || ' ' || description));
CREATE INDEX IF NOT EXISTS idx_movements_search ON movements USING gin(to_tsvector('portuguese', title || ' ' || description || ' ' || COALESCE(content, '')));
CREATE INDEX IF NOT EXISTS idx_parties_search ON parties USING gin(to_tsvector('portuguese', name));
`