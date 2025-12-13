from typing import final

from pydantic import BaseModel, ConfigDict


@final
class HealthResponse(BaseModel):
	model_config = ConfigDict(title='Health response', validate_assignment=True)

	version: str
