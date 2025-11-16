'use client';

import { useFestivalList } from "@/hooks/festivalHook";
import { useSessionStorage } from "@/hooks/sessStorage";
import { Festival } from "@/types/festival";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import { toast } from "react-toastify";

const NewPosterPageClient = () => {
    const { data: festivals } = useFestivalList();
    const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
    const [previewSrc, setPreviewSrc] = useState<string | null>(null);

    const router = useRouter();

    useEffect(() => {
        return () => {
            if (previewSrc) {
                URL.revokeObjectURL(previewSrc);
            }
        };
    }, [previewSrc]);

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formElement = event.currentTarget;
        const formData = new FormData(formElement);
        const name = formData.get("name") as string;
        const description = formData.get("description") as string;
        const imageEntry = formData.get("image");

        if (!(imageEntry instanceof File)) {
            toast.error("画像ファイルが取得できませんでした");
            return;
        }

        const form = new FormData();
        form.append("name", name);
        form.append("description", description);
        form.append("image", imageEntry);
        form.append("festival_id", currentFestivalId);

        const uploadToastId = toast.loading("ポスターを作成中...");
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters`, {
            method: "POST",
            body: form,
        });

        if (res.ok) {
            toast.update(uploadToastId, { render: "ポスターが作成されました", type: "success", isLoading: false, autoClose: 3000 });
            router.push("/poster");
            return;
        }

        toast.update(uploadToastId, { render: `ポスターの作成に失敗しました: ${res.statusText}`, type: "error", isLoading: false, autoClose: 5000 });
    };

    return (
        <main>
            <h1 className="mb-4 text-2xl font-bold text-black">新規ポスター登録</h1>
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
                <div className="mb-4">
                    <form onSubmit={handleSubmit}>
                        <label className="block font-semibold mb-2">ポスター名</label>
                        <input type="text" name="name" placeholder="ポスター名" required className="mb-2 p-2 border border-gray-300 w-full" />
                        <label className="block font-semibold mb-2">ポスターの場所</label>
                        <textarea name="description" placeholder="ポスターの場所" required className="mb-2 p-2 border border-gray-300 w-full"></textarea>
                        <label className="block font-semibold mb-2">ポスターの場所の写真</label>
                        <label className="block mb-2 text-sm">
                            <ul className="list-disc list-inside text-gray-600">
                                <li>ポスターが写っている</li>
                                <li>場所がわかるよう、周辺の様子が引きで写っている</li>
                            </ul>
                        </label>
                        <input type="file" name="image" accept="image/*" required
                            className="p-2 border border-gray-300 w-full mb-2 hover:cursor-pointer"
                            onChange={
                                (e) => {
                                    const file = e.target.files?.[0];
                                    if (file) {
                                        const objectUrl = URL.createObjectURL(file);
                                        setPreviewSrc((prev) => {
                                            if (prev) {
                                                URL.revokeObjectURL(prev);
                                            }
                                            return objectUrl;
                                        });
                                    }
                                    else {
                                        setPreviewSrc((prev) => {
                                            if (prev) {
                                                URL.revokeObjectURL(prev);
                                            }
                                            return null;
                                        });
                                    }
                                }
                            } />
                        <br />
                        {previewSrc && (
                            <Image
                                src={previewSrc}
                                alt="プレビュー画像"
                                width={400}
                                height={300}
                                className="mb-4 max-h-48 object-contain"
                            />
                        )}
                        <br />
                        <button type="submit" className="px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 hover:cursor-pointer">作成</button>
                    </form>
                </div>
            )}
        </main>
    );
}

export default NewPosterPageClient;