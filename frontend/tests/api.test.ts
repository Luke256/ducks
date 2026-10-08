import { afterEach, describe, expect, it, vi } from "vitest";
import { effectScope, nextTick, ref } from "vue";
import { api, apiBase, imageUrl, jsonBody } from "@/lib/api";
import { useResource, listOf } from "@/composables/useResource";
import { saleItems, salesTotal } from "@/lib/sales";
import { festivalSelectionLocked, sessionRef } from "@/state";
import { useMutation } from "@/composables/useMutation";
import type { Stock } from "@/types/stock";
import { flushPromises } from "@vue/test-utils";

afterEach(() => vi.unstubAllGlobals());
describe("APIとデータ取得", () => {
  it("JSONと204を処理し、APIエラー本文を保持する", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response('{"id":"f1"}'))
      .mockResolvedValueOnce(new Response(null, { status: 204 }))
      .mockResolvedValueOnce(
        new Response('{"message":"duplicate"}', { status: 409 }),
      );
    vi.stubGlobal("fetch", fetch);
    expect(await api("/festivals")).toEqual({ id: "f1" });
    expect(await api("/festivals/f1", { method: "DELETE" })).toBeUndefined();
    await expect(api("/posters")).rejects.toThrow("duplicate");
    expect(fetch.mock.calls[0][0]).toBe(`${apiBase}/festivals`);
    expect(jsonBody({ quantity: 2 })).toEqual({
      headers: { "Content-Type": "application/json" },
      body: '{"quantity":2}',
    });
  });
  it("通信失敗を知らせ、AbortErrorはそのまま伝える", async () => {
    vi.stubGlobal(
      "fetch",
      vi
        .fn()
        .mockRejectedValueOnce(new TypeError("offline"))
        .mockRejectedValueOnce(new DOMException("abort", "AbortError")),
    );
    await expect(api("/items")).rejects.toThrow("通信できません");
    await expect(api("/items")).rejects.toMatchObject({ name: "AbortError" });
  });
  it("nullの一覧を空配列へ変換し、画像の相対URLを解決する", () => {
    expect(listOf("items")({ items: null })).toEqual([]);
    expect(imageUrl("/api/v1/images/image1")).toContain(
      "/api/v1/images/image1",
    );
    expect(imageUrl("https://example.com/image.png")).toBe(
      "https://example.com/image.png",
    );
    expect(imageUrl("")).toBe("");
  });
  it("イベント切り替え前の遅い応答で現在の一覧を上書きしない", async () => {
    let resolveFirst!: (value: Response) => void;
    vi.stubGlobal(
      "fetch",
      vi
        .fn()
        .mockImplementationOnce(
          () =>
            new Promise<Response>((resolve) => {
              resolveFirst = resolve;
            }),
        )
        .mockResolvedValueOnce(new Response('{"items":["new"]}')),
    );
    const scope = effectScope();
    const path = ref<string | null>("/first");
    const resource = scope.run(() =>
      useResource(path, listOf<string>("items")),
    )!;
    path.value = "/second";
    await nextTick();
    await flushPromises();
    expect(resource.data.value).toEqual(["new"]);
    resolveFirst(new Response('{"items":["old"]}'));
    await flushPromises();
    expect(resource.data.value).toEqual(["new"]);
    path.value = null;
    await nextTick();
    expect(resource.data.value).toBeNull();
    expect(resource.loading.value).toBe(false);
    scope.stop();
  });
  it("壊れたセッション保存値から復帰する", () => {
    sessionStorage.setItem("test-session", "invalid-json");
    const saved = sessionRef("test-session");
    expect(saved.value).toBe("");
    saved.value = "f1";
    expect(sessionStorage.getItem("test-session")).toBe('"f1"');
  });
});
describe("会計の検証", () => {
  const stocks: Stock[] = [
    {
      id: "s1",
      festival_id: "f1",
      price: 500,
      description: "",
      item: {
        id: "i1",
        name: "商品",
        description: "",
        category: "グッズ",
        image_url: "",
      },
    },
  ];
  it("会計対象のIDと数量を作り、売上金額を合計する", () => {
    expect(saleItems({ s1: 2 }, stocks, "f1")).toEqual([
      { stock_id: "s1", quantity: 2 },
    ]);
    expect(
      salesTotal(
        [{ id: "r1", stock_id: "s1", quantity: 2, created_at: "" }],
        stocks,
      ),
    ).toBe(1000);
  });
  it("空の会計、別イベント、削除済み商品、不正な数量を拒否する", () => {
    expect(() => saleItems({}, stocks, "f1")).toThrow("商品を選択");
    expect(() => saleItems({ s1: 1 }, stocks, "f2")).toThrow("変更");
    expect(() => saleItems({ missing: 1 }, stocks, "f1")).toThrow("変更");
    for (const quantity of [0, -1, 1.5, NaN, Infinity])
      expect(() => saleItems({ s1: quantity }, stocks, "f1")).toThrow("数量");
  });
});

it("並行した保存がすべて終了するまで対象イベントの切り替えを止める", async () => {
  const first = useMutation();
  const second = useMutation();
  let finishFirst!: () => void;
  let finishSecond!: () => void;
  const firstResult = first.run(
    () =>
      new Promise<void>((resolve) => {
        finishFirst = resolve;
      }),
    "保存しました",
  );
  const secondResult = second.run(
    () =>
      new Promise<void>((resolve) => {
        finishSecond = resolve;
      }),
    "保存しました",
  );
  expect(festivalSelectionLocked.value).toBe(true);
  finishFirst();
  await firstResult;
  expect(festivalSelectionLocked.value).toBe(true);
  finishSecond();
  await secondResult;
  expect(festivalSelectionLocked.value).toBe(false);
});
