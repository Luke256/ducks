<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useResource } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { currentFestivalId, loadFestivals } from "@/state";
import { api, jsonBody } from "@/lib/api";
import type { Festival } from "@/types/festival";
import ResourceState from "@/components/ResourceState.vue";
const route = useRoute();
const router = useRouter();
const {
  data: event,
  loading,
  error,
  reload,
} = useResource<Festival>(() => `/festivals/${route.params.eventId}`);
const edit = ref(false);
const form = reactive({ name: "", description: "" });
const { pending, error: saveError, run } = useMutation();
watch(event, (value) => {
  if (value) document.title = `${value.name} | Ducks`;
});
function startEdit() {
  if (event.value) {
    Object.assign(form, event.value);
    edit.value = true;
  }
}
async function save() {
  await run(async () => {
    if (!form.name.trim()) throw new Error("イベント名を入力してください。");
    await api(`/festivals/${route.params.eventId}`, {
      method: "PUT",
      ...jsonBody({ name: form.name.trim(), description: form.description }),
    });
    edit.value = false;
    await Promise.all([reload(), loadFestivals()]);
  }, "イベントを更新しました。");
}
async function remove() {
  if (
    !window.confirm(
      `「${event.value?.name}」を削除しますか？関連するポスター・販売商品・売上記録への影響を確認してください。この操作は取り消せません。`,
    )
  )
    return;
  await run(async () => {
    await api(`/festivals/${route.params.eventId}`, { method: "DELETE" });
    await loadFestivals();
    await router.push("/event");
  }, "イベントを削除しました。");
}
function selectEvent() {
  if (event.value) currentFestivalId.value = event.value.id;
}
</script>
<template>
  <RouterLink to="/event" class="back-link">← イベント一覧</RouterLink>
  <ResourceState :loading="loading" :error="error" @retry="reload()" />
  <template v-if="event">
    <div class="page-heading">
      <div>
        <p class="eyebrow">EVENT DETAIL</p>
        <h1>{{ event.name }}</h1>
        <p class="description muted">
          {{ event.description || "概要は登録されていません。" }}
        </p>
      </div>
      <button
        class="button secondary"
        :disabled="pending"
        @click="edit ? (edit = false) : startEdit()"
      >
        {{ edit ? "編集をキャンセル" : "イベントを編集" }}
      </button>
    </div>
    <div class="quick-links">
      <RouterLink to="/poster" class="panel" @click="selectEvent"
        ><span class="eyebrow">POSTERS</span>
        <h2>ポスター管理 →</h2>
        <p class="muted">設置場所と回収状況を確認</p></RouterLink
      ><RouterLink to="/sales/cashier" class="panel" @click="selectEvent"
        ><span class="eyebrow">CASHIER</span>
        <h2>レジを開く →</h2>
        <p class="muted">商品を選んで会計</p></RouterLink
      ><RouterLink to="/sales/orders" class="panel" @click="selectEvent"
        ><span class="eyebrow">SALES</span>
        <h2>売上履歴 →</h2>
        <p class="muted">販売数と売上を確認</p></RouterLink
      ><RouterLink to="/visitors" class="panel" @click="selectEvent"
        ><span class="eyebrow">VISITORS</span>
        <h2>来場者数 →</h2>
        <p class="muted">人数のカウントと履歴を確認</p></RouterLink
      >
    </div>
    <form v-if="edit" class="panel form-panel" @submit.prevent="save">
      <h2>イベントを編集</h2>
      <fieldset :disabled="pending">
        <label class="field"
          >イベント名<input v-model="form.name" required /></label
        ><label class="field"
          >概要<textarea v-model="form.description" rows="4" />
        </label>
        <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
        <div class="actions">
          <button class="button" type="submit">
            {{ pending ? "保存中…" : "保存する" }}</button
          ><button class="button danger" type="button" @click="remove">
            イベントを削除
          </button>
        </div>
      </fieldset>
    </form>
  </template>
</template>
