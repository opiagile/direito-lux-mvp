# Diretrizes de Desenvolvimento - Direito Lux

Este documento contém todas as convenções, padrões e verificações obrigatórias que devem ser seguidas no desenvolvimento do projeto Direito Lux.

## 📋 Checklist Obrigatório para Novos Serviços

### 1. Estrutura de Imports Go
**SEMPRE verificar e corrigir imports ausentes:**

```go
// Imports padrão sempre necessários
import (
    "context"
    "fmt"
    "time"
    
    // Bibliotecas externas
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go.uber.org/zap"
    
    // Imports internos do projeto
    "github.com/direito-lux/[service-name]/internal/..."
)
```

### 2. Correções Comuns de Imports

**Config packages sempre precisam:**
```go
import (
    "fmt"      // Para fmt.Errorf, fmt.Sprintf
    "time"     // Para time.Duration
)
```

**Logging packages sempre precisam:**
```go
import (
    "context"  // Para context.Context
    "fmt"      // Para formatação
    "time"     // Para timestamps
)
```

**Metrics packages sempre precisam:**
```go
import (
    "runtime"  // Para runtime.NumGoroutine, runtime.MemStats
    "time"     // Para time.Duration
)
```

**Middleware packages sempre precisam:**
```go
import (
    "os"       // Para os.Stdout
    // Remover "time" se não usado
)
```

### 3. Correções de Código Padrão

**Gin middleware retorno:**
```go
// ❌ ERRADO
return gin.Next

// ✅ CORRETO  
return func(c *gin.Context) { c.Next() }
```

**Logger middleware:**
```go
// ❌ ERRADO
return gin.LoggerWithWriter(logger.Sugar().Desugar().Core())

// ✅ CORRETO
return gin.LoggerWithWriter(os.Stdout)
```

**Event imports:**
```go
// SEMPRE adicionar se usar opentracing
"github.com/opentracing/opentracing-go"

// E usar opentracing.StartSpan ao invés de tracing.StartSpan
span := opentracing.StartSpan("message_handler", opentracing.ChildOf(spanCtx))
```

### 4. Verificações de Compilação

**SEMPRE rodar antes de commit:**
```bash
# Para cada serviço
cd services/[service-name]
go mod tidy
go build ./cmd/server

# Se falhar, verificar:
# 1. Imports ausentes (fmt, time, runtime, os)
# 2. Referencias incorretas (gin.Next, logger.Core())
# 3. Função signatures (LogError com err direto)
```

### 5. Configurações de Ambiente

**Todas as variáveis necessárias para desenvolvimento local:**
```bash
# Database
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=direito_lux
export DB_PASSWORD=dev_password_123
export DB_SSL_MODE=disable

# RabbitMQ
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USER=guest
export RABBITMQ_PASSWORD=guest

# Redis
export REDIS_HOST=localhost
export REDIS_PORT=6379

# JWT
export JWT_SECRET=development_jwt_secret_key_change_in_production
export JWT_EXPIRY=24h

# Keycloak (se necessário)
export KEYCLOAK_URL=http://localhost:8080
export KEYCLOAK_CLIENT_SECRET=dev_client_secret

# Service configs
export ENVIRONMENT=development
export LOG_LEVEL=debug
export METRICS_ENABLED=true
export TRACING_ENABLED=false
```

### 6. Template para go.mod

```go
module github.com/direito-lux/[SERVICE-NAME]

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/google/uuid v1.4.0
    github.com/opentracing/opentracing-go v1.2.0
    go.uber.org/fx v1.20.1
    go.uber.org/zap v1.26.0
    // Adicionar outras dependências conforme necessário
)
```

### 7. Estrutura de Main Template

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/direito-lux/[service-name]/internal/infrastructure/config"
    "github.com/direito-lux/[service-name]/internal/infrastructure/logging"
    
    "go.uber.org/fx"
    "go.uber.org/zap"
)

func main() {
    // Carregar configurações
    cfg, err := config.Load()
    if err != nil {
        fmt.Printf("Erro ao carregar configurações: %v\n", err)
        os.Exit(1)
    }

    // Configurar logger
    logger, err := logging.NewLogger(cfg.LogLevel, cfg.Environment)
    if err != nil {
        fmt.Printf("Erro ao configurar logger: %v\n", err)
        os.Exit(1)
    }

    // Resto da implementação...
}
```

## 🧪 Scripts de Teste e Validação

### Scripts Obrigatórios
1. `start-services.sh` - Inicia todos os serviços
2. `stop-services.sh` - Para todos os serviços  
3. `test-local.sh` - Testa tudo localmente
4. `build-all.sh` - Compila todos os serviços

### Comando de Validação Completa
```bash
# SEMPRE executar antes de commit/push
./test-local.sh
```

## 🔧 Comandos de Desenvolvimento

### Build individual
```bash
cd services/[service-name]
go mod tidy
go build ./cmd/server
```

### Migrações
```bash
docker run --rm -v "${PWD}/migrations:/migrations" --network host \
  migrate/migrate -path=/migrations/ \
  -database "postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" up
```

### Logs
```bash
tail -f logs/[service-name].log
```

## ❌ Erros Comuns e Soluções

### "undefined: fmt"
**Solução:** Adicionar `"fmt"` aos imports

### "undefined: time"  
**Solução:** Adicionar `"time"` aos imports

### "undefined: runtime"
**Solução:** Adicionar `"runtime"` aos imports em metrics

### "cannot use ... as io.Writer"
**Solução:** Usar `os.Stdout` ao invés de `logger.Core()`

### "undefined: gin.Next"
**Solução:** Usar `func(c *gin.Context) { c.Next() }`

### "LogError signature mismatch"
**Solução:** Passar `err` diretamente, não `zap.Error(err)`

## 🎯 Workflow de Desenvolvimento

1. **Criar novo serviço:**
   - Usar template de estrutura
   - Verificar go.mod com nome correto
   - Adicionar imports obrigatórios

2. **Durante desenvolvimento:**
   - `go build` frequentemente
   - Corrigir imports imediatamente
   - Testar localmente

3. **Antes de commit:**
   - `./test-local.sh`
   - Verificar todos os serviços funcionando
   - Validar health checks

4. **Deploy:**
   - Só após 100% local funcionando
   - Usar ambiente DEV primeiro

## 📝 Notas Importantes

- **NUNCA** commitar código que não compila
- **SEMPRE** usar o template de configuração
- **SEMPRE** testar localmente antes de deploy
- **SEMPRE** verificar imports ausentes
- Manter este documento atualizado com novos padrões descobertos

---

Este documento deve ser consultado SEMPRE antes de criar novos serviços ou fazer alterações significativas.