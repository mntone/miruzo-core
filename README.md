# miruzo-core

[![License under GPLv3](https://forthebadge.com/api/badges/generate?panels=2&primaryLabel=LICENSE&secondaryLabel=GPL%203%2B&primaryBGColor=%23555555&primaryTextColor=%23FFFFFF&secondaryBGColor=%23007ec6&secondaryTextColor=%23FFFFFF&primaryFontSize=12&primaryFontWeight=300&primaryLetterSpacing=2&primaryFontFamily=Montserrat&primaryTextTransform=uppercase&secondaryFontSize=12&secondaryFontWeight=900&secondaryLetterSpacing=2&secondaryFontFamily=Montserrat&secondaryTextTransform=uppercase&secondaryIconColor=%23FFFFFF&secondaryIconSize=24&secondaryIconPosition=right)](./LICENSE)
[![Made with FastAPI](https://forthebadge.com/api/badges/generate?panels=2&primaryLabel=MADE+WITH&secondaryLabel=Fast+API&primaryBGColor=%23ef4041&primaryTextColor=%23FFFFFF&secondaryBGColor=%23c1282d&secondaryTextColor=%23FFFFFF&primaryFontSize=12&primaryFontWeight=300&primaryLetterSpacing=2&primaryFontFamily=Montserrat&primaryTextTransform=uppercase&secondaryFontSize=12&secondaryFontWeight=900&secondaryLetterSpacing=2&secondaryFontFamily=Montserrat&secondaryTextTransform=uppercase&secondaryIcon=fastapi&secondaryIconColor=%23FFFFFF&secondaryIconSize=24&secondaryIconPosition=right)](https://fastapi.tiangolo.com/)
[![Written by Python](https://forthebadge.com/api/badges/generate?panels=2&primaryLabel=WRITTEN+BY&secondaryLabel=Python&primaryBGColor=%238fc965&primaryTextColor=%23FFFFFF&secondaryBGColor=%23419b5a&secondaryTextColor=%23FFFFFF&primaryFontSize=12&primaryFontWeight=300&primaryLetterSpacing=2&primaryFontFamily=Montserrat&primaryTextTransform=uppercase&secondaryFontSize=12&secondaryFontWeight=900&secondaryLetterSpacing=2&secondaryFontFamily=Montserrat&secondaryTextTransform=uppercase&secondaryIcon=python&secondaryIconColor=%23FFFFFF&secondaryIconSize=24&secondaryIconPosition=right)](https://www.python.org/)

miruzo-core is the FastAPI/SQLModel backend that powers the miruzo photo
archive. It serves REST APIs for image browsing, scoring, and love actions, and
ships an importer pipeline that ingests gataku assets into SQLite or PostgreSQL.


## ✨ Features

- REST API for listing images, fetching contexts, and posting love actions
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
   `python importers/gataku_import.py --help`

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
   python importers/gataku_import.py --jsonl path/to/data.jsonl
   ```

4. Hit `/api/i/latest` or `/api/i/{ingest_id}` to verify data is available.


## ⚙️ Configuration

All configuration flows through [`.env.development`](./.env.development) using
`pydantic-settings`. Key variables include:

- `DATABASE_BACKEND`: `sqlite` (default) or `postgres`
- `DATABASE_URL`: SQLAlchemy DSN (e.g., `sqlite:///var/miruzo.sqlite`,
  `postgresql://user:pass@host/db`)
- `GATAKU_ROOT` / `GATAKU_ASSETS_ROOT` / `MEDIA_ROOT`: filesystem roots for
  importer + media files
- `VARIANT_LAYERS`: loaded from `app/config/variant.py`


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
