from dataclasses import dataclass
from datetime import timedelta

from typing_extensions import final


@dataclass(frozen=True, slots=True)
@final
class PeriodConfig:
	"""
	Configure daily period behavior.

	Priority:
	- Use day_start_offset when set.
	- Otherwise, use the saved location from the settings table.
	- If missing, use initial_location.
	- If still missing, use the host's local timezone.
	"""

	initial_location: str | None = None

	day_start_offset: timedelta | None = None
