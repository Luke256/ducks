"use client";

import { useStockItem } from "@/hooks/itemHook";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "react-toastify";

export default function ItemDetailPageClient({ params }: Readonly<{
    params: { itemId: string };
}>) {
    const { itemId } = params;

    const [editMode, setEditMode] = useState(false);
    const [editedName, setEditedName] = useState("");
    const [editedCategory, setEditedCategory] = useState("");
    const [editedDescription, setEditedDescription] = useState("");

    const router = useRouter();

    const { data: item, error, isLoading, mutate: mutateItem } = useStockItem(itemId);

    return (
        <div>
            {isLoading && <p>Loading...</p>}
            {error && <p>Error loading item data: {error.message}</p>}
            {item && (
                <div>
                    <h2 className="text-2xl font-bold mb-4">アイテム詳細</h2>
                    <button
                        onClick={() => {
                            setEditMode(!editMode)
                            setEditedName(item.name);
                            setEditedCategory(item.category);
                            setEditedDescription(item.description);
                        }}
                        className="mb-4 px-4 py-2 bg-blue-500 text-white hover:cursor-pointer"
                    >
                        {editMode ? "キャンセル" : "編集"}
                    </button>

                    {editMode &&
                        <button className="ml-2 px-4 py-2 bg-red-500 text-white hover:cursor-pointer" onClick={async () => {
                            if (!confirm("本当にこのアイテムを削除しますか？")) {
                                return;
                            }

                            const deleteToastId = toast.loading("アイテムの削除中...");

                            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/items/${item.id}`, {
                                method: "DELETE",
                            });
                            if (res.ok) {
                                toast.update(deleteToastId, { render: "アイテムの削除に成功しました", type: "success", isLoading: false, autoClose: 3000 });
                                router.push("/sales/items");
                            } else {
                                toast.update(deleteToastId, { render: "アイテムの削除に失敗しました", type: "error", isLoading: false, autoClose: 3000 });
                            }
                        }}>
                            アイテムを削除
                        </button>
                    }

                    <h2 className="text-xl font-bold">アイテム名：</h2>
                    {editMode ? (
                        <div>
                            <input
                                type="text"
                                value={editedName}
                                onChange={(e) => setEditedName(e.target.value)}
                                className="w-full h-10 border border-gray-300 p-2 mb-2"
                            />
                        </div>
                    ) : (
                        <p>{item.name}</p>
                    )}

                    <h2 className="text-xl font-bold">カテゴリ：</h2>
                    {editMode ? (
                        <div>
                            <input
                                type="text"
                                value={editedCategory}
                                onChange={(e) => setEditedCategory(e.target.value)}
                                className="w-full h-10 border border-gray-300 p-2 mb-2"
                            />
                        </div>
                    ) : (
                        <p>{item.category}</p>
                    )}

                    <h2 className="text-xl font-bold">説明：</h2>
                    {editMode ? (
                        <div>
                            <textarea
                                value={editedDescription}
                                onChange={(e) => setEditedDescription(e.target.value)}
                                className="w-full border border-gray-300 p-2 mb-2"
                            ></textarea>
                        </div>
                    ) : (
                        <p>{item.description}</p>
                    )}
                    {item.image_url && (
                        <div>
                            <h2 className="text-xl font-bold mt-4">画像：</h2>
                            <Image src={item.image_url} alt={item.name} width={400} height={400} />
                        </div>
                    )}

                    {editMode && (
                        <button
                            className="mt-4 px-4 py-2 bg-green-500 text-white hover:cursor-pointer"
                            onClick={async () => {
                                if (
                                    editedName === item.name &&
                                    editedCategory === item.category &&
                                    editedDescription === item.description
                                ) {
                                    setEditMode(false);
                                    return;
                                }

                                const updateToastId = toast.loading("アイテムを更新中...");

                                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/items/${item.id}`, {
                                    method: "PUT",
                                    headers: {
                                        "Content-Type": "application/json",
                                    },
                                    body: JSON.stringify({
                                        name: editedName,
                                        category: editedCategory,
                                        description: editedDescription,
                                    }),
                                });

                                if (res.ok) {
                                    toast.update(updateToastId, { render: "アイテムが更新されました", type: "success", isLoading: false, autoClose: 3000 });
                                }
                                else {
                                    toast.update(updateToastId, { render: `アイテムの更新に失敗しました: ${res.statusText}`, type: "error", isLoading: false, autoClose: 5000 });
                                }
                                setEditMode(false);
                                await mutateItem();
                            }}
                        >
                            保存
                        </button>
                    )}
                </div>
            )}
        </div>
    )
}