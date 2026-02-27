from datetime import datetime, timezone
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import Session
from starlette.responses import Response

from app.config.environments import env
from app.databases import get_session
from app.domain.activities.daily_period import DailyPeriodResolver
from app.domain.images.cursor import (
	DatetimeImageListCursor,
	ImageListCursorMode,
	UInt8ImageListCursor,
)
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
from app.services.images.cursor_codec import (
	CursorDecodeError,
	decode_datetime_image_list_cursor,
	decode_uint8_image_list_cursor,
	encode_image_list_cursor,
)
from app.services.images.list import ImageListService
from app.services.views.context import ContextService
from app.utils.http.response_builder import build_response


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


def _get_daily_period_resolver() -> DailyPeriodResolver:
	return DailyPeriodResolver(
		base_timezone=env.base_timezone,
		daily_reset_at=env.time.daily_reset_at,
	)


def _get_love_runner(
	resolver: Annotated[DailyPeriodResolver, Depends(_get_daily_period_resolver)],
) -> LoveRunner:
	return LoveRunner(
		period_resolver=resolver,
		daily_love_limit=env.quota.daily_love_limit,
		score_config=env.score,
	)


def _get_love_cancel_runner(
	resolver: Annotated[DailyPeriodResolver, Depends(_get_daily_period_resolver)],
) -> LoveCancelRunner:
	return LoveCancelRunner(
		period_resolver=resolver,
		score_config=env.score,
	)


def _decode_datetime_cursor(
	cursor: str | None,
	*,
	mode: ImageListCursorMode,
) -> DatetimeImageListCursor | None:
	if cursor is None:
		return None

	try:
		return decode_datetime_image_list_cursor(cursor, expected_mode=mode)
	except CursorDecodeError as exception:
		raise HTTPException(400, detail='invalid cursor') from exception


def _decode_uint8_cursor(
	cursor: str | None,
	*,
	mode: ImageListCursorMode,
) -> UInt8ImageListCursor | None:
	if cursor is None:
		return None

	try:
		return decode_uint8_image_list_cursor(cursor, expected_mode=mode)
	except CursorDecodeError as exception:
		raise HTTPException(400, detail='invalid cursor') from exception


def _encode_list_response_cursor(
	response: ImageListResponse[DatetimeImageListCursor] | ImageListResponse[UInt8ImageListCursor],
) -> ImageListResponse:
	return ImageListResponse(
		items=response.items,
		cursor=(encode_image_list_cursor(response.cursor) if response.cursor is not None else None),
	)


router = APIRouter(prefix='/i')


@router.get('/latest', response_model=ImageListResponse[str])
def get_latest(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_latest(
		cursor=_decode_datetime_cursor(query.cursor, mode=ImageListCursorMode.LATEST),
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(_encode_list_response_cursor(response))


@router.get('/chronological', response_model=ImageListResponse[str])
def get_chronological(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_chronological(
		cursor=_decode_datetime_cursor(
			query.cursor,
			mode=ImageListCursorMode.CHRONOLOGICAL,
		),
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(_encode_list_response_cursor(response))


@router.get('/recently', response_model=ImageListResponse[str])
def get_recently(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_recently(
		cursor=_decode_datetime_cursor(query.cursor, mode=ImageListCursorMode.RECENTLY),
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(_encode_list_response_cursor(response))


@router.get('/first_love', response_model=ImageListResponse[str])
def get_first_love(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_first_love(
		cursor=_decode_datetime_cursor(query.cursor, mode=ImageListCursorMode.FIRST_LOVE),
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(_encode_list_response_cursor(response))


@router.get('/hall_of_fame', response_model=ImageListResponse[str])
def get_hall_of_fame(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_hall_of_fame(
		cursor=_decode_datetime_cursor(query.cursor, mode=ImageListCursorMode.HALL_OF_FAME),
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(_encode_list_response_cursor(response))


@router.get('/engaged', response_model=ImageListResponse[str])
def get_engaged(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageListService, Depends(_get_image_list_service)],
) -> Response:
	response = service.get_engaged(
		cursor=_decode_uint8_cursor(query.cursor, mode=ImageListCursorMode.ENGAGED),
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(_encode_list_response_cursor(response))


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
