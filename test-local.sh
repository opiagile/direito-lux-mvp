#!/bin/bash

echo "🚀 Testando Direito Lux - Validação Completa..."
echo "=================================================="

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Por favor, inicie o Docker Desktop."
    exit 1
fi

echo "✅ Docker está rodando"

# Subir os serviços
echo ""
echo "📦 Iniciando serviços com Docker Compose..."
docker-compose up -d

# Aguardar serviços subirem
echo "⏳ Aguardando serviços iniciarem (30 segundos)..."
sleep 30

echo ""
echo "🔍 Verificando status dos containers..."
docker-compose ps

echo ""
echo "🏥 Testando Health Checks dos Microserviços..."
echo "================================================="

# Função para testar health check
test_health() {
    local service_name=$1
    local url=$2
    
    echo -n "Testing $service_name... "
    if curl -f -s "$url" > /dev/null 2>&1; then
        echo "✅ OK"
        return 0
    else
        echo "❌ FALHOU"
        return 1
    fi
}

# Testar cada serviço
test_health "Auth Service" "http://localhost:8081/health"
test_health "Tenant Service" "http://localhost:8082/health" 
test_health "Process Service" "http://localhost:8083/health"
test_health "DataJud Service" "http://localhost:8084/health"

echo ""
echo "🗃️ Testando Infraestrutura..."
echo "=============================="

# Testar PostgreSQL
echo -n "Testing PostgreSQL... "
if docker-compose exec -T postgres psql -U direito_lux -d direito_lux_dev -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ OK"
    
    echo -n "Verificando tabelas principais... "
    table_count=$(docker-compose exec -T postgres psql -U direito_lux -d direito_lux_dev -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | tr -d ' ' | head -1)
    if [ "$table_count" -gt 0 ]; then
        echo "✅ $table_count tabelas encontradas"
    else
        echo "⚠️ Nenhuma tabela encontrada"
    fi
else
    echo "❌ FALHOU"
fi

# Testar Redis
echo -n "Testing Redis... "
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ OK"
else
    echo "❌ FALHOU"
fi

# Testar RabbitMQ
echo -n "Testing RabbitMQ... "
if curl -f -s "http://localhost:15672" > /dev/null 2>&1; then
    echo "✅ OK"
else
    echo "❌ FALHOU"
fi

echo ""
echo "📊 Informações dos Serviços..."
echo "=============================="

# Listar containers ativos
echo "🐳 Containers ativos:"
docker-compose ps --format "table {{.Name}}\t{{.State}}\t{{.Ports}}"

echo ""
echo "📈 URLs de Acesso:"
echo "=================="
echo "🔐 Auth Service:     http://localhost:8081"
echo "🏢 Tenant Service:   http://localhost:8082" 
echo "📋 Process Service:  http://localhost:8083"
echo "🔗 DataJud Service:  http://localhost:8084"
echo "🗄️ PostgreSQL:       localhost:5432"
echo "🚀 Redis:            localhost:6379"
echo "🐰 RabbitMQ:         http://localhost:15672"
echo "🔍 Prometheus:       http://localhost:9090"
echo "📊 Grafana:          http://localhost:3000"

echo ""
echo "🎯 Teste Funcional Básico..."
echo "==========================="

# Testar endpoint de informações (se existir)
echo -n "Testing Auth Service info... "
if curl -f -s "http://localhost:8081/api/v1/info" > /dev/null 2>&1; then
    echo "✅ OK"
else
    echo "⚠️ Endpoint não disponível (normal se não implementado)"
fi

echo ""
echo "📋 Resumo do Teste:"
echo "=================="

# Contar serviços funcionais
functional_services=0
total_services=4

echo "📊 Status dos Microserviços:"
if curl -f -s "http://localhost:8081/health" > /dev/null 2>&1; then
    echo "  ✅ Auth Service"
    ((functional_services++))
else
    echo "  ❌ Auth Service"
fi

if curl -f -s "http://localhost:8082/health" > /dev/null 2>&1; then
    echo "  ✅ Tenant Service"
    ((functional_services++))
else
    echo "  ❌ Tenant Service"
fi

if curl -f -s "http://localhost:8083/health" > /dev/null 2>&1; then
    echo "  ✅ Process Service"
    ((functional_services++))
else
    echo "  ❌ Process Service"
fi

if curl -f -s "http://localhost:8084/health" > /dev/null 2>&1; then
    echo "  ✅ DataJud Service"
    ((functional_services++))
else
    echo "  ❌ DataJud Service"
fi

echo ""
echo "🏆 Resultado Final:"
echo "  📈 $functional_services de $total_services microserviços funcionais"

if [ $functional_services -eq $total_services ]; then
    echo "  🎉 SUCESSO! Todos os serviços estão funcionando"
    echo "  🚀 Pronto para próxima fase de desenvolvimento"
elif [ $functional_services -gt 2 ]; then
    echo "  ⚠️ PARCIAL! Maioria dos serviços funcionando"
    echo "  🔧 Verificar logs dos serviços com falha"
else
    echo "  ❌ FALHA! Muitos serviços com problema"
    echo "  🆘 Verificar configuração do ambiente"
fi

echo ""
echo "📝 Comandos úteis:"
echo "=================="
echo "🔍 Ver logs:           docker-compose logs -f [service-name]"
echo "🗄️ Conectar PostgreSQL: docker-compose exec postgres psql -U direito_lux -d direito_lux_dev"
echo "🛑 Parar tudo:         docker-compose down"
echo "🧹 Limpar volumes:     docker-compose down -v"

echo ""
echo "✅ Teste concluído! $(date)"