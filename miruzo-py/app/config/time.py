from dataclasses import dataclass
from datetime import time

from typing_extensions import final


@dataclass(frozen=True, slots=True)
@final
class TimeConfig:
	daily_reset_at: time = time(5, 0)
