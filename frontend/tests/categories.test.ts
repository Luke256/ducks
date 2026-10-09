import { describe, expect, it } from "vitest";
import { categoryTags } from "@/lib/categories";

describe("カテゴリタグ", () => {
    it.each([
        ["グッズ", ["グッズ"]],
        ["  グッズ /　音楽 /限定 // 音楽\t/グッズ\n", ["グッズ", "音楽", "限定"]],
        ["グッズ 音楽,限定", ["グッズ 音楽,限定"]],
        ["音楽 グッズ / 限定,特典", ["音楽 グッズ", "限定,特典"]],
        [" /　/ ", []],
        ["", []],
    ])("%sを空のタグと重複を除いて分割する", (input, expected) => {
        expect(categoryTags(input)).toEqual(expected);
    });
});
