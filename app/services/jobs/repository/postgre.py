from app.services.jobs.repository.base import BaseJobRepository
from app.utils.database.postgre import PostgreSQLUniqueViolationMixin


class PostgreSQLJobRepository(PostgreSQLUniqueViolationMixin, BaseJobRepository): ...
