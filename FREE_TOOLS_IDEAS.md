# OneOff — Free Tools Strategy

> **Last Updated**: 2025-11-29
> **Status**: Planning
> **Newsletter CTA**: newsletter.meysam.io

## Executive Summary

This document outlines a strategy to build standalone micro-tools that attract developers and DevOps engineers to the OneOff ecosystem. Each tool solves a specific, contained problem related to scheduling, webhooks, or job automation—building trust and brand awareness before users need a full job scheduler. The only CTA is newsletter subscription; no account creation, no paywalls.

---

## Table of Contents

- [Project Context](#project-context)
- [Market Analysis](#market-analysis)
- [Free Tool Ideas](#free-tool-ideas)
- [Prioritization Matrix](#prioritization-matrix)
- [Implementation Roadmap](#implementation-roadmap)
- [Newsletter Integration](#newsletter-integration)
- [Success Metrics](#success-metrics)
- [Anti-Ideas: What We Won't Build](#anti-ideas-what-we-wont-build)

---

## Project Context

**Core Value Proposition**: OneOff is the antidote to over-engineering—a single binary, zero-dependency job scheduler for one-time tasks.

**Primary ICPs**:

- **ICP 1 (Pragmatic DevOps Engineer)**: Senior DevOps/SRE at 10-100 person startup, schedules operational tasks (migrations, cache warmups, deployment triggers) without adding infrastructure complexity
- **ICP 2 (Backend Developer)**: Full-stack developer at early-stage startup, needs to schedule user-facing tasks (trial expirations, reminders, reports) via simple API calls
- **ICP 3 (Small Team Lead)**: Engineering lead or solo developer running small operations, needs reliable task scheduling without enterprise complexity

**Core Problem Space**: Scheduling one-time tasks (HTTP webhooks, shell scripts, Docker containers) without Redis, Postgres, message queues, or cloud vendor lock-in

**Tech Stack**:

- Backend: Go 1.24+ (urfave/cli, SQLite, zerolog)
- Frontend: Vue 3 + Vite
- Landing Page: Astro 5 + Tailwind
- Free tools can be: Static HTML/JS, Astro pages, or minimal Go serverless

**Existing Landing Page**: Yes, at `landing-page/` with job template catalog at `/catalog`

---

## Market Analysis

### Core Market

**Core Market Keywords**:

- "schedule one-time task"
- "webhook scheduler"
- "delayed job execution"
- "cron alternative for one-time"
- "at command alternative"
- "schedule http request"
- "delayed api call"
- "run script at specific time"

**Core Market Pain Points**:

- Setting up Redis/Postgres just to send a webhook at 3 PM
- No visibility with `at` command—did the job run? What was the output?
- Cloud scheduler costs add up for simple tasks
- Vendor lock-in with AWS/GCP schedulers
- Cron is for recurring, not one-time—workarounds feel hacky

**Core Market "I wish I had..."**:

- A simple way to test my webhook endpoint before scheduling it
- A tool to visualize when my scheduled job will actually run
- A cron expression explainer that doesn't require a CS degree
- A quick way to build HTTP request payloads for webhooks

### Adjacent Markets

**Adjacent Market 1**: API Developers

- **Overlap**: Building and testing HTTP requests, webhook integrations
- **Entry point**: Testing webhook endpoints, building API payloads, debugging HTTP calls
- **Size**: Massive—every developer building APIs

**Adjacent Market 2**: DevOps/SRE Learning Docker

- **Overlap**: Running containers, understanding Docker commands
- **Entry point**: Docker run command complexity, container configuration
- **Size**: Growing continuously with Docker adoption

**Adjacent Market 3**: SaaS Founders/Product Managers

- **Overlap**: Need scheduled notifications, time-based triggers
- **Entry point**: Understanding cron expressions, timezone scheduling
- **Size**: Medium—overlaps heavily with ICP 2

**Adjacent Market 4**: System Administrators

- **Overlap**: Scheduling maintenance tasks, shell scripts
- **Entry point**: Shell script validation, cron scheduling
- **Size**: Large, traditional market

**Adjacent Market 5**: Freelance Developers/Consultants

- **Overlap**: Automating client work, scheduled reports
- **Entry point**: Project-based task organization, quick automation
- **Size**: Growing gig economy

---

## Free Tool Ideas

### Tool 1: Cron Expression Generator & Explainer

**One-liner**: Build and understand cron expressions visually

**URL slug**: `oneoff.io/tools/cron` or `cron.oneoff.io`

**Target market**: Core + Adjacent (System Administrators, all developers)

**Search intent**: "cron expression generator" — ~40,000 monthly searches

**The job**:

- User has: A scheduling need they want to express ("every Monday at 9 AM")
- User wants: A valid cron expression and confidence it's correct
- Current alternatives: crontab.guru (good but no export), generic generators lack depth

**Why this works as a lead magnet**:

- Universal developer need—even though OneOff focuses on one-time jobs
- Natural gateway: "For recurring, use cron. For one-time, use OneOff"
- Demonstrates we understand scheduling deeply

**Core functionality**:

1. Visual builder with dropdowns for minute/hour/day/month/weekday
2. Natural language input: "every weekday at 9am" → `0 9 * * 1-5`
3. Cron → English explainer: `*/15 * * * *` → "Every 15 minutes"
4. Next 5 execution times preview
5. Output: Copyable cron expression + explanation

**Newsletter hook**: "Get weekly DevOps tips including advanced scheduling patterns"

**Technical approach**:

- Type: Static / Client-side
- Key tech: Pure JavaScript, cron-parser library for validation
- Complexity: S
- Estimated build: 2-3 days

**Virality mechanics**:

- [x] Shareable output (link with encoded expression)
- [x] "Built with OneOff" subtle footer
- [x] Social-friendly preview (OG image showing the expression)

**Risks/Considerations**:

- Crontab.guru is well-established; need a differentiator (better UX, export options)
- Must handle edge cases in cron syntax

---

### Tool 2: Webhook Tester & Inspector

**One-liner**: Send test webhooks and inspect incoming payloads

**URL slug**: `oneoff.io/tools/webhook-tester` or `webhook.oneoff.io`

**Target market**: Core + Adjacent (API Developers)

**Search intent**: "webhook tester" — ~8,000 monthly searches

**The job**:

- User has: A webhook endpoint they need to test
- User wants: To send test payloads and see what their endpoint receives
- Current alternatives: webhook.site (limited free tier), RequestBin (shut down), Postman (heavyweight)

**Why this works as a lead magnet**:

- Directly adjacent to OneOff's HTTP job type
- User testing webhooks → user scheduling webhooks → user needs OneOff
- Shows we understand the webhook workflow

**Core functionality**:

1. **Send mode**: Build and send HTTP requests to any URL
   - Method selector (GET, POST, PUT, DELETE, PATCH)
   - Headers editor (key-value pairs)
   - Body editor with JSON syntax highlighting
   - Response viewer with timing
2. **Receive mode**: Temporary endpoint that captures incoming webhooks
   - Generate unique URL: `webhook.oneoff.io/catch/abc123`
   - Display incoming requests in real-time (polling, not WebSocket)
   - Show headers, body, timing for each request
3. Output: Shareable capture URL, exportable request/response

**Newsletter hook**: "Get notified about webhook best practices and OneOff updates"

**Technical approach**:

- Type: Serverless backend required (for receive mode)
- Key tech: Go serverless function (Cloudflare Workers or similar), 15-min TTL storage
- Complexity: M
- Estimated build: 5-7 days

**Virality mechanics**:

- [x] Shareable output (capture URL to share with teammates)
- [x] "Powered by OneOff" on capture page
- [x] Social-friendly preview (OG showing "Webhook Inspector")

**Risks/Considerations**:

- Requires backend—adds complexity and cost
- Must prevent abuse (rate limiting, short TTL)
- Competitor landscape is active (webhook.site has premium model)

---

### Tool 3: Unix Timestamp Converter

**One-liner**: Convert between Unix timestamps and human dates instantly

**URL slug**: `oneoff.io/tools/timestamp` or `timestamp.oneoff.io`

**Target market**: Adjacent (All developers) + Core

**Search intent**: "unix timestamp converter" — ~60,000 monthly searches

**The job**:

- User has: A Unix timestamp from logs/APIs or a date they need as timestamp
- User wants: Quick, accurate conversion with timezone awareness
- Current alternatives: epochconverter.com (cluttered), unixtimestamp.com (ads everywhere)

**Why this works as a lead magnet**:

- Extremely high search volume—pure SEO play
- Every developer needs this occasionally
- Ties to OneOff's `scheduled_at` timestamp handling

**Core functionality**:

1. Bidirectional conversion: Unix ↔ Human-readable
2. Live "current time" display with auto-updating timestamp
3. Timezone selector (shows time in user's TZ and UTC)
4. Millisecond/second toggle
5. Relative time display ("3 hours ago", "in 2 days")
6. Output: Copyable formats (ISO 8601, RFC 2822, Unix)

**Newsletter hook**: "Developer tools and tips delivered weekly"

**Technical approach**:

- Type: Static / Client-side
- Key tech: Pure JavaScript, date-fns or Day.js
- Complexity: XS
- Estimated build: 1 day

**Virality mechanics**:

- [x] Shareable output (URL with timestamp encoded: `/timestamp?t=1701234567`)
- [x] "By OneOff" footer
- [x] Social-friendly preview

**Risks/Considerations**:

- Very competitive space—differentiation through cleaner UX
- Low barrier to entry means easy to clone

---

### Tool 4: HTTP Request Builder

**One-liner**: Visually build HTTP requests and get curl/code snippets

**URL slug**: `oneoff.io/tools/http-builder`

**Target market**: Core + Adjacent (API Developers, Backend Developers)

**Search intent**: "curl command generator" — ~15,000 monthly searches

**The job**:

- User has: An API they need to call with specific params
- User wants: A properly formatted request in their preferred format
- Current alternatives: Postman (overkill for quick tests), curl syntax is error-prone

**Why this works as a lead magnet**:

- Directly produces what OneOff HTTP jobs consume
- Export includes "Save as OneOff job" option (JSON config)
- Natural workflow: build request → test it → schedule it

**Core functionality**:

1. Visual request builder:
   - URL input with query param editor
   - Method selector
   - Headers table with common presets (Content-Type, Authorization)
   - Body editor (JSON, form-data, raw)
2. Output formats:
   - curl command
   - JavaScript (fetch)
   - Python (requests)
   - Go (net/http)
   - **OneOff job config JSON** (our differentiator)
3. "Test it" button to execute the request live
4. Response viewer with status, headers, body, timing

**Newsletter hook**: "API development tips and scheduling tricks weekly"

**Technical approach**:

- Type: Static with optional CORS proxy for testing
- Key tech: JavaScript, highlight.js for code output
- Complexity: S-M
- Estimated build: 3-4 days

**Virality mechanics**:

- [x] Shareable output (encoded URL with request config)
- [x] "Export as OneOff job" subtle upsell
- [x] Social-friendly preview

**Risks/Considerations**:

- CORS limitations for "test it" feature—may need proxy
- Competition from Postman, Insomnia, Hoppscotch

---

### Tool 5: Docker Run Command Builder

**One-liner**: Generate docker run commands without memorizing flags

**URL slug**: `oneoff.io/tools/docker-run`

**Target market**: Adjacent (DevOps/Docker learners) + Core

**Search intent**: "docker run command generator" — ~5,000 monthly searches

**The job**:

- User has: A container they need to run with specific config
- User wants: A correct `docker run` command without RTFM
- Current alternatives: Manual docs reading, trial and error

**Why this works as a lead magnet**:

- Maps to OneOff's Docker job type
- Export includes OneOff Docker job config
- DevOps engineers learning Docker → future OneOff users

**Core functionality**:

1. Visual builder for common flags:
   - Image selection with tag
   - Port mappings (-p)
   - Volume mounts (-v)
   - Environment variables (-e)
   - Network settings (--network)
   - Resource limits (--memory, --cpus)
   - Restart policy (--restart)
   - Name (--name)
   - Detached mode (-d)
   - Auto-remove (--rm)
2. Command preview with syntax highlighting
3. Output formats:
   - docker run command
   - docker-compose snippet
   - **OneOff Docker job config**
4. Common presets (nginx, postgres, redis, etc.)

**Newsletter hook**: "Container tips and automation patterns weekly"

**Technical approach**:

- Type: Static / Client-side
- Key tech: Pure JavaScript, no external dependencies
- Complexity: S
- Estimated build: 2-3 days

**Virality mechanics**:

- [x] Shareable output (URL with config encoded)
- [x] "Schedule this container with OneOff" CTA
- [x] Social-friendly preview

**Risks/Considerations**:

- Docker CLI is complex—can't cover everything
- Must keep up with Docker version changes

---

### Tool 6: JSON Formatter & Validator

**One-liner**: Format, validate, and minify JSON instantly

**URL slug**: `oneoff.io/tools/json`

**Target market**: Adjacent (All developers)

**Search intent**: "json formatter" — ~200,000 monthly searches

**The job**:

- User has: Messy JSON from an API response or config
- User wants: Readable, validated JSON they can work with
- Current alternatives: jsonformatter.org (ad-heavy), jsonlint.com (basic)

**Why this works as a lead magnet**:

- Massive search volume—highest SEO potential
- Developers working with JSON → working with APIs → scheduling API calls
- Useful for building OneOff HTTP job payloads

**Core functionality**:

1. Paste/type JSON input
2. Auto-format with customizable indentation (2/4 spaces, tabs)
3. Validation with clear error messages and line highlighting
4. Minify option
5. Tree view for exploring nested structures
6. Path copying (click a key to copy JSON path)
7. Output: Formatted JSON, validation status

**Newsletter hook**: "Developer productivity tips delivered weekly"

**Technical approach**:

- Type: Static / Client-side
- Key tech: Pure JavaScript, JSON.parse for validation
- Complexity: XS-S
- Estimated build: 1-2 days

**Virality mechanics**:

- [x] Shareable output (URL with JSON encoded, though length limits apply)
- [x] "By OneOff" subtle footer
- [x] Social-friendly preview

**Risks/Considerations**:

- Extremely competitive—many established tools
- Differentiation must come from speed and cleanliness
- Large JSON handling (performance for 10MB+ files)

---

### Tool 7: Timezone Scheduler

**One-liner**: Find the best time across timezones for scheduled tasks

**URL slug**: `oneoff.io/tools/timezone`

**Target market**: Core + Adjacent (Remote teams, global operations)

**Search intent**: "timezone converter for scheduling" — ~10,000 monthly searches

**The job**:

- User has: A task to schedule that affects people in multiple timezones
- User wants: To visualize what time it will be everywhere when the job runs
- Current alternatives: worldtimebuddy.com (meeting-focused), timeanddate.com (generic)

**Why this works as a lead magnet**:

- Directly relevant to scheduling jobs
- Shows "schedule this task at X:XX" → "here's what time that is globally"
- Natural lead-in: "Now schedule it with OneOff"

**Core functionality**:

1. Add multiple timezones to compare
2. Select a time in any timezone, see it in all others
3. "Golden hours" overlay—when are business hours across all selected zones?
4. Save timezone presets (stored locally)
5. Output: Shareable comparison view, UTC timestamp

**Newsletter hook**: "Global scheduling tips for distributed teams"

**Technical approach**:

- Type: Static / Client-side
- Key tech: JavaScript, Intl API for timezone handling
- Complexity: S
- Estimated build: 2 days

**Virality mechanics**:

- [x] Shareable output (URL with timezones and time encoded)
- [x] "Schedule this time with OneOff" CTA
- [x] Social-friendly preview (visual showing times across zones)

**Risks/Considerations**:

- Timezone handling edge cases (DST transitions)
- Competition from established time tools

---

### Tool 8: Shell Script Validator

**One-liner**: Check shell scripts for syntax errors before running

**URL slug**: `oneoff.io/tools/shellcheck`

**Target market**: Core + Adjacent (System Administrators)

**Search intent**: "shell script validator" — ~5,000 monthly searches

**The job**:

- User has: A shell script they're about to run/schedule
- User wants: Confidence it won't fail due to syntax errors
- Current alternatives: shellcheck.net (good but basic UI), local shellcheck install

**Why this works as a lead magnet**:

- Directly relevant to OneOff Shell jobs
- "Validate → Schedule" natural workflow
- Saves embarrassment of syntax errors at 3 AM

**Core functionality**:

1. Paste shell script
2. Select shell type (bash, sh, zsh)
3. Run validation (client-side if possible, else serverless)
4. Display errors/warnings with line numbers and explanations
5. Common fix suggestions
6. Output: Validation report, fixed script

**Newsletter hook**: "Shell scripting tips and automation patterns"

**Technical approach**:

- Type: Serverless backend required (shellcheck binary)
- Key tech: Go serverless calling shellcheck, or WASM port
- Complexity: M
- Estimated build: 4-5 days

**Virality mechanics**:

- [x] Shareable output (validation report)
- [x] "Schedule this script with OneOff" CTA
- [x] Social-friendly preview

**Risks/Considerations**:

- Requires backend or WASM—adds complexity
- shellcheck.net already exists and works
- Security: never execute user scripts, only validate

---

### Tool 9: API Response Mocker

**One-liner**: Create mock API endpoints that return custom responses

**URL slug**: `oneoff.io/tools/mock-api`

**Target market**: Adjacent (API Developers, Frontend Developers)

**Search intent**: "mock api online" — ~8,000 monthly searches

**The job**:

- User has: A frontend or integration that needs a fake API endpoint
- User wants: A URL that returns their specified JSON response
- Current alternatives: mockapi.io (requires signup), beeceptor.com (limited)

**Why this works as a lead magnet**:

- Tests webhook receivers before connecting real webhooks
- "Mock it → Test it → Schedule the real thing with OneOff"
- Developers building integrations = OneOff target audience

**Core functionality**:

1. Configure response:
   - Status code
   - Headers
   - JSON body
   - Delay (simulate latency)
2. Get a unique URL: `mock.oneoff.io/m/abc123`
3. URL returns configured response for 24 hours
4. View request history to that mock
5. Output: Mock endpoint URL, request logs

**Newsletter hook**: "API development and testing tips weekly"

**Technical approach**:

- Type: Serverless backend required
- Key tech: Go serverless, KV store for mock configs
- Complexity: M
- Estimated build: 4-5 days

**Virality mechanics**:

- [x] Shareable output (mock URL to share with team)
- [x] "Powered by OneOff" on mock response
- [x] Social-friendly preview

**Risks/Considerations**:

- Requires backend—ongoing cost
- Abuse potential (rate limiting essential)
- TTL management (24h default, then cleanup)

---

### Tool 10: Job Schedule Calculator

**One-liner**: See exactly when your scheduled job will execute

**URL slug**: `oneoff.io/tools/schedule-preview`

**Target market**: Core (OneOff users and prospects)

**Search intent**: "when will my job run" — ~2,000 monthly searches (niche)

**The job**:

- User has: A scheduled time for a job
- User wants: To visualize and confirm when it will actually execute
- Current alternatives: Nothing specific for one-time jobs

**Why this works as a lead magnet**:

- Directly demonstrates OneOff's value proposition
- Answers "what time is 2025-01-15T09:00:00Z in my timezone?"
- Could integrate with "now" special value explanation

**Core functionality**:

1. Enter a datetime (ISO 8601, Unix timestamp, or picker)
2. Select timezone
3. Show:
   - Countdown to execution
   - Time in user's local timezone
   - Time in UTC
   - Relative time ("in 3 days, 4 hours")
4. Comparison: "If you schedule now vs. your input"
5. Output: Formatted times, shareable preview

**Newsletter hook**: "Scheduling tips and OneOff feature updates"

**Technical approach**:

- Type: Static / Client-side
- Key tech: JavaScript, date-fns
- Complexity: XS
- Estimated build: 1 day

**Virality mechanics**:

- [x] Shareable output (URL with schedule encoded)
- [x] "Schedule this with OneOff" strong CTA
- [x] Social-friendly preview

**Risks/Considerations**:

- Low search volume—more of a brand tool
- Only valuable to people already thinking about scheduling

---

## Prioritization Matrix

Scoring criteria:

- **Utility (3x)**: How genuinely useful? Would you use it?
- **Strategic Fit (2x)**: How well does it funnel to OneOff?
- **Discovery (2x)**: Can it rank? Will it spread?
- **Build Effort (1x)**: Inverse scoring—lower effort = higher score

**Score = (Utility×3 + Fit×2 + Discovery×2 + (10-Effort)×1) / 8**

| Rank | Tool                      | Utility | Fit | Discovery | Effort | Score     |
| ---- | ------------------------- | ------- | --- | --------- | ------ | --------- |
| 1    | Cron Expression Generator | 9       | 7   | 10        | 3      | **8.5**   |
| 2    | Unix Timestamp Converter  | 8       | 6   | 10        | 1      | **8.125** |
| 3    | JSON Formatter            | 8       | 5   | 10        | 2      | **7.75**  |
| 4    | HTTP Request Builder      | 9       | 9   | 7         | 4      | **8.125** |
| 5    | Webhook Tester            | 9       | 10  | 7         | 6      | **8.125** |
| 6    | Docker Run Builder        | 8       | 8   | 6         | 3      | **7.375** |
| 7    | Timezone Scheduler        | 7       | 8   | 6         | 3      | **7.0**   |
| 8    | API Response Mocker       | 8       | 7   | 7         | 5      | **7.0**   |
| 9    | Shell Script Validator    | 7       | 8   | 5         | 5      | **6.5**   |
| 10   | Job Schedule Calculator   | 6       | 10  | 3         | 1      | **6.625** |

---

## Implementation Roadmap

### Wave 1: Quick Wins (Build First)

Tools that are high-impact AND low-effort. Get these live fast to start capturing traffic.

- **Unix Timestamp Converter** — XS effort, massive search volume (60K/mo), immediate SEO value. Can be live in 1 day.
- **JSON Formatter** — XS-S effort, highest search volume (200K/mo), generic developer utility. 1-2 days.
- **Job Schedule Calculator** — XS effort, strong brand connection. 1 day. Good internal linking target.

**Wave 1 Total Effort**: 3-4 days
**Expected Outcome**: Foundation of tools subdomain, initial traffic, newsletter signup infrastructure tested

### Wave 2: Strategic Bets

Higher effort but strong strategic value and differentiation potential.

- **Cron Expression Generator** — S effort, 40K/mo searches, natural "cron for recurring, OneOff for one-time" positioning. 2-3 days.
- **HTTP Request Builder** — S-M effort, direct pipeline to OneOff HTTP jobs, export-to-OneOff feature. 3-4 days.
- **Docker Run Builder** — S effort, maps to Docker job type, DevOps learner audience. 2-3 days.
- **Timezone Scheduler** — S effort, scheduling-specific value, global team appeal. 2 days.

**Wave 2 Total Effort**: 9-12 days
**Expected Outcome**: Full scheduling toolkit, strong OneOff integration points, differentiated positioning

### Wave 3: Market Expansion

Adjacent market plays requiring backend infrastructure. Build after Waves 1-2 prove the model.

- **Webhook Tester** — M effort, requires serverless backend, high strategic value but operational cost. 5-7 days.
- **API Response Mocker** — M effort, requires backend, strong developer utility. 4-5 days.
- **Shell Script Validator** — M effort, requires shellcheck backend or WASM, niche but loyal audience. 4-5 days.

**Wave 3 Total Effort**: 13-17 days
**Expected Outcome**: Backend-powered tools, expanded adjacent market reach, enterprise-grade tooling perception

### Parking Lot

Ideas considered but not prioritized for implementation:

- **Dockerfile Validator** — Similar to shell validator but even more niche; Hadolint exists
- **YAML Formatter** — Generic, many alternatives, low strategic fit
- **Base64 Encoder/Decoder** — Useful but zero connection to scheduling
- **Regex Tester** — Excellent search volume but no scheduling connection
- **SSL Certificate Checker** — DevOps adjacent but not scheduling related
- **Port Scanner** — Security concerns, brand risk

---

## Newsletter Integration

### Value Proposition by Tool

| Tool                      | Newsletter Hook                                                | Expected Conversion |
| ------------------------- | -------------------------------------------------------------- | ------------------- |
| Cron Expression Generator | "Weekly DevOps tips including advanced scheduling patterns"    | High                |
| Webhook Tester            | "Get notified about webhook best practices and OneOff updates" | High                |
| Unix Timestamp Converter  | "Developer tools and tips delivered weekly"                    | Medium              |
| HTTP Request Builder      | "API development tips and scheduling tricks weekly"            | High                |
| Docker Run Builder        | "Container tips and automation patterns weekly"                | Medium              |
| JSON Formatter            | "Developer productivity tips delivered weekly"                 | Low                 |
| Timezone Scheduler        | "Global scheduling tips for distributed teams"                 | Medium              |
| Shell Script Validator    | "Shell scripting tips and automation patterns"                 | Medium              |
| API Response Mocker       | "API development and testing tips weekly"                      | Medium              |
| Job Schedule Calculator   | "Scheduling tips and OneOff feature updates"                   | High                |

### CTA Placement

**During use** (non-intrusive):

- Subtle banner at bottom: "Like this tool? Get more developer tips → [Subscribe]"
- Show after 3+ uses in a session (localStorage tracking)
- Never interrupt the primary task

**On output**:

- After successful operation: "Get weekly tips for [relevant topic]"
- On shareable links: Footer "Built by OneOff — Subscribe for more tools"
- On export: Include "Learn more at newsletter.meysam.io"

**Exit intent** (optional, A/B test):

- Modal on leaving: "Before you go — join 1,000+ developers getting weekly tips"
- Only show once per week per user
- Clear dismiss option

### Newsletter Content Strategy

Tools users expect content about:

1. **Scheduling deep dives**: Cron patterns, timezone handling, job chaining
2. **Automation recipes**: Real-world examples of scheduled tasks
3. **Tool tips**: Lesser-known features, keyboard shortcuts
4. **OneOff updates**: New features, improvements (soft sell)
5. **Industry news**: DevOps trends, scheduling best practices

---

## Success Metrics

| Metric                        | Target (6 months)             | Measurement        |
| ----------------------------- | ----------------------------- | ------------------ |
| Tool visits (total)           | 10,000/month                  | Pirsch Analytics   |
| Newsletter signups            | 500 subscribers               | Listmonk           |
| Newsletter conversion rate    | 2-5% of tool visitors         | Pirsch + Listmonk  |
| Referral traffic to oneoff.io | 10% of tool visitors          | Pirsch             |
| Search ranking                | Top 10 for 3+ target keywords | Manual/Ahrefs      |
| Tool shares                   | 100+ shared URLs/month        | URL param tracking |
| GitHub stars from tools       | +200 stars                    | GitHub             |

### Per-Tool Success Criteria

**Wave 1 (after 30 days)**:

- Unix Timestamp: 2,000 visits/month
- JSON Formatter: 5,000 visits/month
- Schedule Calculator: 500 visits/month

**Wave 2 (after 60 days)**:

- Cron Generator: 3,000 visits/month
- HTTP Builder: 1,500 visits/month
- Docker Builder: 1,000 visits/month

---

## Anti-Ideas: What We Won't Build

### 1. Generic Developer Tools Without Scheduling Connection

**Examples**: Regex tester, color picker, lorem ipsum generator
**Why not**: Zero strategic value. Dilutes brand positioning. Competes on a battlefield we can't win.

### 2. Tools Requiring User Accounts

**Why not**: Friction kills conversion. The value proposition is instant utility. Newsletter is the only ask.

### 3. Paid Tiers or Freemium Upsells

**Why not**: Tools are a marketing investment, not a revenue stream. Paywalls create resentment and reduce shareability.

### 4. Tools That Cannibalize OneOff Features

**Examples**: Full job scheduler web app, execution history viewer
**Why not**: These ARE the product. Give away complementary value, not core value.

### 5. Tools Requiring Significant Backend Infrastructure

**Exceptions**: Webhook tester and mock API are strategic enough to justify cost
**Why not**: Static tools are free to host, impossible to break, easy to maintain. Backend adds complexity, cost, and failure modes.

### 6. Mobile Apps

**Why not**: Our audience is developers at desks. Mobile adds dev/maintenance burden for minimal value. Responsive web tools suffice.

### 7. Tools in Saturated Markets Without Differentiation

**Examples**: Yet another JSON formatter without a unique angle
**Why we might**: If we can be 10x cleaner/faster, it's worth it for SEO. If we're just "another one," skip it.

---

## Technical Implementation Notes

### Hosting Strategy

All static tools should be hosted on:

- **Cloudflare Pages** (free tier, global CDN, automatic SSL)
- Alternative: Vercel, Netlify

Backend tools (Wave 3) should use:

- **Cloudflare Workers** (serverless, pay-per-use, global)
- **KV storage** for temporary data (webhook captures, mock configs)

### URL Structure

Option A (subdomain):

- `tools.oneoff.io/cron`
- `tools.oneoff.io/json`

Option B (path):

- `oneoff.io/tools/cron`
- `oneoff.io/tools/json`

**Recommendation**: Option B (path) for SEO value consolidation on main domain.

### Analytics

- **Pirsch Analytics** (privacy-focused, already in use)
- Track: page views, tool usage events, newsletter signup clicks, share button clicks

### Design System

Reuse landing page design system:

- Tailwind CSS configuration
- Dark mode first
- Existing UI components (Button, Card, Badge, Container)
- Consistent typography and colors

---

_This document is a living strategy. Update as tools launch and data comes in._
