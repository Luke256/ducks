import { describe, expect, it } from "vitest";
import { priceColor } from "@/lib/priceColors";

describe("価格の色分け", () => {
    it.each([
        [0, "#595959"],
        [1, "#946200"],
        [500, "#946200"],
        [501, "#005ac7"],
        [1000, "#005ac7"],
        [1001, "#7e459c"],
        [1500, "#7e459c"],
        [1501, "#7e459c"],
        [1999, "#7e459c"],
        [2000, "#b43b37"],
        [10000, "#b43b37"],
    ] as const)("%s円を対応する色で表示する", (price, color) => {
        expect(priceColor(price)).toBe(color);
    });

    it("価格不明の場合は色を付けない", () => {
        expect(priceColor(undefined)).toBeUndefined();
    });
});
