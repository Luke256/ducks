import { describe, expect, it, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { dailyVisitorCounts } from "@/lib/visitors";
import VisitorChart from "@/components/VisitorChart.vue";
import { readFileSync } from "node:fs";

function record(bucket_start: string, count: number) {
    return { festival_id: "f1", bucket_start, count };
}

describe("来場者数の日別グラフ", () => {
    it.each([
        ["2026-10-10T10:25:00+09:00", "2026-10-10T10:00:00+09:00"],
        ["2026-10-10T23:55:00+09:00", "2026-10-10T23:30:00+09:00"],
    ])("%sでは現在の10分枠への区間だけを点線にし、終了すると実線に戻す", async (now, firstStart) => {
        vi.useFakeTimers({ toFake: ["Date", "setTimeout", "clearTimeout"] });
        vi.setSystemTime(new Date(now));
        const start = Date.parse(firstStart);
        const counts = [1, 2, 4].map((count, i) => record(new Date(start + i * 600000).toISOString(), count));
        counts.push(...[5, 3, 2].map((count, i) => record(new Date(start - 86400000 + i * 600000).toISOString(), count)));
        const days = dailyVisitorCounts(counts);
        const wrapper = mount(VisitorChart, { props: { days } });
        try {
            const currentLines = wrapper.findAll(`polyline[data-date="${days[0].date}"]`);
            expect(currentLines).toHaveLength(2);
            expect(currentLines[0].attributes("points")).toBe("60,197.5 410,175");
            expect(currentLines[0].attributes("stroke-dasharray")).toBeUndefined();
            expect(currentLines[1].attributes("points")).toBe("410,175 760,130");
            expect(currentLines[1].attributes("stroke-dasharray")).toBe("4 4");
            expect(currentLines[1].attributes("stroke")).toBe(currentLines[0].attributes("stroke"));
            const previousLines = wrapper.findAll(`polyline[data-date="${days[1].date}"]`);
            expect(previousLines).toHaveLength(1);
            expect(previousLines[0].attributes("stroke-dasharray")).toBeUndefined();
            await vi.advanceTimersByTimeAsync(5 * 60000 - 1);
            expect(wrapper.findAll(".visitor-chart-pending-line")).toHaveLength(1);
            await vi.advanceTimersByTimeAsync(1);
            expect(wrapper.find(".visitor-chart-pending-line").exists()).toBe(false);
            expect(wrapper.find(`polyline[data-date="${days[0].date}"]`).attributes("points")).toBe("60,197.5 410,175 760,130");
        } finally {
            wrapper.unmount();
            expect(vi.getTimerCount()).toBe(0);
            vi.useRealTimers();
        }
    });

    it("現在の1枠だけの日は未確定として扱い、人数が更新されても点線の判定を保つ", async () => {
        vi.useFakeTimers({ toFake: ["Date", "setTimeout", "clearTimeout"] });
        vi.setSystemTime(new Date("2026-10-10T10:25:00+09:00"));
        const wrapper = mount(VisitorChart, { props: {
            days: dailyVisitorCounts([record("2026-10-10T10:20:00+09:00", 1)]),
        } });
        try {
            expect(wrapper.findAll("polyline")).toHaveLength(1);
            expect(wrapper.find(".visitor-chart-pending-line").attributes("points")).toBe("410,175");
            await wrapper.setProps({ days: dailyVisitorCounts([
                record("2026-10-10T10:20:00+09:00", 3),
            ]) });
            expect(wrapper.find(".visitor-chart-pending-line").attributes("points")).toBe("410,85");
        } finally {
            wrapper.unmount();
            vi.useRealTimers();
        }
    });

    it("異なる日を同じ時刻・人数の軸に重ね、各日の範囲を保って色を分ける", async () => {
        const wrapper = mount(VisitorChart, { props: {
            days: dailyVisitorCounts([
                record("2026-10-09T10:10:00+09:00", 8),
                record("2026-10-09T10:30:00+09:00", 4),
                record("2026-10-08T10:00:00+09:00", 2),
                record("2026-10-08T10:20:00+09:00", 8),
            ]),
        } });
        const lines = wrapper.findAll("polyline");
        expect(wrapper.findAll('.visitor-chart-scroll svg')).toHaveLength(1);
        expect(lines).toHaveLength(2);
        expect(lines[0].attributes("stroke")).not.toBe(lines[1].attributes("stroke"));
        const coordinates = lines.map((line) => line.attributes("points")!.split(" ").map((point) => point.split(",").map(Number)));
        expect(coordinates[0]).toHaveLength(3);
        expect(coordinates[1]).toHaveLength(3);
        expect(coordinates[0][0][0]).toBeCloseTo(coordinates[1][1][0]);
        expect(coordinates[0][1][0]).toBeCloseTo(coordinates[1][2][0]);
        expect(coordinates[0][0][1]).toBe(coordinates[1][2][1]);
        expect(coordinates[0][0][0]).toBeGreaterThan(60);
        expect(coordinates[1][2][0]).toBeLessThan(760);
        const areas = wrapper.findAll(".visitor-chart-hover");
        await areas[1].trigger("pointerenter");
        const tooltip = wrapper.find('[role="tooltip"]');
        expect(tooltip.text()).toContain("10:10-10:20");
        expect(tooltip.text()).toContain("2026/10/09：8 人");
        expect(tooltip.text()).toContain("2026/10/08：0 人");
        const markers = wrapper.findAll("circle");
        expect(markers).toHaveLength(2);
        markers.forEach((marker, i) => {
            expect(Number(marker.attributes("cx"))).toBeCloseTo(coordinates[i][i][0]);
            expect(Number(marker.attributes("cy"))).toBeCloseTo(coordinates[i][i][1]);
            expect(marker.attributes("fill")).toBe(lines[i].attributes("stroke"));
        });
        await areas[0].trigger("pointerenter");
        expect(tooltip.text()).toContain("2026/10/08：2 人");
        expect(tooltip.text()).not.toContain("2026/10/09");
        expect(wrapper.findAll("circle")).toHaveLength(1);
        expect(wrapper.find("circle").attributes("data-date")).toBe("2026/10/08");
        await areas[3].trigger("pointerenter");
        expect(tooltip.text()).toContain("2026/10/09：4 人");
        expect(tooltip.text()).not.toContain("2026/10/08");
        expect(wrapper.findAll("circle")).toHaveLength(1);
        expect(wrapper.find("circle").attributes("data-date")).toBe("2026/10/09");
        await wrapper.find(".visitor-chart-scroll svg").trigger("pointerleave");
        expect(wrapper.find("circle").exists()).toBe(false);
        wrapper.unmount();
    });

    it("日本時間の同時刻はUTCの日付が異なっても同じ横位置に表示する", () => {
        const wrapper = mount(VisitorChart, { props: {
            days: dailyVisitorCounts([
                record("2026-10-08T15:10:00Z", 2),
                record("2026-10-09T15:10:00Z", 3),
            ]),
        } });
        expect(wrapper.findAll("polyline").map((line) => line.attributes("points")!.split(",")[0])).toEqual(["410", "410"]);
        expect(wrapper.find(".visitor-chart-legend").text()).toContain("2026/10/10（00:10〜00:10）");
        expect(wrapper.find(".visitor-chart-legend").text()).toContain("2026/10/09（00:10〜00:10）");
        wrapper.unmount();
    });

    it("線とホバー領域は塗りつぶさず、人数表示は半透明の白背景と縁取り文字で表示する", async () => {
        const style = document.createElement("style");
        style.textContent = readFileSync("src/styles.css", "utf8");
        document.head.append(style);
        const wrapper = mount(VisitorChart, { attachTo: document.body, props: {
            days: dailyVisitorCounts([
                record("2026-10-08T10:00:00+09:00", 2),
                record("2026-10-08T10:20:00+09:00", 4),
            ]),
        } });
        try {
            expect(getComputedStyle(wrapper.find(".visitor-chart-scroll svg").element).fill).toBe("");
            expect(getComputedStyle(wrapper.find("polyline").element).fill).toBe("none");
            expect(getComputedStyle(wrapper.find(".visitor-chart-hover").element).fill).toBe("transparent");
            expect(getComputedStyle(wrapper.find("text").element).fill).toBe("#52616e");
            await wrapper.find(".visitor-chart-hover").trigger("pointerenter");
            const backgroundStyle = getComputedStyle(wrapper.find(".visitor-chart-tooltip-background").element);
            expect(backgroundStyle.fill).toBe("white");
            expect(Number(backgroundStyle.fillOpacity)).toBeGreaterThan(0);
            expect(Number(backgroundStyle.fillOpacity)).toBeLessThan(1);
            const textStyle = getComputedStyle(wrapper.find(".visitor-chart-tooltip text").element);
            expect(textStyle.fill).toBe("#182b3b");
            expect(textStyle.stroke).toBe("white");
        } finally {
            wrapper.unmount();
            style.remove();
        }
    });

    it("日本時間で日を分け、正の記録の間だけを10分刻みで補完する", () => {
        const days = dailyVisitorCounts([
            record("2026-10-08T16:00:00Z", 0),
            record("2026-10-08T15:20:00Z", 3),
            record("2026-10-08T14:50:00Z", 2),
            record("2026-10-08T15:10:00Z", 0),
            record("2026-10-08T15:00:00Z", 1),
            record("2026-10-08T14:40:00Z", 0),
            record("2026-10-07T00:00:00Z", 0),
        ]);
        expect(days).toEqual([
            { date: "2026/10/09", buckets: [
                { start: Date.parse("2026-10-08T15:00:00Z"), count: 1 },
                { start: Date.parse("2026-10-08T15:10:00Z"), count: 0 },
                { start: Date.parse("2026-10-08T15:20:00Z"), count: 3 },
            ] },
            { date: "2026/10/08", buckets: [
                { start: Date.parse("2026-10-08T14:50:00Z"), count: 2 },
            ] },
        ]);
        const withGap = dailyVisitorCounts([
            record("2026-10-08T10:00:00+09:00", 2),
            record("2026-10-08T10:20:00+09:00", 4),
        ]);
        expect(withGap[0].buckets.map((b) => b.count)).toEqual([2, 0, 4]);
        expect(dailyVisitorCounts([])).toEqual([]);
        expect(dailyVisitorCounts([record("2026-10-08T10:00:00+09:00", 0)])).toEqual([]);
    });

    it("未選択時は点を表示せず、1枠だけの日も時刻を表示し、更新後は線と軸を再計算する", async () => {
        const wrapper = mount(VisitorChart, { props: {
            days: dailyVisitorCounts([record("2026-10-08T10:00:00+09:00", 1)]),
        } });
        expect(wrapper.find(".visitor-chart-scroll svg").attributes("aria-label")).toContain("日別");
        expect(wrapper.find("circle").exists()).toBe(false);
        expect(wrapper.text()).toContain("10:00〜10:00");
        expect(wrapper.find("polyline").attributes("points")).toBe("410,175");
        expect(wrapper.find("polyline").attributes("points")).not.toMatch(/NaN|Infinity/);
        await wrapper.setProps({ days: dailyVisitorCounts([
            record("2026-10-08T10:00:00+09:00", 1),
            record("2026-10-08T10:20:00+09:00", 8),
        ]) });
        expect(wrapper.find("circle").exists()).toBe(false);
        expect(wrapper.find("polyline").attributes("points")).toBe("60,197.5 410,220 760,40");
        wrapper.unmount();
    });

    it("ホバーとキーボードフォーカスで時刻・人数を表示し、離れたら消す", async () => {
        const wrapper = mount(VisitorChart, { props: {
            days: dailyVisitorCounts([
                record("2026-10-08T10:00:00+09:00", 1234),
                record("2026-10-08T10:20:00+09:00", 8),
            ]),
        } });
        const areas = wrapper.findAll(".visitor-chart-hover");
        expect(wrapper.find('[role="tooltip"]').exists()).toBe(false);
        await areas[0].trigger("pointerenter");
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：1,234 人");
        expect(wrapper.find('[role="tooltip"] > text').text()).toBe("10:00-10:10");
        expect(areas[0].attributes("aria-label")).toContain("10:00-10:10");
        expect(wrapper.findAll("circle")).toHaveLength(1);
        expect(wrapper.find("circle").attributes("cx")).toBe("60");
        await areas[1].trigger("pointerenter");
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：0 人");
        expect(wrapper.find("circle").attributes("cx")).toBe("410");
        expect(wrapper.find("circle").attributes("cy")).toBe("220");
        await wrapper.find(".visitor-chart-scroll svg").trigger("pointerleave");
        expect(wrapper.find('[role="tooltip"]').exists()).toBe(false);
        expect(wrapper.find("circle").exists()).toBe(false);
        await areas[2].trigger("focus");
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：8 人");
        expect(wrapper.find('[role="tooltip"] > text').text()).toBe("10:20-10:30");
        expect(wrapper.find("circle").attributes("cx")).toBe("760");
        await areas[2].trigger("blur");
        expect(wrapper.find('[role="tooltip"]').exists()).toBe(false);
        expect(wrapper.find("circle").exists()).toBe(false);
        await wrapper.setProps({ days: dailyVisitorCounts([
            record("2026-10-08T10:00:00+09:00", 2),
        ]) });
        await wrapper.find(".visitor-chart-hover").trigger("pointerenter");
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：2 人");
        expect(wrapper.find("circle").attributes("cx")).toBe("410");
        expect(wrapper.find("circle").attributes("cy")).toBe("130");
        expect(wrapper.find(".visitor-chart-hover").attributes("width")).toBe("700");
        await wrapper.setProps({ days: dailyVisitorCounts([
            record("2026-10-08T10:00:00+09:00", 3),
        ]) });
        expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：3 人");
        expect(wrapper.find("circle").attributes("cy")).toBe("85");
        wrapper.unmount();
    });

    it("拡縮・スクロール後のカーソル位置に表示し、同じ時間枠内でも追従して端では折り返す", async () => {
        const wrapper = mount(VisitorChart, { props: {
            days: dailyVisitorCounts([
                record("2026-10-08T10:00:00+09:00", 2),
                record("2026-10-08T10:20:00+09:00", 8),
            ]),
        } });
        const bounds = vi.spyOn(wrapper.find(".visitor-chart-scroll svg").element, "getBoundingClientRect")
            .mockReturnValue(new DOMRect(100, 50, 400, 140));
        const areas = wrapper.findAll(".visitor-chart-hover");
        try {
            await areas[0].trigger("pointerenter", { clientX: 140, clientY: 100 });
            expect(wrapper.find('[role="tooltip"]').attributes("transform")).toBe("translate(92, 112)");
            expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：2 人");
            await areas[0].trigger("pointermove", { clientX: 150, clientY: 120 });
            expect(wrapper.find('[role="tooltip"]').attributes("transform")).toBe("translate(112, 152)");
            await areas[2].trigger("pointerenter", { clientX: 480, clientY: 160 });
            expect(wrapper.find('[role="tooltip"]').attributes("transform")).toBe("translate(578, 168)");
            expect(wrapper.find('[role="tooltip"]').text()).toContain("2026/10/08：8 人");
            await areas[2].trigger("focus");
            expect(wrapper.find('[role="tooltip"]').attributes("transform")).toBe("translate(578, 52)");
            bounds.mockReturnValue(new DOMRect(-100, -50, 400, 140));
            await areas[0].trigger("pointermove", { clientX: -50, clientY: 20 });
            expect(wrapper.find('[role="tooltip"]').attributes("transform")).toBe("translate(112, 152)");
        } finally {
            bounds.mockRestore();
            wrapper.unmount();
        }
    });
});
