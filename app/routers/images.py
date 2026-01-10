from datetime import datetime, timezone
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import Session
from starlette.responses import Response

from app.config.environments import env
from app.database import get_session
from app.models.api.activities.responses import LoveStatsResponse
from app.models.api.context.query import ContextQuery
from app.models.api.context.responses import ContextResponse
from app.models.api.images.query import ListQuery
from app.models.api.images.responses import ImageListResponse
from app.persist.actions.factory import create_action_repository
from app.persist.images.factory import create_image_repository
from app.persist.images.list.factory import create_image_list_repository
from app.persist.stats.factory import create_stats_repository
from app.services.activities.love import LoveRunner
from app.services.activities.love_cancel import LoveCancelRunner
from app.services.images.list import ImageListService
from app.services.views.context import ContextService
from app.utils.http.reponse_builder import build_response


def _get_image_list_service(
	session: Annotated[Session, Depends(get_session)],
) -> ImageListService:
	return ImageListService(
		repository=create_image_list_repository(
			session,
			engaged_score_threshold=env.score.engaged_score_threshold,
		),
		variant_layers=env.variant_layers,
	)


def _get_context_service(
	session: Annotated[Session, Depends(get_session)],
) -> ContextService:
	return ContextService(
		session,
		action_repo=create_action_repository(session),
		image_repo=create_image_repository(session),
		stats_repo=create_stats_repository(session),
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


@router.get('/latest', response_model=ImageListResponse[datetime])
def get_latest(
	query: Annotated[ListQuery[datetime], Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_latest(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/chronological', response_model=ImageListResponse[datetime])
def get_chronological(
	query: Annotated[ListQuery[datetime], Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_chronological(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/recently', response_model=ImageListResponse[datetime])
def get_recently(
	query: Annotated[ListQuery[datetime], Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_recently(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/first_love', response_model=ImageListResponse[datetime])
def get_first_love(
	query: Annotated[ListQuery[datetime], Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_first_love(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/hall_of_fame', response_model=ImageListResponse[datetime])
def get_hall_of_fame(
	query: Annotated[ListQuery[datetime], Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_hall_of_fame(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/engaged', response_model=ImageListResponse[int])
def get_engaged(
	query: Annotated[ListQuery[int], Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_engaged(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/{ingest_id}', response_model=ContextResponse)
def get_context(
	ingest_id: int,
	query: Annotated[ContextQuery, Depends()],
	service: Annotated[ContextService, Depends(_get_context_service)],
) -> Response:
	response = service.get_context(ingest_id, query=query)

	if response is None:
		raise HTTPException(404)

	return build_response(response)


@router.post('/{ingest_id}/love', response_model=LoveStatsResponse)
def post_love(
	ingest_id: int,
	session: Annotated[Session, Depends(get_session)],
	runner: Annotated[LoveRunner, Depends(_get_love_runner)],
) -> Response:
	current = datetime.now(timezone.utc)

	with session.begin():
		response = runner.run(session, ingest_id=ingest_id, evaluated_at=current)

	return build_response(response, exclude_none=False)


@router.post('/{ingest_id}/love/cancel', response_model=LoveStatsResponse)
def post_love_cancel(
	ingest_id: int,
	session: Annotated[Session, Depends(get_session)],
	runner: Annotated[LoveCancelRunner, Depends(_get_love_cancel_runner)],
) -> Response:
	current = datetime.now(timezone.utc)

	with session.begin():
		response = runner.run(session, ingest_id=ingest_id, evaluated_at=current)

	return build_response(response, exclude_none=False)
