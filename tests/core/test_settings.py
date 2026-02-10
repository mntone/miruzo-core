import shutil
from pathlib import Path

import pytest

import app.config.environments as settings_module
from app.config.constants import DAILY_LOVE_USED_MAXIMUM
from app.config.environments import Settings
from app.config.quota import QuotaConfig


def _app_base_dir() -> Path:
	return Path(settings_module.__file__).resolve().parent.parent


def _call_normalize_sqlite_url(value: str) -> str:
	proxy = Settings._normalize_sqlite_url.__func__
	method = proxy.wrapped.__get__(None, Settings)
	return method(value)


def _call_validate_quota(value: QuotaConfig) -> QuotaConfig:
	proxy = Settings._validate_quota.__func__
	method = proxy.wrapped.__get__(None, Settings)  # pyright: ignore[reportFunctionMemberAccess]
	return method(value)


def test_normalize_sqlite_url_resolves_relative_path_and_creates_dirs() -> None:
	target_dir = _app_base_dir() / '__tests__' / 'normalize-sqlite-url'
	db_path = target_dir / 'db.sqlite'
	value = 'sqlite:///__tests__/normalize-sqlite-url/db.sqlite'

	try:
		result = _call_normalize_sqlite_url(value)
		assert result == f'sqlite:///{db_path}'
		assert target_dir.exists()
	finally:
		if target_dir.exists():
			shutil.rmtree(target_dir)
		parent = target_dir.parent
		if parent.exists() and not any(parent.iterdir()):
			parent.rmdir()


def test_normalize_sqlite_url_keeps_absolute_path(tmp_path: Path) -> None:
	db_path = tmp_path / 'abs-db.sqlite'
	value = f'sqlite:///{db_path}'
	assert _call_normalize_sqlite_url(value) == value


@pytest.mark.parametrize(
	'value',
	[
		'postgresql://user:pass@localhost:5432/db',
		'sqlite+aiosqlite:///relative/path.sqlite',
	],
)
def test_normalize_sqlite_url_leaves_non_standard_sqlite_values(value: str) -> None:
	assert _call_normalize_sqlite_url(value) == value


def test_validate_quota_accepts_daily_love_limit_at_upper_bound() -> None:
	quota = QuotaConfig(daily_love_limit=DAILY_LOVE_USED_MAXIMUM)
	assert _call_validate_quota(quota) == quota


def test_validate_quota_rejects_daily_love_limit_above_upper_bound() -> None:
	quota = QuotaConfig(daily_love_limit=DAILY_LOVE_USED_MAXIMUM + 1)
	with pytest.raises(
		ValueError,
		match=f'quota.daily_love_limit must be <= {DAILY_LOVE_USED_MAXIMUM}',
	):
		_call_validate_quota(quota)
