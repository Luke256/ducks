"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";

// type PageMode = "items" | "stocks" | "orders" | "cashiers";

const PageTabs = [
    { href: "cashiers", label: "会計処理" },
    { href: "orders", label: "売上管理" },
    { href: "stocks", label: "物販設定" },
    { href: "items", label: "商品管理" },
]

const SalesPageLayout = ({ children }: { children: React.ReactNode }) => {
    const path = usePathname() || "/";

    const isActive = (href: string) => {
        return path.startsWith(`/sales/${href}`);
    };

    return (
        <div>
            <div className="flex flex-row border-b">
                {PageTabs.map((tab) => (
                    <div
                        key={tab.href}
                        className={
                            "hover:cursor-pointer rounded-t-2xl border-2 border-b-0 p-2 pb-1 " +
                            (isActive(tab.href) ? "border-yellow-500" : "border-transparent")
                        }
                    >
                        <Link href={`/sales/${tab.href}`} className={"p-2 text-xl font-semibold hover:cursor-pointer"}>{tab.label}</Link>
                    </div>
                ))}
            </div>
            <div>{children}</div>
        </div>
    )
}

export default SalesPageLayout;