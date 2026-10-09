export function categoryTags(category: string): string[] {
    return [...new Set(category.split("/").map((tag) => tag.trim()).filter(Boolean))];
}
