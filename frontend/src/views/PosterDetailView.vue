<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useResource } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { api, jsonBody, imageUrl } from "@/lib/api";
import { currentFestivalId } from "@/state";
import type { Poster, PosterStatus } from "@/types/poster";
import PosterStatusPicker from "@/components/PosterStatusPicker.vue";
import ResourceState from "@/components/ResourceState.vue";
const route = useRoute();
const router = useRouter();
const {
  data: poster,
  error,
  loading,
  reload,
} = useResource<Poster>(() => `/posters/${route.params.posterId}`);
const edit = ref(false);
const form = reactive({ name: "", description: "" });
const { pending, error: saveError, run } = useMutation();
watch(poster, (value) => {
  if (value) currentFestivalId.value = value.festival.id;
});
function startEdit() {
  if (poster.value) {
    Object.assign(form, poster.value);
    edit.value = true;
  }
}
function updateStatus(status: PosterStatus) {
  if (poster.value) poster.value = { ...poster.value, status };
}
async function save() {
  await run(async () => {
    if (!form.name.trim() || !form.description.trim())
      throw new Error("ポスター名と設置場所を入力してください。");
    await api(`/posters/${route.params.posterId}`, {
      method: "PUT",
      ...jsonBody({
        name: form.name.trim(),
        description: form.description.trim(),
      }),
    });
    edit.value = false;
    await reload();
  }, "ポスターを更新しました。");
}
async function remove() {
  if (
    !window.confirm(
      `「${poster.value?.name}」を削除しますか？この操作は取り消せません。`,
    )
  )
    return;
  await run(async () => {
    await api(`/posters/${route.params.posterId}`, { method: "DELETE" });
    await router.push("/poster");
  }, "ポスターを削除しました。");
}
</script>
<template>
  <RouterLink to="/poster" class="back-link">← ポスター一覧</RouterLink>
  <ResourceState :loading="loading" :error="error" @retry="reload()" />
  <template v-if="poster">
    <div class="page-heading">
      <div>
        <p class="eyebrow">POSTER DETAIL / {{ poster.festival.name }}</p>
        <h1>{{ poster.name }}</h1>
      </div>
      <button class="button secondary" :disabled="pending" @click="edit ? (edit = false) : startEdit()">
        {{ edit ? "編集をキャンセル" : "編集する" }}
      </button>
    </div>
    <div class="detail-grid">
      <section class="panel">
        <h2>設置場所</h2>
        <p class="description">{{ poster.description }}</p>
        <div v-if="poster.image_url.length" class="poster-images">
          <a v-for="(url, index) in poster.image_url" :key="url" :href="imageUrl(url)" target="_blank" rel="noopener"
            :aria-label="`設置場所の写真${index + 1}を拡大`"><img :src="imageUrl(url)"
              :alt="`${poster.name}の設置場所（写真${index + 1}）`" class="detail-image" loading="lazy" /></a>
        </div>
      </section>
      <section class="panel">
        <h2>回収状況</h2>
        <PosterStatusPicker :poster-id="poster.id" :name="poster.name" :status="poster.status" :disabled="pending"
          @updated="updateStatus" />
        <p class="muted small">変更するとすぐに保存されます。</p>
      </section>
    </div>
    <form v-if="edit" class="panel form-panel" @submit.prevent="save">
      <h2>登録情報を編集</h2>
      <fieldset :disabled="pending">
        <label class="field">ポスター名<input v-model="form.name" required maxlength="64" /></label><label
          class="field">設置場所<textarea v-model="form.description" required maxlength="1024" rows="4" />
        </label>
        <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
        <div class="actions">
          <button class="button" type="submit">
            {{ pending ? "保存中…" : "保存する" }}</button><button class="button danger" type="button" @click="remove">
            ポスターを削除
          </button>
        </div>
      </fieldset>
    </form>
  </template>
</template>
