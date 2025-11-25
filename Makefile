.PHONY: build run dev clean test migrate frontend

# Build binary
build: frontend
	@echo "Building OneOff..."
	go build -o oneoff .

# Run the application
run: build
	./oneoff

# Development mode (no frontend build)
dev:
	go run .

# Clean build artifacts
clean:
	rm -f oneoff
	rm -f *.db *.db-shm *.db-wal
	rm -rf dist/
	rm -rf node_modules/

# Run tests
test:
	go test -v ./...

# Run database migrations
migrate:
	go run . migrate

# Build frontend
frontend:
	@echo "Building frontend..."
	@if [ -d "node_modules" ]; then \
		npm run build; \
	else \
		echo "Frontend not set up yet, skipping..."; \
		mkdir -p internal/server/dist; \
	fi

# Install dependencies
deps:
	go mod download
	go mod tidy

# Development setup
setup: deps
	@echo "Installing frontend dependencies..."
	npm install
