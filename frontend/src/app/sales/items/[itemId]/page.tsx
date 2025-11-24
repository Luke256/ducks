import { Metadata } from "next";
import ItemDetailPageClient from "./pageClient";

export async function generateMetadata({ params }: Readonly<{
    params: Promise<{ itemId: string }>;
}>): Promise<Metadata> {
    const { itemId } = await params;

    return {
        title: `アイテム詳細: ${itemId}`,
        description: `アイテム「${itemId}」の詳細ページ`,
    };
}

export default async function ItemDetailPage({ params }: Readonly<{
    params: Promise<{ itemId: string }>;
}>) {
    const { itemId } = await params;

    return <ItemDetailPageClient params={{ itemId }} />;
}