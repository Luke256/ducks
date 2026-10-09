import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { mount, flushPromises, type VueWrapper } from "@vue/test-utils";
import { nextTick } from "vue";
import App from "@/App.vue";
import router from "@/router";
import { apiBase } from "@/lib/api";
import {
    currentFestivalId,
    festivals,
    notices,
    pendingMutations,
    stockFilterCategory,
} from "@/state";
import { resizeImage } from "@/utils/resizeImage";
import type { VisitorCount } from "@/types/visitorCount";
vi.mock("@/utils/resizeImage", () => ({
    resizeImage: vi.fn(async (file: File) => file),
}));

const event = { id: "f1", name: "工大祭 2026", description: "学園祭の運営" };
const item = {
    id: "i1",
    name: "アクリルキーホルダー",
    category: "グッズ",
    description: "オリジナル商品",
    image_url: "/api/v1/images/image1",
};
const stock = {
    id: "s1",
    festival_id: "f1",
    price: 500,
    description: "限定商品",
    item,
};
const poster = {
    id: "p1",
    name: "講義棟 01",
    description: "正面入口",
    status: "uncollected",
    image_url: "/api/v1/images/image1",
    festival: event,
};
const records = [
    { id: "r1", stock_id: "s1", quantity: 1, created_at: "2026-10-08T01:00:00Z" },
    { id: "r2", stock_id: "s1", quantity: 2, created_at: "2026-10-08T02:00:00Z" },
];
let wrapper: VueWrapper;
const requests: { path: string; options: RequestInit }[] = [];
let fail = "";
let visitorCounts: VisitorCount[];
let itemImageUrl: string;
function fixture(input: string, options: RequestInit = {}) {
    const path = input.slice(apiBase.length);
    requests.push({ path, options });
    if (path === fail)
        return Promise.resolve(new Response("test failure", { status: 500 }));
    if (path === "/items/i1/image" && options.method === "PUT")
        itemImageUrl = "/api/v1/images/image2";
    if (path === "/visitors/f1" && options.method === "POST") {
        const { amount } = JSON.parse(options.body as string);
        visitorCounts[0].count = Math.max(0, visitorCounts[0].count + amount);
        return Promise.resolve(new Response(null, { status: 204 }));
    }
    if (options.method && options.method !== "GET") {
        if (options.method === "POST")
            return Promise.resolve(
                Response.json(
                    path === "/items"
                        ? item
                        : path.includes("/stocks")
                            ? stock
                            : path === "/posters"
                                ? poster
                                : path === "/sales"
                                    ? { items: records }
                                    : event,
                    { status: 201 },
                ),
            );
        return Promise.resolve(new Response(null, { status: 204 }));
    }
    const data: Record<string, unknown> = {
        "/festivals": {
            festivals: [event, { id: "f2", name: "別イベント", description: "" }],
        },
        "/festivals/f1": event,
        "/festivals/f2": { id: "f2", name: "別イベント", description: "" },
        "/items": { items: [{ ...item, image_url: itemImageUrl }] },
        "/items/i1": { ...item, image_url: itemImageUrl },
        "/stocks/s1": { ...stock, item: { ...item, image_url: itemImageUrl } },
        "/festivals/f1/stocks": { stocks: [{ ...stock, item: { ...item, image_url: itemImageUrl } }] },
        "/festivals/f2/stocks": { stocks: [] },
        "/festivals/f1/posters": { posters: [poster] },
        "/festivals/f2/posters": { posters: [] },
        "/posters/p1": poster,
        "/sales?festival_id=f1": { sales: records },
        "/sales?festival_id=f2": { sales: [] },
        "/visitors/f1": { festival_id: "f1", counts: visitorCounts },
        "/visitors/f2": { festival_id: "f2", counts: null },
    };
    return Promise.resolve(
        Response.json(data[path] ?? {}, { status: data[path] ? 200 : 404 }),
    );
}
async function open(path: string) {
    await router.push(path);
    wrapper = mount(App, {
        global: { plugins: [router] },
        attachTo: document.body,
    });
    await vi.waitFor(() => expect(wrapper.find("h1").exists()).toBe(true));
    await flushPromises();
}
function button(text: string) {
    return wrapper.findAll("button").find((b) => b.text() === text)!;
}
function mutations(path: string) {
    return requests.filter(
        (r) => r.path === path && r.options.method && r.options.method !== "GET",
    );
}
async function upload() {
    const file = new File(["image"], "photo.png", { type: "image/png" });
    const input = wrapper.find("input[type=file]");
    Object.defineProperty(input.element, "files", {
        configurable: true,
        value: [file],
    });
    await input.trigger("change");
    return file;
}
beforeEach(() => {
    requests.length = 0;
    fail = "";
    stock.price = 500;
    itemImageUrl = item.image_url;
    visitorCounts = [
        { festival_id: "f1", bucket_start: "2026-10-08T10:10:00+09:00", count: 5 },
        { festival_id: "f1", bucket_start: "2026-10-08T10:00:00+09:00", count: 3 },
    ];
    currentFestivalId.value = "f1";
    stockFilterCategory.value = "";
    notices.value = [];
    festivals.value = [];
    vi.stubGlobal("fetch", vi.fn(fixture));
    vi.spyOn(window, "confirm").mockReturnValue(true);
    vi.spyOn(window, "scrollTo").mockImplementation(() => { });
    vi.stubGlobal("URL", URL);
    URL.createObjectURL = vi.fn(() => "blob:preview");
    URL.revokeObjectURL = vi.fn();
});
afterEach(() => {
    wrapper?.unmount();
    document.body.innerHTML = "";
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
});

describe("既存URLと画面", () => {
    it.each([
        ["/", "イベント管理"],
        ["/event/f1", "工大祭 2026"],
        ["/poster", "ポスター管理"],
        ["/poster/new", "ポスターを登録"],
        ["/poster/detail/p1", "講義棟 01"],
        ["/sales", "レジ"],
        ["/sales/cashier", "レジ"],
        ["/sales/orders", "売上履歴"],
        ["/visitors", "来場者数"],
        ["/sales/items", "商品マスター"],
        ["/sales/items/new", "商品を登録"],
        ["/sales/items/i1", item.name],
        ["/sales/stocks", "販売商品"],
        ["/sales/stocks/new", "販売商品を登録"],
        ["/sales/stocks/s1", item.name],
        ["/missing", "ページが見つかりません"],
    ])("%s を直接開ける", async (path, heading) => {
        await open(path);
        expect(wrapper.find("h1").text()).toBe(heading);
    });
});

describe("イベント選択", () => {
    it("一覧で選択したイベントを保持し、運営画面に引き継ぐ", async () => {
        currentFestivalId.value = "";
        await open("/event");
        const cards = wrapper.findAll(".event-card");
        await cards[1].find("button").trigger("click");
        expect(currentFestivalId.value).toBe("f2");
        expect(sessionStorage.getItem("currentFestivalId")).toBe('"f2"');
        expect(router.currentRoute.value.path).toBe("/event");
        expect(cards[1].find(".badge.selected").text()).toBe("選択中");
        expect(cards[1].find("button").attributes("disabled")).toBeDefined();
        expect(cards[0].find(".badge.selected").exists()).toBe(false);
        await cards[0].find("button").trigger("click");
        expect(cards[0].find(".badge.selected").exists()).toBe(true);
        expect(cards[1].find(".badge.selected").exists()).toBe(false);
        await wrapper.find('.sidebar a[href="/visitors"]').trigger("click");
        await flushPromises();
        expect(router.currentRoute.value.path).toBe("/visitors");
        expect(requests.some((r) => r.path === "/visitors/f1")).toBe(true);
        expect(wrapper.find(".festival-picker select").element).toHaveProperty(
            "value",
            "f1",
        );
    });
    it("保存中は一覧からのイベント切り替えも防ぐ", async () => {
        await open("/event");
        const select = wrapper.findAll(".event-card")[1].find("button");
        try {
            pendingMutations.value = 1;
            await nextTick();
            expect(select.attributes("disabled")).toBeDefined();
            await select.trigger("click");
            expect(currentFestivalId.value).toBe("f1");
        } finally {
            pendingMutations.value = 0;
            await nextTick();
        }
        expect(select.attributes("disabled")).toBeUndefined();
        await select.trigger("click");
        expect(currentFestivalId.value).toBe("f2");
    });
});

describe("来場者数", () => {
    it("イベント詳細から対象イベントを選んでカウント画面を開く", async () => {
        currentFestivalId.value = "f2";
        await open("/event/f1");
        await wrapper.find('a[href="/visitors"].panel').trigger("click");
        await flushPromises();
        expect(router.currentRoute.value.path).toBe("/visitors");
        expect(currentFestivalId.value).toBe("f1");
        expect(wrapper.find("h1").text()).toBe("来場者数");
    });
    it("累計と日本時間の履歴を新しい順に表示し、追加・訂正後に再取得する", async () => {
        visitorCounts.reverse();
        await open("/visitors");
        expect(wrapper.find(".summary strong").text()).toBe("8 人");
        expect(wrapper.find("tbody tr").text()).toContain("2026/10/08 10:10");
        await button("＋1人").trigger("click");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("9 人");
        await wrapper.find("input[type=number]").setValue(4);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("13 人");
        await button("−1人（訂正）").trigger("click");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("12 人");
        expect(mutations("/visitors/f1").map((r) => JSON.parse(r.options.body as string))).toEqual([
            { amount: 1 }, { amount: 4 }, { amount: -1 },
        ]);
        expect(requests.filter((r) => r.path === "/visitors/f1" && !r.options.method)).toHaveLength(4);
    });
    it("不正な人数を送信せず、失敗時には入力を保持する", async () => {
        await open("/visitors");
        for (const value of [0, -1, -2, 1.5, ""]) {
            await wrapper.find("input[type=number]").setValue(value);
            await wrapper.find("form").trigger("submit");
            await flushPromises();
        }
        expect(mutations("/visitors/f1")).toHaveLength(0);
        await wrapper.find("input[type=number]").setValue(7);
        fail = "/visitors/f1";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect((wrapper.find("input").element as HTMLInputElement).value).toBe("7");
        expect(wrapper.find(".summary strong").text()).toBe("8 人");
        expect(wrapper.text()).toContain("更新して人数を確認");
        fail = "";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("15 人");
    });
    it("保存中の二重送信と対象イベントの切り替えを防ぐ", async () => {
        await open("/visitors");
        let resolve!: (response: Response) => void;
        vi.mocked(fetch).mockImplementation((input, options) =>
            options?.method === "POST"
                ? new Promise<Response>((done) => {
                    requests.push({ path: "/visitors/f1", options });
                    resolve = done;
                })
                : fixture(String(input), options),
        );
        await button("＋1人").trigger("click");
        await wrapper.find("form").trigger("submit");
        expect(mutations("/visitors/f1")).toHaveLength(1);
        expect(wrapper.find("fieldset").attributes("disabled")).toBeDefined();
        expect(wrapper.find(".festival-picker select").attributes("disabled")).toBeDefined();
        resolve(new Response(null, { status: 204 }));
        await flushPromises();
        expect(wrapper.find("fieldset").attributes("disabled")).toBeUndefined();
    });
    it("イベント切り替えで累計・入力をリセットし、未選択時は取得しない", async () => {
        await open("/visitors");
        await wrapper.find("input[type=number]").setValue(9);
        await wrapper.find(".festival-picker select").setValue("f2");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("0 人");
        expect((wrapper.find("input").element as HTMLInputElement).value).toBe("1");
        expect(wrapper.find("tbody").exists()).toBe(false);
        expect(wrapper.text()).toContain("まだ来場者数が記録されていません");
        const before = requests.length;
        await wrapper.find(".festival-picker select").setValue("");
        await flushPromises();
        expect(requests).toHaveLength(before);
        expect(wrapper.find("form").exists()).toBe(false);
    });
    it("取得失敗時にはカウント操作を止めて再読み込みできる", async () => {
        fail = "/visitors/f1";
        await open("/visitors");
        expect(wrapper.find("fieldset").attributes("disabled")).toBeDefined();
        fail = "";
        await button("再読み込み").trigger("click");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("8 人");
        expect(wrapper.find("fieldset").attributes("disabled")).toBeUndefined();
    });
});
describe("登録・編集", () => {
    it("イベントを作成して詳細を開き、対象イベントを共有する", async () => {
        await open("/event");
        await button("＋ イベントを作成").trigger("click");
        await wrapper.find("form input").setValue("新イベント");
        await wrapper.find("form textarea").setValue("説明");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(
            JSON.parse(mutations("/festivals")[0].options.body as string),
        ).toEqual({ name: "新イベント", description: "説明" });
        expect(router.currentRoute.value.path).toBe("/event/f1");
        expect(currentFestivalId.value).toBe("f1");
    });
    it("編集失敗時に入力を保持し、再試行できる", async () => {
        await open("/event/f1");
        await button("イベントを編集").trigger("click");
        await wrapper.find("form input").setValue("変更したイベント");
        fail = "/festivals/f1";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect((wrapper.find("form input").element as HTMLInputElement).value).toBe(
            "変更したイベント",
        );
        expect(button("保存する").attributes("disabled")).toBeUndefined();
        fail = "";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(wrapper.find("form").exists()).toBe(false);
    });
    it("写真を圧縮してポスターをmultipart形式で登録する", async () => {
        await open("/poster/new");
        await wrapper
            .find("form input[type=text], form input:not([type])")
            .setValue("ポスターA");
        await wrapper.find("textarea").setValue("正面入口");
        const file = await upload();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const body = mutations("/posters")[0].options.body as FormData;
        expect(body.get("festival_id")).toBe("f1");
        expect(body.get("name")).toBe("ポスターA");
        expect(body.get("description")).toBe("正面入口");
        expect(body.get("image")).toBeInstanceOf(File);
        expect(resizeImage).toHaveBeenCalledWith(file);
        expect(router.currentRoute.value.path).toBe("/poster");
    });
    it("商品を写真付きで登録し、商品情報を編集する", async () => {
        await open("/sales/items/new");
        const inputs = wrapper.findAll("form input:not([type=file])");
        await inputs[0].setValue("新商品");
        await inputs[1].setValue("グッズ");
        await upload();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect((mutations("/items")[0].options.body as FormData).get("name")).toBe(
            "新商品",
        );
        expect(router.currentRoute.value.path).toBe("/sales/items/i1");
        await button("編集する").trigger("click");
        await wrapper.find("form input").setValue("改名");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(
            JSON.parse(mutations("/items/i1")[0].options.body as string),
        ).toMatchObject({ name: "改名", category: "グッズ" });
        expect(mutations("/items/i1/image")).toHaveLength(0);
    });
    it.each([750, 0])(
        "販売商品を指定イベントに登録し、価格%s円を数値で送信する",
        async (price) => {
            stock.price = price;
            await open("/sales/stocks/new");
            await wrapper.find("form select[required]").setValue("i1");
            const priceInput = wrapper.find("input[type=number]");
            await priceInput.setValue(price);
            expect((priceInput.element as HTMLInputElement).checkValidity()).toBe(true);
            await wrapper.find("textarea").setValue("補足");
            await wrapper.find("form").trigger("submit");
            await flushPromises();
            expect(
                JSON.parse(mutations("/festivals/f1/stocks")[0].options.body as string),
            ).toEqual({ item_id: "i1", price, description: "補足" });
            expect(router.currentRoute.value.path).toBe("/sales/stocks/s1");
            expect(wrapper.find(".price-large").text()).toBe(`${price} 円`);
        },
    );
    it.each([-1, 0.5, ""])("販売価格%sは登録しない", async (price) => {
        await open("/sales/stocks/new");
        await wrapper.find("form select[required]").setValue("i1");
        const priceInput = wrapper.find("input[type=number]");
        await priceInput.setValue(price);
        expect((priceInput.element as HTMLInputElement).checkValidity()).toBe(false);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/festivals/f1/stocks")).toHaveLength(0);
        expect(wrapper.find("form").text()).toContain("0円以上の整数");
    });
    it("イベント・ポスター・商品・販売商品の削除を維持する", async () => {
        for (const [path, editLabel, deleteLabel, endpoint, destination] of [
            [
                "/event/f1",
                "イベントを編集",
                "イベントを削除",
                "/festivals/f1",
                "/event",
            ],
            [
                "/poster/detail/p1",
                "編集する",
                "ポスターを削除",
                "/posters/p1",
                "/poster",
            ],
            [
                "/sales/items/i1",
                "編集する",
                "商品を削除",
                "/items/i1",
                "/sales/items",
            ],
            [
                "/sales/stocks/s1",
                "説明を編集",
                "販売登録を削除",
                "/stocks/s1",
                "/sales/stocks",
            ],
        ]) {
            await open(path);
            await button(editLabel).trigger("click");
            await button(deleteLabel).trigger("click");
            await flushPromises();
            expect(mutations(endpoint).at(-1)?.options.method).toBe("DELETE");
            expect(router.currentRoute.value.path).toBe(destination);
            wrapper.unmount();
        }
    });
});
describe("商品画像の変更", () => {
    it("差し替え画像を圧縮して送信し、詳細とレジに新画像を表示する", async () => {
        await open("/sales/items/i1");
        await button("編集する").trigger("click");
        expect(wrapper.find("input[type=file]").attributes("required")).toBeUndefined();
        const compressed = new File(["compressed"], "photo.webp", { type: "image/webp" });
        vi.mocked(resizeImage).mockResolvedValueOnce(compressed);
        const file = await upload();
        expect(wrapper.find(".image-preview").attributes("src")).toBe("blob:preview");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const request = mutations("/items/i1/image")[0];
        expect(request.options.method).toBe("PUT");
        expect(request.options.headers).toBeUndefined();
        expect((request.options.body as FormData).get("image")).toBe(compressed);
        expect(resizeImage).toHaveBeenCalledWith(file);
        expect(wrapper.find("form").exists()).toBe(false);
        expect(wrapper.find(".detail-image").attributes("src")).toContain("/images/image2");
        expect(URL.revokeObjectURL).toHaveBeenCalledWith("blob:preview");
        await router.push("/sales/cashier");
        await flushPromises();
        expect(wrapper.find(".product-photo img").attributes("src")).toContain("/images/image2");
    });
    it("編集をキャンセルすると選択画像を破棄し、次の保存で再送しない", async () => {
        await open("/sales/items/i1");
        await button("編集する").trigger("click");
        await upload();
        await button("編集をキャンセル").trigger("click");
        expect(URL.revokeObjectURL).toHaveBeenCalledWith("blob:preview");
        await button("編集する").trigger("click");
        expect(wrapper.find(".image-preview").attributes("src")).toContain("/images/image1");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/items/i1/image")).toHaveLength(0);
    });
    it("画像更新失敗時は入力と選択画像を保持し、再試行できる", async () => {
        await open("/sales/items/i1");
        await button("編集する").trigger("click");
        await wrapper.find("form input").setValue("変更した商品名");
        await upload();
        fail = "/items/i1/image";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect((wrapper.find("form input").element as HTMLInputElement).value).toBe("変更した商品名");
        expect(wrapper.find(".image-preview").attributes("src")).toBe("blob:preview");
        expect(wrapper.find("fieldset").attributes("disabled")).toBeUndefined();
        expect(wrapper.text()).toContain("商品情報だけ保存されている場合があります");
        fail = "";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/items/i1/image")).toHaveLength(2);
        expect(wrapper.find(".detail-image").attributes("src")).toContain("/images/image2");
    });
    it("画像圧縮に失敗した場合は商品情報も送信しない", async () => {
        await open("/sales/items/i1");
        await button("編集する").trigger("click");
        await upload();
        vi.mocked(resizeImage).mockRejectedValueOnce(new Error("画像の圧縮に失敗しました"));
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/items/i1")).toHaveLength(0);
        expect(mutations("/items/i1/image")).toHaveLength(0);
        expect(wrapper.find(".image-preview").attributes("src")).toBe("blob:preview");
        expect(wrapper.find("form").text()).toContain("画像の圧縮に失敗しました");
    });
    it("画像更新中の二重送信と編集キャンセルを防ぐ", async () => {
        await open("/sales/items/i1");
        await button("編集する").trigger("click");
        await upload();
        let resolve!: (response: Response) => void;
        vi.mocked(fetch).mockImplementation((input, options) =>
            String(input).endsWith("/items/i1/image") && options?.method === "PUT"
                ? new Promise<Response>((done) => {
                    requests.push({ path: "/items/i1/image", options });
                    resolve = done;
                })
                : fixture(String(input), options),
        );
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(wrapper.find("fieldset").attributes("disabled")).toBeDefined();
        expect(button("編集をキャンセル").attributes("disabled")).toBeDefined();
        await wrapper.find("form").trigger("submit");
        expect(mutations("/items/i1")).toHaveLength(1);
        expect(mutations("/items/i1/image")).toHaveLength(1);
        resolve(new Response(null, { status: 204 }));
        await flushPromises();
        expect(wrapper.find("form").exists()).toBe(false);
    });
});

describe("ポスター回収", () => {
    it("一覧から状況を変更して件数とフィルタを更新する", async () => {
        await open("/poster");
        const select = wrapper.find("select.status-select");
        await select.setValue("collected");
        await flushPromises();
        expect(
            JSON.parse(mutations("/posters/p1/status")[0].options.body as string),
        ).toEqual({ status: "collected" });
        expect(wrapper.findAll(".summary")[1].text()).toContain("1");
        await wrapper.findAll(".summary")[0].trigger("click");
        expect(wrapper.find("tbody tr").exists()).toBe(false);
    });
    it("状況変更失敗時に選択を戻し、操作を再開できる", async () => {
        await open("/poster");
        fail = "/posters/p1/status";
        const select = wrapper.find("select.status-select");
        await select.setValue("collected");
        await flushPromises();
        expect((select.element as HTMLSelectElement).value).toBe("uncollected");
        expect(select.attributes("disabled")).toBeUndefined();
    });
});
describe("レジ・売上", () => {
    it("0円の商品を表示して会計し、売上履歴にも0円で表示する", async () => {
        stock.price = 0;
        await open("/sales/cashier");
        expect(wrapper.find(".product-price").text()).toBe("0 円");
        await wrapper.find(".cashier-product").trigger("click");
        expect(wrapper.find(".receipt-total").text()).toContain("0 円");
        await wrapper.find(".receipt input").setValue(0);
        expect(button("会計を確定する").attributes("disabled")).toBeUndefined();
        await button("会計を確定する").trigger("click");
        await flushPromises();
        expect(JSON.parse(mutations("/sales")[0].options.body as string)).toEqual({
            items: [{ stock_id: "s1", quantity: 1 }],
        });
        expect(wrapper.find(".receipt-line").exists()).toBe(false);
        await router.push("/sales/orders");
        await flushPromises();
        expect(wrapper.find(".summary strong").text()).toBe("0 円");
        expect(wrapper.find("tbody tr").text()).toContain("0 円");
        expect(wrapper.find("tbody tr").text()).not.toContain("価格不明");
    });
    it("空の会計を防ぎ、数量調整とお釣りを計算し、成功後にクリアする", async () => {
        await open("/sales/cashier");
        expect(button("会計を確定する").attributes("disabled")).toBeDefined();
        await wrapper.find(".cashier-product").trigger("click");
        await wrapper.find(".cashier-product").trigger("click");
        expect(wrapper.find(".receipt-total").text()).toContain("1,000");
        await wrapper.find(".quantity-control button").trigger("click");
        expect(wrapper.find(".receipt-total").text()).toContain("500");
        await wrapper.find(".receipt input").setValue(100);
        expect(button("会計を確定する").attributes("disabled")).toBeDefined();
        await wrapper.find(".receipt input").setValue(1000);
        expect(wrapper.find(".change-amount").text()).toContain("お釣り 500");
        await button("会計を確定する").trigger("click");
        await flushPromises();
        expect(JSON.parse(mutations("/sales")[0].options.body as string)).toEqual({
            items: [{ stock_id: "s1", quantity: 1 }],
        });
        expect(wrapper.find(".receipt-line").exists()).toBe(false);
    });
    it("会計失敗時には選択した商品を保持する", async () => {
        await open("/sales/cashier");
        await wrapper.find(".cashier-product").trigger("click");
        fail = "/sales";
        await button("会計を確定する").trigger("click");
        await flushPromises();
        expect(wrapper.find(".receipt-line").exists()).toBe(true);
        expect(button("会計を確定する").attributes("disabled")).toBeUndefined();
        expect(wrapper.text()).toContain("再送信前に売上履歴");
    });
    it("会計中の二重送信とイベント切り替えを防ぐ", async () => {
        await open("/sales/cashier");
        await wrapper.find(".cashier-product").trigger("click");
        let resolve!: (response: Response) => void;
        vi.mocked(fetch).mockImplementation((input, options) =>
            options?.method === "POST"
                ? new Promise<Response>((done) => {
                    requests.push({ path: "/sales", options });
                    resolve = done;
                })
                : fixture(String(input), options),
        );
        await button("会計を確定する").trigger("click");
        await nextTick();
        expect(
            wrapper.find(".festival-picker select").attributes("disabled"),
        ).toBeDefined();
        await button("会計処理中…").trigger("click");
        expect(mutations("/sales")).toHaveLength(1);
        resolve(Response.json({ items: [] }, { status: 201 }));
        await flushPromises();
        expect(
            wrapper.find(".festival-picker select").attributes("disabled"),
        ).toBeUndefined();
    });
    it("イベント変更で前のイベントの商品をリセットする", async () => {
        await open("/sales/cashier");
        await wrapper.find(".cashier-product").trigger("click");
        await wrapper.find(".festival-picker select").setValue("f2");
        await flushPromises();
        expect(wrapper.find(".receipt-line").exists()).toBe(false);
        expect(requests.some((r) => r.path === "/festivals/f2/stocks")).toBe(true);
        expect(button("会計を確定する").attributes("disabled")).toBeDefined();
    });
    it("売上をイベントで絞り込み、新しい順に表示し、削除後に合計を更新する", async () => {
        await open("/sales/orders");
        expect(requests.some((r) => r.path === "/sales?festival_id=f1")).toBe(true);
        expect(requests.some((r) => r.path === "/sales")).toBe(false);
        expect(wrapper.find(".summary strong").text()).toBe("1,500 円");
        expect(wrapper.find("tbody tr").text()).toContain("1,000 円");
        await button("記録を編集").trigger("click");
        await wrapper.find("tbody button").trigger("click");
        await flushPromises();
        expect(mutations("/sales/r2")[0].options.method).toBe("DELETE");
        expect(wrapper.find(".summary strong").text()).toBe("500 円");
    });
});
