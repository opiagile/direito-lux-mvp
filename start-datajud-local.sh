#!/bin/bash

echo "🚀 Iniciando DataJud Service localmente..."

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

# Entrar no diretório do serviço
cd services/datajud-service

# Remover vendor se existir
if [ -d "vendor" ]; then
    echo "🗑️ Removendo vendor..."
    rm -rf vendor
fi

# Baixar dependências
echo "📦 Baixando dependências..."
go mod download

# Compilar com -mod=mod
echo "🔨 Compilando..."
go build -mod=mod -o datajud-service cmd/server/main.go

# Executar
echo "🎯 Executando DataJud Service..."
./datajud-service