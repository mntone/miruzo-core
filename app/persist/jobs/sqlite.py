from app.persist.jobs.base import BaseJobRepository
from app.utils.database.sqlite import SQLiteUniqueViolationMixin


class SQLiteJobRepository(SQLiteUniqueViolationMixin, BaseJobRepository): ...
