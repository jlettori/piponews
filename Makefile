build:
	go build -ldflags="-X piponews/internal/handlers.BuildVersion=$$(date -u +%Y%m%d%H%M%S)" -o piponews .

run:
	go run -ldflags="-X piponews/internal/handlers.BuildVersion=$$(date -u +%Y%m%d%H%M%S)" .
