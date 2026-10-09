<script setup lang="ts">
import { computed, reactive, ref, shallowRef, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useResource } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { api, jsonBody, imageUrl } from "@/lib/api";
import { currentFestivalId } from "@/state";
import { maxPosterImages, type Poster, type PosterImage, type PosterStatus } from "@/types/poster";
import { resizeImage } from "@/utils/resizeImage";
import ImagesField from "@/components/ImagesField.vue";
import PosterStatusPicker from "@/components/PosterStatusPicker.vue";
import ResourceState from "@/components/ResourceState.vue";
const route = useRoute();
const router = useRouter();
const {
  data: poster,
  error,
  loading,
  reload,
} = useResource<Poster>(() => `/posters/${route.params.posterId}`);
const edit = ref(false);
const form = reactive({ name: "", description: "" });
const images = shallowRef<File[]>([]);
const deleteImageIds = ref<string[]>([]);
const imageCount = computed(() =>
  (poster.value?.image.filter((image) => !deleteImageIds.value.includes(image.id)).length || 0) + images.value.length,
);
const metadataSaved = ref(false);
const { pending, error: saveError, run } = useMutation();
watch(poster, (value) => {
  if (value) currentFestivalId.value = value.festival.id;
});
function startEdit() {
  if (poster.value) {
    images.value = [];
    deleteImageIds.value = [];
    saveError.value = "";
    metadataSaved.value = false;
    Object.assign(form, { name: poster.value.name, description: poster.value.description });
    edit.value = true;
  }
}
function cancelEdit() {
  edit.value = false;
  images.value = [];
  deleteImageIds.value = [];
  saveError.value = "";
}
watch(() => route.params.posterId, cancelEdit);
function updateStatus(status: PosterStatus) {
  if (poster.value) poster.value = { ...poster.value, status };
}
async function save() {
  await run(async () => {
    const currentPoster = poster.value;
    if (!currentPoster) throw new Error("ポスターを読み込んでください。");
    if (!form.name.trim() || !form.description.trim())
      throw new Error("ポスター名と設置場所を入力してください。");
    if (imageCount.value < 1)
      throw new Error("写真を1枚以上残すか、新しい写真を追加してください。");
    if (imageCount.value > maxPosterImages)
      throw new Error(`写真は${maxPosterImages}枚まで登録できます。`);
    const name = form.name.trim();
    const description = form.description.trim();
    const body = new FormData();
    const files = images.value;
    const deleteIds = [...deleteImageIds.value];
    for (const image of files) body.append("image", await resizeImage(image));
    for (const id of deleteIds) body.append("delete_image_ids", id);
    if (name !== currentPoster.name || description !== currentPoster.description) {
      await api(`/posters/${currentPoster.id}`, {
        method: "PUT",
        ...jsonBody({ name, description }),
      });
      if (poster.value?.id === currentPoster.id) {
        poster.value = { ...poster.value, name, description };
        metadataSaved.value = true;
      }
    }
    if (files.length || deleteIds.length) {
      const result = await api<{ image: PosterImage[] }>(`/posters/${currentPoster.id}/images`, {
        method: "PATCH",
        body,
      });
      if (poster.value?.id === currentPoster.id)
        poster.value = { ...poster.value, image: result.image };
    }
    if (route.params.posterId === currentPoster.id) {
      cancelEdit();
      await reload();
    }
  }, "ポスターを更新しました。");
}
async function remove() {
  if (
    !window.confirm(
      `「${poster.value?.name}」を削除しますか？この操作は取り消せません。`,
    )
  )
    return;
  await run(async () => {
    await api(`/posters/${route.params.posterId}`, { method: "DELETE" });
    await router.push("/poster");
  }, "ポスターを削除しました。");
}
</script>
<template>
  <RouterLink to="/poster" class="back-link">← ポスター一覧</RouterLink>
  <ResourceState :loading="loading" :error="error" @retry="reload()" />
  <template v-if="poster">
    <div class="page-heading">
      <div>
        <p class="eyebrow">POSTER DETAIL / {{ poster.festival.name }}</p>
        <h1>{{ poster.name }}</h1>
      </div>
      <button class="button secondary" :disabled="pending" @click="edit ? cancelEdit() : startEdit()">
        {{ edit ? "編集をキャンセル" : "編集する" }}
      </button>
    </div>
    <div class="detail-grid">
      <section class="panel">
        <h2>設置場所</h2>
        <p class="description">{{ poster.description }}</p>
        <div v-if="poster.image.length" class="poster-images">
          <a v-for="(image, index) in poster.image" :key="image.id" :href="imageUrl(image.url)" target="_blank" rel="noopener"
            :aria-label="`設置場所の写真${index + 1}を拡大`"><img :src="imageUrl(image.url)"
              :alt="`${poster.name}の設置場所（写真${index + 1}）`" class="detail-image" loading="lazy" /></a>
        </div>
      </section>
      <section class="panel">
        <h2>回収状況</h2>
        <PosterStatusPicker :poster-id="poster.id" :name="poster.name" :status="poster.status" :disabled="pending"
          @updated="updateStatus" />
        <p class="muted small">変更するとすぐに保存されます。</p>
      </section>
    </div>
    <form v-if="edit" class="panel form-panel" @submit.prevent="save">
      <h2>登録情報を編集</h2>
      <fieldset :disabled="pending">
        <label class="field">ポスター名<input v-model="form.name" required maxlength="64" /></label><label
          class="field">設置場所<textarea v-model="form.description" required maxlength="1024" rows="4" />
        </label>
        <h3>設置場所の写真</h3>
        <div class="image-selection">
          <figure v-for="(image, index) in poster.image" :key="image.id">
            <img :src="imageUrl(image.url)" :alt="`登録済みの写真${index + 1}`" class="image-preview" />
            <figcaption>
              <label class="image-delete"><input v-model="deleteImageIds" type="checkbox" :value="image.id" />写真{{ index + 1 }}を削除する</label>
            </figcaption>
          </figure>
        </div>
        <ImagesField v-model="images" label="写真を追加" />
        <p class="muted small">保存後の写真：{{ imageCount }} / {{ maxPosterImages }}枚。削除する写真を選び、新しい写真を追加すると差し替えられます。画像はアップロード前に圧縮されます。</p>
        <p v-if="imageCount < 1" class="error" role="alert">写真を1枚以上残すか、新しい写真を追加してください。</p>
        <p v-else-if="imageCount > maxPosterImages" class="error" role="alert">写真は{{ maxPosterImages }}枚まで登録できます。</p>
        <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
        <p v-if="saveError && metadataSaved" class="muted small">ポスター名・設置場所は保存済みです。写真の変更を確認し、再度保存してください。</p>
        <div class="actions">
          <button class="button" type="submit" :disabled="imageCount < 1 || imageCount > maxPosterImages">
            {{ pending ? "保存中…" : "保存する" }}</button><button class="button danger" type="button" @click="remove">
            ポスターを削除
          </button>
        </div>
      </fieldset>
    </form>
  </template>
</template>
