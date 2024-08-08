

ENTRY_BUILD_FILE=./cmd/.

BINARY := short-link

up:
	docker compose up

down:
	./stop-gracefully.sh

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

migration_create:
	migrate create -ext sql -dir database/migration/ -seq init_mg

migration_up:
	migrate -path database/migration/ -database "postgresql://postgres:postgres@postgres_db:5432/slink?sslmode=disable" -verbose up

migration_down:
	migrate -path database/migration/ -database "postgresql://default:secret@localhost:5432/slink?sslmode=disable" -verbose down 1

migration_fix:
	migrate -path database/migration/ -database "postgresql://default:secret@localhost:5432/slink?sslmode=disable" force 1

migration_up_v2:
	docker-compose -f docker-compose.yml run --rm app 	migrate -path database/migration/ -database "postgresql://postgres:postgres@postgres_db:5432/slink?sslmode=disable" -verbose up
