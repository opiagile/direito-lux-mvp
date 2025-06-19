#!/bin/bash

# Script para iniciar o ambiente de desenvolvimento do MCP Service

set -e

echo "🚀 Iniciando MCP Service - Ambiente de Desenvolvimento"

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Por favor, inicie o Docker primeiro."
    exit 1
fi

# Diretório do projeto
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

echo "📁 Diretório do projeto: $PROJECT_DIR"

# Carregar variáveis de ambiente
if [ -f ".env.development" ]; then
    echo "📋 Carregando variáveis de ambiente..."
    export $(cat .env.development | grep -v '^#' | xargs)
else
    echo "⚠️  Arquivo .env.development não encontrado. Usando valores padrão."
fi

# Parar containers existentes se estiverem rodando
echo "🛑 Parando containers existentes..."
docker-compose -f docker-compose.dev.yml down --remove-orphans || true

# Limpar volumes se solicitado
if [ "$1" = "--clean" ]; then
    echo "🧹 Removendo volumes existentes..."
    docker-compose -f docker-compose.dev.yml down -v
    docker volume prune -f
fi

# Subir a infraestrutura (PostgreSQL, Redis, RabbitMQ, Jaeger)
echo "🏗️  Subindo infraestrutura de desenvolvimento..."
docker-compose -f docker-compose.dev.yml up -d

# Aguardar os serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."

# Aguardar PostgreSQL
echo "   • PostgreSQL..."
until docker-compose -f docker-compose.dev.yml exec postgres pg_isready -U mcp_user -d direito_lux_mcp; do
    sleep 2
done

# Aguardar Redis
echo "   • Redis..."
until docker-compose -f docker-compose.dev.yml exec redis redis-cli --raw incr ping; do
    sleep 2
done

# Aguardar RabbitMQ
echo "   • RabbitMQ..."
until docker-compose -f docker-compose.dev.yml exec rabbitmq rabbitmq-diagnostics -q ping; do
    sleep 2
done

echo "✅ Infraestrutura pronta!"

# Verificar se o Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado. Por favor, instale o Go primeiro."
    exit 1
fi

# Atualizar dependências
echo "📦 Atualizando dependências Go..."
go mod tidy

# Executar migrações do banco (opcional)
if [ "$2" = "--migrate" ]; then
    echo "🗄️  Executando migrações do banco..."
    # TODO: Implementar migrações com golang-migrate
    echo "   Migrações serão executadas automaticamente no startup"
fi

# Compilar o projeto
echo "🔨 Compilando MCP Service..."
go build -o bin/mcp-service ./cmd/main.go

if [ $? -ne 0 ]; then
    echo "❌ Falha na compilação!"
    exit 1
fi

echo "✅ Compilação concluída!"

# Mostrar informações dos serviços
echo ""
echo "📊 Serviços disponíveis:"
echo "   • PostgreSQL:  localhost:5434"
echo "   • Redis:       localhost:6380"
echo "   • RabbitMQ:    localhost:5673 (management: http://localhost:15673)"
echo "   • Jaeger UI:   http://localhost:16687"
echo "   • Metrics:     http://localhost:9084/metrics"
echo ""

# Iniciar o serviço
if [ "$3" = "--no-start" ]; then
    echo "🏁 Infraestrutura pronta. Use 'go run cmd/main.go' para iniciar o serviço."
else
    echo "🚀 Iniciando MCP Service..."
    ./bin/mcp-service
fi