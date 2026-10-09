import {
    onScopeDispose,
    ref,
    shallowRef,
    toValue,
    watch,
    type MaybeRefOrGetter,
} from "vue";
import { api } from "@/lib/api";
export function useResource<T>(
    path: MaybeRefOrGetter<string | null>,
    select: (body: unknown) => T = (body) => body as T,
) {
    const data = shallowRef<T | null>(null);
    const error = ref("");
    const loading = ref(false);
    let controller: AbortController | undefined;
    async function reload(clear = false) {
        controller?.abort();
        const request = new AbortController();
        controller = request;
        const url = toValue(path);
        if (clear || !url) data.value = null;
        error.value = "";
        loading.value = !!url;
        if (!url) return;
        try {
            const result = await api<unknown>(url, { signal: request.signal });
            if (!request.signal.aborted) data.value = select(result);
        } catch (cause) {
            if (!request.signal.aborted)
                error.value =
                    cause instanceof Error ? cause.message : "読み込みに失敗しました。";
        } finally {
            if (!request.signal.aborted) loading.value = false;
        }
    }
    watch(
        () => toValue(path),
        () => {
            void reload(true);
        },
        { immediate: true },
    );
    const refresh = () => {
        void reload();
    };
    window.addEventListener("focus", refresh);
    onScopeDispose(() => {
        controller?.abort();
        window.removeEventListener("focus", refresh);
    });
    return { data, error, loading, reload };
}
export function listOf<T>(key: string) {
    return (body: unknown): T[] =>
        (body as Record<string, T[] | null>)[key] || [];
}
