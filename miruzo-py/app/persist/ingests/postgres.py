from typing import final

from sqlalchemy import insert

from app.databases.tables import ingest_table
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.protocol import IngestCreateInput


@final
class _IngestRepositoryPostgresImpl(_IngestRepositoryBaseImpl):
	def create(self, entry: IngestCreateInput) -> int:
		stmt = (
			insert(ingest_table)
			.values(
				**entry.model_dump(),
				updated_at=entry.ingested_at,
			)
			.returning(ingest_table.c.id)
		)
		row_id = self._session.execute(stmt).scalar_one()
		return row_id
