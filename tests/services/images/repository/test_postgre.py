import shutil
import socket
import subprocess
import time
import uuid
from collections.abc import Iterator
from datetime import datetime, timedelta, timezone
from typing import Any, Generator

import psycopg2
import pytest
from sqlmodel import Session, SQLModel, create_engine

from tests.services.images.utils import add_image_record

from app.config.constants import DEFAULT_SCORE
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


def test_get_list_orders_desc_and_returns_cursor(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	now = datetime.now(timezone.utc)
	first = add_image_record(session, 1, captured_at=now)
	second = add_image_record(session, 2, captured_at=now - timedelta(hours=1))
	third = add_image_record(session, 3, captured_at=now - timedelta(hours=2))

	items, cursor = repo.get_latest(cursor=None, limit=2)

	assert [item.ingest_id for item in items] == [first.ingest_id, second.ingest_id]
	assert cursor == items[-1].captured_at

	next_items, next_cursor = repo.get_latest(cursor=cursor, limit=2)
	assert [item.ingest_id for item in next_items] == [third.ingest_id]
	assert next_cursor is None


def test_get_detail_with_stats(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 10)
	current = datetime.now(timezone.utc)
	stats = StatsRecord(
		ingest_id=image.ingest_id,
		hall_of_fame_at=current,
		score=42,
		view_count=5,
		last_viewed_at=current,
	)
	session.add(stats)
	session.commit()

	detail = repo.get_context(image.ingest_id)
	assert detail is not None
	assert detail.ingest_id == image.ingest_id

	result = repo.get_detail_with_stats(image.ingest_id)
	assert result is not None
	image_record, stats_record = result
	assert image_record.ingest_id == image.ingest_id
	assert stats_record is not None
	assert stats_record.hall_of_fame_at == current.replace(tzinfo=None)
	assert stats_record.score == 42
	assert stats_record.view_count == 5


def test_upsert_stats_with_increment(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 20)

	stats = repo.upsert_stats_with_increment(image.ingest_id)
	assert stats.view_count == 1
	assert stats.last_viewed_at is not None

	stats_again = repo.upsert_stats_with_increment(image.ingest_id)
	assert stats_again.view_count == 2


def test_update_hall_of_fame(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 30)
	repo.create_stats(image.ingest_id)

	favorite_response = repo.update_favorite(image.ingest_id, True)
	assert favorite_response is not None
	assert favorite_response.favorited_at is not None

	stats = session.get(StatsRecord, image.ingest_id)
	assert stats is not None
	assert stats.hall_of_fame_at == favorite_response.favorited_at


def test_update_score(session: Session) -> None:
	repo = PostgreSQLImageRepository(session)
	image = add_image_record(session, 40)
	repo.create_stats(image.ingest_id)

	score_response = repo.update_score(image.ingest_id, 5)
	assert score_response is not None
	assert score_response.score == DEFAULT_SCORE + 5

	stats = session.get(StatsRecord, image.ingest_id)
	assert stats is not None
	assert stats.score == score_response.score
