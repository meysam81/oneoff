<template>
  <div class="dashboard">
    <n-space vertical :size="24">
      <!-- Stats Cards -->
      <n-grid :cols="4" :x-gap="16">
        <n-gi>
          <n-card title="Scheduled Jobs">
            <n-statistic :value="stats?.total_scheduled || 0">
              <template #prefix>
                <n-icon size="24" color="#6366f1">
                  <TimeOutline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>

        <n-gi>
          <n-card title="Running Now">
            <n-statistic :value="stats?.currently_running || 0">
              <template #prefix>
                <n-icon size="24" color="#10b981">
                  <PlayCircleOutline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>

        <n-gi>
          <n-card title="Completed Today">
            <n-statistic :value="stats?.completed_today || 0">
              <template #prefix>
                <n-icon size="24" color="#3b82f6">
                  <CheckmarkCircleOutline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>

        <n-gi>
          <n-card title="Failed (24h)">
            <n-statistic :value="stats?.failed_recent || 0">
              <template #prefix>
                <n-icon size="24" color="#ef4444">
                  <CloseCircleOutline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
      </n-grid>

      <!-- Quick Actions -->
      <n-card title="Quick Actions">
        <n-space>
          <n-button type="primary" @click="showCreateJobModal = true">
            <template #icon>
              <n-icon><AddOutline /></n-icon>
            </template>
            Create Job
          </n-button>

          <n-button @click="$router.push('/jobs')">
            <template #icon>
              <n-icon><ListOutline /></n-icon>
            </template>
            View All Jobs
          </n-button>

          <n-button @click="$router.push('/executions')">
            <template #icon>
              <n-icon><BarChartOutline /></n-icon>
            </template>
            Execution History
          </n-button>
        </n-space>
      </n-card>

      <!-- Upcoming Jobs -->
      <n-card title="Upcoming Jobs">
        <JobsTable :jobs="scheduledJobs" :loading="loading" />
      </n-card>

      <!-- Worker Status -->
      <n-card title="Worker Status">
        <n-space vertical>
          <n-progress
            type="line"
            :percentage="workerUsagePercent"
            :status="
              workerUsagePercent > 80
                ? 'error'
                : workerUsagePercent > 50
                  ? 'warning'
                  : 'success'
            "
          >
            Workers: {{ workerStatus?.active_workers || 0 }} /
            {{ workerStatus?.total_workers || 0 }}
          </n-progress>

          <n-text v-if="workerStatus?.running_jobs?.length">
            <strong>Running:</strong>
            {{ workerStatus.running_jobs.length }} jobs
          </n-text>

          <n-text>
            <strong>Queue:</strong> {{ workerStatus?.queued_jobs || 0 }} jobs
            waiting
          </n-text>
        </n-space>
      </n-card>
    </n-space>

    <!-- Create Job Modal -->
    <CreateJobModal v-model:show="showCreateJobModal" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import {
  TimeOutline,
  PlayCircleOutline,
  CheckmarkCircleOutline,
  CloseCircleOutline,
  AddOutline,
  ListOutline,
  BarChartOutline,
} from "@vicons/ionicons5";
import { useSystemStore } from "../stores/system";
import { useJobsStore } from "../stores/jobs";
import JobsTable from "../components/JobsTable.vue";
import CreateJobModal from "../components/CreateJobModal.vue";

const systemStore = useSystemStore();
const jobsStore = useJobsStore();

const showCreateJobModal = ref(false);

const stats = computed(() => systemStore.stats);
const workerStatus = computed(() => systemStore.workerStatus);
const scheduledJobs = computed(() => jobsStore.scheduledJobs.slice(0, 10));
const loading = computed(() => jobsStore.loading);

const workerUsagePercent = computed(() => {
  const active = workerStatus.value?.active_workers || 0;
  const total = workerStatus.value?.total_workers || 1;
  return Math.round((active / total) * 100);
});

onMounted(async () => {
  await systemStore.initializeApp();
  await jobsStore.fetchJobs({ limit: 20, status: "scheduled" });
});
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
}
</style>
