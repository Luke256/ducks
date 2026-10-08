<script setup lang="ts">
import { computed, reactive, ref, shallowRef, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useResource } from "@/composables/useResource";
import { useMutation } from "@/composables/useMutation";
import { api, jsonBody, imageUrl } from "@/lib/api";
import { resizeImage } from "@/utils/resizeImage";
import type { StockItem } from "@/types/stockItem";
import ImageField from "@/components/ImageField.vue";
import ResourceState from "@/components/ResourceState.vue";
const route = useRoute();
const router = useRouter();
const isNew = computed(() => !route.params.itemId);
const {
  data: item,
  error,
  loading,
  reload,
} = useResource<StockItem>(() =>
  isNew.value ? null : `/items/${route.params.itemId}`,
);
const edit = ref(isNew.value);
const form = reactive({ name: "", category: "", description: "" });
const image = shallowRef<File | null>(null);
const { pending, error: saveError, run } = useMutation();
watch(item, (value) => {
  if (value && !edit.value)
    Object.assign(form, {
      name: value.name,
      category: value.category,
      description: value.description,
    });
});
function startEdit() {
  if (item.value) {
    image.value = null;
    saveError.value = "";
    Object.assign(form, {
      name: item.value.name,
      category: item.value.category,
      description: item.value.description,
    });
    edit.value = true;
  }
}
async function save() {
  await run(
    async () => {
      if (!form.name.trim() || !form.category.trim())
        throw new Error("商品名とカテゴリを入力してください。");
      if (isNew.value) {
        if (!image.value) throw new Error("商品画像を選択してください。");
        const body = new FormData();
        body.append("name", form.name.trim());
        body.append("category", form.category.trim());
        body.append("description", form.description);
        body.append("image", await resizeImage(image.value));
        const created = await api<StockItem>("/items", {
          method: "POST",
          body,
        });
        await router.push(`/sales/items/${created.id}`);
      } else {
        const body = new FormData();
        if (image.value) body.append("image", await resizeImage(image.value));
        await api(`/items/${route.params.itemId}`, {
          method: "PUT",
          ...jsonBody({
            name: form.name.trim(),
            category: form.category.trim(),
            description: form.description,
          }),
        });
        if (image.value)
          await api(`/items/${route.params.itemId}/image`, {
            method: "PUT",
            body,
          });
        image.value = null;
        edit.value = false;
        await reload();
      }
    },
    isNew.value ? "商品を登録しました。" : "商品を更新しました。",
  );
}
async function remove() {
  if (
    !window.confirm(
      `「${item.value?.name}」を削除しますか？関連する販売商品・売上記録への影響を確認してください。この操作は取り消せません。`,
    )
  )
    return;
  await run(async () => {
    await api(`/items/${route.params.itemId}`, { method: "DELETE" });
    await router.push("/sales/items");
  }, "商品を削除しました。");
}
</script>
<template>
  <RouterLink to="/sales/items" class="back-link">← 商品マスター</RouterLink
  ><ResourceState :loading="loading" :error="error" @retry="reload()" />
  <template v-if="isNew || item"
    ><div class="page-heading">
      <div>
        <p class="eyebrow">{{ isNew ? "NEW ITEM" : "ITEM DETAIL" }}</p>
        <h1>{{ isNew ? "商品を登録" : item?.name }}</h1>
      </div>
      <button
        v-if="!isNew"
        class="button secondary"
        :disabled="pending"
        @click="edit ? (edit = false) : startEdit()"
      >
        {{ edit ? "編集をキャンセル" : "編集する" }}
      </button>
    </div>
    <form v-if="edit" class="panel form-panel" @submit.prevent="save">
      <fieldset :disabled="pending">
        <label class="field"
          >商品名<input v-model="form.name" required maxlength="100" /></label
        ><label class="field"
          >カテゴリ<input
            v-model="form.category"
            required
            maxlength="100"
            placeholder="例：アクリルキーホルダー" /></label
        ><label class="field"
          >説明<textarea v-model="form.description" rows="4" /></label
        ><ImageField
          label="商品画像"
          :required="isNew"
          @change="image = $event"
        /><img
          v-if="!image && item?.image_url"
          :src="imageUrl(item.image_url)"
          :alt="item.name"
          class="image-preview"
        />
        <p v-if="!isNew" class="small muted">
          新しい画像を選ぶと差し替えます。選ばなければ現在の画像を保持します。
          この商品の画像はすべてのイベントで共通です。
        </p>
        <p v-if="saveError" class="error" role="alert">{{ saveError }}</p>
        <p v-if="saveError && !isNew && image" class="small muted">
          商品情報だけ保存されている場合があります。画像を確認し、再度保存してください。
        </p>
        <div class="actions">
          <button class="button" type="submit">
            {{ pending ? "保存中…" : isNew ? "登録する" : "保存する" }}</button
          ><RouterLink v-if="isNew" to="/sales/items" class="button secondary"
            >キャンセル</RouterLink
          ><button v-else class="button danger" type="button" @click="remove">
            商品を削除
          </button>
        </div>
      </fieldset>
    </form>
    <div v-else-if="item" class="detail-grid">
      <section class="panel">
        <a
          v-if="item.image_url"
          :href="imageUrl(item.image_url)"
          target="_blank"
          rel="noopener"
          :aria-label="`${item.name}の写真を拡大`"
          ><img
            :src="imageUrl(item.image_url)"
            :alt="item.name"
            class="detail-image"
        /></a>
      </section>
      <section class="panel">
        <h2>登録情報</h2>
        <dl>
          <dt>カテゴリ</dt>
          <dd>{{ item.category }}</dd>
          <dt>説明</dt>
          <dd class="description">{{ item.description || "説明なし" }}</dd>
        </dl>
        <RouterLink to="/sales/stocks/new" class="button secondary"
          >イベントの販売商品に追加</RouterLink
        >
      </section>
    </div>
  </template>
</template>
