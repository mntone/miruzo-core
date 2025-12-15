from abc import ABC, abstractmethod
from datetime import datetime
from typing import Callable

from sqlalchemy import func
from sqlmodel import Session, select

from app.core.constants import DEFAULT_SCORE
from app.models.api.images.patches import FavoriteResponse, ScoreResponse
from app.models.records import ImageRecord, StatsRecord


class ImageRepository(ABC):
	def __init__(self, session: Session) -> None:
		self._session = session

	def get_list(
		self,
		*,
		cursor: datetime | None,
		limit: int,
	) -> tuple[list[ImageRecord], datetime | None]:
		statement = select(ImageRecord)

		if cursor is not None:
			statement = statement.where(ImageRecord.captured_at < cursor)

		statement = statement.order_by(ImageRecord.captured_at.desc()).limit(limit)
		items = self._session.exec(statement).all()

		next_cursor: datetime | None = None
		if len(items) == limit:
			next_cursor = items[-1].captured_at

		return items, next_cursor

	def get_detail(
		self,
		image_id: int,
	) -> ImageRecord | None:
		image = self._session.get(ImageRecord, image_id)

		return image

	@abstractmethod
	def get_detail_with_stats(
		self,
		image_id: int,
	) -> tuple[ImageRecord, StatsRecord] | None:
		"""Fetch both the image record and stats in a single operation."""

	def create_stats(self, image_id: int) -> StatsRecord:
		stats = StatsRecord(
			image_id=image_id,
			favorite=False,
			score=DEFAULT_SCORE,
			view_count=0,
			last_viewed_at=datetime.now(),
		)
		self._session.add(stats)
		self._session.commit()
		self._session.refresh(stats)

		return stats

	def _upsert_stats_with_increment(
		self,
		insert: Callable,
		image_id: int,
	) -> StatsRecord:
		statement = (
			insert(StatsRecord)
			# INSERT INTO imagestats (image_id, view_count) VALUES (:image_id, 1)
			.values(
				image_id=image_id,
				view_count=1,
				last_viewed_at=func.current_timestamp(),
			)
			# ON CONFLICT(image_id)
			# DO UPDATE SET view_count = view_count + 1
			.on_conflict_do_update(
				index_elements=['image_id'],
				set_={
					'view_count': StatsRecord.view_count + 1,
					'last_viewed_at': func.current_timestamp(),
				},
			)
			# RETURNING image_id, favorite, score, view_count, last_viewed_at
			.returning(
				StatsRecord.image_id,
				StatsRecord.favorite,
				StatsRecord.score,
				StatsRecord.view_count,
				StatsRecord.last_viewed_at,
			)
		)

		row = self._session.exec(statement).first()
		self._session.commit()

		return StatsRecord(**row._asdict())

	@abstractmethod
	def upsert_stats_with_increment(self, image_id: int) -> StatsRecord:
		"""Increment view stats (upserting as needed) and return the latest row."""

	def update_favorite(
		self,
		image_id: int,
		favorite: bool,
	) -> FavoriteResponse | None:
		stats = self._session.get(StatsRecord, image_id)

		if not stats:
			return None

		stats.favorite = favorite
		self._session.add(stats)
		self._session.commit()

		return FavoriteResponse.from_record(stats)

	def update_score(
		self,
		image_id: int,
		delta_score: int,
	) -> ScoreResponse | None:
		stats = self._session.get(StatsRecord, image_id)

		if not stats:
			return None

		stats.score += delta_score
		self._session.add(stats)
		self._session.commit()

		return ScoreResponse.from_record(stats)
