# 🎯 PLANO COMPLETO: DataJud HTTP Client Real

## 📊 Resumo Executivo

**Objetivo**: Substituir implementação mock do DataJud Service por cliente HTTP real integrado com a API oficial do CNJ.

**Impacto**: Unlock completo do ambiente STAGING e validação com dados reais.

**Timeline**: 2-3 dias (16-24 horas de desenvolvimento)

**Prioridade**: 🔥 **MÁXIMA** - Blocker para ambiente produção

---

## 🎯 CONTEXTO E JUSTIFICATIVA

### Por que DataJud Real é Crítico?
1. **Mock não funciona em produção** - APIs CNJ exigem autenticação real
2. **Validação de arquitetura** - Testar circuit breaker, rate limiting, cache real
3. **Dados reais** - Processos jurídicos reais para testes E2E
4. **Compliance legal** - Integração oficial obrigatória para SaaS jurídico
5. **Diferencial competitivo** - Poucos concorrentes têm integração DataJud

### Estado Atual vs Objetivo
```
ANTES (Mock):                    DEPOIS (Real):
┌─────────────────┐             ┌─────────────────┐
│   DataJud       │             │   DataJud       │
│   Service       │             │   Service       │
│                 │             │                 │
│ ┌─────────────┐ │             │ ┌─────────────┐ │
│ │ Mock Client │ │    ═══►     │ │ HTTP Client │ │
│ │ (Fake Data) │ │             │ │ (CNJ API)   │ │
│ └─────────────┘ │             │ └─────────────┘ │
└─────────────────┘             └─────────────────┘
                                         │
                                ┌─────────────────┐
                                │   CNJ DataJud   │
                                │   API (Real)    │
                                └─────────────────┘
```

---

## 🔍 FASE 1: ANÁLISE DA API CNJ DATAJUD

### 1.1 Research Técnico (2-3 horas)

#### Endpoints Principais CNJ
```bash
# API Base URL
https://api-publica.datajud.cnj.jus.br

# Principais endpoints
GET  /api/v1/processos/{numero}           # Consulta processo
GET  /api/v1/processos/{numero}/movs      # Movimentações
GET  /api/v1/processos/bulk               # Consulta em lote
GET  /api/v1/tribunais                    # Lista tribunais
POST /api/v1/auth                         # Autenticação
```

#### Documentação Oficial
- **Portal CNJ**: https://datajud.cnj.jus.br/
- **API Docs**: https://api-publica.datajud.cnj.jus.br/docs
- **Certificação**: https://datajud.cnj.jus.br/certificacao
- **Rate Limits**: 10.000 consultas/dia (produção), 100/dia (desenvolvimento)

#### Autenticação Obrigatória
```bash
# Certificado Digital A1 (arquivo .p12)
DATAJUD_CERTIFICATE_PATH=/certs/certificado.p12
DATAJUD_CERTIFICATE_PASSWORD=senha_certificado

# Ou Certificado A3 (hardware token)  
DATAJUD_A3_PROVIDER=eToken
DATAJUD_A3_PIN=1234
```

### 1.2 Análise de Requisitos

#### Rate Limiting Real
- **Produção**: 10.000 consultas/dia
- **Desenvolvimento**: 100 consultas/dia
- **Burst**: Máximo 10 consultas/minuto
- **Circuit Breaker**: Necessário para falhas

#### Estrutura de Resposta Real
```json
{
  "success": true,
  "data": {
    "numeroProcesso": "1234567-89.2023.8.26.0001",
    "classe": {
      "codigo": "319",
      "nome": "Procedimento Comum"
    },
    "assunto": [
      {
        "codigo": "4391", 
        "nome": "Direito do Consumidor"
      }
    ],
    "orgaoJulgador": {
      "codigo": "26",
      "nome": "1ª Vara Cível"
    },
    "dataAjuizamento": "2023-01-15T00:00:00Z",
    "valorCausa": 15000.50,
    "partes": [...],
    "movimentacoes": [...]
  },
  "errors": [],
  "metadata": {
    "tribunal": "TJSP",
    "instancia": "1G",
    "timestamp": "2025-07-09T15:30:00Z"
  }
}
```

---

## 🏗️ FASE 2: ARQUITETURA DA SOLUÇÃO

### 2.1 Estrutura do Cliente HTTP Real

```
datajud-service/
├── internal/
│   ├── infrastructure/
│   │   └── http/
│   │       ├── datajud_real_client.go      # ✅ Já existe (base)
│   │       ├── datajud_real_client_test.go # ✅ Já existe
│   │       ├── certificate_manager.go      # 🆕 CRIAR
│   │       ├── rate_limiter_real.go        # 🆕 CRIAR  
│   │       ├── circuit_breaker_real.go     # 🆕 CRIAR
│   │       └── response_parser.go          # 🆕 CRIAR
│   ├── domain/
│   │   ├── cnj_types.go                    # 🆕 CRIAR (types CNJ)
│   │   └── datajud_request.go              # ✅ Existe, atualizar
│   └── application/
│       └── datajud_service.go              # ✅ Atualizar para usar real client
```

### 2.2 Interface Unificada (Mock + Real)

```go
// DataJudClient interface unificada
type DataJudClient interface {
    QueryProcess(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error)
    QueryMovements(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error)
    QueryParties(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error)
    BulkQuery(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error)
    TestConnection(ctx context.Context) error
    Close() error
}

// Factory pattern para escolher implementação
func NewDataJudClient(config *config.Config) DataJudClient {
    if config.IsDataJudMockEnabled() {
        return NewMockClient()           // ✅ Existe
    }
    return NewRealHTTPClient(config)     // 🆕 IMPLEMENTAR
}
```

### 2.3 Configuração Flexível

```go
// config/config.go - Adicionar
type DataJudConfig struct {
    MockEnabled         bool   `env:"DATAJUD_MOCK_ENABLED" envDefault:"true"`
    BaseURL            string `env:"DATAJUD_BASE_URL" envDefault:"https://api-publica.datajud.cnj.jus.br"`
    CertificatePath    string `env:"DATAJUD_CERTIFICATE_PATH"`
    CertificatePassword string `env:"DATAJUD_CERTIFICATE_PASSWORD"`
    RateLimit          int    `env:"DATAJUD_RATE_LIMIT" envDefault:"100"` // dev: 100, prod: 10000
    Timeout            int    `env:"DATAJUD_TIMEOUT_SECONDS" envDefault:"30"`
    RetryAttempts      int    `env:"DATAJUD_RETRY_ATTEMPTS" envDefault:"3"`
    CircuitBreakerEnabled bool `env:"DATAJUD_CIRCUIT_BREAKER" envDefault:"true"`
}
```

---

## 🛠️ FASE 3: IMPLEMENTAÇÃO DETALHADA

### 3.1 Certificate Manager (4 horas)

```go
// certificate_manager.go
package http

import (
    "crypto/tls"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "golang.org/x/crypto/pkcs12"
)

type CertificateManager struct {
    certPath     string
    certPassword string
    tlsConfig    *tls.Config
}

func NewCertificateManager(certPath, password string) (*CertificateManager, error) {
    cm := &CertificateManager{
        certPath:     certPath,
        certPassword: password,
    }
    
    if err := cm.loadCertificate(); err != nil {
        return nil, fmt.Errorf("failed to load certificate: %w", err)
    }
    
    return cm, nil
}

func (cm *CertificateManager) loadCertificate() error {
    // Ler arquivo .p12
    certData, err := ioutil.ReadFile(cm.certPath)
    if err != nil {
        return fmt.Errorf("failed to read certificate file: %w", err)
    }
    
    // Decodificar PKCS#12
    privateKey, cert, err := pkcs12.Decode(certData, cm.certPassword)
    if err != nil {
        return fmt.Errorf("failed to decode PKCS#12: %w", err)
    }
    
    // Criar TLS certificate
    tlsCert := tls.Certificate{
        Certificate: [][]byte{cert.Raw},
        PrivateKey:  privateKey,
    }
    
    // Configurar TLS
    cm.tlsConfig = &tls.Config{
        Certificates: []tls.Certificate{tlsCert},
        ClientAuth:   tls.RequireAndVerifyClientCert,
    }
    
    return nil
}

func (cm *CertificateManager) GetTLSConfig() *tls.Config {
    return cm.tlsConfig
}
```

### 3.2 Real HTTP Client (6 horas)

```go
// datajud_real_client.go (expandir existente)
package http

import (
    "bytes"
    "context"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "github.com/direito-lux/datajud-service/internal/domain"
    "github.com/direito-lux/datajud-service/internal/infrastructure/config"
)

type RealHTTPClient struct {
    httpClient    *http.Client
    config        *config.DataJudConfig
    certManager   *CertificateManager
    rateLimiter   *RealRateLimiter
    circuitBreaker *RealCircuitBreaker
    baseURL       string
}

func NewRealHTTPClient(cfg *config.Config) (*RealHTTPClient, error) {
    // Certificate manager
    certManager, err := NewCertificateManager(
        cfg.DataJud.CertificatePath, 
        cfg.DataJud.CertificatePassword,
    )
    if err != nil {
        return nil, fmt.Errorf("certificate manager failed: %w", err)
    }
    
    // HTTP client com certificado
    transport := &http.Transport{
        TLSClientConfig: certManager.GetTLSConfig(),
        // Timeouts otimizados
        IdleConnTimeout:       30 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ResponseHeaderTimeout: 30 * time.Second,
    }
    
    httpClient := &http.Client{
        Transport: transport,
        Timeout:   time.Duration(cfg.DataJud.Timeout) * time.Second,
    }
    
    // Rate limiter real
    rateLimiter := NewRealRateLimiter(cfg.DataJud.RateLimit)
    
    // Circuit breaker
    circuitBreaker := NewRealCircuitBreaker(cfg.DataJud.CircuitBreakerEnabled)
    
    return &RealHTTPClient{
        httpClient:     httpClient,
        config:         &cfg.DataJud,
        certManager:    certManager,
        rateLimiter:    rateLimiter,
        circuitBreaker: circuitBreaker,
        baseURL:        cfg.DataJud.BaseURL,
    }, nil
}

func (c *RealHTTPClient) QueryProcess(ctx context.Context, req *domain.DataJudRequest, provider *domain.CNPJProvider) (*domain.DataJudResponse, error) {
    // Rate limiting check
    if !c.rateLimiter.Allow() {
        return nil, domain.ErrRateLimitExceeded
    }
    
    // Circuit breaker check
    if !c.circuitBreaker.AllowRequest() {
        return nil, domain.ErrCircuitBreakerOpen
    }
    
    // Construir URL
    url := fmt.Sprintf("%s/api/v1/processos/%s", c.baseURL, req.ProcessNumber)
    
    // Fazer requisição
    response, err := c.doRequest(ctx, "GET", url, nil)
    if err != nil {
        c.circuitBreaker.RecordFailure()
        return nil, err
    }
    
    c.circuitBreaker.RecordSuccess()
    return response, nil
}

func (c *RealHTTPClient) doRequest(ctx context.Context, method, url string, body interface{}) (*domain.DataJudResponse, error) {
    var reqBody io.Reader
    
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        reqBody = bytes.NewReader(jsonBody)
    }
    
    // Criar request
    req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Headers obrigatórios CNJ
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("User-Agent", "DireitoLux/1.0")
    
    // Executar request
    startTime := time.Now()
    resp, err := c.httpClient.Do(req)
    duration := time.Since(startTime)
    
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Ler response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }
    
    // Parse response para domain types
    parsedData, err := c.parseResponse(responseBody, req.ProcessNumber)
    if err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    // Construir domain response
    domainResponse := &domain.DataJudResponse{
        ID:         uuid.New(),
        StatusCode: resp.StatusCode,
        Body:       responseBody,
        Headers:    convertHeaders(resp.Header),
        ProcessData: parsedData,
        FromCache:  false,
        ReceivedAt: time.Now(),
        Size:       int64(len(responseBody)),
        Duration:   int(duration.Milliseconds()),
    }
    
    return domainResponse, nil
}
```

### 3.3 Rate Limiter Real (2 horas)

```go
// rate_limiter_real.go
package http

import (
    "sync"
    "time"
    "golang.org/x/time/rate"
)

type RealRateLimiter struct {
    limiter   *rate.Limiter
    daily     *DailyCounter
    burst     *BurstLimiter
    mu        sync.RWMutex
}

type DailyCounter struct {
    count     int
    limit     int
    resetTime time.Time
    mu        sync.RWMutex
}

type BurstLimiter struct {
    requests []time.Time
    limit    int // 10 requests per minute
    mu       sync.RWMutex
}

func NewRealRateLimiter(dailyLimit int) *RealRateLimiter {
    return &RealRateLimiter{
        limiter: rate.NewLimiter(rate.Every(time.Minute), 10), // 10/min burst
        daily: &DailyCounter{
            limit:     dailyLimit,
            resetTime: getNextMidnight(),
        },
        burst: &BurstLimiter{
            limit: 10,
        },
    }
}

func (rl *RealRateLimiter) Allow() bool {
    // Check daily limit
    if !rl.daily.Allow() {
        return false
    }
    
    // Check burst limit
    if !rl.burst.Allow() {
        return false
    }
    
    // Check rate limiter
    return rl.limiter.Allow()
}

func (dc *DailyCounter) Allow() bool {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    // Reset if new day
    if time.Now().After(dc.resetTime) {
        dc.count = 0
        dc.resetTime = getNextMidnight()
    }
    
    if dc.count >= dc.limit {
        return false
    }
    
    dc.count++
    return true
}

func (bl *BurstLimiter) Allow() bool {
    bl.mu.Lock()
    defer bl.mu.Unlock()
    
    now := time.Now()
    cutoff := now.Add(-time.Minute)
    
    // Remove old requests
    var validRequests []time.Time
    for _, reqTime := range bl.requests {
        if reqTime.After(cutoff) {
            validRequests = append(validRequests, reqTime)
        }
    }
    bl.requests = validRequests
    
    // Check limit
    if len(bl.requests) >= bl.limit {
        return false
    }
    
    // Add current request
    bl.requests = append(bl.requests, now)
    return true
}

func getNextMidnight() time.Time {
    now := time.Now()
    return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}
```

### 3.4 Circuit Breaker Real (2 horas)

```go
// circuit_breaker_real.go
package http

import (
    "sync"
    "time"
)

type CircuitState int

const (
    CircuitClosed CircuitState = iota
    CircuitOpen
    CircuitHalfOpen
)

type RealCircuitBreaker struct {
    state           CircuitState
    failureCount    int
    successCount    int
    failureThreshold int
    recoveryTimeout time.Duration
    lastFailure     time.Time
    enabled         bool
    mu              sync.RWMutex
}

func NewRealCircuitBreaker(enabled bool) *RealCircuitBreaker {
    return &RealCircuitBreaker{
        state:           CircuitClosed,
        failureThreshold: 5, // 5 failures = open circuit
        recoveryTimeout: 30 * time.Second,
        enabled:         enabled,
    }
}

func (cb *RealCircuitBreaker) AllowRequest() bool {
    if !cb.enabled {
        return true
    }
    
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    switch cb.state {
    case CircuitClosed:
        return true
    case CircuitOpen:
        // Check if recovery time has passed
        if time.Since(cb.lastFailure) > cb.recoveryTimeout {
            cb.mu.RUnlock()
            cb.mu.Lock()
            cb.state = CircuitHalfOpen
            cb.successCount = 0
            cb.mu.Unlock()
            cb.mu.RLock()
            return true
        }
        return false
    case CircuitHalfOpen:
        return true
    default:
        return false
    }
}

func (cb *RealCircuitBreaker) RecordSuccess() {
    if !cb.enabled {
        return
    }
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    switch cb.state {
    case CircuitHalfOpen:
        cb.successCount++
        if cb.successCount >= 3 { // 3 successes = close circuit
            cb.state = CircuitClosed
            cb.failureCount = 0
        }
    case CircuitClosed:
        cb.failureCount = 0 // Reset failure count on success
    }
}

func (cb *RealCircuitBreaker) RecordFailure() {
    if !cb.enabled {
        return
    }
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    cb.failureCount++
    cb.lastFailure = time.Now()
    
    if cb.failureCount >= cb.failureThreshold {
        cb.state = CircuitOpen
    }
}

func (cb *RealCircuitBreaker) GetState() CircuitState {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    return cb.state
}
```

---

## 🧪 FASE 4: TESTES E VALIDAÇÃO

### 4.1 Testes Unitários (3 horas)

```go
// datajud_real_client_test.go (expandir existente)
package http

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRealHTTPClient_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup real client com certificado de desenvolvimento
    config := &config.Config{
        DataJud: config.DataJudConfig{
            MockEnabled:         false,
            BaseURL:            "https://api-publica.datajud.cnj.jus.br",
            CertificatePath:    "testdata/cert_dev.p12",
            CertificatePassword: "dev123",
            RateLimit:          10, // Limite baixo para testes
        },
    }
    
    client, err := NewRealHTTPClient(config)
    require.NoError(t, err)
    defer client.Close()
    
    t.Run("TestConnection", func(t *testing.T) {
        err := client.TestConnection(context.Background())
        assert.NoError(t, err)
    })
    
    t.Run("QueryProcess_Real", func(t *testing.T) {
        req := &domain.DataJudRequest{
            ProcessNumber: "1234567-89.2023.8.26.0001", // Processo público conhecido
            RequestType:   domain.RequestTypeProcess,
            TenantID:      uuid.New(),
            ClientID:      uuid.New(),
        }
        
        response, err := client.QueryProcess(context.Background(), req, nil)
        require.NoError(t, err)
        
        assert.Equal(t, 200, response.StatusCode)
        assert.NotNil(t, response.ProcessData)
        assert.Equal(t, req.ProcessNumber, response.ProcessData.Number)
        assert.False(t, response.FromCache)
        assert.Greater(t, response.Size, int64(0))
    })
    
    t.Run("RateLimit_Respected", func(t *testing.T) {
        // Test rate limiting
        for i := 0; i < 15; i++ { // Tentar mais que o limite
            req := &domain.DataJudRequest{
                ProcessNumber: fmt.Sprintf("test-%d", i),
                RequestType:   domain.RequestTypeProcess,
            }
            
            _, err := client.QueryProcess(context.Background(), req, nil)
            if i >= 10 { // Após 10 requests
                assert.Error(t, err)
                assert.Contains(t, err.Error(), "rate limit")
            }
        }
    })
    
    t.Run("CircuitBreaker_Opens", func(t *testing.T) {
        // Simular falhas consecutivas
        for i := 0; i < 6; i++ { // Mais que threshold
            req := &domain.DataJudRequest{
                ProcessNumber: "invalid-process-number",
                RequestType:   domain.RequestTypeProcess,
            }
            
            _, err := client.QueryProcess(context.Background(), req, nil)
            // Últimas requests devem falhar por circuit breaker
            if i >= 5 {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), "circuit breaker")
            }
        }
    })
}

func TestCertificateManager(t *testing.T) {
    t.Run("LoadCertificate_Valid", func(t *testing.T) {
        cm, err := NewCertificateManager("testdata/cert_valid.p12", "password")
        require.NoError(t, err)
        
        tlsConfig := cm.GetTLSConfig()
        assert.NotNil(t, tlsConfig)
        assert.Len(t, tlsConfig.Certificates, 1)
    })
    
    t.Run("LoadCertificate_Invalid", func(t *testing.T) {
        _, err := NewCertificateManager("testdata/cert_invalid.p12", "wrong_password")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "failed to decode PKCS#12")
    })
}
```

### 4.2 Testes de Integração (2 horas)

```bash
# integration_test.sh
#!/bin/bash

echo "🧪 Executando testes de integração DataJud Real..."

# Setup certificado de desenvolvimento
export DATAJUD_MOCK_ENABLED=false
export DATAJUD_CERTIFICATE_PATH="./testdata/cert_dev.p12"
export DATAJUD_CERTIFICATE_PASSWORD="dev123"
export DATAJUD_RATE_LIMIT=10

# Testes com curl (validação manual)
echo "📋 Testando conexão direta..."
curl -X GET "https://api-publica.datajud.cnj.jus.br/api/v1/tribunais" \
  --cert ./testdata/cert_dev.p12:dev123 \
  --cert-type P12 \
  -H "Accept: application/json"

# Testes com serviço
echo "📋 Testando via DataJud Service..."
curl -X POST http://localhost:8084/api/v1/process/query \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 550e8400-e29b-41d4-a716-446655440001" \
  -d '{
    "process_number": "1234567-89.2023.8.26.0001",
    "use_cache": false,
    "urgent": false
  }'

# Validar rate limiting
echo "📋 Testando rate limiting..."
for i in {1..15}; do
  echo "Request $i..."
  curl -X GET http://localhost:8084/api/v1/stats \
    -H "X-Tenant-ID: 550e8400-e29b-41d4-a716-446655440001" &
done
wait

echo "✅ Testes de integração concluídos"
```

---

## 📅 TIMELINE DETALHADO

### DIA 1 (8 horas)
- **08:00-10:00**: FASE 1 - Research API CNJ + documentação
- **10:00-12:00**: FASE 2 - Arquitetura e design patterns  
- **13:00-17:00**: FASE 3.1 - Certificate Manager implementation
- **17:00-18:00**: FASE 3.2 - Real HTTP Client (base structure)

### DIA 2 (8 horas)
- **08:00-12:00**: FASE 3.2 - Real HTTP Client (complete implementation)
- **13:00-15:00**: FASE 3.3 - Rate Limiter Real  
- **15:00-17:00**: FASE 3.4 - Circuit Breaker Real
- **17:00-18:00**: Integration setup

### DIA 3 (8 horas)
- **08:00-11:00**: FASE 4.1 - Testes unitários
- **11:00-13:00**: FASE 4.2 - Testes de integração
- **14:00-16:00**: FASE 5 - Validação com dados reais
- **16:00-17:00**: FASE 6 - Setup STAGING
- **17:00-18:00**: Documentação e handover

---

## 🔧 CONFIGURAÇÃO DE AMBIENTE

### Desenvolvimento
```bash
# .env.development
DATAJUD_MOCK_ENABLED=false
DATAJUD_BASE_URL=https://api-publica.datajud.cnj.jus.br
DATAJUD_CERTIFICATE_PATH=./certs/desenvolvimento.p12
DATAJUD_CERTIFICATE_PASSWORD=dev_cert_password
DATAJUD_RATE_LIMIT=100
DATAJUD_TIMEOUT_SECONDS=30
DATAJUD_RETRY_ATTEMPTS=3
DATAJUD_CIRCUIT_BREAKER=true
```

### Staging
```bash
# .env.staging
DATAJUD_MOCK_ENABLED=false
DATAJUD_BASE_URL=https://api-publica.datajud.cnj.jus.br
DATAJUD_CERTIFICATE_PATH=./certs/staging.p12
DATAJUD_CERTIFICATE_PASSWORD=staging_cert_password
DATAJUD_RATE_LIMIT=1000
DATAJUD_TIMEOUT_SECONDS=30
DATAJUD_RETRY_ATTEMPTS=3
DATAJUD_CIRCUIT_BREAKER=true
```

### Produção
```bash
# .env.production
DATAJUD_MOCK_ENABLED=false
DATAJUD_BASE_URL=https://api-publica.datajud.cnj.jus.br
DATAJUD_CERTIFICATE_PATH=./certs/producao.p12
DATAJUD_CERTIFICATE_PASSWORD=${DATAJUD_CERT_PASSWORD} # Secret
DATAJUD_RATE_LIMIT=10000
DATAJUD_TIMEOUT_SECONDS=30
DATAJUD_RETRY_ATTEMPTS=3
DATAJUD_CIRCUIT_BREAKER=true
```

---

## ⚠️ RISCOS E MITIGAÇÕES

### 🚨 Riscos Técnicos

#### 1. Certificado Digital A1/A3
**Risco**: Certificado não reconhecido pelo CNJ
**Mitigação**: 
- Validar certificado com CNJ antes da implementação
- Backup com certificado A3 (hardware token)
- Ambiente de testes com certificado válido

#### 2. Rate Limiting Restritivo
**Risco**: 100 consultas/dia muito limitado para testes
**Mitigação**:
- Implementar cache agressivo
- Usar dados mock para desenvolvimento
- Solicitar quota adicional ao CNJ

#### 3. API CNJ Instável
**Risco**: API pública pode ter instabilidade
**Mitigação**:
- Circuit breaker robusto
- Retry exponential backoff
- Fallback para cache/mock

#### 4. Performance
**Risco**: API externa pode ser lenta
**Mitigação**:
- Timeout configurável
- Connection pooling
- Cache Redis para responses

### 💰 Riscos de Negócio

#### 1. Custo Certificado Digital
**Risco**: Certificado A1/A3 tem custo mensal
**Mitigação**:
- Orçar certificado no MVP
- Compartilhar certificado entre ambientes
- ROI positivo com clientes reais

#### 2. Compliance Legal
**Risco**: Uso incorreto da API CNJ
**Mitigação**:
- Seguir documentação oficial
- Respeitar rate limits
- Logs de auditoria completos

---

## 📊 MÉTRICAS DE SUCESSO

### KPIs Técnicos
- **Disponibilidade**: >99% requests CNJ com sucesso
- **Performance**: <2s response time médio
- **Rate Limiting**: 0 violações de quota CNJ
- **Circuit Breaker**: Recuperação automática em <1min

### KPIs de Produto
- **Dados Reais**: 100% processos vêm da API CNJ
- **Cache Hit**: >80% consultas servidas do cache
- **Uptime**: >99.9% disponibilidade do serviço
- **User Experience**: <3s para mostrar dados reais

---

## 🎯 ENTREGÁVEIS

### Código
1. **RealHTTPClient** - Cliente HTTP completo CNJ
2. **CertificateManager** - Gestão certificados A1/A3  
3. **RealRateLimiter** - Rate limiting 100% compatível CNJ
4. **RealCircuitBreaker** - Circuit breaker production-ready
5. **Integration Tests** - Testes com API real

### Documentação
1. **API Integration Guide** - Como integrar com CNJ
2. **Certificate Setup Guide** - Setup certificados
3. **Troubleshooting Guide** - Debug problemas comuns
4. **Rate Limiting Guide** - Gestão de quotas
5. **Monitoring Guide** - Métricas e alertas

### Infraestrutura
1. **STAGING Environment** - Ambiente com API real
2. **Certificate Storage** - Gestão segura certificados
3. **Monitoring Dashboards** - Grafana dashboards CNJ
4. **Alerting Rules** - Alertas rate limit, circuit breaker
5. **Deployment Scripts** - Deploy automatizado

---

## 🚀 PRÓXIMOS PASSOS PÓS-IMPLEMENTAÇÃO

### Otimizações (Semana seguinte)
1. **Cache Strategy** - Cache inteligente multi-layer
2. **Batch Processing** - Consultas em lote otimizadas
3. **Data Sync** - Sync incremental com CNJ
4. **Performance Tuning** - Otimização connection pooling

### Funcionalidades Avançadas
1. **Real-time Webhooks** - Notificações CNJ em tempo real
2. **Advanced Search** - Busca combinada Elasticsearch + CNJ
3. **ML Integration** - Predição com dados históricos CNJ
4. **Tribunal Specific** - Otimizações por tribunal

---

## 🎉 CONCLUSÃO

**🎯 DataJud HTTP Client Real é o unlock definitivo para STAGING**

Esta implementação remove a última barreira técnica para ambiente de produção, permitindo:

- ✅ **Validação com dados reais** CNJ
- ✅ **Compliance legal** total  
- ✅ **Architecture validation** completa
- ✅ **Customer demos** com dados reais
- ✅ **Production readiness** 100%

**Timeline**: 2-3 dias intensivos de desenvolvimento

**ROI**: Unlock STAGING + Produção + Revenue

**Risco**: Baixo (arquitetura já validada, APIs documentadas)

---

*Plano criado em 09/07/2025 - Ready for execution*

📧 **Próximo passo**: Approval para início da implementação