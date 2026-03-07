"use client";

import { resizeImage } from "@/utils/resizeImage";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { toast } from "react-toastify";

export default function NewItemPageClient() {
    const [previewSrc, setPreviewSrc] = useState<string | null>(null);
    const submitButton = useRef<HTMLButtonElement>(null);

    const router = useRouter();

    useEffect(() => {
        return () => {
            if (previewSrc) {
                URL.revokeObjectURL(previewSrc);
            }
        };
    }, [previewSrc]);

    const submitHandler = async (event: React.FormEvent<HTMLFormElement>) => {
        if (submitButton.current) {
            submitButton.current.disabled = true;
        }
        event.preventDefault();
        const formElement = event.currentTarget;
        const formData = new FormData(formElement);
        const name = formData.get("name") as string;
        const description = formData.get("description") as string;
        const category = formData.get("category") as string;
        const imageEntry = formData.get("image");

        if (!(imageEntry instanceof File)) {
            alert("画像ファイルが取得できませんでした");
            return;
        }

        const uploadToastId = toast.loading("画像を圧縮中...");
        let resizedImage: File;
        try {
            resizedImage = await resizeImage(imageEntry);
        } catch (e) {
            toast.update(uploadToastId, { render: `画像の圧縮に失敗しました: ${e}`, type: "error", isLoading: false, autoClose: 5000 });
            if (submitButton.current) {
                submitButton.current.disabled = false;
            }
            return;
        }

        const form = new FormData();
        form.append("name", name);
        form.append("description", description);
        form.append("category", category);
        form.append("image", resizedImage);

        toast.update(uploadToastId, { render: "アイテムを作成中...", isLoading: true });
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/items`, {
            method: "POST",
            body: form,
        });

        if (res.ok) {
            toast.update(uploadToastId, { render: "アイテムが作成されました", type: "success", isLoading: false, autoClose: 3000 });
            router.push("/sales/items");
            return;
        }

        toast.update(uploadToastId, { render: `アイテムの作成に失敗しました: ${res.statusText}`, type: "error", isLoading: false, autoClose: 5000 });
        if (submitButton.current) {
            submitButton.current.disabled = false;
        }
    }

    return (
        <main>
            <div className="max-w-7xl mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">新規アイテム登録</h1>
                
                <form onSubmit={submitHandler}>
                    <label className="block font-semibold mb-2">アイテム名</label>
                    <input type="text" name="name" placeholder="アイテム名" className="w-full p-2 mb-4 border border-gray-300" required />

                    <label className="block font-semibold mb-2">説明</label>
                    <textarea name="description" placeholder="説明" className="w-full p-2 mb-4 border border-gray-300"></textarea>

                    <label className="block font-semibold mb-2">カテゴリ</label>
                    <input type="text" name="category" placeholder="カテゴリ" className="w-full p-2 mb-4 border border-gray-300" required />

                    <label className="block font-semibold mb-2">画像</label>
                    <input type="file" name="image" accept="image/*" className="w-full p-2 mb-4 border border-gray-300" required 
                        onChange={(e) => {
                            const file = e.target.files?.[0];
                            if (file) {
                                const objectUrl = URL.createObjectURL(file);
                                setPreviewSrc((prev) => {
                                    if (prev) {
                                        URL.revokeObjectURL(prev);
                                    }
                                    return objectUrl;
                                });
                            } else {
                                setPreviewSrc((prev) => {
                                    if (prev) {
                                        URL.revokeObjectURL(prev);
                                    }
                                    return null;
                                });
                            }
                        }}
                    />
                    {previewSrc && (
                        <div className="mb-4">
                            <p className="font-semibold mb-2">画像プレビュー:</p>
                            <Image src={previewSrc} alt="画像プレビュー" className="max-w-xs border border-gray-300" width={200} height={200} />
                        </div>
                    )}

                    <button ref={submitButton} type="submit" className="px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 hover:cursor-pointer">
                        アイテム作成
                    </button>
                </form>
            </div>
        </main>
    );
}