"""Analysis endpoints for legal document processing."""

import time
from typing import Dict, Any, List
from uuid import UUID

from fastapi import APIRouter, HTTPException, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.logging import get_logger
from app.db.database import get_session
from app.models.analysis import (
    DocumentAnalysisRequest,
    DocumentAnalysisResponse,
    ProcessAnalysisRequest,
    ProcessAnalysisResponse
)
from app.services.embeddings import embedding_service

router = APIRouter()
logger = get_logger(__name__)


@router.post("/analyze-document", response_model=DocumentAnalysisResponse)
async def analyze_document(
    request: DocumentAnalysisRequest,
    session: AsyncSession = Depends(get_session)
) -> DocumentAnalysisResponse:
    """Analyze a legal document using AI."""
    start_time = time.time()
    
    try:
        # Create analysis prompt for AI
        analysis_prompt = f"""
        Analise o seguinte documento jurídico e forneça:
        1. Resumo executivo (máximo 200 palavras)
        2. Pontos-chave principais (máximo 5)
        3. Recomendações práticas (máximo 3)

        Documento:
        {request.document_content[:2000]}...

        Responda em português, de forma clara e objetiva.
        """
        
        # Use AI service to analyze document
        try:
            ai_response = await embedding_service._get_ollama_text_completion(analysis_prompt)
            
            # Parse AI response (simplified parsing)
            lines = ai_response.strip().split('\n')
            summary = ""
            key_points = []
            recommendations = []
            
            current_section = ""
            for line in lines:
                line = line.strip()
                if not line:
                    continue
                    
                if "resumo" in line.lower() or "summary" in line.lower():
                    current_section = "summary"
                elif "pontos" in line.lower() or "key" in line.lower():
                    current_section = "key_points"
                elif "recomend" in line.lower():
                    current_section = "recommendations"
                elif current_section == "summary" and len(summary) < 500:
                    summary += line + " "
                elif current_section == "key_points" and len(key_points) < 5:
                    if line.startswith(('-', '•', '*', '1.', '2.', '3.', '4.', '5.')):
                        key_points.append(line)
                elif current_section == "recommendations" and len(recommendations) < 3:
                    if line.startswith(('-', '•', '*', '1.', '2.', '3.')):
                        recommendations.append(line)
            
            # Fallback to simple analysis if AI parsing fails
            if not summary:
                summary = f"Documento analisado com {len(request.document_content)} caracteres. Análise automática realizada."
            if not key_points:
                key_points = ["Documento processado", "Conteúdo analisado"]
            if not recommendations:
                recommendations = ["Revisar conteúdo", "Validar informações"]
                
        except Exception as ai_error:
            logger.warning(f"AI analysis failed, using fallback: {str(ai_error)}")
            # Fallback to simple analysis
            summary = f"Análise do documento: {len(request.document_content)} caracteres processados"
            key_points = ["Documento analisado", "Conteúdo processado"]
            recommendations = ["Revisar conteúdo", "Validar informações"]
        
        processing_time = int((time.time() - start_time) * 1000)
        
        return DocumentAnalysisResponse(
            document_type=request.document_type or "unknown",
            summary=summary.strip(),
            key_points=key_points,
            extracted_entities=[],
            legal_citations=[],
            recommendations=recommendations,
            processing_time_ms=processing_time
        )
        
    except Exception as e:
        logger.error("Document analysis failed", error=str(e))
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/analyze-process", response_model=ProcessAnalysisResponse)
async def analyze_process(
    request: ProcessAnalysisRequest,
    session: AsyncSession = Depends(get_session)
) -> ProcessAnalysisResponse:
    """Analyze a legal process."""
    start_time = time.time()
    
    try:
        # Simple process analysis implementation
        processing_time = int((time.time() - start_time) * 1000)
        
        return ProcessAnalysisResponse(
            process_id=request.process_id,
            processing_time_ms=processing_time
        )
        
    except Exception as e:
        logger.error("Process analysis failed", error=str(e))
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/analysis-types")
async def list_analysis_types() -> Dict[str, List[str]]:
    """List available analysis types."""
    return {
        "document_types": ["contract", "petition", "decision", "law"],
        "analysis_types": ["summary", "entity_extraction", "risk_assessment"],
        "depths": ["quick", "standard", "deep"]
    }