<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { usePlatform } from "../../composables/usePlatform";

var props = defineProps<{
  downloadCommands: Record<string, string>;
  runCommands: Record<string, string>;
  openCommands: Record<string, string>;
}>();

var { platform } = usePlatform();

var downloadCommand = computed(function () {
  return (
    props.downloadCommands[platform.value] ||
    props.downloadCommands["linux-amd64"]
  );
});

var runCommand = computed(function () {
  return props.runCommands[platform.value] || props.runCommands["linux-amd64"];
});

var openCommand = computed(function () {
  return (
    props.openCommands[platform.value] || props.openCommands["linux-amd64"]
  );
});

watch(
  platform,
  function (newPlatform) {
    if (typeof document === "undefined") return;

    var downloadEl = document.getElementById("step-download-cmd");
    var runEl = document.getElementById("step-run-cmd");
    var openEl = document.getElementById("step-open-cmd");

    if (downloadEl)
      downloadEl.textContent =
        props.downloadCommands[newPlatform] ||
        props.downloadCommands["linux-amd64"];
    if (runEl)
      runEl.textContent =
        props.runCommands[newPlatform] || props.runCommands["linux-amd64"];
    if (openEl)
      openEl.textContent =
        props.openCommands[newPlatform] || props.openCommands["linux-amd64"];
  },
  { immediate: true },
);
</script>

<template>
  <div class="steps-commands" style="display: contents">
    <slot
      :downloadCommand="downloadCommand"
      :runCommand="runCommand"
      :openCommand="openCommand"
    />
  </div>
</template>
