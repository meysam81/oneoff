<template>
  <div class="jobs-view">
    <n-space vertical :size="16">
      <n-card>
        <n-space justify="space-between">
          <n-space>
            <n-button type="primary" @click="showCreateModal = true">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              Create Job
            </n-button>

            <n-button @click="fetchJobs">
              <template #icon>
                <n-icon><RefreshOutline /></n-icon>
              </template>
              Refresh
            </n-button>
          </n-space>

          <n-space>
            <n-input
              v-model:value="searchQuery"
              placeholder="Search jobs..."
              clearable
              @update:value="handleSearch"
            >
              <template #prefix>
                <n-icon><SearchOutline /></n-icon>
              </template>
            </n-input>

            <n-select
              v-model:value="statusFilter"
              :options="statusOptions"
              placeholder="Filter by status"
              clearable
              style="width: 180px"
              @update:value="handleFilterChange"
            />

            <n-select
              v-model:value="projectFilter"
              :options="projectOptions"
              placeholder="Filter by project"
              clearable
              style="width: 200px"
              @update:value="handleFilterChange"
            />
          </n-space>
        </n-space>
      </n-card>

      <n-card>
        <JobsTable :jobs="jobs" :loading="loading" />
      </n-card>
    </n-space>

    <CreateJobModal v-model:show="showCreateModal" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { AddOutline, RefreshOutline, SearchOutline } from "@vicons/ionicons5";
import { useJobsStore } from "../stores/jobs";
import { useSystemStore } from "../stores/system";
import { debounce } from "../utils/debounce";
import JobsTable from "../components/JobsTable.vue";
import CreateJobModal from "../components/CreateJobModal.vue";

const jobsStore = useJobsStore();
const systemStore = useSystemStore();

const showCreateModal = ref(false);
const searchQuery = ref("");
const statusFilter = ref(null);
const projectFilter = ref(null);

const jobs = computed(() => jobsStore.jobs);
const loading = computed(() => jobsStore.loading);

const statusOptions = [
  { label: "Scheduled", value: "scheduled" },
  { label: "Running", value: "running" },
  { label: "Completed", value: "completed" },
  { label: "Failed", value: "failed" },
  { label: "Cancelled", value: "cancelled" },
];

const projectOptions = computed(() =>
  systemStore.projects.map((p) => ({ label: p.name, value: p.id })),
);

const fetchJobs = async () => {
  await jobsStore.fetchJobs({
    search: searchQuery.value,
    status: statusFilter.value,
    project_id: projectFilter.value,
  });
};

// Debounce search to avoid excessive API calls
const handleSearch = debounce(() => {
  fetchJobs();
}, 400);

const handleFilterChange = () => {
  fetchJobs();
};

onMounted(async () => {
  await fetchJobs();
});
</script>
