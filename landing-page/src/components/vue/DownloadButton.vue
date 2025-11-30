<script setup lang="ts">
import { usePlatform } from "../../composables/usePlatform";

var { displayName } = usePlatform();
</script>

<template>
  <span class="hidden sm:inline">Download for {{ displayName }}</span>
</template>
