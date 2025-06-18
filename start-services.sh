#!/bin/bash

# Script para inicializar todos os microserviços localmente
# Configurações de ambiente padrão para desenvolvimento

set -e

echo "🚀 Iniciando microserviços Direito Lux..."

# Configurações comuns
export ENVIRONMENT=development
export LOG_LEVEL=debug
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=direito_lux
export DB_PASSWORD=dev_password_123
export DB_SSL_MODE=disable
export REDIS_HOST=localhost
export REDIS_PORT=6379
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USER=guest
export RABBITMQ_PASSWORD=guest
export RABBITMQ_VHOST=/
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Specific service configs
export AUTH_SERVICE_PORT=8081
export TENANT_SERVICE_PORT=8082
export PROCESS_SERVICE_PORT=8083
export DATAJUD_SERVICE_PORT=8084
export NOTIFICATION_SERVICE_PORT=8085

# Database names for each service
export AUTH_DB_NAME=direito_lux_dev
export TENANT_DB_NAME=direito_lux_dev
export PROCESS_DB_NAME=direito_lux_dev
export DATAJUD_DB_NAME=direito_lux_dev
export NOTIFICATION_DB_NAME=direito_lux_dev

# JWT Secrets (for dev only)
export JWT_SECRET=development_jwt_secret_key_change_in_production
export JWT_EXPIRY=24h
export JWT_REFRESH_EXPIRY=168h

# Keycloak configs (dev only)
export KEYCLOAK_URL=http://localhost:8080
export KEYCLOAK_REALM=direito-lux
export KEYCLOAK_CLIENT_ID=direito-lux-api
export KEYCLOAK_CLIENT_SECRET=dev_client_secret
export KEYCLOAK_ADMIN_USER=admin
export KEYCLOAK_ADMIN_PASSWORD=admin

# Outras configurações necessárias
export SERVER_PORT=8081
export CORS_ENABLED=true
export TRACING_ENABLED=false
export METRICS_ENABLED=true

# Start services in background
echo "📦 Building services..."

# Auth Service
echo "🔐 Starting Auth Service..."
cd services/auth-service
go build -o auth-server ./cmd/server
nohup ./auth-server > ../../logs/auth-service.log 2>&1 &
AUTH_PID=$!
echo "Auth Service PID: $AUTH_PID"
cd ../..

# Tenant Service 
echo "🏢 Starting Tenant Service..."
cd services/tenant-service
go build -o tenant-server ./cmd/server
nohup ./tenant-server > ../../logs/tenant-service.log 2>&1 &
TENANT_PID=$!
echo "Tenant Service PID: $TENANT_PID"
cd ../..

# Process Service
echo "📋 Starting Process Service..."
cd services/process-service
go build -o process-server ./cmd/server
nohup ./process-server > ../../logs/process-service.log 2>&1 &
PROCESS_PID=$!
echo "Process Service PID: $PROCESS_PID"
cd ../..

# DataJud Service
echo "🔗 Starting DataJud Service..."
cd services/datajud-service
go build -o datajud-server ./cmd/server
nohup ./datajud-server > ../../logs/datajud-service.log 2>&1 &
DATAJUD_PID=$!
echo "DataJud Service PID: $DATAJUD_PID"
cd ../..

# Notification Service
echo "📧 Starting Notification Service..."
cd services/notification-service
go build -o notification-server ./cmd/server
nohup ./notification-server > ../../logs/notification-service.log 2>&1 &
NOTIFICATION_PID=$!
echo "Notification Service PID: $NOTIFICATION_PID"
cd ../..

# Save PIDs for later cleanup
echo "$AUTH_PID" > .auth.pid
echo "$TENANT_PID" > .tenant.pid
echo "$PROCESS_PID" > .process.pid
echo "$DATAJUD_PID" > .datajud.pid
echo "$NOTIFICATION_PID" > .notification.pid

echo "✅ All services started!"
echo ""
echo "📋 Service URLs:"
echo "🔐 Auth Service:          http://localhost:8081"
echo "🏢 Tenant Service:        http://localhost:8082"
echo "📋 Process Service:       http://localhost:8083"
echo "🔗 DataJud Service:       http://localhost:8084"
echo "📧 Notification Service:  http://localhost:8085"
echo ""
echo "📝 Logs:"
echo "tail -f logs/auth-service.log"
echo "tail -f logs/tenant-service.log"
echo "tail -f logs/process-service.log"
echo "tail -f logs/datajud-service.log"
echo "tail -f logs/notification-service.log"
echo ""
echo "🛑 To stop services: ./stop-services.sh"