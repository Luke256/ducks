import type { VisitorCount } from "@/types/visitorCount";

export interface VisitorDay {
    date: string;
    buckets: { start: number; count: number }[];
}

export function dailyVisitorCounts(counts: VisitorCount[]): VisitorDay[] {
    const days = new Map<string, Map<number, number>>();
    for (const bucket of counts) {
        const start = Date.parse(bucket.bucket_start);
        const date = new Date(start).toLocaleDateString("ja-JP", {
            timeZone: "Asia/Tokyo", year: "numeric", month: "2-digit", day: "2-digit",
        });
        if (!days.has(date)) days.set(date, new Map());
        days.get(date)!.set(start, bucket.count);
    }
    const result: VisitorDay[] = [];
    for (const [date, records] of days) {
        const occupied = [...records.keys()].filter((start) => records.get(start)! >= 1);
        if (!occupied.length) continue;
        const first = Math.min(...occupied);
        const last = Math.max(...occupied);
        const buckets: VisitorDay["buckets"] = [];
        for (let start = first; start <= last; start += 10 * 60 * 1000) {
            buckets.push({ start, count: records.get(start) ?? 0 });
        }
        result.push({ date, buckets });
    }
    return result.sort((a, b) => b.date.localeCompare(a.date));
}
