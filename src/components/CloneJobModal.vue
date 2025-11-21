<template>
  <n-modal
    v-model:show="visible"
    preset="card"
    title="Clone Job"
    style="width: 500px"
    :mask-closable="false"
  >
    <n-form ref="formRef" :model="formValue" :rules="rules">
      <n-form-item label="New Job Name" path="name">
        <n-input v-model:value="formValue.name" placeholder="Enter new job name" />
      </n-form-item>

      <n-form-item label="Scheduled Time" path="scheduled_at">
        <n-date-picker
          v-model:value="scheduledTimestamp"
          type="datetime"
          clearable
          style="width: 100%"
          :is-date-disabled="isDateDisabled"
        />
      </n-form-item>

      <n-alert type="info" :bordered="false" style="margin-bottom: 12px">
        The cloned job will preserve the original job's configuration, project, and tags.
      </n-alert>
    </n-form>

    <template #footer>
      <n-space justify="space-between">
        <n-button @click="visible = false">Cancel</n-button>
        <n-button type="primary" :loading="loading" @click="handleClone">
          Clone Job
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useMessage } from "naive-ui";
import { useJobsStore } from "../stores/jobs";

const props = defineProps({
  show: Boolean,
  job: Object,
});

const emit = defineEmits(["update:show", "cloned"]);

const message = useMessage();
const jobsStore = useJobsStore();

const visible = computed({
  get: () => props.show,
  set: (val) => emit("update:show", val),
});

const formRef = ref(null);
const loading = ref(false);
const scheduledTimestamp = ref(Date.now() + 3600000); // Default 1 hour from now

const formValue = ref({
  name: "",
});

const rules = {
  name: { required: true, message: "Job name is required" },
};

const isDateDisabled = (ts) => {
  return ts < Date.now();
};

// Watch for job prop changes to update the form
watch(
  () => props.job,
  (newJob) => {
    if (newJob) {
      formValue.value.name = `${newJob.name} (Copy)`;
    }
  },
  { immediate: true }
);

const handleClone = async () => {
  try {
    await formRef.value?.validate();

    if (!props.job) {
      message.error("No job selected for cloning");
      return;
    }

    const scheduledAt = new Date(scheduledTimestamp.value).toISOString();

    loading.value = true;

    // Clone the job
    await jobsStore.createJob({
      name: formValue.value.name,
      type: props.job.type,
      config: props.job.config,
      scheduled_at: scheduledAt,
      priority: props.job.priority,
      project_id: props.job.project_id,
      timezone: props.job.timezone,
      tag_ids: props.job.tags?.map(t => t.id) || [],
    });

    message.success("Job cloned successfully");
    visible.value = false;
    emit("cloned");
    resetForm();
  } catch (error) {
    message.error(error.message || "Failed to clone job");
  } finally {
    loading.value = false;
  }
};

const resetForm = () => {
  formValue.value = {
    name: props.job ? `${props.job.name} (Copy)` : "",
  };
  scheduledTimestamp.value = Date.now() + 3600000;
};
</script>
