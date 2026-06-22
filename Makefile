build:
	go build -ldflags="-X piponews/internal/handlers.BuildVersion=$$(date -u +%Y%m%d%H%M%S)" -o piponews .

run:
	go run -ldflags="-X piponews/internal/handlers.BuildVersion=$$(date -u +%Y%m%d%H%M%S)" .

test:
	go test -v -count=1 ./...

coverage:
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"
	@go tool cover -func=coverage.out | tail -1

clean:
	rm -f piponews piponews.db coverage.out coverage.html

migrate/new:
	go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations create $(filter-out $@,$(MAKECMDGOALS)) sql

migrate/status:
	go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations sqlite3 piponews.db status

migrate/up:
	go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations sqlite3 piponews.db up

migrate/down:
	go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations sqlite3 piponews.db down
