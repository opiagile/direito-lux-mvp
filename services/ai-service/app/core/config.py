"""Configuration management for AI Service."""

from typing import List, Optional
from pydantic_settings import BaseSettings
from pydantic import Field, validator


class Settings(BaseSettings):
    """Application settings."""
    
    # Service Info
    service_name: str = Field(default="ai-service", env="SERVICE_NAME")
    version: str = Field(default="1.0.0", env="VERSION")
    environment: str = Field(default="development", env="ENVIRONMENT")
    port: int = Field(default=8000, env="PORT")
    log_level: str = Field(default="info", env="LOG_LEVEL")
    
    # Database
    db_host: str = Field(default="localhost", env="DB_HOST")
    db_port: int = Field(default=5432, env="DB_PORT")
    db_name: str = Field(default="direito_lux_dev", env="DB_NAME")
    db_user: str = Field(default="direito_lux", env="DB_USER")
    db_password: str = Field(env="DB_PASSWORD")
    db_pool_size: int = Field(default=20, env="DB_POOL_SIZE")
    db_max_overflow: int = Field(default=40, env="DB_MAX_OVERFLOW")
    
    # Redis
    redis_host: str = Field(default="localhost", env="REDIS_HOST")
    redis_port: int = Field(default=6379, env="REDIS_PORT")
    redis_password: Optional[str] = Field(default=None, env="REDIS_PASSWORD")
    redis_db: int = Field(default=1, env="REDIS_DB")
    
    # RabbitMQ
    rabbitmq_host: str = Field(default="localhost", env="RABBITMQ_HOST")
    rabbitmq_port: int = Field(default=5672, env="RABBITMQ_PORT")
    rabbitmq_user: str = Field(default="guest", env="RABBITMQ_USER")
    rabbitmq_password: str = Field(default="guest", env="RABBITMQ_PASSWORD")
    rabbitmq_vhost: str = Field(default="/", env="RABBITMQ_VHOST")
    rabbitmq_exchange: str = Field(default="direito_lux.events", env="RABBITMQ_EXCHANGE")
    rabbitmq_queue: str = Field(default="ai-service.events", env="RABBITMQ_QUEUE")
    
    # OpenAI
    openai_api_key: str = Field(env="OPENAI_API_KEY")
    openai_model: str = Field(default="gpt-3.5-turbo", env="OPENAI_MODEL")
    openai_embeddings_model: str = Field(default="text-embedding-ada-002", env="OPENAI_EMBEDDINGS_MODEL")
    openai_max_tokens: int = Field(default=2000, env="OPENAI_MAX_TOKENS")
    openai_temperature: float = Field(default=0.7, env="OPENAI_TEMPERATURE")
    
    # HuggingFace
    huggingface_token: Optional[str] = Field(default=None, env="HUGGINGFACE_TOKEN")
    huggingface_model_name: str = Field(
        default="sentence-transformers/all-MiniLM-L6-v2", 
        env="HUGGINGFACE_MODEL_NAME"
    )
    
    # Vector Store
    vector_store_type: str = Field(default="faiss", env="VECTOR_STORE_TYPE")
    pgvector_table_name: str = Field(default="legal_embeddings", env="PGVECTOR_TABLE_NAME")
    faiss_index_path: str = Field(default="/app/data/faiss_index", env="FAISS_INDEX_PATH")
    
    # AI Settings
    embedding_dimension: int = Field(default=384, env="EMBEDDING_DIMENSION")
    max_search_results: int = Field(default=50, env="MAX_SEARCH_RESULTS")
    similarity_threshold: float = Field(default=0.7, env="SIMILARITY_THRESHOLD")
    cache_ttl: int = Field(default=3600, env="CACHE_TTL")
    
    # Jurisprudence Collectors
    enable_stf_collector: bool = Field(default=True, env="ENABLE_STF_COLLECTOR")
    enable_stj_collector: bool = Field(default=True, env="ENABLE_STJ_COLLECTOR")
    enable_regional_collectors: bool = Field(default=False, env="ENABLE_REGIONAL_COLLECTORS")
    collector_batch_size: int = Field(default=100, env="COLLECTOR_BATCH_SIZE")
    collector_rate_limit: int = Field(default=10, env="COLLECTOR_RATE_LIMIT")
    
    # Security
    jwt_secret_key: str = Field(env="JWT_SECRET_KEY")
    jwt_algorithm: str = Field(default="HS256", env="JWT_ALGORITHM")
    jwt_expiration_minutes: int = Field(default=30, env="JWT_EXPIRATION_MINUTES")
    
    # External Services
    auth_service_url: str = Field(default="http://auth-service:8080", env="AUTH_SERVICE_URL")
    tenant_service_url: str = Field(default="http://tenant-service:8080", env="TENANT_SERVICE_URL")
    process_service_url: str = Field(default="http://process-service:8080", env="PROCESS_SERVICE_URL")
    datajud_service_url: str = Field(default="http://datajud-service:8080", env="DATAJUD_SERVICE_URL")
    
    # Observability
    jaeger_enabled: bool = Field(default=True, env="JAEGER_ENABLED")
    jaeger_agent_host: str = Field(default="jaeger", env="JAEGER_AGENT_HOST")
    jaeger_agent_port: int = Field(default=6831, env="JAEGER_AGENT_PORT")
    prometheus_enabled: bool = Field(default=True, env="PROMETHEUS_ENABLED")
    prometheus_port: int = Field(default=9090, env="PROMETHEUS_PORT")
    
    # CORS
    cors_origins: List[str] = Field(
        default=["http://localhost:3000", "http://localhost:8000"],
        env="CORS_ORIGINS"
    )
    cors_allow_credentials: bool = Field(default=True, env="CORS_ALLOW_CREDENTIALS")
    cors_allow_methods: List[str] = Field(
        default=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
        env="CORS_ALLOW_METHODS"
    )
    cors_allow_headers: List[str] = Field(default=["*"], env="CORS_ALLOW_HEADERS")
    
    @property
    def database_url(self) -> str:
        """Get PostgreSQL connection URL."""
        return f"postgresql+asyncpg://{self.db_user}:{self.db_password}@{self.db_host}:{self.db_port}/{self.db_name}"
    
    @property
    def redis_url(self) -> str:
        """Get Redis connection URL."""
        if self.redis_password:
            return f"redis://:{self.redis_password}@{self.redis_host}:{self.redis_port}/{self.redis_db}"
        return f"redis://{self.redis_host}:{self.redis_port}/{self.redis_db}"
    
    @property
    def rabbitmq_url(self) -> str:
        """Get RabbitMQ connection URL."""
        return f"amqp://{self.rabbitmq_user}:{self.rabbitmq_password}@{self.rabbitmq_host}:{self.rabbitmq_port}/{self.rabbitmq_vhost}"
    
    @validator("log_level")
    def validate_log_level(cls, v):
        """Validate log level."""
        allowed = ["debug", "info", "warning", "error", "critical"]
        if v.lower() not in allowed:
            raise ValueError(f"Log level must be one of {allowed}")
        return v.lower()
    
    @validator("environment")
    def validate_environment(cls, v):
        """Validate environment."""
        allowed = ["development", "staging", "production", "test"]
        if v.lower() not in allowed:
            raise ValueError(f"Environment must be one of {allowed}")
        return v.lower()
    
    class Config:
        """Pydantic config."""
        env_file = ".env"
        case_sensitive = False


# Global settings instance
settings = Settings()