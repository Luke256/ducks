import type { Stock } from "@/types/stock";
import type { SaleRecord } from "@/types/saleRecord";
export function saleItems(
    cart: Record<string, number>,
    stocks: Stock[],
    festivalId: string,
) {
    const items = Object.entries(cart).map(([stock_id, quantity]) => {
        const stock = stocks.find((s) => s.id === stock_id);
        if (!stock || !stock.for_sale || stock.festival_id !== festivalId)
            throw new Error(
                "販売商品が変更されています。一覧を更新して選び直してください。",
            );
        if (!Number.isSafeInteger(quantity) || quantity < 1)
            throw new Error("数量は1以上の整数を指定してください。");
        return { stock_id, quantity };
    });
    if (!items.length) throw new Error("商品を選択してください。");
    return items;
}
export function salesTotal(records: SaleRecord[], stocks: Stock[]) {
    const prices = new Map(stocks.map((s) => [s.id, s.price]));
    return records.reduce(
        (total, r) => total + (prices.get(r.stock_id) || 0) * r.quantity,
        0,
    );
}
