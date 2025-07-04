#!/bin/bash

echo "🚀 Iniciando Report Service..."

# Configurar variáveis de ambiente
export PORT=8087
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=direito_lux_dev
export DB_USER=direito_lux
export DB_PASSWORD=dev_password_123
export DB_SSLMODE=disable
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
export REDIS_PASSWORD=dev_redis_123
export PROCESS_SERVICE_URL=http://localhost:8083
export DATAJUD_SERVICE_URL=http://localhost:8084
export ENV=development

# Parar container se estiver rodando
docker stop direito-lux-report 2>/dev/null || true

# Entrar no diretório do serviço
cd services/report-service

# Verificar se binário existe
if [ ! -f "report-service" ]; then
    echo "🔨 Compilando Report Service..."
    go build -mod=mod -o report-service cmd/server/main.go
fi

# Executar em background
echo "🎯 Executando Report Service em background..."
nohup ./report-service > report.log 2>&1 &

# Salvar PID
echo $! > report.pid

echo "✅ Report Service iniciado com PID: $(cat report.pid)"
echo "📋 Logs em: services/report-service/report.log"
echo ""
echo "Para parar o serviço:"
echo "  kill $(cat report.pid)"