import PosterDetail from "./posterDetail";

export default async function PosterDetailPage({ params }: Readonly<{
    params: Promise<{ posterId: string }>;
}>) {
    const { posterId } = await params;

    return (
        <main className="min-h-screen">
            {/* Additional poster details would go here */}

            <PosterDetail params={{ posterId }} />
        </main>
    );
}
