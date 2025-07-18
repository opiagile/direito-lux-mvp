# 🚀 DEPLOY GCP VIA GITHUB ACTIONS - DIREITO LUX

## 📋 **ESTRATÉGIA DE DEPLOY DEFINIDA**

### **🔄 BASEADA EM FULL CYCLE DEVELOPMENT**

**✅ ARQUITETURA FULL CYCLE OBRIGATÓRIA**
- Pipeline deve seguir todos os conceitos Full Cycle
- Ownership completo: desenvolvedor responsável por código → deploy → monitoramento
- Observabilidade nativa em cada microserviço
- Feedback loops rápidos e automatizados
- Deployment contínuo com rollback automático

### **✅ CONFIRMADO: Deploy GCP via GitHub Actions**

```yaml
PIPELINE DEFINIDO:
├── DEV (Local): Docker Compose - Desenvolvimento
├── CI/CD: GitHub Actions - Automação
├── PROD: Google Cloud Platform - Produção
└── DEPLOY: Automático via push main branch
```

---

## 🔧 **CONFIGURAÇÃO GITHUB ACTIONS**

### **📁 Estrutura de Workflows**
```
.github/workflows/
├── ci-cd.yml                 # Pipeline principal
├── deploy-production.yml     # Deploy produção
├── security-scan.yml         # Scanning segurança
└── tests.yml                 # Testes automatizados
```

### **🔄 PIPELINE FULL CYCLE OBRIGATÓRIO**

#### **📋 Conceitos Full Cycle no Pipeline**
```yaml
FULL_CYCLE_PIPELINE:
├── Code: Desenvolvedor escreve código
├── Test: Testes automatizados obrigatórios
├── Build: Build automático com métricas
├── Deploy: Deploy automático com health checks
├── Monitor: Observabilidade nativa
├── Alert: Alertas para o desenvolvedor
└── Improve: Feedback loop de melhoria
```

#### **🔧 Instrumentação Full Cycle**
```yaml
CADA_DEPLOYMENT_DEVE_TER:
├── Logs estruturados (JSON)
├── Métricas Prometheus expostas
├── Health checks (/health, /ready)
├── Distributed tracing ativo
├── Alertas configurados
└── Dashboards por serviço
```

### **🚀 Workflow Principal (.github/workflows/ci-cd.yml)**
```yaml
name: CI/CD Pipeline - Direito Lux

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  PROJECT_ID: direito-lux-prod
  GKE_CLUSTER: direito-lux-cluster
  GKE_ZONE: us-central1-a
  REGISTRY: gcr.io

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Run Tests
      run: |
        make test-all
        make test-coverage
    
    - name: Security Scan
      run: |
        make security-scan

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Google Cloud CLI
      uses: google-github-actions/setup-gcloud@v1
      with:
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        project_id: ${{ env.PROJECT_ID }}
    
    - name: Configure Docker
      run: gcloud auth configure-docker
    
    - name: Build Images
      run: |
        docker build -t $REGISTRY/$PROJECT_ID/auth-service:$GITHUB_SHA ./services/auth-service
        docker build -t $REGISTRY/$PROJECT_ID/process-service:$GITHUB_SHA ./services/process-service
        docker build -t $REGISTRY/$PROJECT_ID/datajud-service:$GITHUB_SHA ./services/datajud-service
        docker build -t $REGISTRY/$PROJECT_ID/notification-service:$GITHUB_SHA ./services/notification-service
        docker build -t $REGISTRY/$PROJECT_ID/ai-service:$GITHUB_SHA ./services/ai-service
        docker build -t $REGISTRY/$PROJECT_ID/search-service:$GITHUB_SHA ./services/search-service
        docker build -t $REGISTRY/$PROJECT_ID/mcp-service:$GITHUB_SHA ./services/mcp-service
        docker build -t $REGISTRY/$PROJECT_ID/report-service:$GITHUB_SHA ./services/report-service
        docker build -t $REGISTRY/$PROJECT_ID/billing-service:$GITHUB_SHA ./services/billing-service
        docker build -t $REGISTRY/$PROJECT_ID/frontend:$GITHUB_SHA ./frontend
    
    - name: Push Images
      run: |
        docker push $REGISTRY/$PROJECT_ID/auth-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/process-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/datajud-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/notification-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/ai-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/search-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/mcp-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/report-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/billing-service:$GITHUB_SHA
        docker push $REGISTRY/$PROJECT_ID/frontend:$GITHUB_SHA

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Google Cloud CLI
      uses: google-github-actions/setup-gcloud@v1
      with:
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        project_id: ${{ env.PROJECT_ID }}
    
    - name: Get GKE credentials
      run: |
        gcloud container clusters get-credentials $GKE_CLUSTER --zone $GKE_ZONE
    
    - name: Update K8s manifests
      run: |
        sed -i "s|IMAGE_TAG|$GITHUB_SHA|g" k8s/services/*.yaml
        sed -i "s|PROJECT_ID|$PROJECT_ID|g" k8s/services/*.yaml
    
    - name: Deploy to GKE
      run: |
        kubectl apply -f k8s/namespace.yaml
        kubectl apply -f k8s/databases/
        kubectl apply -f k8s/services/
        kubectl apply -f k8s/ingress/
        kubectl apply -f k8s/monitoring/
    
    - name: Verify deployment
      run: |
        kubectl rollout status deployment/auth-service -n direito-lux
        kubectl rollout status deployment/process-service -n direito-lux
        kubectl rollout status deployment/frontend -n direito-lux
    
    - name: Run smoke tests
      run: |
        make test-smoke-production
```

---

## 🔐 **SECRETS GITHUB NECESSÁRIOS**

### **🔑 Configuração de Secrets**
```yaml
GITHUB_SECRETS:
  # GCP
  GCP_SA_KEY: "service-account-key.json"
  GCP_PROJECT_ID: "direito-lux-prod"
  
  # Database
  DB_PASSWORD: "production-db-password"
  DB_CONNECTION_STRING: "postgresql://..."
  
  # APIs Externas
  OPENAI_API_KEY: "sk-real-production-key"
  DATAJUD_API_KEY: "real-cnj-production-key"
  WHATSAPP_ACCESS_TOKEN: "production-whatsapp-token"
  TELEGRAM_BOT_TOKEN: "production-telegram-token"
  
  # JWT e Segurança
  JWT_SECRET: "production-jwt-secret-key"
  ENCRYPTION_KEY: "production-encryption-key"
  
  # Monitoring
  PROMETHEUS_PASSWORD: "monitoring-password"
  GRAFANA_ADMIN_PASSWORD: "grafana-password"
```

### **📝 Como Configurar Secrets**
```bash
# No GitHub, ir em Settings > Secrets and Variables > Actions
# Adicionar cada secret individualmente

# Exemplo para GCP Service Account:
1. Criar service account no GCP
2. Baixar JSON key
3. Codificar em base64: cat key.json | base64
4. Adicionar no GitHub Secrets como GCP_SA_KEY
```

---

## 🗂️ **ESTRUTURA K8S PARA DEPLOY**

### **📁 Manifests Kubernetes**
```
k8s/
├── namespace.yaml
├── databases/
│   ├── postgres.yaml
│   ├── redis.yaml
│   └── elasticsearch.yaml
├── services/
│   ├── auth-service.yaml
│   ├── process-service.yaml
│   ├── datajud-service.yaml
│   ├── notification-service.yaml
│   ├── ai-service.yaml
│   ├── search-service.yaml
│   ├── mcp-service.yaml
│   ├── report-service.yaml
│   ├── billing-service.yaml
│   └── frontend.yaml
├── ingress/
│   └── ingress.yaml
└── monitoring/
    ├── prometheus.yaml
    └── grafana.yaml
```

### **🚀 Deployment Template Example**
```yaml
# k8s/services/auth-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: direito-lux
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: gcr.io/PROJECT_ID/auth-service:IMAGE_TAG
        ports:
        - containerPort: 8081
        env:
        - name: DB_HOST
          value: postgres-service
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: direito-lux
spec:
  selector:
    app: auth-service
  ports:
  - port: 8081
    targetPort: 8081
  type: ClusterIP
```

---

## 🌐 **CONFIGURAÇÃO INGRESS E DOMÍNIO**

### **🔗 Ingress Controller**
```yaml
# k8s/ingress/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: direito-lux-ingress
  namespace: direito-lux
  annotations:
    kubernetes.io/ingress.class: "gce"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    kubernetes.io/ingress.global-static-ip-name: "direito-lux-ip"
spec:
  tls:
  - hosts:
    - app.direitolux.com.br
    - api.direitolux.com.br
    secretName: direito-lux-tls
  rules:
  - host: app.direitolux.com.br
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend
            port:
              number: 3000
  - host: api.direitolux.com.br
    http:
      paths:
      - path: /api/v1/auth
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 8081
      - path: /api/v1/processes
        pathType: Prefix
        backend:
          service:
            name: process-service
            port:
              number: 8083
      - path: /api/v1/datajud
        pathType: Prefix
        backend:
          service:
            name: datajud-service
            port:
              number: 8084
      - path: /api/v1/notifications
        pathType: Prefix
        backend:
          service:
            name: notification-service
            port:
              number: 8085
```

---

## 📊 **MONITORAMENTO E OBSERVABILIDADE**

### **📈 Prometheus + Grafana**
```yaml
# Deploy automático via GitHub Actions
MONITORING_STACK:
├── Prometheus: Coleta métricas
├── Grafana: Dashboards
├── Alertmanager: Alertas
└── Jaeger: Tracing distribuído
```

### **🔔 Alertas Configurados**
```yaml
ALERTAS_PRODUÇÃO:
├── CPU > 80%
├── Memória > 85%
├── Disco > 90%
├── Pods crashando
├── Latência > 1s
└── Erros > 5%
```

---

## 🔄 **FLUXO DE DEPLOY COMPLETO**

### **1. 💻 Desenvolvimento Local**
```bash
# Trabalho em branch feature
git checkout -b feature/nova-funcionalidade
# Desenvolvimento...
git add .
git commit -m "feat: nova funcionalidade"
git push origin feature/nova-funcionalidade
```

### **2. 🔍 Pull Request**
```bash
# Criar PR para main
# GitHub Actions executa:
├── Testes automatizados
├── Security scan
├── Code quality check
└── Build test
```

### **3. ✅ Merge para Main**
```bash
# Após aprovação e merge
git checkout main
git pull origin main
# GitHub Actions executa AUTOMATICAMENTE:
├── ✅ Testes completos
├── 🔨 Build todas as imagens
├── 📤 Push para GCR
├── 🚀 Deploy para GKE
└── 🧪 Smoke tests
```

### **4. 🏃 Deploy Automático**
```bash
# Deploy acontece AUTOMATICAMENTE em ~10 minutos
PIPELINE_EXECUTION:
├── 00:00 - Trigger (push main)
├── 02:00 - Testes completos
├── 04:00 - Build imagens (9 serviços)
├── 07:00 - Push para registry
├── 08:00 - Deploy GKE
├── 10:00 - Verificação saúde
└── 10:30 - ✅ PRODUÇÃO ATUALIZADA
```

---

## 🎯 **CONFIGURAÇÃO INICIAL NECESSÁRIA**

### **1. 🔧 Setup GCP**
```bash
# Criar projeto
gcloud projects create direito-lux-prod

# Habilitar APIs
gcloud services enable container.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable cloudsql.googleapis.com

# Criar cluster GKE
gcloud container clusters create direito-lux-cluster \
    --zone us-central1-a \
    --num-nodes 3 \
    --enable-autoscaling \
    --min-nodes 1 \
    --max-nodes 10
```

### **2. 🔐 Service Account**
```bash
# Criar service account
gcloud iam service-accounts create direito-lux-deployer

# Dar permissões
gcloud projects add-iam-policy-binding direito-lux-prod \
    --member="serviceAccount:direito-lux-deployer@direito-lux-prod.iam.gserviceaccount.com" \
    --role="roles/container.developer"

# Gerar chave
gcloud iam service-accounts keys create key.json \
    --iam-account=direito-lux-deployer@direito-lux-prod.iam.gserviceaccount.com
```

### **3. 🌐 Domínio e DNS**
```bash
# Reservar IP estático
gcloud compute addresses create direito-lux-ip --global

# Configurar DNS
# app.direitolux.com.br -> IP estático
# api.direitolux.com.br -> IP estático
```

---

## 🚀 **COMANDOS ESSENCIAIS**

### **📋 Makefile para Automação**
```makefile
# Makefile
.PHONY: test deploy status logs

test:
	go test ./...
	npm test

deploy:
	git push origin main
	# GitHub Actions faz o resto

status:
	kubectl get pods -n direito-lux
	kubectl get services -n direito-lux

logs:
	kubectl logs -f deployment/auth-service -n direito-lux

rollback:
	kubectl rollout undo deployment/auth-service -n direito-lux

scale:
	kubectl scale deployment/auth-service --replicas=5 -n direito-lux
```

### **🔍 Monitoramento**
```bash
# Ver status deployment
kubectl get deployments -n direito-lux

# Ver logs em tempo real
kubectl logs -f -l app=auth-service -n direito-lux

# Verificar métricas
kubectl top pods -n direito-lux
```

---

## ✅ **CONFIRMAÇÃO: DEPLOY GCP VIA GITHUB ACTIONS**

### **🎯 RESUMO FINAL**
```yaml
ESTRATÉGIA_DEPLOY:
  ✅ Ambiente DEV: Docker Compose (local)
  ✅ CI/CD: GitHub Actions (automático)
  ✅ Produção: Google Cloud Platform (GKE)
  ✅ Trigger: Push para branch main
  ✅ Tempo: ~10 minutos automático
  ✅ Rollback: Comando único
  ✅ Monitoring: Prometheus + Grafana
  ✅ Alertas: Configurados
  ✅ Segurança: Secrets + Service Account
```

### **🔧 PRÓXIMOS PASSOS**
1. **Configurar secrets** no GitHub
2. **Criar service account** no GCP
3. **Configurar domínio** DNS
4. **Testar pipeline** com deploy inicial
5. **Configurar monitoramento** e alertas

**🚀 DEPLOY GCP VIA GITHUB ACTIONS 100% DEFINIDO E DOCUMENTADO!**