"""Cache service for AI operations."""

import json
import hashlib
from typing import Any, Optional, List

import redis.asyncio as redis

from app.core.config import settings
from app.core.logging import get_logger

logger = get_logger(__name__)


class CacheService:
    """Redis-based cache service."""
    
    def __init__(self):
        """Initialize cache service."""
        self._redis: Optional[redis.Redis] = None
    
    async def connect(self) -> None:
        """Connect to Redis."""
        try:
            self._redis = await redis.from_url(
                settings.redis_url,
                decode_responses=True
            )
            await self._redis.ping()
            logger.info("Cache service connected")
        except Exception as e:
            logger.error(f"Failed to connect to Redis", error=str(e))
            # Don't raise - cache is optional
    
    async def disconnect(self) -> None:
        """Disconnect from Redis."""
        if self._redis:
            await self._redis.close()
            logger.info("Cache service disconnected")
    
    def _generate_key(self, prefix: str, data: Any) -> str:
        """Generate cache key from data."""
        # Convert data to string for hashing
        if isinstance(data, (dict, list)):
            data_str = json.dumps(data, sort_keys=True)
        else:
            data_str = str(data)
        
        # Generate hash
        hash_obj = hashlib.sha256(data_str.encode())
        hash_hex = hash_obj.hexdigest()[:16]  # Use first 16 chars
        
        return f"{prefix}:{hash_hex}"
    
    async def get(self, key: str) -> Optional[Any]:
        """Get value from cache."""
        if not self._redis:
            return None
        
        try:
            value = await self._redis.get(key)
            if value:
                return json.loads(value)
            return None
        except Exception as e:
            logger.error(f"Cache get error", key=key, error=str(e))
            return None
    
    async def set(
        self,
        key: str,
        value: Any,
        ttl: Optional[int] = None
    ) -> bool:
        """Set value in cache."""
        if not self._redis:
            return False
        
        try:
            value_str = json.dumps(value)
            if ttl:
                await self._redis.setex(key, ttl, value_str)
            else:
                await self._redis.set(key, value_str)
            return True
        except Exception as e:
            logger.error(f"Cache set error", key=key, error=str(e))
            return False
    
    async def delete(self, key: str) -> bool:
        """Delete value from cache."""
        if not self._redis:
            return False
        
        try:
            result = await self._redis.delete(key)
            return result > 0
        except Exception as e:
            logger.error(f"Cache delete error", key=key, error=str(e))
            return False
    
    async def get_embedding(self, text: str, model: str) -> Optional[List[float]]:
        """Get cached embedding."""
        key = self._generate_key(f"emb:{model}", text)
        result = await self.get(key)
        if result and isinstance(result, list):
            return result
        return None
    
    async def set_embedding(
        self,
        text: str,
        model: str,
        embedding: List[float]
    ) -> bool:
        """Cache embedding."""
        key = self._generate_key(f"emb:{model}", text)
        return await self.set(key, embedding, ttl=settings.cache_ttl)
    
    async def get_search_results(
        self,
        query: str,
        filters: dict
    ) -> Optional[List[dict]]:
        """Get cached search results."""
        key = self._generate_key("search", {"query": query, "filters": filters})
        return await self.get(key)
    
    async def set_search_results(
        self,
        query: str,
        filters: dict,
        results: List[dict]
    ) -> bool:
        """Cache search results."""
        key = self._generate_key("search", {"query": query, "filters": filters})
        return await self.set(key, results, ttl=3600)  # 1 hour
    
    async def get_analysis(
        self,
        process_id: str,
        analysis_type: str
    ) -> Optional[dict]:
        """Get cached analysis result."""
        key = f"analysis:{process_id}:{analysis_type}"
        return await self.get(key)
    
    async def set_analysis(
        self,
        process_id: str,
        analysis_type: str,
        result: dict
    ) -> bool:
        """Cache analysis result."""
        key = f"analysis:{process_id}:{analysis_type}"
        return await self.set(key, result, ttl=7200)  # 2 hours
    
    async def invalidate_pattern(self, pattern: str) -> int:
        """Invalidate all keys matching pattern."""
        if not self._redis:
            return 0
        
        try:
            keys = []
            async for key in self._redis.scan_iter(match=pattern):
                keys.append(key)
            
            if keys:
                return await self._redis.delete(*keys)
            return 0
        except Exception as e:
            logger.error(f"Cache invalidation error", pattern=pattern, error=str(e))
            return 0


# Global cache service instance
cache_service = CacheService()


async def init_cache() -> None:
    """Initialize cache service."""
    await cache_service.connect()


async def close_cache() -> None:
    """Close cache service."""
    await cache_service.disconnect()