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

from app.models.records import ImageRecord, StatsRecord


class ImageListQueryExecutor:
	"""Build and execute list-query statements for images."""

	def __init__(self, session: Session) -> None:
		"""Store the session used to execute list queries."""

		self._statement = None
		self._session = session

	def latest(self, *, cursor: datetime | None) -> 'ImageListQueryExecutor':
		"""Prepare a latest-first images query."""

		# SELECT * FROM images
		self._statement = select(ImageRecord)

		if cursor is not None:
			# WHERE captured_at < ?
			self._statement = self._statement.where(
				ImageRecord.captured_at < cursor,
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
			# ORDER BY stats.hall_of_fame_at DESC
			.order_by(StatsRecord.last_viewed_at.desc())
		)

		if cursor is not None:
			# WHERE stats.last_viewed_at < ?
			self._statement = self._statement.where(
				StatsRecord.last_viewed_at < cursor,
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

	def order_by_latest(self) -> 'ImageListQueryExecutor':
		"""Order by captured time and ingest id, newest first."""

		# ORDER BY images.captured_at DESC, images.ingest_id DESC
		self._statement = self._statement.order_by(
			ImageRecord.captured_at.desc(),
			ImageRecord.ingest_id.desc(),
		)

		return self

	def limit(self, n: int) -> 'ImageListQueryExecutor':
		"""Apply a row limit to the current statement."""

		self._statement = self._statement.limit(n)

		return self

	def execute(self) -> Sequence[Row[Any]]:
		"""Execute the current statement and return rows."""

		rows = self._session.exec(self._statement).all()  # pyright: ignore[reportCallIssue]

		return rows
