import { ref, shallowRef } from "vue";
import { maxPosterImages, type Poster, type PosterImage, type PosterStatus } from "@/types/poster";

export type PosterLocation = Pick<Poster, "id" | "name" | "description"> & Partial<Pick<Poster, "status" | "image">>;
type SavedLocations = {
    festival: { id: string; name: string };
    savedAt: number;
    posters: PosterLocation[];
};
export const posterLocationsKey = "ducks.poster-locations.v1";
// shortcut: 画像本体はブラウザ標準キャッシュに任せる。保持容量の制御が必要なら独自の画像保存を追加する。

function isStatus(value: unknown): value is PosterStatus {
    return value === "uncollected" || value === "collected" || value === "lost";
}
function isImage(value: unknown): value is PosterImage {
    if (!value || typeof value !== "object" || !("id" in value) || typeof value.id !== "string" ||
        !("url" in value) || typeof value.url !== "string" || !value.url) return false;
    try {
        return ["http:", "https:"].includes(new URL(value.url, window.location.origin).protocol);
    } catch {
        return false;
    }
}
function location({ id, name, description, status, image }: PosterLocation): PosterLocation {
    return {
        id, name, description,
        ...(isStatus(status) ? { status } : {}),
        ...(Array.isArray(image) ? { image: image.slice(0, maxPosterImages).filter(isImage).map(({ id, url }) => ({ id, url })) } : {}),
    };
}

function isLocation(value: unknown): value is PosterLocation {
    return !!value && typeof value === "object" &&
        "id" in value && typeof value.id === "string" &&
        "name" in value && typeof value.name === "string" &&
        "description" in value && typeof value.description === "string";
}
function read(): SavedLocations | null {
    try {
        const value = JSON.parse(localStorage.getItem(posterLocationsKey) || "null");
        if (!value || typeof value.festival?.id !== "string" ||
            typeof value.festival.name !== "string" ||
            !Number.isFinite(value.savedAt) || value.savedAt <= 0 ||
            !Array.isArray(value.posters) || !value.posters.every(isLocation)) return null;
        return {
            festival: { id: value.festival.id, name: value.festival.name },
            savedAt: value.savedAt,
            posters: value.posters.map(location),
        };
    } catch {
        return null;
    }
}
export const savedPosterLocations = shallowRef<SavedLocations | null>(read());
export const posterLocationsError = ref("");

export function savePosterLocations(festival: SavedLocations["festival"], posters: readonly PosterLocation[]) {
    const next: SavedLocations = {
        festival: { id: festival.id, name: festival.name },
        savedAt: Date.now(),
        posters: posters.map(location),
    };
    try {
        localStorage.setItem(posterLocationsKey, JSON.stringify(next));
        savedPosterLocations.value = next;
        posterLocationsError.value = "";
    } catch {
        // 保存に失敗しても、以前に保存できたイベントの場所は残す。
        posterLocationsError.value = "場所情報を端末に保存できませんでした。ブラウザの保存設定や空き容量を確認してください。";
    }
}

export function updateSavedPosterLocation(poster: PosterLocation & Pick<Poster, "festival">) {
    const saved = savedPosterLocations.value;
    if (!saved || saved.festival.id !== poster.festival.id) return;
    const previous = saved.posters.find((item) => item.id === poster.id);
    const next = location({ ...poster, status: poster.status ?? previous?.status, image: poster.image ?? previous?.image });
    if (previous?.name === next.name && previous.description === next.description && previous.status === next.status &&
        JSON.stringify(previous.image) === JSON.stringify(next.image) && saved.festival.name === poster.festival.name) return;
    savePosterLocations({ id: saved.festival.id, name: poster.festival.name }, [
        ...saved.posters.filter((item) => item.id !== poster.id),
        next,
    ]);
}
