#!/bin/bash

# Script para compilar todos os microserviços e verificar integridade

set -e

echo "🔨 Compilando todos os microserviços Direito Lux..."

SERVICES=("auth-service" "tenant-service" "process-service" "datajud-service")
BUILD_SUCCESS=true

for service in "${SERVICES[@]}"; do
    echo ""
    echo "📦 Compilando $service..."
    cd "services/$service"
    
    # Verificar se go.mod tem o nome correto
    SERVICE_NAME=$(grep "^module" go.mod | awk '{print $2}' | cut -d'/' -f3)
    if [ "$SERVICE_NAME" != "$service" ]; then
        echo "⚠️  AVISO: go.mod contém '$SERVICE_NAME' mas deveria ser '$service'"
    fi
    
    # Limpar e baixar dependências
    echo "   🧹 Limpando dependências..."
    go mod tidy
    
    # Verificar imports problemáticos
    echo "   🔍 Verificando imports..."
    if grep -r "logger\.Sugar()\.Desugar()\.Core()" . 2>/dev/null; then
        echo "   ❌ ERRO: Logger middleware incorreto encontrado!"
        BUILD_SUCCESS=false
    fi
    
    if grep -r "return gin\.Next$" . 2>/dev/null; then
        echo "   ❌ ERRO: gin.Next incorreto encontrado!"
        BUILD_SUCCESS=false
    fi
    
    # Compilar
    echo "   🔨 Compilando..."
    if go build ./cmd/server; then
        echo "   ✅ $service compilado com sucesso"
        rm -f server  # Limpar executável
    else
        echo "   ❌ FALHA na compilação do $service"
        BUILD_SUCCESS=false
    fi
    
    cd ../..
done

echo ""
if [ "$BUILD_SUCCESS" = true ]; then
    echo "🎉 Todos os serviços compilaram com sucesso!"
    echo ""
    echo "📋 Próximos passos:"
    echo "1. ./start-services.sh - Para iniciar todos os serviços"
    echo "2. ./test-local.sh - Para testar funcionamento"
    exit 0
else
    echo "❌ Falhas na compilação encontradas!"
    echo ""
    echo "🔧 Verificações necessárias:"
    echo "1. Conferir imports ausentes (fmt, time, runtime, os)"
    echo "2. Corrigir middleware incorreto (gin.Next, logger.Core)"
    echo "3. Verificar signatures de função (LogError)"
    echo "4. Consultar DIRETRIZES_DESENVOLVIMENTO.md"
    exit 1
fi