"""Document generation domain models."""

from datetime import datetime
from enum import Enum
from typing import List, Optional, Dict, Any
from uuid import UUID

from pydantic import BaseModel, Field, validator


class DocumentType(str, Enum):
    """Document type enumeration."""
    PETICAO_INICIAL = "peticao_inicial"
    CONTESTACAO = "contestacao"
    RECURSO = "recurso"
    CONTRARRAZOES = "contrarrazoes"
    PARECER = "parecer"
    CONTRATO = "contrato"
    PROCURACAO = "procuracao"
    NOTIFICACAO = "notificacao"
    MEMORANDO = "memorando"
    RELATORIO = "relatorio"


class DocumentStyle(str, Enum):
    """Document writing style."""
    FORMAL = "formal"
    TECHNICAL = "technical"
    SIMPLE = "simple"
    PERSUASIVE = "persuasive"


class DocumentTemplate(BaseModel):
    """Document template model."""
    
    id: UUID
    name: str
    document_type: DocumentType
    description: Optional[str] = None
    template_content: str
    variables: List[Dict[str, str]]  # name, type, description, required
    example_usage: Optional[str] = None
    tags: List[str] = Field(default_factory=list)
    is_active: bool = True
    created_at: datetime
    updated_at: datetime


class DocumentGenerationRequest(BaseModel):
    """Document generation request model."""
    
    document_type: DocumentType
    template_id: Optional[UUID] = None
    case_data: Dict[str, Any]  # Process data, parties, etc.
    variables: Dict[str, Any]  # Template variables
    style: DocumentStyle = Field(default=DocumentStyle.FORMAL)
    include_jurisprudence: bool = False
    jurisprudence_count: int = Field(default=3, ge=0, le=10)
    include_legal_basis: bool = True
    language: str = Field(default="pt-BR")
    output_format: str = Field(default="html")  # html, markdown, docx, pdf
    tenant_id: Optional[UUID] = None
    user_plan: str = Field(default="starter")
    
    @validator("output_format")
    def validate_output_format(cls, v):
        """Validate output format."""
        allowed = ["html", "markdown", "docx", "pdf"]
        if v not in allowed:
            raise ValueError(f"Output format must be one of {allowed}")
        return v


class GeneratedSection(BaseModel):
    """Generated document section."""
    
    title: str
    content: str
    order: int
    is_editable: bool = True
    metadata: Dict[str, Any] = Field(default_factory=dict)


class LegalReference(BaseModel):
    """Legal reference used in document."""
    
    law: str
    article: Optional[str] = None
    paragraph: Optional[str] = None
    full_text: Optional[str] = None
    relevance: float = Field(ge=0.0, le=1.0)


class DocumentGenerationResponse(BaseModel):
    """Document generation response model."""
    
    document_id: UUID
    document_type: DocumentType
    title: str
    sections: List[GeneratedSection]
    legal_references: List[LegalReference] = Field(default_factory=list)
    jurisprudence_citations: List[Dict[str, Any]] = Field(default_factory=list)
    generated_content: str  # Full document in requested format
    word_count: int
    estimated_reading_time: int  # in minutes
    metadata: Dict[str, Any] = Field(default_factory=dict)
    generation_timestamp: datetime = Field(default_factory=datetime.utcnow)
    processing_time_ms: int


class TemplateListRequest(BaseModel):
    """Template list request model."""
    
    document_type: Optional[DocumentType] = None
    tags: Optional[List[str]] = None
    search_query: Optional[str] = None
    is_active: Optional[bool] = True
    page: int = Field(default=1, ge=1)
    page_size: int = Field(default=20, ge=1, le=100)


class TemplateListResponse(BaseModel):
    """Template list response model."""
    
    templates: List[DocumentTemplate]
    total_count: int
    page: int
    page_size: int
    total_pages: int


class DocumentRevisionRequest(BaseModel):
    """Document revision request model."""
    
    document_id: UUID
    current_content: str
    revision_instructions: str
    maintain_style: bool = True
    maintain_structure: bool = False
    tenant_id: Optional[UUID] = None


class DocumentRevisionResponse(BaseModel):
    """Document revision response model."""
    
    document_id: UUID
    revised_content: str
    changes_summary: List[str]
    revision_timestamp: datetime = Field(default_factory=datetime.utcnow)
    processing_time_ms: int