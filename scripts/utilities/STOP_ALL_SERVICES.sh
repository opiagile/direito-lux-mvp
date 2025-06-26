#!/bin/bash

# =============================================================================
# DIREITO LUX - STOP ALL SERVICES
# =============================================================================
# Script para parar todos os serviços de forma limpa
# =============================================================================

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🛑 DIREITO LUX - PARAR SERVIÇOS${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Diretório base
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$BASE_DIR"

echo -e "${YELLOW}📋 Parando todos os serviços Docker...${NC}"

# Parar Docker Compose principal
docker-compose down -v --remove-orphans 2>/dev/null || true

# Parar Docker Compose simples
docker-compose -f docker-compose.simple.yml down -v --remove-orphans 2>/dev/null || true

# Parar Docker Compose de desenvolvimento
docker-compose -f services/docker-compose.dev.yml down -v --remove-orphans 2>/dev/null || true

# Remover containers órfãos
echo -e "${YELLOW}🔍 Removendo containers órfãos...${NC}"
ORPHAN_CONTAINERS=$(docker ps -a -q --filter "label=com.docker.compose.project=direito-lux" 2>/dev/null || true)
if [ ! -z "$ORPHAN_CONTAINERS" ]; then
    docker rm -f $ORPHAN_CONTAINERS 2>/dev/null || true
fi

# Parar frontend se estiver rodando
echo -e "${YELLOW}💻 Verificando frontend...${NC}"
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${YELLOW}Parando frontend na porta 3000...${NC}"
    lsof -ti:3000 | xargs kill -9 2>/dev/null || true
fi

# Limpar rede se não estiver sendo usada
echo -e "${YELLOW}🌐 Limpando rede Docker...${NC}"
docker network rm direito-lux-network 2>/dev/null || true

echo ""
echo -e "${GREEN}✅ Todos os serviços foram parados!${NC}"
echo ""
echo -e "${BLUE}Para iniciar novamente:${NC}"
echo -e "  ${GREEN}./START_CLEAN_ENVIRONMENT.sh${NC}"
echo ""