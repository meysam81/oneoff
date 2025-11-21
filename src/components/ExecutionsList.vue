<template>
  <n-data-table
    :columns="columns"
    :data="executions"
    :loading="loading"
    :row-key="(row) => row.id"
  >
    <template #expand="{ row }">
      <div style="padding: 16px">
        <n-space vertical :size="12">
          <!-- Error Section -->
          <n-alert
            v-if="row.error"
            title="Error"
            type="error"
            :bordered="false"
          >
            <n-code :code="row.error" language="text" />
          </n-alert>

          <!-- Output Section -->
          <n-card
            v-if="row.output"
            title="Output"
            size="small"
            :segmented="{ content: true }"
          >
            <template #header-extra>
              <n-button
                size="small"
                @click="downloadLogs(row, 'output')"
              >
                Download
              </n-button>
            </template>
            <n-scrollbar style="max-height: 400px">
              <n-code
                :code="row.output"
                language="bash"
                word-wrap
              />
            </n-scrollbar>
          </n-card>

          <!-- Metadata Section -->
          <n-descriptions :column="3" size="small" bordered>
            <n-descriptions-item label="Exit Code">
              <n-tag
                v-if="row.exit_code !== null"
                :type="row.exit_code === 0 ? 'success' : 'error'"
                size="small"
              >
                {{ row.exit_code }}
              </n-tag>
              <span v-else>-</span>
            </n-descriptions-item>
            <n-descriptions-item label="Duration">
              {{ row.duration_ms ? `${(row.duration_ms / 1000).toFixed(2)}s` : "-" }}
            </n-descriptions-item>
            <n-descriptions-item label="Completed At">
              <n-time
                v-if="row.completed_at"
                :time="new Date(row.completed_at)"
              />
              <span v-else>-</span>
            </n-descriptions-item>
          </n-descriptions>
        </n-space>
      </div>
    </template>
  </n-data-table>
</template>

<script setup>
import { ref, h, onMounted } from "vue";
import { NButton, NTag, NTime, NSpace, NAlert, NCard, NCode, NScrollbar, NDescriptions, NDescriptionsItem } from "naive-ui";
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
    type: "expand",
    expandable: (row) => !!(row.output || row.error),
  },
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
    render: (row) => {
      if (row.exit_code === null) return "-";
      return h(
        NTag,
        {
          type: row.exit_code === 0 ? "success" : "error",
          size: "small",
        },
        { default: () => row.exit_code }
      );
    },
  },
];

const downloadLogs = (execution, type) => {
  const content = type === "output" ? execution.output : execution.error;
  const blob = new Blob([content], { type: "text/plain" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `execution-${execution.id}-${type}.log`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};

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
