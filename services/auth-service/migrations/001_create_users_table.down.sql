-- =============================================================================
-- Migração Reversa: Remover tabela de usuários
-- =============================================================================

-- Remover trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Remover função
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Remover tabela
DROP TABLE IF EXISTS users;

-- Remover enums
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;