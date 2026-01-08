# pyright: reportArgumentType=false
# pyright: reportAttributeAccessIssue=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

from collections.abc import Sequence
from datetime import datetime
from typing import Any

from sqlalchemy.engine import Row
from sqlmodel import Session, select

from app.models.records import ImageRecord, IngestRecord, StatsRecord


class ImageListQueryExecutor:
	"""Build and execute list-query statements for images."""

	def __init__(self, session: Session, *, engaged_score_threshold: int) -> None:
		"""Store the session used to execute list queries."""

		self._statement = None
		self._session = session
		self._engaged_score_threshold = engaged_score_threshold

	def latest(self, *, cursor: datetime | None) -> 'ImageListQueryExecutor':
		"""Prepare a latest-first images query."""

		self._statement = (
			# SELECT * FROM images
			select(ImageRecord)
			# ORDER BY images.ingested_at DESC
			.order_by(ImageRecord.ingested_at.desc())
		)

		if cursor is not None:
			# WHERE ingested_at < ?
			self._statement = self._statement.where(
				ImageRecord.ingested_at < cursor,
			)

		return self

	def chronological(self, *, cursor: datetime | None) -> 'ImageListQueryExecutor':
		"""Prepare a timeline images query."""

		self._statement = (
			# SELECT images.*, ingests.captured_at FROM images
			select(ImageRecord, IngestRecord.captured_at)
			# JOIN ingests ON ingests.id = images.ingest_id
			.join(
				IngestRecord,
				IngestRecord.id == ImageRecord.ingest_id,
			)
			# ORDER BY ingests.captured_at DESC
			.order_by(IngestRecord.captured_at.desc())
		)

		if cursor is not None:
			# WHERE ingests.captured_at < ?
			self._statement = self._statement.where(
				IngestRecord.captured_at < cursor,
			)

		return self

	def recently(self, *, cursor: datetime | None) -> 'ImageListQueryExecutor':
		"""Prepare a recently-viewed images query."""

		self._statement = (
			# SELECT images.*, stats.last_viewed_at FROM images
			select(ImageRecord, StatsRecord.last_viewed_at)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.last_viewed_at IS NOT NULL
			.where(StatsRecord.last_viewed_at.is_not(None))
			# ORDER BY stats.last_viewed_at DESC
			.order_by(StatsRecord.last_viewed_at.desc())
		)

		if cursor is not None:
			# WHERE stats.last_viewed_at < ?
			self._statement = self._statement.where(
				StatsRecord.last_viewed_at < cursor,
			)

		return self

	def first_love(self, *, cursor: datetime | None) -> 'ImageListQueryExecutor':
		"""Prepare a first-loved images query."""

		self._statement = (
			# SELECT images.*, stats.first_loved_at FROM images
			select(ImageRecord, StatsRecord.first_loved_at)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.first_loved_at IS NOT NULL
			.where(StatsRecord.first_loved_at.is_not(None))
			# ORDER BY stats.first_loved_at DESC
			.order_by(StatsRecord.first_loved_at.desc())
		)

		if cursor is not None:
			# WHERE stats.first_loved_at < ?
			self._statement = self._statement.where(
				StatsRecord.first_loved_at < cursor,
			)

		return self

	def hall_of_fame(self, *, cursor: datetime | None) -> 'ImageListQueryExecutor':
		"""Prepare a hall-of-fame images query."""

		self._statement = (
			# SELECT images.*, stats.hall_of_fame_at FROM images
			select(ImageRecord, StatsRecord.hall_of_fame_at)
			# JOIN stats ON stats.ingest_id = images.ingest_id
			.join(
				StatsRecord,
				StatsRecord.ingest_id == ImageRecord.ingest_id,
			)
			# WHERE stats.hall_of_fame_at IS NOT NULL
			.where(StatsRecord.hall_of_fame_at.is_not(None))
			# ORDER BY stats.hall_of_fame_at DESC
			.order_by(StatsRecord.hall_of_fame_at.desc())
		)

		if cursor is not None:
			# WHERE stats.hall_of_fame_at < ?
			self._statement = self._statement.where(
				StatsRecord.hall_of_fame_at < cursor,
			)

		return self

	def engaged(self, *, cursor: int | None) -> 'ImageListQueryExecutor':
		"""Prepare a engaged images query."""

		self._statement = (
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
			# ORDER BY stats.score_evaluated DESC
			.order_by(StatsRecord.score_evaluated.desc())
		)

		if cursor is not None:
			# WHERE stats.score_evaluated < ?
			self._statement = self._statement.where(
				StatsRecord.score_evaluated < cursor,
			)

		return self

	def order_by_ingest_id(self) -> 'ImageListQueryExecutor':
		"""Append ingest_id DESC ordering as a stable tie-breaker."""

		# ORDER BY images.ingest_id DESC
		self._statement = self._statement.order_by(ImageRecord.ingest_id.desc())

		return self

	def limit(self, n: int) -> 'ImageListQueryExecutor':
		"""Apply a row limit to the current statement."""

		self._statement = self._statement.limit(n)

		return self

	def execute(self) -> Sequence[Row[Any]]:
		"""Execute the current statement and return rows."""

		rows = self._session.exec(self._statement).all()  # pyright: ignore[reportCallIssue]

		return rows
