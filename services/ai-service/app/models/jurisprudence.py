"""Jurisprudence domain models."""

from datetime import datetime
from enum import Enum
from typing import List, Optional, Dict, Any
from uuid import UUID

from pydantic import BaseModel, Field, validator


class CourtType(str, Enum):
    """Court type enumeration."""
    STF = "STF"
    STJ = "STJ"
    TST = "TST"
    TSE = "TSE"
    STM = "STM"
    TRF1 = "TRF1"
    TRF2 = "TRF2"
    TRF3 = "TRF3"
    TRF4 = "TRF4"
    TRF5 = "TRF5"
    TRF6 = "TRF6"
    TJ = "TJ"  # Generic for state courts
    TRT = "TRT"  # Generic for labor courts
    TRE = "TRE"  # Generic for electoral courts


class DecisionType(str, Enum):
    """Decision type enumeration."""
    ACORDAO = "acordao"
    DECISAO_MONOCRATICA = "decisao_monocratica"
    DESPACHO = "despacho"
    SENTENCA = "sentenca"
    SUMULA = "sumula"
    SUMULA_VINCULANTE = "sumula_vinculante"
    REPERCUSSAO_GERAL = "repercussao_geral"
    RECURSO_REPETITIVO = "recurso_repetitivo"


class LegalDecision(BaseModel):
    """Legal decision model."""
    
    id: UUID
    court_name: str
    court_type: CourtType
    decision_date: datetime
    publication_date: Optional[datetime] = None
    process_number: str
    case_type: str
    legal_subject: str
    decision_type: DecisionType
    rapporteur: Optional[str] = None
    decision_text: str
    summary: Optional[str] = None
    legal_references: List[str] = Field(default_factory=list)
    cited_cases: List[str] = Field(default_factory=list)
    keywords: List[str] = Field(default_factory=list)
    embedding: Optional[List[float]] = None
    metadata: Dict[str, Any] = Field(default_factory=dict)
    source_url: Optional[str] = None
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)
    
    class Config:
        """Pydantic configuration."""
        json_encoders = {
            UUID: str,
            datetime: lambda v: v.isoformat()
        }


class JurisprudenceSearchRequest(BaseModel):
    """Jurisprudence search request model."""
    
    query: str = Field(..., min_length=3, max_length=5000)
    court_types: Optional[List[CourtType]] = None
    decision_types: Optional[List[DecisionType]] = None
    date_from: Optional[datetime] = None
    date_to: Optional[datetime] = None
    legal_subjects: Optional[List[str]] = None
    max_results: int = Field(default=20, ge=1, le=100)
    similarity_threshold: float = Field(default=0.7, ge=0.0, le=1.0)
    include_embeddings: bool = False
    tenant_id: Optional[UUID] = None
    
    @validator("date_to")
    def validate_date_range(cls, v, values):
        """Validate date range."""
        if v and "date_from" in values and values["date_from"]:
            if v < values["date_from"]:
                raise ValueError("date_to must be after date_from")
        return v


class JurisprudenceSearchResult(BaseModel):
    """Jurisprudence search result model."""
    
    decision: LegalDecision
    similarity_score: float
    highlights: List[str] = Field(default_factory=list)
    relevance_explanation: Optional[str] = None


class JurisprudenceSearchResponse(BaseModel):
    """Jurisprudence search response model."""
    
    query: str
    total_results: int
    results: List[JurisprudenceSearchResult]
    search_metadata: Dict[str, Any] = Field(default_factory=dict)
    processing_time_ms: int


class SimilarityDimension(str, Enum):
    """Similarity dimension enumeration."""
    SEMANTIC = "semantic"
    LEGAL = "legal"
    PROCEDURAL = "procedural"
    CONTEXTUAL = "contextual"


class CaseSimilarityRequest(BaseModel):
    """Case similarity analysis request."""
    
    case_a: Dict[str, Any]  # Process data from Process Service
    case_b: Optional[Dict[str, Any]] = None  # Optional second case
    case_list: Optional[List[Dict[str, Any]]] = None  # Or list of cases
    similarity_dimensions: List[SimilarityDimension] = Field(
        default=[SimilarityDimension.SEMANTIC, SimilarityDimension.LEGAL]
    )
    include_explanation: bool = True
    tenant_id: Optional[UUID] = None
    
    @validator("case_list")
    def validate_case_input(cls, v, values):
        """Validate case input."""
        if not v and not values.get("case_b"):
            raise ValueError("Either case_b or case_list must be provided")
        return v


class SimilarityScore(BaseModel):
    """Similarity score breakdown."""
    
    dimension: SimilarityDimension
    score: float = Field(ge=0.0, le=1.0)
    explanation: Optional[str] = None
    contributing_factors: List[str] = Field(default_factory=list)


class CaseSimilarityResult(BaseModel):
    """Case similarity analysis result."""
    
    case_id: str
    overall_similarity: float = Field(ge=0.0, le=1.0)
    dimension_scores: List[SimilarityScore]
    legal_precedent_match: bool = False
    recommended_action: Optional[str] = None
    similar_decisions: List[LegalDecision] = Field(default_factory=list)


class CaseSimilarityResponse(BaseModel):
    """Case similarity analysis response."""
    
    base_case_id: str
    similarity_results: List[CaseSimilarityResult]
    analysis_metadata: Dict[str, Any] = Field(default_factory=dict)
    processing_time_ms: int