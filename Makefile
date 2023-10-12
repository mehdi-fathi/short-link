
doc:
	godoc -index

version:
	@go version

run:
	go run cmd/*.go

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out


fmt:
	go fmt ./...

migration_up:
	migrate -path database/migration/ -database "postgresql://default:secret@localhost:5432/slink?sslmode=disable" -verbose up

migration_down:
	migrate -path database/migration/ -database "postgresql://default:secret@localhost:5432/slink?sslmode=disable" -verbose down

migration_fix:
	migrate -path database/migration/ -database "postgresql://default:secret@localhost:5432/slink?sslmode=disable" force 1