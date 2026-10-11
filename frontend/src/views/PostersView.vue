<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { currentFestivalId, festivals } from "@/state";
import { useResource, listOf } from "@/composables/useResource";
import { imageUrl } from "@/lib/api";
import { savedPosterLocations, savePosterLocations, posterLocationsError, type PosterLocation } from "@/lib/posterLocations";
import {
  PosterStatusLabels,
  type Poster,
  type PosterStatus,
} from "@/types/poster";
import PosterStatusPicker from "@/components/PosterStatusPicker.vue";
import ResourceState from "@/components/ResourceState.vue";
if (!currentFestivalId.value && savedPosterLocations.value)
  currentFestivalId.value = savedPosterLocations.value.festival.id;
const {
  data: posters,
  error,
  loading,
  reload,
  cancel,
} = useResource(
  () =>
    currentFestivalId.value
      ? `/festivals/${currentFestivalId.value}/posters`
      : null,
  (body) => {
    const result = listOf<Poster>("posters")(body);
    const id = currentFestivalId.value;
    const name = result[0]?.festival.name || festivals.value.find((festival) => festival.id === id)?.name ||
      (savedPosterLocations.value?.festival.id === id ? savedPosterLocations.value.festival.name : id);
    savePosterLocations({ id, name }, result);
    return result;
  },
);
const saved = computed(() => savedPosterLocations.value?.festival.id === currentFestivalId.value
  ? savedPosterLocations.value : null);
const showingSaved = computed(() => !posters.value && !!saved.value);
const displayPosters = computed<(PosterLocation & Partial<Pick<Poster, "status" | "image">>)[]>(
  () => posters.value ?? saved.value?.posters ?? [],
);
const hasStatuses = computed(() => displayPosters.value.some((poster) => !!poster.status));
const search = ref("");
const status = ref("");
watch(showingSaved, (value) => { if (value) status.value = ""; }, { flush: "sync" });
const filtered = computed(() =>
  displayPosters.value.filter(
    (p) =>
      (!status.value || p.status === status.value) &&
      `${p.name} ${p.description}`.includes(search.value.trim()),
    ).sort((a, b) => {
      const aNum = parseInt(a.name.match(/\d+/)?.[0] || "0", 10);
      const bNum = parseInt(b.name.match(/\d+/)?.[0] || "0", 10);
      return aNum - bNum;
    }),
);
function updateStatus(id: string, next: PosterStatus) {
  // 更新前に取得を始めた応答で、保存成功後の回収状況を戻さない。
  cancel();
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
  <div v-if="savedPosterLocations && !saved && !posters" class="state" role="status">
    <p>選択中のイベントの場所情報は端末に保存されていません。</p>
    <button class="button secondary" @click="currentFestivalId = savedPosterLocations.festival.id">
      保存済みの「{{ savedPosterLocations.festival.name }}」の場所を見る
    </button>
  </div>
  <template v-if="currentFestivalId">
    <p v-if="showingSaved && saved" class="small muted" role="status">
      端末に保存された情報を表示しています。
      <span v-if="loading">最新情報を取得中…</span>
    </p>
    <p v-if="posterLocationsError" class="state error" role="alert">{{ posterLocationsError }}</p>
    <div v-if="posters || hasStatuses" class="summary-grid">
      <button v-for="(label, key) in PosterStatusLabels" :key="key"
        :class="['panel summary', key, { chosen: status === key }]" @click="status = status === key ? '' : key">
        <span>{{ label }}</span><strong>{{displayPosters.filter((p) => p.status === key).length
        }}<small> 枚</small></strong>
      </button>
    </div>
    <div class="toolbar">
      <label class="search-field">ポスター検索<input v-model="search" type="search" placeholder="名前・設置場所で検索" /></label><label v-if="posters || hasStatuses"
        class="filter-field">回収状況<select v-model="status">
          <option value="">すべての回収状況</option>
          <option v-for="(label, key) in PosterStatusLabels" :key="key" :value="key">
            {{ label }}
          </option>
        </select></label><button class="button secondary" :disabled="loading" @click="reload()">
        更新
      </button>
    </div>
    <ResourceState :loading="loading && !showingSaved" :error="error" :empty="!filtered.length" empty-text="条件に一致するポスターがありません。"
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
                <img v-if="poster.image?.[0]" :src="imageUrl(poster.image[0].url)" alt="" class="thumbnail"
                  loading="lazy" />
                <div>
                  <RouterLink :to="`/poster/detail/${poster.id}`" class="text-link">{{ poster.name }}</RouterLink>
                  <p class="muted description">{{ poster.description }}</p>
                </div>
              </div>
            </td>
            <td>
              <span v-if="showingSaved && poster.status" class="badge">{{ PosterStatusLabels[poster.status] }}</span>
              <PosterStatusPicker v-else-if="poster.status" :poster-id="poster.id" :status="poster.status" :name="poster.name"
                @updated="updateStatus(poster.id, $event)" />
              <span v-else class="muted">未取得</span>
            </td>
            <td>
              <RouterLink :to="`/poster/detail/${poster.id}`" class="text-link" :aria-label="`${poster.name}の詳細`">詳細 →
              </RouterLink>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="saved" class="small muted">
      保存日時：{{ new Date(saved.savedAt).toLocaleString('ja-JP') }}
    </p>
  </template>
</template>
