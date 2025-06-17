# Migrações do Process Service

Este diretório contém as migrações do banco de dados para o Process Service do Direito Lux.

## Estrutura das Migrações

### 001_create_processes_table.sql
- **Objetivo**: Criar tabela principal de processos jurídicos
- **Recursos**:
  - Validação automática de número CNJ
  - Trigger para atualização automática de `updated_at`
  - Constraints de integridade
  - Campos JSONB para dados flexíveis (subject, value, monitoring, custom_fields)
  - Suporte a tags e metadados

### 002_create_movements_table.sql  
- **Objetivo**: Criar tabela de movimentações/andamentos
- **Recursos**:
  - Sequência automática por processo
  - Validação de tipos de movimentação
  - Metadados com análise de IA
  - Suporte a anexos (JSONB)
  - Triggers para auto-incremento de sequência

### 003_create_parties_table.sql
- **Objetivo**: Criar tabela de partes dos processos
- **Recursos**:
  - Validação automática de CPF/CNPJ
  - Dados do advogado em JSONB
  - Informações de contato e endereço
  - Constraints de consistência tipo/documento

### 004_create_indexes.sql
- **Objetivo**: Criar índices para performance
- **Recursos**:
  - Índices para queries frequentes
  - Índices compostos para relatórios
  - Índices GIN para busca textual
  - Índices em campos JSONB

### 005_create_functions_and_triggers.sql
- **Objetivo**: Funções de negócio e triggers
- **Recursos**:
  - Formatação e validação CNJ
  - Extração automática de palavras-chave
  - Detecção automática de movimentações importantes
  - Funções de estatísticas
  - Funções de manutenção

### 006_seed_initial_data.sql
- **Objetivo**: Dados iniciais e exemplos
- **Recursos**:
  - Configurações do sistema
  - Dados de exemplo para desenvolvimento
  - Views úteis para desenvolvimento
  - Templates de notificação

## Como Executar as Migrações

### Método 1: Usando o DatabaseMigrator (Recomendado)

```go
package main

import (
    "log"
    "github.com/direito-lux/process-service/internal/infrastructure/config"
)

func main() {
    // Carregar configuração
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Erro ao carregar configuração:", err)
    }

    // Configurar dependências
    deps, err := config.NewDependencies(cfg)
    if err != nil {
        log.Fatal("Erro ao configurar dependências:", err)
    }
    defer deps.Close()

    // Executar migrações
    migrator := config.NewDatabaseMigrator(deps.DB)
    if err := migrator.Migrate(); err != nil {
        log.Fatal("Erro nas migrações:", err)
    }

    log.Println("Migrações executadas com sucesso!")
}
```

### Método 2: Executando Manualmente

```bash
# Conectar ao PostgreSQL
psql -h localhost -U postgres -d process_service

# Executar migrações em ordem
\i migrations/001_create_processes_table.sql
\i migrations/002_create_movements_table.sql
\i migrations/003_create_parties_table.sql
\i migrations/004_create_indexes.sql
\i migrations/005_create_functions_and_triggers.sql
\i migrations/006_seed_initial_data.sql
```

### Método 3: Usando Docker

```bash
# Com volume das migrações
docker run -it --rm \
  -v $(pwd)/migrations:/migrations \
  postgres:15 \
  psql -h host.docker.internal -U postgres -d process_service \
  -f /migrations/001_create_processes_table.sql
```

## Variáveis de Ambiente

Certifique-se de que as seguintes variáveis estejam configuradas:

```bash
# Banco de dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=process_service
DB_SSL_MODE=disable

# Ambiente (para dados de exemplo)
ENVIRONMENT=development
```

## Controle de Migrações

O sistema cria automaticamente uma tabela `schema_migrations` para controlar quais migrações já foram executadas:

```sql
-- Verificar migrações executadas
SELECT * FROM schema_migrations ORDER BY version;

-- Estrutura da tabela
\d schema_migrations
```

## Rollback de Migrações

**Atenção**: Não há rollback automático. Para reverter:

1. **Backup sempre antes de executar migrações**
2. Para desenvolvimento, recrie o banco:
```bash
dropdb process_service
createdb process_service
```

3. Para produção, crie scripts de rollback específicos

## Dados de Exemplo

A migração 006 inclui dados de exemplo que são criados apenas em ambiente de desenvolvimento/teste.

### IDs de Exemplo:
- **Tenant ID**: `12345678-1234-1234-1234-123456789012`
- **Client ID**: `87654321-4321-4321-4321-210987654321`
- **Processo**: `1234567-89.2024.1.01.0001`

### Views Criadas:
- `v_processes_with_stats`: Processos com estatísticas
- `v_pending_notifications`: Notificações pendentes
- `v_dashboard_stats`: Métricas para dashboard

## Troubleshooting

### Erro: "relation already exists"
```sql
-- Verificar se tabela já existe
\dt

-- Dropar tabela se necessário (CUIDADO!)
DROP TABLE IF EXISTS table_name CASCADE;
```

### Erro: "extension does not exist"
```sql
-- Instalar extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

### Performance após migrações
```sql
-- Atualizar estatísticas
ANALYZE;

-- Verificar índices
SELECT schemaname, tablename, indexname, indexdef 
FROM pg_indexes 
WHERE schemaname = 'public'
ORDER BY tablename, indexname;
```

## Manutenção

### Executar manutenção periódica:
```sql
-- Executar função de manutenção
SELECT maintain_database();

-- Limpar dados antigos (cuidado!)
SELECT cleanup_old_data(2555); -- 7 anos
```

### Monitorar tamanho das tabelas:
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Próximos Passos

Após executar as migrações:

1. **Testar conexão** da aplicação
2. **Verificar logs** de migração
3. **Executar testes** básicos
4. **Configurar backup** automático
5. **Implementar monitoring** das tabelas

## Contato

Para dúvidas sobre as migrações, consulte a documentação do projeto ou entre em contato com a equipe de desenvolvimento.