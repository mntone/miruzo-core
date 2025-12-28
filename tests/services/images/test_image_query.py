from tests.services.images.stubs import StubImageRepository
from tests.services.images.utils import build_image_record

from app.services.images.query import ImageQueryService


def test_get_latest_normalizes_variants_and_returns_cursor() -> None:
	image = build_image_record(1)
	image_repo = StubImageRepository()
	image_repo.list_response = ([image], image.captured_at)

	service = ImageQueryService(image_repo)  # type: ignore[arg-type]

	response = service.get_latest(cursor=None, limit=10, exclude_formats=('gif',))

	assert image_repo.list_called_with == {'cursor': None, 'limit': 10}
	assert response.cursor == image.captured_at
	assert len(response.items) == 1
	item = response.items[0]
	assert item.id == image.ingest_id
	assert [variant.format for variant in item.variants[0]] == ['webp', 'webp', 'webp']


def test_get_by_ingest_id() -> None:
	image = build_image_record(10)
	image_repo = StubImageRepository()
	image_repo.one_response = image

	service = ImageQueryService(image_repo)  # type: ignore[arg-type]

	response = service.get_by_ingest_id(image.ingest_id)
	assert response is not None
	assert response.ingest_id == image.ingest_id
