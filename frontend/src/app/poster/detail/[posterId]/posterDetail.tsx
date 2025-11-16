"use client";

import { useFestival } from "@/hooks/festivalHook";
import { usePoster } from "@/hooks/posterHook";
import { useState } from "react";
import { toast } from 'react-toastify';
import StatusPicker from "../../statusPicker";
import { useRouter } from "next/navigation";
import Image from "next/image";

export default function PosterDetail({ params }: Readonly<{
    params: { posterId: string };
}>) {
    const { posterId } = params;

    const [editMode, setEditMode] = useState(false);
    const [editedName, setEditedName] = useState("");
    const [editedStatus, setEditedStatus] = useState("");
    const router = useRouter();

    const { data: poster, error, isLoading, mutate: mutatePoster } = usePoster(posterId);
    const { data: festival, error: festivalError, isLoading: festivalLoading } = useFestival(poster ? poster.festival_id : "");

    return (
        <div>
            {isLoading && <p>Loading...</p>}
            {error && <p>Error loading poster data: {error.message}</p>}
            {poster && (
                <div>
                    <h2 className="text-2xl font-bold mb-4">ポスター詳細</h2>
                    <button
                        onClick={() => {
                            setEditMode(!editMode)
                            setEditedName(poster.name);
                            setEditedStatus(poster.description);
                        }}
                        className="mb-4 px-4 py-2 bg-blue-500 text-white hover:cursor-pointer"
                    >
                        {editMode ? "キャンセル" : "編集"}
                    </button>

                    {editMode &&
                        <button className="ml-2 px-4 py-2 bg-red-500 text-white hover:cursor-pointer" onClick={async () => {
                            if (!confirm("本当にこのポスターを削除しますか？")) {
                                return;
                            }

                            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters/${poster.id}`, {
                                method: "DELETE",
                            });
                            if (res.ok) {
                                toast.success("ポスターが削除されました");
                                router.push("/poster");
                            } else {
                                toast.error(`ポスターの削除に失敗しました: ${res.statusText}`);
                            }
                        }}>
                            ポスターを削除
                        </button>
                    }

                    <h2 className="text-xl font-bold">ポスター名：</h2>
                    {editMode ? (
                        <div>
                            <input
                                type="text"
                                className="w-full h-10 border border-gray-300 p-2 mb-2"
                                defaultValue={poster.name}
                                onChange={(e) => setEditedName(e.target.value)}
                            ></input>
                        </div>
                    ) : (
                        <div>

                            <h2 className="text-xl">{poster.name}</h2>
                        </div>
                    )
                    }
                    <div>
                        <h3 className="text-xl font-semibold mt-4">イベント</h3>
                        {festivalLoading && <p>Loading ...</p>}
                        {festivalError && <p>Error loading festival data: {festivalError.message}</p>}
                        {festival && (
                            <p>{festival.name}</p>
                        )}
                    </div>
                    <div className="mt-4">
                        <h3 className="text-xl font-semibold">ステータス</h3>
                        <StatusPicker status={poster.status} onChange={async (newStatus: string) => {
                            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters/${poster.id}/status`, {
                                method: "PATCH",
                                body: JSON.stringify({ status: newStatus }),
                                headers: { "Content-Type": "application/json" }
                            }
                            );
                            if (res.ok) {
                                toast.success("ポスターのステータスが更新されました");
                            } else {
                                toast.error(`ポスターのステータスの更新に失敗しました: ${res.statusText}`);
                            }
                            await mutatePoster();
                        }} />
                    </div>
                    <div className="mt-4">
                        <h3 className="text-xl font-semibold">設置場所</h3>
                        {editMode ? (
                            <div>
                                <input
                                    type="text"
                                    className="w-full h-10 border border-gray-300 p-2 mb-2"
                                    defaultValue={poster.description}
                                    onChange={(e) => setEditedStatus(e.target.value)}
                                ></input>
                            </div>
                        ) : (
                            <p>{poster.description}</p>
                        )
                        }
                        {poster.image_url ? (
                            <Image src={poster.image_url} alt={poster.name} className="max-w-xs m-auto" width={400} height={400} />
                        ) : (
                            <p>画像がありません。</p>
                        )}
                    </div>

                    {editMode &&
                        <button
                            className="ml-2 mb-4 px-4 py-2 bg-green-500 text-white hover:cursor-pointer"
                            onClick={
                                async () => {
                                    if (editedName == poster.name && editedStatus == poster.description) {
                                        setEditMode(false);
                                        return;
                                    }
                                    const updateToastId = toast.loading("ポスターを更新中...");
                                    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters/${poster.id}`, {
                                        method: "PUT",
                                        headers: {
                                            "Content-Type": "application/json",
                                        },
                                        body: JSON.stringify({
                                            name: editedName,
                                            description: editedStatus,
                                        }),
                                    });

                                    if (res.ok) {
                                        toast.update(updateToastId, { render: "ポスターが更新されました", type: "success", isLoading: false, autoClose: 3000 });
                                    }
                                    else {
                                        toast.update(updateToastId, { render: `ポスターの更新に失敗しました: ${res.statusText}`, type: "error", isLoading: false, autoClose: 5000 });
                                    }
                                    setEditMode(false);
                                    await mutatePoster();
                                }
                            }
                        >
                            保存
                        </button>

                    }
                </div>
            )}
        </div>
    );
}