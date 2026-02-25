from dataclasses import dataclass


@dataclass(slots=True)
class RefInt:
	value: int
