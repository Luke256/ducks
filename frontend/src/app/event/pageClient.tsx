'use client';

import { KeyboardArrowDown, KeyboardArrowUp } from "@mui/icons-material";
import { useEffect, useState } from "react";
import Form from "next/form"
import { redirect } from "next/navigation";
import { Festival } from "@/types/festival";
import { useFestivalList } from "@/hooks/festivalHook";

type FestivalCreateFormData = {
    name: string;
    description: string;
};

const fetchEvents = async () => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals`);
    const data = await response.json();
    return data["festivals"];
};

const createEvent = async (formData: FestivalCreateFormData) => {
    await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
        cache: 'no-store',
    });

    return;
}

export default function EventPageClient() {
    const { data: events, error: eventsError, isLoading: eventsLoading, mutate: reloadEvents } = useFestivalList();
    const [formOpen, setFormOpen] = useState(false);


    const submitAction = async (formData: FormData) => {
        const name = formData.get("name") as string;
        const description = formData.get("description") as string;

        await createEvent({ name, description });
        await reloadEvents();
        setFormOpen(false);
    }

    return (
        <main className="max-w-7xl mx-auto">
            <div className="max-w-3xl bg-white md:p-8 mx-auto">
                <h1 className="mb-4 text-2xl font-bold text-black">イベント管理</h1>
                {eventsLoading && <p>Loading events...</p>}
                {eventsError && <p className="text-red-500">Failed to load events: {eventsError.message}</p>}
                {events &&
                    <table className="min-w-full table-auto border-collapse">
                        <thead>
                            <tr>
                                <th className="px-4 py-2 text-left">イベント名</th>
                                <th className="px-4 py-2 text-left">概要</th>
                            </tr>
                        </thead>
                        <tbody>
                            {events.map((event: Festival) => (
                                <tr
                                    key={event.id}
                                    onClick={() => redirect(`/event/${event.id}`)}
                                    className="hover:cursor-pointer hover:bg-gray-100 transition-colors"
                                >
                                    <td className="border-t border-gray-300 px-4 py-2">{event.name}</td>
                                    <td className="border-t border-gray-300 px-4 py-2">{event.description}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                }

                {/* Create Event Form (Hidden with accordion) */}

                <div className="border border-gray-300 mt-8"></div>
                <div className="mt-4">
                    <button
                        className="px-4 py-2 font-bold text-blue-950 hover:bg-blue-950/10 bg-blue-950/5 w-full hover:cursor-pointer transition-colors"
                        onClick={() => setFormOpen(!formOpen)}
                    >
                        新しいイベントを作成
                        {formOpen ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                    </button>
                </div>

                {formOpen && (
                    <div className="mt-8">
                        <Form action={submitAction} className="mt-4">
                            <label htmlFor="name" className="block mb-2 font-bold text-gray-700">イベント名</label>
                            <input type="text" name="name" placeholder="イベント名" className="w-full mb-4 p-2 border border-gray-300 rounded" required />
                            <label htmlFor="description" className="block mb-2 font-bold text-gray-700">イベントの概要</label>
                            <textarea name="description" placeholder="イベントの概要" className="w-full mb-4 p-2 border border-gray-300 rounded"></textarea>
                            <button type="submit" className="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 transition-colors w-full hover:cursor-pointer">
                                作成
                            </button>
                        </Form>
                    </div>
                )}
            </div>
        </main>
    );
}
