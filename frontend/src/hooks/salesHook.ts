"use client";

import { SaleRecord } from "@/types/saleRecord";
import useSWR from "swr";

// 売上データを取得するカスタムフック
const useSalesDataList = (): {
    data: SaleRecord[] | null;
    error: Error | null;
    isLoading: boolean;
    mutate: () => Promise<void>;
} => {
    const { data, error, isLoading, mutate } = useSWR(
        `${process.env.NEXT_PUBLIC_API_URL}/sales`,
        async (url: string) => {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error("Failed to fetch sales data");
            }
            const data = await res.json();
            const sales = data["sales"];
            
            sales.sort((a: SaleRecord, b: SaleRecord) => {
                return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
            });

            return sales;
        }
    );

    return { data, error, isLoading, mutate };
}

export { useSalesDataList };