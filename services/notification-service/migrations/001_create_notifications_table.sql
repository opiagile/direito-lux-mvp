-- Criar tabela de notificações
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'normal',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    content_html TEXT,
    recipient_id VARCHAR(255) NOT NULL,
    recipient_type VARCHAR(50) NOT NULL,
    recipient_contact VARCHAR(255) NOT NULL,
    tenant_id UUID NOT NULL,
    user_id UUID,
    process_id UUID,
    template_id UUID,
    variables JSONB DEFAULT '{}'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para otimizar buscas
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_channel ON notifications(channel);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient ON notifications(recipient_id, recipient_type);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_process_id ON notifications(process_id) WHERE process_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled_at ON notifications(scheduled_at) WHERE scheduled_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_expires_at ON notifications(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_sent_at ON notifications(sent_at) WHERE sent_at IS NOT NULL;

-- Índice composto para busca eficiente de notificações pendentes por canal
CREATE INDEX IF NOT EXISTS idx_notifications_pending_by_channel ON notifications(channel, priority, created_at) 
WHERE status = 'pending';

-- Índice composto para busca de notificações agendadas
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled ON notifications(scheduled_at, status) 
WHERE scheduled_at IS NOT NULL AND status = 'pending';

-- Índice para busca de notificações para retry
CREATE INDEX IF NOT EXISTS idx_notifications_retry ON notifications(status, retry_count, max_retries, failed_at) 
WHERE status = 'failed';

-- Índice JSONB para busca por variáveis
CREATE INDEX IF NOT EXISTS idx_notifications_variables_gin ON notifications USING GIN(variables);
CREATE INDEX IF NOT EXISTS idx_notifications_metadata_gin ON notifications USING GIN(metadata);

-- Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_notifications_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notifications_updated_at
    BEFORE UPDATE ON notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_notifications_updated_at();

-- Constraints para validar valores dos enums
ALTER TABLE notifications ADD CONSTRAINT check_notification_type 
CHECK (type IN (
    'process_update', 'movement_alert', 'deadline_reminder',
    'trial_expiring', 'subscription_due', 'system_alert',
    'welcome', 'password_reset'
));

ALTER TABLE notifications ADD CONSTRAINT check_notification_channel 
CHECK (channel IN ('whatsapp', 'email', 'telegram', 'push', 'sms'));

ALTER TABLE notifications ADD CONSTRAINT check_notification_priority 
CHECK (priority IN ('low', 'normal', 'high', 'critical'));

ALTER TABLE notifications ADD CONSTRAINT check_notification_status 
CHECK (status IN ('pending', 'processing', 'sent', 'failed', 'cancelled'));

ALTER TABLE notifications ADD CONSTRAINT check_recipient_type 
CHECK (recipient_type IN ('user', 'admin', 'system', 'external'));

-- Constraints de negócio
ALTER TABLE notifications ADD CONSTRAINT check_retry_count_positive 
CHECK (retry_count >= 0);

ALTER TABLE notifications ADD CONSTRAINT check_max_retries_positive 
CHECK (max_retries >= 0);

ALTER TABLE notifications ADD CONSTRAINT check_scheduled_future 
CHECK (scheduled_at IS NULL OR scheduled_at > created_at);

ALTER TABLE notifications ADD CONSTRAINT check_expires_future 
CHECK (expires_at IS NULL OR expires_at > created_at);

-- Comentários para documentação
COMMENT ON TABLE notifications IS 'Tabela principal de notificações do sistema';
COMMENT ON COLUMN notifications.type IS 'Tipo da notificação (process_update, movement_alert, etc.)';
COMMENT ON COLUMN notifications.channel IS 'Canal de envio (whatsapp, email, telegram, push, sms)';
COMMENT ON COLUMN notifications.priority IS 'Prioridade da notificação (low, normal, high, critical)';
COMMENT ON COLUMN notifications.status IS 'Status atual (pending, processing, sent, failed, cancelled)';
COMMENT ON COLUMN notifications.recipient_id IS 'ID do destinatário (user_id, admin_id, etc.)';
COMMENT ON COLUMN notifications.recipient_type IS 'Tipo do destinatário (user, admin, system, external)';
COMMENT ON COLUMN notifications.recipient_contact IS 'Contato do destinatário (email, telefone, etc.)';
COMMENT ON COLUMN notifications.variables IS 'Variáveis para renderização do template (formato JSON)';
COMMENT ON COLUMN notifications.metadata IS 'Metadados adicionais (formato JSON)';
COMMENT ON COLUMN notifications.scheduled_at IS 'Data/hora agendada para envio (NULL = envio imediato)';
COMMENT ON COLUMN notifications.retry_count IS 'Número de tentativas de reenvio já realizadas';
COMMENT ON COLUMN notifications.max_retries IS 'Número máximo de tentativas de reenvio';
COMMENT ON COLUMN notifications.expires_at IS 'Data/hora de expiração (NULL = nunca expira)';