<script setup lang="ts">
import { computed, ref } from "vue";
import { currentFestivalId, stockFilterCategory } from "@/state";
import { useResource, listOf } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { api } from "@/lib/api";
import { salesTotal } from "@/lib/sales";
import { categoryTags } from "@/lib/categories";
import CategoryTags from "@/components/CategoryTags.vue";
import type { Stock } from "@/types/stock";
import type { SaleRecord } from "@/types/saleRecord";
import ResourceState from "@/components/ResourceState.vue";
const {
  data: records,
  loading,
  error,
  reload,
} = useResource(
  () =>
    currentFestivalId.value
      ? `/sales?festival_id=${encodeURIComponent(currentFestivalId.value)}`
      : null,
  listOf<SaleRecord>("sales"),
);
const {
  data: stocks,
  loading: stocksLoading,
  error: stocksError,
  reload: reloadStocks,
} = useResource(
  () =>
    currentFestivalId.value
      ? `/festivals/${currentFestivalId.value}/stocks`
      : null,
  listOf<Stock>("stocks"),
);
const search = ref("");
const edit = ref(false);
const { pending, run } = useMutation();
const stockMap = computed(
  () => new Map((stocks.value || []).map((s) => [s.id, s])),
);
const categories = computed(() =>
  [...new Set((stocks.value || []).flatMap((s) => categoryTags(s.item.category)))].sort(),
);
const filtered = computed(() =>
  (records.value || [])
    .filter((r) => {
      const stock = stockMap.value.get(r.stock_id);
      return (
        (!stockFilterCategory.value ||
          categoryTags(stock?.item.category || "").includes(stockFilterCategory.value)) &&
        (!search.value.trim() || stock?.item.name.includes(search.value.trim()))
      );
    })
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    ),
);
const total = computed(() => salesTotal(filtered.value, stocks.value || []));
const quantity = computed(() =>
  filtered.value.reduce((sum, r) => sum + r.quantity, 0),
);
const unknown = computed(() =>
  filtered.value.some((r) => !stockMap.value.has(r.stock_id)),
);
async function refresh() {
  await Promise.all([reload(), reloadStocks()]);
}
async function remove(record: SaleRecord) {
  if (
    !window.confirm(
      `「${stockMap.value.get(record.stock_id)?.item.name || "商品不明"}」${record.quantity}点の売上記録を削除しますか？この操作は取り消せません。`,
    )
  )
    return;
  await run(async () => {
    await api(`/sales/${record.id}`, { method: "DELETE" });
    records.value = (records.value || []).filter((r) => r.id !== record.id);
  }, "売上記録を削除しました。");
}
function amount(record: SaleRecord) {
  const stock = stockMap.value.get(record.stock_id);
  return stock
    ? `${(stock.price * record.quantity).toLocaleString()} 円`
    : "価格不明";
}
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">SALES HISTORY</p>
      <h1>売上履歴</h1>
      <p class="muted">選択したイベントの売上を、新しい順に表示します。</p>
    </div>
    <button class="button secondary" :disabled="pending" :aria-pressed="edit" @click="edit = !edit">
      {{ edit ? "編集を終了" : "記録を編集" }}
    </button>
  </div>
  <p v-if="!currentFestivalId" class="state">
    上の「対象イベント」からイベントを選択してください。
  </p>
  <template v-else>
    <div v-if="records && stocks && !stocksError" class="summary-grid">
      <div class="panel summary">
        <span>表示中の売上金額</span><strong>{{ total.toLocaleString() }}<small> 円</small></strong>
      </div>
      <div class="panel summary">
        <span>販売数</span><strong>{{ quantity.toLocaleString() }}<small> 点</small></strong>
      </div>
      <div class="panel summary">
        <span>売上記録</span><strong>{{ filtered.length.toLocaleString() }}<small> 件</small></strong>
      </div>
    </div>
    <div class="toolbar">
      <label class="search-field">商品検索<input v-model="search" type="search" placeholder="商品名で検索" /></label><label
        class="filter-field">カテゴリ<select v-model="stockFilterCategory">
          <option value="">すべてのカテゴリ</option>
          <option v-if="
            stockFilterCategory && !categories.includes(stockFilterCategory)
          " :value="stockFilterCategory">
            {{ stockFilterCategory }}
          </option>
          <option v-for="value in categories" :key="value">{{ value }}</option>
        </select></label><button class="button secondary" :disabled="loading || stocksLoading || pending"
        @click="refresh">
        更新
      </button>
    </div>
    <ResourceState :loading="loading || stocksLoading" :error="error || stocksError" :empty="!filtered.length"
      empty-text="条件に一致する売上記録がありません。" @retry="refresh" />
    <p v-if="unknown" class="state">
      販売商品が見つからない記録は「価格不明」と表示し、売上金額の合計から除いています。
    </p>
    <div v-if="filtered.length" class="panel table-scroll">
      <table>
        <thead>
          <tr>
            <th>日時</th>
            <th>商品 / カテゴリ</th>
            <th class="numeric">単価</th>
            <th class="numeric">数量</th>
            <th class="numeric">金額</th>
            <th v-if="edit">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in filtered" :key="record.id">
            <td class="nowrap">
              {{ new Date(record.created_at).toLocaleString("ja-JP") }}
            </td>
            <td>
              <strong>{{
                stockMap.get(record.stock_id)?.item.name || "商品不明"
                }}</strong>
              <p class="small muted">
                <CategoryTags v-if="stockMap.has(record.stock_id)" :category="stockMap.get(record.stock_id)?.item.category || ''" />
                <template v-else>—</template>
              </p>
            </td>
            <td class="numeric">
              {{
                stockMap.get(record.stock_id)?.price.toLocaleString() ?? "不明"
              }}
            </td>
            <td class="numeric">{{ record.quantity }}</td>
            <td class="numeric">{{ amount(record) }}</td>
            <td v-if="edit">
              <button class="button danger compact" :disabled="pending"
                :aria-label="`${new Date(record.created_at).toLocaleString('ja-JP')}の${stockMap.get(record.stock_id)?.item.name || '商品不明'}の記録を削除`"
                @click="remove(record)">
                削除
              </button>
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <th colspan="3">合計</th>
            <td class="numeric">{{ quantity }}</td>
            <td class="numeric">{{ total.toLocaleString() }} 円</td>
            <td v-if="edit" />
          </tr>
        </tfoot>
      </table>
    </div>
  </template>
</template>
