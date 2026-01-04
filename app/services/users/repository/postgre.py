from app.services.users.repository.base import BaseUserRepository
from app.utils.database.postgre import PostgreSQLUniqueViolationMixin


class PostgreSQLUserRepository(PostgreSQLUniqueViolationMixin, BaseUserRepository): ...
