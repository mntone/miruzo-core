# miruzo-core

[![License under GPLv3](https://forthebadge.com/api/badges/generate?panels=2&primaryLabel=LICENSE&secondaryLabel=GPL%203%2B&primaryBGColor=%23555555&primaryTextColor=%23FFFFFF&secondaryBGColor=%23007ec6&secondaryTextColor=%23FFFFFF&primaryFontSize=12&primaryFontWeight=300&primaryLetterSpacing=2&primaryFontFamily=Montserrat&primaryTextTransform=uppercase&secondaryFontSize=12&secondaryFontWeight=900&secondaryLetterSpacing=2&secondaryFontFamily=Montserrat&secondaryTextTransform=uppercase&secondaryIconColor=%23FFFFFF&secondaryIconSize=24&secondaryIconPosition=right)](./LICENSE)
[![Written by Go](https://forthebadge.com/api/badges/generate?panels=2&primaryLabel=WRITTEN+BY&secondaryLabel=Go&primaryBGColor=%238fc965&primaryTextColor=%23FFFFFF&secondaryBGColor=%23419b5a&secondaryTextColor=%23FFFFFF&primaryFontSize=12&primaryFontWeight=300&primaryLetterSpacing=2&primaryFontFamily=Montserrat&primaryTextTransform=uppercase&secondaryFontSize=12&secondaryFontWeight=900&secondaryLetterSpacing=2&secondaryFontFamily=Montserrat&secondaryTextTransform=uppercase&secondaryIcon=go&secondaryIconColor=%23FFFFFF&secondaryIconSize=24&secondaryIconPosition=right)](https://go.dev/)
[![Written by Python](https://forthebadge.com/api/badges/generate?panels=2&primaryLabel=WRITTEN+BY&secondaryLabel=Python&primaryBGColor=%238fc965&primaryTextColor=%23FFFFFF&secondaryBGColor=%23419b5a&secondaryTextColor=%23FFFFFF&primaryFontSize=12&primaryFontWeight=300&primaryLetterSpacing=2&primaryFontFamily=Montserrat&primaryTextTransform=uppercase&secondaryFontSize=12&secondaryFontWeight=900&secondaryLetterSpacing=2&secondaryFontFamily=Montserrat&secondaryTextTransform=uppercase&secondaryIcon=python&secondaryIconColor=%23FFFFFF&secondaryIconSize=24&secondaryIconPosition=right)](https://www.python.org/)

miruzo-core is the backend and ingest core for the miruzo photo archive.
The API serves image listing/context/love endpoints, and the ingest tooling
imports gataku assets into supported database backends.


## ✨ Features

- Go API for image browsing and reaction workflows
- Python ingest pipeline for importing and processing source assets
- Shared support for SQLite and PostgreSQL
- Optional MySQL support in Python ingest
- Generated SQL access via sqlc for Go repositories


## 📏 Why Manbytes?

Human-friendly size units are not universal and are shaped by culture and
convention.

miruzo defines **manbytes**, a size unit based on a 10<sup>4</sup>-byte scale
(the Japanese “man” unit), and uses it in API responses for image file
sizes.

For details, see [`docs/unit.md`](./docs/unit.md).


## 🧩 Repository Layout

- [`miruzo`](./miruzo): Go API application, SQL/migrations, Makefile
- [`miruzo-py`](./miruzo-py): Python ingest services, DB adapters, tests


## 🚀 Requirements

- Git
- Go 1.26.x (see [`miruzo/go.mod`](./miruzo/go.mod))
- Python 3.10+ (3.13 recommended for local development)

Backend/runtime matrix:

| Backend    | Version                         | Go driver                       | Go API    | Python driver             |
| ---------- | ------------------------------- | ------------------------------- | --------- | ------------------------- |
| MySQL      | 8.0.16+ (`CHECK`)               | `go-sql-driver/mysql` (planned) | not yet   | `mysqlclient` (`MySQLdb`) |
| PostgreSQL | 14+                             | `jackc/pgx/v5`                  | supported | `psycopg3` (`psycopg`)    |
| SQLite     | 3.37.0+ (`RETURNING`, `STRICT`) | `mattn/go-sqlite3`              | supported | `sqlite3` (stdlib)        |


## 🛠️ Setup

For source development setup, follow
[`CONTRIBUTING.md`](./CONTRIBUTING.md#prerequisites):

- prerequisites and backend version requirements
- Linux/macOS setup for Go tools and Python DB drivers
- runtime configuration and test environment variables

Prebuilt binary releases are planned. Until then, use source setup.


## ⚙️ Configuration

### Go API (`miruzo/config.yaml`)

- Base file: `miruzo/internal/app/config.sample.yaml`
- Local default file: [`miruzo/config.yaml`](./miruzo/config.yaml)
- Current Go API backends: `sqlite`, `postgres`
- `database.backend=mysql` is not supported yet

### Python ingest (`miruzo-py/.env`)

- Copy [`miruzo-py/.env.development`](./miruzo-py/.env.development) to
  `miruzo-py/.env`
- Set these variables explicitly:
  - `ENVIRONMENT` (`development` or `production`)
  - `DATABASE_BACKEND` (`sqlite`, `postgres`, or `mysql`)
  - `DATABASE_URL`
    - SQLite: `sqlite:///...`
    - PostgreSQL: `postgresql+psycopg://...`
    - MySQL: `mysql+mysqldb://...`
- Path-related variables (`MEDIA_ROOT`, `PUBLIC_MEDIA_ROOT`, `GATAKU_ROOT`,
  `GATAKU_ASSETS_ROOT`, `GATAKU_SYMLINK_DIRNAME`) can be left as defaults on
  first setup, then customized only when needed.


## 🖱️ Run Locally

Start API:

```bash
cd miruzo
make dev
# or
cd miruzo && go run ./cmd/miruzo-api
```

Default API address: `http://127.0.0.1:1360/api`

Run importer help:

```bash
cd miruzo-py
python -m scripts.gataku_import --help
```


## 🧪 Testing

Default suites:

```bash
cd miruzo && go test ./...
cd miruzo-py && pytest
```

Focused suites:

```bash
cd miruzo && go test ./internal/service/...
cd miruzo && go test ./internal/adapter/persistence/contract/...
cd miruzo-py && pytest tests/importers
cd miruzo-py && pytest tests/persist
```

Optional test database DSN environment variables:

- Go tests: `MIRUZO_TEST_POSTGRES_URL`
- Python tests:
  - `MIRUZO_PY_TEST_MYSQL_URL`
  - `MIRUZO_PY_TEST_POSTGRES_URL`


## 📜 License

This project is licensed under GNU GPLv3. See [`LICENSE`](./LICENSE).


## 🤝 Contributing

See:

- [`AGENTS.md`](./AGENTS.md)
- [`CONTRIBUTING.md`](./CONTRIBUTING.md)


## 🔗 Related Projects

- [miruzo-web](https://github.com/mntone/miruzo-web) — Solid.js frontend that consumes the core APIs
- [gataku](https://github.com/mntone/gataku) — Source asset repository used by the importer


## 👤 Contact

miruzo-core is developed and maintained by *mntone*.

- GitHub: https://github.com/mntone
- Mastodon: https://mstdn.jp/@mntone
- X: https://x.com/mntone
