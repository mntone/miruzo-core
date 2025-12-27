from collections.abc import Sequence
from typing import Set

from app.models.records import VariantRecord

_DEFAULT_FALLBACK_FORMATS = {'jpeg', 'png'}
_SUPPORTED_FORMATS = {'gif', 'jpeg', 'png', 'webp'}


def compute_allowed_formats(excluded: tuple[str, ...] | None) -> set[str]:
	"""
	Compute the effective allowed formats by removing the excluded ones.
	Fallback formats are always included to keep legacy clients working.
	"""

	excluded_set: set[str] = {fmt.lower() for fmt in excluded} if excluded else set()
	allowed = _SUPPORTED_FORMATS - excluded_set
	return allowed | _DEFAULT_FALLBACK_FORMATS


def normalize_variants_for_format(
	layers: Sequence[Sequence[VariantRecord]],
	allowed_formats: Set[str],
	*,
	keep_fallback: bool = True,
) -> Sequence[Sequence[VariantRecord]]:
	"""
	Filter and normalize variant layers based on allowed formats.

	Rules:
	- keep variants whose ``format`` value is present in ``allowed_formats``
	- the last layer is treated as a fallback layer (typically JPEG/PNG);
	  it is preserved entirely when ``keep_fallback`` is True
	- layers that become empty after filtering are removed (except fallback)
	- layer ordering is preserved
	"""

	if not layers:
		return layers

	last_layer_index = len(layers) - 1
	normalized_layers: list[Sequence[VariantRecord]] = []

	for index, layer in enumerate(layers):
		is_fallback_layer = index == last_layer_index

		# Fallback layers bypass filtering so legacy clients still work.
		if is_fallback_layer and keep_fallback:
			normalized_layers.append(layer)
			continue

		# Filter non-fallback layers.
		kept: list[VariantRecord] = []

		for spec in layer:
			if spec['format'] in allowed_formats:
				kept.append(spec)

		# Drop layers that have no surviving variants.
		if kept:
			normalized_layers.append(kept)

	# When keep_fallback=False the fallback layer only remains if it kept variants.
	return normalized_layers
