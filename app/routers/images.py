from datetime import datetime, timezone
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
from app.services.activities.love import LoveRunner
from app.services.activities.love_cancel import LoveCancelRunner
from app.services.activities.stats.repository.factory import create_stats_repository
from app.services.images.query_service import ImageQueryService
from app.services.images.repository import ImageRepository
from app.services.views.context import ContextService
from app.utils.http.reponse_builder import build_response


def _get_image_query_service(
	session: Annotated[Session, Depends(get_session)],
) -> ImageQueryService:
	return ImageQueryService(
		session=session,
		repository=ImageRepository(session),
		variant_layers=env.variant_layers,
	)


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


def _get_love_runner() -> LoveRunner:
	return LoveRunner(
		base_timezone=env.base_timezone,
		daily_reset_at=env.time.daily_reset_at,
		daily_love_limit=env.quota.daily_love_limit,
		score_config=env.score,
	)


def _get_love_cancel_runner() -> LoveCancelRunner:
	return LoveCancelRunner(
		base_timezone=env.base_timezone,
		daily_reset_at=env.time.daily_reset_at,
		score_config=env.score,
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


@router.get('/chronological', response_model=ImageListResponse)
def get_chronological(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageQueryService, Depends(_get_image_query_service)],
) -> Response:
	response = service.get_chronological(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/recently', response_model=ImageListResponse)
def get_recently(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageQueryService, Depends(_get_image_query_service)],
) -> Response:
	response = service.get_recently(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/first_love', response_model=ImageListResponse)
def get_first_love(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageQueryService, Depends(_get_image_query_service)],
) -> Response:
	response = service.get_first_love(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/hall_of_fame', response_model=ImageListResponse)
def get_hall_of_fame(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageQueryService, Depends(_get_image_query_service)],
) -> Response:
	response = service.get_hall_of_fame(
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


@router.post('/{ingest_id}/love')
def post_love(
	ingest_id: int,
	session: Annotated[Session, Depends(get_session)],
	runner: Annotated[LoveRunner, Depends(_get_love_runner)],
) -> None:
	current = datetime.now(timezone.utc)

	with session.begin():
		runner.run(session, ingest_id=ingest_id, evaluated_at=current)

	return None


@router.post('/{ingest_id}/love/cancel')
def post_love_cancel(
	ingest_id: int,
	session: Annotated[Session, Depends(get_session)],
	runner: Annotated[LoveCancelRunner, Depends(_get_love_cancel_runner)],
) -> None:
	current = datetime.now(timezone.utc)

	with session.begin():
		runner.run(session, ingest_id=ingest_id, evaluated_at=current)

	return None
