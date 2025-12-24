from collections.abc import Sequence
from pathlib import Path
from typing import Protocol

from app.services.images.variants.types import (
	OriginalFile,
	VariantCommitResult,
	VariantPlan,
	VariantPolicy,
)


class VariantExecutor(Protocol):
	def execute(
		self,
		*,
		media_root: Path,
		file: OriginalFile,
		plan: VariantPlan,
		policy: VariantPolicy,
	) -> Sequence[VariantCommitResult]: ...
