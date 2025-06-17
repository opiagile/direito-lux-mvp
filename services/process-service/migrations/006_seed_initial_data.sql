-- Migration: Seed initial data
-- Description: Dados iniciais e configurações para desenvolvimento/teste
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- DADOS PARA DESENVOLVIMENTO/TESTE
-- ========================================

-- ATENÇÃO: Esta migração só deve ser executada em ambiente de desenvolvimento/teste
-- Em produção, remover ou comentar as inserções de dados de exemplo

-- Verificar se está em ambiente de desenvolvimento
DO $$
DECLARE
    env_name TEXT;
BEGIN
    env_name := COALESCE(current_setting('app.environment', true), 'development');
    
    IF env_name NOT IN ('development', 'dev', 'test', 'testing') THEN
        RAISE NOTICE 'Pulando seed de dados - ambiente: %', env_name;
        RETURN;
    END IF;
    
    RAISE NOTICE 'Inserindo dados de desenvolvimento - ambiente: %', env_name;
END $$;

-- ========================================
-- EXTENSÕES NECESSÁRIAS
-- ========================================

-- Extensão para similaridade de texto (se disponível)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Extensão para funções de hash
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ========================================
-- CONFIGURAÇÕES DO SISTEMA
-- ========================================

-- Tabela de configurações do sistema (se necessário)
CREATE TABLE IF NOT EXISTS system_config (
    key VARCHAR(100) PRIMARY KEY,
    value JSONB NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Configurações padrão
INSERT INTO system_config (key, value, description) VALUES
    ('cnj_validation', '{"enabled": true, "strict": false}', 'Configurações de validação CNJ'),
    ('analysis_keywords', '{
        "importance": ["sentença", "acórdão", "decisão", "despacho", "citação", "intimação", "audiência", "recurso"],
        "positive": ["deferido", "procedente", "ganho", "favorável", "aprovado"],
        "negative": ["indeferido", "improcedente", "negado", "desfavorável", "rejeitado"]
    }', 'Palavras-chave para análise automática'),
    ('sync_intervals', '{"default": 24, "max": 168, "min": 1}', 'Intervalos de sincronização em horas'),
    ('notification_templates', '{"movement_created": "Nova movimentação no processo {process_number}"}', 'Templates de notificação')
ON CONFLICT (key) DO NOTHING;

-- ========================================
-- DADOS DE EXEMPLO (APENAS DESENVOLVIMENTO)
-- ========================================

-- Função para criar dados de exemplo
CREATE OR REPLACE FUNCTION create_sample_data()
RETURNS VOID AS $$
DECLARE
    sample_tenant_id UUID := '12345678-1234-1234-1234-123456789012';
    sample_client_id UUID := '87654321-4321-4321-4321-210987654321';
    sample_process_id UUID;
    sample_party_id UUID;
    i INTEGER;
BEGIN
    -- Verificar se dados já existem
    IF EXISTS (SELECT 1 FROM processes WHERE tenant_id = sample_tenant_id) THEN
        RAISE NOTICE 'Dados de exemplo já existem, pulando criação';
        RETURN;
    END IF;

    -- Criar processo de exemplo
    INSERT INTO processes (
        id, tenant_id, client_id, number, title, description, status, stage,
        subject, court_id, monitoring, tags, created_at, updated_at
    ) VALUES (
        gen_random_uuid(),
        sample_tenant_id,
        sample_client_id,
        '1234567-89.2024.1.01.0001',
        'Ação de Cobrança - Exemplo',
        'Processo de exemplo para demonstração do sistema Direito Lux',
        'active',
        'knowledge',
        '{"code": "1234", "description": "Direito Civil - Cobrança", "parent_code": "12"}',
        'TJSP',
        '{"enabled": true, "notification_channels": ["email", "whatsapp"], "auto_sync": true, "sync_interval_hours": 24}',
        ARRAY['cobrança', 'civil', 'exemplo'],
        CURRENT_TIMESTAMP - INTERVAL '30 days',
        CURRENT_TIMESTAMP
    )
    RETURNING id INTO sample_process_id;

    -- Criar partes do processo
    INSERT INTO parties (
        id, process_id, type, name, document, document_type, role, is_active,
        contact, address, created_at, updated_at
    ) VALUES 
    (
        gen_random_uuid(),
        sample_process_id,
        'legal',
        'Empresa Exemplo LTDA',
        '12345678000195',
        'cnpj',
        'plaintiff',
        true,
        '{"email": "contato@exemplo.com", "phone": "(11) 1234-5678"}',
        '{"street": "Rua Exemplo", "number": "123", "city": "São Paulo", "state": "SP", "zip_code": "01234-567"}',
        CURRENT_TIMESTAMP - INTERVAL '30 days',
        CURRENT_TIMESTAMP
    ),
    (
        gen_random_uuid(),
        sample_process_id,
        'individual',
        'João da Silva Exemplo',
        '12345678901',
        'cpf',
        'defendant',
        true,
        '{"email": "joao@exemplo.com", "phone": "(11) 9876-5432"}',
        '{"street": "Avenida Teste", "number": "456", "city": "São Paulo", "state": "SP", "zip_code": "09876-543"}',
        CURRENT_TIMESTAMP - INTERVAL '30 days',
        CURRENT_TIMESTAMP
    );

    -- Criar algumas movimentações de exemplo
    FOR i IN 1..5 LOOP
        INSERT INTO movements (
            id, process_id, tenant_id, sequence, date, type, code,
            title, description, content, is_important, is_public,
            metadata, created_at, updated_at, synced_at
        ) VALUES (
            gen_random_uuid(),
            sample_process_id,
            sample_tenant_id,
            i,
            CURRENT_TIMESTAMP - INTERVAL '30 days' + (i * INTERVAL '5 days'),
            CASE i 
                WHEN 1 THEN 'filing'
                WHEN 2 THEN 'citation'
                WHEN 3 THEN 'order'
                WHEN 4 THEN 'filing'
                ELSE 'decision'
            END,
            'MOV' || LPAD(i::TEXT, 3, '0'),
            'Movimentação de Exemplo ' || i,
            'Descrição detalhada da movimentação ' || i || ' do processo de exemplo.',
            'Conteúdo completo da movimentação para demonstração do sistema.',
            CASE WHEN i IN (3, 5) THEN true ELSE false END,
            true,
            '{"original_source": "exemplo", "analysis": {"importance": ' || (i % 5 + 1) || ', "confidence": 0.8}}',
            CURRENT_TIMESTAMP - INTERVAL '30 days' + (i * INTERVAL '5 days'),
            CURRENT_TIMESTAMP - INTERVAL '30 days' + (i * INTERVAL '5 days'),
            CURRENT_TIMESTAMP - INTERVAL '30 days' + (i * INTERVAL '5 days')
        );
    END LOOP;

    RAISE NOTICE 'Dados de exemplo criados com sucesso!';
    RAISE NOTICE 'Tenant ID: %', sample_tenant_id;
    RAISE NOTICE 'Process ID: %', sample_process_id;
    
END;
$$ LANGUAGE plpgsql;

-- Executar criação de dados apenas em desenvolvimento
DO $$
DECLARE
    env_name TEXT;
BEGIN
    env_name := COALESCE(current_setting('app.environment', true), 'development');
    
    IF env_name IN ('development', 'dev', 'test', 'testing') THEN
        PERFORM create_sample_data();
    END IF;
END $$;

-- ========================================
-- VIEWS ÚTEIS PARA DESENVOLVIMENTO
-- ========================================

-- View para processos com estatísticas
CREATE OR REPLACE VIEW v_processes_with_stats AS
SELECT 
    p.id,
    p.tenant_id,
    p.client_id,
    p.number,
    p.title,
    p.status,
    p.stage,
    p.created_at,
    p.updated_at,
    p.last_movement_at,
    COALESCE(stats.total_movements, 0) as total_movements,
    COALESCE(stats.important_movements, 0) as important_movements,
    COALESCE(stats.days_since_last_movement, 0) as days_since_last_movement
FROM processes p
LEFT JOIN (
    SELECT 
        process_id,
        COUNT(*) as total_movements,
        COUNT(CASE WHEN is_important THEN 1 END) as important_movements,
        EXTRACT(days FROM CURRENT_TIMESTAMP - MAX(date))::INTEGER as days_since_last_movement
    FROM movements
    GROUP BY process_id
) stats ON p.id = stats.process_id;

-- View para movimentações importantes pendentes
CREATE OR REPLACE VIEW v_pending_notifications AS
SELECT 
    m.id as movement_id,
    m.process_id,
    p.number as process_number,
    p.title as process_title,
    p.tenant_id,
    m.title as movement_title,
    m.date as movement_date,
    m.type as movement_type,
    p.monitoring->>'notification_channels' as notification_channels
FROM movements m
JOIN processes p ON m.process_id = p.id
WHERE m.is_important = true 
  AND m.notification_sent = false
  AND p.monitoring->>'enabled' = 'true'
ORDER BY m.date DESC;

-- View para dashboard de estatísticas
CREATE OR REPLACE VIEW v_dashboard_stats AS
SELECT 
    'processes_total' as metric,
    COUNT(*)::TEXT as value,
    'Total de processos' as description
FROM processes
WHERE status != 'archived'

UNION ALL

SELECT 
    'processes_active' as metric,
    COUNT(*)::TEXT as value,
    'Processos ativos' as description
FROM processes
WHERE status = 'active'

UNION ALL

SELECT 
    'movements_today' as metric,
    COUNT(*)::TEXT as value,
    'Movimentações hoje' as description
FROM movements
WHERE date::date = CURRENT_DATE

UNION ALL

SELECT 
    'notifications_pending' as metric,
    COUNT(*)::TEXT as value,
    'Notificações pendentes' as description
FROM v_pending_notifications;

-- ========================================
-- COMENTÁRIOS NAS VIEWS
-- ========================================

COMMENT ON VIEW v_processes_with_stats IS 'Processos com estatísticas de movimentações';
COMMENT ON VIEW v_pending_notifications IS 'Movimentações importantes aguardando notificação';
COMMENT ON VIEW v_dashboard_stats IS 'Estatísticas agregadas para dashboard';

-- ========================================
-- FINALIZAÇÃO
-- ========================================

-- Atualizar estatísticas após inserções
ANALYZE processes;
ANALYZE movements;
ANALYZE parties;

-- Log de conclusão
DO $$
BEGIN
    RAISE NOTICE 'Migração 006 concluída - dados iniciais criados';
    RAISE NOTICE 'Total de processos: %', (SELECT COUNT(*) FROM processes);
    RAISE NOTICE 'Total de movimentações: %', (SELECT COUNT(*) FROM movements);
    RAISE NOTICE 'Total de partes: %', (SELECT COUNT(*) FROM parties);
END $$;