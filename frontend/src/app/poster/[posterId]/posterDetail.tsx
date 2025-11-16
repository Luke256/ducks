"use client";

import { useFestival } from "@/hooks/festivalHook";
import { usePoster } from "@/hooks/posterHook";
import { PosterStatusLabels } from "@/types/poster";
import { useState } from "react";
import StatusPicker from "../statusPicker";

export default function PosterDetail({ params }: Readonly<{
    params: { posterId: string };
}>) {
    const { posterId } = params;

    const [editMode, setEditMode] = useState(false);
    const [editedName, setEditedName] = useState("");
    const [editedStatus, setEditedStatus] = useState("");

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
                        <button
                            className="ml-2 mb-4 px-4 py-2 bg-green-500 text-white hover:cursor-pointer"
                            onClick={
                                async (_) => {
                                    if (editedName == poster.name && editedStatus == poster.description) {
                                        setEditMode(false);
                                        return;
                                    }
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

                                    if (!res.ok) {
                                        console.error("Failed to update poster");
                                    }
                                    setEditMode(false);
                                    await mutatePoster();
                                }
                            }
                        >
                            保存
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
                            await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters/${poster.id}/status`, {
                                method: "PATCH",
                                body: JSON.stringify({ status: newStatus }),
                                headers: { "Content-Type": "application/json" }
                            }
                            );
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
                            <img src={poster.image_url} alt={poster.name} className="max-w-xs" />
                        ) : (
                            <p>画像がありません。</p>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}