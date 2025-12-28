from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import Session
from starlette.responses import Response

from app.database import get_session
from app.models.api.images.query import ListQuery
from app.models.api.images.responses import ContextResponse, ImageListResponse
from app.services.activities.stats.repository.base import BaseStatsRepository
from app.services.activities.stats.repository.factory import create_stats_repository
from app.services.activities.stats.repository.protocol import StatsRepository
from app.services.activities.stats.service import StatsService
from app.services.images.query import ImageQueryService
from app.services.images.repository import ImageRepository
from app.services.views.context import ContextService
from app.utils.http.reponse_builder import build_response


def _get_image_repository(session: Annotated[Session, Depends(get_session)]) -> ImageRepository:
	return ImageRepository(session)


def _get_stats_repository(session: Annotated[Session, Depends(get_session)]) -> StatsRepository:
	return create_stats_repository(session)


def _get_image_query_service(
	repository: Annotated[ImageRepository, Depends(_get_image_repository)],
) -> ImageQueryService:
	return ImageQueryService(repository)


def _get_stats_service(
	repository: Annotated[BaseStatsRepository, Depends(_get_stats_repository)],
) -> StatsService:
	return StatsService(repository)


def _get_context_service(
	image_query: Annotated[ImageQueryService, Depends(_get_image_query_service)],
	stats: Annotated[StatsService, Depends(_get_stats_service)],
) -> ContextService:
	return ContextService(
		image_query=image_query,
		stats=stats,
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
