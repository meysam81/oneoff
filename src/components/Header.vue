<template>
  <div class="header">
    <div class="header-left">
      <h2>{{ title }}</h2>
    </div>

    <div class="header-right">
      <n-space>
        <n-badge :value="stats?.queue_depth || 0" :max="99" show-zero>
          <n-button text @click="refreshData">
            <template #icon>
              <n-icon size="20">
                <RefreshOutline />
              </n-icon>
            </template>
          </n-button>
        </n-badge>

        <n-popover trigger="hover">
          <template #trigger>
            <n-button text>
              <n-icon size="20" :color="workerStatusColor">
                <ServerOutline />
              </n-icon>
            </n-button>
          </template>
          <div>
            <p>
              <strong>Workers:</strong>
              {{ workerStatus?.active_workers || 0 }} /
              {{ workerStatus?.total_workers || 0 }}
            </p>
            <p>
              <strong>Queued Jobs:</strong> {{ workerStatus?.queued_jobs || 0 }}
            </p>
          </div>
        </n-popover>
      </n-space>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from "vue";
import { useRoute } from "vue-router";
import { RefreshOutline, ServerOutline } from "@vicons/ionicons5";
import { useSystemStore } from "../stores/system";
import { useJobsStore } from "../stores/jobs";

const route = useRoute();
const systemStore = useSystemStore();
const jobsStore = useJobsStore();

const title = computed(() => route.name || "OneOff");
const stats = computed(() => systemStore.stats);
const workerStatus = computed(() => systemStore.workerStatus);

const workerStatusColor = computed(() => {
  const active = workerStatus.value?.active_workers || 0;
  const total = workerStatus.value?.total_workers || 1;
  const ratio = active / total;
  if (ratio > 0.8) return "#ef4444";
  if (ratio > 0.5) return "#f59e0b";
  return "#10b981";
});

const refreshData = async () => {
  await Promise.all([
    systemStore.fetchStats(),
    systemStore.fetchWorkerStatus(),
    jobsStore.fetchJobs(),
  ]);
};

onMounted(() => {
  const interval = setInterval(refreshData, 5000);
});

onUnmounted(() => {
  clearInterval(interval);
});
</script>

<style scoped>
.header {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}
</style>
