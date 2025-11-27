# OneOff Product Roadmap

> **Last Updated**: 2025-11-25
> **Version**: 1.0.1
> **Status**: Active Development

This document outlines the strategic product roadmap for OneOff, a self-hosted job scheduler for one-time task execution.

---

## Table of Contents

- [Executive Summary](#executive-summary)
- [Current State](#current-state)
- [Phased Roadmap](#phased-roadmap)
- [Detailed Feature Analysis](#detailed-feature-analysis)
- [Differentiation Strategy](#differentiation-strategy)
- [Contributing](#contributing)

---

## Executive Summary

OneOff is a well-architected, production-ready job scheduler at v1.0.1. The core is solid—clean Go backend, modern Vue 3 frontend, single-binary deployment. To elevate this from "good OSS project" to "must-have tool," we focus on:

1. **Completing existing foundations** (job chaining, projects UI, job cancellation)
2. **Production essentials** (authentication, webhooks, metrics)
3. **Workflow capabilities** (chains, retries, CLI tooling)
4. **Enterprise readiness** (OAuth, RBAC, Kubernetes)

---

## Current State

### Fully Implemented

- Job CRUD with HTTP, Shell, and Docker executors
- Priority-based worker pool with graceful shutdown
- Project and tag organization system
- Execution history with full audit trail
- Modern Vue 3 UI with dark mode
- SQLite with automatic migrations
- Real-time worker monitoring
- Job cloning and immediate execution

### Partially Implemented

- **Job Chaining**: Database schema exists, API not exposed
- **Projects UI**: Backend complete, frontend shows placeholder
- **Job Cancellation**: Works for scheduled jobs, not running jobs

### Not Yet Implemented

- Authentication/Authorization
- Webhook notifications
- Prometheus metrics
- Recurring jobs (cron)
- Job retry mechanism
- CLI tool
- Test suite

---

## Phased Roadmap

### Phase 1: Foundation Completion (4-6 weeks)

**Theme: "Production Ready"**

| Priority | Feature                       | Effort   | Status  |
| -------- | ----------------------------- | -------- | ------- |
| P0       | API Authentication (API Keys) | Medium   | Planned |
| P0       | Webhook Notifications         | Medium   | Planned |
| P0       | Running Job Cancellation      | Low      | Planned |
| P1       | Prometheus Metrics            | Low      | Planned |
| P1       | Log Retention Enforcement     | Very Low | Planned |
| P1       | Projects UI Completion        | Low      | Planned |

**Outcome**: OneOff becomes truly production-deployable with security, observability, and operational completeness.

---

### Phase 2: Workflow Power (6-8 weeks)

**Theme: "Beyond One-Off"**

| Priority | Feature               | Effort     | Status  |
| -------- | --------------------- | ---------- | ------- |
| P0       | Job Chaining API & UI | Medium     | Planned |
| P1       | Job Retry Mechanism   | Low-Medium | Planned |
| P1       | CLI Tool              | Medium     | Planned |
| P1       | Helm Chart            | Low-Medium | Planned |
| P2       | OpenAPI Documentation | Low        | Planned |

**Outcome**: OneOff becomes a lightweight workflow tool, not just a scheduler.

---

### Phase 3: Enterprise Ready (8-12 weeks)

**Theme: "Scale & Integrate"**

| Priority | Feature                   | Effort | Status  |
| -------- | ------------------------- | ------ | ------- |
| P0       | OAuth2/OIDC Integration   | High   | Planned |
| P1       | Recurring Jobs (Cron)     | High   | Planned |
| P1       | Kubernetes Job Type       | High   | Planned |
| P1       | Comprehensive Test Suite  | High   | Planned |
| P2       | Role-Based Access Control | High   | Planned |

**Outcome**: Enterprise adoption unlocked with proper auth, scale, and cloud-native features.

---

### Phase 4: Ecosystem (Ongoing)

**Theme: "Platform Play"**

| Priority | Feature                    | Effort    | Status |
| -------- | -------------------------- | --------- | ------ |
| P1       | Terraform Provider         | High      | Future |
| P2       | Database Query Job Type    | Medium    | Future |
| P2       | AWS Lambda/Cloud Functions | Medium    | Future |
| P3       | Horizontal Scaling         | Very High | Future |
| P3       | S3/GCS Log Archival        | Medium    | Future |

---

## Detailed Feature Analysis

### Tier 1: Critical (Highest Priority)

#### 1. Webhook Notifications

**Score: 10/10**

Execute HTTP callbacks on job completion, failure, or state changes.

| Pros                             | Cons                                 |
| -------------------------------- | ------------------------------------ |
| Most requested feature class     | Need retry logic for failed webhooks |
| Enables integration ecosystem    | Adds complexity to execution flow    |
| Slack/Discord/Teams via webhooks | Failure handling edge cases          |
| Low implementation complexity    |                                      |

**Implementation Notes:**

- Add `webhooks` table with URL, events, secret for HMAC signing
- Fire webhooks asynchronously after job state transitions
- Include job details, execution output, timestamps in payload
- Implement exponential backoff for failed deliveries

---

#### 2. API Authentication (API Keys)

**Score: 10/10**

Simple API key authentication for all endpoints.

| Pros                                 | Cons                         |
| ------------------------------------ | ---------------------------- |
| Essential for any production use     | Adds operational complexity  |
| Prevents unauthorized job creation   | Need to manage key lifecycle |
| Low complexity to implement          | UI needs auth flow           |
| Table stakes for enterprise adoption |                              |

**Implementation Notes:**

- Add `api_keys` table with hashed keys, scopes, expiration
- Middleware to validate `Authorization: Bearer <key>` header
- UI login flow with key input or generation
- Scope-based permissions (read, write, admin)

---

#### 3. Job Chaining API & UI

**Score: 9/10**

Expose existing job chain backend via API, build UI for sequential job execution.

| Pros                                 | Cons                             |
| ------------------------------------ | -------------------------------- |
| Database & models already exist      | Complex UI to design             |
| High-demand feature (pipelines)      | Needs execution state management |
| Major differentiator vs `at` command | Error handling across chains     |
| Enables workflow automation          |                                  |

**Implementation Notes:**

- Add handlers: `GET/POST/DELETE /api/chains`
- Chain execution service with stop-on-failure logic
- UI: Chain builder with drag-and-drop job ordering
- Execution view showing chain progress

---

#### 4. Running Job Cancellation

**Score: 9/10**

Implement context cancellation for in-progress jobs.

| Pros                           | Cons                             |
| ------------------------------ | -------------------------------- |
| Critical for long-running jobs | Docker container kills need care |
| Users expect cancel to work    | Shell process cleanup tricky     |
| Current behavior is confusing  | May leave orphaned processes     |
| Simple to implement in worker  |                                  |

**Implementation Notes:**

- Pass cancellable context to job executors
- Shell: Kill process group with SIGTERM, then SIGKILL
- Docker: `docker stop` with timeout
- HTTP: Cancel in-flight request via context

---

#### 5. Prometheus Metrics

**Score: 9/10**

Export metrics via `/metrics` endpoint in Prometheus format.

| Pros                            | Cons                             |
| ------------------------------- | -------------------------------- |
| DevOps standard                 | Need to choose metrics carefully |
| Enables alerting pipelines      | Another endpoint to maintain     |
| Dashboard integration (Grafana) | Cardinality concerns with labels |
| SRE adoption driver             |                                  |

**Suggested Metrics:**

```
oneoff_jobs_total{status="scheduled|running|completed|failed"}
oneoff_jobs_scheduled_count
oneoff_executions_total{status, job_type}
oneoff_execution_duration_seconds{job_type}
oneoff_workers_active
oneoff_workers_idle
oneoff_queue_depth
```

---

### Tier 2: High Value (Next Quarter)

#### 6. Job Retry Mechanism

**Score: 8/10**

Configurable automatic retry on job failure.

| Pros                              | Cons                             |
| --------------------------------- | -------------------------------- |
| Resilience for transient failures | Exponential backoff complexity   |
| HTTP jobs have it, others don't   | Max retry limits needed          |
| Common user request               | Retry vs new execution ambiguity |
| Simple configuration              |                                  |

**Implementation Notes:**

- Add `max_retries`, `retry_delay_seconds` to job config
- Track retry count in execution record
- Exponential backoff: delay \* 2^attempt
- Cap at configurable maximum delay

---

#### 7. Projects UI Completion

**Score: 8/10**

Build frontend UI for projects management (backend already complete).

| Pros                            | Cons                               |
| ------------------------------- | ---------------------------------- |
| Backend 100% complete           | Relatively low user impact         |
| Quick win (2 days work)         | Users can use API directly         |
| Removes "coming soon" text      | Basic feature, not differentiating |
| Improves perceived completeness |                                    |

**Implementation Notes:**

- Create `Projects.vue` with CRUD operations
- Project list with color/icon display
- Archive/unarchive functionality
- Job count per project

---

#### 8. CLI Tool

**Score: 8/10**

Command-line tool for job management: `oneoff-cli create`, `list`, `run`.

| Pros                 | Cons                       |
| -------------------- | -------------------------- |
| Developer favorite   | Another binary to maintain |
| CI/CD integration    | API coverage decisions     |
| Scriptable workflows | Auth token handling        |
| Unix philosophy      |                            |

**Suggested Commands:**

```bash
oneoff-cli jobs list [--status=scheduled] [--project=default]
oneoff-cli jobs create --name "Backup" --type shell --config '{...}'
oneoff-cli jobs run <job-id>
oneoff-cli jobs cancel <job-id>
oneoff-cli executions list --job-id <id>
oneoff-cli config get
oneoff-cli config set workers_count 4
```

---

#### 9. Helm Chart

**Score: 8/10**

Official Kubernetes deployment chart.

| Pros                          | Cons                       |
| ----------------------------- | -------------------------- |
| Standard K8s deployment       | Maintenance burden         |
| Easy adoption                 | Version sync with releases |
| Community contribution magnet |                            |

**Chart Features:**

- Configurable replicas (though single-instance for now)
- Persistent volume for SQLite
- Ingress configuration
- Environment variable injection
- Resource limits/requests

---

#### 10. Log Retention Enforcement

**Score: 8/10**

Actually delete old execution logs per `LOG_RETENTION_DAYS` setting.

| Pros                                  | Cons                     |
| ------------------------------------- | ------------------------ |
| Config exists, implementation doesn't | Data loss concerns       |
| Database size management              | User expectation setting |
| Very simple implementation            |                          |

**Implementation Notes:**

- Add cleanup job to worker pool startup
- Run daily: `DELETE FROM job_executions WHERE completed_at < NOW() - retention_days`
- Log cleanup statistics
- Consider soft-delete option

---

### Tier 3: Strategic (Following Quarter)

#### 11. OpenAPI/Swagger Documentation

**Score: 7/10**

Auto-generated API documentation with try-it-out functionality.

| Pros                    | Cons                   |
| ----------------------- | ---------------------- |
| API discoverability     | Annotation maintenance |
| SDK generation possible |                        |
| Standard practice       |                        |
| Interactive testing     |                        |

---

#### 12. Recurring Jobs (Cron-like)

**Score: 7/10**

Schedule jobs to run on cron patterns.

| Pros                       | Cons                            |
| -------------------------- | ------------------------------- |
| Massive market expansion   | Changes core "one-off" identity |
| Competes with cron/systemd | Cron parsing complexity         |
| Higher daily active usage  | Overlap with existing solutions |

**Identity Consideration:** This significantly expands scope. Consider as optional "recurring mode" rather than core feature.

---

#### 13. OAuth2/OIDC Integration

**Score: 7/10**

SSO support (Google, GitHub, Okta, etc.)

| Pros                     | Cons                      |
| ------------------------ | ------------------------- |
| Enterprise requirement   | Significant complexity    |
| Modern auth expectations | Session management        |
| RBAC foundation          | Multiple provider support |

---

#### 14. Kubernetes Job Type

**Score: 7/10**

Execute Kubernetes Jobs as a job type.

| Pros                          | Cons                |
| ----------------------------- | ------------------- |
| Cloud-native adoption         | K8s complexity      |
| Natural evolution from Docker | Kubeconfig handling |
| Enterprise demand             | RBAC concerns       |

---

#### 15. Test Suite

**Score: 7/10**

Comprehensive unit, integration, and E2E tests.

| Pros                  | Cons             |
| --------------------- | ---------------- |
| Quality confidence    | Time investment  |
| Contribution-friendly | Test maintenance |
| Refactoring safety    |                  |
| CI/CD integration     |                  |

**Test Strategy:**

- Unit tests: Services, handlers with mocks
- Integration tests: Repository with in-memory SQLite
- E2E tests: Full API workflow tests
- Frontend: Vitest + Vue Test Utils

---

### Tier 4: Future Consideration

#### 16. Job Templates

**Score: 6/10**

Save job configurations as reusable templates.

---

#### 17. Role-Based Access Control (RBAC)

**Score: 6/10**

User roles (admin, operator, viewer) with granular permissions.

---

#### 18. Structured Audit Logging

**Score: 6/10**

Who did what, when—comprehensive audit trail.

---

#### 19. Database Query Job Type

**Score: 6/10**

Execute SQL queries on schedule (reports, maintenance).

---

#### 20. Terraform Provider

**Score: 6/10**

Manage jobs as infrastructure code.

---

#### 21. Horizontal Scaling

**Score: 4/10**

Multi-node deployment with shared state.

**Note:** This requires PostgreSQL/MySQL and fundamentally changes the architecture. Consider only if single-node limits are reached.

---

## Differentiation Strategy

### Current Unique Value

1. **Single binary** deployment (vs. complex multi-service setups)
2. **Zero dependencies** (vs. Redis/Postgres requirements)
3. **Modern UI** (vs. CLI-only tools)
4. **Plugin architecture** (vs. monolithic schedulers)

### Target Unique Value

1. **"Job Chains Made Simple"** - Workflow capability without Airflow complexity
2. **"Batteries Included"** - Auth, metrics, webhooks out-of-box
3. **"5-Minute Production Deploy"** - Helm/Docker with sensible defaults
4. **"Developer-First"** - CLI, Terraform, API-first design

### Competitive Positioning

| Tool         | Complexity | Use Case        | OneOff Advantage            |
| ------------ | ---------- | --------------- | --------------------------- |
| `at` command | Low        | Simple one-time | Modern UI, API, persistence |
| `cron`       | Low        | Recurring only  | One-time focus, visibility  |
| Celery       | High       | Python apps     | Language agnostic, simpler  |
| Airflow      | Very High  | Complex DAGs    | 10x simpler deployment      |
| Temporal     | Very High  | Workflows       | No learning curve           |

---

## Contributing

We welcome contributions to any roadmap item! Here's how to get involved:

### Picking Up a Feature

1. Check [GitHub Issues](https://github.com/meysam81/oneoff/issues) for existing discussions
2. Comment on the issue or create one referencing this roadmap
3. Discuss approach before implementation
4. Submit PR with tests and documentation

### Proposing New Features

1. Open a GitHub Discussion or Issue
2. Reference this roadmap document
3. Include: problem statement, proposed solution, alternatives considered
4. Community feedback period before acceptance

### Priority Labels

- `P0`: Critical, blocks production use
- `P1`: High value, next release target
- `P2`: Medium value, planned for future
- `P3`: Nice to have, community-driven

---

## Changelog

| Date       | Change                           |
| ---------- | -------------------------------- |
| 2025-11-25 | Initial roadmap document created |

---

_This roadmap is a living document and will be updated as priorities evolve and features are completed._
