import { describe, expect, it } from "vitest";
import { festivalDate, formatFestivalDate } from "@/utils/festivalDate";

describe("イベントの開催日", () => {
    it.each([
        ["2026-10-18T00:00:00Z", "2026-10-18"],
        ["2026-10-18T00:00:00+09:00", "2026-10-18"],
        ["2026-10-18", "2026-10-18"],
        ["0001-01-01T00:00:00Z", ""],
        [null, ""],
        [undefined, ""],
        ["", ""],
    ])("%sを時差なく日付入力・表示に使える", (value, expected) => {
        expect(festivalDate(value)).toBe(expected);
    });

    it.each([
        ["2026-10-18T00:00:00Z", "2026年10月18日"],
        ["2026-01-02T00:00:00+09:00", "2026年1月2日"],
        ["2026-01-02", "2026年1月2日"],
        ["0001-01-01T00:00:00Z", ""],
        [null, ""],
        [undefined, ""],
        ["", ""],
    ])("%sを年月日で表示する", (value, expected) => {
        expect(formatFestivalDate(value)).toBe(expected);
    });
});
