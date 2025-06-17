-- Migration: Create indexes for performance
-- Description: Índices para otimizar consultas nas tabelas principais
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- ÍNDICES PARA A TABELA PROCESSES
-- ========================================

-- Índices básicos de chaves estrangeiras e filtros comuns
CREATE INDEX IF NOT EXISTS idx_processes_tenant_id ON processes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_processes_client_id ON processes(client_id);
CREATE INDEX IF NOT EXISTS idx_processes_status ON processes(status);
CREATE INDEX IF NOT EXISTS idx_processes_stage ON processes(stage);
CREATE INDEX IF NOT EXISTS idx_processes_court_id ON processes(court_id);
CREATE INDEX IF NOT EXISTS idx_processes_judge_id ON processes(judge_id) WHERE judge_id IS NOT NULL;

-- Índices para ordenação e filtragem por data
CREATE INDEX IF NOT EXISTS idx_processes_created_at ON processes(created_at);
CREATE INDEX IF NOT EXISTS idx_processes_updated_at ON processes(updated_at);
CREATE INDEX IF NOT EXISTS idx_processes_last_movement_at ON processes(last_movement_at) WHERE last_movement_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_processes_last_sync_at ON processes(last_sync_at) WHERE last_sync_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_processes_archived_at ON processes(archived_at) WHERE archived_at IS NOT NULL;

-- Índices para monitoramento
CREATE INDEX IF NOT EXISTS idx_processes_monitoring_enabled ON processes((monitoring->>'enabled')) WHERE monitoring->>'enabled' = 'true';
CREATE INDEX IF NOT EXISTS idx_processes_auto_sync ON processes((monitoring->>'auto_sync')) WHERE monitoring->>'auto_sync' = 'true';

-- Índices compostos para queries frequentes
CREATE INDEX IF NOT EXISTS idx_processes_tenant_status ON processes(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_processes_tenant_client ON processes(tenant_id, client_id);
CREATE INDEX IF NOT EXISTS idx_processes_status_updated ON processes(status, updated_at);

-- Índice para busca por número do processo
CREATE INDEX IF NOT EXISTS idx_processes_number_unique ON processes(number);

-- Índice GIN para busca em tags
CREATE INDEX IF NOT EXISTS idx_processes_tags ON processes USING gin(tags);

-- Índice para busca textual
CREATE INDEX IF NOT EXISTS idx_processes_search ON processes USING gin(
    to_tsvector('portuguese', 
        COALESCE(title, '') || ' ' || 
        COALESCE(description, '') || ' ' || 
        COALESCE(number, '')
    )
);

-- Índice para processos que precisam de sincronização
CREATE INDEX IF NOT EXISTS idx_processes_needs_sync ON processes(tenant_id, last_sync_at, (monitoring->>'enabled'), (monitoring->>'auto_sync'))
WHERE monitoring->>'enabled' = 'true' AND monitoring->>'auto_sync' = 'true';

-- ========================================
-- ÍNDICES PARA A TABELA MOVEMENTS
-- ========================================

-- Índices básicos de chaves estrangeiras
CREATE INDEX IF NOT EXISTS idx_movements_process_id ON movements(process_id);
CREATE INDEX IF NOT EXISTS idx_movements_tenant_id ON movements(tenant_id);

-- Índices para filtros comuns
CREATE INDEX IF NOT EXISTS idx_movements_type ON movements(type);
CREATE INDEX IF NOT EXISTS idx_movements_date ON movements(date);
CREATE INDEX IF NOT EXISTS idx_movements_sequence ON movements(process_id, sequence);
CREATE INDEX IF NOT EXISTS idx_movements_external_id ON movements(external_id) WHERE external_id IS NOT NULL;

-- Índices para flags importantes
CREATE INDEX IF NOT EXISTS idx_movements_important ON movements(is_important) WHERE is_important = true;
CREATE INDEX IF NOT EXISTS idx_movements_public ON movements(is_public);
CREATE INDEX IF NOT EXISTS idx_movements_notification_pending ON movements(tenant_id, is_important, notification_sent) 
WHERE is_important = true AND notification_sent = false;

-- Índices para ordenação por data
CREATE INDEX IF NOT EXISTS idx_movements_created_at ON movements(created_at);
CREATE INDEX IF NOT EXISTS idx_movements_updated_at ON movements(updated_at);
CREATE INDEX IF NOT EXISTS idx_movements_synced_at ON movements(synced_at);

-- Índices compostos para queries frequentes
CREATE INDEX IF NOT EXISTS idx_movements_process_date ON movements(process_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_movements_tenant_date ON movements(tenant_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_movements_tenant_important ON movements(tenant_id, is_important) WHERE is_important = true;

-- Índice para juiz
CREATE INDEX IF NOT EXISTS idx_movements_judge ON movements(judge) WHERE judge IS NOT NULL;

-- Índice GIN para tags
CREATE INDEX IF NOT EXISTS idx_movements_tags ON movements USING gin(tags);

-- Índice para busca textual completa
CREATE INDEX IF NOT EXISTS idx_movements_search ON movements USING gin(
    to_tsvector('portuguese', 
        COALESCE(title, '') || ' ' || 
        COALESCE(description, '') || ' ' || 
        COALESCE(content, '')
    )
);

-- Índice para análise por IA (JSONB)
CREATE INDEX IF NOT EXISTS idx_movements_analysis_importance ON movements((metadata->'analysis'->>'importance'));
CREATE INDEX IF NOT EXISTS idx_movements_analysis_sentiment ON movements((metadata->'analysis'->>'sentiment'));
CREATE INDEX IF NOT EXISTS idx_movements_analysis_category ON movements((metadata->'analysis'->>'category'));

-- Índice para movimentações com deadline
CREATE INDEX IF NOT EXISTS idx_movements_deadline ON movements((metadata->'analysis'->>'deadline_date')) 
WHERE (metadata->'analysis'->>'has_deadline')::boolean = true;

-- ========================================
-- ÍNDICES PARA A TABELA PARTIES
-- ========================================

-- Índices básicos
CREATE INDEX IF NOT EXISTS idx_parties_process_id ON parties(process_id);
CREATE INDEX IF NOT EXISTS idx_parties_type ON parties(type);
CREATE INDEX IF NOT EXISTS idx_parties_role ON parties(role);
CREATE INDEX IF NOT EXISTS idx_parties_active ON parties(is_active) WHERE is_active = true;

-- Índice para busca por documento
CREATE INDEX IF NOT EXISTS idx_parties_document ON parties(document) WHERE document IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_parties_document_type ON parties(document_type) WHERE document_type IS NOT NULL;

-- Índice composto para documento único por tipo
CREATE UNIQUE INDEX IF NOT EXISTS idx_parties_document_unique ON parties(document, document_type) 
WHERE document IS NOT NULL AND document_type IS NOT NULL;

-- Índices para datas
CREATE INDEX IF NOT EXISTS idx_parties_created_at ON parties(created_at);
CREATE INDEX IF NOT EXISTS idx_parties_updated_at ON parties(updated_at);

-- Índices compostos frequentes
CREATE INDEX IF NOT EXISTS idx_parties_process_role ON parties(process_id, role);
CREATE INDEX IF NOT EXISTS idx_parties_process_active ON parties(process_id, is_active) WHERE is_active = true;

-- Índice para busca textual em nomes
CREATE INDEX IF NOT EXISTS idx_parties_search ON parties USING gin(
    to_tsvector('portuguese', 
        COALESCE(name, '') || ' ' || 
        COALESCE(document, '')
    )
);

-- Índice para advogados (JSONB)
CREATE INDEX IF NOT EXISTS idx_parties_has_lawyer ON parties((CASE WHEN lawyer IS NOT NULL THEN true ELSE false END));
CREATE INDEX IF NOT EXISTS idx_parties_lawyer_oab ON parties((lawyer->>'oab')) WHERE lawyer IS NOT NULL;

-- ========================================
-- ÍNDICES PARA PERFORMANCE DE RELATÓRIOS
-- ========================================

-- Índice composto para dashboard - processos por tenant e período
CREATE INDEX IF NOT EXISTS idx_processes_dashboard ON processes(tenant_id, status, created_at, updated_at);

-- Índice composto para dashboard - movimentações por tenant e período  
CREATE INDEX IF NOT EXISTS idx_movements_dashboard ON movements(tenant_id, type, date, is_important);

-- Índice para estatísticas mensais
CREATE INDEX IF NOT EXISTS idx_movements_monthly_stats ON movements(tenant_id, date_trunc('month', date), type);

-- Índice para processos mais ativos
CREATE INDEX IF NOT EXISTS idx_movements_process_count ON movements(process_id, date) WHERE date >= CURRENT_DATE - INTERVAL '30 days';

-- ========================================
-- COMENTÁRIOS NOS ÍNDICES
-- ========================================

COMMENT ON INDEX idx_processes_tenant_id IS 'Busca rápida de processos por tenant';
COMMENT ON INDEX idx_processes_search IS 'Busca textual full-text em processos';
COMMENT ON INDEX idx_movements_search IS 'Busca textual full-text em movimentações';
COMMENT ON INDEX idx_movements_notification_pending IS 'Notificações pendentes para envio';
COMMENT ON INDEX idx_processes_needs_sync IS 'Processos que precisam sincronização com DataJud';
COMMENT ON INDEX idx_parties_document_unique IS 'Garante unicidade de documento por tipo';

-- ========================================
-- ESTATÍSTICAS PARA O QUERY PLANNER
-- ========================================

-- Atualizar estatísticas das tabelas
ANALYZE processes;
ANALYZE movements;  
ANALYZE parties;