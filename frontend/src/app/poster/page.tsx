'use client';

import { Festival } from "@/types/festival";
import StatusPicker from "./statusPicker";
import { LaunchTwoTone, RefreshTwoTone } from "@mui/icons-material";
import Link from "next/link";
import { useFestivalList } from "@/hooks/festivalHook";
import { usePosterList } from "@/hooks/posterHook";
import { useRef, useState } from "react";
import { useSessionStorage } from "@/hooks/sessStorage";
import { Poster, PosterStatusLabels } from "@/types/poster";
import { toast } from "react-toastify";

const posterItemBg = {
  "uncollected": "bg-yellow-100",
  "collected": "bg-green-100",
  "lost": "bg-red-100"
}

export default function PosterPage() {
  const { data: festivals } = useFestivalList();
  const [currentFestivalId, setCurrentFestivalId] = useSessionStorage("currentFestivalId", "");
  const { data: posters, mutate: mutatePosters } = usePosterList(currentFestivalId);
  const [filterStatus, setFilterStatus] = useState("");
  const [filterName, setFilterName] = useState("");
  const imagePreview = useRef<HTMLImageElement>(null);

  let filteredPosters = posters;
  if (filterStatus) {
    filteredPosters = filteredPosters.filter((poster: any) => poster.status === filterStatus);
  }
  if (filterName) {
    filteredPosters = filteredPosters.filter((poster: any) => poster.name.includes(filterName));
  }

  return (
    <main className="min-h-screen">
      <div className="max-w-3xl bg-white p-8">
        <h1 className="mb-4 text-2xl font-bold text-black">ポスター</h1>
        <select className="mb-4 p-2 border border-gray-300 hover:cursor-pointer" onChange={(e) => {
          setCurrentFestivalId(e.target.value);
        }} value={currentFestivalId}>
          <option value="">イベントを選択</option>
          {festivals && festivals.map((festival: Festival) => (
            <option key={festival.id} value={festival.id}>
              {festival.name}
            </option>
          ))}
        </select>
        <Link href="/poster/new" className="mb-4 ml-2 px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 hover:cursor-pointer">新しいポスターを作成</Link>

        {/* display posters for the festival */}
        {currentFestivalId && (!posters || posters.length === 0) && (
          <p className="mb-4 text-gray-500">このイベントにはまだポスターがありません。</p>
        )}
        {currentFestivalId && posters && posters.length > 0 && (
          <div>
            <div className="mb-4 flex flex-col md:flex-row md:items-center md:space-x-4">
              <input type="text" placeholder="ポスター名" className="mb-4 p-2 border border-gray-300 w-full" onChange={(e) => {
                setFilterName(e.target.value);
              }} />

              <select className="mb-4 p-2 border border-gray-300 hover:cursor-pointer" onChange={(e) => {
                setFilterStatus(e.target.value);
              }}>
                <option value="">すべてのステータス</option>
                {Object.entries(PosterStatusLabels).map(([key, label]) => (
                  <option key={key} value={key}>
                    {label}
                  </option>
                ))}
              </select>

              <button className="mb-4 px-4 py-2 bg-gray-500 text-white hover:bg-gray-600 hover:cursor-pointer" onClick={async () => {
                // Refresh posters
                await mutatePosters();
              }}>
                <RefreshTwoTone />
              </button>
            </div>

            <table className="mb-8 w-full">
              <thead>
                <tr>
                  <th className="">ポスター名</th>
                  <th className="">説明</th>
                  <th className="">ステータス</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {filteredPosters.map((poster: Poster) => (
                  <tr key={poster.id} className={posterItemBg[poster.status]}>
                    <td className="border-t border-gray-300 p-2 text-center">{poster.name}</td>
                    <td className="border-t border-gray-300 p-2 text-center">{poster.description}</td>
                    <td className="border-t border-gray-300 p-2 text-center">
                      <StatusPicker status={poster.status} onChange={async (newStatus: string) => {
                        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/posters/${poster.id}/status`, {
                          method: "PATCH",
                          body: JSON.stringify({ status: newStatus }),
                          headers: { "Content-Type": "application/json" }
                        }
                        );
                        if (res.ok) {
                          toast.success("ポスターのステータスが更新されました");
                        } else {
                          toast.error(`ポスターのステータスの更新に失敗しました: ${res.statusText}`);
                        }
                        await mutatePosters();
                      }} />
                    </td>
                    <td className="border-t border-gray-300 p-2 text-center">
                      <Link href={`/poster/detail/${poster.id}`}>
                        <LaunchTwoTone className="hover:cursor-pointer" />
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </main>
  );
}
