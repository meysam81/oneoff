<template>
  <n-data-table
    :columns="columns"
    :data="jobs"
    :loading="loading"
    :pagination="pagination"
    :row-key="(row) => row.id"
    striped
  />
</template>

<script setup>
import { h } from "vue";
import { NButton, NSpace, NTag, NTime } from "naive-ui";
import { useRouter } from "vue-router";

const props = defineProps({
  jobs: {
    type: Array,
    default: () => [],
  },
  loading: Boolean,
});

const router = useRouter();

// Moved outside to prevent recreation on every render
const statusColors = {
  scheduled: "info",
  running: "warning",
  completed: "success",
  failed: "error",
  cancelled: "default",
};

// Static column configuration - created once
const createColumns = (router) => [
  {
    title: "Name",
    key: "name",
    ellipsis: { tooltip: true },
  },
  {
    title: "Type",
    key: "type",
    render: (row) => h(NTag, { size: "small" }, { default: () => row.type }),
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
    title: "Scheduled",
    key: "scheduled_at",
    render: (row) =>
      h(NTime, { time: new Date(row.scheduled_at), type: "relative" }),
  },
  {
    title: "Priority",
    key: "priority",
  },
  {
    title: "Actions",
    key: "actions",
    width: 100,
    render: (row) =>
      h(
        NButton,
        {
          size: "small",
          onClick: () => router.push(`/jobs/${row.id}`),
        },
        { default: () => "View" },
      ),
  },
];

const columns = createColumns(router);

const pagination = {
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
};
</script>
