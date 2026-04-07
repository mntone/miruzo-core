import pytest

from app.persist import uow as uow_module
from app.persist.uow import UnitOfWork


class DummySession:
	def __init__(self) -> None:
		self.commit_calls = 0
		self.rollback_calls = 0
		self.close_calls = 0

	def commit(self) -> None:
		self.commit_calls += 1

	def rollback(self) -> None:
		self.rollback_calls += 1

	def close(self) -> None:
		self.close_calls += 1


def test_repos_raises_before_enter() -> None:
	uow = UnitOfWork(session_factory=DummySession)  # pyright: ignore[reportArgumentType]

	with pytest.raises(RuntimeError, match='repositories are unavailable before __enter__'):
		_ = uow.repositories


def test_commit_and_rollback_raise_before_enter() -> None:
	uow = UnitOfWork(session_factory=DummySession)  # pyright: ignore[reportArgumentType]

	with pytest.raises(RuntimeError, match='UnitOfWork is not active'):
		uow.commit()
	with pytest.raises(RuntimeError, match='UnitOfWork is not active'):
		uow.rollback()


def test_enter_initializes_repositories(monkeypatch: pytest.MonkeyPatch) -> None:
	session = DummySession()
	ingest_repo = object()
	image_repo = object()
	stats_repo = object()

	monkeypatch.setattr(uow_module, 'create_ingest_repository', lambda _: ingest_repo)
	monkeypatch.setattr(uow_module, 'create_image_repository', lambda _: image_repo)
	monkeypatch.setattr(uow_module, 'create_stats_repository', lambda _: stats_repo)

	with UnitOfWork(session_factory=lambda: session) as uow:  # pyright: ignore[reportArgumentType]
		assert uow.repositories.ingest is ingest_repo
		assert uow.repositories.image is image_repo
		assert uow.repositories.stats is stats_repo

	assert session.commit_calls == 1
	assert session.rollback_calls == 0
	assert session.close_calls == 1


def test_exit_rolls_back_and_closes_on_exception(monkeypatch: pytest.MonkeyPatch) -> None:
	session = DummySession()
	monkeypatch.setattr(uow_module, 'create_ingest_repository', lambda _: object())
	monkeypatch.setattr(uow_module, 'create_image_repository', lambda _: object())
	monkeypatch.setattr(uow_module, 'create_stats_repository', lambda _: object())

	with pytest.raises(RuntimeError, match='boom'):
		with UnitOfWork(session_factory=lambda: session):  # pyright: ignore[reportArgumentType]
			raise RuntimeError('boom')

	assert session.commit_calls == 0
	assert session.rollback_calls == 1
	assert session.close_calls == 1
