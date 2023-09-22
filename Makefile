
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
