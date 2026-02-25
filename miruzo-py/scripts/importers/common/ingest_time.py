from datetime import datetime, timezone
from logging import getLogger
from pathlib import Path

log = getLogger(__name__)


def resolve_captured_at(
	*,
	created_at_value: str | None,
	src_path: Path,
	warned_fallback: bool,
) -> tuple[datetime, bool, bool]:
	if created_at_value:
		try:
			return datetime.fromisoformat(created_at_value), False, warned_fallback
		except Exception:
			captured_at = datetime.fromtimestamp(src_path.stat().st_mtime, tz=timezone.utc)
			if not warned_fallback:
				log.warning(
					'invalid created_at detected; falling back to file mtime for subsequent entries',
				)
				warned_fallback = True
			return captured_at, True, warned_fallback

	captured_at = datetime.fromtimestamp(src_path.stat().st_mtime, tz=timezone.utc)
	if not warned_fallback:
		log.warning('missing created_at detected; falling back to file mtime for subsequent entries')
		warned_fallback = True
	return captured_at, True, warned_fallback
