"""Document generation endpoints for legal documents."""

import time
from typing import Dict, Any, List, Optional
from uuid import UUID

from fastapi import APIRouter, HTTPException, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.logging import get_logger
from app.db.database import get_session
from app.models.generation import (
    DocumentGenerationRequest,
    DocumentGenerationResponse,
    DocumentType
)

router = APIRouter()
logger = get_logger(__name__)


@router.post("/generate-document", response_model=DocumentGenerationResponse)
async def generate_document(
    request: DocumentGenerationRequest,
    session: AsyncSession = Depends(get_session)
) -> DocumentGenerationResponse:
    """Generate a legal document."""
    start_time = time.time()
    
    try:
        # Simple document generation implementation
        generated_content = f"DOCUMENTO GERADO - {request.document_type.value.upper()}\n\n"
        generated_content += f"Tipo: {request.document_type.value}\n"
        generated_content += f"Gerado em: {time.strftime('%Y-%m-%d %H:%M:%S')}\n\n"
        
        # Add case data if provided
        if hasattr(request, 'case_data') and request.case_data:
            generated_content += "DADOS DO CASO:\n"
            for key, value in request.case_data.items():
                generated_content += f"{key}: {value}\n"
        
        processing_time = int((time.time() - start_time) * 1000)
        
        return DocumentGenerationResponse(
            document_id=UUID("12345678-1234-5678-9012-123456789012"),
            document_type=request.document_type,
            title=f"Documento {request.document_type.value}",
            sections=[],
            generated_content=generated_content,
            word_count=len(generated_content.split()),
            estimated_reading_time=max(1, len(generated_content.split()) // 200),
            processing_time_ms=processing_time
        )
        
    except Exception as e:
        logger.error("Document generation failed", error=str(e))
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/document-types")
async def list_document_types() -> Dict[str, List[str]]:
    """List available document types for generation."""
    return {
        "document_types": [dt.value for dt in DocumentType],
        "styles": ["formal", "technical", "simple", "persuasive"],
        "formats": ["html", "markdown", "docx", "pdf"]
    }


@router.get("/templates")
async def list_templates(
    document_type: Optional[str] = Query(None),
    language: str = Query("pt-BR")
) -> Dict[str, Any]:
    """List available document templates."""
    try:
        # Simple template listing
        templates = [
            {
                "id": "basic_contract",
                "name": "Contrato Básico",
                "document_type": "contrato",
                "description": "Template básico para contratos"
            },
            {
                "id": "initial_petition",
                "name": "Petição Inicial",
                "document_type": "peticao_inicial",
                "description": "Template para petições iniciais"
            }
        ]
        
        if document_type:
            templates = [t for t in templates if t["document_type"] == document_type]
        
        return {
            "templates": templates,
            "total_count": len(templates),
            "language": language
        }
        
    except Exception as e:
        logger.error("Template listing failed", error=str(e))
        raise HTTPException(status_code=500, detail=str(e))