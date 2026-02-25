from app.databases.mixins.postgre import PostgreSQLUniqueViolationMixin
from app.persist.users.base import BaseUserRepository


class PostgreSQLUserRepository(PostgreSQLUniqueViolationMixin, BaseUserRepository): ...
