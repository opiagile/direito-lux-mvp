#!/bin/bash

echo "🤖 AI Service - Local Development"
echo "================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Check if we're in the right directory
if [ ! -f "app/main.py" ]; then
    echo -e "${RED}❌ Execute este script do diretório ai-service${NC}"
    echo "   cd services/ai-service"
    echo "   ./run_local.sh"
    exit 1
fi

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo -e "${YELLOW}📦 Criando ambiente virtual...${NC}"
    python3 -m venv venv
    echo "✅ Ambiente virtual criado"
fi

# Activate virtual environment
echo -e "${YELLOW}🔌 Ativando ambiente virtual...${NC}"
source venv/bin/activate

# Install lightweight dependencies
echo -e "${YELLOW}📚 Instalando dependências leves...${NC}"
pip install --upgrade pip
pip install -r requirements.txt

echo "✅ Dependências instaladas"

# Set environment variables for local development
export ENVIRONMENT=development
export DEPLOYMENT_MODE=local
export PORT=8000
export DEBUG=true

echo ""
echo -e "${GREEN}🚀 Iniciando AI Service (Modo Local)${NC}"
echo -e "${BLUE}   URL: http://localhost:8000${NC}"
echo -e "${BLUE}   Docs: http://localhost:8000/docs${NC}"
echo -e "${BLUE}   Health: http://localhost:8000/health${NC}"
echo ""
echo -e "${YELLOW}💡 Modo: Desenvolvimento Local${NC}"
echo -e "${YELLOW}🔗 AI Pesado: Delegado para GCP${NC}"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Run the service
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload --log-level debug