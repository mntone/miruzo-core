from datetime import datetime, time
from typing import final
from zoneinfo import ZoneInfo

from app.domain.activities.daily_period import resolve_daily_period_range
from app.models.enums import ActionKind
from app.models.records import ActionRecord
from app.services.activities.actions.repository import ActionRepository


@final
class DecayActionCreator:
	"""Create a daily decay action unless one already exists for the period."""

	def __init__(
		self,
		repository: ActionRepository,
		*,
		daily_reset_at: time,
		base_timezone: ZoneInfo | None,
	) -> None:
		self._repository = repository
		self._daily_reset_at = daily_reset_at
		self._base_timezone = base_timezone

	def create(
		self,
		ingest_id: int,
		*,
		occurred_at: datetime,
	) -> ActionRecord | None:
		since, until = resolve_daily_period_range(
			occurred_at,
			daily_reset_at=self._daily_reset_at,
			base_timezone=self._base_timezone,
		)

		action = self._repository.select_one_by(
			ingest_id,
			kind=ActionKind.DECAY,
			since_occurred_at=since,
			until_occurred_at=until,
		)
		if action is not None:
			return None

		action = self._repository.insert(
			ingest_id,
			kind=ActionKind.DECAY,
			occurred_at=occurred_at,
		)

		return action
