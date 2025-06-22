"""
Main FastAPI application for AI Service.

For local development: Uses lightweight version (main_local.py)
For GCP production: Uses full ML implementation (main_gcp.py)
"""

import os

# Determine which version to use based on environment
ENVIRONMENT = os.getenv("ENVIRONMENT", "development").lower()
DEPLOYMENT_MODE = os.getenv("DEPLOYMENT_MODE", "local").lower()

if DEPLOYMENT_MODE == "gcp" or ENVIRONMENT == "production":
    # Use full GCP implementation with all ML libraries
    try:
        from app.main_gcp import app  
        print("🚀 Using GCP AI Service (Full ML Implementation)")
    except ImportError:
        print("⚠️  GCP implementation not available, falling back to local")
        from app.main_local import app
else:
    # Use local lightweight implementation
    from app.main_local import app
    print("⚡ Using Local AI Service (Lightweight Implementation)")

# Export the app for uvicorn
__all__ = ["app"]