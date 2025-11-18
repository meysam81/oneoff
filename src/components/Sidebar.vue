<template>
  <div class="sidebar">
    <div class="logo" :class="{ collapsed }">
      <span v-if="!collapsed" class="logo-text">OneOff</span>
      <span v-else class="logo-icon">🎯</span>
    </div>

    <n-menu
      :collapsed="collapsed"
      :collapsed-width="64"
      :collapsed-icon-size="22"
      :options="menuOptions"
      :value="currentRoute"
      @update:value="handleMenuSelect"
    />
  </div>
</template>

<script setup>
import { computed, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NIcon } from "naive-ui";
import {
  HomeOutline,
  TimeOutline,
  ListOutline,
  FolderOutline,
  SettingsOutline,
} from "@vicons/ionicons5";

const props = defineProps({
  collapsed: Boolean,
});

const router = useRouter();
const route = useRoute();

const currentRoute = computed(() => route.name);

const renderIcon = (icon) => {
  return () => h(NIcon, null, { default: () => h(icon) });
};

const menuOptions = [
  {
    label: "Dashboard",
    key: "Dashboard",
    icon: renderIcon(HomeOutline),
  },
  {
    label: "Jobs",
    key: "Jobs",
    icon: renderIcon(TimeOutline),
  },
  {
    label: "Executions",
    key: "Executions",
    icon: renderIcon(ListOutline),
  },
  {
    label: "Projects",
    key: "Projects",
    icon: renderIcon(FolderOutline),
  },
  {
    label: "Settings",
    key: "Settings",
    icon: renderIcon(SettingsOutline),
  },
];

const handleMenuSelect = (key) => {
  router.push({ name: key });
};
</script>

<style scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 600;
  color: #6366f1;
  border-bottom: 1px solid rgba(255, 255, 255, 0.09);
  transition: all 0.3s;
}

.logo.collapsed {
  font-size: 24px;
}

.logo-text {
  padding: 0 24px;
}

.logo-icon {
  font-size: 28px;
}
</style>
