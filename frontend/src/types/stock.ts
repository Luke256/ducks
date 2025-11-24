import { StockItem } from "./stockItem";

type Stock = {
    id: string;
    festivalId: string;
    price: number;
    description: string;
    item: StockItem;
}

export type { Stock };