"use client";

import { usePoster } from "@/hooks/posterHook";

export default function PosterDetail({ params }: Readonly<{
    params: { posterId: string };
}>) {
    const { posterId } = params;

    const { data, error, isLoading } = usePoster(posterId);

    return (
        <div>
            <h1>Poster Detail Page for Poster ID: {posterId}</h1>
            {/* Additional poster details would go here */}

            {isLoading && <p>Loading...</p>}
            {error && <p>Error loading poster data: {error.message}</p>}
            {data && (
                <div>
                    <h2>{data.name}</h2>
                    <p>{data.description}</p>
                    <img src={data.image_url} alt={data.name} />
                    <p>Status: {data.status}</p>
                </div>
            )}
        </div>
    );
}