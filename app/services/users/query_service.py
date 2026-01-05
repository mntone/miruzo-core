from datetime import datetime, time, timezone
from typing import final
from zoneinfo import ZoneInfo

from app.domain.activities.daily_period import resolve_daily_period_start
from app.models.api.quota import QuotaItem, QuotaResponse
from app.services.users.repository.protocol import UserRepository


@final
class UserQueryService:
	def __init__(
		self,
		repository: UserRepository,
		*,
		daily_love_limit: int,
		daily_reset_at: time,
		base_timezone: ZoneInfo | None,
	) -> None:
		self._repository = repository
		self._daily_love_limit = daily_love_limit
		self._daily_reset_at = daily_reset_at
		self._base_timezone = base_timezone

	def get_quota(self) -> QuotaResponse:
		current = datetime.now(timezone.utc)

		user_record = self._repository.get_or_create_singleton()

		daily_love_limit = self._daily_love_limit
		daily_love_remaining = max(0, daily_love_limit - user_record.daily_love_used)
		daily_love_reset_at = resolve_daily_period_start(
			current,
			daily_reset_at=self._daily_reset_at,
			base_timezone=self._base_timezone,
		)

		response = QuotaResponse(
			love=QuotaItem(
				limit=daily_love_limit,
				remaining=daily_love_remaining,
				period='daily',
				reset_at=daily_love_reset_at,
			),
		)

		return response
