from datetime import datetime, timezone

from app.models.enums import ProcessStatus, VisibilityStatus
from app.models.ingest import Ingest


def make_ingest_fixture(
	ingest_id: int,
	*,
	relative_path: str = 'l0orig/sample.png',
	process: ProcessStatus = ProcessStatus.PROCESSING,
	visibility: VisibilityStatus = VisibilityStatus.PRIVATE,
	ingested_at: datetime | None = None,
	captured_at: datetime | None = None,
	fingerprint: str | None = None,
) -> Ingest:
	now = datetime.now(timezone.utc)
	ingested_at = ingested_at or now
	captured_at = captured_at or ingested_at
	return Ingest(
		id=ingest_id,
		process=process,
		visibility=visibility,
		relative_path=relative_path,
		fingerprint=fingerprint or f'af{ingest_id:062d}',
		ingested_at=ingested_at,
		captured_at=captured_at,
		updated_at=ingested_at,
		executions=[],
	)
