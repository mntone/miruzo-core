from sqlalchemy import RowMapping, select
from sqlalchemy.orm import Session

from app.databases.tables import image_table, ingest_table, stats_table
from app.models.ingest import Ingest


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
