# miruzo-core

[![License under GPLv3](data:image/svg+xml;base64,PHN2ZyBkYXRhLXYtM2M4N2I3YjQ9IiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB3aWR0aD0iMTU1LjE3NjQyMjExOTE0MDYyIiBoZWlnaHQ9IjM1IiB2aWV3Qm94PSIwIDAgMTU1LjE3NjQyMjExOTE0MDYyIDM1IiBjbGFzcz0iYmFkZ2Utc3ZnIj48ZGVmcyBkYXRhLXYtM2M4N2I3YjQ9IiI+PCEtLS0tPjwhLS0tLT48IS0tLS0+PC9kZWZzPjxyZWN0IGRhdGEtdi0zYzg3YjdiND0iIiB3aWR0aD0iODQuMzQxMjI0NjcwNDEwMTYiIGhlaWdodD0iMzUiIGZpbGw9IiM1NTU1NTUiLz48cmVjdCBkYXRhLXYtM2M4N2I3YjQ9IiIgeD0iODQuMzQxMjI0NjcwNDEwMTYiIHdpZHRoPSI3MC44MzUxOTc0NDg3MzA0NyIgaGVpZ2h0PSIzNSIgZmlsbD0iIzAwN2VjNiIvPjwhLS0tLT48dGV4dCBkYXRhLXYtM2M4N2I3YjQ9IiIgeD0iNDIuMTcwNjEyMzM1MjA1MDgiIHk9IjE3LjUiIGR5PSIwLjM1ZW0iIGZvbnQtc2l6ZT0iMTIiIGZvbnQtZmFtaWx5PSJSb2JvdG8sIHNhbnMtc2VyaWYiIGZpbGw9IiNGRkZGRkYiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGxldHRlci1zcGFjaW5nPSIyIiBmb250LXdlaWdodD0iNDAwIiBmb250LXN0eWxlPSJub3JtYWwiIHRleHQtZGVjb3JhdGlvbj0ibm9uZSIgZmlsbC1vcGFjaXR5PSIxIiBmb250LXZhcmlhbnQ9Im5vcm1hbCIgc3R5bGU9InRleHQtdHJhbnNmb3JtOiB1cHBlcmNhc2U7Ij5MSUNFTlNFPC90ZXh0PjwhLS0tLT48dGV4dCBkYXRhLXYtM2M4N2I3YjQ9IiIgeD0iMTE5Ljc1ODgyMzM5NDc3NTM5IiB5PSIxNy41IiBkeT0iMC4zNWVtIiBmb250LXNpemU9IjEyIiBmb250LWZhbWlseT0iUm9ib3RvLCBzYW5zLXNlcmlmIiBmaWxsPSIjRkZGRkZGIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LXdlaWdodD0iOTAwIiBsZXR0ZXItc3BhY2luZz0iMiIgZm9udC1zdHlsZT0ibm9ybWFsIiB0ZXh0LWRlY29yYXRpb249Im5vbmUiIGZpbGwtb3BhY2l0eT0iMSIgZm9udC12YXJpYW50PSJub3JtYWwiIHN0eWxlPSJ0ZXh0LXRyYW5zZm9ybTogdXBwZXJjYXNlOyI+R1BMIDMrPC90ZXh0PjwhLS0tLT48L3N2Zz4=)](./LICENSE)
[![Made with FastAPI](data:image/svg+xml;base64,PHN2ZyBkYXRhLXYtM2M4N2I3YjQ9IiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB3aWR0aD0iMjE1LjE4NzUiIGhlaWdodD0iMzUiIHZpZXdCb3g9IjAgMCAyMTUuMTg3NSAzNSIgY2xhc3M9ImJhZGdlLXN2ZyI+PGRlZnMgZGF0YS12LTNjODdiN2I0PSIiPjwhLS0tLT48IS0tLS0+PCEtLS0tPjwvZGVmcz48cmVjdCBkYXRhLXYtM2M4N2I3YjQ9IiIgd2lkdGg9Ijk4LjUzMTI1IiBoZWlnaHQ9IjM1IiBmaWxsPSIjZWY0MDQxIi8+PHJlY3QgZGF0YS12LTNjODdiN2I0PSIiIHg9Ijk4LjUzMTI1IiB3aWR0aD0iMTE2LjY1NjI1IiBoZWlnaHQ9IjM1IiBmaWxsPSIjYzEyODJkIi8+PCEtLS0tPjx0ZXh0IGRhdGEtdi0zYzg3YjdiND0iIiB4PSI0OS4yNjU2MjUiIHk9IjE3LjUiIGR5PSIwLjM1ZW0iIGZvbnQtc2l6ZT0iMTIiIGZvbnQtZmFtaWx5PSJNb250c2VycmF0LCBzYW5zLXNlcmlmIiBmaWxsPSIjRkZGRkZGIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBsZXR0ZXItc3BhY2luZz0iMS41IiBmb250LXdlaWdodD0iNDAwIiBmb250LXN0eWxlPSJub3JtYWwiIHRleHQtZGVjb3JhdGlvbj0ibm9uZSIgZmlsbC1vcGFjaXR5PSIxIiBmb250LXZhcmlhbnQ9Im5vcm1hbCIgc3R5bGU9InRleHQtdHJhbnNmb3JtOiB1cHBlcmNhc2U7Ij5NQURFIFdJVEg8L3RleHQ+PGcgZGF0YS12LTNjODdiN2I0PSIiIHRyYW5zZm9ybT0idHJhbnNsYXRlKDE4MS4xODc1LCA1LjUpIHNjYWxlKDEpIj48cGF0aCBkYXRhLXYtM2M4N2I3YjQ9IiIgZD0iTTEyIC4wMzg3QzUuMzcyOS4wMzg0LjAwMDMgNS4zOTMxIDAgMTEuOTk4OGMtLjAwMSA2LjYwNjYgNS4zNzIgMTEuOTYyOCAxMiAxMS45NjI1IDYuNjI4LjAwMDMgMTIuMDAxLTUuMzU1OSAxMi0xMS45NjI1LS4wMDAzLTYuNjA1Ny01LjM3MjktMTEuOTYwNC0xMi0xMS45Nm0tLjgyOSA1LjQxNTNoNy41NWwtNy41ODA1IDUuMzI4NGg1LjE4MjhMNS4yNzkgMTguNTQzNnEyLjk0NjYtNi41NDQ0IDUuODkyLTEzLjA4OTYiIGZpbGw9IiNGRkZGRkYiLz48L2c+PHRleHQgZGF0YS12LTNjODdiN2I0PSIiIHg9IjE0Mi44NTkzNzUiIHk9IjE3LjUiIGR5PSIwLjM1ZW0iIGZvbnQtc2l6ZT0iMTIiIGZvbnQtZmFtaWx5PSJNb250c2VycmF0LCBzYW5zLXNlcmlmIiBmaWxsPSIjRkZGRkZGIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LXdlaWdodD0iOTAwIiBsZXR0ZXItc3BhY2luZz0iMiIgZm9udC1zdHlsZT0ibm9ybWFsIiB0ZXh0LWRlY29yYXRpb249Im5vbmUiIGZpbGwtb3BhY2l0eT0iMSIgZm9udC12YXJpYW50PSJub3JtYWwiIHN0eWxlPSJ0ZXh0LXRyYW5zZm9ybTogdXBwZXJjYXNlOyI+RkFTVCBBUEk8L3RleHQ+PCEtLS0tPjwvc3ZnPg==)](https://fastapi.tiangolo.com/)
[![Written by Python](data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyMjEuNTkzNTk5OTk5OTk5OTgiIGhlaWdodD0iMzUiIHZpZXdCb3g9IjAgMCAyMjEuNTkzNTk5OTk5OTk5OTggMzUiPjxyZWN0IHdpZHRoPSIxMDMuNTE2Nzk5OTk5OTk5OTkiIGhlaWdodD0iMzUiIGZpbGw9IiM4ZmM5NjUiIC8+PHJlY3QgeD0iMTAzLjUxNjc5OTk5OTk5OTk5IiB3aWR0aD0iMTE4LjA3NjgiIGhlaWdodD0iMzUiIGZpbGw9IiM0MTliNWEiIC8+PHRleHQgeD0iNTEuNzU4Mzk5OTk5OTk5OTk1IiB5PSIxNy41IiBkeT0iMC4zNWVtIiBmb250LXNpemU9IjEyIiBmb250LWZhbWlseT0iTW9udHNlcnJhdCwgc2Fucy1zZXJpZiIgZmlsbD0iI0ZGRkZGRiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgbGV0dGVyLXNwYWNpbmc9IjEiIGZvbnQtd2VpZ2h0PSI0MDAiIGZpbGwtb3BhY2l0eT0iMSIgc3R5bGU9InRleHQtdHJhbnNmb3JtOiB1cHBlcmNhc2UiPldSSVRURU4gQlk8L3RleHQ+PGcgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMTg3LjU5MzU5OTk5OTk5OTk4LCA1LjUpIHNjYWxlKDEpIj48cGF0aCBkPSJNMTQuMjUuMThsLjkuMi43My4yNi41OS4zLjQ1LjMyLjM0LjM0LjI1LjM0LjE2LjMzLjEuMy4wNC4yNi4wMi4yLS4wMS4xM1Y4LjVsLS4wNS42My0uMTMuNTUtLjIxLjQ2LS4yNi4zOC0uMy4zMS0uMzMuMjUtLjM1LjE5LS4zNS4xNC0uMzMuMS0uMy4wNy0uMjYuMDQtLjIxLjAySDguNzdsLS42OS4wNS0uNTkuMTQtLjUuMjItLjQxLjI3LS4zMy4zMi0uMjcuMzUtLjIuMzYtLjE1LjM3LS4xLjM1LS4wNy4zMi0uMDQuMjctLjAyLjIxdjMuMDZIMy4xN2wtLjIxLS4wMy0uMjgtLjA3LS4zMi0uMTItLjM1LS4xOC0uMzYtLjI2LS4zNi0uMzYtLjM1LS40Ni0uMzItLjU5LS4yOC0uNzMtLjIxLS44OC0uMTQtMS4wNS0uMDUtMS4yMy4wNi0xLjIyLjE2LTEuMDQuMjQtLjg3LjMyLS43MS4zNi0uNTcuNC0uNDQuNDItLjMzLjQyLS4yNC40LS4xNi4zNi0uMS4zMi0uMDUuMjQtLjAxaC4xNmwuMDYuMDFoOC4xNnYtLjgzSDYuMThsLS4wMS0yLjc1LS4wMi0uMzcuMDUtLjM0LjExLS4zMS4xNy0uMjguMjUtLjI2LjMxLS4yMy4zOC0uMi40NC0uMTguNTEtLjE1LjU4LS4xMi42NC0uMS43MS0uMDYuNzctLjA0Ljg0LS4wMiAxLjI3LjA1em0tNi4zIDEuOThsLS4yMy4zMy0uMDguNDEuMDguNDEuMjMuMzQuMzMuMjIuNDEuMDkuNDEtLjA5LjMzLS4yMi4yMy0uMzQuMDgtLjQxLS4wOC0uNDEtLjIzLS4zMy0uMzMtLjIyLS40MS0uMDktLjQxLjA5em0xMy4wOSAzLjk1bC4yOC4wNi4zMi4xMi4zNS4xOC4zNi4yNy4zNi4zNS4zNS40Ny4zMi41OS4yOC43My4yMS44OC4xNCAxLjA0LjA1IDEuMjMtLjA2IDEuMjMtLjE2IDEuMDQtLjI0Ljg2LS4zMi43MS0uMzYuNTctLjQuNDUtLjQyLjMzLS40Mi4yNC0uNC4xNi0uMzYuMDktLjMyLjA1LS4yNC4wMi0uMTYtLjAxaC04LjIydi44Mmg1Ljg0bC4wMSAyLjc2LjAyLjM2LS4wNS4zNC0uMTEuMzEtLjE3LjI5LS4yNS4yNS0uMzEuMjQtLjM4LjItLjQ0LjE3LS41MS4xNS0uNTguMTMtLjY0LjA5LS43MS4wNy0uNzcuMDQtLjg0LjAxLTEuMjctLjA0LTEuMDctLjE0LS45LS4yLS43My0uMjUtLjU5LS4zLS40NS0uMzMtLjM0LS4zNC0uMjUtLjM0LS4xNi0uMzMtLjEtLjMtLjA0LS4yNS0uMDItLjIuMDEtLjEzdi01LjM0bC4wNS0uNjQuMTMtLjU0LjIxLS40Ni4yNi0uMzguMy0uMzIuMzMtLjI0LjM1LS4yLjM1LS4xNC4zMy0uMS4zLS4wNi4yNi0uMDQuMjEtLjAyLjEzLS4wMWg1Ljg0bC42OS0uMDUuNTktLjE0LjUtLjIxLjQxLS4yOC4zMy0uMzIuMjctLjM1LjItLjM2LjE1LS4zNi4xLS4zNS4wNy0uMzIuMDQtLjI4LjAyLS4yMVY2LjA3aDIuMDlsLjE0LjAxem0tNi40NyAxNC4yNWwtLjIzLjMzLS4wOC40MS4wOC40MS4yMy4zMy4zMy4yMy40MS4wOC40MS0uMDguMzMtLjIzLjIzLS4zMy4wOC0uNDEtLjA4LS40MS0uMjMtLjMzLS4zMy0uMjMtLjQxLS4wOC0uNDEuMDh6IiBmaWxsPSIjRkZGRkZGIiAvPjwvZz48dGV4dCB4PSIxNDguNTU1MTk5OTk5OTk5OTkiIHk9IjE3LjUiIGR5PSIwLjM1ZW0iIGZvbnQtc2l6ZT0iMTIiIGZvbnQtZmFtaWx5PSJNb250c2VycmF0LCBzYW5zLXNlcmlmIiBmaWxsPSIjRkZGRkZGIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBsZXR0ZXItc3BhY2luZz0iMiIgZm9udC13ZWlnaHQ9IjkwMCIgZmlsbC1vcGFjaXR5PSIxIiBzdHlsZT0idGV4dC10cmFuc2Zvcm06IHVwcGVyY2FzZSI+UFlUSE9OPC90ZXh0Pjwvc3ZnPg==)](https://fastapi.tiangolo.com/)

miruzo-core is the FastAPI/SQLModel backend that powers the miruzo photo
archive. It serves REST APIs for image browsing, scoring, and favorites, and
ships an importer pipeline that ingests gataku assets into SQLite or PostgreSQL.


## ✨ Features

- REST API for listing images, fetching contexts, and patching scores/favorites
- Importer tooling to populate the database and generate thumbnails/variants
- SQLite (default) and PostgreSQL repositories with shared business logic
- Pure helper modules for variant normalization and query parsing
- OpenAPI metadata automatically generated via Pydantic models


## 📏 Why Manbytes?

Human-friendly size units are not universal — they are shaped by culture and
convention.

miruzo defines **manbytes**, a size unit based on a 10<sup>4</sup>-byte scale
(the Japanese “man” unit), and uses it in its API to represent image file sizes.

While inspired by the Japanese numeric system, manbytes is designed for
practical use: it provides a stable, human-scaled representation that works well
for image-heavy UIs and delivery strategies.

For a detailed explanation of the rationale, design trade-offs, and exact
definitions, see [`docs/unit.md`](./docs/unit.md).


## 🚀 Setup

Run miruzo-core locally by following these steps.

### Requirements
- Python 3.13 (respect [`.python-version`](./.python-version); keep code 3.10+
  compatible)
- Git
- Docker (only if you need to run the PostgreSQL repository tests)

### Steps
1. Clone and install dependencies  
   `git clone https://github.com/mntone/miruzo-core.git && cd miruzo-core && pip install -r requirements.txt`
2. Copy the sample environment file  
   `cp .env.development .env` and adjust paths/DSNs (see [AGENTS.md](./AGENTS.md)).
3. Start the API  
   `uvicorn app.main:app --reload`
4. (Optional) Run importer help to confirm configuration  
   `python importers/import.py --help`

### Common commands
- `pytest`: Run the default test suite (SQLite, service, variants, etc.)
- `pytest tests/services/images/repository/test_postgre.py`: Run Docker-backed
  PostgreSQL repository tests (pulls `postgres:18-alpine`)
- `ruff check app tests`: Lint the codebase


## 🖱️ Usage

1. Start the API (`uvicorn app.main:app --reload`).
2. Point miruzo-web or other clients to the running host (default
   `http://127.0.0.1:1024/api`).
3. Use importer commands to populate the database:

   ```bash
   python importers/import.py --jsonl path/to/data.jsonl --static-dir ./static
   ```

4. Hit `/api/images` or `/api/images/{id}` to verify data is available.


## ⚙️ Configuration

All configuration flows through [`.env.development`](./.env.development) using
`pydantic-settings`. Key variables include:

- `DATABASE_BACKEND`: `sqlite` (default) or `postgres`
- `DATABASE_URL`: SQLAlchemy DSN (e.g., `sqlite:///miruzo.db`,
  `postgresql://user:pass@host/db`)
- `GATAKU_ROOT` / `ASSETS_ROOT` / `STATIC_ROOT`: filesystem roots for importer +
  static files
- `VARIANT_LAYERS`: loaded from `app/core/variant_config.py`


## 📦 Database support

- **SQLite**: default development backend; in-memory mode used for most tests
- **PostgreSQL**: optional backend for production; Docker test suite available


## 📜 License

This project is licensed under the terms of the GNU General Public License v3.0
(GPLv3). See the [`LICENSE`](./LICENSE) file for full details. You are free to
use, modify, and distribute this software under the terms of the GPL, provided
that any derivative work is also distributed under the same license.


## 🤝 Contributing

Interested in contributing? See [`CONTRIBUTING.md`](./CONTRIBUTING.md) and
[`AGENTS.md`](./AGENTS.md).


## 🔗 Related Projects

- [miruzo-web](https://github.com/mntone/miruzo-web) — Solid.js frontend that consumes the core APIs
- [gataku](https://github.com/mntone/gataku) — Source asset repository used by the importer


## 👤 Contact

miruzo-core is developed and maintained by *mntone*.

- GitHub: https://github.com/mntone
- Mastodon: https://mstdn.jp/@mntone
- X: https://x.com/mntone
