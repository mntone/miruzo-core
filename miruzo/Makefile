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

VERSION := $(shell git describe --tags --abbrev=0 --match 'v[0-9]*.[0-9]*.[0-9]*' 2>/dev/null || echo 0.1.0)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_DIRTY := $(shell test -n "$$(git status --porcelain)" && echo +dev)
APP_VERSION := $(VERSION)-$(GIT_COMMIT)$(GIT_DIRTY)

BUILD_DIR ?= ./bin
GO ?= go
COMMON_FLAGS := -Og -g
COMMON_GO_FLAGS := -gcflags="github.com/mntone/miruzo-core/miruzo/...=-N -l" -buildmode=pie -tags='dev libsqlite3'
COMMON_GO_LDFLAGS := -X main.version=$(APP_VERSION)

ifeq ($(shell uname -s),Darwin)
	COMMON_FLAGS += -fstack-protector-strong
	COMMON_LDFLAGS := -Wl,-dead_strip
else
	COMMON_FLAGS += -fno-plt -fstack-clash-protection -fstack-protector-strong -D_FORTIFY_SOURCE=3
	COMMON_LDFLAGS := -Wl,-z,relro,-z,now,-z,noexecstack
endif

define build_debug
	CGO_ENABLED=1 \
	CGO_CFLAGS="$(COMMON_FLAGS)" \
	CGO_CXXFLAGS="$(COMMON_FLAGS)" \
	CGO_FFLAGS="$(COMMON_FLAGS)" \
	CGO_LDFLAGS="$(COMMON_LDFLAGS)" \
	$(GO) build \
		$(COMMON_GO_FLAGS) \
		-ldflags '$(COMMON_GO_LDFLAGS)' \
		-o $(BUILD_DIR)/$(1) \
		./cmd/$(1)
endef

.PHONY: \
	build-api build-cli build \
	init clean dev generate release tools \
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

build-api:
	$(call build_debug,miruzo-api)

build-cli:
	$(call build_debug,miruzo-cli)

build: build-api build-cli

clean:
	rm -rf ./bin
	rm -rf ./internal/database/mysql/gen    ./internal/database/mysql/migrations_min
	rm -rf ./internal/database/postgres/gen ./internal/database/postgres/migrations_min
	rm -f ./internal/database/sqlite/gen/db.go
	rm -f ./internal/database/sqlite/gen/models.go
	rm -f ./internal/database/sqlite/gen/*.sql.go
	rm -rf ./internal/database/sqlite/migrations_min
	go clean

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

init: tools
	go mod download
