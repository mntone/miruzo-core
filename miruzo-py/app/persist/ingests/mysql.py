from typing import final

from sqlalchemy import JSON, BigInteger, DateTime, Integer, bindparam, text
from sqlalchemy.exc import NoResultFound

from app.models.enums import ExecutionStatus
from app.persist.ingests.base import _IngestRepositoryBaseImpl
from app.persist.ingests.protocol import IngestAppendExecutionInput

# keep latest non-success executions, restore chronological order,
# then append new success execution
_APPEND_SUCCESS_STMT = text(
	"""
	UPDATE ingests
	SET
		process=1,
		updated_at=:updated_at,
		executions=(
			SELECT CAST(CONCAT(
				'[', GROUP_CONCAT(CAST(v AS CHAR)ORDER BY i SEPARATOR ','), ']'
			)AS JSON)
			FROM(
				SELECT v, i
				FROM(
					SELECT v, i
					FROM JSON_TABLE(executions, '$[*]' COLUMNS(
						i FOR ORDINALITY,
						v JSON PATH '$',
						status INT PATH '$.status'
					))j
					WHERE status<>0
					ORDER BY i DESC
					LIMIT :max_retained_executions
				)l
				UNION ALL
				SELECT CAST(:execution AS JSON), 2147483647
			)m
		)
	WHERE id=:ingest_id
	""",
).bindparams(
	bindparam('updated_at', type_=DateTime),
	bindparam('ingest_id', type_=BigInteger),
	bindparam('execution', type_=JSON),
	bindparam('max_retained_executions', type_=Integer),
)

# append new error execution, then keep the latest retained executions
_APPEND_ERROR_STMT = text(
	"""
	UPDATE ingests
	SET
		updated_at=:updated_at,
		executions=JSON_EXTRACT(
			JSON_ARRAY_APPEND(executions,'$',CAST(:execution AS JSON)),
			CONCAT('$[last-',:max_retained_executions,' to last]')
		)
	WHERE id=:ingest_id
	""",
).bindparams(
	bindparam('updated_at', type_=DateTime),
	bindparam('ingest_id', type_=BigInteger),
	bindparam('execution', type_=JSON),
	bindparam('max_retained_executions', type_=Integer),
)


@final
class _IngestRepositoryMySQLImpl(_IngestRepositoryBaseImpl):
	def append_execution(self, entry: IngestAppendExecutionInput) -> None:
		params = {
			'ingest_id': entry.ingest_id,
			'updated_at': entry.updated_at,
			'execution': entry.execution.model_dump(mode='json'),
			'max_retained_executions': self._max_executions - 1,
		}
		if entry.execution.status == ExecutionStatus.SUCCESS:
			result = self._session.execute(_APPEND_SUCCESS_STMT, params)
		else:
			result = self._session.execute(_APPEND_ERROR_STMT, params)

		if result.rowcount != 1:  # pyright: ignore[reportAttributeAccessIssue]
			raise NoResultFound('No row was found when one was required')
