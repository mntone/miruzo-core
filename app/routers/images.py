from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import Session
from starlette.responses import Response

from app.config.environments import env
from app.database import get_session
from app.models.api.context.responses import ContextResponse
from app.models.api.images.query import ListQuery
from app.models.api.images.responses import ImageListResponse
from app.services.activities.actions.repository import ActionRepository
from app.services.activities.stats.repository.factory import create_stats_repository
from app.services.images.query import ImageQueryService
from app.services.images.repository import ImageRepository
from app.services.views.context import ContextService
from app.utils.http.reponse_builder import build_response


def _get_image_repository(session: Annotated[Session, Depends(get_session)]) -> ImageRepository:
	return ImageRepository(session)


def _get_image_query_service(
	repository: Annotated[ImageRepository, Depends(_get_image_repository)],
) -> ImageQueryService:
	return ImageQueryService(repository)


def _get_context_service(
	session: Annotated[Session, Depends(get_session)],
	image_query: Annotated[ImageQueryService, Depends(_get_image_query_service)],
) -> ContextService:
	action_repo = ActionRepository(session)
	return ContextService(
		session,
		action=action_repo,
		image_query=image_query,
		stats=create_stats_repository(session),
		env=env,
	)


router = APIRouter(prefix='/i')


@router.get('/latest', response_model=ImageListResponse)
def get_latest(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageQueryService, Depends(_get_image_query_service)],
) -> Response:
	response = service.get_latest(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/{ingest_id}', response_model=ContextResponse)
def get_context(
	ingest_id: int,
	service: Annotated[ContextService, Depends(_get_context_service)],
) -> Response:
	response = service.get_context(ingest_id)

	if response is None:
		raise HTTPException(404)

	return build_response(response)
