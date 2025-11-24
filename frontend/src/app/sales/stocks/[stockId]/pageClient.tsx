"use client";

import { useStock } from "@/hooks/stockHook";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "react-toastify";

export default function StockDetailPageClient({ params }: Readonly<{
    params: { stockId: string };
}>) {
    const { stockId } = params;

    const [editMode, setEditMode] = useState(false);
    const [editedDescription, setEditedDescription] = useState("");

    const router = useRouter();

    const { data: stock, error, isLoading, mutate: mutateStock } = useStock(stockId);

    return (
        <div className="max-w-7xl mx-auto p-4">
            {isLoading && <p>Loading...</p>}
            {error && <p>Error loading item data: {error.message}</p>}
            {stock && (
                <div>
                    <h2 className="text-2xl font-bold mb-4">商品詳細</h2>
                    <button
                        onClick={() => {
                            setEditMode(!editMode)
                            setEditedDescription(stock.description);
                        }}
                        className="mb-4 px-4 py-2 bg-blue-500 text-white hover:cursor-pointer"
                    >
                        {editMode ? "キャンセル" : "編集"}
                    </button>

                    {editMode &&
                        <button className="ml-2 px-4 py-2 bg-red-500 text-white hover:cursor-pointer" onClick={async () => {
                            if (!confirm("本当にこの商品を削除しますか？")) {
                                return;
                            }

                            const deleteToastId = toast.loading("商品の削除中...");

                            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/stocks/${stock.id}`, {
                                method: "DELETE",
                            });
                            if (res.ok) {
                                toast.update(deleteToastId, { render: "商品の削除に成功しました", type: "success", isLoading: false, autoClose: 3000 });
                                router.push("/sales/stocks");
                            } else {
                                toast.update(deleteToastId, { render: "商品の削除に失敗しました", type: "error", isLoading: false, autoClose: 3000 });
                            }
                        }}>
                            商品を削除
                        </button>
                    }

                    <div className="mb-4">
                        <h2 className="text-xl font-bold">アイテム：</h2>
                        <p>{stock.item.name}</p>
                    </div>

                    <div className="mb-4">
                        <h2 className="text-xl font-bold">カテゴリ：</h2>
                        <p>{stock.item.category}</p>
                    </div>

                    <div className="mb-4">
                        <h2 className="text-xl font-bold">説明：</h2>
                        {editMode ? (
                            <textarea className="w-full p-2 border border-gray-300" value={editedDescription} onChange={(e) => setEditedDescription(e.target.value)} />
                        ) : (
                            <p>{stock.description}</p>
                        )}
                    </div>

                    <div className="mb-4">
                        <h2 className="text-xl font-bold">価格：</h2>
                        <p>{stock.price} 円</p>
                    </div>

                    {editMode && (
                        <button
                            onClick={async () => {
                                const updateToastId = toast.loading("商品の更新中...");

                                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/stocks/${stock.id}`, {
                                    method: "PUT",
                                    headers: {
                                        "Content-Type": "application/json",
                                    },
                                    body: JSON.stringify({
                                        description: editedDescription,
                                    }),
                                });

                                if (res.ok) {
                                    toast.update(updateToastId, { render: "商品の更新に成功しました", type: "success", isLoading: false, autoClose: 3000 });
                                    mutateStock();
                                    setEditMode(false);
                                }
                                else {
                                    toast.update(updateToastId, { render: "商品の更新に失敗しました", type: "error", isLoading: false, autoClose: 3000 });
                                }
                            }}
                            className="px-4 py-2 bg-green-500 text-white hover:cursor-pointer"
                        >
                            変更を保存
                        </button>
                    )}
                </div>
            )}
        </div>
    )
}