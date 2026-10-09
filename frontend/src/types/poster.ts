import type { Festival } from "./festival";

type PosterStatus = "uncollected" | "collected" | "lost";

const PosterStatusLabels: { [key in PosterStatus]: string } = {
    uncollected: "未回収",
    collected: "回収済み",
    lost: "消失",
};

type PosterImage = {
    id: string;
    url: string;
};

type Poster = {
    id: string;
    name: string;
    description: string;
    image: PosterImage[];
    status: PosterStatus;
    festival: Festival;
};

export type { Poster, PosterImage, PosterStatus };
export { PosterStatusLabels };
export const maxPosterImages = 10;
