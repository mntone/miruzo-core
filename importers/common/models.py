from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True, slots=True)
class GatakuImageRow:
	filepath: Path
	sha256: str
	created_at: str | None
