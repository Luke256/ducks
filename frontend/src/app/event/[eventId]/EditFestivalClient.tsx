"use client";

import { useState } from "react";
import type { Festival } from "@/types/festival";
import { KeyboardArrowDown, KeyboardArrowUp } from "@mui/icons-material";
import { useRouter } from "next/navigation";
import Form from "next/form";

type Props = {
    event: Festival;
    eventId: string;
};

export default function EditFestivalClient({ event, eventId }: Props) {
    const [formOpen, setFormOpen] = useState(false);

    const router = useRouter();

    const submitEditAction = async (formData: FormData) => {
        const name = formData.get("name") as string;
        const description = formData.get("description") as string;

        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals/${eventId}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ name, description }),
            cache: 'no-store',
        });

        router.refresh();
    }

    const deleteFestival = async () => {
        if (!confirm("本当にこのイベントを削除しますか？")) {
            return;
        }

        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals/${eventId}`, {
            method: "DELETE",
            cache: 'no-store',
        });

        router.push("/event");
    }

    return (
        <div className="mt-8 border-t">
            <div className="mt-4">
                <button
                    className="px-4 py-2 font-bold text-blue-950 hover:bg-blue-950/10 bg-blue-950/2 w-full hover:cursor-pointer transition-colors"
                    onClick={() => setFormOpen(!formOpen)}
                >
                    編集
                    {formOpen ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                </button>
            </div>

            {formOpen && (
                <div>
                    <Form action={submitEditAction} className="mt-4">
                        <label htmlFor="name" className="block mb-2 font-bold text-gray-700">イベント名</label>
                        <input type="text" name="name" placeholder="イベント名" defaultValue={event.name} className="w-full mb-4 p-2 border" />
                        <label htmlFor="description" className="block mb-2 font-bold text-gray-700">イベントの概要</label>
                        <textarea name="description" placeholder="イベントの概要" defaultValue={event.description} className="w-full mb-4 p-2 border" />
                        <button type="submit" className="px-4 py-2 font-bold text-white bg-green-600 hover:bg-green-700 w-full transition-colors hover:cursor-pointer">
                            保存
                        </button>
                    </Form>
                    <button className="mt-4 px-4 py-2 font-bold text-white bg-red-600 hover:bg-red-900 w-full transition-colors hover:cursor-pointer" onClick={() => deleteFestival()}>
                        削除
                    </button>
                </div>
            )}
        </div>
    )
}
