"""Database configuration and session management."""

from typing import AsyncGenerator

from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, async_sessionmaker
from sqlalchemy.orm import declarative_base

from app.core.config import settings
from app.core.logging import get_logger

logger = get_logger(__name__)

# Create async engine (lazily initialized)
engine = None

def get_engine():
    """Get or create database engine."""
    global engine
    if engine is None:
        engine = create_async_engine(
            settings.database_url,
            echo=settings.environment == "development",
            pool_size=settings.db_pool_size,
            max_overflow=settings.db_max_overflow,
            pool_pre_ping=True,
        )
    return engine

# Create async session factory (lazily initialized)
AsyncSessionLocal = None

def get_session_factory():
    """Get or create session factory."""
    global AsyncSessionLocal
    if AsyncSessionLocal is None:
        AsyncSessionLocal = async_sessionmaker(
            get_engine(),
            class_=AsyncSession,
            expire_on_commit=False,
            autocommit=False,
            autoflush=False,
        )
    return AsyncSessionLocal

# Create declarative base
Base = declarative_base()


async def init_db() -> None:
    """Initialize database."""
    try:
        # Test connection
        engine = get_engine()
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        logger.info("Database initialized successfully")
    except Exception as e:
        logger.error("Failed to initialize database", error=str(e))
        raise


async def close_db() -> None:
    """Close database connections."""
    global engine
    if engine is not None:
        await engine.dispose()
        logger.info("Database connections closed")


async def get_session() -> AsyncGenerator[AsyncSession, None]:
    """Get database session."""
    session_factory = get_session_factory()
    async with session_factory() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()