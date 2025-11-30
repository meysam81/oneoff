<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { usePlatform, type Platform } from "../../composables/usePlatform";

var props = defineProps<{
  platformCommands: Record<string, string>;
  terminalTitles: Record<string, string>;
  platforms: Array<{ id: string; name: string; arch: string }>;
}>();

var { platform } = usePlatform();
var currentPlatform = ref<string>("linux-amd64");
var copyFeedback = ref(false);

onMounted(function () {
  currentPlatform.value = platform.value;
});

var currentCommand = computed(function () {
  return props.platformCommands[currentPlatform.value] || "";
});

var currentTerminalTitle = computed(function () {
  return props.terminalTitles[currentPlatform.value] || "Terminal";
});

var formattedCommand = computed(function () {
  var lines = currentCommand.value.split("\n");
  return lines
    .map(function (line) {
      return "$ " + line;
    })
    .join("\n");
});

function selectPlatform(platformId: string) {
  currentPlatform.value = platformId;
}

function isSelected(platformId: string): boolean {
  return currentPlatform.value === platformId;
}

function copyToClipboard() {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(currentCommand.value).then(function () {
      showFeedback();
    });
    return;
  }

  var textarea = document.createElement("textarea");
  textarea.value = currentCommand.value;
  textarea.style.position = "fixed";
  textarea.style.left = "-9999px";
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand("copy");
  document.body.removeChild(textarea);
  showFeedback();
}

function showFeedback() {
  copyFeedback.value = true;
  setTimeout(function () {
    copyFeedback.value = false;
  }, 2000);
}
</script>

<template>
  <div>
    <div
      class="flex items-center justify-center gap-1 sm:gap-2 mb-4 sm:mb-6 overflow-x-auto pb-2"
    >
      <button
        v-for="p in platforms"
        :key="p.id"
        :class="[
          'flex-shrink-0 px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium transition-all duration-200',
          isSelected(p.id)
            ? 'bg-bg-tertiary text-accent-primary border border-accent-primary/30'
            : 'text-fg-muted hover:text-fg-secondary hover:bg-bg-tertiary/50',
        ]"
        @click="selectPlatform(p.id)"
      >
        {{ p.name }}
        <span class="text-[10px] sm:text-xs text-fg-muted ml-1 hidden xs:inline"
          >({{ p.arch }})</span
        >
      </button>
    </div>

    <div class="terminal-window mb-4 sm:mb-6">
      <div class="terminal-header">
        <div class="flex items-center gap-1.5 sm:gap-2">
          <span class="terminal-dot terminal-dot-red"></span>
          <span class="terminal-dot terminal-dot-yellow"></span>
          <span class="terminal-dot terminal-dot-green"></span>
        </div>
        <span class="text-[10px] sm:text-xs text-fg-muted ml-2 sm:ml-4">{{
          currentTerminalTitle
        }}</span>
      </div>
      <div class="terminal-body relative group">
        <pre
          class="text-[10px] sm:text-sm text-fg-secondary break-all leading-relaxed whitespace-pre-wrap"
          >{{ formattedCommand }}</pre
        >
        <button
          class="absolute top-1.5 sm:top-2 right-1.5 sm:right-2 p-1.5 sm:p-2 bg-bg-tertiary rounded-lg text-fg-muted hover:text-fg-primary sm:opacity-0 sm:group-hover:opacity-100 transition-all duration-200"
          title="Copy to clipboard"
          @click="copyToClipboard"
        >
          <svg
            class="w-3.5 h-3.5 sm:w-4 sm:h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path
              d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"
            ></path>
          </svg>
        </button>
      </div>
    </div>

    <div
      v-if="copyFeedback"
      class="text-center text-xs sm:text-sm text-success mb-4 sm:mb-6"
    >
      <svg
        class="w-3.5 h-3.5 sm:w-4 sm:h-4 inline mr-1"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <polyline points="20 6 9 17 4 12"></polyline>
      </svg>
      Copied to clipboard!
    </div>
  </div>
</template>
