-- Migration: Create processes table
-- Description: Tabela principal para armazenar processos jurídicos
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- Criar extension para UUID se não existir
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar tabela de processos
CREATE TABLE IF NOT EXISTS processes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    client_id UUID NOT NULL,
    number VARCHAR(30) UNIQUE NOT NULL,
    original_number VARCHAR(30),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'suspended', 'finished', 'archived', 'canceled')),
    stage VARCHAR(20) NOT NULL CHECK (stage IN ('knowledge', 'execution', 'appeal', 'emergency', 'preventive', 'incidental')),
    subject JSONB NOT NULL,
    value JSONB,
    court_id VARCHAR(50) NOT NULL,
    judge_id VARCHAR(50),
    monitoring JSONB NOT NULL DEFAULT '{"enabled": false, "notification_channels": [], "keywords": [], "auto_sync": false, "sync_interval_hours": 24}',
    tags TEXT[] DEFAULT '{}',
    custom_fields JSONB DEFAULT '{}',
    last_movement_at TIMESTAMP,
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT check_number_format CHECK (
        number ~ '^[0-9]{7}-[0-9]{2}\.[0-9]{4}\.[0-9]{1}\.[0-9]{2}\.[0-9]{4}$'
    ),
    CONSTRAINT check_archived_consistency CHECK (
        (status = 'archived' AND archived_at IS NOT NULL) OR 
        (status != 'archived' AND archived_at IS NULL)
    )
);

-- Comentários nas colunas
COMMENT ON TABLE processes IS 'Tabela principal de processos jurídicos';
COMMENT ON COLUMN processes.id IS 'Identificador único do processo';
COMMENT ON COLUMN processes.tenant_id IS 'Identificador do tenant proprietário';
COMMENT ON COLUMN processes.client_id IS 'Identificador do cliente do processo';
COMMENT ON COLUMN processes.number IS 'Número CNJ do processo (formato: NNNNNNN-DD.AAAA.J.TR.OOOO)';
COMMENT ON COLUMN processes.original_number IS 'Número original antes da unificação CNJ';
COMMENT ON COLUMN processes.title IS 'Título/assunto principal do processo';
COMMENT ON COLUMN processes.description IS 'Descrição detalhada do processo';
COMMENT ON COLUMN processes.status IS 'Status atual do processo';
COMMENT ON COLUMN processes.stage IS 'Fase processual atual';
COMMENT ON COLUMN processes.subject IS 'Assunto do processo baseado na tabela CNJ (JSON)';
COMMENT ON COLUMN processes.value IS 'Valor da causa (JSON com amount e currency)';
COMMENT ON COLUMN processes.court_id IS 'Identificador do tribunal';
COMMENT ON COLUMN processes.judge_id IS 'Identificador do juiz responsável';
COMMENT ON COLUMN processes.monitoring IS 'Configurações de monitoramento (JSON)';
COMMENT ON COLUMN processes.tags IS 'Tags/etiquetas para categorização';
COMMENT ON COLUMN processes.custom_fields IS 'Campos customizados específicos do tenant (JSON)';
COMMENT ON COLUMN processes.last_movement_at IS 'Data da última movimentação';
COMMENT ON COLUMN processes.last_sync_at IS 'Data da última sincronização com DataJud';
COMMENT ON COLUMN processes.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN processes.updated_at IS 'Data da última atualização';
COMMENT ON COLUMN processes.archived_at IS 'Data de arquivamento do processo';

-- Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_processes_updated_at 
    BEFORE UPDATE ON processes 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();