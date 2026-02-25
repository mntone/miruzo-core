from dataclasses import dataclass

from typing_extensions import final


@dataclass(frozen=True, slots=True)
@final
class QuotaConfig:
	daily_love_limit: int = 3
