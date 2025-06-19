"""SQLAlchemy database models."""

from datetime import datetime
from typing import Optional
from uuid import uuid4

try:
    from pgvector.sqlalchemy import Vector
    PGVECTOR_AVAILABLE = True
except ImportError:
    PGVECTOR_AVAILABLE = False
    # Use JSON as fallback when pgvector is not available
    from sqlalchemy import JSON
    Vector = lambda dim: JSON
from sqlalchemy import (
    Boolean, Column, DateTime, Float, ForeignKey, Integer, 
    String, Text, JSON, Index, UniqueConstraint
)
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship

from app.db.database import Base


class LegalDecisionDB(Base):
    """Legal decision database model."""
    
    __tablename__ = "legal_decisions"
    __table_args__ = (
        Index("idx_legal_decisions_court_type", "court_type"),
        Index("idx_legal_decisions_decision_date", "decision_date"),
        Index("idx_legal_decisions_legal_subject", "legal_subject"),
        Index("idx_legal_decisions_embedding", "embedding", postgresql_using="ivfflat"),
    )
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid4)
    court_name = Column(String(200), nullable=False)
    court_type = Column(String(50), nullable=False)
    decision_date = Column(DateTime, nullable=False)
    publication_date = Column(DateTime)
    process_number = Column(String(50), nullable=False, unique=True)
    case_type = Column(String(100))
    legal_subject = Column(String(200))
    decision_type = Column(String(50), nullable=False)
    rapporteur = Column(String(200))
    decision_text = Column(Text, nullable=False)
    summary = Column(Text)
    legal_references = Column(JSON, default=list)
    cited_cases = Column(JSON, default=list)
    keywords = Column(JSON, default=list)
    embedding = Column(Vector(384))  # Dimension based on model
    document_metadata = Column(JSON, default=dict)
    source_url = Column(String(500))
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class DocumentTemplateDB(Base):
    """Document template database model."""
    
    __tablename__ = "document_templates"
    __table_args__ = (
        Index("idx_document_templates_type", "document_type"),
        Index("idx_document_templates_active", "is_active"),
    )
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid4)
    name = Column(String(200), nullable=False)
    document_type = Column(String(50), nullable=False)
    description = Column(Text)
    template_content = Column(Text, nullable=False)
    variables = Column(JSON, default=list)
    example_usage = Column(Text)
    tags = Column(JSON, default=list)
    is_active = Column(Boolean, default=True)
    created_by = Column(UUID(as_uuid=True))
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class AnalysisHistoryDB(Base):
    """Analysis history database model."""
    
    __tablename__ = "analysis_history"
    __table_args__ = (
        Index("idx_analysis_history_tenant", "tenant_id"),
        Index("idx_analysis_history_process", "process_id"),
        Index("idx_analysis_history_type_date", "analysis_type", "created_at"),
    )
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id = Column(UUID(as_uuid=True), nullable=False)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    process_id = Column(UUID(as_uuid=True))
    analysis_type = Column(String(50), nullable=False)
    request_data = Column(JSON, nullable=False)
    response_data = Column(JSON, nullable=False)
    processing_time_ms = Column(Integer)
    tokens_used = Column(Integer)
    cost_estimate = Column(Float)
    created_at = Column(DateTime, default=datetime.utcnow)


class GenerationHistoryDB(Base):
    """Document generation history database model."""
    
    __tablename__ = "generation_history"
    __table_args__ = (
        Index("idx_generation_history_tenant", "tenant_id"),
        Index("idx_generation_history_type_date", "document_type", "created_at"),
    )
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id = Column(UUID(as_uuid=True), nullable=False)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    document_type = Column(String(50), nullable=False)
    template_id = Column(UUID(as_uuid=True), ForeignKey("document_templates.id"))
    request_data = Column(JSON, nullable=False)
    generated_content = Column(Text, nullable=False)
    output_format = Column(String(20))
    word_count = Column(Integer)
    tokens_used = Column(Integer)
    cost_estimate = Column(Float)
    created_at = Column(DateTime, default=datetime.utcnow)
    
    # Relationships
    template = relationship("DocumentTemplateDB", backref="generations")


class JurisprudenceCollectorStatusDB(Base):
    """Jurisprudence collector status tracking."""
    
    __tablename__ = "jurisprudence_collector_status"
    __table_args__ = (
        UniqueConstraint("source_name", "court_type", name="uq_collector_source_court"),
    )
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid4)
    source_name = Column(String(100), nullable=False)
    court_type = Column(String(50), nullable=False)
    last_collection_date = Column(DateTime)
    last_process_number = Column(String(50))
    total_collected = Column(Integer, default=0)
    status = Column(String(50), default="idle")
    error_message = Column(Text)
    collector_metadata = Column(JSON, default=dict)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)