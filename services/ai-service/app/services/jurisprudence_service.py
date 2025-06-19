"""Jurisprudence service for legal decisions analysis."""

import asyncio
from datetime import datetime
from typing import List, Optional, Dict, Any, Tuple
from uuid import UUID

from sqlalchemy import select, and_, or_, func
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.config import settings
from app.core.exceptions import SearchError, AIServiceError
from app.core.logging import get_logger
from app.db.models import LegalDecisionDB
from app.models.jurisprudence import (
    LegalDecision,
    CourtType,
    DecisionType,
    SimilarityDimension,
    CaseSimilarityResult
)
from app.services.embeddings import embedding_service
from app.services.vector_store import vector_store
from app.services.text_processing import LegalTextProcessor
from app.services.cache import cache_service

logger = get_logger(__name__)


class JurisprudenceService:
    """Service for jurisprudence operations."""
    
    def __init__(self):
        """Initialize jurisprudence service."""
        self.embedding_service = embedding_service
        self.vector_store = vector_store
        self.text_processor = LegalTextProcessor()
        self.cache = cache_service
    
    async def search_similar(
        self,
        query: str,
        court_types: Optional[List[CourtType]] = None,
        decision_types: Optional[List[DecisionType]] = None,
        date_from: Optional[datetime] = None,
        date_to: Optional[datetime] = None,
        legal_subjects: Optional[List[str]] = None,
        max_results: int = 10,
        similarity_threshold: float = 0.7,
        session: AsyncSession = None
    ) -> List[LegalDecision]:
        """Search for similar legal decisions."""
        try:
            # Generate embedding for query
            query_embedding = await self.embedding_service.generate_embedding(query)
            
            # Search vector store
            similar_docs = await self.vector_store.search_similar(
                query_embedding=query_embedding,
                k=max_results * 2,  # Get more to filter
                threshold=similarity_threshold
            )
            
            if not similar_docs:
                return []
            
            # Get document IDs
            doc_ids = [doc_id for doc_id, _ in similar_docs]
            
            # Build database query with filters
            db_query = select(LegalDecisionDB).where(
                LegalDecisionDB.id.in_(doc_ids)
            )
            
            # Apply filters
            if court_types:
                db_query = db_query.where(
                    LegalDecisionDB.court_type.in_([ct.value for ct in court_types])
                )
            
            if decision_types:
                db_query = db_query.where(
                    LegalDecisionDB.decision_type.in_([dt.value for dt in decision_types])
                )
            
            if date_from:
                db_query = db_query.where(LegalDecisionDB.decision_date >= date_from)
            
            if date_to:
                db_query = db_query.where(LegalDecisionDB.decision_date <= date_to)
            
            if legal_subjects:
                # Search in legal subjects array
                subject_conditions = [
                    LegalDecisionDB.legal_subjects.contains([subject])
                    for subject in legal_subjects
                ]
                db_query = db_query.where(or_(*subject_conditions))
            
            # Execute query
            result = await session.execute(db_query)
            db_decisions = result.scalars().all()
            
            # Create similarity score map
            similarity_map = {doc_id: score for doc_id, score in similar_docs}
            
            # Convert to domain models and sort by similarity
            decisions = []
            for db_decision in db_decisions:
                similarity_score = similarity_map.get(db_decision.id, 0.0)
                
                decision = LegalDecision(
                    id=db_decision.id,
                    court_name=db_decision.court_name,
                    court_type=CourtType(db_decision.court_type),
                    decision_type=DecisionType(db_decision.decision_type),
                    case_number=db_decision.case_number,
                    decision_date=db_decision.decision_date,
                    summary=db_decision.summary,
                    decision_text=db_decision.decision_text,
                    legal_subjects=db_decision.legal_subjects or [],
                    keywords=db_decision.keywords or [],
                    similarity_score=similarity_score,
                    citation_count=db_decision.citation_count or 0,
                    relevance_score=db_decision.relevance_score or 0.0
                )
                decisions.append(decision)
            
            # Sort by similarity score and limit results
            decisions.sort(key=lambda x: x.similarity_score, reverse=True)
            return decisions[:max_results]
            
        except Exception as e:
            logger.error("Jurisprudence search failed", error=str(e))
            raise SearchError(f"Failed to search jurisprudence: {str(e)}")
    
    async def analyze_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Optional[Dict[str, Any]] = None,
        case_list: Optional[List[Dict[str, Any]]] = None,
        similarity_dimensions: List[SimilarityDimension] = None,
        include_explanation: bool = False,
        session: AsyncSession = None
    ) -> List[CaseSimilarityResult]:
        """Analyze similarity between cases."""
        try:
            if not similarity_dimensions:
                similarity_dimensions = [
                    SimilarityDimension.SEMANTIC,
                    SimilarityDimension.LEGAL,
                    SimilarityDimension.FACTUAL
                ]
            
            results = []
            
            # Compare with single case
            if case_b:
                similarity = await self._calculate_case_similarity(
                    case_a, case_b, similarity_dimensions, include_explanation
                )
                results.append(similarity)
            
            # Compare with multiple cases
            if case_list:
                tasks = []
                for case in case_list:
                    task = self._calculate_case_similarity(
                        case_a, case, similarity_dimensions, include_explanation
                    )
                    tasks.append(task)
                
                similarities = await asyncio.gather(*tasks)
                results.extend(similarities)
            
            return results
            
        except Exception as e:
            logger.error("Case similarity analysis failed", error=str(e))
            raise AIServiceError(f"Failed to analyze case similarity: {str(e)}")
    
    async def _calculate_case_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Dict[str, Any],
        dimensions: List[SimilarityDimension],
        include_explanation: bool
    ) -> SimilarityResult:
        """Calculate similarity between two cases."""
        dimension_scores = {}
        explanations = {}
        
        for dimension in dimensions:
            if dimension == SimilarityDimension.SEMANTIC:
                score, explanation = await self._semantic_similarity(case_a, case_b)
            elif dimension == SimilarityDimension.LEGAL:
                score, explanation = await self._legal_similarity(case_a, case_b)
            elif dimension == SimilarityDimension.FACTUAL:
                score, explanation = await self._factual_similarity(case_a, case_b)
            elif dimension == SimilarityDimension.PROCEDURAL:
                score, explanation = await self._procedural_similarity(case_a, case_b)
            elif dimension == SimilarityDimension.CONTEXTUAL:
                score, explanation = await self._contextual_similarity(case_a, case_b)
            else:
                score, explanation = 0.0, "Unknown dimension"
            
            dimension_scores[dimension.value] = score
            if include_explanation:
                explanations[dimension.value] = explanation
        
        # Calculate overall similarity (weighted average)
        weights = {
            SimilarityDimension.SEMANTIC.value: 0.3,
            SimilarityDimension.LEGAL.value: 0.25,
            SimilarityDimension.FACTUAL.value: 0.25,
            SimilarityDimension.PROCEDURAL.value: 0.1,
            SimilarityDimension.CONTEXTUAL.value: 0.1
        }
        
        overall_score = sum(
            dimension_scores.get(dim, 0.0) * weights.get(dim, 0.0)
            for dim in dimension_scores.keys()
        )
        
        return SimilarityResult(
            case_id=str(case_b.get("id", "unknown")),
            overall_similarity=overall_score,
            dimension_scores=dimension_scores,
            explanations=explanations if include_explanation else None
        )
    
    async def _semantic_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Dict[str, Any]
    ) -> Tuple[float, str]:
        """Calculate semantic similarity."""
        try:
            text_a = case_a.get("decision_text", "") + " " + case_a.get("summary", "")
            text_b = case_b.get("decision_text", "") + " " + case_b.get("summary", "")
            
            if not text_a.strip() or not text_b.strip():
                return 0.0, "Insufficient text for semantic comparison"
            
            # Generate embeddings
            embedding_a = await self.embedding_service.generate_embedding(text_a)
            embedding_b = await self.embedding_service.generate_embedding(text_b)
            
            # Calculate cosine similarity
            import numpy as np
            a = np.array(embedding_a)
            b = np.array(embedding_b)
            
            similarity = np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))
            
            explanation = f"Semantic similarity based on text embeddings: {similarity:.3f}"
            
            return float(similarity), explanation
            
        except Exception as e:
            logger.error("Semantic similarity calculation failed", error=str(e))
            return 0.0, f"Error calculating semantic similarity: {str(e)}"
    
    async def _legal_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Dict[str, Any]
    ) -> Tuple[float, str]:
        """Calculate legal similarity."""
        try:
            subjects_a = set(case_a.get("legal_subjects", []))
            subjects_b = set(case_b.get("legal_subjects", []))
            
            keywords_a = set(case_a.get("keywords", []))
            keywords_b = set(case_b.get("keywords", []))
            
            if not subjects_a and not subjects_b and not keywords_a and not keywords_b:
                return 0.0, "No legal subjects or keywords for comparison"
            
            # Calculate Jaccard similarity for subjects and keywords
            subject_intersection = len(subjects_a & subjects_b)
            subject_union = len(subjects_a | subjects_b)
            subject_similarity = subject_intersection / subject_union if subject_union > 0 else 0.0
            
            keyword_intersection = len(keywords_a & keywords_b)
            keyword_union = len(keywords_a | keywords_b)
            keyword_similarity = keyword_intersection / keyword_union if keyword_union > 0 else 0.0
            
            # Weighted average
            legal_similarity = (subject_similarity * 0.7) + (keyword_similarity * 0.3)
            
            explanation = (
                f"Legal similarity based on subjects ({subject_similarity:.3f}) "
                f"and keywords ({keyword_similarity:.3f}): {legal_similarity:.3f}"
            )
            
            return legal_similarity, explanation
            
        except Exception as e:
            logger.error("Legal similarity calculation failed", error=str(e))
            return 0.0, f"Error calculating legal similarity: {str(e)}"
    
    async def _factual_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Dict[str, Any]
    ) -> Tuple[float, str]:
        """Calculate factual similarity."""
        try:
            # Extract facts from decision text using text processing
            facts_a = await self.text_processor.extract_facts(
                case_a.get("decision_text", "")
            )
            facts_b = await self.text_processor.extract_facts(
                case_b.get("decision_text", "")
            )
            
            if not facts_a or not facts_b:
                return 0.0, "Insufficient factual information for comparison"
            
            # Calculate similarity between extracted facts
            facts_text_a = " ".join(facts_a)
            facts_text_b = " ".join(facts_b)
            
            embedding_a = await self.embedding_service.generate_embedding(facts_text_a)
            embedding_b = await self.embedding_service.generate_embedding(facts_text_b)
            
            import numpy as np
            a = np.array(embedding_a)
            b = np.array(embedding_b)
            
            similarity = np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))
            
            explanation = f"Factual similarity based on extracted facts: {similarity:.3f}"
            
            return float(similarity), explanation
            
        except Exception as e:
            logger.error("Factual similarity calculation failed", error=str(e))
            return 0.0, f"Error calculating factual similarity: {str(e)}"
    
    async def _procedural_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Dict[str, Any]
    ) -> Tuple[float, str]:
        """Calculate procedural similarity."""
        try:
            court_type_a = case_a.get("court_type")
            court_type_b = case_b.get("court_type")
            
            decision_type_a = case_a.get("decision_type")
            decision_type_b = case_b.get("decision_type")
            
            court_similarity = 1.0 if court_type_a == court_type_b else 0.0
            decision_similarity = 1.0 if decision_type_a == decision_type_b else 0.0
            
            procedural_similarity = (court_similarity + decision_similarity) / 2
            
            explanation = (
                f"Procedural similarity based on court type match ({court_similarity}) "
                f"and decision type match ({decision_similarity}): {procedural_similarity:.3f}"
            )
            
            return procedural_similarity, explanation
            
        except Exception as e:
            logger.error("Procedural similarity calculation failed", error=str(e))
            return 0.0, f"Error calculating procedural similarity: {str(e)}"
    
    async def _contextual_similarity(
        self,
        case_a: Dict[str, Any],
        case_b: Dict[str, Any]
    ) -> Tuple[float, str]:
        """Calculate contextual similarity."""
        try:
            # Compare temporal context
            date_a = case_a.get("decision_date")
            date_b = case_b.get("decision_date")
            
            if date_a and date_b:
                if isinstance(date_a, str):
                    date_a = datetime.fromisoformat(date_a.replace('Z', '+00:00'))
                if isinstance(date_b, str):
                    date_b = datetime.fromisoformat(date_b.replace('Z', '+00:00'))
                
                # Calculate temporal similarity (closer dates = higher similarity)
                time_diff = abs((date_a - date_b).days)
                temporal_similarity = max(0.0, 1.0 - (time_diff / 3650))  # 10 years max
            else:
                temporal_similarity = 0.5  # Default if dates unavailable
            
            explanation = f"Contextual similarity based on temporal proximity: {temporal_similarity:.3f}"
            
            return temporal_similarity, explanation
            
        except Exception as e:
            logger.error("Contextual similarity calculation failed", error=str(e))
            return 0.0, f"Error calculating contextual similarity: {str(e)}"
    
    async def find_precedents(
        self,
        process_data: Dict[str, Any],
        max_results: int = 10,
        include_similar_facts: bool = True,
        court_hierarchy: Optional[List[CourtType]] = None,
        session: AsyncSession = None
    ) -> List[Dict[str, Any]]:
        """Find legal precedents for a case."""
        try:
            # Extract key information from process data
            legal_subjects = process_data.get("legal_subjects", [])
            case_summary = process_data.get("summary", "")
            case_facts = process_data.get("facts", "")
            
            # Build search query
            search_text = f"{case_summary} {case_facts}".strip()
            if not search_text:
                search_text = " ".join(legal_subjects)
            
            # Search for similar decisions
            similar_decisions = await self.search_similar(
                query=search_text,
                court_types=court_hierarchy,
                legal_subjects=legal_subjects,
                max_results=max_results * 2,  # Get more to filter
                similarity_threshold=0.6,
                session=session
            )
            
            precedents = []
            for decision in similar_decisions:
                precedent = {
                    "id": str(decision.id),
                    "court_name": decision.court_name,
                    "court_type": decision.court_type.value,
                    "case_number": decision.case_number,
                    "decision_date": decision.decision_date.isoformat() if decision.decision_date else None,
                    "summary": decision.summary,
                    "legal_subjects": decision.legal_subjects,
                    "similarity_score": decision.similarity_score,
                    "citation_count": decision.citation_count,
                    "relevance_score": decision.relevance_score,
                    "precedent_strength": self._calculate_precedent_strength(decision)
                }
                
                if include_similar_facts:
                    precedent["similar_facts"] = await self._extract_similar_facts(
                        case_facts, decision.decision_text
                    )
                
                precedents.append(precedent)
            
            # Sort by precedent strength and similarity
            precedents.sort(
                key=lambda x: (x["precedent_strength"], x["similarity_score"]),
                reverse=True
            )
            
            return precedents[:max_results]
            
        except Exception as e:
            logger.error("Precedent search failed", error=str(e))
            raise SearchError(f"Failed to find precedents: {str(e)}")
    
    def _calculate_precedent_strength(self, decision: LegalDecision) -> float:
        """Calculate the strength of a precedent."""
        strength = 0.0
        
        # Court hierarchy weight
        court_weights = {
            CourtType.SUPREME_COURT: 1.0,
            CourtType.SUPERIOR_COURT: 0.8,
            CourtType.FEDERAL_COURT: 0.6,
            CourtType.STATE_COURT: 0.4,
            CourtType.LABOR_COURT: 0.5,
            CourtType.ELECTORAL_COURT: 0.3
        }
        strength += court_weights.get(decision.court_type, 0.2)
        
        # Citation count influence
        citation_score = min(1.0, (decision.citation_count or 0) / 100)
        strength += citation_score * 0.3
        
        # Relevance score
        strength += (decision.relevance_score or 0.0) * 0.2
        
        return min(1.0, strength)
    
    async def _extract_similar_facts(
        self,
        case_facts: str,
        decision_text: str
    ) -> List[str]:
        """Extract similar facts between case and decision."""
        try:
            # Extract facts from both texts
            case_fact_list = await self.text_processor.extract_facts(case_facts)
            decision_fact_list = await self.text_processor.extract_facts(decision_text)
            
            similar_facts = []
            
            # Find similar facts using embeddings
            for case_fact in case_fact_list[:5]:  # Limit to avoid too many API calls
                case_embedding = await self.embedding_service.generate_embedding(case_fact)
                
                best_similarity = 0.0
                best_fact = ""
                
                for decision_fact in decision_fact_list:
                    decision_embedding = await self.embedding_service.generate_embedding(decision_fact)
                    
                    import numpy as np
                    a = np.array(case_embedding)
                    b = np.array(decision_embedding)
                    
                    similarity = np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))
                    
                    if similarity > best_similarity and similarity > 0.7:
                        best_similarity = similarity
                        best_fact = decision_fact
                
                if best_fact:
                    similar_facts.append(best_fact)
            
            return similar_facts[:3]  # Return top 3 similar facts
            
        except Exception as e:
            logger.error("Similar facts extraction failed", error=str(e))
            return []
    
    async def get_statistics(self, session: AsyncSession) -> Dict[str, Any]:
        """Get jurisprudence database statistics."""
        try:
            # Total decisions
            total_query = select(func.count(LegalDecisionDB.id))
            total_result = await session.execute(total_query)
            total_decisions = total_result.scalar()
            
            # Decisions by court type
            court_query = select(
                LegalDecisionDB.court_type,
                func.count(LegalDecisionDB.id)
            ).group_by(LegalDecisionDB.court_type)
            court_result = await session.execute(court_query)
            court_stats = {row[0]: row[1] for row in court_result}
            
            # Decisions by decision type
            decision_query = select(
                LegalDecisionDB.decision_type,
                func.count(LegalDecisionDB.id)
            ).group_by(LegalDecisionDB.decision_type)
            decision_result = await session.execute(decision_query)
            decision_stats = {row[0]: row[1] for row in decision_result}
            
            # Recent activity (last 30 days)
            from datetime import timedelta
            recent_date = datetime.utcnow() - timedelta(days=30)
            recent_query = select(func.count(LegalDecisionDB.id)).where(
                LegalDecisionDB.created_at >= recent_date
            )
            recent_result = await session.execute(recent_query)
            recent_decisions = recent_result.scalar()
            
            return {
                "total_decisions": total_decisions,
                "by_court_type": court_stats,
                "by_decision_type": decision_stats,
                "recent_decisions_30d": recent_decisions,
                "last_updated": datetime.utcnow().isoformat()
            }
            
        except Exception as e:
            logger.error("Failed to get statistics", error=str(e))
            raise SearchError(f"Failed to get statistics: {str(e)}")