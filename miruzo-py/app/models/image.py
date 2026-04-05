from collections.abc import Sequence
from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, Field

from app.config import constants as c
from app.models.enums import ImageKind
from app.models.types import VariantEntry


@final
class Image(BaseModel):
	ingest_id: Annotated[
		int,
		Field(ge=c.INGEST_ID_MINIMUM, le=c.INGEST_ID_MAXIMUM),
	]
	ingested_at: datetime
	kind: ImageKind
	original: VariantEntry
	fallback: VariantEntry | None
	variants: Annotated[Sequence[VariantEntry], Field(min_length=1)]
