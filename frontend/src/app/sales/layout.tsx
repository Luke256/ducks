"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const PageTabs = [
    { href: "cashier", label: "会計" },
    { href: "orders", label: "売上" },
    { href: "stocks", label: "目録" },
    { href: "items", label: "商品" },
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
                            "hover:cursor-pointer rounded-t-xl border-2 border-b-0 p-1 " +
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