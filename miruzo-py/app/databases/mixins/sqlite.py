# pyright: reportAttributeAccessIssue=false
# pyright: reportUnknownMemberType=false

import sqlite3

from sqlalchemy.exc import IntegrityError


class SQLiteUniqueViolationMixin:
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
