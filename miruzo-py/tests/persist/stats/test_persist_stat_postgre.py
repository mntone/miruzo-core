import shutil
import socket
import subprocess
import time
import uuid
from collections.abc import Iterator
from typing import Any, Generator

import psycopg2
import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

from app.persist.stats.postgre import PostgreSQLStatsRepository

POSTGRES_IMAGE = 'postgres:18-alpine'
POSTGRES_DB = 'miruzo'
POSTGRES_USER = 'miruzo'
POSTGRES_PASSWORD = 'miruzo'


def _find_free_port() -> int:
	with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
		sock.bind(('', 0))
		return sock.getsockname()[1]


@pytest.fixture(scope='session')
def postgres_dsn() -> Iterator[str]:
	if shutil.which('docker') is None:
		pytest.skip('docker binary not available on PATH', allow_module_level=True)

	try:
		host_port = _find_free_port()
	except PermissionError:
		pytest.skip('cannot bind socket to discover free port (insufficient permissions)')

	container_name = f'miruzo-test-{uuid.uuid4().hex[:8]}'
	run_cmd = [
		'docker',
		'run',
		'--rm',
		'--name',
		container_name,
		'-e',
		f'POSTGRES_DB={POSTGRES_DB}',
		'-e',
		f'POSTGRES_USER={POSTGRES_USER}',
		'-e',
		f'POSTGRES_PASSWORD={POSTGRES_PASSWORD}',
		'-p',
		f'{host_port}:5432',
		'-d',
		POSTGRES_IMAGE,
	]
	subprocess.run(run_cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

	dsn = f'postgresql://{POSTGRES_USER}:{POSTGRES_PASSWORD}@localhost:{host_port}/{POSTGRES_DB}'
	for _ in range(30):
		try:
			with psycopg2.connect(dsn, connect_timeout=1):
				break
		except psycopg2.OperationalError:
			time.sleep(1)
	else:
		subprocess.run(['docker', 'rm', '-f', container_name], check=False)
		raise RuntimeError('PostgreSQL container did not become ready in time')

	try:
		yield dsn
	finally:
		subprocess.run(['docker', 'rm', '-f', container_name], check=False)


@pytest.fixture()
def session(postgres_dsn: str) -> Generator[Session, Any, None]:
	engine = create_engine(postgres_dsn, echo=False)
	SQLModel.metadata.create_all(engine)
	with Session(engine) as session:
		yield session


def test_get_or_create(session: Session) -> None:
	repo = PostgreSQLStatsRepository(session)
	image = add_image_record(session, 20)

	stats = repo.get_or_create(image.ingest_id, initial_score=42)
	assert stats.score == 42
	assert stats.view_count == 0
	assert stats.last_viewed_at is None

	stats_again = repo.get_or_create(image.ingest_id, initial_score=99)
	assert stats_again.ingest_id == image.ingest_id
	assert stats_again.score == 42
