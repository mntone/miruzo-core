from typing import Annotated

from fastapi import APIRouter, Depends, Request
from sqlmodel import Session
from starlette.responses import Response

from app.config.environments import env
from app.databases import get_session
from app.models.api.quota import QuotaResponse
from app.persist.users.factory import create_user_repository
from app.services.users.query_service import UserQueryService
from app.utils.http.response_builder import build_response


def _get_user_query_service(
	request: Request,
	session: Annotated[Session, Depends(get_session)],
) -> UserQueryService:
	return UserQueryService(
		create_user_repository(session),
		daily_love_limit=env.quota.daily_love_limit,
		period_resolver=request.app.state['period_resolver'],
	)


router = APIRouter()


@router.get('/quota', response_model=QuotaResponse)
async def get_quota(
	service: Annotated[UserQueryService, Depends(_get_user_query_service)],
) -> Response:
	response = service.get_quota()
	return build_response(response)
