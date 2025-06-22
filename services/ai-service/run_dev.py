#!/usr/bin/env python3
"""
Script de desenvolvimento para AI Service
Executa com reload e configurações de desenvolvimento
"""

import os
import sys
import subprocess
from pathlib import Path

def main():
    # Verificar se está no diretório correto
    if not Path("app/main.py").exists():
        print("❌ Erro: Execute este script do diretório ai-service")
        print("   cd services/ai-service")
        print("   python run_dev.py")
        sys.exit(1)
    
    # Verificar ambiente virtual
    if "venv" not in sys.prefix:
        print("⚠️  Recomendado: ativar ambiente virtual")
        print("   source venv/bin/activate")
        print("")
    
    # Verificar dependências básicas
    try:
        import fastapi
        import uvicorn
        print("✅ Dependências básicas encontradas")
    except ImportError as e:
        print(f"❌ Dependência faltando: {e}")
        print("   pip install fastapi uvicorn")
        sys.exit(1)
    
    # Executar com uvicorn
    print("🚀 Iniciando AI Service em modo desenvolvimento...")
    print("   URL: http://localhost:8000")
    print("   Docs: http://localhost:8000/docs")
    print("")
    
    try:
        subprocess.run([
            "uvicorn", 
            "app.main:app",  # Caminho correto para o app
            "--host", "0.0.0.0",
            "--port", "8000",
            "--reload",
            "--log-level", "debug"
        ])
    except KeyboardInterrupt:
        print("\n👋 AI Service parado")

if __name__ == "__main__":
    main()
