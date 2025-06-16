-- =============================================================================
-- Migração: Criar tabelas de sessões e tokens
-- =============================================================================

-- Tabela de sessões
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    refresh_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT sessions_expires_at_future CHECK (expires_at > created_at),
    CONSTRAINT sessions_refresh_expires_future CHECK (refresh_expires_at > expires_at)
);

-- Tabela de refresh tokens
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    used_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT refresh_tokens_expires_future CHECK (expires_at > created_at),
    CONSTRAINT refresh_tokens_used_at_check CHECK (
        (is_used = false AND used_at IS NULL) OR 
        (is_used = true AND used_at IS NOT NULL)
    )
);

-- Tabela de tentativas de login
CREATE TABLE login_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    tenant_id UUID,
    success BOOLEAN NOT NULL DEFAULT false,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT login_attempts_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

-- Índices para sessions
CREATE UNIQUE INDEX idx_sessions_access_token ON sessions(access_token);
CREATE UNIQUE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_tenant_id ON sessions(tenant_id);
CREATE INDEX idx_sessions_is_active ON sessions(is_active);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_user_active ON sessions(user_id, is_active);

-- Índices para refresh_tokens
CREATE UNIQUE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_tenant_id ON refresh_tokens(tenant_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_is_used ON refresh_tokens(is_used);
CREATE INDEX idx_refresh_tokens_user_unused ON refresh_tokens(user_id, is_used) WHERE is_used = false;

-- Índices para login_attempts
CREATE INDEX idx_login_attempts_email ON login_attempts(email);
CREATE INDEX idx_login_attempts_created_at ON login_attempts(created_at);
CREATE INDEX idx_login_attempts_email_created ON login_attempts(email, created_at);
CREATE INDEX idx_login_attempts_success ON login_attempts(success);
CREATE INDEX idx_login_attempts_ip ON login_attempts(ip_address);

-- Trigger para atualizar updated_at em sessions
CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentários
COMMENT ON TABLE sessions IS 'Sessões ativas de usuários';
COMMENT ON COLUMN sessions.access_token IS 'Token JWT de acesso';
COMMENT ON COLUMN sessions.refresh_token IS 'Token para renovação';
COMMENT ON COLUMN sessions.expires_at IS 'Expiração do access token';
COMMENT ON COLUMN sessions.refresh_expires_at IS 'Expiração do refresh token';

COMMENT ON TABLE refresh_tokens IS 'Tokens de refresh utilizados';
COMMENT ON COLUMN refresh_tokens.token IS 'Token de refresh';
COMMENT ON COLUMN refresh_tokens.is_used IS 'Se o token já foi utilizado';
COMMENT ON COLUMN refresh_tokens.used_at IS 'Quando o token foi utilizado';

COMMENT ON TABLE login_attempts IS 'Tentativas de login (auditoria)';
COMMENT ON COLUMN login_attempts.success IS 'Se a tentativa foi bem-sucedida';
COMMENT ON COLUMN login_attempts.ip_address IS 'Endereço IP da tentativa';