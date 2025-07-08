"""
AI Service - Local Development Version
Lightweight version for local development without heavy ML dependencies
Heavy AI processing is delegated to GCP-deployed AI service
"""

from contextlib import asynccontextmanager
from typing import Dict, List, Optional
import asyncio
import httpx
import os

from fastapi import FastAPI, HTTPException, Request, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from prometheus_client import make_asgi_app

# Simple configuration
class Config:
    SERVICE_NAME = "ai-service-local"
    VERSION = "0.1.0"
    ENVIRONMENT = os.getenv("ENVIRONMENT", "development")
    GCP_AI_SERVICE_URL = os.getenv("GCP_AI_SERVICE_URL", "")
    
    # CORS
    CORS_ORIGINS = ["http://localhost:3000", "http://localhost:8000", "http://localhost:3001"]
    CORS_ALLOW_CREDENTIALS = True
    CORS_ALLOW_METHODS = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    CORS_ALLOW_HEADERS = ["*"]

config = Config()

# Pydantic models
class HealthResponse(BaseModel):
    status: str
    version: str
    environment: str
    mode: str
    gcp_service_available: bool

class AnalysisRequest(BaseModel):
    text: str
    analysis_type: str = "summary"
    tenant_id: Optional[str] = None

class AnalysisResponse(BaseModel):
    result: str
    confidence: float
    analysis_type: str
    processing_mode: str

class JurisprudenceSearchRequest(BaseModel):
    query: str
    limit: int = 10
    similarity_threshold: float = 0.7
    tenant_id: Optional[str] = None

class DocumentGenerationRequest(BaseModel):
    template_type: str
    data: Dict
    tenant_id: Optional[str] = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    print(f"🚀 Starting {config.SERVICE_NAME} v{config.VERSION}")
    print(f"🔧 Environment: {config.ENVIRONMENT}")
    print(f"⚡ Mode: Local Development (Heavy AI → GCP)")
    
    yield
    
    print(f"👋 Shutting down {config.SERVICE_NAME}")


# Create FastAPI app
app = FastAPI(
    title="Direito Lux AI Service (Local)",
    description="AI Service - Local development version with GCP delegation",
    version=config.VERSION,
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=config.CORS_ORIGINS,
    allow_credentials=config.CORS_ALLOW_CREDENTIALS,
    allow_methods=config.CORS_ALLOW_METHODS,
    allow_headers=config.CORS_ALLOW_HEADERS,
)

# HTTP client for GCP communication
async def get_http_client():
    return httpx.AsyncClient(timeout=30.0)

async def check_gcp_service():
    """Check if GCP AI service is available."""
    if not config.GCP_AI_SERVICE_URL:
        return False
    
    try:
        async with httpx.AsyncClient(timeout=5.0) as client:
            response = await client.get(f"{config.GCP_AI_SERVICE_URL}/health")
            return response.status_code == 200
    except:
        return False


# Routes
@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": config.SERVICE_NAME,
        "version": config.VERSION,
        "environment": config.ENVIRONMENT,
        "mode": "local_development",
        "status": "running",
        "ai_processing": "delegated_to_gcp"
    }


@app.get("/health", response_model=HealthResponse)
async def health():
    """Health check endpoint."""
    gcp_available = await check_gcp_service()
    
    return HealthResponse(
        status="healthy",
        version=config.VERSION,
        environment=config.ENVIRONMENT,
        mode="local_development",
        gcp_service_available=gcp_available
    )


@app.post("/api/v1/analysis", response_model=AnalysisResponse)
async def analyze_text(request: AnalysisRequest):
    """
    Text analysis endpoint.
    Delegates to GCP service or provides mock response.
    """
    # Try GCP service first
    if config.GCP_AI_SERVICE_URL:
        try:
            async with await get_http_client() as client:
                response = await client.post(
                    f"{config.GCP_AI_SERVICE_URL}/api/v1/analysis",
                    json=request.dict(),
                    timeout=30.0
                )
                if response.status_code == 200:
                    result = response.json()
                    result["processing_mode"] = "gcp_service"
                    return AnalysisResponse(**result)
        except Exception as e:
            print(f"GCP service error: {e}")
    
    # Fallback to local mock
    await asyncio.sleep(0.1)  # Simulate processing
    
    if request.analysis_type == "summary":
        result = f"Resumo: {request.text[:100]}... (versão mock para desenvolvimento local)"
    elif request.analysis_type == "sentiment":
        result = "Sentimento: Neutro (0.5) - análise mock"
    elif request.analysis_type == "keywords":
        result = "Palavras-chave: direito, processo, análise - mock"
    else:
        result = f"Análise {request.analysis_type} - resultado mock para desenvolvimento"
    
    return AnalysisResponse(
        result=result,
        confidence=0.7,
        analysis_type=request.analysis_type,
        processing_mode="local_mock"
    )


@app.post("/api/v1/jurisprudence/search")
async def search_jurisprudence(request: JurisprudenceSearchRequest):
    """
    Jurisprudence search endpoint.
    Delegates to GCP service or provides mock response.
    """
    # Try GCP service first
    if config.GCP_AI_SERVICE_URL:
        try:
            async with await get_http_client() as client:
                response = await client.post(
                    f"{config.GCP_AI_SERVICE_URL}/api/v1/jurisprudence/search",
                    json=request.dict(),
                    timeout=30.0
                )
                if response.status_code == 200:
                    result = response.json()
                    result["processing_mode"] = "gcp_service"
                    return result
        except Exception as e:
            print(f"GCP service error: {e}")
    
    # Fallback to local mock
    await asyncio.sleep(0.2)
    
    return {
        "results": [
            {
                "id": "mock-1",
                "title": f"Jurisprudência relacionada a: {request.query}",
                "summary": "Resultado mock para desenvolvimento local",
                "court": "STF",
                "date": "2024-01-01",
                "similarity": 0.85,
                "url": "https://example.com/mock-1"
            },
            {
                "id": "mock-2", 
                "title": f"Precedente sobre: {request.query}",
                "summary": "Outro resultado mock para teste",
                "court": "STJ",
                "date": "2024-01-02",
                "similarity": 0.78,
                "url": "https://example.com/mock-2"
            }
        ],
        "total": 2,
        "processing_mode": "local_mock",
        "query": request.query
    }


@app.post("/api/v1/generation/document")
async def generate_document(request: DocumentGenerationRequest):
    """
    Document generation endpoint.
    Delegates to GCP service or provides mock response.
    """
    # Try GCP service first
    if config.GCP_AI_SERVICE_URL:
        try:
            async with await get_http_client() as client:
                response = await client.post(
                    f"{config.GCP_AI_SERVICE_URL}/api/v1/generation/document",
                    json=request.dict(),
                    timeout=60.0
                )
                if response.status_code == 200:
                    result = response.json()
                    result["processing_mode"] = "gcp_service"
                    return result
        except Exception as e:
            print(f"GCP service error: {e}")
    
    # Fallback to local mock
    await asyncio.sleep(0.5)
    
    return {
        "document_id": "mock-doc-123",
        "template_type": request.template_type,
        "status": "generated",
        "content": f"Documento {request.template_type} gerado em modo mock para desenvolvimento local.\n\nDados fornecidos: {request.data}",
        "processing_mode": "local_mock",
        "created_at": "2024-01-01T12:00:00Z"
    }


@app.get("/api/v1/models")
async def list_models():
    """List available AI models."""
    return {
        "models": [
            {"name": "summary-model", "status": "mock", "type": "local_development"},
            {"name": "sentiment-model", "status": "mock", "type": "local_development"},
            {"name": "keyword-model", "status": "mock", "type": "local_development"},
            {"name": "jurisprudence-search", "status": "mock", "type": "local_development"},
            {"name": "document-generator", "status": "mock", "type": "local_development"}
        ],
        "processing_mode": "local_mock",
        "gcp_service_available": await check_gcp_service()
    }


# Mount Prometheus metrics
metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)


# Exception handler
@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    """Handle general exceptions."""
    print(f"Error: {exc}")
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={
            "error": {
                "code": "INTERNAL_SERVER_ERROR",
                "message": "An unexpected error occurred in local AI service",
                "mode": "local_development"
            }
        }
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)