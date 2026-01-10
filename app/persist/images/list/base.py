# pyright: reportArgumentType=false
# pyright: reportAttributeAccessIssue=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from collections.abc import Sequence
from datetime import datetime
from typing import cast, final

from sqlmodel import Session, select

from app.config.constants import DEFAULT_LIMIT
from app.models.records import ImageRecord, IngestRecord, StatsRecord


@final
class BaseImageListRepository:
	def __init__(self, session: Session, *, engaged_score_threshold: int) -> None:
		self._session = session
		self._engaged_score_threshold = engaged_score_threshold

	def select_latest(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[ImageRecord]:
		statement = (
			# SELECT * FROM images
			select(ImageRecord)
			# ORDER BY ingested_at DESC, ingest_id DESC
			.order_by(
				ImageRecord.ingested_at.desc(),
				ImageRecord.ingest_id.desc(),
			)
			# LIMIT ?
			.limit(limit)
		)

		if cursor is not None:
			# WHERE ingested_at < ?
			statement = statement.where(
				ImageRecord.ingested_at < cursor,
			)

		rows = self._session.exec(statement).all()

		return rows

	def select_chronological(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		statement = (
			# SELECT images.*, ingests.captured_at FROM images
			select(ImageRecord, IngestRecord.captured_at)
			# JOIN ingests ON ingests.id = images.ingest_id
			.join(
				IngestRecord,
				IngestRecord.id == ImageRecord.ingest_id,
			)
			# ORDER BY ingests.captured_at DESC, images.ingest_id DESC
			.order_by(
				IngestRecord.captured_at.desc(),
				ImageRecord.ingest_id.desc(),
			)
			# LIMIT ?
			.limit(limit)
		)

		if cursor is not None:
			# WHERE ingests.captured_at < ?
			statement = statement.where(
				IngestRecord.captured_at < cursor,
			)

		rows = self._session.exec(statement).all()

		return rows

	def select_recently(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		statement = (
			# SELECT images.*, stats.last_viewed_at FROM images
			select(ImageRecord, StatsRecord.last_viewed_at)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.last_viewed_at IS NOT NULL
			.where(StatsRecord.last_viewed_at.is_not(None))
			# ORDER BY stats.last_viewed_at DESC, images.ingest_id DESC
			.order_by(
				StatsRecord.last_viewed_at.desc(),
				ImageRecord.ingest_id.desc(),
			)
			# LIMIT ?
			.limit(limit)
		)

		if cursor is not None:
			# WHERE stats.last_viewed_at < ?
			statement = statement.where(
				StatsRecord.last_viewed_at < cursor,
			)

		rows = self._session.exec(statement).all()

		return cast(Sequence[tuple[ImageRecord, datetime]], rows)

	def select_first_love(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		statement = (
			# SELECT images.*, stats.first_loved_at FROM images
			select(ImageRecord, StatsRecord.first_loved_at)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.first_loved_at IS NOT NULL
			.where(StatsRecord.first_loved_at.is_not(None))
			# ORDER BY stats.first_loved_at DESC, images.ingest_id DESC
			.order_by(
				StatsRecord.first_loved_at.desc(),
				ImageRecord.ingest_id.desc(),
			)
			# LIMIT ?
			.limit(limit)
		)

		if cursor is not None:
			# WHERE stats.first_loved_at < ?
			statement = statement.where(
				StatsRecord.first_loved_at < cursor,
			)

		rows = self._session.exec(statement).all()

		return cast(Sequence[tuple[ImageRecord, datetime]], rows)

	def select_hall_of_fame(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		statement = (
			# SELECT images.*, stats.hall_of_fame_at FROM images
			select(ImageRecord, StatsRecord.hall_of_fame_at)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.hall_of_fame_at IS NOT NULL
			.where(StatsRecord.hall_of_fame_at.is_not(None))
			# ORDER BY stats.hall_of_fame_at DESC, images.ingest_id DESC
			.order_by(
				StatsRecord.hall_of_fame_at.desc(),
				ImageRecord.ingest_id.desc(),
			)
			# LIMIT ?
			.limit(limit)
		)

		if cursor is not None:
			# WHERE stats.hall_of_fame_at < ?
			statement = statement.where(
				StatsRecord.hall_of_fame_at < cursor,
			)

		rows = self._session.exec(statement).all()

		return cast(Sequence[tuple[ImageRecord, datetime]], rows)

	def select_engaged(
		self,
		*,
		cursor: int | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, int]]:
		statement = (
			# SELECT images.*, stats.score_evaluated FROM images
			select(ImageRecord, StatsRecord.score_evaluated)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.hall_of_fame_at IS NULL AND stats.score_evaluated >= THRESHOLD
			.where(
				StatsRecord.hall_of_fame_at.is_(None),
				StatsRecord.score_evaluated >= self._engaged_score_threshold,
			)
			# ORDER BY stats.score_evaluated DESC, images.ingest_id DESC
			.order_by(
				StatsRecord.score_evaluated.desc(),
				ImageRecord.ingest_id.desc(),
			)
			# LIMIT ?
			.limit(limit)
		)

		if cursor is not None:
			# WHERE stats.score_evaluated < ?
			statement = statement.where(
				StatsRecord.score_evaluated < cursor,
			)

		rows = self._session.exec(statement).all()

		return rows
