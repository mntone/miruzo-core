import docker
import pytest
from docker.errors import DockerException
from requests.exceptions import ConnectionError as RequestsConnectionError


def _ensure_runtime_api_available() -> None:
	try:
		client = docker.from_env()
		client.ping()  # Verify socket / HTTP API connectivity
	except (DockerException, RequestsConnectionError, OSError, PermissionError) as exc:
		pytest.skip(f'container runtime API unavailable: {exc}')
