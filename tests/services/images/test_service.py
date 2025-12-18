from datetime import datetime, timezone

from tests.services.images.utils import build_image_record

from app.models.records import ImageRecord, StatsRecord
from app.services.images.service import ImageService


class StubImageRepository:
	def __init__(self) -> None:
		self.list_response: tuple[list[ImageRecord], datetime | None] = ([], None)
		self.detail_response: ImageRecord | None = None
		self.stats_response: StatsRecord | None = None

		self.list_called_with: dict[str, object] | None = None
		self.detail_called_with: int | None = None
		self.upsert_called_with: int | None = None

	def get_list(self, *, cursor: datetime | None, limit: int) -> tuple[list[ImageRecord], datetime | None]:
		self.list_called_with = {'cursor': cursor, 'limit': limit}
		return self.list_response

	def get_detail(self, image_id: int) -> ImageRecord | None:
		self.detail_called_with = image_id
		return self.detail_response

	def upsert_stats_with_increment(self, image_id: int) -> StatsRecord:
		self.upsert_called_with = image_id
		if self.stats_response is None:
			raise RuntimeError('stats_response not configured')
		return self.stats_response


def _stats_record(image_id: int) -> StatsRecord:
	return StatsRecord(
		image_id=image_id,
		favorite=False,
		score=5,
		view_count=1,
		last_viewed_at=datetime.now(timezone.utc),
	)


def test_get_latest_normalizes_variants_and_returns_cursor() -> None:
	repo = StubImageRepository()
	image = build_image_record(1, ['webp', 'gif'])
	repo.list_response = ([image], image.captured_at)

	service = ImageService(repo)  # type: ignore[arg-type]

	response = service.get_latest(cursor=None, limit=10, exclude_formats=('gif',))

	assert repo.list_called_with == {'cursor': None, 'limit': 10}
	assert response.cursor == image.captured_at
	assert len(response.items) == 1
	item = response.items[0]
	assert item.id == image.id
	assert [variant.format for variant in item.variants[0]] == ['webp']


def test_get_context_returns_none_when_record_missing() -> None:
	repo = StubImageRepository()
	service = ImageService(repo)  # type: ignore[arg-type]

	result = service.get_context(123)

	assert result is None
	assert repo.detail_called_with == 123
	assert repo.upsert_called_with is None


def test_get_context_returns_summary_and_stats() -> None:
	repo = StubImageRepository()
	image = build_image_record(5, ['webp', 'gif'])
	repo.detail_response = image
	repo.stats_response = _stats_record(5)

	service = ImageService(repo)  # type: ignore[arg-type]

	result = service.get_context(image.id)

	assert repo.detail_called_with == image.id
	assert repo.upsert_called_with == image.id
	assert result is not None
	assert result.image.id == image.id
	assert result.stats is not None
	assert result.stats.view_count == repo.stats_response.view_count
