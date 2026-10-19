// 開催日はタイムゾーン変換せず、Goのゼロ値を未設定として扱う。
export function festivalDate(value: string | null | undefined) {
    const date = value?.slice(0, 10) ?? "";
    return date === "0001-01-01" ? "" : date;
}

export function formatFestivalDate(value: string | null | undefined) {
    const date = festivalDate(value);
    if (!date) return "";
    const [year, month, day] = date.split("-");
    return `${year}年${Number(month)}月${Number(day)}日`;
}
