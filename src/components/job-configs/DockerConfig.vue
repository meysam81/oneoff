<template>
  <n-space vertical>
    <n-form-item label="Docker Image">
      <n-input v-model:value="config.image" placeholder="nginx:latest" />
    </n-form-item>

    <n-form-item label="Command">
      <n-dynamic-input
        v-model:value="config.command"
        placeholder="Add command argument"
      >
        <template #default="{ value, index }">
          <n-input v-model:value="config.command[index]" placeholder="Argument" />
        </template>
      </n-dynamic-input>
    </n-form-item>

    <n-form-item label="Environment Variables">
      <n-input
        v-model:value="envText"
        type="textarea"
        placeholder="KEY1=value1
KEY2=value2
KEY3=value3"
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

const envText = ref("");
const volumesJson = ref("{}");

// Parse .env format (KEY=value) to map
const parseEnvFormat = (text) => {
  const env = {};
  if (!text || text.trim() === "") {
    return env;
  }

  const lines = text.split("\n");
  for (const line of lines) {
    const trimmed = line.trim();
    // Skip empty lines and comments
    if (trimmed === "" || trimmed.startsWith("#")) {
      continue;
    }

    const equalIndex = trimmed.indexOf("=");
    if (equalIndex > 0) {
      const key = trimmed.substring(0, equalIndex).trim();
      const value = trimmed.substring(equalIndex + 1).trim();
      if (key) {
        env[key] = value;
      }
    }
  }

  return env;
};

watch(envText, (val) => {
  config.value.env = parseEnvFormat(val);
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
