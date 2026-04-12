# pyright: reportAttributeAccessIssue=false

import builtins
import sys
from types import ModuleType
from typing import Any

import pytest

import app.databases.database as database_module


@pytest.mark.parametrize(
	('dsn', 'expected_conninfo_arg'),
	[
		('postgresql://user:pass@localhost:5432/app', 'postgresql://user:pass@localhost:5432/app'),
		('postgresql+psycopg://user:pass@localhost:5432/app', 'postgresql://user:pass@localhost:5432/app'),
	],
)
def test_create_postgres_engine_uses_psycopg_pool(
	monkeypatch: pytest.MonkeyPatch,
	dsn: str,
	expected_conninfo_arg: str,
) -> None:
	create_engine_calls: list[tuple[str, dict[str, Any]]] = []
	conninfo_to_dict_calls: list[str] = []
	make_conninfo_calls: list[dict[str, Any]] = []
	pool_ctor_calls: list[dict[str, Any]] = []

	def _conninfo_to_dict(conninfo: str) -> dict[str, Any]:
		conninfo_to_dict_calls.append(conninfo)
		return {'options': '-c statement_timeout=1000'}

	def _make_conninfo(**kwargs: Any) -> str:
		make_conninfo_calls.append(kwargs)
		return 'postgresql://normalized'

	psycopg_conninfo_module = ModuleType('psycopg.conninfo')
	psycopg_conninfo_module.conninfo_to_dict = _conninfo_to_dict
	psycopg_conninfo_module.make_conninfo = _make_conninfo

	psycopg_module = ModuleType('psycopg')
	psycopg_module.__path__ = []  # type: ignore[attr-defined]
	psycopg_module.conninfo = psycopg_conninfo_module

	class _DummyConnectionPool:
		def __init__(self, **kwargs: Any) -> None:
			pool_ctor_calls.append(kwargs)

		def getconn(self) -> object:
			return object()

	psycopg_pool_module = ModuleType('psycopg_pool')
	psycopg_pool_module.ConnectionPool = _DummyConnectionPool

	monkeypatch.setitem(sys.modules, 'psycopg', psycopg_module)
	monkeypatch.setitem(sys.modules, 'psycopg.conninfo', psycopg_conninfo_module)
	monkeypatch.setitem(sys.modules, 'psycopg_pool', psycopg_pool_module)

	def _fake_create_engine(url: str, **kwargs: Any) -> str:
		create_engine_calls.append((url, kwargs))
		return 'engine'

	monkeypatch.setattr(database_module, 'create_engine', _fake_create_engine)
	monkeypatch.setattr(database_module.event, 'listens_for', lambda *_: lambda fn: fn)

	engine = database_module._create_postgres_engine(dsn, pool_size=2, max_overflow=3)

	assert engine == 'engine'
	assert conninfo_to_dict_calls == [expected_conninfo_arg]
	assert make_conninfo_calls == [{'options': '-c statement_timeout=1000 -c TimeZone=UTC'}]
	assert pool_ctor_calls == [
		{
			'conninfo': 'postgresql://normalized',
			'close_returns': True,
			'min_size': 2,
			'max_size': 5,
		},
	]
	assert create_engine_calls[0][0] == 'postgresql+psycopg://'
	assert create_engine_calls[0][1]['poolclass'] is database_module.NullPool
	assert callable(create_engine_calls[0][1]['creator'])


def test_create_postgres_engine_rejects_non_postgresql_dsn() -> None:
	with pytest.raises(RuntimeError, match='Unsupported PostgreSQL DSN'):
		database_module._create_postgres_engine('mysql://user:pass@localhost:3306/app')


def test_create_postgres_engine_raises_runtime_error_when_psycopg_pool_missing(
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	original_import = builtins.__import__

	def _fake_import(
		name: str,
		globals: dict[str, Any] | None = None,
		locals: dict[str, Any] | None = None,
		fromlist: tuple[str, ...] = (),
		level: int = 0,
	) -> Any:
		if name == 'psycopg_pool':
			raise ModuleNotFoundError('No module named psycopg_pool')
		return original_import(name, globals, locals, fromlist, level)

	monkeypatch.setattr(builtins, '__import__', _fake_import)

	with pytest.raises(RuntimeError, match='requires psycopg3 pool'):
		database_module._create_postgres_engine('postgresql://user:pass@localhost:5432/app')


def test_create_postgres_engine_raises_runtime_error_when_psycopg_missing(
	monkeypatch: pytest.MonkeyPatch,
) -> None:
	psycopg_pool_module = ModuleType('psycopg_pool')

	class _DummyConnectionPool:
		def __init__(self, **_: Any) -> None:
			pass

	psycopg_pool_module.ConnectionPool = _DummyConnectionPool

	original_import = builtins.__import__

	def _fake_import(
		name: str,
		globals: dict[str, Any] | None = None,
		locals: dict[str, Any] | None = None,
		fromlist: tuple[str, ...] = (),
		level: int = 0,
	) -> Any:
		if name == 'psycopg_pool':
			return psycopg_pool_module
		if name == 'psycopg.conninfo':
			raise ModuleNotFoundError('No module named psycopg')
		return original_import(name, globals, locals, fromlist, level)

	monkeypatch.setattr(builtins, '__import__', _fake_import)

	with pytest.raises(RuntimeError, match='requires psycopg3'):
		database_module._create_postgres_engine('postgresql://user:pass@localhost:5432/app')
