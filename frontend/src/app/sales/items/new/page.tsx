import { Metadata } from "next";
import NewItemPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "新規アイテム登録",
    description: "新しいアイテムを登録します。",
};

export default function NewItemPage() {
    return <NewItemPageClient />;
}