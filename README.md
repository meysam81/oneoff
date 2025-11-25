<p align="center">
  <img src="assets/logo/oneoff-logo.svg" alt="OneOff Logo" width="120" height="120">
</p>

<h1 align="center">OneOff</h1>

<p align="center">
  <strong>One-time job scheduler. Single binary. Zero dependencies.</strong>
</p>

<p align="center">
  Schedule HTTP requests, shell scripts, and Docker containers to run at specific times.<br>
  No Redis. No Postgres. No message queues. Just download and run.
</p>

<p align="center">
  <a href="https://github.com/meysam81/oneoff/releases"><img src="https://img.shields.io/github/v/release/meysam81/oneoff?style=flat-square&color=6366F1" alt="Release"></a>
  <a href="https://github.com/meysam81/oneoff/blob/main/LICENSE"><img src="https://img.shields.io/github/license/meysam81/oneoff?style=flat-square&color=10B981" alt="License"></a>
  <a href="https://github.com/meysam81/oneoff/stargazers"><img src="https://img.shields.io/github/stars/meysam81/oneoff?style=flat-square&color=F59E0B" alt="Stars"></a>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#api">API</a>
</p>

---

## Why OneOff?

Ever needed to schedule a one-time task and found yourself setting up Redis, Postgres, Celery, and a web of services just to send a webhook at 3 PM?

**OneOff is the antidote to over-engineering.**

```bash
# This is all you need
./oneoff
```

That's it. Open `localhost:8080`, schedule your job, done.

---

## Quick Start

### 30-Second Setup

```bash
# Download (Linux)
curl -fsSL https://github.com/meysam81/oneoff/releases/latest/download/oneoff_Linux_x86_64.tar.gz | tar xz

# Run
./oneoff

# Open http://localhost:8080
```

<details>
<summary><strong>macOS</strong></summary>

```bash
# Intel Mac
curl -fsSL https://github.com/meysam81/oneoff/releases/latest/download/oneoff_Darwin_x86_64.tar.gz | tar xz

# Apple Silicon
curl -fsSL https://github.com/meysam81/oneoff/releases/latest/download/oneoff_Darwin_arm64.tar.gz | tar xz

./oneoff
```
</details>

<details>
<summary><strong>Windows</strong></summary>

Download from [Releases](https://github.com/meysam81/oneoff/releases), extract, and run:
```powershell
.\oneoff.exe
```
</details>

<details>
<summary><strong>Docker</strong></summary>

```bash
docker run -p 8080:8080 -v oneoff-data:/data ghcr.io/meysam81/oneoff:latest
```
</details>

<details>
<summary><strong>Build from Source</strong></summary>

Requires Go 1.23+ and Node.js 18+:

```bash
git clone https://github.com/meysam81/oneoff.git
cd oneoff
make setup && make build
./oneoff
```
</details>

---

## Features

| Feature | Description |
|---------|-------------|
| **Single Binary** | Everything bundled into one ~15MB executable |
| **SQLite Database** | No external database to manage |
| **Modern Web UI** | Dark-mode-first Vue 3 interface |
| **HTTP Jobs** | Schedule webhooks, API calls, notifications |
| **Shell Jobs** | Run scripts, backups, maintenance tasks |
| **Docker Jobs** | Execute containers on demand |
| **Priority Queue** | 1-10 priority levels for job execution |
| **Projects & Tags** | Organize jobs your way |
| **Real-time Monitoring** | Live worker status and execution tracking |
| **Job Chaining** | Create sequences of dependent jobs |

---

## Usage

### Creating Your First Job

1. Open `http://localhost:8080`
2. Click **"Create Job"**
3. Choose job type, set schedule, configure
4. Done

### Job Types

#### HTTP Job
Send a webhook when a trial expires:
```json
{
  "url": "https://api.yourapp.com/webhooks/trial-expired",
  "method": "POST",
  "headers": {
    "Authorization": "Bearer your-token",
    "Content-Type": "application/json"
  },
  "body": "{\"user_id\": \"12345\"}",
  "timeout": 30
}
```

#### Shell Job
Run a database backup at midnight:
```json
{
  "script": "#!/bin/bash\npg_dump $DATABASE_URL > /backups/$(date +%Y%m%d).sql",
  "env": {
    "DATABASE_URL": "postgres://..."
  },
  "timeout": 300
}
```

#### Docker Job
Run a migration container:
```json
{
  "image": "your-app:latest",
  "command": ["npm", "run", "db:migrate"],
  "env": {
    "NODE_ENV": "production"
  },
  "auto_remove": true,
  "timeout": 600
}
```

---

## Configuration

All configuration via environment variables. Zero config files.

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `HOST` | `localhost` | HTTP server host |
| `DB_PATH` | `./oneoff.db` | SQLite database path |
| `WORKERS_COUNT` | `0` | Worker count (0 = CPU cores / 2) |
| `LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `DEFAULT_TIMEZONE` | `UTC` | Default timezone for jobs |
| `DEFAULT_PRIORITY` | `5` | Default job priority (1-10) |

### Example
```bash
PORT=3000 DB_PATH=/var/lib/oneoff/data.db ./oneoff
```

---

## API

OneOff exposes a REST API for programmatic access.

### Jobs
```bash
# List jobs
curl http://localhost:8080/api/jobs

# Create job
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Send reminder",
    "type": "http",
    "scheduled_at": "2025-01-15T09:00:00Z",
    "config": "{\"url\":\"https://...\",\"method\":\"POST\"}"
  }'

# Execute immediately
curl -X POST http://localhost:8080/api/jobs/{id}/execute

# Cancel job
curl -X POST http://localhost:8080/api/jobs/{id}/cancel
```

### Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/jobs` | List all jobs |
| `POST` | `/api/jobs` | Create job |
| `GET` | `/api/jobs/:id` | Get job details |
| `PATCH` | `/api/jobs/:id` | Update job |
| `DELETE` | `/api/jobs/:id` | Delete job |
| `POST` | `/api/jobs/:id/execute` | Execute now |
| `POST` | `/api/jobs/:id/clone` | Clone job |
| `POST` | `/api/jobs/:id/cancel` | Cancel job |
| `GET` | `/api/executions` | List executions |
| `GET` | `/api/projects` | List projects |
| `GET` | `/api/tags` | List tags |
| `GET` | `/api/system/status` | System stats |
| `GET` | `/api/workers/status` | Worker status |

---

## Use Cases

**Trial Expirations** — Send emails when trials end
**Scheduled Reports** — Generate and email weekly reports
**Database Backups** — Run backups at off-peak hours
**Deployment Triggers** — Schedule deployments for maintenance windows
**Webhook Notifications** — Send Slack/Discord alerts at specific times
**Data Cleanup** — Schedule retention policy enforcement
**API Polling** — Check external services at specific times

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    OneOff Binary                    │
├─────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │   Vue 3 UI  │  │  REST API   │  │   Workers   │ │
│  │  (embedded) │  │  (Go net)   │  │   (pool)    │ │
│  └─────────────┘  └─────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────┤
│                    SQLite Database                  │
└─────────────────────────────────────────────────────┘
```

- **Frontend**: Vue 3 + Naive UI, embedded in binary
- **Backend**: Go standard library, no web framework
- **Database**: SQLite with automatic migrations
- **Workers**: Configurable worker pool with priority queue

---

## Comparison

| Feature | OneOff | Celery | Airflow | Rundeck |
|---------|--------|--------|---------|---------|
| Setup time | 30 sec | Hours | Hours | Hours |
| Dependencies | 0 | Redis + DB | DB + Scheduler | Java + DB |
| Memory | ~20MB | 500MB+ | 1GB+ | 1GB+ |
| Binary size | ~15MB | N/A | N/A | N/A |
| Learning curve | Minimal | Steep | Steep | Moderate |
| Best for | One-time jobs | Recurring tasks | Data pipelines | Runbooks |

---

## Development

```bash
# Clone
git clone https://github.com/meysam81/oneoff.git
cd oneoff

# Install dependencies
make setup

# Run backend (with hot reload)
make dev

# Run frontend (separate terminal)
bun run start
# or: npm run start

# Build production binary
make build
```

### Project Structure
```
oneoff/
├── main.go              # CLI entry point
├── internal/
│   ├── config/          # Environment configuration
│   ├── domain/          # Domain models
│   ├── handler/         # HTTP handlers
│   ├── jobs/            # Job executors (HTTP, Shell, Docker)
│   ├── repository/      # SQLite data layer
│   ├── service/         # Business logic
│   ├── server/          # HTTP server + embedded frontend
│   └── worker/          # Worker pool
├── migrations/          # Database migrations
└── src/                 # Vue 3 frontend
```

---

## Contributing

Contributions welcome! Please read the existing code style and open an issue before major changes.

```bash
# Run tests
make test

# Format code
go fmt ./...
```

---

## License

[MIT License](LICENSE) — Use it however you want.

---

<p align="center">
  <strong>Built for developers who value simplicity.</strong><br>
  <a href="https://github.com/meysam81/oneoff">Star on GitHub</a> •
  <a href="https://github.com/meysam81/oneoff/issues">Report Issue</a>
</p>
