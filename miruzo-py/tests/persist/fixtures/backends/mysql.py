import shutil
import socket
import subprocess
import time
import uuid
from collections.abc import Iterator

import MySQLdb as mysql
import pytest
from sqlalchemy.orm import Session

from app.databases.database import MYSQL_CHARSET, MYSQL_COLLATION, _create_mysql_engine
from app.databases.metadata import metadata

MYSQL_IMAGE = 'mysql:9-oracle'
MYSQL_DB = 'miruzo'
MYSQL_USER = 'm'
MYSQL_PASSWORD = 'miruzo1234'


def _find_free_port() -> int:
	with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
		sock.bind(('', 0))
		return sock.getsockname()[1]


@pytest.fixture(scope='session')
def mysql_dsn() -> Iterator[str]:
	if shutil.which('docker') is None:
		pytest.skip('docker binary not available on PATH', allow_module_level=True)

	try:
		host_port = _find_free_port()
	except PermissionError:
		pytest.skip('cannot bind socket to discover free port (insufficient permissions)')

	container_name = f'miruzo-mysql-{uuid.uuid4().hex[:8]}'
	run_cmd = [
		'docker',
		'run',
		'--rm',
		'--name',
		container_name,
		'-e',
		f'MYSQL_ROOT_PASSWORD={MYSQL_PASSWORD}',
		'-e',
		f'MYSQL_DATABASE={MYSQL_DB}',
		'-e',
		f'MYSQL_USER={MYSQL_USER}',
		'-e',
		f'MYSQL_PASSWORD={MYSQL_PASSWORD}',
		'-e',
		'MYSQL_INITDB_SKIP_TZINFO',
		'-p',
		f'{host_port}:3306',
		'-d',
		MYSQL_IMAGE,
		f'--character-set-server={MYSQL_CHARSET}',
		f'--collation-server={MYSQL_COLLATION}',
	]
	subprocess.run(run_cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

	database_config = {
		'host': '127.0.0.1',
		'user': MYSQL_USER,
		'password': MYSQL_PASSWORD,
		'database': MYSQL_DB,
		'port': host_port,
		'connect_timeout': 1,
	}
	for _ in range(30):
		try:
			with mysql.connect(**database_config):
				break
		except mysql.Error:
			time.sleep(1)
	else:
		subprocess.run(['docker', 'rm', '-f', container_name], check=False)
		raise RuntimeError('MySQL container did not become ready in time')

	dsn = f'mysql+mysqldb://{MYSQL_USER}:{MYSQL_PASSWORD}@127.0.0.1:{host_port}/{MYSQL_DB}'
	try:
		yield dsn
	finally:
		subprocess.run(['docker', 'rm', '-f', container_name], check=False)


@pytest.fixture()
def mysql_session(request: pytest.FixtureRequest) -> Iterator[Session]:
	dsn = request.getfixturevalue('mysql_dsn')
	engine = _create_mysql_engine(dsn, pool_size=1, max_overflow=2)
	metadata.create_all(engine)
	with Session(engine) as session:
		yield session
	metadata.drop_all(engine)
