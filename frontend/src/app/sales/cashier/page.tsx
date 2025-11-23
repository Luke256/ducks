import { Metadata } from "next";
import CashiersPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "会計処理",
    description: "会計処理ページ",
};

export default function CashiersPage() {
    return <CashiersPageClient />;
}