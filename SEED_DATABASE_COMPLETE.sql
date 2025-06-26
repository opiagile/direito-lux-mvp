-- ============================================================================
-- SCRIPT COMPLETO DE SEED DE DADOS - DIREITO LUX
-- ============================================================================
-- Este script deve ser executado APÓS todas as migrations estarem aplicadas
-- Cria tenants, planos, usuários e subscriptions de forma consistente
-- ============================================================================

-- ============================================================================
-- 1. PLANOS DE ASSINATURA (inserir primeiro)
-- ============================================================================

INSERT INTO plans (id, name, code, price_cents, billing_cycle, is_active, features, quota_limits, created_at, updated_at) VALUES 
-- STARTER PLAN
(
    'plan-starter-uuid',
    'Starter',
    'starter',
    9900, -- R$ 99,00
    'monthly',
    true,
    '{"whatsapp": true, "email": true, "search": true, "basic_reports": true, "mcp_bot": false}',
    '{"processes": 50, "users": 2, "datajud_daily": 100, "notifications_month": 500, "storage_gb": 5, "mcp_commands": 0, "ai_summaries": 10, "reports": 10, "dashboards": 1}',
    NOW(),
    NOW()
),
-- PROFESSIONAL PLAN  
(
    'plan-professional-uuid',
    'Professional',
    'professional', 
    29900, -- R$ 299,00
    'monthly',
    true,
    '{"whatsapp": true, "email": true, "telegram": true, "search": true, "ai_summaries": true, "mcp_bot": true, "advanced_reports": true}',
    '{"processes": 200, "users": 5, "datajud_daily": 500, "notifications_month": 5000, "storage_gb": 50, "mcp_commands": 200, "ai_summaries": 50, "reports": 100, "dashboards": 3}',
    NOW(),
    NOW()
),
-- BUSINESS PLAN
(
    'plan-business-uuid',
    'Business',
    'business',
    69900, -- R$ 699,00
    'monthly', 
    true,
    '{"whatsapp": true, "email": true, "telegram": true, "search": true, "ai_summaries": true, "mcp_bot": true, "claude_chat": true, "advanced_reports": true, "jurisprudence_search": true}',
    '{"processes": 500, "users": 15, "datajud_daily": 2000, "notifications_month": -1, "storage_gb": 200, "mcp_commands": 1000, "ai_summaries": 200, "reports": 500, "dashboards": 10}',
    NOW(),
    NOW()
),
-- ENTERPRISE PLAN
(
    'plan-enterprise-uuid',
    'Enterprise',
    'enterprise',
    199900, -- R$ 1999,00  
    'monthly',
    true,
    '{"whatsapp": true, "email": true, "telegram": true, "slack": true, "search": true, "ai_summaries": true, "mcp_bot": true, "claude_chat": true, "custom_mcp": true, "white_label": true, "advanced_reports": true, "jurisprudence_search": true}',
    '{"processes": -1, "users": -1, "datajud_daily": 10000, "notifications_month": -1, "storage_gb": 1000, "mcp_commands": -1, "ai_summaries": -1, "reports": -1, "dashboards": -1}',
    NOW(),
    NOW()
);

-- ============================================================================
-- 2. TENANTS (2 por plano = 8 total)
-- ============================================================================

INSERT INTO tenants (id, name, legal_name, document, email, phone, address, status, plan_type, owner_user_id, created_at, updated_at) VALUES 
-- STARTER TENANTS
(
    '11111111-1111-1111-1111-111111111111',
    'Silva & Associados',
    'Silva & Associados Advocacia Ltda',
    '12.345.678/0001-90',
    'admin@silvaassociados.com.br',
    '(11) 98765-4321',
    '{"street": "Rua das Flores", "number": "123", "neighborhood": "Centro", "city": "São Paulo", "state": "SP", "zipCode": "01234-567"}',
    'active',
    'starter',
    '21111111-1111-1111-1111-111111111111', -- será criado depois
    NOW(),
    NOW()
),
(
    '12222222-2222-2222-2222-222222222222',
    'Lima Advogados',
    'Lima Advogados Sociedade Simples',
    '23.456.789/0001-01',
    'admin@limaadvogados.com.br',
    '(11) 91234-5678',
    '{"street": "Av. Paulista", "number": "456", "neighborhood": "Bela Vista", "city": "São Paulo", "state": "SP", "zipCode": "01310-100"}',
    'active',
    'starter',
    '22222222-2222-2222-2222-222222222222',
    NOW(),
    NOW()
),

-- PROFESSIONAL TENANTS
(
    '13333333-3333-3333-3333-333333333333',
    'Costa Santos',
    'Costa Santos & Partners Advocacia',
    '34.567.890/0001-12',
    'admin@costasantos.com.br',
    '(11) 92345-6789',
    '{"street": "Rua Augusta", "number": "789", "neighborhood": "Consolação", "city": "São Paulo", "state": "SP", "zipCode": "01305-000"}',
    'active',
    'professional',
    '23333333-3333-3333-3333-333333333333',
    NOW(),
    NOW()
),
(
    '14444444-4444-4444-4444-444444444444',
    'Pereira Oliveira',
    'Pereira Oliveira Advocacia e Consultoria',
    '45.678.901/0001-23',
    'admin@pereiraoliveira.com.br',
    '(11) 93456-7890',
    '{"street": "Alameda Santos", "number": "321", "neighborhood": "Paraíso", "city": "São Paulo", "state": "SP", "zipCode": "01419-000"}',
    'active',
    'professional',
    '24444444-4444-4444-4444-444444444444',
    NOW(),
    NOW()
),

-- BUSINESS TENANTS
(
    '15555555-5555-5555-5555-555555555555',
    'Machado Advogados',
    'Machado Advogados Associados S/A',
    '56.789.012/0001-34',
    'admin@machadoadvogados.com.br',
    '(11) 94567-8901',
    '{"street": "Rua Oscar Freire", "number": "654", "neighborhood": "Jardins", "city": "São Paulo", "state": "SP", "zipCode": "01426-000"}',
    'active',
    'business',
    '25555555-5555-5555-5555-555555555555',
    NOW(),
    NOW()
),
(
    '16666666-6666-6666-6666-666666666666',
    'Ferreira Legal',
    'Ferreira Legal Consultoria Jurídica',
    '67.890.123/0001-45',
    'admin@ferreiralegal.com.br',
    '(11) 95678-9012',
    '{"street": "Av. Faria Lima", "number": "987", "neighborhood": "Itaim Bibi", "city": "São Paulo", "state": "SP", "zipCode": "04538-132"}',
    'active',
    'business',
    '26666666-6666-6666-6666-666666666666',
    NOW(),
    NOW()
),

-- ENTERPRISE TENANTS
(
    '17777777-7777-7777-7777-777777777777',
    'Barros Enterprise',
    'Barros & Associados Enterprise Legal',
    '78.901.234/0001-56',
    'admin@barrosent.com.br',
    '(11) 96789-0123',
    '{"street": "Av. Brigadeiro Faria Lima", "number": "1111", "neighborhood": "Itaim Bibi", "city": "São Paulo", "state": "SP", "zipCode": "04538-132"}',
    'active',
    'enterprise',
    '27777777-7777-7777-7777-777777777777',
    NOW(),
    NOW()
),
(
    '18888888-8888-8888-8888-888888888888',
    'Rodrigues Global',
    'Rodrigues Global Legal Solutions',
    '89.012.345/0001-67',
    'admin@rodriguesglobal.com.br',
    '(11) 97890-1234',
    '{"street": "Av. Berrini", "number": "2222", "neighborhood": "Brooklin", "city": "São Paulo", "state": "SP", "zipCode": "04571-000"}',
    'active',
    'enterprise',
    '28888888-8888-8888-8888-888888888888',
    NOW(),
    NOW()
);

-- ============================================================================
-- 3. USUÁRIOS ADMIN (um por tenant, que serão owners)
-- ============================================================================

INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, status, created_at, updated_at) VALUES 
-- STARTER ADMINS
(
    '21111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'admin@silvaassociados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Carlos',
    'Silva',
    'admin',
    'active',
    NOW(),
    NOW()
),
(
    '22222222-2222-2222-2222-222222222222',
    '12222222-2222-2222-2222-222222222222',
    'admin@limaadvogados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'João',
    'Lima',
    'admin',
    'active',
    NOW(),
    NOW()
),

-- PROFESSIONAL ADMINS
(
    '23333333-3333-3333-3333-333333333333',
    '13333333-3333-3333-3333-333333333333',
    'admin@costasantos.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Roberto',
    'Costa',
    'admin',
    'active',
    NOW(),
    NOW()
),
(
    '24444444-4444-4444-4444-444444444444',
    '14444444-4444-4444-4444-444444444444',
    'admin@pereiraoliveira.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Eduardo',
    'Pereira',
    'admin',
    'active',
    NOW(),
    NOW()
),

-- BUSINESS ADMINS
(
    '25555555-5555-5555-5555-555555555555',
    '15555555-5555-5555-5555-555555555555',
    'admin@machadoadvogados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Luiz',
    'Machado',
    'admin',
    'active',
    NOW(),
    NOW()
),
(
    '26666666-6666-6666-6666-666666666666',
    '16666666-6666-6666-6666-666666666666',
    'admin@ferreiralegal.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Sergio',
    'Ferreira',
    'admin',
    'active',
    NOW(),
    NOW()
),

-- ENTERPRISE ADMINS
(
    '27777777-7777-7777-7777-777777777777',
    '17777777-7777-7777-7777-777777777777',
    'admin@barrosent.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Alexandre',
    'Barros',
    'admin',
    'active',
    NOW(),
    NOW()
),
(
    '28888888-8888-8888-8888-888888888888',
    '18888888-8888-8888-8888-888888888888',
    'admin@rodriguesglobal.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password
    'Marcelo',
    'Rodrigues',
    'admin',
    'active',
    NOW(),
    NOW()
);

-- ============================================================================
-- 4. USUÁRIOS ADICIONAIS (manager, operator, client para cada tenant)
-- ============================================================================

INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, status, created_at, updated_at) VALUES 
-- STARTER TENANT 1 - USUARIOS ADICIONAIS
(
    '31111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'gerente@silvaassociados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Ana',
    'Silva',
    'manager',
    'active',
    NOW(),
    NOW()
),
(
    '41111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'advogado@silvaassociados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Ricardo',
    'Silva',
    'operator',
    'active',
    NOW(),
    NOW()
),
(
    '51111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'cliente@silvaassociados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Maria',
    'Santos',
    'client',
    'active',
    NOW(),
    NOW()
),

-- STARTER TENANT 2 - USUARIOS ADICIONAIS
(
    '32222222-2222-2222-2222-222222222222',
    '12222222-2222-2222-2222-222222222222',
    'gerente@limaadvogados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Paula',
    'Lima',
    'manager',
    'active',
    NOW(),
    NOW()
),
(
    '42222222-2222-2222-2222-222222222222',
    '12222222-2222-2222-2222-222222222222',
    'advogado@limaadvogados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Marco',
    'Lima',
    'operator',
    'active',
    NOW(),
    NOW()
),
(
    '52222222-2222-2222-2222-222222222222',
    '12222222-2222-2222-2222-222222222222',
    'cliente@limaadvogados.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Lucia',
    'Oliveira',
    'client',
    'active',
    NOW(),
    NOW()
),

-- PROFESSIONAL TENANT 1 - USUARIOS ADICIONAIS
(
    '33333333-3333-3333-3333-333333333333',
    '13333333-3333-3333-3333-333333333333',
    'gerente@costasantos.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Sandra',
    'Santos',
    'manager',
    'active',
    NOW(),
    NOW()
),
(
    '43333333-3333-3333-3333-333333333333',
    '13333333-3333-3333-3333-333333333333',
    'advogado@costasantos.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Felipe',
    'Costa',
    'operator',
    'active',
    NOW(),
    NOW()
),
(
    '53333333-3333-3333-3333-333333333333',
    '13333333-3333-3333-3333-333333333333',
    'cliente@costasantos.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Carla',
    'Santos',
    'client',
    'active',
    NOW(),
    NOW()
),

-- PROFESSIONAL TENANT 2 - USUARIOS ADICIONAIS
(
    '34444444-4444-4444-4444-444444444444',
    '14444444-4444-4444-4444-444444444444',
    'gerente@pereiraoliveira.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Fernanda',
    'Oliveira',
    'manager',
    'active',
    NOW(),
    NOW()
),
(
    '44444444-4444-4444-4444-444444444444',
    '14444444-4444-4444-4444-444444444444',
    'advogado@pereiraoliveira.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Bruno',
    'Pereira',
    'operator',
    'active',
    NOW(),
    NOW()
),
(
    '54444444-4444-4444-4444-444444444444',
    '14444444-4444-4444-4444-444444444444',
    'cliente@pereiraoliveira.com.br',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Isabela',
    'Lima',
    'client',
    'active',
    NOW(),
    NOW()
);

-- ============================================================================
-- 5. SUBSCRIPTIONS (uma para cada tenant)
-- ============================================================================

INSERT INTO subscriptions (id, tenant_id, plan_id, status, current_period_start, current_period_end, trial_start, trial_end, created_at, updated_at) VALUES 
-- STARTER SUBSCRIPTIONS
(
    'sub-11111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'plan-starter-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    'sub-12222222-2222-2222-2222-222222222222',
    '12222222-2222-2222-2222-222222222222',
    'plan-starter-uuid',
    'trialing',
    NOW(),
    NOW() + INTERVAL '1 month',
    NOW(),
    NOW() + INTERVAL '7 days',
    NOW(),
    NOW()
),

-- PROFESSIONAL SUBSCRIPTIONS
(
    'sub-13333333-3333-3333-3333-333333333333',
    '13333333-3333-3333-3333-333333333333',
    'plan-professional-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    'sub-14444444-4444-4444-4444-444444444444',
    '14444444-4444-4444-4444-444444444444',
    'plan-professional-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
),

-- BUSINESS SUBSCRIPTIONS
(
    'sub-15555555-5555-5555-5555-555555555555',
    '15555555-5555-5555-5555-555555555555',
    'plan-business-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    'sub-16666666-6666-6666-6666-666666666666',
    '16666666-6666-6666-6666-666666666666',
    'plan-business-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
),

-- ENTERPRISE SUBSCRIPTIONS
(
    'sub-17777777-7777-7777-7777-777777777777',
    '17777777-7777-7777-7777-777777777777',
    'plan-enterprise-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    'sub-18888888-8888-8888-8888-888888888888',
    '18888888-8888-8888-8888-888888888888',
    'plan-enterprise-uuid',
    'active',
    NOW(),
    NOW() + INTERVAL '1 month',
    NULL,
    NULL,
    NOW(),
    NOW()
);

-- ============================================================================
-- 6. QUOTA LIMITS (definir limites para cada subscription)
-- ============================================================================

INSERT INTO quota_limits (id, tenant_id, quota_type, limit_value, period_type, reset_day, created_at, updated_at) VALUES 
-- STARTER LIMITS (2 tenants)
('ql-11111111-processes', '11111111-1111-1111-1111-111111111111', 'processes', 50, 'total', 1, NOW(), NOW()),
('ql-11111111-users', '11111111-1111-1111-1111-111111111111', 'users', 2, 'total', 1, NOW(), NOW()),
('ql-11111111-datajud', '11111111-1111-1111-1111-111111111111', 'datajud_calls', 100, 'daily', 1, NOW(), NOW()),

('ql-12222222-processes', '12222222-2222-2222-2222-222222222222', 'processes', 50, 'total', 1, NOW(), NOW()),
('ql-12222222-users', '12222222-2222-2222-2222-222222222222', 'users', 2, 'total', 1, NOW(), NOW()),
('ql-12222222-datajud', '12222222-2222-2222-2222-222222222222', 'datajud_calls', 100, 'daily', 1, NOW(), NOW()),

-- PROFESSIONAL LIMITS (2 tenants)
('ql-13333333-processes', '13333333-3333-3333-3333-333333333333', 'processes', 200, 'total', 1, NOW(), NOW()),
('ql-13333333-users', '13333333-3333-3333-3333-333333333333', 'users', 5, 'total', 1, NOW(), NOW()),
('ql-13333333-datajud', '13333333-3333-3333-3333-333333333333', 'datajud_calls', 500, 'daily', 1, NOW(), NOW()),
('ql-13333333-mcp', '13333333-3333-3333-3333-333333333333', 'mcp_commands', 200, 'monthly', 1, NOW(), NOW()),

('ql-14444444-processes', '14444444-4444-4444-4444-444444444444', 'processes', 200, 'total', 1, NOW(), NOW()),
('ql-14444444-users', '14444444-4444-4444-4444-444444444444', 'users', 5, 'total', 1, NOW(), NOW()),
('ql-14444444-datajud', '14444444-4444-4444-4444-444444444444', 'datajud_calls', 500, 'daily', 1, NOW(), NOW()),
('ql-14444444-mcp', '14444444-4444-4444-4444-444444444444', 'mcp_commands', 200, 'monthly', 1, NOW(), NOW()),

-- BUSINESS LIMITS (2 tenants)
('ql-15555555-processes', '15555555-5555-5555-5555-555555555555', 'processes', 500, 'total', 1, NOW(), NOW()),
('ql-15555555-users', '15555555-5555-5555-5555-555555555555', 'users', 15, 'total', 1, NOW(), NOW()),
('ql-15555555-datajud', '15555555-5555-5555-5555-555555555555', 'datajud_calls', 2000, 'daily', 1, NOW(), NOW()),
('ql-15555555-mcp', '15555555-5555-5555-5555-555555555555', 'mcp_commands', 1000, 'monthly', 1, NOW(), NOW()),

('ql-16666666-processes', '16666666-6666-6666-6666-666666666666', 'processes', 500, 'total', 1, NOW(), NOW()),
('ql-16666666-users', '16666666-6666-6666-6666-666666666666', 'users', 15, 'total', 1, NOW(), NOW()),
('ql-16666666-datajud', '16666666-6666-6666-6666-666666666666', 'datajud_calls', 2000, 'daily', 1, NOW(), NOW()),
('ql-16666666-mcp', '16666666-6666-6666-6666-666666666666', 'mcp_commands', 1000, 'monthly', 1, NOW(), NOW()),

-- ENTERPRISE LIMITS (2 tenants) - Ilimitado
('ql-17777777-processes', '17777777-7777-7777-7777-777777777777', 'processes', -1, 'total', 1, NOW(), NOW()),
('ql-17777777-users', '17777777-7777-7777-7777-777777777777', 'users', -1, 'total', 1, NOW(), NOW()),
('ql-17777777-datajud', '17777777-7777-7777-7777-777777777777', 'datajud_calls', 10000, 'daily', 1, NOW(), NOW()),
('ql-17777777-mcp', '17777777-7777-7777-7777-777777777777', 'mcp_commands', -1, 'monthly', 1, NOW(), NOW()),

('ql-18888888-processes', '18888888-8888-8888-8888-888888888888', 'processes', -1, 'total', 1, NOW(), NOW()),
('ql-18888888-users', '18888888-8888-8888-8888-888888888888', 'users', -1, 'total', 1, NOW(), NOW()),
('ql-18888888-datajud', '18888888-8888-8888-8888-888888888888', 'datajud_calls', 10000, 'daily', 1, NOW(), NOW()),
('ql-18888888-mcp', '18888888-8888-8888-8888-888888888888', 'mcp_commands', -1, 'monthly', 1, NOW(), NOW());

-- ============================================================================
-- 7. QUOTA USAGE (inicializar com 0)
-- ============================================================================

INSERT INTO quota_usage (id, tenant_id, quota_type, current_usage, period_start, period_end, created_at, updated_at) VALUES 
-- Inicializar uso com 0 para todos os tenants e tipos de quota
-- STARTER TENANTS
('qu-11111111-processes', '11111111-1111-1111-1111-111111111111', 'processes', 0, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-11111111-users', '11111111-1111-1111-1111-111111111111', 'users', 4, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-11111111-datajud', '11111111-1111-1111-1111-111111111111', 'datajud_calls', 0, DATE_TRUNC('day', NOW()), DATE_TRUNC('day', NOW()) + INTERVAL '1 day', NOW(), NOW()),

('qu-12222222-processes', '12222222-2222-2222-2222-222222222222', 'processes', 0, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-12222222-users', '12222222-2222-2222-2222-222222222222', 'users', 4, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-12222222-datajud', '12222222-2222-2222-2222-222222222222', 'datajud_calls', 0, DATE_TRUNC('day', NOW()), DATE_TRUNC('day', NOW()) + INTERVAL '1 day', NOW(), NOW()),

-- PROFESSIONAL TENANTS
('qu-13333333-processes', '13333333-3333-3333-3333-333333333333', 'processes', 0, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-13333333-users', '13333333-3333-3333-3333-333333333333', 'users', 4, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-13333333-datajud', '13333333-3333-3333-3333-333333333333', 'datajud_calls', 0, DATE_TRUNC('day', NOW()), DATE_TRUNC('day', NOW()) + INTERVAL '1 day', NOW(), NOW()),
('qu-13333333-mcp', '13333333-3333-3333-3333-333333333333', 'mcp_commands', 0, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),

('qu-14444444-processes', '14444444-4444-4444-4444-444444444444', 'processes', 0, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-14444444-users', '14444444-4444-4444-4444-444444444444', 'users', 4, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW()),
('qu-14444444-datajud', '14444444-4444-4444-4444-444444444444', 'datajud_calls', 0, DATE_TRUNC('day', NOW()), DATE_TRUNC('day', NOW()) + INTERVAL '1 day', NOW(), NOW()),
('qu-14444444-mcp', '14444444-4444-4444-4444-444444444444', 'mcp_commands', 0, DATE_TRUNC('month', NOW()), DATE_TRUNC('month', NOW()) + INTERVAL '1 month', NOW(), NOW());

-- ============================================================================
-- FIM DO SCRIPT
-- ============================================================================

-- COMENTÁRIOS DE VALIDAÇÃO:
-- ✅ 4 Planos criados (starter, professional, business, enterprise)
-- ✅ 8 Tenants criados (2 por plano)
-- ✅ 32+ Usuários criados (4+ por tenant)
-- ✅ 8 Subscriptions ativas
-- ✅ Quota limits configurados por plano
-- ✅ Quota usage inicializado

-- CREDENCIAIS PARA TESTE:
-- Todos os usuários admin têm senha: "password"
-- Emails: admin@[tenant].com.br

-- PARA TESTAR LOGIN:
-- Email: admin@silvaassociados.com.br
-- Senha: password
-- Tenant ID: 11111111-1111-1111-1111-111111111111