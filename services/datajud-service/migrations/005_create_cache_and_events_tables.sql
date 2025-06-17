-- Migration: Create cache and events tables
-- Description: Tabelas para cache de consultas e eventos de domínio
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- TABELA DE CACHE
-- ========================================

CREATE TABLE cache_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value JSONB NOT NULL,
    ttl INTEGER NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    size_bytes INTEGER NOT NULL DEFAULT 0,
    hit_count INTEGER NOT NULL DEFAULT 0,
    tenant_id UUID NOT NULL,
    request_type VARCHAR(20) NOT NULL,
    process_number VARCHAR(25),
    court_id VARCHAR(10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_hit_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT chk_cache_ttl_positive CHECK (ttl > 0),
    CONSTRAINT chk_cache_expires_future CHECK (expires_at > created_at),
    CONSTRAINT chk_cache_size_non_negative CHECK (size_bytes >= 0),
    CONSTRAINT chk_cache_hit_count_non_negative CHECK (hit_count >= 0),
    CONSTRAINT chk_cache_request_type CHECK (request_type IN ('process', 'movement', 'party', 'document', 'bulk'))
);

-- ========================================
-- TABELA DE EVENTOS DE DOMÍNIO
-- ========================================

CREATE TABLE domain_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    event_data JSONB NOT NULL,
    event_version INTEGER NOT NULL DEFAULT 1,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE,
    is_processed BOOLEAN NOT NULL DEFAULT false,
    correlation_id UUID,
    causation_id UUID,
    tenant_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_event_version_positive CHECK (event_version > 0),
    CONSTRAINT chk_occurred_at_not_future CHECK (occurred_at <= NOW()),
    CONSTRAINT chk_processed_at_after_occurred CHECK (processed_at IS NULL OR processed_at >= occurred_at)
);

-- ========================================
-- TABELA DE ESTATÍSTICAS DE CACHE
-- ========================================

CREATE TABLE cache_statistics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID,
    request_type VARCHAR(20),
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    total_requests INTEGER NOT NULL DEFAULT 0,
    cache_hits INTEGER NOT NULL DEFAULT 0,
    cache_misses INTEGER NOT NULL DEFAULT 0,
    total_size_bytes BIGINT NOT NULL DEFAULT 0,
    avg_response_time_ms NUMERIC DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_cache_stats_requests_non_negative CHECK (total_requests >= 0),
    CONSTRAINT chk_cache_stats_hits_non_negative CHECK (cache_hits >= 0),
    CONSTRAINT chk_cache_stats_misses_non_negative CHECK (cache_misses >= 0),
    CONSTRAINT chk_cache_stats_hits_not_exceed_total CHECK (cache_hits <= total_requests),
    CONSTRAINT chk_cache_stats_misses_not_exceed_total CHECK (cache_misses <= total_requests),
    CONSTRAINT chk_cache_stats_size_non_negative CHECK (total_size_bytes >= 0),
    CONSTRAINT chk_cache_stats_response_time_non_negative CHECK (avg_response_time_ms >= 0),
    
    -- Unique constraint para evitar duplicatas por dia
    UNIQUE(tenant_id, request_type, date)
);

-- ========================================
-- ÍNDICES PARA CACHE
-- ========================================

-- Índice principal por chave (unique)
CREATE UNIQUE INDEX idx_cache_entries_key ON cache_entries(key);

-- Índice para limpeza de entradas expiradas
CREATE INDEX idx_cache_entries_expires_at ON cache_entries(expires_at);

-- Índice por tenant e tipo para estatísticas
CREATE INDEX idx_cache_entries_tenant_type ON cache_entries(tenant_id, request_type);

-- Índice por número de processo
CREATE INDEX idx_cache_entries_process_number ON cache_entries(process_number)
WHERE process_number IS NOT NULL;

-- Índice por tribunal
CREATE INDEX idx_cache_entries_court_id ON cache_entries(court_id)
WHERE court_id IS NOT NULL;

-- Índice para entradas mais acessadas
CREATE INDEX idx_cache_entries_hit_count ON cache_entries(hit_count DESC, last_hit_at DESC);

-- Índice para limpeza por tamanho
CREATE INDEX idx_cache_entries_size_created ON cache_entries(size_bytes DESC, created_at ASC);

-- ========================================
-- ÍNDICES PARA EVENTOS
-- ========================================

-- Índice por tipo de evento
CREATE INDEX idx_domain_events_type ON domain_events(event_type);

-- Índice por aggregate
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_type, aggregate_id);

-- Índice para processamento de eventos
CREATE INDEX idx_domain_events_processing ON domain_events(is_processed, occurred_at)
WHERE is_processed = false;

-- Índice temporal
CREATE INDEX idx_domain_events_occurred_at ON domain_events(occurred_at);

-- Índice por tenant
CREATE INDEX idx_domain_events_tenant ON domain_events(tenant_id)
WHERE tenant_id IS NOT NULL;

-- Índice para correlação
CREATE INDEX idx_domain_events_correlation ON domain_events(correlation_id)
WHERE correlation_id IS NOT NULL;

-- ========================================
-- ÍNDICES PARA ESTATÍSTICAS
-- ========================================

-- Índice principal para cache statistics
CREATE INDEX idx_cache_statistics_tenant_date ON cache_statistics(tenant_id, date);
CREATE INDEX idx_cache_statistics_type_date ON cache_statistics(request_type, date);
CREATE INDEX idx_cache_statistics_date ON cache_statistics(date);

-- ========================================
-- TRIGGERS
-- ========================================

-- Trigger para atualizar updated_at no cache
CREATE TRIGGER update_cache_entries_updated_at
    BEFORE UPDATE ON cache_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para atualizar updated_at nas estatísticas
CREATE TRIGGER update_cache_statistics_updated_at
    BEFORE UPDATE ON cache_statistics
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para validação de cache
CREATE OR REPLACE FUNCTION validate_cache_entry()
RETURNS TRIGGER AS $$
BEGIN
    -- Calcular tamanho se não foi definido
    IF NEW.size_bytes = 0 THEN
        NEW.size_bytes := LENGTH(NEW.value::text);
    END IF;
    
    -- Definir expires_at baseado no TTL
    IF NEW.expires_at IS NULL OR NEW.expires_at <= NOW() THEN
        NEW.expires_at := NOW() + (NEW.ttl || ' seconds')::INTERVAL;
    END IF;
    
    -- Atualizar last_hit_at se hit_count aumentou
    IF OLD IS NOT NULL AND NEW.hit_count > OLD.hit_count THEN
        NEW.last_hit_at := NOW();
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_cache_entry_trigger
    BEFORE INSERT OR UPDATE ON cache_entries
    FOR EACH ROW
    EXECUTE FUNCTION validate_cache_entry();

-- Trigger para auto-atualizar estatísticas de cache
CREATE OR REPLACE FUNCTION update_cache_statistics()
RETURNS TRIGGER AS $$
DECLARE
    v_is_hit BOOLEAN;
BEGIN
    -- Determinar se foi hit ou miss baseado na mudança de hit_count
    v_is_hit := (OLD IS NOT NULL AND NEW.hit_count > OLD.hit_count);
    
    -- Atualizar estatísticas diárias
    INSERT INTO cache_statistics (
        tenant_id, request_type, date, 
        total_requests, cache_hits, cache_misses, total_size_bytes
    ) VALUES (
        NEW.tenant_id, NEW.request_type, CURRENT_DATE,
        1, 
        CASE WHEN v_is_hit THEN 1 ELSE 0 END,
        CASE WHEN v_is_hit THEN 0 ELSE 1 END,
        NEW.size_bytes
    )
    ON CONFLICT (tenant_id, request_type, date) DO UPDATE SET
        total_requests = cache_statistics.total_requests + 1,
        cache_hits = cache_statistics.cache_hits + CASE WHEN v_is_hit THEN 1 ELSE 0 END,
        cache_misses = cache_statistics.cache_misses + CASE WHEN v_is_hit THEN 0 ELSE 1 END,
        total_size_bytes = cache_statistics.total_size_bytes + NEW.size_bytes,
        updated_at = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_cache_statistics_trigger
    AFTER INSERT OR UPDATE ON cache_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_cache_statistics();

-- ========================================
-- FUNÇÕES UTILITÁRIAS PARA CACHE
-- ========================================

-- Função para buscar entrada no cache
CREATE OR REPLACE FUNCTION cache_get(p_key VARCHAR(255))
RETURNS TABLE(
    cache_value JSONB,
    hit_count INTEGER,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_expired BOOLEAN
) AS $$
DECLARE
    v_entry RECORD;
    v_is_expired BOOLEAN;
BEGIN
    -- Buscar entrada
    SELECT * INTO v_entry
    FROM cache_entries
    WHERE key = p_key;
    
    -- Se não encontrou, retornar vazio
    IF v_entry IS NULL THEN
        RETURN;
    END IF;
    
    -- Verificar se expirou
    v_is_expired := NOW() > v_entry.expires_at;
    
    -- Se expirou, deletar e retornar vazio
    IF v_is_expired THEN
        DELETE FROM cache_entries WHERE key = p_key;
        RETURN;
    END IF;
    
    -- Atualizar hit count
    UPDATE cache_entries
    SET hit_count = hit_count + 1,
        last_hit_at = NOW(),
        updated_at = NOW()
    WHERE key = p_key;
    
    -- Retornar dados
    RETURN QUERY
    SELECT 
        v_entry.value,
        v_entry.hit_count + 1,
        v_entry.expires_at,
        false;
END;
$$ LANGUAGE plpgsql;

-- Função para armazenar no cache
CREATE OR REPLACE FUNCTION cache_set(
    p_key VARCHAR(255),
    p_value JSONB,
    p_ttl INTEGER,
    p_tenant_id UUID,
    p_request_type VARCHAR(20),
    p_process_number VARCHAR(25) DEFAULT NULL,
    p_court_id VARCHAR(10) DEFAULT NULL
)
RETURNS BOOLEAN AS $$
DECLARE
    v_size_bytes INTEGER;
    v_expires_at TIMESTAMP WITH TIME ZONE;
BEGIN
    -- Calcular tamanho e expiração
    v_size_bytes := LENGTH(p_value::text);
    v_expires_at := NOW() + (p_ttl || ' seconds')::INTERVAL;
    
    -- Inserir ou atualizar
    INSERT INTO cache_entries (
        key, value, ttl, expires_at, size_bytes,
        tenant_id, request_type, process_number, court_id
    ) VALUES (
        p_key, p_value, p_ttl, v_expires_at, v_size_bytes,
        p_tenant_id, p_request_type, p_process_number, p_court_id
    )
    ON CONFLICT (key) DO UPDATE SET
        value = EXCLUDED.value,
        ttl = EXCLUDED.ttl,
        expires_at = EXCLUDED.expires_at,
        size_bytes = EXCLUDED.size_bytes,
        tenant_id = EXCLUDED.tenant_id,
        request_type = EXCLUDED.request_type,
        process_number = EXCLUDED.process_number,
        court_id = EXCLUDED.court_id,
        updated_at = NOW();
    
    RETURN true;
END;
$$ LANGUAGE plpgsql;

-- Função para deletar do cache
CREATE OR REPLACE FUNCTION cache_delete(p_key VARCHAR(255))
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    DELETE FROM cache_entries WHERE key = p_key;
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para limpeza de cache expirado
CREATE OR REPLACE FUNCTION cache_cleanup_expired()
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER := 0;
BEGIN
    DELETE FROM cache_entries
    WHERE expires_at <= NOW();
    
    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Função para limpeza por tamanho (LRU)
CREATE OR REPLACE FUNCTION cache_cleanup_by_size(p_max_size_mb INTEGER)
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER := 0;
    v_current_size_mb NUMERIC;
    v_target_size_mb NUMERIC;
BEGIN
    -- Calcular tamanho atual em MB
    SELECT COALESCE(SUM(size_bytes) / 1024.0 / 1024.0, 0)
    INTO v_current_size_mb
    FROM cache_entries;
    
    -- Se não excedeu o limite, não fazer nada
    IF v_current_size_mb <= p_max_size_mb THEN
        RETURN 0;
    END IF;
    
    -- Calcular tamanho alvo (80% do máximo)
    v_target_size_mb := p_max_size_mb * 0.8;
    
    -- Deletar entradas menos usadas até atingir o alvo
    WITH entries_to_delete AS (
        SELECT id
        FROM cache_entries
        ORDER BY 
            hit_count ASC,
            last_hit_at ASC NULLS FIRST,
            created_at ASC
    )
    DELETE FROM cache_entries
    WHERE id IN (
        SELECT id FROM entries_to_delete
        LIMIT (SELECT COUNT(*) * 0.2 FROM cache_entries)::INTEGER  -- Remove 20% das entradas
    );
    
    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Função para estatísticas de cache
CREATE OR REPLACE FUNCTION get_cache_stats(
    p_tenant_id UUID DEFAULT NULL,
    p_request_type VARCHAR(20) DEFAULT NULL,
    p_days INTEGER DEFAULT 7
)
RETURNS TABLE(
    total_entries INTEGER,
    total_size_mb NUMERIC,
    hit_ratio_percentage NUMERIC,
    avg_ttl_seconds NUMERIC,
    most_hit_entries JSONB,
    stats_by_type JSONB,
    daily_stats JSONB
) AS $$
DECLARE
    v_stats RECORD;
    v_most_hit JSONB;
    v_by_type JSONB;
    v_daily JSONB;
BEGIN
    -- Estatísticas gerais
    SELECT 
        COUNT(*)::INTEGER as total,
        ROUND(COALESCE(SUM(size_bytes) / 1024.0 / 1024.0, 0), 2) as size_mb,
        ROUND(AVG(ttl), 2) as avg_ttl
    INTO v_stats
    FROM cache_entries
    WHERE (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
      AND (p_request_type IS NULL OR request_type = p_request_type);
    
    -- Entradas mais acessadas
    SELECT jsonb_agg(
        jsonb_build_object(
            'key', key,
            'hit_count', hit_count,
            'request_type', request_type,
            'process_number', process_number
        )
    ) INTO v_most_hit
    FROM (
        SELECT key, hit_count, request_type, process_number
        FROM cache_entries
        WHERE (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
          AND (p_request_type IS NULL OR request_type = p_request_type)
        ORDER BY hit_count DESC
        LIMIT 10
    ) t;
    
    -- Stats por tipo
    SELECT jsonb_object_agg(request_type, stats) INTO v_by_type
    FROM (
        SELECT 
            request_type,
            jsonb_build_object(
                'count', COUNT(*),
                'size_mb', ROUND(SUM(size_bytes) / 1024.0 / 1024.0, 2),
                'avg_hit_count', ROUND(AVG(hit_count), 2)
            ) as stats
        FROM cache_entries
        WHERE (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
        GROUP BY request_type
    ) t;
    
    -- Stats diárias
    SELECT jsonb_agg(
        jsonb_build_object(
            'date', date,
            'total_requests', total_requests,
            'cache_hits', cache_hits,
            'cache_misses', cache_misses,
            'hit_ratio', CASE 
                WHEN total_requests > 0 THEN 
                    ROUND(cache_hits::NUMERIC / total_requests * 100, 2)
                ELSE 0 
            END
        )
    ) INTO v_daily
    FROM cache_statistics
    WHERE (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
      AND (p_request_type IS NULL OR request_type = p_request_type)
      AND date >= CURRENT_DATE - (p_days || ' days')::INTERVAL
    ORDER BY date;
    
    -- Calcular hit ratio geral das estatísticas
    DECLARE
        v_total_requests BIGINT;
        v_total_hits BIGINT;
        v_hit_ratio NUMERIC;
    BEGIN
        SELECT 
            COALESCE(SUM(total_requests), 0),
            COALESCE(SUM(cache_hits), 0)
        INTO v_total_requests, v_total_hits
        FROM cache_statistics
        WHERE (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
          AND (p_request_type IS NULL OR request_type = p_request_type)
          AND date >= CURRENT_DATE - (p_days || ' days')::INTERVAL;
        
        v_hit_ratio := CASE 
            WHEN v_total_requests > 0 THEN 
                ROUND(v_total_hits::NUMERIC / v_total_requests * 100, 2)
            ELSE 0 
        END;
    END;
    
    RETURN QUERY
    SELECT 
        v_stats.total,
        v_stats.size_mb,
        v_hit_ratio,
        v_stats.avg_ttl,
        COALESCE(v_most_hit, '[]'::jsonb),
        COALESCE(v_by_type, '{}'::jsonb),
        COALESCE(v_daily, '[]'::jsonb);
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- FUNÇÕES UTILITÁRIAS PARA EVENTOS
-- ========================================

-- Função para salvar evento de domínio
CREATE OR REPLACE FUNCTION save_domain_event(
    p_event_type VARCHAR(100),
    p_aggregate_id UUID,
    p_aggregate_type VARCHAR(50),
    p_event_data JSONB,
    p_occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    p_correlation_id UUID DEFAULT NULL,
    p_causation_id UUID DEFAULT NULL,
    p_tenant_id UUID DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_event_id UUID;
BEGIN
    INSERT INTO domain_events (
        event_type, aggregate_id, aggregate_type, event_data,
        occurred_at, correlation_id, causation_id, tenant_id
    ) VALUES (
        p_event_type, p_aggregate_id, p_aggregate_type, p_event_data,
        p_occurred_at, p_correlation_id, p_causation_id, p_tenant_id
    )
    RETURNING id INTO v_event_id;
    
    RETURN v_event_id;
END;
$$ LANGUAGE plpgsql;

-- Função para marcar evento como processado
CREATE OR REPLACE FUNCTION mark_event_processed(p_event_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE domain_events
    SET is_processed = true,
        processed_at = NOW()
    WHERE id = p_event_id AND is_processed = false;
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para obter eventos não processados
CREATE OR REPLACE FUNCTION get_unprocessed_events(p_limit INTEGER DEFAULT 100)
RETURNS TABLE(
    event_id UUID,
    event_type VARCHAR(100),
    aggregate_id UUID,
    aggregate_type VARCHAR(50),
    event_data JSONB,
    occurred_at TIMESTAMP WITH TIME ZONE,
    correlation_id UUID,
    tenant_id UUID
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        de.id,
        de.event_type,
        de.aggregate_id,
        de.aggregate_type,
        de.event_data,
        de.occurred_at,
        de.correlation_id,
        de.tenant_id
    FROM domain_events de
    WHERE de.is_processed = false
    ORDER BY de.occurred_at ASC
    LIMIT p_limit
    FOR UPDATE SKIP LOCKED;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- VIEWS ÚTEIS
-- ========================================

-- View para cache summary
CREATE OR REPLACE VIEW v_cache_summary AS
SELECT 
    tenant_id,
    request_type,
    COUNT(*) as total_entries,
    ROUND(SUM(size_bytes) / 1024.0 / 1024.0, 2) as total_size_mb,
    ROUND(AVG(hit_count), 2) as avg_hit_count,
    ROUND(AVG(ttl), 2) as avg_ttl_seconds,
    MAX(hit_count) as max_hit_count,
    COUNT(CASE WHEN expires_at > NOW() THEN 1 END) as active_entries,
    COUNT(CASE WHEN expires_at <= NOW() THEN 1 END) as expired_entries
FROM cache_entries
GROUP BY tenant_id, request_type;

-- View para eventos recentes
CREATE OR REPLACE VIEW v_recent_events AS
SELECT 
    event_type,
    aggregate_type,
    COUNT(*) as event_count,
    MAX(occurred_at) as last_occurred,
    COUNT(CASE WHEN is_processed THEN 1 END) as processed_count,
    COUNT(CASE WHEN NOT is_processed THEN 1 END) as pending_count
FROM domain_events
WHERE occurred_at >= NOW() - INTERVAL '24 hours'
GROUP BY event_type, aggregate_type
ORDER BY event_count DESC;

-- ========================================
-- COMENTÁRIOS
-- ========================================

COMMENT ON TABLE cache_entries IS 'Entradas de cache para consultas DataJud';
COMMENT ON TABLE domain_events IS 'Eventos de domínio para Event Sourcing';
COMMENT ON TABLE cache_statistics IS 'Estatísticas diárias de uso do cache';

-- ========================================
-- FINALIZAÇÃO
-- ========================================

-- Criar job de limpeza automática (se suportado pelo PostgreSQL)
-- Em produção, seria configurado via cron ou job scheduler externo

-- Atualizar estatísticas
ANALYZE cache_entries;
ANALYZE domain_events;
ANALYZE cache_statistics;

-- Log de conclusão
DO $$
BEGIN
    RAISE NOTICE 'Migração 005 concluída - sistema de cache e eventos criado';
    RAISE NOTICE 'Funções utilitárias e views configuradas';
END $$;