from app.services.users.repository.base import BaseUserRepository
from app.utils.database.sqlite import SQLiteUniqueViolationMixin


class SQLiteUserRepository(SQLiteUniqueViolationMixin, BaseUserRepository): ...
