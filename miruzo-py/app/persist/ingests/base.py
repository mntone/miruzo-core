from sqlalchemy import bindparam, insert, select, update
from sqlalchemy.orm import Session

from app.databases.tables import ingest_table
from app.models.enums import ExecutionStatus, ProcessStatus
from app.models.ingest import Execution, executions_adapter
from app.persist.ingests.protocol import IngestAppendExecutionInput, IngestCreateInput

_EXECUTIONS_SELECT_STATEMENT = select(ingest_table.c.executions).where(
	ingest_table.c.id == bindparam('ingest_id'),
)


class _IngestRepositoryBaseImpl:
	def __init__(self, session: Session, *, max_executions: int) -> None:
		self._session = session
		self._max_executions = max_executions

	def create(self, entry: IngestCreateInput) -> int:
		stmt = insert(ingest_table).values(
			**entry.model_dump(),
			updated_at=entry.ingested_at,
		)
		row_id = self._session.execute(stmt).inserted_primary_key[0]  # pyright: ignore[reportAttributeAccessIssue]
		return row_id

	def append_execution(self, entry: IngestAppendExecutionInput) -> None:
		executions_row = self._session.execute(
			_EXECUTIONS_SELECT_STATEMENT,
			{'ingest_id': entry.ingest_id},
		).scalar_one()
		executions_dto = executions_adapter.validate_python(executions_row)

		executions: list[Execution]
		if entry.execution.status == ExecutionStatus.SUCCESS:
			executions = [e for e in executions_dto if e.status != ExecutionStatus.SUCCESS]
		else:
			executions = executions_dto.copy()
		executions.append(entry.execution)

		values = {
			'updated_at': entry.updated_at,
			'executions': executions_adapter.dump_python(executions[-self._max_executions :], mode='json'),
		}
		if entry.execution.status == ExecutionStatus.SUCCESS:
			values['process'] = ProcessStatus.FINISHED

		stmt = update(ingest_table).where(ingest_table.c.id == entry.ingest_id).values(**values)
		self._session.execute(stmt)
