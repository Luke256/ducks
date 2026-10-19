import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import type { Poster } from "@/types/poster";

let cache: typeof import("@/lib/posterLocations");
const festival = { id: "f1", name: "工大祭", description: "保存不要", date: "2026-10-18T00:00:00Z", locked: false };
const poster: Poster = {
    id: "p1", name: "01", description: "講義棟の入口", festival,
    status: "collected", image: [{ id: "image1", url: "/api/v1/images/image1" }],
};
beforeEach(async () => {
    localStorage.clear();
    vi.resetModules();
    cache = await import("@/lib/posterLocations");
});
afterEach(() => vi.restoreAllMocks());

describe("ポスターの場所の端末保存", () => {
    it("必要な文字情報と画像ID・URLだけを保存し、モジュールを開き直しても読める", async () => {
        cache.savePosterLocations(festival, [poster]);
        expect(JSON.parse(localStorage.getItem(cache.posterLocationsKey)!)).toEqual({
            festival: { id: "f1", name: "工大祭" },
            savedAt: expect.any(Number),
            posters: [{ id: "p1", name: "01", description: "講義棟の入口", status: "collected", image: poster.image }],
        });
        vi.resetModules();
        const reopened = await import("@/lib/posterLocations");
        expect(reopened.savedPosterLocations.value).toEqual(cache.savedPosterLocations.value);
    });
    it("空の一覧も含めて、新しいイベントで1件分の保存内容を置き換える", () => {
        cache.savePosterLocations(festival, [poster]);
        cache.savePosterLocations({ id: "f2", name: "別イベント" }, []);
        expect(localStorage.length).toBe(1);
        expect(cache.savedPosterLocations.value?.festival.id).toBe("f2");
        expect(cache.savedPosterLocations.value?.posters).toEqual([]);
        expect(localStorage.getItem(cache.posterLocationsKey)).not.toContain("p1");
    });
    it.each([
        "broken-json", "null", "[]",
        JSON.stringify({ festival: { id: "f1", name: "祭" }, savedAt: 1, posters: [{ id: "p1" }] }),
        JSON.stringify({ festival: { id: "f1", name: "祭" }, savedAt: "bad", posters: [] }),
    ])("壊れた保存値を無視して通常の取得を続けられる: %s", async (value) => {
        localStorage.setItem(cache.posterLocationsKey, value);
        vi.resetModules();
        expect((await import("@/lib/posterLocations")).savedPosterLocations.value).toBeNull();
    });
    it("保存容量不足でも以前のイベントを失わない", () => {
        cache.savePosterLocations(festival, [poster]);
        const previous = localStorage.getItem(cache.posterLocationsKey);
        vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => { throw new DOMException("full", "QuotaExceededError"); });
        cache.savePosterLocations({ id: "f2", name: "別イベント" }, []);
        expect(cache.savedPosterLocations.value?.festival.id).toBe("f1");
        expect(localStorage.getItem(cache.posterLocationsKey)).toBe(previous);
        expect(cache.posterLocationsError.value).toContain("保存できません");
    });
    it("保存設定が無効でも初期化できる", async () => {
        vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => { throw new Error("blocked"); });
        vi.resetModules();
        expect((await import("@/lib/posterLocations")).savedPosterLocations.value).toBeNull();
    });
    it("同じ詳細の取得では書き直さず、取得した場所の変更を反映する", () => {
        cache.savePosterLocations(festival, [poster]);
        const write = vi.spyOn(Storage.prototype, "setItem");
        cache.updateSavedPosterLocation(poster);
        expect(write).not.toHaveBeenCalled();
        cache.updateSavedPosterLocation({ ...poster, description: "北門" });
        expect(cache.savedPosterLocations.value?.posters[0].description).toBe("北門");
        cache.updateSavedPosterLocation({ ...poster, festival: { ...festival, id: "f2" } });
        expect(cache.savedPosterLocations.value?.festival.id).toBe("f1");
    });
    it.each([undefined, "invalid-status"])("以前の保存情報や不正な回収状況でも場所を失わない: %s", async (status) => {
        localStorage.setItem(cache.posterLocationsKey, JSON.stringify({
            festival: { id: "f1", name: "工大祭" }, savedAt: 1,
            posters: [{ id: "p1", name: "01", description: "講義棟の入口", status }],
        }));
        vi.resetModules();
        const reopened = await import("@/lib/posterLocations");
        expect(reopened.savedPosterLocations.value?.posters).toEqual([{ id: "p1", name: "01", description: "講義棟の入口" }]);
    });
    it("詳細から取得した回収状況を反映し、同じ状態では書き込まない", () => {
        cache.savePosterLocations(festival, [poster]);
        const write = vi.spyOn(Storage.prototype, "setItem");
        cache.updateSavedPosterLocation(poster);
        expect(write).not.toHaveBeenCalled();
        cache.updateSavedPosterLocation({ ...poster, status: "lost" });
        expect(write).toHaveBeenCalledTimes(1);
        expect(cache.savedPosterLocations.value?.posters[0]).toEqual({ id: "p1", name: "01", description: "講義棟の入口", status: "lost", image: poster.image });
        cache.updateSavedPosterLocation({ id: "p1", name: "02", description: "北門", festival });
        expect(cache.savedPosterLocations.value?.posters[0].status).toBe("lost");
    });
    it("画像だけの変更・削除も取得結果に合わせて反映する", () => {
        cache.savePosterLocations(festival, [poster]);
        const image = [{ id: "new-image", url: "https://example.com/api/v1/images/new-image" }];
        cache.updateSavedPosterLocation({ ...poster, image });
        expect(cache.savedPosterLocations.value?.posters[0].image).toEqual(image);
        const write = vi.spyOn(Storage.prototype, "setItem");
        cache.updateSavedPosterLocation({ ...poster, image });
        expect(write).not.toHaveBeenCalled();
        cache.updateSavedPosterLocation({ ...poster, image: [] });
        expect(cache.savedPosterLocations.value?.posters[0].image).toEqual([]);
    });
    it("不正な画像情報は捨てても文字情報を維持する", async () => {
        localStorage.setItem(cache.posterLocationsKey, JSON.stringify({
            festival, savedAt: 1, posters: [{ ...poster, image: [
                null, { id: "invalid", url: 1 }, { id: "script", url: "javascript:alert(1)" },
                { id: "data", url: "data:image/webp;base64,AAAA" }, { id: "broken", url: "https://[" },
                { ...poster.image[0], extra: "保存不要" },
            ] }],
        }));
        vi.resetModules();
        const reopened = await import("@/lib/posterLocations");
        expect(reopened.savedPosterLocations.value?.posters[0]).toEqual({
            id: "p1", name: "01", description: "講義棟の入口", status: "collected", image: poster.image,
        });
    });
    it("画像がない以前の保存形式も読み込める", async () => {
        localStorage.setItem(cache.posterLocationsKey, JSON.stringify({
            festival, savedAt: 1, posters: [{ id: "p1", name: "01", description: "講義棟の入口", status: "collected" }],
        }));
        vi.resetModules();
        const reopened = await import("@/lib/posterLocations");
        expect(reopened.savedPosterLocations.value?.posters[0].status).toBe("collected");
        expect(reopened.savedPosterLocations.value?.posters[0].image).toBeUndefined();
    });
    it("画像情報の保存はポスター1件につき10枚までに制限する", () => {
        cache.savePosterLocations(festival, [{ ...poster, image: Array.from({ length: 12 }, (_, index) => ({
            id: String(index), url: `/api/v1/images/${index}`,
        })) }]);
        expect(cache.savedPosterLocations.value?.posters[0].image).toHaveLength(10);
    });
});
