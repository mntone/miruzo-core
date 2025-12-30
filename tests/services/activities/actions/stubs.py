from datetime import datetime

from app.models.enums import ActionKind
from app.models.records import ActionRecord


class StubActionRepository:
	def __init__(self) -> None:
		self.insert_called_with: dict[str, object] | None = None
		self.select_called_with: int | None = None
		self.actions: list[ActionRecord] = []

	def insert(
		self,
		ingest_id: int,
		*,
		kind: ActionKind,
		occurred_at: datetime,
	) -> ActionRecord:
		self.insert_called_with = {
			'ingest_id': ingest_id,
			'kind': kind,
			'occurred_at': occurred_at,
		}
		action = ActionRecord(
			id=1,
			ingest_id=ingest_id,
			kind=kind,
			occurred_at=occurred_at,
		)
		self.actions.append(action)
		return action

	def select_by_ingest_id(self, ingest_id: int) -> list[ActionRecord]:
		self.select_called_with = ingest_id
		return list(self.actions)
