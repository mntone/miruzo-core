from contextlib import asynccontextmanager
from logging import getLogger
from typing import AsyncGenerator

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware

from app.config.environments import env
from app.databases import init_database
from app.services.images.variants.bootstrap import configure_pillow
from app.services.ingests.bootstrap import ensure_ingest_layout
from app.services.settings.factory import build_daily_period_resolver

log = getLogger('uvicorn.error')


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None]:
	ensure_ingest_layout(env)
	init_database()
	log.info(f'Starting miruzo API in {env.environment.value} mode')

	period_resolver = build_daily_period_resolver(
		day_start_offset=env.period.day_start_offset,
		initial_location=env.period.initial_location,
	)

	app.state['period_resolver'] = period_resolver

	yield


configure_pillow()

app = FastAPI(title='miruzo API', lifespan=lifespan)
app.add_middleware(GZipMiddleware, minimum_size=1000, compresslevel=9)

app.add_middleware(
	CORSMiddleware,
	allow_origins=['*'],
	allow_credentials=True,
	allow_methods=['*'],
	allow_headers=['*'],
)
