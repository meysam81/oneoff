<script setup lang="ts">
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

defineProps<{
  templates: Template[];
}>();

const jobTypeIcons: Record<string, string> = {
  http: "M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9",
  shell:
    "M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z",
  docker:
    "M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4",
};

const jobTypeColors: Record<string, string> = {
  http: "text-accent-primary",
  shell: "text-success",
  docker: "text-warning",
};

const categoryLabels: Record<string, string> = {
  backup: "Backup",
  monitoring: "Monitoring",
  cicd: "CI/CD",
  database: "Database",
  api: "API",
  devops: "DevOps",
  reporting: "Reporting",
  misc: "Misc",
};

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
</script>

<template>
  <div class="catalog-grid">
    <div v-if="templates.length === 0" class="text-center py-12">
      <svg
        class="w-12 h-12 mx-auto text-fg-muted mb-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <p class="text-fg-secondary">
        No templates found matching your criteria.
      </p>
      <p class="text-fg-muted text-sm mt-1">
        Try adjusting your search or filters.
      </p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <a
        v-for="template in templates"
        :key="template.id"
        :href="`/catalog/${template.id}`"
        class="group bg-bg-secondary rounded-xl border border-border-subtle p-6 transition-all duration-300 hover:border-accent-primary/30 hover:shadow-[0_0_30px_rgba(34,211,238,0.05)]"
      >
        <!-- Header -->
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-bg-tertiary rounded-lg">
              <svg
                class="w-5 h-5"
                :class="jobTypeColors[template.job.type] || 'text-fg-muted'"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  :d="jobTypeIcons[template.job.type] || jobTypeIcons.shell"
                />
              </svg>
            </div>
            <span
              class="px-2 py-1 text-xs font-medium rounded-full bg-bg-tertiary text-fg-secondary border border-border-subtle"
            >
              {{ categoryLabels[template.category] || template.category }}
            </span>
          </div>

          <svg
            class="w-5 h-5 text-fg-muted group-hover:text-accent-primary transition-colors"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1.5"
              d="M14 5l7 7m0 0l-7 7m7-7H3"
            />
          </svg>
        </div>

        <!-- Title & Description -->
        <h3
          class="text-lg font-display font-semibold text-fg-primary mb-2 group-hover:text-accent-primary transition-colors"
        >
          {{ template.name }}
        </h3>
        <p class="text-fg-secondary text-sm leading-relaxed mb-4 line-clamp-2">
          {{ template.description }}
        </p>

        <!-- Tags -->
        <div class="flex flex-wrap gap-2 mb-4">
          <span
            v-for="tag in template.tags.slice(0, 3)"
            :key="tag"
            class="px-2 py-0.5 text-xs bg-bg-tertiary text-fg-muted rounded"
          >
            #{{ tag }}
          </span>
          <span
            v-if="template.tags.length > 3"
            class="px-2 py-0.5 text-xs text-fg-muted"
          >
            +{{ template.tags.length - 3 }} more
          </span>
        </div>

        <!-- Footer -->
        <div
          class="flex items-center justify-between text-xs text-fg-muted pt-4 border-t border-border-subtle"
        >
          <a
            :href="`https://github.com/${template.author.github}`"
            target="_blank"
            rel="noopener noreferrer"
            class="hover:text-fg-secondary transition-colors"
            @click.stop
          >
            @{{ template.author.github }}
          </a>
          <span>{{ formatDate(template.created_at) }}</span>
        </div>
      </a>
    </div>
  </div>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
