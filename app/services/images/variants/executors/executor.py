from collections.abc import Iterator
from pathlib import Path
from typing import Protocol

from app.services.images.variants.types import (
	OriginalFile,
	OriginalImage,
	VariantCommitResult,
	VariantPlan,
	VariantPolicy,
)


class VariantExecutor(Protocol):
	def preprocess(self, file: OriginalFile) -> OriginalImage: ...

	def commit(
		self,
		image: OriginalImage,
		*,
		media_root: Path,
		plan: VariantPlan,
		policy: VariantPolicy,
	) -> Iterator[VariantCommitResult]: ...

	def postprocess(self, image: OriginalImage) -> None: ...
