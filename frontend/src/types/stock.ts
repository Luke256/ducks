import type { StockItem } from "./stockItem";

type Stock = {
    id: string;
    festival_id: string;
    price: number;
    description: string;
    item: StockItem;
};

export type { Stock };
