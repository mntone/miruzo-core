from datetime import datetime, timezone

from sqlalchemy import RowMapping, insert, select
from sqlalchemy.orm import Session

from app.databases.tables import image_table, ingest_table, stats_table
from app.models.enums import ProcessStatus, VisibilityStatus
from app.models.ingest import Ingest

next_id = 0


def add_ingest_row(
	session: Session,
	*,
	relative_path: str = 'l0orig/sample.png',
	process: ProcessStatus = ProcessStatus.PROCESSING,
	visibility: VisibilityStatus = VisibilityStatus.PRIVATE,
	fingerprint: str | None = None,
	ingested_at: datetime | None = None,
	captured_at: datetime | None = None,
) -> int:
	global next_id
	next_id += 1

	now = datetime.now(timezone.utc)
	ingested_at = ingested_at or now
	captured_at = captured_at or ingested_at

	stmt = (
		insert(ingest_table)
		.values(
			process=process,
			visibility=visibility,
			relative_path=relative_path,
			fingerprint=fingerprint or f'af{next_id:062d}',
			ingested_at=ingested_at,
			captured_at=captured_at,
			updated_at=ingested_at,
		)
		.returning(ingest_table.c.id)
	)
	row_id = session.execute(stmt).scalar_one()
	return row_id


def get_image_row(session: Session, *, ingest_id: int) -> RowMapping:
	row = (
		session.execute(
			select(image_table).where(image_table.c.ingest_id == ingest_id),
		)
		.mappings()
		.one()
	)
	return row


def get_ingest_row(session: Session, *, ingest_id: int) -> RowMapping:
	row = (
		session.execute(
			select(ingest_table).where(ingest_table.c.id == ingest_id),
		)
		.mappings()
		.one()
	)
	return row


def get_ingest_dto(session: Session, *, ingest_id: int) -> Ingest:
	row = get_ingest_row(session, ingest_id=ingest_id)
	dto = Ingest.model_validate(row)
	return dto


def get_stats_row(session: Session, *, ingest_id: int) -> RowMapping:
	row = (
		session.execute(
			select(stats_table).where(stats_table.c.ingest_id == ingest_id),
		)
		.mappings()
		.one()
	)
	return row
