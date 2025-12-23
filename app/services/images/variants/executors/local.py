from collections.abc import Iterator
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
	def preprocess(self, file: OriginalFile) -> OriginalImage:
		with PILImage.open(file.file_info.absolute_path) as original_image:
			preprocessed_image = OriginalImage(
				image=preprocess_original(original_image, file.image_info),
				info=file.image_info,
			)

		return preprocessed_image

	def commit(
		self,
		image: OriginalImage,
		*,
		media_root: Path,
		plan: VariantPlan,
		policy: VariantPolicy,
	) -> Iterator[VariantCommitResult]:
		results = commit_variant_plan(
			plan=plan,
			policy=policy,
			original=image,
			media_root=media_root,
		)

		return results

	def postprocess(self, image: OriginalImage) -> None:
		pass
