from typing import final

from sqlalchemy import BigInteger, DateTime, Integer, bindparam, insert, text
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.exc import NoResultFound

from app.databases.tables import ingest_table
from app.models.enums import ExecutionStatus
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.protocol import IngestAppendExecutionInput, IngestCreateInput

_APPEND_SUCCESS_SQL = text(
	"""
	UPDATE ingests
	SET
		process=1,
		updated_at=:updated_at,
		executions=(
			SELECT jsonb_agg(v ORDER BY i)
			FROM(
				SELECT v, i
				FROM(
					SELECT v, i
					FROM jsonb_array_elements(executions) WITH ORDINALITY AS t(v, i)
					WHERE(v->>'status')::int<>0
					ORDER BY i DESC
					LIMIT :execution_maximum-1
				)l
				UNION ALL
				SELECT :execution, 2147483647
			)m
		)
	WHERE id=:ingest_id
	""",
).bindparams(
	bindparam('updated_at', type_=DateTime),
	bindparam('execution', type_=JSONB),
	bindparam('execution_maximum', type_=Integer),
	bindparam('ingest_id', type_=BigInteger),
)


_APPEND_ERROR_SQL = text(
	"""
	UPDATE ingests
	SET
		updated_at=:updated_at,
		executions=jsonb_path_query_array(
			executions||:execution,
			('$[last-'||(:execution_maximum-1)||' to last]')::jsonpath
		)
	WHERE id=:ingest_id
	""",
).bindparams(
	bindparam('updated_at', type_=DateTime),
	bindparam('execution', type_=JSONB),
	bindparam('execution_maximum', type_=Integer),
	bindparam('ingest_id', type_=BigInteger),
)


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

	def append_execution(self, entry: IngestAppendExecutionInput) -> None:
		params = {
			'ingest_id': entry.ingest_id,
			'updated_at': entry.updated_at,
			'execution': entry.execution.model_dump(mode='json'),
			'execution_maximum': self._max_executions,
		}
		if entry.execution.status == ExecutionStatus.SUCCESS:
			result = self._session.execute(_APPEND_SUCCESS_SQL, params)
		else:
			result = self._session.execute(_APPEND_ERROR_SQL, params)

		if result.rowcount != 1:  # pyright: ignore[reportAttributeAccessIssue]
			raise NoResultFound('No row was found when one was required')
