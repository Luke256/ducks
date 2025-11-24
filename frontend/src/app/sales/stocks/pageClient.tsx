"use client";

import { useFestivalList } from "@/hooks/festivalHook";
import { useSessionStorage } from "@/hooks/sessStorage";
import { useStockListByFestival as useStockList } from "@/hooks/stockHook";
import { Festival } from "@/types/festival";
import { Stock } from "@/types/stock";
import { LaunchTwoTone, RefreshTwoTone } from "@mui/icons-material";
import Link from "next/link";

export default function StocksPageClient() {
    const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
    const { data: stocks, mutate: mutateStocks } = useStockList(currentFestivalId);
    const { data: festivals } = useFestivalList();

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">物販の設定</h1>
                <select className="mb-4 p-2 border border-gray-300 hover:cursor-pointer" onChange={(e) => {
                    setCurrentFestivalId(e.target.value);
                }} value={currentFestivalId}>
                    <option value="">イベントを選択</option>
                    {festivals && festivals.map((festival: Festival) => (
                        <option key={festival.id} value={festival.id}>
                            {festival.name}
                        </option>
                    ))}
                </select>

                {currentFestivalId && (
                    <Link href="/sales/stocks/new" className="mb-4 inline-block px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 ml-2">
                        物品の登録
                    </Link>
                )}

                {currentFestivalId && (!stocks || stocks.length === 0) && (
                    <p className="text-gray-500">このイベントにはまだ物販が登録されていません。</p>
                )}

                {stocks && (
                    <div>
                        <button className="mb-4 px-4 py-2 bg-gray-500 text-white hover:bg-gray-600 hover:cursor-pointer" onClick={() => mutateStocks()}>
                            <RefreshTwoTone />
                        </button>
                        <table className="table-auto w-full">
                            <thead>
                                <tr>
                                    <th className="text-left">アイテム名</th>
                                    <th className="text-left">カテゴリ</th>
                                    <th className="text-left">説明</th>
                                    <th className="text-right">価格</th>
                                    <th className="text-center"></th>
                                </tr>
                            </thead>
                            <tbody>
                                {stocks.map((stock: Stock) => (
                                    <tr key={stock.id}>
                                        <td className="p-2 border-t text-left">{stock.item.name}</td>
                                        <td className="p-2 border-t text-left">{stock.item.category}</td>
                                        <td className="p-2 border-t text-left">{stock.description}</td>
                                        <td className="p-2 border-t text-right">{stock.price} 円</td>
                                        <td className="p-2 border-t text-center hover:cursor-pointer transition">
                                            <Link href={`/sales/stocks/${stock.id}`} className="hover:text-blue-500">
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