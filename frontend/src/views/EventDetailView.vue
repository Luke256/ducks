<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useResource } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { currentFestivalId, loadFestivals } from "@/state";
import { api, jsonBody } from "@/lib/api";
import type { Festival } from "@/types/festival";
import ResourceState from "@/components/ResourceState.vue";
import EventLockBadge from "@/components/EventLockBadge.vue";
import { festivalDate, formatFestivalDate } from "@/utils/festivalDate";
const route = useRoute();
const router = useRouter();
const {
  data: event,
  loading,
  error,
  reload,
} = useResource<Festival>(() => `/festivals/${route.params.eventId}`);
const edit = ref(false);
const form = reactive({ name: "", description: "", date: "" });
const { pending, error: saveError, run } = useMutation();
watch(event, (value) => {
  if (!value) return;
  document.title = `${value.name} | Ducks`;
  if (value.locked) {
    Object.assign(form, value, { date: festivalDate(value.date) });
    edit.value = false;
  }
});
function startEdit() {
  if (event.value && !event.value.locked && !pending.value && !loading.value) {
    Object.assign(form, event.value, { date: festivalDate(event.value.date) });
    edit.value = true;
  }
}
async function save() {
  await run(async () => {
    if (!form.name.trim()) throw new Error("イベント名を入力してください。");
    if (!form.date) throw new Error("開催日を入力してください。");
    await api(`/festivals/${route.params.eventId}`, {
      method: "PUT",
      ...jsonBody({ name: form.name.trim(), description: form.description, date: form.date }),
    });
    edit.value = false;
    await Promise.all([reload(), loadFestivals()]);
  }, "イベントを更新しました。");
}
async function toggleLock() {
  if (!event.value || pending.value || loading.value) return;
  const locked = !event.value.locked;
  if (locked && !window.confirm(
    `「${event.value.name}」をロックしますか？\nロック中はイベントと関連データを変更・削除できません。${edit.value ? "\n未保存の編集内容は破棄されます。" : ""}`,
  )) return;
  await run(async () => {
    event.value = await api<Festival>(`/festivals/${route.params.eventId}/lock`, {
      method: "PATCH",
      ...jsonBody({ locked }),
    });
    await loadFestivals();
  }, locked ? "イベントをロックしました。" : "イベントのロックを解除しました。");
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
        <h1>{{ event.name }} <EventLockBadge v-if="event.locked" /></h1>
        <p class="muted"><time v-if="festivalDate(event.date)" class="event-date" :datetime="festivalDate(event.date)">{{ formatFestivalDate(event.date) }}</time><span v-else class="event-date">未設定</span></p>
        <p class="description muted">
          {{ event.description || "概要は登録されていません。" }}
        </p>
      </div>
      <div class="actions">
        <button class="button secondary" :disabled="pending || loading || event.locked" @click="edit ? (edit = false) : startEdit()">
          {{ edit ? "編集をキャンセル" : "イベントを編集" }}
        </button>
        <button class="button secondary" type="button" :disabled="pending || loading" @click="toggleLock">
          {{ pending ? "処理中…" : event.locked ? "ロックを解除" : "イベントをロック" }}
        </button>
      </div>
    </div>
    <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
    <div class="quick-links">
      <RouterLink to="/poster" class="panel" @click="selectEvent"><span class="eyebrow">POSTERS</span>
        <h2>ポスター管理 →</h2>
        <p class="muted">設置場所と回収状況を確認</p>
      </RouterLink>
      <RouterLink to="/sales/cashier" class="panel" @click="selectEvent"><span class="eyebrow">CASHIER</span>
        <h2>レジを開く →</h2>
        <p class="muted">商品を選んで会計</p>
      </RouterLink>
      <RouterLink to="/sales/orders" class="panel" @click="selectEvent"><span class="eyebrow">SALES</span>
        <h2>売上履歴 →</h2>
        <p class="muted">販売数と売上を確認</p>
      </RouterLink>
      <RouterLink to="/visitors" class="panel" @click="selectEvent"><span class="eyebrow">VISITORS</span>
        <h2>来場者数 →</h2>
        <p class="muted">人数のカウントと履歴を確認</p>
      </RouterLink>
    </div>
    <form v-if="edit" class="panel form-panel" @submit.prevent="save">
      <h2>イベントを編集</h2>
      <fieldset :disabled="pending || loading || event.locked">
        <label class="field">イベント名<input v-model="form.name" required /></label>
        <label class="field">開催日<input v-model="form.date" type="date" required min="1000-01-01" max="9999-12-31" /></label>
        <label class="field">概要<textarea v-model="form.description" rows="4" />
        </label>
        <div class="actions">
          <button class="button" type="submit">
            {{ pending ? "保存中…" : "保存する" }}</button><button class="button danger" type="button" @click="remove">
            イベントを削除
          </button>
        </div>
      </fieldset>
    </form>
  </template>
</template>
