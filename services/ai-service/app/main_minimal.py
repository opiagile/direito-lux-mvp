"""
Versão minimalista do AI Service para desenvolvimento
Sem dependências pesadas de ML
"""

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import asyncio
import json
from typing import Dict, List, Optional

app = FastAPI(
    title="AI Service - Dev Mode",
    description="Direito Lux AI Service (Modo Desenvolvimento)",
    version="0.1.0"
)

# Models básicos
class AnalysisRequest(BaseModel):
    text: str
    analysis_type: str = "summary"

class AnalysisResponse(BaseModel):
    result: str
    confidence: float
    analysis_type: str

class HealthResponse(BaseModel):
    status: str
    version: str
    dependencies: Dict[str, str]

# Endpoints básicos
@app.get("/")
async def root():
    return {"message": "AI Service - Direito Lux", "status": "development"}

@app.get("/health", response_model=HealthResponse)
async def health():
    return HealthResponse(
        status="healthy",
        version="0.1.0",
        dependencies={
            "fastapi": "installed",
            "database": "not_checked",
            "redis": "not_checked",
            "ml_models": "disabled_in_dev"
        }
    )

@app.post("/api/v1/analyze", response_model=AnalysisResponse)
async def analyze_text(request: AnalysisRequest):
    """
    Análise de texto - versão mock para desenvolvimento
    """
    # Simular processamento
    await asyncio.sleep(0.1)
    
    # Mock response baseado no tipo
    if request.analysis_type == "summary":
        result = f"Resumo do texto: {request.text[:100]}..."
    elif request.analysis_type == "sentiment":
        result = "Sentimento: Neutro (0.5)"
    elif request.analysis_type == "keywords":
        result = "Palavras-chave: direito, processo, análise"
    else:
        result = f"Análise {request.analysis_type} do texto fornecido"
    
    return AnalysisResponse(
        result=result,
        confidence=0.85,
        analysis_type=request.analysis_type
    )

@app.get("/api/v1/models")
async def list_models():
    """Lista modelos disponíveis (mock)"""
    return {
        "models": [
            {"name": "summary-model", "status": "loaded", "type": "mock"},
            {"name": "sentiment-model", "status": "loaded", "type": "mock"},
            {"name": "keyword-model", "status": "loaded", "type": "mock"}
        ]
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
