-- Initialization script for Direito Lux Development Environment
-- Creates databases and initial schemas for all services

-- Database creation is handled by environment variables, 
-- but we'll ensure proper schemas exist

-- Create schemas for different services
CREATE SCHEMA IF NOT EXISTS auth_service;
CREATE SCHEMA IF NOT EXISTS process_service;
CREATE SCHEMA IF NOT EXISTS tenant_service;
CREATE SCHEMA IF NOT EXISTS datajud_service;
CREATE SCHEMA IF NOT EXISTS notification_service;
CREATE SCHEMA IF NOT EXISTS ai_service;
CREATE SCHEMA IF NOT EXISTS search_service;

-- Create common extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create common types and functions
DO $$
BEGIN
    -- Create common audit columns type
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'audit_info') THEN
        CREATE TYPE audit_info AS (
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            created_by UUID,
            updated_by UUID
        );
    END IF;
END$$;

-- Common audit function
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- AI Service specific tables
CREATE TABLE IF NOT EXISTS ai_service.analysis_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    original_text TEXT NOT NULL,
    analysis_type VARCHAR(50) NOT NULL,
    result JSONB NOT NULL,
    confidence_score DECIMAL(5,4),
    processing_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_service.document_embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    document_id UUID NOT NULL,
    chunk_index INTEGER NOT NULL,
    text_chunk TEXT NOT NULL,
    embedding VECTOR(1536), -- OpenAI embedding dimension
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, document_id, chunk_index)
);

-- Search Service specific tables
CREATE TABLE IF NOT EXISTS search_service.search_indices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    index_name VARCHAR(255) NOT NULL,
    index_type VARCHAR(50) NOT NULL,
    settings JSONB DEFAULT '{}',
    mappings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, index_name)
);

CREATE TABLE IF NOT EXISTS search_service.search_queries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    query_text TEXT NOT NULL,
    filters JSONB DEFAULT '{}',
    results_count INTEGER DEFAULT 0,
    response_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_analysis_results_tenant_user ON ai_service.analysis_results(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_analysis_results_type ON ai_service.analysis_results(analysis_type);
CREATE INDEX IF NOT EXISTS idx_analysis_results_created_at ON ai_service.analysis_results(created_at);

CREATE INDEX IF NOT EXISTS idx_document_embeddings_tenant_doc ON ai_service.document_embeddings(tenant_id, document_id);
CREATE INDEX IF NOT EXISTS idx_document_embeddings_created_at ON ai_service.document_embeddings(created_at);

CREATE INDEX IF NOT EXISTS idx_search_indices_tenant ON search_service.search_indices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_search_indices_active ON search_service.search_indices(is_active);

CREATE INDEX IF NOT EXISTS idx_search_queries_tenant ON search_service.search_queries(tenant_id);
CREATE INDEX IF NOT EXISTS idx_search_queries_created_at ON search_service.search_queries(created_at);

-- Create triggers for updated_at
CREATE TRIGGER set_timestamp_analysis_results
    BEFORE UPDATE ON ai_service.analysis_results
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_search_indices
    BEFORE UPDATE ON search_service.search_indices
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Insert sample data for development
INSERT INTO ai_service.analysis_results (
    tenant_id, 
    user_id, 
    document_type, 
    original_text, 
    analysis_type, 
    result, 
    confidence_score
) VALUES 
(
    '11111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'contract',
    'Este é um contrato de prestação de serviços jurídicos.',
    'document_classification',
    '{"category": "contract", "subcategory": "service_agreement", "keywords": ["prestação", "serviços", "jurídicos"]}',
    0.95
),
(
    '11111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'petition',
    'Petição inicial para ação de cobrança.',
    'document_classification',
    '{"category": "petition", "subcategory": "initial_petition", "legal_area": "civil", "action_type": "cobrança"}',
    0.89
)
ON CONFLICT DO NOTHING;

INSERT INTO search_service.search_indices (
    tenant_id,
    index_name,
    index_type,
    settings,
    mappings
) VALUES 
(
    '11111111-1111-1111-1111-111111111111',
    'processes_index',
    'elasticsearch',
    '{"number_of_shards": 1, "number_of_replicas": 0}',
    '{"properties": {"title": {"type": "text"}, "content": {"type": "text"}, "date": {"type": "date"}}}'
),
(
    '11111111-1111-1111-1111-111111111111',
    'documents_index',
    'elasticsearch',
    '{"number_of_shards": 1, "number_of_replicas": 0}',
    '{"properties": {"filename": {"type": "keyword"}, "content": {"type": "text"}, "metadata": {"type": "object"}}}'
)
ON CONFLICT (tenant_id, index_name) DO NOTHING;

-- Grant permissions (simplified for development)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA ai_service TO direito_lux;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA ai_service TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA search_service TO direito_lux;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA search_service TO direito_lux;

-- Create development users/tenants if they don't exist
CREATE TABLE IF NOT EXISTS public.tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(50) NOT NULL DEFAULT 'basic',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES public.tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

-- Insert development tenant and user
INSERT INTO public.tenants (id, name, plan) VALUES 
    ('11111111-1111-1111-1111-111111111111', 'Tenant Dev', 'premium'),
    ('22222222-2222-2222-2222-222222222222', 'Tenant Test', 'basic')
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.users (id, tenant_id, email, name, role) VALUES 
    ('11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'dev@direito-lux.com', 'Dev User', 'admin'),
    ('22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'test@direito-lux.com', 'Test User', 'user')
ON CONFLICT (tenant_id, email) DO NOTHING;