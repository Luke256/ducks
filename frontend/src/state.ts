import { computed, ref, watch } from "vue";
import type { Festival } from "@/types/festival";
import { api } from "@/lib/api";
export function sessionRef(key: string, initial = "") {
  let saved = initial;
  try {
    const parsed: unknown = JSON.parse(sessionStorage.getItem(key) || "null");
    if (typeof parsed === "string") saved = parsed;
  } catch {
    /* 保存が無効でも画面内の操作は続けられる。 */
  }
  const value = ref(saved);
  watch(
    value,
    (next) => {
      try {
        sessionStorage.setItem(key, JSON.stringify(next));
      } catch {
        /* ブラウザが保存を拒否した場合。 */
      }
    },
    { flush: "sync" },
  );
  return value;
}
export const currentFestivalId = sessionRef("currentFestivalId");
export const pendingMutations = ref(0);
export const festivalSelectionLocked = computed(
  () => pendingMutations.value > 0,
);
export const stockFilterCategory = sessionRef("stockFilterCategory");
export const festivals = ref<Festival[]>([]);
export const festivalsLoading = ref(false);
export const festivalsError = ref("");
let festivalRequest = 0;
export async function loadFestivals() {
  const request = ++festivalRequest;
  festivalsLoading.value = true;
  festivalsError.value = "";
  try {
    const body = await api<{ festivals: Festival[] | null }>("/festivals");
    if (request !== festivalRequest) return;
    festivals.value = body.festivals || [];
    if (
      currentFestivalId.value &&
      !festivals.value.some((f) => f.id === currentFestivalId.value)
    )
      currentFestivalId.value = "";
  } catch (cause) {
    if (request !== festivalRequest) return;
    festivalsError.value =
      cause instanceof Error
        ? cause.message
        : "イベントを読み込めませんでした。";
  } finally {
    if (request === festivalRequest) festivalsLoading.value = false;
  }
}
let nextNoticeId = 0;
export const notices = ref<
  { id: number; message: string; type: "success" | "error" }[]
>([]);
export function dismissNotice(id: number) {
  notices.value = notices.value.filter((n) => n.id !== id);
}
export function notify(message: string, type: "success" | "error" = "success") {
  const id = ++nextNoticeId;
  notices.value.push({ id, message, type });
  if (type === "success") window.setTimeout(() => dismissNotice(id), 5000);
}
