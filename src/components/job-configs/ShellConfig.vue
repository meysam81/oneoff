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
        :placeholder="config.is_path ? '/path/to/script.sh' : 'echo &quot;Hello World&quot;'"
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

    <n-form-item label="Environment Variables (JSON)">
      <n-input
        v-model:value="envJson"
        type="textarea"
        placeholder='{"KEY": "value"}'
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
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Object,
})

const emit = defineEmits(['update:modelValue'])

const config = ref({
  script: '',
  is_path: false,
  args: [],
  env: {},
  workdir: '',
  timeout: 60,
})

const envJson = ref('{}')

watch(envJson, (val) => {
  try {
    config.value.env = JSON.parse(val)
  } catch (e) {
    // Invalid JSON
  }
})

watch(config, (val) => {
  emit('update:modelValue', val)
}, { deep: true })
</script>
