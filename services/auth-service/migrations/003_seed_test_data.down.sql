-- Rollback migration: 003_seed_test_data.sql
-- Remove dados de teste

-- ============================================================================
-- REMOVER SESSÕES DE TESTE
-- ============================================================================

DELETE FROM sessions WHERE id IN (
    'session-starter-admin',
    'session-prof-admin', 
    'session-biz-admin',
    'session-ent-admin'
);

-- ============================================================================
-- REMOVER USUÁRIOS DE TESTE
-- ============================================================================

DELETE FROM users WHERE tenant_id IN (
    'tenant-starter-01', 'tenant-starter-02',
    'tenant-prof-01', 'tenant-prof-02',
    'tenant-biz-01', 'tenant-biz-02', 
    'tenant-ent-01', 'tenant-ent-02'
);

-- ============================================================================
-- REMOVER SUBSCRIPTIONS DE TESTE
-- ============================================================================

DELETE FROM subscriptions WHERE id IN (
    'sub-starter-01', 'sub-starter-02',
    'sub-prof-01', 'sub-prof-02',
    'sub-biz-01', 'sub-biz-02',
    'sub-ent-01', 'sub-ent-02'
);

-- ============================================================================
-- REMOVER TENANTS DE TESTE
-- ============================================================================

DELETE FROM tenants WHERE id IN (
    'tenant-starter-01', 'tenant-starter-02',
    'tenant-prof-01', 'tenant-prof-02',
    'tenant-biz-01', 'tenant-biz-02',
    'tenant-ent-01', 'tenant-ent-02'
);