from dataclasses import dataclass, field
from pathlib import Path
from typing import Literal, TypeAlias

from PIL import Image as PILImage

from app.config.variant import VariantSlotkey, VariantSpec
from app.services.images.variants.utils import ImageInfo, parse_variant_slotkey


@dataclass(frozen=True, slots=True)
class OriginalImage:
	image: PILImage.Image
	info: ImageInfo


@dataclass(frozen=True, slots=True)
class VariantFile:
	bytes: int
	info: ImageInfo
	path: Path
	variant_dir: str
	_slotkey_cache: VariantSlotkey | None = field(init=False, default=None, repr=False)

	@property
	def slotkey(self) -> VariantSlotkey:
		if self._slotkey_cache is None:
			slotkey = parse_variant_slotkey(self.variant_dir)
			object.__setattr__(self, '_slotkey_cache', slotkey)
			return slotkey

		return self._slotkey_cache


@dataclass(frozen=True, slots=True)
class VariantComparison:
	spec: VariantSpec
	file: VariantFile


@dataclass(frozen=True, slots=True)
class VariantPlan:
	matched: list[VariantComparison]
	mismatched: list[VariantComparison]
	missing: list[VariantSpec]
	orphaned: list[VariantFile]


@dataclass(frozen=True, slots=True)
class VariantPolicy:
	regenerate_mismatched: bool
	generate_missing: bool
	delete_orphaned: bool


@dataclass(frozen=True, slots=True)
class VariantReport:
	spec: VariantSpec
	file: VariantFile


_VariantCommitAction: TypeAlias = Literal['reuse', 'generate', 'regenerate', 'delete']
_VariantCommitFailureReason: TypeAlias = Literal[
	'file_already_missing',
	'os_error',
	'parent_dir_not_found',
	'permission_denied',
	'save_failed',
]


@dataclass(frozen=True, slots=True)
class VariantCommitResult:
	action: _VariantCommitAction
	result: Literal['success', 'failure']
	reason: _VariantCommitFailureReason | None
	report: VariantReport | None

	@classmethod
	def success(cls, action: _VariantCommitAction, report: VariantReport | None) -> 'VariantCommitResult':
		return cls(action, 'success', None, report)

	@classmethod
	def failure(
		cls,
		action: _VariantCommitAction,
		reason: _VariantCommitFailureReason,
	) -> 'VariantCommitResult':
		return cls(action, 'failure', reason, None)
