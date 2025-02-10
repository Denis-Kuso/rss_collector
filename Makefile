# include environment variables from config file
include .env

## help: display available targets/commands
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ##############################################################################
# GUARDRAILS
# ##############################################################################
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## no-uncommited: test whether there are uncommited changes
.PHONY: no-uncommited
no-uncommited:
	@test -z "$(shell git status --porcelain)"

## psql-connect: can you connect to test db
.PHONY: psql-connect
psql-connect:
	psql ${TEST_DB_DSN}

# ##############################################################################
# DEVELOP
# ##############################################################################
binary_srv_name = rssd
source_srv = "./cmd/api"
source_cli = "./cmd/cli"
binary_cli_name = clipr
cli_version="0.1.0"
commit_hash = $(shell git describe --always --dirty)
migration_dir = ./migrations

## buildSrv: build the server (binaries)
.PHONY: buildSrv
buildSrv:
	@echo "building server binaries..."
	go build -ldflags='-X main.version=dev-${commit_hash}' -o=/tmp/bin/${binary_srv_name} ${source_srv}

## buildCLI: build the client (binaries)
.PHONY: buildCLI
buildCLI:
	@echo "building CLI binaries..."
	go build -ldflags='-X main.version=${cli_version} -X main.commit=dev-${commit_hash}' -o=/tmp/bin/${binary_cli_name} ${source_cli}

## run: run the server
.PHONY: run
run: buildSrv migrate/up
	/tmp/bin/${binary_srv_name} -db-dsn=${TEST_DB_DSN} &

## test: run all the tests
.PHONY: test
test:
	@echo "starting all the tests..."

## migrate/up: run any pending migrations
.PHONY: migrate/up
migrate/up:
	@echo "running migrations..."
	goose -dir $(migration_dir) postgres "${TEST_DB_DSN}" up

# ##############################################################################
# DEPLOYMENT
# ##############################################################################
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.version=${git_description}'

## codeCheck: format, vet and lint code
.PHONY: codeCheck
codeCheck: test
	@echo "- format code"
	go fmt ./...
	@echo "- vet code"
	go vet ./...
	@echo "- verbosely lint"
	revive -formatter friendly

## dependencies: verify dependencies
.PHONY: dependencies
dependencies:
	@echo "- dependencies check"
	go mod verify
	go mod tidy

## production/lint: pass/fail on linting issues
.PHONY: production/lint
production/lint:
	revive -set_exit_status || exit 1

## push: push changes to remote repo
.PHONY: push
push: confirm no-uncommited codeCheck production/lint
	git push

## production/build: build the binaries for deployment
.PHONY: production/build
production/build: confirm no-uncommited codeCheck production/lint
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=/tmp/bin/linux_amd64/${binary_srv_name} ${source_srv}
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=/tmp/bin/linux_amd64/${binary_cli_name} ${source_cli}
