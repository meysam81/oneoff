# OneOff Landing Page UX/UI Roadmap

> **Last Updated:** 2025-11-28
> **Status:** Active Development
> **Goal:** Reduce time-to-aha to near zero, maximize conversion rates

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Current State Analysis](#current-state-analysis)
3. [Core Problems Identified](#core-problems-identified)
4. [Improvement Roadmap](#improvement-roadmap)
   - [P0: Critical (Must Have)](#p0-critical-must-have)
   - [P1: High Impact](#p1-high-impact)
   - [P2: Nice to Have](#p2-nice-to-have)
5. [Design Principles](#design-principles)
6. [Competitive Analysis](#competitive-analysis)
7. [Implementation Notes](#implementation-notes)
8. [Success Metrics](#success-metrics)

---

## Executive Summary

The OneOff landing page has a **solid technical foundation** - clean dark theme, consistent design system, and functional copy. However, it's leaving significant conversions on the table by violating key UX principles:

| Issue            | Impact   | Current State                                   |
| ---------------- | -------- | ----------------------------------------------- |
| Time-to-Aha      | Critical | Users must download and run before seeing value |
| Social Proof     | Critical | Zero testimonials, logos, or usage metrics      |
| Emotional Hook   | High     | Copy is descriptive but not compelling          |
| Interactive Demo | High     | No way to "try before download"                 |
| Visual Proof     | Medium   | Mock UI instead of real screenshots             |

**Primary Objective:** Make developers feel "Finally! Someone gets it." within 3 seconds of landing.

---

## Current State Analysis

### Landing Page Structure (as of 2025-11-28)

```
index.astro
├── Header.astro      - Fixed nav with scroll effect, mobile hamburger
├── Hero.astro        - Badge, headline, subheadline, CTAs, terminal animation
├── Problem.astro     - 3 pain point cards (cron, Airflow, Cloud)
├── Features.astro    - Bento grid with 6 feature cards
├── HowItWorks.astro  - 3-step process with tabbed code examples
├── Comparison.astro  - Table comparing OneOff vs alternatives
├── Installation.astro - Platform tabs with copy commands
├── CTA.astro         - Final conversion section
└── Footer.astro      - Links, social, legal
```

### Current Copy Analysis

**Hero Headlines:**

- Main: "One-time jobs. Zero dependencies."
- Sub: "The modern `at` command. Schedule HTTP webhooks, shell scripts, and Docker containers to run once at a specific time. Single binary. Self-hosted. Done."

**Assessment:** Descriptive but not emotionally compelling. Tells _what_ not _why_.

**Hero CTAs:**

- Primary: "Download for Linux" (static, not OS-aware)
- Secondary: "View on GitHub" (exit link, not conversion)

**Assessment:** Weak CTA hierarchy. No low-friction option to experience value.

### Visual Design System

```
Colors:
- bg-primary: #0a0a0b (main background)
- bg-secondary: #111113
- accent-primary: #22d3ee (cyan)
- fg-primary: #fafafa
- fg-secondary: #a1a1aa

Fonts:
- Display: JetBrains Mono
- Body: Geist
- Code: JetBrains Mono

Animations:
- Terminal typewriter effect
- Fade-in with staggered delays
- Scroll-triggered reveals
```

### What's Working Well

1. **Dark-first aesthetic** - Appropriate for developer audience
2. **Clean visual hierarchy** - Sections are distinct and scannable
3. **Terminal animations** - On-brand for CLI tool
4. **Bento grid features** - Modern, engaging layout
5. **Platform-specific install commands** - Practical and helpful
6. **Job template catalog** - Reduces time-to-value for specific use cases

### What's Not Working

1. **No social proof** - Zero trust signals
2. **Mock UI instead of real product** - Reduces credibility
3. **Static hero CTA** - Doesn't detect user's OS
4. **No interactive demo** - Friction before value
5. **Cognitive overload** - Too much reading before understanding
6. **Exit-focused secondary CTA** - GitHub link competes with conversion

---

## Core Problems Identified

### Problem 1: Time-to-Aha is Too Long

**Current journey:**

```
Land on page → Read about features → Download → Install → Run → Open browser →
Create first job → See it work → AHA!
```

**Time:** 5-10 minutes minimum

**Target journey:**

```
Land on page → See demo/video → AHA! → Download
```

**Time:** 30 seconds

### Problem 2: No Social Proof

Users subconsciously ask:

- "Is anyone else using this?"
- "Is this maintained?"
- "Can I trust it for production?"

**Current answer:** Nothing. Complete silence.

**Required elements:**

- GitHub stars (prominently displayed)
- Download count
- Testimonial quotes
- Company logos (if available)
- "X jobs scheduled" metrics

### Problem 3: Copy Lacks Emotional Weight

**Current approach:** Feature-focused, descriptive
**Required approach:** Pain-focused, empathetic, provocative

**Developer frustrations we should tap into:**

- "I just want to schedule ONE thing, why do I need Redis?"
- "Airflow is massive overkill for what I need"
- "Cloud functions have weird cold starts and surprising bills"
- "Cron has zero visibility into what ran"

### Problem 4: No Way to Experience Value Before Commitment

**Friction hierarchy (low to high):**

1. Watch video / see screenshots ← Currently missing good execution
2. Try interactive demo ← Currently missing entirely
3. Download and run ← Current minimum requirement

---

## Improvement Roadmap

### P0: Critical (Must Have)

These changes will have the highest impact on conversion and should be implemented first.

---

#### P0.1: Add Social Proof Section

**Location:** Immediately after Hero section

**Implementation:**

```astro
<!-- New file: landing-page/src/components/landing/SocialProof.astro -->

<section class="py-12 border-y border-border-subtle bg-bg-secondary/30">
  <Container size="xl">
    <!-- Stats row -->
    <div class="flex flex-wrap justify-center gap-8 mb-8">
      <div class="text-center">
        <div class="text-2xl font-bold text-fg-primary">{stars}+</div>
        <div class="text-sm text-fg-muted">GitHub Stars</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-bold text-fg-primary">{downloads}+</div>
        <div class="text-sm text-fg-muted">Downloads</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-bold text-fg-primary"><1min</div>
        <div class="text-sm text-fg-muted">Setup Time</div>
      </div>
    </div>

    <!-- Testimonial -->
    <blockquote class="max-w-2xl mx-auto text-center">
      <p class="text-lg text-fg-secondary italic mb-4">
        "Finally, a job scheduler that doesn't require a PhD in DevOps.
        I was up and running in 30 seconds."
      </p>
      <cite class="text-sm text-fg-muted">
        — Developer testimonial (collect real ones)
      </cite>
    </blockquote>
  </Container>
</section>
```

**Data to fetch:**

- GitHub stars: Already available via `getGitHubData()`
- Downloads: Add to `getGitHubData()` using GitHub releases API
- Testimonials: Collect from users, GitHub issues, Twitter mentions

**Priority:** P0 - Highest impact, relatively low effort

---

#### P0.2: Rewrite Hero Copy for Emotional Impact

**Current:**

```
One-time jobs.
Zero dependencies.

The modern `at` command. Schedule HTTP webhooks, shell scripts,
and Docker containers to run once at a specific time.
Single binary. Self-hosted. Done.
```

**Option A: Pain-Focused (Recommended)**

```
Stop wrestling with Celery, Airflow, and cron
for one-time jobs.

Schedule HTTP webhooks, shell scripts, and Docker containers
with a single binary. No Redis. No Postgres. No PhD required.
```

**Option B: Outcome-Focused**

```
Schedule that one-off job in 30 seconds,
not 3 hours.

The developer tool for one-time scheduled tasks.
Single binary. Zero dependencies. Actually simple.
```

**Option C: Provocative**

```
You shouldn't need Redis, Postgres, and a message queue
to send one Slack message tomorrow.

OneOff: Download. Run. Schedule. That's literally it.
```

**Recommendation:** Test Option A first - it acknowledges the pain directly.

**Implementation:**

- Update `Hero.astro` lines 33-43
- A/B test different versions if analytics are available

---

#### P0.3: Fix Hero CTA Hierarchy

**Current:**

```
[Download for Linux]  [View on GitHub ⭐]
```

**Problems:**

1. "Download for Linux" is static (what if I'm on Mac?)
2. "View on GitHub" is an exit, not conversion
3. No low-friction option to experience value

**Recommended:**

```
[▶ Try Live Demo]  [⬇ Download for {OS}]

                   ⭐ 1,234 stars on GitHub

No account required • Runs locally • Your data stays yours
```

**Implementation:**

```astro
<!-- Hero.astro CTA section -->

<!-- OS Detection Script -->
<script>
  function detectOS() {
    const userAgent = navigator.userAgent.toLowerCase();
    if (userAgent.includes('mac')) return 'macOS';
    if (userAgent.includes('win')) return 'Windows';
    return 'Linux';
  }

  document.addEventListener('DOMContentLoaded', () => {
    const osBtn = document.getElementById('download-btn');
    const os = detectOS();
    osBtn.textContent = `Download for ${os}`;
    osBtn.href = `#installation`; // or direct download link
  });
</script>

<!-- CTAs -->
<div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-6">
  <Button href="#demo" size="lg" variant="primary">
    <Icon name="play" size={20} />
    Try Live Demo
  </Button>
  <Button href="#installation" size="lg" variant="secondary" id="download-btn">
    <Icon name="download" size={20} />
    Download for Linux
  </Button>
</div>

<!-- Social proof badge -->
<div class="flex items-center justify-center gap-2 text-fg-muted text-sm mb-6">
  <Icon name="github" size={16} />
  <span>{stars} stars on GitHub</span>
</div>

<!-- Trust badges -->
<p class="text-xs text-fg-muted">
  No account required • Runs locally • Your data stays yours
</p>
```

---

#### P0.4: Add Interactive Demo or Video

**Option A: Embedded Demo Instance (Best UX, Most Effort)**

Host a read-only demo instance that users can explore without downloading.

```
┌─────────────────────────────────────────────────────────────┐
│  [Try OneOff Live]                                          │
│                                                             │
│  Opens modal/new tab with:                                  │
│  - Pre-populated sample jobs                                │
│  - Read-only or sandboxed interaction                       │
│  - No login required                                        │
│  - Clear "This is a demo" indicator                        │
└─────────────────────────────────────────────────────────────┘
```

**Implementation considerations:**

- Could use Docker + temporary instances
- Could use WebContainers (like StackBlitz)
- Could embed iframe to demo.oneoff.sh

**Option B: Interactive Simulation (Medium Effort)**

Create a fake but realistic interactive demo in the browser.

```javascript
// Simulated API responses
const mockJobs = [
  { name: "Slack notification", type: "http", status: "completed" },
  { name: "Database backup", type: "shell", status: "running" },
  { name: "Docker cleanup", type: "docker", status: "scheduled" },
];

// Interactive elements:
// - Click to "create" a job (shows form)
// - See jobs "run" with animations
// - View "execution logs"
```

**Option C: Product Video (Lower Effort, Still Effective)**

Create a 30-60 second video showing:

1. Terminal: Download command
2. Terminal: `./oneoff` starts
3. Browser: Dashboard loads
4. Browser: Create a Slack notification job
5. Browser: Job executes, shows success
6. Total time elapsed: "< 60 seconds"

**Placement:** Either:

- Replace terminal animation in Hero
- Add as "See it in action" section after Problem
- Modal triggered by "Watch Demo" CTA

**Recommendation:** Start with Option C (video) as it's quickest to implement, then consider Option A for maximum impact.

---

### P1: High Impact

---

#### P1.1: Replace Mock UI with Real Screenshots

**Current state:** Features section shows a fake UI mockup

**Location:** `Features.astro` lines 143-196

**Required:**

1. Take high-quality screenshots of actual OneOff UI
2. Show dashboard with real (but sanitized) job data
3. Show job creation flow
4. Show execution history

**Implementation:**

```bash
# Screenshot checklist:
1. Dashboard overview (3 workers active, jobs list)
2. Job creation modal (HTTP type selected)
3. Execution details (successful job with output)
4. Jobs list with mixed statuses
```

**Screenshot requirements:**

- 2x resolution for retina displays
- Consistent window size
- Dark mode (matches landing page)
- Interesting but realistic data
- No sensitive information

**File locations:**

```
landing-page/public/screenshots/
├── dashboard.png
├── create-job.png
├── execution-detail.png
└── jobs-list.png
```

---

#### P1.2: Add Use Cases Section

**Purpose:** Help users see themselves in the product

**Location:** After Problem section, before Features

**Implementation:**

```astro
<!-- New file: landing-page/src/components/landing/UseCases.astro -->

<section class="section">
  <Container size="xl">
    <div class="section-header">
      <h2 class="section-title">What developers schedule with OneOff</h2>
    </div>

    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
      {useCases.map(useCase => (
        <div class="flex items-start gap-3 p-4 bg-bg-secondary/50 rounded-lg">
          <span class="text-2xl">{useCase.emoji}</span>
          <div>
            <h3 class="font-semibold text-fg-primary">{useCase.title}</h3>
            <p class="text-sm text-fg-secondary">{useCase.description}</p>
          </div>
        </div>
      ))}
    </div>

    <div class="text-center mt-8">
      <Button href="/catalog" variant="secondary">
        Browse 50+ Job Templates
        <Icon name="arrowRight" size={16} />
      </Button>
    </div>
  </Container>
</section>
```

**Use cases data:**

```javascript
const useCases = [
  {
    emoji: "💬",
    title: "Trial expiration emails",
    description: "Send emails when user trials end",
  },
  {
    emoji: "📊",
    title: "Weekly report generation",
    description: "Generate and email reports every Monday 6 AM",
  },
  {
    emoji: "🗄️",
    title: "Database backups",
    description: "Schedule backups during off-peak hours",
  },
  {
    emoji: "🚀",
    title: "Scheduled deployments",
    description: "Deploy during maintenance windows",
  },
  {
    emoji: "🔔",
    title: "Slack/Discord alerts",
    description: "Send notifications at specific times",
  },
  {
    emoji: "🧹",
    title: "Data cleanup jobs",
    description: "Enforce retention policies automatically",
  },
];
```

---

#### P1.3: Simplify Problem Section Visual

**Current:** 3 text-heavy cards that require reading

**Proposed:** Single visual comparison that communicates instantly

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│     The Old Way              →        The OneOff Way       │
│                                                             │
│  ┌─────────────────┐              ┌─────────────────┐      │
│  │ □ Redis         │              │                 │      │
│  │ □ PostgreSQL    │              │   ./oneoff      │      │
│  │ □ Celery        │              │                 │      │
│  │ □ Docker Compose│              │   ✓ Done.      │      │
│  │ □ 47 YAML files │              │                 │      │
│  │ □ 3 hours setup │              │                 │      │
│  └─────────────────┘              └─────────────────┘      │
│                                                             │
│      ~500MB RAM                      ~20MB RAM             │
│      Multiple services               Single binary         │
│      Complex debugging               Tail one log          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Implementation approach:**

- Keep existing Problem.astro for users who want detail
- Add visual comparison component above it
- Or replace entirely with visual approach

---

#### P1.4: Add FAQ Section

**Purpose:** Preemptively answer objections and reduce support burden

**Location:** Before final CTA section

**Implementation:**

```astro
<!-- New file: landing-page/src/components/landing/FAQ.astro -->

const faqs = [
  {
    q: "Is it really just one binary?",
    a: "Yes. OneOff bundles SQLite, the web UI, and all dependencies into a single executable. No Redis, no Postgres, no Docker required to run."
  },
  {
    q: "Where is my data stored?",
    a: "All data is stored in a local SQLite database file (oneoff.db by default). Your data never leaves your machine."
  },
  {
    q: "Can I use this in production?",
    a: "Absolutely. OneOff is designed for production use with proper error handling, graceful shutdown, and execution logging."
  },
  {
    q: "What about recurring jobs?",
    a: "OneOff is specifically designed for one-time scheduled tasks. For recurring jobs, cron or similar tools are more appropriate."
  },
  {
    q: "How do I upgrade?",
    a: "Download the new binary and replace the old one. Your database and jobs are preserved."
  },
  {
    q: "Is there a hosted/cloud version?",
    a: "Not currently. OneOff is self-hosted only, which means your data stays on your infrastructure."
  }
];
```

---

#### P1.5: Add Pricing/License Clarity

**Purpose:** Answer "what's the catch?" immediately

**Location:** In CTA section or as separate section

**Copy:**

```
┌─────────────────────────────────────────────────────────────┐
│  💰 Free Forever                                            │
│                                                             │
│  OneOff is open source under the Apache 2.0 license.       │
│                                                             │
│  ✓ No usage limits                                         │
│  ✓ No feature gates                                        │
│  ✓ No surprise costs                                       │
│  ✓ Self-hosted on your infrastructure                      │
│                                                             │
│  [View Source Code →]                                      │
└─────────────────────────────────────────────────────────────┘
```

---

### P2: Nice to Have

---

#### P2.1: Homepage Template Preview

**Current:** Users must navigate to `/catalog` to see templates

**Proposed:** Show 3-4 popular templates directly on homepage

```astro
<section class="section">
  <Container>
    <h2>Popular Job Templates</h2>
    <p>Get started instantly with pre-built configurations</p>

    <div class="grid md:grid-cols-3 gap-4">
      <!-- Show top 3 templates -->
      {popularTemplates.map(template => (
        <TemplateCard template={template} />
      ))}
    </div>

    <Button href="/catalog">Browse All Templates →</Button>
  </Container>
</section>
```

---

#### P2.2: Keyboard Navigation Hints

**Purpose:** Signal "built by developers, for developers"

**Implementation:**

```javascript
// Add keyboard shortcuts
document.addEventListener("keydown", (e) => {
  if (e.key === "g")
    window.location.href = "https://github.com/meysam81/oneoff";
  if (e.key === "d") document.getElementById("installation").scrollIntoView();
  if (e.key === "/") window.location.href = "/catalog";
});

// Show hints on hover or in footer
// "Press G for GitHub • D for Download • / for Templates"
```

---

#### P2.3: "What OneOff is NOT" Section

**Purpose:** Set proper expectations, build trust through honesty

```
┌─────────────────────────────────────────────────────────────┐
│  What OneOff is NOT                                         │
│                                                             │
│  ❌ A cron replacement (use cron for recurring tasks)       │
│  ❌ An Airflow competitor (use Airflow for complex DAGs)    │
│  ❌ A distributed queue (use Celery/Redis for that)        │
│                                                             │
│  What OneOff IS                                             │
│                                                             │
│  ✅ The fastest way to schedule one-time tasks             │
│  ✅ A single binary you can run anywhere                   │
│  ✅ The modern `at` command with a beautiful UI            │
└─────────────────────────────────────────────────────────────┘
```

---

#### P2.4: Mobile Sticky CTA

**Purpose:** Always keep conversion action visible on mobile

```css
@media (max-width: 768px) {
  .mobile-sticky-cta {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 1rem;
    background: rgba(10, 10, 11, 0.95);
    backdrop-filter: blur(8px);
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    z-index: 50;
  }
}
```

---

#### P2.5: Real-Time Stats Widget

**If usage tracking is available:**

```astro
<div class="live-stats">
  <span class="pulse"></span>
  <span>2,847 jobs scheduled this week</span>
</div>
```

---

#### P2.6: Comparison Table Improvements

**Current issues:**

- Dense and hard to scan
- All competitors treated equally

**Improvements:**

1. Add row highlighting for OneOff's unique strengths
2. Add "Best for" summary below each column
3. Use color-coded indicators (green checkmarks, red X's)
4. Add hover states with explanatory tooltips

---

## Design Principles

When implementing these changes, follow these principles:

### 1. Reduce Cognitive Load

- Show, don't tell
- One idea per section
- Visual hierarchy guides the eye
- Remove any content that doesn't drive conversion

### 2. Build Trust Progressively

- Social proof early (after hero)
- Technical credibility through real screenshots
- Transparency about limitations
- Clear licensing/pricing

### 3. Match Developer Expectations

- Terminal aesthetics for CLI tool
- Code examples that actually work
- Keyboard shortcuts
- No marketing fluff

### 4. Optimize for Scanning

- Bold headlines that tell the story alone
- Bullet points over paragraphs
- Icons and visual indicators
- Clear section breaks

### 5. Remove All Friction

- OS detection for downloads
- Copy buttons on all code
- Direct links to relevant docs
- No registration required anywhere

---

## Competitive Analysis

### What Top Dev Tools Do Right

#### Linear (linear.app)

- Clean, minimal hero
- Product screenshot immediately visible
- "Get started" CTA is prominent
- Social proof (company logos)
- Video demo available

#### Vercel (vercel.com)

- Interactive deploy demo in hero
- "Deploy Now" is the primary action
- Framework logos build familiarity
- Testimonials from recognizable people

#### Railway (railway.app)

- "Start Now" button prominent
- Live deploy animation
- Template gallery featured
- Pricing clarity upfront

#### Supabase (supabase.com)

- "Start your project" CTA
- Code example in hero
- Feature comparison vs Firebase
- Strong developer community signals

### What OneOff Should Adopt

1. **From Linear:** Clean hero with real product screenshot
2. **From Vercel:** Interactive demo capability
3. **From Railway:** Template gallery prominence
4. **From Supabase:** Direct competitor comparison

---

## Implementation Notes

### File Structure After Changes

```
landing-page/src/
├── components/
│   └── landing/
│       ├── Header.astro
│       ├── Hero.astro          # Updated: New CTAs, OS detection
│       ├── SocialProof.astro   # NEW
│       ├── Problem.astro       # Updated: Visual comparison
│       ├── UseCases.astro      # NEW
│       ├── Features.astro      # Updated: Real screenshots
│       ├── HowItWorks.astro
│       ├── Comparison.astro    # Updated: Better visual hierarchy
│       ├── FAQ.astro           # NEW
│       ├── Installation.astro
│       ├── CTA.astro           # Updated: Pricing clarity
│       └── Footer.astro
├── pages/
│   └── index.astro             # Updated: New section order
└── utils/
    └── github.ts               # Updated: Add download count
```

### New Section Order

```astro
<!-- index.astro -->
<BaseLayout>
  <Header />
  <main>
    <Hero />           <!-- Updated CTAs -->
    <SocialProof />    <!-- NEW: Stars, downloads, testimonial -->
    <Problem />        <!-- Updated: Visual comparison -->
    <UseCases />       <!-- NEW: What developers schedule -->
    <Features />       <!-- Updated: Real screenshots -->
    <HowItWorks />
    <Comparison />     <!-- Updated: Better hierarchy -->
    <FAQ />            <!-- NEW -->
    <Installation />
    <CTA />            <!-- Updated: Pricing clarity -->
  </main>
  <Footer />
</BaseLayout>
```

### GitHub API Updates

```typescript
// landing-page/src/utils/github.ts

export async function getGitHubData() {
  // Existing: version, fullVersion, binarySize
  try {
    // ADD: Download count from releases
    const releases = await fetch(
      "https://api.github.com/repos/meysam81/oneoff/releases",
    ).then((r) => r.json());

    const totalDownloads = releases.reduce((sum, release) => {
      return (
        sum +
        release.assets.reduce((assetSum, asset) => {
          return assetSum + asset.download_count;
        }, 0)
      );
    }, 0);

    // ADD: Star count
    const repo = await fetch(
      "https://api.github.com/repos/meysam81/oneoff",
    ).then((r) => r.json());

    return {
      version,
      fullVersion,
      binarySize,
      stars: repo.stargazers_count, // NEW
      downloads: totalDownloads, // NEW
    };
  } catch (e) {
    // Fallback values if API requests fail
    return {
      version,
      fullVersion,
      binarySize,
      stars: 0,
      downloads: 0,
    };
  }
}
```

---

## Success Metrics

### Primary Metrics

| Metric                  | Current | Target          | Measurement |
| ----------------------- | ------- | --------------- | ----------- |
| Time to first download  | Unknown | < 60 seconds    | Analytics   |
| Bounce rate             | Unknown | < 40%           | Analytics   |
| GitHub star conversions | Unknown | +50%            | GitHub API  |
| Template page visits    | Unknown | 30% of visitors | Analytics   |

### Secondary Metrics

| Metric                 | Target | Measurement    |
| ---------------------- | ------ | -------------- |
| Scroll depth to CTA    | 80%+   | Analytics      |
| Copy button clicks     | Track  | Event tracking |
| Demo interactions      | Track  | Event tracking |
| FAQ section engagement | Track  | Event tracking |

### Qualitative Signals

- User testimonials mentioning ease of setup
- GitHub issues about landing page clarity
- Social media mentions of "simple" or "easy"
- Comparison to competitor landing pages

---

## Changelog

### 2025-11-28 - Initial Roadmap

- Created comprehensive UX/UI audit
- Identified 16+ improvement opportunities
- Prioritized into P0/P1/P2 categories
- Added implementation details for each item
- Defined success metrics

---

## Next Steps

1. **Immediate (This Week)**
   - [ ] Collect testimonials from existing users
   - [ ] Take high-quality screenshots of real UI
   - [ ] Update GitHub data fetching for stars/downloads

2. **Short Term (2 Weeks)**
   - [ ] Implement SocialProof component
   - [ ] Rewrite hero copy (A/B test if possible)
   - [ ] Add OS detection to download CTA
   - [ ] Create product video

3. **Medium Term (1 Month)**
   - [ ] Add UseCases section
   - [ ] Add FAQ section
   - [ ] Replace mock UI with real screenshots
   - [ ] Improve comparison table

4. **Long Term**
   - [ ] Build interactive demo
   - [ ] Add template preview to homepage
   - [ ] Implement analytics tracking
   - [ ] Continuous A/B testing

---

_This roadmap is a living document. Update it as improvements are implemented and new insights are gathered._
