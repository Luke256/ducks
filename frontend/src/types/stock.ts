import type { StockItem } from "./stockItem";

type Stock = {
    id: string;
    festival_id: string;
    price: number;
    description: string;
    for_sale: boolean;
    item: StockItem;
};

export type { Stock };
