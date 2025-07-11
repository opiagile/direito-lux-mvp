-- Migration: Create billing tables
-- Created: 2025-07-11

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    price_monthly BIGINT NOT NULL DEFAULT 0, -- em centavos
    price_yearly BIGINT NOT NULL DEFAULT 0,  -- em centavos
    trial_days INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    
    -- Limites do plano
    max_processes INTEGER NOT NULL DEFAULT 0,
    max_users INTEGER NOT NULL DEFAULT 0,
    max_ai_requests INTEGER NOT NULL DEFAULT 0,
    max_bot_commands INTEGER NOT NULL DEFAULT 0,
    
    -- Funcionalidades (flags booleanas)
    has_whatsapp BOOLEAN NOT NULL DEFAULT false,
    has_telegram BOOLEAN NOT NULL DEFAULT false,
    has_mcp_bot BOOLEAN NOT NULL DEFAULT false,
    has_unlimited_search BOOLEAN NOT NULL DEFAULT false,
    has_ai_analysis BOOLEAN NOT NULL DEFAULT false,
    has_predictions BOOLEAN NOT NULL DEFAULT false,
    has_doc_generation BOOLEAN NOT NULL DEFAULT false,
    has_white_label BOOLEAN NOT NULL DEFAULT false,
    has_api_access BOOLEAN NOT NULL DEFAULT false,
    has_priority_support BOOLEAN NOT NULL DEFAULT false,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    
    -- Dados pessoais
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    document VARCHAR(20) NOT NULL,
    document_type VARCHAR(10) NOT NULL CHECK (document_type IN ('cpf', 'cnpj')),
    
    -- Dados da empresa (para CNPJ)
    company_name VARCHAR(255),
    trading_name VARCHAR(255),
    state_registration VARCHAR(50),
    
    -- Endereço (JSON)
    address JSONB,
    
    -- Dados de cobrança
    billing_email VARCHAR(255),
    billing_phone VARCHAR(20),
    
    -- Integrações externas
    asaas_customer_id VARCHAR(100),
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'blocked')),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES plans(id),
    status VARCHAR(30) NOT NULL DEFAULT 'trial' CHECK (status IN ('trial', 'active', 'past_due', 'suspended', 'cancelled', 'expired', 'payment_pending')),
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly' CHECK (billing_cycle IN ('monthly', 'yearly')),
    
    -- Dados do período
    trial_start_date TIMESTAMP,
    trial_end_date TIMESTAMP,
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    
    -- Dados de pagamento
    amount BIGINT NOT NULL DEFAULT 0, -- em centavos
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('credit_card', 'debit_card', 'pix', 'boleto', 'bitcoin', 'xrp', 'xlm', 'xdc', 'cardano', 'hbar', 'xcn', 'ethereum', 'solana')),
    
    -- Integrações externas
    asaas_subscription_id VARCHAR(100),
    asaas_customer_id VARCHAR(100),
    
    -- Controle de cobrança
    next_billing_date TIMESTAMP,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_payment_attempt TIMESTAMP,
    last_successful_payment TIMESTAMP,
    
    -- Cancelamento
    cancelled_at TIMESTAMP,
    cancel_reason TEXT,
    cancelled_by UUID,
    
    -- Controle de versão
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),
    tenant_id UUID NOT NULL,
    invoice_id UUID,
    
    -- Dados do pagamento
    amount BIGINT NOT NULL DEFAULT 0, -- em centavos
    currency VARCHAR(10) NOT NULL DEFAULT 'BRL',
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('credit_card', 'debit_card', 'pix', 'boleto', 'bitcoin', 'xrp', 'xlm', 'xdc', 'cardano', 'hbar', 'xcn', 'ethereum', 'solana')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed', 'cancelled', 'refunded', 'expired', 'partial')),
    
    -- Integrações externas
    asaas_payment_id VARCHAR(100),
    asaas_charge_id VARCHAR(100),
    now_payment_id VARCHAR(100),
    
    -- Dados específicos do pagamento
    gateway_response TEXT,
    gateway_reference VARCHAR(100),
    transaction_id VARCHAR(100),
    
    -- Dados de cobrança
    due_date TIMESTAMP,
    paid_at TIMESTAMP,
    
    -- Controle de retry
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_attempt_at TIMESTAMP,
    next_retry_at TIMESTAMP,
    
    -- Cancelamento/Estorno
    cancelled_at TIMESTAMP,
    cancel_reason TEXT,
    refunded_at TIMESTAMP,
    refund_amount BIGINT,
    
    -- Cripto específico
    crypto_address VARCHAR(255),
    crypto_amount VARCHAR(50),
    crypto_tx_hash VARCHAR(255),
    exchange_rate DECIMAL(18, 8),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),
    tenant_id UUID NOT NULL,
    payment_id UUID,
    
    -- Dados da fatura
    number VARCHAR(50) NOT NULL UNIQUE,
    amount BIGINT NOT NULL DEFAULT 0,        -- em centavos
    tax_amount BIGINT NOT NULL DEFAULT 0,    -- ISS, etc.
    discount_amount BIGINT NOT NULL DEFAULT 0,
    total_amount BIGINT NOT NULL DEFAULT 0,
    
    -- Período da fatura
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    
    -- Dados de cobrança
    due_date TIMESTAMP NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    paid_at TIMESTAMP,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'issued', 'sent', 'paid', 'overdue', 'cancelled', 'refunded')),
    
    -- Nota fiscal eletrônica
    nfe_number VARCHAR(50),
    nfe_key VARCHAR(100),
    nfe_url TEXT,
    nfe_status VARCHAR(20),
    nfe_issued_at TIMESTAMP,
    
    -- Dados do cliente (para NF-e)
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_document VARCHAR(20) NOT NULL,
    customer_phone VARCHAR(20),
    customer_address JSONB,
    
    -- Dados da empresa (Curitiba)
    company_name VARCHAR(255) NOT NULL,
    company_document VARCHAR(20) NOT NULL,
    company_address JSONB,
    
    -- Integrações externas
    asaas_invoice_id VARCHAR(100),
    
    -- Controle
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_retry_at TIMESTAMP,
    error_message TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_next_billing_date ON subscriptions(next_billing_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_asaas_id ON subscriptions(asaas_subscription_id);

CREATE INDEX IF NOT EXISTS idx_payments_subscription_id ON payments(subscription_id);
CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_due_date ON payments(due_date);
CREATE INDEX IF NOT EXISTS idx_payments_asaas_payment_id ON payments(asaas_payment_id);
CREATE INDEX IF NOT EXISTS idx_payments_now_payment_id ON payments(now_payment_id);

CREATE INDEX IF NOT EXISTS idx_invoices_subscription_id ON invoices(subscription_id);
CREATE INDEX IF NOT EXISTS idx_invoices_tenant_id ON invoices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);
CREATE INDEX IF NOT EXISTS idx_invoices_number ON invoices(number);
CREATE INDEX IF NOT EXISTS idx_invoices_nfe_number ON invoices(nfe_number);

CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_document ON customers(document);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_asaas_id ON customers(asaas_customer_id);

CREATE INDEX IF NOT EXISTS idx_plans_name ON plans(name);
CREATE INDEX IF NOT EXISTS idx_plans_active ON plans(active);

-- Create updated_at triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_plans_updated_at BEFORE UPDATE ON plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add foreign key constraints
ALTER TABLE subscriptions ADD CONSTRAINT fk_subscriptions_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
ALTER TABLE customers ADD CONSTRAINT fk_customers_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
ALTER TABLE payments ADD CONSTRAINT fk_payments_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
ALTER TABLE invoices ADD CONSTRAINT fk_invoices_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
ALTER TABLE invoices ADD CONSTRAINT fk_invoices_payment_id FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE SET NULL;