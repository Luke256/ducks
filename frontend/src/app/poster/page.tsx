'use client';

import { Festival } from "@/types/festival";
import { useEffect, useState } from "react";

const getFestivals = async () => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals`, {
    cache: 'force-cache'
  });
  const data = await res.json();
  return data["festivals"];
}

const getPosters = async (festivalId: string) => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/festivals/${festivalId}/posters`);
  const data = await res.json();
  console.log(data);
  return data["posters"];
}

export default function PosterPage() {
  const [festivals, setFestivals] = useState([]);
  const [posters, setPosters] = useState([]);
  const [currentFestival, setCurrentFestival] = useState<Festival | null>(null);

  useEffect(() => {
    getFestivals().then((data) => setFestivals(data));
  }, []);

  useEffect(() => {
    if (!currentFestival) {
      setPosters([]);
      return;
    }
    getPosters(currentFestival.id).then((data) => {
      setPosters(data)
    });
  }, [currentFestival]);

  console.log(posters);

  return (
    <main className="min-h-screen p-12">
      <div className="max-w-3xl bg-white p-8">
        <h1 className="mb-4 text-2xl font-bold text-black">ポスター</h1>
        <select className="mb-4 p-2 border border-gray-300 hover:cursor-pointer" onChange={(e) => {
          const selectedFestival = festivals.find((festival: Festival) => festival.id === e.target.value);
          setCurrentFestival(selectedFestival || null);
        }}>
          <option value="">イベントを選択</option>
          {festivals.map((festival: Festival) => (
            <option key={festival.id} value={festival.id}>
              {festival.name}
            </option>
          ))}
        </select>

        {/* display posters for the festival */}
        {currentFestival && posters.length === 0 && (
          <p className="mb-4 text-gray-500">このイベントにはまだポスターがありません。</p>
        )}
        {currentFestival && posters.length > 0 && (
          <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
            {posters.map((poster: any) => (
              <div key={poster.id} className="border border-gray-300 p-4">
                <img src={poster.image_url} alt={poster.name} className="mb-2 w-full h-auto" />
                <h2 className="text-lg font-bold text-black">{poster.name}</h2>
                <p className="text-gray-700">{poster.description}</p>
              </div>
            ))}
          </div>
        )}

        {/* create poster for the festival */}
        {currentFestival && (
          <div className="mb-4">
            <form action={async (formData: FormData) => {
              const name = formData.get("name") as string;
              const description = formData.get("description") as string;
              const imageFile = formData.get("image") as File;
              const form = new FormData();
              form.append("name", name);
              form.append("description", description);
              form.append("image", imageFile);
              form.append("festival_id", currentFestival.id);
              await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters`, {
                method: "POST",
                body: form,
              });
              const updatedPosters = await getPosters(currentFestival.id);
              setPosters(updatedPosters);
            }}>
              <h2 className="mb-2 text-xl font-bold text-black">新しいポスターを作成</h2>
              <input type="text" name="name" placeholder="ポスター名" required className="mb-2 p-2 border border-gray-300 w-full" />
              <textarea name="description" placeholder="説明" required className="mb-2 p-2 border border-gray-300 w-full"></textarea>
              <input type="file" name="image" accept="image/*" required className="mb-2" />
              <br />
              <button type="submit" className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">作成</button>
            </form>
          </div>
        )}
      </div>
    </main>
  );
}
