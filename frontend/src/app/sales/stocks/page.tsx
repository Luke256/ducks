import { Metadata } from "next";
import StocksPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "物販設定",
    description: "イベントごとの物販登録",
}

export default function StocksPage() {
    return <StocksPageClient />;
}