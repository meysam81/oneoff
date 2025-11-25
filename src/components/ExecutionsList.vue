<template>
  <n-data-table
    :columns="columns"
    :data="executions"
    :loading="loading"
    :row-key="rowKey"
  />
</template>

<script setup>
import { ref, h, onMounted } from "vue";
import {
  NButton,
  NTag,
  NTime,
  NSpace,
  NAlert,
  NCard,
  NCode,
  NScrollbar,
  NDescriptions,
  NDescriptionsItem,
} from "naive-ui";
import { executionsAPI } from "../utils/api";

var props = defineProps({
  jobId: String,
});

var executions = ref([]);
var loading = ref(false);

var statusColors = {
  running: "warning",
  completed: "success",
  failed: "error",
  cancelled: "default",
};

function rowKey(row) {
  return row.id;
}

function downloadLogs(execution, type) {
  var content = type === "output" ? execution.output : execution.error;
  var blob = new Blob([content], { type: "text/plain" });
  var url = URL.createObjectURL(blob);
  var a = document.createElement("a");
  a.href = url;
  a.download = "execution-" + execution.id + "-" + type + ".log";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

function renderExpand(row) {
  var children = [];

  if (row.error) {
    children.push(
      h(
        NAlert,
        { title: "Error", type: "error", bordered: false },
        {
          default: function () {
            return h(NCode, { code: row.error, language: "text" });
          },
        },
      ),
    );
  }

  if (row.output) {
    children.push(
      h(
        NCard,
        { title: "Output", size: "small", segmented: { content: true } },
        {
          "header-extra": function () {
            return h(
              NButton,
              {
                size: "small",
                onClick: function () {
                  downloadLogs(row, "output");
                },
              },
              {
                default: function () {
                  return "Download";
                },
              },
            );
          },
          default: function () {
            return h(
              NScrollbar,
              { style: "max-height: 400px" },
              {
                default: function () {
                  return h(NCode, {
                    code: row.output,
                    language: "bash",
                    wordWrap: true,
                  });
                },
              },
            );
          },
        },
      ),
    );
  }

  var descriptionItems = [
    h(
      NDescriptionsItem,
      { label: "Exit Code" },
      {
        default: function () {
          if (row.exit_code !== null) {
            return h(
              NTag,
              {
                type: row.exit_code === 0 ? "success" : "error",
                size: "small",
              },
              {
                default: function () {
                  return row.exit_code;
                },
              },
            );
          }
          return "-";
        },
      },
    ),
    h(
      NDescriptionsItem,
      { label: "Duration" },
      {
        default: function () {
          return row.duration_ms
            ? (row.duration_ms / 1000).toFixed(2) + "s"
            : "-";
        },
      },
    ),
    h(
      NDescriptionsItem,
      { label: "Completed At" },
      {
        default: function () {
          if (row.completed_at) {
            return h(NTime, { time: new Date(row.completed_at) });
          }
          return "-";
        },
      },
    ),
  ];

  children.push(
    h(
      NDescriptions,
      { column: 3, size: "small", bordered: true },
      {
        default: function () {
          return descriptionItems;
        },
      },
    ),
  );

  return h("div", { style: "padding: 16px" }, [
    h(
      NSpace,
      { vertical: true, size: 12 },
      {
        default: function () {
          return children;
        },
      },
    ),
  ]);
}

var columns = [
  {
    type: "expand",
    expandable: function (row) {
      return !!(row.output || row.error);
    },
    renderExpand: renderExpand,
  },
  {
    title: "Started",
    key: "started_at",
    render: function (row) {
      return h(NTime, { time: new Date(row.started_at) });
    },
  },
  {
    title: "Status",
    key: "status",
    render: function (row) {
      return h(
        NTag,
        { type: statusColors[row.status], size: "small" },
        {
          default: function () {
            return row.status;
          },
        },
      );
    },
  },
  {
    title: "Duration",
    key: "duration_ms",
    render: function (row) {
      return row.duration_ms ? (row.duration_ms / 1000).toFixed(2) + "s" : "-";
    },
  },
  {
    title: "Exit Code",
    key: "exit_code",
    render: function (row) {
      if (row.exit_code === null) return "-";
      return h(
        NTag,
        { type: row.exit_code === 0 ? "success" : "error", size: "small" },
        {
          default: function () {
            return row.exit_code;
          },
        },
      );
    },
  },
];

async function fetchExecutions() {
  loading.value = true;
  try {
    var response = await executionsAPI.list({ job_id: props.jobId });
    executions.value = response.data || [];
  } catch (error) {
    console.error("Failed to fetch executions:", error);
  } finally {
    loading.value = false;
  }
}

onMounted(fetchExecutions);
</script>
