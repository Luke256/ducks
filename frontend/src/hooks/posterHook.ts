'use client';

import useSWR from "swr"

// IDでポスターを取得するカスタムフック
const usePoster = (posterId: string) => {
    const { data, error, isLoading, isValidating, mutate } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/posters/${posterId}`,
        async (url: string) => {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch poster data");
            }
            return res.json();
        }
    )
    return { data, error, isLoading, isValidating, mutate };
}

// イベントIDでポスター一覧を取得するカスタムフック
const usePosterList = (festivalId: string) => {
    const { data, error, isLoading, isValidating, mutate } = useSWR(
        festivalId ? `${process.env.NEXT_PUBLIC_API_URL}/festivals/${festivalId}/posters` : null,
        async (url: string) => {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch poster list");
            }
            const data = await res.json();
            return data["posters"];
        }
    )
    return { data, error, isLoading, isValidating, mutate };
}

export { usePoster, usePosterList };