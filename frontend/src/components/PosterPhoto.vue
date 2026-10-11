<script setup lang="ts">
import { ref, watch } from "vue";
const props = defineProps<{ src: string; alt: string }>();
const failed = ref(false);
watch(() => props.src, () => { failed.value = false; });
</script>
<template>
  <img v-show="!failed" :src="src" :alt="alt" loading="lazy" @error="failed = true" @load="failed = false" />
</template>
