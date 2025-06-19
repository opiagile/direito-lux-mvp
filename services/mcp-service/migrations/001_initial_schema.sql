-- Initial schema for MCP Service
-- PostgreSQL migration script

-- Extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tenants table (reference)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(50) NOT NULL DEFAULT 'basic',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Users table (reference)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

-- MCP Sessions table
CREATE TABLE mcp_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL, -- whatsapp, telegram, claude, slack
    external_id VARCHAR(255) NOT NULL, -- ID externo do chat
    state VARCHAR(20) NOT NULL DEFAULT 'active',
    context JSONB DEFAULT '{}',
    last_interaction TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    message_count INTEGER DEFAULT 0,
    command_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + INTERVAL '30 minutes'),
    
    -- Indexes para performance
    INDEX idx_mcp_sessions_tenant_user (tenant_id, user_id),
    INDEX idx_mcp_sessions_channel (channel),
    INDEX idx_mcp_sessions_external_id (external_id),
    INDEX idx_mcp_sessions_state (state),
    INDEX idx_mcp_sessions_expires_at (expires_at)
);

-- MCP Messages table (histórico de conversas)
CREATE TABLE mcp_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES mcp_sessions(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- user, assistant, system
    content TEXT NOT NULL,
    tokens_used INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Indexes para busca de histórico
    INDEX idx_mcp_messages_session (session_id),
    INDEX idx_mcp_messages_tenant_user (tenant_id, user_id),
    INDEX idx_mcp_messages_created_at (created_at)
);

-- MCP Tools Executions table (log de execução de ferramentas)
CREATE TABLE mcp_tool_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES mcp_sessions(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tool_name VARCHAR(100) NOT NULL,
    parameters JSONB DEFAULT '{}',
    result JSONB DEFAULT '{}',
    success BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT,
    duration_ms INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Indexes para analytics
    INDEX idx_mcp_tools_session (session_id),
    INDEX idx_mcp_tools_tenant (tenant_id),
    INDEX idx_mcp_tools_name (tool_name),
    INDEX idx_mcp_tools_success (success),
    INDEX idx_mcp_tools_created_at (created_at)
);

-- MCP Events table (eventos de domínio)
CREATE TABLE mcp_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    data JSONB NOT NULL DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    
    -- Indexes para event sourcing
    INDEX idx_mcp_events_type (type),
    INDEX idx_mcp_events_aggregate (aggregate_id),
    INDEX idx_mcp_events_tenant (tenant_id),
    INDEX idx_mcp_events_occurred_at (occurred_at),
    INDEX idx_mcp_events_processed (processed_at)
);

-- MCP Quotas table (controle de uso por tenant)
CREATE TABLE mcp_quotas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    quota_type VARCHAR(50) NOT NULL, -- tokens, requests, sessions
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    limit_value INTEGER NOT NULL,
    used_value INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraint: um quota por tipo por período por tenant
    UNIQUE(tenant_id, quota_type, period_start),
    
    -- Indexes
    INDEX idx_mcp_quotas_tenant (tenant_id),
    INDEX idx_mcp_quotas_type (quota_type),
    INDEX idx_mcp_quotas_period (period_start, period_end)
);

-- Bot Configurations table (configurações por canal)
CREATE TABLE mcp_bot_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL, -- whatsapp, telegram, slack
    config JSONB NOT NULL DEFAULT '{}', -- configurações específicas do bot
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraint: uma config por canal por tenant
    UNIQUE(tenant_id, channel),
    
    -- Index
    INDEX idx_mcp_bot_configs_tenant (tenant_id),
    INDEX idx_mcp_bot_configs_channel (channel)
);

-- Insert sample tenant and user for development
INSERT INTO tenants (id, name, plan) VALUES 
    ('11111111-1111-1111-1111-111111111111', 'Tenant Dev', 'premium'),
    ('22222222-2222-2222-2222-222222222222', 'Tenant Test', 'basic')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, tenant_id, email, name) VALUES 
    ('11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'dev@direito-lux.com', 'Dev User'),
    ('22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'test@direito-lux.com', 'Test User')
ON CONFLICT (tenant_id, email) DO NOTHING;

-- Insert default quota configurations
INSERT INTO mcp_quotas (tenant_id, quota_type, period_start, period_end, limit_value) VALUES 
    ('11111111-1111-1111-1111-111111111111', 'tokens', date_trunc('month', NOW()), date_trunc('month', NOW()) + INTERVAL '1 month', 100000),
    ('11111111-1111-1111-1111-111111111111', 'requests', date_trunc('month', NOW()), date_trunc('month', NOW()) + INTERVAL '1 month', 1000),
    ('11111111-1111-1111-1111-111111111111', 'sessions', date_trunc('month', NOW()), date_trunc('month', NOW()) + INTERVAL '1 month', 50),
    ('22222222-2222-2222-2222-222222222222', 'tokens', date_trunc('month', NOW()), date_trunc('month', NOW()) + INTERVAL '1 month', 10000),
    ('22222222-2222-2222-2222-222222222222', 'requests', date_trunc('month', NOW()), date_trunc('month', NOW()) + INTERVAL '1 month', 100),
    ('22222222-2222-2222-2222-222222222222', 'sessions', date_trunc('month', NOW()), date_trunc('month', NOW()) + INTERVAL '1 month', 10)
ON CONFLICT (tenant_id, quota_type, period_start) DO NOTHING;