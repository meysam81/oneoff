<template>
  <n-modal
    v-model:show="visible"
    preset="card"
    title="Create New Job"
    style="width: 700px"
    :mask-closable="false"
  >
    <n-form ref="formRef" :model="formValue" :rules="rules">
      <n-form-item label="Job Name" path="name">
        <n-input v-model:value="formValue.name" placeholder="Enter job name" />
      </n-form-item>

      <n-form-item label="Job Type" path="type">
        <n-select
          v-model:value="formValue.type"
          :options="jobTypeOptions"
          placeholder="Select job type"
        />
      </n-form-item>

      <n-form-item label="Project" path="project_id">
        <n-select
          v-model:value="formValue.project_id"
          :options="projectOptions"
          placeholder="Select project"
        />
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

      <n-form-item label="Priority" path="priority">
        <n-slider
          v-model:value="formValue.priority"
          :min="1"
          :max="10"
          :step="1"
          :marks="priorityMarks"
        />
      </n-form-item>

      <n-form-item label="Tags" path="tag_ids">
        <n-select
          v-model:value="formValue.tag_ids"
          multiple
          :options="tagOptions"
          placeholder="Select tags"
        />
      </n-form-item>

      <!-- Job type specific config -->
      <n-collapse v-if="formValue.type">
        <n-collapse-item title="Job Configuration" name="config">
          <HTTPConfig v-if="formValue.type === 'http'" v-model="jobConfig" />
          <ShellConfig
            v-else-if="formValue.type === 'shell'"
            v-model="jobConfig"
          />
          <DockerConfig
            v-else-if="formValue.type === 'docker'"
            v-model="jobConfig"
          />
        </n-collapse-item>
      </n-collapse>
    </n-form>

    <template #footer>
      <n-space justify="space-between">
        <n-button @click="visible = false">Cancel</n-button>
        <n-space>
          <n-button :loading="loading" @click="handleSubmit">
            Schedule Job
          </n-button>
          <n-button type="primary" :loading="loading" @click="handleExecuteNow">
            Execute Now
          </n-button>
        </n-space>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useMessage } from "naive-ui";
import { useSystemStore } from "../stores/system";
import { useJobsStore } from "../stores/jobs";
import HTTPConfig from "./job-configs/HTTPConfig.vue";
import ShellConfig from "./job-configs/ShellConfig.vue";
import DockerConfig from "./job-configs/DockerConfig.vue";

const props = defineProps({
  show: Boolean,
});

const emit = defineEmits(["update:show"]);

const message = useMessage();
const systemStore = useSystemStore();
const jobsStore = useJobsStore();

const visible = computed({
  get: () => props.show,
  set: (val) => emit("update:show", val),
});

const formRef = ref(null);
const loading = ref(false);
const scheduledTimestamp = ref(Date.now() + 3600000); // Default 1 hour from now
const jobConfig = ref({});

const formValue = ref({
  name: "",
  type: "",
  project_id: "default",
  priority: 5,
  tag_ids: [],
});

const jobTypeOptions = computed(() =>
  systemStore.jobTypes.map((type) => ({
    label: type.toUpperCase(),
    value: type,
  })),
);

const projectOptions = computed(() =>
  systemStore.projects.map((p) => ({ label: p.name, value: p.id })),
);

const tagOptions = computed(() =>
  systemStore.tags.map((t) => ({ label: t.name, value: t.id })),
);

const priorityMarks = {
  1: "1",
  5: "5",
  10: "10",
};

const rules = {
  name: { required: true, message: "Job name is required" },
  type: { required: true, message: "Job type is required" },
  project_id: { required: true, message: "Project is required" },
};

const isDateDisabled = (ts) => {
  return ts < Date.now();
};

watch(
  () => formValue.value.type,
  () => {
    jobConfig.value = {};
  },
);

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();

    const config = JSON.stringify(jobConfig.value);

    const payload = {
      ...formValue.value,
      config,
      immediate: false,
      scheduled_at: new Date(scheduledTimestamp.value).toISOString(),
    };

    loading.value = true;
    await jobsStore.createJob(payload);
    message.success("Job scheduled successfully");
    visible.value = false;
    resetForm();
  } catch (error) {
    message.error(error.message || "Failed to create job");
  } finally {
    loading.value = false;
  }
};

const handleExecuteNow = async () => {
  try {
    await formRef.value?.validate();

    const config = JSON.stringify(jobConfig.value);

    const payload = {
      ...formValue.value,
      config,
      immediate: true,
    };

    loading.value = true;
    await jobsStore.createJob(payload);
    message.success("Job executed successfully");
    visible.value = false;
    resetForm();
  } catch (error) {
    message.error(error.message || "Failed to execute job");
  } finally {
    loading.value = false;
  }
};

const resetForm = () => {
  formValue.value = {
    name: "",
    type: "",
    project_id: "default",
    priority: 5,
    tag_ids: [],
  };
  jobConfig.value = {};
  scheduledTimestamp.value = Date.now() + 3600000;
};
</script>
