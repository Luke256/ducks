<script setup lang="ts">
import { computed, ref } from "vue";
import { currentFestivalId } from "@/state";
import { useResource, listOf } from "@/composables/useResource";
import { imageUrl } from "@/lib/api";
import {
  PosterStatusLabels,
  type Poster,
  type PosterStatus,
} from "@/types/poster";
import PosterStatusPicker from "@/components/PosterStatusPicker.vue";
import ResourceState from "@/components/ResourceState.vue";
const {
  data: posters,
  error,
  loading,
  reload,
} = useResource(
  () =>
    currentFestivalId.value
      ? `/festivals/${currentFestivalId.value}/posters`
      : null,
  listOf<Poster>("posters"),
);
const search = ref("");
const status = ref("");
const filtered = computed(() =>
  (posters.value || []).filter(
    (p) =>
      (!status.value || p.status === status.value) &&
      `${p.name} ${p.description}`.includes(search.value.trim()),
  ),
);
function updateStatus(id: string, next: PosterStatus) {
  posters.value = (posters.value || []).map((p) =>
    p.id === id ? { ...p, status: next } : p,
  );
}
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">POSTERS</p>
      <h1>ポスター管理</h1>
      <p class="muted">設置場所を確認しながら、回収状況をその場で更新。</p>
    </div>
    <RouterLink to="/poster/new" class="button">＋ ポスターを登録</RouterLink>
  </div>
  <p v-if="!currentFestivalId" class="state">
    上の「対象イベント」からイベントを選択してください。
  </p>
  <template v-else>
    <div v-if="posters" class="summary-grid">
      <button v-for="(label, key) in PosterStatusLabels" :key="key"
        :class="['panel summary', key, { chosen: status === key }]" @click="status = status === key ? '' : key">
        <span>{{ label }}</span><strong>{{posters.filter((p) => p.status === key).length
        }}<small> 枚</small></strong>
      </button>
    </div>
    <div class="toolbar">
      <label class="search-field">ポスター検索<input v-model="search" type="search" placeholder="名前・設置場所で検索" /></label><label
        class="filter-field">回収状況<select v-model="status">
          <option value="">すべての回収状況</option>
          <option v-for="(label, key) in PosterStatusLabels" :key="key" :value="key">
            {{ label }}
          </option>
        </select></label><button class="button secondary" :disabled="loading" @click="reload()">
        更新
      </button>
    </div>
    <ResourceState :loading="loading" :error="error" :empty="!filtered.length" empty-text="条件に一致するポスターがありません。"
      @retry="reload()" />
    <div v-if="filtered.length" class="panel table-scroll">
      <table>
        <thead>
          <tr>
            <th>ポスター / 設置場所</th>
            <th>回収状況</th>
            <th><span class="sr-only">詳細</span></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="poster in filtered" :key="poster.id">
            <td>
              <div class="item-cell">
                <img v-if="poster.image_url" :src="imageUrl(poster.image_url)" alt="" class="thumbnail"
                  loading="lazy" />
                <div>
                  <RouterLink :to="`/poster/detail/${poster.id}`" class="text-link">{{ poster.name }}</RouterLink>
                  <p class="muted description">{{ poster.description }}</p>
                </div>
              </div>
            </td>
            <td>
              <PosterStatusPicker :poster-id="poster.id" :status="poster.status" :name="poster.name"
                @updated="updateStatus(poster.id, $event)" />
            </td>
            <td>
              <RouterLink :to="`/poster/detail/${poster.id}`" class="text-link" :aria-label="`${poster.name}の詳細`">詳細 →
              </RouterLink>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
</template>
