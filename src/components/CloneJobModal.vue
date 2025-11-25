<template>
  <n-modal
    v-model:show="visible"
    preset="card"
    title="Clone Job"
    style="width: 700px; max-height: 90vh"
    :mask-closable="false"
  >
    <n-scrollbar style="max-height: 70vh">
      <n-form ref="formRef" :model="formValue" :rules="rules">
        <n-form-item label="New Job Name" path="name">
          <n-input
            v-model:value="formValue.name"
            placeholder="Enter new job name"
          />
        </n-form-item>

        <n-form-item label="When to Execute">
          <n-space vertical style="width: 100%">
            <n-radio-group v-model:value="executeMode">
              <n-space>
                <n-radio value="immediate">
                  <n-space align="center" :size="4">
                    <span>Execute Immediately</span>
                  </n-space>
                </n-radio>
                <n-radio value="scheduled">
                  <span>Schedule for Later</span>
                </n-radio>
              </n-space>
            </n-radio-group>

            <n-collapse-transition :show="executeMode === 'scheduled'">
              <n-date-picker
                v-model:value="scheduledTimestamp"
                type="datetime"
                clearable
                style="width: 100%; margin-top: 8px"
                :is-date-disabled="isDateDisabled"
                placeholder="Select date and time"
              />
            </n-collapse-transition>

            <n-text
              v-if="executeMode === 'immediate'"
              depth="3"
              style="font-size: 12px"
            >
              The cloned job will start executing as soon as a worker is
              available.
            </n-text>
          </n-space>
        </n-form-item>

        <n-divider title-placement="left">
          <n-text depth="3" style="font-size: 13px">Job Configuration</n-text>
        </n-divider>

        <n-form-item label="Job Type">
          <n-tag :bordered="false" type="info">{{ props.job?.type }}</n-tag>
        </n-form-item>

        <component
          :is="configComponent"
          v-if="configComponent"
          v-model="jobConfig"
        />

        <n-alert type="info" :bordered="false" style="margin-top: 16px">
          <template #icon>
            <n-icon><InformationCircleOutline /></n-icon>
          </template>
          The cloned job will inherit the original job's project, tags, and
          priority settings.
        </n-alert>
      </n-form>
    </n-scrollbar>

    <template #footer>
      <n-space justify="space-between" align="center">
        <n-button @click="visible = false">Cancel</n-button>
        <n-button type="primary" :loading="loading" @click="handleClone">
          <template #icon v-if="executeMode === 'immediate'">
            <n-icon><PlayOutline /></n-icon>
          </template>
          {{
            executeMode === "immediate" ? "Clone & Execute Now" : "Clone Job"
          }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useMessage } from "naive-ui";
import { InformationCircleOutline, PlayOutline } from "@vicons/ionicons5";
import { useJobsStore } from "../stores/jobs";
import HTTPConfig from "./job-configs/HTTPConfig.vue";
import ShellConfig from "./job-configs/ShellConfig.vue";
import DockerConfig from "./job-configs/DockerConfig.vue";

var props = defineProps({
  show: Boolean,
  job: Object,
});

var emit = defineEmits(["update:show", "cloned"]);

var message = useMessage();
var jobsStore = useJobsStore();

var visible = computed({
  get: function getVisible() {
    return props.show;
  },
  set: function setVisible(val) {
    emit("update:show", val);
  },
});

var formRef = ref(null);
var loading = ref(false);
var executeMode = ref("immediate");
var scheduledTimestamp = ref(Date.now() + 3600000);
var jobConfig = ref({});

var formValue = ref({
  name: "",
});

var rules = {
  name: { required: true, message: "Job name is required" },
};

var configComponents = {
  http: HTTPConfig,
  shell: ShellConfig,
  docker: DockerConfig,
};

var configComponent = computed(function getConfigComponent() {
  if (!props.job?.type) return null;
  return configComponents[props.job.type] || null;
});

function isDateDisabled(ts) {
  return ts < Date.now();
}

function parseJobConfig(configStr) {
  try {
    return JSON.parse(configStr);
  } catch (e) {
    return {};
  }
}

function resetForm() {
  if (props.job) {
    formValue.value.name = props.job.name + " (Copy)";
    jobConfig.value = parseJobConfig(props.job.config);
  }
  executeMode.value = "immediate";
  scheduledTimestamp.value = Date.now() + 3600000;
}

watch(
  function watchJob() {
    return props.job;
  },
  function onJobChange(newJob) {
    if (newJob) {
      resetForm();
    }
  },
  { immediate: true },
);

watch(
  function watchShow() {
    return props.show;
  },
  function onShowChange(newShow) {
    if (newShow && props.job) {
      resetForm();
    }
  },
);

async function handleClone() {
  try {
    await formRef.value?.validate();

    if (!props.job) {
      message.error("No job selected for cloning");
      return;
    }

    var scheduledAt;
    if (executeMode.value === "immediate") {
      scheduledAt = "now";
    } else {
      scheduledAt = new Date(scheduledTimestamp.value).toISOString();
    }

    loading.value = true;

    await jobsStore.createJob({
      name: formValue.value.name,
      type: props.job.type,
      config: JSON.stringify(jobConfig.value),
      scheduled_at: scheduledAt,
      priority: props.job.priority,
      project_id: props.job.project_id,
      timezone: props.job.timezone,
      tag_ids:
        props.job.tags?.map(function getTagId(t) {
          return t.id;
        }) || [],
    });

    var successMessage =
      executeMode.value === "immediate"
        ? "Job cloned and queued for immediate execution"
        : "Job cloned successfully";
    message.success(successMessage);
    visible.value = false;
    emit("cloned");
  } catch (error) {
    message.error(error.message || "Failed to clone job");
  } finally {
    loading.value = false;
  }
}
</script>
