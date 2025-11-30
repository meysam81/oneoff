# OneOff Product Roadmap

> **Last Updated**: 2025-11-29
> **Version**: 2.0.0
> **Status**: Active Development

This roadmap outlines the strategic direction for OneOff, a self-hosted, zero-dependency job scheduler for one-time task execution. It reflects the actual implementation state and prioritizes features that strengthen our core value proposition.

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

OneOff has achieved **production readiness** at v1.0.2. The foundation is complete:

- **Security**: API key authentication with scopes, webhook HMAC signing
- **Observability**: Prometheus-compatible metrics, comprehensive execution logging
- **Operations**: Log retention enforcement, running job cancellation, graceful shutdown
- **Developer Experience**: Full-featured Vue 3 UI, RESTful API, landing page with job catalog

**Key Focus Areas Going Forward:**

1. **Complete the partially-implemented job chaining** (database schema exists, needs API/UI)
2. **Developer tooling** (CLI, OpenAPI docs, Helm chart)
3. **Test infrastructure** (zero tests currently exist)
4. **Strategic extensions** (cron support, OAuth, Kubernetes job type)

---

## Core Value Proposition

**OneOff is the antidote to over-engineering scheduled tasks.**

For developers and DevOps teams who need to schedule one-time tasks—webhooks, scripts, containers—OneOff provides a single binary with zero external dependencies. Download, run, schedule. No Redis, no Postgres, no message queues, no Kubernetes manifests for "just send a webhook at 3 PM."

**In one sentence**: OneOff does for scheduled tasks what SQLite did for databases—removes the infrastructure tax.

---

## Ideal Customer Profiles

### ICP 1: The Pragmatic DevOps Engineer

**Who they are**: Senior DevOps/SRE at a 10-100 person startup, responsible for infrastructure and developer tooling. Values simplicity and maintainability over feature completeness.

**Their context**: Manages multiple services, often needs to schedule one-time operational tasks (database migrations, cache warmups, deployment triggers) without adding infrastructure complexity.

**Jobs to be Done (JTBD)**:

1. Schedule operational tasks at specific times without complex setup
2. Integrate job scheduling into CI/CD pipelines via API
3. Monitor job executions and receive alerts on failures

**Pain points with alternatives**:

- Celery/Airflow require Redis/Postgres and significant operational overhead
- `at` command has no visibility, API, or centralized management
- Building custom solutions takes time and creates tech debt

**What "exceptional" looks like**:

- Deploy in <5 minutes, including production hardening
- API-first design that integrates with existing automation
- Slack/webhook notifications without additional services

---

### ICP 2: The Backend Developer

**Who they are**: Full-stack or backend developer at an early-stage startup. Wears many hats, optimizes for shipping speed.

**Their context**: Building a SaaS product, needs to schedule user-facing tasks (trial expirations, reminder emails, report generation) without dedicated infrastructure team.

**Jobs to be Done (JTBD)**:

1. Schedule application-triggered tasks (trial expiry, notifications) via simple API calls
2. Clone and reschedule jobs for similar future events
3. View execution history and debug failed jobs

**Pain points with alternatives**:

- Cloud schedulers (CloudWatch Events, Cloud Scheduler) are vendor lock-in
- Cron jobs lack visibility and require server access
- In-app scheduling adds complexity to the main application

**What "exceptional" looks like**:

- Create jobs programmatically with one API call
- Rich execution logs accessible via UI
- Self-hosted means no data leaves their infrastructure

---

### ICP 3: The Small Team Lead

**Who they are**: Engineering lead or solo developer running a small technical operation (agency, consultancy, internal tools team).

**Their context**: Manages recurring client work or internal automation, needs reliable task scheduling without enterprise complexity or costs.

**Jobs to be Done (JTBD)**:

1. Organize scheduled tasks by project/client
2. Set up job chains for multi-step workflows
3. Share access with team members (when auth is enabled)

**Pain points with alternatives**:

- Enterprise tools (Rundeck, Temporal) are overkill
- Consumer tools (Zapier, IFTTT) have limited customization
- Managing multiple `crontab` files across servers is error-prone

**What "exceptional" looks like**:

- Project-based organization out of the box
- Job chaining for simple workflows
- Web UI accessible to non-technical stakeholders

---

## Current State

### Fully Implemented

| Feature                      | Description                                                      | Quality          |
| ---------------------------- | ---------------------------------------------------------------- | ---------------- |
| **Job CRUD**                 | HTTP, Shell, Docker job types with full lifecycle                | Production-ready |
| **API Key Authentication**   | Scoped keys (read/write/admin), expiration, rotation             | Production-ready |
| **Webhook Notifications**    | Job lifecycle events, HMAC signing, retry with backoff           | Production-ready |
| **Prometheus Metrics**       | Jobs, workers, requests, API key validations, webhook deliveries | Production-ready |
| **Running Job Cancellation** | Context-based cancellation for all job types                     | Production-ready |
| **Log Retention**            | Configurable cleanup with daily scheduler                        | Production-ready |
| **Project Organization**     | Full UI with CRUD, colors, archival                              | Production-ready |
| **Tag System**               | Job categorization with default tags                             | Production-ready |
| **Priority Queue**           | 1-10 priority levels, higher executes first                      | Production-ready |
| **Worker Pool**              | Configurable workers, graceful shutdown                          | Production-ready |
| **Landing Page**             | Astro-based marketing site with job template catalog             | Production-ready |
| **Immediate Execution**      | `scheduled_at: "now"` support                                    | Production-ready |
| **Job Cloning**              | Duplicate jobs with new schedule                                 | Production-ready |

### Partially Implemented

| Feature          | Current State                                                                    | Remaining Work                             |
| ---------------- | -------------------------------------------------------------------------------- | ------------------------------------------ |
| **Job Chaining** | Database schema (`job_chains`, `job_chain_links` tables) and domain models exist | Service layer, API handlers, UI components |

### Not Yet Implemented

| Feature               | Priority | Notes                                     |
| --------------------- | -------- | ----------------------------------------- |
| CLI Tool              | High     | Critical for CI/CD integration            |
| Test Suite            | High     | Zero tests currently—risk for refactoring |
| OpenAPI Documentation | Medium   | Enables SDK generation                    |
| Helm Chart            | Medium   | K8s deployment standard                   |
| Job Retry Mechanism   | Medium   | Application-level retry for all job types |
| Recurring Jobs (Cron) | Medium   | Significant scope expansion               |
| OAuth2/OIDC           | Low      | Enterprise requirement                    |
| Kubernetes Job Type   | Low      | Cloud-native extension                    |
| RBAC                  | Low      | Requires OAuth first                      |

---

## Competitive Positioning

### Current Landscape

| Alternative          | Strengths                  | Weaknesses                                            | OneOff Advantage                           |
| -------------------- | -------------------------- | ----------------------------------------------------- | ------------------------------------------ |
| **`at` command**     | Zero setup, ubiquitous     | No UI, no API, no persistence, no visibility          | Full observability, API-first, modern UI   |
| **cron**             | Battle-tested, everywhere  | Recurring only, no one-time focus, distributed config | One-time optimized, centralized management |
| **Celery**           | Powerful, Python ecosystem | Redis required, complex setup, heavy                  | Zero dependencies, language agnostic       |
| **Airflow**          | Enterprise DAGs, UI        | 1GB+ memory, steep learning curve                     | 10x simpler, <30MB memory                  |
| **Temporal**         | Durable workflows          | Complex, requires expertise                           | Immediate productivity                     |
| **Cloud Schedulers** | Managed, scalable          | Vendor lock-in, cost at scale                         | Self-hosted, data sovereignty              |

### Our Moat

1. **Single binary architecture**: Competitors cannot easily replicate the go-embed approach with Vue UI
2. **SQLite-first design**: No external state store means simpler operations than any alternative
3. **One-time focus**: We optimize for the neglected use case—scheduled tasks that run once
4. **Developer experience**: Modern UI + API + future CLI creates a complete toolchain

### Positioning Statement

> For DevOps engineers and backend developers who need to schedule one-time tasks, OneOff is the self-hosted job scheduler that deploys in 30 seconds with zero dependencies, unlike enterprise schedulers that require database servers and message queues.

---

## Phased Roadmap

### Phase 1: Developer Tooling

**Outcome**: OneOff integrates seamlessly into CI/CD pipelines and developer workflows

**Target ICPs**: DevOps Engineer, Backend Developer

| Priority | Feature                   | Score | Effort | Status  | Notes                                       |
| -------- | ------------------------- | ----- | ------ | ------- | ------------------------------------------- |
| P0       | **CLI Tool**              | 9.2   | M      | Planned | Enables scriptable workflows                |
| P0       | **OpenAPI Documentation** | 8.5   | S      | Planned | Enables SDK generation, API discoverability |
| P1       | **Test Suite (Core)**     | 8.0   | M      | Planned | Unit + integration tests for services       |
| P1       | **Helm Chart**            | 8.0   | S-M    | Planned | K8s deployment standard                     |
| P2       | **Job Templates in UI**   | 7.5   | S      | Planned | Save/load job configurations                |

---

### Phase 2: Workflow Power

**Outcome**: OneOff handles multi-step workflows, not just single jobs

**Target ICPs**: Small Team Lead, Backend Developer

| Priority | Feature                 | Score | Effort | Status  | Notes                                |
| -------- | ----------------------- | ----- | ------ | ------- | ------------------------------------ |
| P0       | **Job Chaining API**    | 9.5   | S      | Partial | Database schema exists               |
| P0       | **Job Chaining UI**     | 9.0   | M      | Planned | Depends on API completion            |
| P1       | **Job Retry Mechanism** | 8.5   | S      | Planned | Configurable retry for all job types |
| P2       | **Conditional Chains**  | 7.0   | M      | Planned | Execute next job based on conditions |

---

### Phase 3: Enterprise Readiness

**Outcome**: Enterprise adoption unlocked with auth, audit, and cloud-native features

**Target ICPs**: DevOps Engineer (enterprise), Small Team Lead (growing)

| Priority | Feature                      | Score | Effort | Status  | Notes                    |
| -------- | ---------------------------- | ----- | ------ | ------- | ------------------------ |
| P1       | **Recurring Jobs (Cron)**    | 7.5   | L      | Planned | Significant scope change |
| P1       | **OAuth2/OIDC**              | 7.0   | L      | Planned | SSO support              |
| P1       | **Test Suite (E2E)**         | 7.0   | M      | Planned | Frontend + API tests     |
| P2       | **Kubernetes Job Type**      | 7.0   | L      | Planned | Execute K8s Jobs         |
| P2       | **RBAC**                     | 6.5   | L      | Planned | Requires OAuth first     |
| P3       | **Structured Audit Logging** | 6.0   | M      | Planned | Who did what, when       |

---

### Phase 4: Ecosystem

**Outcome**: OneOff becomes a platform with infrastructure-as-code and extensibility

**Target ICPs**: DevOps Engineer (advanced)

| Priority | Feature                        | Score | Effort | Status | Notes                |
| -------- | ------------------------------ | ----- | ------ | ------ | -------------------- |
| P2       | **Terraform Provider**         | 6.5   | L      | Future | IaC for jobs         |
| P2       | **Database Query Job Type**    | 6.5   | M      | Future | SQL execution        |
| P3       | **AWS Lambda/Cloud Functions** | 6.0   | M      | Future | Serverless execution |
| P3       | **Plugin System**              | 5.5   | XL     | Future | Custom job types     |
| P3       | **Horizontal Scaling**         | 4.0   | XL     | Future | Requires PostgreSQL  |

---

## Detailed Feature Analysis

### Tier 1: Critical

#### CLI Tool

**Score**: 9.2/10

**Why it matters**: DevOps engineers and backend developers need to integrate job scheduling into scripts, CI/CD pipelines, and automation workflows. A CLI is the standard interface for these use cases.

| Pros                      | Cons                        |
| ------------------------- | --------------------------- |
| Enables CI/CD integration | Another binary to maintain  |
| Scriptable workflows      | API coverage decisions      |
| Unix philosophy alignment | Auth token management       |
| Developer favorite        | Shell completion complexity |

**Current State**: Not implemented

**Implementation Approach**:

- Use same `urfave/cli/v3` framework as main binary
- Subcommands: `jobs`, `executions`, `projects`, `tags`, `config`
- Output formats: table (default), JSON, YAML
- Auth via `--api-key` flag or `ONEOFF_API_KEY` env var
- Config file support: `~/.oneoff/config.yaml`

**Success Criteria**:

- `oneoff jobs create --type http --config '{"url": "..."}' --at "2025-01-15T09:00:00Z"`
- `oneoff jobs list --status scheduled --format json | jq`
- `oneoff executions tail <job-id>` streams logs

---

#### Job Chaining API & UI

**Score**: 9.5/10 (API), 9.0/10 (UI)

**Why it matters**: Multi-step workflows are a natural extension of one-time tasks. The database schema and domain models already exist—this is the highest-ROI feature to complete.

| Pros                      | Cons                           |
| ------------------------- | ------------------------------ |
| Database schema exists    | Complex execution state        |
| Domain models implemented | UI design complexity           |
| High-demand feature       | Error handling across chain    |
| Major differentiator      | Circular dependency prevention |

**Current State**:

- `job_chains` table with `id`, `name`, `project_id`, `created_at`
- `job_chain_links` table with `chain_id`, `job_id`, `sequence_order`, `stop_on_failure`
- `JobChain` and `JobChainLink` domain models in `internal/domain/models.go`

**Implementation Approach**:

1. **Service Layer**: `internal/service/chain_service.go`
   - `CreateChain`, `GetChain`, `ListChains`, `DeleteChain`
   - `AddJobToChain`, `RemoveJobFromChain`, `ReorderChain`
   - `ExecuteChain` with stop-on-failure logic
2. **Handlers**: `internal/handler/chain.go`
   - `GET/POST /api/chains`
   - `GET/PATCH/DELETE /api/chains/:id`
   - `POST /api/chains/:id/execute`
3. **UI Components**:
   - `ChainBuilder.vue` with drag-and-drop ordering
   - Chain execution progress view
   - Chain execution history

**Success Criteria**:

- Create chain via API with ordered job list
- Execute chain with automatic sequential execution
- Chain stops on failure when `stop_on_failure: true`
- UI shows chain progress in real-time

---

#### OpenAPI Documentation

**Score**: 8.5/10

**Why it matters**: API discoverability enables third-party integrations, SDK generation, and reduces support burden.

| Pros                            | Cons                   |
| ------------------------------- | ---------------------- |
| API discoverability             | Annotation maintenance |
| SDK generation (Go, Python, JS) | Swagger UI bundle size |
| Interactive testing             | Schema drift risk      |
| Industry standard               |                        |

**Current State**: Not implemented

**Implementation Approach**:

- Use `swaggo/swag` for Go comment-based OpenAPI generation
- Embed Swagger UI in frontend (separate route `/docs`)
- Generate spec at build time, embed in binary
- Version spec with API version

**Success Criteria**:

- `/api/openapi.json` returns valid OpenAPI 3.0 spec
- `/docs` serves interactive Swagger UI
- All endpoints documented with request/response schemas

---

### Tier 2: High Value

#### Test Suite (Core)

**Score**: 8.0/10

**Why it matters**: Zero tests currently exist. This creates risk for any refactoring and limits contribution confidence.

| Pros                  | Cons                 |
| --------------------- | -------------------- |
| Refactoring safety    | Time investment      |
| Contribution-friendly | Maintenance overhead |
| CI/CD integration     | Mock complexity      |
| Quality confidence    |                      |

**Current State**: No test files exist

**Implementation Approach**:

1. **Unit Tests** (first priority):
   - Service layer tests with mock repository
   - Job executor tests with mock HTTP/Docker
   - Domain model validation tests
2. **Integration Tests**:
   - Repository tests with in-memory SQLite
   - Handler tests with `httptest`
3. **Framework**: Standard `testing` package + `testify` for assertions

**Success Criteria**:

- > 70% coverage on service layer
- All job executors have test coverage
- CI runs tests on every PR

---

#### Job Retry Mechanism

**Score**: 8.5/10

**Why it matters**: Transient failures are common in distributed systems. HTTP jobs already have internal retry via `imroc/req`, but Shell and Docker jobs do not.

| Pros                                 | Cons                      |
| ------------------------------------ | ------------------------- |
| Resilience for transient failures    | Backoff complexity        |
| Consistent behavior across job types | Retry vs new execution UX |
| Common user request                  | Max retry config needed   |

**Current State**: HTTP jobs have internal retry; Shell and Docker do not

**Implementation Approach**:

- Add `max_retries` and `retry_delay_seconds` to job config schema
- Track retry count in execution record
- Exponential backoff: `delay * 2^attempt`
- Cap at configurable maximum delay
- Option: retry on specific exit codes only

**Success Criteria**:

- Job with `max_retries: 3` retries up to 3 times on failure
- Retry attempts visible in execution history
- Exponential backoff between attempts

---

#### Helm Chart

**Score**: 8.0/10

**Why it matters**: Kubernetes is the deployment standard for production workloads. A Helm chart removes friction for K8s users.

| Pros                          | Cons                       |
| ----------------------------- | -------------------------- |
| K8s deployment standard       | Maintenance burden         |
| Easy adoption                 | Version sync with releases |
| Community contribution magnet |                            |

**Current State**: Docker image exists at `ghcr.io/meysam81/oneoff`

**Implementation Approach**:

- Create `charts/oneoff/` directory
- Support: Deployment, Service, Ingress, PVC, ConfigMap, Secret
- Values: replicas (1 for SQLite), resources, ingress config, auth settings
- Publish to artifact hub

**Success Criteria**:

- `helm install oneoff oci://ghcr.io/meysam81/oneoff-chart`
- Supports ingress with TLS
- Configurable persistent volume for SQLite

---

### Tier 3: Strategic

#### Recurring Jobs (Cron)

**Score**: 7.5/10

**Why it matters**: Recurring jobs significantly expand the addressable market, but also dilute the "one-off" identity.

| Pros                             | Cons                        |
| -------------------------------- | --------------------------- |
| Market expansion                 | Identity dilution           |
| Higher daily active usage        | Cron parsing complexity     |
| Feature parity with alternatives | Overlap with existing tools |

**Implementation Notes**:

- Consider as optional mode: `schedule_type: once | cron`
- Cron parsing via `robfig/cron/v3`
- Clear UI separation between one-time and recurring
- May warrant separate "OneOff Scheduler" branding

**Decision**: Defer until Phase 3; evaluate based on user feedback

---

#### OAuth2/OIDC Integration

**Score**: 7.0/10

**Why it matters**: Enterprise users require SSO. API keys alone are insufficient for teams using centralized identity.

| Pros                    | Cons                   |
| ----------------------- | ---------------------- |
| Enterprise requirement  | Significant complexity |
| Modern auth expectation | Session management     |
| RBAC foundation         | Multi-provider support |

**Implementation Notes**:

- Support Google, GitHub, generic OIDC providers
- Use `coreos/go-oidc` library
- Store user identity, link to API key generation
- Consider: is this core, or better as a proxy (oauth2-proxy)?

---

### Tier 4: Future Consideration

| Feature                     | Score | Notes                                      |
| --------------------------- | ----- | ------------------------------------------ |
| **Terraform Provider**      | 6.5   | IaC for jobs, requires stable API          |
| **Database Query Job Type** | 6.5   | SQL execution with result capture          |
| **Kubernetes Job Type**     | 7.0   | Execute K8s Jobs, requires kubeconfig      |
| **AWS Lambda Integration**  | 6.0   | Invoke Lambda functions as jobs            |
| **Horizontal Scaling**      | 4.0   | Requires PostgreSQL migration—major change |
| **Plugin System**           | 5.5   | Custom job type registration at runtime    |

---

## Anti-Roadmap: What We Won't Build

Explicitly stating what we will **not** build protects focus and sets expectations.

### 1. Multi-Database Support

**Why not**: SQLite is a feature, not a limitation. Adding PostgreSQL/MySQL support would:

- Increase operational complexity for users
- Require maintaining multiple code paths
- Dilute our "zero dependencies" value proposition

**Exception**: If horizontal scaling becomes critical, we would add PostgreSQL as a separate deployment mode, not a general option.

### 2. Full Workflow Orchestration (Airflow-style DAGs)

**Why not**: Complex DAG orchestration is a different product. We build job chains (linear sequences), not arbitrary graphs.

**What we will build**: Sequential job chains with stop-on-failure semantics

### 3. Built-in Email Notifications

**Why not**: Email requires SMTP configuration and is easily handled by webhooks + external services (SendGrid, Mailgun, etc.)

**What we will build**: Webhook notifications that can trigger email via external services

### 4. Multi-tenancy

**Why not**: OneOff is designed for single-team use. Multi-tenancy adds complexity (tenant isolation, billing, quotas) that serves a different market.

**Alternative**: Deploy multiple OneOff instances

### 5. Real-time Streaming Logs

**Why not**: Requires WebSocket infrastructure and significantly complicates the frontend. Polling is sufficient for our use case.

**What we will build**: Efficient polling with execution log pagination

### 6. Mobile App

**Why not**: Our users are developers at desks. Mobile adds significant development and maintenance burden for low value.

**What we will build**: Responsive web UI that works on tablets

---

## Contributing

We welcome contributions! Here's how to get involved:

### Picking Up a Roadmap Item

1. Check [GitHub Issues](https://github.com/meysam81/oneoff/issues) for existing work
2. Comment on an issue or create one referencing this roadmap
3. Discuss approach before implementation (especially for M+ effort items)
4. Submit PR with tests and documentation

### Priority Labels

- `P0`: Critical, blocks phase outcome
- `P1`: High value, significant impact
- `P2`: Medium value, nice-to-have
- `P3`: Future consideration, opportunistic

### Effort Estimates

| Size   | Definition                       |
| ------ | -------------------------------- |
| **XS** | < 1 day, isolated change         |
| **S**  | 1-3 days, single module          |
| **M**  | 3-7 days, multiple modules       |
| **L**  | 1-2 weeks, new subsystem         |
| **XL** | 2+ weeks, significant complexity |

### Good First Issues

- OpenAPI annotations for existing endpoints
- Test coverage for domain models
- CLI subcommand implementation
- Helm chart values documentation

---

## Changelog

| Date       | Version | Change                                                                                     |
| ---------- | ------- | ------------------------------------------------------------------------------------------ |
| 2025-11-29 | 2.0.0   | Complete rewrite reflecting actual implementation state; added ICPs, scoring, anti-roadmap |
| 2025-11-25 | 1.0.1   | Initial roadmap document created                                                           |

---

_This roadmap is a living document. It reflects our current understanding and priorities, which will evolve based on user feedback and market conditions._
