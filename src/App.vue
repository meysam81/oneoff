<template>
  <n-config-provider :theme="darkTheme" :hljs="hljs">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <n-layout has-sider style="height: 100vh">
            <n-layout-sider
              v-if="!isLoginPage"
              bordered
              collapse-mode="width"
              :collapsed-width="64"
              :width="240"
              show-trigger
              @collapse="collapsed = true"
              @expand="collapsed = false"
            >
              <Sidebar :collapsed="collapsed" />
            </n-layout-sider>

            <n-layout>
              <n-layout-header
                v-if="!isLoginPage"
                bordered
                style="height: 64px; padding: 0 24px"
              >
                <Header />
              </n-layout-header>

              <n-layout-content :content-style="{ padding: '24px' }">
                <router-view v-slot="{ Component, route }">
                  <transition name="fade" mode="out-in">
                    <keep-alive :max="10">
                      <component
                        :is="Component"
                        :key="route.meta.keepAlive ? undefined : route.fullPath"
                      />
                    </keep-alive>
                  </transition>
                </router-view>
              </n-layout-content>
            </n-layout>
          </n-layout>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup>
import { ref, computed, onMounted, inject } from "vue";
import { darkTheme } from "naive-ui";
import { useRoute } from "vue-router";
import { useSystemStore } from "./stores/system";
import Sidebar from "./components/Sidebar.vue";
import Header from "./components/Header.vue";

const collapsed = ref(false);
const route = useRoute();
const systemStore = useSystemStore();
var hljs = inject("hljs");

const isLoginPage = computed(() => route.path === "/login");

// Initialize app data once on mount
onMounted(async () => {
  await systemStore.initializeApp();
});
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu,
    Cantarell, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  height: 100vh;
}

code {
  font-family: "SF Mono", Monaco, Consolas, "Courier New", monospace;
}

/* Page transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
