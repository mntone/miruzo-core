from enum import Enum
from pathlib import Path
from zoneinfo import ZoneInfo

from pydantic import ValidationInfo, field_validator, model_validator
from pydantic_settings import BaseSettings, SettingsConfigDict

from app.config.constants import DAILY_LOVE_USED_MAXIMUM
from app.config.quota import QuotaConfig
from app.config.score import ScoreConfig
from app.config.time import TimeConfig
from app.config.variant import DEFAULT_VARIANT_LAYERS, VariantLayerSpec
from app.utils.files.permissions import ensure_directory_access


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
	base_timezone: ZoneInfo | None = None  # None=local timezone

	database_backend: DatabaseBackend = DatabaseBackend.SQLITE
	database_url: str = 'sqlite:///../var/miruzo.sqlite'

	media_root: Path = Path('../var/media')
	public_media_root: str = '/media/'

	gataku_root: Path = Path('../../gataku')
	gataku_assets_root: Path = Path('../../gataku/out/downloads')
	gataku_symlink_dirname: str = 'gataku'

	quota: QuotaConfig = QuotaConfig()
	score: ScoreConfig = ScoreConfig()
	time: TimeConfig = TimeConfig()
	variant_layers: tuple[VariantLayerSpec, ...] = DEFAULT_VARIANT_LAYERS

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
	@field_validator('base_timezone', mode='before')
	def _normalize_base_timezone(cls, value: object) -> object:
		if value is None or isinstance(value, ZoneInfo):
			return value
		if isinstance(value, str):
			return ZoneInfo(value)
		raise ValueError('base_timezone must be a zoneinfo')

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

	@classmethod
	@field_validator('media_root', 'gataku_root', 'gataku_assets_root')
	def _validate_path(cls, value: object, info: ValidationInfo) -> Path:
		if isinstance(value, Path):
			path = value
		elif isinstance(value, str):
			path = Path(value)
		else:
			raise ValueError(info.field_name.__str__() + ' must be a path')

		return path

	@classmethod
	@field_validator('quota')
	def _validate_quota(cls, value: QuotaConfig) -> QuotaConfig:
		if value.daily_love_limit > DAILY_LOVE_USED_MAXIMUM:
			raise ValueError(
				f'quota.daily_love_limit must be <= {DAILY_LOVE_USED_MAXIMUM}',
			)

		return value

	@model_validator(mode='after')
	def _normalize(self) -> 'Settings':
		self.media_root = self.media_root.resolve()
		self.gataku_root = self.gataku_root.resolve()

		resolved_gataku_assets_root = self.gataku_assets_root.resolve()
		ensure_directory_access(resolved_gataku_assets_root)
		self.gataku_assets_root = resolved_gataku_assets_root

		return self


env = Settings()
