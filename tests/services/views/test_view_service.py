from tests.services.activities.actions.stubs import StubActionRepository
from tests.services.activities.stats.stubs import StubStatsRepository
from tests.services.images.stubs import StubImageRepository
from tests.services.images.utils import build_image_record

from app.config.environments import env
from app.models.enums import ActionKind
from app.models.records import StatsRecord
from app.services.images.query import ImageQueryService
from app.services.views.context import ContextService


class _StubBegin:
	def __init__(self, session: 'StubSession') -> None:
		self._session = session

	def __enter__(self) -> None:
		self._session.begin_called += 1

	def __exit__(self, exc_type: object, exc: object, tb: object) -> bool:
		return False


class StubSession:
	def __init__(self) -> None:
		self.begin_called = 0

	def begin(self) -> _StubBegin:
		return _StubBegin(self)


def test_get_context_returns_none_when_record_missing() -> None:
	image_repo = StubImageRepository()
	stats_repo = StubStatsRepository()
	action_repo = StubActionRepository()
	session = StubSession()
	service = ContextService(
		session,  # pyright: ignore[reportArgumentType]
		action=action_repo,  # pyright: ignore[reportArgumentType]
		image_query=ImageQueryService(image_repo),  # pyright: ignore[reportArgumentType]
		stats=stats_repo,
		env=env,
	)

	result = service.get_context(123)

	assert result is None
	assert image_repo.one_called_with == 123
	assert stats_repo.get_or_create_called_with is None
	assert action_repo.select_called_with is None
	assert session.begin_called == 1


def test_get_context_returns_summary_and_stats() -> None:
	image = build_image_record(5)
	image_repo = StubImageRepository()
	image_repo.one_response = image

	stats = StatsRecord(
		ingest_id=image.ingest_id,
		score=env.score.initial_score,
		view_count=0,
		last_viewed_at=None,
	)
	stats_repo = StubStatsRepository()
	stats_repo.stats_response = stats
	action_repo = StubActionRepository()
	session = StubSession()

	service = ContextService(
		session,  # pyright: ignore[reportArgumentType]
		action=action_repo,  # pyright: ignore[reportArgumentType]
		image_query=ImageQueryService(image_repo),  # pyright: ignore[reportArgumentType]
		stats=stats_repo,
		env=env,
	)

	result = service.get_context(image.ingest_id)

	assert image_repo.one_called_with == image.ingest_id
	assert stats_repo.get_or_create_called_with == image.ingest_id
	assert stats_repo.get_or_create_initial_score == env.score.initial_score
	assert action_repo.select_called_with == image.ingest_id
	assert session.begin_called == 1

	assert result is not None
	assert result.image.id == image.ingest_id
	assert result.stats is not None
	assert result.stats.view_count == 1
	assert result.actions is not None
	assert len(result.actions) == 1
	assert result.actions[0].type == ActionKind.VIEW.name.lower()
