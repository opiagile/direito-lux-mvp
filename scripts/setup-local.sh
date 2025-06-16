#!/bin/bash

# =============================================================================
# Direito Lux - Setup Ambiente Local
# =============================================================================

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funções auxiliares
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

# Banner
echo -e "${BLUE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                        DIREITO LUX                          ║
║                   Setup Ambiente Local                      ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Verificar pré-requisitos
log_info "Verificando pré-requisitos..."

# Verificar Docker
if ! command -v docker &> /dev/null; then
    log_error "Docker não está instalado. Instale Docker Desktop primeiro."
    exit 1
fi

# Verificar Docker Compose
if ! command -v docker-compose &> /dev/null; then
    log_error "Docker Compose não está instalado."
    exit 1
fi

# Verificar se Docker está rodando
if ! docker info &> /dev/null; then
    log_error "Docker não está rodando. Inicie o Docker Desktop."
    exit 1
fi

log_success "Pré-requisitos verificados ✓"

# Criar arquivo .env se não existir
if [ ! -f .env ]; then
    log_info "Criando arquivo .env a partir do .env.example..."
    cp .env.example .env
    log_warning "IMPORTANTE: Edite o arquivo .env com suas configurações reais"
    log_warning "Especialmente as chaves de API (DataJud, OpenAI, WhatsApp)"
fi

# Criar diretórios necessários
log_info "Criando diretórios de infraestrutura..."

mkdir -p infrastructure/{sql/init,rabbitmq,kong,prometheus,grafana/{datasources,dashboards},pgadmin,mocks/whatsapp,keycloak}
mkdir -p services/{auth-service,tenant-service,process-service,datajud-service,notification-service,ai-service}
mkdir -p storage/{documents,temp}
mkdir -p logs

log_success "Diretórios criados ✓"

# Configurar RabbitMQ
log_info "Configurando RabbitMQ..."

cat > infrastructure/rabbitmq/rabbitmq.conf << EOF
management.load_definitions = /etc/rabbitmq/definitions.json
loopback_users.guest = false
listeners.tcp.default = 5672
management.tcp.port = 15672
management.tcp.ip = 0.0.0.0
log.file.level = info
EOF

cat > infrastructure/rabbitmq/definitions.json << EOF
{
  "rabbit_version": "3.12.0",
  "rabbitmq_version": "3.12.0",
  "product_name": "RabbitMQ",
  "product_version": "3.12.0",
  "users": [
    {
      "name": "direito_lux",
      "password_hash": "gTgMIZ6E5rT8M+w8M7V3TqE6G4V2QZ8m",
      "hashing_algorithm": "rabbit_password_hashing_sha256",
      "tags": ["administrator"]
    }
  ],
  "vhosts": [
    {
      "name": "direito_lux"
    }
  ],
  "permissions": [
    {
      "user": "direito_lux",
      "vhost": "direito_lux",
      "configure": ".*",
      "write": ".*",
      "read": ".*"
    }
  ],
  "exchanges": [
    {
      "name": "direito_lux.events",
      "vhost": "direito_lux",
      "type": "topic",
      "durable": true,
      "auto_delete": false,
      "internal": false,
      "arguments": {}
    },
    {
      "name": "direito_lux.dlx",
      "vhost": "direito_lux",
      "type": "direct",
      "durable": true,
      "auto_delete": false,
      "internal": false,
      "arguments": {}
    }
  ],
  "queues": [
    {
      "name": "auth.events",
      "vhost": "direito_lux",
      "durable": true,
      "auto_delete": false,
      "arguments": {
        "x-dead-letter-exchange": "direito_lux.dlx",
        "x-dead-letter-routing-key": "auth.events.dlq"
      }
    },
    {
      "name": "tenant.events",
      "vhost": "direito_lux",
      "durable": true,
      "auto_delete": false,
      "arguments": {
        "x-dead-letter-exchange": "direito_lux.dlx"
      }
    },
    {
      "name": "process.events",
      "vhost": "direito_lux",
      "durable": true,
      "auto_delete": false,
      "arguments": {
        "x-dead-letter-exchange": "direito_lux.dlx"
      }
    },
    {
      "name": "datajud.events",
      "vhost": "direito_lux",
      "durable": true,
      "auto_delete": false,
      "arguments": {
        "x-dead-letter-exchange": "direito_lux.dlx"
      }
    },
    {
      "name": "notification.events",
      "vhost": "direito_lux",
      "durable": true,
      "auto_delete": false,
      "arguments": {
        "x-dead-letter-exchange": "direito_lux.dlx"
      }
    },
    {
      "name": "ai.events",
      "vhost": "direito_lux",
      "durable": true,
      "auto_delete": false,
      "arguments": {
        "x-dead-letter-exchange": "direito_lux.dlx"
      }
    }
  ],
  "bindings": [
    {
      "source": "direito_lux.events",
      "vhost": "direito_lux",
      "destination": "auth.events",
      "destination_type": "queue",
      "routing_key": "auth.*",
      "arguments": {}
    },
    {
      "source": "direito_lux.events",
      "vhost": "direito_lux",
      "destination": "tenant.events",
      "destination_type": "queue",
      "routing_key": "tenant.*",
      "arguments": {}
    },
    {
      "source": "direito_lux.events",
      "vhost": "direito_lux",
      "destination": "process.events",
      "destination_type": "queue",
      "routing_key": "process.*",
      "arguments": {}
    },
    {
      "source": "direito_lux.events",
      "vhost": "direito_lux",
      "destination": "datajud.events",
      "destination_type": "queue",
      "routing_key": "datajud.*",
      "arguments": {}
    },
    {
      "source": "direito_lux.events",
      "vhost": "direito_lux",
      "destination": "notification.events",
      "destination_type": "queue",
      "routing_key": "notification.*",
      "arguments": {}
    },
    {
      "source": "direito_lux.events",
      "vhost": "direito_lux",
      "destination": "ai.events",
      "destination_type": "queue",
      "routing_key": "ai.*",
      "arguments": {}
    }
  ]
}
EOF

log_success "RabbitMQ configurado ✓"

# Configurar Kong API Gateway
log_info "Configurando Kong API Gateway..."

cat > infrastructure/kong/kong.yml << EOF
_format_version: "3.0"
_transform: true

services:
  - name: auth-service
    url: http://auth-service:8080
    plugins:
      - name: cors
        config:
          origins:
            - http://localhost:3000
            - http://localhost:3001
          methods:
            - GET
            - POST
            - PUT
            - DELETE
            - OPTIONS
          headers:
            - Accept
            - Accept-Version
            - Content-Length
            - Content-MD5
            - Content-Type
            - Date
            - Authorization
            - X-Tenant-ID
          credentials: true

  - name: tenant-service
    url: http://tenant-service:8080

  - name: process-service
    url: http://process-service:8080

  - name: datajud-service
    url: http://datajud-service:8080

  - name: notification-service
    url: http://notification-service:8080

  - name: ai-service
    url: http://ai-service:8000

routes:
  - name: auth-routes
    service: auth-service
    paths:
      - /api/v1/auth

  - name: tenant-routes
    service: tenant-service
    paths:
      - /api/v1/tenants

  - name: process-routes
    service: process-service
    paths:
      - /api/v1/processes

  - name: datajud-routes
    service: datajud-service
    paths:
      - /api/v1/datajud

  - name: notification-routes
    service: notification-service
    paths:
      - /api/v1/notifications

  - name: ai-routes
    service: ai-service
    paths:
      - /api/v1/ai

plugins:
  - name: rate-limiting
    config:
      minute: 100
      hour: 1000
      policy: local

  - name: request-id
    config:
      header_name: X-Request-ID
      generator: uuid

  - name: correlation-id
    config:
      header_name: X-Correlation-ID
      generator: uuid
EOF

log_success "Kong configurado ✓"

# Configurar Prometheus
log_info "Configurando Prometheus..."

cat > infrastructure/prometheus/prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'auth-service'
    static_configs:
      - targets: ['auth-service:9090']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'tenant-service'
    static_configs:
      - targets: ['tenant-service:9090']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'process-service'
    static_configs:
      - targets: ['process-service:9090']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'datajud-service'
    static_configs:
      - targets: ['datajud-service:9090']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'notification-service'
    static_configs:
      - targets: ['notification-service:9090']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'ai-service'
    static_configs:
      - targets: ['ai-service:9090']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'postgres-exporter'
    static_configs:
      - targets: ['postgres:9187']

  - job_name: 'redis-exporter'
    static_configs:
      - targets: ['redis:9121']

  - job_name: 'rabbitmq'
    static_configs:
      - targets: ['rabbitmq:15692']
EOF

log_success "Prometheus configurado ✓"

# Configurar Grafana datasources
log_info "Configurando Grafana..."

cat > infrastructure/grafana/datasources/prometheus.yml << EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
EOF

mkdir -p infrastructure/grafana/dashboards
cat > infrastructure/grafana/dashboards/dashboard.yml << EOF
apiVersion: 1

providers:
  - name: 'default'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
EOF

log_success "Grafana configurado ✓"

# Configurar pgAdmin
log_info "Configurando pgAdmin..."

cat > infrastructure/pgadmin/servers.json << EOF
{
    "Servers": {
        "1": {
            "Name": "Direito Lux Local",
            "Group": "Servers",
            "Host": "postgres",
            "Port": 5432,
            "MaintenanceDB": "direito_lux_dev",
            "Username": "direito_lux",
            "UseSSHTunnel": 0,
            "TunnelPort": "22",
            "TunnelAuthentication": 0
        }
    }
}
EOF

log_success "pgAdmin configurado ✓"

# Configurar Mock WhatsApp API
log_info "Configurando Mock WhatsApp API..."

cat > infrastructure/mocks/whatsapp/mappings/send-message.json << EOF
{
    "request": {
        "method": "POST",
        "urlPattern": "/v18.0/[0-9]+/messages"
    },
    "response": {
        "status": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "jsonBody": {
            "messaging_product": "whatsapp",
            "contacts": [
                {
                    "input": "{{request.body}}",
                    "wa_id": "5511999999999"
                }
            ],
            "messages": [
                {
                    "id": "wamid.{{randomValue length=32 type='ALPHANUMERIC'}}"
                }
            ]
        }
    }
}
EOF

cat > infrastructure/mocks/whatsapp/mappings/webhook.json << EOF
{
    "request": {
        "method": "GET",
        "urlPath": "/webhook"
    },
    "response": {
        "status": 200,
        "headers": {
            "Content-Type": "text/plain"
        },
        "body": "{{request.query.hub.challenge}}"
    }
}
EOF

log_success "Mock WhatsApp API configurado ✓"

# Configurar SQL inicial
log_info "Configurando banco de dados inicial..."

cat > infrastructure/sql/init/01-init.sql << EOF
-- Direito Lux - Database Initialization

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Create schemas por bounded context
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS tenant;
CREATE SCHEMA IF NOT EXISTS process;
CREATE SCHEMA IF NOT EXISTS datajud;
CREATE SCHEMA IF NOT EXISTS notification;
CREATE SCHEMA IF NOT EXISTS ai;
CREATE SCHEMA IF NOT EXISTS document;
CREATE SCHEMA IF NOT EXISTS analytics;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA tenant TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA process TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA datajud TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA notification TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA ai TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA document TO direito_lux;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA analytics TO direito_lux;

-- Create tenant isolation function
CREATE OR REPLACE FUNCTION set_current_tenant_id(tenant_id TEXT)
RETURNS VOID AS \$\$
BEGIN
    PERFORM set_config('app.current_tenant_id', tenant_id, true);
END;
\$\$ LANGUAGE plpgsql;

-- Create audit trigger function
CREATE OR REPLACE FUNCTION audit_trigger()
RETURNS TRIGGER AS \$\$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.created_at = NOW();
        NEW.updated_at = NOW();
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        NEW.updated_at = NOW();
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
\$\$ LANGUAGE plpgsql;

-- Log successful initialization
INSERT INTO public.migration_log (version, description, executed_at) 
VALUES ('001', 'Initial database setup', NOW())
ON CONFLICT (version) DO NOTHING;

-- Create migration log table if not exists
CREATE TABLE IF NOT EXISTS public.migration_log (
    version VARCHAR(10) PRIMARY KEY,
    description TEXT,
    executed_at TIMESTAMP DEFAULT NOW()
);
EOF

log_success "SQL inicial configurado ✓"

# Parar containers existentes se houver
log_info "Parando containers existentes..."
docker-compose down -v 2>/dev/null || true

# Build das imagens (se os Dockerfiles existirem)
log_info "Verificando se existem serviços para build..."
for service in auth-service tenant-service process-service datajud-service notification-service ai-service; do
    if [ -d "services/$service" ]; then
        log_info "Diretório services/$service existe"
    else
        log_warning "Diretório services/$service não existe - será criado quando implementarmos o serviço"
    fi
done

# Subir infraestrutura primeiro
log_info "Subindo infraestrutura (PostgreSQL, Redis, RabbitMQ, Keycloak)..."

docker-compose up -d postgres redis rabbitmq keycloak

# Aguardar serviços ficarem prontos
log_info "Aguardando serviços ficarem prontos..."

# Aguardar PostgreSQL
log_info "Aguardando PostgreSQL..."
while ! docker-compose exec -T postgres pg_isready -U direito_lux 2>/dev/null; do
    sleep 2
done

# Aguardar Redis
log_info "Aguardando Redis..."
while ! docker-compose exec -T redis redis-cli --no-auth-warning -a dev_redis_123 ping 2>/dev/null | grep -q PONG; do
    sleep 2
done

# Aguardar RabbitMQ
log_info "Aguardando RabbitMQ..."
while ! docker-compose exec -T rabbitmq rabbitmq-diagnostics check_port_connectivity 2>/dev/null; do
    sleep 5
done

log_success "Infraestrutura rodando ✓"

# Subir observabilidade
log_info "Subindo stack de observabilidade..."
docker-compose up -d jaeger prometheus grafana

# Subir ferramentas de desenvolvimento
log_info "Subindo ferramentas de desenvolvimento..."
docker-compose up -d pgadmin redis-commander mailhog whatsapp-mock

# Subir API Gateway
log_info "Subindo API Gateway..."
docker-compose up -d kong

log_success "Ambiente local configurado com sucesso! ✓"

# Mostrar URLs importantes
echo ""
echo -e "${GREEN}🎉 AMBIENTE PRONTO! Acesse os serviços:${NC}"
echo ""
echo -e "${BLUE}📊 Dashboards e Monitoramento:${NC}"
echo "  • Grafana:           http://localhost:3000 (admin/dev_grafana_123)"
echo "  • Prometheus:        http://localhost:9090"
echo "  • Jaeger:            http://localhost:16686"
echo ""
echo -e "${BLUE}💾 Gerenciamento de Dados:${NC}"
echo "  • pgAdmin:           http://localhost:5050 (admin@direitolux.com/dev_pgadmin_123)"
echo "  • Redis Commander:   http://localhost:8081 (admin/dev_redis_ui_123)"
echo ""
echo -e "${BLUE}🔧 Infraestrutura:${NC}"
echo "  • RabbitMQ:          http://localhost:15672 (direito_lux/dev_rabbit_123)"
echo "  • Keycloak:          http://localhost:8080 (admin/dev_admin_123)"
echo "  • Kong Admin:        http://localhost:8002"
echo ""
echo -e "${BLUE}🧪 Desenvolvimento:${NC}"
echo "  • MailHog:           http://localhost:8025"
echo "  • WhatsApp Mock:     http://localhost:9000"
echo ""
echo -e "${BLUE}🚀 APIs:${NC}"
echo "  • API Gateway:       http://localhost:8000"
echo "  • Auth Service:      http://localhost:8081 (quando implementado)"
echo "  • Tenant Service:    http://localhost:8082 (quando implementado)"
echo "  • Process Service:   http://localhost:8083 (quando implementado)"
echo ""
echo -e "${YELLOW}⚠️  PRÓXIMOS PASSOS:${NC}"
echo "  1. Edite o arquivo .env com suas chaves de API reais"
echo "  2. Execute: ${BLUE}./scripts/create-sample-data.sh${NC} para dados de teste"
echo "  3. Comece implementando o Auth Service"
echo ""
echo -e "${GREEN}Para parar o ambiente: ${BLUE}docker-compose down${NC}"
echo -e "${GREEN}Para ver logs: ${BLUE}docker-compose logs -f [service_name]${NC}"
echo ""