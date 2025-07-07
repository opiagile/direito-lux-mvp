#!/bin/bash

echo "🚀 Iniciando DataJud Service..."

# Configurar variáveis de ambiente
export PORT=8084
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=direito_lux_dev
export DB_USER=direito_lux
export DB_PASSWORD=dev_password_123
export DB_SSLMODE=disable
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=dev_redis_123
export DATAJUD_API_URL=https://api-publica.datajud.cnj.jus.br
export DATAJUD_API_KEY=demo_key
export RATE_LIMIT_DAILY=10000
export ENV=development

# Parar container se estiver rodando
docker stop direito-lux-datajud 2>/dev/null || true

# Entrar no diretório do serviço
cd services/datajud-service

# Verificar se binário existe
if [ ! -f "datajud-service" ]; then
    echo "🔨 Compilando DataJud Service..."
    go build -mod=mod -o datajud-service cmd/server/main.go
fi

# Executar em background
echo "🎯 Executando DataJud Service em background..."
nohup ./datajud-service > datajud.log 2>&1 &

# Salvar PID
echo $! > datajud.pid

echo "✅ DataJud Service iniciado com PID: $(cat datajud.pid)"
echo "📋 Logs em: services/datajud-service/datajud.log"
echo ""
echo "Para parar o serviço:"
echo "  kill $(cat datajud.pid)"