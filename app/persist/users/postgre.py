from app.persist.mixins.postgre import PostgreSQLUniqueViolationMixin
from app.persist.users.base import BaseUserRepository


class PostgreSQLUserRepository(PostgreSQLUniqueViolationMixin, BaseUserRepository): ...
