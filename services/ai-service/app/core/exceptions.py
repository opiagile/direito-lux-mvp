"""Custom exceptions for AI Service."""

from typing import Any, Dict, Optional


class AIServiceError(Exception):
    """Base exception for AI Service."""
    
    def __init__(
        self,
        message: str,
        code: str = "AI_SERVICE_ERROR",
        details: Optional[Dict[str, Any]] = None
    ):
        self.message = message
        self.code = code
        self.details = details or {}
        super().__init__(self.message)


class ConfigurationError(AIServiceError):
    """Configuration related errors."""
    
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, "CONFIGURATION_ERROR", details)


class EmbeddingError(AIServiceError):
    """Embedding generation errors."""
    
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, "EMBEDDING_ERROR", details)


class SearchError(AIServiceError):
    """Search operation errors."""
    
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, "SEARCH_ERROR", details)


class DocumentGenerationError(AIServiceError):
    """Document generation errors."""
    
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, "DOCUMENT_GENERATION_ERROR", details)


class ExternalServiceError(AIServiceError):
    """External service communication errors."""
    
    def __init__(
        self,
        service: str,
        message: str,
        status_code: Optional[int] = None,
        details: Optional[Dict[str, Any]] = None
    ):
        details = details or {}
        details.update({"service": service, "status_code": status_code})
        super().__init__(message, "EXTERNAL_SERVICE_ERROR", details)


class RateLimitError(AIServiceError):
    """Rate limit exceeded errors."""
    
    def __init__(
        self,
        message: str = "Rate limit exceeded",
        retry_after: Optional[int] = None,
        details: Optional[Dict[str, Any]] = None
    ):
        details = details or {}
        if retry_after:
            details["retry_after"] = retry_after
        super().__init__(message, "RATE_LIMIT_ERROR", details)


class QuotaExceededError(AIServiceError):
    """Quota exceeded errors."""
    
    def __init__(
        self,
        resource: str,
        limit: int,
        current: int,
        details: Optional[Dict[str, Any]] = None
    ):
        message = f"Quota exceeded for {resource}. Limit: {limit}, Current: {current}"
        details = details or {}
        details.update({
            "resource": resource,
            "limit": limit,
            "current": current
        })
        super().__init__(message, "QUOTA_EXCEEDED_ERROR", details)


class ValidationError(AIServiceError):
    """Input validation errors."""
    
    def __init__(self, message: str, field: str, details: Optional[Dict[str, Any]] = None):
        details = details or {}
        details["field"] = field
        super().__init__(message, "VALIDATION_ERROR", details)


class AuthenticationError(AIServiceError):
    """Authentication errors."""
    
    def __init__(self, message: str = "Authentication failed", details: Optional[Dict[str, Any]] = None):
        super().__init__(message, "AUTHENTICATION_ERROR", details)


class AuthorizationError(AIServiceError):
    """Authorization errors."""
    
    def __init__(self, message: str = "Insufficient permissions", details: Optional[Dict[str, Any]] = None):
        super().__init__(message, "AUTHORIZATION_ERROR", details)