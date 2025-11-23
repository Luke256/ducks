import { Metadata } from "next";
import NewStocksPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "物販登録",
    description: "新しい物販アイテムを登録",
}

export default function NewStocksPage() {
    return <NewStocksPageClient />;
}