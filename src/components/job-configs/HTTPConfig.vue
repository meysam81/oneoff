<template>
  <n-space vertical>
    <n-form-item label="URL" required>
      <n-input
        v-model:value="config.url"
        placeholder="https://api.example.com/endpoint"
      />
    </n-form-item>

    <n-form-item label="Method">
      <n-select v-model:value="config.method" :options="methodOptions" />
    </n-form-item>

    <n-form-item label="Headers (JSON)">
      <n-input
        v-model:value="headersJson"
        type="textarea"
        placeholder='{"Content-Type": "application/json"}'
        :rows="3"
      />
    </n-form-item>

    <n-form-item label="Body">
      <n-input
        v-model:value="config.body"
        type="textarea"
        placeholder="Request body (for POST/PUT/PATCH)"
        :rows="4"
      />
    </n-form-item>

    <n-form-item label="Timeout (seconds)">
      <n-input-number v-model:value="config.timeout" :min="1" :max="300" />
    </n-form-item>

    <n-form-item label="Retry Count">
      <n-input-number v-model:value="config.retry_count" :min="0" :max="10" />
    </n-form-item>

    <n-form-item label="Expected Status Code">
      <n-input-number
        v-model:value="config.expected_status"
        :min="100"
        :max="599"
      />
    </n-form-item>
  </n-space>
</template>

<script setup>
import { ref, watch } from "vue";

var props = defineProps({
  modelValue: Object,
});

var emit = defineEmits(["update:modelValue"]);

var methodOptions = [
  { label: "GET", value: "GET" },
  { label: "POST", value: "POST" },
  { label: "PUT", value: "PUT" },
  { label: "PATCH", value: "PATCH" },
  { label: "DELETE", value: "DELETE" },
];

function getDefaultConfig() {
  return {
    url: "",
    method: "GET",
    headers: {},
    body: "",
    timeout: 30,
    retry_count: 0,
    expected_status: 200,
  };
}

function initializeFromProps() {
  var initial = props.modelValue || {};
  config.value = {
    url: initial.url || "",
    method: initial.method || "GET",
    headers: initial.headers || {},
    body: initial.body || "",
    timeout: initial.timeout || 30,
    retry_count: initial.retry_count || 0,
    expected_status: initial.expected_status || 200,
  };
  headersJson.value = JSON.stringify(config.value.headers || {}, null, 2);
}

var config = ref(getDefaultConfig());
var headersJson = ref("{}");

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

watch(headersJson, function onHeadersJsonChange(val) {
  try {
    config.value.headers = JSON.parse(val);
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
