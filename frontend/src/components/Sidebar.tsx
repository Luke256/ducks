"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

export default function Sidebar() {
  const pathname = usePathname() || "/";

  const isActive = (path: string) =>
    pathname === path || pathname.startsWith(path + "/");

  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-yellow-400 p-6 px-0">
      <h2 className="mb-6 text-4xl font-black text-blue-950 px-6">Ducks</h2>
      <nav className="flex flex-col" aria-label="サイドメニュー">
        <Link
          href="/event"
          className={`text-left px-6 py-3 transition-colors cursor-pointer ${
            isActive("/event")
              ? "bg-white text-yellow-600 font-semibold shadow"
              : "text-blue-950/85 hover:bg-white/70"
          }`}
        >
          イベント
        </Link>

        <Link
          href="/poster"
          className={`text-left px-6 py-3 transition-colors cursor-pointer ${
            isActive("/poster")
              ? "bg-white text-yellow-600 font-semibold shadow"
              : "text-blue-950/85 hover:bg-white/70"
          }`}
        >
          ポスター
        </Link>
      </nav>
    </aside>
  );
}
