from dataclasses import dataclass
from typing import final


@final
@dataclass(frozen=True, slots=True)
class ScoreConfig:
	"""
	Configuration for score calculation.
	All values are interpreted by ScoreCalculator.
	"""
