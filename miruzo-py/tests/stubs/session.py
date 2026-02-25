from typing import Callable, final

from tests.stubs.ref import RefInt


@final
class _StubSessionBegin:
	def __init__(self, session: 'StubSession') -> None:
		self._session = session

	def __enter__(self) -> None:
		self._session.begin_called += 1
		return None

	def __exit__(self, exc_type: object, exc: object, tb: object) -> bool | None:  # pyright: ignore[reportMissingParameterType, reportUnknownParameterType]
		return False


@final
class StubSession:
	def __init__(self) -> None:
		self.begin_called = 0

	def __enter__(self) -> 'StubSession':
		return self

	def __exit__(self, exc_type: object, exc: object, tb: object) -> None:  # pyright: ignore[reportMissingParameterType, reportUnknownParameterType]
		return None

	def begin(self) -> '_StubSessionBegin':
		return _StubSessionBegin(self)


def create_stub_session() -> StubSession:
	return StubSession()


def create_stub_session_factory(ref: RefInt) -> Callable[[], StubSession]:
	def factory() -> StubSession:
		ref.value += 1
		return StubSession()

	return factory
