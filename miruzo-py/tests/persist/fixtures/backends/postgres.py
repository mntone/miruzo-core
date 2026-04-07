import shutil
import socket
import subprocess
import time
import uuid
from collections.abc import Iterator

import psycopg2
import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from app.databases.metadata import metadata

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
def postgres_session(request: pytest.FixtureRequest) -> Iterator[Session]:
	dsn = request.getfixturevalue('postgres_dsn')
	engine = create_engine(dsn, echo=False)
	metadata.create_all(engine)
	with Session(engine) as postgres_session:
		yield postgres_session
	metadata.drop_all(engine)
