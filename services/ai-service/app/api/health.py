"""Health check endpoints."""

from datetime import datetime
from typing import Dict, Any

from fastapi import APIRouter

from app.core.config import settings
from app.services.cache import cache_service
from app.services.vector_store import vector_store

router = APIRouter()


@router.get("/health")
async def health_check() -> Dict[str, Any]:
    """Health check endpoint."""
    return {
        "status": "healthy",
        "timestamp": datetime.utcnow().isoformat(),
        "service": settings.service_name,
        "version": settings.version,
        "environment": settings.environment
    }


@router.get("/ready")
async def readiness_check() -> Dict[str, Any]:
    """Readiness check endpoint."""
    checks = {
        "cache": await _check_cache(),
        "vector_store": await _check_vector_store(),
    }
    
    all_ready = all(checks.values())
    
    return {
        "status": "ready" if all_ready else "not_ready",
        "timestamp": datetime.utcnow().isoformat(),
        "service": settings.service_name,
        "checks": checks
    }


async def _check_cache() -> bool:
    """Check cache connectivity."""
    try:
        if cache_service._redis:
            await cache_service._redis.ping()
        return True
    except Exception:
        return False


async def _check_vector_store() -> bool:
    """Check vector store status."""
    try:
        return vector_store.initialized
    except Exception:
        return False