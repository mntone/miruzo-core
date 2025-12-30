# pyright: reportArgumentType=false
# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from sqlalchemy.dialects.postgresql import Insert as PostgreInsert
from sqlalchemy.dialects.postgresql import insert as postgre_insert
from sqlalchemy.exc import IntegrityError
from sqlmodel import SQLModel

from app.services.activities.stats.repository.base import BaseStatsRepository


class PostgreSQLStatsRepository(BaseStatsRepository):
	def _is_unique_violation(self, error: IntegrityError) -> bool:
		orig = error.orig
		if orig is None:
			return False

		pgcode = getattr(orig, 'pgcode', None)
		if pgcode == '23505':
			return True

		try:
			from psycopg2.errors import UniqueViolation
		except ImportError:
			return False

		return isinstance(orig, UniqueViolation)

	def _build_insert(self, model: type[SQLModel]) -> PostgreInsert:
		return postgre_insert(model.__table__)
