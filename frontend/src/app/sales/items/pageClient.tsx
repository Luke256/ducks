'use client';

import { useQueryStockItems } from "@/hooks/itemHook";
import { LaunchTwoTone } from "@mui/icons-material";
import Link from "next/link";

export default function ItemPageClient() {
    const { data: items, error, isLoading } = useQueryStockItems();

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">アイテム一覧</h1>
                {isLoading && <p>Loading...</p>}
                {error && <p className="text-red-500">Error: {error.message}</p>}

                <Link href="/sales/items/new" className="mb-4 inline-block px-4 py-2 bg-blue-500 text-white hover:bg-blue-600">
                    新規アイテム
                </Link>
                {items && (
                    <div>
                        <table className="w-full">
                            <thead>
                                <tr>
                                    <th>アイテム名</th>
                                    <th>カテゴリ</th>
                                    <th>説明</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                {items.map((item: any) => (
                                    <tr key={item.id}>
                                        <td className="p-2 border-t text-center">{item.name}</td>
                                        <td className="p-2 border-t text-center">{item.category}</td>
                                        <td className="p-2 border-t text-center">{item.description}</td>
                                        <td className="p-2 border-t text-center">
                                            <Link href={`/sales/items/${item.id}`}>
                                                <LaunchTwoTone />
                                            </Link>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </main>
    );
}