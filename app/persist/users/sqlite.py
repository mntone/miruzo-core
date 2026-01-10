from app.persist.mixins.sqlite import SQLiteUniqueViolationMixin
from app.persist.users.base import BaseUserRepository


class SQLiteUserRepository(SQLiteUniqueViolationMixin, BaseUserRepository): ...
