from tests.services.activities.stats.stubs import StubStatsRepository
from tests.services.activities.stats.utils import build_stats_record
from tests.services.images.stubs import StubImageRepository
from tests.services.images.utils import build_image_record

from app.services.activities.stats.service import StatsService
from app.services.images.query import ImageQueryService
from app.services.views.context import ContextService


def test_get_context_returns_none_when_record_missing() -> None:
	image_repo = StubImageRepository()
	stats_repo = StubStatsRepository()
	service = ContextService(
		image_query=ImageQueryService(image_repo),  # pyright: ignore[reportArgumentType]
		stats=StatsService(stats_repo),
	)

	result = service.get_context(123)

	assert result is None
	assert image_repo.one_called_with == 123
	assert stats_repo.upsert_called_with is None


def test_get_context_returns_summary_and_stats() -> None:
	image = build_image_record(5)
	image_repo = StubImageRepository()
	image_repo.one_response = image

	stats = build_stats_record(5)
	stats_repo = StubStatsRepository()
	stats_repo.stats_response = stats

	service = ContextService(
		image_query=ImageQueryService(image_repo),  # pyright: ignore[reportArgumentType]
		stats=StatsService(stats_repo),
	)

	result = service.get_context(image.ingest_id)

	assert image_repo.one_called_with == image.ingest_id
	assert stats_repo.upsert_called_with == image.ingest_id

	assert result is not None
	assert result.image.id == image.ingest_id
	assert result.stats is not None
	assert result.stats.score == stats.score
	assert result.stats.view_count == stats.view_count
