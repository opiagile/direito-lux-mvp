"""Analysis domain models."""

from datetime import datetime
from enum import Enum
from typing import List, Optional, Dict, Any
from uuid import UUID

from pydantic import BaseModel, Field


class AnalysisType(str, Enum):
    """Analysis type enumeration."""
    PROCESS_SUMMARY = "process_summary"
    MOVEMENT_ANALYSIS = "movement_analysis"
    RISK_ASSESSMENT = "risk_assessment"
    OUTCOME_PREDICTION = "outcome_prediction"
    TIMELINE_EXTRACTION = "timeline_extraction"
    ENTITY_EXTRACTION = "entity_extraction"
    SENTIMENT_ANALYSIS = "sentiment_analysis"


class RiskLevel(str, Enum):
    """Risk level enumeration."""
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"


class ProcessAnalysisRequest(BaseModel):
    """Process analysis request model."""
    
    process_id: UUID
    process_data: Dict[str, Any]  # From Process Service
    analysis_types: List[AnalysisType]
    include_jurisprudence: bool = True
    include_recommendations: bool = True
    tenant_id: Optional[UUID] = None
    user_plan: str = Field(default="starter")  # For feature gating


class ExtractedEntity(BaseModel):
    """Extracted entity from text."""
    
    entity_type: str  # person, organization, date, value, law_reference
    value: str
    confidence: float = Field(ge=0.0, le=1.0)
    context: Optional[str] = None
    position: Optional[Dict[str, int]] = None  # start, end positions


class ProcessSummary(BaseModel):
    """Process summary result."""
    
    summary: str
    key_points: List[str]
    current_status: str
    next_steps: List[str]
    important_dates: List[Dict[str, Any]]
    parties_involved: List[Dict[str, str]]


class MovementAnalysis(BaseModel):
    """Movement analysis result."""
    
    movement_id: str
    movement_date: datetime
    movement_type: str
    summary: str
    impact_level: str  # low, medium, high
    requires_action: bool
    action_deadline: Optional[datetime] = None
    extracted_entities: List[ExtractedEntity]
    sentiment: Optional[str] = None


class RiskAssessment(BaseModel):
    """Risk assessment result."""
    
    overall_risk: RiskLevel
    risk_factors: List[Dict[str, Any]]
    mitigation_strategies: List[str]
    confidence: float = Field(ge=0.0, le=1.0)
    based_on_precedents: int


class OutcomePrediction(BaseModel):
    """Outcome prediction result."""
    
    predicted_outcome: str
    probability: float = Field(ge=0.0, le=1.0)
    confidence_level: str  # low, medium, high
    supporting_factors: List[str]
    opposing_factors: List[str]
    similar_cases_analyzed: int
    disclaimer: str = Field(
        default="This is an AI prediction based on historical data and should not be considered legal advice."
    )


class ProcessAnalysisResponse(BaseModel):
    """Process analysis response model."""
    
    process_id: UUID
    analysis_timestamp: datetime = Field(default_factory=datetime.utcnow)
    process_summary: Optional[ProcessSummary] = None
    movement_analyses: List[MovementAnalysis] = Field(default_factory=list)
    risk_assessment: Optional[RiskAssessment] = None
    outcome_prediction: Optional[OutcomePrediction] = None
    extracted_entities: List[ExtractedEntity] = Field(default_factory=list)
    timeline: List[Dict[str, Any]] = Field(default_factory=list)
    recommendations: List[str] = Field(default_factory=list)
    related_jurisprudence: List[Dict[str, Any]] = Field(default_factory=list)
    metadata: Dict[str, Any] = Field(default_factory=dict)
    processing_time_ms: int


class DocumentAnalysisRequest(BaseModel):
    """Document analysis request model."""
    
    document_content: str
    document_type: Optional[str] = None
    analysis_depth: str = Field(default="standard")  # quick, standard, deep
    extract_entities: bool = True
    extract_citations: bool = True
    generate_summary: bool = True
    tenant_id: Optional[UUID] = None


class DocumentAnalysisResponse(BaseModel):
    """Document analysis response model."""
    
    document_type: str
    summary: str
    key_points: List[str]
    extracted_entities: List[ExtractedEntity]
    legal_citations: List[Dict[str, str]]
    recommendations: List[str]
    metadata: Dict[str, Any] = Field(default_factory=dict)
    processing_time_ms: int