-- Migration: Create plans table
-- Description: Tabela para armazenar os planos de assinatura disponíveis

CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    description TEXT,
    price BIGINT NOT NULL DEFAULT 0, -- Preço em centavos
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    billing_interval VARCHAR(20) NOT NULL DEFAULT 'monthly',
    features JSONB NOT NULL DEFAULT '{}',
    quotas JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT plans_type_check CHECK (type IN ('starter', 'professional', 'business', 'enterprise')),
    CONSTRAINT plans_billing_interval_check CHECK (billing_interval IN ('monthly', 'yearly')),
    CONSTRAINT plans_currency_check CHECK (currency IN ('BRL', 'USD', 'EUR')),
    CONSTRAINT plans_price_positive CHECK (price >= 0),
    CONSTRAINT plans_name_length CHECK (LENGTH(name) >= 3 AND LENGTH(name) <= 100)
);

-- Indexes
CREATE INDEX idx_plans_type ON plans(type);
CREATE INDEX idx_plans_is_active ON plans(is_active);
CREATE INDEX idx_plans_billing_interval ON plans(billing_interval);
CREATE INDEX idx_plans_price ON plans(price);
CREATE INDEX idx_plans_created_at ON plans(created_at);

-- Unique constraints
CREATE UNIQUE INDEX idx_plans_type_active_unique ON plans(type) WHERE is_active = true;

-- Comments
COMMENT ON TABLE plans IS 'Tabela de planos de assinatura disponíveis';
COMMENT ON COLUMN plans.id IS 'Identificador único do plano';
COMMENT ON COLUMN plans.name IS 'Nome do plano';
COMMENT ON COLUMN plans.type IS 'Tipo do plano (starter, professional, business, enterprise)';
COMMENT ON COLUMN plans.description IS 'Descrição detalhada do plano';
COMMENT ON COLUMN plans.price IS 'Preço do plano em centavos';
COMMENT ON COLUMN plans.currency IS 'Moeda do preço';
COMMENT ON COLUMN plans.billing_interval IS 'Intervalo de cobrança (monthly, yearly)';
COMMENT ON COLUMN plans.features IS 'Funcionalidades disponíveis no plano em formato JSON';
COMMENT ON COLUMN plans.quotas IS 'Limites e quotas do plano em formato JSON';
COMMENT ON COLUMN plans.is_active IS 'Indica se o plano está ativo para novos assinantes';

-- Insert default plans
INSERT INTO plans (id, name, type, description, price, features, quotas) VALUES
(
    uuid_generate_v4(),
    'Starter',
    'starter',
    'Plano ideal para escritórios pequenos que estão começando com automação jurídica.',
    9900,
    '{
        "whatsapp_enabled": true,
        "ai_enabled": false,
        "advanced_ai": false,
        "jurisprudence_enabled": false,
        "white_label_enabled": false,
        "custom_integrations": false,
        "priority_support": false,
        "custom_reports": false,
        "api_access": false,
        "webhooks_enabled": false
    }',
    '{
        "max_processes": 50,
        "max_users": 2,
        "max_clients": 20,
        "datajud_queries_daily": 100,
        "ai_queries_monthly": 10,
        "storage_gb": 1,
        "max_webhooks": 3,
        "max_api_calls_daily": 1000
    }'
),
(
    uuid_generate_v4(),
    'Professional',
    'professional',
    'Plano para escritórios em crescimento que precisam de mais recursos e automação.',
    29900,
    '{
        "whatsapp_enabled": true,
        "ai_enabled": true,
        "advanced_ai": false,
        "jurisprudence_enabled": false,
        "white_label_enabled": false,
        "custom_integrations": false,
        "priority_support": false,
        "custom_reports": true,
        "api_access": true,
        "webhooks_enabled": true
    }',
    '{
        "max_processes": 200,
        "max_users": 5,
        "max_clients": 100,
        "datajud_queries_daily": 500,
        "ai_queries_monthly": 50,
        "storage_gb": 5,
        "max_webhooks": 10,
        "max_api_calls_daily": 5000
    }'
),
(
    uuid_generate_v4(),
    'Business',
    'business',
    'Plano completo para escritórios médios com necessidades avançadas de automação.',
    69900,
    '{
        "whatsapp_enabled": true,
        "ai_enabled": true,
        "advanced_ai": true,
        "jurisprudence_enabled": true,
        "white_label_enabled": false,
        "custom_integrations": true,
        "priority_support": true,
        "custom_reports": true,
        "api_access": true,
        "webhooks_enabled": true
    }',
    '{
        "max_processes": 500,
        "max_users": 15,
        "max_clients": 500,
        "datajud_queries_daily": 2000,
        "ai_queries_monthly": 200,
        "storage_gb": 20,
        "max_webhooks": 25,
        "max_api_calls_daily": 15000
    }'
),
(
    uuid_generate_v4(),
    'Enterprise',
    'enterprise',
    'Plano empresarial com recursos ilimitados e suporte personalizado.',
    199900,
    '{
        "whatsapp_enabled": true,
        "ai_enabled": true,
        "advanced_ai": true,
        "jurisprudence_enabled": true,
        "white_label_enabled": true,
        "custom_integrations": true,
        "priority_support": true,
        "custom_reports": true,
        "api_access": true,
        "webhooks_enabled": true
    }',
    '{
        "max_processes": -1,
        "max_users": -1,
        "max_clients": -1,
        "datajud_queries_daily": 10000,
        "ai_queries_monthly": -1,
        "storage_gb": 100,
        "max_webhooks": -1,
        "max_api_calls_daily": -1
    }'
);