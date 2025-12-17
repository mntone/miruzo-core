from fastapi import APIRouter
from starlette.responses import Response

from app.models.api.health import HealthResponse
from app.utils.http.reponse_builder import build_response
from app.version import __version__

router = APIRouter()


@router.get('/health', response_model=HealthResponse)
async def get_health() -> Response:
	"""Simple endpoint that can be used by monitors or tests."""

	response = HealthResponse(
		version=__version__,
	)

	return build_response(response)
