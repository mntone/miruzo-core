from datetime import timedelta

from tests.stubs.session import StubSession
from tests.stubs.settings import StubSettingsRepository

from app.services.settings import factory as settings_factory


class StubTimezoneProvider:
	def __init__(self, repository: StubSettingsRepository) -> None:
		self._repository = repository
		self._location = 'UTC'

	@property
	def location(self) -> str:
		return self._location

	def ensure_settings(self, initial_location: str | None) -> bool:
		location = self._repository.get('timezone')
		if location is None:
			if initial_location is not None:
				self._location = initial_location
			self._repository.insert('timezone', self._location)
			return True

		self._location = location
		return False


def test_build_daily_period_resolver_uses_offset_without_repo(monkeypatch) -> None:
	def _unexpected_repo(_session: object) -> object:
		raise AssertionError('repository should not be created')

	monkeypatch.setattr(settings_factory, 'create_settings_repository', _unexpected_repo)

	resolver = settings_factory.build_daily_period_resolver(
		initial_location=None,
		day_start_offset=timedelta(hours=3),
	)

	assert resolver.day_start_offset == timedelta(hours=3)


def test_build_daily_period_resolver_initializes_location(monkeypatch) -> None:
	repository = StubSettingsRepository()

	monkeypatch.setattr(settings_factory, 'create_settings_repository', lambda _session: repository)
	monkeypatch.setattr(settings_factory, 'TimezoneProvider', StubTimezoneProvider)

	resolver = settings_factory.build_daily_period_resolver(
		initial_location='UTC',
		day_start_offset=None,
		session=StubSession(),  # pyright: ignore[reportArgumentType]
	)

	assert resolver.day_start_offset == timedelta(hours=5)
	assert repository.inserts == [('timezone', 'UTC')]
