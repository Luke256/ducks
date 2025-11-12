import { Festival } from "@/types/festival";
import EditFestivalClient from "./EditFestivalClient";

const getEvent = async (eventId: string): Promise<Festival> => {
    const endpoint = `${process.env.NEXT_PUBLIC_API_URL}/festivals/${eventId}`;
    const res = await fetch(endpoint);
    const data = await res.json();
    console.log("Fetched event data:", data);
    return data;
}

export default async function EventDetailPage({ params }: Readonly<{
    params: Promise<{ eventId: string }>;
}>) {
    const { eventId } = await params;
    const event = await getEvent(eventId);

    return (
        <main className="min-h-screen p-12 max-w-7xl mx-auto">
            <div className="max-w-3xl bg-white p-8">
                <h1 className="mb-4 text-2xl font-bold text-black">
                    イベント詳細: {event.name}
                </h1>
                <p className="text-gray-600">概要: {event.description}</p>

                <EditFestivalClient event={event} eventId={eventId} />
            </div>
        </main>
    );
}