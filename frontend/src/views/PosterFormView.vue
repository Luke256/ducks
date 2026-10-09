<script setup lang="ts">
import { reactive, shallowRef } from "vue";
import { useRouter } from "vue-router";
import { currentFestivalId } from "@/state";
import { api } from "@/lib/api";
import { resizeImage } from "@/utils/resizeImage";
import { useMutation } from "@/composables/useMutation";
import ImagesField from "@/components/ImagesField.vue";
import { maxPosterImages } from "@/types/poster";
const router = useRouter();
const form = reactive({ name: "", description: "" });
const images = shallowRef<File[]>([]);
const { pending, error, run } = useMutation();
async function save() {
  await run(async () => {
    const festivalId = currentFestivalId.value;
    if (
      !festivalId ||
      !images.value.length ||
      !form.name.trim() ||
      !form.description.trim()
    )
      throw new Error("イベント・名前・設置場所・画像を入力してください。");
    if (images.value.length > maxPosterImages)
      throw new Error(`写真は${maxPosterImages}枚まで登録できます。`);
    const body = new FormData();
    body.append("name", form.name.trim());
    body.append("description", form.description.trim());
    body.append("festival_id", festivalId);
    for (const image of images.value)
      body.append("image", await resizeImage(image));
    await api("/posters", { method: "POST", body });
    currentFestivalId.value = festivalId;
    await router.push("/poster");
  }, "ポスターを登録しました。");
}
</script>
<template>
  <RouterLink to="/poster" class="back-link">← ポスター一覧</RouterLink>
  <div class="page-heading">
    <div>
      <p class="eyebrow">NEW POSTER</p>
      <h1>ポスターを登録</h1>
      <p class="muted">ポスターと周辺の様子が写った写真を登録してください。</p>
    </div>
  </div>
  <p v-if="!currentFestivalId" class="state">
    上の「対象イベント」からイベントを選択してください。
  </p>
  <form class="panel form-panel" @submit.prevent="save">
    <fieldset :disabled="pending">
      <label class="field">ポスター名<input v-model="form.name" required maxlength="64" placeholder="例：42" /></label><label
        class="field">設置場所<textarea v-model="form.description" required maxlength="1024" rows="3"
          placeholder="回収する人が場所を特定できるように記載" /></label>
      <h3>設置場所の写真</h3>
      <ImagesField v-model="images" label="写真を追加" required />
      <p class="muted small">写真は1〜{{ maxPosterImages }}枚登録できます（選択中：{{ images.length }}枚）。画像はアップロード前に圧縮されます。</p>
      <p v-if="error" class="error" role="alert">{{ error }}</p>
      <div class="actions">
        <button class="button" type="submit" :disabled="!currentFestivalId">
          {{ pending ? "画像を圧縮・登録中…" : "登録する" }}</button>
        <RouterLink to="/poster" class="button secondary">キャンセル</RouterLink>
      </div>
    </fieldset>
  </form>
</template>
