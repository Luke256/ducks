const priceBands = [
    { min: 0, max: 0, color: "#595959" },
    { min: 1, max: 500, color: "#946200" },
    { min: 501, max: 1000, color: "#005ac7" },
    // 1,501〜1,999円もこの価格帯の色に含める。
    { min: 1001, max: 1999, color: "#7e459c" },
    { min: 2000, max: Infinity, color: "#b43b37" },
];

export function priceColor(price: number | undefined) {
    return price === undefined
        ? undefined
        : priceBands.find((band) => price >= band.min && price <= band.max)?.color;
}
