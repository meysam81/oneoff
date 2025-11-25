<template>
  <div v-if="job" class="job-details">
    <n-space vertical :size="16">
      <n-card>
        <template #header>
          <n-space justify="space-between" align="center">
            <n-space align="center">
              <n-button text @click="$router.back()">
                <n-icon size="20"><ArrowBackOutline /></n-icon>
              </n-button>
              <h2>{{ job.name }}</h2>
              <n-tag :type="statusColors[job.status]">{{ job.status }}</n-tag>
            </n-space>

            <n-space>
              <n-button
                v-if="job.status === 'scheduled'"
                type="primary"
                @click="executeNow"
              >
                Execute Now
              </n-button>

              <n-button
                v-if="job.status === 'scheduled' || job.status === 'running'"
                type="error"
                @click="cancelJob"
              >
                Cancel
              </n-button>

              <n-button @click="cloneJob">Clone</n-button>

              <n-button
                v-if="job.status !== 'running'"
                type="error"
                @click="deleteJob"
              >
                Delete
              </n-button>
            </n-space>
          </n-space>
        </template>

        <n-descriptions :column="2" bordered>
          <n-descriptions-item label="Type">
            <n-tag size="small">{{ job.type }}</n-tag>
          </n-descriptions-item>

          <n-descriptions-item label="Priority">
            {{ job.priority }}
          </n-descriptions-item>

          <n-descriptions-item label="Scheduled At">
            <n-time :time="new Date(job.scheduled_at)" />
          </n-descriptions-item>

          <n-descriptions-item label="Created At">
            <n-time :time="new Date(job.created_at)" />
          </n-descriptions-item>

          <n-descriptions-item label="Project">
            {{ getProjectName(job.project_id) }}
          </n-descriptions-item>

          <n-descriptions-item label="Timezone">
            {{ job.timezone }}
          </n-descriptions-item>

          <n-descriptions-item label="Tags" :span="2">
            <n-space>
              <n-tag
                v-for="tag in job.tags"
                :key="tag.id"
                size="small"
                :color="{ color: tag.color }"
              >
                {{ tag.name }}
              </n-tag>
            </n-space>
          </n-descriptions-item>

          <n-descriptions-item label="Configuration" :span="2">
            <n-code :code="formatConfig(job.config)" language="json" />
          </n-descriptions-item>
        </n-descriptions>
      </n-card>

      <n-card title="Execution History">
        <ExecutionsList :job-id="job.id" />
      </n-card>
    </n-space>
  </div>

  <n-spin v-else :show="loading" />

  <CloneJobModal
    v-model:show="showCloneModal"
    :job="job"
    @cloned="handleJobCloned"
  />
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useMessage, useDialog } from "naive-ui";
import { ArrowBackOutline } from "@vicons/ionicons5";
import { useJobsStore } from "../stores/jobs";
import { useSystemStore } from "../stores/system";
import ExecutionsList from "../components/ExecutionsList.vue";
import CloneJobModal from "../components/CloneJobModal.vue";

const route = useRoute();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();

const jobsStore = useJobsStore();
const systemStore = useSystemStore();

const job = computed(() => jobsStore.currentJob);
const loading = computed(() => jobsStore.loading);

const statusColors = {
  scheduled: "info",
  running: "warning",
  completed: "success",
  failed: "error",
  cancelled: "default",
};

const getProjectName = (id) => {
  return systemStore.projects.find((p) => p.id === id)?.name || id;
};

const formatConfig = (config) => {
  try {
    return JSON.stringify(JSON.parse(config), null, 2);
  } catch {
    return config;
  }
};

const executeNow = async () => {
  try {
    await jobsStore.executeJob(job.value.id);
    message.success("Job scheduled for immediate execution");
    await jobsStore.fetchJob(job.value.id);
  } catch (error) {
    message.error(error.message);
  }
};

const cancelJob = async () => {
  dialog.warning({
    title: "Cancel Job",
    content: "Are you sure you want to cancel this job?",
    positiveText: "Yes",
    negativeText: "No",
    onPositiveClick: async () => {
      try {
        await jobsStore.cancelJob(job.value.id);
        message.success("Job cancelled");
        await jobsStore.fetchJob(job.value.id);
      } catch (error) {
        message.error(error.message);
      }
    },
  });
};

const showCloneModal = ref(false);

const cloneJob = () => {
  showCloneModal.value = true;
};

const handleJobCloned = async () => {
  // Refresh job list
  await jobsStore.fetchJobs({}, false);

};

const deleteJob = async () => {
  dialog.error({
    title: "Delete Job",
    content:
      "Are you sure you want to delete this job? This action cannot be undone.",
    positiveText: "Delete",
    negativeText: "Cancel",
    onPositiveClick: async () => {
      try {
        await jobsStore.deleteJob(job.value.id);
        message.success("Job deleted");
        router.push("/jobs");
      } catch (error) {
        message.error(error.message);
      }
    },
  });
};

onMounted(async () => {
  await systemStore.initializeApp();
  await jobsStore.fetchJob(route.params.id);
});
</script>
