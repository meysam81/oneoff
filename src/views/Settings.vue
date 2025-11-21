<template>
  <div class="settings-container">
    <n-space vertical :size="20">
      <!-- System Information -->
      <n-card title="System Information">
        <n-spin :show="loading">
          <n-descriptions v-if="systemStats" :column="2" bordered>
            <n-descriptions-item label="Total Jobs">
              {{ systemStats.total_jobs || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Scheduled Jobs">
              {{ systemStats.scheduled_jobs || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Running Jobs">
              {{ systemStats.running_jobs || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Completed Jobs">
              {{ systemStats.completed_jobs || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Failed Jobs">
              {{ systemStats.failed_jobs || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Total Executions">
              {{ systemStats.total_executions || 0 }}
            </n-descriptions-item>
          </n-descriptions>
        </n-spin>
      </n-card>

      <!-- Worker Status -->
      <n-card title="Worker Status">
        <n-spin :show="loading">
          <n-descriptions v-if="workerStatus" :column="2" bordered>
            <n-descriptions-item label="Total Workers">
              {{ workerStatus.total_workers || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Active Workers">
              {{ workerStatus.active_workers || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Idle Workers">
              {{ workerStatus.idle_workers || 0 }}
            </n-descriptions-item>
            <n-descriptions-item label="Pool Status">
              <n-tag :type="workerStatus.is_running ? 'success' : 'error'">
                {{ workerStatus.is_running ? "Running" : "Stopped" }}
              </n-tag>
            </n-descriptions-item>
          </n-descriptions>
        </n-spin>
      </n-card>

      <!-- Configuration Settings -->
      <n-card title="Configuration Settings">
        <n-spin :show="loading">
          <n-form
            ref="formRef"
            :model="configForm"
            label-placement="left"
            label-width="200"
          >
            <n-form-item label="Log Retention Days">
              <n-input-number
                v-model:value="configForm.log_retention_days"
                :min="1"
                :max="365"
                style="width: 200px"
              />
            </n-form-item>

            <n-form-item label="Default Priority">
              <n-slider
                v-model:value="configForm.default_priority"
                :min="1"
                :max="10"
                :step="1"
                :marks="{ 1: '1', 5: '5', 10: '10' }"
                style="width: 300px"
              />
            </n-form-item>

            <n-form-item label="Default Timezone">
              <n-select
                v-model:value="configForm.default_timezone"
                :options="timezoneOptions"
                filterable
                style="width: 300px"
              />
            </n-form-item>
          </n-form>

          <n-space justify="end">
            <n-button @click="resetForm">Reset</n-button>
            <n-button
              type="primary"
              :loading="saving"
              @click="saveSettings"
            >
              Save Changes
            </n-button>
          </n-space>
        </n-spin>
      </n-card>

      <!-- Configuration List (Read-only) -->
      <n-card title="All Configuration Values">
        <n-spin :show="loading">
          <n-data-table
            v-if="configList.length"
            :columns="configColumns"
            :data="configList"
            :pagination="false"
          />
          <n-empty v-else description="No configuration found" />
        </n-spin>
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted, h } from "vue";
import { useMessage } from "naive-ui";
import { NTag } from "naive-ui";
import { systemAPI } from "../utils/api";
import { useSystemStore } from "../stores/system";

const message = useMessage();
const systemStore = useSystemStore();

const loading = ref(false);
const saving = ref(false);
const systemStats = ref(null);
const workerStatus = ref(null);
const configList = ref([]);

const configForm = ref({
  log_retention_days: 90,
  default_priority: 5,
  default_timezone: "UTC",
});

const originalConfig = ref({});

const timezoneOptions = [
  { label: "UTC", value: "UTC" },
  { label: "America/New_York", value: "America/New_York" },
  { label: "America/Chicago", value: "America/Chicago" },
  { label: "America/Denver", value: "America/Denver" },
  { label: "America/Los_Angeles", value: "America/Los_Angeles" },
  { label: "Europe/London", value: "Europe/London" },
  { label: "Europe/Paris", value: "Europe/Paris" },
  { label: "Europe/Berlin", value: "Europe/Berlin" },
  { label: "Asia/Tokyo", value: "Asia/Tokyo" },
  { label: "Asia/Shanghai", value: "Asia/Shanghai" },
  { label: "Asia/Dubai", value: "Asia/Dubai" },
  { label: "Australia/Sydney", value: "Australia/Sydney" },
];

const configColumns = [
  {
    title: "Key",
    key: "key",
    width: 300,
  },
  {
    title: "Value",
    key: "value",
    ellipsis: { tooltip: true },
  },
  {
    title: "Last Updated",
    key: "updated_at",
    render: (row) => new Date(row.updated_at).toLocaleString(),
  },
];

const fetchData = async () => {
  loading.value = true;
  try {
    const [statsResponse, workerResponse, configResponse] = await Promise.all([
      systemAPI.status(),
      systemAPI.workerStatus(),
      systemAPI.config(),
    ]);

    systemStats.value = statsResponse.data;
    workerStatus.value = workerResponse.data;
    configList.value = configResponse.data || [];

    // Parse config values and populate form
    const configMap = {};
    configList.value.forEach((item) => {
      try {
        configMap[item.key] = JSON.parse(item.value);
      } catch {
        configMap[item.key] = item.value;
      }
    });

    configForm.value = {
      log_retention_days: configMap.log_retention_days || 90,
      default_priority: configMap.default_priority || 5,
      default_timezone: configMap.default_timezone || "UTC",
    };

    originalConfig.value = { ...configForm.value };
  } catch (error) {
    console.error("Failed to fetch settings:", error);
    message.error("Failed to load settings");
  } finally {
    loading.value = false;
  }
};

const saveSettings = async () => {
  saving.value = true;
  try {
    const updates = [];

    for (const [key, value] of Object.entries(configForm.value)) {
      if (value !== originalConfig.value[key]) {
        updates.push(
          systemAPI.updateConfig(key, JSON.stringify(value))
        );
      }
    }

    if (updates.length === 0) {
      message.info("No changes to save");
      return;
    }

    await Promise.all(updates);
    message.success("Settings saved successfully");
    originalConfig.value = { ...configForm.value };

    // Refresh data
    await fetchData();

    // Invalidate system store cache
    systemStore.invalidateCache();
  } catch (error) {
    console.error("Failed to save settings:", error);
    message.error("Failed to save settings");
  } finally {
    saving.value = false;
  }
};

const resetForm = () => {
  configForm.value = { ...originalConfig.value };
  message.info("Form reset to last saved values");
};

onMounted(fetchData);
</script>

<style scoped>
.settings-container {
  padding: 20px;
}
</style>
