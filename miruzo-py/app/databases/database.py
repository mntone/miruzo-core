from typing import Any

from sqlalchemy import Engine, NullPool, create_engine, event
from sqlalchemy.engine import make_url
from sqlalchemy.exc import ArgumentError
from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.databases.mysql_version import verify_mysql_supports_check_constraints
from app.databases.sqlite_version import verify_sqlite_supports_returning_and_strict

MYSQL_CHARSET = 'utf8mb4'
MYSQL_COLLATION = 'utf8mb4_0900_bin'


def _create_mysql_engine(
	dsn: str,
	*,
	pool_size: int = 1,
	max_overflow: int = 2,
) -> Engine:
	try:
		parsed_dsn = make_url(dsn)
	except ArgumentError as exc:
		raise RuntimeError('Unsupported MySQL DSN') from exc

	if not parsed_dsn.drivername.startswith('mysql'):
		raise RuntimeError('Unsupported MySQL DSN')

	if parsed_dsn.drivername == 'mysql':
		parsed_dsn = parsed_dsn.set(drivername='mysql+mysqldb')

	if parsed_dsn.drivername == 'mysql+mysqldb' and parsed_dsn.host == 'localhost':
		parsed_dsn = parsed_dsn.set(host='127.0.0.1')

	engine = create_engine(
		parsed_dsn,
		connect_args={
			'init_command': "SET sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY',time_zone='+00:00'",
		},
		echo=False,
		future=True,
		max_overflow=max_overflow,
		pool_size=pool_size,
	)

	@event.listens_for(engine, 'connect')
	def _set_mysql_charset_and_collation(dbapi_connection: Any, _: object) -> None:
		with dbapi_connection.cursor() as cursor:
			cursor.execute(f'SET NAMES {MYSQL_CHARSET} COLLATE {MYSQL_COLLATION}')

	with engine.connect() as conn:
		mysql_version = conn.exec_driver_sql('SELECT VERSION();').scalar_one()
		if not isinstance(mysql_version, str):
			raise RuntimeError('Failed to read MySQL version')
		verify_mysql_supports_check_constraints(mysql_version)

	return engine


def _create_postgres_engine(
	dsn: str,
	*,
	pool_size: int = 1,
	max_overflow: int = 2,
) -> Engine:
	if not dsn.startswith('postgresql'):
		raise RuntimeError('Unsupported PostgreSQL DSN')

	if dsn.startswith('postgresql+psycopg://'):
		dsn = dsn.replace('postgresql+psycopg://', 'postgresql://', 1)

	if dsn.startswith('postgresql://'):

		def _require_psycopg_conninfo() -> Any:
			try:
				import psycopg.conninfo as psycopg_conninfo
			except ModuleNotFoundError as exc:
				raise RuntimeError('PostgreSQL backend requires psycopg3.') from exc
			return psycopg_conninfo

		def _require_psycopg_pool() -> Any:
			try:
				import psycopg_pool
			except ModuleNotFoundError as exc:
				raise RuntimeError('PostgreSQL backend requires psycopg3 pool.') from exc
			return psycopg_pool

		def _ensure_utc_timezone(conninfo: str) -> str:
			psycopg_conninfo = _require_psycopg_conninfo()
			parameters = psycopg_conninfo.conninfo_to_dict(conninfo)

			options = parameters.get('options', '')
			if 'TimeZone=' not in options:
				if options:
					options = f'{options} -c TimeZone=UTC'
				else:
					options = '-c TimeZone=UTC'
				parameters['options'] = options

			return psycopg_conninfo.make_conninfo(**parameters)

		psycopg_pool = _require_psycopg_pool()
		conninfo = _ensure_utc_timezone(dsn)
		pool = psycopg_pool.ConnectionPool(
			conninfo=conninfo,
			close_returns=True,  # Return "closed" active connections to the pool
			min_size=pool_size,
			max_size=pool_size + max_overflow,
		)
		engine = create_engine(
			'postgresql+psycopg://',
			creator=pool.getconn,
			poolclass=NullPool,
		)

		@event.listens_for(engine, 'engine_disposed')
		def _close_psycopg_pool(_: Engine) -> None:
			pool.close()

		return engine

	engine = create_engine(
		dsn,
		connect_args={'options': '-c TimeZone=utc'},
		echo=False,
		future=True,
		max_overflow=max_overflow,
		pool_size=pool_size,
	)
	return engine


def _create_sqlite3_engine(
	dsn: str,
	*,
	pool_size: int = 1,
) -> Engine:
	if not dsn.startswith(('sqlite://', 'sqlite+pysqlite://')):
		raise RuntimeError('Unsupported SQLite DSN')

	from sqlite3 import Connection

	engine = create_engine(
		dsn,
		connect_args={'check_same_thread': False},
		echo=False,
		future=True,
		pool_size=pool_size,
	)

	@event.listens_for(engine, 'connect')
	def _set_sqlite_pragmas(dbapi_connection: Connection, _: object) -> None:
		cursor = dbapi_connection.cursor()
		try:
			cursor.execute('PRAGMA foreign_keys=1;')
			cursor.execute('PRAGMA wal_autocheckpoint=100;')
		finally:
			cursor.close()

	with engine.connect() as conn:
		sqlite_version = conn.exec_driver_sql('SELECT sqlite_version();').scalar_one()
		if not isinstance(sqlite_version, str):
			raise RuntimeError('Failed to read SQLite version')
		verify_sqlite_supports_returning_and_strict(sqlite_version)

		conn.exec_driver_sql('PRAGMA journal_mode=WAL;')

	return engine


match env.database_backend:
	case DatabaseBackend.MYSQL:
		engine = _create_mysql_engine(env.database_url)
	case DatabaseBackend.POSTGRE_SQL:
		engine = _create_postgres_engine(env.database_url)
	case DatabaseBackend.SQLITE:
		engine = _create_sqlite3_engine(env.database_url)
	case _:
		raise ValueError(f'Unsupported database type: {env.database_backend!r}')


def create_session() -> Session:
	session = Session(engine)

	return session
