import { StockItem } from "./stockItem";

type FestivalItem = {
    id: number;
    festivalId: string;
    price: number;
    description: string;
    item: StockItem;
}

export type { FestivalItem };