<script setup lang="ts">
import { onBeforeUnmount, ref } from "vue";
defineProps<{ label: string; required?: boolean }>();
const emit = defineEmits<{ change: [file: File | null] }>();
const preview = ref("");
function choose(event: Event) {
  if (preview.value) URL.revokeObjectURL(preview.value);
  const file = (event.target as HTMLInputElement).files?.[0] || null;
  preview.value = file ? URL.createObjectURL(file) : "";
  emit("change", file);
}
onBeforeUnmount(() => {
  if (preview.value) URL.revokeObjectURL(preview.value);
});
</script>
<template>
  <label class="field"
    >{{ label
    }}<input type="file" accept="image/*" :required="required" @change="choose"
  /></label>
  <img
    v-if="preview"
    :src="preview"
    alt="選択した画像のプレビュー"
    class="image-preview"
  />
</template>
