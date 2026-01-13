from tests.services.activities.stats.factory import build_stats_record
from tests.services.images.utils import build_image_record
from tests.stubs.action import StubActionRepository
from tests.stubs.image import StubImageRepository
from tests.stubs.session import StubSession
from tests.stubs.stats import StubStatsRepository

from app.config.constants import VIEW_MILESTONES
from app.config.environments import env
from app.models.api.context.query import ContextQuery
from app.models.api.images.context import ImageRichModel, ImageSummaryModel
from app.models.enums import ActionKind
from app.services.views.context import ContextService


def test_get_context_returns_none_when_record_missing() -> None:
	image_repo = StubImageRepository()
	stats_repo = StubStatsRepository()
	action_repo = StubActionRepository()
	session = StubSession()
	service = ContextService(
		session,  # pyright: ignore[reportArgumentType]
		action_repo=action_repo,
		image_repo=image_repo,
		stats_repo=stats_repo,
		env=env,
	)

	result = service.get_context(123, query=ContextQuery())

	assert result is None
	assert image_repo.one_called_with == 123
	assert stats_repo.get_or_create_called_with is None
	assert action_repo.select_called_with is None
	assert session.begin_called == 1


def test_get_context_returns_summary_and_stats() -> None:
	image = build_image_record(5)
	image_repo = StubImageRepository()
	image_repo.one_response = image

	stats = build_stats_record(image.ingest_id)
	stats_repo = StubStatsRepository()
	stats_repo.stats_list_response = [stats]
	action_repo = StubActionRepository()
	session = StubSession()

	service = ContextService(
		session,  # pyright: ignore[reportArgumentType]
		action_repo=action_repo,
		image_repo=image_repo,
		stats_repo=stats_repo,
		env=env,
	)

	result = service.get_context(image.ingest_id, query=ContextQuery())

	assert image_repo.one_called_with == image.ingest_id
	assert stats_repo.get_or_create_called_with == image.ingest_id
	assert stats_repo.get_or_create_initial_score == env.score.initial_score
	assert action_repo.select_called_with == image.ingest_id
	assert session.begin_called == 1

	assert result is not None
	assert isinstance(result.image, ImageSummaryModel)
	assert result.image.id == image.ingest_id
	assert result.stats is not None
	assert result.stats.view_count == 1
	assert result.actions is not None
	assert len(result.actions) == 1
	assert result.actions[0].type == ActionKind.VIEW.name.lower()


def test_get_context_returns_rich_when_requested() -> None:
	image = build_image_record(7)
	image_repo = StubImageRepository()
	image_repo.one_response = image

	stats = build_stats_record(image.ingest_id)
	stats_repo = StubStatsRepository()
	stats_repo.stats_list_response = [stats]
	action_repo = StubActionRepository()
	session = StubSession()

	service = ContextService(
		session,  # pyright: ignore[reportArgumentType]
		action_repo=action_repo,
		image_repo=image_repo,
		stats_repo=stats_repo,
		env=env,
	)

	result = service.get_context(image.ingest_id, query=ContextQuery(level='rich'))

	assert result is not None
	assert isinstance(result.image, ImageRichModel)
	assert result.image.id == image.ingest_id
	assert result.image.variants
	assert result.image.variants[0]


def test_get_context_updates_view_milestone() -> None:
	image = build_image_record(9)
	image_repo = StubImageRepository()
	image_repo.one_response = image

	milestone = VIEW_MILESTONES[1]
	stats = build_stats_record(
		image.ingest_id,
		view_count=milestone - 1,
	)
	stats_repo = StubStatsRepository()
	stats_repo.stats_list_response = [stats]
	action_repo = StubActionRepository()
	session = StubSession()

	service = ContextService(
		session,  # pyright: ignore[reportArgumentType]
		action_repo=action_repo,
		image_repo=image_repo,
		stats_repo=stats_repo,
		env=env,
	)

	result = service.get_context(image.ingest_id, query=ContextQuery())

	assert result is not None
	assert result.stats is not None
	assert result.stats.view_milestone_count == milestone
	assert result.stats.view_milestone_archived_at == result.stats.last_viewed_at
