-- Migration: Create rate limiters table
-- Description: Tabela para controle de rate limiting por CNPJ, tenant e global
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- TABELA DE RATE LIMITERS
-- ========================================

CREATE TABLE rate_limiters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL,
    key VARCHAR(255) NOT NULL,
    window_size_seconds INTEGER NOT NULL,
    max_requests INTEGER NOT NULL,
    current_requests INTEGER NOT NULL DEFAULT 0,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT chk_rate_limiter_type CHECK (type IN ('cnpj', 'tenant', 'global')),
    CONSTRAINT chk_window_size_positive CHECK (window_size_seconds > 0),
    CONSTRAINT chk_max_requests_positive CHECK (max_requests > 0),
    CONSTRAINT chk_current_requests_non_negative CHECK (current_requests >= 0),
    CONSTRAINT chk_current_not_exceed_max CHECK (current_requests <= max_requests),
    
    -- Unique constraint para type + key
    UNIQUE(type, key)
);

-- ========================================
-- TABELA DE REQUISIÇÕES DO RATE LIMITER
-- ========================================

CREATE TABLE rate_limiter_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rate_limiter_id UUID NOT NULL REFERENCES rate_limiters(id) ON DELETE CASCADE,
    request_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Índice para limpeza eficiente
    CONSTRAINT chk_request_timestamp_not_future CHECK (request_timestamp <= NOW())
);

-- ========================================
-- ÍNDICES
-- ========================================

-- Índice principal por type e key
CREATE UNIQUE INDEX idx_rate_limiters_type_key ON rate_limiters(type, key);

-- Índice para CNPJs ativos
CREATE INDEX idx_rate_limiters_cnpj_active ON rate_limiters(key, is_active)
WHERE type = 'cnpj' AND is_active = true;

-- Índice para tenants ativos
CREATE INDEX idx_rate_limiters_tenant_active ON rate_limiters(key, is_active)
WHERE type = 'tenant' AND is_active = true;

-- Índice para limpeza por janela
CREATE INDEX idx_rate_limiters_window_cleanup ON rate_limiters(window_start, window_size_seconds)
WHERE is_active = true;

-- Índices para rate_limiter_requests
CREATE INDEX idx_rate_limiter_requests_limiter_id ON rate_limiter_requests(rate_limiter_id);
CREATE INDEX idx_rate_limiter_requests_timestamp ON rate_limiter_requests(request_timestamp);
CREATE INDEX idx_rate_limiter_requests_cleanup ON rate_limiter_requests(rate_limiter_id, request_timestamp);

-- ========================================
-- TRIGGERS
-- ========================================

-- Trigger para atualizar updated_at
CREATE TRIGGER update_rate_limiters_updated_at
    BEFORE UPDATE ON rate_limiters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para limpeza automática de janela
CREATE OR REPLACE FUNCTION cleanup_rate_limiter_window()
RETURNS TRIGGER AS $$
DECLARE
    window_end TIMESTAMP WITH TIME ZONE;
BEGIN
    -- Calcular fim da janela atual
    window_end := NEW.window_start + (NEW.window_size_seconds || ' seconds')::INTERVAL;
    
    -- Se a janela expirou, resetar
    IF NOW() > window_end THEN
        NEW.current_requests := 0;
        NEW.window_start := NOW();
        
        -- Limpar requisições antigas desta janela
        DELETE FROM rate_limiter_requests
        WHERE rate_limiter_id = NEW.id
          AND request_timestamp < NEW.window_start;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cleanup_rate_limiter_window_trigger
    BEFORE UPDATE ON rate_limiters
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_rate_limiter_window();

-- ========================================
-- FUNÇÕES UTILITÁRIAS
-- ========================================

-- Função para obter ou criar rate limiter
CREATE OR REPLACE FUNCTION get_or_create_rate_limiter(
    p_type VARCHAR(20),
    p_key VARCHAR(255),
    p_window_size_seconds INTEGER DEFAULT 3600,
    p_max_requests INTEGER DEFAULT 1000
)
RETURNS UUID AS $$
DECLARE
    v_limiter_id UUID;
BEGIN
    -- Tentar encontrar existente
    SELECT id INTO v_limiter_id
    FROM rate_limiters
    WHERE type = p_type AND key = p_key;
    
    -- Se não encontrou, criar
    IF v_limiter_id IS NULL THEN
        INSERT INTO rate_limiters (type, key, window_size_seconds, max_requests)
        VALUES (p_type, p_key, p_window_size_seconds, p_max_requests)
        RETURNING id INTO v_limiter_id;
    END IF;
    
    RETURN v_limiter_id;
END;
$$ LANGUAGE plpgsql;

-- Função para verificar se requisição pode ser feita
CREATE OR REPLACE FUNCTION check_rate_limit(
    p_type VARCHAR(20),
    p_key VARCHAR(255),
    p_window_size_seconds INTEGER DEFAULT 3600,
    p_max_requests INTEGER DEFAULT 1000
)
RETURNS TABLE(
    allowed BOOLEAN,
    requests_used INTEGER,
    requests_limit INTEGER,
    reset_time TIMESTAMP WITH TIME ZONE,
    retry_after_seconds INTEGER
) AS $$
DECLARE
    v_limiter_id UUID;
    v_limiter RECORD;
    v_window_end TIMESTAMP WITH TIME ZONE;
    v_requests_in_window INTEGER;
BEGIN
    -- Obter ou criar rate limiter
    v_limiter_id := get_or_create_rate_limiter(p_type, p_key, p_window_size_seconds, p_max_requests);
    
    -- Buscar dados atuais do limiter
    SELECT * INTO v_limiter
    FROM rate_limiters
    WHERE id = v_limiter_id;
    
    -- Calcular fim da janela
    v_window_end := v_limiter.window_start + (v_limiter.window_size_seconds || ' seconds')::INTERVAL;
    
    -- Se a janela expirou, resetar
    IF NOW() > v_window_end THEN
        UPDATE rate_limiters
        SET current_requests = 0,
            window_start = NOW(),
            updated_at = NOW()
        WHERE id = v_limiter_id;
        
        v_limiter.current_requests := 0;
        v_limiter.window_start := NOW();
        v_window_end := NOW() + (v_limiter.window_size_seconds || ' seconds')::INTERVAL;
        
        -- Limpar requisições antigas
        DELETE FROM rate_limiter_requests
        WHERE rate_limiter_id = v_limiter_id
          AND request_timestamp < v_limiter.window_start;
    END IF;
    
    -- Contar requisições na janela atual
    SELECT COUNT(*) INTO v_requests_in_window
    FROM rate_limiter_requests
    WHERE rate_limiter_id = v_limiter_id
      AND request_timestamp >= v_limiter.window_start;
    
    -- Se contagem divergir, atualizar
    IF v_requests_in_window != v_limiter.current_requests THEN
        UPDATE rate_limiters
        SET current_requests = v_requests_in_window,
            updated_at = NOW()
        WHERE id = v_limiter_id;
        
        v_limiter.current_requests := v_requests_in_window;
    END IF;
    
    RETURN QUERY
    SELECT 
        (v_limiter.current_requests < v_limiter.max_requests AND v_limiter.is_active) as allowed,
        v_limiter.current_requests as requests_used,
        v_limiter.max_requests as requests_limit,
        v_window_end as reset_time,
        GREATEST(0, EXTRACT(EPOCH FROM (v_window_end - NOW())))::INTEGER as retry_after_seconds;
END;
$$ LANGUAGE plpgsql;

-- Função para consumir uma requisição do rate limiter
CREATE OR REPLACE FUNCTION consume_rate_limit(
    p_type VARCHAR(20),
    p_key VARCHAR(255),
    p_window_size_seconds INTEGER DEFAULT 3600,
    p_max_requests INTEGER DEFAULT 1000
)
RETURNS TABLE(
    success BOOLEAN,
    requests_used INTEGER,
    requests_limit INTEGER,
    reset_time TIMESTAMP WITH TIME ZONE
) AS $$
DECLARE
    v_limiter_id UUID;
    v_allowed BOOLEAN;
    v_requests_used INTEGER;
    v_requests_limit INTEGER;
    v_reset_time TIMESTAMP WITH TIME ZONE;
    v_retry_after INTEGER;
BEGIN
    -- Verificar se pode fazer a requisição
    SELECT * INTO v_allowed, v_requests_used, v_requests_limit, v_reset_time, v_retry_after
    FROM check_rate_limit(p_type, p_key, p_window_size_seconds, p_max_requests);
    
    IF v_allowed THEN
        -- Obter ID do rate limiter
        v_limiter_id := get_or_create_rate_limiter(p_type, p_key, p_window_size_seconds, p_max_requests);
        
        -- Registrar a requisição
        INSERT INTO rate_limiter_requests (rate_limiter_id, request_timestamp)
        VALUES (v_limiter_id, NOW());
        
        -- Atualizar contador
        UPDATE rate_limiters
        SET current_requests = current_requests + 1,
            last_used_at = NOW(),
            updated_at = NOW()
        WHERE id = v_limiter_id;
        
        v_requests_used := v_requests_used + 1;
    END IF;
    
    RETURN QUERY
    SELECT 
        v_allowed as success,
        v_requests_used as requests_used,
        v_requests_limit as requests_limit,
        v_reset_time as reset_time;
END;
$$ LANGUAGE plpgsql;

-- Função para resetar rate limiter
CREATE OR REPLACE FUNCTION reset_rate_limiter(
    p_type VARCHAR(20),
    p_key VARCHAR(255)
)
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE rate_limiters
    SET current_requests = 0,
        window_start = NOW(),
        updated_at = NOW()
    WHERE type = p_type AND key = p_key;
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    IF v_rows_affected > 0 THEN
        -- Limpar requisições antigas
        DELETE FROM rate_limiter_requests
        WHERE rate_limiter_id IN (
            SELECT id FROM rate_limiters 
            WHERE type = p_type AND key = p_key
        );
    END IF;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para obter estatísticas de rate limiting
CREATE OR REPLACE FUNCTION get_rate_limiter_stats(p_type VARCHAR(20) DEFAULT NULL)
RETURNS TABLE(
    limiter_type VARCHAR(20),
    limiter_key VARCHAR(255),
    requests_used INTEGER,
    requests_limit INTEGER,
    usage_percentage NUMERIC,
    window_start TIMESTAMP WITH TIME ZONE,
    window_end TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN,
    last_used_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        rl.type,
        rl.key,
        rl.current_requests,
        rl.max_requests,
        CASE 
            WHEN rl.max_requests > 0 THEN 
                ROUND(rl.current_requests::NUMERIC / rl.max_requests * 100, 2)
            ELSE 0
        END as usage_percentage,
        rl.window_start,
        rl.window_start + (rl.window_size_seconds || ' seconds')::INTERVAL as window_end,
        rl.is_active,
        rl.last_used_at
    FROM rate_limiters rl
    WHERE (p_type IS NULL OR rl.type = p_type)
    ORDER BY 
        rl.type,
        CASE 
            WHEN rl.max_requests > 0 THEN rl.current_requests::NUMERIC / rl.max_requests
            ELSE 0
        END DESC;
END;
$$ LANGUAGE plpgsql;

-- Função para limpeza de dados antigos
CREATE OR REPLACE FUNCTION cleanup_rate_limiter_data(p_days INTEGER DEFAULT 7)
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER := 0;
    v_cutoff_date TIMESTAMP WITH TIME ZONE;
BEGIN
    v_cutoff_date := NOW() - (p_days || ' days')::INTERVAL;
    
    -- Limpar requisições antigas
    DELETE FROM rate_limiter_requests
    WHERE request_timestamp < v_cutoff_date;
    
    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    
    -- Resetar limiters que não foram usados recentemente
    UPDATE rate_limiters
    SET current_requests = 0,
        window_start = NOW(),
        updated_at = NOW()
    WHERE last_used_at < v_cutoff_date
       OR (last_used_at IS NULL AND created_at < v_cutoff_date);
    
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Função para criar rate limiters padrão
CREATE OR REPLACE FUNCTION create_default_rate_limiters()
RETURNS INTEGER AS $$
DECLARE
    v_created_count INTEGER := 0;
BEGIN
    -- Rate limiter global (por hora)
    INSERT INTO rate_limiters (type, key, window_size_seconds, max_requests)
    VALUES ('global', 'global', 3600, 100000)
    ON CONFLICT (type, key) DO NOTHING;
    
    IF FOUND THEN
        v_created_count := v_created_count + 1;
    END IF;
    
    -- Rate limiter global diário
    INSERT INTO rate_limiters (type, key, window_size_seconds, max_requests)
    VALUES ('global', 'global_daily', 86400, 1000000)
    ON CONFLICT (type, key) DO NOTHING;
    
    IF FOUND THEN
        v_created_count := v_created_count + 1;
    END IF;
    
    RETURN v_created_count;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- VIEWS ÚTEIS
-- ========================================

-- View para rate limiters CNPJ com informações de quota
CREATE OR REPLACE VIEW v_cnpj_rate_limiters AS
SELECT 
    rl.id,
    rl.key as cnpj,
    rl.current_requests as daily_usage,
    rl.max_requests as daily_limit,
    (rl.max_requests - rl.current_requests) as available_quota,
    CASE 
        WHEN rl.max_requests > 0 THEN 
            ROUND(rl.current_requests::NUMERIC / rl.max_requests * 100, 2)
        ELSE 0
    END as usage_percentage,
    rl.window_start + (rl.window_size_seconds || ' seconds')::INTERVAL as reset_time,
    rl.is_active,
    rl.last_used_at,
    rl.created_at,
    rl.updated_at
FROM rate_limiters rl
WHERE rl.type = 'cnpj';

-- View para estatísticas gerais de rate limiting
CREATE OR REPLACE VIEW v_rate_limiting_summary AS
SELECT 
    type,
    COUNT(*) as total_limiters,
    COUNT(CASE WHEN is_active THEN 1 END) as active_limiters,
    SUM(max_requests) as total_limit,
    SUM(current_requests) as total_usage,
    SUM(max_requests - current_requests) as total_available,
    CASE 
        WHEN SUM(max_requests) > 0 THEN 
            ROUND(SUM(current_requests)::NUMERIC / SUM(max_requests) * 100, 2)
        ELSE 0
    END as overall_usage_percentage,
    COUNT(CASE WHEN current_requests >= max_requests THEN 1 END) as exhausted_limiters
FROM rate_limiters
WHERE is_active = true
GROUP BY type;

-- ========================================
-- COMENTÁRIOS
-- ========================================

COMMENT ON TABLE rate_limiters IS 'Controladores de rate limiting por CNPJ, tenant e global';
COMMENT ON TABLE rate_limiter_requests IS 'Registro de requisições para controle de janela deslizante';

COMMENT ON COLUMN rate_limiters.type IS 'Tipo de rate limiter: cnpj, tenant, global';
COMMENT ON COLUMN rate_limiters.key IS 'Chave identificadora (CNPJ, tenant_id, etc)';
COMMENT ON COLUMN rate_limiters.window_size_seconds IS 'Tamanho da janela em segundos';
COMMENT ON COLUMN rate_limiters.max_requests IS 'Máximo de requisições na janela';
COMMENT ON COLUMN rate_limiters.current_requests IS 'Requisições atuais na janela';
COMMENT ON COLUMN rate_limiters.window_start IS 'Início da janela atual';

-- ========================================
-- DADOS INICIAIS
-- ========================================

-- Criar rate limiters padrão
SELECT create_default_rate_limiters();

-- ========================================
-- FINALIZAÇÃO
-- ========================================

-- Atualizar estatísticas
ANALYZE rate_limiters;
ANALYZE rate_limiter_requests;

-- Log de conclusão
DO $$
BEGIN
    RAISE NOTICE 'Migração 003 concluída - sistema de rate limiting criado';
    RAISE NOTICE 'Rate limiters padrão: %', (SELECT COUNT(*) FROM rate_limiters);
END $$;