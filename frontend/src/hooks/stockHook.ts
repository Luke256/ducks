'use client';

import { Stock } from "@/types/stock";
import useSWR from "swr";

// IDで物販ストックを取得するカスタムフック
const useStock = (stockId: string) => {
    const { data, error, isLoading, mutate } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/stocks/${stockId}`,
        async (url: string) => {
            if (!stockId) {
                return null;
            }
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch stock data");
            }
            return res.json();
        }
    );
    return { data, error, isLoading, mutate };
}

// イベントIDから物販ストック一覧を取得するカスタムフック
const useStockListByFestival = (festivalId: string) => {
    const { data, error, isLoading, mutate } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/festivals/${festivalId}/stocks`,
        async (url: string) => {
            if (!festivalId) {
                return [];
            }
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch stock list");
            }
            const data = await res.json();
            const stocks = data["stocks"];
            
            if (!stocks) {
                return [];
            }

            // category->nameでソート
            stocks.sort((a: Stock, b: Stock) => {
                if (a.item.category < b.item.category) return -1;
                if (a.item.category > b.item.category) return 1;
                if (a.item.name < b.item.name) return -1;
                if (a.item.name > b.item.name) return 1;
                return 0;
            });

            return stocks;
        }
    );
    return { data, error, isLoading, mutate };
}

export { useStock, useStockListByFestival };