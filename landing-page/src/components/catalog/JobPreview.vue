<script setup lang="ts">
import { ref } from "vue";

interface Template {
  id: string;
  name: string;
  description: string;
  category: string;
  tags: string[];
  author: {
    name: string;
    github: string;
  };
  job: {
    type: string;
    config: Record<string, unknown>;
  };
  created_at: string;
}

const props = defineProps<{
  template: Template;
}>();

const copied = ref(false);

const configJson = JSON.stringify(props.template.job, null, 2);

async function copyConfig() {
  try {
    await navigator.clipboard.writeText(configJson);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (err) {
    console.error("Failed to copy:", err);
  }
}

const categoryLabels: Record<string, string> = {
  backup: "Backup & Recovery",
  monitoring: "Monitoring & Alerts",
  cicd: "CI/CD Integration",
  database: "Database Maintenance",
  api: "API & Webhooks",
  devops: "DevOps Automation",
  reporting: "Reporting",
  misc: "Miscellaneous",
};

const jobTypeColors: Record<string, string> = {
  http: "text-accent-primary bg-accent-primary/10 border-accent-primary/20",
  shell: "text-success bg-success/10 border-success/20",
  docker: "text-warning bg-warning/10 border-warning/20",
};

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}
</script>

<template>
  <div class="job-preview grid lg:grid-cols-2 gap-8">
    <!-- Left: Code Preview -->
    <div class="order-2 lg:order-1">
      <div
        class="bg-bg-secondary rounded-xl border border-border-subtle overflow-hidden"
      >
        <div
          class="flex items-center justify-between px-4 py-3 bg-bg-tertiary border-b border-border-subtle"
        >
          <div class="flex items-center gap-2">
            <span class="w-3 h-3 rounded-full bg-error/80"></span>
            <span class="w-3 h-3 rounded-full bg-warning/80"></span>
            <span class="w-3 h-3 rounded-full bg-success/80"></span>
            <span class="text-xs text-fg-muted ml-2">job-config.json</span>
          </div>
          <button
            @click="copyConfig"
            class="flex items-center gap-2 px-3 py-1.5 text-xs font-medium rounded-lg transition-all duration-200"
            :class="
              copied
                ? 'bg-success/10 text-success'
                : 'bg-bg-secondary text-fg-secondary hover:text-fg-primary'
            "
            aria-label="Copy configuration to clipboard"
          >
            <svg
              v-if="!copied"
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
              />
            </svg>
            <svg
              v-else
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 13l4 4L19 7"
              />
            </svg>
            {{ copied ? "Copied!" : "Copy config" }}
          </button>
        </div>
        <pre
          class="p-4 overflow-x-auto text-sm font-mono text-fg-secondary"
        ><code>{{ configJson }}</code></pre>
      </div>

      <!-- Import instructions -->
      <div
        class="mt-6 bg-bg-secondary rounded-xl border border-border-subtle p-4"
      >
        <h4 class="text-sm font-semibold text-fg-primary mb-3">
          Import via API
        </h4>
        <div
          class="bg-bg-tertiary rounded-lg p-3 font-mono text-xs text-fg-secondary overflow-x-auto"
        >
          <code
            >curl -X POST http://localhost:8080/api/jobs -H "Content-Type:
            application/json" -d '{{ JSON.stringify(template.job) }}'</code
          >
        </div>
        <p class="text-xs text-fg-muted mt-3">
          Or copy the config above and create the job manually via the web UI.
        </p>
      </div>
    </div>

    <!-- Right: Metadata -->
    <div class="order-1 lg:order-2">
      <div class="space-y-6">
        <!-- Title -->
        <div>
          <h1 class="text-2xl font-display font-bold text-fg-primary mb-2">
            {{ template.name }}
          </h1>
          <p class="text-fg-secondary leading-relaxed">
            {{ template.description }}
          </p>
        </div>

        <!-- Metadata cards -->
        <div class="grid grid-cols-2 gap-4">
          <div
            class="bg-bg-secondary rounded-lg border border-border-subtle p-4"
          >
            <span class="text-xs text-fg-muted block mb-1">Job Type</span>
            <span
              class="inline-flex items-center px-2 py-1 text-sm font-medium rounded-full border"
              :class="
                jobTypeColors[template.job.type] ||
                'text-fg-secondary bg-bg-tertiary border-border-subtle'
              "
            >
              {{ template.job.type.toUpperCase() }}
            </span>
          </div>
          <div
            class="bg-bg-secondary rounded-lg border border-border-subtle p-4"
          >
            <span class="text-xs text-fg-muted block mb-1">Category</span>
            <span class="text-sm text-fg-primary font-medium">
              {{ categoryLabels[template.category] || template.category }}
            </span>
          </div>
        </div>

        <!-- Tags -->
        <div class="bg-bg-secondary rounded-lg border border-border-subtle p-4">
          <span class="text-xs text-fg-muted block mb-2">Tags</span>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="tag in template.tags"
              :key="tag"
              class="px-2 py-1 text-xs bg-bg-tertiary text-fg-secondary rounded-full"
            >
              #{{ tag }}
            </span>
          </div>
        </div>

        <!-- Author -->
        <div class="bg-bg-secondary rounded-lg border border-border-subtle p-4">
          <span class="text-xs text-fg-muted block mb-2">Author</span>
          <a
            :href="`https://github.com/${template.author.github}`"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-3 hover:text-accent-primary transition-colors"
          >
            <img
              :src="`https://github.com/${template.author.github}.png?size=40`"
              :alt="template.author.name"
              class="w-10 h-10 rounded-full bg-bg-tertiary"
              loading="lazy"
            />
            <div>
              <span class="text-sm font-medium text-fg-primary block">{{
                template.author.name
              }}</span>
              <span class="text-xs text-fg-muted"
                >@{{ template.author.github }}</span
              >
            </div>
          </a>
        </div>

        <!-- Created date -->
        <div class="text-sm text-fg-muted">
          Added on {{ formatDate(template.created_at) }}
        </div>
      </div>
    </div>
  </div>
</template>
