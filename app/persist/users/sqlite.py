from app.persist.users.base import BaseUserRepository
from app.utils.database.sqlite import SQLiteUniqueViolationMixin


class SQLiteUserRepository(SQLiteUniqueViolationMixin, BaseUserRepository): ...
