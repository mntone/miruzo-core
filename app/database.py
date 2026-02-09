from typing import Any, Generator

from alembic.config import Config
from alembic.runtime.migration import MigrationContext
from alembic.script import ScriptDirectory
from sqlalchemy import create_engine
from sqlmodel import Session, SQLModel

from app.config.environments import DatabaseBackend, env
from app.utils.database.sqlite_version import verify_sqlite_supports_returning


def _verify_schema_version(db_url: str, alembic_ini: str) -> None:
	engine = create_engine(db_url)
	alembic_cfg = Config(alembic_ini)
	script = ScriptDirectory.from_config(alembic_cfg)

	with engine.connect() as conn:
		context = MigrationContext.configure(conn)
		db_revision = context.get_current_revision()

	head_revisions = script.get_heads()

	if db_revision is None:
		raise RuntimeError('database is not under alembic control')

	if db_revision not in head_revisions:
		raise RuntimeError(
			f'schema version mismatch: db={db_revision}, code={head_revisions}',
		)


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
		verify_sqlite_supports_returning(sqlite_version)

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


def init_database() -> None:
	SQLModel.metadata.create_all(engine)

	_verify_schema_version(env.database_url, 'alembic.ini')


def create_session() -> Session:
	session = Session(engine)

	return session


def get_session() -> Generator[Session, Any, None]:
	session = create_session()
	try:
		yield session
	finally:
		session.close()
