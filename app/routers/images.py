from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import Session
from starlette.responses import Response

from app.database import get_session
from app.models.api.images.patches import FavoriteRequest, FavoriteResponse, ScoreRequest, ScoreResponse
from app.models.api.images.query import ListQuery
from app.models.api.images.responses import ContextResponse, ImageListResponse
from app.services.images.repository.base import ImageRepository
from app.services.images.repository.factory import create_image_repository
from app.services.images.service import ImageService
from app.utils.http.reponse_builder import build_response


def get_repository(session: Annotated[Session, Depends(get_session)]) -> ImageRepository:
	return create_image_repository(session)


def get_service(repository: Annotated[ImageRepository, Depends(get_repository)]) -> ImageService:
	return ImageService(repository)


router = APIRouter(prefix='/i')


@router.get('/latest', response_model=ImageListResponse)
def get_latest(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageService, Depends(get_service)],
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
	service: Annotated[ImageService, Depends(get_service)],
) -> Response:
	response = service.get_context(ingest_id)

	if response is None:
		raise HTTPException(404)

	return build_response(response)


@router.patch('/{ingest_id}/favorite', response_model=FavoriteResponse)
def patch_favorite(
	ingest_id: int,
	payload: FavoriteRequest,
	repo: Annotated[ImageRepository, Depends(get_repository)],
) -> Response:
	response = repo.update_favorite(ingest_id, payload.value)

	if response is None:
		raise HTTPException(404)

	return build_response(response)


@router.patch('/{ingest_id}/score', response_model=ScoreResponse)
def patch_score(
	ingest_id: int,
	payload: ScoreRequest,
	repo: Annotated[ImageRepository, Depends(get_repository)],
) -> Response:
	response = repo.update_score(ingest_id, payload.delta)

	if response is None:
		raise HTTPException(404)

	return build_response(response)
