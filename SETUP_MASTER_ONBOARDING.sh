#!/bin/bash

echo "🚀 DIREITO LUX - SETUP MASTER PARA ONBOARDING"
echo "============================================="
echo "Este script configura o ambiente completo do zero"
echo "Ideal para novos desenvolvedores e demonstrações"
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para logging
log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }

# Função para verificar dependências
check_dependencies() {
    log_info "Verificando dependências..."
    
    # Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker não está instalado!"
        echo "Instale Docker Desktop: https://docs.docker.com/desktop/"
        exit 1
    fi
    log_success "Docker encontrado"
    
    # Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "docker-compose não está instalado!"
        echo "Instale docker-compose: https://docs.docker.com/compose/install/"
        exit 1
    fi
    log_success "docker-compose encontrado"
    
    # PostgreSQL client
    if ! command -v psql &> /dev/null; then
        log_warning "psql não encontrado. Instalando..."
        brew install postgresql
    fi
    log_success "psql disponível"
    
    # golang-migrate
    if ! command -v migrate &> /dev/null; then
        log_warning "golang-migrate não encontrado. Instalando..."
        brew install golang-migrate
    fi
    log_success "golang-migrate disponível: $(migrate -version)"
}

# Função para limpar ambiente
clean_environment() {
    log_info "Limpando ambiente anterior..."
    
    # Parar todos os containers
    docker-compose down -v --remove-orphans 2>/dev/null || true
    
    # Remover volumes órfãos
    docker volume prune -f 2>/dev/null || true
    
    # Aguardar containers pararem
    sleep 3
    
    log_success "Ambiente limpo"
}

# Função para iniciar PostgreSQL
start_postgresql() {
    log_info "Iniciando PostgreSQL..."
    
    # Subir apenas PostgreSQL
    docker-compose up -d postgres
    
    # Aguardar PostgreSQL estar pronto
    log_info "Aguardando PostgreSQL inicializar..."
    attempts=0
    max_attempts=30
    
    while [ $attempts -lt $max_attempts ]; do
        if docker exec direito-lux-postgres pg_isready -U postgres &>/dev/null; then
            log_success "PostgreSQL está pronto!"
            break
        fi
        attempts=$((attempts + 1))
        echo -n "."
        sleep 2
    done
    
    if [ $attempts -eq $max_attempts ]; then
        log_error "PostgreSQL não iniciou após ${max_attempts} tentativas"
        docker logs direito-lux-postgres
        exit 1
    fi
    
    # Aguardar scripts de inicialização
    log_info "Aguardando scripts de inicialização (10s)..."
    sleep 10
}

# Função para verificar e criar usuário
setup_database_user() {
    log_info "Configurando usuário e banco..."
    
    # Testar conexão com direito_lux
    if PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -c "SELECT 1;" &>/dev/null; then
        log_success "Usuário direito_lux já existe e funciona"
        return 0
    fi
    
    # Criar usuário manualmente
    log_warning "Criando usuário direito_lux manualmente..."
    PGPASSWORD=postgres psql -h localhost -U postgres << 'EOF'
-- Criar role direito_lux
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'direito_lux') THEN
        CREATE ROLE direito_lux WITH LOGIN PASSWORD 'dev_password_123' CREATEDB SUPERUSER;
    END IF;
END $$;

-- Criar database
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'direito_lux_dev') THEN
        CREATE DATABASE direito_lux_dev OWNER direito_lux;
    END IF;
END $$;

-- Conectar ao database e criar extensões
\c direito_lux_dev
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

SELECT 'Usuario e banco configurados!' as status;
EOF
    
    # Verificar novamente
    if PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -c "SELECT 'Conexão OK!' as status;" &>/dev/null; then
        log_success "Usuário e banco configurados com sucesso"
    else
        log_error "Falha ao configurar usuário"
        exit 1
    fi
}

# Função para executar migrations via SQL direto
execute_migrations() {
    log_info "Executando migrations via SQL direto..."
    
    # Limpar schema_migrations
    PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev << 'EOF'
DROP TABLE IF EXISTS schema_migrations CASCADE;
EOF
    
    # Executar SQL do Tenant Service
    log_info "Criando tabelas do Tenant Service..."
    for sql_file in services/tenant-service/migrations/*.up.sql; do
        if [ -f "$sql_file" ]; then
            log_info "Executando: $(basename $sql_file)"
            PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev < "$sql_file" 2>/dev/null || {
                log_warning "Erro em $sql_file (pode ser esperado se tabela já existe)"
            }
        fi
    done
    
    # Executar SQL do Auth Service
    log_info "Criando tabelas do Auth Service..."
    for sql_file in services/auth-service/migrations/*.up.sql; do
        if [ -f "$sql_file" ]; then
            filename=$(basename "$sql_file")
            if [[ ! "$filename" =~ seed ]]; then # Pular arquivos de seed por enquanto
                log_info "Executando: $filename"
                PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev < "$sql_file" 2>/dev/null || {
                    log_warning "Erro em $sql_file (pode ser esperado se tabela já existe)"
                }
            fi
        fi
    done
    
    # Executar SQL do Process Service
    log_info "Criando tabelas do Process Service..."
    for sql_file in services/process-service/migrations/*.up.sql; do
        if [ -f "$sql_file" ]; then
            filename=$(basename "$sql_file")
            if [[ ! "$filename" =~ seed ]]; then # Pular arquivos de seed por enquanto
                log_info "Executando: $filename"
                PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev < "$sql_file" 2>/dev/null || {
                    log_warning "Erro em $sql_file (pode ser esperado se tabela já existe)"
                }
            fi
        fi
    done
}

# Função para inserir dados de teste
insert_test_data() {
    log_info "Inserindo dados de teste compatíveis..."
    
    PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev << 'EOF'
-- Limpar dados existentes
TRUNCATE TABLE tenants CASCADE;
TRUNCATE TABLE users CASCADE;

-- Inserir 8 tenants (2 por plano)
INSERT INTO tenants (id, name, legal_name, document, email, phone, plan_type, owner_user_id, status, address) VALUES
-- Starter
('11111111-1111-1111-1111-111111111111', 'Silva & Associados', 'Silva & Associados Advogados Ltda', '11111111000111', 'admin@silvaassociados.com.br', '11987654321', 'starter', '11111111-1111-1111-1111-111111111111', 'active', '{"street": "Rua das Flores", "number": "123", "city": "São Paulo", "state": "SP"}'),
('22222222-2222-2222-2222-222222222222', 'Lima Advogados', 'Lima Sociedade de Advogados', '22222222000122', 'admin@limaadvogados.com.br', '21987654321', 'starter', '22222222-2222-2222-2222-222222222222', 'active', '{"street": "Av. Atlântica", "number": "456", "city": "Rio de Janeiro", "state": "RJ"}'),
-- Professional  
('33333333-3333-3333-3333-333333333333', 'Costa Santos', 'Costa Santos Advogados', '33333333000133', 'admin@costasantos.com.br', '31987654321', 'professional', '33333333-3333-3333-3333-333333333333', 'active', '{"street": "Rua da Bahia", "number": "789", "city": "Belo Horizonte", "state": "MG"}'),
('44444444-4444-4444-4444-444444444444', 'Pereira Oliveira', 'Pereira Oliveira Advocacia', '44444444000144', 'admin@pereiraoliveira.com.br', '41987654321', 'professional', '44444444-4444-4444-4444-444444444444', 'active', '{"street": "Rua XV", "number": "321", "city": "Curitiba", "state": "PR"}'),
-- Business
('55555555-5555-5555-5555-555555555555', 'Machado Advogados', 'Machado & Machado S/S', '55555555000155', 'admin@machadoadvogados.com.br', '51987654321', 'business', '55555555-5555-5555-5555-555555555555', 'active', '{"street": "Av. Borges", "number": "654", "city": "Porto Alegre", "state": "RS"}'),
('66666666-6666-6666-6666-666666666666', 'Ferreira Legal', 'Ferreira Advogados', '66666666000166', 'admin@ferreiralegal.com.br', '61987654321', 'business', '66666666-6666-6666-6666-666666666666', 'active', '{"street": "SAS Quadra 1", "number": "Bloco A", "city": "Brasília", "state": "DF"}'),
-- Enterprise
('77777777-7777-7777-7777-777777777777', 'Barros Enterprise', 'Barros Internacional', '77777777000177', 'admin@barrosent.com.br', '11912345678', 'enterprise', '77777777-7777-7777-7777-777777777777', 'active', '{"street": "Av. Faria Lima", "number": "1000", "city": "São Paulo", "state": "SP"}'),
('88888888-8888-8888-8888-888888888888', 'Rodrigues Global', 'Rodrigues Advogados', '88888888000188', 'admin@rodriguesglobal.com.br', '21912345678', 'enterprise', '88888888-8888-8888-8888-888888888888', 'active', '{"street": "Praia de Botafogo", "number": "300", "city": "Rio de Janeiro", "state": "RJ"}');

-- Inserir usuários (32 total - 4 por tenant)
DO $$
DECLARE
    tenant_rec RECORD;
    tenant_domain TEXT;
BEGIN
    FOR tenant_rec IN SELECT id, name, email FROM tenants LOOP
        -- Extrair domínio do email
        tenant_domain := '@' || SUBSTRING(tenant_rec.email FROM '@(.*)$');
        
        -- Admin (mesmo ID do owner_user_id)
        INSERT INTO users (id, email, password_hash, role, tenant_id, status, first_name, last_name)
        VALUES (
            tenant_rec.id, -- Usar mesmo ID para ser o owner
            'admin' || tenant_domain,
            '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- hash para 'password'
            'admin',
            tenant_rec.id,
            'active',
            'Admin',
            COALESCE(NULLIF(SPLIT_PART(tenant_rec.name, ' ', 2), ''), 'User')
        ) ON CONFLICT (id) DO UPDATE SET status = 'active';
        
        -- Manager
        INSERT INTO users (id, email, password_hash, role, tenant_id, status, first_name, last_name)
        VALUES (
            gen_random_uuid(),
            'gerente' || tenant_domain,
            '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
            'manager',
            tenant_rec.id,
            'active',
            'Gerente',
            COALESCE(NULLIF(SPLIT_PART(tenant_rec.name, ' ', 2), ''), 'User')
        );
        
        -- Operator (usando role existente)
        INSERT INTO users (id, email, password_hash, role, tenant_id, status, first_name, last_name)
        VALUES (
            gen_random_uuid(),
            'advogado' || tenant_domain,
            '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
            'operator',
            tenant_rec.id,
            'active',
            'Advogado',
            COALESCE(NULLIF(SPLIT_PART(tenant_rec.name, ' ', 2), ''), 'User')
        );
        
        -- Client (usando role existente)
        INSERT INTO users (id, email, password_hash, role, tenant_id, status, first_name, last_name)
        VALUES (
            gen_random_uuid(),
            'assistente' || tenant_domain,
            '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
            'client',
            tenant_rec.id,
            'active',
            'Assistente',
            COALESCE(NULLIF(SPLIT_PART(tenant_rec.name, ' ', 2), ''), 'User')
        );
    END LOOP;
END $$;

-- Criar tabela de processos se não existir
CREATE TABLE IF NOT EXISTS processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    process_number VARCHAR(50) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    type VARCHAR(100),
    status VARCHAR(50) DEFAULT 'active',
    priority VARCHAR(20) DEFAULT 'medium',
    value NUMERIC(15,2),
    client_name VARCHAR(255),
    responsible_lawyer_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT processes_tenant_number_unique UNIQUE(tenant_id, process_number)
);

-- Criar índices se não existirem
CREATE INDEX IF NOT EXISTS idx_processes_tenant_id ON processes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_processes_status ON processes(status);

-- Inserir processos de teste
TRUNCATE TABLE processes CASCADE;

DO $$
DECLARE
    tenant_rec RECORD;
    process_count INTEGER;
    i INTEGER;
    process_types TEXT[] := ARRAY['Cível', 'Trabalhista', 'Criminal', 'Família', 'Tributário', 'Consumidor'];
    process_status TEXT[] := ARRAY['active', 'concluded', 'suspended'];
    process_priority TEXT[] := ARRAY['low', 'medium', 'high', 'urgent'];
BEGIN
    FOR tenant_rec IN SELECT id, name, plan_type FROM tenants LOOP
        -- Definir quantidade de processos baseado no plano
        CASE tenant_rec.plan_type
            WHEN 'starter' THEN process_count := 5;
            WHEN 'professional' THEN process_count := 10;
            WHEN 'business' THEN process_count := 15;
            WHEN 'enterprise' THEN process_count := 20;
        END CASE;
        
        -- Criar processos
        FOR i IN 1..process_count LOOP
            INSERT INTO processes (
                tenant_id,
                process_number,
                subject,
                type,
                status,
                priority,
                value,
                client_name
            ) VALUES (
                tenant_rec.id,
                'PROC-' || TO_CHAR(NOW(), 'YYYY') || '-' || LPAD(i::TEXT, 6, '0'),
                'Processo ' || i || ' - ' || process_types[1 + (i % array_length(process_types, 1))],
                process_types[1 + (i % array_length(process_types, 1))],
                process_status[1 + (i % array_length(process_status, 1))],
                process_priority[1 + (i % array_length(process_priority, 1))],
                ROUND((RANDOM() * 100000 + 1000)::NUMERIC, 2),
                'Cliente ' || i
            );
        END LOOP;
    END LOOP;
END $$;

SELECT 'Dados de teste inseridos com sucesso!' as status;
EOF
}

# Função para verificar resultado final
verify_setup() {
    log_info "Verificando setup final..."
    
    local result
    result=$(PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -t << 'EOF'
SELECT 
    'Tenants: ' || COUNT(*) FROM tenants
UNION ALL
SELECT 
    'Users: ' || COUNT(*) FROM users  
UNION ALL
SELECT 
    'Processes: ' || COUNT(*) FROM processes;
EOF
)
    
    echo "$result"
    
    # Verificar números esperados
    local tenants=$(echo "$result" | grep "Tenants:" | grep -o '[0-9]\+')
    local users=$(echo "$result" | grep "Users:" | grep -o '[0-9]\+')
    local processes=$(echo "$result" | grep "Processes:" | grep -o '[0-9]\+')
    
    if [ "$tenants" = "8" ] && [ "$users" = "32" ] && [ "$processes" -ge "90" ]; then
        log_success "Setup verificado com sucesso!"
        return 0
    else
        log_warning "Números não conferem: Tenants=$tenants, Users=$users, Processes=$processes"
        return 1
    fi
}

# Função para mostrar credenciais
show_credentials() {
    log_info "Credenciais de teste (senha para todos: password):"
    echo ""
    
    PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev << 'EOF'
SELECT 
    '🔑 ' || u.email as "Email para Login",
    t.name as "Tenant",
    t.plan_type as "Plano"
FROM users u
JOIN tenants t ON u.tenant_id = t.id
WHERE u.role = 'admin'
ORDER BY t.plan_type, t.name;
EOF
}

# Função principal
main() {
    echo "Iniciando setup master em 3 segundos..."
    sleep 3
    
    check_dependencies
    clean_environment
    start_postgresql
    setup_database_user
    execute_migrations
    insert_test_data
    
    if verify_setup; then
        echo ""
        log_success "🎊 SETUP MASTER CONCLUÍDO COM SUCESSO!"
        echo ""
        show_credentials
        echo ""
        log_info "🚀 PRÓXIMOS PASSOS:"
        echo "   1. docker-compose up -d    # Subir todos os serviços"
        echo "   2. cd frontend && npm install && npm run dev   # Iniciar frontend"
        echo "   3. Acessar http://localhost:3000"
        echo "   4. Fazer login com qualquer email admin acima"
        echo ""
        log_info "📚 Documentação: cat DOCUMENTO_TESTE_VALIDACAO.md"
    else
        log_error "Setup falhou na verificação final"
        exit 1
    fi
}

# Executar
main "$@"