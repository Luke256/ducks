'use client';

import { useQueryStockItems } from "@/hooks/itemHook";
import { StockItem } from "@/types/stockItem";
import { LaunchTwoTone, RefreshTwoTone } from "@mui/icons-material";
import Link from "next/link";
import { useState } from "react";

export default function ItemPageClient() {
    const { data: items, error, isLoading, mutate } = useQueryStockItems();
    const [filterCategory, setFilterCategory] = useState("");

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">アイテム一覧</h1>
                {isLoading && <p>Loading...</p>}
                {error && <p className="text-red-500">Error: {error.message}</p>}

                <Link href="/sales/items/new" className="mb-4 inline-block px-4 py-2 bg-blue-500 text-white hover:bg-blue-600">
                    新規アイテム
                </Link>

                <div className="mb-4">
                    <input type="text" placeholder="カテゴリ" value={filterCategory} onChange={(e) => setFilterCategory(e.target.value)} className="p-2 border border-gray-300 w-full" />
                </div>
                <button className="mb-4 px-4 py-2 bg-gray-500 text-white hover:bg-gray-600 hover:cursor-pointer" onClick={() => mutate(filterCategory)}>
                    <RefreshTwoTone />
                </button>
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
                                {items.map((item: StockItem) => (
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