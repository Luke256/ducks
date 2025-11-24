"use client";

import { useFestivalList } from "@/hooks/festivalHook";
import { useQueryStockItems } from "@/hooks/itemHook";
import { useSessionStorage } from "@/hooks/sessStorage";
import { Festival } from "@/types/festival";
import { StockItem } from "@/types/stockItem";
import { useRouter } from "next/navigation";
import { useRef, useState } from "react";
import { toast } from "react-toastify";


export default function NewStocksPageClient() {
    const { data: festivals } = useFestivalList();
    const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
    const [selectedCategory, setSelectedCategory] = useState("");
    const { data: stockItems } = useQueryStockItems();
    const submitButton = useRef<HTMLButtonElement>(null);

    const router = useRouter();

    const submitHandler = async (e: React.FormEvent<HTMLFormElement>) => {
        if (submitButton.current) {
            submitButton.current.disabled = true;
        }
        e.preventDefault();

        const formElement = e.currentTarget;
        const formData = new FormData(formElement);
        const itemId = formData.get("item_id") as string;
        const description = formData.get("description") as string;
        const price = formData.get("price") as string;

        const uploadToasId = toast.loading("物販アイテムを登録中...");

        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals/${currentFestivalId}/stocks`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                item_id: itemId,
                price: Number(price),
                description: description,
            }),
        });

        if (res.ok) {
            toast.update(uploadToasId, { render: "物販アイテムが登録されました", type: "success", isLoading: false, autoClose: 3000 });
            router.push("/sales/stocks");
            return;
        }

        toast.update(uploadToasId, { render: `物販アイテムの登録に失敗しました: ${res.statusText}`, type: "error", isLoading: false, autoClose: 5000 });
        if (submitButton.current) {
            submitButton.current.disabled = false;
        }
    }

    const categories: string[] = Array.from(new Set(stockItems?.map((item: StockItem) => item.category)));

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">新規物販アイテムの登録</h1>

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

                <form onSubmit={submitHandler}>
                    <div className="mb-4">
                        <label className="block mb-2">アイテム選択</label>
                        <div className="flex">
                            <select name="category" className="p-2 border border-gray-300" onChange={(e) => setSelectedCategory(e.target.value)} value={selectedCategory}>
                                <option value="">カテゴリ</option>
                                {categories && categories.map((category: string) => (
                                    <option key={category} value={category}>
                                        {category}
                                    </option>
                                ))}
                            </select>
                            <select name="item_id" className="ml-2 p-2 border border-gray-300 w-full" required onChange={(e) => {
                                const selectedItem = stockItems?.find((item: StockItem) => item.id === e.target.value);
                                if (selectedItem) {
                                    setSelectedCategory(selectedItem.category);
                                }
                            }}>
                                <option value="">アイテムを選択</option>
                                {stockItems && 
                                    stockItems.filter((item: StockItem) => item.category === selectedCategory || selectedCategory === "").map((item: StockItem) => (
                                        <option key={item.id} value={item.id}>
                                            {item.name} ({item.category})
                                        </option>
                                    ))
                                }
                            </select>
                        </div>
                    </div>
                    <div className="mb-4">
                        <label className="block mb-2">説明</label>
                        <input name="description" type="text" className="p-2 border border-gray-300 w-full" />
                    </div>
                    <div className="mb-4">
                        <label className="block mb-2">価格 (円)</label>
                        <input name="price" type="number" className="p-2 border border-gray-300 w-full" required />
                    </div>
                    <button ref={submitButton} type="submit" className="px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 hover:cursor-pointer">
                        登録
                    </button>
                </form>
            </div>
        </main>
    );
}