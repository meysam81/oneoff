# Directory Submission Materials

> Ready-to-use descriptions and metadata for directory submissions

---

## Table of Contents

1. [AlternativeTo](#alternativeto)
2. [Awesome Self-Hosted](#awesome-self-hosted)
3. [Awesome Go](#awesome-go)
4. [GitHub Topics](#github-topics)
5. [DevHunt](#devhunt)
6. [Hacker News](#hacker-news)
7. [Reddit](#reddit)
8. [Dev.to](#devto)
9. [Hashnode](#hashnode)
10. [Discord Communities](#discord-communities)

---

## AlternativeTo

**URL**: https://alternativeto.net/

### Product Name

OneOff

### Short Description (160 chars)

```
One-time job scheduler. Single binary, SQLite, Vue 3 UI. Schedule HTTP requests, shell scripts, Docker containers. Zero dependencies.
```

### Full Description

```
OneOff is a self-hosted, developer-focused job scheduler designed specifically for one-time tasks. Unlike complex workflow engines, OneOff embraces simplicity: download a single binary, run it, and schedule your jobs through a modern web interface.

Key Features:
• Single binary deployment (~15MB) with embedded frontend
• SQLite database requiring zero configuration
• Modern Vue 3 dark-mode web interface
• Three job types: HTTP webhooks, shell scripts, Docker containers
• Priority queue system (1-10 levels)
• Project organization and tag system
• Real-time worker monitoring
• Job execution history and logs
• Automatic database migrations
• Graceful shutdown handling

Perfect for scheduling trial expirations, database backups, deployment triggers, webhook notifications, and any task that needs to run at a specific time without the overhead of enterprise schedulers.

Built with Go backend and Vue 3 frontend. MIT licensed and fully open source.
```

### Alternatives To

- Celery
- Airflow
- Rundeck
- Temporal
- cron
- at (Linux command)
- Sidekiq
- Bull (Node.js)

### Tags

```
Job Scheduler, Task Scheduler, Self-Hosted, Open Source, Developer Tools, DevOps, Webhook, Automation, SQLite, Go, Vue.js
```

### Platforms

- Linux
- macOS
- Windows
- Docker
- Self-Hosted

---

## Awesome Self-Hosted

**URL**: https://github.com/awesome-selfhosted/awesome-selfhosted

### Category

`Software Development - Task Runners and Job Schedulers`

### Entry Format

```markdown
- [OneOff](https://github.com/meysam81/oneoff) - One-time job scheduler with web UI for scheduling HTTP requests, shell scripts, and Docker containers. Single binary, SQLite database, zero dependencies. `MIT` `Go/Docker`
```

### Pull Request Description

```markdown
## Add OneOff to Task Runners and Job Schedulers

OneOff is a self-hosted job scheduler specifically designed for one-time tasks (as opposed to recurring cron jobs).

### Key differentiators:

- Single binary deployment (~15MB)
- Embedded SQLite database (no external DB required)
- Modern Vue 3 web interface
- Three job types: HTTP, Shell, Docker
- Zero dependencies

### Why it belongs in awesome-selfhosted:

1. Fully self-hosted with no external dependencies
2. Single binary makes deployment trivial
3. Fills a specific niche (one-time scheduled tasks)
4. Active development, MIT licensed

Website: https://github.com/meysam81/oneoff
Demo: N/A (self-hosted only)
```

---

## Awesome Go

**URL**: https://github.com/avelino/awesome-go

### Category

`Job Scheduler`

### Entry Format

```markdown
- [OneOff](https://github.com/meysam81/oneoff) - One-time job scheduler with embedded Vue 3 web UI, SQLite database, and support for HTTP, shell, and Docker jobs.
```

---

## GitHub Topics

Add these topics to the repository:

```
job-scheduler
task-scheduler
scheduler
one-time-jobs
cron-alternative
self-hosted
sqlite
vue3
golang
devops
automation
webhook-scheduler
docker
open-source
developer-tools
```

---

## DevHunt

**URL**: https://devhunt.org/

### Name

OneOff

### Tagline

```
One-time job scheduler. Single binary. Zero dependencies.
```

### Description

```
OneOff is what happens when you're tired of setting up Redis + Postgres + Celery just to schedule a webhook.

Download a single binary. Run it. Schedule your job. Done.

Features:
- Single binary (~15MB) with embedded Vue 3 UI
- SQLite database (zero configuration)
- HTTP jobs (webhooks, API calls)
- Shell jobs (scripts, commands)
- Docker jobs (run containers on demand)
- Priority queue (1-10 levels)
- Projects and tags for organization

Built for developers who value simplicity over enterprise complexity.
```

### Categories

- Developer Tools
- Open Source
- Productivity

---

## Hacker News

**URL**: https://news.ycombinator.com/

### Show HN Post

**Title**:

```
Show HN: OneOff – One-time job scheduler with zero dependencies
```

**Post Text**:

```
Hey HN,

I built OneOff because I was tired of setting up Redis, Postgres, and Celery just to schedule a single webhook.

OneOff is a job scheduler specifically for one-time tasks:
- Single binary (~15MB) with embedded Vue 3 frontend
- SQLite database (no external DB)
- Three job types: HTTP, Shell, Docker
- Modern web UI for scheduling and monitoring

The idea is simple: download, run, schedule your job, move on.

Some use cases:
- Trial expiration notifications
- Scheduled database backups
- Deployment triggers at maintenance windows
- One-time webhook notifications

Tech stack: Go backend (no framework, just net/http), Vue 3 + Naive UI frontend (embedded via go:embed), SQLite with golang-migrate.

What do you think? I'd love feedback on whether this fills a real gap or if I'm solving a problem nobody has.

GitHub: https://github.com/meysam81/oneoff
```

---

## Reddit

### r/selfhosted

**Title**: `OneOff - Self-hosted one-time job scheduler (single binary, SQLite, Vue 3 UI)`

**Post**:

```markdown
Hey r/selfhosted!

I built **OneOff**, a job scheduler designed specifically for one-time tasks.

**The problem**: Every time I needed to schedule a webhook or run a script at a specific time, my options were either:

1. Set up a full Celery/Redis stack (overkill)
2. Write a cron job and remember to delete it (messy)
3. Set a phone alarm (embarrassing)

**The solution**: A single binary you download and run. That's it.

**Features**:

- Single binary (~15MB) with embedded frontend
- SQLite database (no external DB needed)
- Modern Vue 3 web interface
- HTTP jobs (webhooks, API calls)
- Shell jobs (scripts, backups)
- Docker jobs (run containers)
- Priority queue

**Self-hosting details**:

- No dependencies
- Works offline
- Data stored in a single SQLite file
- Easy backup (just copy the .db file)

GitHub: https://github.com/meysam81/oneoff

Would love your feedback!
```

### r/golang

**Title**: `OneOff - Job scheduler built with Go, featuring embedded Vue 3 frontend via go:embed`

**Post**:

```markdown
Built a one-time job scheduler in Go with some interesting patterns:

**Architecture highlights**:

- Frontend embedded using `go:embed` - single binary distribution
- SQLite with `golang-migrate` for automatic migrations
- Worker pool pattern with priority queue
- Job registry with plugin-like architecture for job types
- Standard library `net/http` for routing (no framework)

**Tech stack**:

- Go 1.23+ backend
- Vue 3 + Pinia + Naive UI frontend
- SQLite with `mattn/go-sqlite3`
- `zerolog` for structured logging
- `urfave/cli/v3` for CLI

The goal was maximum simplicity: no Redis, no Postgres, no message queues. Just a binary that works.

GitHub: https://github.com/meysam81/oneoff

Interested in feedback on the architecture and any Go patterns I could improve!
```

### r/webdev

**Title**: `Built a Vue 3 frontend embedded in a Go binary - one-time job scheduler`

---

## Dev.to

### Article Title

```
I built a job scheduler that doesn't require a PhD in DevOps
```

### Tags

`go`, `vue`, `opensource`, `devops`, `showdev`

### Article Outline

```markdown
# I built a job scheduler that doesn't require a PhD in DevOps

## The Problem

[Story about needing to schedule a simple webhook and facing complex solutions]

## What I Built

[OneOff introduction with screenshot]

## Key Design Decisions

### 1. Single Binary

[Explain go:embed and why it matters]

### 2. SQLite > Everything

[Why SQLite is perfect for this use case]

### 3. No Framework

[Using Go's standard library for HTTP]

## How It Works

[Architecture diagram and explanation]

## Try It

[Quick start instructions]

## What's Next

[Roadmap and call for feedback]
```

---

## Hashnode

Same content as Dev.to, adapted for Hashnode's format.

### Tags

`Go`, `Vue.js`, `Open Source`, `DevOps`, `Self-Hosted`

---

## Discord Communities

### Servers to Post

1. **Go Discord** (https://discord.gg/golang)
   - Channel: #showcase

2. **Vue Land** (https://chat.vuejs.org/)
   - Channel: #show-and-tell

3. **Self-Hosted** (https://discord.gg/selfhosted)
   - Channel: #tools-and-applications

4. **DevOps** communities
   - Various servers, look for showcase channels

### Template Message

```
Hey everyone! 👋

Just released **OneOff** - an open-source one-time job scheduler.

**What it does**: Schedule HTTP requests, shell scripts, and Docker containers to run at specific times.

**Why it's different**: Single binary (~15MB), SQLite database, no external dependencies. Download, run, done.

Built with Go backend + Vue 3 frontend (embedded in binary).

GitHub: https://github.com/meysam81/oneoff

Would love any feedback!
```

---

## Directory Submission Checklist

### Week 1

- [ ] Submit to ProductHunt
- [ ] Post on Hacker News (Show HN)
- [ ] Submit to AlternativeTo
- [ ] Post on r/selfhosted
- [ ] Post on r/golang

### Week 2

- [ ] Submit PR to awesome-selfhosted
- [ ] Submit PR to awesome-go
- [ ] Publish Dev.to article
- [ ] Publish Hashnode article
- [ ] Share in Discord communities

### Week 3

- [ ] Submit to DevHunt
- [ ] Post on r/webdev
- [ ] Share on Twitter/X with relevant hashtags
- [ ] Reach out to newsletter curators

### Ongoing

- [ ] Monitor and respond to all comments
- [ ] Update descriptions based on feedback
- [ ] Track which sources drive the most traffic/stars

---

## SEO Keywords

Include these in descriptions where natural:

- one-time job scheduler
- self-hosted job scheduler
- simple job scheduler
- no dependency job scheduler
- single binary scheduler
- sqlite job scheduler
- webhook scheduler
- cron alternative
- celery alternative
- task scheduler open source
- developer job scheduler
- devops automation tool
