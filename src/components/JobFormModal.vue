<template>
  <n-modal
    v-model:show="visible"
    preset="card"
    :title="modalTitle"
    style="width: 700px; max-height: 90vh"
    :mask-closable="false"
  >
    <n-scrollbar style="max-height: 70vh">
      <n-form ref="formRef" :model="formValue" :rules="rules">
        <n-form-item label="Job Name" path="name">
          <n-input
            v-model:value="formValue.name"
            placeholder="Enter job name"
          />
        </n-form-item>

        <!-- Job Type - only shown in create mode -->
        <n-form-item v-if="mode === 'create'" label="Job Type" path="type">
          <n-select
            v-model:value="formValue.type"
            :options="jobTypeOptions"
            placeholder="Select job type"
          />
        </n-form-item>

        <!-- Job Type display - only shown in clone mode -->
        <n-form-item v-if="mode === 'clone'" label="Job Type">
          <n-tag :bordered="false" type="info">{{ props.job?.type }}</n-tag>
        </n-form-item>

        <!-- Project - only shown in create mode -->
        <n-form-item v-if="mode === 'create'" label="Project" path="project_id">
          <n-select
            v-model:value="formValue.project_id"
            :options="projectOptions"
            placeholder="Select project"
          />
        </n-form-item>

        <!-- Execution Mode - radio button approach -->
        <n-form-item label="When to Execute">
          <n-space vertical style="width: 100%">
            <n-radio-group v-model:value="executeMode">
              <n-space>
                <n-radio value="immediate">
                  <span>Execute Immediately</span>
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
              The job will start executing as soon as a worker is available.
            </n-text>
          </n-space>
        </n-form-item>

        <!-- Priority - only shown in create mode -->
        <n-form-item v-if="mode === 'create'" label="Priority" path="priority">
          <n-slider
            v-model:value="formValue.priority"
            :min="1"
            :max="10"
            :step="1"
            :marks="priorityMarks"
          />
        </n-form-item>

        <!-- Tags - only shown in create mode -->
        <n-form-item v-if="mode === 'create'" label="Tags" path="tag_ids">
          <n-select
            v-model:value="formValue.tag_ids"
            multiple
            :options="tagOptions"
            placeholder="Select tags"
          />
        </n-form-item>

        <!-- Job Configuration Section -->
        <n-divider title-placement="left">
          <n-text depth="3" style="font-size: 13px">Job Configuration</n-text>
        </n-divider>

        <component
          :is="configComponent"
          v-if="configComponent"
          v-model="jobConfig"
        />

        <!-- Info alert for clone mode -->
        <n-alert
          v-if="mode === 'clone'"
          type="info"
          :bordered="false"
          style="margin-top: 16px"
        >
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
        <n-button type="primary" :loading="loading" @click="handleSubmit">
          <template #icon v-if="executeMode === 'immediate'">
            <n-icon><PlayOutline /></n-icon>
          </template>
          {{ submitButtonText }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useMessage } from "naive-ui";
import { InformationCircleOutline, PlayOutline } from "@vicons/ionicons5";
import { useSystemStore } from "../stores/system";
import { useJobsStore } from "../stores/jobs";
import HTTPConfig from "./job-configs/HTTPConfig.vue";
import ShellConfig from "./job-configs/ShellConfig.vue";
import DockerConfig from "./job-configs/DockerConfig.vue";

var props = defineProps({
  show: Boolean,
  mode: {
    type: String,
    default: "create",
    validator: function (value) {
      return ["create", "clone"].includes(value);
    },
  },
  job: {
    type: Object,
    default: null,
  },
});

var emit = defineEmits(["update:show", "submitted"]);

var message = useMessage();
var systemStore = useSystemStore();
var jobsStore = useJobsStore();

var visible = computed({
  get: function () {
    return props.show;
  },
  set: function (val) {
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
  type: "",
  project_id: "default",
  priority: 5,
  tag_ids: [],
});

var modalTitle = computed(function () {
  return props.mode === "clone" ? "Clone Job" : "Create New Job";
});

var submitButtonText = computed(function () {
  if (props.mode === "clone") {
    return executeMode.value === "immediate"
      ? "Clone & Execute Now"
      : "Clone Job";
  }
  return executeMode.value === "immediate" ? "Create & Execute Now" : "Create Job";
});

var jobTypeOptions = computed(function () {
  return systemStore.jobTypes.map(function (type) {
    return { label: type.toUpperCase(), value: type };
  });
});

var projectOptions = computed(function () {
  return systemStore.projects.map(function (p) {
    return { label: p.name, value: p.id };
  });
});

var tagOptions = computed(function () {
  return systemStore.tags.map(function (t) {
    return { label: t.name, value: t.id };
  });
});

var priorityMarks = {
  1: "1",
  5: "5",
  10: "10",
};

var rules = computed(function () {
  var baseRules = {
    name: { required: true, message: "Job name is required" },
  };

  if (props.mode === "create") {
    baseRules.type = { required: true, message: "Job type is required" };
    baseRules.project_id = { required: true, message: "Project is required" };
  }

  return baseRules;
});

var configComponents = {
  http: HTTPConfig,
  shell: ShellConfig,
  docker: DockerConfig,
};

var configComponent = computed(function () {
  var jobType = props.mode === "clone" ? props.job?.type : formValue.value.type;
  if (!jobType) return null;
  return configComponents[jobType] || null;
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
  if (props.mode === "clone" && props.job) {
    formValue.value.name = props.job.name + " (Copy)";
    jobConfig.value = parseJobConfig(props.job.config);
  } else {
    formValue.value = {
      name: "",
      type: "",
      project_id: "default",
      priority: 5,
      tag_ids: [],
    };
    jobConfig.value = {};
  }
  executeMode.value = "immediate";
  scheduledTimestamp.value = Date.now() + 3600000;
}

// Watch for job changes in clone mode
watch(
  function () {
    return props.job;
  },
  function (newJob) {
    if (newJob && props.mode === "clone") {
      resetForm();
    }
  },
  { immediate: true }
);

// Watch for show changes to reset form
watch(
  function () {
    return props.show;
  },
  function (newShow) {
    if (newShow) {
      resetForm();
    }
  }
);

// Watch for type changes in create mode to reset config
watch(
  function () {
    return formValue.value.type;
  },
  function () {
    if (props.mode === "create") {
      jobConfig.value = {};
    }
  }
);

async function handleSubmit() {
  try {
    await formRef.value?.validate();

    var scheduledAt;
    if (executeMode.value === "immediate") {
      scheduledAt = "now";
    } else {
      scheduledAt = new Date(scheduledTimestamp.value).toISOString();
    }

    var payload;

    if (props.mode === "clone") {
      if (!props.job) {
        message.error("No job selected for cloning");
        return;
      }

      payload = {
        name: formValue.value.name,
        type: props.job.type,
        config: JSON.stringify(jobConfig.value),
        scheduled_at: scheduledAt,
        priority: props.job.priority,
        project_id: props.job.project_id,
        timezone: props.job.timezone,
        tag_ids:
          props.job.tags?.map(function (t) {
            return t.id;
          }) || [],
      };
    } else {
      payload = {
        ...formValue.value,
        config: JSON.stringify(jobConfig.value),
        scheduled_at: scheduledAt,
      };
    }

    loading.value = true;
    await jobsStore.createJob(payload);

    var successMessage;
    if (props.mode === "clone") {
      successMessage =
        executeMode.value === "immediate"
          ? "Job cloned and queued for immediate execution"
          : "Job cloned successfully";
    } else {
      successMessage =
        executeMode.value === "immediate"
          ? "Job created and queued for immediate execution"
          : "Job scheduled successfully";
    }

    message.success(successMessage);
    visible.value = false;
    emit("submitted");
  } catch (error) {
    message.error(error.message || "Failed to " + (props.mode === "clone" ? "clone" : "create") + " job");
  } finally {
    loading.value = false;
  }
}
</script>
