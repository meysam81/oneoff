<template>
  <n-space vertical>
    <n-form-item label="URL">
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
        placeholder="Request body"
        :rows="4"
      />
    </n-form-item>

    <n-form-item label="Timeout (seconds)">
      <n-input-number v-model:value="config.timeout" :min="1" :max="300" />
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
  url: "",
  method: "GET",
  headers: {},
  body: "",
  timeout: 30,
});

const headersJson = ref("{}");

const methodOptions = [
  { label: "GET", value: "GET" },
  { label: "POST", value: "POST" },
  { label: "PUT", value: "PUT" },
  { label: "PATCH", value: "PATCH" },
  { label: "DELETE", value: "DELETE" },
];

watch(headersJson, (val) => {
  try {
    config.value.headers = JSON.parse(val);
  } catch (e) {
    // Invalid JSON, ignore
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
