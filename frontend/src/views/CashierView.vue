<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from "vue";
import { onBeforeRouteLeave } from "vue-router";
import { currentFestivalId, notify } from "@/state";
import { useResource, listOf } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { api, jsonBody, imageUrl } from "@/lib/api";
import { saleItems } from "@/lib/sales";
import { categoryTags } from "@/lib/categories";
import type { Stock } from "@/types/stock";
import ResourceState from "@/components/ResourceState.vue";
const {
  data: stocks,
  error,
  loading,
  reload,
} = useResource(
  () =>
    currentFestivalId.value
      ? `/festivals/${currentFestivalId.value}/stocks`
      : null,
  listOf<Stock>("stocks"),
);
const cart = ref<Record<string, number>>({});
const search = ref("");
const category = ref("");
const received = ref<number | "">("");
const { pending, error: saleError, run } = useMutation();
const categories = computed(() =>
  [...new Set((stocks.value || []).flatMap((s) => categoryTags(s.item.category)))].sort(),
);
const groups = computed(() => {
  const result: Record<string, Stock[]> = Object.create(null);
  for (const stock of stocks.value || []) {
    const tags = categoryTags(stock.item.category);
    if (
      (category.value && !tags.includes(category.value)) ||
      !`${stock.item.name} ${stock.description}`.includes(search.value.trim())
    )
      continue;
    (result[tags[0] || "未分類"] ||= []).push(stock);
  }
  return Object.entries(result)
    .sort(([a], [b]) => a.localeCompare(b, "ja"))
    .map(
      ([name, products]) =>
        [
          name,
          products.sort((a, b) => a.item.name.localeCompare(b.item.name, "ja")),
        ] as const,
    );
});
const lines = computed(() =>
  Object.entries(cart.value).map(([id, quantity]) => ({
    id,
    quantity,
    stock: stocks.value?.find((s) => s.id === id),
  })),
);
const total = computed(() =>
  lines.value.reduce((sum, l) => sum + (l.stock?.price || 0) * l.quantity, 0),
);
const count = computed(() =>
  lines.value.reduce((sum, l) => sum + l.quantity, 0),
);
const validCash = computed(
  () =>
    received.value === "" ||
    (Number.isSafeInteger(received.value) && received.value >= total.value),
);
const ready = computed(
  () =>
    !!currentFestivalId.value &&
    lines.value.length > 0 &&
    lines.value.every((l) => !!l.stock) &&
    validCash.value &&
    !loading.value &&
    !error.value &&
    !pending.value,
);
function adjust(id: string, delta: number) {
  if (pending.value || loading.value || error.value) return;
  const quantity = (cart.value[id] || 0) + delta;
  if (quantity > 0) cart.value = { ...cart.value, [id]: quantity };
  else {
    const next = { ...cart.value };
    delete next[id];
    cart.value = next;
  }
}
function clear() {
  cart.value = {};
  received.value = "";
}
watch(
  currentFestivalId,
  () => {
    if (lines.value.length)
      notify("イベントを切り替えたため、会計中の商品をリセットしました。");
    clear();
    category.value = "";
  },
  { flush: "sync" },
);
async function checkout() {
  if (!ready.value) return;
  await run(async () => {
    const items = saleItems(
      cart.value,
      stocks.value || [],
      currentFestivalId.value,
    );
    await api("/sales", { method: "POST", ...jsonBody({ items }) });
    clear();
  }, "会計が完了しました");
}
onBeforeRouteLeave(() =>
  pending.value
    ? false
    : !lines.value.length ||
    window.confirm("会計中の商品があります。選択を破棄して移動しますか？"),
);
function beforeUnload(event: BeforeUnloadEvent) {
  if (lines.value.length || pending.value) {
    event.preventDefault();
    event.returnValue = "";
  }
}
window.addEventListener("beforeunload", beforeUnload);
onBeforeUnmount(() => {
  window.removeEventListener("beforeunload", beforeUnload);
});
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">CASHIER</p>
      <h1>レジ</h1>
      <p class="muted">商品をタップして追加。右の明細で数量を調整できます。</p>
    </div>
    <RouterLink to="/sales/orders" class="button secondary">売上履歴 →</RouterLink>
  </div>
  <p v-if="!currentFestivalId" class="state">
    上の「対象イベント」からイベントを選択してください。
  </p>
  <div v-else class="cashier-layout">
    <section>
      <div class="toolbar">
        <label class="search-field">商品検索<input v-model="search" type="search" placeholder="名前・説明で検索" /></label><label
          class="filter-field">カテゴリ<select v-model="category">
            <option value="">すべてのカテゴリ</option>
            <option v-for="value in categories" :key="value">
              {{ value }}
            </option>
          </select></label><button class="button secondary" :disabled="loading || pending" @click="reload()">
          更新
        </button>
      </div>
      <ResourceState :loading="loading" :error="error" :empty="!groups.length" empty-text="販売商品がありません。販売商品タブから登録してください。"
        @retry="reload()" />
      <section v-for="[name, products] in groups" :key="name" class="category-section">
        <h2>
          {{ name }}<span class="muted small"> {{ products.length }} 商品</span>
        </h2>
        <div class="cashier-products">
          <button v-for="stock in products" :key="stock.id" type="button"
            :class="['panel cashier-product', { selected: cart[stock.id] }]" :disabled="pending || loading || !!error"
            :aria-label="`${stock.item.name}を1点追加、${stock.price}円`" @click="adjust(stock.id, 1)">
            <div class="product-photo">
              <img v-if="stock.item.image_url" :src="imageUrl(stock.item.image_url)" alt="" loading="lazy" /><span
                v-if="cart[stock.id]" class="quantity-badge">{{
                  cart[stock.id]
                }}</span>
            </div>
            <strong>{{ stock.item.name }}</strong><span class="product-price">{{ stock.price.toLocaleString() }}
              円</span><span v-if="stock.description" class="muted small">{{
                stock.description
              }}</span>
          </button>
        </div>
      </section>
      <RouterLink v-if="!loading && !error && !stocks?.length" to="/sales/stocks/new" class="button secondary">販売商品を登録
      </RouterLink>
    </section>
    <aside class="panel receipt" aria-label="会計明細">
      <div class="receipt-heading">
        <h2>会計明細</h2>
        <span class="badge">{{ count }} 点</span>
      </div>
      <p v-if="!lines.length" class="receipt-empty muted">
        商品を選ぶと、ここに明細が表示されます。
      </p>
      <div v-for="line in lines" :key="line.id" class="receipt-line">
        <strong>{{
          line.stock?.item.name || "販売商品を再確認してください"
          }}</strong>
        <div class="line-bottom">
          <div class="quantity-control">
            <button type="button" :aria-label="`${line.stock?.item.name}の数量を減らす`"
              :disabled="pending || loading || !!error" @click="adjust(line.id, -1)">
              −</button><span>{{ line.quantity }}</span><button type="button"
              :aria-label="`${line.stock?.item.name}の数量を増やす`" :disabled="pending || loading || !!error"
              @click="adjust(line.id, 1)">
              ＋
            </button>
          </div>
          <span>{{
            ((line.stock?.price || 0) * line.quantity).toLocaleString()
          }}
            円</span>
        </div>
      </div>
      <div class="receipt-total">
        <span>合計</span><strong>{{ total.toLocaleString() }}<small> 円</small></strong>
      </div>
      <label class="field">お預かり金額（任意）<input v-model.number="received" type="number" min="0" step="1" inputmode="numeric"
          placeholder="金額を入力" :disabled="pending" /></label>
      <p v-if="received !== ''" :class="['change-amount', { error: !validCash }]">
        {{
          validCash
            ? `お釣り ${(Number(received) - total).toLocaleString()} 円`
            : "お預かり金額が不足しているか、整数ではありません。"
        }}
      </p>
      <p v-if="saleError" class="error small" role="alert">
        {{ saleError
        }}<br />通信が途切れた場合は、再送信前に売上履歴で登録済みか確認してください。
      </p>
      <button class="button checkout-button" type="button" :disabled="!ready" @click="checkout">
        {{ pending ? "会計処理中…" : "会計を確定する" }}</button><button v-if="lines.length" type="button"
        class="button secondary clear-button" :disabled="pending" @click="clear">
        選択をすべてクリア
      </button>
    </aside>
  </div>
</template>
