"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";

const navItems = [
  { href: "/event", label: "イベント" },
  { href: "/poster", label: "ポスター" },
];

export default function Sidebar() {
  const pathname = usePathname() || "/";

  const isActive = (path: string) =>
    pathname === path || pathname.startsWith(path + "/");

  const linkClasses = (path: string) =>
    `text-left px-6 py-3 transition-colors cursor-pointer ${isActive(path)
      ? "bg-white text-yellow-600 font-semibold shadow"
      : "text-blue-950/85 hover:bg-white/70"
    }`;

  const bottomLinkClasses = (path: string) =>
    `flex flex-1 flex-col items-center gap-1 py-2 text-xs font-semibold transition-colors ${isActive(path)
      ? "text-blue-950"
      : "text-blue-950/70 hover:text-blue-950"
    }`;

  return (
    <>
      {/* Navigation for PC */}
      <aside className="invisible md:visible fixed left-0 top-0 h-full w-64 bg-yellow-400 p-6 px-0">
        <h2 className="mb-6 text-4xl font-black text-blue-950 px-6">Ducks</h2>
        <nav className="flex flex-col" aria-label="サイドメニュー">
          {navItems.map(({ href, label }) => (
            <Link key={href} href={href} className={linkClasses(href)}>
              {label}
            </Link>
          ))}
        </nav>
      </aside>

      {/* Navigation for Mobile */}
      <div className="py-2 md:hidden sticky top-0 z-30 bg-yellow-400 px-4">
        <h2 className="text-3xl font-black text-blue-950">Ducks</h2>
      </div>

      <nav
        aria-label="モバイル下部メニュー"
        className="md:hidden fixed bottom-0 left-0 right-0 z-30 border-t border-yellow-200 bg-white/95 backdrop-blur px-2"
      >
        <div className="flex">
          {navItems.map(({ href, label }) => (
            <Link key={href} href={href} className={bottomLinkClasses(href)}>
              <span>{label}</span>
              <span
                className={`h-1 w-12 rounded-full transition-colors ${isActive(href) ? "bg-yellow-500" : "bg-transparent"
                  }`}
                aria-hidden="true"
              />
            </Link>
          ))}
        </div>
      </nav>
    </>
  );
}
