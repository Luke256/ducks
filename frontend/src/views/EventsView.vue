<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  festivals,
  festivalsLoading,
  festivalsError,
  loadFestivals,
  currentFestivalId,
  festivalSelectionLocked,
} from "@/state";
import { api, jsonBody } from "@/lib/api";
import type { Festival } from "@/types/festival";
import { useMutation } from "@/composables/useMutation";
import ResourceState from "@/components/ResourceState.vue";
const router = useRouter();
const formOpen = ref(false);
const search = ref("");
const form = reactive({ name: "", description: "" });
const filtered = computed(() =>
  festivals.value.filter((f) =>
    `${f.name} ${f.description}`.includes(search.value.trim()),
  ),
);
const { pending, error, run } = useMutation();
async function create() {
  await run(async () => {
    if (!form.name.trim()) throw new Error("イベント名を入力してください。");
    const event = await api<Festival>("/festivals", {
      method: "POST",
      ...jsonBody({ name: form.name.trim(), description: form.description }),
    });
    currentFestivalId.value = event.id;
    await loadFestivals();
    await router.push(`/event/${event.id}`);
  }, "イベントを作成しました。");
}
</script>
<template>
  <div class="page-heading">
    <div>
      <p class="eyebrow">EVENTS</p>
      <h1>イベント管理</h1>
      <p class="muted">準備から当日の運営まで、イベントごとにまとめて管理。</p>
    </div>
    <button
      class="button"
      :aria-expanded="formOpen"
      @click="formOpen = !formOpen"
    >
      {{ formOpen ? "閉じる" : "＋ イベントを作成" }}
    </button>
  </div>
  <form v-if="formOpen" class="panel form-panel" @submit.prevent="create">
    <h2>新しいイベント</h2>
    <fieldset :disabled="pending">
      <label class="field"
        >イベント名<input
          v-model="form.name"
          required
          placeholder="例：工大祭 2026" /></label
      ><label class="field"
        >概要<textarea
          v-model="form.description"
          rows="3"
          placeholder="開催内容や運営メモ"
        />
      </label>
      <p v-if="error" class="error" role="alert">{{ error }}</p>
      <button class="button" type="submit">
        {{ pending ? "作成中…" : "作成する" }}
      </button>
    </fieldset>
  </form>
  <div class="toolbar">
    <label class="search-field"
      >イベント検索<input
        v-model="search"
        type="search"
        placeholder="名前・概要で検索" /></label
    ><span class="muted">{{ filtered.length }} 件</span
    ><button
      class="button secondary"
      :disabled="festivalsLoading"
      @click="loadFestivals"
    >
      更新
    </button>
  </div>
  <ResourceState
    :loading="festivalsLoading"
    :error="festivalsError"
    :empty="!filtered.length"
    empty-text="イベントがありません。新しいイベントを作成してください。"
    @retry="loadFestivals"
  />
  <div v-if="filtered.length" class="event-grid">
    <article v-for="event in filtered" :key="event.id" class="panel event-card">
      <div class="card-top">
        <span class="badge">EVENT</span
        ><span v-if="currentFestivalId === event.id" class="badge selected"
          >選択中</span
        >
      </div>
      <h2>
        <RouterLink :to="`/event/${event.id}`">{{ event.name }}</RouterLink>
      </h2>
      <p class="description muted">
        {{ event.description || "概要は登録されていません。" }}
      </p>
    < div class="actions" >
        <div>
            <button
            type="button"
            class="button secondary"
            :disabled="
            festivalsLoading ||
            festivalSelectionLocked ||
            currentFestivalId === event.id
            "
            @click="currentFestivalId = event.id"
            >
                {{ currentFestivalId === event.id ? "選択中" : "このイベントを選択" }}
            </button>
        </div>
        <div>
            <RouterLink :to="`/event/${event.id}`" class="text-link"
            >詳細・運営ツールを開く <span aria-hidden="true">→</span></RouterLink
              >
        </div>
      </div>
    </article>
  </div>
</template>
