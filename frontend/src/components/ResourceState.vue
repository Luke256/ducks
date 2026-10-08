<script setup lang="ts">
defineProps<{
  loading?: boolean;
  error?: string;
  empty?: boolean;
  emptyText?: string;
}>();
defineEmits<{ retry: [] }>();
</script>
<template>
  <p v-if="loading" class="state" role="status">読み込み中…</p>
  <div v-else-if="error" class="state error" role="alert">
    <p>{{ error }}</p>
    <button class="button secondary" type="button" @click="$emit('retry')">
      再読み込み
    </button>
  </div>
  <p v-else-if="empty" class="state">
    {{ emptyText || "まだ登録されていません。" }}
  </p>
</template>
