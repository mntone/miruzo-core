# pyright: reportArgumentType=false
# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownArgumentType=false
# pyright: reportUnknownMemberType=false

import sqlite3

from sqlalchemy.dialects.sqlite import Insert as SQLiteInsert
from sqlalchemy.dialects.sqlite import insert as sqlite_insert
from sqlalchemy.exc import IntegrityError
from sqlmodel import SQLModel

from app.services.activities.stats.repository.base import BaseStatsRepository


class SQLiteStatsRepository(BaseStatsRepository):
	def _is_unique_violation(self, error: IntegrityError) -> bool:
		orig = error.orig
		if not isinstance(orig, sqlite3.IntegrityError):
			return False

		if hasattr(orig, 'sqlite_errorcode'):
			return orig.sqlite_errorcode in {
				sqlite3.SQLITE_CONSTRAINT_UNIQUE,
				sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY,
			}

		message = str(orig)
		return message.startswith('UNIQUE constraint failed')

	def _build_insert(self, model: type[SQLModel]) -> SQLiteInsert:
		return sqlite_insert(model.__table__)
