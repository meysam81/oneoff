<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { usePlatform } from "../../composables/usePlatform";

var props = defineProps<{
  downloadCommands: Record<string, string>;
  shortCommands: Record<string, string>;
  runCommands: Record<string, string>;
  promptChars: Record<string, string>;
}>();

var { platform, displayName } = usePlatform();

var line2Visible = ref(false);
var line3Visible = ref(false);
var line4Visible = ref(false);

var downloadCommand = computed(function () {
  return (
    props.downloadCommands[platform.value] ||
    props.downloadCommands["linux-amd64"]
  );
});

var shortCommand = computed(function () {
  return (
    props.shortCommands[platform.value] || props.shortCommands["linux-amd64"]
  );
});

var runCommand = computed(function () {
  return props.runCommands[platform.value] || props.runCommands["linux-amd64"];
});

var promptChar = computed(function () {
  const prmpt = props.promptChars[platform.value] || "$";
  return prmpt + " ";
});

function animateTerminal() {
  setTimeout(function () {
    line2Visible.value = true;
  }, 1500);
  setTimeout(function () {
    line3Visible.value = true;
  }, 2200);
  setTimeout(function () {
    line4Visible.value = true;
  }, 3000);
}

onMounted(function () {
  var observer = new IntersectionObserver(function (entries) {
    entries.forEach(function (entry) {
      if (entry.isIntersecting) {
        animateTerminal();
        observer.disconnect();
      }
    });
  });

  var el = document.getElementById("hero-terminal-vue");
  if (el) observer.observe(el);
});

defineExpose({ displayName });
</script>

<template>
  <div
    class="terminal-window max-w-2xl lg:max-w-4xl xl:max-w-5xl mx-auto text-left shadow-2xl"
  >
    <div class="terminal-header">
      <div class="flex items-center gap-1.5 sm:gap-2">
        <span class="terminal-dot terminal-dot-red"></span>
        <span class="terminal-dot terminal-dot-yellow"></span>
        <span class="terminal-dot terminal-dot-green"></span>
      </div>
      <span class="text-xs text-fg-muted ml-2 sm:ml-4">Terminal</span>
    </div>
    <div class="terminal-body overflow-x-auto" id="hero-terminal-vue">
      <div class="code-line whitespace-nowrap">
        <span class="code-prompt">{{ promptChar }}</span>
        <span class="code-command text-xs sm:text-sm lg:text-base">
          <span class="hidden lg:inline">{{ downloadCommand }}</span>
          <span class="hidden sm:inline lg:hidden">{{ shortCommand }}</span>
          <span class="sm:hidden">{{ shortCommand }}</span>
        </span>
      </div>
      <div
        class="code-line"
        :class="{ 'opacity-0': !line2Visible, 'animate-fade-in': line2Visible }"
      >
        <span class="code-prompt">{{ promptChar }}</span>
        <span class="code-command text-xs sm:text-sm lg:text-base">
          {{ runCommand }}</span
        >
      </div>
      <div
        class="code-line"
        :class="{ 'opacity-0': !line3Visible, 'animate-fade-in': line3Visible }"
      >
        <span class="code-output text-xs sm:text-sm lg:text-base"
          >→ Server running at http://localhost:8080</span
        >
      </div>
      <div
        class="code-line mt-2"
        :class="{ 'opacity-0': !line4Visible, 'animate-fade-in': line4Visible }"
      >
        <span class="text-success text-xs sm:text-sm lg:text-base"
          >✓ Ready to schedule jobs!</span
        >
        <span class="cursor-blink"></span>
      </div>
    </div>
  </div>
</template>
