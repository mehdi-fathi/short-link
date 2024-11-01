
ENTRY_BUILD_FILE=./cmd/.

.PHONY: up open build-docker create_kind_cluster down doc run build run_build fmt migration_create migration_up migration_down migration_fix migration_up_v2 create_test_db

BINARY := short-link

DB_HOST = postgres
DB_USER = postgres
DB_PASSWORD = postgres
DB_PORT = 5432
DB_NAME = slink

# Define variables
DOCKER_COMPOSE = docker compose -f deploy/dev/docker-compose.yml --env-file .env.local
DB_URL = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_PATH = -path database/migration/
MIGRATE_CMD = run --rm app migrate $(MIGRATION_PATH) -database "$(DB_URL)"

up:
	$(DOCKER_COMPOSE) up

open:
	./deploy/dev/start-and-open.sh

build-docker:
	$(DOCKER_COMPOSE) build

down:
	./deploy/dev/stop-gracefully.sh

doc:
	godoc -index

version:
	@go version

run:
	go run cmd/*.go

build:
	go build -o ./bin/$(BINARY) $(ENTRY_BUILD_FILE) && chmod +x bin/$(BINARY)

run_build:
	bin/$(BINARY)

test:
	go test ./... -v

test_coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

fmt:
	go fmt ./...

# Targets
migration_create:
	$(DOCKER_COMPOSE) $(MIGRATE_CMD) create -ext sql -dir database/migration/ -seq init_mg

migration_up:
	$(DOCKER_COMPOSE) $(MIGRATE_CMD) -verbose up

migration_down:
	$(DOCKER_COMPOSE) $(MIGRATE_CMD) -verbose down 1

migration_fix:
	$(DOCKER_COMPOSE) $(MIGRATE_CMD) force 1

migration_up_v2:
	$(DOCKER_COMPOSE) $(MIGRATE_CMD) -verbose up

create_test_db:
	$(DOCKER_COMPOSE) exec db psql -U $(DB_USER) -c "CREATE DATABASE $(DB_TEST_NAME);"

build_docker:
	docker build -f deploy/kuber/Dockerfile -t go_app .
