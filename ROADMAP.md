# OneOff Product Roadmap

> **Last Updated**: 2025-11-29
> **Version**: 2.0.0
> **Status**: Active Development

This document outlines the strategic product roadmap for OneOff, a self-hosted, developer-focused job scheduler for executing one-time tasks at specific future times. Every roadmap item strengthens our core value proposition: **zero-dependency simplicity with production-grade capabilities**.

---

## Table of Contents

- [Executive Summary](#executive-summary)
- [Core Value Proposition](#core-value-proposition)
- [Ideal Customer Profiles](#ideal-customer-profiles)
- [Current State](#current-state)
- [Competitive Positioning](#competitive-positioning)
- [Phased Roadmap](#phased-roadmap)
- [Detailed Feature Analysis](#detailed-feature-analysis)
- [Anti-Roadmap: What We Won't Build](#anti-roadmap-what-we-wont-build)
- [Contributing](#contributing)
- [Changelog](#changelog)

---

## Executive Summary

OneOff is a production-ready job scheduler at v1.0.2 with a clean Go backend, modern Vue 3 frontend, and single-binary deployment. The architecture is solid—layered services, plugin-based job executors, and proper domain separation.

**Current state:**
- Core job scheduling is complete and production-ready
- Three job types (HTTP, Shell, Docker) fully functional
- Projects, Tags, and Execution history fully implemented
- API Authentication and Webhooks have complete backends but missing frontend UI
- Job Chaining has database schema but incomplete implementation

**Strategic focus:**
1. **Complete partial implementations** - API Keys UI, Webhooks UI, Job Chaining (highest ROI)
2. **Production essentials** - Metrics, Log retention, Running job cancellation
3. **Developer experience** - CLI tool, OpenAPI docs
4. **Platform expansion** - Helm chart, recurring jobs (optional mode)

---

## Core Value Proposition

**OneOff is the antidote to over-engineering for one-time task scheduling.**

When you need to send a webhook at 3 PM, run a database migration tomorrow at midnight, or trigger a deployment during a maintenance window, you shouldn't need Redis, Postgres, Celery, and a web of services. You need to download one binary, run it, and schedule your job.

OneOff delivers:
- **30-second setup**: Download, run, done
- **Zero dependencies**: No Redis, no Postgres, no message queues
- **Modern experience**: Dark-mode UI, real-time monitoring, API-first
- **Production-grade**: Priority queues, execution history, graceful shutdown

Every feature we add must pass this test: *Does this make OneOff simpler to use while maintaining production reliability?*

---

## Ideal Customer Profiles

### ICP 1: The Solo DevOps Engineer

**Who they are**: DevOps/Platform engineer at a startup (1-20 employees), wearing many hats, managing infrastructure for 3-10 services.

**Their context**: They have limited time and hate maintaining infrastructure. Every new system is another thing to monitor, upgrade, and debug at 3 AM.

**Jobs to be Done (JTBD)**:
1. Schedule database backups before major deployments
2. Send webhook notifications at specific business times
3. Run cleanup scripts during off-peak hours
4. Trigger deployment rollbacks if health checks fail after N minutes

**Pain points with alternatives**:
- Celery/Airflow require Redis/Postgres and days of setup
- Cron is CLI-only with no visibility into failures
- `at` command has no persistence, UI, or history

**What "exceptional" looks like for them**:
- Install in under 1 minute
- Create first job through UI in under 2 minutes
- Know immediately when jobs fail via webhooks
- Never think about OneOff's infrastructure

---

### ICP 2: The Backend Developer

**Who they are**: Full-stack or backend developer who needs to schedule one-off tasks as part of application workflows.

**Their context**: Building features that require scheduled actions—trial expirations, scheduled reports, delayed notifications, data sync jobs.

**Jobs to be Done (JTBD)**:
1. Schedule HTTP webhooks when user trials expire
2. Run data export jobs at customer-specified times
3. Execute migration scripts during deployment windows
4. Chain multiple tasks in sequence (backup → migrate → notify)

**Pain points with alternatives**:
- Application-embedded scheduling adds complexity and state
- External schedulers are overkill for simple webhooks
- No good middle-ground between `setTimeout` and Airflow

**What "exceptional" looks like for them**:
- Simple REST API for creating jobs programmatically
- SDK/CLI for scripting job creation
- Job chaining for multi-step workflows
- Clear execution logs for debugging

---

### ICP 3: The Platform Engineer (Mid-size Company)

**Who they are**: Platform/Infrastructure engineer at a 50-500 person company, responsible for internal developer tools and operational automation.

**Their context**: They have Airflow/Temporal for complex workflows but need something lighter for simple scheduled tasks that don't warrant full DAG infrastructure.

**Jobs to be Done (JTBD)**:
1. Provide self-service scheduling to developers without Airflow overhead
2. Run operational tasks (certificate rotation, log cleanup) on schedule
3. Monitor scheduled job execution with existing observability stack
4. Integrate with CI/CD for scheduled deployments

**Pain points with alternatives**:
- Airflow is overkill for single-task jobs
- Managing multiple cron servers is fragmented
- No visibility across all scheduled tasks

**What "exceptional" looks like for them**:
- Prometheus metrics for integration with Grafana dashboards
- Helm chart for Kubernetes deployment
- API keys for team-level access control
- Webhook integration with Slack/PagerDuty

---

### ICP 4: The SRE/Infrastructure Engineer

**Who they are**: Site Reliability Engineer responsible for operational excellence, incident response, and automation.

**Their context**: Needs reliable task execution with full observability, audit trails, and integration with existing monitoring infrastructure.

**Jobs to be Done (JTBD)**:
1. Schedule maintenance tasks with guaranteed execution tracking
2. Create runbooks that execute automatically at specified times
3. Set up escalation chains (alert → wait 5 min → page on-call)
4. Audit all scheduled actions for compliance

**Pain points with alternatives**:
- Cron jobs are invisible to monitoring
- Rundeck is heavyweight and Java-based
- PagerDuty scheduled actions are limited

**What "exceptional" looks like for them**:
- Full execution audit trail with output capture
- Prometheus metrics for job success/failure rates
- API for programmatic job management
- Integration with existing auth systems (future: OAuth)

---

## Current State

### Fully Implemented

| Feature | Description | Location |
|---------|-------------|----------|
| **Job CRUD** | Create, read, update, delete jobs | `internal/handler/job.go` |
| **HTTP Job Executor** | Execute HTTP requests with retry | `internal/jobs/http.go` |
| **Shell Job Executor** | Execute shell scripts with env vars | `internal/jobs/shell.go` |
| **Docker Job Executor** | Run Docker containers | `internal/jobs/docker.go` |
| **Priority Queue** | 1-10 priority levels, FIFO within priority | `internal/worker/pool.go` |
| **Worker Pool** | Configurable workers with graceful shutdown | `internal/worker/pool.go` |
| **Projects System** | Full CRUD with UI | `src/views/Projects.vue` |
| **Tags System** | Job categorization with colors | `internal/service/tag_service.go` |
| **Execution History** | Full audit trail with output | `internal/service/execution_service.go` |
| **Job Cloning** | Clone existing jobs | `internal/handler/job.go` |
| **Immediate Execution** | `scheduled_at: "now"` support | `internal/service/job_service.go` |
| **API Key Auth (Backend)** | Full CRUD, scopes, rotation, expiration | `internal/handler/apikey.go` |
| **Webhooks (Backend)** | Full CRUD, events, HMAC signing, retry | `internal/handler/webhook.go` |
| **Dark Mode UI** | Vue 3 with Naive UI | `src/App.vue` |
| **Landing Page** | Astro-based marketing site | `landing-page/` |
| **Job Template Catalog** | 5 reusable job templates | `landing-page/src/content/catalog/` |

### Partially Implemented

| Feature | What Exists | What's Missing | Effort to Complete |
|---------|-------------|----------------|-------------------|
| **API Key Management UI** | Complete backend with CRUD, scopes, rotation | Frontend page and components | S (2-3 days) |
| **Webhook Management UI** | Complete backend with CRUD, events, testing | Frontend page and components | S (3-4 days) |
| **Job Chaining** | Database schema (`job_chains`, `job_chain_links`), domain models | Service layer, handlers, execution logic, UI | M (5-7 days) |
| **Running Job Cancellation** | Works for scheduled jobs | Context cancellation for in-progress jobs | S (2-3 days) |

### Not Yet Implemented

| Feature | Priority | Notes |
|---------|----------|-------|
| Prometheus Metrics | P0 | Critical for ICP 3/4 |
| Log Retention Enforcement | P1 | Config exists, logic doesn't |
| CLI Tool | P1 | Developer experience |
| OpenAPI Documentation | P2 | API discoverability |
| Helm Chart | P1 | Kubernetes deployment |
| Test Suite | P1 | Quality confidence |
| Recurring Jobs (Cron) | P2 | Scope expansion, optional |
| OAuth2/OIDC | P3 | Enterprise feature |
| RBAC | P3 | Enterprise feature |
| Kubernetes Job Type | P3 | Cloud-native expansion |

---

## Competitive Positioning

### Current Landscape

| Alternative | Strengths | Weaknesses | OneOff Advantage |
|-------------|-----------|------------|------------------|
| **`at` command** | Simple, built-in | No UI, no persistence, no history | Modern UI, API, full audit trail |
| **`cron`** | Ubiquitous, reliable | Recurring-only, no visibility | One-time focus, real-time monitoring |
| **Celery** | Powerful, Python ecosystem | Redis/DB required, complex setup | Zero dependencies, 30-sec setup |
| **Airflow** | Enterprise-grade DAGs | 1GB+ memory, steep learning curve | 20MB memory, minimal learning curve |
| **Temporal** | Workflow orchestration | Complex, requires cluster | Single binary, no cluster needed |
| **Rundeck** | Runbook automation | Java-based, heavyweight | Go-native, lightweight |

### Our Moat

1. **Single-binary distribution**: No dependency management, no version conflicts
2. **SQLite persistence**: No database to manage, automatic migrations
3. **Focused scope**: One-time jobs done exceptionally well
4. **Modern stack**: Go + Vue 3 + Naive UI (vs. Java or legacy frameworks)
5. **API-first design**: Everything scriptable, extensible

### Positioning Statement

> For **developers and DevOps engineers** who need to **schedule one-time tasks**, **OneOff** is the **self-hosted job scheduler** that **eliminates infrastructure overhead with a single binary**, unlike **Celery, Airflow, and cron** which **require external dependencies or lack visibility**.

---

## Phased Roadmap

### Phase 1: Complete the Core

**Outcome**: All partially-implemented features are complete. Users get full value from existing backend capabilities.

**Target ICPs**: All (foundation for everyone)

| Priority | Feature | Score | Effort | Status | Notes |
|----------|---------|-------|--------|--------|-------|
| P0 | API Key Management UI | 9.5 | S | Partial | Backend 100% done |
| P0 | Webhook Management UI | 9.5 | S | Partial | Backend 100% done |
| P0 | Running Job Cancellation | 9.2 | S | Partial | Context propagation |
| P1 | Prometheus Metrics | 9.0 | S | Planned | `/metrics` endpoint |
| P1 | Log Retention Enforcement | 8.5 | XS | Planned | Config exists |

**Success Criteria**:
- Users can create/manage API keys through the UI
- Users can configure webhooks for job events through the UI
- Running jobs can be cancelled with proper cleanup
- Prometheus can scrape job metrics
- Old execution logs are automatically cleaned up

---

### Phase 2: Developer Experience

**Outcome**: OneOff is a joy to use programmatically. Developers can script, automate, and integrate.

**Target ICPs**: ICP 2 (Backend Developer), ICP 3 (Platform Engineer)

| Priority | Feature | Score | Effort | Status | Notes |
|----------|---------|-------|--------|--------|-------|
| P0 | Job Chaining API & UI | 9.0 | M | Partial | Schema exists |
| P1 | CLI Tool | 8.5 | M | Planned | `oneoff jobs create` |
| P1 | OpenAPI Documentation | 8.0 | S | Planned | Auto-generated |
| P2 | Job Templates in App | 7.5 | S | Planned | Import from catalog |

**Success Criteria**:
- Users can create and execute job chains through UI
- Developers can manage jobs entirely via CLI
- API documentation is available at `/api/docs`
- Job templates can be imported with one click

---

### Phase 3: Platform Integration

**Outcome**: OneOff deploys easily to Kubernetes and integrates with existing infrastructure.

**Target ICPs**: ICP 3 (Platform Engineer), ICP 4 (SRE)

| Priority | Feature | Score | Effort | Status | Notes |
|----------|---------|-------|--------|--------|-------|
| P0 | Helm Chart | 8.5 | S | Planned | Official chart |
| P1 | Comprehensive Test Suite | 8.0 | L | Planned | Unit + Integration + E2E |
| P1 | Job Retry Policies | 8.0 | S | Planned | Max retries, backoff |
| P2 | Structured Audit Logging | 7.5 | S | Planned | JSON audit trail |

**Success Criteria**:
- OneOff can be deployed to Kubernetes via `helm install`
- Test coverage reaches 70%+
- Failed jobs can automatically retry with exponential backoff
- All actions are audit-logged in structured format

---

### Phase 4: Capability Expansion

**Outcome**: OneOff handles more use cases while maintaining simplicity.

**Target ICPs**: ICP 3, ICP 4 (advanced users)

| Priority | Feature | Score | Effort | Status | Notes |
|----------|---------|-------|--------|--------|-------|
| P1 | Recurring Jobs (Cron Mode) | 7.5 | L | Planned | Optional, preserves identity |
| P2 | OAuth2/OIDC Integration | 7.0 | L | Planned | Enterprise SSO |
| P2 | Kubernetes Job Type | 7.0 | M | Planned | K8s native execution |
| P3 | RBAC | 6.5 | L | Future | Role-based permissions |

**Success Criteria**:
- Users can optionally enable cron-like scheduling
- Enterprise users can authenticate via SSO
- Platform teams can execute Kubernetes Jobs

---

### Phase 5: Ecosystem (Future)

**Outcome**: OneOff becomes a platform with integrations and extensibility.

| Priority | Feature | Score | Effort | Status |
|----------|---------|-------|--------|--------|
| P2 | Terraform Provider | 6.5 | L | Future |
| P2 | Database Query Job Type | 6.5 | M | Future |
| P3 | AWS Lambda/Cloud Functions | 6.0 | M | Future |
| P3 | S3/GCS Log Archival | 5.5 | M | Future |
| P4 | Horizontal Scaling | 4.0 | XL | Future |

---

## Detailed Feature Analysis

### Tier 1: Critical (Complete the Core)

#### API Key Management UI

**Score**: 9.5/10 (Weighted: ICP=9, Core=10, Leverage=10, Diff=8, Clarity=10)

**Why it matters**: API keys are fully implemented in the backend with sophisticated features (scopes, rotation, expiration, revocation). Users currently cannot access this functionality without direct API calls. This is the highest-ROI feature—zero new backend work required.

| Pros | Cons |
|------|------|
| Backend 100% complete | Adds another settings page |
| Unblocks production API usage | Need key display UX (show once) |
| Security feature (scopes) already built | |
| Table stakes for production | |

**Current State**:
- `internal/handler/apikey.go`: Full CRUD, revoke, rotate handlers
- `internal/service/apikey_service.go`: Business logic complete
- `migrations/000002_add_api_keys.up.sql`: Schema exists
- `src/views/`: No API keys page

**Implementation Approach**:
1. Create `src/views/ApiKeys.vue` with table of keys
2. Add create modal with name, scopes selector, expiration picker
3. Show raw key ONLY on creation (security best practice)
4. Add revoke/rotate actions with confirmation
5. Add route to `src/router.js`
6. Add navigation item to `src/components/Sidebar.vue`

**Success Criteria**:
- Users can create API keys with custom scopes
- Raw key shown once on creation
- Users can revoke/rotate keys
- Key last-used time visible

---

#### Webhook Management UI

**Score**: 9.5/10 (Weighted: ICP=9, Core=10, Leverage=10, Diff=9, Clarity=10)

**Why it matters**: Webhooks enable the integration ecosystem (Slack, Discord, PagerDuty). Backend is 100% complete with HMAC signing, retry logic, and delivery tracking. Users need UI to configure.

| Pros | Cons |
|------|------|
| Backend 100% complete | Webhook testing UX needed |
| Enables Slack/Discord/PagerDuty | Event selection UI complexity |
| Most requested integration type | |
| HMAC signing already implemented | |

**Current State**:
- `internal/handler/webhook.go`: Full CRUD, test, delivery history
- `internal/service/webhook_service.go`: Business logic with retry
- `migrations/000002_add_api_keys.up.sql`: Schema exists (webhooks table)
- `src/views/`: No webhooks page

**Implementation Approach**:
1. Create `src/views/Webhooks.vue` with webhook list
2. Add create/edit modal with URL, events checkboxes, secret field
3. Add "Test Webhook" button that fires test event
4. Show delivery history with status, response code, retry count
5. Add route and navigation

**Success Criteria**:
- Users can configure webhooks for job events
- Webhook secrets can be set for HMAC verification
- Test button validates webhook endpoint
- Delivery history shows success/failure

---

#### Running Job Cancellation

**Score**: 9.2/10 (Weighted: ICP=9, Core=10, Leverage=8, Diff=9, Clarity=9)

**Why it matters**: Users expect "Cancel" to work for running jobs. Currently, cancel only works for scheduled (not yet running) jobs. Long-running jobs cannot be stopped.

| Pros | Cons |
|------|------|
| Critical for long-running jobs | Docker container cleanup tricky |
| User expectation mismatch today | Shell process group handling |
| Context pattern already in Go | May leave orphaned processes |

**Current State**:
- Job status updates to "cancelled" but execution continues
- Workers don't check context cancellation during execution
- No process group handling for shell jobs

**Implementation Approach**:
1. Store `context.CancelFunc` in worker when job starts
2. Add job ID → cancel function map in worker pool
3. On cancel request, call cancel function and update status
4. Shell executor: Kill process group with SIGTERM, then SIGKILL after timeout
5. Docker executor: `docker stop` with timeout
6. HTTP executor: Context already supports cancellation

**Success Criteria**:
- Running HTTP jobs stop immediately on cancel
- Running shell jobs terminate process group
- Running Docker jobs stop container
- Job status reflects cancelled state

---

#### Prometheus Metrics

**Score**: 9.0/10 (Weighted: ICP=10, Core=8, Leverage=7, Diff=9, Clarity=9)

**Why it matters**: DevOps and SRE teams (ICP 3, 4) expect Prometheus metrics for integration with existing observability stacks (Grafana, Alertmanager).

| Pros | Cons |
|------|------|
| Industry standard observability | Cardinality management needed |
| Enables Grafana dashboards | Another endpoint to maintain |
| SRE adoption driver | Metrics selection decisions |
| Alerting pipeline integration | |

**Current State**: No metrics implementation.

**Implementation Approach**:
1. Add `prometheus/client_golang` dependency
2. Create `internal/metrics/metrics.go` with metric definitions
3. Expose `/metrics` endpoint in server
4. Instrument: job counts, execution duration, worker status, queue depth

**Suggested Metrics**:
```
oneoff_jobs_total{status="scheduled|running|completed|failed|cancelled"}
oneoff_jobs_created_total{type="http|shell|docker"}
oneoff_executions_total{status,job_type}
oneoff_execution_duration_seconds{job_type} (histogram)
oneoff_workers_active
oneoff_workers_idle
oneoff_queue_depth
oneoff_api_requests_total{method,path,status}
```

**Success Criteria**:
- `/metrics` returns Prometheus-format metrics
- Grafana can visualize job success/failure rates
- Alert rules can be created for job failures

---

#### Log Retention Enforcement

**Score**: 8.5/10 (Weighted: ICP=8, Core=9, Leverage=9, Diff=6, Clarity=10)

**Why it matters**: `LOG_RETENTION_DAYS` config exists but isn't enforced. Database grows unbounded. This is a config promise not being kept.

| Pros | Cons |
|------|------|
| Config already exists | Data loss concerns |
| Database size management | User expectation setting |
| Trivial implementation | |
| Production hygiene | |

**Current State**:
- `system_config` table has `log_retention_days` key
- No cleanup job runs

**Implementation Approach**:
1. Add cleanup goroutine in worker pool startup
2. Run daily: `DELETE FROM job_executions WHERE completed_at < datetime('now', '-N days')`
3. Log cleanup statistics
4. Respect retention setting from system config

**Success Criteria**:
- Execution logs older than retention period are deleted
- Cleanup runs daily without user intervention
- Cleanup statistics logged

---

### Tier 2: High Value (Developer Experience)

#### Job Chaining API & UI

**Score**: 9.0/10 (Weighted: ICP=9, Core=9, Leverage=8, Diff=10, Clarity=7)

**Why it matters**: Job chaining enables workflow automation—backup → migrate → notify. Database schema exists, making this high-leverage. Major differentiator vs `at` command.

| Pros | Cons |
|------|------|
| Schema and models exist | Complex execution state management |
| High-demand feature | UI design complexity |
| Major differentiator | Error handling across chain |
| Enables workflow automation | |

**Current State**:
- `job_chains` and `job_chain_links` tables exist (migration 000001)
- `domain.JobChain` and `domain.JobChainLink` models defined
- No service layer, no handlers, no UI

**Implementation Approach**:
1. Create `internal/service/chain_service.go` with CRUD and execution logic
2. Add `internal/handler/chain.go` with REST endpoints
3. Chain execution: execute jobs in sequence, stop on failure if configured
4. Create `src/views/Chains.vue` with chain list
5. Create `src/components/ChainBuilder.vue` with drag-and-drop job ordering
6. Add execution view showing chain progress

**Success Criteria**:
- Users can create chains of existing jobs
- Chains execute in sequence with proper status tracking
- Stop-on-failure behavior works correctly
- UI shows chain progress in real-time

---

#### CLI Tool

**Score**: 8.5/10 (Weighted: ICP=9, Core=8, Leverage=6, Diff=8, Clarity=9)

**Why it matters**: Developers love CLI tools. Enables CI/CD integration, scripting, and Unix-philosophy workflows.

| Pros | Cons |
|------|------|
| Developer favorite | Another binary to maintain |
| CI/CD integration | API key management needed |
| Scriptable workflows | Shell completion complexity |
| Unix philosophy | |

**Current State**: No CLI tool exists.

**Implementation Approach**:
1. Create `cmd/oneoff-cli/main.go` using `urfave/cli/v3` (already in go.mod)
2. Commands: `jobs`, `executions`, `projects`, `tags`, `config`
3. Support `--api-key` flag and `ONEOFF_API_KEY` env var
4. Output formats: table (default), JSON, YAML
5. Consider building into main binary with `oneoff cli` subcommand

**Suggested Commands**:
```bash
oneoff jobs list [--status=scheduled] [--project=default]
oneoff jobs create --name "Backup" --type shell --config '{...}'
oneoff jobs run <job-id>
oneoff jobs cancel <job-id>
oneoff executions list --job-id <id>
oneoff config get
oneoff config set workers_count 4
```

**Success Criteria**:
- All major operations available via CLI
- JSON output for scripting
- Works with API key authentication
- Included in release artifacts

---

#### OpenAPI Documentation

**Score**: 8.0/10 (Weighted: ICP=8, Core=7, Leverage=6, Diff=7, Clarity=9)

**Why it matters**: API discoverability enables integrations, SDK generation, and reduces support burden.

| Pros | Cons |
|------|------|
| API discoverability | Annotation maintenance burden |
| SDK generation possible | Documentation drift risk |
| Interactive testing | |
| Standard practice | |

**Current State**: No OpenAPI spec.

**Implementation Approach**:
1. Create `api/openapi.yaml` manually (cleaner than annotation-based)
2. Embed in binary and serve at `/api/docs`
3. Use Swagger UI or Redoc for interactive documentation
4. Generate spec from handler signatures where possible

**Success Criteria**:
- OpenAPI 3.0 spec available at `/api/openapi.yaml`
- Interactive docs at `/api/docs`
- All endpoints documented with examples

---

### Tier 3: Strategic (Platform Integration)

#### Helm Chart

**Score**: 8.5/10 (Weighted: ICP=9, Core=7, Leverage=6, Diff=8, Clarity=9)

**Why it matters**: Kubernetes is the deployment target for ICP 3/4. Official Helm chart removes friction.

| Pros | Cons |
|------|------|
| Standard K8s deployment | Version sync burden |
| Easy adoption | Chart maintenance |
| Community contribution magnet | |

**Current State**: No Helm chart.

**Implementation Approach**:
1. Create `charts/oneoff/` directory
2. Basic templates: Deployment, Service, PVC, Ingress, ConfigMap
3. Values: replicas, resources, persistence, ingress, env vars
4. Publish to artifact hub or GitHub Pages

**Chart Features**:
- Configurable resource limits
- Persistent volume for SQLite
- Ingress with TLS
- Environment variable injection
- Prometheus ServiceMonitor (optional)

**Success Criteria**:
- `helm install oneoff ./charts/oneoff` works
- All configuration via values.yaml
- Documentation for common scenarios

---

#### Comprehensive Test Suite

**Score**: 8.0/10 (Weighted: ICP=6, Core=9, Leverage=5, Diff=6, Clarity=8)

**Why it matters**: Tests enable confident refactoring, catch regressions, and welcome contributors.

| Pros | Cons |
|------|------|
| Quality confidence | Time investment |
| Contribution-friendly | Test maintenance |
| Refactoring safety | |
| CI/CD integration | |

**Current State**: No tests.

**Implementation Approach**:
1. Unit tests: Services with mocked repositories
2. Integration tests: Repository with in-memory SQLite
3. E2E tests: Full API workflow tests with `net/http/httptest`
4. Frontend tests: Vitest + Vue Test Utils for critical flows

**Test Strategy**:
```
internal/service/*_test.go     - Unit tests with mocks
internal/repository/*_test.go  - Integration tests with :memory: SQLite
internal/handler/*_test.go     - HTTP handler tests with httptest
src/**/__tests__/*.spec.js     - Vue component tests
```

**Success Criteria**:
- 70%+ code coverage
- Tests run in CI on every PR
- Critical paths covered (job creation, execution, webhooks)

---

#### Job Retry Policies

**Score**: 8.0/10 (Weighted: ICP=9, Core=8, Leverage=6, Diff=7, Clarity=8)

**Why it matters**: Transient failures (network blips, rate limits) shouldn't require manual re-scheduling.

| Pros | Cons |
|------|------|
| Resilience for transient failures | Retry vs new execution semantics |
| HTTP executor has retry internally | Max retry limits needed |
| Common user request | Exponential backoff complexity |

**Current State**:
- HTTP executor has hardcoded 3 retries internally
- No job-level retry configuration
- No retry tracking in execution records

**Implementation Approach**:
1. Add `max_retries`, `retry_delay_seconds` to job config
2. Track `retry_count` in execution record
3. On failure, check if retries remaining
4. Schedule retry with exponential backoff: `delay * 2^attempt`
5. Cap at configurable maximum delay

**Success Criteria**:
- Jobs can be configured with retry policy
- Retries use exponential backoff
- Retry count visible in execution history
- Max retries respected

---

### Tier 4: Future Consideration

#### Recurring Jobs (Cron Mode)

**Score**: 7.5/10

**Why it matters**: Expands addressable market significantly. Many users want both one-time and recurring.

**Identity Consideration**: This changes OneOff's core identity. Implement as optional "recurring mode" that can be disabled. Keep one-time as the primary experience.

**Implementation Notes**:
- Add `cron_expression` field to jobs (nullable)
- Add `next_run_at` computed field
- After execution, compute next run if cron expression exists
- UI toggle between "One-time" and "Recurring"

---

#### OAuth2/OIDC Integration

**Score**: 7.0/10

**Why it matters**: Enterprise requirement for SSO.

**Implementation Notes**:
- Support Google, GitHub, Okta providers
- Session management with secure cookies
- Map OAuth claims to internal permissions
- Consider Dex for unified OIDC

---

#### Kubernetes Job Type

**Score**: 7.0/10

**Why it matters**: Cloud-native execution for platform teams.

**Implementation Notes**:
- New job executor using `client-go`
- Kubeconfig handling (in-cluster or file)
- Job spec as config (image, command, resources)
- Log streaming from pods

---

## Anti-Roadmap: What We Won't Build

Explicitly stating what we won't build protects focus and sets expectations.

### Multi-Database Support

**Why not**: SQLite is a feature, not a limitation. It enables zero-dependency deployment. Adding PostgreSQL/MySQL would require:
- Connection pool management
- Migration framework changes
- Configuration complexity
- Testing matrix explosion

If you need PostgreSQL, you probably need Airflow.

### Real-Time Log Streaming

**Why not**: Job executions are typically seconds to minutes. Full logs available after completion is sufficient. Real-time streaming adds:
- WebSocket infrastructure
- Partial log state management
- Connection handling complexity

### Visual Workflow Designer (DAG Builder)

**Why not**: OneOff is for one-time jobs, not complex DAGs. Job chaining is sequential, not arbitrary graphs. If you need DAGs, use Airflow or Temporal.

### Multi-Region/Global Deployment

**Why not**: OneOff is a single-instance scheduler. SQLite doesn't support multi-writer. Horizontal scaling would require:
- Database migration to PostgreSQL
- Distributed locking
- Leader election
- Significant architecture changes

### Plugin System for Custom Job Types

**Why not**: The three job types (HTTP, Shell, Docker) cover 99% of use cases. A plugin system adds:
- ABI compatibility concerns
- Security review requirements
- Documentation burden
- Testing complexity

If you need custom job types, fork and add them directly.

### Email/SMS Notifications

**Why not**: Webhooks are the universal integration point. Users can:
- Send webhooks to Zapier for email
- Use Slack/Discord webhooks directly
- Build notification logic in receiving service

Building email/SMS infrastructure is out of scope.

### Job Versioning/Rollback

**Why not**: Jobs are ephemeral, not versioned artifacts. If you need versioned job definitions, use:
- Git for job configs
- Terraform for infrastructure-as-code
- CLI scripts with version control

---

## Contributing

We welcome contributions to roadmap items! Here's how to get involved:

### Picking Up a Feature

1. Check [GitHub Issues](https://github.com/meysam81/oneoff/issues) for existing discussions
2. Comment on the issue or create one referencing this roadmap
3. Discuss approach before implementation
4. Submit PR with tests and documentation

### Contribution Guidelines

1. **Read CLAUDE.md** before making changes
2. **Follow existing patterns** - layered architecture, domain separation
3. **Add tests** for new functionality
4. **Update documentation** if APIs change
5. **Small PRs** - one logical change per PR

### Priority Labels

- `P0`: Blocks production use, critical
- `P1`: High value, next release target
- `P2`: Medium value, planned
- `P3`: Nice to have, community-driven

### Feature Proposal Process

1. Open a GitHub Discussion with:
   - Problem statement (what user pain does this solve?)
   - Proposed solution
   - Alternatives considered
   - Which ICP does this serve?
2. Community feedback period (1 week)
3. Maintainer decision with rationale
4. If accepted, create tracking issue

---

## Changelog

| Date | Version | Change |
|------|---------|--------|
| 2025-11-29 | 2.0.0 | Complete rewrite with ICP analysis, feature scoring, anti-roadmap |
| 2025-11-25 | 1.0.1 | Initial roadmap document |

---

*This roadmap is a living document. It reflects current understanding of user needs and technical state. Priorities may shift as we learn more from users and the market.*
