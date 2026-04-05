from sqlalchemy import RowMapping, select
from sqlalchemy.orm import Session

from app.databases.tables import image_table, stats_table


def get_image_row(session: Session, *, ingest_id: int) -> RowMapping:
	row = (
		session.execute(
			select(image_table).where(image_table.c.ingest_id == ingest_id),
		)
		.mappings()
		.one()
	)
	return row


def get_stats_row(session: Session, *, ingest_id: int) -> RowMapping:
	row = (
		session.execute(
			select(stats_table).where(stats_table.c.ingest_id == ingest_id),
		)
		.mappings()
		.one()
	)
	return row
