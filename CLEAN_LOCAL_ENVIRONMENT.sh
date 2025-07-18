#!/bin/bash

# 🧹 LIMPEZA COMPLETA DO AMBIENTE LOCAL
# Prepara para desenvolvimento do novo MicroSaaS

set -e

echo "🧹 LIMPEZA DO AMBIENTE LOCAL"
echo "============================"
echo ""

# Cores
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${YELLOW}Esta ação irá:${NC}"
echo "   • Parar todos os containers Docker"
echo "   • Remover todos os containers e imagens"
echo "   • Limpar volumes e networks Docker"
echo "   • Remover arquivos temporários"
echo "   • Preparar ambiente para novo desenvolvimento"
echo ""

read -p "Continuar com a limpeza? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Operação cancelada."
    exit 1
fi

echo -e "${YELLOW}ETAPA 1: Parando todos os containers${NC}"
docker stop $(docker ps -aq) 2>/dev/null || true

echo -e "${YELLOW}ETAPA 2: Removendo containers${NC}"
docker rm $(docker ps -aq) 2>/dev/null || true

echo -e "${YELLOW}ETAPA 3: Removendo imagens do projeto${NC}"
docker rmi $(docker images | grep direito-lux | awk '{print $3}') 2>/dev/null || true

echo -e "${YELLOW}ETAPA 4: Limpando volumes Docker${NC}"
docker volume prune -f 2>/dev/null || true

echo -e "${YELLOW}ETAPA 5: Limpando networks Docker${NC}"
docker network prune -f 2>/dev/null || true

echo -e "${YELLOW}ETAPA 6: Limpando arquivos temporários${NC}"
# Remove logs
find . -name "*.log" -type f -delete 2>/dev/null || true
find . -name "*.pid" -type f -delete 2>/dev/null || true

# Remove binários Go
find . -name "server" -type f -delete 2>/dev/null || true
find . -name "*-service" -type f -delete 2>/dev/null || true
find . -name "*-test" -type f -delete 2>/dev/null || true

# Remove node_modules (se existir)
find . -name "node_modules" -type d -exec rm -rf {} + 2>/dev/null || true

# Remove coverage files
find . -name "coverage.out" -type f -delete 2>/dev/null || true
find . -name "coverage.html" -type f -delete 2>/dev/null || true

echo -e "${YELLOW}ETAPA 7: Limpando cache local${NC}"
# Go mod cache
go clean -modcache 2>/dev/null || true

# NPM cache
npm cache clean --force 2>/dev/null || true

echo -e "${YELLOW}ETAPA 8: Resetando Docker${NC}"
docker system prune -af --volumes 2>/dev/null || true

echo ""
echo -e "${GREEN}✅ LIMPEZA COMPLETA!${NC}"
echo ""
echo "Status:"
echo "   • Docker: Limpo e resetado"
echo "   • Arquivos temporários: Removidos"
echo "   • Cache: Limpo"
echo "   • Ambiente: Pronto para novo desenvolvimento"
echo ""
echo -e "${GREEN}🚀 Pronto para começar o ProcessAlert MicroSaaS!${NC}"