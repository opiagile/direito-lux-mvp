-- Migration: Create circuit breakers table
-- Description: Tabela para controle de circuit breakers para proteção contra falhas
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- TABELA DE CIRCUIT BREAKERS
-- ========================================

CREATE TABLE circuit_breakers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    state VARCHAR(20) NOT NULL DEFAULT 'closed',
    failure_threshold INTEGER NOT NULL DEFAULT 5,
    success_threshold INTEGER NOT NULL DEFAULT 3,
    timeout_seconds INTEGER NOT NULL DEFAULT 30,
    max_requests INTEGER NOT NULL DEFAULT 5,
    failure_count INTEGER NOT NULL DEFAULT 0,
    success_count INTEGER NOT NULL DEFAULT 0,
    request_count INTEGER NOT NULL DEFAULT 0,
    last_failure_time TIMESTAMP WITH TIME ZONE,
    last_success_time TIMESTAMP WITH TIME ZONE,
    state_changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_circuit_breaker_state CHECK (state IN ('closed', 'open', 'half_open')),
    CONSTRAINT chk_failure_threshold_positive CHECK (failure_threshold > 0),
    CONSTRAINT chk_success_threshold_positive CHECK (success_threshold > 0),
    CONSTRAINT chk_timeout_positive CHECK (timeout_seconds > 0),
    CONSTRAINT chk_max_requests_positive CHECK (max_requests > 0),
    CONSTRAINT chk_failure_count_non_negative CHECK (failure_count >= 0),
    CONSTRAINT chk_success_count_non_negative CHECK (success_count >= 0),
    CONSTRAINT chk_request_count_non_negative CHECK (request_count >= 0),
    CONSTRAINT chk_failure_count_not_exceed_threshold CHECK (failure_count <= failure_threshold + 10), -- Margem para burst
    CONSTRAINT chk_request_count_reasonable CHECK (request_count <= max_requests + 10) -- Margem para concorrência
);

-- ========================================
-- TABELA DE EXECUÇÕES DO CIRCUIT BREAKER
-- ========================================

CREATE TABLE circuit_breaker_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    circuit_breaker_id UUID NOT NULL REFERENCES circuit_breakers(id) ON DELETE CASCADE,
    success BOOLEAN NOT NULL,
    duration_ms INTEGER NOT NULL,
    error_message TEXT,
    error_code VARCHAR(50),
    state_before VARCHAR(20) NOT NULL,
    state_after VARCHAR(20) NOT NULL,
    executed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_duration_non_negative CHECK (duration_ms >= 0),
    CONSTRAINT chk_state_before_valid CHECK (state_before IN ('closed', 'open', 'half_open')),
    CONSTRAINT chk_state_after_valid CHECK (state_after IN ('closed', 'open', 'half_open'))
);

-- ========================================
-- ÍNDICES
-- ========================================

-- Índice principal por nome
CREATE UNIQUE INDEX idx_circuit_breakers_name ON circuit_breakers(name);

-- Índice por estado para consultas de saúde
CREATE INDEX idx_circuit_breakers_state ON circuit_breakers(state, is_active);

-- Índice para circuit breakers que podem tentar reset
CREATE INDEX idx_circuit_breakers_reset_check ON circuit_breakers(state, state_changed_at, timeout_seconds)
WHERE state = 'open' AND is_active = true;

-- Índice para circuit breakers ativos
CREATE INDEX idx_circuit_breakers_active ON circuit_breakers(is_active, name)
WHERE is_active = true;

-- Índices para executions
CREATE INDEX idx_circuit_breaker_executions_cb_id ON circuit_breaker_executions(circuit_breaker_id);
CREATE INDEX idx_circuit_breaker_executions_executed_at ON circuit_breaker_executions(executed_at);
CREATE INDEX idx_circuit_breaker_executions_success ON circuit_breaker_executions(circuit_breaker_id, success, executed_at);

-- ========================================
-- TRIGGERS
-- ========================================

-- Trigger para atualizar updated_at
CREATE TRIGGER update_circuit_breakers_updated_at
    BEFORE UPDATE ON circuit_breakers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para validação e auto-reset
CREATE OR REPLACE FUNCTION validate_circuit_breaker()
RETURNS TRIGGER AS $$
BEGIN
    -- Auto-reset para half-open se passou do timeout
    IF NEW.state = 'open' AND 
       (NOW() - NEW.state_changed_at) >= (NEW.timeout_seconds || ' seconds')::INTERVAL THEN
        NEW.state := 'half_open';
        NEW.state_changed_at := NOW();
        NEW.success_count := 0;
        NEW.failure_count := 0;
        NEW.request_count := 0;
    END IF;
    
    -- Reset contadores quando muda de estado
    IF OLD.state IS NOT NULL AND NEW.state != OLD.state THEN
        NEW.state_changed_at := NOW();
        IF NEW.state = 'closed' THEN
            NEW.failure_count := 0;
            NEW.success_count := 0;
            NEW.request_count := 0;
        ELSIF NEW.state = 'half_open' THEN
            NEW.success_count := 0;
            NEW.failure_count := 0;
            NEW.request_count := 0;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_circuit_breaker_trigger
    BEFORE UPDATE ON circuit_breakers
    FOR EACH ROW
    EXECUTE FUNCTION validate_circuit_breaker();

-- ========================================
-- FUNÇÕES UTILITÁRIAS
-- ========================================

-- Função para obter ou criar circuit breaker
CREATE OR REPLACE FUNCTION get_or_create_circuit_breaker(
    p_name VARCHAR(100),
    p_failure_threshold INTEGER DEFAULT 5,
    p_success_threshold INTEGER DEFAULT 3,
    p_timeout_seconds INTEGER DEFAULT 30,
    p_max_requests INTEGER DEFAULT 5
)
RETURNS UUID AS $$
DECLARE
    v_breaker_id UUID;
BEGIN
    -- Tentar encontrar existente
    SELECT id INTO v_breaker_id
    FROM circuit_breakers
    WHERE name = p_name;
    
    -- Se não encontrou, criar
    IF v_breaker_id IS NULL THEN
        INSERT INTO circuit_breakers (
            name, failure_threshold, success_threshold, 
            timeout_seconds, max_requests
        )
        VALUES (
            p_name, p_failure_threshold, p_success_threshold,
            p_timeout_seconds, p_max_requests
        )
        RETURNING id INTO v_breaker_id;
    END IF;
    
    RETURN v_breaker_id;
END;
$$ LANGUAGE plpgsql;

-- Função para verificar se circuit breaker pode executar
CREATE OR REPLACE FUNCTION can_execute_circuit_breaker(p_name VARCHAR(100))
RETURNS TABLE(
    can_execute BOOLEAN,
    current_state VARCHAR(20),
    reason TEXT
) AS $$
DECLARE
    v_breaker RECORD;
    v_timeout_passed BOOLEAN;
BEGIN
    -- Buscar circuit breaker
    SELECT * INTO v_breaker
    FROM circuit_breakers
    WHERE name = p_name AND is_active = true;
    
    -- Se não encontrou, pode executar (será criado)
    IF v_breaker IS NULL THEN
        RETURN QUERY SELECT true, 'closed'::VARCHAR(20), 'Circuit breaker will be created'::TEXT;
        RETURN;
    END IF;
    
    -- Verificar se timeout passou para estado open
    v_timeout_passed := (NOW() - v_breaker.state_changed_at) >= (v_breaker.timeout_seconds || ' seconds')::INTERVAL;
    
    -- Lógica baseada no estado
    CASE v_breaker.state
        WHEN 'closed' THEN
            RETURN QUERY SELECT true, v_breaker.state, 'Circuit breaker is closed'::TEXT;
            
        WHEN 'open' THEN
            IF v_timeout_passed THEN
                -- Pode tentar (vai para half-open)
                RETURN QUERY SELECT true, 'half_open'::VARCHAR(20), 'Timeout passed, transitioning to half-open'::TEXT;
            ELSE
                RETURN QUERY SELECT false, v_breaker.state, 'Circuit breaker is open'::TEXT;
            END IF;
            
        WHEN 'half_open' THEN
            IF v_breaker.request_count < v_breaker.max_requests THEN
                RETURN QUERY SELECT true, v_breaker.state, 'Half-open with available requests'::TEXT;
            ELSE
                RETURN QUERY SELECT false, v_breaker.state, 'Half-open request limit reached'::TEXT;
            END IF;
            
        ELSE
            RETURN QUERY SELECT false, v_breaker.state, 'Unknown state'::TEXT;
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Função para executar com circuit breaker
CREATE OR REPLACE FUNCTION execute_with_circuit_breaker(
    p_name VARCHAR(100),
    p_success BOOLEAN,
    p_duration_ms INTEGER,
    p_error_message TEXT DEFAULT NULL,
    p_error_code VARCHAR(50) DEFAULT NULL
)
RETURNS TABLE(
    execution_allowed BOOLEAN,
    new_state VARCHAR(20),
    state_changed BOOLEAN
) AS $$
DECLARE
    v_breaker_id UUID;
    v_breaker RECORD;
    v_old_state VARCHAR(20);
    v_new_state VARCHAR(20);
    v_state_changed BOOLEAN := false;
    v_can_execute BOOLEAN;
    v_reason TEXT;
BEGIN
    -- Verificar se pode executar
    SELECT * INTO v_can_execute, v_new_state, v_reason
    FROM can_execute_circuit_breaker(p_name);
    
    IF NOT v_can_execute THEN
        RETURN QUERY SELECT false, v_new_state, false;
        RETURN;
    END IF;
    
    -- Obter ou criar circuit breaker
    v_breaker_id := get_or_create_circuit_breaker(p_name);
    
    -- Buscar estado atual
    SELECT * INTO v_breaker
    FROM circuit_breakers
    WHERE id = v_breaker_id;
    
    v_old_state := v_breaker.state;
    v_new_state := v_breaker.state;
    
    -- Atualizar estado baseado no resultado
    IF p_success THEN
        -- Sucesso
        UPDATE circuit_breakers
        SET success_count = success_count + 1,
            request_count = request_count + 1,
            last_success_time = NOW(),
            updated_at = NOW()
        WHERE id = v_breaker_id;
        
        -- Verificar transições de estado para sucesso
        CASE v_breaker.state
            WHEN 'half_open' THEN
                IF v_breaker.success_count + 1 >= v_breaker.success_threshold THEN
                    v_new_state := 'closed';
                    v_state_changed := true;
                END IF;
        END CASE;
        
    ELSE
        -- Falha
        UPDATE circuit_breakers
        SET failure_count = failure_count + 1,
            request_count = request_count + 1,
            last_failure_time = NOW(),
            updated_at = NOW()
        WHERE id = v_breaker_id;
        
        -- Verificar transições de estado para falha
        CASE v_breaker.state
            WHEN 'closed' THEN
                IF v_breaker.failure_count + 1 >= v_breaker.failure_threshold THEN
                    v_new_state := 'open';
                    v_state_changed := true;
                END IF;
                
            WHEN 'half_open' THEN
                v_new_state := 'open';
                v_state_changed := true;
        END CASE;
    END IF;
    
    -- Atualizar estado se mudou
    IF v_state_changed THEN
        UPDATE circuit_breakers
        SET state = v_new_state,
            state_changed_at = NOW(),
            updated_at = NOW()
        WHERE id = v_breaker_id;
    END IF;
    
    -- Registrar execução
    INSERT INTO circuit_breaker_executions (
        circuit_breaker_id, success, duration_ms, error_message, error_code,
        state_before, state_after
    ) VALUES (
        v_breaker_id, p_success, p_duration_ms, p_error_message, p_error_code,
        v_old_state, v_new_state
    );
    
    RETURN QUERY SELECT true, v_new_state, v_state_changed;
END;
$$ LANGUAGE plpgsql;

-- Função para forçar reset de circuit breaker
CREATE OR REPLACE FUNCTION reset_circuit_breaker(p_name VARCHAR(100))
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE circuit_breakers
    SET state = 'closed',
        failure_count = 0,
        success_count = 0,
        request_count = 0,
        state_changed_at = NOW(),
        updated_at = NOW()
    WHERE name = p_name AND is_active = true;
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para forçar abertura de circuit breaker
CREATE OR REPLACE FUNCTION force_open_circuit_breaker(p_name VARCHAR(100))
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE circuit_breakers
    SET state = 'open',
        state_changed_at = NOW(),
        updated_at = NOW()
    WHERE name = p_name AND is_active = true;
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para obter estatísticas de circuit breaker
CREATE OR REPLACE FUNCTION get_circuit_breaker_stats(p_name VARCHAR(100) DEFAULT NULL)
RETURNS TABLE(
    breaker_name VARCHAR(100),
    current_state VARCHAR(20),
    failure_count INTEGER,
    success_count INTEGER,
    failure_threshold INTEGER,
    success_threshold INTEGER,
    failure_rate NUMERIC,
    uptime_seconds BIGINT,
    time_in_current_state_seconds BIGINT,
    last_failure_time TIMESTAMP WITH TIME ZONE,
    last_success_time TIMESTAMP WITH TIME ZONE,
    total_executions_24h BIGINT,
    success_executions_24h BIGINT,
    failure_executions_24h BIGINT,
    avg_response_time_24h NUMERIC,
    is_healthy BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        cb.name,
        cb.state,
        cb.failure_count,
        cb.success_count,
        cb.failure_threshold,
        cb.success_threshold,
        CASE 
            WHEN (cb.failure_count + cb.success_count) > 0 THEN
                ROUND(cb.failure_count::NUMERIC / (cb.failure_count + cb.success_count) * 100, 2)
            ELSE 0
        END as failure_rate,
        EXTRACT(EPOCH FROM (NOW() - cb.created_at))::BIGINT as uptime_seconds,
        EXTRACT(EPOCH FROM (NOW() - cb.state_changed_at))::BIGINT as time_in_current_state_seconds,
        cb.last_failure_time,
        cb.last_success_time,
        COALESCE(exec_stats.total_executions, 0) as total_executions_24h,
        COALESCE(exec_stats.success_executions, 0) as success_executions_24h,
        COALESCE(exec_stats.failure_executions, 0) as failure_executions_24h,
        COALESCE(exec_stats.avg_response_time, 0) as avg_response_time_24h,
        (cb.state = 'closed' OR (cb.state = 'half_open' AND cb.success_count > 0)) as is_healthy
    FROM circuit_breakers cb
    LEFT JOIN (
        SELECT 
            circuit_breaker_id,
            COUNT(*) as total_executions,
            COUNT(CASE WHEN success THEN 1 END) as success_executions,
            COUNT(CASE WHEN NOT success THEN 1 END) as failure_executions,
            ROUND(AVG(duration_ms), 2) as avg_response_time
        FROM circuit_breaker_executions
        WHERE executed_at >= NOW() - INTERVAL '24 hours'
        GROUP BY circuit_breaker_id
    ) exec_stats ON cb.id = exec_stats.circuit_breaker_id
    WHERE (p_name IS NULL OR cb.name = p_name)
      AND cb.is_active = true
    ORDER BY cb.name;
END;
$$ LANGUAGE plpgsql;

-- Função para verificação de saúde geral
CREATE OR REPLACE FUNCTION get_circuit_breakers_health()
RETURNS TABLE(
    total_breakers INTEGER,
    healthy_breakers INTEGER,
    unhealthy_breakers INTEGER,
    open_breakers INTEGER,
    half_open_breakers INTEGER,
    closed_breakers INTEGER,
    overall_health_percentage NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*)::INTEGER as total_breakers,
        COUNT(CASE WHEN state = 'closed' OR (state = 'half_open' AND success_count > 0) THEN 1 END)::INTEGER as healthy_breakers,
        COUNT(CASE WHEN state = 'open' OR (state = 'half_open' AND success_count = 0) THEN 1 END)::INTEGER as unhealthy_breakers,
        COUNT(CASE WHEN state = 'open' THEN 1 END)::INTEGER as open_breakers,
        COUNT(CASE WHEN state = 'half_open' THEN 1 END)::INTEGER as half_open_breakers,
        COUNT(CASE WHEN state = 'closed' THEN 1 END)::INTEGER as closed_breakers,
        CASE 
            WHEN COUNT(*) > 0 THEN
                ROUND(COUNT(CASE WHEN state = 'closed' OR (state = 'half_open' AND success_count > 0) THEN 1 END)::NUMERIC / COUNT(*) * 100, 2)
            ELSE 100
        END as overall_health_percentage
    FROM circuit_breakers
    WHERE is_active = true;
END;
$$ LANGUAGE plpgsql;

-- Função para limpeza de dados antigos
CREATE OR REPLACE FUNCTION cleanup_circuit_breaker_data(p_days INTEGER DEFAULT 30)
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER := 0;
    v_cutoff_date TIMESTAMP WITH TIME ZONE;
BEGIN
    v_cutoff_date := NOW() - (p_days || ' days')::INTERVAL;
    
    -- Limpar execuções antigas
    DELETE FROM circuit_breaker_executions
    WHERE executed_at < v_cutoff_date;
    
    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- VIEWS ÚTEIS
-- ========================================

-- View para circuit breakers com estatísticas recentes
CREATE OR REPLACE VIEW v_circuit_breakers_status AS
SELECT 
    cb.id,
    cb.name,
    cb.state,
    cb.failure_count,
    cb.success_count,
    cb.failure_threshold,
    cb.success_threshold,
    cb.state_changed_at,
    EXTRACT(EPOCH FROM (NOW() - cb.state_changed_at))::INTEGER as seconds_in_current_state,
    cb.timeout_seconds,
    (cb.state = 'closed' OR (cb.state = 'half_open' AND cb.success_count > 0)) as is_healthy,
    cb.is_active,
    cb.last_failure_time,
    cb.last_success_time,
    recent_stats.executions_1h,
    recent_stats.failures_1h,
    recent_stats.success_rate_1h
FROM circuit_breakers cb
LEFT JOIN (
    SELECT 
        circuit_breaker_id,
        COUNT(*) as executions_1h,
        COUNT(CASE WHEN NOT success THEN 1 END) as failures_1h,
        CASE 
            WHEN COUNT(*) > 0 THEN
                ROUND(COUNT(CASE WHEN success THEN 1 END)::NUMERIC / COUNT(*) * 100, 2)
            ELSE 0
        END as success_rate_1h
    FROM circuit_breaker_executions
    WHERE executed_at >= NOW() - INTERVAL '1 hour'
    GROUP BY circuit_breaker_id
) recent_stats ON cb.id = recent_stats.circuit_breaker_id
WHERE cb.is_active = true;

-- ========================================
-- DADOS INICIAIS
-- ========================================

-- Criar circuit breakers padrão para DataJud Service
DO $$
BEGIN
    -- Circuit breaker para API DataJud geral
    PERFORM get_or_create_circuit_breaker('datajud-api', 5, 3, 30, 5);
    
    -- Circuit breaker para autenticação DataJud
    PERFORM get_or_create_circuit_breaker('datajud-auth', 3, 2, 60, 3);
    
    -- Circuit breaker para banco de dados
    PERFORM get_or_create_circuit_breaker('database', 10, 5, 15, 10);
    
    -- Circuit breaker para cache (Redis)
    PERFORM get_or_create_circuit_breaker('cache', 5, 3, 10, 8);
    
    -- Circuit breaker para message queue (RabbitMQ)
    PERFORM get_or_create_circuit_breaker('rabbitmq', 5, 3, 20, 5);
    
    RAISE NOTICE 'Circuit breakers padrão criados';
END $$;

-- ========================================
-- COMENTÁRIOS
-- ========================================

COMMENT ON TABLE circuit_breakers IS 'Circuit breakers para proteção contra falhas em cascata';
COMMENT ON TABLE circuit_breaker_executions IS 'Histórico de execuções dos circuit breakers';

COMMENT ON COLUMN circuit_breakers.name IS 'Nome único do circuit breaker';
COMMENT ON COLUMN circuit_breakers.state IS 'Estado atual: closed, open, half_open';
COMMENT ON COLUMN circuit_breakers.failure_threshold IS 'Número de falhas para abrir o circuit breaker';
COMMENT ON COLUMN circuit_breakers.success_threshold IS 'Sucessos para fechar quando half-open';
COMMENT ON COLUMN circuit_breakers.timeout_seconds IS 'Tempo para tentar half-open quando open';
COMMENT ON COLUMN circuit_breakers.max_requests IS 'Max requests quando half-open';

-- ========================================
-- FINALIZAÇÃO
-- ========================================

-- Atualizar estatísticas
ANALYZE circuit_breakers;
ANALYZE circuit_breaker_executions;

-- Log de conclusão
DO $$
BEGIN
    RAISE NOTICE 'Migração 004 concluída - sistema de circuit breakers criado';
    RAISE NOTICE 'Circuit breakers padrão: %', (SELECT COUNT(*) FROM circuit_breakers);
END $$;