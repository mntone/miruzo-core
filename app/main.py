from contextlib import asynccontextmanager
from logging import getLogger
from typing import AsyncGenerator

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from fastapi.staticfiles import StaticFiles

from app.core.settings import settings
from app.database import init_database
from app.routers.health import router as health
from app.routers.images import router as images

log = getLogger('uvicorn.error')


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None]:
	init_database()
	log.info(f'Starting miruzo API in {settings.environment.value} mode')
	yield


app = FastAPI(title='miruzo API', lifespan=lifespan)
app.add_middleware(GZipMiddleware, minimum_size=1000, compresslevel=9)
app.include_router(images, prefix='/api')
app.include_router(health, prefix='/api')

app.add_middleware(
	CORSMiddleware,
	allow_origins=['*'],
	allow_credentials=True,
	allow_methods=['*'],
	allow_headers=['*'],
)
app.mount(
	settings.public_media_root,
	StaticFiles(directory=settings.media_root, follow_symlink=True),
	name='media',
)
