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
          <n-input
            v-model:value="config.command[index]"
            placeholder="Argument"
          />
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

var props = defineProps({
  modelValue: Object,
});

var emit = defineEmits(["update:modelValue"]);

function getDefaultConfig() {
  return {
    image: "",
    command: [],
    env: {},
    volumes: {},
    workdir: "",
    auto_remove: true,
    timeout: 300,
  };
}

function envToText(envObj) {
  if (!envObj || typeof envObj !== "object") {
    return "";
  }
  var lines = [];
  var keys = Object.keys(envObj);
  for (var i = 0; i < keys.length; i++) {
    lines.push(keys[i] + "=" + envObj[keys[i]]);
  }
  return lines.join("\n");
}

function parseEnvFormat(text) {
  var env = {};
  if (!text || text.trim() === "") {
    return env;
  }

  var lines = text.split("\n");
  for (var i = 0; i < lines.length; i++) {
    var trimmed = lines[i].trim();
    if (trimmed === "" || trimmed.startsWith("#")) {
      continue;
    }

    var equalIndex = trimmed.indexOf("=");
    if (equalIndex > 0) {
      var key = trimmed.substring(0, equalIndex).trim();
      var value = trimmed.substring(equalIndex + 1).trim();
      if (key) {
        env[key] = value;
      }
    }
  }

  return env;
}

function initializeFromProps() {
  var initial = props.modelValue || {};
  config.value = {
    image: initial.image || "",
    command: initial.command || [],
    env: initial.env || {},
    volumes: initial.volumes || {},
    workdir: initial.workdir || "",
    auto_remove: initial.auto_remove !== undefined ? initial.auto_remove : true,
    timeout: initial.timeout || 300,
  };
  envText.value = envToText(config.value.env);
  volumesJson.value = JSON.stringify(config.value.volumes || {}, null, 2);
}

var config = ref(getDefaultConfig());
var envText = ref("");
var volumesJson = ref("{}");

initializeFromProps();

watch(
  function watchModelValue() {
    return props.modelValue;
  },
  function onModelValueChange(newVal) {
    if (newVal) {
      initializeFromProps();
    }
  },
  { deep: true },
);

watch(envText, function onEnvTextChange(val) {
  config.value.env = parseEnvFormat(val);
});

watch(volumesJson, function onVolumesJsonChange(val) {
  try {
    config.value.volumes = JSON.parse(val);
  } catch (e) {}
});

watch(
  config,
  function onConfigChange(val) {
    emit("update:modelValue", val);
  },
  { deep: true },
);
</script>
