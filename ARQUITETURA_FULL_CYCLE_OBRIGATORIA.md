# 🔄 ARQUITETURA FULL CYCLE OBRIGATÓRIA - DIREITO LUX

## 📋 **FULL CYCLE DEVELOPMENT - CONCEITOS OBRIGATÓRIOS**

### **⚠️ IMPORTANTE: ARQUITETURA BASEADA EM FULL CYCLE**

**✅ TODOS os conceitos de Full Cycle Development DEVEM ser aplicados**
- Ownership completo do desenvolvedor
- Observabilidade nativa em cada serviço
- Deployment contínuo com feedback loops
- Responsabilidade end-to-end: código → deploy → monitoramento → suporte

---

## 🎯 **PRINCÍPIOS FULL CYCLE OBRIGATÓRIOS**

### **1. 👨‍💻 OWNERSHIP COMPLETO**

```yaml
DESENVOLVEDOR_FULL_CYCLE_DIREITO_LUX:
├── 💻 Código: Desenvolve funcionalidade completa
├── 🧪 Testes: Escreve testes automatizados
├── 🚀 Deploy: Responsável pelo deployment
├── 📊 Monitoramento: Acompanha métricas em produção
├── 🚨 Alertas: Recebe e resolve alertas
├── 🐛 Bugs: Corrige problemas em produção
├── 📈 Performance: Otimiza baseado em métricas
└── 🔄 Melhoria: Itera baseado em feedback
```

### **2. 📊 OBSERVABILIDADE NATIVA**

```yaml
CADA_MICROSERVIÇO_OBRIGATÓRIO:
├── 📝 Logs estruturados (JSON + Zap)
├── 📈 Métricas Prometheus expostas
├── 🏥 Health checks (/health, /ready)
├── 🔍 Distributed tracing (Jaeger)
├── 🚨 Alertas configurados (AlertManager)
├── 📊 Dashboards Grafana
└── 🔔 Notificações Slack/PagerDuty
```

### **3. 🚀 DEPLOYMENT CONTÍNUO**

```yaml
PIPELINE_FULL_CYCLE_OBRIGATÓRIO:
├── 📝 Commit: Código + testes
├── 🔧 CI: Testes automatizados
├── 🏗️ Build: Docker images
├── 🚀 Deploy: Automático GKE
├── 📊 Monitor: Métricas em tempo real
├── 🚨 Alert: Notificação problemas
└── 🔄 Feedback: Melhoria contínua
```

### **4. 🔄 FEEDBACK LOOPS RÁPIDOS**

```yaml
CICLOS_FEEDBACK_OBRIGATÓRIOS:
├── 🧪 Desenvolvimento: Testes locais (<1min)
├── 🔧 Integração: CI/CD pipeline (<5min)
├── 🚀 Deployment: Deploy automático (<10min)
├── 📊 Produção: Métricas em tempo real
├── 🚨 Alertas: Notificação imediata
└── 🔄 Correção: Hotfix rápido (<30min)
```

---

## 🛠️ **IMPLEMENTAÇÃO FULL CYCLE POR SERVIÇO**

### **📋 CHECKLIST OBRIGATÓRIO PARA CADA MICROSERVIÇO**

```yaml
✅ CÓDIGO:
├── ✅ Arquitetura Hexagonal
├── ✅ Testes unitários (>80% coverage)
├── ✅ Testes de integração
├── ✅ Validação de entrada
└── ✅ Tratamento de erros

✅ OBSERVABILIDADE:
├── ✅ Logs estruturados (Zap)
├── ✅ Métricas Prometheus
├── ✅ Tracing distribuído
├── ✅ Health checks
└── ✅ Graceful shutdown

✅ DEPLOY:
├── ✅ Dockerfile otimizado
├── ✅ Kubernetes manifests
├── ✅ ConfigMaps/Secrets
├── ✅ Liveness/Readiness probes
└── ✅ Resource limits

✅ MONITORAMENTO:
├── ✅ Dashboard Grafana
├── ✅ Alertas configurados
├── ✅ SLIs/SLOs definidos
├── ✅ Runbooks documentados
└── ✅ Escalação configurada
```

### **🔧 TEMPLATE FULL CYCLE OBRIGATÓRIO**

```go
// template-service/internal/infrastructure/observability/
package observability

import (
    "context"
    "net/http"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.uber.org/zap"
    "github.com/opentracing/opentracing-go"
)

// ObservabilityManager - OBRIGATÓRIO EM TODOS OS SERVIÇOS
type ObservabilityManager struct {
    logger  *zap.Logger
    metrics *prometheus.Registry
    tracer  opentracing.Tracer
    
    // Métricas padrão obrigatórias
    requestsTotal    *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    activeConnections prometheus.Gauge
    errorRate        *prometheus.CounterVec
    dependencyUp     *prometheus.GaugeVec
}

// NewObservabilityManager - OBRIGATÓRIO
func NewObservabilityManager(serviceName string) *ObservabilityManager {
    logger, _ := zap.NewProduction()
    
    metrics := prometheus.NewRegistry()
    
    // Métricas padrão obrigatórias
    requestsTotal := prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: serviceName + "_requests_total",
            Help: "Total number of requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration := prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    serviceName + "_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    activeConnections := prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: serviceName + "_active_connections",
            Help: "Number of active connections",
        },
    )
    
    errorRate := prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: serviceName + "_errors_total",
            Help: "Total number of errors",
        },
        []string{"type", "operation"},
    )
    
    dependencyUp := prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: serviceName + "_dependency_up",
            Help: "Dependency availability",
        },
        []string{"dependency"},
    )
    
    // Registrar métricas
    metrics.MustRegister(requestsTotal)
    metrics.MustRegister(requestDuration)
    metrics.MustRegister(activeConnections)
    metrics.MustRegister(errorRate)
    metrics.MustRegister(dependencyUp)
    
    return &ObservabilityManager{
        logger:            logger,
        metrics:           metrics,
        requestsTotal:     requestsTotal,
        requestDuration:   requestDuration,
        activeConnections: activeConnections,
        errorRate:         errorRate,
        dependencyUp:      dependencyUp,
    }
}

// HTTPMiddleware - OBRIGATÓRIO EM TODOS OS HANDLERS
func (o *ObservabilityManager) HTTPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Tracing obrigatório
        span := o.tracer.StartSpan(r.URL.Path)
        defer span.Finish()
        
        // Context com span
        ctx := opentracing.ContextWithSpan(r.Context(), span)
        r = r.WithContext(ctx)
        
        // Logging obrigatório
        o.logger.Info("request started",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("trace_id", span.Context().TraceID()),
        )
        
        // Métricas obrigatórias
        o.activeConnections.Inc()
        defer o.activeConnections.Dec()
        
        // Executar handler
        next.ServeHTTP(w, r)
        
        // Métricas pós-request
        duration := time.Since(start)
        o.requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
        o.requestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
        
        // Logging obrigatório
        o.logger.Info("request completed",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Duration("duration", duration),
        )
    })
}

// HealthCheck - OBRIGATÓRIO EM TODOS OS SERVIÇOS
func (o *ObservabilityManager) HealthCheck(dependencies map[string]func() error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        health := map[string]string{
            "status": "healthy",
            "service": "auth-service",
            "timestamp": time.Now().Format(time.RFC3339),
        }
        
        // Verificar dependências
        for name, check := range dependencies {
            if err := check(); err != nil {
                o.logger.Error("dependency check failed",
                    zap.String("dependency", name),
                    zap.Error(err),
                )
                
                // Atualizar métrica
                o.dependencyUp.WithLabelValues(name).Set(0)
                
                health["status"] = "unhealthy"
                health[name] = "down"
                
                w.WriteHeader(http.StatusServiceUnavailable)
                json.NewEncoder(w).Encode(health)
                return
            }
            
            // Atualizar métrica
            o.dependencyUp.WithLabelValues(name).Set(1)
            health[name] = "up"
        }
        
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(health)
    }
}

// MetricsHandler - OBRIGATÓRIO EM TODOS OS SERVIÇOS
func (o *ObservabilityManager) MetricsHandler() http.Handler {
    return promhttp.HandlerFor(o.metrics, promhttp.HandlerOpts{})
}
```

---

## 🚨 **ALERTAS FULL CYCLE OBRIGATÓRIOS**

### **📊 MÉTRICAS OBRIGATÓRIAS POR SERVIÇO**

```yaml
# prometheus/alerts/direito-lux.yml
groups:
- name: direito-lux-full-cycle
  rules:
  
  # Alertas de latência (SLI)
  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(service_request_duration_seconds_bucket[5m])) > 1
    for: 2m
    labels:
      severity: warning
      team: full-cycle
    annotations:
      summary: "High latency detected"
      description: "Service {{ $labels.service }} has P95 latency > 1s"
      runbook: "https://wiki.direitolux.com/runbooks/latency"
  
  # Alertas de erro (SLI)
  - alert: HighErrorRate
    expr: rate(service_errors_total[5m]) / rate(service_requests_total[5m]) > 0.05
    for: 1m
    labels:
      severity: critical
      team: full-cycle
    annotations:
      summary: "High error rate detected"
      description: "Service {{ $labels.service }} has error rate > 5%"
      runbook: "https://wiki.direitolux.com/runbooks/errors"
  
  # Alertas de dependência
  - alert: DependencyDown
    expr: service_dependency_up == 0
    for: 30s
    labels:
      severity: critical
      team: full-cycle
    annotations:
      summary: "Dependency down"
      description: "Service {{ $labels.service }} dependency {{ $labels.dependency }} is down"
      runbook: "https://wiki.direitolux.com/runbooks/dependencies"
  
  # Alertas de resource
  - alert: HighCPUUsage
    expr: rate(container_cpu_usage_seconds_total[5m]) > 0.8
    for: 5m
    labels:
      severity: warning
      team: full-cycle
    annotations:
      summary: "High CPU usage"
      description: "Service {{ $labels.service }} CPU usage > 80%"
      runbook: "https://wiki.direitolux.com/runbooks/cpu"
```

### **🔔 ESCALAÇÃO FULL CYCLE**

```yaml
# alertmanager/config.yml
route:
  group_by: ['service', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'full-cycle-team'
  routes:
  - match:
      severity: critical
    receiver: 'pagerduty-critical'
  - match:
      severity: warning
    receiver: 'slack-warnings'

receivers:
- name: 'full-cycle-team'
  slack_configs:
  - channel: '#direito-lux-alerts'
    title: 'Full Cycle Alert'
    text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'

- name: 'pagerduty-critical'
  pagerduty_configs:
  - service_key: 'YOUR_PAGERDUTY_KEY'
    description: 'Critical alert in Direito Lux'

- name: 'slack-warnings'
  slack_configs:
  - channel: '#direito-lux-warnings'
    title: 'Warning Alert'
    text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
```

---

## 📊 **DASHBOARDS FULL CYCLE OBRIGATÓRIOS**

### **🎯 DASHBOARD POR SERVIÇO**

```json
{
  "dashboard": {
    "title": "Direito Lux - Auth Service Full Cycle",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(auth_service_requests_total[5m])",
            "legendFormat": "{{ method }} {{ endpoint }}"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(auth_service_errors_total[5m]) / rate(auth_service_requests_total[5m])",
            "legendFormat": "Error Rate"
          }
        ]
      },
      {
        "title": "Latency (P95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(auth_service_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P95 Latency"
          }
        ]
      },
      {
        "title": "Active Users",
        "type": "singlestat",
        "targets": [
          {
            "expr": "auth_service_active_sessions",
            "legendFormat": "Active Sessions"
          }
        ]
      }
    ]
  }
}
```

---

## 🎯 **CHECKLIST FULL CYCLE FINAL**

### **✅ OBRIGATÓRIO PARA CADA MICROSERVIÇO**

```yaml
ANTES_DO_DEPLOY:
├── ✅ Logs estruturados implementados
├── ✅ Métricas Prometheus expostas
├── ✅ Health checks funcionando
├── ✅ Tracing distribuído ativo
├── ✅ Alertas configurados
├── ✅ Dashboard Grafana criado
├── ✅ Runbook documentado
├── ✅ Escalação configurada
├── ✅ Testes >80% coverage
└── ✅ SLIs/SLOs definidos

RESPONSABILIDADES_DESENVOLVEDOR:
├── ✅ Código + testes
├── ✅ Deploy via GitHub Actions
├── ✅ Monitoramento ativo
├── ✅ Resposta a alertas
├── ✅ Correção de bugs
├── ✅ Otimização baseada em métricas
└── ✅ Melhoria contínua
```

---

## 🏆 **RESULTADO FINAL**

### **✅ DIREITO LUX COM FULL CYCLE DEVELOPMENT**

**TODOS os microserviços DEVEM seguir os princípios Full Cycle:**
- **Ownership completo** do desenvolvedor
- **Observabilidade nativa** em cada serviço
- **Deployment contínuo** automatizado
- **Feedback loops** rápidos e eficientes
- **Responsabilidade end-to-end** do código à produção

**🔄 ARQUITETURA FULL CYCLE 100% OBRIGATÓRIA NO DIREITO LUX!**