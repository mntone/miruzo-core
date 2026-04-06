from typing import Any, Generator

from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.databases.sqlite_version import verify_sqlite_supports_returning_and_strict

if env.database_backend == DatabaseBackend.SQLITE:
	engine = create_engine(
		env.database_url,
		connect_args={'check_same_thread': False},
		echo=False,
		future=True,
	)

	with engine.connect() as conn:
		sqlite_version = conn.exec_driver_sql('SELECT sqlite_version();').scalar_one()
		if not isinstance(sqlite_version, str):
			raise RuntimeError('Failed to read SQLite version')
		verify_sqlite_supports_returning_and_strict(sqlite_version)

		conn.exec_driver_sql('PRAGMA journal_mode=WAL;')
		conn.exec_driver_sql('PRAGMA synchronous=NORMAL;')
		conn.exec_driver_sql('PRAGMA wal_autocheckpoint=1000;')

elif env.database_backend == DatabaseBackend.POSTGRE_SQL:
	engine = create_engine(
		env.database_url,
		connect_args={'options': '-c timezone=utc'},
		echo=False,
		future=True,
		max_overflow=20,
		pool_size=10,
	)

else:
	raise ValueError(f'Unsupported database type: {env.database_backend}')


def create_session() -> Session:
	session = Session(engine)

	return session


def get_session() -> Generator[Session, Any, None]:
	session = create_session()
	try:
		yield session
	finally:
		session.close()
