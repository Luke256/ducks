import { Metadata } from "next";
import EventPageClient from "./pageClient";

export const metadata: Metadata = {
    title: "イベント管理",
    description: "イベントの管理ページ",
}



export default function EventPage() {
    return <EventPageClient />;
}