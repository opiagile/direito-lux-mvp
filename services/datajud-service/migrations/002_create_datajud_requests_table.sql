-- Migration: Create DataJud requests table
-- Description: Tabela para gerenciar requisições para a API DataJud
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- TABELA DE REQUISIÇÕES DATAJUD
-- ========================================

CREATE TABLE datajud_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    client_id UUID NOT NULL,
    process_id UUID,
    type VARCHAR(20) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 2,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    cnpj_provider_id UUID,
    
    -- Parâmetros da requisição
    process_number VARCHAR(25),
    court_id VARCHAR(10),
    parameters JSONB DEFAULT '{}',
    
    -- Controle de cache
    cache_key VARCHAR(255),
    cache_ttl INTEGER DEFAULT 3600,
    use_cache BOOLEAN DEFAULT true,
    
    -- Controle de retry
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    retry_after TIMESTAMP WITH TIME ZONE,
    
    -- Circuit breaker
    circuit_breaker_key VARCHAR(100),
    
    -- Timestamps
    requested_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processing_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Resultado
    response_data JSONB,
    error_message TEXT,
    error_code VARCHAR(50),
    
    -- Constraints
    CONSTRAINT chk_type_valid CHECK (type IN ('process', 'movement', 'party', 'document', 'bulk')),
    CONSTRAINT chk_priority_range CHECK (priority BETWEEN 1 AND 4),
    CONSTRAINT chk_status_valid CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cached', 'retrying')),
    CONSTRAINT chk_retry_count_non_negative CHECK (retry_count >= 0),
    CONSTRAINT chk_max_retries_positive CHECK (max_retries > 0),
    CONSTRAINT chk_retry_count_not_exceed_max CHECK (retry_count <= max_retries),
    CONSTRAINT chk_cache_ttl_positive CHECK (cache_ttl > 0),
    CONSTRAINT chk_processing_after_requested CHECK (processing_at IS NULL OR processing_at >= requested_at),
    CONSTRAINT chk_completed_after_processing CHECK (completed_at IS NULL OR completed_at >= COALESCE(processing_at, requested_at)),
    CONSTRAINT chk_process_number_format CHECK (process_number IS NULL OR process_number ~ '^[0-9]{7}-[0-9]{2}\.[0-9]{4}\.[0-9]\.[0-9]{2}\.[0-9]{4}$')
);

-- ========================================
-- ÍNDICES
-- ========================================

-- Índice principal por tenant
CREATE INDEX idx_datajud_requests_tenant_id ON datajud_requests(tenant_id);

-- Índice para requisições pendentes (usado pela queue)
CREATE INDEX idx_datajud_requests_pending ON datajud_requests(status, priority DESC, requested_at ASC)
WHERE status IN ('pending', 'retrying');

-- Índice para retry (requisições que podem ser reprocessadas)
CREATE INDEX idx_datajud_requests_retry ON datajud_requests(status, retry_after)
WHERE status = 'retrying';

-- Índice por status e data para cleanup
CREATE INDEX idx_datajud_requests_status_completed ON datajud_requests(status, completed_at)
WHERE status IN ('completed', 'failed');

-- Índice por processo (para consultar histórico)
CREATE INDEX idx_datajud_requests_process_number ON datajud_requests(process_number)
WHERE process_number IS NOT NULL;

-- Índice por CNPJ provider (para estatísticas)
CREATE INDEX idx_datajud_requests_cnpj_provider ON datajud_requests(cnpj_provider_id)
WHERE cnpj_provider_id IS NOT NULL;

-- Índice para cache key (busca rápida no cache)
CREATE INDEX idx_datajud_requests_cache_key ON datajud_requests(cache_key)
WHERE cache_key IS NOT NULL;

-- Índice para circuit breaker
CREATE INDEX idx_datajud_requests_circuit_breaker ON datajud_requests(circuit_breaker_key, status, completed_at)
WHERE circuit_breaker_key IS NOT NULL;

-- Índice composto para estatísticas
CREATE INDEX idx_datajud_requests_stats ON datajud_requests(tenant_id, type, status, created_at);

-- Índice para cleanup por data
CREATE INDEX idx_datajud_requests_cleanup ON datajud_requests(status, created_at)
WHERE status IN ('completed', 'failed');

-- ========================================
-- TRIGGERS
-- ========================================

-- Trigger para atualizar updated_at
CREATE TRIGGER update_datajud_requests_updated_at
    BEFORE UPDATE ON datajud_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para validação e auto-preenchimento
CREATE OR REPLACE FUNCTION validate_datajud_request()
RETURNS TRIGGER AS $$
BEGIN
    -- Gerar cache key se não foi definida
    IF NEW.cache_key IS NULL AND NEW.use_cache = true THEN
        NEW.cache_key := 'datajud:' || NEW.type || ':' || 
                        COALESCE(NEW.process_number, '') || ':' ||
                        COALESCE(NEW.court_id, '') || ':' ||
                        md5(NEW.parameters::text);
    END IF;
    
    -- Validar número de processo CNJ se fornecido
    IF NEW.process_number IS NOT NULL THEN
        -- Remover formatação e verificar se tem 20 dígitos
        IF length(regexp_replace(NEW.process_number, '[^0-9]', '', 'g')) != 20 THEN
            RAISE EXCEPTION 'Número de processo CNJ deve ter 20 dígitos: %', NEW.process_number;
        END IF;
    END IF;
    
    -- Definir timestamps baseado no status
    IF NEW.status = 'processing' AND OLD.status != 'processing' THEN
        NEW.processing_at := NOW();
    END IF;
    
    IF NEW.status IN ('completed', 'failed') AND OLD.status NOT IN ('completed', 'failed') THEN
        NEW.completed_at := NOW();
    END IF;
    
    -- Limpar retry_after se não está mais retrying
    IF NEW.status != 'retrying' THEN
        NEW.retry_after := NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_datajud_request_trigger
    BEFORE INSERT OR UPDATE ON datajud_requests
    FOR EACH ROW
    EXECUTE FUNCTION validate_datajud_request();

-- ========================================
-- FUNÇÕES UTILITÁRIAS
-- ========================================

-- Função para obter próxima requisição da fila
CREATE OR REPLACE FUNCTION get_next_pending_request()
RETURNS TABLE(
    request_id UUID,
    tenant_id UUID,
    type VARCHAR(20),
    priority INTEGER,
    process_number VARCHAR(25),
    court_id VARCHAR(10),
    parameters JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        dr.id,
        dr.tenant_id,
        dr.type,
        dr.priority,
        dr.process_number,
        dr.court_id,
        dr.parameters
    FROM datajud_requests dr
    WHERE dr.status IN ('pending', 'retrying')
      AND (dr.retry_after IS NULL OR dr.retry_after <= NOW())
    ORDER BY 
        dr.priority DESC,
        dr.requested_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED; -- Evita lock contention entre workers
END;
$$ LANGUAGE plpgsql;

-- Função para marcar requisição como processando
CREATE OR REPLACE FUNCTION start_processing_request(p_request_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE datajud_requests
    SET status = 'processing',
        processing_at = NOW(),
        updated_at = NOW()
    WHERE id = p_request_id 
      AND status IN ('pending', 'retrying');
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para completar requisição com sucesso
CREATE OR REPLACE FUNCTION complete_request(
    p_request_id UUID,
    p_response_data JSONB
)
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE datajud_requests
    SET status = 'completed',
        completed_at = NOW(),
        response_data = p_response_data,
        error_message = NULL,
        error_code = NULL,
        updated_at = NOW()
    WHERE id = p_request_id 
      AND status = 'processing';
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para falhar requisição
CREATE OR REPLACE FUNCTION fail_request(
    p_request_id UUID,
    p_error_code VARCHAR(50),
    p_error_message TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
BEGIN
    UPDATE datajud_requests
    SET status = 'failed',
        completed_at = NOW(),
        error_code = p_error_code,
        error_message = p_error_message,
        updated_at = NOW()
    WHERE id = p_request_id 
      AND status IN ('processing', 'retrying');
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para agendar retry
CREATE OR REPLACE FUNCTION schedule_retry(
    p_request_id UUID,
    p_retry_after TIMESTAMP WITH TIME ZONE
)
RETURNS BOOLEAN AS $$
DECLARE
    v_rows_affected INTEGER;
    v_current_retry_count INTEGER;
    v_max_retries INTEGER;
BEGIN
    -- Verificar se pode fazer retry
    SELECT retry_count, max_retries 
    INTO v_current_retry_count, v_max_retries
    FROM datajud_requests
    WHERE id = p_request_id;
    
    IF v_current_retry_count >= v_max_retries THEN
        -- Marcar como falhada definitivamente
        UPDATE datajud_requests
        SET status = 'failed',
            completed_at = NOW(),
            error_code = 'MAX_RETRIES_EXCEEDED',
            error_message = 'Máximo de tentativas excedido',
            updated_at = NOW()
        WHERE id = p_request_id;
        
        RETURN FALSE;
    END IF;
    
    -- Agendar retry
    UPDATE datajud_requests
    SET status = 'retrying',
        retry_count = retry_count + 1,
        retry_after = p_retry_after,
        updated_at = NOW()
    WHERE id = p_request_id 
      AND status IN ('processing', 'failed');
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para estatísticas de requisições
CREATE OR REPLACE FUNCTION get_request_statistics(
    p_tenant_id UUID DEFAULT NULL,
    p_hours INTEGER DEFAULT 24
)
RETURNS TABLE(
    total_requests BIGINT,
    completed_requests BIGINT,
    failed_requests BIGINT,
    pending_requests BIGINT,
    processing_requests BIGINT,
    retrying_requests BIGINT,
    success_rate NUMERIC,
    avg_response_time_ms NUMERIC,
    requests_by_type JSONB,
    hourly_stats JSONB
) AS $$
DECLARE
    v_start_time TIMESTAMP WITH TIME ZONE;
    v_stats RECORD;
    v_hourly JSONB;
    v_by_type JSONB;
BEGIN
    v_start_time := NOW() - (p_hours || ' hours')::INTERVAL;
    
    -- Estatísticas gerais
    SELECT 
        COUNT(*) as total,
        COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
        COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
        COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
        COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing,
        COUNT(CASE WHEN status = 'retrying' THEN 1 END) as retrying,
        AVG(EXTRACT(EPOCH FROM (completed_at - requested_at)) * 1000) as avg_time
    INTO v_stats
    FROM datajud_requests
    WHERE created_at >= v_start_time
      AND (p_tenant_id IS NULL OR tenant_id = p_tenant_id);
    
    -- Estatísticas por tipo
    SELECT jsonb_object_agg(type, count) INTO v_by_type
    FROM (
        SELECT type, COUNT(*) as count
        FROM datajud_requests
        WHERE created_at >= v_start_time
          AND (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
        GROUP BY type
    ) t;
    
    -- Estatísticas por hora
    SELECT jsonb_agg(
        jsonb_build_object(
            'hour', hour,
            'total', total,
            'completed', completed,
            'failed', failed
        )
    ) INTO v_hourly
    FROM (
        SELECT 
            DATE_TRUNC('hour', created_at) as hour,
            COUNT(*) as total,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
            COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed
        FROM datajud_requests
        WHERE created_at >= v_start_time
          AND (p_tenant_id IS NULL OR tenant_id = p_tenant_id)
        GROUP BY DATE_TRUNC('hour', created_at)
        ORDER BY hour
    ) h;
    
    RETURN QUERY
    SELECT 
        v_stats.total,
        v_stats.completed,
        v_stats.failed,
        v_stats.pending,
        v_stats.processing,
        v_stats.retrying,
        CASE 
            WHEN v_stats.total > 0 THEN 
                ROUND(v_stats.completed::NUMERIC / v_stats.total * 100, 2)
            ELSE 0
        END,
        ROUND(COALESCE(v_stats.avg_time, 0), 2),
        COALESCE(v_by_type, '{}'::jsonb),
        COALESCE(v_hourly, '[]'::jsonb);
END;
$$ LANGUAGE plpgsql;

-- Função para limpeza de requisições antigas
CREATE OR REPLACE FUNCTION cleanup_old_requests(p_days INTEGER DEFAULT 30)
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER := 0;
    v_cutoff_date TIMESTAMP WITH TIME ZONE;
BEGIN
    v_cutoff_date := NOW() - (p_days || ' days')::INTERVAL;
    
    -- Deletar requisições completadas/falhadas antigas
    DELETE FROM datajud_requests
    WHERE status IN ('completed', 'failed')
      AND completed_at < v_cutoff_date;
    
    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- COMENTÁRIOS
-- ========================================

COMMENT ON TABLE datajud_requests IS 'Requisições para a API DataJud com controle de fila, retry e cache';

COMMENT ON COLUMN datajud_requests.id IS 'Identificador único da requisição';
COMMENT ON COLUMN datajud_requests.tenant_id IS 'ID do tenant que fez a requisição';
COMMENT ON COLUMN datajud_requests.client_id IS 'ID do cliente/usuário que fez a requisição';
COMMENT ON COLUMN datajud_requests.process_id IS 'ID do processo relacionado (se aplicável)';
COMMENT ON COLUMN datajud_requests.type IS 'Tipo de consulta: process, movement, party, document, bulk';
COMMENT ON COLUMN datajud_requests.priority IS 'Prioridade da requisição (1=urgente, 4=baixa)';
COMMENT ON COLUMN datajud_requests.status IS 'Status atual: pending, processing, completed, failed, cached, retrying';
COMMENT ON COLUMN datajud_requests.cnpj_provider_id IS 'CNPJ provider usado para a requisição';
COMMENT ON COLUMN datajud_requests.process_number IS 'Número do processo CNJ formatado';
COMMENT ON COLUMN datajud_requests.court_id IS 'ID do tribunal (ex: TJSP, TJRJ)';
COMMENT ON COLUMN datajud_requests.parameters IS 'Parâmetros adicionais da requisição';
COMMENT ON COLUMN datajud_requests.cache_key IS 'Chave para busca no cache';
COMMENT ON COLUMN datajud_requests.cache_ttl IS 'TTL do cache em segundos';
COMMENT ON COLUMN datajud_requests.use_cache IS 'Se deve usar cache para esta requisição';
COMMENT ON COLUMN datajud_requests.retry_count IS 'Número de tentativas já realizadas';
COMMENT ON COLUMN datajud_requests.max_retries IS 'Máximo de tentativas permitidas';
COMMENT ON COLUMN datajud_requests.retry_after IS 'Quando pode tentar novamente';
COMMENT ON COLUMN datajud_requests.circuit_breaker_key IS 'Chave do circuit breaker';
COMMENT ON COLUMN datajud_requests.response_data IS 'Dados da resposta da API';
COMMENT ON COLUMN datajud_requests.error_message IS 'Mensagem de erro se falhou';
COMMENT ON COLUMN datajud_requests.error_code IS 'Código do erro se falhou';

-- ========================================
-- FINALIZAÇÃO
-- ========================================

-- Atualizar estatísticas
ANALYZE datajud_requests;

-- Log de conclusão
DO $$
BEGIN
    RAISE NOTICE 'Migração 002 concluída - tabela datajud_requests criada';
    RAISE NOTICE 'Triggers e funções utilitárias configuradas';
END $$;