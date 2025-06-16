#!/bin/bash

# =============================================================================
# Direito Lux - Criar Dados de Exemplo
# =============================================================================

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo -e "${BLUE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                     DIREITO LUX                             ║
║                 Criando Dados de Exemplo                    ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Verificar se o PostgreSQL está rodando
log_info "Verificando se PostgreSQL está rodando..."
if ! docker-compose exec -T postgres pg_isready -U direito_lux 2>/dev/null; then
    log_error "PostgreSQL não está rodando. Execute primeiro: ./scripts/setup-local.sh"
    exit 1
fi

log_success "PostgreSQL está rodando ✓"

# Executar SQL para criar dados de exemplo
log_info "Criando dados de exemplo..."

docker-compose exec -T postgres psql -U direito_lux -d direito_lux_dev << 'EOF'

-- =============================================================================
-- DIREITO LUX - DADOS DE EXEMPLO
-- =============================================================================

-- Tenants de exemplo
INSERT INTO tenant.tenants (
    id, name, document, plan, status, created_at, updated_at
) VALUES 
(
    'tenant-001',
    'Silva & Associados Advogados',
    '12.345.678/0001-90',
    'professional',
    'active',
    NOW(),
    NOW()
),
(
    'tenant-002', 
    'Departamento Jurídico Empresa XYZ',
    '98.765.432/0001-10',
    'business',
    'active',
    NOW(),
    NOW()
),
(
    'tenant-003',
    'Dr. João Advogado Autônomo',
    '123.456.789-00',
    'starter',
    'trial',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Quotas dos tenants
INSERT INTO tenant.tenant_quotas (
    tenant_id, process_limit, users_limit, datajud_daily_quota, 
    storage_quota_gb, notifications_monthly, created_at, updated_at
) VALUES
(
    'tenant-001', 200, 5, 500, 50, 5000, NOW(), NOW()
),
(
    'tenant-002', 500, 20, 2000, 200, 999999, NOW(), NOW()
),
(
    'tenant-003', 50, 1, 100, 5, 500, NOW(), NOW()
)
ON CONFLICT (tenant_id) DO NOTHING;

-- Usuários de exemplo
INSERT INTO auth.users (
    id, email, tenant_id, role, first_name, last_name, 
    status, created_at, updated_at
) VALUES
(
    'user-001',
    'admin@silva-advogados.com',
    'tenant-001',
    'admin',
    'Carlos',
    'Silva',
    'active',
    NOW(),
    NOW()
),
(
    'user-002',
    'maria@silva-advogados.com', 
    'tenant-001',
    'lawyer',
    'Maria',
    'Santos',
    'active',
    NOW(),
    NOW()
),
(
    'user-003',
    'juridico@empresaxyz.com',
    'tenant-002',
    'admin',
    'Ana',
    'Costa',
    'active',
    NOW(),
    NOW()
),
(
    'user-004',
    'joao@advogado.com',
    'tenant-003',
    'admin',
    'João',
    'Oliveira',
    'active',
    NOW(),
    NOW()
),
(
    'user-005',
    'cliente1@empresa.com',
    'tenant-001',
    'client',
    'Roberto',
    'Cliente',
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Clientes
INSERT INTO tenant.clients (
    id, tenant_id, name, document, email, phone, 
    type, created_at, updated_at
) VALUES
(
    'client-001',
    'tenant-001',
    'Empresa ABC Ltda',
    '11.222.333/0001-44',
    'contato@empresaabc.com',
    '+5511999888777',
    'legal_entity',
    NOW(),
    NOW()
),
(
    'client-002',
    'tenant-001',
    'José da Silva',
    '123.456.789-10',
    'jose.silva@email.com',
    '+5511888777666',
    'individual',
    NOW(),
    NOW()
),
(
    'client-003',
    'tenant-002',
    'Fornecedor XYZ S.A.',
    '55.666.777/0001-88',
    'juridico@fornecedorxyz.com',
    '+5511777666555',
    'legal_entity',
    NOW(),
    NOW()
),
(
    'client-004',
    'tenant-003',
    'Maria Souza',
    '987.654.321-00',
    'maria.souza@email.com',
    '+5511666555444',
    'individual',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Processos de exemplo
INSERT INTO process.processes (
    id, number, tenant_id, client_id, court_code, court_name,
    class_code, class_name, subject_code, subject_name,
    status, monitoring_status, created_at, updated_at
) VALUES
(
    'process-001',
    '1234567-89.2023.8.26.0001',
    'tenant-001',
    'client-001',
    'TJSP',
    'Tribunal de Justiça de São Paulo',
    '319',
    'Execução Fiscal',
    '1834',
    'IPTU / Imposto Predial e Territorial Urbano',
    'active',
    'monitoring',
    NOW(),
    NOW()
),
(
    'process-002',
    '9876543-21.2023.5.02.0001',
    'tenant-001', 
    'client-002',
    'TRT2',
    'Tribunal Regional do Trabalho 2ª Região',
    '1394',
    'Reclamação Trabalhista',
    '7547',
    'Aviso Prévio',
    'active',
    'monitoring',
    NOW(),
    NOW()
),
(
    'process-003',
    '5555555-55.2023.8.26.0100',
    'tenant-002',
    'client-003',
    'TJSP',
    'Tribunal de Justiça de São Paulo',
    '275',
    'Ação de Cobrança',
    '1234',
    'Contratos Bancários',
    'active',
    'monitoring',
    NOW(),
    NOW()
),
(
    'process-004',
    '1111111-11.2023.4.03.6100',
    'tenant-003',
    'client-004',
    'TRF3',
    'Tribunal Regional Federal 3ª Região',
    '319',
    'Execução Fiscal',
    '6185',
    'Imposto de Renda Pessoa Física',
    'active',
    'paused',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Movimentações dos processos
INSERT INTO process.movements (
    id, process_id, date, description, type, 
    created_at, updated_at
) VALUES
(
    'movement-001',
    'process-001',
    '2023-12-01 14:30:00',
    'Distribuição - Distribuído por sorteio',
    'distribution',
    NOW(),
    NOW()
),
(
    'movement-002',
    'process-001',
    '2023-12-05 10:15:00',
    'Citação - Citação realizada via postal',
    'citation',
    NOW(),
    NOW()
),
(
    'movement-003',
    'process-001',
    '2023-12-20 16:45:00',
    'Petição - Embargos à execução apresentados',
    'petition',
    NOW(),
    NOW()
),
(
    'movement-004',
    'process-002',
    '2023-11-15 09:00:00',
    'Distribuição - Distribuído para a 5ª Vara do Trabalho',
    'distribution',
    NOW(),
    NOW()
),
(
    'movement-005',
    'process-002', 
    '2023-11-28 14:20:00',
    'Audiência - Designada audiência para 15/01/2024',
    'hearing',
    NOW(),
    NOW()
),
(
    'movement-006',
    'process-003',
    '2023-12-10 11:30:00',
    'Decisão - Determinada citação do requerido',
    'decision',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Prazos
INSERT INTO process.deadlines (
    id, process_id, date, description, type, priority,
    status, created_at, updated_at
) VALUES
(
    'deadline-001',
    'process-001',
    '2024-01-20',
    'Prazo para impugnação aos embargos à execução',
    'counter_argument',
    'high',
    'pending',
    NOW(),
    NOW()
),
(
    'deadline-002',
    'process-002',
    '2024-01-15',
    'Audiência de instrução e julgamento',
    'hearing',
    'critical',
    'pending',
    NOW(),
    NOW()
),
(
    'deadline-003',
    'process-003',
    '2024-01-10',
    'Prazo para contestação',
    'defense',
    'high',
    'pending',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Templates de notificação
INSERT INTO notification.templates (
    id, tenant_id, name, type, channel, subject, body,
    is_active, created_at, updated_at
) VALUES
(
    'template-001',
    'tenant-001',
    'Nova Movimentação',
    'process_movement',
    'whatsapp',
    '',
    'Olá {{client_name}}! Seu processo {{process_number}} teve uma nova movimentação: {{movement_description}}. Acesse o painel para mais detalhes.',
    true,
    NOW(),
    NOW()
),
(
    'template-002',
    'tenant-001',
    'Prazo se Aproximando',
    'deadline_reminder',
    'email',
    'Prazo importante - Processo {{process_number}}',
    'Prezado(a) {{client_name}},\n\nLembramos que o prazo "{{deadline_description}}" do processo {{process_number}} vence em {{days_remaining}} dias ({{deadline_date}}).\n\nAtenciosamente,\nEquipe Jurídica',
    true,
    NOW(),
    NOW()
),
(
    'template-003',
    NULL,
    'Quota Excedida',
    'quota_warning',
    'email',
    'Limite de quota atingido',
    'Sua quota de {{quota_type}} foi excedida. Atual: {{current_usage}}/{{limit}}. Considere fazer upgrade do plano.',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Histórico de consultas DataJud
INSERT INTO datajud.queries (
    id, tenant_id, process_number, query_type, status,
    response_time_ms, cache_hit, created_at, updated_at
) VALUES
(
    'query-001',
    'tenant-001',
    '1234567-89.2023.8.26.0001',
    'full',
    'success',
    1250,
    false,
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '2 hours'
),
(
    'query-002',
    'tenant-001',
    '1234567-89.2023.8.26.0001',
    'movements',
    'success',
    450,
    true,
    NOW() - INTERVAL '1 hour',
    NOW() - INTERVAL '1 hour'
),
(
    'query-003',
    'tenant-002',
    '5555555-55.2023.8.26.0100',
    'full',
    'success',
    2100,
    false,
    NOW() - INTERVAL '30 minutes',
    NOW() - INTERVAL '30 minutes'
)
ON CONFLICT (id) DO NOTHING;

-- Uso de quotas atual
INSERT INTO tenant.quota_usage (
    tenant_id, process_count, user_count, datajud_queries_today,
    storage_used_gb, notifications_sent_month, last_updated
) VALUES
(
    'tenant-001', 2, 3, 15, 2.5, 45, NOW()
),
(
    'tenant-002', 1, 1, 8, 1.2, 12, NOW()  
),
(
    'tenant-003', 1, 1, 3, 0.5, 5, NOW()
)
ON CONFLICT (tenant_id) DO UPDATE SET
    process_count = EXCLUDED.process_count,
    user_count = EXCLUDED.user_count,
    datajud_queries_today = EXCLUDED.datajud_queries_today,
    storage_used_gb = EXCLUDED.storage_used_gb,
    notifications_sent_month = EXCLUDED.notifications_sent_month,
    last_updated = EXCLUDED.last_updated;

-- Log de eventos para analytics
INSERT INTO analytics.events (
    id, tenant_id, event_type, entity_type, entity_id,
    metadata, occurred_at
) VALUES
(
    gen_random_uuid(),
    'tenant-001',
    'process_registered',
    'process',
    'process-001',
    '{"user_id": "user-001", "client_id": "client-001"}',
    NOW() - INTERVAL '7 days'
),
(
    gen_random_uuid(),
    'tenant-001',
    'movement_detected',
    'process',
    'process-001',
    '{"movement_id": "movement-003", "movement_type": "petition"}',
    NOW() - INTERVAL '1 day'
),
(
    gen_random_uuid(),
    'tenant-001',
    'notification_sent',
    'notification',
    'notification-001',
    '{"channel": "whatsapp", "success": true}',
    NOW() - INTERVAL '1 day'
),
(
    gen_random_uuid(),
    'tenant-002',
    'user_login',
    'user',
    'user-003',
    '{"ip_address": "192.168.1.100", "user_agent": "Mozilla/5.0"}',
    NOW() - INTERVAL '2 hours'
),
(
    gen_random_uuid(),
    'tenant-003',
    'quota_warning',
    'tenant',
    'tenant-003',
    '{"quota_type": "processes", "usage": 45, "limit": 50}',
    NOW() - INTERVAL '1 hour'
)
ON CONFLICT (id) DO NOTHING;

COMMIT;

EOF

if [ $? -eq 0 ]; then
    log_success "Dados de exemplo criados com sucesso! ✓"
else
    log_error "Erro ao criar dados de exemplo"
    exit 1
fi

echo ""
echo -e "${GREEN}📊 DADOS CRIADOS:${NC}"
echo ""
echo -e "${BLUE}👥 Tenants:${NC}"
echo "  • Silva & Associados (Professional) - tenant-001"
echo "  • Dept. Jurídico XYZ (Business) - tenant-002"  
echo "  • Dr. João Autônomo (Starter) - tenant-003"
echo ""
echo -e "${BLUE}👤 Usuários:${NC}"
echo "  • admin@silva-advogados.com (Admin)"
echo "  • maria@silva-advogados.com (Advogada)"
echo "  • juridico@empresaxyz.com (Admin)"
echo "  • joao@advogado.com (Admin)"
echo "  • cliente1@empresa.com (Cliente)"
echo ""
echo -e "${BLUE}⚖️  Processos:${NC}"
echo "  • 4 processos de exemplo"
echo "  • 6 movimentações"  
echo "  • 3 prazos pendentes"
echo ""
echo -e "${BLUE}📧 Templates:${NC}"
echo "  • Template WhatsApp para movimentações"
echo "  • Template email para prazos"
echo "  • Template para alertas de quota"
echo ""
echo -e "${BLUE}📈 Analytics:${NC}"
echo "  • Eventos de exemplo dos últimos 7 dias"
echo "  • Métricas de uso de quota"
echo "  • Histórico de consultas DataJud"
echo ""
echo -e "${YELLOW}🔗 Para testar via API (quando os serviços estiverem implementados):${NC}"
echo ""
echo -e "${BLUE}# Listar processos do tenant Silva & Associados${NC}"
echo "curl -H 'X-Tenant-ID: tenant-001' http://localhost:8000/api/v1/processes"
echo ""
echo -e "${BLUE}# Buscar processo específico${NC}"
echo "curl -H 'X-Tenant-ID: tenant-001' http://localhost:8000/api/v1/processes/process-001"
echo ""
echo -e "${BLUE}# Consultar quota atual${NC}"
echo "curl -H 'X-Tenant-ID: tenant-001' http://localhost:8000/api/v1/tenants/tenant-001/quotas"
echo ""