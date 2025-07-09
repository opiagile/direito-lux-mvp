"""Embedding generation service."""

from typing import List, Optional, Dict, Any
import numpy as np
from tenacity import retry, stop_after_attempt, wait_exponential

try:
    from sentence_transformers import SentenceTransformer
    SENTENCE_TRANSFORMERS_AVAILABLE = True
except ImportError:
    SENTENCE_TRANSFORMERS_AVAILABLE = False
    SentenceTransformer = None

try:
    import openai
    from openai import OpenAI
    OPENAI_AVAILABLE = True
except ImportError:
    OPENAI_AVAILABLE = False
    openai = None
    OpenAI = None

try:
    import httpx
    HTTPX_AVAILABLE = True
except ImportError:
    HTTPX_AVAILABLE = False
    httpx = None

from app.core.config import settings
from app.core.exceptions import EmbeddingError
from app.core.logging import get_logger
from app.services.text_processing import LegalTextProcessor

logger = get_logger(__name__)


class EmbeddingService:
    """Service for generating embeddings from text."""
    
    def __init__(self):
        """Initialize embedding service."""
        self.text_processor = LegalTextProcessor()
        self._models = {}
        self._openai_client = None
        
    def _get_model(self, model_name: str) -> Any:
        """Get or load a model."""
        if not SENTENCE_TRANSFORMERS_AVAILABLE:
            raise EmbeddingError("sentence-transformers not available")
            
        if model_name not in self._models:
            try:
                if model_name.startswith("sentence-transformers/"):
                    logger.info(f"Loading model: {model_name}")
                    self._models[model_name] = SentenceTransformer(model_name)
                else:
                    raise ValueError(f"Unsupported model: {model_name}")
            except Exception as e:
                logger.error(f"Failed to load model {model_name}", error=str(e))
                raise EmbeddingError(f"Failed to load model: {model_name}")
        
        return self._models[model_name]
    
    def _get_openai_client(self) -> OpenAI:
        """Get OpenAI client."""
        if not OPENAI_AVAILABLE:
            raise EmbeddingError("OpenAI not available")
            
        if not self._openai_client:
            self._openai_client = OpenAI(api_key=settings.openai_api_key)
        return self._openai_client
    
    async def _get_ollama_embedding(self, text: str, model: str = None) -> List[float]:
        """Get embedding from Ollama."""
        if not HTTPX_AVAILABLE:
            raise EmbeddingError("httpx not available for Ollama")
        
        if not model:
            model = settings.ollama_embeddings_model
            
        url = f"{settings.ollama_base_url}/api/embeddings"
        
        async with httpx.AsyncClient() as client:
            try:
                response = await client.post(
                    url,
                    json={
                        "model": model,
                        "prompt": text
                    },
                    timeout=30.0
                )
                response.raise_for_status()
                data = response.json()
                return data.get("embedding", [])
            except Exception as e:
                logger.error(f"Ollama embedding failed: {str(e)}")
                raise EmbeddingError(f"Ollama embedding failed: {str(e)}")
    
    async def _get_ollama_text_completion(self, prompt: str, model: str = None) -> str:
        """Get text completion from Ollama."""
        if not HTTPX_AVAILABLE:
            raise EmbeddingError("httpx not available for Ollama")
        
        if not model:
            model = settings.ollama_model
            
        url = f"{settings.ollama_base_url}/api/generate"
        
        async with httpx.AsyncClient() as client:
            try:
                response = await client.post(
                    url,
                    json={
                        "model": model,
                        "prompt": prompt,
                        "stream": False,
                        "options": {
                            "temperature": settings.ollama_temperature,
                            "num_predict": settings.ollama_max_tokens
                        }
                    },
                    timeout=60.0
                )
                response.raise_for_status()
                data = response.json()
                return data.get("response", "")
            except Exception as e:
                logger.error(f"Ollama text completion failed: {str(e)}")
                raise EmbeddingError(f"Ollama text completion failed: {str(e)}")
    
    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=4, max=10)
    )
    async def generate_embedding(
        self,
        text: str,
        model_name: Optional[str] = None,
        preprocess: bool = True
    ) -> List[float]:
        """Generate embedding for a single text."""
        try:
            # Preprocess text if requested
            if preprocess:
                text = self.text_processor.clean_legal_text(text)
            
            # Generate embedding based on AI provider configuration
            if settings.ai_provider == "ollama":
                # Use Ollama for embeddings
                embedding = await self._get_ollama_embedding(text, model_name)
            elif settings.ai_provider == "openai":
                # Use OpenAI embeddings
                if not model_name:
                    model_name = settings.openai_embeddings_model
                client = self._get_openai_client()
                response = await client.embeddings.create(
                    model=model_name,
                    input=text
                )
                embedding = response.data[0].embedding
            else:
                # Default to local models (sentence-transformers)
                if not model_name:
                    model_name = settings.huggingface_model_name
                model = self._get_model(model_name)
                embedding = model.encode(text, convert_to_numpy=True).tolist()
            
            logger.debug(f"Generated embedding with dimension: {len(embedding)} using {settings.ai_provider}")
            return embedding
            
        except Exception as e:
            logger.error(f"Embedding generation failed", error=str(e))
            raise EmbeddingError(f"Failed to generate embedding: {str(e)}")
    
    async def generate_embeddings_batch(
        self,
        texts: List[str],
        model_name: Optional[str] = None,
        preprocess: bool = True,
        batch_size: int = 32
    ) -> List[List[float]]:
        """Generate embeddings for multiple texts."""
        try:
            # Use default model if not specified
            if not model_name:
                model_name = settings.huggingface_model_name
            
            # Preprocess texts if requested
            if preprocess:
                texts = [self.text_processor.clean_legal_text(text) for text in texts]
            
            embeddings = []
            
            # Process in batches
            for i in range(0, len(texts), batch_size):
                batch = texts[i:i + batch_size]
                
                if model_name.startswith("text-embedding-"):
                    # OpenAI embeddings
                    client = self._get_openai_client()
                    response = await client.embeddings.create(
                        model=model_name,
                        input=batch
                    )
                    batch_embeddings = [item.embedding for item in response.data]
                else:
                    # Local models
                    model = self._get_model(model_name)
                    batch_embeddings = model.encode(batch, convert_to_numpy=True).tolist()
                
                embeddings.extend(batch_embeddings)
            
            logger.info(f"Generated {len(embeddings)} embeddings")
            return embeddings
            
        except Exception as e:
            logger.error(f"Batch embedding generation failed", error=str(e))
            raise EmbeddingError(f"Failed to generate embeddings: {str(e)}")
    
    def compute_similarity(
        self,
        embedding1: List[float],
        embedding2: List[float]
    ) -> float:
        """Compute cosine similarity between two embeddings."""
        try:
            # Convert to numpy arrays
            vec1 = np.array(embedding1)
            vec2 = np.array(embedding2)
            
            # Compute cosine similarity
            dot_product = np.dot(vec1, vec2)
            norm1 = np.linalg.norm(vec1)
            norm2 = np.linalg.norm(vec2)
            
            if norm1 == 0 or norm2 == 0:
                return 0.0
            
            similarity = dot_product / (norm1 * norm2)
            return float(similarity)
            
        except Exception as e:
            logger.error(f"Similarity computation failed", error=str(e))
            raise EmbeddingError(f"Failed to compute similarity: {str(e)}")
    
    def get_model_info(self, model_name: Optional[str] = None) -> Dict[str, Any]:
        """Get information about the embedding model."""
        if not model_name:
            model_name = settings.huggingface_model_name
        
        info = {
            "model_name": model_name,
            "type": "openai" if model_name.startswith("text-embedding-") else "local",
            "dimension": settings.embedding_dimension
        }
        
        if model_name.startswith("sentence-transformers/"):
            try:
                model = self._get_model(model_name)
                info["dimension"] = model.get_sentence_embedding_dimension()
                info["max_seq_length"] = model.max_seq_length
            except Exception:
                pass
        
        return info


# Global embedding service instance
embedding_service = EmbeddingService()