from datetime import timedelta

from sqlmodel import Session

from app.databases.database import create_session
from app.domain.activities.daily_period import DailyPeriodResolver
from app.persist.settings.factory import create_settings_repository
from app.services.settings.timezone import TimezoneProvider


def _get_location(session: Session, initial_location: str | None) -> tuple[str, bool]:
	repository = create_settings_repository(session)
	timezone_provider = TimezoneProvider(repository)
	initial = timezone_provider.ensure_settings(initial_location)
	location = timezone_provider.location
	return location, initial


def build_daily_period_resolver(
	*,
	initial_location: str | None,
	day_start_offset: timedelta | None,
	session: Session | None = None,
) -> DailyPeriodResolver:
	if day_start_offset is not None:
		return DailyPeriodResolver(day_start_offset)

	if session is None:
		with create_session() as session:
			location, initial = _get_location(session, initial_location)
			if initial:
				session.commit()
	else:
		location, _ = _get_location(session, initial_location)

	return DailyPeriodResolver.from_location(location)
