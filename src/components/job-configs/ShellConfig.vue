<template>
  <n-space vertical>
    <n-form-item label="Script Type">
      <n-radio-group v-model:value="config.is_path">
        <n-radio :value="false">Inline Script</n-radio>
        <n-radio :value="true">Script File Path</n-radio>
      </n-radio-group>
    </n-form-item>

    <n-form-item :label="config.is_path ? 'Script Path' : 'Script Content'">
      <n-input
        v-model:value="config.script"
        type="textarea"
        :placeholder="
          config.is_path ? '/path/to/script.sh' : 'echo &quot;Hello World&quot;'
        "
        :rows="config.is_path ? 1 : 6"
      />
    </n-form-item>

    <n-form-item label="Arguments">
      <n-dynamic-input
        v-model:value="config.args"
        placeholder="Add argument"
        #="{ value }"
      >
        <n-input v-model:value="value.value" placeholder="Argument" />
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

    <n-form-item label="Working Directory">
      <n-input v-model:value="config.workdir" placeholder="/path/to/workdir" />
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
    script: "",
    is_path: false,
    args: [],
    env: {},
    workdir: "",
    timeout: 60,
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
    script: initial.script || "",
    is_path: initial.is_path !== undefined ? initial.is_path : false,
    args: initial.args || [],
    env: initial.env || {},
    workdir: initial.workdir || "",
    timeout: initial.timeout || 60,
  };
  envText.value = envToText(config.value.env);
}

var config = ref(getDefaultConfig());
var envText = ref("");

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

watch(
  config,
  function onConfigChange(val) {
    emit("update:modelValue", val);
  },
  { deep: true },
);
</script>
