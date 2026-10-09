<script setup lang="ts">
import { onBeforeUnmount, shallowRef, watch } from "vue";
defineProps<{ label: string; required?: boolean }>();
const files = defineModel<File[]>({ required: true });
const previews = shallowRef<{ file: File; url: string }[]>([]);
function clearPreviews() {
  for (const preview of previews.value) URL.revokeObjectURL(preview.url);
}
watch(files, (value) => {
  clearPreviews();
  previews.value = value.map((file) => ({ file, url: URL.createObjectURL(file) }));
}, { immediate: true });
function choose(event: Event) {
  const input = event.target as HTMLInputElement;
  const selected = Array.from(input.files || []);
  if (selected.length) files.value = [...files.value, ...selected];
  input.value = "";
}
function remove(index: number) {
  files.value = files.value.filter((_, i) => i !== index);
}
onBeforeUnmount(clearPreviews);
</script>
<template>
  <label class="field">{{ label }}<input type="file" accept="image/*" multiple
    :required="required && !files.length" @change="choose" /></label>
  <div v-if="previews.length" class="image-selection">
    <figure v-for="(preview, index) in previews" :key="index">
      <img :src="preview.url" :alt="`選択した写真${index + 1}のプレビュー`" class="image-preview" />
      <figcaption>
        <p class="small">{{ preview.file.name }}</p>
        <button type="button" class="button secondary" :aria-label="`選択した写真${index + 1}を取り消す`"
          @click="remove(index)">選択を取り消す</button>
      </figcaption>
    </figure>
  </div>
</template>
