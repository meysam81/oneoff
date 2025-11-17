# CLAUDE.md - AI Assistant Guide for OneOff

This document provides comprehensive guidance for AI assistants working on the OneOff codebase. It explains the architecture, conventions, workflows, and important considerations for making changes.

**Last Updated**: 2025-11-17
**Codebase Size**: ~5,300 lines across 49 source files
**Version**: Development (pre-1.0)

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture & Design Patterns](#architecture--design-patterns)
3. [Codebase Structure](#codebase-structure)
4. [Technology Stack](#technology-stack)
5. [Development Workflow](#development-workflow)
6. [Coding Conventions](#coding-conventions)
7. [Common Tasks & Patterns](#common-tasks--patterns)
8. [Testing Strategy](#testing-strategy)
9. [Important Considerations](#important-considerations)
10. [Troubleshooting](#troubleshooting)

---

## Project Overview

**OneOff** is a self-hosted, developer-focused job scheduler for executing one-time tasks at specific future times. Think "Linux `at` command meets modern web UI" with a plugin-based architecture.

### Core Features
- Single binary deployment with embedded Vue 3 frontend
- SQLite database with automatic migrations
- Multiple job types: HTTP requests, shell scripts, Docker containers
- Priority queue system (1-10 priority levels)
- Project organization and tag system
- Real-time worker monitoring
- Job chaining for sequential execution

### Key Design Goals
1. **Zero Dependencies**: Everything bundled into one executable
2. **Easy Deployment**: No external database or complex setup
3. **Developer-Friendly**: Clean API, modern UI, extensible architecture
4. **Production-Ready**: Graceful shutdown, error handling, observability

---

## Architecture & Design Patterns

### Backend Architecture (Go)

#### 1. Layered Architecture
```
┌─────────────────────────────────────────────────┐
│  HTTP Handlers (internal/handler/)              │
│  - Parse requests, validate input               │
│  - Call service layer                           │
│  - Format responses                             │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│  Service Layer (internal/service/)              │
│  - Business logic                               │
│  - Orchestration                                │
│  - Domain rule enforcement                      │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│  Repository Layer (internal/repository/)        │
│  - Database abstraction                         │
│  - CRUD operations                              │
│  - Query building                               │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│  SQLite Database                                │
└─────────────────────────────────────────────────┘
```

#### 2. Domain-Driven Design
- **Domain Models** (`internal/domain/models.go`): Core business entities
- **Domain Interfaces** (`internal/domain/job.go`): Contracts for job executors
- **Domain Errors** (`internal/domain/errors.go`): Business rule violations
- **Repository Interface** (`internal/domain/repository.go`): Data access contract

#### 3. Plugin Architecture (Job Registry)
```go
// Job types are registered via factory pattern
type JobExecutor interface {
    Execute(ctx context.Context, config string) (*ExecutionResult, error)
    Validate(config string) error
}

// Registry pattern for extensibility
registry := domain.NewJobRegistry()
registry.Register("http", jobs.NewHTTPJobFactory())
registry.Register("shell", jobs.NewShellJobFactory())
registry.Register("docker", jobs.NewDockerJobFactory())
```

#### 4. Worker Pool Pattern
- Configurable number of workers (default: CPU cores / 2)
- Job scheduler polls database every 5 seconds
- Jobs dispatched to workers via buffered channel
- Priority-based job selection (higher priority first, then FIFO)
- Graceful shutdown with context cancellation

#### 5. Embedded Frontend
```go
//go:embed dist/*
var staticFiles embed.FS

// Serve frontend from embedded filesystem
http.FileServer(http.FS(distFS))
```

### Frontend Architecture (Vue 3)

#### 1. Component Hierarchy
```
App.vue (root layout)
├── Sidebar.vue (navigation)
├── Header.vue (top bar)
└── <router-view> (dynamic content)
    ├── Dashboard.vue (stats overview)
    ├── Jobs.vue (job list + CRUD)
    │   ├── CreateJobModal.vue
    │   ├── JobsTable.vue
    │   └── job-configs/ (type-specific forms)
    ├── JobDetails.vue (single job view)
    ├── Executions.vue (execution history)
    ├── Projects.vue (project management)
    └── Settings.vue (system config)
```

#### 2. State Management (Pinia)
- **jobs.js**: Job list, current job, filters, CRUD operations
- **system.js**: System stats, worker status, configuration

#### 3. API Client Architecture
```javascript
// Centralized API client with retry logic
const api = ky.create({
  prefixUrl: '/api',
  timeout: 30000,
  retry: {
    limit: 2,
    methods: ['get', 'post', 'patch', 'delete']
  }
})

// Modular API functions
export const jobsAPI = { ... }
export const executionsAPI = { ... }
export const projectsAPI = { ... }
```

---

## Codebase Structure

### Backend Directory Layout

```
internal/
├── config/              # Configuration management
│   └── config.go       # Env var parsing with caarlos0/env
│
├── domain/             # Domain models and interfaces
│   ├── models.go       # Core entities (Job, Execution, Project, etc.)
│   ├── job.go          # JobExecutor interface
│   ├── errors.go       # Domain errors (ErrNotFound, ErrInvalidPriority, etc.)
│   └── registry.go     # Job type registry
│
├── handler/            # HTTP request handlers
│   ├── handler.go      # Base handler with common utilities
│   ├── job.go          # Job endpoints (CRUD, execute, clone, cancel)
│   ├── execution.go    # Execution endpoints
│   ├── project.go      # Project endpoints
│   ├── tag.go          # Tag endpoints
│   └── misc.go         # System/worker/job-type endpoints
│
├── jobs/               # Job executor implementations
│   ├── http.go         # HTTP request job (with retry)
│   ├── shell.go        # Shell script execution
│   ├── docker.go       # Docker container execution
│   └── registry.go     # Factory registration
│
├── repository/         # Database layer
│   ├── repository.go   # Repository interface
│   ├── sqlite.go       # SQLite implementation (job CRUD)
│   ├── sqlite_execution.go  # Execution-related queries
│   ├── sqlite_misc.go       # Project/tag/config queries
│   └── migrations.go   # Migration runner
│
├── service/            # Business logic
│   ├── job_service.go       # Job business logic
│   ├── execution_service.go # Execution tracking
│   ├── project_service.go   # Project management
│   ├── tag_service.go       # Tag management
│   └── system_service.go    # System stats and config
│
├── worker/             # Job execution workers
│   └── pool.go         # Worker pool with scheduler
│
└── server/             # HTTP server
    ├── server.go       # Server setup, routing, middleware
    └── dist/           # Embedded frontend (generated)
```

### Frontend Directory Layout

```
src/
├── components/         # Reusable Vue components
│   ├── CreateJobModal.vue
│   ├── ExecutionsList.vue
│   ├── Header.vue
│   ├── JobsTable.vue
│   ├── Sidebar.vue
│   └── job-configs/    # Job type-specific config forms
│       ├── DockerConfig.vue
│       ├── HTTPConfig.vue
│       └── ShellConfig.vue
│
├── views/              # Page-level components
│   ├── Dashboard.vue
│   ├── Executions.vue
│   ├── JobDetails.vue
│   ├── Jobs.vue
│   ├── Projects.vue
│   └── Settings.vue
│
├── stores/             # Pinia state management
│   ├── jobs.js         # Job state and operations
│   └── system.js       # System state and config
│
├── utils/              # Utility functions
│   └── api.js          # API client factory
│
├── App.vue             # Root component
├── main.js             # Application entry point
└── router.js           # Vue Router configuration
```

### Configuration Files

| File | Purpose |
|------|---------|
| `go.mod` | Go module dependencies |
| `package.json` | Node.js dependencies and scripts |
| `vite.config.js` | Vite build config (output: `internal/server/dist/`) |
| `.env.example` | Environment variable template |
| `Makefile` | Build automation and dev commands |
| `migrations/` | Database schema migrations |

---

## Technology Stack

### Backend (Go 1.23+)

| Package | Purpose | Notes |
|---------|---------|-------|
| `urfave/cli/v3` | CLI framework | Command-line interface and flags |
| `caarlos0/env/v11` | Config parsing | 12-factor app pattern, env vars |
| `mattn/go-sqlite3` | SQLite driver | CGO required for compilation |
| `golang-migrate/migrate/v4` | Schema migrations | Automatic migration on startup |
| `rs/zerolog` | Structured logging | JSON logging, leveled output |
| `imroc/req/v3` | HTTP client | Retry logic for HTTP jobs |
| Standard library | HTTP server, context, sync | No external web framework |

### Frontend (Vue 3 + Vite)

| Package | Purpose | Notes |
|---------|---------|-------|
| `vue@^3.5` | Framework | Composition API, SFCs |
| `vue-router@^4.5` | Routing | SPA navigation |
| `pinia@^2.3` | State management | Vue store |
| `naive-ui@^2.40` | UI components | Tree-shakable, dark theme |
| `ky@^1.7` | HTTP client | Lightweight, retry logic |
| `date-fns@^4.1` | Date utilities | Date formatting and manipulation |
| `@vicons/ionicons5` | Icons | Icon components |
| `vite@^6.0` | Build tool | Fast HMR, optimized bundling |

---

## Development Workflow

### Initial Setup

```bash
# 1. Clone repository
git clone https://github.com/meysam81/oneoff.git
cd oneoff

# 2. Install all dependencies (Go + Node.js)
make setup

# 3. Build application
make build

# 4. Run application
./oneoff

# Access at http://localhost:8080
```

### Development Commands

```bash
# Build everything (frontend + backend)
make build

# Build frontend only
make frontend

# Build backend only (requires frontend to be built first)
go build -o oneoff ./cmd/oneoff

# Run in development mode (no frontend build, Go only)
make dev

# Run frontend with hot reload (separate terminal)
npm run dev
# Access Vite dev server at http://localhost:5173
# API calls proxy to http://localhost:8080

# Run tests
make test

# Clean build artifacts
make clean

# Run migrations manually
./oneoff migrate --direction up
```

### Git Workflow

**Branch Convention**: All development branches should start with `claude/` prefix when working with Claude Code.

```bash
# Current development branch
git branch  # claude/claude-md-mi3buu09qni1mf1e-014jfuZoVSknEKpY2kQGntoH

# Create commits with clear, descriptive messages
git add .
git commit -m "feat: add job retry mechanism for failed executions"

# Push to remote (always use -u for new branches)
git push -u origin <branch-name>
```

### Configuration

All configuration is done via environment variables (12-factor app):

```bash
# Copy example config
cp .env.example .env

# Edit configuration
vim .env

# Key variables:
PORT=8080                    # HTTP server port
HOST=localhost               # HTTP server host
DB_PATH=./oneoff.db          # SQLite database path
WORKERS_COUNT=0              # Worker count (0 = auto: CPU cores / 2)
LOG_LEVEL=info               # Logging level (debug, info, warn, error)
DEFAULT_TIMEZONE=UTC         # Default timezone for jobs
LOG_RETENTION_DAYS=90        # How long to keep execution logs
DEFAULT_PRIORITY=5           # Default job priority (1-10)
ENVIRONMENT=production       # Environment name
```

---

## Coding Conventions

### Go Conventions

#### 1. Error Handling
```go
// Always wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to create job: %w", err)
}

// Use domain errors for business rule violations
if priority < 1 || priority > 10 {
    return nil, domain.ErrInvalidPriority
}

// Define domain errors in domain/errors.go
var (
    ErrNotFound = errors.New("not found")
    ErrInvalidPriority = errors.New("priority must be between 1 and 10")
)
```

#### 2. Context Propagation
```go
// Always accept context as first parameter
func (s *JobService) CreateJob(ctx context.Context, req domain.CreateJobRequest) (*domain.Job, error) {
    // Pass context to all downstream calls
    if err := s.repo.CreateJob(ctx, job, req.TagIDs); err != nil {
        return nil, err
    }
}
```

#### 3. Struct Initialization
```go
// Use designated initializers for clarity
job := &domain.Job{
    Name:        req.Name,
    Type:        req.Type,
    Config:      req.Config,
    ScheduledAt: scheduledAt.UTC(),
    Priority:    priority,
    Status:      domain.JobStatusScheduled,
}
```

#### 4. JSON Handling
```go
// Use struct tags for JSON marshaling
type Job struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    ScheduledAt time.Time `json:"scheduled_at"`
    Tags        []Tag     `json:"tags,omitempty"` // omitempty for optional fields
}

// For nullable fields, use pointers
type JobExecution struct {
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    ExitCode    *int       `json:"exit_code,omitempty"`
}
```

#### 5. Repository Pattern
```go
// Define interfaces in domain package
type Repository interface {
    CreateJob(ctx context.Context, job *Job, tagIDs []string) error
    GetJob(ctx context.Context, id string) (*Job, error)
    // ...
}

// Implement in repository package
type SQLiteRepository struct {
    db *sql.DB
}

func (r *SQLiteRepository) CreateJob(ctx context.Context, job *Job, tagIDs []string) error {
    // Implementation
}
```

#### 6. Service Layer Pattern
```go
// Services orchestrate business logic
type JobService struct {
    repo     repository.Repository  // Data access
    registry *domain.JobRegistry    // Job type registry
    pool     *worker.Pool           // Worker pool
}

// Validate inputs, enforce business rules, coordinate operations
func (s *JobService) CreateJob(ctx context.Context, req domain.CreateJobRequest) (*domain.Job, error) {
    // 1. Validate job type
    if _, err := s.registry.Create(req.Type, req.Config); err != nil {
        return nil, fmt.Errorf("invalid job type or config: %w", err)
    }

    // 2. Validate business rules
    if scheduledAt.Before(time.Now().UTC()) {
        return nil, domain.ErrInvalidScheduleTime
    }

    // 3. Set defaults
    // 4. Call repository
    // 5. Return result
}
```

### Vue 3 Conventions

#### 1. Composition API
```vue
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useJobsStore } from '@/stores/jobs'

// Use stores
const jobsStore = useJobsStore()

// Reactive state
const loading = ref(false)
const selectedJob = ref(null)

// Computed properties
const hasJobs = computed(() => jobsStore.jobs.length > 0)

// Lifecycle
onMounted(async () => {
  await jobsStore.fetchJobs()
})
</script>
```

#### 2. Component Structure
```vue
<!-- Template -->
<template>
  <div class="container">
    <n-card>
      <!-- Content -->
    </n-card>
  </div>
</template>

<!-- Script -->
<script setup>
// Imports
// Props/emits
// Reactive state
// Computed
// Methods
// Lifecycle hooks
</script>

<!-- Styles (scoped) -->
<style scoped>
.container {
  padding: 20px;
}
</style>
```

#### 3. Pinia Store Pattern
```javascript
import { defineStore } from 'pinia'
import { jobsAPI } from '@/utils/api'

export const useJobsStore = defineStore('jobs', {
  state: () => ({
    jobs: [],
    currentJob: null,
    loading: false,
    error: null
  }),

  getters: {
    scheduledJobs: (state) => state.jobs.filter(j => j.status === 'scheduled')
  },

  actions: {
    async fetchJobs() {
      this.loading = true
      try {
        this.jobs = await jobsAPI.list()
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    }
  }
})
```

#### 4. API Client Pattern
```javascript
import ky from 'ky'

const api = ky.create({
  prefixUrl: '/api',
  timeout: 30000,
  retry: { limit: 2 }
})

export const jobsAPI = {
  list: async (params) => await api.get('jobs', { searchParams: params }).json(),
  get: async (id) => await api.get(`jobs/${id}`).json(),
  create: async (data) => await api.post('jobs', { json: data }).json(),
  update: async (id, data) => await api.patch(`jobs/${id}`, { json: data }).json(),
  delete: async (id) => await api.delete(`jobs/${id}`)
}
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Go packages | lowercase, single word | `handler`, `service`, `repository` |
| Go types | PascalCase | `JobService`, `JobExecutor` |
| Go interfaces | PascalCase, often ends in -er | `Repository`, `JobExecutor` |
| Go functions | PascalCase (exported), camelCase (private) | `CreateJob()`, `validateInput()` |
| Go constants | PascalCase or SCREAMING_SNAKE_CASE | `JobStatusScheduled`, `DEFAULT_PRIORITY` |
| Vue components | PascalCase | `CreateJobModal.vue`, `JobsTable.vue` |
| Vue files | PascalCase for components, lowercase for utils | `Dashboard.vue`, `api.js` |
| Pinia stores | camelCase with "use" prefix | `useJobsStore`, `useSystemStore` |
| JavaScript functions | camelCase | `fetchJobs()`, `formatDate()` |
| CSS classes | kebab-case | `.job-container`, `.status-badge` |

---

## Common Tasks & Patterns

### Adding a New Job Type

1. **Create job executor** (`internal/jobs/newtype.go`):
```go
package jobs

import (
    "context"
    "encoding/json"
    "github.com/meysam81/oneoff/internal/domain"
)

// NewTypeConfig defines configuration for new job type
type NewTypeConfig struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

// NewTypeJob implements the new job type
type NewTypeJob struct {
    config NewTypeConfig
}

func (j *NewTypeJob) Execute(ctx context.Context) (*domain.ExecutionResult, error) {
    // Implementation
}

func (j *NewTypeJob) Validate() error {
    // Validation
}

// Factory function
func NewNewTypeJobFactory() domain.JobFactory {
    return func(config string) (domain.JobExecutor, error) {
        var cfg NewTypeConfig
        if err := json.Unmarshal([]byte(config), &cfg); err != nil {
            return nil, err
        }

        job := &NewTypeJob{config: cfg}
        if err := job.Validate(); err != nil {
            return nil, err
        }
        return job, nil
    }
}
```

2. **Register in server** (`internal/server/server.go`):
```go
registry := domain.NewJobRegistry()
registry.Register("http", jobs.NewHTTPJobFactory())
registry.Register("shell", jobs.NewShellJobFactory())
registry.Register("docker", jobs.NewDockerJobFactory())
registry.Register("newtype", jobs.NewNewTypeJobFactory())  // Add here
```

3. **Create Vue config component** (`src/components/job-configs/NewTypeConfig.vue`):
```vue
<template>
  <n-form-item label="Field 1" required>
    <n-input v-model:value="config.field1" placeholder="Enter field 1" />
  </n-form-item>
  <n-form-item label="Field 2" required>
    <n-input-number v-model:value="config.field2" :min="0" />
  </n-form-item>
</template>

<script setup>
import { defineProps, defineEmits, watch } from 'vue'

const props = defineProps({
  modelValue: Object
})

const emit = defineEmits(['update:modelValue'])

const config = props.modelValue || { field1: '', field2: 0 }

watch(config, (newVal) => {
  emit('update:modelValue', newVal)
}, { deep: true })
</script>
```

4. **Update CreateJobModal** to include new config component

### Adding a New API Endpoint

1. **Add handler method** (`internal/handler/job.go` or new file):
```go
func (h *Handler) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Parse request
    var req SomeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Call service
    result, err := h.someService.DoSomething(ctx, req)
    if err != nil {
        h.respondError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond
    h.respondJSON(w, result, http.StatusOK)
}
```

2. **Register route** (`internal/server/server.go`):
```go
router.HandleFunc("POST /api/something", handler.HandleNewEndpoint)
```

3. **Add API client function** (`src/utils/api.js`):
```javascript
export const someAPI = {
  doSomething: async (data) =>
    await api.post('something', { json: data }).json()
}
```

4. **Use in Vue component/store**

### Adding a Database Migration

1. **Create migration files**:
```bash
# Up migration
cat > migrations/000002_add_new_table.up.sql << 'EOF'
CREATE TABLE IF NOT EXISTS new_table (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF

# Down migration
cat > migrations/000002_add_new_table.down.sql << 'EOF'
DROP TABLE IF EXISTS new_table;
EOF
```

2. **Run migration**:
```bash
./oneoff migrate --direction up
```

3. **Update repository** to add new queries

---

## Testing Strategy

### Current State
- **No tests implemented yet** (test infrastructure needs to be added)
- Go testing framework is standard `testing` package
- Consider using `testify` for assertions and mocking

### Recommended Testing Approach

#### Backend Testing

1. **Unit Tests**
```go
// internal/service/job_service_test.go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestJobService_CreateJob(t *testing.T) {
    // Setup
    mockRepo := new(MockRepository)
    mockRegistry := domain.NewJobRegistry()
    service := NewJobService(mockRepo, mockRegistry, nil)

    // Test case
    req := domain.CreateJobRequest{
        Name: "test job",
        Type: "http",
        // ...
    }

    mockRepo.On("CreateJob", mock.Anything, mock.Anything, mock.Anything).
        Return(nil)

    // Execute
    job, err := service.CreateJob(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, job)
    mockRepo.AssertExpectations(t)
}
```

2. **Integration Tests**
```go
// internal/repository/sqlite_test.go
func TestSQLiteRepository_CreateAndGetJob(t *testing.T) {
    // Setup in-memory SQLite
    db, _ := sql.Open("sqlite3", ":memory:")
    defer db.Close()

    // Run migrations
    // Create repository
    // Test actual database operations
}
```

3. **Handler Tests**
```go
// internal/handler/job_test.go
func TestHandler_CreateJob(t *testing.T) {
    // Use httptest.NewRecorder() and httptest.NewRequest()
    // Test HTTP request/response handling
}
```

#### Frontend Testing

1. **Component Tests** (with Vitest)
```javascript
// src/components/__tests__/JobsTable.spec.js
import { mount } from '@vue/test-utils'
import JobsTable from '../JobsTable.vue'

describe('JobsTable', () => {
  it('renders job rows', () => {
    const wrapper = mount(JobsTable, {
      props: { jobs: [/* test data */] }
    })
    expect(wrapper.findAll('.job-row')).toHaveLength(1)
  })
})
```

2. **Store Tests**
```javascript
// src/stores/__tests__/jobs.spec.js
import { setActivePinia, createPinia } from 'pinia'
import { useJobsStore } from '../jobs'

describe('Jobs Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('fetches jobs', async () => {
    const store = useJobsStore()
    await store.fetchJobs()
    expect(store.jobs).toBeDefined()
  })
})
```

---

## Important Considerations

### Security

1. **SQL Injection Prevention**
   - Always use parameterized queries
   - Never concatenate user input into SQL strings
   ```go
   // GOOD
   db.QueryContext(ctx, "SELECT * FROM jobs WHERE id = ?", jobID)

   // BAD
   db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM jobs WHERE id = '%s'", jobID))
   ```

2. **Command Injection** (Shell/Docker jobs)
   - Shell jobs execute arbitrary commands - document security implications
   - Docker jobs require Docker socket access - document risks
   - Consider adding execution restrictions in production

3. **XSS Prevention**
   - Vue automatically escapes template interpolations
   - Be careful with `v-html` directive
   - Sanitize job outputs before display

4. **CORS Configuration**
   - Currently configured in `internal/server/server.go`
   - Adjust for production deployment

### Performance

1. **Database Queries**
   - Use indexes for frequently queried fields (already in migrations)
   - Limit result sets with pagination
   - Use prepared statements (already implemented)

2. **Worker Pool Sizing**
   - Default: CPU cores / 2
   - Adjust based on workload type (CPU vs I/O bound)
   - Monitor worker utilization

3. **Frontend Bundle Size**
   - Naive UI is tree-shakable - import only needed components
   - Vite automatically code-splits vendor bundles
   - Lazy-load routes if app grows

### Data Integrity

1. **Time Zones**
   - All times stored in UTC in database
   - User-specified timezone stored in `jobs.timezone` field
   - Convert to UTC for storage, convert back for display

2. **Job Configuration**
   - Stored as JSON string in `jobs.config` field
   - Validate on creation with job executor's `Validate()` method
   - Parse/marshal carefully

3. **Cascading Deletes**
   - Deleting a project cascades to jobs (check migrations)
   - Deleting a job cascades to executions
   - Document implications for users

### Operational Concerns

1. **Database Migrations**
   - Migrations run automatically on server startup
   - Always create both up and down migrations
   - Test migrations before deployment

2. **Graceful Shutdown**
   - Server supports context cancellation
   - Workers should respect context timeouts
   - Ensure in-flight jobs complete or are marked properly

3. **Logging**
   - Use structured logging with zerolog
   - Log levels: debug, info, warn, error
   - Avoid logging sensitive data (tokens, passwords in job configs)

4. **Configuration Validation**
   - Server fails fast on invalid config (see `internal/config/config.go`)
   - Validate required environment variables at startup
   - Provide clear error messages

---

## Troubleshooting

### Common Issues

#### 1. Build Failures

**Problem**: `go build` fails with SQLite errors
```
# github.com/mattn/go-sqlite3
cgo: C compiler not found
```

**Solution**: Install GCC/build tools
```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# macOS
xcode-select --install

# Alpine (for Docker)
apk add gcc musl-dev
```

#### 2. Frontend Build Issues

**Problem**: `npm run build` fails
```
ENOENT: no such file or directory, mkdir 'internal/server/dist'
```

**Solution**: Ensure parent directories exist
```bash
mkdir -p internal/server/dist
npm run build
```

#### 3. Migration Failures

**Problem**: Migrations fail on startup
```
migration failed: table already exists
```

**Solution**: Check database schema version
```bash
# Remove database and start fresh (dev only!)
rm oneoff.db oneoff.db-shm oneoff.db-wal
./oneoff
```

#### 4. Worker Pool Not Executing Jobs

**Problem**: Jobs stay in "scheduled" status

**Checklist**:
- Check worker pool is running: `GET /api/workers/status`
- Verify scheduled time is in the past
- Check worker logs for errors
- Ensure `WORKERS_COUNT > 0` or uses default

#### 5. Frontend Can't Reach API

**Problem**: API calls fail with 404

**Solutions**:
- Check server is running on correct port
- Verify Vite proxy configuration in `vite.config.js`
- Check CORS settings in `internal/server/server.go`
- Ensure API routes are registered correctly

### Development Tips

1. **Use Development Mode**
   ```bash
   # Terminal 1: Backend with hot reload
   make dev

   # Terminal 2: Frontend with hot reload
   npm run dev
   ```

2. **Enable Debug Logging**
   ```bash
   LOG_LEVEL=debug ./oneoff
   ```

3. **Inspect Database**
   ```bash
   sqlite3 oneoff.db
   .schema  # Show schema
   .tables  # List tables
   SELECT * FROM jobs;  # Query data
   ```

4. **Test API Endpoints**
   ```bash
   # List jobs
   curl http://localhost:8080/api/jobs

   # Create job
   curl -X POST http://localhost:8080/api/jobs \
     -H "Content-Type: application/json" \
     -d '{"name":"test","type":"http","config":"...","scheduled_at":"2025-12-01T00:00:00Z"}'
   ```

5. **Monitor Worker Status**
   ```bash
   curl http://localhost:8080/api/workers/status | jq
   ```

### Debugging Workflow

1. **Backend Issues**
   - Add debug logs with `log.Debug().Msgf("...")`
   - Use Go debugger (Delve): `dlv debug ./cmd/oneoff`
   - Check error wrapping for context

2. **Frontend Issues**
   - Use Vue DevTools browser extension
   - Check browser console for errors
   - Inspect Pinia store state
   - Use network tab to debug API calls

3. **Database Issues**
   - Enable SQL query logging (add to repository layer)
   - Use SQLite CLI to inspect data
   - Check foreign key constraints

---

## Quick Reference

### File Locations

| What | Where |
|------|-------|
| Entry point | `cmd/oneoff/main.go` (needs to be created) |
| Server setup | `internal/server/server.go` |
| HTTP handlers | `internal/handler/*.go` |
| Business logic | `internal/service/*.go` |
| Database | `internal/repository/*.go` |
| Job executors | `internal/jobs/*.go` |
| Domain models | `internal/domain/*.go` |
| Migrations | `migrations/*.sql` |
| Frontend entry | `src/main.js` |
| Vue components | `src/components/*.vue` |
| Vue pages | `src/views/*.vue` |
| API client | `src/utils/api.js` |
| Stores | `src/stores/*.js` |

### Key Commands

```bash
make setup      # First-time setup
make build      # Production build
make dev        # Development mode (backend)
npm run dev     # Development mode (frontend)
make test       # Run tests
make clean      # Clean artifacts
./oneoff        # Run application
```

### Environment Variables

```bash
PORT=8080                    # Server port
HOST=localhost               # Server host
DB_PATH=./oneoff.db          # Database file
WORKERS_COUNT=0              # Worker count (0=auto)
LOG_LEVEL=info               # Log level
DEFAULT_TIMEZONE=UTC         # Default TZ
LOG_RETENTION_DAYS=90        # Log retention
DEFAULT_PRIORITY=5           # Default priority
ENVIRONMENT=production       # Environment
```

### API Endpoints

```
Jobs:        GET/POST/PATCH/DELETE /api/jobs
             POST /api/jobs/:id/execute
             POST /api/jobs/:id/clone
             POST /api/jobs/:id/cancel
Executions:  GET /api/executions
Projects:    GET/POST/PATCH/DELETE /api/projects
Tags:        GET/POST/PATCH/DELETE /api/tags
System:      GET /api/system/status
             GET /api/system/config
             PATCH /api/system/config
Workers:     GET /api/workers/status
Job Types:   GET /api/job-types
```

---

## Additional Resources

- **README.md**: User-facing documentation and quick start
- **migrations/**: Database schema and evolution
- **package.json**: Frontend dependencies and scripts
- **go.mod**: Backend dependencies
- **Makefile**: Build commands and automation

---

## Contributing Guidelines for AI Assistants

When making changes to this codebase:

1. **Always read relevant files before editing** - Understand existing patterns
2. **Follow established conventions** - Match existing code style
3. **Add proper error handling** - Wrap errors with context
4. **Update documentation** - Keep this file and README current
5. **Consider security implications** - Especially for job execution
6. **Test thoroughly** - Even if automated tests don't exist yet
7. **Use small, focused commits** - One logical change per commit
8. **Update CLAUDE.md** - If you change architecture or patterns

### Before Committing

- [ ] Code follows Go/Vue conventions
- [ ] Error handling is comprehensive
- [ ] No security vulnerabilities introduced
- [ ] Documentation is updated if needed
- [ ] Build succeeds (`make build`)
- [ ] Application runs without errors

---

**Document Version**: 1.0
**Last Updated**: 2025-11-17
**Maintainer**: OneOff Development Team
