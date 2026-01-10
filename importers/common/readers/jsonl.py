import json
from collections.abc import Iterator
from pathlib import Path

from importers.common.models import GatakuImageRow


class JsonlReader:
	def __init__(self, path: Path) -> None:
		self._path = path

	def read(self, *, limit: int | None = None) -> Iterator[GatakuImageRow | None]:
		"""Yield parsed rows; invalid lines are emitted as None."""

		with self._path.open('r') as f:
			for i, line in enumerate(f):
				if limit is not None and i >= limit:
					break
				yield self._parse_line(line)

	def _parse_line(self, line: str) -> GatakuImageRow | None:
		try:
			record = json.loads(line)
		except Exception:
			return None

		return GatakuImageRow(
			filepath=Path(record['filepath']),
			sha256=record['sha256'],
			created_at=record.get('created_at'),
		)
