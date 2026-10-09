import { ref } from "vue";
import { notify, pendingMutations } from "@/state";
export function useMutation() {
    const pending = ref(false);
    const error = ref("");
    async function run(action: () => Promise<void>, success: string) {
        if (pending.value) return false;
        pending.value = true;
        pendingMutations.value++;
        error.value = "";
        try {
            await action();
            notify(success);
            return true;
        } catch (cause) {
            error.value =
                cause instanceof Error ? cause.message : "処理に失敗しました。";
            notify(error.value, "error");
            return false;
        } finally {
            pending.value = false;
            pendingMutations.value--;
        }
    }
    return { pending, error, run };
}
