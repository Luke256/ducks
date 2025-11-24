"use client";

import { useFestivalList } from "@/hooks/festivalHook";
import { useSessionStorage } from "@/hooks/sessStorage";
import { useStockListByFestival } from "@/hooks/stockHook";
import { Festival } from "@/types/festival";
import { Stock } from "@/types/stock";
import { Clear } from "@mui/icons-material";
import Image from "next/image";
import { useRef, useState } from "react";
import { toast } from "react-toastify";

const hsl2hex = (h: number, s: number, l: number): string => {
    s /= 100;
    l /= 100;
    const k = (n: number) => (n + h / 30) % 12;
    const a = s * Math.min(l, 1 - l);
    const f = (n: number) => {
        const color = l - a * Math.max(-1, Math.min(k(n) - 3, Math.min(9 - k(n), 1)));
        return Math.round(255 * color).toString(16).padStart(2, "0");
    };
    return `#${f(0)}${f(8)}${f(4)}`;
}

export default function CashiersPageClient() {
    const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
    const { data: festivals } = useFestivalList();
    const { data: stocks, error: stocksError, isLoading: stocksLoading } = useStockListByFestival(currentFestivalId);
    const [selectedStocks, setSelectedStockIds] = useState<{ [key: string]: number }>({});
    const submitButton = useRef<HTMLButtonElement>(null);

    const displayStocksMap: { [key: string]: Stock[] } = {};
    if (stocks) {
        stocks.forEach((stock: Stock) => {
            const categoryName = stock.item.category;
            if (!displayStocksMap[categoryName]) {
                displayStocksMap[categoryName] = [];
            }
            displayStocksMap[categoryName].push(stock);
        });
    }

    const priceMin = stocks ? Math.min(...stocks.map((stock: Stock) => stock.price)) : 0;
    const priceMax = stocks ? Math.max(...stocks.map((stock: Stock) => stock.price)) : 0;
    const priceColorClass = (price: number) => {
        // HSLの色相を計算（赤から緑へのグラデーション）
        const hue = priceMax === priceMin ? 120 : 120 - ((price - priceMin) / (priceMax - priceMin)) * 120;
        return hsl2hex(hue, 60, 50);
    }

    const saleHandler = async () => {
        submitButton.current?.setAttribute("disabled", "true");

        const items = Object.entries(selectedStocks).map(([stockId, quantity]) => ({
            stock_id: stockId,
            quantity,
        }));

        const saleToastId = toast.loading("会計処理中...");

        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/sales`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ items }),
        });

        if (res.ok) {
            toast.update(saleToastId, { render: "会計処理が完了しました。", type: "success", isLoading: false, autoClose: 3000 });
            setSelectedStockIds({});
        }
        else {
            toast.update(saleToastId, { render: `会計処理中にエラーが発生しました: ${res.statusText}`, type: "error", isLoading: false, autoClose: 3000 });
        }
        submitButton.current?.removeAttribute("disabled");
    }

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">レジ</h1>

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

                {stocksLoading && <p>商品を読み込み中...</p>}
                {stocksError && <p className="text-red-500">商品の取得中にエラーが発生しました。</p>}
                {stocks && (
                    <div>
                        {Object.entries(displayStocksMap).map(([itemName, stocks]) => (
                            <div key={itemName}>
                                <h2 className="text-xl font-semibold">{itemName}</h2>
                                <div className="grid grid-cols-6">
                                    {stocks.map((stock) => (
                                        <div key={stock.id} className="p-2 m-2 text-center hover:cursor-pointer relative" onClick={() => {
                                            // increment
                                            setSelectedStockIds((prev) => {
                                                const newSelected = { ...prev };
                                                if (newSelected[stock.id]) {
                                                    newSelected[stock.id] += 1;
                                                } else {
                                                    newSelected[stock.id] = 1;
                                                }
                                                return newSelected;
                                            });
                                        }}>
                                            {/* 選択中のものは半透明の白い幕を被せる */}
                                            <div className="relative">
                                                <Image src={stock.item.image_url} alt={stock.item.name} width={160} height={160}
                                                    className={`mx-auto shadow-md ${selectedStocks[stock.id] ? "opacity-50" : ""}`} />

                                                {selectedStocks[stock.id] && (
                                                    <div>
                                                        <div className="flex text-green-600 text-5xl font-black absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
                                                            {selectedStocks[stock.id]}
                                                        </div>
                                                        <div className="rounded-full bg-white border border-red-600 p-2 hover:cursor-pointer absolute top-0 right-0 transform translate-x-1/2 -translate-y-1/2" onClick={(e) => {
                                                            e.stopPropagation();
                                                            setSelectedStockIds((prev) => {
                                                                const newSelected = { ...prev };
                                                                delete newSelected[stock.id];
                                                                return newSelected;
                                                            });
                                                        }}>
                                                            <Clear fontSize="medium" className="text-red-600" />
                                                        </div>
                                                    </div>
                                                )}
                                            </div>
                                            <p>{stock.item.name}</p>
                                            <p className={`font-bold`} style={{ color: priceColorClass(stock.price) }}>{stock.price.toLocaleString()} 円</p>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
                <div>
                    <h2 className="text-xl font-semibold mt-4">合計</h2>
                    <p className="text-2xl font-bold">
                        {Object.entries(selectedStocks).reduce((sum, [stockId, quantity]) => {
                            const stock = stocks?.find((s: Stock) => s.id === stockId);
                            if (stock) {
                                return sum + stock.price * quantity;
                            }
                            return sum;
                        }, 0).toLocaleString()} 円
                    </p>
                    <button ref={submitButton} className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 hover:cursor-pointer" onClick={saleHandler}>
                        会計する
                    </button>
                </div>
            </div>
        </main>
    );
}