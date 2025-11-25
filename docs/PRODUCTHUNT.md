# ProductHunt Launch Materials

> Complete submission materials for ProductHunt launch

---

## Product Information

### Name
**OneOff**

### Tagline (60 characters max)
```
One-time job scheduler. Single binary. Zero dependencies.
```

Alternative taglines:
- `Schedule once, execute perfectly. No Redis required.`
- `The job scheduler that respects your time and sanity.`
- `Linux "at" command meets modern web UI.`

### Description (260 characters for short description)
```
OneOff is a self-hosted job scheduler for one-time tasks. Schedule HTTP webhooks, shell scripts, and Docker containers to run at specific times. Single binary deployment, SQLite database, modern Vue 3 UI. No Redis, no Postgres, no complexity.
```

---

## Full Description (for ProductHunt post)

### The Problem

Every developer has been there: you need to schedule a one-time task—maybe send a webhook when a trial expires, run a database migration at 2 AM, or trigger a notification at a specific time.

Your options?
- Set up Celery with Redis and Postgres (hours of work)
- Deploy Airflow (overkill for simple tasks)
- Write a cron job and hope you remember to delete it
- Set an alarm on your phone and do it manually

**There had to be a better way.**

### The Solution

OneOff is the antidote to over-engineered job schedulers.

```bash
./oneoff
# That's it. Open localhost:8080. Done.
```

**What you get:**
- Single binary (~15MB) with everything bundled
- SQLite database (no external DB to manage)
- Beautiful Vue 3 dark-mode UI
- Three job types: HTTP, Shell, Docker
- Priority queue for job execution
- Projects and tags for organization
- Real-time execution monitoring

**What you DON'T need:**
- Redis
- PostgreSQL
- Message queues
- Container orchestration
- A PhD in DevOps

### Who is it for?

- **Indie hackers** scheduling trial expirations, payment reminders
- **DevOps engineers** running maintenance tasks at specific times
- **Developers** who need "fire and forget" scheduled jobs
- **Self-hosters** who want simple, dependency-free tools

### The Maker Story

I built OneOff because I was tired of spinning up entire infrastructure stacks just to send a webhook at a specific time.

Every project I worked on had the same pattern: "We just need to schedule this one thing..." and suddenly we're configuring Redis connections, setting up Celery workers, managing database migrations for a task queue we'll use three times a year.

OneOff is what I wished existed: download a binary, run it, schedule your job, move on with your life.

### Key Features

1. **30-Second Setup** — Download, run, done
2. **Single Binary** — Everything embedded, nothing to install
3. **SQLite Database** — No external database required
4. **Modern Web UI** — Dark-mode first, Vue 3 + Naive UI
5. **HTTP Jobs** — Schedule webhooks and API calls
6. **Shell Jobs** — Run scripts and commands
7. **Docker Jobs** — Execute containers on demand
8. **Priority Queue** — Important jobs run first
9. **Open Source** — MIT licensed, self-hosted

---

## Gallery Images

Upload these SVG files (or convert to PNG):

1. **Hero Banner** (`assets/producthunt/hero-banner.svg`)
   - Main product showcase with UI mockup

2. **Feature: Single Binary** (`assets/producthunt/feature-single-binary.svg`)
   - Comparison: complex setup vs OneOff

3. **Feature: Job Types** (`assets/producthunt/feature-job-types.svg`)
   - HTTP, Shell, and Docker job examples

4. **Feature: Quick Start** (`assets/producthunt/feature-quickstart.svg`)
   - Terminal showing 60-second setup

---

## First Comment (Maker Comment)

```
Hey Product Hunt! 👋

I'm excited to share OneOff—a job scheduler built for developers who value simplicity.

**The backstory:** I was working on a SaaS project and needed to send webhook notifications when user trials expired. My options were either a full Celery/Redis setup or hacking together cron jobs. Neither felt right for one-time scheduled tasks.

So I built OneOff: a single binary you can download and run. That's it.

**What makes it different:**
- No Redis, no Postgres, no message queues
- SQLite database embedded—zero configuration
- Modern Vue 3 UI because CLI-only is so 2010
- Three job types: HTTP webhooks, shell scripts, Docker containers

**Some use cases:**
- Schedule trial expiration emails
- Run database backups at 2 AM
- Trigger deployment webhooks at maintenance windows
- Send yourself reminders via Slack/Discord

It's MIT licensed and fully open source. I'd love your feedback!

Try it: `curl -fsSL https://github.com/meysam81/oneoff/releases/latest/download/oneoff_Linux_x86_64.tar.gz | tar xz && ./oneoff`

What one-time tasks would you schedule?
```

---

## Topics/Categories

Primary: **Developer Tools**

Secondary:
- Open Source
- Productivity
- Self-Hosted
- Task Management

---

## Links

- **Website/Landing**: `https://github.com/meysam81/oneoff`
- **GitHub**: `https://github.com/meysam81/oneoff`

---

## Promotional Tweets

### Launch Tweet
```
🚀 Just launched OneOff on Product Hunt!

A one-time job scheduler that respects your sanity.

✅ Single binary (~15MB)
✅ SQLite database
✅ Modern Vue 3 UI
✅ HTTP, Shell, Docker jobs
✅ Zero dependencies

No Redis. No Postgres. No complexity.

https://producthunt.com/posts/oneoff
```

### Technical Tweet
```
Built a job scheduler in Go + Vue 3 that:

• Bundles frontend in the binary (go:embed)
• Uses SQLite for zero-config storage
• Runs a worker pool with priority queue
• Has automatic migrations
• Works offline

One command: ./oneoff

Open source: github.com/meysam81/oneoff
```

### Problem/Solution Tweet
```
The old way to schedule a webhook:
1. Set up Redis
2. Configure Postgres
3. Deploy Celery workers
4. Write task code
5. Pray it works

The OneOff way:
1. ./oneoff
2. Schedule job in UI
3. Done

github.com/meysam81/oneoff
```

---

## Upvote Request Template

```
Hey [Name]!

I just launched my open-source project OneOff on Product Hunt. It's a simple job scheduler for one-time tasks—think "cron but with a UI and no setup."

Would really appreciate an upvote if you have a moment:
[ProductHunt Link]

No pressure at all—just wanted to share what I've been working on!
```

---

## Post-Launch Checklist

- [ ] Respond to every comment within 1 hour
- [ ] Share launch on Twitter/X
- [ ] Post in relevant Discord servers (Developer communities)
- [ ] Share on Reddit (r/selfhosted, r/golang, r/webdev)
- [ ] Post on Hacker News (Show HN)
- [ ] Update GitHub README with ProductHunt badge
- [ ] Thank everyone who supported

---

## FAQ Responses

**Q: How is this different from cron?**
> OneOff is for one-time scheduled tasks with a modern UI. Cron is great for recurring jobs, but managing one-off cron entries is tedious and error-prone. OneOff gives you a visual interface, execution history, and automatic cleanup.

**Q: Why SQLite instead of a "real" database?**
> SQLite is a real database! It's perfect for single-node deployments. No connection strings to manage, no separate process to run, and it handles thousands of jobs without breaking a sweat. OneOff is designed to be simple.

**Q: Can I use this in production?**
> Absolutely. SQLite handles write-ahead logging, OneOff has graceful shutdown, and the worker pool is battle-tested. For high-availability, you'd want a clustered solution, but for 99% of one-time job needs, OneOff is production-ready.

**Q: What about recurring jobs?**
> OneOff is specifically designed for one-time scheduled tasks. For recurring jobs, cron or dedicated schedulers like Temporal are better fits. We believe in doing one thing well.

**Q: Is there a hosted version?**
> Not yet! OneOff is designed for self-hosting. A hosted version might come in the future based on demand.
