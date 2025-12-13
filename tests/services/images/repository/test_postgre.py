from __future__ import annotations

import shutil
import socket
import subprocess
import time
import uuid
from datetime import datetime, timezone
from typing import Any, Generator

import psycopg2
import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

from app.core.constants import DEFAULT_SCORE
from app.models.records import StatsRecord
from app.services.images.repository.postgre import PostgreSQLImageRepository

POSTGRES_IMAGE = 'postgres:18-alpine'
POSTGRES_DB = 'miruzo'
POSTGRES_USER = 'miruzo'
POSTGRES_PASSWORD = 'miruzo'


def _find_free_port() -> int:
	with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
		sock.bind(('', 0))
		return sock.getsockname()[1]


@pytest.fixture(scope='session')
def postgres_dsn() -> str:
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


def test_get_detail_with_stats(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 1)
	stats = StatsRecord(
		image_id=image.id,
		favorite=True,
		score=DEFAULT_SCORE,
		view_count=2,
		last_viewed_at=datetime.now(timezone.utc),
	)
	session.add(stats)
	session.commit()

	result = repo.get_detail_with_stats(image.id)

	assert result is not None
	image_record, stats_record = result
	assert image_record.id == image.id
	assert stats_record.favorite is True
	assert stats_record.view_count == 2


def test_upsert_stats_with_increment(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 2)

	stats = repo.upsert_stats_with_increment(image.id)
	assert stats.view_count == 1

	stats_again = repo.upsert_stats_with_increment(image.id)
	assert stats_again.view_count == 2


def test_update_score_and_favorite(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 3)
	repo.create_stats(image.id)

	score_response = repo.update_score(image.id, 5)
	assert score_response is not None
	assert score_response.score == DEFAULT_SCORE + 5

	favorite_response = repo.update_favorite(image.id, True)
	assert favorite_response is not None
	stats = session.get(StatsRecord, image.id)
	assert stats is not None
	assert stats.favorite is True
