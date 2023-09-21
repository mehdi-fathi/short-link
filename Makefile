
doc:
	godoc -index

version:
	@go version

db:
		docker run -d -p 5433:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_USER=root -e POSTGRES_DB=app --name aparat-dsp_db  postgres:12.2 && echo "sleeping 10 seconds to make database ready" && sleep 10 && docker exec aparat-dsp_db sh -c "createdb app_test"
run:
	go run cmd/*.go

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out