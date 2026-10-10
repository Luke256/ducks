import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { mount, flushPromises, type VueWrapper } from "@vue/test-utils";
import { nextTick } from "vue";
import App from "@/App.vue";
import router from "@/router";
import { apiBase, imageUrl } from "@/lib/api";
import {
    currentFestivalId,
    festivals,
    notices,
    pendingMutations,
    stockFilterCategory,
} from "@/state";
import { resizeImage } from "@/utils/resizeImage";
import type { VisitorCount } from "@/types/visitorCount";
import type { Poster } from "@/types/poster";
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
    for_sale: true,
    item,
};
const poster: Poster = {
    id: "p1",
    name: "講義棟 01",
    description: "正面入口",
    status: "uncollected",
    image: [
        { id: "image1", url: "/api/v1/images/image1" },
        { id: "image2", url: "/api/v1/images/image2" },
    ],
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
let posterImageNumber = 0;
function fixture(input: string, options: RequestInit = {}) {
    const path = input.slice(apiBase.length);
    requests.push({ path, options });
    if (path === fail)
        return Promise.resolve(new Response("test failure", { status: 500 }));
    if (path === "/posters/p1" && options.method === "PUT")
        Object.assign(poster, JSON.parse(options.body as string));
    if (path === "/stocks/s1" && options.method === "PUT")
        Object.assign(stock, JSON.parse(options.body as string));
    if (path === "/festivals/f1/stocks" && options.method === "POST")
        Object.assign(stock, JSON.parse(options.body as string));
    if (path === "/posters/p1/images" && options.method === "PATCH") {
        const body = options.body as FormData;
        const deleted = body.getAll("delete_image_ids");
        poster.image = [
            ...poster.image.filter((image) => !deleted.includes(image.id)),
            ...body.getAll("image").map(() => {
                const id = `uploaded-${++posterImageNumber}`;
                return { id, url: `/api/v1/images/${id}` };
            }),
        ];
        return Promise.resolve(Response.json({ image: poster.image }));
    }
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
        "/festivals/f1/stocks?only_for_sale=true": { stocks: stock.for_sale ? [{ ...stock, item: { ...item, image_url: itemImageUrl } }] : [] },
        "/festivals/f2/stocks?only_for_sale=true": { stocks: [] },
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
async function upload(files = [new File(["image"], "photo.png", { type: "image/png" })]) {
    const input = wrapper.find("input[type=file]");
    Object.defineProperty(input.element, "files", {
        configurable: true,
        value: files,
    });
    await input.trigger("change");
    return files[0];
}
beforeEach(() => {
    requests.length = 0;
    fail = "";
    stock.price = 500;
    stock.description = "限定商品";
    stock.for_sale = true;
    item.category = "グッズ";
    itemImageUrl = item.image_url;
    poster.name = "講義棟 01";
    poster.description = "正面入口";
    poster.image = [
        { id: "image1", url: "/api/v1/images/image1" },
        { id: "image2", url: "/api/v1/images/image2" },
    ];
    posterImageNumber = 0;
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
    it("複数日の来場者数を1つのグラフに重ねて表示する", async () => {
        visitorCounts.push({ festival_id: "f1", bucket_start: "2026-10-09T10:10:00+09:00", count: 4 });
        await open("/visitors");
        expect(wrapper.findAll(".visitor-chart")).toHaveLength(1);
        expect(wrapper.findAll(".visitor-chart polyline")).toHaveLength(2);
        expect(wrapper.find(".visitor-chart-legend").text()).toContain("2026/10/09");
        expect(wrapper.find(".visitor-chart-legend").text()).toContain("2026/10/08");
        await wrapper.findAll(".visitor-chart-hover")[1].trigger("pointerenter");
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/09：4 人");
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：5 人");
    });
    it("詳細を12件ずつ切り替え、グラフは全件から作り、更新とイベント切り替えでページを補正する", async () => {
        visitorCounts = Array.from({ length: 25 }, (_, i) => ({
            festival_id: "f1",
            bucket_start: new Date(Date.parse("2026-10-08T10:00:00+09:00") + i * 600000).toISOString(),
            count: 1,
        }));
        await open("/visitors");
        expect(wrapper.findAll("tbody tr")).toHaveLength(12);
        expect(wrapper.find("tbody tr").text()).toContain("14:00");
        expect(wrapper.find(".visitor-chart polyline").attributes("points")!.split(" ")).toHaveLength(25);
        expect(button("前へ").attributes("disabled")).toBeDefined();
        await button("次へ").trigger("click");
        expect(wrapper.findAll("tbody tr")).toHaveLength(12);
        expect(wrapper.find("tbody tr").text()).toContain("12:00");
        await button("次へ").trigger("click");
        expect(wrapper.findAll("tbody tr")).toHaveLength(1);
        expect(wrapper.find("tbody tr").text()).toContain("10:00");
        expect(button("次へ").attributes("disabled")).toBeDefined();
        await button("前へ").trigger("click");
        expect(wrapper.find(".visitor-pagination").text()).toContain("2 / 3");
        visitorCounts = visitorCounts.slice(0, 1);
        await button("更新").trigger("click");
        await flushPromises();
        expect(wrapper.find(".visitor-pagination").text()).toContain("1 / 1");
        expect(wrapper.findAll("tbody tr")).toHaveLength(1);
        await wrapper.find(".festival-picker select").setValue("f2");
        await flushPromises();
        expect(wrapper.find(".visitor-chart").exists()).toBe(false);
        expect(wrapper.find(".visitor-pagination").exists()).toBe(false);
    });
    it("0人だけの日はグラフを表示せず、詳細には記録を残す", async () => {
        visitorCounts.forEach((bucket) => { bucket.count = 0; });
        await open("/visitors");
        expect(wrapper.find("svg[role=img]").exists()).toBe(false);
        expect(wrapper.findAll("tbody tr")).toHaveLength(2);
        expect(wrapper.text()).toContain("1人以上の記録がないため");
    });
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
            ).toEqual({ item_id: "i1", price, description: "補足", for_sale: true });
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
    it("販売商品を非表示で登録し、falseを真偽値で送信する", async () => {
        await open("/sales/stocks/new");
        const status = wrapper.find('input[type="checkbox"][name="for_sale"]');
        expect((status.element as HTMLInputElement).checked).toBe(true);
        await wrapper.find("form select[required]").setValue("i1");
        await status.setValue(false);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(JSON.parse(mutations("/festivals/f1/stocks")[0].options.body as string)).toEqual({
            item_id: "i1", price: 500, description: "", for_sale: false,
        });
        expect(wrapper.find(".detail-grid p > .badge").text()).toBe("レジ非表示");
    });
    it.each([true, false])("販売状態%sを保って説明を編集する", async (forSale) => {
        stock.for_sale = forSale;
        await open("/sales/stocks/s1");
        await button("販売情報を編集").trigger("click");
        expect((wrapper.find('input[name="for_sale"]').element as HTMLInputElement).checked).toBe(forSale);
        await wrapper.find("textarea").setValue("更新した説明");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(JSON.parse(mutations("/stocks/s1")[0].options.body as string)).toEqual({
            description: "更新した説明", for_sale: forSale,
        });
        expect(wrapper.find(".detail-grid p > .badge").exists()).toBe(!forSale);
    });
    it.each([true, false])("販売状態を%sへ変更し、保存後の詳細に反映する", async (forSale) => {
        stock.for_sale = !forSale;
        await open("/sales/stocks/s1");
        await button("販売情報を編集").trigger("click");
        await wrapper.find('input[name="for_sale"]').setValue(forSale);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(JSON.parse(mutations("/stocks/s1")[0].options.body as string).for_sale).toBe(forSale);
        expect(wrapper.find(".detail-grid p > .badge").exists()).toBe(!forSale);
    });
    it("編集のキャンセル後に元の販売状態で編集を再開する", async () => {
        stock.for_sale = false;
        await open("/sales/stocks/s1");
        await button("販売情報を編集").trigger("click");
        await wrapper.find('input[name="for_sale"]').setValue(true);
        await button("編集をキャンセル").trigger("click");
        await button("販売情報を編集").trigger("click");
        expect((wrapper.find('input[name="for_sale"]').element as HTMLInputElement).checked).toBe(false);
        expect(mutations("/stocks/s1")).toHaveLength(0);
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
                "販売情報を編集",
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
describe("カテゴリタグ", () => {
    beforeEach(() => {
        item.category = " グッズ /　音楽 / 限定 /音楽 ";
        const otherItem = { ...item, id: "i2", name: "別商品", category: "グ" };
        const otherStock = { ...stock, id: "s2", item: otherItem };
        const taggedData: Record<string, unknown> = {
            "/items": { items: [item, otherItem] },
            "/festivals/f1/stocks": { stocks: [stock, otherStock] },
            "/festivals/f1/stocks?only_for_sale=true": { stocks: [stock, otherStock] },
            "/sales?festival_id=f1": { sales: [...records, { ...records[0], id: "r3", stock_id: "s2" }] },
        };
        vi.mocked(fetch).mockImplementation((input, options) => {
            const path = String(input).slice(apiBase.length);
            if ((!options?.method || options.method === "GET") && taggedData[path]) {
                requests.push({ path, options: options || {} });
                return Promise.resolve(Response.json(taggedData[path]));
            }
            return fixture(String(input), options);
        });
    });
    it.each([
        ["/sales/items", ".product-card"],
        ["/sales/stocks", "tbody tr"],
        ["/sales/orders", "tbody tr"],
        ["/sales/cashier", ".cashier-product"],
        ["/sales/stocks/new", "form select[required] option[value]"],
    ])("%sで各タグを選択でき、部分一致するタグとは区別する", async (path, selector) => {
        await open(path);
        const filter = wrapper.find(path.endsWith("/new") ? "form select" : ".filter-field select");
        expect(filter.findAll("option").map((option) => option.text())).toEqual([
            "すべてのカテゴリ", "グ", "グッズ", "限定", "音楽",
        ]);
        for (const tag of ["グッズ", "音楽", "限定"]) {
            await filter.setValue(tag);
            const text = wrapper.findAll(selector).map((row) => row.text()).join(" ");
            expect(text).toContain(item.name);
            expect(text).not.toContain("別商品");
        }
        await filter.setValue("グ");
        const text = wrapper.findAll(selector).map((row) => row.text()).join(" ");
        expect(text).toContain("別商品");
        expect(text).not.toContain(item.name);
    });
    it("販売商品の選択は一致する別タグへの切り替えで保持し、不一致なら解除する", async () => {
        await open("/sales/stocks/new");
        const filter = wrapper.find("form select");
        const selected = wrapper.find("form select[required]");
        await selected.setValue("i1");
        await filter.setValue("音楽");
        expect((selected.element as HTMLSelectElement).value).toBe("i1");
        await filter.setValue("限定");
        expect((selected.element as HTMLSelectElement).value).toBe("i1");
        await filter.setValue("グ");
        expect((selected.element as HTMLSelectElement).value).toBe("");
    });
    it.each(["/sales/items", "/sales/stocks", "/sales/orders", "/sales/items/i1", "/sales/stocks/s1"])("%sで個別のタグを表示し、最初のタグだけを強調する", async (path) => {
        await open(path);
        const badges = wrapper.find(".category-tags").findAll(".badge");
        expect(badges.map((badge) => badge.text())).toEqual(["グッズ", "音楽", "限定"]);
        expect(badges.map((badge) => badge.classes().includes("primary"))).toEqual([true, false, false]);
        expect(badges[0].attributes("title")).toBe("主タグ");
    });
    it("レジはメイン一覧の下にサブカテゴリ別一覧を表示し、数量を共有して一度だけ会計する", async () => {
        await open("/sales/cashier");
        const cards = wrapper.findAll(".main-category-list .cashier-product").filter((card) => card.text().includes(item.name));
        expect(cards).toHaveLength(1);
        expect(wrapper.findAll(".main-category-list .category-heading").map((heading) => heading.element.firstChild?.textContent?.trim())).toEqual(["グ", "グッズ"]);
        const lists = wrapper.findAll(".main-category-list, .subcategory-list");
        expect(lists.map((list) => list.classes()[0])).toEqual(["main-category-list", "subcategory-list"]);
        expect(wrapper.find(".subcategory-list").text()).not.toContain("サブカテゴリ");
        expect(wrapper.findAll(".subcategory-list .category-heading").map((heading) => heading.element.firstChild?.textContent?.trim())).toEqual(["音楽", "限定"]);
        expect(wrapper.findAll(".subcategory-list .cashier-product")).toHaveLength(2);
        expect(wrapper.find(".subcategory-list").text()).not.toContain("別商品");
        await cards[0].trigger("click");
        await wrapper.find(".filter-field select").setValue("音楽");
        expect(wrapper.findAll(".main-category-list .category-section")).toHaveLength(1);
        expect(wrapper.find(".main-category-list .category-heading").element.firstChild?.textContent?.trim()).toBe("グッズ");
        expect(wrapper.findAll(".cashier-product")).toHaveLength(3);
        expect(wrapper.findAll(".quantity-badge").map((badge) => badge.text())).toEqual(["1", "1", "1"]);
        await wrapper.find(".subcategory-list .cashier-product").trigger("click");
        expect(wrapper.findAll(".receipt-line")).toHaveLength(1);
        expect(wrapper.findAll(".quantity-badge").map((badge) => badge.text())).toEqual(["2", "2", "2"]);
        expect(wrapper.find(".receipt-total").text()).toContain("1,000 円");
        await button("会計を確定する").trigger("click");
        await flushPromises();
        expect(JSON.parse(mutations("/sales")[0].options.body as string)).toEqual({ items: [{ stock_id: "s1", quantity: 2 }] });
    });
    it("レジの検索とカテゴリ絞り込みは両一覧に反映され、該当するサブカテゴリがなければ区切りを表示しない", async () => {
        await open("/sales/cashier");
        const search = wrapper.find('input[type="search"]');
        await search.setValue("限定商品");
        expect(wrapper.findAll(".subcategory-list .cashier-product")).toHaveLength(2);
        await search.setValue("別商品");
        expect(wrapper.findAll(".main-category-list .cashier-product")).toHaveLength(1);
        expect(wrapper.find(".subcategory-list").exists()).toBe(false);
        await search.setValue("");
        await wrapper.find(".filter-field select").setValue("グ");
        expect(wrapper.find(".subcategory-list").exists()).toBe(false);
        await search.setValue("該当なし");
        expect(wrapper.findAll(".cashier-product")).toHaveLength(0);
        expect(wrapper.find(".main-category-list").exists()).toBe(false);
    });
    it("区切り文字しかないカテゴリでは商品を登録しない", async () => {
        await open("/sales/items/new");
        const inputs = wrapper.findAll("form input:not([type=file])");
        await inputs[0].setValue("新商品");
        await inputs[1].setValue(" /　/ ");
        await upload();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/items")).toHaveLength(0);
        expect(wrapper.find("form .error").text()).toContain("商品名とカテゴリを入力してください");
    });
});
describe("商品画像の変更", () => {
    it("画像の選択と取り消しをポスターと同じデザインで行い、取り消すと現在の画像に戻す", async () => {
        await open("/sales/items/i1");
        await button("編集する").trigger("click");
        expect(wrapper.find(".image-add-button").text()).toBe("＋画像を変更");
        expect(wrapper.find("input[type=file]").attributes("aria-label")).toBe("画像を変更");
        expect(wrapper.find(".image-selection .image-preview").attributes("src")).toContain("/images/image1");
        expect(wrapper.find(".image-selection button").exists()).toBe(false);
        await upload();
        const removeImage = wrapper.find('button[aria-label="追加する商品画像の選択を取り消す"]');
        expect(removeImage.text()).toBe("取り消し");
        expect(removeImage.classes()).toContain("secondary");
        vi.mocked(window.confirm).mockReturnValueOnce(false);
        await removeImage.trigger("click");
        expect(window.confirm).toHaveBeenCalledWith("「photo.png」の選択を取り消しますか？");
        expect(wrapper.find(".image-preview").attributes("src")).toBe("blob:preview");
        expect(URL.revokeObjectURL).not.toHaveBeenCalled();
        await removeImage.trigger("click");
        expect(wrapper.find(".image-selection .image-preview").attributes("src")).toContain("/images/image1");
        expect(URL.revokeObjectURL).toHaveBeenCalledWith("blob:preview");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/items/i1/image")).toHaveLength(0);
    });
    it("新規商品では選択画像の取り消し後に画像を必須とし、同じ画像を再選択できる", async () => {
        await open("/sales/items/new");
        const inputs = wrapper.findAll("form input:not([type=file])");
        await inputs[0].setValue("新商品");
        await inputs[1].setValue("グッズ");
        expect(wrapper.find(".image-add-button").text()).toBe("＋画像を追加");
        const file = await upload();
        expect(wrapper.find("input[type=file]").attributes("required")).toBeUndefined();
        await wrapper.find('button[aria-label="追加する商品画像の選択を取り消す"]').trigger("click");
        expect(wrapper.find(".image-selection").exists()).toBe(false);
        expect(wrapper.find("input[type=file]").attributes("required")).toBeDefined();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/items")).toHaveLength(0);
        await upload([file]);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect((mutations("/items")[0].options.body as FormData).get("image")).toBe(file);
    });
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

describe("ポスター画像", () => {
    it.each([
        ["画像なし", []],
        ["1枚", ["/api/v1/images/image1"]],
        ["複数枚", ["/api/v1/images/image1", "https://example.com/image2.jpg"]],
    ] as [string, string[]][])("%sの配列を受け取り、一覧と詳細に表示する", async (_, urls) => {
        poster.image = urls.map((url, index) => ({ id: `image${index + 1}`, url }));
        await open("/poster");
        const thumbnails = wrapper.findAll("tbody .thumbnail");
        expect(thumbnails).toHaveLength(urls.length ? 1 : 0);
        if (urls.length) expect(thumbnails[0].attributes("src")).toBe(imageUrl(urls[0]));

        await router.push("/poster/detail/p1");
        await flushPromises();
        expect(wrapper.findAll(".detail-image").map((img) => img.attributes("src"))).toEqual(urls.map(imageUrl));
        const links = wrapper.findAll(".poster-images a");
        expect(links.map((link) => link.attributes("href"))).toEqual(urls.map(imageUrl));
        for (const link of links) {
            expect(link.attributes("target")).toBe("_blank");
            expect(link.attributes("rel")).toBe("noopener");
        }
    });

    it("複数の写真を圧縮して同じimageフィールドで登録する", async () => {
        await open("/poster/new");
        await wrapper.find("form input:not([type])").setValue("ポスターA");
        await wrapper.find("textarea").setValue("正面入口");
        expect(wrapper.find("input[type=file]").attributes("multiple")).toBeDefined();
        expect(wrapper.find(".image-add-button").text()).toBe("＋写真を追加");
        expect(wrapper.find("input[type=file]").attributes("aria-label")).toBe("写真を追加");
        const files = [
            new File(["first"], "first.png", { type: "image/png" }),
            new File(["second"], "second.png", { type: "image/png" }),
        ];
        const compressed = files.map((file) => new File(["compressed"], `${file.name}.webp`, { type: "image/webp" }));
        vi.mocked(resizeImage).mockResolvedValueOnce(compressed[0]).mockResolvedValueOnce(compressed[1]);
        await upload(files);
        expect(wrapper.findAll(".image-selection .image-preview")).toHaveLength(2);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const request = mutations("/posters")[0];
        expect(request.options.headers).toBeUndefined();
        const uploaded = (request.options.body as FormData).getAll("image");
        expect(uploaded).toHaveLength(2);
        for (let i = 0; i < compressed.length; i++) expect(uploaded[i]).toBe(compressed[i]);
        expect(resizeImage).toHaveBeenNthCalledWith(1, files[0]);
        expect(resizeImage).toHaveBeenNthCalledWith(2, files[1]);
        expect(router.currentRoute.value.path).toBe("/poster");
        expect(URL.revokeObjectURL).toHaveBeenCalled();
    });

    it.each([0, 11])("登録時に写真が%s枚の場合は送信しない", async (count) => {
        await open("/poster/new");
        await wrapper.find("form input:not([type])").setValue("ポスターA");
        await wrapper.find("textarea").setValue("正面入口");
        if (count) await upload(Array.from({ length: count }, (_, i) => new File(["image"], `${i}.png`, { type: "image/png" })));
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/posters")).toHaveLength(0);
        expect(resizeImage).not.toHaveBeenCalled();
        expect(wrapper.find("form .error").exists()).toBe(true);
    });

    it("選択した写真を追加・削除でき、上限の10枚を登録できる", async () => {
        await open("/poster/new");
        await wrapper.find("form input:not([type])").setValue("ポスターA");
        await wrapper.find("textarea").setValue("正面入口");
        const removed = await upload();
        await wrapper.find('.image-selection button').trigger("click");
        expect(wrapper.find("input[type=file]").attributes("required")).toBeDefined();
        const files = Array.from({ length: 10 }, (_, i) => new File(["image"], `${i}.png`, { type: "image/png" }));
        await upload(files.slice(0, 5));
        await upload(files.slice(5));
        expect(wrapper.findAll(".image-selection .image-preview")).toHaveLength(10);
        expect(wrapper.find("input[type=file]").attributes("required")).toBeUndefined();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect((mutations("/posters")[0].options.body as FormData).getAll("image")).toEqual(files);
        expect(vi.mocked(resizeImage).mock.calls.some(([file]) => file === removed)).toBe(false);
    });

    it("既存写真は削除対象の表示を切り替え、追加写真は確認後に一覧から取り除く", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        const files = [
            new File(["first"], "first.png", { type: "image/png" }),
            new File(["second"], "second.png", { type: "image/png" }),
        ];
        await upload(files);
        expect(wrapper.findAll("form .image-selection")).toHaveLength(1);
        expect(wrapper.findAll("form .image-selection figure")).toHaveLength(4);
        const controls = wrapper.findAll(".image-selection button");
        expect(controls.map((control) => control.text())).toEqual(Array(4).fill("削除"));
        expect(controls.every((control) => control.classes().includes("secondary"))).toBe(true);
        const existing = wrapper.find('button[aria-label="登録済みの写真1を削除する"]');
        await existing.trigger("click");
        expect(wrapper.findAll("form .image-selection figure")).toHaveLength(4);
        expect(wrapper.findAll(".marked-for-deletion")).toHaveLength(1);
        expect(wrapper.text()).not.toContain("削除予定");
        expect(existing.text()).toBe("元に戻す");
        await existing.trigger("click");
        expect(wrapper.find(".marked-for-deletion").exists()).toBe(false);
        expect(existing.text()).toBe("削除");
        await existing.trigger("click");
        const cancelAddition = wrapper.findAll("form .image-selection figure")[2].find("button");
        vi.mocked(window.confirm).mockReturnValueOnce(false);
        await cancelAddition.trigger("click");
        expect(window.confirm).toHaveBeenCalledWith("「first.png」を削除しますか？");
        expect(wrapper.findAll("form .image-selection figure")).toHaveLength(4);
        expect(URL.revokeObjectURL).not.toHaveBeenCalled();
        await cancelAddition.trigger("click");
        const photos = wrapper.findAll("form .image-selection figure");
        expect(photos).toHaveLength(3);
        expect(photos[2].text()).toContain("second.png");
        expect(photos[2].find("button").text()).toBe("削除");
        expect(photos[2].classes()).not.toContain("marked-for-deletion");
        expect(URL.revokeObjectURL).toHaveBeenCalled();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const body = mutations("/posters/p1/images")[0].options.body as FormData;
        expect(body.getAll("delete_image_ids")).toEqual(["image1"]);
        expect(body.getAll("image")).toHaveLength(1);
        expect(body.get("image")).toBe(files[1]);
        expect(vi.mocked(resizeImage).mock.calls.some(([file]) => file === files[0])).toBe(false);
    });

    it("名前と設置場所だけを編集すると写真を変更しない", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await wrapper.find("form input:not([type])").setValue("変更したポスター");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(JSON.parse(mutations("/posters/p1")[0].options.body as string)).toEqual({ name: "変更したポスター", description: "正面入口" });
        expect(mutations("/posters/p1/images")).toHaveLength(0);
        expect(wrapper.findAll(".detail-image")).toHaveLength(2);
    });

    it("写真だけを複数追加し、既存の写真を保持する", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        const files = Array.from({ length: 8 }, (_, i) => new File(["image"], `${i}.png`, { type: "image/png" }));
        await upload(files);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const request = mutations("/posters/p1/images")[0];
        const body = request.options.body as FormData;
        expect(request.options.method).toBe("PATCH");
        expect(request.options.headers).toBeUndefined();
        expect(body.getAll("image")).toEqual(files);
        expect(body.getAll("delete_image_ids")).toEqual([]);
        expect(mutations("/posters/p1")).toHaveLength(0);
        const urls = wrapper.findAll(".detail-image").map((img) => img.attributes("src"));
        expect(urls).toHaveLength(10);
        expect(urls.slice(0, 2)).toEqual([imageUrl("/api/v1/images/image1"), imageUrl("/api/v1/images/image2")]);
        expect(urls[2]).toContain("/images/uploaded-1");
    });

    it("指定した写真だけを削除し、残りの写真を保持する", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await wrapper.find('button[aria-label="登録済みの写真1を削除する"]').trigger("click");
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const body = mutations("/posters/p1/images")[0].options.body as FormData;
        expect(body.getAll("delete_image_ids")).toEqual(["image1"]);
        expect(body.getAll("image")).toEqual([]);
        expect(resizeImage).not.toHaveBeenCalled();
        expect(wrapper.findAll(".detail-image").map((img) => img.attributes("src"))).toEqual([imageUrl("/api/v1/images/image2")]);
    });

    it("全写真の削除と新しい写真の追加を同時に送って差し替える", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        for (const control of wrapper.findAll('form button[aria-label^="登録済みの写真"]')) await control.trigger("click");
        const file = await upload();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        const body = mutations("/posters/p1/images")[0].options.body as FormData;
        expect(body.getAll("delete_image_ids")).toEqual(["image1", "image2"]);
        expect(body.getAll("image")).toEqual([file]);
        expect(wrapper.findAll(".detail-image")).toHaveLength(1);
        expect(wrapper.find(".detail-image").attributes("src")).toContain("/images/uploaded-1");
    });

    it.each([1, 2])("残っている%s枚の写真をすべて削除する変更は送信しない", async (count) => {
        poster.image = poster.image.slice(0, count);
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await wrapper.find("form input:not([type])").setValue("変更したポスター");
        for (const control of wrapper.findAll('form button[aria-label^="登録済みの写真"]')) await control.trigger("click");
        expect(button("保存する").attributes("disabled")).toBeDefined();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/posters/p1")).toHaveLength(0);
        expect(mutations("/posters/p1/images")).toHaveLength(0);
        expect(wrapper.find("form").text()).toContain("写真を1枚以上残す");
    });

    it("編集後の写真が11枚になる追加は送信しない", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await upload(Array.from({ length: 9 }, (_, i) => new File(["image"], `${i}.png`, { type: "image/png" })));
        expect(button("保存する").attributes("disabled")).toBeDefined();
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/posters/p1/images")).toHaveLength(0);
        expect(resizeImage).not.toHaveBeenCalled();
        expect(wrapper.find("form").text()).toContain("写真は10枚まで");
    });

    it("キャンセルすると写真の追加・削除指定を破棄する", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await wrapper.find('button[aria-label="登録済みの写真1を削除する"]').trigger("click");
        await upload();
        await button("編集をキャンセル").trigger("click");
        expect(URL.revokeObjectURL).toHaveBeenCalledWith("blob:preview");
        await button("編集する").trigger("click");
        expect(wrapper.find('form button[aria-label^="登録済みの写真"]').text()).toBe("削除");
        expect(wrapper.findAll("form .image-preview")).toHaveLength(2);
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/posters/p1/images")).toHaveLength(0);
    });

    it("圧縮失敗時はポスター情報も送信せず、選択を保持する", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await wrapper.find("form input:not([type])").setValue("変更したポスター");
        await upload();
        vi.mocked(resizeImage).mockRejectedValueOnce(new Error("画像の圧縮に失敗しました"));
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/posters/p1")).toHaveLength(0);
        expect(mutations("/posters/p1/images")).toHaveLength(0);
        expect(wrapper.find("form .image-preview[src='blob:preview']").exists()).toBe(true);
        expect(wrapper.find("form").text()).toContain("画像の圧縮に失敗しました");
    });

    it("写真の保存失敗時に選択を保持し、保存済みの名前を再送せず再試行する", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await wrapper.find("form input:not([type])").setValue("変更したポスター");
        await wrapper.find('button[aria-label="登録済みの写真1を削除する"]').trigger("click");
        const file = await upload();
        fail = "/posters/p1/images";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(wrapper.find('form button[aria-label^="登録済みの写真"]').text()).toBe("元に戻す");
        expect(wrapper.find(".marked-for-deletion").exists()).toBe(true);
        expect(wrapper.find("form .image-preview[src='blob:preview']").exists()).toBe(true);
        expect(wrapper.find("form").text()).toContain("ポスター名・設置場所は保存済みです");
        fail = "";
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(mutations("/posters/p1")).toHaveLength(1);
        expect(mutations("/posters/p1/images")).toHaveLength(2);
        const body = mutations("/posters/p1/images")[1].options.body as FormData;
        expect(body.getAll("image")).toEqual([file]);
        expect(body.getAll("delete_image_ids")).toEqual(["image1"]);
        expect(wrapper.findAll(".detail-image").map((img) => img.attributes("src"))).toEqual([imageUrl("/api/v1/images/image2"), imageUrl("/api/v1/images/uploaded-1")]);
    });

    it("写真の保存中は二重送信とキャンセルを防ぐ", async () => {
        await open("/poster/detail/p1");
        await button("編集する").trigger("click");
        await upload();
        let resolve!: (response: Response) => void;
        vi.mocked(fetch).mockImplementation((input, options) =>
            String(input).endsWith("/posters/p1/images") && options?.method === "PATCH"
                ? new Promise<Response>((done) => {
                    requests.push({ path: "/posters/p1/images", options });
                    resolve = done;
                })
                : fixture(String(input), options),
        );
        await wrapper.find("form").trigger("submit");
        await flushPromises();
        expect(wrapper.find("fieldset").attributes("disabled")).toBeDefined();
        expect(button("編集をキャンセル").attributes("disabled")).toBeDefined();
        await wrapper.find("form").trigger("submit");
        expect(mutations("/posters/p1/images")).toHaveLength(1);
        resolve(Response.json({ image: poster.image }));
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
    it("レジは販売中だけを要求し、非表示の商品を表示しない", async () => {
        stock.for_sale = false;
        await open("/sales/cashier");
        expect(requests.some((r) => r.path === "/festivals/f1/stocks?only_for_sale=true")).toBe(true);
        expect(wrapper.find(".cashier-product").exists()).toBe(false);
        expect(button("会計を確定する").attributes("disabled")).toBeDefined();
    });
    it("選択後に非表示になった商品は更新後に会計できず、明細から削除できる", async () => {
        await open("/sales/cashier");
        await wrapper.find(".cashier-product").trigger("click");
        stock.for_sale = false;
        await button("更新").trigger("click");
        await flushPromises();
        expect(wrapper.find(".cashier-product").exists()).toBe(false);
        expect(wrapper.find(".receipt-line").text()).toContain("販売商品を再確認してください");
        expect(wrapper.findAll(".quantity-control button")[1].attributes("disabled")).toBeDefined();
        expect(button("会計を確定する").attributes("disabled")).toBeDefined();
        expect(mutations("/sales")).toHaveLength(0);
        await wrapper.find(".quantity-control button").trigger("click");
        expect(wrapper.find(".receipt-line").exists()).toBe(false);
    });
    it("管理一覧と売上履歴は非表示の商品も取得し、過去の売上金額を維持する", async () => {
        stock.for_sale = false;
        await open("/sales/stocks");
        expect(wrapper.findAll("tbody tr")).toHaveLength(1);
        expect(wrapper.find("tbody tr").text()).toContain("レジ非表示");
        expect(wrapper.find("tbody tr").text()).toContain(item.name);
        await router.push("/sales/orders");
        await flushPromises();
        expect(wrapper.findAll("tbody tr")).toHaveLength(2);
        expect(wrapper.find("tbody tr").text()).toContain(item.name);
        expect(wrapper.find("tbody tr").text()).toContain("レジ非表示");
        expect(wrapper.find(".summary strong").text()).toBe("1,500 円");
        expect(requests.filter((r) => r.path.includes("/stocks")).every((r) => r.path === "/festivals/f1/stocks")).toBe(true);
    });
    it.each([
        ["/sales/cashier", ".product-price"],
        ["/sales/stocks", "tbody tr td.numeric"],
        ["/sales/stocks/s1", ".price-large"],
        ["/sales/orders", "tbody tr td.numeric"],
    ])("%sの単価に価格帯の色を反映する", async (path, selector) => {
        await open(path);
        expect((wrapper.find(selector).element as HTMLElement).style.color).toBe("rgb(148, 98, 0)");
    });

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
        expect(requests.some((r) => r.path === "/festivals/f2/stocks?only_for_sale=true")).toBe(true);
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
