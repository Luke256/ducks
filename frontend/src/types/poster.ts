
type PosterStatus = 'uncollected' | 'collected' | 'lost';

const PosterStatusLabels: { [key in PosterStatus]: string } = {
    'uncollected': '未回収',
    'collected': '回収済み',
    'lost': '消失',
};

type Poster = {
    id: string;
    name: string;
    festival_id: string;
    description: string;
    image_url: string;
    status: PosterStatus;
}

export type { Poster, PosterStatus };
export { PosterStatusLabels };