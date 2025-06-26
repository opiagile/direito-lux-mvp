#!/bin/bash

echo "🧹 LIMPEZA TOTAL DO AMBIENTE DIREITO LUX"
echo "========================================"
echo ""
echo "⚠️  ATENÇÃO: Este script vai REMOVER TUDO!"
echo "- Todos os containers Docker"
echo "- Todos os volumes (dados serão perdidos)"
echo "- Todas as redes Docker"
echo "- Processos PostgreSQL locais"
echo ""
read -p "Tem certeza que deseja continuar? (sim/não): " resposta

if [ "$resposta" != "sim" ]; then
    echo "❌ Operação cancelada"
    exit 0
fi

echo ""
echo "🛑 1. Parando todos os containers..."
docker stop $(docker ps -aq) 2>/dev/null || echo "Nenhum container rodando"

echo ""
echo "🗑️  2. Removendo todos os containers..."
docker rm -f $(docker ps -aq) 2>/dev/null || echo "Nenhum container para remover"

echo ""
echo "💾 3. Removendo todos os volumes..."
docker volume rm $(docker volume ls -q) 2>/dev/null || echo "Nenhum volume para remover"

echo ""
echo "🌐 4. Removendo redes customizadas..."
docker network prune -f 2>/dev/null

echo ""
echo "🧹 5. Limpeza geral do Docker..."
docker system prune -af --volumes

echo ""
echo "🐘 6. Verificando PostgreSQL local..."
# Para macOS
if command -v brew &> /dev/null; then
    if brew services list | grep -q "postgresql.*started"; then
        echo "Parando PostgreSQL local..."
        brew services stop postgresql
    fi
fi

# Para Linux
if command -v systemctl &> /dev/null; then
    if systemctl is-active --quiet postgresql; then
        echo "Parando PostgreSQL local..."
        sudo systemctl stop postgresql
    fi
fi

# Matar processos PostgreSQL órfãos
echo "Procurando processos PostgreSQL..."
pkill -f postgres 2>/dev/null || echo "Nenhum processo PostgreSQL encontrado"

echo ""
echo "🔍 7. Verificando portas em uso..."
echo "Portas que podem estar em uso:"
lsof -i :5432 2>/dev/null || echo "Porta 5432 (PostgreSQL) livre"
lsof -i :6379 2>/dev/null || echo "Porta 6379 (Redis) livre"
lsof -i :5672 2>/dev/null || echo "Porta 5672 (RabbitMQ) livre"
lsof -i :15672 2>/dev/null || echo "Porta 15672 (RabbitMQ Management) livre"
lsof -i :3000 2>/dev/null || echo "Porta 3000 (Frontend) livre"
lsof -i :8080 2>/dev/null || echo "Porta 8080 (API) livre"

echo ""
echo "📁 8. Limpando arquivos temporários..."
# Limpar logs
rm -rf logs/*.log 2>/dev/null
# Limpar node_modules do frontend se existir
rm -rf frontend/node_modules 2>/dev/null
rm -rf frontend/.next 2>/dev/null
# Limpar vendor Go se existir
find . -name "vendor" -type d -exec rm -rf {} + 2>/dev/null

echo ""
echo "✅ LIMPEZA COMPLETA FINALIZADA!"
echo ""
echo "📊 Status do Docker:"
docker ps -a
echo ""
echo "💾 Volumes Docker:"
docker volume ls
echo ""
echo "🌐 Redes Docker:"
docker network ls
echo ""
echo "🎯 Próximos passos:"
echo "1. Execute o novo script de setup quando estiver pronto"
echo "2. Aguarde todos os serviços subirem completamente"
echo "3. Verifique os logs se houver problemas"
echo ""
echo "✨ Ambiente limpo e pronto para um novo início!"