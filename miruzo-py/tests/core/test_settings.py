import shutil
from pathlib import Path

import pytest

import app.config.environments as settings_module
from app.config.environments import Settings


def _app_base_dir() -> Path:
	return Path(settings_module.__file__).resolve().parent.parent


def _call_normalize_sqlite_url(value: str) -> str:
	proxy = Settings._normalize_sqlite_url.__func__
	method = proxy.wrapped.__get__(None, Settings)
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
