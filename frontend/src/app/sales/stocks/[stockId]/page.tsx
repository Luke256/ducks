import { Metadata } from "next";
import StockDetailPageClient from "./pageClient";

export async function generateMetadata({ params }: Readonly<{
    params: Promise<{ stockId: string }>;
}>): Promise<Metadata> {
    const { stockId } = await params;

    return {
        title: `アイテム詳細: ${stockId}`,
        description: `アイテム「${stockId}」の詳細ページ`,
    };
}

export default async function ItemDetailPage({ params }: Readonly<{
    params: Promise<{ stockId: string }>;
}>) {
    const { stockId } = await params;

    return <StockDetailPageClient params={{ stockId }} />;
}