"""
Configuration Module

Handles application configuration using Pydantic settings management.
Supports environment variables, validation, and multiple environments.
"""

import os
from typing import List, Optional, Dict, Any
from pathlib import Path
from functools import lru_cache

from pydantic import (
    BaseSettings,
    PostgresDsn,
    RedisDsn,
    AnyHttpUrl,
    validator,
    Field
)


class Settings(BaseSettings):
    """
    Application settings with validation and environment variable support.
    """
    
    # Application
    APP_NAME: str = "Python Complex API"
    APP_VERSION: str = "1.0.0"
    APP_ENV: str = Field("development", env="APP_ENV")
    DEBUG: bool = Field(False, env="DEBUG")
    
    # Server
    HOST: str = Field("0.0.0.0", env="HOST")
    PORT: int = Field(8000, env="PORT")
    WORKERS: int = Field(4, env="WORKERS")
    
    # Security
    SECRET_KEY: str = Field(..., env="SECRET_KEY")
    JWT_ALGORITHM: str = "HS256"
    JWT_EXPIRATION: int = 3600  # 1 hour
    JWT_REFRESH_EXPIRATION: int = 86400  # 24 hours
    
    # CORS
    CORS_ORIGINS: List[AnyHttpUrl] = Field(
        ["http://localhost:3000"],
        env="CORS_ORIGINS"
    )
    
    @validator("CORS_ORIGINS", pre=True)
    def parse_cors_origins(cls, v):
        """Parse CORS origins from string or list."""
        if isinstance(v, str):
            return [origin.strip() for origin in v.split(",")]
        return v
    
    # Database
    POSTGRES_HOST: str = Field("localhost", env="POSTGRES_HOST")
    POSTGRES_PORT: int = Field(5432, env="POSTGRES_PORT")
    POSTGRES_USER: str = Field("postgres", env="POSTGRES_USER")
    POSTGRES_PASSWORD: str = Field("postgres", env="POSTGRES_PASSWORD")
    POSTGRES_DB: str = Field("app_db", env="POSTGRES_DB")
    
    @property
    def DATABASE_URL(self) -> str:
        """Build database URL from components."""
        return f"postgresql://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"
    
    # Redis
    REDIS_HOST: str = Field("localhost", env="REDIS_HOST")
    REDIS_PORT: int = Field(6379, env="REDIS_PORT")
    REDIS_DB: int = Field(0, env="REDIS_DB")
    REDIS_PASSWORD: Optional[str] = Field(None, env="REDIS_PASSWORD")
    
    @property
    def REDIS_URL(self) -> str:
        """Build Redis URL from components."""
        if self.REDIS_PASSWORD:
            return f"redis://:{self.REDIS_PASSWORD}@{self.REDIS_HOST}:{self.REDIS_PORT}/{self.REDIS_DB}"
        return f"redis://{self.REDIS_HOST}:{self.REDIS_PORT}/{self.REDIS_DB}"
    
    # Email
    SMTP_HOST: str = Field("smtp.gmail.com", env="SMTP_HOST")
    SMTP_PORT: int = Field(587, env="SMTP_PORT")
    SMTP_USER: str = Field(..., env="SMTP_USER")
    SMTP_PASSWORD: str = Field(..., env="SMTP_PASSWORD")
    EMAIL_FROM: str = Field("noreply@example.com", env="EMAIL_FROM")
    
    # External APIs
    STRIPE_API_KEY: Optional[str] = Field(None, env="STRIPE_API_KEY")
    STRIPE_WEBHOOK_SECRET: Optional[str] = Field(None, env="STRIPE_WEBHOOK_SECRET")
    
    AWS_ACCESS_KEY_ID: Optional[str] = Field(None, env="AWS_ACCESS_KEY_ID")
    AWS_SECRET_ACCESS_KEY: Optional[str] = Field(None, env="AWS_SECRET_ACCESS_KEY")
    AWS_REGION: str = Field("us-east-1", env="AWS_REGION")
    S3_BUCKET: Optional[str] = Field(None, env="S3_BUCKET")
    
    # Rate Limiting
    RATE_LIMIT_ENABLED: bool = Field(True, env="RATE_LIMIT_ENABLED")
    RATE_LIMIT_REQUESTS: int = Field(100, env="RATE_LIMIT_REQUESTS")
    RATE_LIMIT_PERIOD: int = Field(60, env="RATE_LIMIT_PERIOD")  # seconds
    
    # Logging
    LOG_LEVEL: str = Field("INFO", env="LOG_LEVEL")
    LOG_FORMAT: str = Field("json", env="LOG_FORMAT")
    LOG_FILE: Optional[str] = Field(None, env="LOG_FILE")
    
    # Monitoring
    SENTRY_DSN: Optional[str] = Field(None, env="SENTRY_DSN")
    PROMETHEUS_ENABLED: bool = Field(True, env="PROMETHEUS_ENABLED")
    
    # Feature Flags
    FEATURE_SIGNUP: bool = Field(True, env="FEATURE_SIGNUP")
    FEATURE_NOTIFICATIONS: bool = Field(True, env="FEATURE_NOTIFICATIONS")
    FEATURE_BETA: bool = Field(False, env="FEATURE_BETA")
    
    class Config:
        """Pydantic configuration."""
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = True
    
    @validator("APP_ENV")
    def validate_environment(cls, v):
        """Validate environment name."""
        allowed = {"development", "staging", "production", "test"}
        if v not in allowed:
            raise ValueError(f"Environment must be one of {allowed}")
        return v
    
    @validator("LOG_LEVEL")
    def validate_log_level(cls, v):
        """Validate log level."""
        allowed = {"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}
        if v.upper() not in allowed:
            raise ValueError(f"Log level must be one of {allowed}")
        return v.upper()
    
    @property
    def is_development(self) -> bool:
        """Check if running in development."""
        return self.APP_ENV == "development"
    
    @property
    def is_production(self) -> bool:
        """Check if running in production."""
        return self.APP_ENV == "production"
    
    @property
    def is_testing(self) -> bool:
        """Check if running in testing."""
        return self.APP_ENV == "test"


@lru_cache()
def get_settings() -> Settings:
    """
    Get cached settings instance.
    Use this function to access settings throughout the application.
    """
    return Settings()


settings = get_settings()