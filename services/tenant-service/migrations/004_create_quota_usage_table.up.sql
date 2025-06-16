-- Migration: Create quota_usage table
-- Description: Tabela para armazenar o uso atual de quotas dos tenants

CREATE TABLE quota_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    processes_count INTEGER NOT NULL DEFAULT 0,
    users_count INTEGER NOT NULL DEFAULT 0,
    clients_count INTEGER NOT NULL DEFAULT 0,
    datajud_queries_daily INTEGER NOT NULL DEFAULT 0,
    datajud_queries_month INTEGER NOT NULL DEFAULT 0,
    ai_queries_monthly INTEGER NOT NULL DEFAULT 0,
    storage_used_gb DECIMAL(10,3) NOT NULL DEFAULT 0,
    webhooks_count INTEGER NOT NULL DEFAULT 0,
    api_calls_daily INTEGER NOT NULL DEFAULT 0,
    api_calls_monthly INTEGER NOT NULL DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_reset_daily TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_reset_monthly TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT quota_usage_positive_counts CHECK (
        processes_count >= 0 AND
        users_count >= 0 AND
        clients_count >= 0 AND
        datajud_queries_daily >= 0 AND
        datajud_queries_month >= 0 AND
        ai_queries_monthly >= 0 AND
        storage_used_gb >= 0 AND
        webhooks_count >= 0 AND
        api_calls_daily >= 0 AND
        api_calls_monthly >= 0
    ),

    -- Foreign keys
    CONSTRAINT fk_quota_usage_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

-- Indexes
CREATE UNIQUE INDEX idx_quota_usage_tenant_id_unique ON quota_usage(tenant_id);
CREATE INDEX idx_quota_usage_last_updated ON quota_usage(last_updated);
CREATE INDEX idx_quota_usage_last_reset_daily ON quota_usage(last_reset_daily);
CREATE INDEX idx_quota_usage_last_reset_monthly ON quota_usage(last_reset_monthly);

-- Comments
COMMENT ON TABLE quota_usage IS 'Tabela de uso atual de quotas dos tenants';
COMMENT ON COLUMN quota_usage.id IS 'Identificador único do registro';
COMMENT ON COLUMN quota_usage.tenant_id IS 'ID do tenant';
COMMENT ON COLUMN quota_usage.processes_count IS 'Número atual de processos';
COMMENT ON COLUMN quota_usage.users_count IS 'Número atual de usuários';
COMMENT ON COLUMN quota_usage.clients_count IS 'Número atual de clientes';
COMMENT ON COLUMN quota_usage.datajud_queries_daily IS 'Consultas DataJud realizadas hoje';
COMMENT ON COLUMN quota_usage.datajud_queries_month IS 'Consultas DataJud realizadas no mês';
COMMENT ON COLUMN quota_usage.ai_queries_monthly IS 'Consultas IA realizadas no mês';
COMMENT ON COLUMN quota_usage.storage_used_gb IS 'Armazenamento utilizado em GB';
COMMENT ON COLUMN quota_usage.webhooks_count IS 'Número atual de webhooks';
COMMENT ON COLUMN quota_usage.api_calls_daily IS 'Chamadas API realizadas hoje';
COMMENT ON COLUMN quota_usage.api_calls_monthly IS 'Chamadas API realizadas no mês';
COMMENT ON COLUMN quota_usage.last_updated IS 'Data da última atualização';
COMMENT ON COLUMN quota_usage.last_reset_daily IS 'Data do último reset de contadores diários';
COMMENT ON COLUMN quota_usage.last_reset_monthly IS 'Data do último reset de contadores mensais';