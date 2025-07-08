-- =============================================================================
-- Migração: Criar tabela de tokens de recuperação de senha
-- =============================================================================

-- Tabela de tokens de recuperação de senha
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    used_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT password_reset_tokens_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$'),
    CONSTRAINT password_reset_tokens_token_check CHECK (length(token) >= 32),
    CONSTRAINT password_reset_tokens_expires_at_check CHECK (expires_at > created_at),
    CONSTRAINT password_reset_tokens_used_at_check CHECK (used_at IS NULL OR (is_used = TRUE AND used_at >= created_at))
);

-- Índices
CREATE UNIQUE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_email ON password_reset_tokens(email);
CREATE INDEX idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);
CREATE INDEX idx_password_reset_tokens_is_used ON password_reset_tokens(is_used);
CREATE INDEX idx_password_reset_tokens_created_at ON password_reset_tokens(created_at);

-- Índice composto para limpeza de tokens expirados
CREATE INDEX idx_password_reset_tokens_expired ON password_reset_tokens(expires_at, is_used) 
WHERE expires_at < NOW() OR is_used = TRUE;

-- Comentários
COMMENT ON TABLE password_reset_tokens IS 'Tabela de tokens para recuperação de senha';
COMMENT ON COLUMN password_reset_tokens.id IS 'Identificador único do token';
COMMENT ON COLUMN password_reset_tokens.user_id IS 'Identificador do usuário';
COMMENT ON COLUMN password_reset_tokens.tenant_id IS 'Identificador do tenant';
COMMENT ON COLUMN password_reset_tokens.token IS 'Token de recuperação (hash único)';
COMMENT ON COLUMN password_reset_tokens.email IS 'Email do usuário para verificação';
COMMENT ON COLUMN password_reset_tokens.expires_at IS 'Data/hora de expiração do token';
COMMENT ON COLUMN password_reset_tokens.is_used IS 'Se o token já foi utilizado';
COMMENT ON COLUMN password_reset_tokens.created_at IS 'Data/hora de criação';
COMMENT ON COLUMN password_reset_tokens.used_at IS 'Data/hora de utilização do token';