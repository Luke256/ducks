<script setup lang="ts">
import { computed, ref } from "vue";
import { currentFestivalId, stockFilterCategory } from "@/state";
import { useResource, listOf } from "@/composables/useResource";
import { imageUrl } from "@/lib/api";
import { categoryTags } from "@/lib/categories";
import CategoryTags from "@/components/CategoryTags.vue";
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
const search = ref("");
const categories = computed(() =>
  [...new Set((stocks.value || []).flatMap((s) => categoryTags(s.item.category)))].sort(),
);
const filtered = computed(() =>
  (stocks.value || [])
    .filter(
      (s) =>
        (!stockFilterCategory.value ||
          categoryTags(s.item.category).includes(stockFilterCategory.value)) &&
        `${s.item.name} ${s.description}`.includes(search.value.trim()),
    )
    .sort(
      (a, b) =>
        a.item.category.localeCompare(b.item.category, "ja") ||
        a.item.name.localeCompare(b.item.name, "ja"),
    ),
);
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">CATALOG</p>
      <h1>販売商品</h1>
      <p class="muted">このイベントで販売する商品と価格を管理します。</p>
    </div>
    <RouterLink to="/sales/stocks/new" class="button">＋ 販売商品を登録</RouterLink>
  </div>
  <p v-if="!currentFestivalId" class="state">
    上の「対象イベント」からイベントを選択してください。
  </p>
  <template v-else>
    <div class="toolbar">
      <label class="search-field">販売商品検索<input v-model="search" type="search" placeholder="名前・説明で検索" /></label><label
        class="filter-field">カテゴリ<select v-model="stockFilterCategory">
          <option value="">すべてのカテゴリ</option>
          <option v-if="
            stockFilterCategory && !categories.includes(stockFilterCategory)
          " :value="stockFilterCategory">
            {{ stockFilterCategory }}
          </option>
          <option v-for="value in categories" :key="value">{{ value }}</option>
        </select></label><button class="button secondary" :disabled="loading" @click="reload()">
        更新
      </button>
    </div>
    <ResourceState :loading="loading" :error="error" :empty="!filtered.length" empty-text="条件に一致する販売商品がありません。"
      @retry="reload()" />
    <div v-if="filtered.length" class="panel table-scroll">
      <table>
        <thead>
          <tr>
            <th>商品</th>
            <th>説明</th>
            <th class="numeric">価格</th>
            <th><span class="sr-only">詳細</span></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in filtered" :key="stock.id">
            <td>
              <div class="item-cell">
                <img v-if="stock.item.image_url" :src="imageUrl(stock.item.image_url)" alt="" class="thumbnail"
                  loading="lazy" />
                <div>
                  <RouterLink :to="`/sales/stocks/${stock.id}`" class="text-link">{{ stock.item.name
                  }}</RouterLink>
                  <p><CategoryTags :category="stock.item.category" /></p>
                </div>
              </div>
            </td>
            <td class="description">{{ stock.description || "—" }}</td>
            <td class="numeric">{{ stock.price.toLocaleString() }} 円</td>
            <td>
              <RouterLink :to="`/sales/stocks/${stock.id}`" class="text-link" :aria-label="`${stock.item.name}の販売詳細`">詳細
                →</RouterLink>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
</template>
