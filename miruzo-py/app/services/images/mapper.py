from collections.abc import Sequence

from app.config.variant import VariantLayerSpec
from app.domain.images.cursor import TImageListCursor
from app.models.api.images.list import ImageListModel
from app.models.api.images.responses import ImageListResponse
from app.models.records import ImageRecord
from app.services.images.variants.api import compute_allowed_formats, normalize_variants_for_format
from app.services.images.variants.mapper import map_variants_to_layers


def map_image_records_to_list_response(
	images: Sequence[ImageRecord],
	*,
	next_cursor: TImageListCursor | None,
	exclude_formats: tuple[str, ...],
	variant_layers: Sequence[VariantLayerSpec],
) -> ImageListResponse[TImageListCursor]:
	"""Build an ImageListResponse from records with normalized variants."""

	allowed_formats = compute_allowed_formats(exclude_formats)

	items: list[ImageListModel] | None = None
	if len(images) != 0:
		items = []
		for image in images:
			layers = map_variants_to_layers(image.variants, spec=variant_layers)
			normalized_layers = normalize_variants_for_format(layers, allowed_formats)
			image_model = ImageListModel.from_record(image, normalized_layers)
			items.append(image_model)

	response = ImageListResponse(
		items=items,
		cursor=next_cursor,
	)

	return response
