import { StockItem } from "./stockItem";

type FestivalItem = {
    id: number;
    festivalId: string;
    price: number;
    item: StockItem;
}

export type { FestivalItem };