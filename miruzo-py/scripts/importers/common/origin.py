from dataclasses import dataclass
from logging import getLogger
from pathlib import Path

from scripts.importers.common.models import GatakuImageRow

log = getLogger(__name__)


@dataclass(frozen=True, slots=True)
class _OriginResolution:
	src_path: Path
	origin_relative_path: Path


class OriginResolver:
	def __init__(self, *, gataku_root: Path, gataku_assets_root: Path) -> None:
		self._gataku_root = gataku_root
		self._gataku_assets_root = gataku_assets_root

	def resolve(self, row: GatakuImageRow) -> _OriginResolution | None:
		raw_path = row.filepath
		src_path = raw_path if raw_path.is_absolute() else self._gataku_root / raw_path
		if not src_path.exists():
			log.warning(f'missing file: {src_path}')
			return None

		try:
			origin_relative_path = src_path.relative_to(self._gataku_assets_root)
		except ValueError:
			log.warning(f'file outside assets root: {src_path}')
			return None

		return _OriginResolution(
			src_path=src_path,
			origin_relative_path=origin_relative_path,
		)
