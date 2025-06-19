-- Migration: 001_create_search_indices_table
-- Description: Criar tabelas para gerenciar índices de busca

-- Tabela para gerenciar índices do Elasticsearch
CREATE TABLE IF NOT EXISTS search_indices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    mapping JSONB NOT NULL,
    settings JSONB DEFAULT '{}',
    aliases TEXT[] DEFAULT '{}',
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Estatísticas
    document_count BIGINT DEFAULT 0,
    size_in_bytes BIGINT DEFAULT 0,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadados
    tenant_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID,
    updated_by UUID
);

-- Tabela para logs de indexação
CREATE TABLE IF NOT EXISTS search_indexing_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    index_name VARCHAR(255) NOT NULL,
    operation VARCHAR(50) NOT NULL, -- 'index', 'update', 'delete', 'bulk'
    document_id VARCHAR(255),
    document_type VARCHAR(100),
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'success', 'error'
    error_message TEXT,
    
    -- Métricas
    processing_time_ms INTEGER,
    documents_processed INTEGER DEFAULT 1,
    
    -- Dados do documento (opcional, para reprocessamento)
    document_data JSONB,
    
    -- Metadados
    tenant_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Tabela para estatísticas de busca
CREATE TABLE IF NOT EXISTS search_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    hour INTEGER NOT NULL CHECK (hour >= 0 AND hour <= 23),
    
    -- Métricas de busca
    total_searches BIGINT DEFAULT 0,
    successful_searches BIGINT DEFAULT 0,
    failed_searches BIGINT DEFAULT 0,
    avg_response_time_ms NUMERIC(10,2) DEFAULT 0,
    
    -- Tipos de busca
    basic_searches BIGINT DEFAULT 0,
    advanced_searches BIGINT DEFAULT 0,
    suggestions_requests BIGINT DEFAULT 0,
    
    -- Por índice
    index_name VARCHAR(255),
    
    -- Metadados
    tenant_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE (date, hour, index_name, tenant_id)
);

-- Tabela para cache de buscas
CREATE TABLE IF NOT EXISTS search_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    query_hash VARCHAR(64) NOT NULL,
    
    -- Dados da busca
    search_query JSONB NOT NULL,
    search_results JSONB NOT NULL,
    total_results INTEGER DEFAULT 0,
    
    -- Metadados de cache
    hit_count INTEGER DEFAULT 0,
    last_accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Contexto
    tenant_id UUID,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    INDEX idx_search_cache_key (cache_key),
    INDEX idx_search_cache_expires (expires_at),
    INDEX idx_search_cache_tenant (tenant_id)
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_search_indices_name ON search_indices(name);
CREATE INDEX IF NOT EXISTS idx_search_indices_tenant ON search_indices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_search_indices_active ON search_indices(is_active);

CREATE INDEX IF NOT EXISTS idx_search_indexing_logs_index ON search_indexing_logs(index_name);
CREATE INDEX IF NOT EXISTS idx_search_indexing_logs_status ON search_indexing_logs(status);
CREATE INDEX IF NOT EXISTS idx_search_indexing_logs_created ON search_indexing_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_search_indexing_logs_tenant ON search_indexing_logs(tenant_id);

CREATE INDEX IF NOT EXISTS idx_search_statistics_date ON search_statistics(date);
CREATE INDEX IF NOT EXISTS idx_search_statistics_tenant ON search_statistics(tenant_id);
CREATE INDEX IF NOT EXISTS idx_search_statistics_index ON search_statistics(index_name);

-- Triggers para updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_search_indices_updated_at 
    BEFORE UPDATE ON search_indices
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_search_statistics_updated_at 
    BEFORE UPDATE ON search_statistics
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Função para limpeza automática de cache expirado
CREATE OR REPLACE FUNCTION cleanup_expired_search_cache()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM search_cache WHERE expires_at < NOW();
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Comentários nas tabelas
COMMENT ON TABLE search_indices IS 'Gerenciamento de índices do Elasticsearch';
COMMENT ON TABLE search_indexing_logs IS 'Logs de operações de indexação';
COMMENT ON TABLE search_statistics IS 'Estatísticas de uso do search service';
COMMENT ON TABLE search_cache IS 'Cache de resultados de busca';