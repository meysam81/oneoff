<script setup lang="ts">
import { ref, computed, watch } from "vue";

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
  templates: Template[];
}>();

const emit = defineEmits<{
  (e: "filter", templates: Template[]): void;
}>();

const search = ref("");
const selectedCategories = ref<string[]>([]);
const sortBy = ref<"newest" | "alpha">("newest");

const categories = [
  { value: "backup", label: "Backup & Recovery" },
  { value: "monitoring", label: "Monitoring & Alerts" },
  { value: "cicd", label: "CI/CD Integration" },
  { value: "database", label: "Database Maintenance" },
  { value: "api", label: "API & Webhooks" },
  { value: "devops", label: "DevOps Automation" },
  { value: "reporting", label: "Reporting" },
  { value: "misc", label: "Miscellaneous" },
];

const filtered = computed(() => {
  let result = [...props.templates];

  // Search filter
  if (search.value) {
    const q = search.value.toLowerCase();
    result = result.filter(
      (t) =>
        t.name.toLowerCase().includes(q) ||
        t.description.toLowerCase().includes(q) ||
        t.tags.some((tag) => tag.toLowerCase().includes(q)),
    );
  }

  // Category filter
  if (selectedCategories.value.length > 0) {
    result = result.filter((t) =>
      selectedCategories.value.includes(t.category),
    );
  }

  // Sort
  if (sortBy.value === "alpha") {
    result.sort((a, b) => a.name.localeCompare(b.name));
  } else {
    result.sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    );
  }

  return result;
});

// Emit filtered results on change
watch(filtered, (val) => emit("filter", val), { immediate: true });

function toggleCategory(category: string) {
  const index = selectedCategories.value.indexOf(category);
  if (index === -1) {
    selectedCategories.value.push(category);
  } else {
    selectedCategories.value.splice(index, 1);
  }
}

function clearFilters() {
  search.value = "";
  selectedCategories.value = [];
  sortBy.value = "newest";
}
</script>

<template>
  <div class="catalog-search space-y-3 sm:space-y-4">
    <!-- Search input -->
    <div class="relative">
      <svg
        class="absolute left-2.5 sm:left-3 top-1/2 -translate-y-1/2 w-4 h-4 sm:w-5 sm:h-5 text-fg-muted"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
        />
      </svg>
      <input
        v-model="search"
        type="search"
        placeholder="Search templates..."
        class="w-full pl-8 sm:pl-10 pr-3 sm:pr-4 py-2.5 sm:py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-sm sm:text-base text-fg-primary placeholder:text-fg-muted focus:outline-none focus:border-accent-primary/50 focus:ring-1 focus:ring-accent-primary/50 transition-colors"
      />
    </div>

    <!-- Filters row -->
    <div
      class="flex flex-col gap-3 sm:gap-4"
    >
      <!-- Category filters -->
      <div class="overflow-x-auto -mx-2 px-2 pb-2 sm:mx-0 sm:px-0 sm:pb-0">
        <div class="flex items-center gap-1.5 sm:gap-2 sm:flex-wrap min-w-max sm:min-w-0">
          <span class="text-xs sm:text-sm text-fg-muted whitespace-nowrap">Filter:</span>
          <button
            v-for="cat in categories"
            :key="cat.value"
            @click="toggleCategory(cat.value)"
            :class="[
              'px-2 sm:px-3 py-1 sm:py-1.5 text-[10px] sm:text-xs font-medium rounded-full transition-all duration-200 whitespace-nowrap',
              selectedCategories.includes(cat.value)
                ? 'bg-accent-primary/20 text-accent-primary border border-accent-primary/30'
                : 'bg-bg-tertiary text-fg-secondary border border-border-subtle hover:border-border-default',
            ]"
            role="button"
            :aria-pressed="selectedCategories.includes(cat.value)"
          >
            {{ cat.label }}
          </button>
        </div>
      </div>

      <!-- Sort & Clear -->
      <div class="flex items-center justify-between sm:justify-end gap-3 sm:gap-4">
        <select
          v-model="sortBy"
          class="px-2 sm:px-3 py-1 sm:py-1.5 bg-bg-tertiary border border-border-subtle rounded-lg text-xs sm:text-sm text-fg-secondary focus:outline-none focus:border-accent-primary/50"
          aria-label="Sort by"
        >
          <option value="newest">Sort by: Newest first</option>
          <option value="alpha">Sort by: A-Z</option>
        </select>

        <button
          v-if="search || selectedCategories.length > 0"
          @click="clearFilters"
          class="text-xs sm:text-sm text-fg-muted hover:text-fg-secondary transition-colors"
        >
          Clear filters
        </button>
      </div>
    </div>

    <!-- Results count -->
    <div class="text-xs sm:text-sm text-fg-muted">
      Showing {{ filtered.length }} of {{ templates.length }} templates
    </div>
  </div>
</template>

<style scoped>
/* Custom styles for this component */
input[type="search"]::-webkit-search-cancel-button {
  -webkit-appearance: none;
  appearance: none;
}
</style>
