"use client";

import { useFestivalList } from "@/hooks/festivalHook";
import { useSalesDataList } from "@/hooks/salesHook";
import { useSessionStorage } from "@/hooks/sessStorage";
import { useStockListByFestival } from "@/hooks/stockHook";
import { Festival } from "@/types/festival";
import { SaleRecord } from "@/types/saleRecord";
import { Stock } from "@/types/stock";
import { useState } from "react";

export default function OrdersPageClient() {
    const { data, error, isLoading } = useSalesDataList();
    const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
    const { data: festivals } = useFestivalList();
    const { data: stocks, error: stockError, isLoading: stockIsLoading } = useStockListByFestival(currentFestivalId);

    const festivalStocks: { [key: string]: Stock } = {};
    if (stocks) {
        stocks.forEach((stock: Stock) => {
            festivalStocks[stock.id] = stock;
        });
    }

    console.log("sales data:", data);
    console.log("festival stocks:", festivalStocks);

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
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
                <h1 className="text-2xl font-bold mb-4">売上管理</h1>
                <div>
                    {isLoading && <p>Loading...</p>}
                    {error && <p>Error loading sales data.</p>}
                    {data && (
                        <table className="w-full table table-auto">
                            <thead>
                                <tr>
                                    <th className="text-center p-2 px-4">時刻</th>
                                    <th className="text-center p-2 px-4">カテゴリ</th>
                                    <th className="text-center p-2 px-4">商品名</th>
                                    <th className="text-right p-2 px-4">単価</th>
                                    <th className="text-right p-2 px-4">数量</th>
                                    <th className="text-right p-2 px-4">金額</th>
                                </tr>
                            </thead>
                            <tbody>
                                {data.map((order: SaleRecord, index: number) => (
                                    <tr key={order.id} className={"border-t border-gray-300 " + (index % 2 === 0 ? "bg-white" : "bg-gray-100")}>
                                        <td className="text-center p-2 px-4">
                                            {new Date(order.created_at).toLocaleString()}
                                        </td>
                                        <td className="text-center p-2 px-4">
                                            {festivalStocks[order.stock_id]?.item.category || "不明"}
                                        </td>
                                        <td className="text-center p-2 px-4">
                                            {festivalStocks[order.stock_id]?.item.name || "不明"}
                                        </td>
                                        <td className="text-right p-2 px-4">
                                            {festivalStocks[order.stock_id]?.price.toLocaleString() || "不明"}
                                        </td>
                                        <td className="text-right p-2 px-4">{order.quantity}</td>
                                        <td className="text-right p-2 px-4">
                                            {festivalStocks[order.stock_id]
                                                ? (festivalStocks[order.stock_id].price * order.quantity).toLocaleString()
                                                : "不明"}
                                        </td>
                                    </tr>
                                ))}
                                <tr className="font-bold bg-green-100 border-t border-gray-300">
                                    <td className="p-2 px-4">
                                        合計
                                    </td>
                                    <td></td>
                                    <td></td>
                                    <td></td>
                                    <td></td>
                                    <td className="p-2 px-4 text-right">
                                        {data.reduce((total: number, order: SaleRecord) => {
                                            const stock = festivalStocks[order.stock_id];
                                            if (stock) {
                                                return total + stock.price * order.quantity;
                                            }
                                            return total;
                                        }, 0).toLocaleString()}
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    )}
                </div>
            </div>
        </main>
    );
}