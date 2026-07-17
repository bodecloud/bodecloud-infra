"""Platform settings from environment (R-ODR-6, R-LDR-1/2)."""

from __future__ import annotations

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="SYNTHORA_", env_file=".env", extra="ignore")

    database_url: str = "sqlite+aiosqlite:///./synthora.db"
    redis_url: str = "redis://localhost:6379/0"
    auth_mode: str = "none"  # none | session
    secret_key: str = "change-me"
    token_ttl_seconds: int = 60 * 60 * 24 * 7
    allow_registrations: bool = True
    max_concurrent_researches: int = 3
    cors_origins: str = "*"


settings = Settings()
