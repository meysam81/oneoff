# OneOff - Modern One-Time Job Scheduler

OneOff is a self-hosted, developer-focused job scheduler for executing one-time tasks at specific future times. Think "Linux `at` command meets modern web UI" with a plugin-based architecture.

## Features

- **Single Binary Deployment**: Zero external dependencies, everything bundled into one executable
- **SQLite Database**: Minimal operational complexity with built-in migrations
- **Modern Vue 3 UI**: Beautiful, dark-mode-first interface built with Naive UI
- **Multiple Job Types**: HTTP requests, shell scripts, Docker containers
- **Priority Queue System**: 1-10 priority levels for job execution
- **Project Organization**: Group and manage jobs by project
- **Tag System**: Categorize jobs with customizable tags
- **Real-time Monitoring**: Live worker status and job execution tracking
- **Full Observability**: Execution logs, duration tracking, error reporting
- **Job Chaining**: Create sequences of jobs that execute in order

## Quick Start

### Prerequisites

- Go 1.23+ (for building from source)
- Node.js 18+ (for building frontend)
- Docker (optional, for Docker job type)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/meysam81/oneoff.git
cd oneoff
```

2. Build the application:

```bash
make setup  # Install dependencies
make build  # Build frontend and backend
```

3. Run OneOff:

```bash
./oneoff
```

The application will be available at `http://localhost:8080`

### Configuration

OneOff is configured via environment variables (12-factor app compliant):

```bash
# Server configuration
PORT=8080
HOST=localhost

# Database
DB_PATH=./oneoff.db

# Workers
WORKERS_COUNT=0  # 0 = auto-detect (N/2 cores)

# Logging
LOG_LEVEL=info  # debug, info, warn, error

# Defaults
DEFAULT_TIMEZONE=UTC
LOG_RETENTION_DAYS=90
DEFAULT_PRIORITY=5

# Environment
ENVIRONMENT=production
```

## Usage

### Creating a Job

1. Navigate to the Dashboard
2. Click "Create Job"
3. Fill in the job details:
   - **Name**: Descriptive name for your job
   - **Type**: HTTP, Shell, or Docker
   - **Scheduled Time**: When to execute (datetime picker)
   - **Priority**: 1-10 (higher = more important)
   - **Project**: Organize by project
   - **Tags**: Add categorization tags
4. Configure job-specific settings
5. Click "Create Job"

### Job Types

#### HTTP Job

Execute HTTP requests at scheduled times:

```json
{
  "url": "https://api.example.com/webhook",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer token"
  },
  "body": "{\"event\": \"scheduled\"}",
  "timeout": 30
}
```

#### Shell Job

Run shell scripts or commands:

```json
{
  "script": "#!/bin/bash\\necho 'Hello World'",
  "is_path": false,
  "args": [],
  "env": {
    "VAR": "value"
  },
  "workdir": "/tmp",
  "timeout": 60
}
```

#### Docker Job

Execute Docker containers:

```json
{
  "image": "alpine:latest",
  "command": ["echo", "hello"],
  "env": {
    "ENV_VAR": "value"
  },
  "volumes": {
    "/host/path": "/container/path"
  },
  "auto_remove": true,
  "timeout": 300
}
```

## Development

### Project Structure

```
oneoff/
├── main.go          # CLI entry point
├── internal/
│   ├── config/          # Configuration
│   ├── domain/          # Domain models and interfaces
│   ├── repository/      # Database layer
│   ├── service/         # Business logic
│   ├── handler/         # HTTP handlers
│   ├── worker/          # Worker pool
│   ├── jobs/            # Job executors
│   └── server/          # HTTP server
├── migrations/          # Database migrations
├── src/                 # Vue 3 frontend
│   ├── components/      # Vue components
│   ├── views/           # Page views
│   ├── stores/          # Pinia stores
│   ├── utils/           # Utilities
│   └── assets/          # Static assets
└── Makefile            # Build commands
```

### Building

```bash
# Build everything
make build

# Build frontend only
make frontend

# Build backend only
go build -o oneoff .

# Run in development mode
make dev

# Run frontend in dev mode (with hot reload)
bun start
```

### Running Migrations

```bash
./oneoff migrate --direction up
./oneoff migrate --direction down
```

## API Documentation

### Jobs

- `GET /api/jobs` - List jobs
- `POST /api/jobs` - Create job
- `GET /api/jobs/:id` - Get job details
- `PATCH /api/jobs/:id` - Update job
- `DELETE /api/jobs/:id` - Delete job
- `POST /api/jobs/:id/execute` - Execute job immediately
- `POST /api/jobs/:id/clone` - Clone job
- `POST /api/jobs/:id/cancel` - Cancel job

### Executions

- `GET /api/executions` - List executions
- `GET /api/executions/:id` - Get execution details

### Projects

- `GET /api/projects` - List projects
- `POST /api/projects` - Create project
- `PATCH /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

### Tags

- `GET /api/tags` - List tags
- `POST /api/tags` - Create tag
- `PATCH /api/tags/:id` - Update tag
- `DELETE /api/tags/:id` - Delete tag

### System

- `GET /api/system/status` - Get system stats
- `GET /api/system/config` - Get configuration
- `PATCH /api/system/config` - Update configuration
- `GET /api/workers/status` - Get worker status
- `GET /api/job-types` - Get available job types

## Architecture

### Backend (Go)

- **CLI**: urfave/cli/v3
- **Config**: caarlos0/env/v11
- **HTTP Client**: imroc/req/v3
- **Database**: SQLite with golang-migrate
- **Logging**: zerolog

### Frontend (Vue 3)

- **Framework**: Vue 3 with Composition API
- **Build Tool**: Vite (optimized bundling)
- **UI Components**: Naive UI (tree-shakable)
- **HTTP Client**: ky (with retry)
- **State Management**: Pinia
- **Routing**: Vue Router

### Key Design Decisions

1. **Single Binary**: Frontend embedded into Go binary for easy deployment
2. **SQLite**: No external database required, perfect for single-node deployments
3. **Worker Pool**: Configurable concurrency with graceful shutdown
4. **Job Registry**: Plugin-based architecture for extensible job types
5. **12-Factor App**: All configuration via environment variables
6. **Fail Fast**: Missing required configuration fails at startup

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Support

For issues and questions:

- GitHub Issues: https://github.com/meysam81/oneoff/issues
