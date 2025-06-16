-- Migration: Create tenants table
-- Description: Tabela principal para armazenar informações dos tenants/organizações

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    document VARCHAR(20), -- CNPJ
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    website VARCHAR(255),
    address JSONB DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    plan_type VARCHAR(20) NOT NULL DEFAULT 'starter',
    owner_user_id UUID NOT NULL,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMP WITH TIME ZONE,
    suspended_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT tenants_status_check CHECK (status IN ('pending', 'active', 'suspended', 'canceled', 'blocked')),
    CONSTRAINT tenants_plan_type_check CHECK (plan_type IN ('starter', 'professional', 'business', 'enterprise')),
    CONSTRAINT tenants_email_valid CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT tenants_name_length CHECK (LENGTH(name) >= 3 AND LENGTH(name) <= 50)
);

-- Indexes
CREATE INDEX idx_tenants_email ON tenants(email);
CREATE INDEX idx_tenants_document ON tenants(document) WHERE document IS NOT NULL;
CREATE INDEX idx_tenants_owner_user_id ON tenants(owner_user_id);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_plan_type ON tenants(plan_type);
CREATE INDEX idx_tenants_created_at ON tenants(created_at);

-- Unique constraints
CREATE UNIQUE INDEX idx_tenants_email_unique ON tenants(email);
CREATE UNIQUE INDEX idx_tenants_document_unique ON tenants(document) WHERE document IS NOT NULL;

-- Comments
COMMENT ON TABLE tenants IS 'Tabela principal de tenants/organizações do sistema';
COMMENT ON COLUMN tenants.id IS 'Identificador único do tenant';
COMMENT ON COLUMN tenants.name IS 'Nome fantasia do tenant';
COMMENT ON COLUMN tenants.legal_name IS 'Razão social do tenant';
COMMENT ON COLUMN tenants.document IS 'CNPJ do tenant';
COMMENT ON COLUMN tenants.email IS 'Email principal do tenant';
COMMENT ON COLUMN tenants.address IS 'Endereço do tenant em formato JSON';
COMMENT ON COLUMN tenants.status IS 'Status atual do tenant';
COMMENT ON COLUMN tenants.plan_type IS 'Tipo de plano contratado';
COMMENT ON COLUMN tenants.owner_user_id IS 'ID do usuário proprietário do tenant';
COMMENT ON COLUMN tenants.settings IS 'Configurações específicas do tenant em formato JSON';