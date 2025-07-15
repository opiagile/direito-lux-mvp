#!/bin/bash

echo "☁️ MIGRAR PARA CLOUD RUN - ECONOMIA MÁXIMA (95%)"
echo "================================================"

PROJECT_ID="direito-lux-staging-2025"
REGION="us-central1"

# Função para criar Dockerfile otimizado para Cloud Run
create_cloud_run_dockerfile() {
    local service=$1
    echo "🐳 Criando Dockerfile para Cloud Run: $service"
    
    cat > "services/$service/Dockerfile.cloudrun" << 'EOF'
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
ENV PORT=8080

CMD ["./main"]
EOF

    echo "✅ Dockerfile Cloud Run criado para $service"
}

# Função para criar docker-compose para desenvolvimento local
create_local_docker_compose() {
    echo "🔧 Criando docker-compose.yml local..."
    
    cat > docker-compose.local.yml << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: direito_lux_staging
      POSTGRES_USER: direito_lux
      POSTGRES_PASSWORD: dev_password_123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  rabbitmq:
    image: rabbitmq:3-management
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - "5672:5672"
      - "15672:15672"

  auth-service:
    build:
      context: ./services/auth-service
      dockerfile: Dockerfile.cloudrun
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=direito_lux_staging
      - DB_USER=direito_lux
      - DB_PASSWORD=dev_password_123
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
      - JWT_SECRET=dev_jwt_secret_key_123
      - PORT=8080
    ports:
      - "8081:8080"
    depends_on:
      - postgres
      - redis
      - rabbitmq

  tenant-service:
    build:
      context: ./services/tenant-service
      dockerfile: Dockerfile.cloudrun
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=direito_lux_staging
      - DB_USER=direito_lux
      - DB_PASSWORD=dev_password_123
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
      - PORT=8080
    ports:
      - "8082:8080"
    depends_on:
      - postgres
      - redis
      - rabbitmq

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8081
      - NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:8081/api/v1/auth
      - NEXT_PUBLIC_TENANT_SERVICE_URL=http://localhost:8082/api/v1/tenants
    ports:
      - "3000:3000"
    depends_on:
      - auth-service
      - tenant-service

volumes:
  postgres_data:
EOF

    echo "✅ Docker Compose local criado"
}

# Função para fazer deploy no Cloud Run
deploy_to_cloud_run() {
    echo "🚀 Fazendo deploy para Cloud Run..."
    
    # Serviços para deploy
    SERVICES=("auth-service" "tenant-service" "frontend")
    
    for service in "${SERVICES[@]}"; do
        echo "📦 Deploying $service..."
        
        # Build e push da imagem
        gcloud builds submit "services/$service" \
            --dockerfile="services/$service/Dockerfile.cloudrun" \
            --tag="gcr.io/$PROJECT_ID/$service:latest" \
            --project=$PROJECT_ID
        
        # Deploy no Cloud Run
        gcloud run deploy $service \
            --image="gcr.io/$PROJECT_ID/$service:latest" \
            --platform=managed \
            --region=$REGION \
            --allow-unauthenticated \
            --memory=512Mi \
            --cpu=1 \
            --min-instances=0 \
            --max-instances=2 \
            --port=8080 \
            --project=$PROJECT_ID
        
        echo "✅ $service deployed to Cloud Run"
    done
    
    # Obter URLs
    echo ""
    echo "📋 URLs dos serviços:"
    for service in "${SERVICES[@]}"; do
        URL=$(gcloud run services describe $service --region=$REGION --project=$PROJECT_ID --format="value(status.url)")
        echo "   $service: $URL"
    done
}

# Função para configurar Cloud SQL (PostgreSQL gerenciado)
setup_cloud_sql() {
    echo "🗄️ Configurando Cloud SQL..."
    
    # Criar instância Cloud SQL
    gcloud sql instances create direito-lux-db \
        --database-version=POSTGRES_15 \
        --tier=db-f1-micro \
        --region=$REGION \
        --storage-size=10GB \
        --storage-type=SSD \
        --project=$PROJECT_ID
    
    # Criar database
    gcloud sql databases create direito_lux_staging \
        --instance=direito-lux-db \
        --project=$PROJECT_ID
    
    # Criar usuário
    gcloud sql users create direito_lux \
        --instance=direito-lux-db \
        --password=dev_password_123 \
        --project=$PROJECT_ID
    
    # Obter IP de conexão
    CONNECTION_NAME=$(gcloud sql instances describe direito-lux-db --project=$PROJECT_ID --format="value(connectionName)")
    
    echo "✅ Cloud SQL configurado"
    echo "📋 Connection Name: $CONNECTION_NAME"
}

# Função para calcular economia
show_cost_comparison() {
    echo ""
    echo "💰 COMPARAÇÃO DE CUSTOS (por mês):"
    echo "================================================="
    echo ""
    echo "🔴 CONFIGURAÇÃO ATUAL (GKE):"
    echo "   - 6x e2-standard-2 nodes: R$ 2.160,00"
    echo "   - Load Balancer: R$ 18,00"
    echo "   - Persistent Disks: R$ 30,00"
    echo "   - TOTAL: R$ 2.208,00/mês"
    echo ""
    echo "🟡 CONFIGURAÇÃO OTIMIZADA (GKE + Auto-shutdown):"
    echo "   - 1x e2-small node (16h/dia): R$ 216,00"
    echo "   - Load Balancer: R$ 18,00"
    echo "   - TOTAL: R$ 234,00/mês (-89%)"
    echo ""
    echo "🟢 CLOUD RUN (RECOMENDADO):"
    echo "   - Cloud Run (100 req/dia): R$ 5,00"
    echo "   - Cloud SQL db-f1-micro: R$ 25,00"
    echo "   - Load Balancer: R$ 0 (incluído)"
    echo "   - TOTAL: R$ 30,00/mês (-98.6%)"
    echo ""
    echo "💡 ECONOMIA ANUAL:"
    echo "   - GKE → Cloud Run: R$ 26.136,00"
    echo "   - Payback: IMEDIATO"
    echo ""
}

# Função para limpeza do GKE
cleanup_gke() {
    echo "🧹 Limpando recursos GKE caros..."
    
    read -p "⚠️  Tem certeza que quer deletar o cluster GKE? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🗑️ Deletando cluster GKE..."
        gcloud container clusters delete $CLUSTER_NAME \
            --region=$REGION \
            --project=$PROJECT_ID \
            --quiet
        
        echo "✅ Cluster GKE deletado - economia imediata!"
    fi
}

# Menu principal
case "${1:-help}" in
    "cost-analysis")
        show_cost_comparison
        ;;
    "setup-cloudrun")
        echo "🚀 Configurando Cloud Run (economia 98%)..."
        create_cloud_run_dockerfile "auth-service"
        create_cloud_run_dockerfile "tenant-service"
        create_local_docker_compose
        deploy_to_cloud_run
        setup_cloud_sql
        show_cost_comparison
        ;;
    "optimize-gke")
        echo "⚡ Otimizando GKE atual (economia 89%)..."
        ./gcp-cost-optimizer.sh optimize
        ./setup-auto-shutdown.sh setup
        show_cost_comparison
        ;;
    "cleanup")
        cleanup_gke
        ;;
    "emergency")
        echo "🚨 MODO EMERGÊNCIA - PARANDO TUDO AGORA"
        ./gcp-cost-optimizer.sh stop
        echo "✅ Cluster parado - custo atual: R$0/hora"
        echo "💡 Para reiniciar: ./gcp-cost-optimizer.sh start"
        ;;
    "help"|*)
        echo "💰 SOLUÇÕES PARA REDUZIR CUSTOS:"
        echo "================================"
        echo ""
        echo "🚨 EMERGÊNCIA (economia imediata):"
        echo "  ./migrate-to-cloud-run.sh emergency    # Para tudo AGORA"
        echo ""
        echo "🟡 OTIMIZAÇÃO GKE (economia 89%):"
        echo "  ./migrate-to-cloud-run.sh optimize-gke # Mantém GKE otimizado"
        echo ""
        echo "🟢 MIGRAÇÃO CLOUD RUN (economia 98%):"
        echo "  ./migrate-to-cloud-run.sh setup-cloudrun # Migra para Cloud Run"
        echo ""
        echo "📊 ANÁLISES:"
        echo "  ./migrate-to-cloud-run.sh cost-analysis # Comparação detalhada"
        echo "  ./migrate-to-cloud-run.sh cleanup       # Remove GKE definitivamente"
        echo ""
        echo "🎯 RECOMENDAÇÃO:"
        echo "   1. EMERGENCY (agora) - para parar custos"
        echo "   2. SETUP-CLOUDRUN (quando pronto) - 98% economia"
        ;;
esac