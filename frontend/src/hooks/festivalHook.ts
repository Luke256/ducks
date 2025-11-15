'use client';

import useSWR from "swr";

// IDでイベントを取得するカスタムフック
const useFestival = (festivalId: string) => {
    const { data, error, isLoading } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/festivals/${festivalId}`,
        async (url: string) => {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch festival data");
            }
            return res.json();
        }
    );
    return { data, error, isLoading };
};

// 全イベント一覧を取得するカスタムフック
const useFestivalList = () => {
    const { data, error, isLoading } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/festivals`,
        async (url: string) => {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch festival list");
            }
            const data = await res.json();
            return data["festivals"];
        }
    );
    return { data, error, isLoading };
};

export { useFestival, useFestivalList };