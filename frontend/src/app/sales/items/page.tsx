import { Metadata } from "next";
import ItemPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "商品管理ページ",
    description: "物販対象の商品を管理するページ"
}

export default function ItemsPage() {
    return <ItemPageClient />;
}