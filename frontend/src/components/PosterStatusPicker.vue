<script setup lang="ts">
import { api, jsonBody } from "@/lib/api";
import { useMutation } from "@/composables/useMutation";
import { PosterStatusLabels, type PosterStatus } from "@/types/poster";
const props = defineProps<{
  posterId: string;
  status: PosterStatus;
  name: string;
  disabled?: boolean;
}>();
const emit = defineEmits<{ updated: [status: PosterStatus] }>();
const { pending, run } = useMutation();
async function change(event: Event) {
  const target = event.target as HTMLSelectElement;
  const status = target.value as PosterStatus;
  const saved = await run(async () => {
    await api(`/posters/${props.posterId}/status`, {
      method: "PATCH",
      ...jsonBody({ status }),
    });
    emit("updated", status);
  }, "回収状況を更新しました。");
  if (!saved) target.value = props.status;
}
</script>
<template>
  <select :value="status" :disabled="pending || disabled" :aria-label="`${name}の回収状況`"
    :class="['status-select', status]" @change="change">
    <option v-for="(label, value) in PosterStatusLabels" :key="value" :value="value">
      {{ label }}
    </option>
  </select>
</template>
