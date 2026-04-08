from sqlalchemy import Engine, create_engine, event
from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.databases.mysql_version import verify_mysql_supports_check_constraints
from app.databases.sqlite_version import verify_sqlite_supports_returning_and_strict

MYSQL_CHARSET = 'utf8mb4'
MYSQL_COLLATION = 'utf8mb4_0900_bin'


def _create_mysql_engine(
	dsn: str,
	*,
	pool_size: int = 4,
	max_overflow: int = 8,
) -> Engine:
	from MySQLdb import Connection

	engine = create_engine(
		dsn,
		connect_args={
			'init_command': "SET sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY',time_zone='+00:00'",
		},
		echo=False,
		future=True,
		max_overflow=max_overflow,
		pool_size=pool_size,
	)

	@event.listens_for(engine, 'connect')
	def _set_mysql_charset_and_collation(dbapi_connection: Connection, _: object) -> None:
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
	pool_size: int = 4,
	max_overflow: int = 8,
) -> Engine:
	engine = create_engine(
		dsn,
		connect_args={'options': '-c timezone=utc'},
		echo=False,
		future=True,
		max_overflow=max_overflow,
		pool_size=pool_size,
	)
	return engine


def _create_sqlite3_engine(dsn: str) -> Engine:
	from sqlite3 import Connection

	engine = create_engine(
		dsn,
		connect_args={'check_same_thread': False},
		echo=False,
		future=True,
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
