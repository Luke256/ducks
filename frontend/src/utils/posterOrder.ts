import type { Poster } from "@/types/poster";

export function comparePosterNumbers(a: Pick<Poster, "name">, b: Pick<Poster, "name">) {
    const aNum = parseInt(a.name.match(/\d+/)?.[0] || "0", 10);
    const bNum = parseInt(b.name.match(/\d+/)?.[0] || "0", 10);
    return aNum - bNum;
}
