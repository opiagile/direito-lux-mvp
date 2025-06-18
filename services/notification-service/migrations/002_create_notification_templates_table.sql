-- Criar tabela de templates de notificação
CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    content_html TEXT,
    variables TEXT[] DEFAULT ARRAY[]::TEXT[],
    tenant_id UUID,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para otimizar buscas
CREATE INDEX IF NOT EXISTS idx_notification_templates_tenant_id ON notification_templates(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notification_templates_type ON notification_templates(type);
CREATE INDEX IF NOT EXISTS idx_notification_templates_channel ON notification_templates(channel);
CREATE INDEX IF NOT EXISTS idx_notification_templates_status ON notification_templates(status);
CREATE INDEX IF NOT EXISTS idx_notification_templates_is_system ON notification_templates(is_system);
CREATE INDEX IF NOT EXISTS idx_notification_templates_name ON notification_templates(name);

-- Índice composto para busca eficiente por tipo e canal
CREATE INDEX IF NOT EXISTS idx_notification_templates_type_channel ON notification_templates(type, channel, status);

-- Índice composto para templates ativos por tenant
CREATE INDEX IF NOT EXISTS idx_notification_templates_active_tenant ON notification_templates(tenant_id, status, type, channel) 
WHERE status = 'active';

-- Índice para templates do sistema ativos
CREATE INDEX IF NOT EXISTS idx_notification_templates_system_active ON notification_templates(type, channel) 
WHERE is_system = TRUE AND status = 'active';

-- Trigger para atualizar updated_at automaticamente
CREATE TRIGGER trigger_notification_templates_updated_at
    BEFORE UPDATE ON notification_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_notifications_updated_at();

-- Constraints para validar valores dos enums
ALTER TABLE notification_templates ADD CONSTRAINT check_template_type 
CHECK (type IN (
    'process_update', 'movement_alert', 'deadline_reminder',
    'trial_expiring', 'subscription_due', 'system_alert',
    'welcome', 'password_reset'
));

ALTER TABLE notification_templates ADD CONSTRAINT check_template_channel 
CHECK (channel IN ('whatsapp', 'email', 'telegram', 'push', 'sms'));

ALTER TABLE notification_templates ADD CONSTRAINT check_template_status 
CHECK (status IN ('active', 'inactive', 'draft'));

-- Constraint para garantir que templates do sistema não tenham tenant
ALTER TABLE notification_templates ADD CONSTRAINT check_system_template_no_tenant 
CHECK ((is_system = TRUE AND tenant_id IS NULL) OR (is_system = FALSE AND tenant_id IS NOT NULL));

-- Constraint de unicidade: apenas um template ativo por tipo/canal por tenant
CREATE UNIQUE INDEX IF NOT EXISTS idx_notification_templates_unique_active_tenant 
ON notification_templates(tenant_id, type, channel) 
WHERE status = 'active' AND tenant_id IS NOT NULL;

-- Constraint de unicidade: apenas um template de sistema ativo por tipo/canal
CREATE UNIQUE INDEX IF NOT EXISTS idx_notification_templates_unique_active_system 
ON notification_templates(type, channel) 
WHERE status = 'active' AND is_system = TRUE;

-- Comentários para documentação
COMMENT ON TABLE notification_templates IS 'Templates de notificação reutilizáveis';
COMMENT ON COLUMN notification_templates.name IS 'Nome descritivo do template';
COMMENT ON COLUMN notification_templates.type IS 'Tipo da notificação que o template serve';
COMMENT ON COLUMN notification_templates.channel IS 'Canal para o qual o template é destinado';
COMMENT ON COLUMN notification_templates.status IS 'Status do template (active, inactive, draft)';
COMMENT ON COLUMN notification_templates.variables IS 'Lista de variáveis que o template aceita';
COMMENT ON COLUMN notification_templates.tenant_id IS 'ID do tenant proprietário (NULL para templates do sistema)';
COMMENT ON COLUMN notification_templates.is_system IS 'Indica se é um template padrão do sistema';

-- Inserir templates padrão do sistema
INSERT INTO notification_templates (id, name, type, channel, status, subject, content, variables, is_system, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'Bem-vindo',
    'welcome',
    'email',
    'active',
    'Bem-vindo ao Direito Lux, {{nome}}!',
    'Olá {{nome}}, seja bem-vindo ao Direito Lux! Sua conta foi criada com sucesso.',
    ARRAY['nome'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Movimentação Processual - WhatsApp',
    'process_update',
    'whatsapp',
    'active',
    'Nova movimentação no processo {{numero_processo}}',
    'Olá {{nome}}, houve uma nova movimentação no processo {{numero_processo}}: {{descricao_movimento}}. Data: {{data_movimento}}',
    ARRAY['nome', 'numero_processo', 'descricao_movimento', 'data_movimento'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Movimentação Processual - Email',
    'process_update',
    'email',
    'active',
    'Nova movimentação no processo {{numero_processo}}',
    'Prezado(a) {{nome}},

Informamos que houve uma nova movimentação no processo {{numero_processo}}.

Detalhes:
- Processo: {{numero_processo}}
- Movimento: {{descricao_movimento}}
- Data: {{data_movimento}}

Para mais detalhes, acesse sua conta no Direito Lux.

Atenciosamente,
Equipe Direito Lux',
    ARRAY['nome', 'numero_processo', 'descricao_movimento', 'data_movimento'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Lembrete de Prazo',
    'deadline_reminder',
    'whatsapp',
    'active',
    'Prazo vencendo em {{dias}} dias - Processo {{numero_processo}}',
    'Atenção {{nome}}, o prazo para o processo {{numero_processo}} vence em {{dias}} dias. Descrição: {{descricao_prazo}}',
    ARRAY['nome', 'numero_processo', 'dias', 'descricao_prazo'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Trial Expirando',
    'trial_expiring',
    'email',
    'active',
    'Seu período de teste expira em {{dias}} dias',
    'Olá {{nome}}, seu período de teste do Direito Lux expira em {{dias}} dias. Não perca tempo e assine já!',
    ARRAY['nome', 'dias'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Alerta de Sistema',
    'system_alert',
    'email',
    'active',
    'Alerta do Sistema - {{titulo}}',
    'Caro usuário,

Detectamos a seguinte situação que requer sua atenção:

{{mensagem}}

Data/Hora: {{data_hora}}
Severidade: {{severidade}}

Por favor, tome as ações necessárias.

Atenciosamente,
Sistema Direito Lux',
    ARRAY['titulo', 'mensagem', 'data_hora', 'severidade'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Reset de Senha',
    'password_reset',
    'email',
    'active',
    'Redefinição de senha - Direito Lux',
    'Olá {{nome}},

Recebemos uma solicitação para redefinir a senha da sua conta.

Clique no link abaixo para redefinir sua senha:
{{link_reset}}

Este link expira em 24 horas.

Se você não solicitou esta redefinição, ignore este email.

Atenciosamente,
Equipe Direito Lux',
    ARRAY['nome', 'link_reset'],
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);