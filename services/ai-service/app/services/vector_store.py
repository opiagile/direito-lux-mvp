"""Vector store service for similarity search."""

import os
import pickle
from typing import List, Optional, Tuple, Dict, Any
from uuid import UUID

import numpy as np
try:
    import faiss
    FAISS_AVAILABLE = True
except ImportError:
    FAISS_AVAILABLE = False
    faiss = None
from sqlalchemy import select, text
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.config import settings
from app.core.exceptions import SearchError
from app.core.logging import get_logger
from app.db.database import get_session
from app.db.models import LegalDecisionDB
from app.models.jurisprudence import LegalDecision

logger = get_logger(__name__)


class VectorStore:
    """Vector store for similarity search."""
    
    def __init__(self):
        """Initialize vector store."""
        self.dimension = settings.embedding_dimension
        self.index_path = settings.faiss_index_path
        self.index: Optional[faiss.Index] = None
        self.id_map: Dict[int, UUID] = {}  # FAISS index -> document ID
        self.initialized = False
    
    async def initialize(self) -> None:
        """Initialize vector store."""
        try:
            if settings.vector_store_type == "faiss":
                await self._init_faiss()
            elif settings.vector_store_type == "pgvector":
                await self._init_pgvector()
            else:
                raise ValueError(f"Unknown vector store type: {settings.vector_store_type}")
            
            self.initialized = True
            logger.info(f"Vector store initialized", type=settings.vector_store_type)
        except Exception as e:
            logger.error(f"Failed to initialize vector store", error=str(e))
            raise SearchError(f"Vector store initialization failed: {str(e)}")
    
    async def _init_faiss(self) -> None:
        """Initialize FAISS index."""
        if not FAISS_AVAILABLE:
            logger.warning("FAISS not available, skipping FAISS initialization")
            return
            
        # Try to load existing index
        if os.path.exists(self.index_path):
            try:
                self.index = faiss.read_index(self.index_path)
                
                # Load ID map
                id_map_path = self.index_path + ".ids"
                if os.path.exists(id_map_path):
                    with open(id_map_path, "rb") as f:
                        self.id_map = pickle.load(f)
                
                logger.info(f"Loaded FAISS index", 
                          num_vectors=self.index.ntotal,
                          dimension=self.dimension)
            except Exception as e:
                logger.warning(f"Failed to load existing index", error=str(e))
                self._create_new_index()
        else:
            self._create_new_index()
    
    def _create_new_index(self) -> None:
        """Create new FAISS index."""
        if not FAISS_AVAILABLE:
            logger.warning("FAISS not available, cannot create index")
            return
            
        # Create index with cosine similarity
        self.index = faiss.IndexFlatIP(self.dimension)
        # Normalize vectors for cosine similarity
        self.index = faiss.IndexIDMap(self.index)
        logger.info(f"Created new FAISS index", dimension=self.dimension)
    
    async def _init_pgvector(self) -> None:
        """Initialize pgvector (PostgreSQL extension)."""
        # pgvector is initialized via SQLAlchemy models
        logger.info("Using pgvector for similarity search")
    
    async def add_embeddings(
        self,
        embeddings: List[List[float]],
        document_ids: List[UUID]
    ) -> None:
        """Add embeddings to vector store."""
        if not self.initialized:
            await self.initialize()
        
        if settings.vector_store_type == "faiss":
            await self._add_to_faiss(embeddings, document_ids)
        else:
            # pgvector stores embeddings directly in PostgreSQL
            pass
    
    async def _add_to_faiss(
        self,
        embeddings: List[List[float]],
        document_ids: List[UUID]
    ) -> None:
        """Add embeddings to FAISS index."""
        try:
            # Convert to numpy array
            vectors = np.array(embeddings, dtype=np.float32)
            
            # Normalize for cosine similarity
            faiss.normalize_L2(vectors)
            
            # Generate IDs for new vectors
            start_id = len(self.id_map)
            ids = np.arange(start_id, start_id + len(vectors), dtype=np.int64)
            
            # Add to index
            self.index.add_with_ids(vectors, ids)
            
            # Update ID map
            for i, doc_id in enumerate(document_ids):
                self.id_map[start_id + i] = doc_id
            
            # Save index and ID map
            await self._save_faiss_index()
            
            logger.info(f"Added {len(embeddings)} vectors to FAISS index")
        except Exception as e:
            logger.error(f"Failed to add embeddings to FAISS", error=str(e))
            raise SearchError(f"Failed to add embeddings: {str(e)}")
    
    async def _save_faiss_index(self) -> None:
        """Save FAISS index to disk."""
        try:
            # Create directory if needed
            os.makedirs(os.path.dirname(self.index_path), exist_ok=True)
            
            # Save index
            faiss.write_index(self.index, self.index_path)
            
            # Save ID map
            id_map_path = self.index_path + ".ids"
            with open(id_map_path, "wb") as f:
                pickle.dump(self.id_map, f)
            
            logger.debug("Saved FAISS index to disk")
        except Exception as e:
            logger.error(f"Failed to save FAISS index", error=str(e))
    
    async def search_similar(
        self,
        query_embedding: List[float],
        k: int = 10,
        threshold: float = 0.0
    ) -> List[Tuple[UUID, float]]:
        """Search for similar documents."""
        if not self.initialized:
            await self.initialize()
        
        if settings.vector_store_type == "faiss":
            return await self._search_faiss(query_embedding, k, threshold)
        else:
            return await self._search_pgvector(query_embedding, k, threshold)
    
    async def _search_faiss(
        self,
        query_embedding: List[float],
        k: int,
        threshold: float
    ) -> List[Tuple[UUID, float]]:
        """Search using FAISS index."""
        try:
            if self.index.ntotal == 0:
                return []
            
            # Convert to numpy array
            query_vector = np.array([query_embedding], dtype=np.float32)
            
            # Normalize for cosine similarity
            faiss.normalize_L2(query_vector)
            
            # Search
            k = min(k, self.index.ntotal)
            distances, indices = self.index.search(query_vector, k)
            
            # Filter by threshold and convert to results
            results = []
            for i, (dist, idx) in enumerate(zip(distances[0], indices[0])):
                if idx >= 0 and dist >= threshold:  # Valid index and above threshold
                    doc_id = self.id_map.get(int(idx))
                    if doc_id:
                        results.append((doc_id, float(dist)))
            
            return results
        except Exception as e:
            logger.error(f"FAISS search failed", error=str(e))
            raise SearchError(f"Search failed: {str(e)}")
    
    async def _search_pgvector(
        self,
        query_embedding: List[float],
        k: int,
        threshold: float
    ) -> List[Tuple[UUID, float]]:
        """Search using pgvector."""
        try:
            async with get_session() as session:
                # Build query
                query = text("""
                    SELECT id, 1 - (embedding <=> :query_embedding) as similarity
                    FROM legal_decisions
                    WHERE embedding IS NOT NULL
                    AND 1 - (embedding <=> :query_embedding) >= :threshold
                    ORDER BY embedding <=> :query_embedding
                    LIMIT :k
                """)
                
                result = await session.execute(
                    query,
                    {
                        "query_embedding": query_embedding,
                        "threshold": threshold,
                        "k": k
                    }
                )
                
                return [(row.id, row.similarity) for row in result]
        except Exception as e:
            logger.error(f"pgvector search failed", error=str(e))
            raise SearchError(f"Search failed: {str(e)}")
    
    async def update_embedding(
        self,
        document_id: UUID,
        embedding: List[float]
    ) -> None:
        """Update embedding for a document."""
        if settings.vector_store_type == "faiss":
            # For FAISS, we need to rebuild or use a different strategy
            logger.warning("FAISS index update not implemented - requires rebuild")
        else:
            # pgvector updates are handled via SQLAlchemy
            pass
    
    async def remove_document(self, document_id: UUID) -> None:
        """Remove document from vector store."""
        if settings.vector_store_type == "faiss":
            # For FAISS, we need to rebuild or use a different strategy
            logger.warning("FAISS document removal not implemented - requires rebuild")
        else:
            # pgvector removals are handled via SQLAlchemy
            pass
    
    async def get_stats(self) -> Dict[str, Any]:
        """Get vector store statistics."""
        stats = {
            "type": settings.vector_store_type,
            "initialized": self.initialized,
            "dimension": self.dimension
        }
        
        if settings.vector_store_type == "faiss" and self.index:
            stats["num_vectors"] = self.index.ntotal
            stats["index_type"] = type(self.index).__name__
        elif settings.vector_store_type == "pgvector":
            async with get_session() as session:
                result = await session.execute(
                    select(LegalDecisionDB.id).where(
                        LegalDecisionDB.embedding.isnot(None)
                    )
                )
                stats["num_vectors"] = len(result.all())
        
        return stats


# Global vector store instance
vector_store = VectorStore()


async def init_vector_store() -> None:
    """Initialize vector store."""
    await vector_store.initialize()


async def close_vector_store() -> None:
    """Close vector store."""
    # No specific cleanup needed for current implementations
    pass