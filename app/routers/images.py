from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import Session

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


@router.get('/latest')
def get_latest(
	query: Annotated[ListQuery, Depends()],
	service: Annotated[ImageService, Depends(get_service)],
) -> ImageListResponse:
	response = service.get_latest(
		cursor=query.cursor,
		limit=query.limit,
		exclude_formats=query.exclude_formats,
	)
	return build_response(response)


@router.get('/{image_id}')
def get_context(
	image_id: int,
	service: Annotated[ImageService, Depends(get_service)],
) -> ContextResponse:
	response = service.get_context(image_id)

	if response is None:
		raise HTTPException(404)

	return build_response(response)


@router.patch('/{image_id}/favorite')
def patch_favorite(
	image_id: int,
	payload: FavoriteRequest,
	repo: Annotated[ImageRepository, Depends(get_repository)],
) -> FavoriteResponse:
	response = repo.update_favorite(image_id, payload.value)

	if response is None:
		raise HTTPException(404)

	return build_response(response)


@router.patch('/{image_id}/score')
def patch_score(
	image_id: int,
	payload: ScoreRequest,
	repo: Annotated[ImageRepository, Depends(get_repository)],
) -> ScoreResponse:
	response = repo.update_score(image_id, payload.delta)

	if response is None:
		raise HTTPException(404)

	return build_response(response)
