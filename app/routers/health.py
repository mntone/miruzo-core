from fastapi import APIRouter

from app.models.api.health import HealthResponse
from app.version import __version__

router = APIRouter()


@router.get('/health')
async def get_health() -> HealthResponse:
	"""Simple endpoint that can be used by monitors or tests."""
	return {
		'version': __version__,
	}
