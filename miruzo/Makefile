#
# miruzo Makefile
#
# Examples:
#   make clean
#   make build
#   make dev
#   make generate
#

TOOLS_HINT := command found. Install with: make tools (Linux) or brew install

.PHONY: \
	build clean dev generate release tools \
	air-check sqlc-check goreleaser-check

# -----------------------------------------------------------------------------
# Tool checks
# -----------------------------------------------------------------------------

air-check:
	@command -v air >/dev/null 2>&1 || (echo "No air $(TOOLS_HINT) go-air (macOS)."; exit 1)

sqlc-check:
	@command -v sqlc >/dev/null 2>&1 || (echo "No sqlc $(TOOLS_HINT) sqlc (macOS)."; exit 1)

goreleaser-check:
	@command -v goreleaser >/dev/null 2>&1 || (echo "No goreleaser $(TOOLS_HINT) goreleaser (macOS)."; exit 1)

# -----------------------------------------------------------------------------
# Core commands
# -----------------------------------------------------------------------------

build:
	go build -gcflags="all=-N -l" -tags=dev -o ./bin/miruzo-api ./cmd/miruzo-api
	go build -gcflags="all=-N -l" -tags=dev -o ./bin/miruzo-cli ./cmd/miruzo-cli

clean:
	rm -rf ./bin
	rm -rf ./internal/database/mysql/gen    ./internal/database/mysql/migrations_min
	rm -rf ./internal/database/postgres/gen ./internal/database/postgres/migrations_min
	rm -f ./internal/database/sqlite/gen/db.go
	rm -f ./internal/database/sqlite/gen/models.go
	rm -f ./internal/database/sqlite/gen/*.sql.go
	rm -rf ./internal/database/sqlite/migrations_min

dev: air-check
	air

generate: sqlc-check
	cd ./internal/database/mysql/migrations && go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=mysql
	cd ./internal/database/postgres/migrations && go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=postgres
	cd ./internal/database/sqlite/migrations && go run ../../../../../tools/sql_minify/main.go ../migrations ../migrations_min --dialect=sqlite
	cd ./internal/database && sqlc generate

release: goreleaser-check
	goreleaser release --snapshot --clean

tools:
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install golang.org/x/tools/gopls@latest
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/goreleaser/goreleaser/v2@latest
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
