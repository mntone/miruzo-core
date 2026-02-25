from collections.abc import Sequence
from pathlib import Path

from PIL import Image as PILImage

from app.services.images.variants.commit import commit_variant_plan
from app.services.images.variants.executors.executor import VariantExecutor
from app.services.images.variants.preprocess import preprocess_original
from app.services.images.variants.types import (
	OriginalFile,
	OriginalImage,
	VariantCommitResult,
	VariantPlan,
	VariantPolicy,
)


class LocalVariantExecutor(VariantExecutor):
	def execute(
		self,
		*,
		media_root: Path,
		file: OriginalFile,
		plan: VariantPlan,
		policy: VariantPolicy,
	) -> Sequence[VariantCommitResult]:
		with PILImage.open(file.file_info.absolute_path) as original_image:
			preprocessed_image = OriginalImage(
				image=preprocess_original(original_image, file.image_info),
				info=file.image_info,
			)

			results = commit_variant_plan(
				plan=plan,
				policy=policy,
				original=preprocessed_image,
				media_root=media_root,
			)

			return results
