<template>
  <n-data-table
    :columns="columns"
    :data="executions"
    :loading="loading"
    :row-key="(row) => row.id"
  />
</template>

<script setup>
import { ref, h, onMounted } from "vue";
import { NButton, NTag, NTime, NSpace } from "naive-ui";
import { executionsAPI } from "../utils/api";

const props = defineProps({
  jobId: String,
});

const executions = ref([]);
const loading = ref(false);

const statusColors = {
  running: "warning",
  completed: "success",
  failed: "error",
  cancelled: "default",
};

const columns = [
  {
    title: "Started",
    key: "started_at",
    render: (row) => h(NTime, { time: new Date(row.started_at) }),
  },
  {
    title: "Status",
    key: "status",
    render: (row) =>
      h(
        NTag,
        { type: statusColors[row.status], size: "small" },
        { default: () => row.status },
      ),
  },
  {
    title: "Duration",
    key: "duration_ms",
    render: (row) =>
      row.duration_ms ? `${(row.duration_ms / 1000).toFixed(2)}s` : "-",
  },
  {
    title: "Exit Code",
    key: "exit_code",
    render: (row) => (row.exit_code !== null ? row.exit_code : "-"),
  },
];

const fetchExecutions = async () => {
  loading.value = true;
  try {
    const response = await executionsAPI.list({ job_id: props.jobId });
    executions.value = response.data || [];
  } catch (error) {
    console.error("Failed to fetch executions:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchExecutions);
</script>
