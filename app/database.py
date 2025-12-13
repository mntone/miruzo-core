from typing import Any, Generator

from sqlalchemy import create_engine
from sqlmodel import Session, SQLModel

from app.core.settings import DatabaseBackend, settings

if settings.database_backend == DatabaseBackend.SQLITE:
	engine = create_engine(
		settings.database_url,
		connect_args={'check_same_thread': False},
		echo=False,
		future=True,
	)

	with engine.connect() as conn:
		conn.exec_driver_sql('PRAGMA journal_mode=WAL;')
		conn.exec_driver_sql('PRAGMA synchronous=NORMAL;')
		conn.exec_driver_sql('PRAGMA wal_autocheckpoint=1000;')

elif settings.database_backend == DatabaseBackend.POSTGRE_SQL:
	engine = create_engine(
		settings.database_url,
		connect_args={'options': '-c timezone=utc'},
		echo=False,
		future=True,
		max_overflow=20,
		pool_size=10,
	)

else:
	raise ValueError(f'Unsupported database type: {settings.db_type}')


def init_database() -> None:
	SQLModel.metadata.create_all(engine)


def get_session() -> Generator[Session, Any, None]:
	session = Session(engine)
	try:
		yield session
	finally:
		session.close()
