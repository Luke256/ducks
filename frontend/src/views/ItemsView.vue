<script setup lang="ts">
import { computed, ref } from "vue";
import { useResource, listOf } from "@/composables/useResource";
import { imageUrl } from "@/lib/api";
import { categoryTags } from "@/lib/categories";
import CategoryTags from "@/components/CategoryTags.vue";
import type { StockItem } from "@/types/stockItem";
import ResourceState from "@/components/ResourceState.vue";
const {
  data: items,
  loading,
  error,
  reload,
} = useResource("/items", listOf<StockItem>("items"));
const search = ref("");
const category = ref("");
const categories = computed(() =>
  [...new Set((items.value || []).flatMap((i) => categoryTags(i.category)))].sort(),
);
const filtered = computed(() =>
  (items.value || [])
    .filter(
      (i) =>
        (!category.value || categoryTags(i.category).includes(category.value)) &&
        `${i.name} ${i.description}`.includes(search.value.trim()),
    )
    .sort(
      (a, b) =>
        a.category.localeCompare(b.category, "ja") ||
        a.name.localeCompare(b.name, "ja"),
    ),
);
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">ITEM MASTER</p>
      <h1>商品マスター</h1>
      <p class="muted">
        商品名・写真・カテゴリを登録。イベント共通で使用します。
      </p>
    </div>
    <RouterLink to="/sales/items/new" class="button">＋ 商品を登録</RouterLink>
  </div>
  <div class="toolbar">
    <label class="search-field">商品検索<input v-model="search" type="search" placeholder="名前・説明で検索" /></label><label
      class="filter-field">カテゴリ<select v-model="category">
        <option value="">すべてのカテゴリ</option>
        <option v-for="value in categories" :key="value">{{ value }}</option>
      </select></label><button class="button secondary" :disabled="loading" @click="reload()">
      更新
    </button>
  </div>
  <ResourceState :loading="loading" :error="error" :empty="!filtered.length" empty-text="条件に一致する商品がありません。"
    @retry="reload()" />
  <div class="product-grid">
    <RouterLink v-for="item in filtered" :key="item.id" :to="`/sales/items/${item.id}`" class="panel product-card"><img
        v-if="item.image_url" :src="imageUrl(item.image_url)" :alt="item.name" class="product-image" loading="lazy" />
      <div class="product-info">
        <CategoryTags :category="item.category" />
        <h2>{{ item.name }}</h2>
        <p class="description muted">{{ item.description || "説明なし" }}</p>
        <span class="text-link">詳細・編集 →</span>
      </div>
    </RouterLink>
  </div>
</template>
