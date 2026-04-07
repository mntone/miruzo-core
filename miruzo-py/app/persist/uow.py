from contextlib import AbstractContextManager
from dataclasses import dataclass
from types import TracebackType
from typing import Callable, final

from sqlalchemy.orm import Session

from app.persist.images.implementation import create_image_repository
from app.persist.images.protocol import ImageRepository
from app.persist.ingests.factory import create_ingest_repository
from app.persist.ingests.protocol import IngestRepository
from app.persist.stats.implementation import create_stats_repository
from app.persist.stats.protocol import StatsRepository


@dataclass(frozen=True, slots=True)
@final
class Repositories:
	ingest: IngestRepository
	image: ImageRepository
	stats: StatsRepository


@final
class UnitOfWork(AbstractContextManager['UnitOfWork']):
	def __init__(self, *, session_factory: Callable[[], Session]) -> None:
		self._session_factory = session_factory
		self._session: Session | None = None
		self._repos: Repositories | None = None

	def __enter__(self) -> 'UnitOfWork':
		session = self._session_factory()
		self._session = session
		self._repos = Repositories(
			ingest=create_ingest_repository(session),
			image=create_image_repository(session),
			stats=create_stats_repository(session),
		)

		return self

	def __exit__(
		self,
		exc_type: type[BaseException] | None,
		exc: BaseException | None,
		tb: TracebackType | None,
	) -> None:
		session = self._session
		if session is None:
			return

		try:
			if exc_type is None:
				session.commit()
			else:
				session.rollback()
		finally:
			session.close()
			self._session = None
			self._repos = None

	def commit(self) -> None:
		session = self._session
		if session is None:
			raise RuntimeError('UnitOfWork is not active. Use within "with UnitOfWork(...)".')

		session.commit()

	def rollback(self) -> None:
		session = self._session
		if session is None:
			raise RuntimeError('UnitOfWork is not active. Use within "with UnitOfWork(...)".')

		session.rollback()

	@property
	def repositories(self) -> Repositories:
		repos = self._repos
		if repos is None:
			raise RuntimeError(
				'UnitOfWork repositories are unavailable before __enter__. Use within "with UnitOfWork(...)".',
			)

		return repos
