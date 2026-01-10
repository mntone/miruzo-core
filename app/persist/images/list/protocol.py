from collections.abc import Sequence
from datetime import datetime
from typing import Protocol

from app.config.constants import DEFAULT_LIMIT
from app.models.records import ImageRecord


class ImageListRepository(Protocol):
	"""Define list-query entry points and cursor types for image lists."""

	def select_latest(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[ImageRecord]:
		"""Return a paginated list of latest images."""
		...

	def select_chronological(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		"""Return a paginated list of timeline images."""
		...

	def select_recently(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		"""Return a paginated list of recently viewed images."""
		...

	def select_first_love(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		"""Return a paginated list of first-loved images."""
		...

	def select_hall_of_fame(
		self,
		*,
		cursor: datetime | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, datetime]]:
		"""Return a paginated list of hall-of-fame images."""
		...

	def select_engaged(
		self,
		*,
		cursor: int | None = None,
		limit: int = DEFAULT_LIMIT,
	) -> Sequence[tuple[ImageRecord, int]]:
		"""Return a paginated list of engaged images."""
		...
