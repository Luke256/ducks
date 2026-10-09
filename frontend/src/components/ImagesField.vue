<script setup lang="ts">
import { computed, onBeforeUnmount, shallowRef, watch } from "vue";
import { imageUrl } from "@/lib/api";
import type { PosterImage } from "@/types/poster";
const props = withDefaults(defineProps<{ label: string; required?: boolean; existing?: PosterImage[] }>(), {
  existing: () => [],
});
const files = defineModel<File[]>({ required: true });
const deleteImageIds = defineModel<string[]>("deleteImageIds", { default: () => [] });
const previews = shallowRef<{ file: File; url: string }[]>([]);
const photos = computed(() => [
  ...props.existing.map((image, index) => ({ id: image.id, url: imageUrl(image.url), name: `写真${index + 1}`, index })),
  ...previews.value.map((preview, index) => ({ id: "", url: preview.url, name: preview.file.name, index })),
]);
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
function cancelExisting(id: string) {
  deleteImageIds.value = deleteImageIds.value.includes(id)
    ? deleteImageIds.value.filter((imageId) => imageId !== id)
    : [...deleteImageIds.value, id];
}
function remove(index: number) {
  if (!window.confirm(`「${files.value[index].name}」を削除しますか？`)) return;
  files.value = files.value.filter((_, i) => i !== index);
}
onBeforeUnmount(clearPreviews);
</script>
<template>
  <div v-if="photos.length" class="image-selection">
    <figure v-for="photo in photos" :key="photo.id ? `existing-${photo.id}` : `new-${photo.index}`"
      :class="{ 'marked-for-deletion': photo.id && deleteImageIds.includes(photo.id) }">
      <img :src="photo.url" :alt="`${photo.id ? '登録済み' : '追加する'}の写真${photo.index + 1}`" class="image-preview" />
      <figcaption>
        <p class="small">{{ photo.name }}</p>
        <button v-if="photo.id" type="button" class="button secondary"
          :aria-label="`登録済みの写真${photo.index + 1}${deleteImageIds.includes(photo.id) ? 'の削除を元に戻す' : 'を削除する'}`"
          @click="cancelExisting(photo.id)">{{ deleteImageIds.includes(photo.id) ? "元に戻す" : "削除" }}</button>
        <button v-else type="button" class="button secondary" :aria-label="`追加する写真${photo.index + 1}を削除する`"
          @click="remove(photo.index)">削除</button>
      </figcaption>
    </figure>
  </div>
  <label class="button secondary image-add-button"><span aria-hidden="true">＋</span>{{ label }}
    <input type="file" accept="image/*" multiple :aria-label="label"
      :required="required && !files.length" @change="choose" />
  </label>
</template>
