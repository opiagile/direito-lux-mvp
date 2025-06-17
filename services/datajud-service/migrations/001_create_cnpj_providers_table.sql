-- Migration: Create CNPJ providers table
-- Description: Tabela para gerenciar múltiplos CNPJs para consultas DataJud
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- EXTENSÕES NECESSÁRIAS
-- ========================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========================================
-- TABELA DE CNPJ PROVIDERS
-- ========================================

CREATE TABLE cnpj_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    cnpj VARCHAR(18) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    api_key TEXT NOT NULL,
    certificate TEXT,
    certificate_pass TEXT,
    daily_limit INTEGER NOT NULL DEFAULT 10000,
    daily_usage INTEGER NOT NULL DEFAULT 0,
    usage_reset_time TIMESTAMP WITH TIME ZONE DEFAULT (DATE_TRUNC('day', NOW()) + INTERVAL '1 day'),
    is_active BOOLEAN NOT NULL DEFAULT true,
    priority INTEGER NOT NULL DEFAULT 1,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deactivated_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT chk_cnpj_format CHECK (cnpj ~ '^\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}$'),
    CONSTRAINT chk_email_format CHECK (email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT chk_daily_limit_positive CHECK (daily_limit > 0),
    CONSTRAINT chk_daily_usage_non_negative CHECK (daily_usage >= 0),
    CONSTRAINT chk_daily_usage_not_exceed_limit CHECK (daily_usage <= daily_limit),
    CONSTRAINT chk_priority_range CHECK (priority BETWEEN 1 AND 10),
    CONSTRAINT chk_usage_reset_future CHECK (usage_reset_time > created_at)
);

-- ========================================
-- ÍNDICES
-- ========================================

-- Índice principal por tenant para busca rápida
CREATE INDEX idx_cnpj_providers_tenant_id ON cnpj_providers(tenant_id);

-- Índice para providers ativos ordenados por prioridade
CREATE INDEX idx_cnpj_providers_active_priority ON cnpj_providers(is_active, priority, last_used_at) 
WHERE is_active = true;

-- Índice para busca por CNPJ (único)
CREATE UNIQUE INDEX idx_cnpj_providers_cnpj ON cnpj_providers(cnpj);

-- Índice para providers com quota disponível
CREATE INDEX idx_cnpj_providers_available_quota ON cnpj_providers(is_active, (daily_limit - daily_usage))
WHERE is_active = true AND (daily_limit - daily_usage) > 0;

-- Índice para limpeza por data de desativação
CREATE INDEX idx_cnpj_providers_deactivated_at ON cnpj_providers(deactivated_at)
WHERE deactivated_at IS NOT NULL;

-- Índice composto para otimizar consultas de pool
CREATE INDEX idx_cnpj_providers_pool_selection ON cnpj_providers(tenant_id, is_active, priority, daily_usage);

-- ========================================
-- TRIGGERS
-- ========================================

-- Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_cnpj_providers_updated_at
    BEFORE UPDATE ON cnpj_providers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para validar CNPJ automaticamente
CREATE OR REPLACE FUNCTION validate_cnpj_provider()
RETURNS TRIGGER AS $$
DECLARE
    clean_cnpj TEXT;
BEGIN
    -- Limpar e formatar CNPJ
    clean_cnpj := regexp_replace(NEW.cnpj, '[^0-9]', '', 'g');
    
    -- Verificar se tem 14 dígitos
    IF length(clean_cnpj) != 14 THEN
        RAISE EXCEPTION 'CNPJ deve ter exatamente 14 dígitos: %', NEW.cnpj;
    END IF;
    
    -- Formatar no padrão XX.XXX.XXX/XXXX-XX
    NEW.cnpj := substring(clean_cnpj, 1, 2) || '.' || 
                substring(clean_cnpj, 3, 3) || '.' || 
                substring(clean_cnpj, 6, 3) || '/' || 
                substring(clean_cnpj, 9, 4) || '-' || 
                substring(clean_cnpj, 13, 2);
    
    -- Verificar se não são todos dígitos iguais
    IF clean_cnpj ~ '^(.)\1{13}$' THEN
        RAISE EXCEPTION 'CNPJ não pode ter todos os dígitos iguais: %', NEW.cnpj;
    END IF;
    
    -- Definir data de reset se não foi definida
    IF NEW.usage_reset_time IS NULL THEN
        NEW.usage_reset_time := DATE_TRUNC('day', NOW()) + INTERVAL '1 day';
    END IF;
    
    -- Garantir que daily_usage não exceda daily_limit
    IF NEW.daily_usage > NEW.daily_limit THEN
        NEW.daily_usage := NEW.daily_limit;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_cnpj_provider_trigger
    BEFORE INSERT OR UPDATE ON cnpj_providers
    FOR EACH ROW
    EXECUTE FUNCTION validate_cnpj_provider();

-- Trigger para reset automático de quota diária
CREATE OR REPLACE FUNCTION auto_reset_daily_usage()
RETURNS TRIGGER AS $$
BEGIN
    -- Se passou do tempo de reset, zerar o uso diário
    IF NEW.usage_reset_time <= NOW() THEN
        NEW.daily_usage := 0;
        NEW.usage_reset_time := DATE_TRUNC('day', NOW()) + INTERVAL '1 day';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_reset_daily_usage_trigger
    BEFORE UPDATE ON cnpj_providers
    FOR EACH ROW
    EXECUTE FUNCTION auto_reset_daily_usage();

-- ========================================
-- FUNÇÕES UTILITÁRIAS
-- ========================================

-- Função para obter provider disponível com quota
CREATE OR REPLACE FUNCTION get_available_cnpj_provider(
    p_tenant_id UUID DEFAULT NULL,
    p_min_quota INTEGER DEFAULT 1
)
RETURNS TABLE(
    provider_id UUID,
    cnpj VARCHAR(18),
    available_quota INTEGER,
    priority INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        cp.id,
        cp.cnpj,
        (cp.daily_limit - cp.daily_usage) as available_quota,
        cp.priority
    FROM cnpj_providers cp
    WHERE cp.is_active = true
      AND (cp.daily_limit - cp.daily_usage) >= p_min_quota
      AND (p_tenant_id IS NULL OR cp.tenant_id = p_tenant_id)
      AND (cp.usage_reset_time > NOW() OR cp.usage_reset_time IS NULL)
    ORDER BY 
        cp.priority ASC,
        (cp.daily_limit - cp.daily_usage) DESC,
        cp.last_used_at ASC NULLS FIRST
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Função para usar quota de um provider
CREATE OR REPLACE FUNCTION use_cnpj_quota(
    p_provider_id UUID,
    p_amount INTEGER DEFAULT 1
)
RETURNS BOOLEAN AS $$
DECLARE
    v_available_quota INTEGER;
    v_rows_affected INTEGER;
BEGIN
    -- Verificar quota disponível
    SELECT (daily_limit - daily_usage) INTO v_available_quota
    FROM cnpj_providers
    WHERE id = p_provider_id AND is_active = true;
    
    -- Se não encontrou o provider ou não tem quota
    IF v_available_quota IS NULL OR v_available_quota < p_amount THEN
        RETURN FALSE;
    END IF;
    
    -- Usar a quota
    UPDATE cnpj_providers
    SET daily_usage = daily_usage + p_amount,
        last_used_at = NOW(),
        updated_at = NOW()
    WHERE id = p_provider_id AND is_active = true;
    
    GET DIAGNOSTICS v_rows_affected = ROW_COUNT;
    
    RETURN v_rows_affected > 0;
END;
$$ LANGUAGE plpgsql;

-- Função para reset diário de todos os providers
CREATE OR REPLACE FUNCTION reset_all_daily_usage()
RETURNS INTEGER AS $$
DECLARE
    v_reset_count INTEGER := 0;
BEGIN
    UPDATE cnpj_providers
    SET daily_usage = 0,
        usage_reset_time = DATE_TRUNC('day', NOW()) + INTERVAL '1 day',
        updated_at = NOW()
    WHERE usage_reset_time <= NOW() OR usage_reset_time IS NULL;
    
    GET DIAGNOSTICS v_reset_count = ROW_COUNT;
    
    RETURN v_reset_count;
END;
$$ LANGUAGE plpgsql;

-- Função para obter estatísticas de uso
CREATE OR REPLACE FUNCTION get_cnpj_usage_stats(p_tenant_id UUID DEFAULT NULL)
RETURNS TABLE(
    total_providers INTEGER,
    active_providers INTEGER,
    total_daily_limit INTEGER,
    total_daily_usage INTEGER,
    available_quota INTEGER,
    usage_percentage NUMERIC,
    exhausted_providers INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*)::INTEGER as total_providers,
        COUNT(CASE WHEN is_active THEN 1 END)::INTEGER as active_providers,
        COALESCE(SUM(daily_limit), 0)::INTEGER as total_daily_limit,
        COALESCE(SUM(daily_usage), 0)::INTEGER as total_daily_usage,
        COALESCE(SUM(daily_limit - daily_usage), 0)::INTEGER as available_quota,
        CASE 
            WHEN SUM(daily_limit) > 0 THEN 
                ROUND(SUM(daily_usage)::NUMERIC / SUM(daily_limit) * 100, 2)
            ELSE 0
        END as usage_percentage,
        COUNT(CASE WHEN daily_usage >= daily_limit THEN 1 END)::INTEGER as exhausted_providers
    FROM cnpj_providers
    WHERE (p_tenant_id IS NULL OR tenant_id = p_tenant_id);
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- COMENTÁRIOS NAS TABELAS E COLUNAS
-- ========================================

COMMENT ON TABLE cnpj_providers IS 'Providers CNPJ para consultas na API DataJud com rate limiting por CNPJ';

COMMENT ON COLUMN cnpj_providers.id IS 'Identificador único do provider';
COMMENT ON COLUMN cnpj_providers.tenant_id IS 'ID do tenant proprietário do CNPJ';
COMMENT ON COLUMN cnpj_providers.cnpj IS 'CNPJ formatado (XX.XXX.XXX/XXXX-XX)';
COMMENT ON COLUMN cnpj_providers.name IS 'Nome da empresa/pessoa jurídica';
COMMENT ON COLUMN cnpj_providers.email IS 'Email de contato do provider';
COMMENT ON COLUMN cnpj_providers.api_key IS 'Chave de API para autenticação no DataJud';
COMMENT ON COLUMN cnpj_providers.certificate IS 'Certificado digital para autenticação';
COMMENT ON COLUMN cnpj_providers.certificate_pass IS 'Senha do certificado digital';
COMMENT ON COLUMN cnpj_providers.daily_limit IS 'Limite diário de consultas (padrão 10.000)';
COMMENT ON COLUMN cnpj_providers.daily_usage IS 'Uso atual do dia';
COMMENT ON COLUMN cnpj_providers.usage_reset_time IS 'Quando o contador diário será resetado';
COMMENT ON COLUMN cnpj_providers.is_active IS 'Se o provider está ativo para uso';
COMMENT ON COLUMN cnpj_providers.priority IS 'Prioridade de uso (1=alta, 10=baixa)';
COMMENT ON COLUMN cnpj_providers.last_used_at IS 'Última vez que foi usado';
COMMENT ON COLUMN cnpj_providers.deactivated_at IS 'Data de desativação se aplicável';

-- ========================================
-- DADOS DE EXEMPLO (APENAS DESENVOLVIMENTO)
-- ========================================

-- Inserir dados de exemplo apenas em ambiente de desenvolvimento
DO $$
DECLARE
    env_name TEXT;
    sample_tenant_id UUID := '12345678-1234-1234-1234-123456789012';
BEGIN
    env_name := COALESCE(current_setting('app.environment', true), 'development');
    
    IF env_name IN ('development', 'dev', 'test', 'testing') THEN
        -- Provider de exemplo 1
        INSERT INTO cnpj_providers (
            id, tenant_id, cnpj, name, email, api_key, daily_limit, priority
        ) VALUES (
            '11111111-1111-1111-1111-111111111111',
            sample_tenant_id,
            '12345678000195',
            'Empresa Exemplo 1 LTDA',
            'contato1@exemplo.com',
            'api_key_exemplo_1_secreto',
            10000,
            1
        ) ON CONFLICT (cnpj) DO NOTHING;
        
        -- Provider de exemplo 2  
        INSERT INTO cnpj_providers (
            id, tenant_id, cnpj, name, email, api_key, daily_limit, priority
        ) VALUES (
            '22222222-2222-2222-2222-222222222222',
            sample_tenant_id,
            '98765432000187',
            'Empresa Exemplo 2 LTDA',
            'contato2@exemplo.com',
            'api_key_exemplo_2_secreto',
            10000,
            2
        ) ON CONFLICT (cnpj) DO NOTHING;
        
        -- Provider global para casos sem tenant específico
        INSERT INTO cnpj_providers (
            id, tenant_id, cnpj, name, email, api_key, daily_limit, priority
        ) VALUES (
            '33333333-3333-3333-3333-333333333333',
            '00000000-0000-0000-0000-000000000000',
            '11222333000144',
            'Provider Global Direito Lux',
            'global@direitolux.com',
            'api_key_global_direito_lux',
            50000,
            5
        ) ON CONFLICT (cnpj) DO NOTHING;
        
        RAISE NOTICE 'CNPJs de exemplo criados para desenvolvimento';
    ELSE
        RAISE NOTICE 'Pulando criação de dados de exemplo - ambiente: %', env_name;
    END IF;
END $$;

-- ========================================
-- FINALIZAÇÃO
-- ========================================

-- Atualizar estatísticas após criação
ANALYZE cnpj_providers;

-- Log de conclusão
DO $$
BEGIN
    RAISE NOTICE 'Migração 001 concluída - tabela cnpj_providers criada';
    RAISE NOTICE 'Total de providers: %', (SELECT COUNT(*) FROM cnpj_providers);
END $$;