import { Metadata } from "next";
import OrdersPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "売上管理",
    description: "売上管理ページ",
};

export default function CashiersPage() {
    return <OrdersPageClient />;
}