'use client';

import { Festival } from "@mui/icons-material";
import useSWR from "swr";

// IDでアイテムを取得するカスタムフック
const useStockItem = (itemId: string) => {
    const { data, error, isLoading, mutate } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/items/${itemId}`,
        async (url: string) => {
            if (!itemId) {
                return null;
            }
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch item data");
            }
            return res.json();
        }
    );
    return { data, error, isLoading, mutate };
}

// 全アイテム一覧を取得するカスタムフック
const useQueryStockItems = (category?: string) => {
    const queryParam = category ? `?category=${category}` : '';
    const { data, error, isLoading, mutate } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/items${queryParam}`,
        async (url: string) => {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch item list");
            }
            const data = await res.json();
            return data["items"];
        }
    );

    const mutateWithFilter = async (category?: string) => {
        const queryParam = category ? `?category=${category}` : '';
        await mutate(
            async () => {
                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/items${queryParam}`);
                if (!res.ok) {
                    throw new Error("Failed to fetch item list");
                }
                const data = await res.json();
                return data["items"];
            },
            false
        );
    }

    return { data, error, isLoading, mutate: mutateWithFilter };
}

export { useStockItem, useQueryStockItems };