from sqlalchemy import create_engine, event
from sqlalchemy.orm import Session

from app.config.environments import DatabaseBackend, env
from app.databases.sqlite_version import verify_sqlite_supports_returning_and_strict

if env.database_backend == DatabaseBackend.SQLITE:
	from sqlite3 import Connection

	engine = create_engine(
		env.database_url,
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
