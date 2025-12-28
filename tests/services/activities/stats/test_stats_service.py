from tests.services.activities.stats.stubs import StubStatsRepository
from tests.services.activities.stats.utils import build_stats_record

from app.services.activities.stats.service import StatsService


def test_get_by_ingest_id() -> None:
	image = build_stats_record(10)
	image_repo = StubStatsRepository()
	image_repo.stats_response = image

	service = StatsService(image_repo)

	response = service.get_by_ingest_id(image.ingest_id)
	assert response is not None
	assert response.ingest_id == image.ingest_id
	assert response.view_count == 1
