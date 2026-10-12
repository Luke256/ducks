<script setup lang="ts">
defineProps<{ iconOnly?: boolean }>();
</script>

<template>
  <span class="event-lock-badge" :class="{ badge: !iconOnly }"
    :role="iconOnly ? 'img' : undefined" :aria-label="iconOnly ? 'ロック中' : undefined"
    :title="iconOnly ? 'ロック中' : undefined">
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"
      stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <rect x="4" y="10" width="16" height="12" rx="2" />
      <path d="M8 10V6a4 4 0 0 1 8 0v4" />
    </svg>
    <template v-if="!iconOnly">ロック中</template>
  </span>
</template>

<style scoped>
.event-lock-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  vertical-align: middle;
  white-space: nowrap;
}
</style>
