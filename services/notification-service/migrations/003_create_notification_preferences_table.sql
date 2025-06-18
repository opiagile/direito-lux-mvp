-- Criar tabela de preferências de notificação dos usuários
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    channels TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para otimizar buscas
CREATE INDEX IF NOT EXISTS idx_notification_preferences_tenant_id ON notification_preferences(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notification_preferences_user_id ON notification_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_preferences_type ON notification_preferences(type);
CREATE INDEX IF NOT EXISTS idx_notification_preferences_enabled ON notification_preferences(enabled);

-- Índice composto para busca eficiente por usuário
CREATE INDEX IF NOT EXISTS idx_notification_preferences_user_type ON notification_preferences(user_id, type);

-- Índice para preferências ativas por tenant
CREATE INDEX IF NOT EXISTS idx_notification_preferences_tenant_enabled ON notification_preferences(tenant_id, enabled) 
WHERE enabled = TRUE;

-- Trigger para atualizar updated_at automaticamente
CREATE TRIGGER trigger_notification_preferences_updated_at
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_notifications_updated_at();

-- Constraints para validar valores dos enums
ALTER TABLE notification_preferences ADD CONSTRAINT check_preference_type 
CHECK (type IN (
    'process_update', 'movement_alert', 'deadline_reminder',
    'trial_expiring', 'subscription_due', 'system_alert',
    'welcome', 'password_reset'
));

-- Constraint para validar canais
ALTER TABLE notification_preferences ADD CONSTRAINT check_preference_channels 
CHECK (
    channels <@ ARRAY['whatsapp', 'email', 'telegram', 'push', 'sms']::TEXT[] 
    AND array_length(channels, 1) > 0
);

-- Constraint de unicidade: apenas uma preferência por usuário/tipo
CREATE UNIQUE INDEX IF NOT EXISTS idx_notification_preferences_unique_user_type 
ON notification_preferences(user_id, type);

-- Comentários para documentação
COMMENT ON TABLE notification_preferences IS 'Preferências de notificação por usuário e tipo';
COMMENT ON COLUMN notification_preferences.tenant_id IS 'ID do tenant do usuário';
COMMENT ON COLUMN notification_preferences.user_id IS 'ID do usuário';
COMMENT ON COLUMN notification_preferences.type IS 'Tipo de notificação';
COMMENT ON COLUMN notification_preferences.channels IS 'Lista de canais preferidos para este tipo';
COMMENT ON COLUMN notification_preferences.enabled IS 'Se as notificações deste tipo estão habilitadas';

-- Criar tabela de configurações globais de notificação por tenant
CREATE TABLE IF NOT EXISTS notification_tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL UNIQUE,
    whatsapp_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    email_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    telegram_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sms_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    whatsapp_business_token VARCHAR(500),
    whatsapp_phone_number VARCHAR(20),
    whatsapp_webhook_url VARCHAR(500),
    telegram_bot_token VARCHAR(500),
    telegram_chat_id VARCHAR(100),
    smtp_host VARCHAR(255),
    smtp_port INTEGER,
    smtp_username VARCHAR(255),
    smtp_password VARCHAR(500),
    smtp_from_email VARCHAR(255),
    smtp_from_name VARCHAR(255),
    daily_limit_whatsapp INTEGER DEFAULT 1000,
    daily_limit_email INTEGER DEFAULT 5000,
    daily_limit_telegram INTEGER DEFAULT 1000,
    daily_limit_sms INTEGER DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para configurações do tenant
CREATE INDEX IF NOT EXISTS idx_notification_tenant_settings_tenant_id ON notification_tenant_settings(tenant_id);

-- Trigger para atualizar updated_at automaticamente
CREATE TRIGGER trigger_notification_tenant_settings_updated_at
    BEFORE UPDATE ON notification_tenant_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_notifications_updated_at();

-- Constraints de validação
ALTER TABLE notification_tenant_settings ADD CONSTRAINT check_daily_limits_positive 
CHECK (
    daily_limit_whatsapp >= 0 AND 
    daily_limit_email >= 0 AND 
    daily_limit_telegram >= 0 AND 
    daily_limit_sms >= 0
);

ALTER TABLE notification_tenant_settings ADD CONSTRAINT check_smtp_config 
CHECK (
    (email_enabled = FALSE) OR 
    (email_enabled = TRUE AND smtp_host IS NOT NULL AND smtp_port IS NOT NULL)
);

-- Comentários para documentação
COMMENT ON TABLE notification_tenant_settings IS 'Configurações de notificação por tenant';
COMMENT ON COLUMN notification_tenant_settings.tenant_id IS 'ID único do tenant';
COMMENT ON COLUMN notification_tenant_settings.whatsapp_business_token IS 'Token da WhatsApp Business API';
COMMENT ON COLUMN notification_tenant_settings.whatsapp_phone_number IS 'Número de telefone WhatsApp Business';
COMMENT ON COLUMN notification_tenant_settings.telegram_bot_token IS 'Token do bot Telegram';
COMMENT ON COLUMN notification_tenant_settings.daily_limit_whatsapp IS 'Limite diário de mensagens WhatsApp';

-- Criar tabela de histórico de entrega
CREATE TABLE IF NOT EXISTS notification_delivery_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    channel VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    provider_id VARCHAR(100),
    provider_response TEXT,
    delivered_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    error_code VARCHAR(50),
    error_message TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para histórico de entrega
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_notification_id ON notification_delivery_history(notification_id);
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_channel ON notification_delivery_history(channel);
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_status ON notification_delivery_history(status);
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_provider_id ON notification_delivery_history(provider_id) WHERE provider_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_delivered_at ON notification_delivery_history(delivered_at) WHERE delivered_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_created_at ON notification_delivery_history(created_at);

-- Índice JSONB para metadados
CREATE INDEX IF NOT EXISTS idx_notification_delivery_history_metadata_gin ON notification_delivery_history USING GIN(metadata);

-- Constraints para histórico de entrega
ALTER TABLE notification_delivery_history ADD CONSTRAINT check_delivery_channel 
CHECK (channel IN ('whatsapp', 'email', 'telegram', 'push', 'sms'));

ALTER TABLE notification_delivery_history ADD CONSTRAINT check_delivery_status 
CHECK (status IN ('sent', 'delivered', 'read', 'failed', 'bounced', 'complained'));

-- Comentários para documentação
COMMENT ON TABLE notification_delivery_history IS 'Histórico detalhado de entrega de notificações';
COMMENT ON COLUMN notification_delivery_history.notification_id IS 'ID da notificação relacionada';
COMMENT ON COLUMN notification_delivery_history.provider_id IS 'ID do provedor (ex: message_id do WhatsApp)';
COMMENT ON COLUMN notification_delivery_history.provider_response IS 'Resposta completa do provedor';
COMMENT ON COLUMN notification_delivery_history.error_code IS 'Código de erro do provedor';
COMMENT ON COLUMN notification_delivery_history.metadata IS 'Metadados adicionais do provedor';