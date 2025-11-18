default:
	@just --list

build: frontend
	@echo "Building OneOff..."
	go build -o oneoff ./cmd/oneoff

run: build
	./oneoff

dev-frontend:
	@echo "Starting frontend in development mode..."
	bun start

dev:
	air

clean:
	rm -f oneoff
	rm -f *.db *.db-shm *.db-wal
	rm -rf dist/
	rm -rf node_modules/

test:
	go test -v ./...

migrate:
	go run ./cmd/oneoff migrate

frontend:
	@echo "Building frontend..."
	if [ -d "node_modules" ]; then \
		bun run build; \
	else \
		echo "Frontend not set up yet, skipping..."; \
		mkdir -p internal/server/dist; \
	fi

deps:
	go mod download
	go mod tidy

setup: deps
	@echo "Installing frontend dependencies..."
	bun install
