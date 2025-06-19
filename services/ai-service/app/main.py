"""Main FastAPI application."""

from contextlib import asynccontextmanager
from typing import Any

from fastapi import FastAPI, Request, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from prometheus_client import make_asgi_app

from app.api import health, jurisprudence, analysis, generation
from app.core.config import settings
from app.core.exceptions import AIServiceError
from app.core.logging import get_logger, setup_logging
from app.db.database import init_db, close_db
from app.services.cache import init_cache, close_cache
from app.services.vector_store import init_vector_store, close_vector_store

# Setup logging
setup_logging()
logger = get_logger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    logger.info("Starting AI Service", 
                service_name=settings.service_name,
                version=settings.version,
                environment=settings.environment)
    
    # Initialize services
    await init_db()
    await init_cache()
    await init_vector_store()
    
    logger.info("AI Service started successfully")
    
    yield
    
    # Cleanup
    logger.info("Shutting down AI Service")
    await close_vector_store()
    await close_cache()
    await close_db()
    logger.info("AI Service shutdown complete")


# Create FastAPI app
app = FastAPI(
    title="Direito Lux AI Service",
    description="AI-powered jurisprudence analysis and document generation service",
    version=settings.version,
    lifespan=lifespan,
    docs_url="/docs" if settings.environment != "production" else None,
    redoc_url="/redoc" if settings.environment != "production" else None,
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.cors_origins,
    allow_credentials=settings.cors_allow_credentials,
    allow_methods=settings.cors_allow_methods,
    allow_headers=settings.cors_allow_headers,
)


# Exception handlers
@app.exception_handler(AIServiceError)
async def ai_service_error_handler(request: Request, exc: AIServiceError):
    """Handle AI Service specific errors."""
    return JSONResponse(
        status_code=status.HTTP_400_BAD_REQUEST,
        content={
            "error": {
                "code": exc.code,
                "message": exc.message,
                "details": exc.details,
            }
        }
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    """Handle general exceptions."""
    logger.error("Unhandled exception", exc_info=exc)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={
            "error": {
                "code": "INTERNAL_SERVER_ERROR",
                "message": "An unexpected error occurred",
            }
        }
    )


# Mount Prometheus metrics
if settings.prometheus_enabled:
    metrics_app = make_asgi_app()
    app.mount("/metrics", metrics_app)


# Include routers
app.include_router(health.router, tags=["health"])
app.include_router(
    jurisprudence.router,
    prefix="/api/v1/jurisprudence",
    tags=["jurisprudence"]
)
app.include_router(
    analysis.router,
    prefix="/api/v1/analysis",
    tags=["analysis"]
)
app.include_router(
    generation.router,
    prefix="/api/v1/generation",
    tags=["generation"]
)


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": settings.service_name,
        "version": settings.version,
        "environment": settings.environment,
        "status": "running"
    }