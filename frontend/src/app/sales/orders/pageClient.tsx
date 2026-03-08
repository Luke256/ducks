"use client";

import { useFestivalList } from "@/hooks/festivalHook";
import { useSalesDataList } from "@/hooks/salesHook";
import { useSessionStorage } from "@/hooks/sessStorage";
import { useStockListByFestival } from "@/hooks/stockHook";
import { Festival } from "@/types/festival";
import { SaleRecord } from "@/types/saleRecord";
import { Stock } from "@/types/stock";
import { DeleteTwoTone } from "@mui/icons-material";
import { useState } from "react";
import { toast } from "react-toastify";

export default function OrdersPageClient() {
    const { data, error, isLoading, mutate: mutateSalesData } = useSalesDataList();
    const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
    const [filterCategory, setFilterCategory] = useSessionStorage("stockFilterCategory", "");
    const { data: festivals } = useFestivalList();
    const { data: stocks } = useStockListByFestival(currentFestivalId);
    const [editMode, setEditMode] = useState(false);

    const festivalStocks: { [key: string]: Stock } = {};
    if (stocks) {
        stocks.forEach((stock: Stock) => {
            festivalStocks[stock.id] = stock;
        });
    }

    // collect unique categories
    const categories: string[] = [];
    if (stocks) {
        stocks.forEach((stock: Stock) => {
            if (stock.item.category && !categories.includes(stock.item.category)) {
                categories.push(stock.item.category);
            }
        });
        categories.sort();
    }

    // filter by category
    let filteredData = data || [];
    if (filterCategory) {
        filteredData = filteredData.filter((order: SaleRecord) => {
            const stock = festivalStocks[order.stock_id];
            return stock && stock.item.category === filterCategory;
        });
    }

    // sort by category, name, created_at
    if (filteredData) {
        filteredData.sort((a: SaleRecord, b: SaleRecord) => {
            const stockA = festivalStocks[a.stock_id];
            const stockB = festivalStocks[b.stock_id];

            if (!stockA || !stockB) {
                return 0;
            }

            const dateA = new Date(a.created_at);
            const dateB = new Date(b.created_at);

            if (stockA.item.category < stockB.item.category) return -1;
            if (stockA.item.category > stockB.item.category) return 1;
            
            if (stockA.item.name < stockB.item.name) return -1;
            if (stockA.item.name > stockB.item.name) return 1;

            if (dateA < dateB) return -1;
            if (dateA > dateB) return 1;

            return 0;
        });
    }

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
                <select className="mb-4 ml-4 p-2 border border-gray-300 hover:cursor-pointer" onChange={(e) => {
                    setFilterCategory(e.target.value);
                }} value={filterCategory}>
                    <option value="">全てのカテゴリ</option>
                    {categories.map((category: string) => (
                        <option key={category} value={category}>
                            {category}
                        </option>
                    ))}
                </select>
                <h1 className="text-2xl font-bold mb-4">売上管理</h1>
                <div>
                    {isLoading && <p>Loading...</p>}
                    {error && <p>Error loading sales data.</p>}
                    {filteredData && data && (
                        <div>
                            
                            <button className="mb-4 py-2 px-4 bg-blue-500 text-white hover:bg-blue-600 hover:cursor-pointer" onClick={() => setEditMode(!editMode)}>
                                {editMode ? "終了" : "編集"}
                            </button>
                            <table className="w-full table table-auto">
                                <thead>
                                    <tr>
                                        <th className="text-left p-2 px-4">時刻</th>
                                        <th className="text-left p-2 px-4">カテゴリ</th>
                                        <th className="text-left p-2 px-4">商品名</th>
                                        <th className="text-right p-2 px-4">単価</th>
                                        <th className="text-right p-2 px-4">数量</th>
                                        <th className="text-right p-2 px-4">金額</th>
                                        {editMode && <th>削除</th>}
                                    </tr>
                                </thead>
                                <tbody>
                                    {filteredData.map((order: SaleRecord, index: number) => (
                                        <tr key={order.id} className={"border-t border-gray-300 " + (index % 2 === 0 ? "bg-white" : "bg-gray-100")}>
                                            <td className="text-left p-2 px-4">
                                                {new Date(order.created_at).toLocaleString()}
                                            </td>
                                            <td className="text-left p-2 px-4">
                                                {festivalStocks[order.stock_id]?.item.category || "不明"}
                                            </td>
                                            <td className="text-left p-2 px-4">
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
                                            {editMode && (<td className="text-center p-2 px-4">
                                                <button className="hover:cursor-pointer" onClick={async ()=>{
                                                    if (!confirm("注文を削除しますか？\nこの操作は取り消せません")) {
                                                        return;
                                                    }

                                                    const deleteToastId = toast.loading("注文を削除中...");

                                                    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/sales/${order.id}`, {
                                                        method: "DELETE",
                                                    })

                                                    if (res.ok) {
                                                        toast.update(deleteToastId, { render: "注文を削除しました", type: "success", isLoading: false, autoClose: 3000 });
                                                        mutateSalesData();
                                                    } else {
                                                        toast.update(deleteToastId, { render: "注文の削除に失敗しました", type: "error", isLoading: false, autoClose: 5000 });
                                                    }
                                                }}>
                                                    <DeleteTwoTone color="error" fontSize="small"/>
                                                </button>
                                            </td>)}
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
                                            {filteredData.reduce((total: number, order: SaleRecord) => {
                                                const stock = festivalStocks[order.stock_id];
                                                if (stock) {
                                                    return total + stock.price * order.quantity;
                                                }
                                                return total;
                                            }, 0).toLocaleString()}
                                        </td>
                                        {editMode && <td></td>}
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            </div>
        </main>
    );
}