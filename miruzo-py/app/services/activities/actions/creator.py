from datetime import datetime
from typing import final

from app.models.enums import ActionKind
from app.models.records import ActionRecord
from app.persist.actions.protocol import ActionRepository


@final
class ActionCreator:
	def __init__(self, repository: ActionRepository) -> None:
		self._repository = repository

	def view(self, ingest_id: int, *, occurred_at: datetime) -> ActionRecord:
		action = self._repository.insert(
			ingest_id,
			kind=ActionKind.VIEW,
			occurred_at=occurred_at,
		)

		return action

	def love(self, ingest_id: int, *, occurred_at: datetime) -> ActionRecord:
		action = self._repository.insert(
			ingest_id,
			kind=ActionKind.LOVE,
			occurred_at=occurred_at,
		)

		return action

	def cancel_love(self, ingest_id: int, *, occurred_at: datetime) -> ActionRecord:
		action = self._repository.insert(
			ingest_id,
			kind=ActionKind.LOVE_CANCELED,
			occurred_at=occurred_at,
		)

		return action

	def hall_of_fame_added(self, ingest_id: int, *, occurred_at: datetime) -> None:
		self._repository.insert(
			ingest_id,
			kind=ActionKind.HALL_OF_FAME_ADDED,
			occurred_at=occurred_at,
		)

	def hall_of_fame_removed(self, ingest_id: int, *, occurred_at: datetime) -> None:
		self._repository.insert(
			ingest_id,
			kind=ActionKind.HALL_OF_FAME_REMOVED,
			occurred_at=occurred_at,
		)
