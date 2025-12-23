from collections.abc import Iterator
from pathlib import Path

import pytest
from PIL import UnidentifiedImageError as PILUnidentifiedImageError
from sqlalchemy.exc import DataError

from app.models.enums import ExecutionStatus
from app.services.images.variants.pipeline_execution import VariantPipelineExecutionSession
from app.services.images.variants.types import (
	OriginalFile,
	OriginalImage,
	VariantCommitResult,
	VariantPlan,
	VariantPolicy,
)


class DummyExecutor:
	def preprocess(self, file: OriginalFile) -> OriginalImage:  # noqa: ARG002
		raise AssertionError('preprocess should not be called')

	def commit(  # pragma: no cover - not executed
		self,
		image: OriginalImage,  # noqa: ARG002
		*,
		media_root: Path,  # noqa: ARG002
		plan: VariantPlan,  # noqa: ARG002
		policy: VariantPolicy,  # noqa: ARG002
	) -> Iterator[VariantCommitResult]:
		raise AssertionError('commit should not be called')

	def postprocess(self, image: OriginalImage) -> None:  # noqa: ARG002
		raise AssertionError('postprocess should not be called')


def test_execution_session_to_entry_records_success() -> None:
	session = VariantPipelineExecutionSession(DummyExecutor())

	with session:
		with session.phase('inspect'):
			pass

	entry = session.to_entry()

	assert entry['status'] == ExecutionStatus.SUCCESS
	assert entry['error_type'] is None
	assert entry['error_message'] is None
	assert entry['inspect'] is not None
	assert entry['overall'] is not None


@pytest.mark.parametrize(
	('exc', 'status', 'error_type', 'error_message', 'swallow'),
	[
		(
			PILUnidentifiedImageError('broken'),
			ExecutionStatus.IMAGE_ERROR,
			'UnidentifiedImageError',
			'broken',
			True,
		),
		(
			DataError('stmt', {}, Exception('orig')),
			ExecutionStatus.DB_ERROR,
			'DataError',
			None,
			True,
		),
		(
			ValueError('boom'),
			ExecutionStatus.UNKNOWN_ERROR,
			'ValueError',
			'boom',
			False,
		),
	],
)
def test_execution_session_to_entry_records_errors(
	exc: Exception,
	status: ExecutionStatus,
	error_type: str,
	error_message: str | None,
	swallow: bool,
) -> None:
	session = VariantPipelineExecutionSession(DummyExecutor())

	def raise_in_session() -> None:
		with session:
			with session.phase('inspect'):
				pass
			raise exc

	if swallow:
		raise_in_session()
	else:
		with pytest.raises(type(exc)):
			raise_in_session()

	entry = session.to_entry()

	assert entry['status'] == status
	assert entry['error_type'] == error_type
	if error_message is None:
		assert entry['error_message']
	else:
		assert entry['error_message'] == error_message
	assert entry['inspect'] is not None
	assert entry['overall'] is not None
