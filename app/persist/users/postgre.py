from app.persist.users.base import BaseUserRepository
from app.utils.database.postgre import PostgreSQLUniqueViolationMixin


class PostgreSQLUserRepository(PostgreSQLUniqueViolationMixin, BaseUserRepository): ...
