-- =============================================================================
-- Rollback: Remover tabela de tokens de recuperação de senha
-- =============================================================================

-- Remover tabela de tokens de recuperação de senha
DROP TABLE IF EXISTS password_reset_tokens;