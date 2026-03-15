from tzlocal import get_localzone

from app.persist.settings.protocol import SettingsRepository

_TIMEZONE_SETTINGS_KEY = 'timezone'


def _get_default_location() -> str:
	localzone = get_localzone()
	return localzone.key


class TimezoneProvider:
	def __init__(self, repository: SettingsRepository) -> None:
		self._repository = repository
		self._location = _get_default_location()

	@property
	def location(self) -> str:
		return self._location

	def set_location(self, value: str) -> None:
		self._location = value

	def ensure_settings(self, initial_location: str | None) -> bool:
		location = self._repository.get(_TIMEZONE_SETTINGS_KEY)
		if location is None:
			if initial_location is not None:
				self._location = initial_location
			self._repository.insert(_TIMEZONE_SETTINGS_KEY, self._location)
			return True

		self._location = location
		return False
