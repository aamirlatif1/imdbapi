
include .env
export

# clean: cleans generated files
.PHONY: clean
clean:
	@rm -Rf bin/
	@rm -Rf pgdata/
	@rm -f coverage.txt

## build: builds the binary
.PHONY: build
build:
	go build -o ./bin/imdb ./cmd/api

## test: run tests
.PHONY: test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

## migration_status: checks the database migration status
.PHONY: migration_status
migration_status:
	docker compose run migrator status

## migration: creates a new database migration. Usage: make migration name=test
.PHONY: migration
migration:
	docker compose run migrator create $(name) sql

## migrate: runs the the migrations
.PHONY: migrate
migrate:
	docker compose run migrator up

## migrate_down: downs the migration
.PHONY: migrate_down
migrate_down:
	docker compose run migrator down

## docker_up: boostrap the environment via Docker
.PHONY: docker_up
docker_up:
	@docker compose up --build

.PHONY: help
help : Makefile
	@sed -n 's/^##//p' $<
	@echo ""