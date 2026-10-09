<script setup lang="ts">
import { onBeforeUnmount, ref, shallowRef } from "vue";
import { imageUrl } from "@/lib/api";
defineProps<{ label: string; required?: boolean; existing?: string }>();
const emit = defineEmits<{ change: [file: File | null] }>();
const file = shallowRef<File | null>(null);
const preview = ref("");
function choose(event: Event) {
  const input = event.target as HTMLInputElement;
  const selected = input.files?.[0];
  if (!selected) return;
  if (preview.value) URL.revokeObjectURL(preview.value);
  file.value = selected;
  preview.value = URL.createObjectURL(selected);
  input.value = "";
  emit("change", selected);
}
function remove() {
  if (!file.value || !window.confirm(`「${file.value.name}」の選択を取り消しますか？`)) return;
  URL.revokeObjectURL(preview.value);
  preview.value = "";
  file.value = null;
  emit("change", null);
}
onBeforeUnmount(() => {
  if (preview.value) URL.revokeObjectURL(preview.value);
});
</script>
<template>
  <div v-if="preview || existing" class="image-selection">
    <figure>
      <img :src="preview || imageUrl(existing || '')" :alt="file ? '追加する商品画像' : '登録済みの商品画像'" class="image-preview" />
      <figcaption>
        <p class="small">{{ file ? file.name : "現在の画像" }}</p>
        <button v-if="file" type="button" class="button secondary" aria-label="追加する商品画像の選択を取り消す"
          @click="remove">取り消し</button>
      </figcaption>
    </figure>
  </div>
  <label class="button secondary image-add-button"><span aria-hidden="true">＋</span>{{ label }}
    <input type="file" accept="image/*" :aria-label="label" :required="required && !file" @change="choose" />
  </label>
</template>
