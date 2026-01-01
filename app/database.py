from typing import Any, Generator

from alembic.config import Config
from alembic.runtime.migration import MigrationContext
from alembic.script import ScriptDirectory
from sqlalchemy import create_engine
from sqlmodel import Session, SQLModel

from app.config.environments import DatabaseBackend, env


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
	_verify_schema_version(env.database_url, 'alembic.ini')

	SQLModel.metadata.create_all(engine)


def get_session() -> Generator[Session, Any, None]:
	session = Session(engine)
	try:
		yield session
	finally:
		session.close()
