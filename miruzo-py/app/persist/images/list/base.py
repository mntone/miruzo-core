# pyright: reportArgumentType=false
# pyright: reportAttributeAccessIssue=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from collections.abc import Sequence
from datetime import datetime
from typing import cast, final

from sqlmodel import Session, and_, or_, select

from app.config.constants import DEFAULT_LIMIT
from app.domain.images.cursor import DatetimeImageListCursor, UInt8ImageListCursor
from app.models.records import ImageRecord, IngestRecord, StatsRecord


@final
class BaseImageListRepository:
	def __init__(self, session: Session, *, engaged_score_threshold: int) -> None:
		self._session = session
		self._engaged_score_threshold = engaged_score_threshold

	def select_latest(
		self,
		*,
		cursor: DatetimeImageListCursor | None = None,
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
			# WHERE (ingested_at < ?) OR (ingested_at == ? AND ingest_id < ?)
			statement = statement.where(
				or_(
					ImageRecord.ingested_at < cursor.value,
					and_(
						ImageRecord.ingested_at == cursor.value,
						ImageRecord.ingest_id < cursor.ingest_id,
					),
				),
			)

		rows = self._session.exec(statement).all()

		return rows

	def select_chronological(
		self,
		*,
		cursor: DatetimeImageListCursor | None = None,
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
			# WHERE (ingests.captured_at < ?) OR (ingests.captured_at == ? AND images.ingest_id < ?)
			statement = statement.where(
				or_(
					IngestRecord.captured_at < cursor.value,
					and_(
						IngestRecord.captured_at == cursor.value,
						ImageRecord.ingest_id < cursor.ingest_id,
					),
				),
			)

		rows = self._session.exec(statement).all()

		return rows

	def select_recently(
		self,
		*,
		cursor: DatetimeImageListCursor | None = None,
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
			# WHERE (stats.last_viewed_at < ?) OR (stats.last_viewed_at == ? AND images.ingest_id < ?)
			statement = statement.where(
				or_(
					StatsRecord.last_viewed_at < cursor.value,
					and_(
						StatsRecord.last_viewed_at == cursor.value,
						ImageRecord.ingest_id < cursor.ingest_id,
					),
				),
			)

		rows = self._session.exec(statement).all()

		return cast(Sequence[tuple[ImageRecord, datetime]], rows)

	def select_first_love(
		self,
		*,
		cursor: DatetimeImageListCursor | None = None,
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
			# WHERE (stats.first_loved_at < ?) OR (stats.first_loved_at == ? AND images.ingest_id < ?)
			statement = statement.where(
				or_(
					StatsRecord.first_loved_at < cursor.value,
					and_(
						StatsRecord.first_loved_at == cursor.value,
						ImageRecord.ingest_id < cursor.ingest_id,
					),
				),
			)

		rows = self._session.exec(statement).all()

		return cast(Sequence[tuple[ImageRecord, datetime]], rows)

	def select_hall_of_fame(
		self,
		*,
		cursor: DatetimeImageListCursor | None = None,
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
			# WHERE (stats.hall_of_fame_at < ?) OR (stats.hall_of_fame_at == ? AND images.ingest_id < ?)
			statement = statement.where(
				or_(
					StatsRecord.hall_of_fame_at < cursor.value,
					and_(
						StatsRecord.hall_of_fame_at == cursor.value,
						ImageRecord.ingest_id < cursor.ingest_id,
					),
				),
			)

		rows = self._session.exec(statement).all()

		return cast(Sequence[tuple[ImageRecord, datetime]], rows)

	def select_engaged(
		self,
		*,
		cursor: UInt8ImageListCursor | None = None,
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
			# WHERE (stats.score_evaluated < ?) OR (stats.score_evaluated == ? AND images.ingest_id < ?)
			statement = statement.where(
				or_(
					StatsRecord.score_evaluated < cursor.value,
					and_(
						StatsRecord.score_evaluated == cursor.value,
						ImageRecord.ingest_id < cursor.ingest_id,
					),
				),
			)

		rows = self._session.exec(statement).all()

		return rows
