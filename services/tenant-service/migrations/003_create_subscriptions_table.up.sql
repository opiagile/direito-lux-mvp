-- Migration: Create subscriptions table
-- Description: Tabela para armazenar as assinaturas dos tenants

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    plan_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'trialing',
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT false,
    trial_start TIMESTAMP WITH TIME ZONE,
    trial_end TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    canceled_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT subscriptions_status_check CHECK (status IN ('active', 'trialing', 'past_due', 'canceled', 'unpaid')),
    CONSTRAINT subscriptions_period_valid CHECK (current_period_end > current_period_start),
    CONSTRAINT subscriptions_trial_valid CHECK (
        (trial_start IS NULL AND trial_end IS NULL) OR 
        (trial_start IS NOT NULL AND trial_end IS NOT NULL AND trial_end > trial_start)
    ),

    -- Foreign keys
    CONSTRAINT fk_subscriptions_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscriptions_plan FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE RESTRICT
);

-- Indexes
CREATE INDEX idx_subscriptions_tenant_id ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_current_period_end ON subscriptions(current_period_end);
CREATE INDEX idx_subscriptions_trial_end ON subscriptions(trial_end) WHERE trial_end IS NOT NULL;
CREATE INDEX idx_subscriptions_created_at ON subscriptions(created_at);

-- Unique constraints
CREATE UNIQUE INDEX idx_subscriptions_tenant_active_unique ON subscriptions(tenant_id) 
WHERE status IN ('active', 'trialing', 'past_due');

-- Comments
COMMENT ON TABLE subscriptions IS 'Tabela de assinaturas dos tenants';
COMMENT ON COLUMN subscriptions.id IS 'Identificador único da assinatura';
COMMENT ON COLUMN subscriptions.tenant_id IS 'ID do tenant proprietário da assinatura';
COMMENT ON COLUMN subscriptions.plan_id IS 'ID do plano contratado';
COMMENT ON COLUMN subscriptions.status IS 'Status atual da assinatura';
COMMENT ON COLUMN subscriptions.current_period_start IS 'Início do período atual de cobrança';
COMMENT ON COLUMN subscriptions.current_period_end IS 'Fim do período atual de cobrança';
COMMENT ON COLUMN subscriptions.cancel_at_period_end IS 'Indica se a assinatura será cancelada no final do período';
COMMENT ON COLUMN subscriptions.trial_start IS 'Início do período de teste';
COMMENT ON COLUMN subscriptions.trial_end IS 'Fim do período de teste';
COMMENT ON COLUMN subscriptions.canceled_at IS 'Data de cancelamento da assinatura';