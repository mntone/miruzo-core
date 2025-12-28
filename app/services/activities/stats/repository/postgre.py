# pyright: reportArgumentType=false
# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

from sqlalchemy.dialects.postgresql import Insert as PostgreInsert
from sqlalchemy.dialects.postgresql import insert as postgre_insert
from sqlmodel import SQLModel

from app.services.activities.stats.repository.base import BaseStatsRepository


class PostgreSQLStatsRepository(BaseStatsRepository):
	def _build_insert(self, model: type[SQLModel]) -> PostgreInsert:
		return postgre_insert(model.__table__)
