<script setup lang="ts">
import { categoryTags } from "@/lib/categories";
defineProps<{ category: string }>();
</script>
<template>
  <span class="category-tags">
    <span v-for="(tag, index) in categoryTags(category)" :key="tag" class="badge"
      :class="{ primary: index === 0 }" :title="index === 0 ? '主タグ' : undefined">{{ tag }}</span>
  </span>
</template>
