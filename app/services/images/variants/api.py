from collections.abc import Sequence
from typing import Set

from app.config.variant import is_variant_fallback_id
from app.models.types import VariantEntry

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
	layers: Sequence[Sequence[VariantEntry]],
	allowed_formats: Set[str],
	*,
	keep_fallback: bool = True,
) -> Sequence[Sequence[VariantEntry]]:
	"""
	Filter and normalize variant layers based on allowed formats.

	Rules:
	- keep variants whose ``format`` value is present in ``allowed_formats``
	- layers whose ``layer_id`` is identified as fallback (typically JPEG/PNG)
	  are preserved entirely when ``keep_fallback`` is True
	- layers that become empty after filtering are removed (except fallback)
	- layer ordering is preserved
	"""

	if not layers:
		return layers

	normalized_layers: list[Sequence[VariantEntry]] = []
	for layer in layers:
		if len(layer) == 0:
			continue

		is_fallback_layer = is_variant_fallback_id(layer[0]['layer_id'])

		# Fallback layers bypass filtering so legacy clients still work.
		if is_fallback_layer and keep_fallback:
			normalized_layers.append(layer)
			continue

		# Filter non-fallback layers.
		kept: list[VariantEntry] = []

		for spec in layer:
			if spec['format'] in allowed_formats:
				kept.append(spec)

		# Drop layers that have no surviving variants.
		if kept:
			normalized_layers.append(kept)

	# When keep_fallback=False the fallback layer only remains if it kept variants.
	return normalized_layers
