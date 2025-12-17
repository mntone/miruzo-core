from enum import Enum
from pathlib import Path

from pydantic import field_validator
from pydantic_settings import BaseSettings, SettingsConfigDict

from app.config.variant import DEFAULT_VARIANT_LAYERS, VariantLayer


class DatabaseBackend(str, Enum):
	POSTGRE_SQL = 'postgres'
	SQLITE = 'sqlite'


class Environment(str, Enum):
	DEVELOPMENT = 'development'
	PRODUCTION = 'production'


_ALLOWED_ENVIRONMENTS = {env.value for env in Environment}


class Settings(BaseSettings):
	model_config = SettingsConfigDict(
		env_file=('.env', '.env.production'),
		env_file_encoding='utf-8',
	)

	environment: Environment = Environment.PRODUCTION

	database_backend: DatabaseBackend = DatabaseBackend.SQLITE
	database_url: str = 'sqlite:///var/miruzo.sqlite'

	media_root: Path = Path('./var/media')
	public_media_root: str = '/media'

	gataku_root: Path = Path('../gataku')
	gataku_assets_root: Path = Path('../gataku/out/downloads')

	variant_layers: tuple[VariantLayer, ...] = DEFAULT_VARIANT_LAYERS

	@property
	def debug(self) -> bool:
		return self.environment == Environment.DEVELOPMENT

	@classmethod
	@field_validator('environment', mode='before')
	def _normalize_environment(cls, value: object) -> object:
		if isinstance(value, Environment) or value is None:
			return value
		if isinstance(value, str):
			normalized = value.lower()
			if normalized in _ALLOWED_ENVIRONMENTS:
				return normalized
		raise ValueError("environment must be 'development' or 'production'")

	@classmethod
	@field_validator('database_url')
	def _normalize_sqlite_url(cls, value: str) -> str:
		if not value.startswith('sqlite:///'):
			return value

		raw_path = value.removeprefix('sqlite:///')
		path = Path(raw_path)

		if not path.is_absolute():
			baseDir = Path(__file__).resolve().parent.parent
			path = baseDir / path

		path.parent.mkdir(parents=True, exist_ok=True)

		return f'sqlite:///{path}'


env = Settings()
