# pyright: reportAttributeAccessIssue=false
# pyright: reportOptionalMemberAccess=false
# pyright: reportOptionalOperand=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from collections.abc import Sequence
from datetime import datetime

from sqlmodel import Session, select

from app.models.records import ImageRecord


class ImageRepository:
	def __init__(self, session: Session) -> None:
		self._session = session

	def select_latest(
		self,
		*,
		cursor: datetime | None,
		limit: int,
	) -> tuple[Sequence[ImageRecord], datetime | None]:
		statement = select(ImageRecord)

		if cursor is not None:
			statement = statement.where(ImageRecord.captured_at < cursor)

		statement = statement.order_by(ImageRecord.captured_at.desc()).limit(limit)
		items = self._session.exec(statement).all()

		next_cursor: datetime | None = None
		if len(items) == limit:
			next_cursor = items[-1].captured_at

		return items, next_cursor

	def select_by_ingest_id(
		self,
		ingest_id: int,
	) -> ImageRecord | None:
		image = self._session.get(ImageRecord, ingest_id)

		return image

	def insert(self, image: ImageRecord) -> None:
		self._session.add(image)
		self._session.commit()
		self._session.refresh(image)
