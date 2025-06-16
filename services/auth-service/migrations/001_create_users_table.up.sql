-- =============================================================================
-- Migração: Criar tabela de usuários
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum para roles de usuários
CREATE TYPE user_role AS ENUM (
    'admin',
    'manager', 
    'operator',
    'client',
    'readonly'
);

-- Enum para status de usuários
CREATE TYPE user_status AS ENUM (
    'active',
    'inactive',
    'pending',
    'suspended',
    'blocked'
);

-- Tabela de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'client',
    status user_status NOT NULL DEFAULT 'pending',
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT users_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$'),
    CONSTRAINT users_first_name_check CHECK (length(first_name) >= 1),
    CONSTRAINT users_last_name_check CHECK (length(last_name) >= 1)
);

-- Índices
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_tenant_status ON users(tenant_id, status);

-- Função para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para atualizar updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentários
COMMENT ON TABLE users IS 'Tabela de usuários do sistema';
COMMENT ON COLUMN users.id IS 'Identificador único do usuário';
COMMENT ON COLUMN users.tenant_id IS 'Identificador do tenant';
COMMENT ON COLUMN users.email IS 'Email do usuário (único)';
COMMENT ON COLUMN users.password_hash IS 'Hash da senha do usuário';
COMMENT ON COLUMN users.first_name IS 'Primeiro nome do usuário';
COMMENT ON COLUMN users.last_name IS 'Último nome do usuário';
COMMENT ON COLUMN users.role IS 'Papel do usuário no sistema';
COMMENT ON COLUMN users.status IS 'Status atual do usuário';
COMMENT ON COLUMN users.last_login_at IS 'Data/hora do último login';
COMMENT ON COLUMN users.created_at IS 'Data/hora de criação';
COMMENT ON COLUMN users.updated_at IS 'Data/hora da última atualização';