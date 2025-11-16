import { Metadata } from "next";
import PosterDetail from "./posterDetail";

export async function generateMetadata({ params }: Readonly<{
    params: Promise<{ posterId: string }>;
}>): Promise<Metadata> {
    const { posterId } = await params;

    return {
        title: `ポスター詳細: ${posterId}`,
        description: `ポスター「${posterId}」の詳細ページ`,
    };
}

export default async function PosterDetailPage({ params }: Readonly<{
    params: Promise<{ posterId: string }>;
}>) {
    const { posterId } = await params;

    return (
        <main>
            {/* Additional poster details would go here */}

            <PosterDetail params={{ posterId }} />
        </main>
    );
}
