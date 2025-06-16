-- Migration: Create quota_limits table
-- Description: Tabela para armazenar os limites de quota dos tenants baseados em seus planos

CREATE TABLE quota_limits (
    tenant_id UUID PRIMARY KEY,
    max_processes INTEGER NOT NULL DEFAULT 0,
    max_users INTEGER NOT NULL DEFAULT 0,
    max_clients INTEGER NOT NULL DEFAULT 0,
    datajud_queries_daily INTEGER NOT NULL DEFAULT 0,
    ai_queries_monthly INTEGER NOT NULL DEFAULT 0,
    storage_gb INTEGER NOT NULL DEFAULT 0,
    max_webhooks INTEGER NOT NULL DEFAULT 0,
    max_api_calls_daily INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT quota_limits_positive_or_unlimited CHECK (
        max_processes >= -1 AND
        max_users >= -1 AND
        max_clients >= -1 AND
        datajud_queries_daily >= -1 AND
        ai_queries_monthly >= -1 AND
        storage_gb >= -1 AND
        max_webhooks >= -1 AND
        max_api_calls_daily >= -1
    ),

    -- Foreign keys
    CONSTRAINT fk_quota_limits_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_quota_limits_updated_at ON quota_limits(updated_at);

-- Comments
COMMENT ON TABLE quota_limits IS 'Tabela de limites de quota dos tenants baseados em seus planos';
COMMENT ON COLUMN quota_limits.tenant_id IS 'ID do tenant (chave primária)';
COMMENT ON COLUMN quota_limits.max_processes IS 'Limite máximo de processos (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.max_users IS 'Limite máximo de usuários (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.max_clients IS 'Limite máximo de clientes (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.datajud_queries_daily IS 'Limite diário de consultas DataJud (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.ai_queries_monthly IS 'Limite mensal de consultas IA (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.storage_gb IS 'Limite de armazenamento em GB (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.max_webhooks IS 'Limite máximo de webhooks (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.max_api_calls_daily IS 'Limite diário de chamadas API (-1 = ilimitado)';
COMMENT ON COLUMN quota_limits.updated_at IS 'Data da última atualização dos limites';