<template>
  <n-space vertical>
    <n-form-item label="Docker Image">
      <n-input v-model:value="config.image" placeholder="nginx:latest" />
    </n-form-item>

    <n-form-item label="Command">
      <n-dynamic-input
        v-model:value="config.command"
        placeholder="Add command argument"
        #="{ value }"
      >
        <n-input v-model:value="value.value" placeholder="Argument" />
      </n-dynamic-input>
    </n-form-item>

    <n-form-item label="Environment Variables (JSON)">
      <n-input
        v-model:value="envJson"
        type="textarea"
        placeholder='{"KEY": "value"}'
        :rows="3"
      />
    </n-form-item>

    <n-form-item label="Volumes (JSON)">
      <n-input
        v-model:value="volumesJson"
        type="textarea"
        placeholder='{"/host/path": "/container/path"}'
        :rows="3"
      />
    </n-form-item>

    <n-form-item label="Working Directory">
      <n-input v-model:value="config.workdir" placeholder="/app" />
    </n-form-item>

    <n-form-item label="Auto Remove">
      <n-switch v-model:value="config.auto_remove" />
    </n-form-item>

    <n-form-item label="Timeout (seconds)">
      <n-input-number v-model:value="config.timeout" :min="1" :max="3600" />
    </n-form-item>
  </n-space>
</template>

<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  modelValue: Object,
});

const emit = defineEmits(["update:modelValue"]);

const config = ref({
  image: "",
  command: [],
  env: {},
  volumes: {},
  workdir: "",
  auto_remove: true,
  timeout: 300,
});

const envJson = ref("{}");
const volumesJson = ref("{}");

watch(envJson, (val) => {
  try {
    config.value.env = JSON.parse(val);
  } catch (e) {
    // Invalid JSON
  }
});

watch(volumesJson, (val) => {
  try {
    config.value.volumes = JSON.parse(val);
  } catch (e) {
    // Invalid JSON
  }
});

watch(
  config,
  (val) => {
    emit("update:modelValue", val);
  },
  { deep: true },
);
</script>
