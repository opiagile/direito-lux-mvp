-- =============================================================================
-- Migração Reversa: Remover tabelas de sessões e tokens
-- =============================================================================

-- Remover triggers
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

-- Remover tabelas
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS sessions;