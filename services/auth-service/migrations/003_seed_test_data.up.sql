-- Migration: 003_seed_test_data.sql
-- Cria dados de teste para validação de funcionalidades
-- 2 tenants por plano (8 total) + usuários com diferentes roles

-- ============================================================================
-- PLANOS (Assumindo que já existem da migration do tenant-service)
-- ============================================================================

-- ============================================================================
-- TENANTS DE TESTE (2 por plano = 8 total)
-- ============================================================================

-- STARTER PLAN TENANTS
INSERT INTO tenants (id, name, cnpj, email, phone, address, plan, is_active, created_at, updated_at) VALUES
(
    'tenant-starter-01',
    'Advocacia Silva & Associados',
    '12.345.678/0001-90',
    'contato@silvaassociados.com.br',
    '(11) 98765-4321',
    '{"street": "Rua das Flores", "number": "123", "neighborhood": "Centro", "city": "São Paulo", "state": "SP", "zipCode": "01234-567"}',
    'starter',
    true,
    NOW(),
    NOW()
),
(
    'tenant-starter-02', 
    'Escritório Oliveira Advocacia',
    '23.456.789/0001-01',
    'escritorio@oliveira.adv.br',
    '(11) 91234-5678',
    '{"street": "Av. Paulista", "number": "456", "neighborhood": "Bela Vista", "city": "São Paulo", "state": "SP", "zipCode": "01310-100"}',
    'starter',
    true,
    NOW(),
    NOW()
),

-- PROFESSIONAL PLAN TENANTS
(
    'tenant-prof-01',
    'Costa, Santos & Partners',
    '34.567.890/0001-12', 
    'contato@costasantos.com.br',
    '(11) 92345-6789',
    '{"street": "Rua Augusta", "number": "789", "neighborhood": "Consolação", "city": "São Paulo", "state": "SP", "zipCode": "01305-000"}',
    'professional',
    true,
    NOW(),
    NOW()
),
(
    'tenant-prof-02',
    'Advocacia Pereira & Lima',
    '45.678.901/0001-23',
    'admin@pereiralima.adv.br', 
    '(11) 93456-7890',
    '{"street": "Alameda Santos", "number": "321", "neighborhood": "Paraíso", "city": "São Paulo", "state": "SP", "zipCode": "01419-000"}',
    'professional',
    true,
    NOW(),
    NOW()
),

-- BUSINESS PLAN TENANTS
(
    'tenant-biz-01',
    'Machado Advogados Associados',
    '56.789.012/0001-34',
    'contato@machadoadvogados.com.br',
    '(11) 94567-8901',
    '{"street": "Rua Oscar Freire", "number": "654", "neighborhood": "Jardins", "city": "São Paulo", "state": "SP", "zipCode": "01426-000"}',
    'business',
    true,
    NOW(),
    NOW()
),
(
    'tenant-biz-02',
    'Tribunal Consultoria Jurídica',
    '67.890.123/0001-45',
    'juridico@tribunal.com.br',
    '(11) 95678-9012',
    '{"street": "Av. Faria Lima", "number": "987", "neighborhood": "Itaim Bibi", "city": "São Paulo", "state": "SP", "zipCode": "04538-132"}',
    'business', 
    true,
    NOW(),
    NOW()
),

-- ENTERPRISE PLAN TENANTS
(
    'tenant-ent-01',
    'Barros & Associados Enterprise',
    '78.901.234/0001-56',
    'contato@barrosent.com.br',
    '(11) 96789-0123',
    '{"street": "Av. Brigadeiro Faria Lima", "number": "1111", "neighborhood": "Itaim Bibi", "city": "São Paulo", "state": "SP", "zipCode": "04538-132"}',
    'enterprise',
    true,
    NOW(),
    NOW()
),
(
    'tenant-ent-02',
    'Mega Advocacia Corporate',
    '89.012.345/0001-67',
    'admin@megaadv.com.br',
    '(11) 97890-1234',
    '{"street": "Av. Berrini", "number": "2222", "neighborhood": "Brooklin", "city": "São Paulo", "state": "SP", "zipCode": "04571-000"}',
    'enterprise',
    true,
    NOW(),
    NOW()
);

-- ============================================================================
-- SUBSCRIPTIONS (uma para cada tenant)
-- ============================================================================

INSERT INTO subscriptions (id, tenant_id, plan, status, start_date, trial, quotas, created_at, updated_at) VALUES
-- STARTER SUBSCRIPTIONS
(
    'sub-starter-01',
    'tenant-starter-01', 
    'starter',
    'active',
    NOW(),
    false,
    '{"processes": 50, "users": 2, "mcpCommands": 0, "aiSummaries": 10, "reports": 10, "dashboards": 1, "widgetsPerDashboard": 5, "schedules": 2}',
    NOW(),
    NOW()
),
(
    'sub-starter-02',
    'tenant-starter-02',
    'starter', 
    'trial',
    NOW(),
    true,
    '{"processes": 50, "users": 2, "mcpCommands": 0, "aiSummaries": 10, "reports": 10, "dashboards": 1, "widgetsPerDashboard": 5, "schedules": 2}',
    NOW(),
    NOW()
),

-- PROFESSIONAL SUBSCRIPTIONS  
(
    'sub-prof-01',
    'tenant-prof-01',
    'professional',
    'active', 
    NOW(),
    false,
    '{"processes": 200, "users": 5, "mcpCommands": 200, "aiSummaries": 50, "reports": 100, "dashboards": 3, "widgetsPerDashboard": 10, "schedules": 10}',
    NOW(),
    NOW()
),
(
    'sub-prof-02',
    'tenant-prof-02',
    'professional',
    'active',
    NOW(),
    false, 
    '{"processes": 200, "users": 5, "mcpCommands": 200, "aiSummaries": 50, "reports": 100, "dashboards": 3, "widgetsPerDashboard": 10, "schedules": 10}',
    NOW(),
    NOW()
),

-- BUSINESS SUBSCRIPTIONS
(
    'sub-biz-01',
    'tenant-biz-01',
    'business',
    'active',
    NOW(),
    false,
    '{"processes": 500, "users": 15, "mcpCommands": 1000, "aiSummaries": 200, "reports": 500, "dashboards": 10, "widgetsPerDashboard": 20, "schedules": 50}',
    NOW(),
    NOW()
),
(
    'sub-biz-02',
    'tenant-biz-02',
    'business',
    'active',
    NOW(),
    false,
    '{"processes": 500, "users": 15, "mcpCommands": 1000, "aiSummaries": 200, "reports": 500, "dashboards": 10, "widgetsPerDashboard": 20, "schedules": 50}',
    NOW(),
    NOW()
),

-- ENTERPRISE SUBSCRIPTIONS
(
    'sub-ent-01',
    'tenant-ent-01',
    'enterprise',
    'active',
    NOW(),
    false, 
    '{"processes": -1, "users": -1, "mcpCommands": -1, "aiSummaries": -1, "reports": -1, "dashboards": -1, "widgetsPerDashboard": -1, "schedules": -1}',
    NOW(),
    NOW()
),
(
    'sub-ent-02',
    'tenant-ent-02',
    'enterprise',
    'active',
    NOW(),
    false,
    '{"processes": -1, "users": -1, "mcpCommands": -1, "aiSummaries": -1, "reports": -1, "dashboards": -1, "widgetsPerDashboard": -1, "schedules": -1}',
    NOW(),
    NOW()
);

-- ============================================================================
-- USUÁRIOS DE TESTE (4 roles por tenant = 32 total)
-- ============================================================================

-- HASH da senha "123456" para todos os usuários de teste
-- bcrypt hash: $2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG

-- STARTER TENANT 01 USERS
INSERT INTO users (id, email, name, role, tenant_id, password_hash, is_active, created_at, updated_at) VALUES
-- Admin - acesso total ao tenant
('user-starter-01-admin', 'admin@silvaassociados.com.br', 'Carlos Silva', 'admin', 'tenant-starter-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
-- Manager - gerenciamento de usuários e processos
('user-starter-01-manager', 'gerente@silvaassociados.com.br', 'Ana Silva', 'manager', 'tenant-starter-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
-- Lawyer - advocado com acesso a processos e clientes  
('user-starter-01-lawyer', 'advogado@silvaassociados.com.br', 'Ricardo Silva', 'lawyer', 'tenant-starter-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
-- Assistant - acesso limitado, operações básicas
('user-starter-01-assistant', 'assistente@silvaassociados.com.br', 'Mariana Silva', 'assistant', 'tenant-starter-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- STARTER TENANT 02 USERS  
('user-starter-02-admin', 'admin@oliveira.adv.br', 'João Oliveira', 'admin', 'tenant-starter-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-starter-02-manager', 'gerente@oliveira.adv.br', 'Paula Oliveira', 'manager', 'tenant-starter-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-starter-02-lawyer', 'advogado@oliveira.adv.br', 'Marco Oliveira', 'lawyer', 'tenant-starter-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-starter-02-assistant', 'assistente@oliveira.adv.br', 'Lucia Oliveira', 'assistant', 'tenant-starter-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- PROFESSIONAL TENANT 01 USERS
('user-prof-01-admin', 'admin@costasantos.com.br', 'Roberto Costa', 'admin', 'tenant-prof-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-prof-01-manager', 'gerente@costasantos.com.br', 'Sandra Santos', 'manager', 'tenant-prof-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-prof-01-lawyer', 'advogado@costasantos.com.br', 'Felipe Costa', 'lawyer', 'tenant-prof-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-prof-01-assistant', 'assistente@costasantos.com.br', 'Carla Santos', 'assistant', 'tenant-prof-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- PROFESSIONAL TENANT 02 USERS
('user-prof-02-admin', 'admin@pereiralima.adv.br', 'Eduardo Pereira', 'admin', 'tenant-prof-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-prof-02-manager', 'gerente@pereiralima.adv.br', 'Fernanda Lima', 'manager', 'tenant-prof-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-prof-02-lawyer', 'advogado@pereiralima.adv.br', 'Bruno Pereira', 'lawyer', 'tenant-prof-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-prof-02-assistant', 'assistente@pereiralima.adv.br', 'Isabela Lima', 'assistant', 'tenant-prof-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- BUSINESS TENANT 01 USERS
('user-biz-01-admin', 'admin@machadoadvogados.com.br', 'Luiz Machado', 'admin', 'tenant-biz-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-biz-01-manager', 'gerente@machadoadvogados.com.br', 'Patricia Machado', 'manager', 'tenant-biz-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-biz-01-lawyer', 'advogado@machadoadvogados.com.br', 'André Machado', 'lawyer', 'tenant-biz-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-biz-01-assistant', 'assistente@machadoadvogados.com.br', 'Camila Machado', 'assistant', 'tenant-biz-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- BUSINESS TENANT 02 USERS
('user-biz-02-admin', 'admin@tribunal.com.br', 'Sergio Tribunal', 'admin', 'tenant-biz-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-biz-02-manager', 'gerente@tribunal.com.br', 'Renata Tribunal', 'manager', 'tenant-biz-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-biz-02-lawyer', 'advogado@tribunal.com.br', 'Gustavo Tribunal', 'lawyer', 'tenant-biz-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-biz-02-assistant', 'assistente@tribunal.com.br', 'Juliana Tribunal', 'assistant', 'tenant-biz-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- ENTERPRISE TENANT 01 USERS
('user-ent-01-admin', 'admin@barrosent.com.br', 'Alexandre Barros', 'admin', 'tenant-ent-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-ent-01-manager', 'gerente@barrosent.com.br', 'Claudia Barros', 'manager', 'tenant-ent-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-ent-01-lawyer', 'advogado@barrosent.com.br', 'Rodrigo Barros', 'lawyer', 'tenant-ent-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-ent-01-assistant', 'assistente@barrosent.com.br', 'Vanessa Barros', 'assistant', 'tenant-ent-01', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),

-- ENTERPRISE TENANT 02 USERS  
('user-ent-02-admin', 'admin@megaadv.com.br', 'Marcelo Mega', 'admin', 'tenant-ent-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-ent-02-manager', 'gerente@megaadv.com.br', 'Simone Mega', 'manager', 'tenant-ent-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-ent-02-lawyer', 'advogado@megaadv.com.br', 'Thiago Mega', 'lawyer', 'tenant-ent-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW()),
('user-ent-02-assistant', 'assistente@megaadv.com.br', 'Leticia Mega', 'assistant', 'tenant-ent-02', '$2a$10$N9qo8uLOickgx2ZMRZoMye0IXW4fFkXpEklQzb9VQYm8wUaQ8YTYG', true, NOW(), NOW());

-- ============================================================================
-- SESSÕES ATIVAS PARA ALGUNS USUÁRIOS (para teste de auth)
-- ============================================================================

INSERT INTO sessions (id, user_id, token, refresh_token, expires_at, created_at, updated_at) VALUES
-- Sessões ativas para um usuário de cada plano
('session-starter-admin', 'user-starter-01-admin', 'jwt-token-starter-admin-active', 'refresh-token-starter-admin', NOW() + INTERVAL '24 hours', NOW(), NOW()),
('session-prof-admin', 'user-prof-01-admin', 'jwt-token-prof-admin-active', 'refresh-token-prof-admin', NOW() + INTERVAL '24 hours', NOW(), NOW()),
('session-biz-admin', 'user-biz-01-admin', 'jwt-token-biz-admin-active', 'refresh-token-biz-admin', NOW() + INTERVAL '24 hours', NOW(), NOW()),
('session-ent-admin', 'user-ent-01-admin', 'jwt-token-ent-admin-active', 'refresh-token-ent-admin', NOW() + INTERVAL '24 hours', NOW(), NOW());

-- ============================================================================
-- COMENTÁRIOS DE TESTE
-- ============================================================================

-- CREDENCIAIS DE TESTE:
-- Todos os usuários têm senha: 123456
-- 
-- TENANTS POR PLANO:
-- Starter: tenant-starter-01, tenant-starter-02
-- Professional: tenant-prof-01, tenant-prof-02  
-- Business: tenant-biz-01, tenant-biz-02
-- Enterprise: tenant-ent-01, tenant-ent-02
--
-- ROLES POR TENANT:
-- admin - Acesso total ao tenant
-- manager - Gerenciamento de usuários e processos
-- lawyer - Processos, clientes, relatórios
-- assistant - Acesso limitado, operações básicas
--
-- QUOTAS TESTÁVEIS:
-- Starter: 50 processos, 2 usuários, sem MCP
-- Professional: 200 processos, 5 usuários, 200 MCP commands
-- Business: 500 processos, 15 usuários, 1000 MCP commands  
-- Enterprise: Ilimitado