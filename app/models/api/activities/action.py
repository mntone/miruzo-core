from datetime import datetime
from typing import Annotated, final

from pydantic import BaseModel, ConfigDict, Field

from app.models.enums import ActionKind
from app.models.records import ActionRecord


@final
class ActionModel(BaseModel):
	"""Single action entry recorded for an image."""

	model_config = ConfigDict(
		title='Action model',
		extra='forbid',
		frozen=True,
	)

	type: Annotated[
		str,
		Field(
			title='Action type',
			description='action type identifier (mapped from ActionKind)',
		),
	]
	"""action kind identifier"""

	occurred_at: Annotated[
		datetime,
		Field(
			title='Occurred timestamp',
			description='timestamp when the action occurred',
		),
	]
	"""timestamp when the action occurred"""

	@classmethod
	def from_record(cls, action: ActionRecord) -> 'ActionModel':
		return cls(
			type=ActionKind(action.kind).name.lower(),
			occurred_at=action.occurred_at,
		)
