<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { currentFestivalId } from "@/state";
import { useResource, listOf } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { api, jsonBody } from "@/lib/api";
import type { VisitorCount } from "@/types/visitorCount";
import ResourceState from "@/components/ResourceState.vue";
import VisitorChart from "@/components/VisitorChart.vue";
import { dailyVisitorCounts } from "@/lib/visitors";

const {
  data: counts,
  loading,
  error,
  reload,
} = useResource(
  () =>
    currentFestivalId.value ? `/visitors/${currentFestivalId.value}` : null,
  listOf<VisitorCount>("counts"),
);
const { pending, error: saveError, run } = useMutation();
const amount = ref<number | string>(1);
const total = computed(() =>
  (counts.value || []).reduce((sum, bucket) => sum + bucket.count, 0),
);
const history = computed(() =>
  [...(counts.value || [])].sort(
    (a, b) => Date.parse(b.bucket_start) - Date.parse(a.bucket_start),
  ),
);
const days = computed(() => dailyVisitorCounts(counts.value || []));
const pageSize = 12;
const page = ref(1);
const pageCount = computed(() => Math.max(1, Math.ceil(history.value.length / pageSize)));
const pageHistory = computed(() => history.value.slice((page.value - 1) * pageSize, page.value * pageSize));
watch(pageCount, (pages) => {
  page.value = Math.min(page.value, pages);
});
watch(currentFestivalId, () => {
  page.value = 1;
  amount.value = 1;
  saveError.value = "";
});
async function add(value: number, correction = false) {
  await run(async () => {
    if (!currentFestivalId.value)
      throw new Error("イベントを選択してください。");
    if (
      !Number.isSafeInteger(value) ||
      (value < 1 && !(correction && value === -1))
    )
      throw new Error("人数は1以上の整数で入力してください。");
    await api(`/visitors/${currentFestivalId.value}`, {
      method: "POST",
      ...jsonBody({ amount: value }),
    });
    await reload();
  }, value === -1 ? "現在の10分間の人数を1人減らしました。" : `${value}人を追加しました。`);
}
function bucketLabel(start: string) {
  return new Date(start).toLocaleString("ja-JP", {
    timeZone: "Asia/Tokyo",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">VISITOR COUNTER</p>
      <h1>来場者数</h1>
      <p class="muted">選択したイベントの来場者数をカウントします。</p>
    </div>
    <button v-if="currentFestivalId" class="button secondary" :disabled="loading || pending"
      @click="reload()">更新</button>
  </div>
  <p v-if="!currentFestivalId" class="state">
    上の「対象イベント」からイベントを選択してください。
  </p>
  <template v-else>
    <ResourceState :loading="loading" :error="error" @retry="reload()" />
    <div v-if="counts && !error" class="panel summary visitor-total" aria-live="polite">
      <span>累計来場者数</span>
      <strong>{{ total.toLocaleString() }}<small> 人</small></strong>
    </div>
    <form class="panel form-panel" @submit.prevent="add(Number(amount))">
      <h2>来場者を追加</h2>
      <fieldset :disabled="pending || loading || !counts || !!error">
        <div class="actions">
          <button class="button visitor-increment" type="button" @click="add(1)">＋1人</button>
          <button class="button secondary" type="button" @click="add(-1, true)">−1人（訂正）</button>
        </div>
        <div class="actions visitor-quick-add">
          <button v-for="count in [2, 3, 4, 5]" :key="count" class="button compact" type="button"
            @click="add(count)">＋{{ count }}人</button>
        </div>
        <p class="small muted">訂正は現在の10分間の人数から減らします。過去の時間帯は変更できません。</p>
        <label class="field">まとめて追加
          <input v-model="amount" type="number" min="1" step="1" required />
        </label>
        <button class="button" type="submit">{{ pending ? "保存中…" : "人数を追加" }}</button>
      </fieldset>
      <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
      <p v-if="saveError" class="small muted">通信が途切れた場合は、更新して人数を確認してから再送信してください。</p>
    </form>
    <h2>10分ごとの来場者数</h2>
    <p v-if="counts && !error && !history.length" class="state">まだ来場者数が記録されていません。</p>
    <template v-if="history.length && !error">
      <VisitorChart v-if="days.length" :days="days" />
      <p v-if="!days.length" class="state">1人以上の記録がないため、グラフは表示されません。</p>
      <h2>詳細</h2>
      <div class="panel table-scroll">
        <table>
          <caption class="sr-only">10分ごとの来場者数の詳細（新しい順）</caption>
          <thead>
            <tr>
              <th>開始日時</th>
              <th class="numeric">来場者数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="bucket in pageHistory" :key="bucket.bucket_start">
              <td class="nowrap">{{ bucketLabel(bucket.bucket_start) }}</td>
              <td class="numeric">{{ bucket.count.toLocaleString() }} 人</td>
            </tr>
          </tbody>
        </table>
      </div>
      <nav class="actions visitor-pagination" aria-label="来場者数の詳細ページ">
        <button class="button secondary" :disabled="page === 1" @click="page--">前へ</button>
        <span class="small" aria-live="polite">{{ page }} / {{ pageCount }} ページ（{{ (page - 1) * pageSize + 1 }}〜{{ Math.min(page * pageSize, history.length) }} / {{ history.length }}件）</span>
        <button class="button secondary" :disabled="page === pageCount" @click="page++">次へ</button>
      </nav>
    </template>
  </template>
</template>
