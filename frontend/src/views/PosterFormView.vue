<script setup lang="ts">
import { reactive, shallowRef } from "vue";
import { useRouter } from "vue-router";
import { currentFestivalId } from "@/state";
import { api } from "@/lib/api";
import { resizeImage } from "@/utils/resizeImage";
import { useMutation } from "@/composables/useMutation";
import ImageField from "@/components/ImageField.vue";
const router = useRouter();
const form = reactive({ name: "", description: "" });
const image = shallowRef<File | null>(null);
const { pending, error, run } = useMutation();
async function save() {
  await run(async () => {
    const festivalId = currentFestivalId.value;
    if (
      !festivalId ||
      !image.value ||
      !form.name.trim() ||
      !form.description.trim()
    )
      throw new Error("イベント・名前・設置場所・画像を入力してください。");
    const body = new FormData();
    body.append("name", form.name.trim());
    body.append("description", form.description.trim());
    body.append("festival_id", festivalId);
    body.append("image", await resizeImage(image.value));
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
      <label class="field"
        >ポスター名<input
          v-model="form.name"
          required
          maxlength="64"
          placeholder="例：講義棟入口 01" /></label
      ><label class="field"
        >設置場所<textarea
          v-model="form.description"
          required
          maxlength="1024"
          rows="3"
          placeholder="回収する人が場所を特定できるように記載"
        /></label
      ><ImageField label="設置場所の写真" required @change="image = $event" />
      <p class="muted small">画像はアップロード前に圧縮されます。</p>
      <p v-if="error" class="error" role="alert">{{ error }}</p>
      <div class="actions">
        <button class="button" type="submit" :disabled="!currentFestivalId">
          {{ pending ? "画像を圧縮・登録中…" : "登録する" }}</button
        ><RouterLink to="/poster" class="button secondary"
          >キャンセル</RouterLink
        >
      </div>
    </fieldset>
  </form>
</template>
