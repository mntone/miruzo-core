# pyright: reportAttributeAccessIssue=false
# pyright: reportArgumentType=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false
# pyright: reportUnknownVariableType=false

from sqlalchemy import select
from sqlalchemy.dialects.postgresql import Insert as PostgreInsert
from sqlalchemy.dialects.postgresql import insert as postgre_insert
from sqlmodel import SQLModel

from app.models.records import StatsRecord
from app.persist.mixins.postgre import PostgreSQLUniqueViolationMixin
from app.persist.stats.base import BaseStatsRepository


class PostgreSQLStatsRepository(PostgreSQLUniqueViolationMixin, BaseStatsRepository):
	def get_or_create(
		self,
		ingest_id: int,
		*,
		initial_score: int,
	) -> StatsRecord:
		stats_table = StatsRecord.__table__

		insert_statement = (
			# INSERT INTO stats
			postgre_insert(stats_table)
			# (ingest_id, score, score_evaluated)
			# VALUES (:ingest_id, :initial_score, :initial_score)
			.values(
				ingest_id=ingest_id,
				score=initial_score,
				score_evaluated=initial_score,
			)
			# ON CONFLICT (ingest_id)
			# DO NOTHING
			.on_conflict_do_nothing(
				index_elements=[stats_table.c.ingest_id],
			)
			# RETURNING *
			.returning(stats_table)
		)

		# WITH inserted AS
		insert_cte = insert_statement.cte('inserted')

		select_statement = (
			# SELECT * FROM inserted
			select(insert_cte)
			# UNION ALL
			.union_all(
				# SELECT * FROM stats
				select(stats_table)
				# WHERE ingest_id = :ingest_id
				.where(stats_table.c.ingest_id == ingest_id),
			)
			# LIMIT 1
			.limit(1)
		)

		stats = self._session.exec(select(StatsRecord).from_statement(select_statement)).scalar_one()  # pyright: ignore[reportCallIssue]

		return stats

	def _build_insert(self, model: type[SQLModel]) -> PostgreInsert:
		return postgre_insert(model.__table__)
