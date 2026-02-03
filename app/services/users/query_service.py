from datetime import datetime, timezone
from typing import final

from app.domain.activities.daily_period import DailyPeriodResolver
from app.models.api.quota import QuotaItem, QuotaResponse
from app.persist.users.protocol import UserRepository


@final
class UserQueryService:
	def __init__(
		self,
		repository: UserRepository,
		*,
		daily_love_limit: int,
		period_resolver: DailyPeriodResolver,
	) -> None:
		self._repository = repository
		self._daily_love_limit = daily_love_limit
		self._period_resolver = period_resolver

	def get_quota(self) -> QuotaResponse:
		current = datetime.now(timezone.utc)

		user_record = self._repository.get_or_create_singleton()

		daily_love_limit = self._daily_love_limit
		daily_love_remaining = max(0, daily_love_limit - user_record.daily_love_used)
		daily_love_reset_at = self._period_resolver.resolve_period_start(current)

		response = QuotaResponse(
			love=QuotaItem(
				limit=daily_love_limit,
				remaining=daily_love_remaining,
				period='daily',
				reset_at=daily_love_reset_at,
			),
		)

		return response
