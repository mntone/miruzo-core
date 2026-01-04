# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from sqlalchemy.dialects.sqlite import Insert as SQLiteInsert
from sqlalchemy.dialects.sqlite import insert as sqlite_insert
from sqlmodel import SQLModel

from app.services.activities.stats.repository.base import BaseStatsRepository
from app.utils.database.sqlite import SQLiteUniqueViolationMixin


class SQLiteStatsRepository(SQLiteUniqueViolationMixin, BaseStatsRepository):
	def _build_insert(self, model: type[SQLModel]) -> SQLiteInsert:
		return sqlite_insert(model.__table__)
