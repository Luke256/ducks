<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useResource, listOf } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { currentFestivalId } from "@/state";
import { api, jsonBody, imageUrl } from "@/lib/api";
import { priceColor } from "@/lib/priceColors";
import { categoryTags } from "@/lib/categories";
import CategoryTags from "@/components/CategoryTags.vue";
import type { Stock } from "@/types/stock";
import type { StockItem } from "@/types/stockItem";
import ResourceState from "@/components/ResourceState.vue";
const route = useRoute();
const router = useRouter();
const isNew = computed(() => !route.params.stockId);
const {
  data: stock,
  error,
  loading,
  reload,
} = useResource<Stock>(() =>
  isNew.value ? null : `/stocks/${route.params.stockId}`,
);
const {
  data: items,
  error: itemsError,
  loading: itemsLoading,
  reload: reloadItems,
} = useResource(
  () => (isNew.value ? "/items" : null),
  listOf<StockItem>("items"),
);
const form = reactive({ item_id: "", description: "", price: 500 });
const category = ref("");
const edit = ref(isNew.value);
const categories = computed(() =>
  [...new Set((items.value || []).flatMap((i) => categoryTags(i.category)))].sort(),
);
const filtered = computed(() =>
  (items.value || []).filter(
    (i) => !category.value || categoryTags(i.category).includes(category.value),
  ),
);
const selectedItem = computed(() =>
  items.value?.find((i) => i.id === form.item_id),
);
const { pending, error: saveError, run } = useMutation();
watch(category, () => {
  if (
    selectedItem.value &&
    category.value &&
    !categoryTags(selectedItem.value.category).includes(category.value)
  )
    form.item_id = "";
});
watch(stock, (value) => {
  if (value) {
    currentFestivalId.value = value.festival_id;
    if (!edit.value) form.description = value.description;
  }
});
function startEdit() {
  if (stock.value) {
    form.description = stock.value.description;
    edit.value = true;
  }
}
async function save() {
  await run(
    async () => {
      if (isNew.value) {
        const festivalId = currentFestivalId.value;
        if (
          !festivalId ||
          !form.item_id ||
          !Number.isSafeInteger(form.price) ||
          form.price < 0
        )
          throw new Error(
            "イベント・商品・0円以上の整数の価格を入力してください。",
          );
        const created = await api<Stock>(`/festivals/${festivalId}/stocks`, {
          method: "POST",
          ...jsonBody(form),
        });
        currentFestivalId.value = festivalId;
        await router.push(`/sales/stocks/${created.id}`);
      } else {
        await api(`/stocks/${route.params.stockId}`, {
          method: "PUT",
          ...jsonBody({ description: form.description }),
        });
        edit.value = false;
        await reload();
      }
    },
    isNew.value ? "販売商品を登録しました。" : "販売商品を更新しました。",
  );
}
async function remove() {
  if (
    !window.confirm(
      `「${stock.value?.item.name}」の販売登録を削除しますか？関連する売上記録への影響を確認してください。この操作は取り消せません。`,
    )
  )
    return;
  await run(async () => {
    await api(`/stocks/${route.params.stockId}`, { method: "DELETE" });
    await router.push("/sales/stocks");
  }, "販売商品を削除しました。");
}
</script>
<template>
  <RouterLink to="/sales/stocks" class="back-link">← 販売商品一覧</RouterLink>
  <ResourceState :loading="loading || itemsLoading" :error="error || itemsError"
    @retry="isNew ? reloadItems() : reload()" />
  <template v-if="isNew || stock">
    <div class="page-heading">
      <div>
        <p class="eyebrow">
          {{ isNew ? "NEW CATALOG ITEM" : "CATALOG DETAIL" }}
        </p>
        <h1>{{ isNew ? "販売商品を登録" : stock?.item.name }}</h1>
        <p v-if="isNew" class="muted">
          商品マスターから選び、このイベントでの販売価格を設定します。
        </p>
      </div>
      <button v-if="!isNew" class="button secondary" :disabled="pending" @click="edit ? (edit = false) : startEdit()">
        {{ edit ? "編集をキャンセル" : "説明を編集" }}
      </button>
    </div>
    <p v-if="isNew && !currentFestivalId" class="state">
      上の「対象イベント」からイベントを選択してください。
    </p>
    <form v-if="edit" class="panel form-panel" @submit.prevent="save">
      <fieldset :disabled="pending">
        <template v-if="isNew"><label class="field">カテゴリで絞り込み<select v-model="category">
              <option value="">すべてのカテゴリ</option>
              <option v-for="value in categories" :key="value">
                {{ value }}
              </option>
            </select></label><label class="field">商品<select v-model="form.item_id" required>
              <option value="">商品を選択</option>
              <option v-for="item in filtered" :key="item.id" :value="item.id">
                {{ item.name }} / {{ categoryTags(item.category).join(" / ") }}
              </option>
            </select></label><img v-if="selectedItem?.image_url" :src="imageUrl(selectedItem.image_url)"
            :alt="selectedItem.name" class="image-preview" />
          <RouterLink v-if="!itemsLoading && !items?.length" to="/sales/items/new" class="text-link">
            先に商品マスターを登録 →</RouterLink><label class="field">販売価格（円）<input v-model.number="form.price" type="number"
              required min="0" step="1" inputmode="numeric" /></label>
        </template><label class="field">販売時の説明<textarea v-model="form.description" rows="3"
            placeholder="レジにも表示される補足情報" />
        </label>
        <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
        <div class="actions">
          <button class="button" type="submit" :disabled="isNew && !currentFestivalId">
            {{ pending ? "保存中…" : isNew ? "登録する" : "保存する" }}</button>
          <RouterLink v-if="isNew" to="/sales/stocks" class="button secondary">キャンセル</RouterLink><button v-else
            class="button danger" type="button" @click="remove">
            販売登録を削除
          </button>
        </div>
      </fieldset>
    </form>
    <div v-else-if="stock" class="detail-grid">
      <section class="panel">
        <img v-if="stock.item.image_url" :src="imageUrl(stock.item.image_url)" :alt="stock.item.name"
          class="detail-image" />
      </section>
      <section class="panel">
        <CategoryTags :category="stock.item.category" />
        <p class="price-large" :style="{ color: priceColor(stock.price) }">
          {{ stock.price.toLocaleString() }}<small> 円</small>
        </p>
        <p class="description">{{ stock.description || "説明なし" }}</p>
        <RouterLink :to="`/sales/items/${stock.item.id}`" class="text-link">商品マスターを開く →</RouterLink>
        <p class="muted small">
          販売価格は登録時の価格です。変更する場合は新しい販売登録を作成してください。
        </p>
      </section>
    </div>
  </template>
</template>
